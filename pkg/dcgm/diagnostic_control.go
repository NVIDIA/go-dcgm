package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// DiagnosticRequest selects a diagnostic level or named tests and their parameters.
// A nil Group with no Entities selects all GPUs. Group and Entities are mutually
// exclusive. TestParameters use testName.parameterName=value syntax; DCGM owns
// parameter validation. Timeout is zero for DCGM's default or a whole number of
// seconds. DCGM converts this value to a uint32 millisecond timeout, limiting it
// to 4,294,967 seconds to avoid overflow during that conversion.
//
// EnableHeartbeat requires DiagSendHeartbeat from another goroutine before
// the host engine's configured timeout expires (ten minutes by default).
// Host engines supporting DCGM_DIAG_HEARTBEAT_TIMEOUT_SEC use that value in seconds.
type DiagnosticRequest struct {
	Level           DiagType
	TestNames       []string
	TestParameters  []string
	Group           *GroupHandle
	Entities        []GroupEntityPair
	Timeout         time.Duration
	EnableHeartbeat bool
}

const diagnosticFeature = "DCGM diagnostic and profiling control API"

// requireDiagnosticSymbol requires an active session and the requested optional native function.
func requireDiagnosticSymbol(symbol string) error {
	if !dcgmSessionSymbolAvailable(symbol) {
		return &FeatureUnavailableError{Feature: diagnosticFeature, Symbol: symbol}
	}

	return nil
}

// requireDiagnosticSymbol checks diagnostic availability before a native call.
func (cgoAdapter) requireDiagnosticSymbol(symbol string) error {
	return requireDiagnosticSymbol(symbol)
}

// diagnosticResult preserves the native failure code with diagnostic context.
func diagnosticResult(symbol string, result C.dcgmReturn_t) error {
	if result == C.DCGM_ST_OK {
		return nil
	}

	return dcgmFeatureResult(diagnosticFeature, symbol, result, dcgmErrorText(result))
}

// diagnosticEntitySelector encodes supported entities in the selector syntax accepted by DCGM.
func diagnosticEntitySelector(entities []GroupEntityPair) (string, error) {
	if len(entities) > int(C.DCGM_DIAG_RESPONSE_ENTITIES_MAX) {
		return "", fmt.Errorf("diagnostic has %d entities; maximum is %d", len(entities), C.DCGM_DIAG_RESPONSE_ENTITIES_MAX)
	}

	selectors := make([]string, len(entities))
	for i, entity := range entities {
		var prefix string
		switch entity.EntityGroupId {
		case FE_GPU:
			prefix = "gpu"
		case FE_SWITCH:
			prefix = "nvswitch"
		case FE_VGPU:
			prefix = "vgpu"
		case FE_GPU_I:
			prefix = "instance"
		case FE_GPU_CI:
			prefix = "ci"
		case FE_CPU:
			prefix = "cpu"
		case FE_CPU_CORE:
			prefix = "core"
		case FE_CONNECTX:
			prefix = "cx"
		case FE_LINK:
			prefix = "link"
		default:
			return "", fmt.Errorf("diagnostic entity %d has unsupported group %d", i, entity.EntityGroupId)
		}
		if entity.EntityId > math.MaxUint32 {
			return "", fmt.Errorf("diagnostic entity %d ID exceeds uint32", i)
		}
		selectors[i] = prefix + ":" + strconv.FormatUint(uint64(entity.EntityId), 10)
	}

	selector := strings.Join(selectors, ",")
	if len(selector) >= int(C.DCGM_ENTITY_ID_LIST_LEN) {
		return "", fmt.Errorf("diagnostic entity selector exceeds %d bytes", C.DCGM_ENTITY_ID_LIST_LEN-1)
	}

	return selector, nil
}

// validateDiagnosticRequest rejects conflicting selections and values that cannot fit the native request.
func validateDiagnosticRequest(request DiagnosticRequest) error {
	if request.Level < 0 || request.Level > DiagExtended {
		return errors.New("invalid diagnostic level")
	}

	if request.Level == 0 && len(request.TestNames) == 0 {
		return errors.New("diagnostic needs a level or at least one test name")
	}

	if request.Group != nil && len(request.Entities) != 0 {
		return errors.New("diagnostic group and entities are mutually exclusive")
	}

	if err := validateDiagnosticTests(request.TestNames, request.TestParameters); err != nil {
		return err
	}
	// DCGM multiplies timeoutSeconds by unsigned int 1000 before sending the request.
	if request.Timeout < 0 || request.Timeout%time.Second != 0 || request.Timeout/time.Second > math.MaxUint32/1000 {
		return errors.New("diagnostic timeout must be whole seconds whose conversion to milliseconds fits uint32")
	}

	_, err := diagnosticEntitySelector(request.Entities)
	return err
}

// validateDiagnosticTests checks test names and parameters against native array and string limits.
func validateDiagnosticTests(names, parameters []string) error {
	if len(names) > int(C.DCGM_MAX_TEST_NAMES) {
		return fmt.Errorf("diagnostic has more than %d test names", C.DCGM_MAX_TEST_NAMES)
	}

	if len(parameters) > int(C.DCGM_MAX_TEST_PARMS) {
		return fmt.Errorf("diagnostic has more than %d parameters", C.DCGM_MAX_TEST_PARMS)
	}

	for _, name := range names {
		if err := validateCString("diagnostic test name", name, int(C.DCGM_MAX_TEST_NAMES_LEN)); err != nil {
			return err
		}
	}

	for _, parameter := range parameters {
		if err := validateCString("diagnostic parameter", parameter, int(C.DCGM_MAX_TEST_PARMS_LEN_V2)); err != nil {
			return err
		}
	}

	return nil
}

