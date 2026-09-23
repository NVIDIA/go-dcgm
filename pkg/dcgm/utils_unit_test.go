package dcgm

import "testing"

func TestIsInt32BlankBoundary(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"below", dcgmInt32Blank - 1, false},
		{"blank", dcgmInt32Blank, true},
		{"above", dcgmInt32Blank + 1, true},
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
		{"below", dcgmInt64Blank - 1, false},
		{"blank", dcgmInt64Blank, true},
		{"above", dcgmInt64Blank + 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt64Blank(tt.value); got != tt.want {
				t.Fatalf("IsInt64Blank(%d) = %t, want %t", tt.value, got, tt.want)
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
