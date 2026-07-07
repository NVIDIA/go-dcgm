package handlers

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/gorilla/mux"
)

func TestLogRequestErrorIncludesSanitizedContext(t *testing.T) {
	var output bytes.Buffer
	previousOutput, previousFlags := log.Writer(), log.Flags()
	log.SetOutput(&output)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(previousOutput)
		log.SetFlags(previousFlags)
	})

	req := httptest.NewRequest(http.MethodGet, "/devices/0?token=secret", http.NoBody)
	logRequestError(req, errors.New("failed\nforged"))

	got := output.String()
	for _, want := range []string{`method="GET"`, `path="/devices/0"`, `message="failed\nforged"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("log output %q does not contain %q", got, want)
		}
	}
	if strings.Contains(got, "token=secret") {
		t.Fatalf("log output contains query data: %q", got)
	}
}

func TestGetProcessInfoDestroysWatchedGroup(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/process/123", http.NoBody)
	req = mux.SetURLVars(req, map[string]string{"pid": "123"})

	destroyCalled := false
	got := getProcessInfoWithDeps(httptest.NewRecorder(), req, processInfoDeps{
		watchPidFields: func() (dcgm.GroupHandle, error) {
			return dcgm.GroupHandle{}, nil
		},
		getProcessInfo: func(dcgm.GroupHandle, uint) ([]dcgm.ProcessInfo, error) {
			return []dcgm.ProcessInfo{{PID: 123}}, nil
		},
		destroyGroup: func(dcgm.GroupHandle) error {
			destroyCalled = true
			return nil
		},
		sleep: func(time.Duration) {},
	})

	if !destroyCalled {
		t.Fatal("expected DestroyGroup to be called")
	}
	if len(got) != 1 {
		t.Fatalf("expected one process info entry, got %d", len(got))
	}
	if got[0].PID != 123 {
		t.Fatalf("expected process info for PID 123, got %d", got[0].PID)
	}
}

func TestIsJsonUsesPathSuffix(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{url: "/", want: false},
		{url: "/dcgm/status/json", want: true},
		{url: "/dcgm/status/notjson", want: false},
		{url: "/dcgm/status", want: false},
		{url: "/dcgm/status?format=json", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)
			if got := isJson(req); got != tt.want {
				t.Fatalf("isJson(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestGetProcessInfoDestroysGroupWhenGetProcessInfoFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/process/123", http.NoBody)
	req = mux.SetURLVars(req, map[string]string{"pid": "123"})
	rr := httptest.NewRecorder()

	destroyCalled := false
	got := getProcessInfoWithDeps(rr, req, processInfoDeps{
		watchPidFields: func() (dcgm.GroupHandle, error) {
			return dcgm.GroupHandle{}, nil
		},
		getProcessInfo: func(dcgm.GroupHandle, uint) ([]dcgm.ProcessInfo, error) {
			return nil, errors.New("forced GetProcessInfo failure")
		},
		destroyGroup: func(dcgm.GroupHandle) error {
			destroyCalled = true
			return nil
		},
		sleep: func(time.Duration) {},
	})

	if !destroyCalled {
		t.Fatal("expected DestroyGroup to be called even when GetProcessInfo fails")
	}
	if got != nil {
		t.Fatalf("expected nil process info on error, got %v", got)
	}
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestGetProcessInfoDoesNotDestroyGroupWhenWatchFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/process/123", http.NoBody)
	req = mux.SetURLVars(req, map[string]string{"pid": "123"})
	rr := httptest.NewRecorder()

	destroyCalled := false
	getProcessInfoWithDeps(rr, req, processInfoDeps{
		watchPidFields: func() (dcgm.GroupHandle, error) {
			return dcgm.GroupHandle{}, errors.New("forced WatchPidFields failure")
		},
		getProcessInfo: func(dcgm.GroupHandle, uint) ([]dcgm.ProcessInfo, error) {
			t.Fatal("GetProcessInfo must not be called when watch fails")
			return nil, nil
		},
		destroyGroup: func(dcgm.GroupHandle) error {
			destroyCalled = true
			return nil
		},
		sleep: func(time.Duration) {
			t.Fatal("sleep must not be called when watch fails")
		},
	})

	if destroyCalled {
		t.Fatal("DestroyGroup must not be called when WatchPidFields returned no group")
	}
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}
