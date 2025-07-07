package tests

import (
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestDiagnostics demonstrates running DCGM diagnostics
// This is equivalent to the diag sample
func TestDiagnostics(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Run quick diagnostics on all GPUs
	dr, err := dcgm.RunDiag(dcgm.DiagQuick, dcgm.GroupAllGPUs())
	if err != nil {
		t.Fatalf("Failed to run diagnostics: %v", err)
	}

	// Log software test results
	t.Logf("Software Tests:")
	for _, test := range dr.Software {
		t.Logf("  %-50s %s\t%s", test.TestName, test.Status, test.TestOutput)
	}

	// Basic validation - we should have some results
	if len(dr.Software) == 0 {
		t.Error("No diagnostic results returned")
	}

	// Check for any failed tests
	failedTests := 0
	for _, test := range dr.Software {
		if test.Status == "fail" {
			failedTests++
			t.Logf("Software test failed: %s - %s", test.TestName, test.TestOutput)
		}
	}

	if failedTests > 0 {
		t.Logf("Total failed tests: %d", failedTests)
	} else {
		t.Log("All diagnostic tests passed")
	}
}

// TestDiagnosticsLong demonstrates running longer diagnostics
func TestDiagnosticsLong(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long diagnostics test in short mode")
	}

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Get supported devices first
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}

	if len(gpus) == 0 {
		t.Skip("No supported GPUs found for diagnostics")
	}

	// Run diagnostics on first GPU only for time efficiency
	group, err := dcgm.CreateGroup("test-group")
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}
	defer func() {
		if err = dcgm.DestroyGroup(group); err != nil {
			t.Logf("Failed to destroy group: %v", err)
		}
	}()

	err = dcgm.AddToGroup(group, gpus[0])
	if err != nil {
		t.Fatalf("Failed to add GPU to group: %v", err)
	}

	// Run medium-level diagnostics
	dr, err := dcgm.RunDiag(dcgm.DiagMedium, group)
	if err != nil {
		t.Fatalf("Failed to run medium diagnostics: %v", err)
	}

	t.Logf("Medium diagnostics completed for GPU %d", gpus[0])

	// Log results
	for _, test := range dr.Software {
		t.Logf("  %s: %s", test.TestName, test.Status)
	}
}
