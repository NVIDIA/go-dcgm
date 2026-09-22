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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateMultiNodeDiagnosticRequest(t *testing.T) {
	valid := MultiNodeDiagnosticRequest{
		Hosts:          []MultiNodeDiagnosticHost{{Address: "host-a"}},
		TestName:       "mnubergemm",
		TestParameters: []string{"mnubergemm.time_to_run=60"},
	}
	require.NoError(t, validateMultiNodeDiagnosticRequest(valid))

	tests := []struct {
		name    string
		request MultiNodeDiagnosticRequest
	}{
		{name: "no hosts", request: MultiNodeDiagnosticRequest{TestName: "mnubergemm"}},
		{name: "too many hosts", request: MultiNodeDiagnosticRequest{
			Hosts:    make([]MultiNodeDiagnosticHost, MaxMultiNodeDiagnosticHosts+1),
			TestName: "mnubergemm",
		}},
		{name: "too many parameters", request: MultiNodeDiagnosticRequest{
			Hosts:          []MultiNodeDiagnosticHost{{Address: "host-a"}},
			TestName:       "mnubergemm",
			TestParameters: make([]string, maxMultiNodeDiagnosticTestParameters+1),
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Error(t, validateMultiNodeDiagnosticRequest(test.request))
		})
	}
}

func TestValidateMultiNodeDiagnosticHostLimit(t *testing.T) {
	hosts := make([]MultiNodeDiagnosticHost, MaxMultiNodeDiagnosticHosts)
	for i := range hosts {
		hosts[i].Address = "host"
	}
	assert.NoError(t, validateMultiNodeDiagnosticRequest(MultiNodeDiagnosticRequest{
		Hosts:    hosts,
		TestName: "mnubergemm",
	}))
}

func TestMultiNodeDiagnosticV2ResponseEntityLimit(t *testing.T) {
	assert.Equal(t, MaxMultiNodeDiagnosticEntities, multiNodeDiagnosticResponseEntityLimitForTest())
}

func TestRunMultiNodeDiagnosticCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := RunMultiNodeDiagnostic(ctx, MultiNodeDiagnosticRequest{})
	require.ErrorIs(t, err, context.Canceled)
}

func TestStopMultiNodeDiagnosticAndWait(t *testing.T) {
	done := make(chan multiNodeDiagnosticCallResult, 1)
	stopCalled := make(chan struct{})
	completed := make(chan error, 1)

	go func() {
		completed <- stopMultiNodeDiagnosticAndWait(done, func() error {
			close(stopCalled)
			return nil
		})
	}()

	<-stopCalled
	select {
	case err := <-completed:
		t.Fatalf("returned before the diagnostic call finished: %v", err)
	default:
	}

	done <- multiNodeDiagnosticCallResult{}
	require.NoError(t, <-completed)
}

func TestStopMultiNodeDiagnosticAndWaitWaitsAfterStopError(t *testing.T) {
	done := make(chan multiNodeDiagnosticCallResult, 1)
	stopCalled := make(chan struct{})
	completed := make(chan error, 1)

	go func() {
		completed <- stopMultiNodeDiagnosticAndWait(done, func() error {
			close(stopCalled)
			return assert.AnError
		})
	}()

	<-stopCalled
	select {
	case err := <-completed:
		t.Fatalf("returned before the diagnostic call finished: %v", err)
	default:
	}

	done <- multiNodeDiagnosticCallResult{}
	require.ErrorIs(t, <-completed, assert.AnError)
}

func TestNewCMultiNodeDiagnosticRequest(t *testing.T) {
	requestStrings, err := multiNodeDiagnosticRequestStringsForTest(MultiNodeDiagnosticRequest{
		Hosts: []MultiNodeDiagnosticHost{
			{Address: "host-a"},
			{Address: "host-b"},
		},
		TestName:       "mnubergemm",
		TestParameters: []string{"mnubergemm.time_to_run=60"},
	})
	require.NoError(t, err)
	assert.Equal(t, multiNodeDiagnosticVersion2ForTest(), requestStrings.version)
	assert.Equal(t, requestStrings.expectedVersion, requestStrings.version)
	assert.Equal(t, []string{"host-a", "host-b"}, requestStrings.hosts)
	assert.Equal(t, "mnubergemm", requestStrings.testName)
	assert.Equal(t, []string{"mnubergemm.time_to_run=60"}, requestStrings.parameters)
}

func TestMultiNodeDiagnosticVersionMismatchPreservesDCGMCode(t *testing.T) {
	err := multiNodeDiagnosticVersionMismatchErrorForTest()
	var dcgmErr *Error
	require.ErrorAs(t, err, &dcgmErr)
	assert.Equal(t, multiNodeDiagnosticVersionMismatchCodeForTest(), int(dcgmErr.Code))
	assert.False(t, IsFeatureUnavailable(err))
}

func TestMultiNodeDiagnosticResultsMapping(t *testing.T) {
	result := createTestMultiNodeDiagnosticResults()
	require.Len(t, result.Hosts, 1)
	assert.Equal(t, "host-a", result.Hosts[0].Hostname)
	assert.Equal(t, []uint{0}, result.Hosts[0].EntityIndices)
	require.Len(t, result.Entities, 1)
	assert.Equal(t, GroupEntityPair{EntityGroupId: FE_GPU, EntityId: 5}, result.Entities[0].Entity)
	require.Len(t, result.Tests, 1)
	assert.Equal(t, "mnubergemm", result.Tests[0].Name)
	assert.Equal(t, "pass", result.Tests[0].Status)
	require.Len(t, result.Tests[0].Results, 1)
	assert.Equal(t, "pass", result.Tests[0].Results[0].Status)
}
