package handlers

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/gorilla/mux"
)

func baseDeviceDeps() deviceDeps {
	return deviceDeps{
		introspect: func() (dcgm.Status, error) {
			return dcgm.Status{}, nil
		},
		getAllDeviceCount: func() (uint, error) {
			return 2, nil
		},
		getSupportedDevices: func() ([]uint, error) {
			return []uint{0, 1}, nil
		},
		getDeviceInfo: func(id uint) (dcgm.Device, error) {
			return dcgm.Device{GPU: id}, nil
		},
		getDeviceStatus: func(uint) (dcgm.DeviceStatus, error) {
			return dcgm.DeviceStatus{}, nil
		},
		healthCheckByGPUID: func(id uint) (dcgm.DeviceHealth, error) {
			return dcgm.DeviceHealth{GPU: id}, nil
		},
	}
}

func deviceRequest(vars map[string]string) (*httptest.ResponseRecorder, *http.Request) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/devices/0", http.NoBody)
	return resp, mux.SetURLVars(req, vars)
}

func TestGetStatusWithDeps(t *testing.T) {
	tests := []struct {
		name       string
		status     dcgm.Status
		err        error
		wantNil    bool
		wantStatus int
		wantText   string
	}{
		{name: "success", status: dcgm.Status{Memory: 42}, wantStatus: http.StatusOK},
		{name: "error", err: errors.New("introspection failed"), wantNil: true, wantStatus: http.StatusInternalServerError, wantText: "introspection failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := baseDeviceDeps()
			deps.introspect = func() (dcgm.Status, error) {
				return tt.status, tt.err
			}
			resp, req := deviceRequest(nil)

			got := getStatusWithDeps(resp, req, deps)

			if (got == nil) != tt.wantNil {
				t.Fatalf("getStatusWithDeps() = %v, want nil %t", got, tt.wantNil)
			}
			if got != nil && *got != tt.status {
				t.Fatalf("getStatusWithDeps() = %+v, want %+v", *got, tt.status)
			}
			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.Code, tt.wantStatus)
			}
			if !strings.Contains(resp.Body.String(), tt.wantText) {
				t.Fatalf("body %q does not contain %q", resp.Body.String(), tt.wantText)
			}
		})
	}
}

func TestGetDeviceInfoWithDeps(t *testing.T) {
	previousUUIDs := uuids
	t.Cleanup(func() {
		uuids = previousUUIDs
	})

	tests := []struct {
		name           string
		vars           map[string]string
		uuidMap        map[string]uint
		count          uint
		countErr       error
		infoErr        error
		wantNil        bool
		wantStatus     int
		wantID         uint
		wantCountCalls int
		wantInfoCalls  int
		wantText       string
	}{
		{name: "numeric ID", vars: map[string]string{"id": "0"}, count: 2, wantStatus: http.StatusOK, wantID: 0, wantCountCalls: 1, wantInfoCalls: 1},
		{name: "invalid numeric ID", vars: map[string]string{"id": "not-a-number"}, wantNil: true, wantStatus: http.StatusBadRequest},
		{name: "overflowing numeric ID", vars: map[string]string{"id": "4294967296"}, wantNil: true, wantStatus: http.StatusBadRequest},
		{name: "known UUID", vars: map[string]string{"uuid": "GPU-known"}, uuidMap: map[string]uint{"GPU-known": 1}, count: 2, wantStatus: http.StatusOK, wantID: 1, wantCountCalls: 1, wantInfoCalls: 1},
		{name: "missing UUID", vars: map[string]string{"uuid": "GPU-missing"}, uuidMap: map[string]uint{}, wantNil: true, wantStatus: http.StatusNotFound},
		{name: "count error", vars: map[string]string{"id": "0"}, countErr: errors.New("count failed"), wantNil: true, wantStatus: http.StatusInternalServerError, wantCountCalls: 1, wantText: "count failed"},
		{name: "out of range", vars: map[string]string{"id": "2"}, count: 2, wantNil: true, wantStatus: http.StatusNotFound, wantCountCalls: 1},
		{name: "device info error", vars: map[string]string{"id": "0"}, count: 2, infoErr: errors.New("device info failed"), wantNil: true, wantStatus: http.StatusInternalServerError, wantCountCalls: 1, wantInfoCalls: 1, wantText: "device info failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids = tt.uuidMap
			countCalls := 0
			infoCalls := 0
			deps := baseDeviceDeps()
			deps.getAllDeviceCount = func() (uint, error) {
				countCalls++
				return tt.count, tt.countErr
			}
			deps.getDeviceInfo = func(id uint) (dcgm.Device, error) {
				infoCalls++
				return dcgm.Device{GPU: id}, tt.infoErr
			}
			resp, req := deviceRequest(tt.vars)

			got := getDeviceInfoWithDeps(resp, req, deps)

			if (got == nil) != tt.wantNil {
				t.Fatalf("getDeviceInfoWithDeps() = %v, want nil %t", got, tt.wantNil)
			}
			if got != nil && got.GPU != tt.wantID {
				t.Fatalf("device GPU = %d, want %d", got.GPU, tt.wantID)
			}
			if countCalls != tt.wantCountCalls {
				t.Fatalf("GetAllDeviceCount calls = %d, want %d", countCalls, tt.wantCountCalls)
			}
			if infoCalls != tt.wantInfoCalls {
				t.Fatalf("GetDeviceInfo calls = %d, want %d", infoCalls, tt.wantInfoCalls)
			}
			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.Code, tt.wantStatus)
			}
			if !strings.Contains(resp.Body.String(), tt.wantText) {
				t.Fatalf("body %q does not contain %q", resp.Body.String(), tt.wantText)
			}
		})
	}
}

