package dcgm

/*
#include "dcgm_agent.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

// ModuleID identifies a DCGM module.
type ModuleID uint32

// DCGM module identifiers.
const (
	ModuleCore       ModuleID = C.DcgmModuleIdCore
	ModuleNvSwitch   ModuleID = C.DcgmModuleIdNvSwitch
	ModuleVGPU       ModuleID = C.DcgmModuleIdVGPU
	ModuleIntrospect ModuleID = C.DcgmModuleIdIntrospect
	ModuleHealth     ModuleID = C.DcgmModuleIdHealth
	ModulePolicy     ModuleID = C.DcgmModuleIdPolicy
	ModuleConfig     ModuleID = C.DcgmModuleIdConfig
	ModuleDiag       ModuleID = C.DcgmModuleIdDiag
	ModuleProfiling  ModuleID = C.DcgmModuleIdProfiling
	ModuleSysmon     ModuleID = C.DcgmModuleIdSysmon
	ModuleMnDiag     ModuleID = C.DcgmModuleIdMnDiag
)

// ModuleState describes whether a module can be used.
type ModuleState uint32

// DCGM module states.
const (
	ModuleNotLoaded  ModuleState = C.DcgmModuleStatusNotLoaded
	ModuleDenylisted ModuleState = C.DcgmModuleStatusDenylisted
	ModuleFailed     ModuleState = C.DcgmModuleStatusFailed
	ModuleLoaded     ModuleState = C.DcgmModuleStatusLoaded
	ModuleUnloaded   ModuleState = C.DcgmModuleStatusUnloaded
	ModulePaused     ModuleState = C.DcgmModuleStatusPaused
	ModuleReloadable ModuleState = C.DcgmModuleStatusReloadable
)

// ModuleStatus pairs a module identifier with its current state.
type ModuleStatus struct {
	ID    ModuleID
	State ModuleState
}

// LoggingSeverity controls DCGM logging verbosity.
type LoggingSeverity int32

// DCGM logging severities.
const (
	LoggingUnspecified LoggingSeverity = C.DcgmLoggingSeverityUnspecified
	LoggingNone        LoggingSeverity = C.DcgmLoggingSeverityNone
	LoggingFatal       LoggingSeverity = C.DcgmLoggingSeverityFatal
	LoggingError       LoggingSeverity = C.DcgmLoggingSeverityError
	LoggingWarning     LoggingSeverity = C.DcgmLoggingSeverityWarning
	LoggingInfo        LoggingSeverity = C.DcgmLoggingSeverityInfo
	LoggingDebug       LoggingSeverity = C.DcgmLoggingSeverityDebug
	LoggingVerbose     LoggingSeverity = C.DcgmLoggingSeverityVerbose
)

// EmbeddedOptions configures an automatically updating embedded host engine.
// An empty LogFile uses DCGM_LOG_FILE or DCGM's default log path.
// "-" writes to stdout. LoggingNone disables logging.
type EmbeddedOptions struct {
	LogFile  string
	Severity LoggingSeverity
	Denylist []ModuleID
}

// InitEmbedded starts an embedded host engine with typed options and returns
// cleanup as Init does. It requires that no session is already initialized.
// Subsequent Init calls share this session until every caller shuts down.
func InitEmbedded(options EmbeddedOptions) (func(), error) {
	mux.Lock()
	defer mux.Unlock()
	if dcgmInitCounter != 0 {
		return nil, errors.New("InitEmbedded requires no active DCGM session")
	}
	if err := initEmbeddedWithCalls(cgoAdapter{}, options); err != nil {
		return nil, err
	}
	stopMode = Embedded
	dcgmInitCounter = 1
	return func() {
		if err := Shutdown(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to shutdown DCGM: %v\n", err)
		}
	}, nil
}

// initEmbeddedWithCalls loads the library and releases it if embedded startup fails.
func initEmbeddedWithCalls(api embeddedSessionCalls, options EmbeddedOptions) error {
	if err := validateEmbeddedOptions(options); err != nil {
		return err
	}
	if err := api.loadLibrary(); err != nil {
		return err
	}
	if err := api.startEmbeddedWithOptions(options); err != nil {
		api.closeLibrary()
		return err
	}

	return nil
}

// validateEmbeddedOptions rejects log paths, severities and denylists the native request cannot represent.
func validateEmbeddedOptions(options EmbeddedOptions) error {
	if strings.IndexByte(options.LogFile, 0) >= 0 {
		return errors.New("embedded log file must not contain a NUL byte")
	}
	if options.Severity < LoggingUnspecified || options.Severity > LoggingVerbose {
		return newBadParameterError()
	}
	if len(options.Denylist) > int(C.DCGM_MODULE_ID_COUNT_V2) {
		return newBadParameterError()
	}
	for _, id := range options.Denylist {
		if id >= C.DcgmModuleIdCount {
			return newBadParameterError()
		}
	}

	return nil
}

// startEmbeddedWithOptions starts the native session through the replaceable call boundary.
func (api cgoAdapter) startEmbeddedWithOptions(options EmbeddedOptions) error {
	return startEmbeddedOptionsWithCalls(api, options)
}

// startEmbeddedOptionsWithCalls rolls back DCGM initialization when starting the hostengine fails.
func startEmbeddedOptionsWithCalls(api embeddedStartCalls, options EmbeddedOptions) error {
	if err := api.initializeDCGM(); err != nil {
		return err
	}
	if err := api.startEmbeddedNative(options); err != nil {
		return rollbackDCGMInit(err, api.shutdownDCGM)
	}

	return nil
}

// initializeDCGM initializes DCGM only when its library and startup symbol are loaded.
func (cgoAdapter) initializeDCGM() error {
	if err := requireLoadedAdministrationSymbol("dcgmInit"); err != nil {
		return err
	}

	return adminResult("dcgmInit", C.dcgmInit())
}

// shutdownDCGM shuts down DCGM only when its library and shutdown symbol are loaded.
func (cgoAdapter) shutdownDCGM() error {
	if err := requireLoadedAdministrationSymbol("dcgmShutdown"); err != nil {
		return err
	}

	return adminResult("dcgmShutdown", C.dcgmShutdown())
}

//nolint:gocritic // cgo generates identical pointer-check expressions for this parameter struct.
func (cgoAdapter) startEmbeddedNative(options EmbeddedOptions) error {
	if err := requireLoadedAdministrationSymbol("dcgmStartEmbedded_v2"); err != nil {
		return err
	}
	var logFile *C.char
	if options.LogFile != "" {
		logFile = C.CString(options.LogFile)
		defer freeCString(logFile)
	}

	params := embeddedParameters(options, logFile)
	result := C.dcgmStartEmbedded_v2((*C.dcgmStartEmbeddedV2Params_v1)(unsafe.Pointer(&params)))
	if err := adminResult("dcgmStartEmbedded_v2", result); err != nil {
		return err
	}
	handle = dcgmHandle{params.dcgmHandle}
	return nil
}

// embeddedParameters builds the native request while borrowing the supplied C log path.
func embeddedParameters(options EmbeddedOptions, logFile *C.char) C.dcgmStartEmbeddedV2Params_v3 {
	params := C.dcgmStartEmbeddedV2Params_v3{
		version:       C.dcgmStartEmbeddedV2Params_version3,
		opMode:        C.DCGM_OPERATION_MODE_AUTO,
		severity:      C.DcgmLoggingSeverity_t(options.Severity),
		denyListCount: C.uint(len(options.Denylist)),
		logFile:       logFile,
	}
	for i, id := range options.Denylist {
		params.denyList[i] = C.uint(id)
	}

	return params
}

// HostengineIsHealthy reports whether DCGM considers its host engine healthy.
func HostengineIsHealthy() (bool, error) { return cgoAdapter{}.hostengineIsHealthy() }

// hostengineIsHealthy uses the native adapter to read native hostengine health.
func (api cgoAdapter) hostengineIsHealthy() (bool, error) {
	return hostengineIsHealthyWithCalls(api)
}

// hostengineIsHealthyWithCalls checks the session and converts the native health response.
func hostengineIsHealthyWithCalls(api administrationNativeCalls) (bool, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmHostengineIsHealthy"); err != nil {
		return false, err
	}

	response := C.dcgmHostengineHealth_t{version: C.dcgmHostengineHealth_version}
	if err := administrationNativeResult(api, "dcgmHostengineIsHealthy", api.hostengineHealthNative(&response)); err != nil {
		return false, err
	}

	return response.overallHealth == 0, nil
}

// HostengineSetLoggingSeverity changes the selected logger's severity.
// target is the DCGM logger ID (0 for the base logger, 1 for syslog).
// This affects all clients sharing the host engine.
func HostengineSetLoggingSeverity(target int, severity LoggingSeverity) error {
	if target < 0 || target > 1 || severity < LoggingNone || severity > LoggingVerbose {
		return newBadParameterError()
	}

	return cgoAdapter{}.hostengineSetLoggingSeverity(target, severity)
}

// hostengineSetLoggingSeverity uses the native adapter to update a validated logging target.
func (api cgoAdapter) hostengineSetLoggingSeverity(target int, severity LoggingSeverity) error {
	return hostengineSetLoggingSeverityWithCalls(api, target, severity)
}

// hostengineSetLoggingSeverityWithCalls initializes the logging request and preserves native failures.
func hostengineSetLoggingSeverityWithCalls(api administrationNativeCalls, target int, severity LoggingSeverity) error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmHostengineSetLoggingSeverity"); err != nil {
		return err
	}

	request := C.dcgmSettingsSetLoggingSeverity_t{
		version:        C.dcgmSettingsSetLoggingSeverity_version,
		targetLogger:   C.int(target),
		targetSeverity: C.DcgmLoggingSeverity_t(severity),
	}

	return administrationNativeResult(api, "dcgmHostengineSetLoggingSeverity", api.loggingSeverityNative(&request))
}

// HostengineEnvironmentVariable reads one variable from the host engine's
// environment. DCGM restricts which names it allows and returns
// DCGM_ST_NOT_CONFIGURED when an allowed variable is unset.
func HostengineEnvironmentVariable(name string) (string, error) {
	if err := validateCString("environment variable name", name, int(C.DCGM_MAX_STR_LENGTH)); err != nil {
		return "", err
	}
	if name == "" {
		return "", newBadParameterError()
	}

	return cgoAdapter{}.hostengineEnvironmentVariable(name)
}

// hostengineEnvironmentVariable uses the native adapter to read an allowed environment variable.
func (api cgoAdapter) hostengineEnvironmentVariable(name string) (string, error) {
	return hostengineEnvironmentVariableWithCalls(api, name)
}

// hostengineEnvironmentVariableWithCalls checks both call and payload statuses before copying the value.
func hostengineEnvironmentVariableWithCalls(api administrationNativeCalls, name string) (string, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmHostengineEnvironmentVariableInfo"); err != nil {
		return "", err
	}

	request := C.dcgmEnvVarInfo_t{version: C.dcgmEnvVarInfo_version}
	for i := range name {
		request.envVarName[i] = C.char(name[i])
	}
	if err := administrationNativeResult(api, "dcgmHostengineEnvironmentVariableInfo", api.environmentVariableNative(&request)); err != nil {
		return "", err
	}
	if err := administrationNativeResult(api, "dcgmHostengineEnvironmentVariableInfo", int32(request.ret)); err != nil {
		return "", err
	}

	return boundedCString(request.envVarValue[:]), nil
}

// ModuleGetStatuses returns all host engine module states.
func ModuleGetStatuses() ([]ModuleStatus, error) { return cgoAdapter{}.moduleGetStatuses() }

// moduleGetStatuses uses the native adapter to read module states.
func (api cgoAdapter) moduleGetStatuses() ([]ModuleStatus, error) {
	return moduleGetStatusesWithCalls(api)
}

// moduleGetStatusesWithCalls initializes the inventory request and converts bounded module results.
func moduleGetStatusesWithCalls(api administrationNativeCalls) ([]ModuleStatus, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmModuleGetStatuses"); err != nil {
		return nil, err
	}

	response := C.dcgmModuleGetStatuses_t{version: C.dcgmModuleGetStatuses_version}
	if err := administrationNativeResult(api, "dcgmModuleGetStatuses", api.moduleStatusesNative(&response)); err != nil {
		return nil, err
	}

	return moduleStatusesFromNative(response)
}

// moduleStatusesFromNative rejects oversized native counts before copying module states.
func moduleStatusesFromNative(response C.dcgmModuleGetStatuses_t) ([]ModuleStatus, error) {
	if uint64(response.numStatuses) > uint64(len(response.statuses)) {
		return nil, fmt.Errorf("module status count %d exceeds capacity", response.numStatuses)
	}

	statuses := make([]ModuleStatus, int(response.numStatuses))
	for i := range statuses {
		statuses[i] = ModuleStatus{ModuleID(response.statuses[i].id), ModuleState(response.statuses[i].status)}
	}

	return statuses, nil
}

// ModuleIDToName returns the DCGM name of a module.
func ModuleIDToName(id ModuleID) (string, error) { return cgoAdapter{}.moduleIDToName(id) }

// moduleIDToName uses the native adapter to read a module name.
func (api cgoAdapter) moduleIDToName(id ModuleID) (string, error) {
	return moduleIDToNameWithCalls(api, id)
}

// moduleIDToNameWithCalls checks the session and preserves a failed native name lookup.
func moduleIDToNameWithCalls(api administrationNativeCalls, id ModuleID) (string, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmModuleIdToName"); err != nil {
		return "", err
	}

	name, status := api.moduleNameNative(id)
	if err := administrationNativeResult(api, "dcgmModuleIdToName", status); err != nil {
		return "", err
	}

	return name, nil
}

// ModuleDenylist prevents a module from loading. Call it before the module is
// used. This changes the host engine shared by every connected client.
func ModuleDenylist(id ModuleID) error { return cgoAdapter{}.moduleDenylist(id) }

// moduleDenylist uses the native adapter to prevent a module from loading.
func (api cgoAdapter) moduleDenylist(id ModuleID) error {
	return moduleDenylistWithCalls(api, id)
}

// moduleDenylistWithCalls requires an active session before changing the shared module denylist.
func moduleDenylistWithCalls(api administrationNativeCalls, id ModuleID) error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmModuleDenylist"); err != nil {
		return err
	}

	return administrationNativeResult(api, "dcgmModuleDenylist", api.moduleDenylistNative(id))
}

// adminResult converts an administration status into the existing typed error.
func adminResult(symbol string, result C.dcgmReturn_t) error {
	return dcgmFeatureResult("administration", symbol, result, dcgmErrorText(result))
}

// boundedCString copies at most the supplied buffer, stopping at its first NUL byte.
func boundedCString(value []C.char) string {
	for i, c := range value {
		if c == 0 {
			return C.GoStringN(&value[0], C.int(i))
		}
	}
	if len(value) == 0 {
		return ""
	}

	return C.GoStringN(&value[0], C.int(len(value)))
}
