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

// healthResponseFromValuesForTest constructs a native response without writing
// beyond the incident array, even when its reported count is invalid.
func healthResponseFromValuesForTest(overall HealthResult, incidents []Incident, reportedCount uint32) (HealthResponse, error) {
	var response C.dcgmHealthResponse_v5
	response.overallHealth = C.dcgmHealthWatchResults_t(overall)
	response.incidentCount = C.uint(reportedCount)
	for i := 0; i < min(len(incidents), len(response.incidents)); i++ {
		input := incidents[i]
		native := &response.incidents[i]
		native.system = C.dcgmHealthSystems_t(input.System)
		native.health = C.dcgmHealthWatchResults_t(input.Health)
		native.error.code = C.uint(input.Error.Code)
		for j := 0; j < min(len(input.Error.Message), len(native.error.msg)-1); j++ {
			native.error.msg[j] = C.char(input.Error.Message[j])
		}
		native.entityInfo.entityGroupId = C.dcgm_field_entity_group_t(input.EntityInfo.EntityGroupId)
		native.entityInfo.entityId = C.uint(input.EntityInfo.EntityId)
	}
	return healthResponseFromNative(&response)
}

func healthIncidentCapacityForTest() uint32 {
	var response C.dcgmHealthResponse_v5
	return uint32(len(response.incidents))
}
