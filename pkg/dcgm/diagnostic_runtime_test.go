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
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
