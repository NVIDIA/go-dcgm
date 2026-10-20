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

// metricGroupInputForTest describes one metric group in plain Go values so unit
// tests do not need to construct C arrays directly.
type metricGroupInputForTest struct {
	major    uint
	minor    uint
	fieldIDs []uint
}

// metricGroupsFromValuesForTest builds a native metric-group response from Go
// test values, then passes it through the production response converter.
func metricGroupsFromValuesForTest(groups []metricGroupInputForTest, reportedGroupCount uint) ([]MetricGroup, error) {
	var response C.dcgmProfGetMetricGroups_t
	response.numMetricGroups = C.uint(reportedGroupCount)

	groupCount := min(len(groups), len(response.metricGroups))
	for i := 0; i < groupCount; i++ {
		response.metricGroups[i].majorId = C.ushort(groups[i].major)
		response.metricGroups[i].minorId = C.ushort(groups[i].minor)
		response.metricGroups[i].numFieldIds = C.uint(len(groups[i].fieldIDs))

		fieldCount := min(len(groups[i].fieldIDs), len(response.metricGroups[i].fieldIds))
		for j := 0; j < fieldCount; j++ {
			response.metricGroups[i].fieldIds[j] = C.ushort(groups[i].fieldIDs[j])
		}
	}

	return metricGroupsFromResponse(&response)
}

// metricGroupLimitsForTest reports the fixed native array sizes used to test
// valid boundary counts and rejected overflow counts.
func metricGroupLimitsForTest() (groupCount, fieldCount int) {
	var response C.dcgmProfGetMetricGroups_t
	return len(response.metricGroups), len(response.metricGroups[0].fieldIds)
}
