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
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestGetDeviceInfoWithOps checks Device fields and read errors.
// It also checks that extra reads require a supported GPU with EntityStatusOk.
func TestGetDeviceInfoWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")
	attrs := deviceAttributes{
		busID:   "00000000:31:00.0",
		uuid:    "GPU-test",
		power:   72,
		bar1:    32768,
		fbTotal: 23034,
		identifiers: DeviceIdentifiers{
			Brand:               "NVIDIA",
			Model:               "Test GPU",
			Serial:              "serial",
			Vbios:               "vbios",
			InforomImageVersion: "inforom",
			DriverVersion:       "driver",
		},
	}
	topology := []P2PLink{{GPU: 4}}

	tests := []struct {
		name      string
		fail      string
		gpus      []uint
		status    EntityStatus
		want      Device
		wantError bool
	}{
		{name: "attribute failure", fail: "attributes", wantError: true},
		{name: "supported list failure", fail: "supported", wantError: true},
		{
			name:   "GPU absent from supported list",
			gpus:   []uint{4},
			status: EntityStatusOk,
			want: Device{
				GPU: 3, DCGMSupported: "No", UUID: attrs.uuid, Power: attrs.power,
				PCI:         PCIInfo{BusID: attrs.busID, BAR1: attrs.bar1, FBTotal: attrs.fbTotal},
				Identifiers: attrs.identifiers,
			},
		},
		{
			name:   "non-OK status",
			gpus:   []uint{3},
			status: EntityStatusDisabled,
			want: Device{
				GPU: 3, DCGMSupported: "No", UUID: attrs.uuid, Power: attrs.power,
				PCI:         PCIInfo{BusID: attrs.busID, BAR1: attrs.bar1, FBTotal: attrs.fbTotal},
				Identifiers: attrs.identifiers,
			},
		},
		{name: "affinity failure", fail: "affinity", gpus: []uint{3}, status: EntityStatusOk, wantError: true},
		{name: "topology failure", fail: "topology", gpus: []uint{3}, status: EntityStatusOk, wantError: true},
		{name: "bandwidth failure", fail: "bandwidth", gpus: []uint{3}, status: EntityStatusOk, wantError: true},
		{
			name:   "success",
			gpus:   []uint{3},
			status: EntityStatusOk,
			want: Device{
				GPU: 3, DCGMSupported: "Yes", UUID: attrs.uuid, Power: attrs.power,
				PCI:         PCIInfo{BusID: attrs.busID, BAR1: attrs.bar1, FBTotal: attrs.fbTotal, Bandwidth: 32000},
				Identifiers: attrs.identifiers, Topology: topology, CPUAffinity: "{0-7}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockDeviceInfoOps(gomock.NewController(t))
			attributesCall := mock.EXPECT().readDeviceAttributes(uint(3))
			if tt.fail == "attributes" {
				attributesCall.Return(deviceAttributes{}, stageErr)
			} else {
				attributesCall.Return(attrs, nil)
				supportedCall := mock.EXPECT().getSupportedDevices().After(attributesCall)
				if tt.fail == "supported" {
					supportedCall.Return(nil, stageErr)
				} else {
					supportedCall.Return(tt.gpus, nil)
					statusCall := mock.EXPECT().getGPUStatus(uint(3)).Return(tt.status).After(supportedCall)
					if slices.Contains(tt.gpus, uint(3)) && tt.status == EntityStatusOk {
						affinityCall := mock.EXPECT().getCPUAffinity(uint(3)).After(statusCall)
						if tt.fail == "affinity" {
							affinityCall.Return("", stageErr)
						} else {
							affinityCall.Return("{0-7}", nil)
							topologyCall := mock.EXPECT().getDeviceTopology(uint(3)).After(affinityCall)
							if tt.fail == "topology" {
								topologyCall.Return(nil, stageErr)
							} else {
								topologyCall.Return(topology, nil)
								bandwidthCall := mock.EXPECT().getPciBandwidth(uint(3)).After(topologyCall)
								if tt.fail == "bandwidth" {
									bandwidthCall.Return(int64(0), stageErr)
								} else {
									bandwidthCall.Return(int64(32000), nil)
								}
							}
						}
					}
				}
			}

			got, err := getDeviceInfoWithOps(mock, 3)
			if tt.wantError {
				require.ErrorIs(t, err, stageErr)
				require.Equal(t, Device{}, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
