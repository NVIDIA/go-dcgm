package tests

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func configurationTestGroup(t *testing.T) dcgm.GroupHandle {
	t.Helper()
	if os.Getenv("DCGM_ENABLE_CONFIGURATION_TESTS") != "1" {
		t.Skip("configuration checks require opt-in task test:configuration")
	}
	if os.Geteuid() != 0 {
		t.Skip("configuration checks require root")
	}
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(cleanup)
	return newPidWatchTestGroup(t)
}

func configurationSnapshot(t *testing.T, group dcgm.GroupHandle) (current, target dcgm.GPUConfig) {
	t.Helper()
	live, err := dcgm.ConfigGet(group, dcgm.ConfigCurrentState)
	if err != nil {
		t.Fatalf("read live configuration: %v", err)
	}
	desired, err := dcgm.ConfigGet(group, dcgm.ConfigTargetState)
	if err != nil {
		t.Fatalf("read target configuration: %v", err)
	}
	if len(live) != 1 || len(desired) != 1 {
		t.Fatalf("expected one live and target config, got %d and %d", len(live), len(desired))
	}
	return live[0], desired[0]
}

func TestConfigurationRead(t *testing.T) {
	group := configurationTestGroup(t)
	current, target := configurationSnapshot(t, group)
	if current.GPU != target.GPU {
		t.Fatalf("live GPU %d differs from target GPU %d", current.GPU, target.GPU)
	}
	t.Logf("GPU %d configuration v3 read: power=%v, target power=%v", current.GPU, current.PowerLimit, target.PowerLimit)
}

// TestConfigurationPowerLimit records and enforces the existing power limit without changing its watts.
func TestConfigurationPowerLimit(t *testing.T) {
	group := configurationTestGroup(t)
	beforeLive, beforeTarget := configurationPowerLimitSetup(t, group)
	if err := dcgm.ConfigSet(group, dcgm.GPUConfig{PowerLimit: beforeLive.PowerLimit}); err != nil {
		t.Fatal(err)
	}
	configuredLive, configuredTarget := configurationSnapshot(t, group)
	if !reflect.DeepEqual(configuredLive.PowerLimit, beforeLive.PowerLimit) || !reflect.DeepEqual(configuredTarget.PowerLimit, beforeLive.PowerLimit) {
		t.Fatal("power limit was not read back in live and target state")
	}
	if err := dcgm.ConfigEnforce(group); err != nil {
		t.Fatal(err)
	}
	if err := dcgm.ConfigReset(group); err != nil {
		t.Fatal(err)
	}
	resetLive, resetTarget := configurationSnapshot(t, group)
	if !reflect.DeepEqual(resetLive, beforeLive) {
		t.Fatal("ConfigReset changed live configuration")
	}
	if !reflect.DeepEqual(resetTarget, dcgm.GPUConfig{GPU: beforeTarget.GPU}) {
		t.Fatalf("ConfigReset left a target: %+v", resetTarget)
	}
}

func configurationPowerLimitSetup(t *testing.T, group dcgm.GroupHandle) (beforeLive, beforeTarget dcgm.GPUConfig) {
	t.Helper()
	beforeLive, beforeTarget = configurationSnapshot(t, group)
	if beforeLive.PowerLimit == nil {
		t.Skip("GPU does not report a writable power limit")
	}
	powerOnlyTarget := beforeTarget
	powerOnlyTarget.PowerLimit = nil
	if !reflect.DeepEqual(powerOnlyTarget, dcgm.GPUConfig{GPU: beforeTarget.GPU}) {
		t.Skip("prior target contains non-power settings; enforcing it could change unrelated GPU state")
	}
	if beforeTarget.PowerLimit != nil && !reflect.DeepEqual(beforeLive.PowerLimit, beforeTarget.PowerLimit) {
		t.Skip("prior target power differs from live power; cannot safely enforce and restore it")
	}
	t.Cleanup(func() { restoreConfigurationPowerLimit(t, group, beforeLive, beforeTarget) })
	return beforeLive, beforeTarget
}

func restoreConfigurationPowerLimit(t *testing.T, group dcgm.GroupHandle, beforeLive, beforeTarget dcgm.GPUConfig) {
	t.Helper()
	// Only power was changed. Restore live watts, clear our target, then restore an existing power target.
	if restoreErr := dcgm.ConfigSet(group, dcgm.GPUConfig{PowerLimit: beforeLive.PowerLimit}); restoreErr != nil {
		t.Errorf("restore live power: %v", restoreErr)
	}
	if resetErr := dcgm.ConfigReset(group); resetErr != nil {
		t.Errorf("clear test target: %v", resetErr)
	}
	if beforeTarget.PowerLimit != nil {
		if restoreErr := dcgm.ConfigSet(group, dcgm.GPUConfig{PowerLimit: beforeTarget.PowerLimit}); restoreErr != nil {
			t.Errorf("restore target power: %v", restoreErr)
		}
	}
	afterLive, afterTarget := configurationSnapshot(t, group)
	if !reflect.DeepEqual(beforeLive, afterLive) {
		t.Errorf("live configuration was not restored: before=%+v after=%+v", beforeLive, afterLive)
	}
	if !reflect.DeepEqual(beforeTarget, afterTarget) {
		t.Errorf("target configuration was not restored: before=%+v after=%+v", beforeTarget, afterTarget)
	}
}

