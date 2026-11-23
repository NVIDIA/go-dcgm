//go:build linux && cgo

/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dcgm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func buildErrorLibrary(t *testing.T, flags ...string) string {
	t.Helper()
	dir := t.TempDir()
	args := append([]string{"-shared", "-fPIC", "-I.", "testdata/error_handling.c", "-o", filepath.Join(dir, dcgmLibraryName)}, flags...)
	cmd := exec.Command("cc", args...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	return dir
}

func runErrorChild(t *testing.T, testName, dir string, extra ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^"+testName+"$", "-test.v") //nolint:gosec // Runs this test binary with a fixed test selector.
	// This target leaves libdcgm unlinked so Init's dlopen selects the fixture.
	cmd.Env = append(os.Environ(), "LD_LIBRARY_PATH="+dir, "LD_PRELOAD=", "LD_BIND_NOW=", "GO_DCGM_VERSION_FD=")
	cmd.Env = append(cmd.Env, extra...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", out)
	require.Contains(t, string(out), "--- PASS: "+testName, "child must execute the selected test")
}

// TestLoadedLibraryDiagnostics selects a fixture in a fresh loader process for
// each combination of missing symbols, bad build info, and native errors.
func TestLoadedLibraryDiagnostics(t *testing.T) {
	if os.Getenv("GO_DCGM_LOADER_CHILD") == "1" {
		checkLoadedLibraryDiagnostics(t)
		return
	}
	for _, flags := range [][]string{
		nil,
		{"-DNO_CLIENT"},
		{"-DNO_HOST"},
		{"-DNO_CLIENT", "-DNO_HOST"},
		{"-DUNTERMINATED_VERSION"},
		{"-DVERSION_ERROR"},
		{"-DVERSION_ERROR", "-DNO_ERROR_TEXT"},
	} {
		t.Run(fmt.Sprint(flags), func(t *testing.T) {
			runErrorChild(t, "TestLoadedLibraryDiagnostics", buildErrorLibrary(t, flags...),
				"GO_DCGM_LOADER_CHILD=1", "GO_DCGM_FLAGS="+strings.Join(flags, " "))
		})
	}

	t.Run("invalid handshake descriptors", func(t *testing.T) {
		dir := buildErrorLibrary(t)
		for _, fd := range []string{"invalid", "0", "2", strings.Repeat("9", 32)} {
			t.Run(fd, func(t *testing.T) {
				runErrorChild(t, "TestLoadedLibraryDiagnostics", dir,
					"GO_DCGM_LOADER_CHILD=1", "GO_DCGM_FLAGS=", "GO_DCGM_VERSION_FD="+fd)
			})
		}
	})
}

// checkLoadedLibraryDiagnostics verifies the loaded path and queries before and after closing it.
func checkLoadedLibraryDiagnostics(t *testing.T) {
	t.Helper()
	path, err := GetLoadedLibraryPath()
	require.Error(t, err)
	require.Empty(t, path)
	require.NoError(t, cgoAdapter{}.loadLibrary())
	defer cgoAdapter{}.closeLibrary()
	path, err = GetLoadedLibraryPath()
	require.NoError(t, err)
	want := filepath.Join(os.Getenv("LD_LIBRARY_PATH"), dcgmLibraryName)
	require.Equal(t, want, path)

	// Check the backing address while the native storage is valid.
	nativePath := loadedLibraryPathCString()
	require.NotNil(t, nativePath)
	require.NotEqual(t, unsafe.Pointer(nativePath), unsafe.Pointer(unsafe.StringData(path)))

	// The fixture covers symbol lookup. The real-runtime test owns session startup.
	dcgmInitCounter = 1
	flags := strings.Fields(os.Getenv("GO_DCGM_FLAGS"))
	for _, query := range []struct {
		name, absent, want string
		read               func() (VersionInfo, error)
	}{
		{"client", "-DNO_CLIENT", "version:4.8.0;", GetVersionInfo},
		{"hostengine", "-DNO_HOST", "version:4.7.0;", GetHostengineVersionInfo},
	} {
		checkFixtureVersion(t, flags, query.name, query.absent, query.want, query.read)
	}

	// New queries report the missing library after dlclose.
	dcgmInitCounter = 0
	cgoAdapter{}.closeLibrary()
	require.Equal(t, want, path)
	closedPath, err := GetLoadedLibraryPath()
	require.Error(t, err)
	require.Empty(t, closedPath)
	_, err = GetVersionInfo()
	require.True(t, IsFeatureUnavailable(err))
	_, err = GetHostengineVersionInfo()
	require.True(t, IsFeatureUnavailable(err))
}

// checkFixtureVersion checks the result selected by the fixture's compiler flags.
func checkFixtureVersion(t *testing.T, flags []string, name, absent, want string, read func() (VersionInfo, error)) {
	t.Helper()
	info, queryErr := read()
	switch {
	case slices.Contains(flags, absent):
		var unavailable *FeatureUnavailableError
		require.ErrorAs(t, queryErr, &unavailable, name)
		require.Empty(t, info.RawBuildInfoString)
	case slices.Contains(flags, "-DUNTERMINATED_VERSION"):
		require.ErrorContains(t, queryErr, "unterminated version string", name)
		require.Empty(t, info.RawBuildInfoString)
	case slices.Contains(flags, "-DVERSION_ERROR"):
		var native *Error
		require.ErrorAs(t, queryErr, &native, name)
		require.Equal(t, DCGM_ST_VER_MISMATCH, int(native.Code), name)
		wantMessage := "fixture version mismatch"
		if slices.Contains(flags, "-DNO_ERROR_TEXT") {
			wantMessage = fmt.Sprintf("DCGM version query failed with status %d", DCGM_ST_VER_MISMATCH)
		}
		require.Equal(t, wantMessage, native.Error(), name)
		require.Empty(t, info.RawBuildInfoString)
	default:
		require.NoError(t, queryErr)
		require.Equal(t, want, info.RawBuildInfoString)
	}
}

// TestVersionQueriesHoldSessionLock checks that the library cannot close while
// a version query is executing inside the native fixture.
func TestVersionQueriesHoldSessionLock(t *testing.T) {
	if os.Getenv("GO_DCGM_VERSION_CHILD") != "1" {
		runErrorChild(t, "TestVersionQueriesHoldSessionLock", buildErrorLibrary(t), "GO_DCGM_VERSION_CHILD=1")
		return
	}
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	require.NoError(t, err)
	signal := os.NewFile(uintptr(fds[0]), "version-signal")
	// The fixture uses fds[1] directly; this file only keeps it open and closes it.
	native := os.NewFile(uintptr(fds[1]), "version-native")
	t.Cleanup(func() {
		require.NoError(t, signal.Close())
		require.NoError(t, native.Close())
	})
	t.Setenv("GO_DCGM_VERSION_FD", strconv.Itoa(fds[1]))
	require.NoError(t, cgoAdapter{}.loadLibrary())
	defer cgoAdapter{}.closeLibrary()
	dcgmInitCounter = 1
	defer func() { dcgmInitCounter = 0 }()
	v10, _, _ := diagnosticCompatibilityVersionsForTest()
	for _, query := range []struct {
		name  string
		calls int
		read  func() (VersionInfo, error)
	}{
		{"client", 1, GetVersionInfo},
		{"hostengine", 1, GetHostengineVersionInfo},
		{"diagnostic selection", 2, func() (VersionInfo, error) {
			if got := uint(selectedDiagnosticRequestVersion()); got != v10 {
				return VersionInfo{}, fmt.Errorf("diagnostic request version %d, want %d", got, v10)
			}
			return VersionInfo{}, nil
		}},
	} {
		t.Run(query.name, func(t *testing.T) {
			done := make(chan error, 1)
			go func() {
				_, err := query.read()
				done <- err
			}()
			for range query.calls {
				_, err := signal.Read(make([]byte, 1))
				require.NoError(t, err)
				if mux.TryLock() {
					mux.Unlock()
					t.Error("version query must hold the session lock inside native code")
				}
				_, err = signal.Write([]byte{1})
				require.NoError(t, err)
			}
			require.NoError(t, <-done)
		})
	}
}
