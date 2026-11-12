//go:build linux && cgo

package dcgm

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// diagnosticProcess bounds native diagnostics and kills only this test's process group on exit.
func diagnosticProcess(t *testing.T, run func()) {
	t.Helper()
	if os.Getenv("DCGM_TEST_DIAGNOSTICS") != "1" {
		t.Skip("requires task test:diagnostics on a selected diagnostic host")
	}

	if os.Getenv("GO_DCGM_DIAGNOSTIC_CHILD") == t.Name() {
		run()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	executable, err := os.Executable()
	require.NoError(t, err)
	cmd := diagnosticCommand(ctx, executable, "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.v")
	cmd.Env = append(os.Environ(), "GO_DCGM_DIAGNOSTIC_CHILD="+t.Name())
	output, err := cmd.CombinedOutput()
	if cmd.Process != nil && err != nil {
		if killErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); killErr != nil && !errors.Is(killErr, syscall.ESRCH) {
			t.Errorf("kill diagnostic child group: %v", killErr)
		}
	}
	t.Logf("diagnostic child:\n%s", output)
	require.NoError(t, ctx.Err(), "diagnostic child exceeded three minutes")
	require.NoError(t, err)
	if diagnosticChildSkipped(output, t.Name()) {
		t.Skip("diagnostic child reported unavailable prerequisites; see its output")
	}
}

// diagnosticCommand owns a process group so cancellation also stops children holding output pipes.
func diagnosticCommand(ctx context.Context, executable string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, executable, args...) //nolint:gosec // Runs only the test-owned executable and arguments without a shell.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return os.ErrProcessDone
		}
		return err
	}

	cmd.WaitDelay = 2 * time.Second
	return cmd
}

// diagnosticPipeHelper runs the subprocess roles used to test pipe-holder cancellation.
func diagnosticPipeHelper(t *testing.T) bool {
	t.Helper()
	args := flag.Args()
	if len(args) != 2 {
		return false
	}
	switch args[0] {
	case "diagnostic-pipe-parent":
		cmd := exec.Command(os.Args[0], "-test.run=^TestDiagnosticProcessCancellation$", "--", "diagnostic-pipe-holder", args[1]) //nolint:gosec // Reexecutes this test binary with its private readiness file.
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		require.NoError(t, cmd.Run())
		return true
	case "diagnostic-pipe-holder":
		ready := fmt.Sprintf("%d %d", syscall.Getpgrp(), os.Getpid())
		require.NoError(t, os.WriteFile(args[1], []byte(ready), 0o600))
		time.Sleep(15 * time.Second)
		return true
	default:
		return false
	}
}

// diagnosticPIDRunning reports whether a process still exists and has not become a zombie.
func diagnosticPIDRunning(pid int) (bool, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	fields := strings.Fields(string(data[bytes.LastIndexByte(data, ')')+1:]))
	if len(fields) == 0 {
		return false, fmt.Errorf("missing process state for PID %d", pid)
	}

	return fields[0] != "Z" && fields[0] != "X", nil
}

// TestDiagnosticProcessCancellation stops a pipe-holding descendant before waiting for command output.
func TestDiagnosticProcessCancellation(t *testing.T) {
	if diagnosticPipeHelper(t) {
		return
	}

	ready := filepath.Join(t.TempDir(), "ready")
	ctx, cancel := context.WithCancel(context.Background())
	cmd := diagnosticCommand(ctx, os.Args[0], "-test.run=^TestDiagnosticProcessCancellation$", "--", "diagnostic-pipe-parent", ready)
	done := make(chan struct{})
	var output []byte
	var runErr error
	var groupID, descendant int
	t.Cleanup(func() {
		cancel()
		if groupID > 0 {
			if err := syscall.Kill(-groupID, syscall.SIGKILL); err != nil && !errors.Is(err, syscall.ESRCH) {
				t.Errorf("kill pipe-holder group: %v", err)
			}
		}
		select {
		case <-done:
		case <-time.After(20 * time.Second):
			t.Error("pipe-holder command did not finish during cleanup")
		}
	})
	go func() {
		output, runErr = cmd.CombinedOutput()
		close(done)
	}()
	require.Eventually(t, func() bool {
		data, err := os.ReadFile(ready)
		if err != nil {
			return false
		}
		var group, pid int
		_, err = fmt.Sscanf(string(data), "%d %d", &group, &pid)
		return err == nil && group > 0 && pid > 0
	}, 5*time.Second, 10*time.Millisecond, "pipe-holding descendant did not become ready")
	data, err := os.ReadFile(ready)
	require.NoError(t, err)
	_, err = fmt.Sscanf(string(data), "%d %d", &groupID, &descendant)
	require.NoError(t, err)
	require.NotEqual(t, syscall.Getpgrp(), groupID, "fixture must own a separate process group")
	cancel()
	select {
	case <-done:
		t.Logf("cancelled pipe-holder output: %s", output)
		require.Error(t, runErr)
	case <-time.After(3 * time.Second):
		t.Error("command cancellation waited for its descendant's output pipe")
		return
	}
	require.Eventually(t, func() bool {
		running, err := diagnosticPIDRunning(descendant)
		return err == nil && !running
	}, time.Second, 10*time.Millisecond, "cancelled command left its descendant running")
}

