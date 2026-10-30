package dcgm

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestConnectStandaloneV2RejectsInvalidSocketFlagBeforeInit(t *testing.T) {
	if err := connectStandaloneV2("127.0.0.1:5555", "invalid"); !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("error = %v, want socket flag parse error", err)
	}
}

// shortTempDir returns a new directory under /tmp so tests do not depend on the length of TMPDIR.
func shortTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "h")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestCreateHostengineSocketDir(t *testing.T) {
	t.Setenv("TMPDIR", shortTempDir(t))

	dir, socketPath, err := createHostengineSocketDir()
	if err != nil {
		t.Fatalf("createHostengineSocketDir() error = %v", err)
	}
	defer os.RemoveAll(dir)

	if want := filepath.Join(dir, hostengineSocketName); socketPath != want {
		t.Fatalf("socket path = %q, want %q", socketPath, want)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat %q: %v", dir, err)
	}
	if !info.IsDir() || info.Mode().Perm() != 0o700 {
		t.Fatalf("directory mode = %v, want a directory with mode 0700", info.Mode())
	}
}

func TestCheckSocketPathLengthBoundary(t *testing.T) {
	const longestUsable = maxUnixSocketPathBytes - 1
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{name: "longest usable path", length: longestUsable},
		{name: "path with no room for the NUL", length: maxUnixSocketPathBytes, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkSocketPathLength(strings.Repeat("a", tt.length))
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateHostengineSocketDirRejectsLongPathAndRemovesDirectory(t *testing.T) {
	parent := filepath.Join(shortTempDir(t), strings.Repeat("a", maxUnixSocketPathBytes))
	if err := os.Mkdir(parent, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", parent)

	dir, socketPath, err := createHostengineSocketDir()
	if err == nil {
		os.RemoveAll(dir)
		t.Fatalf("createHostengineSocketDir() = %q, want an error for a path over the limit", socketPath)
	}
	if dir != "" || socketPath != "" {
		t.Fatalf("returned dir %q and socket path %q with an error, want empty strings", dir, socketPath)
	}
	entries, readErr := os.ReadDir(parent)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("directory left behind after error: %v", entries)
	}
}