func TestConfigurationWorkloadPowerProfiles(t *testing.T) {
	group := configurationTestGroup(t)
	beforeLive, beforeTarget := configurationSnapshot(t, group)
	if !reflect.DeepEqual(beforeTarget, dcgm.GPUConfig{GPU: beforeTarget.GPU}) {
		t.Skip("workload profile test requires an empty private-engine target")
	}
	before := configurationWorkloadProfiles(t, beforeLive.GPU)
	restoreProfiles := true
	t.Cleanup(func() {
		restoreConfigurationWorkloadProfiles(t, group, beforeLive, beforeTarget, before, restoreProfiles)
	})
	// Reapply the existing requested mask so no new workload policy is selected.
	err := dcgm.ConfigSetWorkloadPowerProfile(group, dcgm.WorkloadProfileOverwrite, before.Requested)
	if err != nil {
		var native *dcgm.Error
		if errors.As(err, &native) && int(native.Code) == dcgm.DCGM_ST_NOT_SUPPORTED {
			// DCGM returns NOT_SUPPORTED before merging the target for this single-GPU operation.
			restoreProfiles = false
			t.Skipf("GPU cannot set workload profiles: %v", err)
		}
		t.Fatal(err)
	}
	after, err := dcgm.GetDeviceWorkloadPowerProfileInfo(beforeLive.GPU)
	if err != nil {
		t.Fatal(err)
	}
	if before.Requested != after.Requested || before.Enforced != after.Enforced {
		t.Fatalf("reapplying requested profiles changed GPU policy: before=%+v after=%+v", before, after)
	}
}

func configurationWorkloadProfiles(t *testing.T, gpu uint) dcgm.WorkloadPowerProfileInfo {
	t.Helper()
	before, err := dcgm.GetDeviceWorkloadPowerProfileInfo(gpu)
	if err != nil {
		var native *dcgm.Error
		if dcgm.IsFeatureUnavailable(err) || (errors.As(err, &native) && int(native.Code) == dcgm.DCGM_ST_NOT_SUPPORTED) {
			t.Skipf("workload power profiles unavailable: %v", err)
		}
		t.Fatal(err)
	}
	if len(before.Profiles) == 0 || before.Supported == (dcgm.PowerProfileMask{}) {
		t.Skip("GPU exposes no supported workload power profiles")
	}
	requireConfigurationProfileName(t, dcgm.PowerProfileMaxP)
	for _, profile := range before.Profiles {
		requireConfigurationProfileName(t, profile.ID)
	}
	return before
}

func requireConfigurationProfileName(t *testing.T, id dcgm.PowerProfileID) {
	t.Helper()
	name, err := dcgm.PowerProfileName(id)
	if err != nil {
		var native *dcgm.Error
		if id > dcgm.PowerProfileMIG && errors.As(err, &native) && int(native.Code) == dcgm.DCGM_ST_BADPARAM {
			t.Logf("driver workload profile ID %d has no DCGM name: %v", id, err)
			return
		}
		t.Fatalf("name workload profile ID %d: %v", id, err)
	}
	if name == "" {
		t.Fatalf("empty workload profile name for ID %d", id)
	}
}

func restoreConfigurationWorkloadProfiles(t *testing.T, group dcgm.GroupHandle, beforeLive, beforeTarget dcgm.GPUConfig, before dcgm.WorkloadPowerProfileInfo, restoreProfiles bool) {
	t.Helper()
	if restoreProfiles {
		if restoreErr := dcgm.ConfigSetWorkloadPowerProfile(group, dcgm.WorkloadProfileOverwrite, before.Requested); restoreErr != nil {
			t.Errorf("restore requested power profiles: %v", restoreErr)
		}
	}
	if resetErr := dcgm.ConfigReset(group); resetErr != nil {
		t.Errorf("clear workload profile target: %v", resetErr)
	}
	after, readErr := dcgm.GetDeviceWorkloadPowerProfileInfo(beforeLive.GPU)
	if readErr != nil {
		t.Errorf("read restored power profiles: %v", readErr)
	} else if !reflect.DeepEqual(before, after) {
		t.Errorf("workload power profiles were not restored: before=%+v after=%+v", before, after)
	}
	afterLive, afterTarget := configurationSnapshot(t, group)
	if !reflect.DeepEqual(beforeLive, afterLive) || !reflect.DeepEqual(beforeTarget, afterTarget) {
		t.Errorf("configuration changed after workload profile restoration: live=%+v target=%+v", afterLive, afterTarget)
	}
}
