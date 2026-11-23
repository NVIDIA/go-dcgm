package dcgm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	dcgmInitCounter int
	mux             sync.Mutex
)

// Init starts DCGM in the specified mode
// Mode can be:
// - Embedded: Start hostengine within this process
// - Standalone: Connect to an already running nv-hostengine
// - StartHostengine: Start and connect to nv-hostengine, terminate before exiting
// Returns a cleanup function on success. On error, cleanup is nil. If the DCGM
// library cannot be loaded, the error is a *LibraryLoadError. Invalid modes
// return ErrInvalidMode, including during an active session.
func Init(m mode, args ...string) (cleanup func(), err error) {
	mux.Lock()
	defer mux.Unlock()
	if m != Embedded && m != Standalone && m != StartHostengine {
		return nil, ErrInvalidMode
	}

	if dcgmInitCounter < 0 {
		return nil, fmt.Errorf("shutdown() is called %d times, before init()", -dcgmInitCounter)
	}

	if dcgmInitCounter == 0 {
		err = initDCGM(m, args...)
		if err != nil {
			return nil, err
		}
	}

	dcgmInitCounter += 1

	return func() {
		if shutdownErr := Shutdown(); shutdownErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to shutdown DCGM with error: `%v`", shutdownErr)
		}
	}, err
}

// Shutdown stops DCGM and destroys all connections
// Returns an error if DCGM is not initialized
func Shutdown() (err error) {
	mux.Lock()
	defer mux.Unlock()

	if dcgmInitCounter <= 0 {
		return errors.New("init() needs to be called before shutdown()")
	}

	if dcgmInitCounter == 1 {
		err = shutdown()
	}

	dcgmInitCounter -= 1

	return
}

// dcgmSessionSymbolAvailable requires an active session as well as a loaded symbol.
func dcgmSessionSymbolAvailable(symbol string) bool {
	mux.Lock()
	defer mux.Unlock()
	return dcgmInitCounter > 0 && dcgmSymbolAvailable(symbol)
}

// GetAllDeviceCount returns the count of all GPUs in the system
func GetAllDeviceCount() (uint, error) {
	return getAllDeviceCount()
}

// GetEntityGroupEntities returns the entities of the specified group type, up to
// DCGM_ENTITY_GROUP_ENTITIES_MAX of them.
//
// The client library and hostengine must use compatible enumeration protocols.
// Incompatible pairings can return DCGM_ST_VER_MISMATCH. The DCGM 4.8 client
// retries the v2 protocol when a hostengine rejects v3; v2 holds at most 1,024 entities.
func GetEntityGroupEntities(entityGroup Field_Entity_Group) ([]uint, error) {
	return getEntityGroupEntities(entityGroup)
}

// GetSupportedDevices returns a list of DCGM-supported GPU IDs
func GetSupportedDevices() ([]uint, error) {
	return getSupportedDevices()
}

// GetDeviceInfo returns detailed information about the specified GPU.
// A successful read can return a populated Device alongside a temporary-group cleanup error.
func GetDeviceInfo(gpuID uint) (Device, error) {
	return getDeviceInfo(gpuID)
}

// GetGPUStatus returns the entity status of the specified GPU
func GetGPUStatus(gpuID uint) EntityStatus {
	return getGPUStatus(gpuID)
}

// GetDeviceStatus returns current status information about the specified GPU.
// Unavailable readings carry DCGM blank markers. Other per-field failures
// and cleanup failures can return an error alongside the populated status.
// A malformed field response returns an error.
func GetDeviceStatus(gpuID uint) (DeviceStatus, error) {
	return latestValuesForDevice(gpuID)
}

// GetDeviceTopology returns the topology (connectivity) information for the specified GPU
func GetDeviceTopology(gpuID uint) ([]P2PLink, error) {
	return getDeviceTopology(gpuID)
}

// WatchPidFields configures DCGM to start recording stats for GPU processes
// Must be called before GetProcessInfo.
//
// Important: The returned GroupHandle should be cleaned up by calling DestroyGroup
// when monitoring is no longer needed to prevent resource leaks.
//
// Example:
//
//	group, err := dcgm.WatchPidFields()
//	if err != nil {
//	    return err
//	}
//	defer dcgm.DestroyGroup(group)
//
//	// Use GetProcessInfo with the group...
func WatchPidFields() (GroupHandle, error) {
	return watchPidFields(time.Microsecond*time.Duration(defaultUpdateFreq), time.Second*time.Duration(defaultMaxKeepAge), defaultMaxKeepSamples)
}

// WatchPidFieldsForGroup is like WatchPidFields, but records process stats for
// the GPUs in a group the caller already created, instead of creating a new
// group. Call it before GetProcessInfo with the same group.
//
// The caller keeps ownership of the group and destroys it with DestroyGroup.
//
// Example:
//
//	group, err := dcgm.CreateGroup("pid-watch-gpus")
//	if err != nil {
//	    return err
//	}
//	defer dcgm.DestroyGroup(group)
//
//	if err := dcgm.AddToGroup(group, gpuID); err != nil {
//	    return err
//	}
//	if err := dcgm.WatchPidFieldsForGroup(group); err != nil {
//	    return err
//	}
//
//	// Use GetProcessInfo with the group...
func WatchPidFieldsForGroup(group GroupHandle) error {
	return WatchPidFieldsForGroupEx(group, time.Microsecond*time.Duration(defaultUpdateFreq), time.Second*time.Duration(defaultMaxKeepAge), defaultMaxKeepSamples)
}

// GetProcessInfo returns detailed per-GPU statistics for the specified process
func GetProcessInfo(group GroupHandle, pid uint) ([]ProcessInfo, error) {
	return getProcessInfo(group, pid)
}

