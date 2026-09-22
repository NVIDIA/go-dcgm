package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import "strings"

// Package dcgm provides bindings for NVIDIA's Data Center GPU Manager (DCGM)

// DIAG_RESULT_STRING_SIZE represents the maximum size of diagnostic result strings.
const DIAG_RESULT_STRING_SIZE = 1024

// DiagType represents the type of diagnostic test to run
type DiagType int

const (
	// DiagQuick represents a quick diagnostic test that performs basic health checks
	DiagQuick DiagType = 1

	// DiagMedium represents a medium-length diagnostic test that performs more comprehensive checks
	DiagMedium DiagType = 2

	// DiagLong represents a long diagnostic test that performs extensive health checks
	DiagLong DiagType = 3

	// DiagExtended represents an extended diagnostic test that performs the most thorough system checks
	DiagExtended DiagType = 4
)

// DiagResult represents the legacy flattened result of one diagnostic test for one entity.
//
// Deprecated: Use DiagTest and DiagEntityResult.
type DiagResult struct {
	Status       string
	TestName     string
	TestOutput   string
	ErrorCode    uint
	ErrorMessage string
	SerialNumber string
	EntityID     uint
}

// DiagEntity contains metadata for an entity included in a diagnostic response.
type DiagEntity struct {
	Entity       GroupEntityPair
	SerialNumber string
	SKUDeviceID  string
}

// DiagEntityResult contains one diagnostic test result for one entity.
type DiagEntityResult struct {
	Entity GroupEntityPair
	Status string
}

// DiagError contains an error reported by a diagnostic test.
type DiagError struct {
	Entity   GroupEntityPair
	Code     uint
	Category ErrorCategory
	Severity ErrorSeverity
	Message  string
}

// DiagInfo contains an informational message reported by a diagnostic test.
type DiagInfo struct {
	Entity  GroupEntityPair
	Message string
}

// DiagTest contains the overall and per-entity results of one diagnostic test.
type DiagTest struct {
	Name           string
	PluginName     string
	Category       string
	Status         string
	Errors         []DiagError
	Info           []DiagInfo
	Results        []DiagEntityResult
	AuxDataVersion uint
	AuxData        string
}

// DiagResults contains the complete diagnostic response.
type DiagResults struct {
	// Software contains a flattened GPU-only compatibility view of Tests.
	//
	// Deprecated: Use Tests for diagnostic test results.
	Software      []DiagResult
	DCGMVersion   string
	DriverVersion string
	Categories    []string
	Entities      []DiagEntity
	SystemErrors  []DiagError
	Tests         []DiagTest
}

// diagResultString converts a diagnostic result code to its string representation
func diagResultString(r int) string {
	switch r {
	case C.DCGM_DIAG_RESULT_PASS:
		return "pass"
	case C.DCGM_DIAG_RESULT_SKIP:
		return "skipped"
	case C.DCGM_DIAG_RESULT_WARN:
		return "warn"
	case C.DCGM_DIAG_RESULT_FAIL:
		return "fail"
	case C.DCGM_DIAG_RESULT_NOT_RUN:
		return "notrun"
	}
	return ""
}

func groupEntityPair(entity C.dcgmGroupEntityPair_t) GroupEntityPair {
	return GroupEntityPair{
		EntityGroupId: Field_Entity_Group(entity.entityGroupId),
		EntityId:      uint(entity.entityId),
	}
}

func newDiagError(diagError *C.dcgmDiagError_v1) DiagError {
	return DiagError{
		Entity:   groupEntityPair(diagError.entity),
		Code:     uint(diagError.code),
		Category: ErrorCategory(diagError.category),
		Severity: ErrorSeverity(diagError.severity),
		Message:  C.GoString(&diagError.msg[0]),
	}
}

func legacyDiagResults(result DiagResults) []DiagResult {
	legacy := make([]DiagResult, 0)
	for i := range result.Tests {
		test := &result.Tests[i]
		for _, entityResult := range test.Results {
			if entityResult.Entity.EntityGroupId != FE_GPU {
				continue
			}
			diagResult := DiagResult{
				Status:   entityResult.Status,
				TestName: test.Name,
				EntityID: entityResult.Entity.EntityId,
			}

			for _, entity := range result.Entities {
				if entity.Entity == entityResult.Entity {
					diagResult.SerialNumber = entity.SerialNumber
					break
				}
			}

			var info []string
			for _, diagInfo := range test.Info {
				if diagInfo.Entity == entityResult.Entity || diagInfo.Entity.EntityGroupId == FE_NONE {
					info = append(info, diagInfo.Message)
				}
			}
			diagResult.TestOutput = strings.Join(info, " | ")

			for _, diagError := range test.Errors {
				if diagError.Entity != entityResult.Entity && diagError.Entity.EntityGroupId != FE_NONE {
					continue
				}
				diagResult.ErrorCode = diagError.Code
				diagResult.ErrorMessage = diagError.Message
				if diagError.Entity == entityResult.Entity {
					break
				}
			}

			legacy = append(legacy, diagResult)
		}
	}
	return legacy
}