// cDiagnosticRequest builds a native v10 request after checking every caller-supplied value.
func cDiagnosticRequest(request DiagnosticRequest) (C.dcgmRunDiag_v10, error) {
	if err := validateDiagnosticRequest(request); err != nil {
		return C.dcgmRunDiag_v10{}, err
	}
	var native C.dcgmRunDiag_v10
	native.version = C.dcgmRunDiag_version10
	native.groupId = C.DCGM_GROUP_ALL_GPUS
	if request.Group != nil {
		native.groupId = request.Group.handle
	}

	if len(request.Entities) != 0 {
		native.groupId = C.DCGM_GROUP_NULL
		selector, err := diagnosticEntitySelector(request.Entities)
		if err != nil {
			return C.dcgmRunDiag_v10{}, err
		}
		copyGoStringToCChars(native.entityIds[:], selector)
	}

	native.validate = C.dcgmPolicyValidation_t(request.Level)
	for i, name := range request.TestNames {
		copyGoStringToCChars(native.testNames[i][:], name)
	}

	for i, parameter := range request.TestParameters {
		copyGoStringToCChars(native.testParms[i][:], parameter)
	}

	native.timeoutSeconds = C.uint(request.Timeout / time.Second)
	if request.EnableHeartbeat {
		native.flags = C.DCGM_RUN_FLAGS_ENABLE_HEARTBEAT
	}

	return native, nil
}

// ActionValidate runs a diagnostic using named tests, parameters, and entity selection.
// It blocks until DCGM finishes; StopDiag may run concurrently. Native errors
// retain their DCGM code, including DCGM_ST_DIAG_STOPPED on cancellation.
// Non-GPU test results, errors, and auxiliary data are available in Tests.
func ActionValidate(request DiagnosticRequest) (DiagResults, error) {
	return actionValidateWithCalls(cgoAdapter{}, request)
}

// actionValidateWithCalls checks the request and session before invoking the selected native adapter.
func actionValidateWithCalls(api diagnosticCalls, request DiagnosticRequest) (DiagResults, error) {
	if err := validateDiagnosticRequest(request); err != nil {
		return DiagResults{}, err
	}

	if err := api.requireDiagnosticSymbol("dcgmActionValidate_v2"); err != nil {
		return DiagResults{}, err
	}

	return api.actionValidateNative(request)
}

// actionValidateNative runs a prepared diagnostic and retains results even when DCGM returns an error.
func (cgoAdapter) actionValidateNative(request DiagnosticRequest) (DiagResults, error) {
	native, err := cDiagnosticRequest(request)
	if err != nil {
		return DiagResults{}, err
	}

	response := new(C.dcgmDiagResponse_v12)
	response.version = C.dcgmDiagResponse_version12
	result := C.dcgmActionValidate_v2(handle.handle, &native, response)
	return newDiagResults(response), diagnosticResult("dcgmActionValidate_v2", result)
}

// StopDiag asks DCGM to stop the active diagnostic. It may run concurrently with
// RunDiag or ActionValidate; keep the DCGM session open until both calls return.
func StopDiag() error { return diagnosticControlWithCalls(cgoAdapter{}, "dcgmStopDiagnostic") }

// DiagSendHeartbeat keeps a diagnostic started with EnableHeartbeat alive.
func DiagSendHeartbeat() error {
	return diagnosticControlWithCalls(cgoAdapter{}, "dcgmDiagSendHeartbeat")
}

// ProfPause pauses DCGM profiling while another profiler uses the GPU counters.
// Call ProfResume when that work finishes.
func ProfPause() error { return diagnosticControlWithCalls(cgoAdapter{}, "dcgmProfPause") }

// ProfResume resumes profiling after ProfPause.
func ProfResume() error { return diagnosticControlWithCalls(cgoAdapter{}, "dcgmProfResume") }

// diagnosticControlWithCalls checks availability before running a diagnostic or profiling control.
func diagnosticControlWithCalls(api diagnosticCalls, symbol string) error {
	if err := api.requireDiagnosticSymbol(symbol); err != nil {
		return err
	}

	return api.diagnosticControlNative(symbol)
}

// diagnosticControlNative dispatches the supported control symbols and preserves native errors.
func (cgoAdapter) diagnosticControlNative(symbol string) error {
	var result C.dcgmReturn_t
	switch symbol {
	case "dcgmStopDiagnostic":
		result = C.dcgmStopDiagnostic(handle.handle)
	case "dcgmDiagSendHeartbeat":
		result = C.dcgmDiagSendHeartbeat(handle.handle)
	case "dcgmProfPause":
		result = C.dcgmProfPause(handle.handle)
	case "dcgmProfResume":
		result = C.dcgmProfResume(handle.handle)
	default:
		return fmt.Errorf("unknown diagnostic control symbol %q", symbol)
	}

	return diagnosticResult(symbol, result)
}

// runDiagWithCalls checks availability before starting a diagnostic at the requested level.
func runDiagWithCalls(api diagnosticCalls, level DiagType, group GroupHandle) (DiagResults, error) {
	if err := api.requireDiagnosticSymbol("dcgmRunDiagnostic"); err != nil {
		return DiagResults{}, err
	}

	return api.runDiagNative(level, group)
}
