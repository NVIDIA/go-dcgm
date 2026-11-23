package dcgm

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// buildInfoForTest formats a version in the build string DCGM returns.
func buildInfoForTest(version string) VersionInfo {
	return VersionInfo{RawBuildInfoString: "version:" + version + ";arch:x86_64"}
}

// TestRuntimeVersionSelection checks version parsing and which clients need the v2 config layout.
func TestRuntimeVersionSelection(t *testing.T) {
	for _, raw := range []string{
		"",
		"version:",
		"version:4.8",
		"version:4.8.0.1",
		"version:4.08.0",
		"version:4.-8.0",
		"version:4.8.0-rc1",
		"version:4.8.0;version:4.7.1",
		"version:;version:4.8.0",
		"build:4.8.0",
		"version:4294967296.8.0",
	} {
		t.Run(raw, func(t *testing.T) {
			info := VersionInfo{RawBuildInfoString: raw}
			_, known := runtimeVersion(info)
			require.False(t, known)
			require.False(t, configUsesV2(info))
		})
	}

	for _, version := range []string{"4.6.1", "4.7.0", "4.7.1", "4.7.2"} {
		require.True(t, configUsesV2(buildInfoForTest(version)), version)
	}

	for _, version := range []string{"4.8.0", "4.8.1", "4.9.0", "5.0.0", "3.3.9"} {
		require.False(t, configUsesV2(buildInfoForTest(version)), version)
	}

	version, known := runtimeVersion(VersionInfo{RawBuildInfoString: " arch:x86_64 ; version: 4.8.0 ; build:123"})
	require.True(t, known)
	require.Equal(t, [3]uint64{4, 8, 0}, version)
}

// TestDiagnosticCompatibilityRequest checks request versions and encoded options.
func TestDiagnosticCompatibilityRequest(t *testing.T) {
	v10, v11, responseVersion := diagnosticCompatibilityVersionsForTest()
	for _, tc := range []struct {
		client, peer   string
		requestVersion uint
	}{
		{"4.8.0", "4.7.1", v11},
		{"4.8.1", "4.7.1", v11},
		{"4.8.0", "4.8.0", v10},
		{"4.8.0", "4.7.0", v10},
		{"4.7.1", "4.7.0", v10},
		{"4.7.1", "4.8.0", v10},
		{"4.8.0", "", v10},
		{"", "4.7.1", v10},
	} {
		t.Run(tc.client+"/"+tc.peer, func(t *testing.T) {
			group := GroupAllGPUs()
			for _, request := range []DiagnosticRequest{
				{Level: DiagQuick, Group: &group},
				{
					TestNames:       []string{"software"},
					TestParameters:  []string{"software.test_duration=2"},
					Group:           &group,
					Timeout:         3 * time.Second,
					EnableHeartbeat: true,
				},
			} {
				sent, nativeVersion, response, err := diagnosticRequestSentForTest(request, buildInfoForTest(tc.client), buildInfoForTest(tc.peer))
				require.NoError(t, err)
				require.Equal(t, tc.requestVersion, nativeVersion)
				require.Equal(t, responseVersion, response)
				require.Equal(t, request, sent)
			}
		})
	}
}

// TestConfigV2PreservesUnitsAndFields checks MHz encoding and verifies that converting back preserves the original config.
func TestConfigV2PreservesUnitsAndFields(t *testing.T) {
	config := nativeConfigFixtureForTest()
	converted, err := configV2RoundTripForTest(config)
	require.NoError(t, err)
	expectedVersion, expectedClock := configV2VersionsForTest()
	require.Equal(t, expectedVersion, converted.version)
	require.Equal(t, expectedClock, converted.clockVersion)
	require.Equal(t, uint(1200), converted.memMHz)
	require.Equal(t, uint(900), converted.smMHz)
	require.Equal(t, config, converted.config)
}

