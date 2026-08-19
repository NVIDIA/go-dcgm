/*
 * Copyright (c) 2025, NVIDIA CORPORATION.  All rights reserved.
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

// Pure lookup tests against the generated dcgmFields / legacyDCGMFields maps.
// No setupTest(t): these do not require a DCGM runtime. Precedent:
// TestGetFieldIDBindUnbindEvent in bind_unbind_test.go.

// Deprecated alias resolves via legacyDCGMFields to the same ID as its
// canonical THROUGHPUT counterpart.
func TestAliasResolution_NVLinkBandwidthTotalResolvesToThroughputTotal(t *testing.T) {
	aliasID, aliasOK := GetFieldID("DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL")
	require.True(t, aliasOK, "deprecated alias must be resolvable")

	canonicalID, canonicalOK := GetFieldID("DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL")
	require.True(t, canonicalOK, "canonical field must be resolvable")

	assert.Equal(t, canonicalID, aliasID,
		"deprecated alias must resolve to the same ID as its canonical target")
}

// The deprecated alias lives in the legacy map only.
func TestAliasResolution_NVLinkBandwidthTotalIsLegacyOnly(t *testing.T) {
	assert.True(t, IsLegacyField("DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL"),
		"deprecated alias must be marked legacy")
	assert.False(t, IsCurrentField("DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL"),
		"deprecated alias must not appear in the current-fields map")
}

// The canonical field is current, not legacy.
func TestAliasResolution_CanonicalThroughputTotalIsCurrent(t *testing.T) {
	assert.True(t, IsCurrentField("DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL"),
		"canonical field must be marked current")
	assert.False(t, IsLegacyField("DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL"),
		"canonical field must not appear in the legacy map")
}

// Range sentinels like DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST are alias-style
// #defines but are NOT deprecated. They must not be resolvable via the
// public lookup API. A future change that naively resolves all aliases
// would break this assertion.
func TestAliasResolution_RangeSentinelIsNotResolvable(t *testing.T) {
	_, ok := GetFieldID("DCGM_FI_PROF_NVLINK_THROUGHPUT_FIRST")
	assert.False(t, ok,
		"range sentinel (non-deprecated alias) must not resolve as a field name")
}

// Outside-the-#ifdef-block deprecated alias (CLOCK_THROTTLE_REASONS)
// resolves via legacyDCGMFields to its canonical counterpart. Pins the
// "Deprecated:" comment-heuristic end-to-end.
func TestAliasResolution_ClockThrottleReasonsResolvesToClocksEventReasons(t *testing.T) {
	aliasID, aliasOK := GetFieldID("DCGM_FI_DEV_CLOCK_THROTTLE_REASONS")
	require.True(t, aliasOK, "outside-block deprecated alias must be resolvable")

	canonicalID, canonicalOK := GetFieldID("DCGM_FI_DEV_CLOCKS_EVENT_REASONS")
	require.True(t, canonicalOK, "canonical field must be resolvable")

	assert.Equal(t, canonicalID, aliasID)
	assert.True(t, IsLegacyField("DCGM_FI_DEV_CLOCK_THROTTLE_REASONS"),
		"outside-block deprecated alias must be marked legacy")
}

// GetFieldIDOrPanic must resolve a deprecated alias without panicking
// (the existing panic path is reserved for unknown names).
func TestAliasResolution_GetFieldIDOrPanicDoesNotPanicOnAlias(t *testing.T) {
	assert.NotPanics(t, func() {
		id := GetFieldIDOrPanic("DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL")
		assert.Equal(t, GetFieldIDOrPanic("DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL"), id,
			"GetFieldIDOrPanic on alias must return same ID as its canonical counterpart")
	})
}

func TestAliasResolution_LowercaseLegacyFieldsResolve(t *testing.T) {
	tests := []struct {
		legacy   string
		current  string
		expected Short
	}{
		{
			legacy:   "dcgm_gpu_temp",
			current:  "DCGM_FI_DEV_GPU_TEMP_CELSIUS",
			expected: DCGM_FI_DEV_GPU_TEMP_CELSIUS,
		},
		{
			legacy:   "dcgm_power_usage",
			current:  "DCGM_FI_DEV_BOARD_POWER_WATTS",
			expected: DCGM_FI_DEV_BOARD_POWER_WATTS,
		},
		{
			legacy:   "dcgm_xid_errors",
			current:  "DCGM_FI_DEV_XID_ERROR",
			expected: DCGM_FI_DEV_XID_ERROR,
		},
		{
			legacy:   "dcgm_nvlink_bandwidth_total",
			current:  "DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL",
			expected: DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.legacy, func(t *testing.T) {
			legacyID, legacyOK := GetFieldID(tt.legacy)
			require.True(t, legacyOK, "lowercase legacy field must be resolvable")

			currentID, currentOK := GetFieldID(tt.current)
			require.True(t, currentOK, "canonical field must be resolvable")

			assert.Equal(t, tt.expected, currentID)
			assert.Equal(t, currentID, legacyID)
			assert.True(t, IsLegacyField(tt.legacy), "lowercase legacy field must be marked legacy")
			assert.False(t, IsCurrentField(tt.legacy), "lowercase legacy field must not be marked current")
			assert.True(t, IsCurrentField(tt.current), "canonical field must be marked current")
		})
	}
}

func TestAliasResolution_HeaderDeprecatedAliasesAreLegacyOnlyAndSourceCompatible(t *testing.T) {
	tests := []struct {
		deprecated   string
		current      string
		deprecatedID Short
		currentID    Short
	}{
		{
			deprecated:   "DCGM_FI_DEV_GPU_TEMP",
			current:      "DCGM_FI_DEV_GPU_TEMP_CELSIUS",
			deprecatedID: DCGM_FI_DEV_GPU_TEMP,
			currentID:    DCGM_FI_DEV_GPU_TEMP_CELSIUS,
		},
		{
			deprecated:   "DCGM_FI_DEV_POWER_USAGE",
			current:      "DCGM_FI_DEV_BOARD_POWER_WATTS",
			deprecatedID: DCGM_FI_DEV_POWER_USAGE,
			currentID:    DCGM_FI_DEV_BOARD_POWER_WATTS,
		},
		{
			deprecated:   "DCGM_FI_DEV_XID_ERRORS",
			current:      "DCGM_FI_DEV_XID_ERROR",
			deprecatedID: DCGM_FI_DEV_XID_ERRORS,
			currentID:    DCGM_FI_DEV_XID_ERROR,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL",
			current:      "DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL,
			currentID:    DCGM_FI_DEV_NVLINK_THROUGHPUT_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_PROF_GR_ENGINE_ACTIVE",
			current:      "DCGM_FI_PROF_GR_ENGINE_UTIL_RATIO",
			deprecatedID: DCGM_FI_PROF_GR_ENGINE_ACTIVE,
			currentID:    DCGM_FI_PROF_GR_ENGINE_UTIL_RATIO,
		},
		{
			deprecated:   "DCGM_FI_DEV_CPU_UTIL_TOTAL",
			current:      "DCGM_FI_DEV_CPU_UTIL_RATIO",
			deprecatedID: DCGM_FI_DEV_CPU_UTIL_TOTAL,
			currentID:    DCGM_FI_DEV_CPU_UTIL_RATIO,
		},
		{
			deprecated:   "DCGM_FI_DEV_CPU_UTIL_USER",
			current:      "DCGM_FI_DEV_CPU_UTIL_USER_RATIO",
			deprecatedID: DCGM_FI_DEV_CPU_UTIL_USER,
			currentID:    DCGM_FI_DEV_CPU_UTIL_USER_RATIO,
		},
		{
			deprecated:   "DCGM_FI_DEV_CPU_UTIL_NICE",
			current:      "DCGM_FI_DEV_CPU_UTIL_NICE_RATIO",
			deprecatedID: DCGM_FI_DEV_CPU_UTIL_NICE,
			currentID:    DCGM_FI_DEV_CPU_UTIL_NICE_RATIO,
		},
		{
			deprecated:   "DCGM_FI_DEV_CPU_UTIL_SYS",
			current:      "DCGM_FI_DEV_CPU_UTIL_SYS_RATIO",
			deprecatedID: DCGM_FI_DEV_CPU_UTIL_SYS,
			currentID:    DCGM_FI_DEV_CPU_UTIL_SYS_RATIO,
		},
		{
			deprecated:   "DCGM_FI_DEV_CPU_UTIL_IRQ",
			current:      "DCGM_FI_DEV_CPU_UTIL_IRQ_RATIO",
			deprecatedID: DCGM_FI_DEV_CPU_UTIL_IRQ,
			currentID:    DCGM_FI_DEV_CPU_UTIL_IRQ_RATIO,
		},
		{
			deprecated:   "DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT",
			current:      "DCGM_FI_DEV_SYSIO_POWER_WATTS",
			deprecatedID: DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT,
			currentID:    DCGM_FI_DEV_SYSIO_POWER_WATTS,
		},
		{
			deprecated:   "DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT",
			current:      "DCGM_FI_DEV_MODULE_POWER_WATTS",
			deprecatedID: DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT,
			currentID:    DCGM_FI_DEV_MODULE_POWER_WATTS,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS",
			current:      "DCGM_FI_DEV_NVLINK_TX_DISCARD_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS,
			currentID:    DCGM_FI_DEV_NVLINK_TX_DISCARD_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS",
			current:      "DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERROR_STATUS",
			deprecatedID: DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS,
			currentID:    DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERROR_STATUS,
		},
		{
			deprecated:   "DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK",
			current:      "DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERROR_MASK",
			deprecatedID: DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK,
			currentID:    DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERROR_MASK,
		},
		{
			deprecated:   "DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE",
			current:      "DCGM_FI_DEV_CONNECTX_TEMP_CELSIUS",
			deprecatedID: DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE,
			currentID:    DCGM_FI_DEV_CONNECTX_TEMP_CELSIUS,
		},
		{
			deprecated:   "DCGM_FI_DEV_C2C_LINK_ERROR_INTR",
			current:      "DCGM_FI_DEV_C2C_LINK_ERROR_TOTAL",
			deprecatedID: DCGM_FI_DEV_C2C_LINK_ERROR_INTR,
			currentID:    DCGM_FI_DEV_C2C_LINK_ERROR_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY",
			current:      "DCGM_FI_DEV_C2C_LINK_REPLAY_ERROR_TOTAL",
			deprecatedID: DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY,
			currentID:    DCGM_FI_DEV_C2C_LINK_REPLAY_ERROR_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY_B2B",
			current:      "DCGM_FI_DEV_C2C_LINK_REPLAY_ERROR_B2B_TOTAL",
			deprecatedID: DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY_B2B,
			currentID:    DCGM_FI_DEV_C2C_LINK_REPLAY_ERROR_B2B_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_C2C_LINK_POWER_STATE",
			current:      "DCGM_FI_DEV_C2C_LINK_POWER_STATUS",
			deprecatedID: DCGM_FI_DEV_C2C_LINK_POWER_STATE,
			currentID:    DCGM_FI_DEV_C2C_LINK_POWER_STATUS,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_0",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_0_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_0,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_0_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_1",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_1_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_1,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_1_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_2",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_2_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_2,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_2_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_3",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_3_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_3,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_3_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_4",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_4_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_4,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_4_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_5",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_5_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_5,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_5_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_6",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_6_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_6,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_6_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_7",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_7_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_7,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_7_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_8",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_8_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_8,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_8_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_9",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_9_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_9,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_9_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_10",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_10_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_10,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_10_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_11",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_11_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_11,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_11_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_12",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_12_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_12,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_12_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_13",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_13_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_13,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_13_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_14",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_14_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_14,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_14_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_15",
			current:      "DCGM_FI_DEV_NVLINK_FEC_HISTORY_15_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_15,
			currentID:    DCGM_FI_DEV_NVLINK_FEC_HISTORY_15_TOTAL,
		},
		{
			deprecated:   "DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_CEIL",
			current:      "DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_CEIL_WATTS",
			deprecatedID: DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_CEIL,
			currentID:    DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_CEIL_WATTS,
		},
		{
			deprecated:   "DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_FLOOR",
			current:      "DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_FLOOR_WATTS",
			deprecatedID: DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_FLOOR,
			currentID:    DCGM_FI_DEV_PWR_SMOOTHING_APPLIED_TMP_FLOOR_WATTS,
		},
		{
			deprecated:   "DCGM_FI_DEV_PWR_SMOOTHING_MAX_NUM_PRESET_PROFILES",
			current:      "DCGM_FI_DEV_PWR_SMOOTHING_MAX_PRESET_PROFILES",
			deprecatedID: DCGM_FI_DEV_PWR_SMOOTHING_MAX_NUM_PRESET_PROFILES,
			currentID:    DCGM_FI_DEV_PWR_SMOOTHING_MAX_PRESET_PROFILES,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_GET_STATE",
			current:      "DCGM_FI_DEV_NVLINK_LINK_STATUS",
			deprecatedID: DCGM_FI_DEV_NVLINK_GET_STATE,
			currentID:    DCGM_FI_DEV_NVLINK_LINK_STATUS,
		},
		{
			deprecated:   "DCGM_FI_DEV_NVLINK_PPCNT_IBPC_PORT_XMIT_WAIT",
			current:      "DCGM_FI_DEV_NVLINK_PPCNT_IBPC_TX_WAIT_TOTAL",
			deprecatedID: DCGM_FI_DEV_NVLINK_PPCNT_IBPC_PORT_XMIT_WAIT,
			currentID:    DCGM_FI_DEV_NVLINK_PPCNT_IBPC_TX_WAIT_TOTAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.deprecated, func(t *testing.T) {
			deprecatedID, deprecatedOK := GetFieldID(tt.deprecated)
			require.True(t, deprecatedOK, "deprecated header alias must be resolvable")

			currentID, currentOK := GetFieldID(tt.current)
			require.True(t, currentOK, "canonical target must be resolvable")

			assert.Equal(t, currentID, deprecatedID)
			assert.Equal(t, tt.currentID, tt.deprecatedID, "deprecated Go const must match canonical const")
			assert.True(t, IsLegacyField(tt.deprecated), "deprecated header alias must be marked legacy")
			assert.False(t, IsCurrentField(tt.deprecated), "deprecated header alias must not be marked current")
			assert.True(t, IsCurrentField(tt.current), "canonical target must be marked current")
			assert.False(t, IsLegacyField(tt.current), "canonical target must not be marked legacy")
		})
	}
}

func TestAliasResolution_DCGM46RenamedFieldsResolve(t *testing.T) {
	tests := []struct {
		legacy    string
		current   string
		legacyID  Short
		currentID Short
	}{
		{
			legacy:    "DCGM_FI_DEV_GPU_TEMP",
			current:   "DCGM_FI_DEV_GPU_TEMP_CELSIUS",
			legacyID:  DCGM_FI_DEV_GPU_TEMP,
			currentID: DCGM_FI_DEV_GPU_TEMP_CELSIUS,
		},
		{
			legacy:    "DCGM_FI_DEV_POWER_USAGE",
			current:   "DCGM_FI_DEV_BOARD_POWER_WATTS",
			legacyID:  DCGM_FI_DEV_POWER_USAGE,
			currentID: DCGM_FI_DEV_BOARD_POWER_WATTS,
		},
		{
			legacy:    "DCGM_FI_DEV_XID_ERRORS",
			current:   "DCGM_FI_DEV_XID_ERROR",
			legacyID:  DCGM_FI_DEV_XID_ERRORS,
			currentID: DCGM_FI_DEV_XID_ERROR,
		},
		{
			legacy:    "DCGM_FI_PROF_GR_ENGINE_ACTIVE",
			current:   "DCGM_FI_PROF_GR_ENGINE_UTIL_RATIO",
			legacyID:  DCGM_FI_PROF_GR_ENGINE_ACTIVE,
			currentID: DCGM_FI_PROF_GR_ENGINE_UTIL_RATIO,
		},
		{
			legacy:    "DCGM_FI_PROF_PIPE_TENSOR_ACTIVE",
			current:   "DCGM_FI_PROF_TENSOR_UTIL_RATIO",
			legacyID:  DCGM_FI_PROF_PIPE_TENSOR_ACTIVE,
			currentID: DCGM_FI_PROF_TENSOR_UTIL_RATIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.legacy, func(t *testing.T) {
			legacyID, legacyOK := GetFieldID(tt.legacy)
			require.True(t, legacyOK, "legacy field must be resolvable")

			currentID, currentOK := GetFieldID(tt.current)
			require.True(t, currentOK, "current field must be resolvable")

			assert.Equal(t, currentID, legacyID)
			assert.Equal(t, tt.currentID, tt.legacyID, "Go alias constant must match canonical constant")
			assert.True(t, IsLegacyField(tt.legacy), "legacy field must be marked legacy")
			assert.True(t, IsCurrentField(tt.current), "current field must be marked current")
		})
	}
}

func TestAliasResolution_DCGM46ProfilingTotalsResolve(t *testing.T) {
	tests := []struct {
		name string
		id   Short
	}{
		{"DCGM_FI_PROF_SM_CYCLES_ELAPSED_TOTAL", 1084},
		{"DCGM_FI_PROF_SM_CYCLES_ACTIVE_TOTAL", 1085},
		{"DCGM_FI_PROF_MMA_CYCLES_ACTIVE_TOTAL", 1086},
		{"DCGM_FI_PROF_DMMA_CYCLES_ACTIVE_TOTAL", 1087},
		{"DCGM_FI_PROF_HMMA_CYCLES_ACTIVE_TOTAL", 1088},
		{"DCGM_FI_PROF_IMMA_CYCLES_ACTIVE_TOTAL", 1089},
		{"DCGM_FI_PROF_DFMA_CYCLES_ACTIVE_TOTAL", 1090},
		{"DCGM_FI_PROF_PCIE_TX_BYTES_TOTAL", 1091},
		{"DCGM_FI_PROF_PCIE_RX_BYTES_TOTAL", 1092},
		{"DCGM_FI_PROF_INT_CYCLES_ACTIVE_TOTAL", 1093},
		{"DCGM_FI_PROF_FP64_CYCLES_ACTIVE_TOTAL", 1094},
		{"DCGM_FI_PROF_FP32_CYCLES_ACTIVE_TOTAL", 1095},
		{"DCGM_FI_PROF_FP16_CYCLES_ACTIVE_TOTAL", 1096},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, ok := GetFieldID(tt.name)
			require.True(t, ok, "new DCGM 4.6 profiling field must be resolvable")
			assert.Equal(t, tt.id, id)
			assert.True(t, IsCurrentField(tt.name), "new profiling field must be current")
		})
	}
}