func TestGetDeviceStatusWithDeps(t *testing.T) {
	tests := []struct {
		name         string
		supported    []uint
		supportedErr error
		statusErr    error
		wantNil      bool
		wantStatus   int
		wantCalls    int
		wantText     string
	}{
		{name: "success", supported: []uint{0}, wantStatus: http.StatusOK, wantCalls: 1},
		{name: "supported devices error", supportedErr: errors.New("supported failed"), wantNil: true, wantStatus: http.StatusInternalServerError, wantText: "supported failed"},
		{name: "unsupported GPU", supported: []uint{1}, wantNil: true, wantStatus: http.StatusInternalServerError, wantText: "not supported by dcgm"},
		{name: "device status error", supported: []uint{0}, statusErr: errors.New("status failed"), wantNil: true, wantStatus: http.StatusInternalServerError, wantCalls: 1, wantText: "status failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCalls := 0
			deps := baseDeviceDeps()
			deps.getSupportedDevices = func() ([]uint, error) {
				return tt.supported, tt.supportedErr
			}
			deps.getDeviceStatus = func(uint) (dcgm.DeviceStatus, error) {
				statusCalls++
				return dcgm.DeviceStatus{Temperature: 70}, tt.statusErr
			}
			resp, req := deviceRequest(map[string]string{"id": "0"})

			got := getDeviceStatusWithDeps(resp, req, deps)

			if (got == nil) != tt.wantNil {
				t.Fatalf("getDeviceStatusWithDeps() = %v, want nil %t", got, tt.wantNil)
			}
			if got != nil && got.Temperature != 70 {
				t.Fatalf("temperature = %d, want 70", got.Temperature)
			}
			if statusCalls != tt.wantCalls {
				t.Fatalf("GetDeviceStatus calls = %d, want %d", statusCalls, tt.wantCalls)
			}
			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.Code, tt.wantStatus)
			}
			if !strings.Contains(resp.Body.String(), tt.wantText) {
				t.Fatalf("body %q does not contain %q", resp.Body.String(), tt.wantText)
			}
		})
	}
}

func TestGetHealthWithDeps(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantNil    bool
		wantStatus int
		wantText   string
	}{
		{name: "success", wantStatus: http.StatusOK},
		{name: "error", err: errors.New("health failed"), wantNil: true, wantStatus: http.StatusInternalServerError, wantText: "health failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := baseDeviceDeps()
			deps.healthCheckByGPUID = func(id uint) (dcgm.DeviceHealth, error) {
				return dcgm.DeviceHealth{GPU: id, Status: "Healthy"}, tt.err
			}
			resp, req := deviceRequest(map[string]string{"id": "0"})

			got := getHealthWithDeps(resp, req, deps)

			if (got == nil) != tt.wantNil {
				t.Fatalf("getHealthWithDeps() = %v, want nil %t", got, tt.wantNil)
			}
			if got != nil && (got.GPU != 0 || got.Status != "Healthy") {
				t.Fatalf("health = %+v, want GPU 0 Healthy", *got)
			}
			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.Code, tt.wantStatus)
			}
			if !strings.Contains(resp.Body.String(), tt.wantText) {
				t.Fatalf("body %q does not contain %q", resp.Body.String(), tt.wantText)
			}
		})
	}
}

