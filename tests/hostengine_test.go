package tests

import (
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestHostEngineStatus demonstrates DCGM host engine introspection
// This is equivalent to the hostengineStatus sample
func TestHostEngineStatus(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	st, err := dcgm.Introspect()
	if err != nil {
		t.Fatalf("Failed to introspect host engine: %v", err)
	}

	t.Logf("Host Engine Status:")
	t.Logf("  Memory: %v KB", st.Memory)
	t.Logf("  CPU: %.2f%%", st.CPU)

	// Basic validation
	if st.Memory < 0 {
		t.Error("Memory usage cannot be negative")
	}
	if st.CPU < 0 || st.CPU > 100 {
		t.Errorf("CPU usage out of expected range: %.2f%%", st.CPU)
	}

	// Log some insights
	if st.Memory > 100000 { // > 100MB
		t.Logf("Host engine is using significant memory: %v KB", st.Memory)
	}
	if st.CPU > 50 {
		t.Logf("Host engine is using significant CPU: %.2f%%", st.CPU)
	}
}

// TestHostEngineStatusMultipleSamples demonstrates taking multiple introspection samples
func TestHostEngineStatusMultipleSamples(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping multiple samples test in short mode")
	}

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	samples := 3
	memoryUsages := make([]int64, 0, samples)
	cpuUsages := make([]float64, 0, samples)

	for i := 0; i < samples; i++ {
		st, err := dcgm.Introspect()
		if err != nil {
			t.Errorf("Failed to introspect host engine sample %d: %v", i+1, err)
			continue
		}

		memoryUsages = append(memoryUsages, st.Memory)
		cpuUsages = append(cpuUsages, st.CPU)

		t.Logf("Sample %d - Memory: %v KB, CPU: %.2f%%", i+1, st.Memory, st.CPU)
	}

	if len(memoryUsages) > 1 {
		// Check for significant memory changes
		minMem := memoryUsages[0]
		maxMem := memoryUsages[0]

		for _, mem := range memoryUsages[1:] {
			if mem < minMem {
				minMem = mem
			}
			if mem > maxMem {
				maxMem = mem
			}
		}

		if maxMem-minMem > 1000 { // More than 1MB difference
			t.Logf("Memory usage varied significantly: %v KB to %v KB", minMem, maxMem)
		} else {
			t.Logf("Memory usage remained stable around %v KB", memoryUsages[0])
		}
	}

	if len(cpuUsages) > 1 {
		// Check for significant CPU changes
		minCPU := cpuUsages[0]
		maxCPU := cpuUsages[0]

		for _, cpu := range cpuUsages[1:] {
			if cpu < minCPU {
				minCPU = cpu
			}
			if cpu > maxCPU {
				maxCPU = cpu
			}
		}

		if maxCPU-minCPU > 10 { // More than 10% difference
			t.Logf("CPU usage varied significantly: %.2f%% to %.2f%%", minCPU, maxCPU)
		} else {
			t.Logf("CPU usage remained stable around %.2f%%", cpuUsages[0])
		}
	}
}

// TestHostEngineStatusWithLoad demonstrates introspection while performing operations
func TestHostEngineStatusWithLoad(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	// Get baseline status
	baselineSt, err := dcgm.Introspect()
	if err != nil {
		t.Fatalf("Failed to get baseline introspection: %v", err)
	}

	t.Logf("Baseline - Memory: %v KB, CPU: %.2f%%", baselineSt.Memory, baselineSt.CPU)

	// Perform some operations to potentially increase load
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Logf("Failed to get supported devices: %v", err)
	} else {
		// Get device info for all GPUs
		for _, gpu := range gpus {
			_, err = dcgm.GetDeviceInfo(gpu)
			if err != nil {
				t.Logf("Failed to get device info for GPU %d: %v", gpu, err)
			}
		}
	}

	// Get status after operations
	loadedSt, err := dcgm.Introspect()
	if err != nil {
		t.Fatalf("Failed to get loaded introspection: %v", err)
	}

	t.Logf("After load - Memory: %v KB, CPU: %.2f%%", loadedSt.Memory, loadedSt.CPU)

	// Compare baseline vs loaded
	memoryDiff := loadedSt.Memory - baselineSt.Memory
	cpuDiff := loadedSt.CPU - baselineSt.CPU

	t.Logf("Differences - Memory: %+d KB, CPU: %+.2f%%", memoryDiff, cpuDiff)

	// Basic checks
	if loadedSt.Memory == 0 {
		t.Error("Memory usage should not be zero")
	}
	if loadedSt.CPU < 0 {
		t.Error("CPU usage should not be negative")
	}
}
