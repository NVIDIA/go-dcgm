package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

// TestProcessInfoWithGPUFixture requires a live GPU process record, rather than
// accepting the missing records used by the process-info examples.
func TestProcessInfoWithGPUFixture(t *testing.T) {
	load := jobStatsLoadExecutable(t)
	group, device := jobStatsSession(t)
	if mode := jobNvidiaSMI(t, "--id="+device.UUID, "--query-gpu=accounting.mode", "--format=csv,noheader,nounits"); mode == "N/A" || mode == "[N/A]" {
		t.Skip("GPU does not support process accounting")
	}
	if err := dcgm.WatchPidFieldsForGroupEx(group, 100*time.Millisecond, time.Minute, 0); err != nil {
		skipIfPidWatchRequiresRoot(t, err)
		var nativeErr *dcgm.Error
		if errors.As(err, &nativeErr) && int(nativeErr.Code) == dcgm.DCGM_ST_NOT_SUPPORTED {
			t.Skipf("GPU does not support PID field watches: %v", err)
		}
		t.Fatalf("watch PID fields: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, load) // #nosec G204,G702 -- explicitly supplied test-only executable, no shell.
	cmd.Env = append(os.Environ(), "DCGM_JOB_GPU_ID="+strconv.FormatUint(uint64(device.GPU), 10), "DCGM_JOB_GPU_UUID="+device.UUID, "DCGM_JOB_LOAD_SECONDS=10")
	cmd.WaitDelay = time.Second
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	started := time.Now()
	if err := cmd.Start(); err != nil {
		t.Fatalf("start GPU fixture: %v", err)
	}
	var waitErr error
	done := make(chan struct{})
	go func() {
		waitErr = cmd.Wait()
		close(done)
	}()
	t.Cleanup(func() {
		cancel()
		<-done
		if t.Failed() {
			t.Logf("GPU fixture exit: %v; output: %s", waitErr, output.String())
		}
	})

	pid := uint(cmd.Process.Pid) // #nosec G115 -- started process IDs are positive.
	waitForLiveProcessInfo(t, ctx, done, group, dcgm.ProcessInfo{GPU: device.GPU, PID: pid}, started)
	<-done
	if waitErr != nil || ctx.Err() != nil {
		t.Fatalf("GPU fixture did not exit successfully: %v (context: %v): %s", waitErr, ctx.Err(), output.String())
	}
}

func waitForLiveProcessInfo(t *testing.T, ctx context.Context, done <-chan struct{}, group dcgm.GroupHandle, want dcgm.ProcessInfo, started time.Time) {
	t.Helper()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			t.Fatal("GPU fixture exited before a live PID record was observed")
		case <-ctx.Done():
			t.Fatalf("no live GPU record for PID %d before deadline: %v", want.PID, ctx.Err())
		case <-ticker.C:
		}
		info, err := dcgm.GetProcessInfo(group, want.PID)
		if err != nil {
			var nativeErr *dcgm.Error
			if errors.As(err, &nativeErr) && int(nativeErr.Code) == dcgm.DCGM_ST_NO_DATA {
				continue // Wait for DCGM to observe the fixture's CUDA context.
			}
			t.Fatalf("query live PID %d: %v", want.PID, err)
		}
		if len(info) == 0 {
			continue
		}
		checkLiveProcessInfo(t, info, want, started)
		select {
		case <-done:
			t.Fatal("GPU fixture exited before its PID response could be verified live")
		default:
		}
		t.Logf("verified live PID %d on GPU %d with start time %s", want.PID, want.GPU, info[0].ProcessUtilization.StartTime)
		return
	}
}

func checkLiveProcessInfo(t *testing.T, info []dcgm.ProcessInfo, want dcgm.ProcessInfo, started time.Time) {
	t.Helper()
	if len(info) != 1 || info[0].GPU != want.GPU || info[0].PID != want.PID {
		t.Fatalf("expected PID %d on GPU %d, got %+v", want.PID, want.GPU, info)
	}
	util := info[0].ProcessUtilization
	if util.StartTime < dcgm.Time(started.Unix()) || util.StartTime > dcgm.Time(time.Now().Unix()) || util.EndTime != 0 { // #nosec G115 -- contemporary Unix timestamps are positive.
		t.Fatalf("invalid live process interval: %+v", util)
	}
}