// HealthCheckByGpuId performs a health check on the specified GPU.
// A successful check can return a result alongside a temporary-group cleanup error.
func HealthCheckByGpuId(gpuID uint) (DeviceHealth, error) {
	return healthCheckByGpuId(gpuID)
}

// ListenForPolicyViolations sets up monitoring for the specified policy conditions on all GPUs.
// Returns a channel that receives policy violations and any error encountered. Delivery is
// best-effort: callers must drain the returned channel promptly, or matching violations may be
// dropped for that caller. Use PolicyViolationDropCount to observe local drops.
//
// Important: The context MUST be cancelled when monitoring is no longer needed to properly
// clean up resources and prevent goroutine leaks. When the context is cancelled, the returned
// channel will be closed after local cleanup completes. Concurrent listeners are independent
// subscribers; each active listener receives matching violations instead of sharing one queue.
// Empty condition lists and unknown policy conditions return an error before registering with DCGM.
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel() // Ensures cleanup happens
//
//	violations, err := dcgm.ListenForPolicyViolations(ctx, dcgm.XidPolicy)
//	if err != nil {
//	    return err
//	}
//
//	for violation := range violations {
//	    // Handle violation...
//	}
func ListenForPolicyViolations(ctx context.Context, typ ...policyCondition) (<-chan PolicyViolation, error) {
	groupID := GroupAllGPUs()
	return ListenForPolicyViolationsForGroup(ctx, groupID, typ...)
}

// ListenForPolicyViolationsForGroup sets up policy monitoring for the specified GPU group.
// Returns a best-effort channel that receives policy violations and any error encountered.
//
// Important: The context MUST be cancelled when monitoring is no longer needed to properly
// clean up resources and prevent goroutine leaks. Canceling one listener only closes that
// listener's channel; surviving listeners remain registered until their contexts are canceled.
// Empty condition lists and unknown policy conditions return an error before registering with DCGM.
// See ListenForPolicyViolations for usage example.
func ListenForPolicyViolationsForGroup(ctx context.Context, group GroupHandle, typ ...policyCondition) (<-chan PolicyViolation, error) {
	return registerPolicy(ctx, group, typ...)
}

// Introspect returns memory and CPU usage statistics for the DCGM hostengine
func Introspect() (Status, error) {
	return introspect()
}

// GetVersionInfo returns build information for the DCGM client library.
// If the library or version function is unavailable, IsFeatureUnavailable
// returns true for the returned error.
func GetVersionInfo() (VersionInfo, error) {
	mux.Lock()
	defer mux.Unlock()
	return versionInfo()
}

// GetLoadedLibraryPath returns the path of the loaded DCGM client library as
// recorded by the Linux dynamic loader. The path may be relative or contain
// symlinks. The returned string remains valid after Shutdown.
// GetLoadedLibraryPath returns an error when the library has not been loaded,
// after final Shutdown, or when the platform cannot report the path.
func GetLoadedLibraryPath() (string, error) {
	mux.Lock()
	defer mux.Unlock()
	return loadedLibraryPath()
}

// GetHostengineVersionInfo returns build information for the active DCGM hostengine session.
// If the session or version function is unavailable, IsFeatureUnavailable
// returns true for the returned error.
func GetHostengineVersionInfo() (VersionInfo, error) {
	mux.Lock()
	defer mux.Unlock()
	return hostengineVersionInfo()
}

// GetErrorMeta returns metadata for a DCGM health or diagnostic error code.
// It returns nil when code is outside the range recognized by DCGM.
// DCGM must be initialized before this function is called.
func GetErrorMeta(code HealthCheckErrorCode) *ErrorMeta {
	return getErrorMeta(code)
}

// GetSupportedMetricGroups returns all supported metric groups for the specified GPU
func GetSupportedMetricGroups(gpuID uint) ([]MetricGroup, error) {
	return getSupportedMetricGroups(gpuID)
}

// GetNvLinkLinkStatus returns the status of all NVLink connections
func GetNvLinkLinkStatus() ([]NvLinkStatus, error) {
	return getNvLinkLinkStatus()
}

// GetNvLinkP2PStatus returns the status of NvLinks between GPU pairs
func GetNvLinkP2PStatus() (NvLinkP2PStatus, error) {
	return getNvLinkP2PStatus()
}

// SetPolicyForGroup configures policies with optional custom thresholds and actions for a GPU group
func SetPolicyForGroup(group GroupHandle, configs ...PolicyConfig) error {
	return setPolicyForGroupWithConfig(group, configs...)
}

// GetPolicyForGroup retrieves the current policy configuration for a GPU group
func GetPolicyForGroup(group GroupHandle) (*PolicyStatus, error) {
	return getPolicyForGroup(group)
}

// ClearPolicyForGroup clears all policy conditions for a GPU group
func ClearPolicyForGroup(group GroupHandle) error {
	return clearPolicyForGroup(group)
}

// WatchPolicyViolationsForGroup registers to receive violation notifications for a specific GPU group.
// Unlike ListenForPolicyViolationsForGroup, it does not set policy thresholds first. Delivery is
// best-effort, the context must be canceled to release resources, and canceling one watcher does
// not stop surviving watchers. Empty condition lists and unknown policy conditions return an error
// before registering with DCGM.
func WatchPolicyViolationsForGroup(ctx context.Context, group GroupHandle, typ ...PolicyCondition) (<-chan PolicyViolation, error) {
	return registerPolicyOnly(ctx, group, typ...)
}

// PolicyViolationDropCount returns the number of local policy violations dropped because
// listener channels were full. The counter is process-wide and monotonically increasing.
func PolicyViolationDropCount() uint64 {
	return policyCallbacks.dropped()
}
