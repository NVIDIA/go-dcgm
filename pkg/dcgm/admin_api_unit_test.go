package dcgm

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestAdministrationConversions checks empty, valid and oversized native administration responses.
func TestAdministrationConversions(t *testing.T) {
	t.Run("modules", func(t *testing.T) {
		statuses, err := moduleStatusesFixture(1)
		require.NoError(t, err)
		require.Equal(t, []ModuleStatus{{ModuleHealth, ModuleDenylisted}}, statuses)

		statuses, err = moduleStatusesFixture(0)
		require.NoError(t, err)
		require.Empty(t, statuses)

		_, err = moduleStatusesFixture(17)
		require.Error(t, err)
	})

	t.Run("groups", func(t *testing.T) {
		groups, err := groupIDsFixture(2)
		require.NoError(t, err)
		require.Equal(t, uintptr(7), groups[0].GetHandle())
		require.Equal(t, uintptr(8), groups[1].GetHandle())

		groups, err = groupIDsFixture(0)
		require.NoError(t, err)
		require.Empty(t, groups)

		_, err = groupIDsFixture(3)
		require.Error(t, err)
	})

	t.Run("field groups", func(t *testing.T) {
		fields, err := fieldGroupsFixture(1, 2)
		require.NoError(t, err)
		require.Equal(t, "g", fields[0].Name)
		require.Equal(t, uintptr(9), fields[0].Handle.GetHandle())
		require.Equal(t, []Short{100, 200}, fields[0].FieldIDs)

		fields, err = fieldGroupsFixture(0, 0)
		require.NoError(t, err)
		require.Empty(t, fields)

		_, err = fieldGroupsFixture(65, 0)
		require.Error(t, err)
		_, err = fieldGroupsFixture(1, 129)
		require.Error(t, err)
	})

	t.Run("topology", func(t *testing.T) {
		topology := groupTopologyFixture()
		require.True(t, topology.NUMAOptimal)
		require.Equal(t, uint64(3), topology.CPUAffinityMask[0])
		require.Equal(t, uint64(1), topology.CPUAffinityMask[len(topology.CPUAffinityMask)-1])
		require.Equal(t, uint64(0x220), topology.SlowestPath)
	})
}

// TestAdministrationErrors checks successful and typed administration status conversion.
func TestAdministrationErrors(t *testing.T) {
	require.NoError(t, administrationErrorFixture(0))
	err := administrationErrorFixture(DCGM_ST_BADPARAM)
	var native *Error
	require.ErrorAs(t, err, &native)
	require.Equal(t, DCGM_ST_BADPARAM, int(native.Code))
	require.True(t, IsFeatureUnavailable(administrationErrorFixture(DCGM_ST_FUNCTION_NOT_FOUND)))
}

// TestEmbeddedInitializationCleanup checks library cleanup after each embedded-startup outcome.
func TestEmbeddedInitializationCleanup(t *testing.T) {
	options := EmbeddedOptions{LogFile: "-", Severity: LoggingWarning, Denylist: []ModuleID{ModuleDiag}}
	failure := errors.New("start failed")
	for _, stage := range []string{"load", "start", "success"} {
		t.Run(stage, func(t *testing.T) {
			api := NewMockEmbeddedSessionCalls(gomock.NewController(t))
			switch stage {
			case "load":
				api.EXPECT().loadLibrary().Return(failure)
			case "start":
				gomock.InOrder(
					api.EXPECT().loadLibrary().Return(nil),
					api.EXPECT().startEmbeddedWithOptions(options).Return(failure),
					api.EXPECT().closeLibrary(),
				)
			case "success":
				gomock.InOrder(
					api.EXPECT().loadLibrary().Return(nil),
					api.EXPECT().startEmbeddedWithOptions(options).Return(nil),
				)
			}
			err := initEmbeddedWithCalls(api, options)
			if stage == "success" {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, failure)
			}
		})
	}
}

// TestEmbeddedNativeRollback checks DCGM rollback when native startup fails.
func TestEmbeddedNativeRollback(t *testing.T) {
	failure := &Error{msg: "native failed", Code: DCGM_ST_BADPARAM}
	cleanupFailure := &Error{msg: "shutdown failed", Code: DCGM_ST_GENERIC_ERROR}
	for _, stage := range []string{"init", "start", "start and rollback fail", "success"} {
		t.Run(stage, func(t *testing.T) {
			api := NewMockEmbeddedStartCalls(gomock.NewController(t))
			options := EmbeddedOptions{}
			switch stage {
			case "init":
				api.EXPECT().initializeDCGM().Return(failure)
			case "start", "start and rollback fail":
				gomock.InOrder(
					api.EXPECT().initializeDCGM().Return(nil),
					api.EXPECT().startEmbeddedNative(options).Return(failure),
					api.EXPECT().shutdownDCGM().Return(pick[error](stage == "start and rollback fail", cleanupFailure, nil)),
				)
			case "success":
				gomock.InOrder(
					api.EXPECT().initializeDCGM().Return(nil),
					api.EXPECT().startEmbeddedNative(options).Return(nil),
				)
			}
			err := startEmbeddedOptionsWithCalls(api, options)
			if stage == "success" {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, failure)
			}
			switch stage {
			case "start":
				require.Same(t, failure, err)
			case "start and rollback fail":
				require.ErrorIs(t, err, cleanupFailure)
				require.ErrorContains(t, err, "shutting down DCGM after failed startup")
			}
		})
	}
}

