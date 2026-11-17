//go:build linux && cgo

package dcgm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestFieldPolicyAdapterLifecycle exercises policy creation, readback, modification, and cleanup against DCGM.
func TestFieldPolicyAdapterLifecycle(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("DCGM_POLICY_ALLOWED_DIRS", configDir)
	t.Setenv("DCGM_POLICY_CONFIG_PATH", filepath.Join(configDir, "policies.yaml"))
	teardown := setupTest(t)
	defer teardown(t)
	runOnlyWithLiveGPUs(t)
	gpus, err := GetSupportedDevices()
	require.NoError(t, err)

	id, err := CreateFieldPolicy(FieldPolicySpec{
		FieldID:   DCGM_FI_DEV_GPU_TEMP,
		Name:      fmt.Sprintf("go-dcgm-%d", time.Now().UnixNano()),
		Threshold: 90,
		Operator:  FieldPolicyGreaterThan,
		Channels:  FieldPolicyChannelConsole,
		Entities:  []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpus[0]}},
		Enabled:   true,
	})
	if IsFeatureUnavailable(err) {
		t.Skipf("field-policy API unavailable: %v", err)
	}
	require.NoError(t, err)
	defer func() { require.NoError(t, DeleteFieldPolicy(DCGM_FI_DEV_GPU_TEMP, id)) }()

	policy, err := GetFieldPolicy(DCGM_FI_DEV_GPU_TEMP, id)
	require.NoError(t, err)
	require.Equal(t, id, policy.PolicyID)
	policy.Threshold = 85
	require.NoError(t, ModifyFieldPolicy(policy))
	policies, err := GetAllFieldPolicies()
	require.NoError(t, err)
	found := false
	for _, listed := range policies {
		if listed.PolicyID == id {
			require.Equal(t, float64(85), listed.Threshold)
			found = true
			break
		}
	}
	require.True(t, found, "created field policy is absent from the list")
}

// TestFieldPolicyUUIDImport verifies import atomicity, unresolved identity, and persistence.
func TestFieldPolicyUUIDImport(t *testing.T) {
	for _, format := range []string{"legacy-flat", "grouped"} {
		t.Run(format, func(t *testing.T) { testFieldPolicyUUIDImport(t, format == "grouped") })
	}
}

// fieldPolicyUUIDTargets renders UUID targets in either supported YAML shape.
func fieldPolicyUUIDTargets(grouped bool, uuids ...string) string {
	if grouped {
		return "{GPU: {componentUuids: [" + strings.Join(uuids, ", ") + "]}}"
	}
	targets := make([]string, len(uuids))
	for i, uuid := range uuids {
		targets[i] = "{entityGroupId: 1, componentUuid: " + uuid + "}"
	}
	return "[" + strings.Join(targets, ", ") + "]"
}

// testFieldPolicyUUIDImport checks import contracts and restores a saved policy after restart.
func testFieldPolicyUUIDImport(t *testing.T, grouped bool) {
	t.Helper()

	dir := t.TempDir()
	config := filepath.Join(dir, "saved.yaml")
	t.Setenv("DCGM_POLICY_ALLOWED_DIRS", dir)
	t.Setenv("DCGM_POLICY_CONFIG_PATH", config)
	// The checks inject field values, which DCGM 4.8 allows only when the engine starts with this set.
	t.Setenv("DCGM_ENABLE_INJECTION_APIS", "1")
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	running := true
	defer func() {
		if running {
			cleanup()
		}
	}()
	if !dcgmSymbolAvailable("dcgmPolicyImport_v2") {
		t.Skip("counted import requires a DCGM 4.8 library")
	}
	runOnlyWithLiveGPUs(t)
	gpus, err := GetSupportedDevices()
	require.NoError(t, err)
	device, err := GetDeviceInfo(gpus[0])
	require.NoError(t, err)

	valid := checkFieldPolicyImportAtomicity(t, dir, grouped, device.UUID, gpus[0])
	checkUnresolvedFieldPolicy(t, dir, grouped, gpus[0])

	// Save the valid policy and verify that a fresh session restores its UUID target.
	_, _, err = ImportFieldPoliciesWithCounts(writeFieldPolicyFixture(t, dir, "restore-valid.yaml", valid))
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		data, readErr := os.ReadFile(config)
		return readErr == nil && strings.Contains(string(data), device.UUID) &&
			strings.Contains(string(data), "componentUuids:") && strings.Contains(string(data), "major: 2") &&
			strings.Contains(string(data), "minor: 0")
	}, 5*time.Second, 50*time.Millisecond)
	require.NoError(t, Shutdown())
	running = false
	cleanup, err = Init(Embedded)
	require.NoError(t, err)
	running = true
	restarted, err := GetAllFieldPolicies()
	require.NoError(t, err)
	var restored bool
	for i := range restarted {
		policy := &restarted[i]
		if policy.Name != "uuid-target" {
			continue
		}
		require.Len(t, policy.Entities, 1)
		restoredDevice, deviceErr := GetDeviceInfo(policy.Entities[0].EntityId)
		require.NoError(t, deviceErr)
		require.Equal(t, device.UUID, restoredDevice.UUID)
		restored = true
	}

	require.True(t, restored, "saved UUID policy was not restored")
}

