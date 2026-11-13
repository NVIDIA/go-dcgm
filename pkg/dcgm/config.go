package dcgm

/*
#include "dcgm_agent.h"
#include <dlfcn.h>

// Resolve the 4.8-only function at runtime even when the binary linked an older libdcgm.
typedef dcgmReturn_t (*goConfigResetFn)(dcgmHandle_t, dcgmGpuGrp_t, dcgmStatus_t);
static dcgmReturn_t go_dcgmConfigReset(void *library, dcgmHandle_t handle, dcgmGpuGrp_t group, dcgmStatus_t status) {
    goConfigResetFn reset = (goConfigResetFn)dlsym(library, "dcgmConfigReset");
    if (reset == NULL) return DCGM_ST_FUNCTION_NOT_FOUND;
    return reset(handle, group, status);
}
*/
import "C"

import (
	"errors"
	"fmt"
)

// ConfigType selects the desired target or the live configuration.
type ConfigType int

// Configuration state selectors.
const (
	ConfigTargetState  ConfigType = C.DCGM_CONFIG_TARGET_STATE
	ConfigCurrentState ConfigType = C.DCGM_CONFIG_CURRENT_STATE
)

// ConfigPowerLimitType selects a per-GPU limit or a budget shared by the group.
type ConfigPowerLimitType int

// Power cap allocation modes.
const (
	ConfigPowerCapIndividual ConfigPowerLimitType = C.DCGM_CONFIG_POWER_CAP_INDIVIDUAL
	ConfigPowerBudgetGroup   ConfigPowerLimitType = C.DCGM_CONFIG_POWER_BUDGET_GROUP
)

// ConfigPowerLimit specifies watts, without changing the unit used by DCGM.
type ConfigPowerLimit struct {
	Type  ConfigPowerLimitType
	Watts uint32
}

// PowerProfileMaskSize is the number of 32-bit words in a workload profile mask.
const PowerProfileMaskSize = C.DCGM_POWER_PROFILE_ARRAY_SIZE

// PowerProfileMask contains profile bits. Profile n is bit n%32 of word n/32.
type PowerProfileMask [PowerProfileMaskSize]uint32

// GPUConfig contains version-3 configuration with clocks in Hz and power in watts.
// Nil pointers represent a value that is not set, unavailable or unsupported.
// ConfigSet ignores nil fields; to clear a stored target use ConfigReset.
// GPU identifies a ConfigGet result; ConfigSet applies settings to the whole group.
type GPUConfig struct {
	GPU                             uint
	ECCMode, ComputeMode, SyncBoost *uint32
	MemoryClockHz, SMClockHz        *int64
	PowerLimit                      *ConfigPowerLimit
	WorkloadPowerProfiles           [PowerProfileMaskSize]*uint32
}

// WorkloadProfileAction selects how a mask changes the configured profiles.
type WorkloadProfileAction int

// Workload profile update actions.
const (
	WorkloadProfileSet       WorkloadProfileAction = C.DCGM_WORKLOAD_PROFILE_ACTION_SET
	WorkloadProfileOverwrite WorkloadProfileAction = C.DCGM_WORKLOAD_PROFILE_ACTION_SET_AND_OVERWRITE
	WorkloadProfileClear     WorkloadProfileAction = C.DCGM_WORKLOAD_PROFILE_ACTION_CLEAR
)

// PowerProfileID identifies a DCGM workload power profile.
type PowerProfileID uint32

// Workload power profile identifiers.
const (
	PowerProfileMaxP         PowerProfileID = C.DCGM_POWER_PROFILE_MAX_P
	PowerProfileMaxQ         PowerProfileID = C.DCGM_POWER_PROFILE_MAX_Q
	PowerProfileCompute      PowerProfileID = C.DCGM_POWER_PROFILE_COMPUTE
	PowerProfileMemoryBound  PowerProfileID = C.DCGM_POWER_PROFILE_MEMORY_BOUND
	PowerProfileNetwork      PowerProfileID = C.DCGM_POWER_PROFILE_NETWORK
	PowerProfileBalanced     PowerProfileID = C.DCGM_POWER_PROFILE_BALANCED
	PowerProfileLLMInference PowerProfileID = C.DCGM_POWER_PROFILE_LLM_INFERENCE
	PowerProfileLLMTraining  PowerProfileID = C.DCGM_POWER_PROFILE_LLM_TRAINING
	PowerProfileRBM          PowerProfileID = C.DCGM_POWER_PROFILE_RBM
	PowerProfileDCPCIe       PowerProfileID = C.DCGM_POWER_PROFILE_DCPCIE
	PowerProfileHMMASparse   PowerProfileID = C.DCGM_POWER_PROFILE_HMMA_SPARSE
	PowerProfileHMMADense    PowerProfileID = C.DCGM_POWER_PROFILE_HMMA_DENSE
	PowerProfileSyncBalanced PowerProfileID = C.DCGM_POWER_PROFILE_SYNC_BALANCED
	PowerProfileHPC          PowerProfileID = C.DCGM_POWER_PROFILE_HPC
	PowerProfileMIG          PowerProfileID = C.DCGM_POWER_PROFILE_MIG
)

