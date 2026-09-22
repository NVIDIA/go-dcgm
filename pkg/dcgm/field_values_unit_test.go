package dcgm

import (
	"strconv"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValueWatchRejectsInvalidArguments(t *testing.T) {
	tests := []struct {
		name            string
		updateFrequency time.Duration
		maxKeepAge      time.Duration
		maxKeepSamples  int
	}{
		{name: "zero update frequency", maxKeepSamples: 1},
		{name: "sub-microsecond update frequency", updateFrequency: time.Nanosecond, maxKeepSamples: 1},
		{name: "no retention bound", updateFrequency: time.Second},
		{name: "negative retention age", updateFrequency: time.Second, maxKeepAge: -time.Second, maxKeepSamples: 1},
		{name: "negative retention sample count", updateFrequency: time.Second, maxKeepAge: time.Second, maxKeepSamples: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WatchFieldValue(0, DCGM_FI_SYSTEM_GPU_BIND_EVENT,
				tt.updateFrequency, tt.maxKeepAge, tt.maxKeepSamples)
			var dcgmErr *Error
			require.ErrorAs(t, err, &dcgmErr)
			assert.Equal(t, DCGM_ST_BADPARAM, int(dcgmErr.Code))
		})
	}

	if strconv.IntSize > 32 {
		t.Run("retention sample count exceeds C int", func(t *testing.T) {
			err := WatchFieldValue(
				0, DCGM_FI_SYSTEM_GPU_BIND_EVENT, time.Second, 0, maxWatchFieldValueSamples+1,
			)
			var dcgmErr *Error
			require.ErrorAs(t, err, &dcgmErr)
			assert.Equal(t, DCGM_ST_BADPARAM, int(dcgmErr.Code))
		})
	}

	for _, maxSamples := range []int{0, maxFieldValueHistorySamples + 1} {
		_, err := GetMultipleValuesForField(0, DCGM_FI_SYSTEM_GPU_BIND_EVENT, maxSamples, time.Time{}, time.Time{})
		var dcgmErr *Error
		require.ErrorAs(t, err, &dcgmErr)
		assert.Equal(t, DCGM_ST_BADPARAM, int(dcgmErr.Code))
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

func TestTimestampUSECToTime(t *testing.T) {
	tests := []struct {
		name  string
		input int64
	}{
		{"epoch", 0},
		{"one microsecond", 1},
		{"fractional second", 1_500_001},
		{"before epoch", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timestampUSECToTime(tt.input).UnixMicro(); got != tt.input {
				t.Fatalf("timestampUSECToTime(%d).UnixMicro() = %d", tt.input, got)
			}
		})
	}
}
