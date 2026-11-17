package dcgm

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestAdministrationRuntime checks administration APIs in one live embedded session.
func TestAdministrationRuntime(t *testing.T) {
	t.Setenv("DCGM_DIAG_REMOVED", "go-dcgm-admin-environment-value")
	cleanup, err := InitEmbedded(EmbeddedOptions{Severity: LoggingNone})
	require.NoError(t, err)
	defer cleanup()
	runOnlyWithLiveGPUs(t)

	gpus, err := getSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus)

	t.Run("hostengine", checkHostengineAdministration)
	t.Run("environment", checkHostengineEnvironment)
	t.Run("groups", func(t *testing.T) { checkGroupAdministration(t, gpus[0]) })
	t.Run("field groups", checkFieldGroupAdministration)
	t.Run("GPUs", func(t *testing.T) { checkGPUAdministration(t, gpus) })
}

// checkHostengineAdministration reads health and modules without changing loaded modules.
func checkHostengineAdministration(t *testing.T) {
	t.Helper()
	healthy, err := HostengineIsHealthy()
	require.NoError(t, err)
	require.True(t, healthy)
	require.NoError(t, HostengineSetLoggingSeverity(0, LoggingNone))

	statuses, err := ModuleGetStatuses()
	require.NoError(t, err)
	require.NotEmpty(t, statuses)
	for _, status := range statuses {
		name, err := ModuleIDToName(status.ID)
		require.NoError(t, err)
		require.NotEmpty(t, name)
	}
}

// checkHostengineEnvironment distinguishes a set, empty and unset allowed variable.
func checkHostengineEnvironment(t *testing.T) {
	t.Helper()
	value, err := HostengineEnvironmentVariable("DCGM_DIAG_REMOVED")
	require.NoError(t, err)
	require.Equal(t, "go-dcgm-admin-environment-value", value)

	require.NoError(t, os.Setenv("DCGM_DIAG_REMOVED", ""))
	value, err = HostengineEnvironmentVariable("DCGM_DIAG_REMOVED")
	require.NoError(t, err)
	require.Empty(t, value)

	require.NoError(t, os.Unsetenv("DCGM_DIAG_REMOVED"))
	_, err = HostengineEnvironmentVariable("DCGM_DIAG_REMOVED")
	var native *Error
	require.ErrorAs(t, err, &native)
	require.Equal(t, DCGM_ST_NOT_CONFIGURED, int(native.Code))
}

// checkGroupAdministration creates an owned group and verifies membership and topology reads.
func checkGroupAdministration(t *testing.T, gpu uint) {
	t.Helper()
	group, err := CreateGroup("admin-runtime")
	require.NoError(t, err)
	defer func() { require.NoError(t, DestroyGroup(group)) }()

	require.NoError(t, AddToGroup(group, gpu))
	groups, err := GroupGetAllIDs()
	require.NoError(t, err)
	require.Contains(t, groups, group)

	t.Run("topology", func(t *testing.T) {
		_, topologyErr := GetGroupTopology(group)
		var native *Error
		if errors.As(topologyErr, &native) && int(native.Code) == DCGM_ST_NOT_SUPPORTED {
			t.Skipf("group topology is not supported by this GPU: %v", topologyErr)
		}

		require.NoError(t, topologyErr)
	})

	require.NoError(t, RemoveFromGroup(group, gpu))
	info, err := GetGroupInfo(group)
	require.NoError(t, err)
	require.Empty(t, info.EntityList)
}

// checkFieldGroupAdministration reads an owned field group through both inventory APIs.
func checkFieldGroupAdministration(t *testing.T) {
	t.Helper()
	fields, err := FieldGroupCreate("admin-runtime-fields", []Short{DCGM_FI_DEV_GPU_TEMP})
	require.NoError(t, err)
	defer func() { require.NoError(t, FieldGroupDestroy(fields)) }()

	info, err := FieldGroupGetInfo(fields)
	require.NoError(t, err)
	require.Equal(t, []Short{DCGM_FI_DEV_GPU_TEMP}, info.FieldIDs)

	all, err := FieldGroupGetAll()
	require.NoError(t, err)
	require.Contains(t, all, info)
}

