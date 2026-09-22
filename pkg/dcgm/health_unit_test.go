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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIMEXHealthWatchConstant(t *testing.T) {
	const want HealthSystem = 0x2000

	require.Equal(t, want, DCGM_HEALTH_WATCH_IMEX)
	require.Equal(t, want, DCGM_HEALTH_WATCH_ALL&want)
}

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
