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
	"context"
	"errors"
	"fmt"
	"sync"
)

const (
	// MaxMultiNodeDiagnosticHosts is the DCGM v2 multi-node diagnostic host limit.
	MaxMultiNodeDiagnosticHosts = int(C.DCGM_MNDIAG_MAX_NUM_HOSTS)

	// MaxMultiNodeDiagnosticEntities is the maximum number of entities in a DCGM v2 multi-node diagnostic response.
	MaxMultiNodeDiagnosticEntities = int(C.DCGM_MN_DIAG_RESPONSE_ENTITIES_MAX)

	maxMultiNodeDiagnosticTestParameters = int(C.DCGM_MAX_TEST_PARMS)
)

// MultiNodeDiagnosticHost selects one DCGM host.
type MultiNodeDiagnosticHost struct {
	Address string
}

// MultiNodeDiagnosticRequest describes a DCGM v2 multi-node diagnostic request.
type MultiNodeDiagnosticRequest struct {
	Hosts          []MultiNodeDiagnosticHost
	TestName       string
	TestParameters []string
}

// MultiNodeDiagnosticHostResult describes one host reported by a multi-node diagnostic.
type MultiNodeDiagnosticHostResult struct {
	Hostname      string
	DCGMVersion   string
	DriverVersion string
	EntityIndices []uint
}

// MultiNodeDiagnosticEntity describes one entity in a multi-node diagnostic response.
type MultiNodeDiagnosticEntity struct {
	Entity       GroupEntityPair
	HostID       uint
	SerialNumber string
	SKUDeviceID  string
}

// MultiNodeDiagnosticError describes one diagnostic error.
type MultiNodeDiagnosticError struct {
	Entity   GroupEntityPair
	HostID   uint
	TestID   uint
	Code     uint
	Category ErrorCategory
	Severity ErrorSeverity
	Message  string
}

// MultiNodeDiagnosticInfo describes one diagnostic information message.
type MultiNodeDiagnosticInfo struct {
	Entity  GroupEntityPair
	HostID  uint
	TestID  uint
	Message string
}

// MultiNodeDiagnosticEntityResult describes one test result for one entity.
type MultiNodeDiagnosticEntityResult struct {
	Entity GroupEntityPair
	Status string
	HostID uint
	TestID uint
}

// MultiNodeDiagnosticTest describes one completed multi-node diagnostic test.
type MultiNodeDiagnosticTest struct {
	Name       string
	PluginName string
	Status     string
	Errors     []MultiNodeDiagnosticError
	Info       []MultiNodeDiagnosticInfo
	Results    []MultiNodeDiagnosticEntityResult
}

// MultiNodeDiagnosticResults contains the complete DCGM v2 multi-node diagnostic response.
type MultiNodeDiagnosticResults struct {
	Hosts    []MultiNodeDiagnosticHostResult
	Entities []MultiNodeDiagnosticEntity
	Tests    []MultiNodeDiagnosticTest
	Errors   []MultiNodeDiagnosticError
	Info     []MultiNodeDiagnosticInfo
	Results  []MultiNodeDiagnosticEntityResult
}

const multiNodeDiagnosticFeature = "DCGM multi-node diagnostic API"

type multiNodeDiagnosticCallResult struct {
	result C.dcgmReturn_t
}

// DCGM allows one active multi-node diagnostic for this handle.
var multiNodeDiagnosticMu sync.Mutex

func requireMultiNodeDiagnosticSymbol(symbol string) error {
	if dcgmSymbolAvailable(symbol) {
		return nil
	}
	return &FeatureUnavailableError{Feature: multiNodeDiagnosticFeature, Symbol: symbol}
}

func multiNodeDiagnosticResult(symbol string, result C.dcgmReturn_t) error {
	return dcgmFeatureResult(multiNodeDiagnosticFeature, symbol, result, C.GoString(C.errorString(result)))
}

func validateMultiNodeDiagnosticRequest(request MultiNodeDiagnosticRequest) error {
	if len(request.Hosts) == 0 {
		return errors.New("multi-node diagnostic needs at least one host")
	}
	if len(request.Hosts) > MaxMultiNodeDiagnosticHosts {
		return fmt.Errorf("multi-node diagnostic has %d hosts; DCGM supports at most %d", len(request.Hosts), MaxMultiNodeDiagnosticHosts)
	}
	if err := validateCString("multi-node diagnostic test name", request.TestName, int(C.DCGM_MAX_STR_LENGTH)); err != nil {
		return err
	}
	if len(request.TestParameters) > maxMultiNodeDiagnosticTestParameters {
		return fmt.Errorf("multi-node diagnostic has %d test parameters; DCGM supports at most %d", len(request.TestParameters), maxMultiNodeDiagnosticTestParameters)
	}

	for i, host := range request.Hosts {
		if err := validateCString(fmt.Sprintf("multi-node diagnostic host %d", i), host.Address, int(C.DCGM_MAX_STR_LENGTH)); err != nil {
			return err
		}
	}

	for i, parameter := range request.TestParameters {
		if err := validateCString(fmt.Sprintf("multi-node diagnostic test parameter %d", i), parameter, int(C.DCGM_MAX_TEST_PARMS_LEN_V2)); err != nil {
			return err
		}
	}
	return nil
}

