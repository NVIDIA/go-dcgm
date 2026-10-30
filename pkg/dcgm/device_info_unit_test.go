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

// TestGetPciBandwidthWithOps checks setup failures, value reads, and temporary-resource cleanup.
func TestGetPciBandwidthWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")
	cleanupErr := errors.New("cleanup failed")
	pciFields := []Short{DCGM_FI_DEV_PCIE_MAX_LINK_GEN, DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH}
	var fieldsID FieldHandle
	fieldsID.SetHandle(17)
	var groupID GroupHandle
	groupID.SetHandle(23)

	tests := []struct {
		name string
		fail string
		want int64
	}{
		{name: "field group creation failure", fail: "create"},
		{name: "watch setup failure", fail: "watch"},
		{name: "value collection failure", fail: "collect"},
		{name: "success", want: 31504},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockDeviceStatusOps(gomock.NewController(t))
			create := mock.EXPECT().fieldGroupCreate(gomock.Any(), pciFields)
			if tt.fail == "create" {
				create.Return(FieldHandle{}, stageErr)
			} else {
				create.Return(fieldsID, nil)
				watch := mock.EXPECT().watchFields(uint(3), fieldsID, gomock.Any()).After(create)
				if tt.fail == "watch" {
					watch.Return(GroupHandle{}, stageErr)
					mock.EXPECT().fieldGroupDestroy(fieldsID).Return(cleanupErr).After(watch)
				} else {
					watch.Return(groupID, nil)
					collect := mock.EXPECT().getLatestValuesForFields(uint(3), pciFields).After(watch)
					if tt.fail == "collect" {
						collect.Return(nil, stageErr)
					} else {
						collect.Return([]FieldValue_v1{testInt64FieldValue(4), testInt64FieldValue(16)}, nil)
					}
					destroyFields := mock.EXPECT().fieldGroupDestroy(fieldsID).Return(cleanupErr).After(collect)
					mock.EXPECT().destroyGroup(groupID).Return(cleanupErr).After(destroyFields)
				}
			}

			got, err := getPciBandwidthWithOps(mock, 3)
			switch tt.fail {
			case "create", "watch":
				require.ErrorIs(t, err, stageErr)
				require.Zero(t, got)
			case "collect":
				require.EqualError(t, err, "error getting Pcie bandwidth: stage failed")
				require.Zero(t, got)
			default:
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// TestGetCPUAffinityWithOps checks setup failures, bitmask formatting, and deferred cleanup.
func TestGetCPUAffinityWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")
	cleanupErr := errors.New("cleanup failed")
	affinityFields := []Short{
		DCGM_FI_DEV_CPU_AFFINITY_0,
		DCGM_FI_DEV_CPU_AFFINITY_1,
		DCGM_FI_DEV_CPU_AFFINITY_2,
		DCGM_FI_DEV_CPU_AFFINITY_3,
	}
	var fieldsID FieldHandle
	fieldsID.SetHandle(29)
	var groupID GroupHandle
	groupID.SetHandle(31)

	tests := []struct {
		name           string
		fail           string
		cleanupFailure bool
		values         []int64
		want           string
	}{
		{name: "field group creation failure", fail: "create"},
		{name: "watch setup failure", fail: "watch"},
		{name: "value collection failure", fail: "collect"},
		{
			name:   "formats words across the complete affinity mask",
			values: []int64{1<<0 | 1<<3, 1, -1 << 63, 1 << 1},
			want:   "{0,3,64,191,193}",
		},
		{
			name:           "cleanup failures do not replace a successful result",
			cleanupFailure: true,
			values:         []int64{1, 0, 0, 0},
			want:           "{0}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockDeviceStatusOps(gomock.NewController(t))
			create := mock.EXPECT().fieldGroupCreate(gomock.Any(), affinityFields)
			if tt.fail == "create" {
				create.Return(FieldHandle{}, stageErr)
			} else {
				create.Return(fieldsID, nil)
				watch := mock.EXPECT().watchFields(uint(5), fieldsID, gomock.Any()).After(create)
				if tt.fail == "watch" {
					watch.Return(GroupHandle{}, stageErr)
					mock.EXPECT().fieldGroupDestroy(fieldsID).Return(nil).After(watch)
				} else {
					watch.Return(groupID, nil)
					collect := mock.EXPECT().getLatestValuesForFields(uint(5), affinityFields).After(watch)
					if tt.fail == "collect" {
						collect.Return(nil, stageErr)
					} else {
						values := make([]FieldValue_v1, len(tt.values))
						for i, value := range tt.values {
							values[i] = testInt64FieldValue(value)
						}
						collect.Return(values, nil)
					}

					cleanupResult := error(nil)
					if tt.cleanupFailure {
						cleanupResult = cleanupErr
					}
					destroyGroup := mock.EXPECT().destroyGroup(groupID).Return(cleanupResult).After(collect)
					mock.EXPECT().fieldGroupDestroy(fieldsID).Return(cleanupResult).After(destroyGroup)
				}
			}

			got, err := getCPUAffinityWithOps(mock, 5)
			switch tt.fail {
			case "create", "watch":
				require.ErrorIs(t, err, stageErr)
				require.Equal(t, "N/A", got)
			case "collect":
				require.EqualError(t, err, "error getting cpu affinity: stage failed")
				require.Equal(t, "N/A", got)
			default:
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// TestSupportedDevicesFromNative checks supported-device conversion at valid native count boundaries.
func TestSupportedDevicesFromNative(t *testing.T) {
	got, err := supportedDevicesFromValuesForTest(nil, 0)
	require.NoError(t, err)
	require.Empty(t, got)

	got, err = supportedDevicesFromValuesForTest([]uint{2, 7, 11}, 3)
	require.NoError(t, err)
	require.Equal(t, []uint{2, 7, 11}, got)

	capacity := supportedDeviceCapacityForTest()
	deviceIDs := make([]uint, capacity)
	for i := range deviceIDs {
		deviceIDs[i] = uint(i + 1)
	}
	got, err = supportedDevicesFromValuesForTest(deviceIDs, capacity)
	require.NoError(t, err)
	require.Equal(t, deviceIDs, got)
}

// TestSupportedDevicesFromNativeRejectsInvalidCounts checks negative and oversized native counts.
func TestSupportedDevicesFromNativeRejectsInvalidCounts(t *testing.T) {
	capacity := supportedDeviceCapacityForTest()
	for _, count := range []int{-1, capacity + 1} {
		t.Run(fmt.Sprintf("reported count %d", count), func(t *testing.T) {
			got, err := supportedDevicesFromValuesForTest(nil, count)
			require.Nil(t, got)
			require.EqualError(t, err, fmt.Sprintf("invalid supported GPU count %d, expected 0 through %d", count, capacity))
		})
	}
}

// TestEntityGroupEntitiesFromNative checks entity conversion at valid native count boundaries.
func TestEntityGroupEntitiesFromNative(t *testing.T) {
	got, err := entityGroupEntitiesFromValuesForTest(nil, 0)
	require.NoError(t, err)
	require.Empty(t, got)

	got, err = entityGroupEntitiesFromValuesForTest([]uint{4, 8, 15}, 3)
	require.NoError(t, err)
	require.Equal(t, []uint{4, 8, 15}, got)

	capacity := entityGroupCapacityForTest()
	entityIDs := make([]uint, capacity)
	for i := range entityIDs {
		entityIDs[i] = uint(i + 1)
	}
	got, err = entityGroupEntitiesFromValuesForTest(entityIDs, capacity)
	require.NoError(t, err)
	require.Equal(t, entityIDs, got)
}

// TestEntityGroupEntitiesFromNativeLargeResponse checks a response larger than DCGM_GROUP_MAX_ENTITIES_V2 (1,024 entities).
func TestEntityGroupEntitiesFromNativeLargeResponse(t *testing.T) {
	const entityCount = 2048
	require.GreaterOrEqual(t, entityGroupCapacityForTest(), entityCount)

	entityIDs := make([]uint, entityCount)
	for i := range entityIDs {
		entityIDs[i] = uint(i + 100)
	}
	got, err := entityGroupEntitiesFromValuesForTest(entityIDs, entityCount)
	require.NoError(t, err)
	require.Len(t, got, entityCount)
	require.Equal(t, entityIDs, got)
}

// TestEntityGroupEntitiesFromNativeRejectsInvalidCounts checks negative and oversized native counts.
func TestEntityGroupEntitiesFromNativeRejectsInvalidCounts(t *testing.T) {
	capacity := entityGroupCapacityForTest()
	for _, count := range []int{-1, capacity + 1} {
		t.Run(fmt.Sprintf("reported count %d", count), func(t *testing.T) {
			got, err := entityGroupEntitiesFromValuesForTest(nil, count)
			require.Nil(t, got)
			require.EqualError(t, err, fmt.Sprintf("invalid entity count %d, expected 0 through %d", count, capacity))
		})
	}
}

// TestNvLinkP2PStatusFromNative checks native link-matrix conversion and its maximum valid size.
func TestNvLinkP2PStatusFromNative(t *testing.T) {
	statuses := [][]Link_State{
		{Link_State(0), Link_State(1)},
		{Link_State(2), Link_State(3)},
	}
	got, err := nvLinkP2PStatusFromValuesForTest(statuses, 2)
	require.NoError(t, err)
	require.Equal(t, NvLinkP2PStatus{numGpus: 2, Gpus: statuses}, got)

	capacity := nvLinkGPUCapacityForTest()
	got, err = nvLinkP2PStatusFromValuesForTest(nil, uint(capacity))
	require.NoError(t, err)
	require.Equal(t, uint(capacity), got.numGpus)
	require.Len(t, got.Gpus, capacity)
	for _, row := range got.Gpus {
		require.Len(t, row, capacity)
	}
}

// TestNvLinkP2PStatusFromNativeRejectsInvalidCounts checks oversized native GPU counts.
func TestNvLinkP2PStatusFromNativeRejectsInvalidCounts(t *testing.T) {
	capacity := nvLinkGPUCapacityForTest()
	for _, count := range []uint{uint(capacity + 1), uint(^uint32(0))} {
		t.Run(fmt.Sprintf("reported count %d", count), func(t *testing.T) {
			got, err := nvLinkP2PStatusFromValuesForTest(nil, count)
			require.Equal(t, NvLinkP2PStatus{}, got)
			require.EqualError(t, err, fmt.Sprintf("invalid NVLink GPU count %d, expected 0 through %d", count, capacity))
		})
	}
}

// TestDeviceAttributesFromNative checks conversion of every Device field sourced from native attributes.
func TestDeviceAttributesFromNative(t *testing.T) {
	want := deviceAttributes{
		busID:   "00000000:17:00.0",
		uuid:    "GPU-device-info-test",
		power:   325,
		bar1:    65536,
		fbTotal: 81920,
		identifiers: DeviceIdentifiers{
			Brand:               "NVIDIA",
			Model:               "Test GPU",
			Serial:              "serial-17",
			Vbios:               "vbios-17",
			InforomImageVersion: "inforom-17",
			DriverVersion:       "driver-17",
		},
	}

	require.Equal(t, want, deviceAttributesFromValuesForTest(want))
}
