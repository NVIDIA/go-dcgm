package dcgm

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestRepresentativeMetricCompatibility checks device attributes and fresh native
// samples through the same public binding calls on each supported DCGM runtime.
func TestRepresentativeMetricCompatibility(t *testing.T) {
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()

	gpus, err := GetSupportedDevices()
	require.NoError(t, err)
	if len(gpus) == 0 {
		t.Fatal("representative metric reads require a physical GPU")
	}

	fields := []Short{DCGM_FI_DEV_SM_CLOCK, DCGM_FI_DEV_GPU_TEMP_CELSIUS, DCGM_FI_DEV_POWER_USAGE, DCGM_FI_DEV_FB_USED}
	group, err := CreateGroup("metric-compatibility")
	require.NoError(t, err)
	defer func() { require.NoError(t, DestroyGroup(group)) }()
	for _, gpu := range gpus {
		device, infoErr := GetDeviceInfo(gpu)
		require.NoError(t, infoErr)
		require.Equal(t, gpu, device.GPU)
		require.NotEmpty(t, device.UUID)
		require.NotEmpty(t, device.Identifiers.Model)
		require.Positive(t, device.PCI.FBTotal)
		require.NoError(t, AddToGroup(group, gpu))
	}

	fieldGroup, err := FieldGroupCreate("metric-compatibility", fields)
	require.NoError(t, err)
	defer func() { require.NoError(t, FieldGroupDestroy(fieldGroup)) }()
	require.NoError(t, WatchFieldsWithGroupExNoRefresh(fieldGroup, group, 100_000, 5, 0))
	defer func() { require.NoError(t, UnwatchFields(fieldGroup, group)) }()

	// A second observation must advance every sample, rather than merely read
	// a value left in the cache. Poll on the test goroutine so failures are safe.
	before := compatibilityMetricTimestamps(t, gpus, fields)
	deadline := time.Now().Add(15 * time.Second)
	for {
		after := compatibilityMetricTimestamps(t, gpus, fields)
		advanced := true
		for key, timestamp := range before {
			advanced = advanced && after[key] > timestamp
		}
		if advanced {
			t.Logf("fresh samples for %d GPUs and %d fields", len(gpus), len(fields))
			return
		}

		require.True(t, time.Now().Before(deadline), "native metric timestamps did not advance: before=%v after=%v", before, after)
		time.Sleep(100 * time.Millisecond)
	}
}

// compatibilityMetricTimestamps checks native identity, status, type and value
// bounds before returning the latest timestamp for each GPU and field.
func compatibilityMetricTimestamps(t *testing.T, gpus []uint, fields []Short) map[string]int64 {
	t.Helper()
	timestamps := make(map[string]int64, len(gpus)*len(fields))
	for _, gpu := range gpus {
		values, err := GetLatestValuesForFields(gpu, fields)
		require.NoError(t, err)
		require.Len(t, values, len(fields))
		for i := range values {
			value := &values[i]
			require.Equal(t, fields[i], value.FieldID)
			require.EqualValues(t, DCGM_ST_OK, value.Status)
			require.Positive(t, value.TS)
			if fields[i] == DCGM_FI_DEV_POWER_USAGE {
				require.EqualValues(t, DCGM_FT_DOUBLE, value.FieldType)
				number := value.Float64()
				require.Less(t, number, float64(DCGM_FT_FP64_BLANK))
				require.False(t, math.IsNaN(number) || math.IsInf(number, 0))
				require.GreaterOrEqual(t, number, float64(0))
			} else {
				require.EqualValues(t, DCGM_FT_INT64, value.FieldType)
				number := value.Int64()
				require.False(t, IsInt64Blank(number))
				require.GreaterOrEqual(t, number, int64(0))
			}
			timestamps[fmt.Sprintf("%d:%d", gpu, value.FieldID)] = value.TS
		}
	}
	return timestamps
}