// WorkloadPowerProfile describes a supported profile and its conflicting profiles.
type WorkloadPowerProfile struct {
	ID              PowerProfileID
	Priority        uint32
	ConflictingMask PowerProfileMask
}

// WorkloadPowerProfileInfo contains profile metadata and the GPU's profile masks.
type WorkloadPowerProfileInfo struct {
	Profiles                       []WorkloadPowerProfile
	Supported, Requested, Enforced PowerProfileMask
}

// ConfigFieldError identifies the GPU, field and underlying DCGM error for a failed setting.
type ConfigFieldError struct {
	GPU     uint
	FieldID Short
	Err     error
}

func (e ConfigFieldError) Error() string {
	return fmt.Sprintf("GPU %d field %d: %v", e.GPU, e.FieldID, e.Err)
}
func (e ConfigFieldError) Unwrap() error { return e.Err }

// ConfigError preserves both the operation error and every per-GPU error.
// A configuration operation can partly succeed; read the configuration again to see what applied.
type ConfigError struct {
	Cause  error
	Fields []ConfigFieldError
}

func (e *ConfigError) Error() string {
	if err := errors.Join(e.Unwrap()...); err != nil {
		return "configuration: " + err.Error()
	}
	return "configuration failed"
}

func (e *ConfigError) Unwrap() []error {
	errs := make([]error, 0, len(e.Fields)+1)
	if e.Cause != nil {
		errs = append(errs, e.Cause)
	}
	for _, field := range e.Fields {
		errs = append(errs, field)
	}
	return errs
}

// ConfigGet reads current or target configuration. It returns partial results alongside errors.
// This uses dcgmConfig_v3 (Hz), and requires a compatible client library and hostengine.
func ConfigGet(group GroupHandle, kind ConfigType) ([]GPUConfig, error) {
	return configGetWithCalls(cgoAdapter{}, group, kind)
}

func configGetWithCalls(api configCalls, group GroupHandle, kind ConfigType) ([]GPUConfig, error) {
	if kind != ConfigTargetState && kind != ConfigCurrentState {
		return nil, errors.New("invalid configuration type")
	}
	configs := make([]nativeGPUConfig, C.DCGM_MAX_NUM_DEVICES)
	for i := range configs {
		configs[i] = blankNativeConfig()
		configs[i].gpuId = C.DCGM_INT32_BLANK
	}
	err := withConfigStatus(api, func(status configStatusHandle) error { return api.configGetNative(group, kind, configs, status) })
	result := make([]GPUConfig, 0)
	for i := range configs {
		if configs[i].gpuId == C.DCGM_INT32_BLANK {
			continue
		}
		result = append(result, configFromNative(&configs[i]))
	}
	return result, err
}

// ConfigSet applies the specified settings to every GPU in group and stores them as targets.
// The target persists beyond the group's lifetime. Clocks must be whole MHz expressed in Hz.
func ConfigSet(group GroupHandle, config GPUConfig) error {
	return configSetWithCalls(cgoAdapter{}, group, config)
}

func configSetWithCalls(api configCalls, group GroupHandle, config GPUConfig) error {
	native, err := configToNative(config)
	if err != nil {
		return err
	}
	return withConfigStatus(api, func(status configStatusHandle) error { return api.configSetNative(group, &native, status) })
}

// ConfigEnforce reapplies the stored target to each GPU in group.
func ConfigEnforce(group GroupHandle) error { return configEnforceWithCalls(cgoAdapter{}, group) }

