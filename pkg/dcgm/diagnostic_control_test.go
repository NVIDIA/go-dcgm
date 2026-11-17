package dcgm

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestDiagnosticRequestNativeEncoding preserves explicit requests and applies native defaults for omitted options.
func TestDiagnosticRequestNativeEncoding(t *testing.T) {
	group := GroupAllGPUs()
	request := DiagnosticRequest{
		Level:           DiagMedium,
		TestNames:       []string{"memory", "nic_eud"},
		TestParameters:  []string{"memory.test_duration=1", "nic_eud.test_duration=2"},
		Group:           &group,
		Timeout:         5 * time.Minute,
		EnableHeartbeat: true,
	}

	decoded, selector, version, expected, err := diagnosticRequestForTest(request)
	require.NoError(t, err)
	require.Equal(t, expected, version)
	// DCGM requires empty logging and statistics paths for v10 requests.
	debugLogFile, statsPath, err := diagnosticRequestPathsForTest(request)
	require.NoError(t, err)
	require.Empty(t, debugLogFile)
	require.Empty(t, statsPath)
	require.Equal(t, request, decoded)
	require.Empty(t, selector)

	request.Group = nil
	request.Entities = []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: 7}, {EntityGroupId: FE_CONNECTX, EntityId: 17}}
	decoded, selector, _, _, err = diagnosticRequestForTest(request)
	require.NoError(t, err)
	require.Equal(t, "gpu:7,cx:17", selector)
	require.Equal(t, diagnosticNullGroupForTest(), decoded.Group.GetHandle())
	require.Equal(t, request.TestNames, decoded.TestNames)
	require.Equal(t, request.TestParameters, decoded.TestParameters)

	decoded, selector, _, _, err = diagnosticRequestForTest(DiagnosticRequest{Level: DiagQuick})
	require.NoError(t, err)
	require.Equal(t, GroupAllGPUs(), *decoded.Group)
	require.Empty(t, selector)
	require.Empty(t, decoded.TestNames)
	require.Empty(t, decoded.TestParameters)
	require.Zero(t, decoded.Timeout)
	require.False(t, decoded.EnableHeartbeat)
}

// TestDiagnosticEntitySelector matches native entity-group prefixes and retains each selected ID.
func TestDiagnosticEntitySelector(t *testing.T) {
	groups := []Field_Entity_Group{FE_GPU, FE_SWITCH, FE_VGPU, FE_GPU_I, FE_GPU_CI, FE_CPU, FE_CPU_CORE, FE_CONNECTX, FE_LINK}
	entities := make([]GroupEntityPair, len(groups))
	for i, group := range groups {
		entities[i] = GroupEntityPair{EntityGroupId: group, EntityId: uint(i)}
	}

	selector, err := diagnosticEntitySelector(entities)
	require.NoError(t, err)
	// EntityListHelpers.cpp dispatches these prefixes to the corresponding groups.
	require.Equal(t, "gpu:0,nvswitch:1,vgpu:2,instance:3,ci:4,cpu:5,core:6,cx:7,link:8", selector)
}

