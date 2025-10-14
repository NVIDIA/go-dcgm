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

// TestGetInfoMsg_NoMessages verifies getInfoMsg returns empty string when no info messages exist
func TestGetInfoMsg_NoMessages(t *testing.T) {
	response := createTestDiagResponse()

	result := getInfoMsg(0, 0, response)

	assert.Empty(t, result, "expected empty string when no info messages exist")
}

// TestGetInfoMsg_SingleMessage verifies getInfoMsg returns the single message without separator
func TestGetInfoMsg_SingleMessage(t *testing.T) {
	response := createTestDiagResponse()

	expectedMsg := "Allocated 83618558100 bytes (98.4%)"
	addInfoMessage(&response, 0, testMemoryIndex, expectedMsg)

	result := getInfoMsg(0, testMemoryIndex, response)

	assert.Equal(t, expectedMsg, result, "expected single message to be returned as-is")
}

// TestGetInfoMsg_MultipleMessages verifies all matching info messages are concatenated
func TestGetInfoMsg_MultipleMessages(t *testing.T) {
	response := createTestDiagResponse()

	entityID := uint(0)
	testID := uint(testPCIIndex)

	messages := []string{
		"GPU to Host bandwidth: 28.27 GB/s",
		"Host to GPU bandwidth: 27.65 GB/s",
		"bidirectional bandwidth: 50.59 GB/s",
		"GPU to Host latency: 1.305 us",
		"Host to GPU latency: 2.097 us",
		"bidirectional latency: 2.666 us",
	}

	for _, msg := range messages {
		addInfoMessage(&response, entityID, testID, msg)
	}

	result := getInfoMsg(entityID, testID, response)

	expected := "GPU to Host bandwidth: 28.27 GB/s | Host to GPU bandwidth: 27.65 GB/s | bidirectional bandwidth: 50.59 GB/s | GPU to Host latency: 1.305 us | Host to GPU latency: 2.097 us | bidirectional latency: 2.666 us"
	assert.Equal(t, expected, result, "expected all messages to be concatenated with ' | ' separator")
}

// TestGetInfoMsg_FiltersByEntityID verifies only messages matching entityId are returned
func TestGetInfoMsg_FiltersByEntityID(t *testing.T) {
	response := createTestDiagResponse()

	targetEntityID := uint(0)
	testID := uint(testMemoryIndex)

	// Add messages for different entities
	addInfoMessage(&response, targetEntityID, testID, "Message for entity 0")
	addInfoMessage(&response, 1, testID, "Message for entity 1")
	addInfoMessage(&response, targetEntityID, testID, "Another message for entity 0")

	result := getInfoMsg(targetEntityID, testID, response)

	expected := "Message for entity 0 | Another message for entity 0"
	assert.Equal(t, expected, result, "expected only messages matching entityId to be included")
	assert.NotContains(t, result, "entity 1", "should not contain messages from different entity")
}

// TestGetInfoMsg_FiltersByTestID verifies only messages matching testId are returned
func TestGetInfoMsg_FiltersByTestID(t *testing.T) {
	response := createTestDiagResponse()

	entityID := uint(0)
	targetTestID := uint(testMemoryIndex)

	// Add messages for different test IDs
	addInfoMessage(&response, entityID, targetTestID, "Memory test message 1")
	addInfoMessage(&response, entityID, testPCIIndex, "PCIe test message")
	addInfoMessage(&response, entityID, targetTestID, "Memory test message 2")

	result := getInfoMsg(entityID, targetTestID, response)

	expected := "Memory test message 1 | Memory test message 2"
	assert.Equal(t, expected, result, "expected only messages matching testId to be included")
	assert.NotContains(t, result, "PCIe", "should not contain messages from different test")
}

// TestGetInfoMsg_NoMatchingMessages verifies empty string when no messages match filters
func TestGetInfoMsg_NoMatchingMessages(t *testing.T) {
	response := createTestDiagResponse()

	// Add messages that don't match the query
	addInfoMessage(&response, 0, testMemoryIndex, "Some message")
	addInfoMessage(&response, 1, testPCIIndex, "Another message")

	// Query with different entityId and testId
	result := getInfoMsg(99, 99, response)

	assert.Empty(t, result, "expected empty string when no messages match the filters")
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

// TestGpuTestName verifies gpuTestName conversion
func TestGpuTestName(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"memory", testMemoryIndex, "memory"},
		{"diagnostic", testDiagnosticIndex, "diagnostic"},
		{"pcie", testPCIIndex, "pcie"},
		{"sm stress", testSMStressIndex, "sm stress"},
		{"targeted stress", testTargetedStressIndex, "targeted stress"},
		{"targeted power", testTargetedPowerIndex, "targeted power"},
		{"memory bandwidth", testMemoryBandwidthIndex, "memory bandwidth"},
		{"memtest", testMemtestIndex, "memtest"},
		{"pulse", testPulseTestIndex, "pulse"},
		{"eud", testEUDTestIndex, "eud"},
		{"software", testSoftwareIndex, "software"},
		{"context create", testContextCreateIndex, "context create"},
		{"unknown", 999, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gpuTestName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestNewDiagResult verifies DiagResult construction with multiple info messages
func TestNewDiagResult(t *testing.T) {
	response := createTestDiagResponse()

	entityID := uint(0)
	testID := uint(testPCIIndex)
	serialNumber := "1652923033635"

	// Setup result
	addDiagResult(&response, entityID, testID, testDiagResultPass)

	// Setup multiple info messages
	messages := []string{
		"GPU to Host bandwidth: 28.27 GB/s",
		"Host to GPU bandwidth: 27.65 GB/s",
		"bidirectional bandwidth: 50.59 GB/s",
	}
	for _, msg := range messages {
		addInfoMessage(&response, entityID, testID, msg)
	}

	// Setup entity with serial number
	addEntityWithSerial(&response, entityID, serialNumber)

	result := newDiagResult(0, response)

	require.NotNil(t, result)
	assert.Equal(t, "pass", result.Status)
	assert.Equal(t, "pcie", result.TestName)
	assert.Equal(t, "GPU to Host bandwidth: 28.27 GB/s | Host to GPU bandwidth: 27.65 GB/s | bidirectional bandwidth: 50.59 GB/s", result.TestOutput)
	assert.Equal(t, uint(0), result.ErrorCode)
	assert.Empty(t, result.ErrorMessage)
	assert.Equal(t, serialNumber, result.SerialNumber)
	assert.Equal(t, entityID, result.EntityID)
}