func configEnforceWithCalls(api configCalls, group GroupHandle) error {
	return withConfigStatus(api, func(status configStatusHandle) error { return api.configEnforceNative(group, status) })
}

// ConfigReset clears the stored target without changing the GPUs' live settings.
// If the library lacks this DCGM 4.8 function, IsFeatureUnavailable reports true for the returned error.
func ConfigReset(group GroupHandle) error { return configResetWithCalls(cgoAdapter{}, group) }

func configResetWithCalls(api configCalls, group GroupHandle) error {
	return withConfigStatus(api, func(status configStatusHandle) error { return api.configResetNative(group, status) })
}

// ConfigSetWorkloadPowerProfile updates workload power profiles for group.
// Set merges mask bits, Overwrite replaces the mask, and Clear removes its bits.
func ConfigSetWorkloadPowerProfile(group GroupHandle, action WorkloadProfileAction, mask PowerProfileMask) error {
	return configSetWorkloadPowerProfileWithCalls(cgoAdapter{}, group, action, mask)
}

func configSetWorkloadPowerProfileWithCalls(api configCalls, group GroupHandle, action WorkloadProfileAction, mask PowerProfileMask) error {
	mux.Lock()
	defer mux.Unlock()
	if action != WorkloadProfileSet && action != WorkloadProfileOverwrite && action != WorkloadProfileClear {
		return errors.New("invalid workload profile action")
	}
	request := nativeWorkloadPowerProfile{version: C.dcgmWorkloadPowerProfile_version1, groupId: group.handle, action: C.dcgmWorkloadProfileAction_t(action)}
	for i := range mask {
		request.profileMask[i] = C.uint(mask[i])
	}
	return api.configSetWorkloadPowerProfileNative(&request)
}

// GetDeviceWorkloadPowerProfileInfo reads supported profiles and current masks.
func GetDeviceWorkloadPowerProfileInfo(gpuID uint) (WorkloadPowerProfileInfo, error) {
	return workloadPowerProfileInfoWithCalls(cgoAdapter{}, gpuID)
}

func workloadPowerProfileInfoWithCalls(api configCalls, gpuID uint) (WorkloadPowerProfileInfo, error) {
	mux.Lock()
	defer mux.Unlock()
	if uint64(gpuID) > uint64(^uint32(0)) {
		return WorkloadPowerProfileInfo{}, errors.New("GPU ID exceeds uint32")
	}
	var profiles nativeWorkloadProfilesInfo
	var status nativeWorkloadProfilesStatus
	profiles.version = C.dcgmWorkloadPowerProfileProfilesInfo_version1
	status.version = C.dcgmDeviceWorkloadPowerProfilesStatus_version1
	for i := range profiles.workloadPowerProfile {
		profiles.workloadPowerProfile[i].version = C.dcgmWorkloadPowerProfileInfo_version1
	}
	if err := api.workloadPowerProfileInfoNative(gpuID, &profiles, &status); err != nil {
		return WorkloadPowerProfileInfo{}, err
	}
	return workloadProfileInfoFromNative(&profiles, &status)
}

// PowerProfileName returns DCGM's friendly name for a workload power profile ID.
func PowerProfileName(id PowerProfileID) (string, error) {
	return powerProfileNameWithCalls(cgoAdapter{}, id)
}

func powerProfileNameWithCalls(api configCalls, id PowerProfileID) (string, error) {
	mux.Lock()
	defer mux.Unlock()
	return api.powerProfileNameNative(id)
}

func withConfigStatus(api configCalls, operation func(configStatusHandle) error) (err error) {
	mux.Lock()
	defer mux.Unlock()
	status, err := api.configStatusCreate()
	if err != nil {
		return fmt.Errorf("create configuration status: %w", err)
	}
	defer func() {
		if destroyErr := api.configStatusDestroy(status); destroyErr != nil {
			err = errors.Join(err, fmt.Errorf("destroy configuration status: %w", destroyErr))
		}
	}()
	if err = api.configStatusClear(status); err != nil {
		return fmt.Errorf("clear configuration status: %w", err)
	}
	operationErr := operation(status)
	count, countErr := api.configStatusCount(status)
	if countErr != nil {
		return errors.Join(operationErr, fmt.Errorf("count configuration errors: %w", countErr))
	}
	// Each native configuration response holds at most DCGM_MAX_NUM_DEVICES status entries.
	if count > C.DCGM_MAX_NUM_DEVICES {
		return errors.Join(operationErr, fmt.Errorf("configuration status count %d exceeds %d", count, C.DCGM_MAX_NUM_DEVICES))
	}
	fields := make([]ConfigFieldError, 0, count)
	for i := uint(0); i < count; i++ {
		field, popErr := api.configStatusPop(status)
		if popErr != nil {
			operationErr = errors.Join(operationErr, fmt.Errorf("pop configuration error: %w", popErr))
			break
		}
		fields = append(fields, field)
	}
	if operationErr != nil || len(fields) > 0 {
		return &ConfigError{Cause: operationErr, Fields: fields}
	}
	return nil
}