func TestDevicesUuidsWithDeps(t *testing.T) {
	previousUUIDs := uuids
	t.Cleanup(func() {
		uuids = previousUUIDs
	})

	t.Run("success", func(t *testing.T) {
		deps := baseDeviceDeps()
		deps.getDeviceInfo = func(id uint) (dcgm.Device, error) {
			if id == 0 {
				return dcgm.Device{UUID: "GPU-zero"}, nil
			}
			return dcgm.Device{UUID: "GPU-one"}, nil
		}

		devicesUuidsWithDeps(deps)

		if len(uuids) != 2 || uuids["GPU-zero"] != 0 || uuids["GPU-one"] != 1 {
			t.Fatalf("uuids = %#v, want two mapped devices", uuids)
		}
	})

	t.Run("count error", func(t *testing.T) {
		deps := baseDeviceDeps()
		deps.getAllDeviceCount = func() (uint, error) {
			return 0, errors.New("count failed")
		}

		devicesUuidsWithDeps(deps)

		if len(uuids) != 0 {
			t.Fatalf("uuids = %#v, want empty map", uuids)
		}
	})

	t.Run("device info error keeps completed entries", func(t *testing.T) {
		deps := baseDeviceDeps()
		deps.getDeviceInfo = func(id uint) (dcgm.Device, error) {
			if id == 0 {
				return dcgm.Device{UUID: "GPU-zero"}, nil
			}
			return dcgm.Device{}, errors.New("device info failed")
		}

		devicesUuidsWithDeps(deps)

		if len(uuids) != 1 || uuids["GPU-zero"] != 0 {
			t.Fatalf("uuids = %#v, want only completed device", uuids)
		}
	})
}

func TestGetIDRejectsMaxUint32Sentinel(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/devices/4294967295", http.NoBody)

	if got := getId(resp, req, "4294967295"); got != math.MaxUint32 {
		t.Fatalf("getId() = %d, want %d", got, uint(math.MaxUint32))
	}
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusBadRequest)
	}
}

func TestHandlerRenderingWithGetters(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		write           func(http.ResponseWriter, *http.Request)
		wantBody        string
		wantContentType string
		wantEmpty       bool
	}{
		{
			name: "device info JSON",
			path: "/json",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleDeviceInfo(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.Device {
					return &dcgm.Device{GPU: 7}
				})
			},
			wantBody:        `"GPU":7`,
			wantContentType: "application/json",
		},
		{
			name: "device info text",
			path: "/",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleDeviceInfo(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.Device {
					return &dcgm.Device{GPU: 7}
				})
			},
			wantBody: "GPU                    : 7",
		},
		{
			name:      "device info empty",
			path:      "/",
			wantEmpty: true,
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleDeviceInfo(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.Device {
					return nil
				})
			},
		},
		{
			name: "device status JSON",
			path: "/json",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleDeviceStatus(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.DeviceStatus {
					return &dcgm.DeviceStatus{Temperature: 70}
				})
			},
			wantBody:        `"Temperature":70`,
			wantContentType: "application/json",
		},
		{
			name: "device status text",
			path: "/",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleDeviceStatus(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.DeviceStatus {
					return &dcgm.DeviceStatus{Temperature: 70}
				})
			},
			wantBody: "Temperature (°C)        : 70",
		},
		{
			name:      "device status empty",
			path:      "/",
			wantEmpty: true,
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleDeviceStatus(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.DeviceStatus {
					return nil
				})
			},
		},
		{
			name: "health JSON",
			path: "/json",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleHealth(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.DeviceHealth {
					return &dcgm.DeviceHealth{Status: "Healthy"}
				})
			},
			wantBody:        `"Status":"Healthy"`,
			wantContentType: "application/json",
		},
		{
			name: "health text",
			path: "/",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleHealth(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.DeviceHealth {
					return &dcgm.DeviceHealth{Status: "Healthy"}
				})
			},
			wantBody: "Status             : Healthy",
		},
		{
			name:      "health empty",
			path:      "/",
			wantEmpty: true,
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleHealth(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.DeviceHealth {
					return nil
				})
			},
		},
		{
			name: "status JSON",
			path: "/json",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleStatus(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.Status {
					return &dcgm.Status{Memory: 42}
				})
			},
			wantBody:        `"Memory":42`,
			wantContentType: "application/json",
		},
		{
			name: "status text",
			path: "/",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleStatus(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.Status {
					return &dcgm.Status{Memory: 42}
				})
			},
			wantBody: "Memory(KB)      : 42",
		},
		{
			name:      "status empty",
			path:      "/",
			wantEmpty: true,
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleStatus(resp, req, func(http.ResponseWriter, *http.Request) *dcgm.Status {
					return nil
				})
			},
		},
		{
			name: "process JSON",
			path: "/json",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleProcessInfo(resp, req, func(http.ResponseWriter, *http.Request) []dcgm.ProcessInfo {
					return []dcgm.ProcessInfo{{PID: 7}}
				})
			},
			wantBody:        `"PID":7`,
			wantContentType: "application/json",
		},
		{
			name: "process text",
			path: "/",
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleProcessInfo(resp, req, func(http.ResponseWriter, *http.Request) []dcgm.ProcessInfo {
					return []dcgm.ProcessInfo{{PID: 7}}
				})
			},
			wantBody: "PID                          : 7",
		},
		{
			name:      "process empty",
			path:      "/",
			wantEmpty: true,
			write: func(resp http.ResponseWriter, req *http.Request) {
				handleProcessInfo(resp, req, func(http.ResponseWriter, *http.Request) []dcgm.ProcessInfo {
					return nil
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, http.NoBody)

			tt.write(resp, req)

			if tt.wantEmpty {
				if resp.Body.Len() != 0 {
					t.Fatalf("body = %q, want empty", resp.Body.String())
				}
				if got := resp.Header().Get("Content-Type"); got != "" {
					t.Fatalf("Content-Type = %q, want empty", got)
				}
				return
			}

			if !strings.Contains(resp.Body.String(), tt.wantBody) {
				t.Fatalf("body %q does not contain %q", resp.Body.String(), tt.wantBody)
			}
			if tt.wantContentType != "" {
				if got := resp.Header().Get("Content-Type"); got != tt.wantContentType {
					t.Fatalf("Content-Type = %q, want %q", got, tt.wantContentType)
				}
			}
		})
	}
}

