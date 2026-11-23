package dcgm

import "time"

// cgoAdapter implements the native call interfaces using DCGM's C API.
type cgoAdapter struct{}

// deviceAttributes holds converted native attributes used to assemble a Device.
type deviceAttributes struct {
	busID       string
	uuid        string
	power       uint
	bar1        uint
	fbTotal     uint
	identifiers DeviceIdentifiers
}

// sessionCalls supplies library and session operations for initialization and shutdown.
type sessionCalls interface {
	loadLibrary() error
	closeLibrary()
	startEmbedded() error
	stopEmbedded() error
	connectStandalone(...string) error
	disconnectStandalone() error
	startHostengine() error
	stopHostengine() error
}

// watchGroupOps supplies group and field operations for setting up a GPU field watch.
type watchGroupOps interface {
	createGroup(string, bool) (GroupHandle, error)
	addToGroup(GroupHandle, uint) error
	destroyGroup(GroupHandle) error
	watchFieldsNative(FieldHandle, GroupHandle, int64, float64, int32) error
	updateAllFields() error
}

// pidWatchOps supplies discovery, group, and field operations for setting up PID watches.
type pidWatchOps interface {
	createGroup(string, bool) (GroupHandle, error)
	getSupportedDevices() ([]uint, error)
	addToGroup(GroupHandle, uint) error
	destroyGroup(GroupHandle) error
	watchPidFieldsNative(GroupHandle, time.Duration, time.Duration, int) error
	updateAllFields() error
}

// deviceStatusOps supplies temporary watches and value queries for device status reads.
type deviceStatusOps interface {
	fieldGroupCreate(string, []Short) (FieldHandle, error)
	watchFields(uint, FieldHandle, string) (GroupHandle, error)
	getLatestValuesForFields(uint, []Short) ([]FieldValue_v1, error)
	fieldGroupDestroy(FieldHandle) error
	destroyGroup(GroupHandle) error
}

// healthGroupOps supplies temporary groups and health queries for a single GPU health check.
type healthGroupOps interface {
	createGroup(string, bool) (GroupHandle, error)
	addToGroup(GroupHandle, uint) error
	destroyGroup(GroupHandle) error
	healthSet(GroupHandle, HealthSystem) error
	healthCheck(GroupHandle) (HealthResponse, error)
}

// policyRegistrationCalls defines methods to register and unregister policy callbacks.
type policyRegistrationCalls interface {
	registerPolicyNative(GroupHandle, uint32, uint64) error
	unregisterPolicyNative(GroupHandle, uint32) error
}

// fieldPolicyListCalls supplies the count and list queries used to retrieve field policies.
type fieldPolicyListCalls interface {
	fieldPolicyCount() (int, error)
	fieldPolicyList(int) ([]FieldPolicy, int, error)
}

// fieldPolicyEventCalls defines methods to register and unregister field policy callbacks.
type fieldPolicyEventCalls interface {
	registerFieldPolicyNative(Short, uint64, uint64) error
	unregisterFieldPolicyNative(Short, uint64) error
}

// multiNodeDiagnosticCalls supplies the native run and stop operations used for cancellation.
type multiNodeDiagnosticCalls interface {
	runMultiNodeDiagnosticNative(MultiNodeDiagnosticRequest) (MultiNodeDiagnosticResults, error)
	stopMultiNodeDiagnosticNative() error
}

// deviceInfoOps provides the GPU reads used to build a Device.
type deviceInfoOps interface {
	readDeviceAttributes(uint) (deviceAttributes, error)
	getSupportedDevices() ([]uint, error)
	getGPUStatus(uint) EntityStatus
	getCPUAffinity(uint) (affinity string, cleanupErr, err error)
	getDeviceTopology(uint) ([]P2PLink, error)
	getPciBandwidth(uint) (bandwidth int64, cleanupErr, err error)
}

// policyReadOps supplies group membership and policy data used to build a PolicyStatus.
type policyReadOps interface {
	getGroupInfo(GroupHandle) (*GroupInfo, error)
	readPolicy(GroupHandle, int) (policySnapshot, error)
}

