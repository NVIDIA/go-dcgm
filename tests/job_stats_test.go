package tests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestJobStats uses a private embedded engine so removing all jobs cannot affect another client.
func TestJobStats(t *testing.T) {
	load := jobStatsLoadExecutable(t)
	group, device := jobStatsSession(t)
	if err := dcgm.WatchJobFields(group, 100*time.Millisecond, time.Minute, 0); err != nil {
		t.Fatal(err)
	}
	job := startOwnedJobStats(t, group)
	runJobStatsLoad(t, load, device)
	runningEnd := checkRunningJobStats(t, job.id)
	if err := dcgm.JobStop(job.id); err != nil {
		t.Fatal(err)
	}
	job.stopped = true
	checkStoppedJobStats(t, job.id, device.GPU, runningEnd)
	if err := dcgm.JobRemove(job.id); err != nil {
		t.Fatal(err)
	}
	job.removed = true
	reuseJobStatsAndRemoveAll(t, group, job)
}

func jobStatsLoadExecutable(t *testing.T) string {
	t.Helper()
	if os.Geteuid() != 0 {
		t.Skip("job-state lifecycle test requires root")
	}
	load := os.Getenv("DCGM_JOB_LOAD_COMMAND")
	if load == "" {
		t.Skip("set DCGM_JOB_LOAD_COMMAND to a bounded GPU load executable; see tests/README.md")
	}
	return load
}

func jobStatsSession(t *testing.T) (group dcgm.GroupHandle, device dcgm.Device) {
	t.Helper()
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatal(err)
	}
	restoreAccounting := func() {}
	t.Cleanup(func() { defer restoreAccounting(); cleanup() })
	group = newPidWatchTestGroup(t)
	devices, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) == 0 {
		t.Fatal("supported GPU disappeared during job setup")
	}
	device, err = dcgm.GetDeviceInfo(devices[0])
	if err != nil {
		t.Fatal(err)
	}
	restoreAccounting = jobAccountingRestoration(t, device.UUID)
	return group, device
}

func jobAccountingRestoration(t *testing.T, uuid string) func() {
	t.Helper()
	original := jobNvidiaSMI(t, "--id="+uuid, "--query-gpu=accounting.mode", "--format=csv,noheader,nounits")
	setting := ""
	switch original {
	case "Disabled":
		setting = "0"
	case "Enabled":
		setting = "1"
	case "N/A", "[N/A]":
	default:
		t.Fatalf("unrecognized accounting mode %q", original)
	}
	return func() {
		current := jobNvidiaSMI(t, "--id="+uuid, "--query-gpu=accounting.mode", "--format=csv,noheader,nounits")
		if current != original && setting != "" {
			if _, restoreErr := jobNvidiaSMIOutput("--id="+uuid, "--accounting-mode="+setting); restoreErr != nil {
				t.Errorf("restore accounting mode: %v", restoreErr)
			}
			current = jobNvidiaSMI(t, "--id="+uuid, "--query-gpu=accounting.mode", "--format=csv,noheader,nounits")
		}
		if current != original {
			t.Errorf("accounting mode not restored: got %q, want %q", current, original)
		}
	}
}

type ownedJobStats struct {
	id               string
	stopped, removed bool
}

func startOwnedJobStats(t *testing.T, group dcgm.GroupHandle) *ownedJobStats {
	t.Helper()
	job := &ownedJobStats{id: fmt.Sprintf("go-dcgm-job-%d", time.Now().UnixNano())}
	if err := dcgm.JobStart(group, job.id); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { job.cleanup(t) })
	return job
}

func (job *ownedJobStats) cleanup(t *testing.T) {
	t.Helper()
	if !job.stopped {
		if err := dcgm.JobStop(job.id); err != nil {
			t.Errorf("stop owned job: %v", err)
		}
	}
	if !job.removed {
		if err := dcgm.JobRemove(job.id); err != nil {
			t.Errorf("remove owned job: %v", err)
		}
	}
}

func runJobStatsLoad(t *testing.T, load string, device dcgm.Device) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, load) // #nosec G204,G702 -- explicitly supplied test-only executable, no shell.
	cmd.Env = append(os.Environ(), "DCGM_JOB_GPU_ID="+strconv.FormatUint(uint64(device.GPU), 10), "DCGM_JOB_LOAD_SECONDS=3", "DCGM_JOB_GPU_UUID="+device.UUID)
	cmd.WaitDelay = time.Second
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("GPU load: %v (context: %v): %s", err, ctx.Err(), output)
	}
	time.Sleep(200 * time.Millisecond)
}

func checkRunningJobStats(t *testing.T, id string) int64 {
	t.Helper()
	beforeRead := time.Now().UnixMicro()
	running, err := dcgm.JobGetStats(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(running.GPUs) != 1 {
		t.Fatalf("expected one running GPU, got %d", len(running.GPUs))
	}
	gpu := &running.GPUs[0]
	if gpu.StartTime <= 0 || gpu.EndTime <= gpu.StartTime || gpu.EndTime < beforeRead || gpu.EndTime > time.Now().UnixMicro() {
		t.Fatalf("running job has no current collection endpoint: %+v", running)
	}
	return gpu.EndTime
}

func checkStoppedJobStats(t *testing.T, id string, gpuID uint, runningEnd int64) {
	t.Helper()
	stats, err := dcgm.JobGetStats(id)
	if err != nil {
		t.Fatal(err)
	}
	if len(stats.GPUs) != 1 {
		t.Fatalf("expected one GPU, got %d", len(stats.GPUs))
	}
	gpu := &stats.GPUs[0]
	if gpu.GPU != gpuID || gpu.StartTime <= 0 || gpu.EndTime <= gpu.StartTime || gpu.EndTime < runningEnd {
		t.Fatalf("invalid job interval or GPU: %+v", gpu)
	}
	if dcgm.IsInt32Blank(int(gpu.SMUtilization.Max)) || gpu.SMUtilization.Max <= 0 {
		t.Fatalf("load produced no SM utilization: %+v", gpu.SMUtilization)
	}
}

func reuseJobStatsAndRemoveAll(t *testing.T, group dcgm.GroupHandle, job *ownedJobStats) {
	t.Helper()
	// Reuse proves removal, then exercise RemoveAll only on this private engine.
	if err := dcgm.JobStart(group, job.id); err != nil {
		t.Fatal(err)
	}
	job.removed = false
	job.stopped = false
	if err := dcgm.JobStop(job.id); err != nil {
		t.Fatal(err)
	}
	job.stopped = true
	if err := dcgm.JobRemoveAll(); err != nil {
		t.Fatal(err)
	}
	job.removed = true
	_, err := dcgm.JobGetStats(job.id)
	var nativeErr *dcgm.Error
	if !errors.As(err, &nativeErr) || int(nativeErr.Code) != dcgm.DCGM_ST_NO_DATA {
		t.Fatalf("removed job: expected DCGM_ST_NO_DATA, got %v", err)
	}
}

func jobNvidiaSMI(t *testing.T, args ...string) string {
	t.Helper()
	output, err := jobNvidiaSMIOutput(args...)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

func jobNvidiaSMIOutput(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "nvidia-smi", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("nvidia-smi %v: %w: %s", args, err, output)
	}
	return strings.TrimSpace(string(output)), nil
}
