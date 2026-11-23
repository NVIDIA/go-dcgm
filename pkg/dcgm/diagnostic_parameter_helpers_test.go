package dcgm

// diagnosticParameterRejected reports whether the plugin rejected a parameter with the expected message.
func diagnosticParameterRejected(result DiagResults, err error, message string) bool {
	if err != nil {
		return false
	}
	for i := range result.Tests {
		test := &result.Tests[i]
		if test.Name != "diagnostic" || test.Status != "fail" {
			continue
		}
		for _, testError := range test.Errors {
			if testError.Code == uint(DCGM_FR_INTERNAL) &&
				testError.Message == "Caught an exception while trying to execute the test: "+message {
				return true
			}
		}
	}
	return false
}
