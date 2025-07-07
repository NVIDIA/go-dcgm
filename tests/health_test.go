package tests

import (
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestHealthCheck demonstrates GPU health checking functionality
// This is equivalent to the health sample but runs for a limited time
func TestHealthCheck(t *testing.T) {
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
		t.Skip("No supported GPUs found for health checking")
	}

	// Monitor health for a few seconds instead of indefinitely
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	timeout := time.After(6 * time.Second)
	checkCount := 0

	for {
		select {
		case <-ticker.C:
			for _, gpu := range gpus {
				h, err := dcgm.HealthCheckByGpuId(gpu)
				if err != nil {
					t.Errorf("Failed to get health status for GPU %d: %v", gpu, err)
					continue
				}

				t.Logf("GPU %d Health Check:", gpu)
				t.Logf("  Status: %s", h.Status)

				for _, watch := range h.Watches {
					t.Logf("  Watch Type: %s", watch.Type)
					t.Logf("  Watch Status: %s", watch.Status)
					if watch.Error != "" {
						t.Logf("  Watch Error: %s", watch.Error)
					}
				}

				// Basic validation
				if h.Status == "" {
					t.Errorf("GPU %d has empty health status", gpu)
				}
			}
			checkCount++

		case <-timeout:
			t.Logf("Health monitoring completed after %d checks", checkCount)
			return
		}
	}
}

// TestHealthCheckSingle demonstrates a single health check
func TestHealthCheckSingle(t *testing.T) {
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

	// Test first GPU
	gpu := gpus[0]
	h, err := dcgm.HealthCheckByGpuId(gpu)
	if err != nil {
		t.Fatalf("Failed to get health status for GPU %d: %v", gpu, err)
	}

	t.Logf("GPU %d Health Status: %s", gpu, h.Status)

	if len(h.Watches) == 0 {
		t.Logf("No health watches configured for GPU %d", gpu)
	} else {
		t.Logf("Health watches for GPU %d:", gpu)
		for i, watch := range h.Watches {
			t.Logf("  Watch %d:", i+1)
			t.Logf("    Type: %s", watch.Type)
			t.Logf("    Status: %s", watch.Status)
			if watch.Error != "" {
				t.Logf("    Error: %s", watch.Error)
			}
		}
	}

	// Basic assertions
	if h.Status == "" {
		t.Error("Health status is empty")
	}
}

// TestHealthCheckAllGPUs demonstrates health checking for all GPUs
func TestHealthCheckAllGPUs(t *testing.T) {
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

	healthyGPUs := 0
	unhealthyGPUs := 0

	for _, gpu := range gpus {
		h, err := dcgm.HealthCheckByGpuId(gpu)
		if err != nil {
			t.Errorf("Failed to get health status for GPU %d: %v", gpu, err)
			continue
		}

		t.Logf("GPU %d: %s", gpu, h.Status)

		// Count healthy vs unhealthy
		if h.Status == "Healthy" || h.Status == "OK" {
			healthyGPUs++
		} else {
			unhealthyGPUs++
			t.Logf("GPU %d is not healthy: %s", gpu, h.Status)

			// Log any watch errors
			for _, watch := range h.Watches {
				if watch.Error != "" {
					t.Logf("  Watch %s error: %s", watch.Type, watch.Error)
				}
			}
		}
	}

	t.Logf("Health summary: %d healthy, %d unhealthy GPUs", healthyGPUs, unhealthyGPUs)

	// We expect at least some GPUs to be available
	if healthyGPUs == 0 && unhealthyGPUs == 0 {
		t.Error("No GPU health status could be determined")
	}
}