// TestDiagnosticCleanupWaits checks session cleanup after an early test failure.
func TestDiagnosticCleanupWaits(t *testing.T) {
	if os.Getenv("GO_DCGM_DIAGNOSTIC_CLEANUP_CHILD") == "1" {
		finished := make(chan struct{})
		t.Cleanup(func() {
			select {
			case <-finished:
				t.Log("diagnostic returned before session cleanup")
			default:
				t.Error("session cleanup ran while diagnostic was active")
			}
		})
		startDiagnostic(t, func() error {
			time.Sleep(100 * time.Millisecond)
			close(finished)
			return nil
		})
		t.Fatal("exercise cleanup after an early failure")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := diagnosticCommand(ctx, os.Args[0], "-test.run=^TestDiagnosticCleanupWaits$", "-test.v")
	cmd.Env = append(os.Environ(), "GO_DCGM_DIAGNOSTIC_CLEANUP_CHILD=1")
	output, err := cmd.CombinedOutput()
	require.NoError(t, ctx.Err())
	require.Error(t, err, "child must report its intentional failure")
	require.Contains(t, string(output), "diagnostic returned before session cleanup")
	require.NotContains(t, string(output), "session cleanup ran while diagnostic was active")
}

// startDiagnostic joins the native call before earlier session cleanup can unload DCGM.
func startDiagnostic(t *testing.T, run func() error) <-chan error {
	t.Helper()
	done := make(chan error, 1)
	finished := make(chan struct{})
	t.Cleanup(func() {
		select {
		case <-finished:
			return
		default:
		}
		if err := StopDiag(); err != nil {
			t.Errorf("stop diagnostic during cleanup: %v", err)
		}
		// The parent process bounds this wait if the native call never returns.
		<-finished
	})
	go func() {
		defer close(finished)
		done <- run()
		close(done)
	}()
	return done
}

// diagnosticChildSkipped recognizes only a top-level skip for the selected test.
func diagnosticChildSkipped(output []byte, testName string) bool {
	return strings.Contains(string(output), "--- SKIP: "+testName+" (")
}

// TestDiagnosticChildSkipClassification keeps subtest skips from hiding a successful selected test.
func TestDiagnosticChildSkipClassification(t *testing.T) {
	for _, tc := range []struct {
		name   string
		output string
		skip   bool
	}{
		{"top-level skip", "--- SKIP: TestDiagnosticNamedParameters (0.01s)\n", true},
		{"subtest skip", "--- PASS: TestDiagnosticNamedParameters (0.02s)\n    --- SKIP: TestDiagnosticNamedParameters/DCGM48StrictNumericParameters (0.00s)\n", false},
		{"similar test name", "--- SKIP: TestDiagnosticNamedParametersOther (0.01s)\n", false},
		{"top-level pass", "--- PASS: TestDiagnosticNamedParameters (0.01s)\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.skip, diagnosticChildSkipped([]byte(tc.output), "TestDiagnosticNamedParameters"))
		})
	}
}

// TestDiagnosticAPIsAfterFailedInit rejects diagnostic calls when the library is loaded without a session.
func TestDiagnosticAPIsAfterFailedInit(t *testing.T) {
	requireLoadedLibraryWithoutSession(t)
	for _, tc := range []struct {
		symbol string
		call   func() error
	}{
		{"dcgmStopDiagnostic", StopDiag},
		{"dcgmDiagSendHeartbeat", DiagSendHeartbeat},
		{"dcgmProfPause", ProfPause},
		{"dcgmProfResume", ProfResume},
		{"dcgmActionValidate_v2", func() error {
			_, err := ActionValidate(DiagnosticRequest{Level: DiagQuick})
			return err
		}},
	} {
		t.Run(tc.symbol, func(t *testing.T) {
			require.True(t, dcgmSymbolAvailable(tc.symbol), "regression requires a loaded native symbol")
			err := tc.call()
			var unavailable *FeatureUnavailableError
			require.ErrorAs(t, err, &unavailable)
			require.Equal(t, tc.symbol, unavailable.Symbol)
		})
	}
}

