//go:build linux && cgo

/*
 * Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dcgm

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthWhenInvalidGroupID(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	runOnlyWithLiveGPUs(t)

	var invalidGroupID uintptr = 99
	gh := GroupHandle{}
	gh.SetHandle(invalidGroupID)
	err := HealthSet(gh, DCGM_HEALTH_WATCH_PCIE)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Setting not configured")

	_, err = HealthGet(gh)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Setting not configured")

	_, err = HealthGet(gh)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Setting not configured")
}

func TestHealthCheckPCIE(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	runOnlyWithLiveGPUs(t)
	gpus, err := withInjectionGPUs(t, 1)
	require.NoError(t, err)

	type testCase struct {
		name              string
		pcieGen           int
		pcieGenSpeed      float64 // in Gbps
		pcieLanes         int
		pcieReplayCounter int
		expectingIncident bool
	}

	pcieGenSpeeds := []float64{
		2.5,  // Gen1 speed in Gbps
		5.0,  // Gen2
		8.0,  // Gen3
		16.0, // Gen4
		32.0, // Gen5
		64.0, // Gen6
	}

	var tests []testCase
	// Generate test cases
	for i := 0; i < 1; i++ { // Run multiple iterations
		for gen, speed := range pcieGenSpeeds {
			pcieGen := gen + 1
			pcieLanes := rand.Intn(16) + 1
			ratePerLane := speed / 1000 * 60 // Convert to errors/min per lane
			expectedLimit := math.Ceil(ratePerLane * float64(pcieLanes))
			pcieReplayCounter := rand.Intn(2*int(expectedLimit)) + 1
			expectingIncident := pcieReplayCounter > int(expectedLimit)

			tests = append(tests, testCase{
				name:              fmt.Sprintf("PCIe_Gen%d_%dLanes_Counter%d", pcieGen, pcieLanes, pcieReplayCounter),
				pcieGen:           pcieGen,
				pcieGenSpeed:      speed,
				pcieLanes:         pcieLanes,
				pcieReplayCounter: pcieReplayCounter,
				expectingIncident: expectingIncident,
			})
		}
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ratePerLane := tc.pcieGenSpeed / 1000 * 60
			expectedLimit := math.Ceil(ratePerLane * float64(tc.pcieLanes))

			errMsg := fmt.Sprintf("pcieGen=%d pcieGenSpeed=%f pcieLanes=%d expectedLimit=%f pcieReplayCounter=%d expectingIncident=%v",
				tc.pcieGen, tc.pcieGenSpeed, tc.pcieLanes, expectedLimit, tc.pcieReplayCounter, tc.expectingIncident)

			healthCheckPCIE(t, gpus, tc.pcieGen, tc.pcieLanes, tc.pcieReplayCounter, tc.expectingIncident, errMsg)
			defer resetPCICReplayCounter(t, gpus)
		})
	}
}

func resetPCICReplayCounter(t *testing.T, gpuIDs []uint) {
	gpuID := gpuIDs[0]
	err := InjectFieldValue(gpuID,
		DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		DCGM_FT_INT64,
		0,
		time.Now().Add(100*time.Second).UnixMicro(),
		int64(0),
	)
	require.NoError(t, err)
}

func healthCheckPCIE(t *testing.T, gpuIDs []uint, pcieGen int, pcieLanes int, pcieReplayCounter int, expectingPCIEIncident bool, errMessage string) {
	gpuID := gpuIDs[0]

	groupID, err := CreateGroup("test1")
	require.NoError(t, err)
	defer func() {
		_ = DestroyGroup(groupID)
	}()
	err = AddEntityToGroup(groupID, FE_GPU, gpuID)
	require.NoError(t, err)

	err = HealthSet(groupID, DCGM_HEALTH_WATCH_PCIE)
	require.NoError(t, err)

	system, err := HealthGet(groupID)
	require.NoError(t, err)
	require.Equal(t, DCGM_HEALTH_WATCH_PCIE, system)

	skipTestIfUnhealthy(t, groupID)

	// inject PCIe Gen and width/lanes
	err = InjectFieldValue(gpuID,
		DCGM_FI_DEV_PCIE_LINK_GEN,
		DCGM_FT_INT64,
		0,
		0,
		int64(pcieGen),
	)
	require.NoError(t, err)

	err = InjectFieldValue(gpuID,
		DCGM_FI_DEV_PCIE_LINK_WIDTH,
		DCGM_FT_INT64,
		0,
		0,
		int64(pcieLanes),
	)
	require.NoError(t, err)

	err = InjectFieldValue(gpuID,
		DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		DCGM_FT_INT64,
		0,
		time.Now().Add(-50*time.Second).UnixMicro(),
		int64(0),
	)
	require.NoError(t, err)

	// we expect that there will be no data here
	response, err := HealthCheck(groupID)
	require.NoError(t, err)
	require.Equal(t, DCGM_HEALTH_RESULT_PASS, response.OverallHealth)

	// inject an error into PCIe
	err = InjectFieldValue(gpuID,
		DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		DCGM_FT_INT64,
		0,
		time.Now().Add(100*time.Second).UnixMicro(),
		int64(pcieReplayCounter),
	) // set the injected data into the future
	require.NoError(t, err)

	response, err = HealthCheck(groupID)
	require.NoError(t, err)
	if expectingPCIEIncident {
		require.Len(t, response.Incidents, 1, errMessage)
		require.Equal(t, gpuID, response.Incidents[0].EntityInfo.EntityId)
		require.Equal(t, DCGM_HEALTH_WATCH_PCIE, response.Incidents[0].System)
		require.Equal(t, DCGM_FR_PCI_REPLAY_RATE, response.Incidents[0].Error.Code)
	} else {
		require.Len(t, response.Incidents, 0, errMessage)
	}
}

func skipTestIfUnhealthy(t *testing.T, groupId GroupHandle) {
	health, err := HealthCheck(groupId)
	require.NoError(t, err)
	if health.OverallHealth != DCGM_HEALTH_RESULT_PASS {
		msg := "Skipping health check test because we are already unhealthy: "
		incidents := []string{}
		for _, incident := range health.Incidents {
			incidents = append(incidents, incident.Error.Message)
		}

		t.Skip(msg + strings.Join(incidents, ", "))
	}
}
