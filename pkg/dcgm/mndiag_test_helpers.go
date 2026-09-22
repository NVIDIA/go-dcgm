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

type multiNodeDiagnosticRequestStrings struct {
	version         uint
	expectedVersion uint
	hosts           []string
	testName        string
	parameters      []string
}

func multiNodeDiagnosticRequestStringsForTest(request MultiNodeDiagnosticRequest) (multiNodeDiagnosticRequestStrings, error) {
	cRequest, err := newCMultiNodeDiagnosticRequest(request)
	if err != nil {
		return multiNodeDiagnosticRequestStrings{}, err
	}
	hosts := make([]string, len(request.Hosts))
	for i := range hosts {
		hosts[i] = C.GoString(&cRequest.hostList[i][0])
	}
	parameters := make([]string, len(request.TestParameters))
	for i := range parameters {
		parameters[i] = C.GoString(&cRequest.testParms[i][0])
	}
	return multiNodeDiagnosticRequestStrings{
		version:         uint(cRequest.version),
		expectedVersion: uint(makeVersion2(unsafe.Sizeof(cRequest))),
		hosts:           hosts,
		testName:        C.GoString(&cRequest.testName[0]),
		parameters:      parameters,
	}, nil
}

func multiNodeDiagnosticVersion2ForTest() uint {
	return uint(C.dcgmRunMnDiag_version2)
}

func multiNodeDiagnosticResponseEntityLimitForTest() int {
	return int(C.DCGM_MN_DIAG_RESPONSE_HOSTS_MAX * C.DCGM_MN_DIAG_RESPONSE_ENTITIES_PER_HOST_MAX)
}

func multiNodeDiagnosticVersionMismatchErrorForTest() error {
	return dcgmFeatureResult(multiNodeDiagnosticFeature, "dcgmRunMnDiagnostic", C.DCGM_ST_VER_MISMATCH, "version mismatch")
}

func multiNodeDiagnosticVersionMismatchCodeForTest() int {
	return int(C.DCGM_ST_VER_MISMATCH)
}

func createTestMultiNodeDiagnosticResults() MultiNodeDiagnosticResults {
	response := new(C.dcgmMnDiagResponse_v2)
	response.version = C.dcgmMnDiagResponse_version2
	response.numHosts = 1
	copyGoStringToCChars(response.hosts[0].hostname[:], "host-a")
	copyGoStringToCChars(response.hosts[0].dcgmVersion[:], "4.7.0")
	copyGoStringToCChars(response.hosts[0].driverVersion[:], "580.1")
	response.hosts[0].numEntities = 1
	response.hosts[0].entityIndices[0] = 0
	response.numEntities = 1
	response.entities[0].entity.entityGroupId = C.DCGM_FE_GPU
	response.entities[0].entity.entityId = 5
	response.entities[0].hostId = 0
	copyGoStringToCChars(response.entities[0].serialNum[:], "GPU-SERIAL")
	response.numResults = 1
	response.results[0].entity = response.entities[0].entity
	response.results[0].hostId = 0
	response.results[0].testId = 0
	response.results[0].result = C.DCGM_DIAG_RESULT_PASS
	response.numTests = 1
	copyGoStringToCChars(response.tests[0].name[:], "mnubergemm")
	copyGoStringToCChars(response.tests[0].pluginName[:], "mnubergemm")
	response.tests[0].result = C.DCGM_DIAG_RESULT_PASS
	response.tests[0].numResults = 1
	response.tests[0].resultIndices[0] = 0
	return multiNodeDiagnosticResultsFromC(response)
}