func newCMultiNodeDiagnosticRequest(request MultiNodeDiagnosticRequest) (C.dcgmRunMnDiag_v2, error) {
	if err := validateMultiNodeDiagnosticRequest(request); err != nil {
		return C.dcgmRunMnDiag_v2{}, err
	}

	var cRequest C.dcgmRunMnDiag_v2
	cRequest.version = C.dcgmRunMnDiag_version2
	for i, host := range request.Hosts {
		copyGoStringToCChars(cRequest.hostList[i][:], host.Address)
	}
	copyGoStringToCChars(cRequest.testName[:], request.TestName)
	for i, parameter := range request.TestParameters {
		copyGoStringToCChars(cRequest.testParms[i][:], parameter)
	}
	return cRequest, nil
}

// RunMultiNodeDiagnostic runs a DCGM v2 multi-node diagnostic.
// If ctx is cancelled while the diagnostic runs, it asks DCGM to stop the run,
// waits for that run to finish, and returns ctx.Err().
func RunMultiNodeDiagnostic(ctx context.Context, request MultiNodeDiagnosticRequest) (MultiNodeDiagnosticResults, error) {
	const symbol = "dcgmRunMnDiagnostic"
	if ctx == nil {
		return MultiNodeDiagnosticResults{}, errors.New("multi-node diagnostic context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return MultiNodeDiagnosticResults{}, err
	}
	if err := requireMultiNodeDiagnosticSymbol(symbol); err != nil {
		return MultiNodeDiagnosticResults{}, err
	}
	cRequest, err := newCMultiNodeDiagnosticRequest(request)
	if err != nil {
		return MultiNodeDiagnosticResults{}, err
	}
	multiNodeDiagnosticMu.Lock()
	defer multiNodeDiagnosticMu.Unlock()
	if err := ctx.Err(); err != nil {
		return MultiNodeDiagnosticResults{}, err
	}

	cResponse := new(C.dcgmMnDiagResponse_v2)
	cResponse.version = C.dcgmMnDiagResponse_version2
	done := make(chan multiNodeDiagnosticCallResult, 1)
	go func() {
		done <- multiNodeDiagnosticCallResult{result: C.dcgmRunMnDiagnostic(handle.handle, &cRequest, cResponse)}
	}()

	select {
	case call := <-done:
		if err := multiNodeDiagnosticResult(symbol, call.result); err != nil {
			return MultiNodeDiagnosticResults{}, err
		}
		return multiNodeDiagnosticResultsFromC(cResponse), nil
	case <-ctx.Done():
		if err := stopMultiNodeDiagnosticAndWait(done, StopMultiNodeDiagnostic); err != nil {
			return MultiNodeDiagnosticResults{}, errors.Join(ctx.Err(), err)
		}
		return MultiNodeDiagnosticResults{}, ctx.Err()
	}
}

func stopMultiNodeDiagnosticAndWait(done <-chan multiNodeDiagnosticCallResult, stop func() error) error {
	err := stop()
	<-done
	return err
}

// StopMultiNodeDiagnostic asks DCGM to stop the active multi-node diagnostic.
func StopMultiNodeDiagnostic() error {
	const symbol = "dcgmStopMnDiagnostic"
	if err := requireMultiNodeDiagnosticSymbol(symbol); err != nil {
		return err
	}
	return multiNodeDiagnosticResult(symbol, C.dcgmStopMnDiagnostic(handle.handle))
}

func multiNodeDiagnosticResultsFromC(response *C.dcgmMnDiagResponse_v2) MultiNodeDiagnosticResults {
	if response == nil {
		return MultiNodeDiagnosticResults{}
	}

	result := MultiNodeDiagnosticResults{}
	hostCount := min(int(response.numHosts), len(response.hosts))
	result.Hosts = make([]MultiNodeDiagnosticHostResult, hostCount)
	for i := range result.Hosts {
		host := &response.hosts[i]
		entityCount := min(int(host.numEntities), len(host.entityIndices))
		entityIndices := make([]uint, entityCount)
		for j := range entityIndices {
			entityIndices[j] = uint(host.entityIndices[j])
		}
		result.Hosts[i] = MultiNodeDiagnosticHostResult{
			Hostname:      C.GoString(&host.hostname[0]),
			DCGMVersion:   C.GoString(&host.dcgmVersion[0]),
			DriverVersion: C.GoString(&host.driverVersion[0]),
			EntityIndices: entityIndices,
		}
	}

	entityCount := min(int(response.numEntities), len(response.entities))
	result.Entities = make([]MultiNodeDiagnosticEntity, entityCount)
	for i := range result.Entities {
		entity := &response.entities[i]
		result.Entities[i] = MultiNodeDiagnosticEntity{
			Entity:       groupEntityPair(entity.entity),
			HostID:       uint(entity.hostId),
			SerialNumber: C.GoString(&entity.serialNum[0]),
			SKUDeviceID:  C.GoString(&entity.skuDeviceId[0]),
		}
	}

	errorCount := min(int(response.numErrors), len(response.errors))
	result.Errors = make([]MultiNodeDiagnosticError, errorCount)
	for i := range result.Errors {
		result.Errors[i] = multiNodeDiagnosticErrorFromC(response.errors[i])
	}
	infoCount := min(int(response.numInfos), len(response.info))
	result.Info = make([]MultiNodeDiagnosticInfo, infoCount)
	for i := range result.Info {
		result.Info[i] = multiNodeDiagnosticInfoFromC(response.info[i])
	}
	resultCount := min(int(response.numResults), len(response.results))
	result.Results = make([]MultiNodeDiagnosticEntityResult, resultCount)
	for i := range result.Results {
		result.Results[i] = multiNodeDiagnosticEntityResultFromC(response.results[i])
	}

	testCount := min(int(response.numTests), len(response.tests))
	result.Tests = make([]MultiNodeDiagnosticTest, testCount)
	for i := range result.Tests {
		test := &response.tests[i]
		result.Tests[i] = MultiNodeDiagnosticTest{
			Name:       C.GoString(&test.name[0]),
			PluginName: C.GoString(&test.pluginName[0]),
			Status:     diagResultString(int(test.result)),
			Errors:     multiNodeDiagnosticTestErrors(test, result.Errors),
			Info:       multiNodeDiagnosticTestInfo(test, result.Info),
			Results:    multiNodeDiagnosticTestResults(test, result.Results),
		}
	}
	return result
}

func multiNodeDiagnosticErrorFromC(cError C.dcgmMnDiagError_v1) MultiNodeDiagnosticError {
	return MultiNodeDiagnosticError{
		Entity:   groupEntityPair(cError.entity),
		HostID:   uint(cError.hostId),
		TestID:   uint(cError.testId),
		Code:     uint(cError.code),
		Category: ErrorCategory(cError.category),
		Severity: ErrorSeverity(cError.severity),
		Message:  C.GoString(&cError.msg[0]),
	}
}

func multiNodeDiagnosticInfoFromC(cInfo C.dcgmMnDiagInfo_v1) MultiNodeDiagnosticInfo {
	return MultiNodeDiagnosticInfo{
		Entity:  groupEntityPair(cInfo.entity),
		HostID:  uint(cInfo.hostId),
		TestID:  uint(cInfo.testId),
		Message: C.GoString(&cInfo.msg[0]),
	}
}

func multiNodeDiagnosticEntityResultFromC(cResult C.dcgmMnDiagEntityResult_v1) MultiNodeDiagnosticEntityResult {
	return MultiNodeDiagnosticEntityResult{
		Entity: groupEntityPair(cResult.entity),
		Status: diagResultString(int(cResult.result)),
		HostID: uint(cResult.hostId),
		TestID: uint(cResult.testId),
	}
}

func multiNodeDiagnosticTestErrors(test *C.dcgmMnDiagTestRun_v2, diagnosticErrors []MultiNodeDiagnosticError) []MultiNodeDiagnosticError {
	count := min(int(test.numErrors), len(test.errorIndices))
	result := make([]MultiNodeDiagnosticError, 0, count)
	for i := 0; i < count; i++ {
		index := int(test.errorIndices[i])
		if index < len(diagnosticErrors) {
			result = append(result, diagnosticErrors[index])
		}
	}
	return result
}

func multiNodeDiagnosticTestInfo(test *C.dcgmMnDiagTestRun_v2, info []MultiNodeDiagnosticInfo) []MultiNodeDiagnosticInfo {
	count := min(int(test.numInfo), len(test.infoIndices))
	result := make([]MultiNodeDiagnosticInfo, 0, count)
	for i := 0; i < count; i++ {
		index := int(test.infoIndices[i])
		if index < len(info) {
			result = append(result, info[index])
		}
	}
	return result
}

func multiNodeDiagnosticTestResults(test *C.dcgmMnDiagTestRun_v2, results []MultiNodeDiagnosticEntityResult) []MultiNodeDiagnosticEntityResult {
	count := min(int(test.numResults), len(test.resultIndices))
	result := make([]MultiNodeDiagnosticEntityResult, 0, count)
	for i := 0; i < count; i++ {
		index := int(test.resultIndices[i])
		if index < len(results) {
			result = append(result, results[index])
		}
	}
	return result
}