type (
	nativeGPUConfig              = C.dcgmConfig_v3
	nativeWorkloadPowerProfile   = C.dcgmWorkloadPowerProfile_v1
	nativeWorkloadProfilesInfo   = C.dcgmWorkloadPowerProfileProfilesInfo_v1
	nativeWorkloadProfilesStatus = C.dcgmDeviceWorkloadPowerProfilesStatus_v1
	configStatusHandle           struct{ handle C.dcgmStatus_t }
)

func blankNativeConfig() nativeGPUConfig {
	var r nativeGPUConfig
	r.version = C.dcgmConfig_version3
	r.eccMode = C.DCGM_INT32_BLANK
	r.computeMode = C.DCGM_INT32_BLANK
	r.perfState.syncBoost = C.DCGM_INT32_BLANK
	r.perfState.targetClocks.version = C.dcgmClockSet_version2
	r.perfState.targetClocks.memClock = C.DCGM_INT64_BLANK
	r.perfState.targetClocks.smClock = C.DCGM_INT64_BLANK
	r.powerLimit._type = C.DCGM_CONFIG_POWER_CAP_INDIVIDUAL
	r.powerLimit.val = C.DCGM_INT32_BLANK
	for i := range r.workloadPowerProfiles {
		r.workloadPowerProfiles[i] = C.DCGM_INT32_BLANK
	}
	return r
}

func configToNative(config GPUConfig) (nativeGPUConfig, error) {
	r := blankNativeConfig()
	if err := configModesToNative(config, &r); err != nil {
		return r, err
	}
	if err := configClocksToNative(config, &r); err != nil {
		return r, err
	}
	if config.PowerLimit != nil {
		power := config.PowerLimit
		if (power.Type != ConfigPowerCapIndividual && power.Type != ConfigPowerBudgetGroup) || power.Watts >= C.DCGM_INT32_BLANK {
			return r, errors.New("invalid power limit type or watts")
		}
		r.powerLimit._type = C.dcgmConfigPowerLimitType_t(power.Type)
		r.powerLimit.val = C.uint(power.Watts)
	}
	for i, word := range config.WorkloadPowerProfiles {
		if word != nil {
			if *word >= C.DCGM_INT32_BLANK {
				return r, fmt.Errorf("workload profile word %d contains a reserved DCGM sentinel", i)
			}
			r.workloadPowerProfiles[i] = C.uint(*word)
		}
	}
	return r, nil
}

func configModesToNative(config GPUConfig, r *nativeGPUConfig) error {
	for _, field := range []struct {
		name   string
		value  *uint32
		target *C.uint
		max    uint32
	}{
		{"ECC mode", config.ECCMode, &r.eccMode, 1},
		{"compute mode", config.ComputeMode, &r.computeMode, C.DCGM_CONFIG_COMPUTEMODE_EXCLUSIVE_PROCESS},
		{"sync boost", config.SyncBoost, &r.perfState.syncBoost, 1},
	} {
		if field.value != nil {
			if *field.value > field.max {
				return fmt.Errorf("invalid %s: %d", field.name, *field.value)
			}
			*field.target = C.uint(*field.value)
		}
	}
	return nil
}

func configClocksToNative(config GPUConfig, r *nativeGPUConfig) error {
	for _, clock := range []struct {
		name   string
		value  *int64
		target *C.longlong
	}{
		{"memory clock", config.MemoryClockHz, &r.perfState.targetClocks.memClock},
		{"SM clock", config.SMClockHz, &r.perfState.targetClocks.smClock},
	} {
		if clock.value != nil {
			hz := *clock.value
			if hz < 0 || hz%1000000 != 0 || hz/1000000 >= C.DCGM_INT32_BLANK {
				return fmt.Errorf("%s must be nonnegative whole MHz in Hz below DCGM's limit", clock.name)
			}
			*clock.target = C.longlong(hz)
		}
	}
	return nil
}