// fieldPolicyYAMLEntry renders one disabled policy fixture with the supplied targets.
func fieldPolicyYAMLEntry(name, target string) string {
	return fmt.Sprintf(`  - fieldId: %d
    name: %s
    operator: GREATER_THAN
    threshold: 90
    entities: %s
    channels: [CONSOLE]
    status: DISABLED
    rateLimitSeconds: 0
`, DCGM_FI_DEV_GPU_TEMP, name, target)
}

// writeFieldPolicyFixture saves a policy fixture inside the test's allowed directory.
func writeFieldPolicyFixture(t *testing.T, dir, name, body string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))

	return path
}

// checkFieldPolicyImportAtomicity checks UUID targeting and rejects invalid imports without changing policies.
func checkFieldPolicyImportAtomicity(t *testing.T, dir string, grouped bool, uuid string, gpuID uint) string {
	t.Helper()

	schema := "schema_version: {major: 2, minor: 0}\npolicies:\n"
	valid := schema + fieldPolicyYAMLEntry("uuid-target", fieldPolicyUUIDTargets(grouped, fmt.Sprintf("%q", uuid)))
	parsed, skipped, err := ImportFieldPoliciesWithCounts(writeFieldPolicyFixture(t, dir, "valid.yaml", valid))
	require.NoError(t, err)
	require.Equal(t, uint(1), parsed)
	require.Zero(t, skipped)
	before, err := GetAllFieldPolicies()
	require.NoError(t, err)
	require.Len(t, before, 1)
	require.Equal(t, []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpuID}}, before[0].Entities)
	for _, schema := range []int{1, 2} {
		body := fmt.Sprintf("schema_version: %d\npolicies:\n", schema) + fieldPolicyYAMLEntry("old-id", fmt.Sprintf("[{entityGroupId: 1, entityId: %d}]", gpuID))
		parsed, skipped, err = ImportFieldPoliciesWithCounts(writeFieldPolicyFixture(t, dir, fmt.Sprintf("old-%d.yaml", schema), body))
		requireDCGMCode(t, err, DCGM_ST_BADPARAM)
		require.Zero(t, parsed)
		require.Equal(t, uint(1), skipped)
		after, listErr := GetAllFieldPolicies()
		require.NoError(t, listErr)
		require.Equal(t, before, after)
	}

	invalidTarget := "[{entityGroupId: 1, entityId: 0}]"
	if grouped {
		invalidTarget = fmt.Sprintf("{GPU: {all: true, componentUuids: [%q]}}", uuid)
	}

	// DCGM replaces every policy on import, so a new name makes a partial import visible.
	mixed := schema + fieldPolicyYAMLEntry("mixed-atomic", fieldPolicyUUIDTargets(grouped, fmt.Sprintf("%q", uuid))) +
		fieldPolicyYAMLEntry("invalid", invalidTarget)
	parsed, skipped, err = ImportFieldPoliciesWithCounts(writeFieldPolicyFixture(t, dir, "mixed.yaml", mixed))
	requireDCGMCode(t, err, DCGM_ST_BADPARAM)
	require.Equal(t, uint(1), parsed)
	require.Equal(t, uint(1), skipped)
	after, err := GetAllFieldPolicies()
	require.NoError(t, err)
	require.Equal(t, before, after)

	return valid
}

