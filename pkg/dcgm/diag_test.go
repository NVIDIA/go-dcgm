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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagResultStringSizeCompatibility(t *testing.T) {
	assert.Equal(t, 1024, DIAG_RESULT_STRING_SIZE)
}

// TestDiagResultString verifies diagResultString conversion
func TestDiagResultString(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"pass", testDiagResultPass, "pass"},
		{"skip", testDiagResultSkip, "skipped"},
		{"warn", testDiagResultWarn, "warn"},
		{"fail", testDiagResultFail, "fail"},
		{"not run", testDiagResultNotRun, "notrun"},
		{"unknown", 999, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := diagResultString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDiagLevel(t *testing.T) {
	tests := []struct {
		name  string
		input DiagType
		want  int
	}{
		{"quick", DiagQuick, int(testDiagLevelShort)},
		{"medium", DiagMedium, int(testDiagLevelMedium)},
		{"long", DiagLong, int(testDiagLevelLong)},
		{"extended", DiagExtended, int(testDiagLevelXLong)},
		{"invalid", DiagType(99), int(testDiagLevelInvalid)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := int(diagLevel(tt.input)); got != tt.want {
				t.Fatalf("diagLevel(%d) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewDiagResultsPreservesHierarchy(t *testing.T) {
	response := createFullTestDiagResponse()

	result := newDiagResults(&response)

	assert.Equal(t, "4.6.1", result.DCGMVersion)
	assert.Equal(t, "580.1", result.DriverVersion)
	assert.Equal(t, []string{"Deployment", "Hardware"}, result.Categories)
	require.Len(t, result.Entities, 1)
	assert.Equal(t, GroupEntityPair{EntityGroupId: FE_GPU, EntityId: 7}, result.Entities[0].Entity)
	assert.Equal(t, "GPU-SERIAL", result.Entities[0].SerialNumber)
	assert.Equal(t, "27B8", result.Entities[0].SKUDeviceID)
	require.Len(t, result.SystemErrors, 1)
	assert.Equal(t, uint(99), result.SystemErrors[0].Code)
	assert.Equal(t, "system failure", result.SystemErrors[0].Message)

	require.Len(t, result.Tests, 2)
	software := result.Tests[0]
	assert.Equal(t, "software", software.Name)
	assert.Equal(t, "software", software.PluginName)
	assert.Equal(t, "Deployment", software.Category)
	assert.Equal(t, "warn", software.Status)
	assert.Equal(t, uint(1), software.AuxDataVersion)
	assert.Equal(t, `{"version":"1"}`, software.AuxData)
	require.Len(t, software.Errors, 1)
	assert.Equal(t, GroupEntityPair{EntityGroupId: FE_NONE}, software.Errors[0].Entity)
	assert.Equal(t, uint(42), software.Errors[0].Code)
	assert.Equal(t, DCGM_FR_EC_SOFTWARE_CONFIG, software.Errors[0].Category)
	assert.Equal(t, DCGM_ERROR_CONFIG, software.Errors[0].Severity)
	assert.Equal(t, "software warning", software.Errors[0].Message)
	require.Len(t, software.Info, 1)
	assert.Equal(t, "software checked", software.Info[0].Message)

	memory := result.Tests[1]
	assert.Equal(t, "memory", memory.Name)
	assert.Equal(t, "Hardware", memory.Category)
	assert.Equal(t, "pass", memory.Status)
	require.Len(t, memory.Results, 1)
	assert.Equal(t, GroupEntityPair{EntityGroupId: FE_GPU, EntityId: 7}, memory.Results[0].Entity)
	assert.Equal(t, "pass", memory.Results[0].Status)
	require.Len(t, memory.Info, 1)
	assert.Equal(t, GroupEntityPair{EntityGroupId: FE_GPU, EntityId: 7}, memory.Info[0].Entity)
	assert.Equal(t, "memory ok", memory.Info[0].Message)
}

func TestNewDiagResultsIgnoresOutOfRangeCategoryIndex(t *testing.T) {
	response := createFullTestDiagResponse()
	response.tests[0].categoryIndex = response.numCategories

	result := newDiagResults(&response)

	require.NotEmpty(t, result.Tests)
	assert.Empty(t, result.Tests[0].Category)
}

func TestNewDiagResultsPreservesLegacyView(t *testing.T) {
	response := createFullTestDiagResponse()

	result := newDiagResults(&response)

	require.Len(t, result.Software, 2)
	assert.Equal(t, DiagResult{
		Status:       "warn",
		TestName:     "software",
		TestOutput:   "software checked",
		ErrorCode:    42,
		ErrorMessage: "software warning",
		SerialNumber: "GPU-SERIAL",
		EntityID:     7,
	}, result.Software[0])
	assert.Equal(t, DiagResult{
		Status:       "pass",
		TestName:     "memory",
		TestOutput:   "memory ok",
		SerialNumber: "GPU-SERIAL",
		EntityID:     7,
	}, result.Software[1])
}

func TestLegacyDiagResultsExcludesNonGPUEntities(t *testing.T) {
	result := DiagResults{Tests: []DiagTest{{
		Name: "cpu",
		Results: []DiagEntityResult{{
			Entity: GroupEntityPair{EntityGroupId: FE_CPU, EntityId: 7},
			Status: "fail",
		}},
	}}}

	assert.Empty(t, legacyDiagResults(result))
}
