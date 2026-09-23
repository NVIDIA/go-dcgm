package dcgm

import (
	"encoding/binary"
	"math"
	"testing"
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
	values := make([]FieldValue_v1, deviceStatusFieldCount)
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

	tests := []struct {
		name   string
		values []FieldValue_v1
		want   DeviceStatus
	}{
		{
			name:   "zero values",
			values: make([]FieldValue_v1, deviceStatusFieldCount),
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
			if got := deviceStatusFromValues(tt.values); got != tt.want {
				t.Fatalf("deviceStatusFromValues() = %+v, want %+v", got, tt.want)
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
	var field FieldValue_v1
	copy(field.Value[:], int64Bytes(value))
	return field
}

func testFloat64FieldValue(value float64) FieldValue_v1 {
	var field FieldValue_v1
	binary.NativeEndian.PutUint64(field.Value[:], math.Float64bits(value))
	return field
}
