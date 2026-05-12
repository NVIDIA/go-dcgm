package dcgm

import (
	"os"
	"regexp"
	"strconv"
	"testing"
)

func TestSelectedStatusCodeConstantsMatchCHeader(t *testing.T) {
	tests := []struct {
		name string
		got  int
	}{
		{"DCGM_ST_OK", DCGM_ST_OK},
		{"DCGM_ST_BADPARAM", DCGM_ST_BADPARAM},
		{"DCGM_ST_PENDING", DCGM_ST_PENDING},
		{"DCGM_ST_UNINITIALIZED", DCGM_ST_UNINITIALIZED},
		{"DCGM_ST_TIMEOUT", DCGM_ST_TIMEOUT},
		{"DCGM_ST_VER_MISMATCH", DCGM_ST_VER_MISMATCH},
		{"DCGM_ST_FUNCTION_NOT_FOUND", DCGM_ST_FUNCTION_NOT_FOUND},
		{"DCGM_ST_CONNECTION_NOT_VALID", DCGM_ST_CONNECTION_NOT_VALID},
		{"DCGM_ST_LIBRARY_NOT_FOUND", DCGM_ST_LIBRARY_NOT_FOUND},
		{"DCGM_ST_INIT_ERROR", DCGM_ST_INIT_ERROR},
		{"DCGM_ST_NVML_ERROR", DCGM_ST_NVML_ERROR},
		{"DCGM_ST_NVML_NOT_LOADED", DCGM_ST_NVML_NOT_LOADED},
	}

	headerStatusCodes := dcgmStructsStatusCodes(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, ok := headerStatusCodes[tt.name]
			if !ok {
				t.Fatalf("%s was not found in dcgm_structs.h", tt.name)
			}
			if tt.got != want {
				t.Fatalf("%s = %d, want %d", tt.name, tt.got, want)
			}
		})
	}
}

func dcgmStructsStatusCodes(t *testing.T) map[string]int {
	t.Helper()

	content, err := os.ReadFile("dcgm_structs.h")
	if err != nil {
		t.Fatalf("read dcgm_structs.h: %v", err)
	}

	statusCodePattern := regexp.MustCompile(`(?m)^\s*(DCGM_ST_[A-Z0-9_]+)\s*=\s*(-?\d+),`)
	matches := statusCodePattern.FindAllStringSubmatch(string(content), -1)
	if len(matches) == 0 {
		t.Fatal("no DCGM_ST_* status codes found in dcgm_structs.h")
	}

	statusCodes := make(map[string]int, len(matches))
	for _, match := range matches {
		value, err := strconv.Atoi(match[2])
		if err != nil {
			t.Fatalf("parse %s value %q: %v", match[1], match[2], err)
		}
		statusCodes[match[1]] = value
	}

	return statusCodes
}
