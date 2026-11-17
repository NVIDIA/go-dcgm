package dcgm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func configPointerForTest[T any](value T) *T { return &value }

func configStatusHandleForTest() configStatusHandle {
	return configStatusHandle{handle: 0xC0F19}
}

func checkConfigSessionLocked(t *testing.T, stage string) {
	t.Helper()
	if mux.TryLock() {
		mux.Unlock()
		t.Errorf("shutdown can acquire the session lock during %s", stage)
	}
}

func TestConfigSessionLifetime(t *testing.T) {
	for _, operation := range []string{"get", "set", "enforce", "reset", "profile set", "profile info", "profile name"} {
		t.Run(operation, func(t *testing.T) {
			api := NewMockConfigCalls(gomock.NewController(t))
			status := configStatusHandleForTest()
			if operation == "get" || operation == "set" || operation == "enforce" || operation == "reset" {
				api.EXPECT().configStatusCreate().DoAndReturn(func() (configStatusHandle, error) {
					checkConfigSessionLocked(t, "status creation")
					return status, nil
				})
				api.EXPECT().configStatusClear(status).DoAndReturn(func(configStatusHandle) error {
					checkConfigSessionLocked(t, "status clear")
					return nil
				})
				api.EXPECT().configStatusCount(status).DoAndReturn(func(configStatusHandle) (uint, error) {
					checkConfigSessionLocked(t, "status count")
					return 1, nil
				})
				api.EXPECT().configStatusPop(status).DoAndReturn(func(configStatusHandle) (ConfigFieldError, error) {
					checkConfigSessionLocked(t, "status pop")
					return ConfigFieldError{GPU: 7, Err: configErrorForTest(DCGM_ST_BADPARAM)}, nil
				})
				api.EXPECT().configStatusDestroy(status).DoAndReturn(func(configStatusHandle) error {
					checkConfigSessionLocked(t, "status destruction")
					return nil
				})
			}
			var err error
			switch operation {
			case "get":
				api.EXPECT().configGetNative(GroupHandle{}, ConfigCurrentState, gomock.Any(), status).DoAndReturn(func(GroupHandle, ConfigType, []nativeGPUConfig, configStatusHandle) error {
					checkConfigSessionLocked(t, "native get")
					return nil
				})
				_, err = configGetWithCalls(api, GroupHandle{}, ConfigCurrentState)
			case "set":
				api.EXPECT().configSetNative(GroupHandle{}, gomock.Any(), status).DoAndReturn(func(GroupHandle, *nativeGPUConfig, configStatusHandle) error {
					checkConfigSessionLocked(t, "native set")
					return nil
				})
				err = configSetWithCalls(api, GroupHandle{}, GPUConfig{})
			case "enforce":
				api.EXPECT().configEnforceNative(GroupHandle{}, status).DoAndReturn(func(GroupHandle, configStatusHandle) error {
					checkConfigSessionLocked(t, "native enforce")
					return nil
				})
				err = configEnforceWithCalls(api, GroupHandle{})
			case "reset":
				api.EXPECT().configResetNative(GroupHandle{}, status).DoAndReturn(func(GroupHandle, configStatusHandle) error {
					checkConfigSessionLocked(t, "native reset")
					return nil
				})
				err = configResetWithCalls(api, GroupHandle{})
			case "profile set":
				api.EXPECT().configSetWorkloadPowerProfileNative(gomock.Any()).DoAndReturn(func(*nativeWorkloadPowerProfile) error {
					checkConfigSessionLocked(t, "native profile set")
					return nil
				})
				err = configSetWorkloadPowerProfileWithCalls(api, GroupHandle{}, WorkloadProfileSet, PowerProfileMask{1})
			case "profile info":
				api.EXPECT().workloadPowerProfileInfoNative(uint(7), gomock.Any(), gomock.Any()).DoAndReturn(func(uint, *nativeWorkloadProfilesInfo, *nativeWorkloadProfilesStatus) error {
					checkConfigSessionLocked(t, "native profile info")
					return nil
				})
				_, err = workloadPowerProfileInfoWithCalls(api, 7)
			case "profile name":
				api.EXPECT().powerProfileNameNative(PowerProfileCompute).DoAndReturn(func(PowerProfileID) (string, error) {
					checkConfigSessionLocked(t, "native profile name")
					return "Compute", nil
				})
				_, err = powerProfileNameWithCalls(api, PowerProfileCompute)
			}
			if operation == "get" || operation == "set" || operation == "enforce" || operation == "reset" {
				var configErr *ConfigError
				require.ErrorAs(t, err, &configErr)
				require.Len(t, configErr.Fields, 1)
			} else {
				require.NoError(t, err)
			}
			require.True(t, mux.TryLock(), "operation must release the session lock")
			mux.Unlock()
		})
	}
}

