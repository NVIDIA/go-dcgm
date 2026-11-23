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
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestGetDeviceInfoWithOps checks Device fields and read errors.
// It also checks that extra reads require a supported GPU with EntityStatusOk.
func TestGetDeviceInfoWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")
	affinityErr := errors.New("affinity cleanup failed")
	bandwidthErr := errors.New("bandwidth cleanup failed")
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

	supportedDevice := Device{
		GPU: 3, DCGMSupported: "Yes", UUID: attrs.uuid, Power: attrs.power,
		PCI:         PCIInfo{BusID: attrs.busID, BAR1: attrs.bar1, FBTotal: attrs.fbTotal, Bandwidth: 32000},
		Identifiers: attrs.identifiers, Topology: topology, CPUAffinity: "{0-7}",
	}

	tests := []deviceInfoTestCase{
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
		{name: "affinity cleanup and topology failure", cleanup: true, fail: "topology", gpus: []uint{3}, status: EntityStatusOk, wantError: true},
		{name: "bandwidth failure", fail: "bandwidth", gpus: []uint{3}, status: EntityStatusOk, wantError: true},
		{name: "success", gpus: []uint{3}, status: EntityStatusOk, want: supportedDevice},
		{name: "successful reads with cleanup failures", gpus: []uint{3}, status: EntityStatusOk, want: supportedDevice, cleanup: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockDeviceInfoOps(gomock.NewController(t))
			expectDeviceInfoReads(mock, tt, attrs, topology, stageErr, affinityErr, bandwidthErr)

			got, err := getDeviceInfoWithOps(mock, 3)
			if tt.wantError {
				require.ErrorIs(t, err, stageErr)
				if tt.cleanup {
					require.ErrorIs(t, err, affinityErr)
				}
				require.Equal(t, Device{}, got)
				return
			}
			if tt.cleanup {
				require.ErrorIs(t, err, affinityErr)
				require.ErrorIs(t, err, bandwidthErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

type deviceInfoTestCase struct {
	name      string
	fail      string
	gpus      []uint
	status    EntityStatus
	want      Device
	wantError bool
	cleanup   bool
}

// expectDeviceInfoReads expects the reads up to the selected failure or the end
// of the successful path. Gomock fails if a later read runs after an error.
func expectDeviceInfoReads(
	mock *MockDeviceInfoOps, tt deviceInfoTestCase, attrs deviceAttributes, topology []P2PLink,
	stageErr, affinityErr, bandwidthErr error,
) {
	attributes := mock.EXPECT().readDeviceAttributes(uint(3))
	if tt.fail == "attributes" {
		attributes.Return(deviceAttributes{}, stageErr)
		return
	}
	attributes.Return(attrs, nil)

	supported := mock.EXPECT().getSupportedDevices().After(attributes)
	if tt.fail == "supported" {
		supported.Return(nil, stageErr)
		return
	}
	supported.Return(tt.gpus, nil)

	status := mock.EXPECT().getGPUStatus(uint(3)).Return(tt.status).After(supported)
	if !slices.Contains(tt.gpus, uint(3)) || tt.status != EntityStatusOk {
		return
	}
	affinity := mock.EXPECT().getCPUAffinity(uint(3)).After(status)
	if tt.fail == "affinity" {
		affinity.Return("", nil, stageErr)
		return
	}
	affinity.Return("{0-7}", pick(tt.cleanup, affinityErr, nil), nil)

	links := mock.EXPECT().getDeviceTopology(uint(3)).After(affinity)
	if tt.fail == "topology" {
		links.Return(nil, stageErr)
		return
	}
	links.Return(topology, nil)

	bandwidth := mock.EXPECT().getPciBandwidth(uint(3)).After(links)
	if tt.fail == "bandwidth" {
		bandwidth.Return(int64(0), nil, stageErr)
		return
	}
	bandwidth.Return(int64(32000), pick(tt.cleanup, bandwidthErr, nil), nil)
}

// TestGetPciBandwidthWithOps checks setup failures, value reads, and temporary-resource cleanup.
func TestGetPciBandwidthWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")
	cleanupErr := errors.New("cleanup failed")
	pciFields := []Short{DCGM_FI_DEV_PCIE_MAX_LINK_GEN, DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH}
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
			values := derivedInt64Samples(pciFields, 4, 16)
			expectDerivedFieldRead(mock, 3, pciFields, values, derivedReadErrors{
				stage: tt.fail, operation: stageErr,
				fieldCleanup: pick(tt.fail != "", cleanupErr, nil),
				groupCleanup: pick(tt.fail != "", cleanupErr, nil),
			})

			got, cleanup, err := getPciBandwidthWithOps(mock, 3)
			switch tt.fail {
			case "create", "watch":
				require.ErrorIs(t, err, stageErr)
				if tt.fail == "watch" {
					require.ErrorIs(t, cleanup, cleanupErr)
				}
				require.Zero(t, got)
			case "collect":
				require.ErrorIs(t, err, stageErr)
				require.ErrorIs(t, cleanup, cleanupErr)
				require.Zero(t, got)
			default:
				require.NoError(t, err)
				require.NoError(t, cleanup)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// TestGetCPUAffinityWithOps checks setup failures, bitmask formatting, and deferred cleanup.
func TestGetCPUAffinityWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")
	affinityFields := []Short{
		DCGM_FI_DEV_CPU_AFFINITY_0,
		DCGM_FI_DEV_CPU_AFFINITY_1,
		DCGM_FI_DEV_CPU_AFFINITY_2,
		DCGM_FI_DEV_CPU_AFFINITY_3,
	}
	tests := []struct {
		name   string
		fail   string
		values []int64
		want   string
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
			name:   "preserves a mask overlapping the numeric blank range",
			values: []int64{0x7fffffffffffffff, 0, 0, 0},
			want:   "{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockDeviceStatusOps(gomock.NewController(t))
			values := derivedInt64Samples(affinityFields, tt.values...)
			expectDerivedFieldRead(mock, 5, affinityFields, values, derivedReadErrors{stage: tt.fail, operation: stageErr})

			got, cleanup, err := getCPUAffinityWithOps(mock, 5)
			require.NoError(t, cleanup)
			switch tt.fail {
			case "create", "watch":
				require.ErrorIs(t, err, stageErr)
				require.Equal(t, "N/A", got)
			case "collect":
				require.ErrorIs(t, err, stageErr)
				require.EqualError(t, err, "get CPU affinity for GPU 5: stage failed")
				require.Equal(t, "N/A", got)
			default:
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

// TestDerivedFieldsValidateRequiredSamples checks that each read validates fields
// by ID, status and type before interpreting the values.
func TestDerivedFieldsValidateRequiredSamples(t *testing.T) {
	for _, affinity := range []bool{false, true} {
		for _, shape := range []string{"missing", "wrong type", "failed", "blank sentinel bit pattern", "out of order"} {
			t.Run(pick(affinity, "CPU affinity/", "PCIe bandwidth/")+shape, func(t *testing.T) {
				checkRequiredDerivedSample(t, affinity, shape)
			})
		}
	}
}

func checkRequiredDerivedSample(t *testing.T, affinity bool, shape string) {
	t.Helper()
	fields := []Short{DCGM_FI_DEV_PCIE_MAX_LINK_GEN, DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH}
	integers := []int64{4, 16}
	if affinity {
		fields = []Short{DCGM_FI_DEV_CPU_AFFINITY_0, DCGM_FI_DEV_CPU_AFFINITY_1, DCGM_FI_DEV_CPU_AFFINITY_2, DCGM_FI_DEV_CPU_AFFINITY_3}
		integers = []int64{1, 2, 4, 8}
	}
	values := requiredSampleWithShape(derivedInt64Samples(fields, integers...), shape)
	m := NewMockDeviceStatusOps(gomock.NewController(t))
	expectDerivedFieldRead(m, 5, fields, values, derivedReadErrors{})
	var got, want any
	var cleanup, err error
	if affinity {
		got, cleanup, err = getCPUAffinityWithOps(m, 5)
		want = "{0,65,130,195}"
	} else {
		got, cleanup, err = getPciBandwidthWithOps(m, 5)
		want = int64(31504)
	}
	require.NoError(t, cleanup)

	// A numeric sentinel can also be a valid CPU affinity bitmask.
	if shape == "out of order" || (affinity && shape == "blank sentinel bit pattern") {
		require.NoError(t, err)
		if shape == "blank sentinel bit pattern" {
			want = "{4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,65,130,195}"
		}
		require.Equal(t, want, got)
		return
	}
	require.Error(t, err)
	if affinity {
		require.Equal(t, "N/A", got)
	} else {
		require.Equal(t, int64(0), got)
	}
	if shape == "failed" {
		var native *Error
		require.ErrorAs(t, err, &native)
		require.Equal(t, DCGM_ST_NO_DATA, int(native.Code))
	}
}

// requiredSampleWithShape changes one property of a valid response for each case.
func requiredSampleWithShape(values []FieldValue_v1, shape string) []FieldValue_v1 {
	switch shape {
	case "missing":
		values = values[:1]
	case "wrong type":
		values[0].FieldType = DCGM_FT_STRING
	case "failed":
		values[0].Status = DCGM_ST_NO_DATA
	case "blank sentinel bit pattern":
		fieldID := values[0].FieldID
		values[0] = testInt64FieldValue(DCGM_FT_INT64_BLANK)
		values[0].FieldID = fieldID
	case "out of order":
		slices.Reverse(values)
	}
	return values
}

func derivedInt64Samples(fields []Short, integers ...int64) []FieldValue_v1 {
	values := make([]FieldValue_v1, len(integers))
	for i, value := range integers {
		values[i] = testInt64FieldValue(value)
		values[i].FieldID = fields[i]
	}
	return values
}

type derivedReadErrors struct {
	stage        string
	operation    error
	fieldCleanup error
	groupCleanup error
}

// expectDerivedFieldRead keeps acquisition, read and reverse-order cleanup expectations together.
func expectDerivedFieldRead(m *MockDeviceStatusOps, gpuID uint, fields []Short, values []FieldValue_v1, failures derivedReadErrors) {
	var fieldGroup FieldHandle
	fieldGroup.SetHandle(17)
	var group GroupHandle
	group.SetHandle(23)
	create := m.EXPECT().fieldGroupCreate(gomock.Any(), fields)
	if failures.stage == "create" {
		create.Return(FieldHandle{}, failures.operation)
		return
	}
	create.Return(fieldGroup, nil)
	watch := m.EXPECT().watchFields(gpuID, fieldGroup, gomock.Any()).After(create)
	if failures.stage == "watch" {
		watch.Return(GroupHandle{}, failures.operation)
		m.EXPECT().fieldGroupDestroy(fieldGroup).Return(failures.fieldCleanup).After(watch)
		return
	}
	watch.Return(group, nil)
	collect := m.EXPECT().getLatestValuesForFields(gpuID, fields).After(watch)
	if failures.stage == "collect" {
		collect.Return(nil, failures.operation)
	} else {
		collect.Return(values, nil)
	}
	destroy := m.EXPECT().destroyGroup(group).Return(failures.groupCleanup).After(collect)
	m.EXPECT().fieldGroupDestroy(fieldGroup).Return(failures.fieldCleanup).After(destroy)
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

// TestDeviceMemoryCapacity checks valid MiB values and rejects blank, negative, and overflowing capacities.
func TestDeviceMemoryCapacity(t *testing.T) {
	for _, tc := range []struct {
		name      string
		field     Short
		value     int64
		wantError string
	}{
		{name: "zero BAR1", field: DCGM_FI_DEV_BAR1_TOTAL},
		{name: "framebuffer MiB", field: DCGM_FI_DEV_FB_TOTAL, value: 8192},
		{name: "largest legacy capacity", field: DCGM_FI_DEV_FB_TOTAL, value: math.MaxUint32},
		{name: "blank BAR1", field: DCGM_FI_DEV_BAR1_TOTAL, value: DCGM_FT_INT64_BLANK, wantError: "blank capacity"},
		{name: "unsupported framebuffer", field: DCGM_FI_DEV_FB_TOTAL, value: DCGM_FT_INT64_NOT_SUPPORTED, wantError: "blank capacity"},
		{name: "negative capacity", field: DCGM_FI_DEV_FB_TOTAL, value: -1, wantError: "out of range"},
		{name: "capacity overflow", field: DCGM_FI_DEV_FB_TOTAL, value: math.MaxUint32 + 1, wantError: "out of range"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			value := FieldValue_v2{FieldID: tc.field, EntityGroupId: FE_GPU, EntityID: 7, FieldType: DCGM_FT_INT64}
			_, err := binary.Encode(value.Value[:], binary.NativeEndian, tc.value)
			require.NoError(t, err)

			capacity, err := deviceMemoryCapacity(&value, tc.field, 7)
			if tc.wantError != "" {
				require.ErrorContains(t, err, tc.wantError)
				require.Zero(t, capacity)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tc.value, capacity)
		})
	}
}

// TestReadLegacyDeviceMemoryFields checks the field request and what the fallback writes.
func TestReadLegacyDeviceMemoryFields(t *testing.T) {
	queryErr := errors.New("cannot read fields")
	bar1, framebuffer := DCGM_FI_DEV_BAR1_TOTAL, DCGM_FI_DEV_FB_TOTAL
	capacities := map[Short]int64{bar1: 256, framebuffer: 8192}
	for _, tc := range []struct {
		name                      string
		bar1, framebuffer         uint
		fields                    []Short
		reply                     string
		wantBAR1, wantFramebuffer uint
	}{
		{name: "both absent", fields: []Short{bar1, framebuffer}, wantBAR1: 256, wantFramebuffer: 8192},
		{name: "BAR1 absent", framebuffer: 8192, fields: []Short{bar1}, wantBAR1: 256, wantFramebuffer: 8192},
		{name: "framebuffer absent", bar1: 256, fields: []Short{framebuffer}, wantBAR1: 256, wantFramebuffer: 8192},
		{name: "query failed", fields: []Short{bar1, framebuffer}, reply: "transport"},
		{name: "missing value", fields: []Short{bar1, framebuffer}, reply: "short"},
		{name: "wrong GPU", fields: []Short{bar1, framebuffer}, reply: "wrong GPU"},
		{name: "wrong type", fields: []Short{bar1, framebuffer}, reply: "wrong type"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			attrs := deviceAttributes{bar1: tc.bar1, fbTotal: tc.framebuffer}
			err := readLegacyDeviceMemoryFields(7, &attrs, func(entities []GroupEntityPair, fields []Short, flags uint) ([]FieldValue_v2, error) {
				require.Equal(t, []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: 7}}, entities)
				require.Equal(t, tc.fields, fields)
				require.Equal(t, DCGM_FV_FLAG_LIVE_DATA, flags)

				values := make([]FieldValue_v2, len(fields))
				for i, field := range fields {
					values[i] = FieldValue_v2{FieldID: field, EntityGroupId: FE_GPU, EntityID: 7, FieldType: DCGM_FT_INT64}
					_, err := binary.Encode(values[i].Value[:], binary.NativeEndian, capacities[field])
					require.NoError(t, err)
				}
				switch tc.reply {
				case "transport":
					return nil, queryErr
				case "short":
					return values[:1], nil
				case "wrong GPU":
					values[0].EntityID = 8
				case "wrong type":
					values[0].FieldType = DCGM_FT_DOUBLE
				}
				return values, nil
			})

			switch tc.reply {
			case "":
				require.NoError(t, err)
			case "transport":
				require.ErrorIs(t, err, queryErr)
			case "short":
				require.ErrorContains(t, err, "returned 1 of 2 memory fields")
			case "wrong GPU":
				require.ErrorContains(t, err, "for GPU 7")
			case "wrong type":
				require.ErrorContains(t, err, "expected an integer value")
			}
			require.Equal(t, tc.wantBAR1, attrs.bar1)
			require.Equal(t, tc.wantFramebuffer, attrs.fbTotal)
		})
	}
}
