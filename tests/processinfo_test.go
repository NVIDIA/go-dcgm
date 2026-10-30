package tests

import (
	"errors"
	"fmt"
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
	currentPid := uint(os.Getpid()) // #nosec G115 -- process IDs are non-negative.
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

// newPidWatchTestGroup returns a caller-owned group holding the first supported
// GPU. The test destroys it at cleanup and fails if that destroy fails, which
// also shows the watch call left the group in place.
func newPidWatchTestGroup(t *testing.T) dcgm.GroupHandle {
	t.Helper()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}
	if len(gpus) == 0 {
		t.Skip("No supported GPUs found")
	}

	group, err := dcgm.CreateGroup(fmt.Sprintf("pid-watch-%d", time.Now().UnixNano()))
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}
	t.Cleanup(func() {
		if err := dcgm.DestroyGroup(group); err != nil {
			t.Errorf("Failed to destroy group: %v", err)
		}
	})

	if err := dcgm.AddToGroup(group, gpus[0]); err != nil {
		t.Fatalf("Failed to add GPU %d to group: %v", gpus[0], err)
	}

	return group
}

// TestWatchPidFieldsForGroup checks that PID watches start on a group the caller
// created, and that GetProcessInfo accepts that group afterward.
func TestWatchPidFieldsForGroup(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	t.Cleanup(cleanup) // runs after the group cleanup registered later

	group := newPidWatchTestGroup(t)

	if err = dcgm.WatchPidFieldsForGroup(group); err != nil {
		skipIfPidWatchRequiresRoot(t, err)
		t.Fatalf("WatchPidFieldsForGroup failed: %v", err)
	}

	// This process does not use the GPU, so DCGM has no record of it. A
	// missing PID is reported as DCGM_ST_NO_DATA; any other error means the
	// group was not usable.
	_, err = dcgm.GetProcessInfo(group, uint(os.Getpid())) // #nosec G115 -- process IDs are non-negative.
	var dcgmErr *dcgm.Error
	noRecord := errors.As(err, &dcgmErr) && int(dcgmErr.Code) == dcgm.DCGM_ST_NO_DATA
	if err != nil && !noRecord {
		t.Fatalf("GetProcessInfo on the watched group failed: %v", err)
	}
}
