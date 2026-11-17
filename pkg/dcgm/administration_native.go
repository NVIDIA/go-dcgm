package dcgm

/*
#include "dcgm_agent.h"
*/
import "C"

// Native aliases keep request layouts in the adapter while allowing test fakes
// to populate the exact responses passed through the conversion functions.
type (
	nativeHostengineHealth    = C.dcgmHostengineHealth_t
	nativeLoggingSeverity     = C.dcgmSettingsSetLoggingSeverity_t
	nativeEnvironmentVariable = C.dcgmEnvVarInfo_t
	nativeModuleStatuses      = C.dcgmModuleGetStatuses_t
	nativeGroupID             = C.dcgmGpuGrp_t
	nativeCount               = C.uint
	nativeFieldGroupInfo      = C.dcgmFieldGroupInfo_t
	nativeFieldGroups         = C.dcgmAllFieldGroup_t
	nativeGroupTopology       = C.dcgmGroupTopology_v2
)

// requireLoadedAdministrationSymbol rejects optional startup calls when their native symbol is absent.
func requireLoadedAdministrationSymbol(symbol string) error {
	if !dcgmSymbolAvailable(symbol) {
		return &FeatureUnavailableError{Feature: "DCGM administration startup", Symbol: symbol}
	}

	return nil
}

// requireAdministrationNativeSymbol requires an active session and the requested native symbol.
func requireAdministrationNativeSymbol(api administrationNativeCalls, symbol string) error {
	if !api.administrationSymbolAvailable(symbol) {
		return &FeatureUnavailableError{Feature: "DCGM administration (requires an initialized session)", Symbol: symbol}
	}

	return nil
}

// administrationNativeResult preserves a native failure status with its administration context.
func administrationNativeResult(api administrationNativeCalls, symbol string, status int32) error {
	if status == 0 {
		return nil
	}

	return dcgmFeatureResult("administration", symbol, C.dcgmReturn_t(status), api.administrationErrorText(status))
}

// administrationSymbolAvailable requires the caller to hold mux through the operation and error lookup.
func (cgoAdapter) administrationSymbolAvailable(symbol string) bool {
	return dcgmInitCounter > 0 && dcgmSymbolAvailable(symbol)
}

// administrationErrorText reads the native description for a returned status.
func (cgoAdapter) administrationErrorText(status int32) string {
	return dcgmErrorText(C.dcgmReturn_t(status))
}

// hostengineHealthNative passes the caller-owned health response to DCGM.
func (cgoAdapter) hostengineHealthNative(response *nativeHostengineHealth) int32 {
	return int32(C.dcgmHostengineIsHealthy(handle.handle, response))
}

// loggingSeverityNative passes the validated logging request to DCGM.
func (cgoAdapter) loggingSeverityNative(request *nativeLoggingSeverity) int32 {
	return int32(C.dcgmHostengineSetLoggingSeverity(handle.handle, request))
}

// environmentVariableNative passes the caller-owned environment request to DCGM.
func (cgoAdapter) environmentVariableNative(request *nativeEnvironmentVariable) int32 {
	return int32(C.dcgmHostengineEnvironmentVariableInfo(handle.handle, request))
}

// moduleStatusesNative asks DCGM to fill the caller-owned module response.
func (cgoAdapter) moduleStatusesNative(response *nativeModuleStatuses) int32 {
	return int32(C.dcgmModuleGetStatuses(handle.handle, response))
}

// moduleNameNative copies the borrowed native module name into Go memory.
//
//nolint:gocritic // cgo generates identical pointer-check expressions for the borrowed C string.
func (cgoAdapter) moduleNameNative(id ModuleID) (string, int32) {
	var name *C.char
	status := C.dcgmModuleIdToName(C.dcgmModuleId_t(id), &name)
	if status != C.DCGM_ST_OK {
		return "", int32(status)
	}

	return C.GoString(name), int32(status)
}

// moduleDenylistNative denylists one module in the active hostengine.
func (cgoAdapter) moduleDenylistNative(id ModuleID) int32 {
	return int32(C.dcgmModuleDenylist(handle.handle, C.dcgmModuleId_t(id)))
}

// groupIDsNative fills a nonempty ID buffer and its returned count.
func (cgoAdapter) groupIDsNative(ids []nativeGroupID, count *nativeCount) int32 {
	return int32(C.dcgmGroupGetAllIds(handle.handle, &ids[0], count))
}

// removeGPUFromGroupNative removes one validated GPU ID from the supplied group.
func (cgoAdapter) removeGPUFromGroupNative(group GroupHandle, gpuID uint) int32 {
	return int32(C.dcgmGroupRemoveDevice(handle.handle, group.handle, C.uint(gpuID)))
}

// fieldGroupInfoNative fills the caller-owned field-group response.
func (cgoAdapter) fieldGroupInfoNative(response *nativeFieldGroupInfo) int32 {
	return int32(C.dcgmFieldGroupGetInfo(handle.handle, response))
}

// fieldGroupsNative fills the caller-owned field-group inventory.
func (cgoAdapter) fieldGroupsNative(response *nativeFieldGroups) int32 {
	return int32(C.dcgmFieldGroupGetAll(handle.handle, response))
}

// groupTopologyNative reads topology into the caller-owned native response.
func (cgoAdapter) groupTopologyNative(group GroupHandle, response *nativeGroupTopology) int32 {
	return int32(C.dcgmGetGroupTopology(handle.handle, group.handle, response))
}

// chipArchitectureNative returns the architecture and native status without interpreting either.
func (cgoAdapter) chipArchitectureNative(gpuID uint) (architecture ChipArchitecture, status int32) {
	var nativeArchitecture C.dcgmChipArchitecture_t
	nativeStatus := C.dcgmGetGpuChipArchitecture(handle.handle, C.uint(gpuID), &nativeArchitecture)
	return ChipArchitecture(nativeArchitecture), int32(nativeStatus)
}

// topologySelectionNative returns the native selection mask and status.
func (cgoAdapter) topologySelectionNative(input uint64, count uint32, hints TopologyHint) (output uint64, status int32) {
	var nativeOutput C.uint64_t
	nativeStatus := C.dcgmSelectGpusByTopology(handle.handle, C.uint64_t(input), C.uint32_t(count), &nativeOutput, C.uint64_t(hints))
	return uint64(nativeOutput), int32(nativeStatus)
}