// checkUnresolvedFieldPolicy checks that missing UUIDs keep their identity and never match a live GPU.
func checkUnresolvedFieldPolicy(t *testing.T, dir string, grouped bool, gpuID uint) {
	t.Helper()

	// Retaining all unknown targets preserves their private UUID identities.
	unknown := "GPU-00000000-0000-0000-0000-000000000001"
	unknown2 := "GPU-00000000-0000-0000-0000-000000000002"
	unresolved := "schema_version: {major: 2, minor: 0}\npolicies:\n" + fieldPolicyYAMLEntry("unresolved", fieldPolicyUUIDTargets(grouped, unknown, unknown2))
	parsed, skipped, err := ImportFieldPoliciesWithCounts(writeFieldPolicyFixture(t, dir, "unresolved.yaml", unresolved))
	require.NoError(t, err)
	require.Equal(t, uint(1), parsed)
	require.Zero(t, skipped)
	policies, err := GetAllFieldPolicies()
	require.NoError(t, err)
	var missing FieldPolicy
	for i := range policies {
		if policies[i].Name == "unresolved" {
			missing = policies[i]
		}
	}

	require.Equal(t, "unresolved", missing.Name)
	require.Len(t, missing.Entities, 2)
	for _, entity := range missing.Entities {
		require.Equal(t, UnresolvedFieldPolicyEntityID, entity.EntityId)
	}

	// Require a fresh save so the original import cannot satisfy the UUID check.
	config := filepath.Join(dir, "saved.yaml")
	require.NoError(t, os.Remove(config))
	missing.Threshold = 95
	require.NoError(t, ModifyFieldPolicy(missing))
	readback, err := GetFieldPolicy(missing.FieldID, missing.PolicyID)
	require.NoError(t, err)
	require.Equal(t, missing.Entities, readback.Entities)
	require.Equal(t, float64(95), readback.Threshold)
	require.Eventually(t, func() bool {
		data, readErr := os.ReadFile(config)
		return readErr == nil && strings.Contains(string(data), unknown) && strings.Contains(string(data), unknown2)
	}, 5*time.Second, 50*time.Millisecond)
	ambiguous := readback
	ambiguous.Entities = ambiguous.Entities[:1]
	requireDCGMCode(t, ModifyFieldPolicy(ambiguous), DCGM_ST_BADPARAM)
	retained, err := GetFieldPolicy(missing.FieldID, missing.PolicyID)
	require.NoError(t, err)
	require.Equal(t, readback, retained)
	_, err = CreateFieldPolicy(FieldPolicySpec{
		FieldID:  missing.FieldID,
		Name:     "new-unresolved",
		Operator: FieldPolicyGreaterThan,
		Entities: missing.Entities,
	})
	requireDCGMCode(t, err, DCGM_ST_BADPARAM)

	// An unresolved policy must not be evaluated against a different live GPU.
	readback.Enabled = true
	require.NoError(t, ModifyFieldPolicy(readback))
	require.NoError(t, InjectFieldValue(gpuID, DCGM_FI_DEV_GPU_TEMP, DCGM_FT_INT64, DCGM_ST_OK, time.Now().UnixMicro(), int64(150)))
	require.NoError(t, UpdateAllFields())
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		state, readErr := GetFieldPolicy(missing.FieldID, missing.PolicyID)
		require.NoError(t, readErr)
		require.Zero(t, state.ViolationCount, "unresolved target was evaluated against a live GPU")

		time.Sleep(100 * time.Millisecond)
	}
}

// requireDCGMCode checks the native status wrapped by a Go error.
func requireDCGMCode(t *testing.T, err error, code int) {
	t.Helper()

	var native *Error
	require.ErrorAs(t, err, &native)
	require.Equal(t, code, int(native.Code))
}

// TestFieldPolicyRediscovery changes driver attachment only on explicitly selected hosts.
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
	// The checks inject field values, which DCGM 4.8 allows only when the engine starts with this set.
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

// requireRediscoveredPolicyWatch checks that a rebound entity produces fresh policy violations.
func requireRediscoveredPolicyWatch(t *testing.T, restored FieldPolicy) {
	t.Helper()
	injectedAt := time.UnixMicro(time.Now().UnixMicro())
	entity := restored.Entities[0]
	// Injection notifies existing watches directly; a forced refresh could replace the sample.
	require.NoError(t, InjectFieldValue(entity.EntityId, restored.FieldID, DCGM_FT_INT64, DCGM_ST_OK, injectedAt.UnixMicro(), int64(250)))
	require.Eventually(t, func() bool {
		policy, readErr := GetFieldPolicy(restored.FieldID, restored.PolicyID)
		return readErr == nil && policy.ViolationCount > restored.ViolationCount &&
			policy.LastTriggeredTimestamp.After(restored.LastTriggeredTimestamp) && !policy.LastTriggeredTimestamp.Before(injectedAt)
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

// TestCountedPolicyImportOlderLibrary checks both the public guard and native dispatch.
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

// TestCountedPolicyImportAfterFailedInit checks that a loaded library without a session cannot import.
func TestCountedPolicyImportAfterFailedInit(t *testing.T) {
	requireLoadedLibraryWithoutSession(t)
	if !dcgmSymbolAvailable("dcgmPolicyImport_v2") {
		t.Skip("counted policy import session regression requires the DCGM 4.8 symbol")
	}
	parsed, skipped, err := ImportFieldPoliciesWithCounts("policies.yaml")
	require.Zero(t, parsed)
	require.Zero(t, skipped)
	var unavailable *FeatureUnavailableError
	require.ErrorAs(t, err, &unavailable)
	require.Equal(t, "dcgmPolicyImport_v2", unavailable.Symbol)
}
