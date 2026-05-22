package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/gorilla/mux"
)

func TestGetProcessInfoDestroysWatchedGroup(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/process/123", nil)
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

func TestGetProcessInfoDestroysGroupWhenGetProcessInfoFails(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/process/123", nil)
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
	req := httptest.NewRequest(http.MethodGet, "/process/123", nil)
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
