package tests

import (
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestDeviceMonitoring demonstrates device monitoring functionality
// This is equivalent to the dmon sample but runs for a limited time
func TestDeviceMonitoring(t *testing.T) {
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
		t.Skip("No supported GPUs found for monitoring")
	}

	t.Log("# gpu   pwr  temp    sm   mem   enc   dec  mclk  pclk")
	t.Log("# Idx     W     C     %     %     %     %   MHz   MHz")

	// Monitor for a few seconds instead of indefinitely
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	timeout := time.After(5 * time.Second)
	sampleCount := 0

	for {
		select {
		case <-ticker.C:
			for _, gpu := range gpus {
				st, err := dcgm.GetDeviceStatus(gpu)
				if err != nil {
					t.Errorf("Failed to get device status for GPU %d: %v", gpu, err)
					continue
				}

				t.Logf("%5d %5d %5d %5d %5d %5d %5d %5d %5d",
					gpu, int64(st.Power), st.Temperature, st.Utilization.GPU, st.Utilization.Memory,
					st.Utilization.Encoder, st.Utilization.Decoder, st.Clocks.Memory, st.Clocks.Cores)

				// Basic validation
				if st.Temperature < 0 || st.Temperature > 150 {
					t.Errorf("GPU %d temperature out of expected range: %d°C", gpu, st.Temperature)
				}
				if st.Utilization.GPU < 0 || st.Utilization.GPU > 100 {
					t.Errorf("GPU %d utilization out of range: %d%%", gpu, st.Utilization.GPU)
				}
			}
			sampleCount++

		case <-timeout:
			t.Logf("Monitoring completed after %d samples", sampleCount)
			return
		}
	}
}

// TestDeviceStatusSingle demonstrates getting device status for a single GPU
func TestDeviceStatusSingle(t *testing.T) {
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
	st, err := dcgm.GetDeviceStatus(gpu)
	if err != nil {
		t.Fatalf("Failed to get device status for GPU %d: %v", gpu, err)
	}

	t.Logf("GPU %d Status:", gpu)
	t.Logf("  Power: %d W", int64(st.Power))
	t.Logf("  Temperature: %d°C", st.Temperature)
	t.Logf("  GPU Utilization: %d%%", st.Utilization.GPU)
	t.Logf("  Memory Utilization: %d%%", st.Utilization.Memory)
	t.Logf("  Encoder Utilization: %d%%", st.Utilization.Encoder)
	t.Logf("  Decoder Utilization: %d%%", st.Utilization.Decoder)
	t.Logf("  Memory Clock: %d MHz", st.Clocks.Memory)
	t.Logf("  Core Clock: %d MHz", st.Clocks.Cores)

	// Validate ranges
	if st.Temperature < 0 || st.Temperature > 150 {
		t.Errorf("Temperature out of expected range: %d°C", st.Temperature)
	}
	if st.Utilization.GPU < 0 || st.Utilization.GPU > 100 {
		t.Errorf("GPU utilization out of range: %d%%", st.Utilization.GPU)
	}
	if st.Utilization.Memory < 0 || st.Utilization.Memory > 100 {
		t.Errorf("Memory utilization out of range: %d%%", st.Utilization.Memory)
	}
}

// TestDeviceStatusMultipleSamples demonstrates taking multiple samples over time
func TestDeviceStatusMultipleSamples(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping multiple samples test in short mode")
	}

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

	// Take samples every 500ms for 3 seconds
	gpu := gpus[0]
	samples := make([]dcgm.DeviceStatus, 0, 6)

	for i := 0; i < 6; i++ {
		st, err := dcgm.GetDeviceStatus(gpu)
		if err != nil {
			t.Errorf("Failed to get device status sample %d: %v", i, err)
			continue
		}
		samples = append(samples, st)
		time.Sleep(500 * time.Millisecond)
	}

	t.Logf("Collected %d samples for GPU %d", len(samples), gpu)

	// Analyze samples for consistency
	if len(samples) > 1 {
		firstTemp := samples[0].Temperature
		tempVariation := false
		for _, sample := range samples[1:] {
			if abs64(sample.Temperature-firstTemp) > 5 { // Allow 5°C variation
				tempVariation = true
				break
			}
		}

		if !tempVariation {
			t.Logf("Temperature remained stable around %d°C", firstTemp)
		} else {
			t.Logf("Temperature variation detected across samples")
		}
	}
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
