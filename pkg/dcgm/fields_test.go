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

func TestFieldHandle(t *testing.T) {
	fh := FieldHandle{}
	assert.Equal(t, uintptr(0), fh.GetHandle(), "value mismatch")

	inputs := []uintptr{1000, 0, 1, 10, 11, 50, 100, 1939902, 9992932938239, 999999999999999999}

	for _, input := range inputs {
		fh.SetHandle(input)
		assert.Equal(t, input, fh.GetHandle(), "values mismatch")
	}
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
		err := DestroyGroup(groupId)
		require.NoError(t, err)
	}()

	// Setup field group
	fieldId := Short(DCGM_FI_DEV_XID_ERRORS)
	n, err := crand.Int(crand.Reader, big.NewInt(1000000))
	require.NoError(t, err)
	fieldGroupName := fmt.Sprintf("fieldGroupName%d", n.Int64())
	fieldsGroup, err := FieldGroupCreate(fieldGroupName, []Short{fieldId})
	require.NoError(t, err)
	defer func() {
		err := FieldGroupDestroy(fieldsGroup)
		require.NoError(t, err)
	}()

	// Inject test value
	err = InjectFieldValue(gpuId,
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
		Short(DCGM_FI_DEV_XID_ERRORS),
		Short(DCGM_FI_DEV_DIAG_MEMORY_RESULT),
		Short(DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION),
		Short(DCGM_FI_DEV_GPU_TEMP),
		Short(DCGM_FI_DEV_MEMORY_TEMP),
		Short(DCGM_FI_DEV_GPU_UTIL),
		Short(DCGM_FI_DEV_MEM_COPY_UTIL),
		Short(DCGM_FI_DEV_ENC_UTIL),
		Short(DCGM_FI_DEV_DEC_UTIL),
		Short(DCGM_FI_DEV_FB_FREE),
		Short(DCGM_FI_DEV_FB_USED),
		Short(DCGM_FI_DEV_PCIE_REPLAY_COUNTER),
		Short(DCGM_FI_DEV_SM_CLOCK),
		Short(DCGM_FI_DEV_RETIRED_PENDING),
		Short(DCGM_FI_DEV_RETIRED_SBE),
		Short(DCGM_FI_DEV_RETIRED_DBE),
		Short(DCGM_FI_DEV_POWER_VIOLATION),
		Short(DCGM_FI_DEV_THERMAL_VIOLATION),
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
				err := FieldGroupDestroy(fieldsGroup)
				require.NoError(b, err)
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
				err = InjectFieldValue(gpuId,
					uint(fieldId),
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
