//go:build linux && cgo

package dcgm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// diagnosticProcess runs the test in a child process with a three-minute limit.
// If the command fails, it cleans up that process group.
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

// nvvsChildRunning reports whether this test process started an NVVS child.
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

// waitDiagnosticStart waits for this test's NVVS process to appear.
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

// stopAndWaitDiagnostic stops the diagnostic and waits up to 30 seconds for it to return.
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

// diagnostic48Runtime reports whether the host engine is running DCGM 4.8 or later.
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

// assertDiagnosticStopped checks for a DCGM error. DCGM 4.8 and later must return DIAG_STOPPED.
func assertDiagnosticStopped(t *testing.T, err error) {
	t.Helper()
	var dcgmErr *Error
	require.ErrorAs(t, err, &dcgmErr)
	t.Logf("cancelled diagnostic: code=%d message=%s", dcgmErr.Code, dcgmErr.Error())

	strict := diagnostic48Runtime(t)
	t.Run("DCGM48StoppedCode", func(t *testing.T) {
		if !strict {
			t.Skip("checking for DIAG_STOPPED requires DCGM 4.8 or later")
		}
		require.Equal(t, DCGM_ST_DIAG_STOPPED, int(dcgmErr.Code))
	})
}

// TestDiagnosticCancellationAndRestart checks that a new diagnostic can start after cancellation.
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
		require.True(t, slices.ContainsFunc(restarted.Tests, func(result DiagTest) bool {
			return result.Status == "pass"
		}), "DCGM did not report a passing check after diagnostics restarted")
	})
}

// TestDiagnosticHeartbeat checks that a diagnostic keeps running while it
// receives heartbeats and stops when they end.
func TestDiagnosticHeartbeat(t *testing.T) {
	t.Setenv("DCGM_DIAG_HEARTBEAT_TIMEOUT_SEC", "2")
	diagnosticProcess(t, func() {
		diagnosticSession(t)
		gpu, selected := diagnosticGPU(t)

		request := DiagnosticRequest{
			TestNames:       []string{"diagnostic"},
			TestParameters:  []string{"diagnostic.test_duration=60"},
			Timeout:         2 * time.Minute,
			EnableHeartbeat: true,
		}
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
				t.Skip("heartbeat timing checks require DCGM 4.8 or later")
			}

			// Keep sending heartbeats past the next eight-second expiry check.
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
				t.Skip("heartbeat timing checks require DCGM 4.8 or later")
			}

			done := startDiagnostic(t, func() error {
				_, err := ActionValidate(request)
				return err
			})
			waitDiagnosticStart(t, done)
			select {
			case err := <-done:
				assertDiagnosticStopped(t, err)
			// Allow the eight-second heartbeat poll and up to forty seconds of NVVS teardown.
			case <-time.After(time.Minute):
				t.Fatal("diagnostic did not expire without heartbeats")
			}
		})
	})
}

// TestDiagnosticNamedParameters checks that DCGM accepts a valid duration.
// DCGM 4.8 and later must also reject invalid numbers.
func TestDiagnosticNamedParameters(t *testing.T) {
	diagnosticProcess(t, func() {
		diagnosticSession(t)
		gpu, _ := diagnosticGPU(t)

		request := DiagnosticRequest{
			TestNames:      []string{"diagnostic"},
			Entities:       []GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpu}},
			Timeout:        2 * time.Minute,
			TestParameters: []string{"diagnostic.test_duration=1"},
		}
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
				t.Skip("checking that DCGM rejects invalid numbers requires DCGM 4.8 or later")
			})
			return
		}

		for _, input := range []struct {
			value, message string
		}{
			{"1suffix", "Invalid numeric test parameter"},
			{" 1", "Invalid numeric test parameter"},
			{"1 ", "Invalid numeric test parameter"},
			{"NaN", "Invalid numeric test parameter"},
			{"inf", "Invalid numeric test parameter"},
			{"1e10000", "Numeric test parameter is outside the double range"},
		} {
			value := input.value
			t.Run(value, func(t *testing.T) {
				request.TestParameters = []string{"diagnostic.test_duration=" + value}
				result, err := ActionValidate(request)
				t.Logf("invalid numeric parameter %q: error=%v results=%+v", value, err, result)
				require.True(t, diagnosticParameterRejected(result, err, input.message), "expected numeric parameter rejection")
			})
		}
	})
}

// TestDiagnosticConnectXRuntime checks the ConnectX results returned by the optional NIC-EUD plugin.
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

		// NIC-EUD checks all NICs and rejects explicit device lists.
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
