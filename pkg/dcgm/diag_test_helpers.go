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
	"unsafe"
)

// createTestDiagResponse creates a dcgmDiagResponse_v12 for testing
func createTestDiagResponse() C.dcgmDiagResponse_v12 {
	var response C.dcgmDiagResponse_v12
	response.version = C.dcgmDiagResponse_version12
	return response
}

func setDiagCString(dst *C.char, value string) {
	cStr := C.CString(value)
	defer C.free(unsafe.Pointer(cStr))
	C.strcpy(dst, cStr)
}

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
	response.numErrors = 2

	response.info[0].entity.entityGroupId = C.DCGM_FE_NONE
	response.info[0].testId = 0
	setDiagCString(&response.info[0].msg[0], "software checked")

	response.results[0].entity.entityGroupId = C.DCGM_FE_GPU
	response.results[0].entity.entityId = 7
	response.results[0].result = C.DCGM_DIAG_RESULT_WARN
	response.results[0].testId = 0

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
	response.numInfo = 2

	response.results[1].entity.entityGroupId = C.DCGM_FE_GPU
	response.results[1].entity.entityId = 7
	response.results[1].result = C.DCGM_DIAG_RESULT_PASS
	response.results[1].testId = 1

	response.numTests = 2
	response.numResults = 2
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
