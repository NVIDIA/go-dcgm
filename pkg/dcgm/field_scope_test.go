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
	"fmt"
	"testing"
)

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
			name:      "earlier Go name for the ratio field is a legacy name for field 1613",
			fieldName: "DCGM_FI_DEV_GPU_UTIL_RATIO_V2",
			wantID:    1613,
			current:   false,
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

	if DCGM_FI_DEV_GPU_UTIL_RATIO != 1613 {
		t.Fatalf("DCGM_FI_DEV_GPU_UTIL_RATIO = %d, want the ratio field 1613", DCGM_FI_DEV_GPU_UTIL_RATIO)
	}
	if DCGM_FI_DEV_GPU_UTIL_RATIO_V2 != DCGM_FI_DEV_GPU_UTIL_RATIO {
		t.Fatal("DCGM_FI_DEV_GPU_UTIL_RATIO_V2 must remain an alias of DCGM_FI_DEV_GPU_UTIL_RATIO")
	}
	if !IsLegacyField("DCGM_FI_DEV_GPU_UTIL_RATIO_V2") {
		t.Fatal("DCGM_FI_DEV_GPU_UTIL_RATIO_V2 must remain a compatibility lookup for field 1613")
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

// linkFieldExpectations lists the binding's link-capable fields in DCGM 4.8.0 metadata
// (dcgmlib/src/dcgm_fields.cpp), plus the GPU-declared fields accepted by
// ReadAndCacheNvLinkFields and DcgmFieldIsNvLinkCountField. Parent expectations
// come from the readers' parent checks and NVSwitch module ownership, not from
// the compatibility helper. Field additions outside this binding are checked
// when its generated field definitions are updated.
func linkFieldExpectations() map[Short]struct{ declared, parent Field_Entity_Group } {
	return map[Short]struct{ declared, parent Field_Entity_Group }{
		DCGM_FI_DEV_NVLINK_LINK_STATUS:                              {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_IBPC_TX_WAIT_TOTAL:                 {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_PER_LINK_TOTAL:            {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_PER_LINK_TOTAL:            {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_REPLAY_ERROR_PER_LINK_TOTAL:              {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_PER_LINK_TOTAL:            {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_THROUGHPUT_PER_LINK:                      {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_TX_THROUGHPUT_PER_LINK:                   {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_THROUGHPUT_PER_LINK:                   {FE_LINK, FE_GPU},
		DCGM_FI_PROF_NVLINK_TX_BYTES_PER_LINK:                       {FE_LINK, FE_GPU},
		DCGM_FI_PROF_NVLINK_RX_BYTES_PER_LINK:                       {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPRM_OPER_RECOVERY:                       {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_SINCE_LAST_SECONDS:        {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_TIME_BETWEEN_LAST_TWO:     {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_SUCCESSFUL_TOTAL:          {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PHYSICAL_RECOVERY_SUCCESSFUL_TOTAL: {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PHYSICAL_LINK_DOWN_TOTAL:           {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PLR_RX_CODE_TOTAL:                  {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PLR_RX_CODE_ERROR_TOTAL:            {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PLR_RX_CODE_UNCORRECTABLE_TOTAL:    {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PLR_TX_CODE_TOTAL:                  {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PLR_TX_RETRY_CODE_TOTAL:            {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PLR_TX_RETRY_EVENT_TOTAL:           {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_PLR_SYNC_EVENT_TOTAL:               {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_LAST_GAP_SECONDS:          {FE_LINK, FE_GPU},
		DCGM_FI_DEV_NVSWITCH_LINK_ID:                                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_STATUS:                            {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_TYPE:                              {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION:              {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_LINK_ID:                    {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_LINK_SID:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_UUID:                                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX:                     {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX:                     {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS:                      {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS:                  {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERROR_TOTAL:              {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERROR_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERROR_TOTAL:                  {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_TOTAL:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_TOTAL:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3:                   {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0:                  {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1:                  {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2:                  {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3:                  {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0:                 {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1:                 {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2:                 {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3:                 {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC0_TOTAL:          {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC1_TOTAL:          {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC2_TOTAL:          {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC3_TOTAL:          {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L0_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L1_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L2_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L3_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L4_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L5_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L6_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L7_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L0_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L1_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L2_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L3_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L4_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L5_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L6_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L7_TOTAL:                {FE_LINK, FE_SWITCH},
		DCGM_FI_PROF_NVLINK_TX_BYTES:                                {FE_GPU, FE_GPU},
		DCGM_FI_PROF_NVLINK_RX_BYTES:                                {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_TX_PACKET_TOTAL:                          {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_TX_BYTES_TOTAL:                           {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_PACKET_TOTAL:                          {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_BYTES_TOTAL:                           {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_PACKET_MALFORMED_TOTAL:                {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_PACKET_DROPPED_TOTAL:                  {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_ERROR_TOTAL:                           {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_REMOTE_ERROR_TOTAL:                    {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_GENERAL_ERROR_TOTAL:                   {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_INTEGRITY_ERROR_TOTAL:                    {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_TX_DISCARD_TOTAL:                         {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RECOVERY_SUCCESSFUL_TOTAL:                {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RECOVERY_FAILED_TOTAL:                    {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RECOVERY_EVENT_TOTAL:                     {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_RX_SYMBOL_ERROR_TOTAL:                    {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_SYMBOL_BER_RAW:                           {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_SYMBOL_BER_RATIO:                         {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_EFFECTIVE_BER_RAW:                        {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_EFFECTIVE_BER_RATIO:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_EFFECTIVE_ERROR_TOTAL:                    {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_0_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_1_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_2_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_3_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_4_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_5_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_6_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_7_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_8_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_9_TOTAL:                      {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_10_TOTAL:                     {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_11_TOTAL:                     {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_12_TOTAL:                     {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_13_TOTAL:                     {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_14_TOTAL:                     {FE_GPU, FE_GPU},
		DCGM_FI_DEV_NVLINK_FEC_HISTORY_15_TOTAL:                     {FE_GPU, FE_GPU},
	}
}

// TestFieldSupportsEntityGroupWithParent checks allowed and forbidden pairs for
// every inventoried field, including preservation of each declared scope.
func TestFieldSupportsEntityGroupWithParent(t *testing.T) {
	for fieldID, want := range linkFieldExpectations() {
		for _, parent := range []Field_Entity_Group{FE_GPU, FE_SWITCH, FE_NONE, FE_GPU_I, FE_GPU_CI, FE_CPU, FE_LINK, FE_COUNT, 999} {
			t.Run(fmt.Sprintf("field_%d/parent_%d", fieldID, parent), func(t *testing.T) {
				got := FieldSupportsEntityGroupWithParent(fieldID, want.declared, FE_LINK, parent)
				if got != (parent == want.parent) {
					t.Fatalf("link field %d with parent %s: got %t, want %t", fieldID, parent, got, parent == want.parent)
				}
			})
		}

		if got := FieldSupportsEntityGroupWithParent(fieldID, want.declared, FE_GPU, FE_NONE); got != (want.declared == FE_GPU) {
			t.Errorf("field %d GPU support = %t, want %t", fieldID, got, want.declared == FE_GPU)
		}
	}
}

// TestLinkFieldInventory rejects missing or extra classifications for every
// field known to this binding, including fields outside the link inventory.
func TestLinkFieldInventory(t *testing.T) {
	inventory := linkFieldExpectations()
	for _, fieldID := range dcgmFields {
		want, listed := inventory[fieldID]
		for _, declared := range []Field_Entity_Group{FE_LINK, FE_GPU} {
			for _, parent := range []Field_Entity_Group{FE_GPU, FE_SWITCH} {
				expected := listed && want.declared == declared && want.parent == parent
				if got := FieldSupportsEntityGroupWithParent(fieldID, declared, FE_LINK, parent); got != expected {
					t.Errorf("field %d declared %s parent %s: got %t, want %t", fieldID, declared, parent, got, expected)
				}
			}
		}
	}

	for fieldID := range inventory {
		if !isKnownFieldID(fieldID) {
			t.Errorf("inventoried field %d is unknown to the binding", fieldID)
		}
	}
}

// TestFieldParentCompatibilityOtherScopes checks global, unrelated and unknown
// fields and verifies that adding a parent leaves non-link behavior unchanged.
func TestFieldParentCompatibilityOtherScopes(t *testing.T) {
	tests := []struct {
		name                     string
		fieldID                  Short
		declared, entity, parent Field_Entity_Group
		want                     bool
	}{
		{"global GPU link", DCGM_FI_DRIVER_VERSION, FE_NONE, FE_LINK, FE_GPU, true},
		{"global switch link", DCGM_FI_DRIVER_VERSION, FE_NONE, FE_LINK, FE_SWITCH, true},
		{"global unknown parent", DCGM_FI_DRIVER_VERSION, FE_NONE, FE_LINK, FE_NONE, false},
		{"unknown global field", 65535, FE_NONE, FE_LINK, FE_GPU, false},
		{"unknown link field", 65535, FE_LINK, FE_LINK, FE_SWITCH, false},
		{"unrelated GPU field", DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_LINK, FE_GPU, false},
		{"unclassified declared link", DCGM_FI_DEV_GPU_TEMP, FE_LINK, FE_LINK, FE_GPU, false},
		{"switch declared field", DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX, FE_SWITCH, FE_LINK, FE_SWITCH, false},
		{"declared GPU", DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_GPU, 999, true},
		{"declared instance", DCGM_FI_DEV_XID_ERRORS, FE_GPU_CI, FE_GPU_CI, FE_NONE, true},
		{"unrelated scope", DCGM_FI_DEV_GPU_TEMP, FE_GPU, FE_CPU, FE_GPU, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FieldSupportsEntityGroupWithParent(tt.fieldID, tt.declared, tt.entity, tt.parent); got != tt.want {
				t.Fatalf("got %t, want %t", got, tt.want)
			}
		})
	}

	for _, fieldID := range dcgmFields {
		for _, declared := range []Field_Entity_Group{FE_NONE, FE_GPU, FE_GPU_I, FE_GPU_CI, FE_SWITCH, FE_LINK} {
			for _, entity := range []Field_Entity_Group{FE_NONE, FE_GPU, FE_GPU_I, FE_GPU_CI, FE_SWITCH, FE_CPU} {
				if got, want := FieldSupportsEntityGroupWithParent(fieldID, declared, entity, 999), FieldSupportsEntityGroup(fieldID, declared, entity); got != want {
					t.Errorf("non-link compatibility changed for field %d declared %s target %s", fieldID, declared, entity)
				}
			}
		}
	}
}
