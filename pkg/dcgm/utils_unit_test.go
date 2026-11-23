package dcgm

import (
	"math"
	"testing"
)

func TestIsInt32BlankBoundary(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"below", int(DCGM_FT_INT32_BLANK) - 1, false},
		{"blank", int(DCGM_FT_INT32_BLANK), true},
		{"above", int(DCGM_FT_INT32_BLANK) + 1, true},
		{"past named markers", int(DCGM_FT_INT32_BLANK) + 4, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt32Blank(tt.value); got != tt.want {
				t.Fatalf("IsInt32Blank(%d) = %t, want %t", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsInt64BlankBoundary(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  bool
	}{
		{"below", DCGM_FT_INT64_BLANK - 1, false},
		{"blank", DCGM_FT_INT64_BLANK, true},
		{"above", DCGM_FT_INT64_BLANK + 1, true},
		{"past named markers", DCGM_FT_INT64_BLANK + 4, true},
		{"32-bit marker", DCGM_FT_INT32_BLANK, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt64Blank(tt.value); got != tt.want {
				t.Fatalf("IsInt64Blank(%d) = %t, want %t", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsStringBlank(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value string
		want  bool
	}{
		{"blank", DCGM_FT_STR_BLANK, true},
		{"not found", DCGM_FT_STR_NOT_FOUND, true},
		{"not supported", DCGM_FT_STR_NOT_SUPPORTED, true},
		{"not permissioned", DCGM_FT_STR_NOT_PERMISSIONED, true},
		{"unknown marker", "<<<FUTURE>>>", true},
		{"trailing text", "<<<FUTURE>>>detail", true},
		{"empty", "", false},
		{"ordinary value", "GPU-1234", false},
		{"later marker", "value<<<NULL>>>", false},
		{"missing end", "<<<NULL", false},
		{"NUL before end", "<<<NULL\x00>>>", false},
		{"NUL after end", "<<<NULL>>>\x00detail", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStringBlank(tt.value); got != tt.want {
				t.Errorf("IsStringBlank(%q) = %t, want %t", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsFloat64Blank(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value float64
		want  bool
	}{
		{"below", math.Nextafter(DCGM_FT_FP64_BLANK, math.Inf(-1)), false},
		{"blank", DCGM_FT_FP64_BLANK, true},
		{"not found", DCGM_FT_FP64_NOT_FOUND, true},
		{"not supported", DCGM_FT_FP64_NOT_SUPPORTED, true},
		{"not permissioned", DCGM_FT_FP64_NOT_PERMISSIONED, true},
		{"past named markers", DCGM_FT_FP64_BLANK + 4, true},
		{"zero", 0, false},
		{"NaN", math.NaN(), false},
		{"positive infinity", math.Inf(1), true},
		{"negative infinity", math.Inf(-1), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFloat64Blank(tt.value); got != tt.want {
				t.Errorf("IsFloat64Blank(%g) = %t, want %t", tt.value, got, tt.want)
			}
		})
	}
}

func TestMakeVersion(t *testing.T) {
	tests := []struct {
		name  string
		build func(uintptr) uint
		want  uint
	}{
		{"v1", func(n uintptr) uint { return uint(makeVersion1(n)) }, 0x01001234},
		{"v2", func(n uintptr) uint { return uint(makeVersion2(n)) }, 0x02001234},
		{"v3", func(n uintptr) uint { return uint(makeVersion3(n)) }, 0x03001234},
		{"v4", func(n uintptr) uint { return uint(makeVersion4(n)) }, 0x04001234},
		{"v5", func(n uintptr) uint { return uint(makeVersion5(n)) }, 0x05001234},
		{"v12", func(n uintptr) uint { return uint(makeVersion12(n)) }, 0x0c001234},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.build(0x1234); got != tt.want {
				t.Fatalf("%s(0x1234) = %#x, want %#x", tt.name, got, tt.want)
			}
		})
	}
}

func TestRoundFloat(t *testing.T) {
	positiveDown := 1.4
	positiveHalf := 1.5
	negativeDown := -1.4
	negativeHalf := -1.5
	tests := []struct {
		name  string
		input *float64
		want  float64
	}{
		{"nil", nil, 0},
		{"positive down", &positiveDown, 1},
		{"positive half", &positiveHalf, 2},
		{"negative down", &negativeDown, -1},
		{"negative half", &negativeHalf, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roundFloat(tt.input)
			if got == nil || *got != tt.want {
				t.Fatalf("roundFloat(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
