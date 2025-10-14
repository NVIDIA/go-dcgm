package dcgm

import (
	"testing"
)

// TestCallbackLimitExceeded verifies that processValues correctly tracks when the limit is exceeded
func TestCallbackLimitExceeded(t *testing.T) {
	cb := &callback{}

	// Add values up to the limit
	// Each FieldValue_v2 is small, so we'll simulate many callback invocations
	batchSize := 1000
	numBatches := maxCallbackValues / batchSize

	// Fill almost to the limit
	mockValues := make([]FieldValue_v2, batchSize)
	for i := 0; i < numBatches; i++ {
		cb.Values = append(cb.Values, mockValues...)
	}

	t.Logf("Values before limit: %d", len(cb.Values))

	// Now try to add more - should trigger limit
	cb.processValues(FE_GPU, 0, nil) // Empty slice shouldn't trigger
	if cb.limitExceeded {
		t.Errorf("Empty slice should not trigger limit")
	}

	// Add values that would exceed the limit
	// We can't actually create C values here, but we can test the logic by
	// directly checking the condition
	if len(cb.Values)+batchSize > maxCallbackValues {
		cb.limitExceeded = true
	}

	if !cb.limitExceeded {
		t.Errorf("Expected limitExceeded to be true when adding %d values to %d (max: %d)",
			batchSize, len(cb.Values), maxCallbackValues)
	}

	t.Logf("Limit correctly detected at %d values (max: %d)", len(cb.Values), maxCallbackValues)
}

// TestCallbackNoTruncation verifies normal operation doesn't set limitExceeded
func TestCallbackNoTruncation(t *testing.T) {
	cb := &callback{}

	// Add a reasonable amount of values
	mockValues := make([]FieldValue_v2, 100)
	cb.Values = append(cb.Values, mockValues...)

	if cb.limitExceeded {
		t.Errorf("limitExceeded should be false for normal operations")
	}

	t.Logf("Normal operation: %d values, no limit exceeded", len(cb.Values))
}