// topologyOps supplies peer links and field values used to add peer GPU bus IDs.
type topologyOps interface {
	readDeviceTopology(uint) ([]P2PLink, error)
	entitiesGetLatestValues([]GroupEntityPair, []Short, uint) ([]FieldValue_v2, error)
}

// embeddedSessionCalls supplies library ownership and typed embedded startup.
type embeddedSessionCalls interface {
	loadLibrary() error
	closeLibrary()
	startEmbeddedWithOptions(EmbeddedOptions) error
}

// gpuSelector supplies DCGM's topology mask selection.
type gpuSelector interface {
	selectGPUsByTopology(uint64, uint32, TopologyHint) (uint64, error)
}

// embeddedStartCalls supplies native initialization, startup and failure rollback.
type embeddedStartCalls interface {
	initializeDCGM() error
	startEmbeddedNative(EmbeddedOptions) error
	shutdownDCGM() error
}

// administrationNativeCalls supplies native requests, responses and status codes.
type administrationNativeCalls interface {
	administrationSymbolAvailable(string) bool
	administrationErrorText(int32) string
	hostengineHealthNative(*nativeHostengineHealth) int32
	loggingSeverityNative(*nativeLoggingSeverity) int32
	environmentVariableNative(*nativeEnvironmentVariable) int32
	moduleStatusesNative(*nativeModuleStatuses) int32
	moduleNameNative(ModuleID) (string, int32)
	moduleDenylistNative(ModuleID) int32
	groupIDsNative([]nativeGroupID, *nativeCount) int32
	removeGPUFromGroupNative(GroupHandle, uint) int32
	fieldGroupInfoNative(*nativeFieldGroupInfo) int32
	fieldGroupsNative(*nativeFieldGroups) int32
	groupTopologyNative(GroupHandle, *nativeGroupTopology) int32
	chipArchitectureNative(uint) (ChipArchitecture, int32)
	topologySelectionNative(uint64, uint32, TopologyHint) (uint64, int32)
}

// jobStatsCalls supplies job recording and statistics operations.
type jobStatsCalls interface {
	watchJobFieldsNative(GroupHandle, time.Duration, time.Duration, int) error
	jobStartNative(GroupHandle, string) error
	jobStopNative(string) error
	jobGetStatsNative(string) (JobInfo, error)
	jobRemoveNative(string) error
	jobRemoveAllNative() error
}

// diagnosticCalls supplies guarded diagnostic runs and concurrent control calls.
type diagnosticCalls interface {
	requireDiagnosticSymbol(string) error
	actionValidateNative(DiagnosticRequest) (DiagResults, error)
	runDiagNative(DiagType, GroupHandle) (DiagResults, error)
	diagnosticControlNative(string) error
}

// fieldWatchCalls supplies raw watch results so DCGM status codes remain available.
type fieldWatchCalls interface {
	nativeErrorText(int) string
	watchFieldsResult(FieldHandle, GroupHandle, int64, float64, int32) int
	watchFieldValueResult(uint, Short, time.Duration, time.Duration, int) int
}

// fieldPolicyImportCalls supplies counted imports and library capability checks.
type fieldPolicyImportCalls interface {
	nativeErrorText(int) string
	requireFieldPolicySymbol(string) error
	importFieldPoliciesNative(string) (uint, uint, int)
}

// configCalls supplies native configuration buffers and per-GPU status operations.
type configCalls interface {
	configGetNative(GroupHandle, ConfigType, []nativeGPUConfig, configStatusHandle) error
	configSetNative(GroupHandle, *nativeGPUConfig, configStatusHandle) error
	configEnforceNative(GroupHandle, configStatusHandle) error
	configResetNative(GroupHandle, configStatusHandle) error
	configSetWorkloadPowerProfileNative(*nativeWorkloadPowerProfile) error
	workloadPowerProfileInfoNative(uint, *nativeWorkloadProfilesInfo, *nativeWorkloadProfilesStatus) error
	powerProfileNameNative(PowerProfileID) (string, error)
	configStatusCreate() (configStatusHandle, error)
	configStatusDestroy(configStatusHandle) error
	configStatusClear(configStatusHandle) error
	configStatusCount(configStatusHandle) (uint, error)
	configStatusPop(configStatusHandle) (ConfigFieldError, error)
}