func TestConfigConversion(t *testing.T) {
	source := nativeConfigFixtureForTest()
	got := configFromNative(&source)
	require.Equal(t, uint(7), got.GPU)
	require.Equal(t, uint32(1), *got.ECCMode)
	require.Equal(t, uint32(2), *got.ComputeMode)
	require.Zero(t, *got.SyncBoost)
	require.Equal(t, int64(1200000000), *got.MemoryClockHz)
	require.Equal(t, int64(900000000), *got.SMClockHz)
	require.Equal(t, &ConfigPowerLimit{Type: ConfigPowerBudgetGroup, Watts: 250}, got.PowerLimit)
	for i, word := range got.WorkloadPowerProfiles {
		require.Equal(t, uint32(i+1), *word)
	}
	source.powerLimit.val = 99
	require.Equal(t, uint32(250), got.PowerLimit.Watts)
	roundTrip, err := configToNative(got)
	require.NoError(t, err)
	require.Equal(t, uint32(1), uint32(roundTrip.eccMode))
	require.Equal(t, uint32(2), uint32(roundTrip.computeMode))
	require.Zero(t, uint32(roundTrip.perfState.syncBoost))
	version, clockVersion := configVersionsForTest()
	require.Equal(t, version, uint32(roundTrip.version))
	require.Equal(t, clockVersion, int(roundTrip.perfState.targetClocks.version))
	require.Equal(t, int64(1200000000), int64(roundTrip.perfState.targetClocks.memClock))
	require.Equal(t, int64(900000000), int64(roundTrip.perfState.targetClocks.smClock))
	require.Equal(t, int(ConfigPowerBudgetGroup), int(roundTrip.powerLimit._type))
	require.Equal(t, uint32(250), uint32(roundTrip.powerLimit.val))
	for i, word := range roundTrip.workloadPowerProfiles {
		require.Equal(t, uint32(i+1), uint32(word))
	}
	blank, err := configToNative(GPUConfig{})
	require.NoError(t, err)
	blank32, blank64 := configNativeBlanksForTest()
	require.Equal(t, blank32, uint32(blank.eccMode))
	require.Equal(t, blank32, uint32(blank.computeMode))
	require.Equal(t, blank32, uint32(blank.perfState.syncBoost))
	require.Equal(t, blank64, int64(blank.perfState.targetClocks.memClock))
	require.Equal(t, blank64, int64(blank.perfState.targetClocks.smClock))
	require.Equal(t, blank32, uint32(blank.powerLimit.val))
	for _, word := range blank.workloadPowerProfiles {
		require.Equal(t, blank32, uint32(word))
	}
	require.Equal(t, GPUConfig{}, configFromNative(&blank))
	require.Equal(t, GPUConfig{GPU: 5}, configNativeUnavailableForTest())
}