// diagnosticSession starts an embedded session and requires a live GPU for the runtime fixture.
func diagnosticSession(t *testing.T) {
	t.Helper()
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	t.Cleanup(cleanup)
	runOnlyWithLiveGPUs(t)
}

// parseDiagnosticGPU parses the explicitly selected diagnostic GPU identifier.
func parseDiagnosticGPU(value string, supported []uint) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid DCGM_TEST_DIAGNOSTIC_GPU_ID %q: %w", value, err)
	}
	for _, gpu := range supported {
		if gpu == uint(id) {
			return gpu, nil
		}
	}
	return 0, fmt.Errorf("DCGM_TEST_DIAGNOSTIC_GPU_ID %d is not a supported GPU", id)
}

// diagnosticGPU selects the test GPU and checks that it is supported.
func diagnosticGPU(t *testing.T) (uint, bool) {
	t.Helper()
	gpus, err := GetSupportedDevices()
	require.NoError(t, err)
	require.NotEmpty(t, gpus)
	value, selected := os.LookupEnv("DCGM_TEST_DIAGNOSTIC_GPU_ID")
	if !selected {
		return gpus[0], false
	}
	gpu, err := parseDiagnosticGPU(value, gpus)
	require.NoError(t, err)
	t.Logf("selected diagnostic GPU: %d", gpu)
	return gpu, true
}

// TestDiagnosticGPUSelection checks explicit GPU selection and invalid identifiers.
func TestDiagnosticGPUSelection(t *testing.T) {
	for _, tc := range []struct {
		value string
		gpu   uint
		valid bool
	}{
		{"0", 0, true},
		{"1", 1, true},
		{"2", 0, false},
		{"", 0, false},
		{"-1", 0, false},
		{" 1", 0, false},
		{"gpu:1", 0, false},
		{"4294967296", 0, false},
	} {
		t.Run(tc.value, func(t *testing.T) {
			gpu, err := parseDiagnosticGPU(tc.value, []uint{0, 1})
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.gpu, gpu)
		})
	}
}

// diagnosticUnavailable identifies missing diagnostic prerequisites that runtime tests may skip.
func diagnosticUnavailable(err error) bool {
	if IsFeatureUnavailable(err) {
		return true
	}
	var native *Error
	if !errors.As(err, &native) {
		return false
	}
	switch int(native.Code) {
	case DCGM_ST_NOT_SUPPORTED, DCGM_ST_REQUIRES_ROOT, DCGM_ST_NVVS_BINARY_NOT_FOUND:
		return true
	}

	return false
}

// nvvsChildRunning finds an NVVS process directly owned by this test executable.
func nvvsChildRunning() bool {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if _, err := strconv.Atoi(entry.Name()); err != nil {
			continue
		}
		root := "/proc/" + entry.Name()
		comm, err := os.ReadFile(root + "/comm")
		if err != nil || strings.TrimSpace(string(comm)) != "nvvs" {
			continue
		}
		status, err := os.ReadFile(root + "/status")
		if err == nil && strings.Contains(string(status), fmt.Sprintf("PPid:\t%d\n", os.Getpid())) {
			return true
		}
	}

	return false
}

// waitDiagnosticStart waits for the owned NVVS child or reports why the diagnostic ended early.
func waitDiagnosticStart(t *testing.T, done <-chan error) {
	t.Helper()
	deadline := time.NewTimer(30 * time.Second)
	defer deadline.Stop()
	tick := time.NewTicker(50 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case err := <-done:
			if diagnosticUnavailable(err) {
				t.Skipf("diagnostic prerequisite unavailable: %v", err)
			}
			t.Fatalf("diagnostic returned before NVVS could be stopped: %v", err)
		case <-deadline.C:
			t.Fatal("NVVS did not start within 30 seconds")
		case <-tick.C:
			if nvvsChildRunning() {
				return
			}
		}
	}
}

// stopAndWaitDiagnostic requests cancellation and bounds the wait for its result.
func stopAndWaitDiagnostic(t *testing.T, done <-chan error) error {
	t.Helper()
	require.NoError(t, StopDiag())
	select {
	case err := <-done:
		return err
	case <-time.After(30 * time.Second):
		t.Fatal("diagnostic did not return within 30 seconds of StopDiag")
		return nil
	}
}

