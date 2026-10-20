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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMetricGroupsFromResponse verifies that empty and populated native metric
// groups are copied into independent Go values in their original order.
func TestMetricGroupsFromResponse(t *testing.T) {
	t.Run("empty response", func(t *testing.T) {
		groups, err := metricGroupsFromValuesForTest(nil, 0)

		require.NoError(t, err)
		require.NotNil(t, groups)
		assert.Empty(t, groups)
	})

	t.Run("preserves groups and field order", func(t *testing.T) {
		input := []metricGroupInputForTest{
			{major: 1, minor: 2},
			{major: 3, minor: 4, fieldIDs: []uint{1001, 1002, 1003}},
		}

		groups, err := metricGroupsFromValuesForTest(input, uint(len(input)))

		require.NoError(t, err)
		require.Equal(t, []MetricGroup{
			{Major: 1, Minor: 2, FieldIds: []uint{}},
			{Major: 3, Minor: 4, FieldIds: []uint{1001, 1002, 1003}},
		}, groups)

		input[1].fieldIDs[0] = 9999
		assert.Equal(t, uint(1001), groups[1].FieldIds[0])
	})
}

// TestMetricGroupsFromResponseAcceptsLimits verifies that responses using every
// available native group and field slot are accepted.
func TestMetricGroupsFromResponseAcceptsLimits(t *testing.T) {
	maxGroups, maxFields := metricGroupLimitsForTest()
	input := make([]metricGroupInputForTest, maxGroups)
	for i := range input {
		input[i] = metricGroupInputForTest{
			major:    uint(i + 1),
			minor:    uint(i + 11),
			fieldIDs: make([]uint, maxFields),
		}
		for j := range input[i].fieldIDs {
			input[i].fieldIDs[j] = uint(i*maxFields + j)
		}
	}

	groups, err := metricGroupsFromValuesForTest(input, uint(maxGroups))

	require.NoError(t, err)
	require.Len(t, groups, maxGroups)
	assert.Equal(t, uint(1), groups[0].Major)
	assert.Equal(t, uint(maxGroups+10), groups[maxGroups-1].Minor)
	require.Len(t, groups[maxGroups-1].FieldIds, maxFields)
	assert.Equal(t, uint(maxGroups*maxFields-1), groups[maxGroups-1].FieldIds[maxFields-1])
}

// TestMetricGroupsFromResponseRejectsInvalidCounts verifies that reported
// counts larger than the native arrays are rejected before they are indexed.
func TestMetricGroupsFromResponseRejectsInvalidCounts(t *testing.T) {
	maxGroups, maxFields := metricGroupLimitsForTest()

	t.Run("too many groups", func(t *testing.T) {
		groups, err := metricGroupsFromValuesForTest(nil, uint(maxGroups+1))

		require.Nil(t, groups)
		require.EqualError(t, err, fmt.Sprintf(
			"metric group response contains %d groups, maximum %d",
			maxGroups+1,
			maxGroups,
		))
	})

	t.Run("too many fields", func(t *testing.T) {
		input := []metricGroupInputForTest{{fieldIDs: make([]uint, maxFields+1)}}

		groups, err := metricGroupsFromValuesForTest(input, uint(len(input)))

		require.Nil(t, groups)
		require.EqualError(t, err, fmt.Sprintf(
			"metric group 0 contains %d field IDs, maximum %d",
			maxFields+1,
			maxFields,
		))
	})
}
