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
	"crypto/rand"
	"encoding/binary"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// secureRandomUint returns a random uint in the range [1, max]
func secureRandomUint(maxValue uint) (uint, error) {
	var buf [8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return 0, err
	}

	// Convert to uint64 and reduce to our range
	n := binary.BigEndian.Uint64(buf[:])
	// Add 1 to shift range from [0, max-1] to [1, max]
	return uint(n%uint64(maxValue)) + 1, nil
}

func TestPolicyErrors(t *testing.T) {
	type testCase struct {
		policy      []policyCondition
		numErrors   int
		injectError func() error
		assert      func(cb PolicyViolation, en int)
	}

	tests := []testCase{
		{
			policy:    []policyCondition{DbePolicy},
			numErrors: 1,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_ECC_DBE_VOL_DEV", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_ECC_DBE_VOL_DEV,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)
				assert.Equal(t, DbePolicy, cb.Condition)
				require.IsType(t, DbePolicyCondition{}, cb.Data)
				policyCondition := cb.Data.(DbePolicyCondition)
				assert.Equal(t, uint(1), policyCondition.NumErrors)
				assert.Equal(t, "Device", policyCondition.Location)
			},
		},
		{
			policy:    []policyCondition{PowerPolicy},
			numErrors: 1,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_POWER_USAGE", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_POWER_USAGE,
					DCGM_FT_DOUBLE,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					float64(300.0),
				)
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)
				assert.Equal(t, PowerPolicy, cb.Condition)
				require.IsType(t, PowerPolicyCondition{}, cb.Data)
				policyCondition := cb.Data.(PowerPolicyCondition)
				assert.Equal(t, uint(300), policyCondition.PowerViolation)
			},
		},
		{
			policy:    []policyCondition{PCIePolicy},
			numErrors: 1,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_POWER_USAGE", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)
				assert.Equal(t, PCIePolicy, cb.Condition)
				require.IsType(t, PciPolicyCondition{}, cb.Data)
				pciPolicyCondition := cb.Data.(PciPolicyCondition)
				assert.Equal(t, uint(1), pciPolicyCondition.ReplayCounter)
			},
		},
		{
			policy:    []policyCondition{MaxRtPgPolicy},
			numErrors: 1,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_RETIRED_DBE", gpu)
				err := InjectFieldValue(gpu,
					DCGM_FI_DEV_RETIRED_DBE,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(10),
				)
				if err == nil {
					// inject a SBE too so that the health check code gets past its internal checks
					t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_RETIRED_SBE", gpu)
					err = InjectFieldValue(gpu,
						DCGM_FI_DEV_RETIRED_SBE,
						DCGM_FT_INT64,
						0,
						time.Now().Add(60*time.Second).UnixMicro(),
						int64(10),
					)
				}
				return err
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)
				assert.Equal(t, MaxRtPgPolicy, cb.Condition)
				require.IsType(t, RetiredPagesPolicyCondition{}, cb.Data)
				retiredPagesPolicyCondition := cb.Data.(RetiredPagesPolicyCondition)
				assert.Equal(t, uint(10), retiredPagesPolicyCondition.DbePages)
			},
		},
		{
			policy:    []policyCondition{ThermalPolicy},
			numErrors: 1,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_GPU_TEMP", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_GPU_TEMP,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(101),
				)
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)
				assert.Equal(t, ThermalPolicy, cb.Condition)
				require.IsType(t, ThermalPolicyCondition{}, cb.Data)
				thermalPolicyCondition := cb.Data.(ThermalPolicyCondition)
				assert.Equal(t, uint(101), thermalPolicyCondition.ThermalViolation)
			},
		},
		{
			policy:    []policyCondition{NvlinkPolicy},
			numErrors: 1,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)
				assert.Equal(t, NvlinkPolicy, cb.Condition)
				require.IsType(t, NvlinkPolicyCondition{}, cb.Data)
				nvlinkPolicyCondition := cb.Data.(NvlinkPolicyCondition)
				assert.Equal(t, uint(1), nvlinkPolicyCondition.Counter)
			},
		},
		{
			policy:    []policyCondition{XidPolicy},
			numErrors: 1,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_XID_ERRORS", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_XID_ERRORS,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(16),
				)
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)
				assert.Equal(t, XidPolicy, cb.Condition)
				require.IsType(t, XidPolicyCondition{}, cb.Data)
				xidPolicyCondition := cb.Data.(XidPolicyCondition)
				assert.Equal(t, uint(16), xidPolicyCondition.ErrNum)
			},
		},
		{
			// testcase: register multiple policy conditions
			policy:    []policyCondition{NvlinkPolicy, XidPolicy},
			numErrors: 2,
			injectError: func() error {
				gpu, _ := secureRandomUint(8)
				// Inject a DBE error; since it has not registered DBEPolicy it will not get this event.
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_ECC_DBE_VOL_DEV", gpu)
				err := InjectFieldValue(gpu,
					DCGM_FI_DEV_ECC_DBE_VOL_DEV,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
				if err != nil {
					return err
				}

				gpu, _ = secureRandomUint(8)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_XID_ERRORS", gpu)
				err = InjectFieldValue(gpu,
					DCGM_FI_DEV_XID_ERRORS,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(16),
				)
				if err != nil {
					return err
				}

				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL", gpu)
				err = InjectFieldValue(gpu,
					DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
				return err
			},
			assert: func(cb PolicyViolation, _ int) {
				require.NotNil(t, cb)

				switch cb.Condition {
				case XidPolicy:
					require.IsType(t, XidPolicyCondition{}, cb.Data)
					xidPolicyCondition := cb.Data.(XidPolicyCondition)
					assert.Equal(t, uint(16), xidPolicyCondition.ErrNum)
				case NvlinkPolicy:
					require.IsType(t, NvlinkPolicyCondition{}, cb.Data)
					nvlinkPolicyCondition := cb.Data.(NvlinkPolicyCondition)
					assert.Equal(t, uint(1), nvlinkPolicyCondition.Counter)
				default:
					require.FailNowf(
						t,
						"unexpected condition %s",
						string(cb.Condition),
					)
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(joinPolicy(tc.policy, "|"), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cleanup, err := Init(Embedded)
			require.NoError(t, err)

			defer func() {
				log.Printf("Cleaning up %s \n", t.Name())
				cleanup()
				cancel()
				time.Sleep(100 * time.Millisecond)
			}()

			numGPUs, err := GetAllDeviceCount()
			require.NoError(t, err)

			if numGPUs+1 > MAX_NUM_DEVICES {
				t.Skipf("Unable to add fake GPU with more than %d gpus", MAX_NUM_DEVICES)
			}

			entityList := []MigHierarchyInfo{
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
				{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
			}

			_, err = CreateFakeEntities(entityList)
			require.NoError(t, err)

			callback, err := ListenForPolicyViolations(ctx, tc.policy...)
			require.NoError(t, err)

			err = tc.injectError()
			require.NoError(t, err)

			numCb := 0
			select {
			case callbackData := <-callback:
				require.NotNil(t, callbackData)

				numCb++
				tc.assert(callbackData, numCb)

				if numCb == tc.numErrors {
					break
				}
			case <-time.After(20 * time.Second):
				require.Fail(t, "policy callback never happened")
			}
		})
	}
}

func joinPolicy(policy []policyCondition, sep string) string {
	var result strings.Builder

	for i, v := range policy {
		if i > 0 {
			result.WriteString(sep)
		}

		result.WriteString(string(v))
	}

	return result.String()
}

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

	err = SetPolicyForGroup(group,
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

	numGPUs, err := GetAllDeviceCount()
	require.NoError(t, err)
	t.Logf("Found %d GPU(s) in the system", numGPUs)

	if numGPUs+1 > MAX_NUM_DEVICES {
		t.Skipf("Unable to add fake GPU with more than %d gpus", MAX_NUM_DEVICES)
	}

	// Create fake GPUs for testing
	t.Log("Creating fake GPU entities for testing...")
	entityList := []MigHierarchyInfo{
		{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
		{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
		{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
		{Entity: GroupEntityPair{EntityGroupId: FE_GPU}},
	}
	_, err = CreateFakeEntities(entityList)
	require.NoError(t, err)
	t.Log("Fake GPU entities created")

	group := GroupAllGPUs()
	t.Logf("Created group handle for all GPUs: %+v", group)

	// Set up policies with thresholds using SetPolicyForGroup
	action := PolicyActionNone
	validation := PolicyValidationNone
	thermalThreshold := uint32(100)
	powerThreshold := uint32(250)

	t.Log("Setting thermal and power policies with SetPolicyForGroup...")
	err = SetPolicyForGroup(group,
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
		gpu, _ := secureRandomUint(4)
		t.Logf("Injecting thermal violation for GPU %d (threshold: %d°C)", gpu, thermalThreshold)

		err := InjectFieldValue(gpu,
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
		gpu, _ := secureRandomUint(4)
		t.Logf("Injecting power violation for GPU %d (threshold: %dW)", gpu, powerThreshold)

		err := InjectFieldValue(gpu,
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

	err = SetPolicyForGroup(group,
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
