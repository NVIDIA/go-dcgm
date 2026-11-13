package dcgm

/*
#include "dcgm_agent.h"
*/
import "C"

func nativeConfigFixtureForTest() nativeGPUConfig {
	r := blankNativeConfig()
	r.gpuId = 7
	r.eccMode = 1
	r.computeMode = 2
	r.perfState.syncBoost = 0
	r.perfState.targetClocks.memClock = 1200000000
	r.perfState.targetClocks.smClock = 900000000
	r.powerLimit._type = C.DCGM_CONFIG_POWER_BUDGET_GROUP
	r.powerLimit.val = 250
	for i := range r.workloadPowerProfiles {
		r.workloadPowerProfiles[i] = C.uint(i + 1)
	}
	return r
}

func configVersionsForTest() (configVersion uint32, clockVersion int) {
	return C.dcgmConfig_version3, C.dcgmClockSet_version2
}

func configNativeBlanksForTest() (blank32 uint32, blank64 int64) {
	return C.DCGM_INT32_BLANK, C.DCGM_INT64_BLANK
}

func configNativeUnavailableForTest() GPUConfig {
	r := blankNativeConfig()
	r.gpuId = 5
	r.eccMode = C.DCGM_INT32_NOT_SUPPORTED
	r.computeMode = C.DCGM_INT32_NOT_FOUND
	r.perfState.syncBoost = C.DCGM_INT32_NOT_PERMISSIONED
	r.perfState.targetClocks.memClock = C.DCGM_INT64_NOT_SUPPORTED
	r.perfState.targetClocks.smClock = C.DCGM_INT64_NOT_PERMISSIONED
	r.powerLimit.val = C.DCGM_INT32_NOT_FOUND
	return configFromNative(&r)
}

func configFieldErrorForTest(gpu uint, field Short, code int) ConfigFieldError {
	return configFieldErrorFromNative(C.dcgmErrorInfo_t{gpuId: C.uint(gpu), fieldId: C.short(field), status: C.int(code)}, "native configuration failure")
}

func configErrorForTest(code int) error {
	return &Error{Code: C.dcgmReturn_t(code), msg: "native configuration failure"}
}

func populateProfilesForTest(profiles *nativeWorkloadProfilesInfo, status *nativeWorkloadProfilesStatus, count int) {
	profiles.profileCount = C.uint(count)
	for i := range profiles.workloadPowerProfile {
		profiles.workloadPowerProfile[i].profileId = C.dcgmPowerProfileType_t(i)
		profiles.workloadPowerProfile[i].priority = C.uint(i + 10)
		for j := range profiles.workloadPowerProfile[i].conflictingMask {
			profiles.workloadPowerProfile[i].conflictingMask[j] = C.uint(i + j + 20)
		}
	}
	for i := range status.profileMask {
		status.profileMask[i] = C.uint(i + 30)
		status.requestedProfileMask[i] = C.uint(i + 40)
		status.enforcedProfileMask[i] = C.uint(i + 50)
	}
}

func profileVersionsForTest() (profilesVersion, infoVersion, statusVersion uint32) {
	return C.dcgmWorkloadPowerProfileProfilesInfo_version1, C.dcgmWorkloadPowerProfileInfo_version1, C.dcgmDeviceWorkloadPowerProfilesStatus_version1
}
func workloadRequestVersionForTest() uint32 { return C.dcgmWorkloadPowerProfile_version1 }
