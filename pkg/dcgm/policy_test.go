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
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				gpu := uint(rand.Intn(8) + 1)
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
				gpu := uint(rand.Intn(8) + 1)
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
				gpu := uint(rand.Intn(8) + 1)
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
				gpu := uint(rand.Intn(8) + 1)
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
				gpu := uint(rand.Intn(8) + 1)
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
				gpu := uint(rand.Intn(8) + 1)
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
				gpu := uint(rand.Intn(8) + 1)
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
			//testcase: register multiple policy conditions
			policy:    []policyCondition{NvlinkPolicy, XidPolicy},
			numErrors: 2,
			injectError: func() error {
				gpu := uint(rand.Intn(8) + 1)
				//Inject a DBE error; since it has not registered DBEPolicy it will not get this event.
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_ECC_DBE_VOL_DEV", gpu)
				err := InjectFieldValue(gpu,
					DCGM_FI_DEV_ECC_DBE_VOL_DEV,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(1),
				)
				gpu = uint(rand.Intn(8) + 1)
				t.Logf("injecting %s for gpuId %d", "DCGM_FI_DEV_XID_ERRORS", gpu)
				err = InjectFieldValue(gpu,
					DCGM_FI_DEV_XID_ERRORS,
					DCGM_FT_INT64,
					0,
					time.Now().Add(60*time.Second).UnixMicro(),
					int64(16),
				)
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
			assert: func(cb PolicyViolation, en int) {
				switch en {
				case 1:
					require.NotNil(t, cb)
					assert.Equal(t, XidPolicy, cb.Condition)
					require.IsType(t, XidPolicyCondition{}, cb.Data)
					xidPolicyCondition := cb.Data.(XidPolicyCondition)
					assert.Equal(t, uint(16), xidPolicyCondition.ErrNum)
				case 2:
					require.NotNil(t, cb)
					assert.Equal(t, NvlinkPolicy, cb.Condition)
					require.IsType(t, NvlinkPolicyCondition{}, cb.Data)
					nvlinkPolicyCondition := cb.Data.(NvlinkPolicyCondition)
					assert.Equal(t, uint(1), nvlinkPolicyCondition.Counter)
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
