//go:build linux && cgo

/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
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
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestHealthCheckByGPUWithOpsCleansUp checks that the temporary group is
// destroyed whichever step fails, and that the failing step's error is returned.
func TestHealthCheckByGPUWithOpsCleansUp(t *testing.T) {
	for _, tt := range []struct {
		name string
		fail string
	}{
		{name: "create failure", fail: "create"},
		{name: "add failure", fail: "add"},
		{name: "watch setup failure", fail: "set"},
		{name: "check failure", fail: "check"},
		{name: "success"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockHealthGroupOps(gomock.NewController(t))
			stageErr := errors.New(tt.fail + " failed")
			expectHealthCheckSteps(m, tt.fail, stageErr)

			got, err := healthCheckByGPUWithOps(m, 5)

			if tt.fail == "" {
				require.NoError(t, err)
				require.Equal(t, uint(5), got.GPU)
			} else {
				require.ErrorIs(t, err, stageErr)
			}
		})
	}
}

// expectHealthCheckSteps registers, in order, the calls healthCheckByGPUWithOps
// makes for GPU 5. The step named by fail ("create", "add", "set" or "check")
// returns stageErr and no later step is expected; an empty fail means every step
// succeeds. Once the group exists, destroyGroup must be called.
func expectHealthCheckSteps(m *MockHealthGroupOps, fail string, stageErr error) {
	var group GroupHandle
	group.SetHandle(11)

	create := m.EXPECT().createGroup(gomock.Any(), false)
	if fail == "create" {
		create.Return(GroupHandle{}, stageErr)
		return
	}
	create.Return(group, nil)
	m.EXPECT().destroyGroup(group).After(create)

	add := m.EXPECT().addToGroup(group, uint(5)).After(create)
	if fail == "add" {
		add.Return(stageErr)
		return
	}
	add.Return(nil)

	set := m.EXPECT().healthSet(group, DCGM_HEALTH_WATCH_ALL).After(add)
	if fail == "set" {
		set.Return(stageErr)
		return
	}
	set.Return(nil)

	check := m.EXPECT().healthCheck(group).After(set)
	if fail == "check" {
		check.Return(HealthResponse{}, stageErr)
		return
	}
	check.Return(HealthResponse{}, nil)
}

// TestIMEXHealthWatchConstant checks the IMEX watch bit matches the DCGM header and is
// part of DCGM_HEALTH_WATCH_ALL.
func TestIMEXHealthWatchConstant(t *testing.T) {
	const want HealthSystem = 0x2000

	require.Equal(t, want, DCGM_HEALTH_WATCH_IMEX)
	require.Equal(t, want, DCGM_HEALTH_WATCH_ALL&want)
}

// TestHealthStatus checks the display name of each health result, including an unknown one.
func TestHealthStatus(t *testing.T) {
	tests := []struct {
		name   string
		status HealthResult
		want   string
	}{
		{name: "healthy", status: 0, want: "Healthy"},
		{name: "warning", status: 10, want: "Warning"},
		{name: "failure", status: 20, want: "Failure"},
		{name: "unknown", status: 99, want: "N/A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, healthStatus(tt.status))
		})
	}
}

// TestSystemWatch checks the display name of each health watch system, including an unknown one.
func TestSystemWatch(t *testing.T) {
	tests := []struct {
		name  string
		watch HealthSystem
		want  string
	}{
		{name: "PCIe", watch: DCGM_HEALTH_WATCH_PCIE, want: "PCIe watches"},
		{name: "NVLink", watch: DCGM_HEALTH_WATCH_NVLINK, want: "NVLINK watches"},
		{name: "PMU", watch: DCGM_HEALTH_WATCH_PMU, want: "Power Management unit watches"},
		{name: "MCU", watch: DCGM_HEALTH_WATCH_MCU, want: "Microcontroller unit watches"},
		{name: "memory", watch: DCGM_HEALTH_WATCH_MEM, want: "Memory watches"},
		{name: "SM", watch: DCGM_HEALTH_WATCH_SM, want: "Streaming Multiprocessor watches"},
		{name: "Inforom", watch: DCGM_HEALTH_WATCH_INFOROM, want: "Inforom watches"},
		{name: "thermal", watch: DCGM_HEALTH_WATCH_THERMAL, want: "Temperature watches"},
		{name: "power", watch: DCGM_HEALTH_WATCH_POWER, want: "Power watches"},
		{name: "driver", watch: DCGM_HEALTH_WATCH_DRIVER, want: "Driver-related watches"},
		{name: "NVSwitch non-fatal", watch: DCGM_HEALTH_WATCH_NVSWITCH_NONFATAL, want: "NVSwitch non-fatal watches"},
		{name: "NVSwitch fatal", watch: DCGM_HEALTH_WATCH_NVSWITCH_FATAL, want: "NVSwitch fatal watches"},
		{name: "ConnectX", watch: DCGM_HEALTH_WATCH_CONNECTX, want: "ConnectX watches"},
		{name: "IMEX", watch: DCGM_HEALTH_WATCH_IMEX, want: "IMEX watches"},
		{name: "unknown", watch: HealthSystem(1 << 30), want: "N/A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, systemWatch(tt.watch))
		})
	}
}

