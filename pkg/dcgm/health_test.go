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

	gpuID := gpus[0]

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

	err = InjectFieldValue(gpuID,
		DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		DCGM_FT_INT64,
		0,
		time.Now().Add(-50*time.Second).UnixMicro(),
		int64(0),
	)
	require.NoError(t, err)

	response, err := HealthCheck(groupID)
	require.NoError(t, err)
	require.Equal(t, DCGM_HEALTH_RESULT_PASS, response.OverallHealth)

	// inject an error into PCI
	err = InjectFieldValue(gpuID,
		DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		DCGM_FT_INT64,
		0,
		time.Now().Add(100*time.Second).UnixMicro(),
		int64(10),
	)
	require.NoError(t, err)
	response, err = HealthCheck(groupID)
	require.NoError(t, err)
	require.Equal(t, DCGM_HEALTH_RESULT_WARN, response.OverallHealth)
	require.Len(t, response.Incidents, 1)
	assert.Equal(t, gpuID, response.Incidents[0].EntityInfo.EntityId)
	assert.Equal(t, DCGM_HEALTH_WATCH_PCIE, response.Incidents[0].System)
	assert.Equal(t, DCGM_FR_PCI_REPLAY_RATE, response.Incidents[0].Error.Code)
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
