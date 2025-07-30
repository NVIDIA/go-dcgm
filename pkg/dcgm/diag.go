package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"unsafe"
)

// Package dcgm provides bindings for NVIDIA's Data Center GPU Manager (DCGM)

// DIAG_RESULT_STRING_SIZE represents the maximum size of diagnostic result strings
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

// DiagResult represents the result of a single diagnostic test
type DiagResult struct {
	// Status indicates the test result: "pass", "fail", "warn", "skip", or "notrun"
	Status string
	// TestName is the name of the diagnostic test that was run
	TestName string
	// TestOutput contains any additional output or messages from the test
	TestOutput string
	// ErrorCode is the numeric error code if the test failed
	ErrorCode uint
	// ErrorMessage contains a detailed error message if the test failed
	ErrorMessage string
}

// DiagResults contains the results of all diagnostic tests
type DiagResults struct {
	// Software contains the results of software-related diagnostic tests
	Software []DiagResult
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

func swTestName(t int) string {
	switch t {
	case C.DCGM_SWTEST_DENYLIST:
		return "presence of drivers on the denylist (e.g. nouveau)"
	case C.DCGM_SWTEST_NVML_LIBRARY:
		return "presence (and version) of NVML lib"
	case C.DCGM_SWTEST_CUDA_MAIN_LIBRARY:
		return "presence (and version) of CUDA lib"
	case C.DCGM_SWTEST_CUDA_RUNTIME_LIBRARY:
		return "presence (and version) of CUDA RT lib"
	case C.DCGM_SWTEST_PERMISSIONS:
		return "character device permissions"
	case C.DCGM_SWTEST_PERSISTENCE_MODE:
		return "persistence mode enabled"
	case C.DCGM_SWTEST_ENVIRONMENT:
		return "CUDA environment vars that may slow tests"
	case C.DCGM_SWTEST_PAGE_RETIREMENT:
		return "pending frame buffer page retirement"
	case C.DCGM_SWTEST_GRAPHICS_PROCESSES:
		return "graphics processes running"
	case C.DCGM_SWTEST_INFOROM:
		return "inforom corruption"
	}

	return ""
}

func gpuTestName(t int) string {
	switch t {
	case C.DCGM_MEMORY_INDEX:
		return "Memory"
	case C.DCGM_DIAGNOSTIC_INDEX:
		return "Diagnostic"
	case C.DCGM_PCI_INDEX:
		return "PCIe"
	case C.DCGM_SM_STRESS_INDEX:
		return "SM Stress"
	case C.DCGM_TARGETED_STRESS_INDEX:
		return "Targeted Stress"
	case C.DCGM_TARGETED_POWER_INDEX:
		return "Targeted Power"
	case C.DCGM_MEMORY_BANDWIDTH_INDEX:
		return "Memory bandwidth"
	case C.DCGM_MEMTEST_INDEX:
		return "Memtest"
	case C.DCGM_PULSE_TEST_INDEX:
		return "Pulse"
	case C.DCGM_EUD_TEST_INDEX:
		return "EUD"
	case C.DCGM_SOFTWARE_INDEX:
		return "Software"
	case C.DCGM_CONTEXT_CREATE_INDEX:
		return "Context create"
	}
	return ""
}

func getErrorMsg(entityId uint, response C.dcgmDiagResponse_v12) (msg string, code uint) {
	for i := 0; i < int(response.numErrors); i++ {
		if uint(response.errors[i].entity.entityId) != entityId {
			continue
		}

		msg = C.GoString((*C.char)(unsafe.Pointer(&response.errors[i].msg)))
		code = uint(response.errors[i].code)
		return
	}

	return
}

func getInfoMsg(entityId uint, response C.dcgmDiagResponse_v12) string {
	for i := 0; i < int(response.numInfo); i++ {
		if uint(response.info[i].entity.entityId) != entityId {
			continue
		}

		msg := C.GoString((*C.char)(unsafe.Pointer(&response.info[i].msg)))
		return msg
	}

	return ""
}

func newDiagResult(resultIndex uint, response C.dcgmDiagResponse_v12) DiagResult {
	entityId := uint(response.results[resultIndex].entity.entityId)

	msg, code := getErrorMsg(entityId, response)
	info := getInfoMsg(entityId, response)
	testName := swTestName(int(response.results[resultIndex].testId))

	return DiagResult{
		Status:       diagResultString(int(response.results[resultIndex].result)),
		TestName:     testName,
		TestOutput:   info,
		ErrorCode:    code,
		ErrorMessage: msg,
	}
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
//   - groupId: The group of GPUs to run diagnostics on
//
// Returns:
//   - DiagResults containing the results of all diagnostic tests
//   - error if the diagnostics failed to run
func RunDiag(diagType DiagType, groupID GroupHandle) (DiagResults, error) {
	var diagResults C.dcgmDiagResponse_v12
	diagResults.version = makeVersion12(unsafe.Sizeof(diagResults))

	result := C.dcgmRunDiagnostic(handle.handle, groupID.handle, diagLevel(diagType), (*C.dcgmDiagResponse_v12)(unsafe.Pointer(&diagResults)))
	if err := errorString(result); err != nil {
		return DiagResults{}, &Error{msg: C.GoString(C.errorString(result)), Code: result}
	}

	var diagRun DiagResults
	diagRun.Software = make([]DiagResult, diagResults.numResults)
	for i := 0; i < int(diagResults.numResults); i++ {
		diagRun.Software[i] = newDiagResult(uint(diagResults.results[i].entity.entityId), diagResults)
	}

	return diagRun, nil
}