// TestConfigV2ClockBoundsAndBlanks checks that blank clocks survive conversion and invalid Hz values fail.
func TestConfigV2ClockBoundsAndBlanks(t *testing.T) {
	blank32, blank64 := configNativeBlanksForTest()
	for _, hz := range []int64{0, 1000000, (int64(blank32) - 1) * 1000000, blank64, blank64 + 1, blank64 + 2, blank64 + 3, blank64 + 15} {
		native := configClocksForTest(hz, hz)
		converted, err := configV2RoundTripForTest(native)
		require.NoError(t, err)
		require.Equal(t, native, converted.config)
	}
	for _, hz := range []int64{-1, 1, 999999, 1000001, int64(blank32) * 1000000} {
		for _, config := range []nativeGPUConfig{configClocksForTest(hz, 0), configClocksForTest(0, hz)} {
			_, err := configV2RoundTripForTest(config)
			require.ErrorContains(t, err, "cannot be represented")
		}
	}
}

// TestMissingFieldMetadataClassification checks that a missing field reports unsupported only when the catalog is available.
func TestMissingFieldMetadataClassification(t *testing.T) {
	field := DCGM_FI_DEV_GPU_UTIL_RATIO
	for _, present := range []bool{false, true} {
		for _, initialized := range []bool{false, true} {
			metadata, lookups, err := fieldMetadataForTest(present, 0, initialized)
			require.Empty(t, metadata)
			require.ErrorContains(t, err, fmt.Sprintf("field ID %d", field))
			require.Equal(t, []Short{field, DCGM_FI_DEV_GPU_TEMP}, lookups)
			require.Equal(t, initialized, errors.Is(err, errors.ErrUnsupported))
			var native *Error
			require.NotErrorAs(t, err, &native)
		}
	}
}

// TestConfigV2UnsignedBlankClocks checks that extended v2 blank codes map to v3 values without overflowing.
func TestConfigV2UnsignedBlankClocks(t *testing.T) {
	blank32, blank64 := configNativeBlanksForTest()
	for _, tc := range []struct {
		mhz uint32
		hz  int64
	}{
		{blank32, blank64},
		{blank32 + 3, blank64 + 3},
		{blank32 + 15, math.MaxInt64},
		{blank32 + 16, blank64},
		{math.MaxUint32, blank64},
	} {
		config := configFromV2ClocksForTest(tc.mhz)
		require.EqualValues(t, tc.hz, config.perfState.targetClocks.memClock)
		require.EqualValues(t, tc.hz, config.perfState.targetClocks.smClock)
	}
}

// TestRecognizedFieldMetadata checks that a populated entry needs only one lookup.
func TestRecognizedFieldMetadata(t *testing.T) {
	field := DCGM_FI_DEV_GPU_UTIL_RATIO
	metadata, lookups, err := fieldMetadataForTest(true, byte(DCGM_FT_DOUBLE), true)
	require.NoError(t, err)
	require.Equal(t, field, metadata.FieldID)
	require.EqualValues(t, DCGM_FT_DOUBLE, metadata.FieldType)
	require.Equal(t, "gpu_util_ratio", metadata.Tag)
	require.Equal(t, []Short{field}, lookups)
}

// TestRunDiagVersionDispatch checks both native routes and the empty result on failure.
func TestRunDiagVersionDispatch(t *testing.T) {
	_, v11, _ := diagnosticCompatibilityVersionsForTest()
	group := GroupAllGPUs()
	nativeFailure := errors.New("diagnostic failed")
	for _, tc := range []struct {
		name, client, peer string
		version            uint
		err                error
	}{
		{"matched", "4.8.0", "4.8.0", 0, nil},
		{"mixed", "4.8.0", "4.7.1", v11, nil},
		{"mixed failure", "4.8.0", "4.7.1", v11, nativeFailure},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, sent, fullVersion, err := runDiagDispatchForTest(
				DiagMedium, group, buildInfoForTest(tc.client), buildInfoForTest(tc.peer), tc.err,
			)

			require.Equal(t, diagnosticDispatch{calls: 1, level: DiagMedium, group: group, version: tc.version}, sent)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				require.Empty(t, result)
				return
			}
			require.NoError(t, err)
			require.Equal(t, fullVersion, result.DCGMVersion)
			require.Equal(t, "driver", result.DriverVersion)
		})
	}
}