func newDiagResults(response *C.dcgmDiagResponse_v12) DiagResults {
	result := DiagResults{
		DCGMVersion:   C.GoString(&response.dcgmVersion[0]),
		DriverVersion: C.GoString(&response.driverVersion[0]),
	}

	categoryCount := min(int(response.numCategories), len(response.categories))
	result.Categories = make([]string, categoryCount)
	for i := range categoryCount {
		result.Categories[i] = C.GoString(&response.categories[i][0])
	}

	entityCount := min(int(response.numEntities), len(response.entities))
	result.Entities = make([]DiagEntity, entityCount)
	for i := range entityCount {
		result.Entities[i] = DiagEntity{
			Entity:       groupEntityPair(response.entities[i].entity),
			SerialNumber: C.GoString(&response.entities[i].serialNum[0]),
			SKUDeviceID:  C.GoString(&response.entities[i].skuDeviceId[0]),
		}
	}

	for i := 0; i < min(int(response.numErrors), len(response.errors)); i++ {
		if response.errors[i].testId == C.DCGM_DIAG_RESPONSE_SYSTEM_ERROR {
			result.SystemErrors = append(result.SystemErrors, newDiagError(&response.errors[i]))
		}
	}

	testCount := min(int(response.numTests), len(response.tests))
	result.Tests = make([]DiagTest, testCount)
	for i := range testCount {
		test := &response.tests[i]
		diagTest := DiagTest{
			Name:           C.GoString(&test.name[0]),
			PluginName:     C.GoString(&test.pluginName[0]),
			Status:         diagResultString(int(test.result)),
			AuxDataVersion: uint(test.auxData.version >> 24),
			AuxData:        C.GoString(&test.auxData.data[0]),
		}
		if int(test.categoryIndex) < len(result.Categories) {
			diagTest.Category = result.Categories[test.categoryIndex]
		}

		for j := 0; j < min(int(test.numErrors), len(test.errorIndices)); j++ {
			index := int(test.errorIndices[j])
			if index >= min(int(response.numErrors), len(response.errors)) {
				continue
			}
			diagTest.Errors = append(diagTest.Errors, newDiagError(&response.errors[index]))
		}

		for j := 0; j < min(int(test.numInfo), len(test.infoIndices)); j++ {
			index := int(test.infoIndices[j])
			if index >= min(int(response.numInfo), len(response.info)) {
				continue
			}
			diagInfo := &response.info[index]
			diagTest.Info = append(diagTest.Info, DiagInfo{
				Entity:  groupEntityPair(diagInfo.entity),
				Message: C.GoString(&diagInfo.msg[0]),
			})
		}

		for j := 0; j < min(int(test.numResults), len(test.resultIndices)); j++ {
			index := uint(test.resultIndices[j])
			if index >= uint(min(int(response.numResults), len(response.results))) {
				continue
			}
			entityResult := response.results[index]
			diagTest.Results = append(diagTest.Results, DiagEntityResult{
				Entity: groupEntityPair(entityResult.entity),
				Status: diagResultString(int(entityResult.result)),
			})
		}
		result.Tests[i] = diagTest
	}

	result.Software = legacyDiagResults(result)
	return result
}

func diagLevel(diagType DiagType) C.dcgmDiagnosticLevel_t {
	switch diagType {
	case DiagQuick:
		return C.DCGM_DIAG_LVL_SHORT
	case DiagMedium:
		return C.DCGM_DIAG_LVL_MED
	case DiagLong:
		return C.DCGM_DIAG_LVL_LONG
	case DiagExtended:
		return C.DCGM_DIAG_LVL_XLONG
	}
	return C.DCGM_DIAG_LVL_INVALID
}

// RunDiag runs diagnostic tests on a group of GPUs with the specified diagnostic level.
// Parameters:
//   - diagType: The type/level of diagnostic test to run (Quick, Medium, Long, or Extended)
//   - groupID: The group of GPUs to run diagnostics on
//
// Returns:
//   - DiagResults containing the results of all diagnostic tests
//   - error if the diagnostics failed to run
func RunDiag(diagType DiagType, groupID GroupHandle) (DiagResults, error) {
	diagResults := new(C.dcgmDiagResponse_v12)
	diagResults.version = C.dcgmDiagResponse_version12

	result := C.dcgmRunDiagnostic(handle.handle, groupID.handle, diagLevel(diagType), diagResults)
	if err := errorString(result); err != nil {
		return DiagResults{}, &Error{msg: C.GoString(C.errorString(result)), Code: result}
	}

	return newDiagResults(diagResults), nil
}
