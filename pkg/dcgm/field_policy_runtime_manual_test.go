//go:build linux && cgo

package dcgm

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestFieldPolicyRediscovery checks that a GPU policy still works after the
// driver is detached and reattached.
func TestFieldPolicyRediscovery(t *testing.T) {
	if os.Getenv("DCGM_TEST_POLICY_REDISCOVERY") != "1" {
		t.Skip("requires DCGM_TEST_POLICY_REDISCOVERY=1 on a dedicated GPU host")
	}
	if os.Geteuid() != 0 {
		t.Skip("driver rediscovery requires root")
	}

	dir := t.TempDir()
	t.Setenv("DCGM_POLICY_ALLOWED_DIRS", dir)
	t.Setenv("DCGM_POLICY_CONFIG_PATH", filepath.Join(dir, "saved.yaml"))
	// DCGM 4.8 reads this setting at startup to allow field injection.
	t.Setenv("DCGM_ENABLE_INJECTION_APIS", "1")

	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()

	if !dcgmSymbolAvailable("dcgmPolicyImport_v2") {
		t.Skip("UUID policy rediscovery requires DCGM 4.8")
	}
	runOnlyWithLiveGPUs(t)
	gpus, err := GetSupportedDevices()
	require.NoError(t, err)
	original, err := GetDeviceInfo(gpus[0])
	require.NoError(t, err)

	id, err := CreateFieldPolicy(FieldPolicySpec{
		FieldID:   DCGM_FI_DEV_GPU_TEMP,
		Name:      "rediscovery",
		Threshold: 200,
		Operator:  FieldPolicyGreaterThan,
		Entities:  []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpus[0]}},
		Channels:  FieldPolicyChannelConsole,
		Enabled:   true,
	})
	require.NoError(t, err)
	defer func() { require.NoError(t, DeleteFieldPolicy(DCGM_FI_DEV_GPU_TEMP, id)) }()

	detached := false
	defer func() {
		if detached {
			require.NoError(t, AttachDriver())
		}
	}()

	require.NoError(t, DetachDriver())
	detached = true

	missing, err := GetFieldPolicy(DCGM_FI_DEV_GPU_TEMP, id)
	require.NoError(t, err)
	require.Len(t, missing.Entities, 1)
	require.Equal(t, UnresolvedFieldPolicyEntityID, missing.Entities[0].EntityId)
	require.Zero(t, missing.ViolationCount)

	require.NoError(t, AttachDriver())
	detached = false
	require.Eventually(t, func() bool {
		policy, readErr := GetFieldPolicy(DCGM_FI_DEV_GPU_TEMP, id)
		return readErr == nil && len(policy.Entities) == 1 && policy.Entities[0].EntityId != UnresolvedFieldPolicyEntityID
	}, 10*time.Second, 100*time.Millisecond)

	restored, err := GetFieldPolicy(DCGM_FI_DEV_GPU_TEMP, id)
	require.NoError(t, err)
	current, err := GetDeviceInfo(restored.Entities[0].EntityId)
	require.NoError(t, err)
	require.Equal(t, original.UUID, current.UUID)
	requireRediscoveredPolicyWatch(t, restored)
}

// requireRediscoveredPolicyWatch checks that the reattached GPU can trigger a new policy violation.
func requireRediscoveredPolicyWatch(t *testing.T, restored FieldPolicy) {
	t.Helper()
	injectedAt := time.UnixMicro(time.Now().UnixMicro())
	entity := restored.Entities[0]

	// A forced refresh could overwrite the injected sample.
	require.NoError(t, InjectFieldValue(entity.EntityId, restored.FieldID, DCGM_FT_INT64, DCGM_ST_OK, injectedAt.UnixMicro(), int64(250)))

	require.Eventually(t, func() bool {
		policy, readErr := GetFieldPolicy(restored.FieldID, restored.PolicyID)
		return readErr == nil &&
			policy.ViolationCount > restored.ViolationCount &&
			policy.LastTriggeredTimestamp.After(restored.LastTriggeredTimestamp) &&
			!policy.LastTriggeredTimestamp.Before(injectedAt)
	}, 10*time.Second, 100*time.Millisecond, "policy watch did not record a fresh violation after reattachment")

	var recovered FieldPolicyViolation
	require.Eventually(t, func() bool {
		violations, readErr := GetFieldPolicyViolations(injectedAt)
		if readErr != nil {
			return false
		}
		for _, violation := range violations {
			if violation.PolicyID == restored.PolicyID && violation.FieldID == restored.FieldID && violation.Entity == entity {
				recovered = violation
				return true
			}
		}
		return false
	}, 10*time.Second, 100*time.Millisecond, "restored policy watch did not emit a violation for the rebound entity")

	require.False(t, recovered.Timestamp.Before(injectedAt))
	require.Equal(t, float64(250), recovered.CurrentValue)
	require.Equal(t, restored.Threshold, recovered.Threshold)
	require.Equal(t, restored.Operator, recovered.Operator)
}

// TestCountedPolicyImportOlderLibrary checks how the Go API and C adapter report a missing import function.
func TestCountedPolicyImportOlderLibrary(t *testing.T) {
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()

	if dcgmSymbolAvailable("dcgmPolicyImport_v2") {
		t.Skip("requires a library without dcgmPolicyImport_v2")
	}

	parsed, skipped, err := ImportFieldPoliciesWithCounts("unused-policy-file.yaml")
	require.Zero(t, parsed)
	require.Zero(t, skipped)
	require.True(t, IsFeatureUnavailable(err))

	parsed, skipped, status := cgoAdapter{}.importFieldPoliciesNative("unused-policy-file.yaml")
	require.Zero(t, parsed)
	require.Zero(t, skipped)
	require.Equal(t, DCGM_ST_FUNCTION_NOT_FOUND, status)
}
