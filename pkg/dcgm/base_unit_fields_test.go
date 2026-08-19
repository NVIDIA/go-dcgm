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

type baseUnitFieldCase struct {
	name       string
	fieldID    Short
	wantID     Short
	legacyName string
	legacyID   Short
}

var baseUnitFieldCases = []baseUnitFieldCase{
	{"DCGM_FI_DEV_BAR1_CAPACITY_BYTES", DCGM_FI_DEV_BAR1_CAPACITY_BYTES, 1600, "DCGM_FI_DEV_BAR1_TOTAL", 90},
	{"DCGM_FI_DEV_BAR1_USED_BYTES", DCGM_FI_DEV_BAR1_USED_BYTES, 1601, "DCGM_FI_DEV_BAR1_USED", 92},
	{"DCGM_FI_DEV_BAR1_FREE_BYTES", DCGM_FI_DEV_BAR1_FREE_BYTES, 1602, "DCGM_FI_DEV_BAR1_FREE", 93},
	{"DCGM_FI_DEV_SM_CLOCK_HERTZ", DCGM_FI_DEV_SM_CLOCK_HERTZ, 1603, "DCGM_FI_DEV_SM_CLOCK", 100},
	{"DCGM_FI_DEV_MEMORY_CLOCK_HERTZ", DCGM_FI_DEV_MEMORY_CLOCK_HERTZ, 1604, "DCGM_FI_DEV_MEM_CLOCK", 101},
	{"DCGM_FI_DEV_VIDEO_CLOCK_HERTZ", DCGM_FI_DEV_VIDEO_CLOCK_HERTZ, 1605, "DCGM_FI_DEV_VIDEO_CLOCK", 102},
	{"DCGM_FI_DEV_SM_APP_CLOCK_HERTZ", DCGM_FI_DEV_SM_APP_CLOCK_HERTZ, 1606, "DCGM_FI_DEV_APP_SM_CLOCK", 110},
	{"DCGM_FI_DEV_MEMORY_APP_CLOCK_HERTZ", DCGM_FI_DEV_MEMORY_APP_CLOCK_HERTZ, 1607, "DCGM_FI_DEV_APP_MEM_CLOCK", 111},
	{"DCGM_FI_DEV_SM_MAX_CLOCK_HERTZ", DCGM_FI_DEV_SM_MAX_CLOCK_HERTZ, 1608, "DCGM_FI_DEV_MAX_SM_CLOCK", 113},
	{"DCGM_FI_DEV_MEMORY_MAX_CLOCK_HERTZ", DCGM_FI_DEV_MEMORY_MAX_CLOCK_HERTZ, 1609, "DCGM_FI_DEV_MAX_MEM_CLOCK", 114},
	{"DCGM_FI_DEV_VIDEO_MAX_CLOCK_HERTZ", DCGM_FI_DEV_VIDEO_MAX_CLOCK_HERTZ, 1610, "DCGM_FI_DEV_MAX_VIDEO_CLOCK", 115},
	{"DCGM_FI_DEV_GPU_ENERGY_JOULES_TOTAL", DCGM_FI_DEV_GPU_ENERGY_JOULES_TOTAL, 1611, "DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION", 156},
	{"DCGM_FI_DEV_FAN_SPEED_RATIO", DCGM_FI_DEV_FAN_SPEED_RATIO, 1612, "DCGM_FI_DEV_FAN_SPEED", 191},
}

func TestBaseUnitFieldCanonicalNamesResolve(t *testing.T) {
	for _, tt := range baseUnitFieldCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantID, tt.fieldID, "exported Go constant must retain the upstream ID")

			got, ok := GetFieldID(tt.name)
			require.True(t, ok, "canonical base-unit name must resolve")
			assert.Equal(t, tt.wantID, got)
			assert.True(t, IsCurrentField(tt.name))
			assert.False(t, IsLegacyField(tt.name))
		})
	}
}

func TestBaseUnitFieldLegacyNumericNamesRemainIndependent(t *testing.T) {
	for _, tt := range baseUnitFieldCases {
		t.Run(tt.legacyName, func(t *testing.T) {
			got, ok := GetFieldID(tt.legacyName)
			require.True(t, ok, "legacy numeric field name must remain resolvable")
			assert.Equal(t, tt.legacyID, got)
			assert.NotEqual(t, tt.fieldID, got, "different units must retain different field IDs")
			assert.True(t, IsCurrentField(tt.legacyName), "numeric field definitions remain in the current map")
			assert.False(t, IsLegacyField(tt.legacyName), "only aliases and curated names use the legacy map")
		})
	}
}

func TestBaseUnitFieldLowercaseLegacyNamesResolve(t *testing.T) {
	tests := []struct {
		name string
		id   Short
	}{
		{"dcgm_sm_clock", DCGM_FI_DEV_SM_CLOCK},
		{"dcgm_memory_clock", DCGM_FI_DEV_MEM_CLOCK},
		{"dcgm_total_energy_consumption", DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetFieldID(tt.name)
			require.True(t, ok, "curated lowercase legacy name must remain resolvable")
			assert.Equal(t, tt.id, got)
			assert.True(t, IsLegacyField(tt.name))
			assert.False(t, IsCurrentField(tt.name))
		})
	}
}
