/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
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

/*
#include "dcgm_structs.h"
*/
import "C"

import "unsafe"

func maxFieldPolicyEntitiesForTest() int {
	return int(C.DCGM_MAX_NUM_DEVICES)
}

func createTestFieldPolicyReadback() FieldPolicy {
	var cPolicy C.dcgmPolicyInfo_t
	cPolicy.version = C.dcgmPolicyInfo_version
	cPolicy.policyId = 41
	cPolicy.fieldId = C.ushort(DCGM_FI_DEV_GPU_TEMP)
	copyGoStringToCChars(cPolicy.name[:], "temperature")
	cPolicy.currentValue = 81.5
	cPolicy.threshold = 80
	cPolicy.policyOperator = C.DCGM_POLICY_OP_GT
	cPolicy.entities[0].entityGroupId = C.DCGM_FE_GPU
	cPolicy.entities[0].entityId = 3
	cPolicy.entityCount = 1
	cPolicy.channelMask = C.DCGM_POLICY_CHANNEL_CALLBACK
	cPolicy.state = C.DCGM_POLICY_STATE_ENABLED
	cPolicy.rateLimitSec = 7
	cPolicy.createdAt = 1_700_000_000_123_456
	cPolicy.violationCount = 9
	cPolicy.lastTriggeredTimestamp = 1_700_000_001_123_456
	return fieldPolicyFromC(cPolicy)
}

func fieldPolicyNameForTest(name string) string {
	var cPolicy C.dcgmPolicyInfo_t
	copyGoStringToCChars(cPolicy.name[:], name)
	return C.GoString(&cPolicy.name[0])
}

func fieldPolicyFunctionNotFoundErrorForTest() error {
	return dcgmFeatureResult(fieldPolicyFeature, "dcgmPolicyGetAll", C.DCGM_ST_FUNCTION_NOT_FOUND, "not found")
}

func fieldPolicyUninitializedErrorForTest() error {
	return &Error{Code: C.DCGM_ST_UNINITIALIZED}
}

func fieldPolicyFunctionNotFoundCodeForTest() int {
	return int(C.DCGM_ST_FUNCTION_NOT_FOUND)
}

func fieldPolicyLayoutForTest() (specVersion, expectedSpecVersion, infoVersion, expectedInfoVersion uint) {
	var spec C.dcgmPolicySpec_t
	var info C.dcgmPolicyInfo_t
	return uint(C.dcgmPolicySpec_version),
		uint(makeVersion1(unsafe.Sizeof(spec))),
		uint(C.dcgmPolicyInfo_version),
		uint(makeVersion1(unsafe.Sizeof(info)))
}

func fieldPolicyRateLimitForTest() uint {
	var policy C.dcgmPolicyInfo_t
	policy.rateLimitSec = C.uint(4_294_967_295)
	return uint(policy.rateLimitSec)
}