// diagnostic48Runtime checks whether the runtime includes the diagnostic fixes introduced in DCGM 4.8.
func diagnostic48Runtime(t *testing.T) bool {
	t.Helper()
	info, err := GetHostengineVersionInfo()
	require.NoError(t, err)
	t.Logf("host engine build: %s", info.RawBuildInfoString)
	for _, pair := range strings.Split(info.RawBuildInfoString, ";") {
		key, version, found := strings.Cut(pair, ":")
		if !found || strings.TrimSpace(key) != "version" {
			continue
		}
		parts := strings.Split(strings.TrimSpace(version), ".")
		require.GreaterOrEqual(t, len(parts), 2)
		major, err := strconv.Atoi(parts[0])
		require.NoError(t, err)
		minor, err := strconv.Atoi(parts[1])
		require.NoError(t, err)
		return major > 4 || (major == 4 && minor >= 8)
	}
	t.Fatal("host engine did not report its version")
	return false
}

// assertDiagnosticStopped requires a typed cancellation result and checks the fixed native status when supported.
func assertDiagnosticStopped(t *testing.T, err error) {
	t.Helper()
	var native *Error
	require.ErrorAs(t, err, &native)
	t.Logf("cancelled diagnostic: code=%d message=%s", native.Code, native.Error())
	strict := diagnostic48Runtime(t)
	t.Run("DCGM48StoppedCode", func(t *testing.T) {
		if !strict {
			t.Skip("exact DIAG_STOPPED acceptance requires a DCGM 4.8 runtime containing the cancellation fix")
		}
		require.Equal(t, DCGM_ST_DIAG_STOPPED, int(native.Code))
	})
}

// TestDiagnosticCancellationAndRestart restarts a diagnostic after the previous invocation has stopped.
func TestDiagnosticCancellationAndRestart(t *testing.T) {
	diagnosticProcess(t, func() {
		diagnosticSession(t)
		gpu, selected := diagnosticGPU(t)
		group := GroupAllGPUs()
		if selected {
			var err error
			group, err = CreateGroup(fmt.Sprintf("diagnostic-%d", time.Now().UnixNano()))
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, DestroyGroup(group)) })
			require.NoError(t, AddToGroup(group, gpu))
		}
		done := startDiagnostic(t, func() error {
			_, err := RunDiag(DiagLong, group)
			return err
		})
		waitDiagnosticStart(t, done)
		err := stopAndWaitDiagnostic(t, done)
		assertDiagnosticStopped(t, err)
		restarted, err := RunDiag(DiagQuick, group)
		require.NoError(t, err)
		require.NotEmpty(t, restarted.Tests, "next diagnostic did not produce results")
	})
}

// TestDiagnosticHeartbeat keeps an opted-in diagnostic alive until its stop request.
func TestDiagnosticHeartbeat(t *testing.T) {
	t.Setenv("DCGM_DIAG_HEARTBEAT_TIMEOUT_SEC", "2")
	diagnosticProcess(t, func() {
		diagnosticSession(t)
		gpu, selected := diagnosticGPU(t)
		request := DiagnosticRequest{TestNames: []string{"diagnostic"}, TestParameters: []string{"diagnostic.test_duration=60"}, Timeout: 2 * time.Minute, EnableHeartbeat: true}
		if selected {
			request.Entities = []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpu}}
		}
		done := startDiagnostic(t, func() error {
			_, err := ActionValidate(request)
			return err
		})
		waitDiagnosticStart(t, done)
		require.NoError(t, DiagSendHeartbeat())
		require.NoError(t, DiagSendHeartbeat())
		strict := diagnostic48Runtime(t)
		t.Run("DCGM48HeartbeatSurvival", func(t *testing.T) {
			if !strict {
				t.Skip("heartbeat expiry checks require a DCGM 4.8 runtime supporting the timeout override")
			}
			// DCGM checks expiry every eight seconds; outlast that interval with regular heartbeats.
			deadline := time.NewTimer(12 * time.Second)
			defer deadline.Stop()
			tick := time.NewTicker(500 * time.Millisecond)
			defer tick.Stop()
			for {
				select {
				case err := <-done:
					t.Fatalf("diagnostic ended while heartbeats were sent: %v", err)
				case <-tick.C:
					require.NoError(t, DiagSendHeartbeat())
				case <-deadline.C:
					return
				}
			}
		})
		err := stopAndWaitDiagnostic(t, done)
		assertDiagnosticStopped(t, err)
		t.Run("DCGM48HeartbeatExpiry", func(t *testing.T) {
			if !strict {
				t.Skip("heartbeat expiry checks require a DCGM 4.8 runtime supporting the timeout override")
			}
			done := startDiagnostic(t, func() error {
				_, err := ActionValidate(request)
				return err
			})
			waitDiagnosticStart(t, done)
			select {
			case err := <-done:
				assertDiagnosticStopped(t, err)
			case <-time.After(20 * time.Second):
				t.Fatal("diagnostic did not expire without heartbeats")
			}
		})
	})
}