// TestDeviceHealthFromResponse checks that a group health result becomes the
// right per-GPU summary, with and without incidents.
func TestDeviceHealthFromResponse(t *testing.T) {
	tests := []struct {
		name     string
		gpuID    uint
		response HealthResponse
		want     DeviceHealth
	}{
		{
			name:     "healthy without incidents",
			gpuID:    2,
			response: HealthResponse{OverallHealth: 0},
			want:     DeviceHealth{GPU: 2, Status: "Healthy", Watches: []SystemWatch{}},
		},
		{
			name:  "maps incidents",
			gpuID: 7,
			response: HealthResponse{
				OverallHealth: 10,
				Incidents: []Incident{
					{System: DCGM_HEALTH_WATCH_PCIE, Health: 10, Error: DiagErrorDetail{Message: "degraded"}},
					{System: HealthSystem(1 << 30), Health: 99, Error: DiagErrorDetail{Message: "unknown"}},
				},
			},
			want: DeviceHealth{
				GPU:    7,
				Status: "Warning",
				Watches: []SystemWatch{
					{Type: "PCIe watches", Status: "Warning", Error: "degraded"},
					{Type: "N/A", Status: "N/A", Error: "unknown"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, deviceHealthFromResponse(tt.gpuID, tt.response))
		})
	}
}

// TestHealthResponseFromNative checks the empty, single, and full native arrays,
// including every incident field and the aggregate health result.
func TestHealthResponseFromNative(t *testing.T) {
	capacity := healthIncidentCapacityForTest()
	for _, count := range []uint32{0, 1, capacity} {
		t.Run(fmt.Sprintf("count %d", count), func(t *testing.T) {
			incidents := make([]Incident, count)
			for i := range incidents {
				incidents[i] = Incident{
					System: DCGM_HEALTH_WATCH_PCIE,
					Health: HealthResult(10),
					Error: DiagErrorDetail{
						Message: fmt.Sprintf("incident %d", i),
						Code:    HealthCheckErrorCode(i + 1),
					},
					EntityInfo: GroupEntityPair{EntityGroupId: FE_GPU, EntityId: uint(i + 7)},
				}
			}

			response, err := healthResponseFromValuesForTest(HealthResult(20), incidents, count)

			require.NoError(t, err)
			require.Equal(t, HealthResponse{OverallHealth: HealthResult(20), Incidents: incidents}, response)
		})
	}
}

// TestHealthResponseFromNativeRejectsInvalidCounts exercises the actual native
// decoder repeatedly, with a valid response after each rejected count.
func TestHealthResponseFromNativeRejectsInvalidCounts(t *testing.T) {
	capacity := healthIncidentCapacityForTest()
	for _, count := range []uint32{capacity + 1, math.MaxUint32} {
		t.Run(fmt.Sprintf("count %d", count), func(t *testing.T) {
			for range 3 {
				response, err := healthResponseFromValuesForTest(HealthResult(20), nil, count)

				require.EqualError(t, err, fmt.Sprintf("health response contains %d incidents, maximum %d", count, capacity))
				require.Equal(t, HealthResponse{}, response)

				valid := []Incident{{
					System: DCGM_HEALTH_WATCH_MEM,
					Health: HealthResult(10),
					Error:  DiagErrorDetail{Message: "recovered", Code: 42},
				}}
				response, err = healthResponseFromValuesForTest(HealthResult(10), valid, 1)

				require.NoError(t, err)
				require.Equal(t, HealthResponse{OverallHealth: HealthResult(10), Incidents: valid}, response)
			}
		})
	}
}

// TestHealthSetWithParamsStrictRejectsWideMask checks that bits above DCGM's 32-bit
// mask are rejected before DCGM is called, because they would otherwise be dropped.
func TestHealthSetWithParamsStrictRejectsWideMask(t *testing.T) {
	err := HealthSetWithParamsStrict(GroupHandle{}, HealthSystem(1)<<32, time.Second, time.Minute)

	require.False(t, IsFeatureUnavailable(err))

	var dcgmErr *Error
	require.ErrorAs(t, err, &dcgmErr)
	require.Equal(t, testDCGMReturn(testDCGMStatusBadParam), dcgmErr.Code)
}

// TestHealthV3CallsWithoutLibrary checks the DCGM 4.7.1 health functions when no
// DCGM library is loaded: each reports its unavailable symbol, with no DCGM code.
func TestHealthV3CallsWithoutLibrary(t *testing.T) {
	oldLibHandle := dcgmLibHandle
	dcgmLibHandle = nil
	t.Cleanup(func() {
		dcgmLibHandle = oldLibHandle
	})

	requireUnavailable := func(t *testing.T, err error, feature, symbol string) {
		t.Helper()

		var unavailable *FeatureUnavailableError
		require.ErrorAs(t, err, &unavailable)
		require.Equal(t, feature, unavailable.Feature)
		require.Equal(t, symbol, unavailable.Symbol)

		var dcgmErr *Error
		require.NotErrorAs(t, err, &dcgmErr, "a missing symbol has no DCGM code")
	}

	systems, err := HealthGetSupportedSystems()
	require.Zero(t, systems)
	requireUnavailable(t, err, healthSupportedSystemsFeature, healthSupportedSystemsSymbol)

	err = HealthSetWithParamsStrict(GroupHandle{}, DCGM_HEALTH_WATCH_ALL, time.Second, time.Minute)
	requireUnavailable(t, err, healthSetV3Feature, healthSetV3Symbol)
}