func TestConfigValidation(t *testing.T) {
	blank32, blank64 := configNativeBlanksForTest()
	for _, cfg := range []GPUConfig{
		{ECCMode: configPointerForTest(uint32(2))},
		{ComputeMode: configPointerForTest(uint32(3))},
		{SyncBoost: configPointerForTest(uint32(2))},
		{MemoryClockHz: configPointerForTest(int64(-1))},
		{SMClockHz: configPointerForTest(int64(1000001))},
		{SMClockHz: configPointerForTest(int64(blank32) * 1000000)},
		{MemoryClockHz: configPointerForTest(blank64)},
		{PowerLimit: &ConfigPowerLimit{Type: ConfigPowerLimitType(99), Watts: 200}},
		{PowerLimit: &ConfigPowerLimit{Watts: blank32}},
		{WorkloadPowerProfiles: [PowerProfileMaskSize]*uint32{configPointerForTest(blank32)}},
	} {
		api := NewMockConfigCalls(gomock.NewController(t))
		require.Error(t, configSetWithCalls(api, GroupHandle{}, cfg))
	}
	api := NewMockConfigCalls(gomock.NewController(t))
	_, err := configGetWithCalls(api, GroupHandle{}, ConfigType(99))
	require.Error(t, err)
	require.Error(t, configSetWorkloadPowerProfileWithCalls(api, GroupHandle{}, WorkloadProfileAction(99), PowerProfileMask{}))
	for _, hz := range []int64{0, 1000000, (int64(blank32) - 1) * 1000000} {
		_, err := configToNative(GPUConfig{SMClockHz: &hz})
		require.NoError(t, err)
	}
}

type configStatusScenario struct {
	name   string
	fail   string
	perGPU bool
}

func TestConfigStatusLifecycle(t *testing.T) {
	for _, tc := range []configStatusScenario{
		{name: "success"},
		{name: "create", fail: "create"},
		{name: "clear", fail: "clear"},
		{name: "operation", fail: "operation"},
		{name: "count", fail: "count"},
		{name: "pop", fail: "pop", perGPU: true},
		{name: "destroy", fail: "destroy"},
		{name: "OK with GPU errors", perGPU: true},
		{name: "partial failure with GPU errors", fail: "operation", perGPU: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api := NewMockConfigCalls(gomock.NewController(t))
			nativeErr := configErrorForTest(DCGM_ST_BADPARAM)
			expectConfigStatusLifecycle(api, tc, nativeErr)
			err := configEnforceWithCalls(api, GroupHandle{})
			requireConfigStatusResult(t, tc, nativeErr, err)
		})
	}
}

func expectConfigStatusLifecycle(api *MockConfigCalls, tc configStatusScenario, nativeErr error) {
	fail := func(stage string) error {
		if tc.fail == stage {
			return nativeErr
		}
		return nil
	}
	status := configStatusHandleForTest()
	create := api.EXPECT().configStatusCreate().Return(status, fail("create"))
	if tc.fail == "create" {
		return
	}
	clearCall := api.EXPECT().configStatusClear(status).Return(fail("clear")).After(create)
	last := clearCall
	if tc.fail != "clear" {
		last = expectConfigOperationStatus(api, tc, clearCall, fail)
	}
	api.EXPECT().configStatusDestroy(status).Return(fail("destroy")).After(last)
}

func expectConfigOperationStatus(api *MockConfigCalls, tc configStatusScenario, clearCall *gomock.Call, fail func(string) error) *gomock.Call {
	status := configStatusHandleForTest()
	op := api.EXPECT().configEnforceNative(GroupHandle{}, status).Return(fail("operation")).After(clearCall)
	count := uint(0)
	if tc.perGPU {
		count = 2
	}
	countCall := api.EXPECT().configStatusCount(status).Return(count, fail("count")).After(op)
	if !tc.perGPU || tc.fail == "count" {
		return countCall
	}
	pop := api.EXPECT().configStatusPop(status).Return(configFieldErrorForTest(7, DCGM_FI_DEV_POWER_MGMT_LIMIT, DCGM_ST_NOT_SUPPORTED), nil).After(countCall)
	return api.EXPECT().configStatusPop(status).Return(configFieldErrorForTest(8, DCGM_FI_DEV_ECC_CURRENT, DCGM_ST_REQUIRES_ROOT), fail("pop")).After(pop)
}