// TestDiagnosticRequestBounds rejects requests that exceed native limits before any API call.
func TestDiagnosticRequestBounds(t *testing.T) {
	maxNames, nameBytes, maxParameters, parameterBytes, maxEntities := diagnosticBoundsForTest()
	group := GroupAllGPUs()
	tests := []struct {
		name    string
		request DiagnosticRequest
	}{
		{"no tests", DiagnosticRequest{}},
		{"negative level", DiagnosticRequest{Level: -1}},
		{"large level", DiagnosticRequest{Level: DiagExtended + 1}},
		{"selector conflict", DiagnosticRequest{Level: DiagQuick, Group: &group, Entities: []GroupEntityPair{{EntityGroupId: FE_GPU}}}},
		{"too many tests", DiagnosticRequest{TestNames: make([]string, maxNames+1)}},
		{"too many parameters", DiagnosticRequest{Level: DiagQuick, TestParameters: make([]string, maxParameters+1)}},
		{"empty name", DiagnosticRequest{TestNames: []string{""}}},
		{"name NUL", DiagnosticRequest{TestNames: []string{"memory\x00other"}}},
		{"name too long", DiagnosticRequest{TestNames: []string{strings.Repeat("x", nameBytes)}}},
		{"name byte count", DiagnosticRequest{TestNames: []string{strings.Repeat("é", nameBytes/2+1)}}},
		{"parameter NUL", DiagnosticRequest{Level: DiagQuick, TestParameters: []string{"memory.test_duration=1\x00"}}},
		{"empty parameter", DiagnosticRequest{Level: DiagQuick, TestParameters: []string{""}}},
		{"parameter too long", DiagnosticRequest{Level: DiagQuick, TestParameters: []string{strings.Repeat("x", parameterBytes)}}},
		{"negative timeout", DiagnosticRequest{Level: DiagQuick, Timeout: -time.Second}},
		{"fractional timeout", DiagnosticRequest{Level: DiagQuick, Timeout: time.Second + time.Nanosecond}},
		{"timeout overflow", DiagnosticRequest{Level: DiagQuick, Timeout: time.Duration(math.MaxInt64)}},
		{"millisecond conversion overflow", DiagnosticRequest{Level: DiagQuick, Timeout: (time.Duration(math.MaxUint32/1000) + 1) * time.Second}},
		{"unknown entity group", DiagnosticRequest{Level: DiagQuick, Entities: []GroupEntityPair{{EntityGroupId: FE_NONE}}}},
		{"too many entities", DiagnosticRequest{Level: DiagQuick, Entities: make([]GroupEntityPair, maxEntities+1)}},
		{"selector too long", DiagnosticRequest{Level: DiagQuick, Entities: make([]GroupEntityPair, maxEntities)}},
	}
	for i := range tests[len(tests)-1].request.Entities {
		tests[len(tests)-1].request.Entities[i] = GroupEntityPair{EntityGroupId: FE_CONNECTX, EntityId: math.MaxUint32}
	}
	if uint64(^uint(0)) > math.MaxUint32 {
		tooLarge := uint(math.MaxUint32)
		tooLarge++
		tests = append(tests, struct {
			name    string
			request DiagnosticRequest
		}{"entity ID overflow", DiagnosticRequest{Level: DiagQuick, Entities: []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: tooLarge}}}})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewMockdiagnosticCalls(gomock.NewController(t))
			_, err := actionValidateWithCalls(api, tt.request)
			require.Error(t, err)
			_, _, _, _, err = diagnosticRequestForTest(tt.request)
			require.Error(t, err)
		})
	}
}

// TestDiagnosticRequestAtNativeLimits accepts the largest representable names, parameters and timeout.
func TestDiagnosticRequestAtNativeLimits(t *testing.T) {
	maxNames, nameBytes, maxParameters, parameterBytes, _ := diagnosticBoundsForTest()
	valid := DiagnosticRequest{
		TestNames:      make([]string, maxNames),
		TestParameters: make([]string, maxParameters),
		Timeout:        time.Duration(math.MaxUint32/1000) * time.Second,
	}
	for i := range valid.TestNames {
		valid.TestNames[i] = strings.Repeat("x", nameBytes-1)
	}
	for i := range valid.TestParameters {
		valid.TestParameters[i] = strings.Repeat("x", parameterBytes-1)
	}
	decoded, _, _, _, err := diagnosticRequestForTest(valid)
	require.NoError(t, err)
	require.Equal(t, valid.Timeout, decoded.Timeout)
	require.Equal(t, valid.TestNames, decoded.TestNames)
	require.Equal(t, valid.TestParameters, decoded.TestParameters)
}

// TestActionValidateNativeResults preserves decoded results and typed native failures.
func TestActionValidateNativeResults(t *testing.T) {
	for _, code := range []int{DCGM_ST_OK, DCGM_ST_BADPARAM, DCGM_ST_DIAG_STOPPED, DCGM_ST_FUNCTION_NOT_FOUND} {
		api := NewMockdiagnosticCalls(gomock.NewController(t))
		request := DiagnosticRequest{TestNames: []string{"nic_eud"}}
		response := createConnectXDiagResponseForTest()
		var nativeErr error
		if code != DCGM_ST_OK {
			nativeErr = diagnosticNativeErrorForTest(code)
		}
		api.EXPECT().requireDiagnosticSymbol("dcgmActionValidate_v2").Return(nil)
		api.EXPECT().actionValidateNative(request).Return(response, nativeErr)
		got, err := actionValidateWithCalls(api, request)
		require.Equal(t, response, got)
		if code == DCGM_ST_OK {
			require.NoError(t, err)
			continue
		}
		var native *Error
		require.ErrorAs(t, err, &native)
		require.Equal(t, code, int(native.Code))
		require.Equal(t, code == DCGM_ST_FUNCTION_NOT_FOUND, IsFeatureUnavailable(err))
	}
}

