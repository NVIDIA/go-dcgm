package dcgm

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
	"unsafe"

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

func int64Bytes(value int64) []byte {
	payload := *(*[8]byte)(unsafe.Pointer(&value))
	return payload[:]
}

func float64Bytes(value float64) []byte {
	payload := *(*[8]byte)(unsafe.Pointer(&value))
	return payload[:]
}

func TestFieldValueAllocationPaths(t *testing.T) {
	t.Run("conversion", func(t *testing.T) {
		const (
			defaultFieldID   = Short(100)
			defaultStatus    = -3
			defaultTimestamp = int64(1700000000000000)
		)

		tests := []struct {
			name      string
			fieldID   Short
			fieldType uint
			payload   []byte
			entities  int
		}{
			{name: "int64", fieldType: DCGM_FT_INT64, payload: int64Bytes(42), entities: 1},
			{name: "double-profiling", fieldID: DCGM_FI_PROF_GR_ENGINE_ACTIVE, fieldType: DCGM_FT_DOUBLE, payload: float64Bytes(0.5), entities: 1},
			{name: "string", fieldType: DCGM_FT_STRING, payload: append([]byte("GPU-test"), 0), entities: 1},
			{name: "string-blank", fieldType: DCGM_FT_STRING, payload: append([]byte("<<<NULL>>>"), 0), entities: 1},
			{name: "string-error", fieldType: DCGM_FT_STRING, payload: append([]byte("<<<NOT_SUPPORTED>>>"), 0), entities: 1},
			{name: "numeric-blank", fieldType: DCGM_FT_INT64, payload: int64Bytes(DCGM_FT_INT64_BLANK), entities: 1},
			{name: "numeric-error", fieldType: DCGM_FT_INT64, payload: int64Bytes(DCGM_FT_INT64_NOT_SUPPORTED), entities: 1},
			{name: "timestamp", fieldType: DCGM_FT_TIMESTAMP, payload: int64Bytes(defaultTimestamp), entities: 1},
			{name: "binary", fieldType: DCGM_FT_BINARY, payload: []byte{0xde, 0xad, 0xbe, 0xef}, entities: 1},
			{name: "multi-entity", fieldType: DCGM_FT_INT64, payload: int64Bytes(7), entities: 8},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fieldID := tt.fieldID
				if fieldID == 0 {
					fieldID = defaultFieldID
				}
				spec := testFieldValueSpec{
					fieldID:   fieldID,
					fieldType: tt.fieldType,
					status:    defaultStatus,
					timestamp: defaultTimestamp,
					payload:   tt.payload,
				}
				var wantPayload [4096]byte
				copy(wantPayload[:], tt.payload)

				assertV1 := func(value FieldValue_v1) {
					assert.NotZero(t, value.Version)
					assert.Equal(t, fieldID, value.FieldID)
					assert.Equal(t, tt.fieldType, value.FieldType)
					assert.Equal(t, defaultStatus, value.Status)
					assert.Equal(t, defaultTimestamp, value.TS)
					assert.Equal(t, wantPayload, value.Value)
					if tt.fieldType == DCGM_FT_STRING {
						assert.Equal(t, string(tt.payload[:len(tt.payload)-1]), value.String())
					}
				}
				assertV2 := func(value FieldValue_v2, entityID uint, timestamp int64) {
					assert.NotZero(t, value.Version)
					assert.Equal(t, FE_GPU, value.EntityGroupId)
					assert.Equal(t, entityID, value.EntityID)
					assert.Equal(t, fieldID, value.FieldID)
					assert.Equal(t, tt.fieldType, value.FieldType)
					assert.Equal(t, defaultStatus, value.Status)
					assert.Equal(t, timestamp, value.TS)
					assert.Equal(t, wantPayload, value.Value)
					if tt.fieldType == DCGM_FT_STRING {
						require.NotNil(t, value.StringValue)
						want := string(tt.payload[:len(tt.payload)-1])
						assert.Equal(t, want, *value.StringValue)
						assert.Equal(t, want, value.String())
					}
				}

				v1Fields := makeTestCFieldsFromSpec(1, spec)
				assertV1(toFieldValue(v1Fields)[0])
				assertV2(dcgmFieldValue_v1ToFieldValue_v2(FE_GPU, 0, v1Fields)[0], 0, defaultTimestamp)

				values := toFieldValue_v2(makeTestCFieldsV2FromSpec(tt.entities, 1, spec))
				require.Len(t, values, tt.entities)
				expectedEntityID := uint(0)
				expectedTimestamp := defaultTimestamp
				for _, value := range values {
					assertV2(value, expectedEntityID, expectedTimestamp)
					expectedEntityID++
					expectedTimestamp++
				}
			})
		}
	})

	t.Run("ownership", func(t *testing.T) {
		spec := testFieldValueSpec{
			fieldID:   100,
			fieldType: DCGM_FT_STRING,
			payload:   append([]byte("owned"), 0),
		}
		native := acquireFieldValueSlice(1)
		native.values[0] = makeTestCFieldsFromSpec(1, spec)[0]
		value := toFieldValue(native.values)[0]
		want := value
		for i := range native.values[0].value {
			native.values[0].value[i] = 0xff
		}
		releaseFieldValueSlice(native)

		assert.Equal(t, want, value)
		assert.Equal(t, "owned", value.String())
	})
	t.Run("pool", func(t *testing.T) {
		tests := []struct {
			name     string
			size     int
			capacity int
			want     bool
		}{
			{name: "exact", size: 1, capacity: 1, want: true},
			{name: "compatible", size: 8, capacity: 32, want: true},
			{name: "too-small", size: 32, capacity: 8},
			{name: "more-than-four-times", size: 1, capacity: 32},
			{name: "pool-boundary", size: 256, capacity: 256, want: true},
			{name: "above-pool-boundary", size: 257, capacity: 512},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, reusableSliceCapacity(tt.size, tt.capacity))
			})
		}

		t.Run("release-retains-holder", func(t *testing.T) {
			var pool sync.Pool
			pooled := acquirePooledSlice[int](&pool, 3)
			copy(pooled.values, []int{1, 2, 3})
			holder := pooled.holder

			releasePooledSlice(&pool, pooled)

			require.NotNil(t, holder)
			assert.Empty(t, *holder)
			assert.Equal(t, []int{0, 0, 0}, pooled.values)
		})

		var pool sync.Pool
		assert.Equal(t, 8, cap(acquirePooledSlice[int](&pool, 8).values))
		assert.Equal(t, poolCapacityThreshold+1, cap(acquirePooledSlice[int](&pool, poolCapacityThreshold+1).values))
	})

	t.Run("field-ids", func(t *testing.T) {
		for _, size := range []int{1, 8, 32, 128} {
			t.Run(strconv.Itoa(size), func(t *testing.T) {
				fields := make([]Short, size)
				fieldID := Short(100)
				for i := range fields {
					fields[i] = fieldID
					fieldID++
				}

				view := unsafe.Slice(fieldIDPointer(fields), len(fields))
				fields[0]++
				for i := range fields {
					assert.Equal(t, fields[i], Short(view[i]))
				}
			})
		}
	})
}