// checkGPUAdministration reads architecture and checks multi-GPU topology when available.
func checkGPUAdministration(t *testing.T, gpus []uint) {
	t.Helper()
	architecture, err := GetGPUChipArchitecture(gpus[0])
	require.NoError(t, err)
	require.NotZero(t, architecture)

	t.Run("selection", func(t *testing.T) {
		if len(gpus) < 2 {
			t.Skip("topology selection of multiple GPUs requires at least two GPUs")
		}

		selected, err := SelectGPUsByTopology(gpus, 2, 0)
		require.NoError(t, err)
		require.NotEmpty(t, selected)
		require.LessOrEqual(t, len(selected), 2)
	})
}

// TestModuleDenylistIsolated changes the module denylist only in a separate test process.
func TestModuleDenylistIsolated(t *testing.T) {
	if os.Getenv("GO_DCGM_DENYLIST_CHILD") == "1" {
		cleanup, err := InitEmbedded(EmbeddedOptions{Denylist: []ModuleID{ModuleHealth}})
		require.NoError(t, err)
		defer cleanup()
		statuses, err := ModuleGetStatuses()
		require.NoError(t, err)
		require.Contains(t, statuses, ModuleStatus{ModuleHealth, ModuleDenylisted})
		require.NoError(t, ModuleDenylist(ModuleDiag))
		statuses, err = ModuleGetStatuses()
		require.NoError(t, err)
		require.Contains(t, statuses, ModuleStatus{ModuleDiag, ModuleDenylisted})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	//nolint:gosec // The subprocess runs this test executable with fixed test arguments.
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestModuleDenylistIsolated$", "-test.v")
	cmd.Env = append(os.Environ(), "GO_DCGM_DENYLIST_CHILD=1")
	output, err := cmd.CombinedOutput()
	require.NoError(t, ctx.Err(), "denylist subprocess timed out: %s", output)
	require.NoError(t, err, "denylist subprocess failed: %s", output)
}

// TestAdministrationNativeErrors checks typed errors returned by invalid native requests.
func TestAdministrationNativeErrors(t *testing.T) {
	cleanup, err := InitEmbedded(EmbeddedOptions{})
	require.NoError(t, err)
	defer cleanup()
	_, err = ModuleIDToName(ModuleID(^uint32(0)))
	var dcgmError *Error
	require.ErrorAs(t, err, &dcgmError)
	require.Equal(t, DCGM_ST_BADPARAM, int(dcgmError.Code))
	_, err = HostengineEnvironmentVariable("GO_DCGM_UNSUPPORTED_ENVIRONMENT_VARIABLE")
	require.ErrorAs(t, err, &dcgmError)
	require.Equal(t, DCGM_ST_BADPARAM, int(dcgmError.Code))
}

// TestAdministrationAfterShutdown checks that public calls reject a closed session.
func TestAdministrationAfterShutdown(t *testing.T) {
	cleanup, err := InitEmbedded(EmbeddedOptions{})
	require.NoError(t, err)
	cleanup()
	require.Nil(t, dcgmLibHandle)
	checkAdministrationWithoutSession(t)
}

// requireLoadedLibraryWithoutSession retains a failed-init library reference and restores its cleanup state.
func requireLoadedLibraryWithoutSession(t *testing.T) {
	t.Helper()
	require.Zero(t, dcgmInitCounter)
	cleanup, err := Init(Standalone)
	require.Error(t, err)
	require.Nil(t, cleanup)
	require.Zero(t, dcgmInitCounter)
	require.NotNil(t, dcgmLibHandle)
	failedLibrary := dcgmLibHandle
	t.Cleanup(func() {
		// Failed Init retains a library reference; release it after recovery's cleanup.
		mux.Lock()
		defer mux.Unlock()
		require.Zero(t, dcgmInitCounter)
		dcgmLibHandle = failedLibrary
		cgoAdapter{}.closeLibrary()
	})
}

// TestAdministrationAfterFailedInit checks rejection after failed initialization and successful recovery.
func TestAdministrationAfterFailedInit(t *testing.T) {
	requireLoadedLibraryWithoutSession(t)
	require.True(t, dcgmSymbolAvailable("dcgmModuleIdToName"))
	_, err := ModuleIDToName(ModuleHealth)
	require.True(t, IsFeatureUnavailable(err), "loaded library without a session must reject administration calls")
	checkAdministrationWithoutSession(t)
	cleanup, err := InitEmbedded(EmbeddedOptions{Severity: LoggingNone})
	require.NoError(t, err)
	t.Cleanup(cleanup)
	name, err := ModuleIDToName(ModuleHealth)
	require.NoError(t, err)
	require.NotEmpty(t, name)
	require.Equal(t, 1, dcgmInitCounter)
}
