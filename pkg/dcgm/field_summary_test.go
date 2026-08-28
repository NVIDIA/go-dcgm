//go:build linux && cgo

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityGetFieldSummaryRejectsInvalidSummaryTypes(t *testing.T) {
	tests := []struct {
		name         string
		summaryTypes FieldSummaryType
	}{
		{name: "empty", summaryTypes: 0},
		{name: "unknown", summaryTypes: 1 << 7},
		{name: "known and unknown", summaryTypes: DCGM_SUMMARY_MIN | 1<<7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EntityGetFieldSummary(FE_GPU, 0, DCGM_FI_DEV_XID_ERRORS, tt.summaryTypes, 0, 0)
			require.EqualError(t, err, fmt.Sprintf("invalid field summary types %#x", tt.summaryTypes))
		})
	}
}

func TestGetFieldSummary(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	gpus, err := withInjectionGPUs(t, 1)
	require.NoError(t, err)
	gpuID := gpus[0]

	groupID, err := CreateGroup(fmt.Sprintf("field-summary-%d", time.Now().UnixNano()))
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, DestroyGroup(groupID))
	}()
	require.NoError(t, AddToGroup(groupID, gpuID))

	fieldIDs := []Short{DCGM_FI_DEV_XID_ERRORS, DCGM_FI_DEV_BOARD_POWER_WATTS}
	fieldGroup, err := FieldGroupCreate(
		fmt.Sprintf("field-summary-fields-%d", time.Now().UnixNano()),
		fieldIDs,
	)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, FieldGroupDestroy(fieldGroup))
	}()

	require.NoError(t, WatchFieldsWithGroupEx(fieldGroup, groupID, 1_000_000, 60, 10))
	defer func() {
		assert.NoError(t, UnwatchFields(fieldGroup, groupID))
	}()

	startTime := time.Now().Add(-3 * time.Second).UnixMicro()
	for i, value := range []int64{10, 20, 40} {
		require.NoError(t, InjectFieldValue(
			gpuID,
			DCGM_FI_DEV_XID_ERRORS,
			DCGM_FT_INT64,
			0,
			startTime+int64(i)*int64(time.Second/time.Microsecond),
			value,
		))
	}
	for i, value := range []float64{1.5, 2.5, 4} {
		require.NoError(t, InjectFieldValue(
			gpuID,
			DCGM_FI_DEV_BOARD_POWER_WATTS,
			DCGM_FT_DOUBLE,
			0,
			startTime+int64(i)*int64(time.Second/time.Microsecond),
			value,
		))
	}

	intSummary, err := GetFieldSummary(
		gpuID,
		DCGM_FI_DEV_XID_ERRORS,
		allFieldSummaryTypes,
		uint64(startTime),
		uint64(startTime+2*int64(time.Second/time.Microsecond)),
	)
	require.NoError(t, err)
	require.Equal(t, DCGM_FT_INT64, intSummary.FieldType)
	require.Nil(t, intSummary.Float64Values)
	wantIntSummaryValues := map[FieldSummaryType]int64{
		DCGM_SUMMARY_MIN:      10,
		DCGM_SUMMARY_MAX:      40,
		DCGM_SUMMARY_AVG:      23,
		DCGM_SUMMARY_SUM:      70,
		DCGM_SUMMARY_COUNT:    3,
		DCGM_SUMMARY_INTEGRAL: 45_000_000,
		DCGM_SUMMARY_DIFF:     30,
	}
	require.Equal(t, wantIntSummaryValues, intSummary.Int64Values)

	queries := []struct {
		name string
		get  func(uint64, uint64) (FieldSummary, error)
	}{
		{
			name: "GetFieldSummary",
			get: func(startTime, endTime uint64) (FieldSummary, error) {
				return GetFieldSummary(
					gpuID,
					DCGM_FI_DEV_XID_ERRORS,
					allFieldSummaryTypes,
					startTime,
					endTime,
				)
			},
		},
		{
			name: "EntityGetFieldSummary",
			get: func(startTime, endTime uint64) (FieldSummary, error) {
				return EntityGetFieldSummary(
					FE_GPU,
					gpuID,
					DCGM_FI_DEV_XID_ERRORS,
					allFieldSummaryTypes,
					startTime,
					endTime,
				)
			},
		},
	}
	windows := []struct {
		name      string
		startTime uint64
		endTime   uint64
	}{
		{
			name:    "unbounded start",
			endTime: uint64(startTime + 2*int64(time.Second/time.Microsecond)),
		},
		{
			name:      "unbounded end",
			startTime: uint64(startTime),
		},
		{name: "fully unbounded"},
	}
	for _, query := range queries {
		for _, window := range windows {
			t.Run(query.name+"/"+window.name, func(t *testing.T) {
				summary, summaryErr := query.get(window.startTime, window.endTime)
				require.NoError(t, summaryErr)
				require.Equal(t, DCGM_FT_INT64, summary.FieldType)
				require.Nil(t, summary.Float64Values)
				require.Equal(t, wantIntSummaryValues, summary.Int64Values)
			})
		}
	}

	sparseSummary, err := GetFieldSummary(
		gpuID,
		DCGM_FI_DEV_XID_ERRORS,
		DCGM_SUMMARY_MIN|DCGM_SUMMARY_AVG|DCGM_SUMMARY_DIFF,
		uint64(startTime+int64(time.Second/time.Microsecond)),
		uint64(startTime+2*int64(time.Second/time.Microsecond)),
	)
	require.NoError(t, err)
	require.Equal(t, map[FieldSummaryType]int64{
		DCGM_SUMMARY_MIN:  20,
		DCGM_SUMMARY_AVG:  30,
		DCGM_SUMMARY_DIFF: 20,
	}, sparseSummary.Int64Values)

	floatSummary, err := EntityGetFieldSummary(
		FE_GPU,
		gpuID,
		DCGM_FI_DEV_BOARD_POWER_WATTS,
		allFieldSummaryTypes,
		uint64(startTime),
		uint64(startTime+2*int64(time.Second/time.Microsecond)),
	)
	require.NoError(t, err)
	require.Equal(t, DCGM_FT_DOUBLE, floatSummary.FieldType)
	require.Nil(t, floatSummary.Int64Values)
	require.Equal(t, 1.5, floatSummary.Float64Values[DCGM_SUMMARY_MIN])
	require.Equal(t, 4.0, floatSummary.Float64Values[DCGM_SUMMARY_MAX])
	require.InDelta(t, 8.0/3.0, floatSummary.Float64Values[DCGM_SUMMARY_AVG], 1e-12)
	require.Equal(t, 8.0, floatSummary.Float64Values[DCGM_SUMMARY_SUM])
	require.Equal(t, 3.0, floatSummary.Float64Values[DCGM_SUMMARY_COUNT])
	require.Equal(t, 5_250_000.0, floatSummary.Float64Values[DCGM_SUMMARY_INTEGRAL])
	require.Equal(t, 2.5, floatSummary.Float64Values[DCGM_SUMMARY_DIFF])
}
