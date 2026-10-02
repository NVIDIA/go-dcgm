//go:build linux && cgo

/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dcgm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestGetPolicyForGroupWithOps checks group validation, policy decoding, and native read failures.
func TestGetPolicyForGroupWithOps(t *testing.T) {
	stageErr := errors.New("stage failed")
	dcgmCode := testDCGMReturn(testDCGMStatusVersionMismatch)
	dcgmErr := &Error{msg: "native policy failed", Code: dcgmCode}
	groupID := GroupHandle{}
	mixedGroup := &GroupInfo{EntityList: []GroupEntityPair{
		{EntityGroupId: FE_CPU, EntityId: 8},
		{EntityGroupId: FE_GPU, EntityId: 2},
		{EntityGroupId: FE_GPU, EntityId: 7},
	}}
	allConditions := policySnapshot{
		mode:       3,
		action:     PolicyActionGPUReset,
		validation: PolicyValidationMedium,
		condition:  0x7f,
		params: [xidPolicyIndex + 1]policyConditionParam{
			maxRtPgPolicyIndex: {typ: 1, value: 11},
			thermalPolicyIndex: {typ: 1, value: 82},
			powerPolicyIndex:   {typ: 1, value: 275},
		},
	}

	tests := []struct {
		name        string
		groupInfo   *GroupInfo
		groupErr    error
		gpuCount    int
		snapshot    policySnapshot
		readErr     error
		want        *PolicyStatus
		wantErrIs   error
		wantErrText string
	}{
		{
			name:        "wraps group info failure",
			groupErr:    stageErr,
			wantErrIs:   stageErr,
			wantErrText: "error getting group info",
		},
		{
			name:        "rejects group without GPUs",
			groupInfo:   &GroupInfo{EntityList: []GroupEntityPair{{EntityGroupId: FE_CPU, EntityId: 8}}},
			wantErrText: "cannot get policy for a group with no GPUs",
		},
		{
			name:      "preserves native DCGM error",
			groupInfo: mixedGroup,
			gpuCount:  2,
			readErr:   dcgmErr,
			wantErrIs: dcgmErr,
		},
		{
			name:      "decodes every condition and threshold",
			groupInfo: mixedGroup,
			gpuCount:  2,
			snapshot:  allConditions,
			want: &PolicyStatus{
				Mode:       3,
				Action:     PolicyActionGPUReset,
				Validation: PolicyValidationMedium,
				Conditions: map[PolicyCondition]interface{}{
					DbePolicy:     true,
					PCIePolicy:    true,
					MaxRtPgPolicy: uint32(11),
					ThermalPolicy: uint32(82),
					PowerPolicy:   uint32(275),
					NvlinkPolicy:  true,
					XidPolicy:     true,
				},
			},
		},
		{
			name:      "ignores non-LLONG thresholds",
			groupInfo: mixedGroup,
			gpuCount:  2,
			snapshot: policySnapshot{
				condition: 0x7f,
				params: [xidPolicyIndex + 1]policyConditionParam{
					maxRtPgPolicyIndex: {typ: 2, value: 11},
					thermalPolicyIndex: {typ: 2, value: 82},
					powerPolicyIndex:   {typ: 2, value: 275},
				},
			},
			want: &PolicyStatus{Conditions: map[PolicyCondition]interface{}{
				DbePolicy:    true,
				PCIePolicy:   true,
				NvlinkPolicy: true,
				XidPolicy:    true,
			}},
		},
		{
			name:      "returns empty conditions for zero mask",
			groupInfo: mixedGroup,
			gpuCount:  2,
			want:      &PolicyStatus{Conditions: map[PolicyCondition]interface{}{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewMockPolicyReadOps(gomock.NewController(t))
			api.EXPECT().getGroupInfo(groupID).Return(tt.groupInfo, tt.groupErr)
			if tt.groupErr == nil && tt.gpuCount > 0 {
				api.EXPECT().readPolicy(groupID, tt.gpuCount).Return(tt.snapshot, tt.readErr)
			}

			got, err := getPolicyForGroupWithOps(api, groupID)
			if tt.wantErrText != "" || tt.wantErrIs != nil {
				require.Error(t, err)
				require.Nil(t, got)
				if tt.wantErrText != "" {
					require.ErrorContains(t, err, tt.wantErrText)
				}
				if tt.wantErrIs != nil {
					require.ErrorIs(t, err, tt.wantErrIs)
				}
				if tt.readErr != nil {
					var gotDCGMErr *Error
					require.ErrorAs(t, err, &gotDCGMErr)
					require.Equal(t, dcgmCode, gotDCGMErr.Code)
				}
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
