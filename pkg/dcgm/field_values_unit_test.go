package dcgm

import (
	"strconv"
	"strings"
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
	t.Run("conversion", testFieldValueConversions)

	t.Run("bounded native string member", func(t *testing.T) {
		want := strings.Repeat("x", fieldStringCapacity)
		spec := testFieldValueSpec{
			fieldID:   DCGM_FI_DEV_PCI_BUS_ID,
			fieldType: DCGM_FT_STRING,
			payload:   []byte(want + "outside the string member"),
		}
		v1Fields := makeTestCFieldsFromSpec(1, spec)
		require.Equal(t, want, toFieldValue(v1Fields)[0].String())

		for _, value := range []FieldValue_v2{
			dcgmFieldValue_v1ToFieldValue_v2(FE_GPU, 0, v1Fields)[0],
			toFieldValue_v2(makeTestCFieldsV2FromSpec(1, 1, spec))[0],
			appendConvertedValues(nil, FE_GPU, 0, v1Fields)[0],
		} {
			require.Equal(t, want, value.String())
			require.NotNil(t, value.StringValue)
			require.Equal(t, want, *value.StringValue)
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

	t.Run("field-ids", testFieldIDPointerView)
}

// testFieldIDPointerView checks that the native pointer sees the Go slice's updates.
func testFieldIDPointerView(t *testing.T) {
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
}

// testFieldValueConversions checks the Go values produced by each native converter.
func testFieldValueConversions(t *testing.T) {
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

			v1Fields := makeTestCFieldsFromSpec(1, spec)
			assertConvertedFieldV1(t, &toFieldValue(v1Fields)[0], spec, &wantPayload)
			assertConvertedFieldV2(t, &dcgmFieldValue_v1ToFieldValue_v2(FE_GPU, 0, v1Fields)[0],
				spec, &wantPayload, 0, defaultTimestamp)

			values := toFieldValue_v2(makeTestCFieldsV2FromSpec(tt.entities, 1, spec))
			require.Len(t, values, tt.entities)
			for i := range values {
				assertConvertedFieldV2(t, &values[i], spec, &wantPayload, uint(i), defaultTimestamp+int64(i))
			}
		})
	}
}

func assertConvertedFieldV1(t *testing.T, value *FieldValue_v1, spec testFieldValueSpec, wantPayload *[4096]byte) {
	t.Helper()
	assert.NotZero(t, value.Version)
	assert.Equal(t, spec.fieldID, value.FieldID)
	assert.Equal(t, spec.fieldType, value.FieldType)
	assert.Equal(t, spec.status, value.Status)
	assert.Equal(t, spec.timestamp, value.TS)
	assert.Equal(t, *wantPayload, value.Value)
	if spec.fieldType == DCGM_FT_STRING {
		assert.Equal(t, string(spec.payload[:len(spec.payload)-1]), value.String())
	}
}

func assertConvertedFieldV2(
	t *testing.T, value *FieldValue_v2, spec testFieldValueSpec, wantPayload *[4096]byte, entityID uint, timestamp int64,
) {
	t.Helper()
	assert.NotZero(t, value.Version)
	assert.Equal(t, FE_GPU, value.EntityGroupId)
	assert.Equal(t, entityID, value.EntityID)
	assert.Equal(t, spec.fieldID, value.FieldID)
	assert.Equal(t, spec.fieldType, value.FieldType)
	assert.Equal(t, spec.status, value.Status)
	assert.Equal(t, timestamp, value.TS)
	assert.Equal(t, *wantPayload, value.Value)
	if spec.fieldType == DCGM_FT_STRING {
		require.NotNil(t, value.StringValue)
		want := string(spec.payload[:len(spec.payload)-1])
		assert.Equal(t, want, *value.StringValue)
		assert.Equal(t, want, value.String())
	}
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

func TestCheckFieldValueCount(t *testing.T) {
	const maxSamples = 4
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{name: "negative", count: -1, wantErr: true},
		{name: "zero", count: 0},
		{name: "max", count: maxSamples},
		{name: "above max", count: maxSamples + 1, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkFieldValueCount(tt.count, maxSamples)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
