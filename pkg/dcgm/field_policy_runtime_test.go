//go:build linux && cgo

package dcgm

import (
	"fmt"
	"path/filepath"
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