func configFromNative(r *nativeGPUConfig) GPUConfig {
	result := GPUConfig{GPU: uint(r.gpuId), ECCMode: configUint32(r.eccMode), ComputeMode: configUint32(r.computeMode), SyncBoost: configUint32(r.perfState.syncBoost), MemoryClockHz: configInt64(r.perfState.targetClocks.memClock), SMClockHz: configInt64(r.perfState.targetClocks.smClock)}
	if value := configUint32(r.powerLimit.val); value != nil {
		result.PowerLimit = &ConfigPowerLimit{Type: ConfigPowerLimitType(r.powerLimit._type), Watts: *value}
	}
	for i := range result.WorkloadPowerProfiles {
		result.WorkloadPowerProfiles[i] = configUint32(r.workloadPowerProfiles[i])
	}
	return result
}

func configUint32(value C.uint) *uint32 {
	if value >= C.DCGM_INT32_BLANK {
		return nil
	}
	result := uint32(value)
	return &result
}

func configInt64(value C.longlong) *int64 {
	if value >= C.DCGM_INT64_BLANK {
		return nil
	}
	result := int64(value)
	return &result
}

func workloadProfileInfoFromNative(profiles *nativeWorkloadProfilesInfo, status *nativeWorkloadProfilesStatus) (WorkloadPowerProfileInfo, error) {
	if uint64(profiles.profileCount) > uint64(len(profiles.workloadPowerProfile)) {
		return WorkloadPowerProfileInfo{}, fmt.Errorf("workload profile count %d exceeds %d", profiles.profileCount, len(profiles.workloadPowerProfile))
	}
	result := WorkloadPowerProfileInfo{Profiles: make([]WorkloadPowerProfile, int(profiles.profileCount))}
	for i := range result.Supported {
		result.Supported[i] = uint32(status.profileMask[i])
		result.Requested[i] = uint32(status.requestedProfileMask[i])
		result.Enforced[i] = uint32(status.enforcedProfileMask[i])
	}
	for i := range result.Profiles {
		native := &profiles.workloadPowerProfile[i]
		result.Profiles[i].ID = PowerProfileID(native.profileId)
		result.Profiles[i].Priority = uint32(native.priority)
		for j := range result.Profiles[i].ConflictingMask {
			result.Profiles[i].ConflictingMask[j] = uint32(native.conflictingMask[j])
		}
	}
	return result, nil
}

// requireConfigSymbol requires the caller to hold mux through the native call and error lookup.
func requireConfigSymbol(symbol string) error {
	if dcgmInitCounter > 0 && dcgmSymbolAvailable(symbol) {
		return nil
	}
	return &FeatureUnavailableError{Feature: "DCGM configuration API", Symbol: symbol}
}

func configResult(symbol string, code C.dcgmReturn_t) error {
	if code == C.DCGM_ST_OK {
		return nil
	}
	return dcgmFeatureResult("DCGM configuration API", symbol, code, dcgmErrorText(code))
}

func (cgoAdapter) configGetNative(group GroupHandle, kind ConfigType, configs []nativeGPUConfig, status configStatusHandle) error {
	if err := requireConfigSymbol("dcgmConfigGet"); err != nil {
		return err
	}
	if len(configs) == 0 || len(configs) > C.DCGM_MAX_NUM_DEVICES {
		return errors.New("invalid configuration buffer size")
	}
	return configResult("dcgmConfigGet", C.dcgmConfigGet(handle.handle, group.handle, C.dcgmConfigType_t(kind), C.int(len(configs)), &configs[0], status.handle))
}

func (cgoAdapter) configSetNative(group GroupHandle, config *nativeGPUConfig, status configStatusHandle) error {
	if err := requireConfigSymbol("dcgmConfigSet"); err != nil {
		return err
	}
	return configResult("dcgmConfigSet", C.dcgmConfigSet(handle.handle, group.handle, config, status.handle))
}

func (cgoAdapter) configEnforceNative(group GroupHandle, status configStatusHandle) error {
	if err := requireConfigSymbol("dcgmConfigEnforce"); err != nil {
		return err
	}
	return configResult("dcgmConfigEnforce", C.dcgmConfigEnforce(handle.handle, group.handle, status.handle))
}

