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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCPUHierarchyV1Conversion checks CPU IDs and owned-core bitmasks in native hierarchy conversion.
func TestCPUHierarchyV1Conversion(t *testing.T) {
	lastBitmaskIndex := int(MAX_CPU_CORE_BITMASK_COUNT - 1)
	tests := []struct {
		name      string
		cpus      []testCPUHierarchyV1CPU
		wantIDs   []uint
		wantMasks []map[int]uint64
	}{
		{name: "empty"},
		{
			name: "multiple CPUs and boundary masks",
			cpus: []testCPUHierarchyV1CPU{
				{cpuID: 7, ownedCoreBitmasks: map[int]uint64{0: 0b101, lastBitmaskIndex: 0b10}},
				{cpuID: 8, ownedCoreBitmasks: map[int]uint64{0: 0b10000}},
			},
			wantIDs: []uint{7, 8},
			wantMasks: []map[int]uint64{
				{0: 0b101, lastBitmaskIndex: 0b10},
				{0: 0b10000},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hierarchy := createTestCPUHierarchyV1(tt.cpus)

			assert.Equal(t, uint(testCPUHierarchyVersion1), hierarchy.Version)
			assert.Equal(t, uint(len(tt.wantIDs)), hierarchy.NumCPUs)
			if len(tt.wantIDs) == 0 {
				assert.Equal(t, CPUHierarchyCPU_v1{}, hierarchy.CPUs[0])
				return
			}

			for i, wantID := range tt.wantIDs {
				assert.Equal(t, wantID, hierarchy.CPUs[i].CPUID)
				require.Len(t, hierarchy.CPUs[i].OwnedCores, int(MAX_CPU_CORE_BITMASK_COUNT))
				for index, wantMask := range tt.wantMasks[i] {
					assert.Equal(t, wantMask, hierarchy.CPUs[i].OwnedCores[index])
				}
				assert.Zero(t, hierarchy.CPUs[i].OwnedCores[1])
			}
		})
	}
}

func TestCPUHierarchyV2IncludesSerial(t *testing.T) {
	lastBitmaskIndex := int(MAX_CPU_CORE_BITMASK_COUNT - 1)
	hierarchy := createTestCPUHierarchyV2([]testCPUHierarchyV2CPU{
		{
			cpuID:  7,
			serial: "GRACE-SERIAL-0007",
			ownedCoreBitmasks: map[int]uint64{
				0:                0b101,
				1:                0b1000,
				lastBitmaskIndex: 0b10,
			},
		},
		{
			cpuID:  8,
			serial: "",
			ownedCoreBitmasks: map[int]uint64{
				0: 0b10000,
			},
		},
	})

	assert.Equal(t, uint(testCPUHierarchyVersion2), hierarchy.Version)
	assert.Equal(t, uint(2), hierarchy.NumCPUs)
	assert.Equal(t, uint(7), hierarchy.CPUs[0].CPUID)
	assert.Len(t, hierarchy.CPUs[0].OwnedCores, int(MAX_CPU_CORE_BITMASK_COUNT))
	assert.Equal(t, uint64(0b101), hierarchy.CPUs[0].OwnedCores[0])
	assert.Equal(t, uint64(0b1000), hierarchy.CPUs[0].OwnedCores[1])
	assert.Equal(t, uint64(0b10), hierarchy.CPUs[0].OwnedCores[lastBitmaskIndex])
	assert.Equal(t, "GRACE-SERIAL-0007", hierarchy.CPUs[0].Serial)

	assert.Equal(t, uint(8), hierarchy.CPUs[1].CPUID)
	assert.Len(t, hierarchy.CPUs[1].OwnedCores, int(MAX_CPU_CORE_BITMASK_COUNT))
	assert.Equal(t, uint64(0b10000), hierarchy.CPUs[1].OwnedCores[0])
	assert.Empty(t, hierarchy.CPUs[1].Serial)
}

func TestCPUHierarchyV2HandlesEmptyHierarchy(t *testing.T) {
	hierarchy := createTestCPUHierarchyV2(nil)

	assert.Equal(t, uint(testCPUHierarchyVersion2), hierarchy.Version)
	assert.Equal(t, uint(0), hierarchy.NumCPUs)
}

func TestCPUHierarchyErrorWrapsDCGMError(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		result    testDCGMReturn
	}{
		{
			name:      "v1 hierarchy error",
			operation: "error retrieving DCGM CPU hierarchy",
			result:    testDCGMStatusVersionMismatch,
		},
		{
			name:      "v2 hierarchy error",
			operation: "error retrieving DCGM CPU hierarchy v2",
			result:    testDCGMStatusFunctionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cpuHierarchyError(tt.operation, &Error{msg: "dcgm error", Code: tt.result})

			require.Error(t, err)
			assert.ErrorContains(t, err, tt.operation)

			var dcgmErr *Error
			require.ErrorAs(t, err, &dcgmErr)
			assert.Equal(t, tt.result, dcgmErr.Code)
		})
	}
}
