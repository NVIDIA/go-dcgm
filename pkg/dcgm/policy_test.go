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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicyErrors(t *testing.T) {
	type testCase struct {
		policy      policyCondition
		injectError func(gpu uint) error
		assert      func(cb PolicyViolation)
	}
	tests := []testCase{
		{
			policy: DbePolicy,
			injectError: func(gpu uint) error {
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_ECC_DBE_VOL_DEV", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_ECC_DBE_VOL_DEV,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
			},
			assert: func(cb PolicyViolation) {
				require.NotNil(t, cb)
				assert.Equal(t, DbePolicy, cb.Condition)
				require.IsType(t, dbePolicyCondition{}, cb.Data)
				policyCondition := cb.Data.(dbePolicyCondition)
				assert.Equal(t, uint(1), policyCondition.NumErrors)
				assert.Equal(t, "Device", policyCondition.Location)
			},
		},
		{
			policy: PowerPolicy,
			injectError: func(gpu uint) error {
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_POWER_USAGE", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_POWER_USAGE,
					DCGM_FT_DOUBLE,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					float64(300.0),
				)
			},
			assert: func(cb PolicyViolation) {
				require.NotNil(t, cb)
				assert.Equal(t, PowerPolicy, cb.Condition)
				require.IsType(t, powerPolicyCondition{}, cb.Data)
				policyCondition := cb.Data.(powerPolicyCondition)
				assert.Equal(t, uint(300), policyCondition.PowerViolation)
			},
		},
		{
			policy: PCIePolicy,
			injectError: func(gpu uint) error {
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_POWER_USAGE", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
			},
			assert: func(cb PolicyViolation) {
				require.NotNil(t, cb)
				assert.Equal(t, PCIePolicy, cb.Condition)
				require.IsType(t, pciPolicyCondition{}, cb.Data)
				pciPolicyCondition := cb.Data.(pciPolicyCondition)
				assert.Equal(t, uint(1), pciPolicyCondition.ReplayCounter)
			},
		},
		{
			policy: MaxRtPgPolicy,
			injectError: func(gpu uint) error {
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_RETIRED_DBE", gpu)
				err := InjectFieldValue(gpu,
					DCGM_FI_DEV_RETIRED_DBE,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(10),
				)
				if err == nil {
					//inject a SBE too so that the health check code gets past its internal checks
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
			assert: func(cb PolicyViolation) {
				require.NotNil(t, cb)
				assert.Equal(t, MaxRtPgPolicy, cb.Condition)
				require.IsType(t, retiredPagesPolicyCondition{}, cb.Data)
				retiredPagesPolicyCondition := cb.Data.(retiredPagesPolicyCondition)
				assert.Equal(t, uint(10), retiredPagesPolicyCondition.DbePages)
			},
		},
		{
			policy: ThermalPolicy,
			injectError: func(gpu uint) error {
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_GPU_TEMP", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_GPU_TEMP,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(101),
				)
			},
			assert: func(cb PolicyViolation) {
				require.NotNil(t, cb)
				assert.Equal(t, ThermalPolicy, cb.Condition)
				require.IsType(t, thermalPolicyCondition{}, cb.Data)
				thermalPolicyCondition := cb.Data.(thermalPolicyCondition)
				assert.Equal(t, uint(101), thermalPolicyCondition.ThermalViolation)
			},
		},
		{
			policy: NvlinkPolicy,
			injectError: func(gpu uint) error {
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
			},
			assert: func(cb PolicyViolation) {
				require.NotNil(t, cb)
				assert.Equal(t, NvlinkPolicy, cb.Condition)
				require.IsType(t, nvlinkPolicyCondition{}, cb.Data)
				nvlinkPolicyCondition := cb.Data.(nvlinkPolicyCondition)
				assert.Equal(t, uint(1), nvlinkPolicyCondition.Counter)
			},
		},
		{
			policy: XidPolicy,
			injectError: func(gpu uint) error {
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_XID_ERRORS", gpu)
				return InjectFieldValue(gpu,
					DCGM_FI_DEV_XID_ERRORS,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(16),
				)
			},
			assert: func(cb PolicyViolation) {
				require.NotNil(t, cb)
				assert.Equal(t, XidPolicy, cb.Condition)
				require.IsType(t, xidPolicyCondition{}, cb.Data)
				xidPolicyCondition := cb.Data.(xidPolicyCondition)
				assert.Equal(t, uint(16), xidPolicyCondition.ErrNum)
			},
		},
	}
	for _, tc := range tests {
		t.Run(string(tc.policy), func(t *testing.T) {
			cleanup, err := Init(Embedded)
			require.NoError(t, err)

			defer cleanup()

			numGPUs, err := GetAllDeviceCount()
			require.NoError(t, err)

			if numGPUs+1 > MAX_NUM_DEVICES {
				t.Skipf("Unable to add fake GPU with more than %d gpus", MAX_NUM_DEVICES)
			}

			entityList := []MigHierarchyInfo{
				{
					Entity: GroupEntityPair{EntityGroupId: FE_GPU},
				},
			}

			gpuIDs, err := CreateFakeEntities(entityList)
			require.NoError(t, err)

			gpu := gpuIDs[0]

			callback, err := Policy(gpu, tc.policy)
			require.NoError(t, err)

			err = tc.injectError(gpu)
			require.NoError(t, err)

			select {
			case callbackData := <-callback:
				require.NotNil(t, callbackData)
				tc.assert(callbackData)
				break
			case <-time.After(20 * time.Second):
				require.Fail(t, "policy callback never happened")
			}
		})
	}
}