func TestLatestValuesRejectEmptyInputs(t *testing.T) {
	entities := []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: 0}}
	fields := []Short{DCGM_FI_DEV_GPU_TEMP}

	tests := []struct {
		name string
		call func() (bool, error)
	}{
		{
			name: "GetLatestValuesForFields/empty-fields",
			call: func() (bool, error) {
				values, err := GetLatestValuesForFields(0, nil)
				return values == nil, err
			},
		},
		{
			name: "EntityGetLatestValues/empty-fields",
			call: func() (bool, error) {
				values, err := EntityGetLatestValues(FE_GPU, 0, nil)
				return values == nil, err
			},
		},
		{
			name: "LinkGetLatestValues/empty-fields",
			call: func() (bool, error) {
				values, err := LinkGetLatestValues(0, FE_GPU, 0, nil)
				return values == nil, err
			},
		},
		{
			name: "EntitiesGetLatestValues/empty-fields",
			call: func() (bool, error) {
				values, err := EntitiesGetLatestValues(entities, nil, 0)
				return values == nil, err
			},
		},
		{
			name: "EntitiesGetLatestValues/empty-entities",
			call: func() (bool, error) {
				values, err := EntitiesGetLatestValues(nil, fields, 0)
				return values == nil, err
			},
		},
		{
			name: "EntitiesGetLatestValues/empty-entities-and-fields",
			call: func() (bool, error) {
				values, err := EntitiesGetLatestValues(nil, nil, 0)
				return values == nil, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultIsNil, err := tt.call()

			require.True(t, resultIsNil)
			var dcgmErr *Error
			require.ErrorAs(t, err, &dcgmErr)
			assert.Equal(t, DCGM_ST_BADPARAM, int(dcgmErr.Code))
		})
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
