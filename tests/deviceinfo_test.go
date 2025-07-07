package tests

import (
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestDeviceInfo demonstrates getting device information from all GPUs
// This is equivalent to the deviceInfo sample
func TestDeviceInfoTest(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		t.Fatalf("Failed to get device count: %v", err)
	}

	t.Logf("Found %d devices", count)

	for i := uint(0); i < count; i++ {
		deviceInfo, err := dcgm.GetDeviceInfo(i)
		if err != nil {
			t.Errorf("Failed to get device info for GPU %d: %v", i, err)
			continue
		}

		// Log device information
		t.Logf("Device %d Information:", i)
		t.Logf("  Driver Version: %s", deviceInfo.Identifiers.DriverVersion)
		t.Logf("  GPU: %d", deviceInfo.GPU)
		t.Logf("  DCGMSupported: %v", deviceInfo.DCGMSupported)
		t.Logf("  UUID: %s", deviceInfo.UUID)
		t.Logf("  Brand: %s", deviceInfo.Identifiers.Brand)
		t.Logf("  Model: %s", deviceInfo.Identifiers.Model)
		t.Logf("  Serial Number: %s", deviceInfo.Identifiers.Serial)

		if deviceInfo.Identifiers.Vbios != "" {
			t.Logf("  Vbios: %s", deviceInfo.Identifiers.Vbios)
		}

		t.Logf("  InforomImage Version: %s", deviceInfo.Identifiers.InforomImageVersion)
		t.Logf("  Bus ID: %s", deviceInfo.PCI.BusID)

		if deviceInfo.PCI.BAR1 != 0 {
			t.Logf("  BAR1 (MB): %d", deviceInfo.PCI.BAR1)
		}

		if deviceInfo.PCI.FBTotal != 0 {
			t.Logf("  FrameBuffer Memory (MB): %d", deviceInfo.PCI.FBTotal)
		}

		if deviceInfo.PCI.Bandwidth != 0 {
			t.Logf("  Bandwidth (MB/s): %d", deviceInfo.PCI.Bandwidth)
		}

		if deviceInfo.Power != 0 {
			t.Logf("  Power (W): %d", deviceInfo.Power)
		}

		if deviceInfo.CPUAffinity != "" {
			t.Logf("  CPUAffinity: %s", deviceInfo.CPUAffinity)
		}

		// Log P2P topology if available
		if len(deviceInfo.Topology) > 0 {
			t.Logf("  P2P Available:")
			for _, topo := range deviceInfo.Topology {
				t.Logf("    GPU%d - (BusID)%s - %p", topo.GPU, topo.BusID, topo.Link.PCIPaths)
			}
		} else {
			t.Logf("  P2P Available: None")
		}

		// Basic assertions to ensure we got valid data
		if deviceInfo.UUID == "" {
			t.Errorf("Device %d has empty UUID", i)
		}
		if deviceInfo.Identifiers.Brand == "" {
			t.Errorf("Device %d has empty brand", i)
		}
		if deviceInfo.PCI.BusID == "" {
			t.Errorf("Device %d has empty bus ID", i)
		}
	}
}

// TestDeviceInfoWithConnection demonstrates connecting to a standalone hostengine
func TestDeviceInfoWithConnection(t *testing.T) {
	// Skip this test if we're not testing with a specific connection
	if testing.Short() {
		t.Skip("Skipping connection test in short mode")
	}

	connectAddr := "localhost"
	isSocket := "0"

	cleanup, err := dcgm.Init(dcgm.Standalone, connectAddr, isSocket)
	if err != nil {
		t.Skipf("Failed to connect to standalone hostengine at %s: %v", connectAddr, err)
	}
	defer cleanup()

	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		t.Fatalf("Failed to get device count: %v", err)
	}

	t.Logf("Connected to standalone hostengine, found %d devices", count)

	// Just test first device if available
	if count > 0 {
		deviceInfo, err := dcgm.GetDeviceInfo(0)
		if err != nil {
			t.Errorf("Failed to get device info for GPU 0: %v", err)
		} else {
			t.Logf("First device UUID: %s", deviceInfo.UUID)
		}
	}
}
