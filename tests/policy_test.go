package tests

import (
	"context"
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestPolicyViolations demonstrates listening for policy violations
// This is equivalent to the policy sample but runs for a limited time
func TestPolicyViolations(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Listen for policy violations (DBE and XID errors)
	c, err := dcgm.ListenForPolicyViolations(ctx, dcgm.DbePolicy, dcgm.XidPolicy)
	if err != nil {
		t.Fatalf("Failed to start listening for policy violations: %v", err)
	}

	t.Log("Listening for policy violations (DBE and XID errors) for 10 seconds...")

	violationCount := 0
	timeout := time.After(10 * time.Second)

	for {
		select {
		case pe := <-c:
			violationCount++
			t.Logf("Policy Violation %d:", violationCount)
			t.Logf("  Condition: %v", pe.Condition)
			t.Logf("  Timestamp: %v", pe.Timestamp)
			t.Logf("  Data: %v", pe.Data)

		case <-ctx.Done():
			t.Logf("Policy violation monitoring completed")
			t.Logf("Total violations detected: %d", violationCount)
			return

		case <-timeout:
			t.Logf("Policy violation monitoring timed out")
			t.Logf("Total violations detected: %d", violationCount)
			return
		}
	}
}

// TestPolicyViolationsSingleType demonstrates listening for a specific type of policy violation
func TestPolicyViolationsSingleType(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping single type policy test in short mode")
	}

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Listen for only XID policy violations
	c, err := dcgm.ListenForPolicyViolations(ctx, dcgm.XidPolicy)
	if err != nil {
		t.Fatalf("Failed to start listening for XID policy violations: %v", err)
	}

	t.Log("Listening for XID policy violations for 5 seconds...")

	xidCount := 0
	timeout := time.After(5 * time.Second)

	for {
		select {
		case pe := <-c:
			xidCount++
			t.Logf("XID Policy Violation %d:", xidCount)
			t.Logf("  Condition: %v", pe.Condition)
			t.Logf("  Timestamp: %v", pe.Timestamp)
			t.Logf("  Data: %v", pe.Data)

		case <-ctx.Done():
			t.Logf("XID policy violation monitoring completed")
			t.Logf("Total XID violations detected: %d", xidCount)
			return

		case <-timeout:
			t.Logf("XID policy violation monitoring timed out")
			t.Logf("Total XID violations detected: %d", xidCount)
			return
		}
	}
}

// TestPolicyViolationsMultipleTypes demonstrates listening for multiple types of policy violations
func TestPolicyViolationsMultipleTypes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping multiple types policy test in short mode")
	}

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	// Listen for multiple types of policy violations
	// Note: Some policies may require root privileges

	c, err := dcgm.ListenForPolicyViolations(ctx, dcgm.DbePolicy, dcgm.XidPolicy, dcgm.ThermalPolicy, dcgm.PowerPolicy)
	if err != nil {
		t.Logf("Failed to start listening for all policy violations (may need root): %v", err)
		// Try with just basic policies
		c, err = dcgm.ListenForPolicyViolations(ctx, dcgm.DbePolicy, dcgm.XidPolicy)
		if err != nil {
			t.Fatalf("Failed to start listening for basic policy violations: %v", err)
		}
		t.Log("Listening for basic policy violations (DBE and XID) for 8 seconds...")
	} else {
		t.Log("Listening for multiple policy violations for 8 seconds...")
	}

	violationsByType := make(map[string]int)
	timeout := time.After(8 * time.Second)

	for {
		select {
		case pe := <-c:
			conditionStr := string(pe.Condition)
			violationsByType[conditionStr]++

			t.Logf("Policy Violation:")
			t.Logf("  Type: %s", conditionStr)
			t.Logf("  Timestamp: %v", pe.Timestamp)
			t.Logf("  Data: %v", pe.Data)

		case <-ctx.Done():
			t.Log("Multi-type policy violation monitoring completed")
			for policyType, count := range violationsByType {
				t.Logf("  %s violations: %d", policyType, count)
			}
			return

		case <-timeout:
			t.Log("Multi-type policy violation monitoring timed out")
			for policyType, count := range violationsByType {
				t.Logf("  %s violations: %d", policyType, count)
			}
			return
		}
	}
}

// TestPolicyViolationsContextCancellation demonstrates proper context cancellation
func TestPolicyViolationsContextCancellation(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())

	c, err := dcgm.ListenForPolicyViolations(ctx, dcgm.DbePolicy)
	if err != nil {
		t.Fatalf("Failed to start listening for policy violations: %v", err)
	}

	t.Log("Starting policy violation monitoring, will cancel after 2 seconds...")

	// Cancel after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		t.Log("Cancelling policy violation monitoring...")
		cancel()
	}()

	violationCount := 0
	startTime := time.Now()

	for {
		select {
		case pe := <-c:
			violationCount++
			t.Logf("Policy violation %d: %v", violationCount, pe.Condition)

		case <-ctx.Done():
			elapsed := time.Since(startTime)
			t.Logf("Policy violation monitoring stopped after %v", elapsed)
			t.Logf("Total violations detected: %d", violationCount)

			// Should have stopped within reasonable time after cancellation
			if elapsed > 3*time.Second {
				t.Errorf("Context cancellation took too long: %v", elapsed)
			}
			return
		}
	}
}