func TestPublicHandlersRejectInvalidIdentifiersWithoutDCGM(t *testing.T) {
	previousUUIDs := uuids
	uuids = map[string]uint{}
	t.Cleanup(func() {
		uuids = previousUUIDs
	})

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		vars       map[string]string
		wantStatus int
	}{
		{name: "device info ID", handler: DeviceInfo, vars: map[string]string{"id": "invalid"}, wantStatus: http.StatusBadRequest},
		{name: "device status ID", handler: DeviceStatus, vars: map[string]string{"id": "invalid"}, wantStatus: http.StatusBadRequest},
		{name: "process ID", handler: ProcessInfo, vars: map[string]string{"pid": "invalid"}, wantStatus: http.StatusBadRequest},
		{name: "health ID", handler: Health, vars: map[string]string{"id": "invalid"}, wantStatus: http.StatusBadRequest},
		{name: "device info UUID", handler: DeviceInfoByUuid, vars: map[string]string{"uuid": "GPU-missing"}, wantStatus: http.StatusNotFound},
		{name: "device status UUID", handler: DeviceStatusByUuid, vars: map[string]string{"uuid": "GPU-missing"}, wantStatus: http.StatusNotFound},
		{name: "health UUID", handler: HealthByUuid, vars: map[string]string{"uuid": "GPU-missing"}, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, req := deviceRequest(tt.vars)

			tt.handler(resp, req)

			if resp.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeviceOperationsByUUIDWithDeps(t *testing.T) {
	previousUUIDs := uuids
	t.Cleanup(func() {
		uuids = previousUUIDs
	})
	uuids = map[string]uint{"GPU-known": 1}

	tests := []struct {
		name string
		call func(http.ResponseWriter, *http.Request, deviceDeps) bool
	}{
		{
			name: "device status",
			call: func(resp http.ResponseWriter, req *http.Request, deps deviceDeps) bool {
				return getDeviceStatusWithDeps(resp, req, deps) != nil
			},
		},
		{
			name: "health",
			call: func(resp http.ResponseWriter, req *http.Request, deps deviceDeps) bool {
				return getHealthWithDeps(resp, req, deps) != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, req := deviceRequest(map[string]string{"uuid": "GPU-known"})

			if !tt.call(resp, req, baseDeviceDeps()) {
				t.Fatal("operation returned nil for a known UUID")
			}
			if resp.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", resp.Code, http.StatusOK)
			}
		})
	}
}