func requireConfigStatusResult(t *testing.T, tc configStatusScenario, nativeErr, err error) {
	t.Helper()
	if tc.fail == "" && !tc.perGPU {
		require.NoError(t, err)
		return
	}
	require.Error(t, err)
	if tc.fail != "" {
		require.ErrorIs(t, err, nativeErr)
	}
	if !tc.perGPU {
		return
	}
	var cfgErr *ConfigError
	require.ErrorAs(t, err, &cfgErr)
	require.Equal(t, uint(7), cfgErr.Fields[0].GPU)
	require.Equal(t, DCGM_FI_DEV_POWER_MGMT_LIMIT, cfgErr.Fields[0].FieldID)
	var typed *Error
	require.ErrorAs(t, cfgErr.Fields[0], &typed)
	require.Equal(t, DCGM_ST_NOT_SUPPORTED, int(typed.Code))
	require.Contains(t, err.Error(), "GPU 7 field")
	if tc.fail != "pop" {
		require.Len(t, cfgErr.Fields, 2)
		require.Equal(t, uint(8), cfgErr.Fields[1].GPU)
	} else {
		require.Len(t, cfgErr.Fields, 1)
	}
}

func expectConfigStatusForTest(api *MockConfigCalls) {
	api.EXPECT().configStatusCreate().Return(configStatusHandleForTest(), nil)
	api.EXPECT().configStatusClear(configStatusHandleForTest())
	api.EXPECT().configStatusCount(configStatusHandleForTest()).Return(uint(0), nil)
	api.EXPECT().configStatusDestroy(configStatusHandleForTest())
}

func TestConfigNativeRequests(t *testing.T) {
	api := NewMockConfigCalls(gomock.NewController(t))
	expectConfigStatusForTest(api)
	api.EXPECT().configSetNative(GroupHandle{}, gomock.Any(), configStatusHandleForTest()).DoAndReturn(func(_ GroupHandle, r *nativeGPUConfig, _ configStatusHandle) error {
		version, clockVersion := configVersionsForTest()
		blank32, blank64 := configNativeBlanksForTest()
		require.Equal(t, version, uint32(r.version))
		require.Equal(t, clockVersion, int(r.perfState.targetClocks.version))
		require.Equal(t, uint32(225), uint32(r.powerLimit.val))
		require.Equal(t, int(ConfigPowerCapIndividual), int(r.powerLimit._type))
		require.Equal(t, blank32, uint32(r.eccMode))
		require.Equal(t, blank32, uint32(r.computeMode))
		require.Equal(t, blank32, uint32(r.perfState.syncBoost))
		require.Equal(t, blank64, int64(r.perfState.targetClocks.memClock))
		require.Equal(t, blank64, int64(r.perfState.targetClocks.smClock))
		for _, word := range r.workloadPowerProfiles {
			require.Equal(t, blank32, uint32(word))
		}
		return nil
	})
	require.NoError(t, configSetWithCalls(api, GroupHandle{}, GPUConfig{PowerLimit: &ConfigPowerLimit{Watts: 225}}))
	for _, kind := range []ConfigType{ConfigCurrentState, ConfigTargetState} {
		expectConfigStatusForTest(api)
		api.EXPECT().configGetNative(GroupHandle{}, kind, gomock.Any(), configStatusHandleForTest()).DoAndReturn(func(_ GroupHandle, _ ConfigType, r []nativeGPUConfig, _ configStatusHandle) error {
			require.Len(t, r, int(MAX_NUM_DEVICES))
			version, _ := configVersionsForTest()
			blank32, _ := configNativeBlanksForTest()
			for i := range r {
				require.Equal(t, version, uint32(r[i].version))
				require.Equal(t, blank32, uint32(r[i].gpuId))
			}
			r[0] = nativeConfigFixtureForTest()
			return nil
		})
		got, err := configGetWithCalls(api, GroupHandle{}, kind)
		require.NoError(t, err)
		require.Len(t, got, 1)
		require.Equal(t, uint(7), got[0].GPU)
	}
	expectConfigStatusForTest(api)
	api.EXPECT().configResetNative(GroupHandle{}, configStatusHandleForTest()).Return(configErrorForTest(DCGM_ST_NOT_SUPPORTED))
	require.Error(t, configResetWithCalls(api, GroupHandle{}))
}