// TestLegacyHostengineCompatibility checks memory reads from a 4.6.1 hostengine
// with the 4.6.1 and 4.8.0 clients. The matched pair also checks config writes.
func TestLegacyHostengineCompatibility(t *testing.T) {
	socket := os.Getenv("DCGM_TEST_LEGACY_HOSTENGINE_SOCKET")
	expectedClient := os.Getenv("DCGM_TEST_LEGACY_CLIENT_VERSION")
	if socket == "" && expectedClient == "" {
		t.Skip("set DCGM_TEST_LEGACY_HOSTENGINE_SOCKET and DCGM_TEST_LEGACY_CLIENT_VERSION")
	}
	require.NotEmpty(t, socket, "set DCGM_TEST_LEGACY_HOSTENGINE_SOCKET to a standalone DCGM 4.6.1 socket")
	require.Contains(t, []string{"4.6.1", "4.8.0"}, expectedClient,
		"set DCGM_TEST_LEGACY_CLIENT_VERSION to the expected client version")

	cleanup, err := Init(Standalone, socket, "1")
	require.NoError(t, err)
	defer cleanup()

	clientInfo, err := GetVersionInfo()
	require.NoError(t, err)

	client, known := runtimeVersion(clientInfo)
	require.True(t, known, clientInfo.RawBuildInfoString)
	require.Equal(t, expectedClient, fmt.Sprintf("%d.%d.%d", client[0], client[1], client[2]))

	peerInfo, err := GetHostengineVersionInfo()
	require.NoError(t, err)
	peer, known := runtimeVersion(peerInfo)
	require.True(t, known, peerInfo.RawBuildInfoString)
	require.Equal(t, [3]uint64{4, 6, 1}, peer)
	t.Logf("DCGM client %v, hostengine %v", client, peer)

	gpus, err := GetSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus, "the hostengine must serve a supported GPU")

	for _, gpu := range gpus {
		t.Run(fmt.Sprintf("memory_GPU_%d", gpu), func(t *testing.T) {
			device, err := GetDeviceInfo(gpu)
			require.NoError(t, err)
			require.NotEmpty(t, device.UUID)

			fields := []Short{DCGM_FI_DEV_BAR1_TOTAL, DCGM_FI_DEV_FB_TOTAL}
			values, err := EntitiesGetLatestValues(
				[]GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpu}}, fields, DCGM_FV_FLAG_LIVE_DATA,
			)
			require.NoError(t, err)
			require.Len(t, values, len(fields))

			for i, capacity := range []uint{device.PCI.BAR1, device.PCI.FBTotal} {
				require.Equal(t, fields[i], values[i].FieldID)
				require.Equal(t, DCGM_ST_OK, values[i].Status)
				require.Equal(t, DCGM_FT_INT64, values[i].FieldType)

				want := values[i].Int64()
				require.Positive(t, want)
				require.False(t, IsInt64Blank(want))
				require.EqualValues(t, want, capacity)
			}

			t.Logf("BAR1 %d MiB, framebuffer %d MiB", device.PCI.BAR1, device.PCI.FBTotal)
		})
	}

	if client == peer {
		t.Run("legacy_config_get", func(t *testing.T) {
			checkLegacyConfigGet(t, gpus[0])
		})

		t.Run("legacy_config_set", func(t *testing.T) {
			group, err := CreateGroup("legacy-config-empty")
			require.NoError(t, err)
			defer func() { require.NoError(t, DestroyGroup(group)) }()

			// Exercise the native write path with an empty group and no settings.
			var native *Error
			require.ErrorAs(t, ConfigSet(group, GPUConfig{}), &native)
			require.EqualValues(t, DCGM_ST_NOT_CONFIGURED, native.Code)
		})
	}
}

// checkLegacyConfigGet verifies the v2 read on a supported GPU.
func checkLegacyConfigGet(t *testing.T, gpuID uint) {
	t.Helper()

	group, err := CreateGroup("legacy-config-read")
	require.NoError(t, err)
	defer func() { require.NoError(t, DestroyGroup(group)) }()
	require.NoError(t, AddToGroup(group, gpuID))

	configs, err := ConfigGet(group, ConfigCurrentState)
	if err != nil {
		// The older hostengine requires root to serve current config.
		var native *Error
		require.ErrorAs(t, err, &native)
		require.EqualValues(t, DCGM_ST_REQUIRES_ROOT, native.Code)
		require.Empty(t, configs)
		t.Log("ConfigGet reached the native read and reported that the hostengine requires root")
		return
	}
	require.Len(t, configs, 1)
	require.Equal(t, gpuID, configs[0].GPU)
	for _, clock := range []*int64{configs[0].MemoryClockHz, configs[0].SMClockHz} {
		if clock != nil {
			require.Zero(t, *clock%1000000, "legacy clocks are whole MHz expressed in Hz")
		}
	}
}
