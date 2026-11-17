package dcgm

/*
#include "dcgm_agent.h"
*/
import "C"

// administrationErrorFixture creates a typed result without calling the native library.
func administrationErrorFixture(code int) error {
	return dcgmFeatureResult("administration", "dcgmExample", C.dcgmReturn_t(code), "fixture status")
}

// moduleStatusesFixture builds a native module response with the requested count, including invalid-count cases.
func moduleStatusesFixture(count uint32) ([]ModuleStatus, error) {
	response := C.dcgmModuleGetStatuses_t{numStatuses: C.uint(count)}
	response.statuses[0].id = C.DcgmModuleIdHealth
	response.statuses[0].status = C.DcgmModuleStatusDenylisted
	return moduleStatusesFromNative(response)
}

// groupIDsFixture converts a fixed native ID buffer using the supplied native count.
func groupIDsFixture(count uint) ([]GroupHandle, error) {
	return groupIDsFromNative([]C.dcgmGpuGrp_t{7, 8}, count)
}

// fieldGroupsFixture builds group and field counts for conversion and bounds checks.
func fieldGroupsFixture(groups, fields uint32) ([]FieldGroupInfo, error) {
	response := C.dcgmAllFieldGroup_t{numFieldGroups: C.uint(groups)}
	response.fieldGroups[0].numFieldIds = C.uint(fields)
	response.fieldGroups[0].fieldGroupId = 9
	response.fieldGroups[0].fieldGroupName[0] = 'g'
	response.fieldGroups[0].fieldIds[0] = 100
	response.fieldGroups[0].fieldIds[1] = 200
	return fieldGroupsFromNative(response)
}

// groupTopologyFixture keeps both ends of the affinity mask and combined path bits observable.
func groupTopologyFixture() GroupTopology {
	response := C.dcgmGroupTopology_v2{numaOptimalFlag: 1, slowestPath: C.DCGM_TOPOLOGY_NVLINK2 | C.DCGM_TOPOLOGY_SYSTEM}
	response.groupCpuAffinityMask[0] = 3
	response.groupCpuAffinityMask[len(response.groupCpuAffinityMask)-1] = 1
	return groupTopologyFromNative(response)
}

// embeddedParametersFixture reads the native initialization request and frees its temporary log path.
func embeddedParametersFixture(options EmbeddedOptions) (version, operationMode uint32, severity int32, file string, denylist []ModuleID) {
	var logFile *C.char
	if options.LogFile != "" {
		logFile = C.CString(options.LogFile)
		defer freeCString(logFile)
	}

	params := embeddedParameters(options, logFile)
	ids := make([]ModuleID, int(params.denyListCount))
	for i := range ids {
		ids[i] = ModuleID(params.denyList[i])
	}

	return uint32(params.version), uint32(params.opMode), int32(params.severity), C.GoString(params.logFile), ids
}

// administrationNativeVersions returns header-defined request versions for adapter assertions.
func administrationNativeVersions() map[string]uint32 {
	return map[string]uint32{
		"embedded":    C.dcgmStartEmbeddedV2Params_version3,
		"operation":   C.DCGM_OPERATION_MODE_AUTO,
		"health":      C.dcgmHostengineHealth_version,
		"logging":     C.dcgmSettingsSetLoggingSeverity_version,
		"environment": C.dcgmEnvVarInfo_version,
		"modules":     C.dcgmModuleGetStatuses_version,
		"fields":      C.dcgmFieldGroupInfo_version,
		"fieldGroups": C.dcgmAllFieldGroup_version,
		"topology":    C.dcgmGroupTopology_version2,
	}
}

// setEnvironmentVariableFixture fills a native environment response and its payload status.
func setEnvironmentVariableFixture(response *nativeEnvironmentVariable, value string, status int32) {
	for i := range value {
		response.envVarValue[i] = C.char(value[i])
	}
	response.ret = C.dcgmReturn_t(status)
}

// environmentVariableNameFixture reads the bounded variable name sent to the native adapter.
func environmentVariableNameFixture(response *nativeEnvironmentVariable) string {
	return boundedCString(response.envVarName[:])
}

// setModuleStatusFixture fills one module entry in a native response.
func setModuleStatusFixture(response *nativeModuleStatuses, id ModuleID, state ModuleState) {
	response.numStatuses = 1
	response.statuses[0].id = C.dcgmModuleId_t(id)
	response.statuses[0].status = C.dcgmModuleStatus_t(state)
}
