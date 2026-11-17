package dcgm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestEmbeddedParameters checks that native startup fields preserve the supplied options.
func TestEmbeddedParameters(t *testing.T) {
	for _, options := range []EmbeddedOptions{
		{},
		{LogFile: "/tmp/unique-dcgm-test.log", Severity: LoggingDebug, Denylist: []ModuleID{ModuleHealth, ModuleDiag}},
	} {
		version, mode, severity, logFile, denylist := embeddedParametersFixture(options)
		expected := administrationNativeVersions()
		require.Equal(t, expected["embedded"], version)
		require.Equal(t, expected["operation"], mode)
		require.Equal(t, int32(options.Severity), severity)
		require.Equal(t, options.LogFile, logFile)
		require.Len(t, denylist, len(options.Denylist))
		require.Equal(t, append([]ModuleID{}, options.Denylist...), denylist)
	}
}

// administrationNativeCallCase pairs one API with its native request and expected Go value.
type administrationNativeCallCase struct {
	name   string
	symbol string
	expect func(*testing.T, *MockAdministrationNativeCalls, int32)
	call   func(administrationNativeCalls) (any, error)
	want   any
}

// administrationNativeCallCases keeps request assertions beside the API and its expected result.
func administrationNativeCallCases() []administrationNativeCallCase {
	versions := administrationNativeVersions()
	group := GroupHandle{}
	group.SetHandle(7)
	fields := FieldHandle{}
	fields.SetHandle(9)
	return []administrationNativeCallCase{
		{
			name:   "health",
			symbol: "dcgmHostengineIsHealthy",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().hostengineHealthNative(gomock.Any()).DoAndReturn(func(response *nativeHostengineHealth) int32 {
					requireAdministrationSessionLocked(t)
					require.Equal(t, versions["health"], uint32(response.version))
					response.overallHealth = 0
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return hostengineIsHealthyWithCalls(api) },
			want: true,
		},
		{
			name:   "logging",
			symbol: "dcgmHostengineSetLoggingSeverity",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().loggingSeverityNative(gomock.Any()).DoAndReturn(func(request *nativeLoggingSeverity) int32 {
					requireAdministrationSessionLocked(t)
					require.Equal(t, versions["logging"], uint32(request.version))
					require.Equal(t, 1, int(request.targetLogger))
					require.Equal(t, int32(LoggingDebug), int32(request.targetSeverity))
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) {
				return nil, hostengineSetLoggingSeverityWithCalls(api, 1, LoggingDebug)
			},
			want: nil,
		},
		{
			name:   "environment",
			symbol: "dcgmHostengineEnvironmentVariableInfo",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().environmentVariableNative(gomock.Any()).DoAndReturn(func(response *nativeEnvironmentVariable) int32 {
					requireAdministrationSessionLocked(t)
					require.Equal(t, versions["environment"], uint32(response.version))
					require.Equal(t, "DCGM_DIAG_REMOVED", environmentVariableNameFixture(response))
					setEnvironmentVariableFixture(response, "distinct-value", 0)
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) {
				return hostengineEnvironmentVariableWithCalls(api, "DCGM_DIAG_REMOVED")
			},
			want: "distinct-value",
		},
		{
			name:   "modules",
			symbol: "dcgmModuleGetStatuses",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().moduleStatusesNative(gomock.Any()).DoAndReturn(func(response *nativeModuleStatuses) int32 {
					requireAdministrationSessionLocked(t)
					require.Equal(t, versions["modules"], uint32(response.version))
					setModuleStatusFixture(response, ModuleHealth, ModuleDenylisted)
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return moduleGetStatusesWithCalls(api) },
			want: []ModuleStatus{{ModuleHealth, ModuleDenylisted}},
		},
		{
			name:   "module name",
			symbol: "dcgmModuleIdToName",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().moduleNameNative(ModuleHealth).DoAndReturn(func(ModuleID) (string, int32) {
					requireAdministrationSessionLocked(t)
					return "Health", status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return moduleIDToNameWithCalls(api, ModuleHealth) },
			want: "Health",
		},
		{
			name:   "denylist",
			symbol: "dcgmModuleDenylist",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().moduleDenylistNative(ModuleHealth).DoAndReturn(func(ModuleID) int32 {
					requireAdministrationSessionLocked(t)
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) {
				return nil, moduleDenylistWithCalls(api, ModuleHealth)
			},
			want: nil,
		},
		{
			name:   "group IDs",
			symbol: "dcgmGroupGetAllIds",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().groupIDsNative(gomock.Any(), gomock.Any()).DoAndReturn(func(ids []nativeGroupID, count *nativeCount) int32 {
					requireAdministrationSessionLocked(t)
					require.GreaterOrEqual(t, len(ids), 1)
					require.Equal(t, len(ids), int(*count))
					ids[0] = 7
					*count = 1
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return groupGetAllIDsWithCalls(api) },
			want: []GroupHandle{group},
		},
		{
			name:   "remove GPU",
			symbol: "dcgmGroupRemoveDevice",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().removeGPUFromGroupNative(group, uint(3)).DoAndReturn(func(GroupHandle, uint) int32 {
					requireAdministrationSessionLocked(t)
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return nil, removeFromGroupWithCalls(api, group, 3) },
			want: nil,
		},
		{
			name:   "field group",
			symbol: "dcgmFieldGroupGetInfo",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().fieldGroupInfoNative(gomock.Any()).DoAndReturn(func(response *nativeFieldGroupInfo) int32 {
					requireAdministrationSessionLocked(t)
					require.Equal(t, versions["fields"], uint32(response.version))
					require.Equal(t, uintptr(9), uintptr(response.fieldGroupId))
					response.numFieldIds = 1
					response.fieldIds[0] = 100
					response.fieldGroupName[0] = 'g'
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return fieldGroupGetInfoWithCalls(api, fields) },
			want: FieldGroupInfo{Handle: fields, Name: "g", FieldIDs: []Short{100}},
		},
		{
			name:   "field groups",
			symbol: "dcgmFieldGroupGetAll",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().fieldGroupsNative(gomock.Any()).DoAndReturn(func(response *nativeFieldGroups) int32 {
					requireAdministrationSessionLocked(t)
					require.Equal(t, versions["fieldGroups"], uint32(response.version))
					response.numFieldGroups = 1
					response.fieldGroups[0].fieldGroupId = 9
					response.fieldGroups[0].numFieldIds = 1
					response.fieldGroups[0].fieldIds[0] = 100
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return fieldGroupGetAllWithCalls(api) },
			want: []FieldGroupInfo{{Handle: fields, FieldIDs: []Short{100}}},
		},
		{
			name:   "group topology",
			symbol: "dcgmGetGroupTopology",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().groupTopologyNative(group, gomock.Any()).DoAndReturn(func(_ GroupHandle, response *nativeGroupTopology) int32 {
					requireAdministrationSessionLocked(t)
					require.Equal(t, versions["topology"], uint32(response.version))
					response.numaOptimalFlag = 1
					return status
				})
			},
			call: func(api administrationNativeCalls) (any, error) {
				topology, err := getGroupTopologyWithCalls(api, group)
				return topology.NUMAOptimal, err
			},
			want: true,
		},
		{
			name:   "architecture",
			symbol: "dcgmGetGpuChipArchitecture",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().chipArchitectureNative(uint(3)).DoAndReturn(func(uint) (ChipArchitecture, int32) {
					requireAdministrationSessionLocked(t)
					return ChipHopper, status
				})
			},
			call: func(api administrationNativeCalls) (any, error) { return getGPUChipArchitectureWithCalls(api, 3) },
			want: ChipHopper,
		},
		{
			name:   "selection",
			symbol: "dcgmSelectGpusByTopology",
			expect: func(t *testing.T, api *MockAdministrationNativeCalls, status int32) {
				api.EXPECT().topologySelectionNative(uint64(5), uint32(1), TopologyIgnoreHealth).DoAndReturn(func(uint64, uint32, TopologyHint) (uint64, int32) {
					requireAdministrationSessionLocked(t)
					return uint64(4), status
				})
			},
			call: func(api administrationNativeCalls) (any, error) {
				return selectGPUMaskWithCalls(api, 5, 1, TopologyIgnoreHealth)
			},
			want: uint64(4),
		},
	}
}

// TestAdministrationNativeCalls checks request data and status handling for every administration API.
func TestAdministrationNativeCalls(t *testing.T) {
	for _, test := range administrationNativeCallCases() {
		t.Run(test.name, func(t *testing.T) {
			for _, status := range []int32{0, DCGM_ST_BADPARAM, DCGM_ST_FUNCTION_NOT_FOUND} {
				t.Run(statusName(status), func(t *testing.T) {
					checkAdministrationNativeCall(t, test, status)
				})
			}

			t.Run("missing symbol", func(t *testing.T) {
				api := NewMockAdministrationNativeCalls(gomock.NewController(t))
				api.EXPECT().administrationSymbolAvailable(test.symbol).DoAndReturn(func(string) bool {
					requireAdministrationSessionLocked(t)
					return false
				})

				_, err := test.call(api)
				requireAdministrationSessionUnlocked(t)
				require.True(t, IsFeatureUnavailable(err))
			})
		})
	}
}

// checkAdministrationNativeCall verifies a request and the status returned through its Go API.
func checkAdministrationNativeCall(t *testing.T, test administrationNativeCallCase, status int32) {
	t.Helper()
	api := NewMockAdministrationNativeCalls(gomock.NewController(t))
	api.EXPECT().administrationSymbolAvailable(test.symbol).DoAndReturn(func(string) bool {
		requireAdministrationSessionLocked(t)
		return true
	})
	test.expect(t, api, status)
	if status != 0 {
		api.EXPECT().administrationErrorText(status).DoAndReturn(func(int32) string {
			requireAdministrationSessionLocked(t)
			return "fake native error"
		})
	}

	got, err := test.call(api)
	requireAdministrationSessionUnlocked(t)
	if status == 0 {
		require.NoError(t, err)
		require.Equal(t, test.want, got)
		return
	}

	var native *Error
	require.ErrorAs(t, err, &native)
	require.Equal(t, status, int32(native.Code))
	require.Equal(t, "fake native error", native.Error())
	require.Equal(t, status == DCGM_ST_FUNCTION_NOT_FOUND, IsFeatureUnavailable(err))
}

// statusName gives native result cases readable subtest names.
func statusName(status int32) string {
	switch status {
	case 0:
		return "success"
	case DCGM_ST_BADPARAM:
		return "typed error"
	default:
		return "unavailable"
	}
}

// TestAdministrationEnvironmentPayloadStatus checks the payload status independently of the call status.
func TestAdministrationEnvironmentPayloadStatus(t *testing.T) {
	for _, status := range []int32{0, DCGM_ST_BADPARAM, DCGM_ST_NOT_CONFIGURED} {
		t.Run(statusName(status), func(t *testing.T) {
			api := NewMockAdministrationNativeCalls(gomock.NewController(t))
			api.EXPECT().administrationSymbolAvailable("dcgmHostengineEnvironmentVariableInfo").Return(true)
			api.EXPECT().environmentVariableNative(gomock.Any()).DoAndReturn(func(response *nativeEnvironmentVariable) int32 {
				requireAdministrationSessionLocked(t)
				setEnvironmentVariableFixture(response, "", status)
				return 0
			})
			if status != 0 {
				api.EXPECT().administrationErrorText(status).Return("payload error")
			}
			value, err := hostengineEnvironmentVariableWithCalls(api, "DCGM_DIAG_REMOVED")
			if status == 0 {
				require.NoError(t, err)
				require.Empty(t, value)
			} else {
				var native *Error
				require.ErrorAs(t, err, &native)
				require.Equal(t, status, int32(native.Code))
			}
		})
	}
}

// TestAdministrationBeforeInit checks that unloaded and sessionless native operations are rejected.
func TestAdministrationBeforeInit(t *testing.T) {
	require.Nil(t, dcgmLibHandle)
	require.True(t, IsFeatureUnavailable(cgoAdapter{}.initializeDCGM()))
	require.True(t, IsFeatureUnavailable(cgoAdapter{}.startEmbeddedNative(EmbeddedOptions{})))
	require.True(t, IsFeatureUnavailable(cgoAdapter{}.shutdownDCGM()))
	checkAdministrationWithoutSession(t)
}

// TestAdministrationUnhealthyResponse preserves an unhealthy result from a successful native health call.
func TestAdministrationUnhealthyResponse(t *testing.T) {
	api := NewMockAdministrationNativeCalls(gomock.NewController(t))
	api.EXPECT().administrationSymbolAvailable("dcgmHostengineIsHealthy").Return(true)
	api.EXPECT().hostengineHealthNative(gomock.Any()).DoAndReturn(func(response *nativeHostengineHealth) int32 {
		requireAdministrationSessionLocked(t)
		response.overallHealth = 1
		return 0
	})
	healthy, err := hostengineIsHealthyWithCalls(api)
	require.NoError(t, err)
	require.False(t, healthy)
}

// TestAdministrationReturnedDataIsOwned checks that later native buffer changes cannot alter results.
func TestAdministrationReturnedDataIsOwned(t *testing.T) {
	t.Run("module statuses", func(t *testing.T) {
		api := NewMockAdministrationNativeCalls(gomock.NewController(t))
		api.EXPECT().administrationSymbolAvailable("dcgmModuleGetStatuses").Return(true)
		var native *nativeModuleStatuses
		api.EXPECT().moduleStatusesNative(gomock.Any()).DoAndReturn(func(response *nativeModuleStatuses) int32 {
			native = response
			setModuleStatusFixture(response, ModuleHealth, ModuleLoaded)
			return 0
		})
		got, err := moduleGetStatusesWithCalls(api)
		require.NoError(t, err)

		setModuleStatusFixture(native, ModuleDiag, ModuleFailed)
		require.Equal(t, []ModuleStatus{{ModuleHealth, ModuleLoaded}}, got)
	})

	t.Run("field group", func(t *testing.T) {
		api := NewMockAdministrationNativeCalls(gomock.NewController(t))
		api.EXPECT().administrationSymbolAvailable("dcgmFieldGroupGetInfo").Return(true)
		var native *nativeFieldGroupInfo
		api.EXPECT().fieldGroupInfoNative(gomock.Any()).DoAndReturn(func(response *nativeFieldGroupInfo) int32 {
			native = response
			response.numFieldIds = 1
			response.fieldIds[0] = 100
			response.fieldGroupName[0] = 'g'
			return 0
		})
		got, err := fieldGroupGetInfoWithCalls(api, FieldHandle{})
		require.NoError(t, err)

		native.fieldIds[0] = 200
		native.fieldGroupName[0] = 'h'
		require.Equal(t, []Short{100}, got.FieldIDs)
		require.Equal(t, "g", got.Name)
	})

	t.Run("group topology", func(t *testing.T) {
		api := NewMockAdministrationNativeCalls(gomock.NewController(t))
		api.EXPECT().administrationSymbolAvailable("dcgmGetGroupTopology").Return(true)
		var native *nativeGroupTopology
		api.EXPECT().groupTopologyNative(GroupHandle{}, gomock.Any()).DoAndReturn(func(_ GroupHandle, response *nativeGroupTopology) int32 {
			native = response
			response.groupCpuAffinityMask[0] = 3
			return 0
		})
		got, err := getGroupTopologyWithCalls(api, GroupHandle{})
		require.NoError(t, err)

		native.groupCpuAffinityMask[0] = 7
		require.Equal(t, uint64(3), got.CPUAffinityMask[0])
	})
}

// requireAdministrationSessionLocked checks that shutdown cannot unload the library during a call.
func requireAdministrationSessionLocked(t *testing.T) {
	t.Helper()
	if mux.TryLock() {
		mux.Unlock()
		t.Error("shutdown can acquire the session lock during administration")
	}
}

// requireAdministrationSessionUnlocked restores the mutex before reporting a leaked lock.
func requireAdministrationSessionUnlocked(t *testing.T) {
	t.Helper()
	unlocked := mux.TryLock()
	mux.Unlock()
	require.True(t, unlocked, "operation must release the session lock")
}