func TestConfigGetPartialResults(t *testing.T) {
	api := NewMockConfigCalls(gomock.NewController(t))
	expectConfigStatusForTest(api)
	want := configErrorForTest(DCGM_ST_GENERIC_ERROR)
	api.EXPECT().configGetNative(GroupHandle{}, ConfigCurrentState, gomock.Any(), configStatusHandleForTest()).DoAndReturn(func(_ GroupHandle, _ ConfigType, r []nativeGPUConfig, _ configStatusHandle) error {
		r[0] = nativeConfigFixtureForTest()
		return want
	})
	got, err := configGetWithCalls(api, GroupHandle{}, ConfigCurrentState)
	require.ErrorIs(t, err, want)
	require.Len(t, got, 1)
	require.Equal(t, uint(7), got[0].GPU)
}

func TestConfigStatusCountBounds(t *testing.T) {
	api := NewMockConfigCalls(gomock.NewController(t))
	api.EXPECT().configStatusCreate().Return(configStatusHandleForTest(), nil)
	api.EXPECT().configStatusClear(configStatusHandleForTest())
	api.EXPECT().configStatusCount(configStatusHandleForTest()).Return(MAX_NUM_DEVICES+1, nil)
	api.EXPECT().configStatusDestroy(configStatusHandleForTest())
	require.ErrorContains(t, withConfigStatus(api, func(status configStatusHandle) error {
		require.Equal(t, configStatusHandleForTest(), status)
		return nil
	}), "exceeds")
}

func TestWorkloadProfileNativeRequests(t *testing.T) {
	api := NewMockConfigCalls(gomock.NewController(t))
	group := GroupHandle{}
	group.SetHandle(123)
	for _, action := range []WorkloadProfileAction{WorkloadProfileSet, WorkloadProfileOverwrite, WorkloadProfileClear} {
		mask := PowerProfileMask{1, 2, 3, 4, 5, 6, 7, 8}
		api.EXPECT().configSetWorkloadPowerProfileNative(gomock.Any()).DoAndReturn(func(r *nativeWorkloadPowerProfile) error {
			require.Equal(t, workloadRequestVersionForTest(), uint32(r.version))
			require.Equal(t, uintptr(123), uintptr(r.groupId))
			require.Equal(t, int(action), int(r.action))
			for i, word := range r.profileMask {
				require.Equal(t, mask[i], uint32(word))
			}
			return nil
		})
		require.NoError(t, configSetWorkloadPowerProfileWithCalls(api, group, action, mask))
	}
	api.EXPECT().workloadPowerProfileInfoNative(uint(7), gomock.Any(), gomock.Any()).DoAndReturn(func(_ uint, p *nativeWorkloadProfilesInfo, s *nativeWorkloadProfilesStatus) error {
		pv, iv, sv := profileVersionsForTest()
		require.Equal(t, pv, uint32(p.version))
		require.Equal(t, sv, uint32(s.version))
		for i := range p.workloadPowerProfile {
			require.Equal(t, iv, uint32(p.workloadPowerProfile[i].version))
		}
		populateProfilesForTest(p, s, 2)
		return nil
	})
	got, err := workloadPowerProfileInfoWithCalls(api, 7)
	require.NoError(t, err)
	require.Len(t, got.Profiles, 2)
	require.Equal(t, PowerProfileID(0), got.Profiles[0].ID)
	require.Equal(t, PowerProfileID(1), got.Profiles[1].ID)
	require.Equal(t, uint32(10), got.Profiles[0].Priority)
	require.Equal(t, uint32(11), got.Profiles[1].Priority)
	require.Equal(t, PowerProfileMask{20, 21, 22, 23, 24, 25, 26, 27}, got.Profiles[0].ConflictingMask)
	require.Equal(t, PowerProfileMask{21, 22, 23, 24, 25, 26, 27, 28}, got.Profiles[1].ConflictingMask)
	require.Equal(t, PowerProfileMask{30, 31, 32, 33, 34, 35, 36, 37}, got.Supported)
	require.Equal(t, PowerProfileMask{40, 41, 42, 43, 44, 45, 46, 47}, got.Requested)
	require.Equal(t, PowerProfileMask{50, 51, 52, 53, 54, 55, 56, 57}, got.Enforced)
	api.EXPECT().powerProfileNameNative(PowerProfileCompute).Return("Compute", nil)
	name, err := powerProfileNameWithCalls(api, PowerProfileCompute)
	require.NoError(t, err)
	require.Equal(t, "Compute", name)
}

