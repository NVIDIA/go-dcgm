package dcgm

/*
#include "dcgm_agent.h"
*/
import "C"

import (
	"strings"
	"time"
)

// diagnosticRequestSentForTest captures the request fields and both versions without
// starting a diagnostic.
func diagnosticRequestSentForTest(request DiagnosticRequest, client, peer VersionInfo) (sent DiagnosticRequest, requestVersion, responseVersion uint, err error) {
	_, err = runDiagnosticRequestWithCall(request, diagnosticRequestVersion(client, peer), func(native *C.dcgmRunDiag_v10, response *C.dcgmDiagResponse_v12) C.dcgmReturn_t {
		requestVersion, responseVersion = uint(native.version), uint(response.version)
		sent.Level = DiagType(native.validate)
		sent.Timeout = time.Duration(native.timeoutSeconds) * time.Second
		sent.EnableHeartbeat = native.flags&C.DCGM_RUN_FLAGS_ENABLE_HEARTBEAT != 0
		group := GroupHandle{native.groupId}
		sent.Group = &group
		for i := range native.testNames {
			if name := C.GoString(&native.testNames[i][0]); name != "" {
				sent.TestNames = append(sent.TestNames, name)
			}
		}
		for i := range native.testParms {
			if parameter := C.GoString(&native.testParms[i][0]); parameter != "" {
				sent.TestParameters = append(sent.TestParameters, parameter)
			}
		}
		return C.DCGM_ST_OK
	})
	return sent, requestVersion, responseVersion, err
}

// diagnosticCompatibilityVersionsForTest returns the native request and response versions used by the tests.
func diagnosticCompatibilityVersionsForTest() (v10, v11, response uint) {
	return C.dcgmRunDiag_version10, C.dcgmRunDiag_version11, C.dcgmDiagResponse_version12
}

type configV2RoundTrip struct {
	config        nativeGPUConfig
	version       uint
	clockVersion  int
	memMHz, smMHz uint
}

// configV2RoundTripForTest records the v2 layout and MHz values before converting the config back.
func configV2RoundTripForTest(config nativeGPUConfig) (configV2RoundTrip, error) {
	legacy, err := configToV2(&config)
	if err != nil {
		return configV2RoundTrip{}, err
	}
	return configV2RoundTrip{
		config:       configFromV2(&legacy),
		version:      uint(legacy.version),
		clockVersion: int(legacy.perfState.targetClocks.version),
		memMHz:       uint(legacy.perfState.targetClocks.memClock),
		smMHz:        uint(legacy.perfState.targetClocks.smClock),
	}, nil
}

// configV2VersionsForTest returns the native versions for the legacy config and its clock fields.
func configV2VersionsForTest() (config uint, clock int) {
	return C.dcgmConfig_version2, C.dcgmClockSet_version1
}

// configClocksForTest creates a config with the given memory and SM clocks in Hz.
func configClocksForTest(mem, sm int64) nativeGPUConfig {
	config := blankNativeConfig()
	config.perfState.targetClocks.memClock = C.longlong(mem)
	config.perfState.targetClocks.smClock = C.longlong(sm)
	return config
}

// fieldMetadataForTest supplies a missing, empty, or populated catalog entry.
func fieldMetadataForTest(present bool, fieldType byte, initialized bool) (FieldMeta, []Short, error) {
	field := DCGM_FI_DEV_GPU_UTIL_RATIO
	var metadata C.dcgm_field_meta_t
	if fieldType != 0 {
		metadata.fieldId = C.ushort(field)
		metadata.fieldType = C.char(fieldType)
		copyGoStringToCChars(metadata.tag[:], "gpu_util_ratio")
	}
	var known C.dcgm_field_meta_t
	if initialized {
		known.fieldType = C.char(DCGM_FT_INT64)
	}

	var lookups []Short
	value, err := fieldMetadataByID(field, func(id Short) C.dcgm_field_meta_p {
		lookups = append(lookups, id)
		if id == DCGM_FI_DEV_GPU_TEMP {
			return &known
		}
		if present {
			return &metadata
		}
		return nil
	})
	return value, lookups, err
}

// configFromV2ClocksForTest sets both v2 clocks to the same MHz value before converting them to Hz.
func configFromV2ClocksForTest(mhz uint32) nativeGPUConfig {
	legacy := C.dcgmConfig_v2{}
	legacy.perfState.targetClocks.memClock = C.uint(mhz)
	legacy.perfState.targetClocks.smClock = C.uint(mhz)
	return configFromV2(&legacy)
}

// diagnosticDispatch records the native entry point and request made by RunDiag.
type diagnosticDispatch struct {
	calls   int
	level   DiagType
	group   GroupHandle
	version uint
}

// runDiagDispatchForTest exercises the production routing with canned native responses.
func runDiagDispatchForTest(level DiagType, group GroupHandle, client, peer VersionInfo, nativeErr error) (DiagResults, diagnosticDispatch, string, error) {
	var sent diagnosticDispatch
	response := new(C.dcgmDiagResponse_v12)
	for i := range response.dcgmVersion {
		response.dcgmVersion[i] = 'v'
	}
	copyGoStringToCChars(response.driverVersion[:], "driver")
	result, err := runDiagForVersion(level, group, diagnosticRequestVersion(client, peer),
		func(level DiagType, group GroupHandle) (*C.dcgmDiagResponse_v12, error) {
			sent.calls++
			sent.level, sent.group = level, group
			return response, nativeErr
		},
		func(request DiagnosticRequest, version C.uint) (*C.dcgmDiagResponse_v12, error) {
			sent.calls++
			sent.level, sent.group, sent.version = request.Level, *request.Group, uint(version)
			return response, nativeErr
		},
	)
	return result, sent, strings.Repeat("v", len(response.dcgmVersion)), err
}
