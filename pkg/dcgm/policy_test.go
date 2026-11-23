//go:build linux && cgo

/*
 * Copyright (c) 2023, NVIDIA CORPORATION.  All rights reserved.
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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAndGetPolicy(t *testing.T) {
	t.Log("Initializing DCGM in Embedded mode...")
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()
	t.Log("DCGM initialized successfully")

	group := GroupAllGPUs()
	t.Logf("Created group handle for all GPUs: %+v", group)

	// Check how many GPUs we have
	gpuCount, err := GetAllDeviceCount()
	require.NoError(t, err)
	t.Logf("Found %d GPU(s) in the system", gpuCount)

	action := PolicyActionNone
	validation := PolicyValidationNone

	// Test cases for each policy type
	testCases := []struct {
		name        string
		config      PolicyConfig
		expected    interface{}
		conditionID PolicyCondition
	}{
		{
			name: "ThermalPolicy",
			config: PolicyConfig{
				Condition:      ThermalPolicy,
				Action:         &action,
				Validation:     &validation,
				MaxTemperature: ptrUint32(85),
			},
			expected:    uint32(85),
			conditionID: ThermalPolicy,
		},
		{
			name: "PowerPolicy",
			config: PolicyConfig{
				Condition:  PowerPolicy,
				Action:     &action,
				Validation: &validation,
				MaxPower:   ptrUint32(300),
			},
			expected:    uint32(300),
			conditionID: PowerPolicy,
		},
		{
			name: "MaxRtPgPolicy",
			config: PolicyConfig{
				Condition:       MaxRtPgPolicy,
				Action:          &action,
				Validation:      &validation,
				MaxRetiredPages: ptrUint32(15),
			},
			expected:    uint32(15),
			conditionID: MaxRtPgPolicy,
		},
		{
			name: "DbePolicy",
			config: PolicyConfig{
				Condition:  DbePolicy,
				Action:     &action,
				Validation: &validation,
			},
			expected:    true,
			conditionID: DbePolicy,
		},
		{
			name: "PCIePolicy",
			config: PolicyConfig{
				Condition:  PCIePolicy,
				Action:     &action,
				Validation: &validation,
			},
			expected:    true,
			conditionID: PCIePolicy,
		},
		{
			name: "NvlinkPolicy",
			config: PolicyConfig{
				Condition:  NvlinkPolicy,
				Action:     &action,
				Validation: &validation,
			},
			expected:    true,
			conditionID: NvlinkPolicy,
		},
		{
			name: "XidPolicy",
			config: PolicyConfig{
				Condition:  XidPolicy,
				Action:     &action,
				Validation: &validation,
			},
			expected:    true,
			conditionID: XidPolicy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Setting %s policy...", tc.name)

			err := SetPolicyForGroup(group, tc.config)
			require.NoError(t, err)
			t.Logf("%s policy set successfully", tc.name)

			// Get the policy and verify it was set correctly
			t.Log("Retrieving policy configuration...")
			status, err := GetPolicyForGroup(group)
			require.NoError(t, err)
			require.NotNil(t, status)
			t.Logf("Policy retrieved - Mode: %d, Action: %v, Validation: %v, Conditions: %v",
				status.Mode, status.Action, status.Validation, status.Conditions)

			// Verify the policy is set
			assert.Contains(t, status.Conditions, tc.conditionID)
			assert.Equal(t, tc.expected, status.Conditions[tc.conditionID])
			assert.Equal(t, action, status.Action)
			assert.Equal(t, validation, status.Validation)

			t.Logf("%s policy assertions passed", tc.name)
		})
	}
}

func TestSetAndGetMultiplePolicies(t *testing.T) {
	t.Log("Initializing DCGM in Embedded mode...")
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()
	t.Log("DCGM initialized successfully")

	group := GroupAllGPUs()
	t.Logf("Created group handle for all GPUs: %+v", group)

	// Check how many GPUs we have
	gpuCount, err := GetAllDeviceCount()
	require.NoError(t, err)
	t.Logf("Found %d GPU(s) in the system", gpuCount)

	action := PolicyActionNone
	validation := PolicyValidationNone

	// Set multiple policies at once
	t.Log("Setting multiple policies simultaneously...")
	thermalThreshold := uint32(90)
	powerThreshold := uint32(350)
	maxRetiredPages := uint32(20)

	err = SetPolicyForGroup(
		group,
		PolicyConfig{
			Condition:      ThermalPolicy,
			Action:         &action,
			Validation:     &validation,
			MaxTemperature: &thermalThreshold,
		},
		PolicyConfig{
			Condition:  PowerPolicy,
			Action:     &action,
			Validation: &validation,
			MaxPower:   &powerThreshold,
		},
		PolicyConfig{
			Condition:       MaxRtPgPolicy,
			Action:          &action,
			Validation:      &validation,
			MaxRetiredPages: &maxRetiredPages,
		},
		PolicyConfig{
			Condition:  DbePolicy,
			Action:     &action,
			Validation: &validation,
		},
		PolicyConfig{
			Condition:  XidPolicy,
			Action:     &action,
			Validation: &validation,
		},
	)
	require.NoError(t, err)
	t.Log("Multiple policies set successfully")

	// Get the policy and verify all were set correctly
	t.Log("Retrieving policy configuration...")
	status, err := GetPolicyForGroup(group)
	require.NoError(t, err)
	require.NotNil(t, status)
	t.Logf("Policy retrieved - Mode: %d, Action: %v, Validation: %v, Conditions: %v",
		status.Mode, status.Action, status.Validation, status.Conditions)

	// Verify all policies are present
	t.Log("Verifying all policies were set correctly...")
	require.Len(t, status.Conditions, 5, "Expected 5 policy conditions to be set")

	// Verify each policy individually
	assert.Contains(t, status.Conditions, ThermalPolicy)
	assert.Equal(t, thermalThreshold, status.Conditions[ThermalPolicy])
	t.Logf("✓ ThermalPolicy: %d°C", thermalThreshold)

	assert.Contains(t, status.Conditions, PowerPolicy)
	assert.Equal(t, powerThreshold, status.Conditions[PowerPolicy])
	t.Logf("✓ PowerPolicy: %dW", powerThreshold)

	assert.Contains(t, status.Conditions, MaxRtPgPolicy)
	assert.Equal(t, maxRetiredPages, status.Conditions[MaxRtPgPolicy])
	t.Logf("✓ MaxRtPgPolicy: %d pages", maxRetiredPages)

	assert.Contains(t, status.Conditions, DbePolicy)
	assert.Equal(t, true, status.Conditions[DbePolicy])
	t.Log("✓ DbePolicy: enabled")

	assert.Contains(t, status.Conditions, XidPolicy)
	assert.Equal(t, true, status.Conditions[XidPolicy])
	t.Log("✓ XidPolicy: enabled")

	// Verify action and validation apply to all
	assert.Equal(t, action, status.Action)
	assert.Equal(t, validation, status.Validation)

	t.Log("All multiple policy assertions passed")
}

func TestSetPolicyAndWatchViolations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(100 * time.Millisecond)
	}()

	t.Log("Initializing DCGM in Embedded mode...")
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()
	t.Log("DCGM initialized successfully")

	t.Log("Creating fake GPU entities for testing...")
	fakeGPUs, err := withInjectionGPUs(t, 1)
	require.NoError(t, err)
	require.NotEmpty(t, fakeGPUs)
	gpu := fakeGPUs[0]
	t.Log("Fake GPU entities created")

	group := GroupAllGPUs()
	t.Logf("Created group handle for all GPUs: %+v", group)

	// Set up policies with thresholds using SetPolicyForGroup
	action := PolicyActionNone
	validation := PolicyValidationNone
	thermalThreshold := uint32(100)
	powerThreshold := uint32(250)

	t.Log("Setting thermal and power policies with SetPolicyForGroup...")
	err = SetPolicyForGroup(
		group,
		PolicyConfig{
			Condition:      ThermalPolicy,
			Action:         &action,
			Validation:     &validation,
			MaxTemperature: &thermalThreshold,
		},
		PolicyConfig{
			Condition:  PowerPolicy,
			Action:     &action,
			Validation: &validation,
			MaxPower:   &powerThreshold,
		},
	)
	require.NoError(t, err)
	t.Log("Policies set successfully with SetPolicyForGroup")

	// Watch for policy violations using WatchPolicyViolationsForGroup
	t.Log("Starting to watch for policy violations with WatchPolicyViolationsForGroup...")
	violations, err := WatchPolicyViolationsForGroup(ctx, group, ThermalPolicy, PowerPolicy)
	require.NoError(t, err)
	t.Log("Watching for violations")

	// Test 1: Inject thermal violation
	t.Run("ThermalViolation", func(t *testing.T) {
		t.Logf("Injecting thermal violation for GPU %d (threshold: %d°C)", gpu, thermalThreshold)

		err := InjectFieldValue(
			gpu,
			DCGM_FI_DEV_GPU_TEMP,
			DCGM_FT_INT64,
			0,
			time.Now().Add(60*time.Second).UnixMicro(),
			int64(thermalThreshold+1), // Exceed threshold
		)
		require.NoError(t, err)

		// Wait for violation
		select {
		case violation := <-violations:
			t.Logf("Received violation: %+v", violation)
			assert.Equal(t, ThermalPolicy, violation.Condition)
			require.IsType(t, ThermalPolicyCondition{}, violation.Data)
			thermalData := violation.Data.(ThermalPolicyCondition)
			assert.Equal(t, uint(thermalThreshold+1), thermalData.ThermalViolation)
			t.Logf("✓ Thermal violation detected: %d°C", thermalData.ThermalViolation)
		case <-time.After(20 * time.Second):
			t.Fatal("Timeout waiting for thermal violation")
		}
	})

	// Test 2: Inject power violation
	t.Run("PowerViolation", func(t *testing.T) {
		t.Logf("Injecting power violation for GPU %d (threshold: %dW)", gpu, powerThreshold)

		err := InjectFieldValue(
			gpu,
			DCGM_FI_DEV_POWER_USAGE,
			DCGM_FT_DOUBLE,
			0,
			time.Now().Add(60*time.Second).UnixMicro(),
			float64(powerThreshold+50), // Exceed threshold
		)
		require.NoError(t, err)

		// Wait for violation
		select {
		case violation := <-violations:
			t.Logf("Received violation: %+v", violation)
			assert.Equal(t, PowerPolicy, violation.Condition)
			require.IsType(t, PowerPolicyCondition{}, violation.Data)
			powerData := violation.Data.(PowerPolicyCondition)
			assert.Equal(t, uint(powerThreshold+50), powerData.PowerViolation)
			t.Logf("✓ Power violation detected: %dW", powerData.PowerViolation)
		case <-time.After(20 * time.Second):
			t.Fatal("Timeout waiting for power violation")
		}
	})

	t.Log("All SetPolicyForGroup + WatchPolicyViolationsForGroup tests passed")
}

func TestClearPolicyForGroup(t *testing.T) {
	t.Log("Initializing DCGM in Embedded mode...")
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()
	t.Log("DCGM initialized successfully")

	group := GroupAllGPUs()
	t.Logf("Created group handle for all GPUs: %+v", group)

	// Check how many GPUs we have
	gpuCount, err := GetAllDeviceCount()
	require.NoError(t, err)
	t.Logf("Found %d GPU(s) in the system", gpuCount)

	action := PolicyActionNone
	validation := PolicyValidationNone

	// Step 1: Set some policies
	t.Log("Step 1: Setting multiple policies...")
	thermalThreshold := uint32(90)
	powerThreshold := uint32(350)

	err = SetPolicyForGroup(
		group,
		PolicyConfig{
			Condition:      ThermalPolicy,
			Action:         &action,
			Validation:     &validation,
			MaxTemperature: &thermalThreshold,
		},
		PolicyConfig{
			Condition:  PowerPolicy,
			Action:     &action,
			Validation: &validation,
			MaxPower:   &powerThreshold,
		},
		PolicyConfig{
			Condition:  DbePolicy,
			Action:     &action,
			Validation: &validation,
		},
	)
	require.NoError(t, err)
	t.Log("Policies set successfully")

	// Step 2: Verify policies were set
	t.Log("Step 2: Verifying policies were set...")
	status, err := GetPolicyForGroup(group)
	require.NoError(t, err)
	require.NotNil(t, status)
	require.Len(t, status.Conditions, 3, "Expected 3 policies to be set")
	assert.Contains(t, status.Conditions, ThermalPolicy)
	assert.Contains(t, status.Conditions, PowerPolicy)
	assert.Contains(t, status.Conditions, DbePolicy)
	t.Logf("Verified 3 policies are active: %v", status.Conditions)

	// Step 3: Clear all policies
	t.Log("Step 3: Clearing all policies...")
	err = ClearPolicyForGroup(group)
	require.NoError(t, err)
	t.Log("Policies cleared successfully")

	// Step 4: Verify policies were cleared
	t.Log("Step 4: Verifying policies were cleared...")
	status, err = GetPolicyForGroup(group)
	require.NoError(t, err)
	require.NotNil(t, status)
	assert.Empty(t, status.Conditions, "Expected no policies after clear")
	t.Logf("Verified all policies cleared. Conditions map: %v", status.Conditions)

	t.Log("All clear policy tests passed")
}

// Helper function to create pointer to uint32
func ptrUint32(v uint32) *uint32 {
	return &v
}
