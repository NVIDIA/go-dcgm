package tests

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/stretchr/testify/assert"
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

	requireDiagnosticResults(t, dr)
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
		t.Fatal("No supported GPUs found for diagnostics")
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

	requireDiagnosticResults(t, dr)
}

// TestDiagTestNameFormat checks that DCGM returns short, lowercase test names.
func TestDiagTestNameFormat(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	dr, err := dcgm.RunDiag(dcgm.DiagQuick, dcgm.GroupAllGPUs())
	if err != nil {
		t.Fatalf("Failed to run diagnostics: %v", err)
	}

	requireDiagnosticResults(t, dr)

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

// requireDiagnosticResults logs the results from DCGM and fails unless at least one diagnostic passed.
func requireDiagnosticResults(t *testing.T, result dcgm.DiagResults) {
	t.Helper()
	for i := range result.Tests {
		test := &result.Tests[i]
		t.Logf("Diagnostic %s: %s; info=%+v; errors=%+v", test.Name, test.Status, test.Info, test.Errors)
	}

	if err := diagnosticResultsError(result.Tests); err != nil {
		t.Fatal(err)
	}
}

// diagnosticResultsError reports failures and requires at least one passing result.
func diagnosticResultsError(tests []dcgm.DiagTest) error {
	passed := false
	for i := range tests {
		test := &tests[i]
		switch test.Status {
		case "fail":
			return fmt.Errorf("diagnostic test failed: %s", test.Name)
		case "pass":
			passed = true
		}
	}

	if !passed {
		return errors.New("DCGM returned no passed diagnostic tests")
	}
	return nil
}

// TestDiagnosticResultAcceptance checks which diagnostic results pass or fail the suite.
func TestDiagnosticResultAcceptance(t *testing.T) {
	for _, tc := range []struct {
		name      string
		statuses  []string
		wantError bool
	}{
		{name: "empty", wantError: true},
		{name: "skipped only", statuses: []string{"skipped", "notrun"}, wantError: true},
		{name: "warning only", statuses: []string{"warn"}, wantError: true},
		{name: "passed", statuses: []string{"pass"}},
		{name: "passed with optional skip", statuses: []string{"pass", "skipped", "notrun"}},
		{name: "failed with pass", statuses: []string{"pass", "fail"}, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			results := make([]dcgm.DiagTest, len(tc.statuses))
			for i, status := range tc.statuses {
				results[i] = dcgm.DiagTest{Name: fmt.Sprintf("diagnostic-%d", i), Status: status}
			}

			err := diagnosticResultsError(results)
			if (err != nil) != tc.wantError {
				t.Fatalf("diagnosticResultsError() = %v, want error=%v", err, tc.wantError)
			}
		})
	}
}