// diagnosticRejected detects native errors and failed diagnostic results for invalid-parameter checks.
func diagnosticRejected(result DiagResults, err error) bool {
	if err != nil || len(result.SystemErrors) != 0 {
		return true
	}

	for i := range result.Tests {
		test := &result.Tests[i]
		if test.Status == "fail" {
			return true
		}
	}

	return false
}

// TestDiagnosticNamedParameters accepts a valid duration and checks strict rejection on runtimes with the parser fix.
func TestDiagnosticNamedParameters(t *testing.T) {
	diagnosticProcess(t, func() {
		diagnosticSession(t)
		gpu, _ := diagnosticGPU(t)
		request := DiagnosticRequest{TestNames: []string{"diagnostic"}, Entities: []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpu}}, Timeout: 2 * time.Minute}
		request.TestParameters = []string{"diagnostic.test_duration=1"}
		result, err := ActionValidate(request)
		if diagnosticUnavailable(err) {
			t.Skipf("named diagnostic unavailable: %v", err)
		}
		t.Logf("valid numeric parameter: error=%v results=%+v", err, result)
		require.NoError(t, err)
		require.NotEmpty(t, result.Tests)
		var found bool
		for i := range result.Tests {
			test := &result.Tests[i]
			if test.Name == "diagnostic" {
				found = true
				if test.Status == "skipped" {
					t.Skip("named diagnostic was skipped by DCGM")
				}
				require.Equal(t, "pass", test.Status, "named diagnostic did not execute successfully")
			}
		}
		require.True(t, found, "named diagnostic was not run")
		if !diagnostic48Runtime(t) {
			t.Run("DCGM48StrictNumericParameters", func(t *testing.T) {
				t.Skip("strict numeric parameter acceptance requires a DCGM 4.8 runtime containing the parser fix")
			})
			return
		}
		for _, value := range []string{"1suffix", " 1", "1 ", "NaN", "inf", "1e10000"} {
			t.Run(value, func(t *testing.T) {
				request.TestParameters = []string{"diagnostic.test_duration=" + value}
				result, err := ActionValidate(request)
				t.Logf("invalid numeric parameter %q: error=%v results=%+v", value, err, result)
				require.True(t, diagnosticRejected(result, err), "DCGM accepted the invalid numeric parameter")
			})
		}
	})
}

// TestDiagnosticConnectXRuntime checks named network-device results when the optional NIC-EUD plugin is available.
func TestDiagnosticConnectXRuntime(t *testing.T) {
	if os.Getenv("DCGM_TEST_NIC_EUD") != "1" {
		t.Skip("requires DCGM_TEST_NIC_EUD=1 with supported ConnectX hardware and the matching proprietary plugin")
	}

	diagnosticProcess(t, func() {
		diagnosticSession(t)
		entities, err := GetEntityGroupEntities(FE_CONNECTX)
		require.NoError(t, err)
		if len(entities) == 0 {
			t.Skip("NIC-EUD requires supported ConnectX hardware and the proprietary plugin")
		}
		// NIC-EUD uses all NICs; the runtime rejects explicit entity selectors.
		result, err := ActionValidate(DiagnosticRequest{
			TestNames: []string{"nic_eud"},
			Timeout:   2 * time.Minute,
		})
		t.Logf("NIC-EUD: error=%v results=%+v", err, result)
		if diagnosticUnavailable(err) {
			t.Skipf("NIC-EUD prerequisite unavailable: %v", err)
		}
		require.NoError(t, err)
		var found bool
		for i := range result.Tests {
			test := &result.Tests[i]
			if test.Name != "nic_eud" {
				continue
			}
			found = true
			require.NotEmpty(t, test.Results)
			for _, entity := range test.Results {
				require.Equal(t, FE_CONNECTX, entity.Entity.EntityGroupId)
			}
			require.NotEmpty(t, test.AuxData, "NIC-EUD did not supply auxiliary data for conversion acceptance")
		}
		require.True(t, found, "NIC-EUD plugin did not produce a named test")
	})
}

// TestProfilingPauseResume restores profiling after checking both control calls.
func TestProfilingPauseResume(t *testing.T) {
	diagnosticSession(t)
	err := ProfPause()
	if diagnosticUnavailable(err) {
		t.Skipf("profiling control unavailable: %v", err)
	}
	require.NoError(t, err)
	paused := true
	defer func() {
		if paused {
			require.NoError(t, ProfResume())
		}
	}()
	require.NoError(t, ProfResume())
	paused = false
}
