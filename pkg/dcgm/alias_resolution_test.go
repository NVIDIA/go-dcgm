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
