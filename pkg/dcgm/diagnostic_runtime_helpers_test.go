//go:build linux && cgo

package dcgm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// diagnosticCommand uses a separate process group so cancellation also stops
// child processes that still have the output pipes open.
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

// startDiagnostic runs the diagnostic in the background. Its cleanup waits
// for the call to return before DCGM is unloaded.
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
		// The parent process times out if the diagnostic never returns.
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

// diagnosticSession starts DCGM in embedded mode and checks for a live GPU.
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

// diagnosticUnavailable reports errors that let a hardware test skip.
func diagnosticUnavailable(err error) bool {
	if IsFeatureUnavailable(err) {
		return true
	}
	var dcgmErr *Error
	if !errors.As(err, &dcgmErr) {
		return false
	}
	switch int(dcgmErr.Code) {
	case DCGM_ST_NOT_SUPPORTED, DCGM_ST_REQUIRES_ROOT, DCGM_ST_NVVS_BINARY_NOT_FOUND:
		return true
	}

	return false
}
