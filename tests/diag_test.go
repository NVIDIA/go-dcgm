package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/stretchr/testify/assert"
)

func skipIfDiagnosticsDisabled(t *testing.T) {
	t.Helper()
	if os.Getenv("DCGM_SKIP_DIAGNOSTICS") == "1" {
		t.Skip("diagnostics disabled by DCGM_SKIP_DIAGNOSTICS")
	}
}

// TestDiagnostics demonstrates running DCGM diagnostics
// This is equivalent to the diag sample
func TestDiagnostics(t *testing.T) {
	skipIfDiagnosticsDisabled(t)

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

	t.Log("Diagnostic Tests:")
	for _, test := range dr.Tests {
		t.Logf("  %-50s %s", test.Name, test.Status)
	}

	if len(dr.Tests) == 0 {
		t.Error("No diagnostic results returned")
	}

	failedTests := 0
	for _, test := range dr.Tests {
		if test.Status == "fail" {
			failedTests++
			t.Logf("Diagnostic test failed: %s", test.Name)
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
	skipIfDiagnosticsDisabled(t)

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

	for _, test := range dr.Tests {
		t.Logf("  %s: %s", test.Name, test.Status)
	}
}

// TestDiagTestNameFormat validates that Name contains a test name,
// not detailed test descriptions (issue #97)
func TestDiagTestNameFormat(t *testing.T) {
	skipIfDiagnosticsDisabled(t)

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	dr, err := dcgm.RunDiag(dcgm.DiagQuick, dcgm.GroupAllGPUs())
	if err != nil {
		t.Fatalf("Failed to run diagnostics: %v", err)
	}

	assert.NotEmpty(t, dr.Tests, "diagnostic results should not be empty")

	// Invalid strings that should NOT appear in TestName
	// These are detailed descriptions that were incorrectly returned before fix
	invalidPatterns := []string{
		"presence of drivers on the denylist",
		"(e.g. nouveau)",
		"Allocated",
		"bytes",
		"presence (and version)",
	}

	for i, test := range dr.Tests {
		t.Logf("Result %d: Name=%q, Status=%s", i, test.Name, test.Status)
		assert.NotEmpty(t, test.Name, "diagnostic test name should not be empty")

		// TestName should NOT contain detailed descriptions
		for _, invalid := range invalidPatterns {
			assert.NotContains(
				t,
				test.Name,
				invalid,
				"Name should not contain detailed descriptions, got: %q",
				test.Name,
			)
		}

		// TestName should be lowercase
		assert.Equal(
			t,
			strings.ToLower(test.Name),
			test.Name,
			"Name should be lowercase, got: %q",
			test.Name,
		)
	}
}
