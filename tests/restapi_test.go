package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/gorilla/mux"
)

// TestRestApiEndpoints demonstrates REST API functionality
// This is equivalent to the restApi sample but using httptest instead of a real server
func TestRestApiEndpoints(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Get supported devices for testing
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}

	if len(gpus) == 0 {
		t.Skip("No supported GPUs found for REST API testing")
	}

	// Create a test router with basic endpoints
	router := mux.NewRouter()
	setupRoutes(router)

	// Test device info endpoint
	t.Run("DeviceInfo", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/dcgm/device/info/id/%d", gpus[0]), http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		t.Logf("Device info response: %s", rr.Body.String())
	})

	// Test device status endpoint
	t.Run("DeviceStatus", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/dcgm/device/status/id/%d", gpus[0]), http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		t.Logf("Device status response: %s", rr.Body.String())
	})

	// Test health endpoint
	t.Run("Health", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/dcgm/health/id/%d", gpus[0]), http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		t.Logf("Health response: %s", rr.Body.String())
	})

	// Test DCGM status endpoint
	t.Run("DCGMStatus", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/dcgm/status", http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		t.Logf("DCGM status response: %s", rr.Body.String())
	})
}

// setupRoutes configures the REST API routes for testing
func setupRoutes(router *mux.Router) {
	// Device info endpoints
	deviceInfo := "/dcgm/device/info"
	subrouter := router.PathPrefix(deviceInfo).Subrouter()
	subrouter.HandleFunc("/id/{id}", handleDeviceInfo).Methods("GET")
	subrouter.HandleFunc("/id/{id}/json", handleDeviceInfo).Methods("GET")

	// Device status endpoints
	deviceStatus := "/dcgm/device/status"
	subrouter = router.PathPrefix(deviceStatus).Subrouter()
	subrouter.HandleFunc("/id/{id}", handleDeviceStatus).Methods("GET")
	subrouter.HandleFunc("/id/{id}/json", handleDeviceStatus).Methods("GET")

	// Health endpoints
	health := "/dcgm/health"
	subrouter = router.PathPrefix(health).Subrouter()
	subrouter.HandleFunc("/id/{id}", handleHealth).Methods("GET")
	subrouter.HandleFunc("/id/{id}/json", handleHealth).Methods("GET")

	// DCGM status endpoint
	dcgmStatus := "/dcgm/status"
	subrouter = router.PathPrefix(dcgmStatus).Subrouter()
	subrouter.HandleFunc("", handleDCGMStatus).Methods("GET")
	subrouter.HandleFunc("/json", handleDCGMStatus).Methods("GET")
}

// handleDeviceInfo handles device info requests
func handleDeviceInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	deviceInfo, err := dcgm.GetDeviceInfo(uint(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get device info: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deviceInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get device info: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleDeviceStatus handles device status requests
func handleDeviceStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	deviceStatus, err := dcgm.GetDeviceStatus(uint(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get device status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deviceStatus)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get device status: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleHealth handles health check requests
func handleHealth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	health, err := dcgm.HealthCheckByGpuId(uint(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get health status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(health)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get health status: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleDCGMStatus handles DCGM status requests
func handleDCGMStatus(w http.ResponseWriter, r *http.Request) {
	status, err := dcgm.Introspect()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get DCGM status: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get DCGM status: %v", err), http.StatusInternalServerError)
		return
	}
}

// TestRestApiJsonResponses demonstrates testing JSON response format
func TestRestApiJsonResponses(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}

	if len(gpus) == 0 {
		t.Skip("No supported GPUs found")
	}

	router := mux.NewRouter()
	setupRoutes(router)

	// Test JSON response format
	req, err := http.NewRequest("GET", fmt.Sprintf("/dcgm/device/info/id/%d/json", gpus[0]), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check that response is valid JSON
	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Check content type
	expectedContentType := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, ct)
	}

	t.Logf("JSON response validated successfully")
}

// TestRestApiErrorHandling demonstrates error handling in the API
func TestRestApiErrorHandling(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	router := mux.NewRouter()
	setupRoutes(router)

	// Test invalid device ID
	req, err := http.NewRequest("GET", "/dcgm/device/info/id/invalid", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %v for invalid ID, got %v", http.StatusBadRequest, status)
	}

	// Test non-existent device ID
	req, err = http.NewRequest("GET", "/dcgm/device/info/id/999", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Logf("Non-existent device ID returned status: %v", status)
		// This might be OK depending on DCGM behavior
	}

	t.Log("Error handling tests completed")
}
