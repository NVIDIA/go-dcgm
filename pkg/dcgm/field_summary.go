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
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"fmt"
	"math/bits"
	"unsafe"
)

// FieldSummaryType identifies a statistic that DCGM can calculate over
// watched field samples. Values can be combined as a bitmask.
type FieldSummaryType uint32

const (
	// DCGM_SUMMARY_MIN requests the minimum value.
	DCGM_SUMMARY_MIN FieldSummaryType = 1 << iota
	// DCGM_SUMMARY_MAX requests the maximum value.
	DCGM_SUMMARY_MAX
	// DCGM_SUMMARY_AVG requests the arithmetic mean.
	DCGM_SUMMARY_AVG
	// DCGM_SUMMARY_SUM requests the sum of the values.
	DCGM_SUMMARY_SUM
	// DCGM_SUMMARY_COUNT requests the number of values.
	DCGM_SUMMARY_COUNT
	// DCGM_SUMMARY_INTEGRAL requests the time integral using microseconds as the time unit.
	DCGM_SUMMARY_INTEGRAL
	// DCGM_SUMMARY_DIFF requests the last value minus the first value.
	DCGM_SUMMARY_DIFF

	allFieldSummaryTypes = DCGM_SUMMARY_MIN |
		DCGM_SUMMARY_MAX |
		DCGM_SUMMARY_AVG |
		DCGM_SUMMARY_SUM |
		DCGM_SUMMARY_COUNT |
		DCGM_SUMMARY_INTEGRAL |
		DCGM_SUMMARY_DIFF
)

var orderedFieldSummaryTypes = [...]FieldSummaryType{
	DCGM_SUMMARY_MIN,
	DCGM_SUMMARY_MAX,
	DCGM_SUMMARY_AVG,
	DCGM_SUMMARY_SUM,
	DCGM_SUMMARY_COUNT,
	DCGM_SUMMARY_INTEGRAL,
	DCGM_SUMMARY_DIFF,
}

// FieldSummary contains summary statistics for one watched field. DCGM stores
// every statistic using the field's numeric type, so exactly one values map is
// populated according to FieldType.
type FieldSummary struct {
	FieldType     uint
	Int64Values   map[FieldSummaryType]int64
	Float64Values map[FieldSummaryType]float64
}

// GetFieldSummary returns summary statistics for a watched GPU field.
// startTime and endTime are Unix timestamps in microseconds. Zero means the
// corresponding end of the time range is unbounded.
func GetFieldSummary(
	gpuID uint,
	fieldID Short,
	summaryTypes FieldSummaryType,
	startTime uint64,
	endTime uint64,
) (FieldSummary, error) {
	return EntityGetFieldSummary(FE_GPU, gpuID, fieldID, summaryTypes, startTime, endTime)
}

// EntityGetFieldSummary returns summary statistics for a watched entity field.
// startTime and endTime are Unix timestamps in microseconds. Zero means the
// corresponding end of the time range is unbounded.
func EntityGetFieldSummary(
	entityGroup Field_Entity_Group,
	entityID uint,
	fieldID Short,
	summaryTypes FieldSummaryType,
	startTime uint64,
	endTime uint64,
) (FieldSummary, error) {
	if summaryTypes == 0 || summaryTypes&^allFieldSummaryTypes != 0 {
		return FieldSummary{}, fmt.Errorf("invalid field summary types %#x", summaryTypes)
	}

	request := C.dcgmFieldSummaryRequest_v1{
		version:         makeVersion1(unsafe.Sizeof(C.dcgmFieldSummaryRequest_v1{})),
		fieldId:         C.ushort(fieldID),
		entityGroupId:   C.dcgm_field_entity_group_t(entityGroup),
		entityId:        C.dcgm_field_eid_t(entityID),
		summaryTypeMask: C.uint32_t(summaryTypes),
		startTime:       C.uint64_t(startTime),
		endTime:         C.uint64_t(endTime),
	}

	result := C.dcgmGetFieldSummary(handle.handle, &request)
	if err := errorString(result); err != nil {
		return FieldSummary{}, fmt.Errorf("error getting field summary: %w", err)
	}

	requestedCount := bits.OnesCount32(uint32(summaryTypes))
	if int(request.response.summaryCount) != requestedCount {
		return FieldSummary{}, fmt.Errorf(
			"field summary returned %d values, expected %d",
			request.response.summaryCount,
			requestedCount,
		)
	}

	summary := FieldSummary{FieldType: uint(request.response.fieldType)}
	switch summary.FieldType {
	case DCGM_FT_INT64:
		summary.Int64Values = make(map[FieldSummaryType]int64, requestedCount)
	case DCGM_FT_DOUBLE:
		summary.Float64Values = make(map[FieldSummaryType]float64, requestedCount)
	default:
		return FieldSummary{}, fmt.Errorf("unsupported field summary type %#x", summary.FieldType)
	}

	responseIndex := 0
	for _, summaryType := range orderedFieldSummaryTypes {
		if summaryTypes&summaryType == 0 {
			continue
		}

		value := unsafe.Pointer(&request.response.values[responseIndex])
		if summary.FieldType == DCGM_FT_INT64 {
			summary.Int64Values[summaryType] = int64(*(*C.int64_t)(value))
		} else {
			summary.Float64Values[summaryType] = float64(*(*C.double)(value))
		}
		responseIndex++
	}

	return summary, nil
}