func (cgoAdapter) configResetNative(group GroupHandle, status configStatusHandle) error {
	if err := requireConfigSymbol("dcgmConfigReset"); err != nil {
		return err
	}
	return configResult("dcgmConfigReset", C.go_dcgmConfigReset(dcgmLibHandle, handle.handle, group.handle, status.handle))
}

func (cgoAdapter) configSetWorkloadPowerProfileNative(request *nativeWorkloadPowerProfile) error {
	if err := requireConfigSymbol("dcgmConfigSetWorkloadPowerProfile"); err != nil {
		return err
	}
	return configResult("dcgmConfigSetWorkloadPowerProfile", C.dcgmConfigSetWorkloadPowerProfile(handle.handle, request))
}

func (cgoAdapter) workloadPowerProfileInfoNative(gpuID uint, profiles *nativeWorkloadProfilesInfo, status *nativeWorkloadProfilesStatus) error {
	if err := requireConfigSymbol("dcgmGetDeviceWorkloadPowerProfileInfo"); err != nil {
		return err
	}
	return configResult("dcgmGetDeviceWorkloadPowerProfileInfo", C.dcgmGetDeviceWorkloadPowerProfileInfo(handle.handle, C.uint(gpuID), profiles, status))
}

func (cgoAdapter) powerProfileNameNative(id PowerProfileID) (string, error) {
	if err := requireConfigSymbol("dcgmPowerProfileIdToName"); err != nil {
		return "", err
	}
	var name *C.char
	code := C.dcgmPowerProfileIdToName(C.dcgmPowerProfileType_t(id), &name) //nolint:gocritic // Cgo generates a duplicate comparison for this char** argument.
	if err := configResult("dcgmPowerProfileIdToName", code); err != nil {
		return "", err
	}
	if name == nil {
		return "", errors.New("DCGM returned a nil power profile name")
	}
	// DCGM owns the static name string.
	return C.GoString(name), nil
}

func (cgoAdapter) configStatusCreate() (configStatusHandle, error) {
	for _, symbol := range []string{"dcgmStatusCreate", "dcgmStatusDestroy", "dcgmStatusClear", "dcgmStatusGetCount", "dcgmStatusPopError"} {
		if err := requireConfigSymbol(symbol); err != nil {
			return configStatusHandle{}, err
		}
	}
	var status configStatusHandle
	err := configResult("dcgmStatusCreate", C.dcgmStatusCreate(&status.handle))
	return status, err
}

func (cgoAdapter) configStatusDestroy(status configStatusHandle) error {
	if err := requireConfigSymbol("dcgmStatusDestroy"); err != nil {
		return err
	}
	return configResult("dcgmStatusDestroy", C.dcgmStatusDestroy(status.handle))
}

func (cgoAdapter) configStatusClear(status configStatusHandle) error {
	if err := requireConfigSymbol("dcgmStatusClear"); err != nil {
		return err
	}
	return configResult("dcgmStatusClear", C.dcgmStatusClear(status.handle))
}

func (cgoAdapter) configStatusCount(status configStatusHandle) (uint, error) {
	if err := requireConfigSymbol("dcgmStatusGetCount"); err != nil {
		return 0, err
	}
	var count C.uint
	err := configResult("dcgmStatusGetCount", C.dcgmStatusGetCount(status.handle, &count))
	return uint(count), err
}

func (cgoAdapter) configStatusPop(status configStatusHandle) (ConfigFieldError, error) {
	if err := requireConfigSymbol("dcgmStatusPopError"); err != nil {
		return ConfigFieldError{}, err
	}
	var info C.dcgmErrorInfo_t
	if err := configResult("dcgmStatusPopError", C.dcgmStatusPopError(status.handle, &info)); err != nil {
		return ConfigFieldError{}, err
	}
	return configFieldErrorFromNative(info, dcgmErrorText(C.dcgmReturn_t(info.status))), nil
}

func configFieldErrorFromNative(info C.dcgmErrorInfo_t, message string) ConfigFieldError {
	return ConfigFieldError{GPU: uint(info.gpuId), FieldID: Short(info.fieldId), Err: &Error{Code: C.dcgmReturn_t(info.status), msg: message}}
}
