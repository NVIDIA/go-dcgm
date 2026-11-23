package dcgm

import (
	"errors"
	"testing"
)

func TestDiagnosticParameterRejection(t *testing.T) {
	const invalid = "Invalid numeric test parameter"
	const overflow = "Numeric test parameter is outside the double range"
	const invalidException = "Caught an exception while trying to execute the test: Invalid numeric test parameter"
	const overflowException = "Caught an exception while trying to execute the test: Numeric test parameter is outside the double range"
	pluginResult := func(name, status string, code uint, message string) DiagResults {
		return DiagResults{Tests: []DiagTest{{
			Name:   name,
			Status: status,
			Errors: []DiagError{{Code: code, Message: message}},
		}}}
	}
	invalidResult := pluginResult("diagnostic", "fail", uint(DCGM_FR_INTERNAL), invalidException)
	overflowResult := pluginResult("diagnostic", "fail", uint(DCGM_FR_INTERNAL), overflowException)
	for _, tc := range []struct {
		name     string
		result   DiagResults
		err      error
		message  string
		rejected bool
	}{
		{"native plugin response", invalidResult, nil, invalid, true},
		{"native overflow response", overflowResult, nil, overflow, true},
		{"accepted input", DiagResults{}, nil, invalid, false},
		{"transport failure", invalidResult, errors.New("connection timed out"), invalid, false},
		{"native API failure", invalidResult, &Error{Code: DCGM_ST_NVVS_ERROR}, invalid, false},
		{"unrelated system failure", DiagResults{SystemErrors: []DiagError{{Code: uint(DCGM_FR_INTERNAL), Message: "GPU is lost"}}}, nil, invalid, false},
		{"system error is not plugin evidence", DiagResults{SystemErrors: []DiagError{{Code: uint(DCGM_FR_INTERNAL), Message: invalidException}}}, nil, invalid, false},
		{"wrong plugin", pluginResult("software", "fail", uint(DCGM_FR_INTERNAL), invalidException), nil, invalid, false},
		{"plugin passed", pluginResult("diagnostic", "pass", uint(DCGM_FR_INTERNAL), invalidException), nil, invalid, false},
		{"wrong error code", pluginResult("diagnostic", "fail", 0, invalidException), nil, invalid, false},
		{"wrong parser failure", overflowResult, nil, invalid, false},
		{"unrelated plugin error", pluginResult("diagnostic", "fail", uint(DCGM_FR_INTERNAL), "GPU is lost"), nil, invalid, false},
		{"missing exception context", pluginResult("diagnostic", "fail", uint(DCGM_FR_INTERNAL), invalid), nil, invalid, false},
		{"message substring is insufficient", pluginResult("diagnostic", "fail", uint(DCGM_FR_INTERNAL), invalidException+" (unrelated)"), nil, invalid, false},
		{"GPU test failure", DiagResults{Tests: []DiagTest{{Name: "diagnostic", Status: "fail"}}}, nil, invalid, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := diagnosticParameterRejected(tc.result, tc.err, tc.message); got != tc.rejected {
				t.Fatalf("rejected=%t, want %t", got, tc.rejected)
			}
		})
	}
}