// TestEmbeddedOptionsValidation rejects options that do not fit the native startup request.
func TestEmbeddedOptionsValidation(t *testing.T) {
	require.NoError(t, validateEmbeddedOptions(EmbeddedOptions{}))
	for _, options := range []EmbeddedOptions{
		{LogFile: "a\x00b"},
		{Severity: LoggingVerbose + 1},
		{Severity: LoggingUnspecified - 1},
		{Denylist: make([]ModuleID, 17)},
		{Denylist: []ModuleID{99}},
	} {
		api := NewMockEmbeddedSessionCalls(gomock.NewController(t))
		require.Error(t, initEmbeddedWithCalls(api, options))
	}
}

// TestInitEmbeddedRejectsActiveSession checks shared-session reference counts and explicit-start rejection.
func TestInitEmbeddedRejectsActiveSession(t *testing.T) {
	mux.Lock()
	saved := dcgmInitCounter
	dcgmInitCounter = 1
	mux.Unlock()
	t.Cleanup(func() {
		mux.Lock()
		dcgmInitCounter = saved
		mux.Unlock()
	})
	cleanup, err := InitEmbedded(EmbeddedOptions{})
	require.Error(t, err)
	require.Nil(t, cleanup)
	cleanup, err = Init(Embedded)
	require.NoError(t, err)
	require.NotNil(t, cleanup)
	require.Equal(t, 2, dcgmInitCounter)
	require.NoError(t, Shutdown())
	require.Equal(t, 1, dcgmInitCounter)
}

// TestAdministrationInputValidation rejects invalid names, log targets and severity values.
func TestAdministrationInputValidation(t *testing.T) {
	for _, test := range []struct {
		name string
		call func() error
	}{
		{"empty environment name", func() error { _, err := HostengineEnvironmentVariable(""); return err }},
		{"NUL in environment name", func() error { _, err := HostengineEnvironmentVariable("a\x00b"); return err }},
		{"long environment name", func() error {
			_, err := HostengineEnvironmentVariable(strings.Repeat("a", 256))
			return err
		}},
		{"negative logging target", func() error { return HostengineSetLoggingSeverity(-1, LoggingInfo) }},
		{"invalid logging severity", func() error { return HostengineSetLoggingSeverity(0, LoggingVerbose+1) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := test.call()
			require.Error(t, err)
			require.False(t, IsFeatureUnavailable(err), "invalid input must be rejected before session lookup")
		})
	}
}

// TestTopologySelection checks mask conversion, empty results, native errors and invalid counts.
func TestTopologySelection(t *testing.T) {
	api := NewMockGPUSelector(gomock.NewController(t))
	api.EXPECT().selectGPUsByTopology(uint64(1)<<63|1, uint32(2), TopologyIgnoreHealth).Return(uint64(1)<<63, nil)
	selected, err := selectGPUsByTopologyWithCalls(api, []uint{0, 63, 63}, 2, TopologyIgnoreHealth)
	require.NoError(t, err)
	require.Equal(t, []uint{63}, selected)
	api.EXPECT().selectGPUsByTopology(uint64(0), uint32(1), TopologyHint(0)).Return(uint64(0), nil)
	selected, err = selectGPUsByTopologyWithCalls(api, nil, 1, 0)
	require.NoError(t, err)
	require.Empty(t, selected)
	failure := errors.New("native topology failed")
	api.EXPECT().selectGPUsByTopology(uint64(2), uint32(1), TopologyHint(0)).Return(uint64(0), failure)
	_, err = selectGPUsByTopologyWithCalls(api, []uint{1}, 1, 0)
	require.ErrorIs(t, err, failure)
	for _, input := range []struct {
		ids   []uint
		count uint
	}{{[]uint{64}, 1}, {nil, 0}, {nil, 65}} {
		_, err := selectGPUsByTopologyWithCalls(api, input.ids, input.count, 0)
		require.Error(t, err)
	}
}
