package dcgm

import (
	"encoding/binary"
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntityStatusString(t *testing.T) {
	tests := []struct {
		name  string
		input EntityStatus
		want  string
	}{
		{"unknown", EntityStatusUnknown, "Unknown"},
		{"ok", EntityStatusOk, "OK"},
		{"unsupported", EntityStatusUnsupported, "Unsupported"},
		{"inaccessible", EntityStatusInaccessible, "Inaccessible"},
		{"lost", EntityStatusLost, "Lost"},
		{"fake", EntityStatusFake, "Fake"},
		{"disabled", EntityStatusDisabled, "Disabled"},
		{"detached", EntityStatusDetached, "Detached"},
		{"unrecognized", EntityStatus(99), "Unknown(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Fatalf("EntityStatus(%d).String() = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPerfStateString(t *testing.T) {
	tests := []struct {
		name  string
		input PerfState
		want  string
	}{
		{"maximum", 0, "P0"},
		{"middle", 7, "P7"},
		{"minimum", 15, "P15"},
		{"above range", 16, "Unknown"},
		{"unknown sentinel", 32, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.want {
				t.Fatalf("PerfState(%d).String() = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDeviceStatusFromValues(t *testing.T) {
	values := deviceStatusSamples()
	values[devicePower] = testFloat64FieldValue(250.5)
	values[deviceTemperature] = testInt64FieldValue(70)
	values[deviceGPUUtilization] = testInt64FieldValue(80)
	values[deviceMemoryUtilization] = testInt64FieldValue(60)
	values[deviceEncoderUtilization] = testInt64FieldValue(40)
	values[deviceDecoderUtilization] = testInt64FieldValue(30)
	values[deviceSMClock] = testInt64FieldValue(1200)
	values[deviceMemoryClock] = testInt64FieldValue(5000)
	values[deviceBAR1Used] = testInt64FieldValue(100)
	values[devicePCIeRxThroughput] = testInt64FieldValue(10)
	values[devicePCIeTxThroughput] = testInt64FieldValue(20)
	values[devicePCIeReplay] = testInt64FieldValue(2)
	values[deviceFBUsed] = testInt64FieldValue(8000)
	values[deviceSingleBitErrors] = testInt64FieldValue(3)
	values[deviceDoubleBitErrors] = testInt64FieldValue(4)
	values[devicePerformanceState] = testInt64FieldValue(7)
	values[deviceFanSpeed] = testInt64FieldValue(55)

	for i, id := range deviceStatusFields {
		values[i].FieldID = id
	}
	slices.Reverse(values)

	tests := []struct {
		name   string
		values []FieldValue_v1
		want   DeviceStatus
	}{
		{
			name:   "zero values",
			values: deviceStatusSamples(),
		},
		{
			name:   "maps every field",
			values: values,
			want: DeviceStatus{
				Power:       250.5,
				Temperature: 70,
				Utilization: UtilizationInfo{GPU: 80, Memory: 60, Encoder: 40, Decoder: 30},
				Memory:      MemoryInfo{ECCErrors: ECCErrorsInfo{SingleBit: 3, DoubleBit: 4}},
				Clocks:      ClockInfo{Cores: 1200, Memory: 5000},
				PCI: PCIStatusInfo{
					BAR1Used:   100,
					Throughput: PCIThroughputInfo{Rx: 10, Tx: 20, Replays: 2},
					FBUsed:     8000,
				},
				Performance: 7,
				FanSpeed:    55,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := deviceStatusFromValues(tt.values); err != nil || got != tt.want {
				t.Fatalf("deviceStatusFromValues() = %+v, %v, want %+v, nil", got, err, tt.want)
			}
		})
	}
}

func deviceStatusSamples() []FieldValue_v1 {
	values := make([]FieldValue_v1, len(deviceStatusFields))
	for i, fieldID := range deviceStatusFields {
		values[i] = FieldValue_v1{FieldID: fieldID, FieldType: DCGM_FT_INT64}
	}
	values[devicePower].FieldType = DCGM_FT_DOUBLE
	return values
}

func TestDeviceStatusRejectsMalformedResponse(t *testing.T) {
	for _, tt := range []struct {
		name   string
		change func([]FieldValue_v1) []FieldValue_v1
		want   string
	}{
		{"short response", func(v []FieldValue_v1) []FieldValue_v1 { return v[:1] }, "missing GPU status field"},
		{"duplicate field", func(v []FieldValue_v1) []FieldValue_v1 { v[1] = v[0]; return v }, "duplicate GPU status field"},
		{"unknown field", func(v []FieldValue_v1) []FieldValue_v1 { v[1].FieldID = 0; return v }, "unexpected GPU status field"},
		{"wrong field type", func(v []FieldValue_v1) []FieldValue_v1 { v[0].FieldType = DCGM_FT_BINARY; return v }, "unexpected type"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := deviceStatusFromValues(tt.change(deviceStatusSamples()))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestDeviceStatusPreservesAvailableFields(t *testing.T) {
	for _, tt := range []struct {
		name      string
		status    int
		wantInt   int64
		wantFloat float64
		wantError bool
	}{
		{"no data", DCGM_ST_NO_DATA, DCGM_FT_INT64_BLANK, DCGM_FT_FP64_BLANK, false},
		{"not found", DCGM_ST_UNKNOWN_FIELD, DCGM_FT_INT64_NOT_FOUND, DCGM_FT_FP64_NOT_FOUND, false},
		{"unsupported", DCGM_ST_NOT_SUPPORTED, DCGM_FT_INT64_NOT_SUPPORTED, DCGM_FT_FP64_NOT_SUPPORTED, false},
		{"permission", DCGM_ST_NO_PERMISSION, DCGM_FT_INT64_NOT_PERMISSIONED, DCGM_FT_FP64_NOT_PERMISSIONED, false},
		{"GPU lost", DCGM_ST_GPU_IS_LOST, DCGM_FT_INT64_BLANK, DCGM_FT_FP64_BLANK, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			values := deviceStatusSamples()
			for _, i := range []int{devicePower, deviceEncoderUtilization} {
				values[i].Status = tt.status
				values[i].FieldType = DCGM_FT_BINARY
				binary.NativeEndian.PutUint64(values[i].Value[:], math.Float64bits(38))
			}
			binary.NativeEndian.PutUint64(values[deviceTemperature].Value[:], 42)
			got, err := deviceStatusFromValues(values)
			require.Equal(t, tt.wantFloat, got.Power)
			require.Equal(t, tt.wantInt, got.Utilization.Encoder)
			require.Equal(t, int64(42), got.Temperature)
			if tt.wantError {
				var native *Error
				require.ErrorAs(t, err, &native)
				require.Equal(t, tt.status, int(native.Code))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPCIeBandwidth(t *testing.T) {
	tests := []struct {
		name  string
		gen   int64
		width int64
		want  int64
	}{
		{name: "generation 1", gen: 1, width: 16, want: 4000},
		{name: "generation 2", gen: 2, width: 16, want: 8000},
		{name: "generation 3", gen: 3, width: 16, want: 15760},
		{name: "generation 4", gen: 4, width: 16, want: 31504},
		{name: "unknown generation", gen: 5, width: 16},
		{name: "negative generation", gen: -1, width: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pcieBandwidth(tt.gen, tt.width); got != tt.want {
				t.Fatalf("pcieBandwidth(%d, %d) = %d, want %d", tt.gen, tt.width, got, tt.want)
			}
		})
	}
}

func testInt64FieldValue(value int64) FieldValue_v1 {
	field := FieldValue_v1{FieldType: DCGM_FT_INT64}
	copy(field.Value[:], int64Bytes(value))
	return field
}

func testFloat64FieldValue(value float64) FieldValue_v1 {
	field := FieldValue_v1{FieldType: DCGM_FT_DOUBLE}
	binary.NativeEndian.PutUint64(field.Value[:], math.Float64bits(value))
	return field
}