// TestDiagnosticControlsAndAvailability keeps control errors distinct from unavailable functions.
func TestDiagnosticControlsAndAvailability(t *testing.T) {
	for _, symbol := range []string{"dcgmStopDiagnostic", "dcgmDiagSendHeartbeat", "dcgmProfPause", "dcgmProfResume"} {
		t.Run(symbol, func(t *testing.T) {
			api := NewMockdiagnosticCalls(gomock.NewController(t))
			api.EXPECT().requireDiagnosticSymbol(symbol).Return(nil).Times(2)
			api.EXPECT().diagnosticControlNative(symbol).Return(nil)
			nativeErr := diagnosticNativeErrorForTest(DCGM_ST_BADPARAM)
			api.EXPECT().diagnosticControlNative(symbol).Return(nativeErr)
			require.NoError(t, diagnosticControlWithCalls(api, symbol))
			require.ErrorIs(t, diagnosticControlWithCalls(api, symbol), nativeErr)
			missing := &FeatureUnavailableError{Feature: diagnosticFeature, Symbol: symbol}
			api.EXPECT().requireDiagnosticSymbol(symbol).Return(missing)
			require.ErrorIs(t, diagnosticControlWithCalls(api, symbol), missing)
		})
	}
	api := NewMockdiagnosticCalls(gomock.NewController(t))
	api.EXPECT().requireDiagnosticSymbol("dcgmActionValidate_v2").Return(&FeatureUnavailableError{Feature: diagnosticFeature, Symbol: "dcgmActionValidate_v2"})
	_, err := actionValidateWithCalls(api, DiagnosticRequest{Level: DiagQuick})
	require.True(t, IsFeatureUnavailable(err))
}

// TestDiagnosticAPIsWithoutInit rejects diagnostic operations without an active session.
func TestDiagnosticAPIsWithoutInit(t *testing.T) {
	for _, call := range []func() error{
		StopDiag, DiagSendHeartbeat, ProfPause, ProfResume,
		func() error { _, err := ActionValidate(DiagnosticRequest{Level: DiagQuick}); return err },
		func() error { _, err := RunDiag(DiagQuick, GroupAllGPUs()); return err },
	} {
		require.True(t, IsFeatureUnavailable(call()))
	}
}

// TestDiagnosticRunAllowsConcurrentStop proves StopDiag can interrupt a blocked diagnostic call.
func TestDiagnosticRunAllowsConcurrentStop(t *testing.T) {
	api := NewMockdiagnosticCalls(gomock.NewController(t))

	// The diagnostic holds until a separate stop call reaches the native boundary.
	started := make(chan struct{})
	stopped := make(chan struct{})
	api.EXPECT().requireDiagnosticSymbol("dcgmRunDiagnostic").Return(nil)
	api.EXPECT().runDiagNative(DiagLong, GroupAllGPUs()).DoAndReturn(func(DiagType, GroupHandle) (DiagResults, error) {
		close(started)
		select {
		case <-stopped:
			return DiagResults{}, diagnosticNativeErrorForTest(DCGM_ST_DIAG_STOPPED)
		case <-time.After(time.Second):
			return DiagResults{}, errors.New("stop blocked behind run")
		}
	})

	api.EXPECT().requireDiagnosticSymbol("dcgmStopDiagnostic").Return(nil)
	api.EXPECT().diagnosticControlNative("dcgmStopDiagnostic").DoAndReturn(func(string) error {
		close(stopped)
		return nil
	})
	done := make(chan error, 1)
	go func() {
		_, err := runDiagWithCalls(api, DiagLong, GroupAllGPUs())
		done <- err
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("run did not start")
	}
	// Stop must complete while the original native call is still in flight.
	require.NoError(t, diagnosticControlWithCalls(api, "dcgmStopDiagnostic"))

	select {
	case err := <-done:
		var native *Error
		require.ErrorAs(t, err, &native)
		require.Equal(t, DCGM_ST_DIAG_STOPPED, int(native.Code))
	case <-time.After(time.Second):
		t.Fatal("run did not stop")
	}
}

// TestDiagnosticConnectXResults retains network-device identity in structured diagnostic results.
func TestDiagnosticConnectXResults(t *testing.T) {
	result := createConnectXDiagResponseForTest()
	require.Equal(t, FE_CONNECTX, result.Entities[0].Entity.EntityGroupId)
	require.Equal(t, "nic_eud", result.Tests[0].Name)
	require.Equal(t, FE_CONNECTX, result.Tests[0].Results[0].Entity.EntityGroupId)
	require.Equal(t, uint(17), result.Tests[0].Results[0].Entity.EntityId)
	require.Equal(t, FE_CONNECTX, result.Tests[0].Errors[0].Entity.EntityGroupId)
	require.Equal(t, uint(42), result.Tests[0].Errors[0].Code)
	require.Equal(t, uint(1), result.Tests[0].AuxDataVersion)
	require.Equal(t, `{"version":"1"}`, result.Tests[0].AuxData)
	for _, legacy := range result.Software {
		require.NotEqual(t, uint(17), legacy.EntityID)
	}
}
