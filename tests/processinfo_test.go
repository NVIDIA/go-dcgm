package tests

import (
	"os"
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestProcessInfo demonstrates getting process information for GPU processes
// This is equivalent to the processInfo sample
func TestProcessInfo(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Request DCGM to start recording stats for GPU process fields
	group, err := dcgm.WatchPidFields()
	if err != nil {
		t.Fatalf("Failed to watch PID fields: %v", err)
	}

	// Wait for watches to be enabled and collect data
	t.Log("Enabling DCGM watches to start collecting process stats. This may take a few seconds...")
	time.Sleep(3000 * time.Millisecond)

	// Get current process ID as an example
	//nolint:gosec // disable G115
	currentPid := uint(os.Getpid())
	t.Logf("Testing with current process PID: %d", currentPid)

	pidInfo, err := dcgm.GetProcessInfo(group, currentPid)
	if err != nil {
		t.Logf("Failed to get process info for PID %d: %v", currentPid, err)
		t.Log("This is expected if the current process is not using GPU")
		return
	}

	if len(pidInfo) == 0 {
		t.Logf("No process information found for PID %d", currentPid)
		return
	}

	// Log basic process information
	for i, info := range pidInfo {
		t.Logf("Process Info %d:", i+1)
		t.Logf("  GPU ID: %d", info.GPU)
		t.Logf("  PID: %d", info.PID)
		if info.Name != "" {
			t.Logf("  Name: %s", info.Name)
		}
		t.Logf("  Start Time: %s", info.ProcessUtilization.StartTime.String())
		t.Logf("  End Time: %s", info.ProcessUtilization.EndTime.String())
		t.Logf("  Critical XID Errors: %d", info.XIDErrors.NumErrors)
	}
}

// TestProcessInfoWithSpecificPID demonstrates getting process info for a specific PID
func TestProcessInfoWithSpecificPID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping specific PID test in short mode")
	}

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Request DCGM to start recording stats for GPU process fields
	group, err := dcgm.WatchPidFields()
	if err != nil {
		t.Fatalf("Failed to watch PID fields: %v", err)
	}

	// Wait for watches to be enabled and collect data
	time.Sleep(3000 * time.Millisecond)

	// Test with PID 1 (init process) - should not have GPU usage
	testPid := uint(1)
	pidInfo, err := dcgm.GetProcessInfo(group, testPid)
	if err != nil {
		t.Logf("Expected: No process info found for PID %d: %v", testPid, err)
	} else if len(pidInfo) == 0 {
		t.Logf("Expected: No GPU usage found for PID %d", testPid)
	} else {
		t.Logf("Unexpected: Found GPU usage for PID %d", testPid)
	}
}

// TestWatchPidFields demonstrates the WatchPidFields functionality
func TestWatchPidFields(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Test WatchPidFields function
	group, err := dcgm.WatchPidFields()
	if err != nil {
		t.Fatalf("Failed to watch PID fields: %v", err)
	}

	t.Logf("Successfully created group for watching PID fields: %v", group)

	// Wait a bit to ensure watches are properly set up
	time.Sleep(1000 * time.Millisecond)
	t.Log("PID field watches enabled successfully")
}
