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
	getCPUAffinity(uint) (string, error)
	getDeviceTopology(uint) ([]P2PLink, error)
	getPciBandwidth(uint) (int64, error)
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
