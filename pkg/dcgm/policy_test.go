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
	"errors"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
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

// Test helpers for policy fan-in invariants.
// Key assumptions:
//   - ViolationRegistration must stay bounded even when callback queues are full.
//   - pinPolicyState prevents leaked cleanup goroutines from niling callbacks mid-test.
//   - goleak.IgnoreCurrent() excludes DCGM background goroutines started by Init.

// resetPolicyChannels rebuilds the per-condition channels and fails fast if a prior test leaked a listener.
func resetPolicyChannels(t *testing.T) {
	t.Helper()
	policyCleanupMux.Lock()
	listeners := activeListeners
	policyCleanupMux.Unlock()
	require.Zero(t, listeners,
		"policy channels reset requires activeListeners==0; got %d (likely a leak from a prior test)",
		listeners)
	cleanupPolicyChannels()
	makePolicyChannels()
}

// fillPolicyChannel saturates the per-condition channel for `key` with zero-valued sentinels.
func fillPolicyChannel(t *testing.T, key string) {
	t.Helper()
	for i := 0; i < cap(callbacks[key]); i++ {
		select {
		case callbacks[key] <- PolicyViolation{}:
		default:
			t.Fatalf("fillPolicyChannel: %s already full at iteration %d (precondition broken)", key, i)
		}
	}
}

// TestViolationRegistration_BoundedTime asserts the cgo callback never blocks on a full per-condition channel.
func TestViolationRegistration_BoundedTime(t *testing.T) {
	conditions := []string{"dbe", "pcie", "maxrtpg", "thermal", "power", "nvlink", "xid"}
	for _, key := range conditions {
		t.Run(key, func(t *testing.T) {
			resetPolicyChannels(t)
			fillPolicyChannel(t, key)

			done := make(chan struct{})
			go func() {
				fireFakePolicyCallback(key)
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("ViolationRegistration blocked on full %s channel", key)
			}
		})
	}
}

// TestViolationRegistration_DropCounterAccuracy asserts droppedPolicyViolations increments exactly once per drop.
func TestViolationRegistration_DropCounterAccuracy(t *testing.T) {
	resetPolicyChannels(t)
	fillPolicyChannel(t, "xid")

	const N = 5
	before := droppedPolicyViolations.Load()
	for i := 0; i < N; i++ {
		fireFakePolicyCallback("xid")
	}
	got := droppedPolicyViolations.Load() - before
	require.Equal(t, uint64(N), got, "expected exactly %d drops, got %d", N, got)
}

// TestRegisterPolicy_SetPolicyFailureNoLeak asserts activeListeners net-zeroes when setPolicy fails before cgo registration.
func TestRegisterPolicy_SetPolicyFailureNoLeak(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup(t)

	policyCleanupMux.Lock()
	listenersBefore := activeListeners
	policyCleanupMux.Unlock()

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	var bogus GroupHandle
	bogus.SetHandle(^uintptr(0)) // out-of-range handle forces dcgmPolicySet -> DCGM_ST_BADPARAM

	ch, err := registerPolicy(ctx, bogus, XidPolicy)
	require.Error(t, err, "registerPolicy must fail with an out-of-range GroupHandle")
	require.Nil(t, ch)

	policyCleanupMux.Lock()
	listenersAfter := activeListeners
	policyCleanupMux.Unlock()
	require.Equal(t, listenersBefore, listenersAfter,
		"activeListeners must net-zero after failed register; before=%d after=%d",
		listenersBefore, listenersAfter)
}

// TestRegisterPolicy_LifecycleNoLeak asserts both fan-in and cleanup goroutines exit when the caller cancels.
func TestRegisterPolicy_LifecycleNoLeak(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup(t)

	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	ctx, cancel := context.WithCancel(context.Background())

	group := GroupAllGPUs()
	ch, err := registerPolicy(ctx, group, XidPolicy)
	require.NoError(t, err)
	require.NotNil(t, ch)

	cancel()

	for range ch {
	}

	waitForNoActiveListeners(t)
}

// TestViolationRegistration_RecoversAfterDrop asserts post-drain sends succeed after a non-blocking drop.
func TestViolationRegistration_RecoversAfterDrop(t *testing.T) {
	resetPolicyChannels(t)
	fillPolicyChannel(t, "xid")

	fireFakePolicyCallback("xid")

	for len(callbacks["xid"]) > 0 {
		<-callbacks["xid"]
	}

	fireFakePolicyCallback("xid")
	require.Equal(t, 1, len(callbacks["xid"]), "post-drain send must enqueue")
}

