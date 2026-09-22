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

import "testing"

func TestDCGM47GPUUtilFieldSemantics(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		wantID    Short
		current   bool
	}{
		{
			name:      "percentage GPU utilization remains field 203",
			fieldName: "DCGM_FI_DEV_GPU_UTIL",
			wantID:    203,
			current:   true,
		},
		{
			name:      "new ratio GPU utilization is field 1613",
			fieldName: "DCGM_FI_DEV_GPU_UTIL_RATIO_V2",
			wantID:    1613,
			current:   true,
		},
		{
			name:      "current header ratio GPU utilization is field 1613",
			fieldName: "DCGM_FI_DEV_GPU_UTIL_RATIO",
			wantID:    1613,
			current:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldID, ok := GetFieldID(tt.fieldName)
			if !ok {
				t.Fatalf("GetFieldID(%q) did not find the field", tt.fieldName)
			}
			if fieldID != tt.wantID {
				t.Fatalf("GetFieldID(%q) = %d, want %d", tt.fieldName, fieldID, tt.wantID)
			}
			if got := IsCurrentField(tt.fieldName); got != tt.current {
				t.Fatalf("IsCurrentField(%q) = %t, want %t", tt.fieldName, got, tt.current)
			}
		})
	}

	if DCGM_FI_DEV_GPU_UTIL == DCGM_FI_DEV_GPU_UTIL_RATIO_V2 {
		t.Fatal("percentage and new ratio GPU utilization fields must have different IDs")
	}
	if DCGM_FI_DEV_GPU_UTIL_RATIO != DCGM_FI_DEV_GPU_UTIL {
		t.Fatal("legacy GPU utilization ratio constant must remain field 203")
	}
	if !IsLegacyField("DCGM_FI_DEV_GPU_UTIL") {
		t.Fatal("DCGM_FI_DEV_GPU_UTIL must remain a compatibility lookup for field 203")
	}
	if IsLegacyField("DCGM_FI_DEV_GPU_UTIL_RATIO") {
		t.Fatal("DCGM_FI_DEV_GPU_UTIL_RATIO must resolve to the current DCGM header field")
	}
}

func TestFieldSupportsEntityGroup(t *testing.T) {
	type fieldScopeCase struct {
		name          string
		fieldID       Short
		declaredLevel Field_Entity_Group
		entityGroup   Field_Entity_Group
		want          bool
	}

	tests := []fieldScopeCase{
		{
			name:          "declared link field remains available to links",
			fieldID:       DCGM_FI_DEV_NVLINK_LINK_STATUS,
			declaredLevel: FE_LINK,
			entityGroup:   FE_LINK,
			want:          true,
		},
		{
			name:          "declared link field is not promoted to GPUs",
			fieldID:       DCGM_FI_DEV_NVLINK_LINK_STATUS,
			declaredLevel: FE_LINK,
			entityGroup:   FE_GPU,
			want:          false,
		},
		{
			name:          "overloaded NVLink TX profiling bytes remains available to GPUs",
			fieldID:       DCGM_FI_PROF_NVLINK_TX_BYTES,
			declaredLevel: FE_GPU,
			entityGroup:   FE_GPU,
			want:          true,
		},
		{
			name:          "overloaded NVLink TX profiling bytes is available to links",
			fieldID:       DCGM_FI_PROF_NVLINK_TX_BYTES,
			declaredLevel: FE_GPU,
			entityGroup:   FE_LINK,
			want:          true,
		},
		{
			name:          "overloaded NVLink RX profiling bytes remains available to GPUs",
			fieldID:       DCGM_FI_PROF_NVLINK_RX_BYTES,
			declaredLevel: FE_GPU,
			entityGroup:   FE_GPU,
			want:          true,
		},
		{
			name:          "overloaded NVLink RX profiling bytes is available to links",
			fieldID:       DCGM_FI_PROF_NVLINK_RX_BYTES,
			declaredLevel: FE_GPU,
			entityGroup:   FE_LINK,
			want:          true,
		},
		{
			name:          "unrelated GPU field is not promoted to links",
			fieldID:       DCGM_FI_DEV_GPU_TEMP_CELSIUS,
			declaredLevel: FE_GPU,
			entityGroup:   FE_LINK,
			want:          false,
		},
		{
			name:          "field after the NVLink COUNT range is not promoted to links",
			fieldID:       DCGM_FI_DEV_NVLINK_ECC_ERROR_TOTAL,
			declaredLevel: FE_GPU,
			entityGroup:   FE_LINK,
			want:          false,
		},
		{
			name:          "field before the NVLink FEC COUNT range is not promoted to links",
			fieldID:       DCGM_FI_DEV_C2C_LINK_POWER_STATUS,
			declaredLevel: FE_GPU,
			entityGroup:   FE_LINK,
			want:          false,
		},
		{
			name:          "field after the NVLink FEC COUNT range is not promoted to links",
			fieldID:       DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_POWER_CAP_NS,
			declaredLevel: FE_GPU,
			entityGroup:   FE_LINK,
			want:          false,
		},
		{
			name:          "unknown field is rejected",
			fieldID:       Short(65535),
			declaredLevel: FE_GPU,
			entityGroup:   FE_GPU,
			want:          false,
		},
	}

	for fieldID := DCGM_FI_DEV_NVLINK_TX_PACKET_TOTAL; fieldID <= DCGM_FI_DEV_NVLINK_EFFECTIVE_ERROR_TOTAL; fieldID++ {
		tests = append(tests,
			fieldScopeCase{
				name:          "NVLink COUNT field retains GPU scope",
				fieldID:       fieldID,
				declaredLevel: FE_GPU,
				entityGroup:   FE_GPU,
				want:          true,
			},
			fieldScopeCase{
				name:          "NVLink COUNT field is available to links",
				fieldID:       fieldID,
				declaredLevel: FE_GPU,
				entityGroup:   FE_LINK,
				want:          true,
			},
		)
	}
	for fieldID := DCGM_FI_DEV_NVLINK_FEC_HISTORY_0_TOTAL; fieldID <= DCGM_FI_DEV_NVLINK_FEC_HISTORY_15_TOTAL; fieldID++ {
		tests = append(tests,
			fieldScopeCase{
				name:          "NVLink FEC COUNT field retains GPU scope",
				fieldID:       fieldID,
				declaredLevel: FE_GPU,
				entityGroup:   FE_GPU,
				want:          true,
			},
			fieldScopeCase{
				name:          "NVLink FEC COUNT field is available to links",
				fieldID:       fieldID,
				declaredLevel: FE_GPU,
				entityGroup:   FE_LINK,
				want:          true,
			},
		)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FieldSupportsEntityGroup(tt.fieldID, tt.declaredLevel, tt.entityGroup)
			if got != tt.want {
				t.Fatalf("FieldSupportsEntityGroup(%d, %s, %s) = %t, want %t", tt.fieldID, tt.declaredLevel, tt.entityGroup, got, tt.want)
			}
		})
	}
}
