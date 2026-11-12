/*
 * Copyright (c) 2025, NVIDIA CORPORATION.  All rights reserved.
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
#include <stdlib.h>
#include <string.h>
#include "dcgm_agent.h"
#include "dcgm_errors.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"time"
	"unsafe"
)

// createTestDiagResponse creates a dcgmDiagResponse_v12 for testing
func createTestDiagResponse() C.dcgmDiagResponse_v12 {
	var response C.dcgmDiagResponse_v12
	response.version = C.dcgmDiagResponse_version12
	return response
}

// setDiagCString copies value into a fixed-size C string field of a test response.
func setDiagCString(dst *C.char, value string) {
	cStr := C.CString(value)
	defer C.free(unsafe.Pointer(cStr))
	C.strcpy(dst, cStr)
}

// createFullTestDiagResponse builds a diag response with versions, two
// categories, one GPU entity, and test results, for decoder tests.
func createFullTestDiagResponse() C.dcgmDiagResponse_v12 {
	response := createTestDiagResponse()
	setDiagCString(&response.dcgmVersion[0], "4.6.1")
	setDiagCString(&response.driverVersion[0], "580.1")

	setDiagCString(&response.categories[0][0], "Deployment")
	setDiagCString(&response.categories[1][0], "Hardware")
	response.numCategories = 2

	response.entities[0].entity.entityGroupId = C.DCGM_FE_GPU
	response.entities[0].entity.entityId = 7
	setDiagCString(&response.entities[0].serialNum[0], "GPU-SERIAL")
	setDiagCString(&response.entities[0].skuDeviceId[0], "27B8")
	response.numEntities = 1

	addSoftwareDiagnosticResult(&response)
	addMemoryDiagnosticResult(&response)
	response.numTests = 2
	response.numErrors = 2
	response.numInfo = 2
	response.numResults = 2
	return response
}

// addSoftwareDiagnosticResult fills the software warning and system error in a decoder fixture.
func addSoftwareDiagnosticResult(response *C.dcgmDiagResponse_v12) {
	software := &response.tests[0]
	setDiagCString(&software.name[0], "software")
	setDiagCString(&software.pluginName[0], "software")
	software.result = C.DCGM_DIAG_RESULT_WARN
	software.categoryIndex = 0
	software.numErrors = 1
	software.errorIndices[0] = 0
	software.numInfo = 1
	software.infoIndices[0] = 0
	software.numResults = 1
	software.resultIndices[0] = 0
	software.auxData.version = C.dcgmDiagTestAuxData_version1
	setDiagCString(&software.auxData.data[0], `{"version":"1"}`)

	response.errors[0].entity.entityGroupId = C.DCGM_FE_NONE
	response.errors[0].code = 42
	response.errors[0].category = C.DCGM_FR_EC_SOFTWARE_CONFIG
	response.errors[0].severity = C.DCGM_ERROR_CONFIG
	response.errors[0].testId = 0
	setDiagCString(&response.errors[0].msg[0], "software warning")

	response.errors[1].entity.entityGroupId = C.DCGM_FE_NONE
	response.errors[1].code = 99
	response.errors[1].category = C.DCGM_FR_EC_INTERNAL_OTHER
	response.errors[1].severity = C.DCGM_ERROR_ISOLATE
	response.errors[1].testId = C.DCGM_DIAG_RESPONSE_SYSTEM_ERROR
	setDiagCString(&response.errors[1].msg[0], "system failure")

	response.info[0].entity.entityGroupId = C.DCGM_FE_NONE
	response.info[0].testId = 0
	setDiagCString(&response.info[0].msg[0], "software checked")

	response.results[0].entity.entityGroupId = C.DCGM_FE_GPU
	response.results[0].entity.entityId = 7
	response.results[0].result = C.DCGM_DIAG_RESULT_WARN
	response.results[0].testId = 0
}

// addMemoryDiagnosticResult adds a successful GPU memory test beside the software warning.
func addMemoryDiagnosticResult(response *C.dcgmDiagResponse_v12) {
	memory := &response.tests[1]
	setDiagCString(&memory.name[0], "memory")
	setDiagCString(&memory.pluginName[0], "memory")
	memory.result = C.DCGM_DIAG_RESULT_PASS
	memory.categoryIndex = 1
	memory.numInfo = 1
	memory.infoIndices[0] = 1
	memory.numResults = 1
	memory.resultIndices[0] = 1

	response.info[1].entity.entityGroupId = C.DCGM_FE_GPU
	response.info[1].entity.entityId = 7
	response.info[1].testId = 1
	setDiagCString(&response.info[1].msg[0], "memory ok")

	response.results[1].entity.entityGroupId = C.DCGM_FE_GPU
	response.results[1].entity.entityId = 7
	response.results[1].result = C.DCGM_DIAG_RESULT_PASS
	response.results[1].testId = 1
}

// createDiagResponseWithRuntimeInfo adds a per-test runtime message to the memory result.
// DCGM reports elapsed time as "Test runtime: <s.mmm> seconds".
func createDiagResponseWithRuntimeInfo(message string, entityGroup Field_Entity_Group, entityID uint) C.dcgmDiagResponse_v12 {
	response := createFullTestDiagResponse()
	index := response.numInfo
	response.info[index].entity.entityGroupId = C.dcgm_field_entity_group_t(entityGroup)
	response.info[index].entity.entityId = C.uint(entityID)
	response.info[index].testId = 1
	setDiagCString(&response.info[index].msg[0], message)
	memory := &response.tests[1]
	memory.infoIndices[memory.numInfo] = index
	memory.numInfo++
	response.numInfo++
	return response
}

// Test constants exposed for testing
const (
	testDiagResultPass   = C.DCGM_DIAG_RESULT_PASS
	testDiagResultSkip   = C.DCGM_DIAG_RESULT_SKIP
	testDiagResultWarn   = C.DCGM_DIAG_RESULT_WARN
	testDiagResultFail   = C.DCGM_DIAG_RESULT_FAIL
	testDiagResultNotRun = C.DCGM_DIAG_RESULT_NOT_RUN
	testDiagLevelShort   = C.DCGM_DIAG_LVL_SHORT
	testDiagLevelMedium  = C.DCGM_DIAG_LVL_MED
	testDiagLevelLong    = C.DCGM_DIAG_LVL_LONG
	testDiagLevelXLong   = C.DCGM_DIAG_LVL_XLONG
	testDiagLevelInvalid = C.DCGM_DIAG_LVL_INVALID
)

// diagnosticRequestForTest converts a request through the C layout and returns its decoded fields for comparison.
func diagnosticRequestForTest(request DiagnosticRequest) (converted DiagnosticRequest, entitySelector string, nativeVersion, expectedVersion uint, conversionErr error) {
	native, err := cDiagnosticRequest(request)
	if err != nil {
		return DiagnosticRequest{}, "", 0, 0, err
	}
	decoded := DiagnosticRequest{
		Level:           DiagType(native.validate),
		Timeout:         time.Duration(native.timeoutSeconds) * time.Second,
		EnableHeartbeat: native.flags&C.DCGM_RUN_FLAGS_ENABLE_HEARTBEAT != 0,
	}
	group := GroupHandle{native.groupId}
	decoded.Group = &group

	for i := range native.testNames {
		if name := C.GoString(&native.testNames[i][0]); name != "" {
			decoded.TestNames = append(decoded.TestNames, name)
		}
	}

	for i := range native.testParms {
		if parameter := C.GoString(&native.testParms[i][0]); parameter != "" {
			decoded.TestParameters = append(decoded.TestParameters, parameter)
		}
	}

	return decoded, C.GoString(&native.entityIds[0]), uint(native.version), uint(C.dcgmRunDiag_version10), nil
}

// diagnosticBoundsForTest returns the native request limits used by boundary cases.
func diagnosticBoundsForTest() (names, nameBytes, parameters, parameterBytes, entities int) {
	return int(C.DCGM_MAX_TEST_NAMES), int(C.DCGM_MAX_TEST_NAMES_LEN), int(C.DCGM_MAX_TEST_PARMS), int(C.DCGM_MAX_TEST_PARMS_LEN_V2), int(C.DCGM_DIAG_RESPONSE_ENTITIES_MAX)
}

// diagnosticNativeErrorForTest builds a typed native failure without starting DCGM.
func diagnosticNativeErrorForTest(code int) error {
	return dcgmFeatureResult(diagnosticFeature, "dcgmActionValidate_v2", C.dcgmReturn_t(code), "native diagnostic error")
}

// diagnosticNullGroupForTest returns the native sentinel for an explicit entity selector.
func diagnosticNullGroupForTest() uintptr { return uintptr(C.DCGM_GROUP_NULL) }

// createConnectXDiagResponseForTest retargets the decoder fixture to a network device while preserving its result data.
func createConnectXDiagResponseForTest() DiagResults {
	response := createFullTestDiagResponse()
	response.entities[0].entity.entityGroupId = C.DCGM_FE_CONNECTX
	response.entities[0].entity.entityId = 17
	response.results[0].entity = response.entities[0].entity
	response.errors[0].entity = response.entities[0].entity
	setDiagCString(&response.tests[0].name[0], "nic_eud")
	setDiagCString(&response.tests[0].pluginName[0], "nic_eud")
	return newDiagResults(&response)
}