// policyRegisterer is the shared signature of registerPolicy and registerPolicyOnly.
type policyRegisterer func(context.Context, GroupHandle, ...PolicyCondition) (<-chan PolicyViolation, error)

// withPolicyRegisterStub swaps dcgmPolicyRegisterFn for `stub` for the test's lifetime. Not parallel-safe.
func withPolicyRegisterStub(t *testing.T, stub func(GroupHandle, uint64) error) {
	t.Helper()
	orig := dcgmPolicyRegisterFn
	dcgmPolicyRegisterFn = stub
	t.Cleanup(func() {
		dcgmPolicyRegisterFn = orig
	})
}

// waitForNoActiveListeners polls activeListeners until it returns to zero or the deadline elapses.
func waitForNoActiveListeners(t *testing.T) {
	t.Helper()
	require.Eventually(t, func() bool {
		policyCleanupMux.Lock()
		defer policyCleanupMux.Unlock()
		return activeListeners == 0
	}, 2*time.Second, 10*time.Millisecond, "expected activeListeners to return to zero")
}

// pinPolicyState parks activeListeners at a sentinel so leaked cleanup goroutines cannot nil the callbacks map. See header.
func pinPolicyState(t *testing.T) {
	t.Helper()
	const pin = 1 << 20

	policyCleanupMux.Lock()
	callbacks = nil
	policyChannelsInitialized = false
	activeListeners = pin
	policyCleanupMux.Unlock()

	t.Cleanup(func() {
		policyCleanupMux.Lock()
		callbacks = nil
		policyChannelsInitialized = false
		activeListeners = 0
		policyCleanupMux.Unlock()
	})
}

// waitForPolicyChannelDrain blocks until the per-condition channel for `key` is empty.
func waitForPolicyChannelDrain(t *testing.T, key string) {
	t.Helper()
	require.Eventually(t, func() bool {
		ch, ok := callbacks[key]
		return ok && len(ch) == 0
	}, time.Second, 10*time.Millisecond, "expected %s callback queue to drain", key)
}

// TestRunPolicyFanIn_CancellationUnblocksPendingSend asserts a fan-in parked on a blocking send exits within 2s of cancel.
func TestRunPolicyFanIn_CancellationUnblocksPendingSend(t *testing.T) {
	resetPolicyChannels(t)

	violation := make(chan PolicyViolation, 1)
	violation <- PolicyViolation{}

	callbacks["xid"] <- PolicyViolation{}
	callbacks["xid"] <- PolicyViolation{}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		runPolicyFanIn(ctx, violation)
	}()

	require.Eventually(t, func() bool {
		return len(callbacks["xid"]) == 1
	}, time.Second, 10*time.Millisecond,
		"fan-in did not park on the blocking send; precondition broken")

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("runPolicyFanIn did not exit within 2s of cancellation; " +
			"blocking send is not cancellation-aware")
	}
}

// assertRegisterFailureNoLeak drives the registration-failure cleanup path through the supplied registerer.
func assertRegisterFailureNoLeak(t *testing.T, register policyRegisterer) {
	t.Helper()
	cleanup := setupTest(t)
	defer cleanup(t)

	pinPolicyState(t)

	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	forcedErr := errors.New("forced dcgmPolicyRegister_v2 failure")
	withPolicyRegisterStub(t, func(_ GroupHandle, _ uint64) error {
		fireFakePolicyCallback("xid")
		waitForPolicyChannelDrain(t, "xid")
		fireFakePolicyCallback("xid")
		waitForPolicyChannelDrain(t, "xid")
		return forcedErr
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	var (
		ch  <-chan PolicyViolation
		err error
	)
	go func() {
		ch, err = register(ctx, GroupAllGPUs(), XidPolicy)
		close(done)
	}()

	require.Eventually(t, func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond,
		"register path blocked during registration failure cleanup")

	require.ErrorIs(t, err, forcedErr)
	require.Nil(t, ch)
}

// TestRegisterPolicy_RegisterFailureNoLeak covers the failure-cleanup path of registerPolicy.
func TestRegisterPolicy_RegisterFailureNoLeak(t *testing.T) {
	assertRegisterFailureNoLeak(t, registerPolicy)
}

// TestRegisterPolicyOnly_RegisterFailureNoLeak covers the failure-cleanup path of registerPolicyOnly.
func TestRegisterPolicyOnly_RegisterFailureNoLeak(t *testing.T) {
	assertRegisterFailureNoLeak(t, registerPolicyOnly)
}
