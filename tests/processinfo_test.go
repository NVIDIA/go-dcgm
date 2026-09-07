package tests

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func skipIfPidWatchRequiresRoot(t *testing.T, err error) {
	t.Helper()

	var dcgmErr *dcgm.Error
	if errors.As(err, &dcgmErr) && int(dcgmErr.Code) == dcgm.DCGM_ST_REQUIRES_ROOT {
		t.Skipf("PID field watches require root on this DCGM host: %v", err)
	}
}

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
		skipIfPidWatchRequiresRoot(t, err)
		t.Fatalf("Failed to watch PID fields: %v", err)
	}

	// Wait for watches to be enabled and collect data
	t.Log("Enabling DCGM watches to start collecting process stats. This may take a few seconds...")
	time.Sleep(3000 * time.Millisecond)

	// Get current process ID as an example
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
		skipIfPidWatchRequiresRoot(t, err)
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
		skipIfPidWatchRequiresRoot(t, err)
		t.Fatalf("Failed to watch PID fields: %v", err)
	}

	t.Logf("Successfully created group for watching PID fields: %v", group)

	// Wait a bit to ensure watches are properly set up
	time.Sleep(1000 * time.Millisecond)
	t.Log("PID field watches enabled successfully")
}

// TestWatchPidFieldsForGroup demonstrates the WatchPidFieldsForGroup functionality
func TestWatchPidFieldsForGroup(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Create a group first
	group, err := dcgm.CreateGroup("test-group")
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}
	defer func() {
		if err := dcgm.DestroyGroup(group); err != nil {
			t.Logf("Warning: failed to destroy group: %v", err)
		}
	}()

	// Add at least one supported GPU to the group
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}
	if len(gpus) == 0 {
		t.Skip("No supported GPUs found")
	}
	err = dcgm.AddToGroup(group, gpus[0])
	if err != nil {
		t.Fatalf("Failed to add GPU to group: %v", err)
	}

	// Test WatchPidFieldsForGroup function
	err = dcgm.WatchPidFieldsForGroup(group)
	if err != nil {
		skipIfPidWatchRequiresRoot(t, err)
		t.Fatalf("Failed to watch PID fields for group: %v", err)
	}

	t.Logf("Successfully created PID field watcher for group: %v", group)

	// Wait a bit to ensure watches are properly set up
	time.Sleep(1000 * time.Millisecond)
	t.Log("PID field watches enabled successfully")
}

// TestWatchPidFieldsForGroupEx demonstrates the WatchPidFieldsForGroupEx functionality
func TestWatchPidFieldsForGroupEx(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Create a group first
	group, err := dcgm.CreateGroup("test-group-ex")
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}
	defer func() {
		if err := dcgm.DestroyGroup(group); err != nil {
			t.Logf("Warning: failed to destroy group: %v", err)
		}
	}()

	// Add at least one supported GPU to the group
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}
	if len(gpus) == 0 {
		t.Skip("No supported GPUs found")
	}
	err = dcgm.AddToGroup(group, gpus[0])
	if err != nil {
		t.Fatalf("Failed to add GPU to group: %v", err)
	}

	// Test WatchPidFieldsForGroupEx function with custom parameters
	err = dcgm.WatchPidFieldsForGroupEx(group, time.Microsecond*1000000, time.Second*60, 5)
	if err != nil {
		skipIfPidWatchRequiresRoot(t, err)
		t.Fatalf("Failed to watch PID fields for group with custom params: %v", err)
	}

	t.Logf("Successfully created PID field watcher for group with custom params: %v", group)

	// Wait a bit to ensure watches are properly set up
	time.Sleep(1000 * time.Millisecond)
	t.Log("PID field watches enabled successfully with custom params")
}
