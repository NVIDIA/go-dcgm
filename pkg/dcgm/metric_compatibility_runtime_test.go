package dcgm

import (
	"fmt"
	"math"
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
		t.Skip("representative metric reads require a physical GPU")
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
