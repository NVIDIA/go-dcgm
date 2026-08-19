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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests deliberately omit setupTest: empty inputs must be rejected before any DCGM call.
func TestSelectedEmptyInputFieldGroupCreate(t *testing.T) {
	tests := map[string][]Short{
		"nil":   nil,
		"empty": {},
	}

	for name, fields := range tests {
		t.Run(name, func(t *testing.T) {
			fieldGroup, err := FieldGroupCreate("empty-field-group", fields)

			require.EqualError(t, err, "at least one field must be provided")
			assert.Zero(t, fieldGroup.GetHandle())
		})
	}
}

func TestSelectedEmptyInputGetLatestValuesForFields(t *testing.T) {
	tests := []struct {
		name   string
		fields []Short
	}{
		{name: "nil"},
		{name: "empty", fields: []Short{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := GetLatestValuesForFields(0, tt.fields)

			requireRejectedEmptyLatestValues(t, values, err)
		})
	}
}

func TestSelectedEmptyInputEntityGetLatestValues(t *testing.T) {
	tests := []struct {
		name   string
		fields []Short
	}{
		{name: "nil"},
		{name: "empty", fields: []Short{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := EntityGetLatestValues(FE_GPU, 0, tt.fields)

			requireRejectedEmptyLatestValues(t, values, err)
		})
	}
}

func TestSelectedEmptyInputEntitiesGetLatestValues(t *testing.T) {
	entity := GroupEntityPair{EntityGroupId: FE_GPU, EntityId: 0}
	field := DCGM_FI_DEV_NAME
	tests := []struct {
		name     string
		entities []GroupEntityPair
		fields   []Short
	}{
		{name: "nil entities", fields: []Short{field}},
		{name: "empty entities", entities: []GroupEntityPair{}, fields: []Short{field}},
		{name: "nil fields", entities: []GroupEntityPair{entity}},
		{name: "empty fields", entities: []GroupEntityPair{entity}, fields: []Short{}},
		{name: "nil entities and fields"},
		{name: "empty entities and fields", entities: []GroupEntityPair{}, fields: []Short{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := EntitiesGetLatestValues(tt.entities, tt.fields, DCGM_FV_FLAG_LIVE_DATA)

			requireRejectedEmptyLatestValues(t, values, err)
		})
	}
}

func requireRejectedEmptyLatestValues[T any](t *testing.T, values []T, err error) {
	t.Helper()

	require.Nil(t, values)
	var dcgmErr *Error
	require.ErrorAs(t, err, &dcgmErr)
	assert.Equal(t, DCGM_ST_BADPARAM, int(dcgmErr.Code))
}
