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

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFieldSummaryAPIsRejectInvalidSummaryTypes verifies validation before native calls.
func TestFieldSummaryAPIsRejectInvalidSummaryTypes(t *testing.T) {
	queries := []struct {
		name string
		get  func(FieldSummaryType) (FieldSummary, error)
	}{
		{
			name: "GPU",
			get: func(summaryTypes FieldSummaryType) (FieldSummary, error) {
				return GetFieldSummary(0, DCGM_FI_DEV_XID_ERRORS, summaryTypes, 0, 0)
			},
		},
		{
			name: "entity",
			get: func(summaryTypes FieldSummaryType) (FieldSummary, error) {
				return EntityGetFieldSummary(FE_GPU, 0, DCGM_FI_DEV_XID_ERRORS, summaryTypes, 0, 0)
			},
		},
	}
	tests := []struct {
		name         string
		summaryTypes FieldSummaryType
	}{
		{name: "empty", summaryTypes: 0},
		{name: "unknown", summaryTypes: 1 << 7},
		{name: "known and unknown", summaryTypes: DCGM_SUMMARY_MIN | 1<<7},
	}

	for _, query := range queries {
		for _, tt := range tests {
			t.Run(query.name+"/"+tt.name, func(t *testing.T) {
				summary, err := query.get(tt.summaryTypes)

				require.EqualError(t, err, fmt.Sprintf("invalid field summary types %#x", tt.summaryTypes))
				assert.Equal(t, FieldSummary{}, summary)
			})
		}
	}
}

// TestFieldSummaryFromResponseRejectsInvalidResponses verifies count and type validation.
func TestFieldSummaryFromResponseRejectsInvalidResponses(t *testing.T) {
	tests := []struct {
		name         string
		summaryTypes FieldSummaryType
		response     fieldSummaryResponse
		wantError    string
	}{
		{
			name:         "too few values",
			summaryTypes: DCGM_SUMMARY_MIN | DCGM_SUMMARY_MAX,
			response:     fieldSummaryResponse{fieldType: DCGM_FT_INT64, summaryCount: 1},
			wantError:    "field summary returned 1 values, expected 2",
		},
		{
			name:         "too many values",
			summaryTypes: DCGM_SUMMARY_MIN | DCGM_SUMMARY_MAX,
			response:     fieldSummaryResponse{fieldType: DCGM_FT_INT64, summaryCount: 3},
			wantError:    "field summary returned 3 values, expected 2",
		},
		{
			name:         "unsupported field type",
			summaryTypes: DCGM_SUMMARY_MIN,
			response:     fieldSummaryResponse{fieldType: DCGM_FT_STRING, summaryCount: 1},
			wantError:    fmt.Sprintf("unsupported field summary type %#x", DCGM_FT_STRING),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary, err := fieldSummaryFromResponse(tt.summaryTypes, tt.response)

			require.EqualError(t, err, tt.wantError)
			assert.Equal(t, FieldSummary{}, summary)
		})
	}
}

// TestFieldSummaryFromResponseConvertsInt64Values verifies every integer summary mapping.
func TestFieldSummaryFromResponseConvertsInt64Values(t *testing.T) {
	values := [...]uint64{^uint64(4), 0, 7, 42, 3, 1 << 40, ^uint64(10)}
	response := fieldSummaryResponse{
		fieldType:    DCGM_FT_INT64,
		summaryCount: uint(len(values)),
		values:       values,
	}

	summary, err := fieldSummaryFromResponse(allFieldSummaryTypes, response)

	require.NoError(t, err)
	assert.Equal(t, DCGM_FT_INT64, summary.FieldType)
	assert.Equal(t, map[FieldSummaryType]int64{
		DCGM_SUMMARY_MIN:      -5,
		DCGM_SUMMARY_MAX:      0,
		DCGM_SUMMARY_AVG:      7,
		DCGM_SUMMARY_SUM:      42,
		DCGM_SUMMARY_COUNT:    3,
		DCGM_SUMMARY_INTEGRAL: 1 << 40,
		DCGM_SUMMARY_DIFF:     -11,
	}, summary.Int64Values)
	assert.Nil(t, summary.Float64Values)
}

// TestFieldSummaryFromResponseConvertsFloat64Values verifies every floating-point summary mapping.
func TestFieldSummaryFromResponseConvertsFloat64Values(t *testing.T) {
	values := [...]float64{-1.5, 0, 2.25, 3.5, 4, 5.75, -6.125}
	response := fieldSummaryResponse{
		fieldType:    DCGM_FT_DOUBLE,
		summaryCount: uint(len(values)),
	}
	for i, value := range values {
		response.values[i] = math.Float64bits(value)
	}

	summary, err := fieldSummaryFromResponse(allFieldSummaryTypes, response)

	require.NoError(t, err)
	assert.Equal(t, DCGM_FT_DOUBLE, summary.FieldType)
	assert.Equal(t, map[FieldSummaryType]float64{
		DCGM_SUMMARY_MIN:      -1.5,
		DCGM_SUMMARY_MAX:      0,
		DCGM_SUMMARY_AVG:      2.25,
		DCGM_SUMMARY_SUM:      3.5,
		DCGM_SUMMARY_COUNT:    4,
		DCGM_SUMMARY_INTEGRAL: 5.75,
		DCGM_SUMMARY_DIFF:     -6.125,
	}, summary.Float64Values)
	assert.Nil(t, summary.Int64Values)
}

// TestFieldSummaryFromResponsePreservesPackedOrder verifies sparse native response ordering.
func TestFieldSummaryFromResponsePreservesPackedOrder(t *testing.T) {
	response := fieldSummaryResponse{
		fieldType:    DCGM_FT_INT64,
		summaryCount: 3,
		values:       [len(orderedFieldSummaryTypes)]uint64{10, 20, 30},
	}

	summary, err := fieldSummaryFromResponse(
		DCGM_SUMMARY_MIN|DCGM_SUMMARY_AVG|DCGM_SUMMARY_DIFF,
		response,
	)

	require.NoError(t, err)
	assert.Equal(t, map[FieldSummaryType]int64{
		DCGM_SUMMARY_MIN:  10,
		DCGM_SUMMARY_AVG:  20,
		DCGM_SUMMARY_DIFF: 30,
	}, summary.Int64Values)
	assert.NotContains(t, summary.Int64Values, DCGM_SUMMARY_MAX)
	assert.Nil(t, summary.Float64Values)
}
