package dcgm

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalFieldWatchDoesNotRequireGPUs(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	gpus, err := GetSupportedDevices()
	require.NoError(t, err)
	if len(gpus) != 0 {
		t.Skipf("requires a DCGM instance without physical GPUs; found %d", len(gpus))
	}
	require.Empty(t, gpus, "test requires a DCGM instance without physical GPUs")

	require.NoError(t, WatchFieldValue(
		0,
		DCGM_FI_SYSTEM_GPU_BIND_EVENT,
		time.Second,
		0,
		2,
	))
	defer func() {
		assert.NoError(t, UnwatchFieldValue(0, DCGM_FI_SYSTEM_GPU_BIND_EVENT, true))
	}()

	values, err := EntityGetLatestValues(FE_NONE, 0, []Short{DCGM_FI_SYSTEM_GPU_BIND_EVENT})
	require.NoError(t, err)
	require.Len(t, values, 1)
	assert.Equal(t, DCGM_FI_SYSTEM_GPU_BIND_EVENT, values[0].FieldID)

	history, err := GetMultipleValuesForField(
		0, DCGM_FI_SYSTEM_GPU_BIND_EVENT, 2, time.Time{}, time.Time{},
	)
	require.NoError(t, err)
	require.NotEmpty(t, history)
	assert.Equal(t, DCGM_FI_SYSTEM_GPU_BIND_EVENT, history[0].FieldID)
}

func TestGetLatestValuesForFields(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	runOnlyWithLiveGPUs(t)

	// Setup test GPU
	gpus, err := withInjectionGPUs(t, 1)
	require.NoError(t, err)
	gpuId := gpus[0]

	// Setup test group
	groupId, err := NewDefaultGroup("mygroup")
	require.NoError(t, err)
	defer func() {
		destroyGroupErr := DestroyGroup(groupId)
		require.NoError(t, destroyGroupErr)
	}()

	// Setup field group
	fieldId := DCGM_FI_DEV_XID_ERRORS
	n, err := crand.Int(crand.Reader, big.NewInt(1000000))
	require.NoError(t, err)
	fieldGroupName := fmt.Sprintf("fieldGroupName%d", n.Int64())
	fieldsGroup, err := FieldGroupCreate(fieldGroupName, []Short{fieldId})
	require.NoError(t, err)
	defer func() {
		destroyFieldsGroupErr := FieldGroupDestroy(fieldsGroup)
		require.NoError(t, destroyFieldsGroupErr)
	}()

	// Inject test value
	err = InjectFieldValue(
		gpuId,
		DCGM_FI_DEV_XID_ERRORS,
		DCGM_FT_INT64,
		0,
		time.Now().Add(-time.Duration(5)*time.Second).UnixMicro(),
		int64(10),
	)
	require.NoError(t, err)

	// Setup field watching
	err = WatchFieldsWithGroupEx(
		fieldsGroup,
		groupId,
		defaultUpdateFreq,
		defaultMaxKeepAge,
		defaultMaxKeepSamples,
	)
	require.NoError(t, err)

	err = UpdateAllFields()
	require.NoError(t, err)

	// Test
	values, err := GetLatestValuesForFields(gpuId, []Short{fieldId})
	require.NoError(t, err)

	// Verify results
	assert.Len(t, values, 1)
	assert.NotEmpty(t, values[0].String())
	assert.Equal(t, int64(10), values[0].Int64())
}

func BenchmarkGetLatestValuesForFieldsVariousSize(b *testing.B) {
	teardownTest := setupTest(b)
	defer teardownTest(b)

	// Setup test GPU
	gpus, err := withInjectionGPUs(b, 1)
	require.NoError(b, err)
	gpuId := gpus[0]

	// Setup test group
	groupId, err := NewDefaultGroup("mygroup")
	require.NoError(b, err)
	defer func() {
		err := DestroyGroup(groupId)
		require.NoError(b, err)
	}()

	// Use the same fields as in the main benchmark
	allFieldIds := []Short{
		DCGM_FI_DEV_XID_ERRORS,
		DCGM_FI_DEV_DIAG_MEMORY_RESULT,
		DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION,
		DCGM_FI_DEV_GPU_TEMP,
		DCGM_FI_DEV_MEMORY_TEMP,
		DCGM_FI_DEV_GPU_UTIL,
		DCGM_FI_DEV_MEM_COPY_UTIL,
		DCGM_FI_DEV_ENC_UTIL,
		DCGM_FI_DEV_DEC_UTIL,
		DCGM_FI_DEV_FB_FREE,
		DCGM_FI_DEV_FB_USED,
		DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
		DCGM_FI_DEV_SM_CLOCK,
		DCGM_FI_DEV_RETIRED_PENDING,
		DCGM_FI_DEV_RETIRED_SBE,
		DCGM_FI_DEV_RETIRED_DBE,
		DCGM_FI_DEV_POWER_VIOLATION,
		DCGM_FI_DEV_THERMAL_VIOLATION,
	}

	// Test different field counts
	fieldCounts := []int{1, 5, 10, len(allFieldIds)}

	for _, count := range fieldCounts {
		b.Run(fmt.Sprintf("Fields-%d", count), func(b *testing.B) {
			fieldIds := allFieldIds[:count] // Take first 'count' fields

			// Setup field group
			fieldGroupName := fmt.Sprintf("fieldGroup-%d", count)
			fieldsGroup, err := FieldGroupCreate(fieldGroupName, fieldIds)
			require.NoError(b, err)
			defer func() {
				destroyFieldsGroupErr := FieldGroupDestroy(fieldsGroup)
				require.NoError(b, destroyFieldsGroupErr)
			}()

			// Setup field watching
			err = WatchFieldsWithGroupEx(
				fieldsGroup,
				groupId,
				defaultUpdateFreq,
				defaultMaxKeepAge,
				defaultMaxKeepSamples,
			)
			require.NoError(b, err)

			// Inject values for all fields
			for _, fieldId := range fieldIds {
				err = InjectFieldValue(
					gpuId,
					fieldId,
					DCGM_FT_INT64,
					0,
					time.Now().Add(-time.Duration(5)*time.Second).UnixMicro(),
					int64(10),
				)
				require.NoError(b, err)
			}

			err = UpdateAllFields()
			require.NoError(b, err)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				values, err := GetLatestValuesForFields(gpuId, fieldIds)
				require.NoError(b, err)
				require.Len(b, values, len(fieldIds), "expected %d values, got %d", len(fieldIds), len(values))
				runtime.KeepAlive(values)
			}
		})
	}
}