func TestWorkloadProfileErrorsAndBounds(t *testing.T) {
	api := NewMockConfigCalls(gomock.NewController(t))
	want := configErrorForTest(DCGM_ST_NOT_SUPPORTED)
	api.EXPECT().configSetWorkloadPowerProfileNative(gomock.Any()).Return(want)
	require.ErrorIs(t, configSetWorkloadPowerProfileWithCalls(api, GroupHandle{}, WorkloadProfileSet, PowerProfileMask{}), want)
	api.EXPECT().workloadPowerProfileInfoNative(uint(0), gomock.Any(), gomock.Any()).Return(want)
	_, err := workloadPowerProfileInfoWithCalls(api, 0)
	require.ErrorIs(t, err, want)
	api.EXPECT().powerProfileNameNative(PowerProfileID(99)).Return("", want)
	_, err = powerProfileNameWithCalls(api, 99)
	require.ErrorIs(t, err, want)
	for _, count := range []int{0, 255, 256} {
		var p nativeWorkloadProfilesInfo
		var s nativeWorkloadProfilesStatus
		populateProfilesForTest(&p, &s, count)
		got, convErr := workloadProfileInfoFromNative(&p, &s)
		if count == 256 {
			require.Error(t, convErr)
		} else {
			require.NoError(t, convErr)
			require.Len(t, got.Profiles, count)
		}
	}
}

func TestConfigAPIsWithoutLibrary(t *testing.T) {
	previous := dcgmLibHandle
	dcgmLibHandle = nil
	t.Cleanup(func() { dcgmLibHandle = previous })
	calls := []struct {
		name string
		call func() error
	}{
		{"get", func() error { _, err := ConfigGet(GroupHandle{}, ConfigCurrentState); return err }},
		{"set", func() error { return ConfigSet(GroupHandle{}, GPUConfig{PowerLimit: &ConfigPowerLimit{Watts: 200}}) }},
		{"enforce", func() error { return ConfigEnforce(GroupHandle{}) }},
		{"reset", func() error { return ConfigReset(GroupHandle{}) }},
		{"profile set", func() error {
			return ConfigSetWorkloadPowerProfile(GroupHandle{}, WorkloadProfileSet, PowerProfileMask{1})
		}},
		{"profile info", func() error { _, err := GetDeviceWorkloadPowerProfileInfo(0); return err }},
		{"profile name", func() error { _, err := PowerProfileName(PowerProfileCompute); return err }},
	}
	for _, tc := range calls {
		t.Run(tc.name, func(t *testing.T) { require.True(t, IsFeatureUnavailable(tc.call())) })
	}
}

func TestConfigNativeAPIsWithoutLibrary(t *testing.T) {
	previous := dcgmLibHandle
	dcgmLibHandle = nil
	t.Cleanup(func() { dcgmLibHandle = previous })
	api := cgoAdapter{}
	config := blankNativeConfig()
	var request nativeWorkloadPowerProfile
	var profiles nativeWorkloadProfilesInfo
	var profilesStatus nativeWorkloadProfilesStatus
	status := configStatusHandleForTest()
	calls := []struct {
		symbol string
		call   func() error
	}{
		{"dcgmConfigGet", func() error {
			return api.configGetNative(GroupHandle{}, ConfigCurrentState, []nativeGPUConfig{config}, status)
		}},
		{"dcgmConfigSet", func() error { return api.configSetNative(GroupHandle{}, &config, status) }},
		{"dcgmConfigEnforce", func() error { return api.configEnforceNative(GroupHandle{}, status) }},
		{"dcgmConfigReset", func() error { return api.configResetNative(GroupHandle{}, status) }},
		{"dcgmConfigSetWorkloadPowerProfile", func() error { return api.configSetWorkloadPowerProfileNative(&request) }},
		{"dcgmGetDeviceWorkloadPowerProfileInfo", func() error { return api.workloadPowerProfileInfoNative(0, &profiles, &profilesStatus) }},
		{"dcgmPowerProfileIdToName", func() error { _, err := api.powerProfileNameNative(PowerProfileMaxP); return err }},
		{"dcgmStatusCreate", func() error { _, err := api.configStatusCreate(); return err }},
		{"dcgmStatusClear", func() error { return api.configStatusClear(status) }},
		{"dcgmStatusGetCount", func() error { _, err := api.configStatusCount(status); return err }},
		{"dcgmStatusPopError", func() error { _, err := api.configStatusPop(status); return err }},
		{"dcgmStatusDestroy", func() error { return api.configStatusDestroy(status) }},
	}
	for _, tc := range calls {
		t.Run(tc.symbol, func(t *testing.T) {
			err := tc.call()
			require.True(t, IsFeatureUnavailable(err))
			var unavailable *FeatureUnavailableError
			require.ErrorAs(t, err, &unavailable)
			require.Equal(t, tc.symbol, unavailable.Symbol)
		})
	}
}

func TestConfigSetAndResetNativeResults(t *testing.T) {
	for _, operation := range []string{"set", "reset"} {
		for _, failed := range []bool{false, true} {
			t.Run(operation+map[bool]string{false: " success", true: " error"}[failed], func(t *testing.T) {
				checkConfigSetAndResetNativeResult(t, operation, failed)
			})
		}
	}
}

func checkConfigSetAndResetNativeResult(t *testing.T, operation string, failed bool) {
	t.Helper()
	api := NewMockConfigCalls(gomock.NewController(t))
	expectConfigStatusForTest(api)
	var want error
	if failed {
		want = configErrorForTest(DCGM_ST_BADPARAM)
	}
	var err error
	if operation == "set" {
		api.EXPECT().configSetNative(GroupHandle{}, gomock.Any(), configStatusHandleForTest()).Return(want)
		err = configSetWithCalls(api, GroupHandle{}, GPUConfig{})
	} else {
		api.EXPECT().configResetNative(GroupHandle{}, configStatusHandleForTest()).Return(want)
		err = configResetWithCalls(api, GroupHandle{})
	}
	if want == nil {
		require.NoError(t, err)
	} else {
		require.ErrorIs(t, err, want)
	}
}

func TestConfigStatusRetainsOperationAndCleanupErrors(t *testing.T) {
	api := NewMockConfigCalls(gomock.NewController(t))
	operationErr := configErrorForTest(DCGM_ST_BADPARAM)
	cleanupErr := configErrorForTest(DCGM_ST_UNINITIALIZED)
	api.EXPECT().configStatusCreate().Return(configStatusHandleForTest(), nil)
	api.EXPECT().configStatusClear(configStatusHandleForTest())
	api.EXPECT().configStatusCount(configStatusHandleForTest()).Return(uint(1), nil)
	api.EXPECT().configStatusPop(configStatusHandleForTest()).Return(configFieldErrorForTest(9, DCGM_FI_DEV_POWER_MGMT_LIMIT, DCGM_ST_NOT_SUPPORTED), nil)
	api.EXPECT().configStatusDestroy(configStatusHandleForTest()).Return(cleanupErr)
	err := withConfigStatus(api, func(status configStatusHandle) error {
		require.Equal(t, configStatusHandleForTest(), status)
		return operationErr
	})
	require.ErrorIs(t, err, operationErr)
	require.ErrorIs(t, err, cleanupErr)
	var fields *ConfigError
	require.ErrorAs(t, err, &fields)
	require.Len(t, fields.Fields, 1)
	require.Equal(t, uint(9), fields.Fields[0].GPU)
}
