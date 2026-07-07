package tests

import (
	"strconv"
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestTopology demonstrates getting GPU topology information
// This is equivalent to the topology sample
func TestTopology(t *testing.T) {
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
		t.Skip("No supported GPUs found for topology testing")
	}

	t.Log("GPU Topology Matrix:")
	t.Log("Legend:")
	t.Log(" X    = Self")
	t.Log(" SYS  = Connection traversing PCIe as well as the SMP interconnect between NUMA nodes (e.g., QPI/UPI)")
	t.Log(" NODE = Connection traversing PCIe as well as the interconnect between PCIe Host Bridges within a NUMA node")
	t.Log(" PHB  = Connection traversing PCIe as well as a PCIe Host Bridge (typically the CPU)")
	t.Log(" PXB  = Connection traversing multiple PCIe switches (without traversing the PCIe Host Bridge)")
	t.Log(" PIX  = Connection traversing a single PCIe switch")
	t.Log(" PSB  = Connection traversing a single on-board PCIe switch")
	t.Log(" NV#  = Connection traversing a bonded set of # NVLinks")

	// Print header
	var header strings.Builder
	header.WriteString("     ")
	for _, gpu := range gpus {
		header.WriteString(" GPU")
		header.WriteString(strconv.FormatUint(uint64(gpu), 10))
	}
	header.WriteString(" CPUAffinity")
	t.Log(header.String())

	for i := range gpus {
		topo, err := dcgm.GetDeviceTopology(gpus[i])
		if err != nil {
			t.Errorf("Failed to get topology for GPU %d: %v", gpus[i], err)
			continue
		}

		gpuTopo := make(map[uint]string, len(gpus))
		for _, gpu := range gpus {
			gpuTopo[gpu] = ""
		}

		// Fill topology information
		for _, topoInfo := range topo {
			if _, ok := gpuTopo[topoInfo.GPU]; ok {
				gpuTopo[topoInfo.GPU] = topoInfo.Link.PCIPaths()
			}
		}

		// Self connection
		gpuTopo[gpus[i]] = "X"

		// Add topology info to row
		var row strings.Builder
		row.WriteString("GPU")
		row.WriteString(strconv.FormatUint(uint64(gpus[i]), 10))
		for _, gpu := range gpus {
			path := gpuTopo[gpu]
			if path == "" {
				path = "N/A"
			}
			row.WriteString("  ")
			row.WriteString(path)
		}

		// Get device info for CPU affinity
		deviceInfo, err := dcgm.GetDeviceInfo(gpus[i])
		if err != nil {
			t.Errorf("Failed to get device info for GPU %d: %v", gpus[i], err)
			row.WriteString("  N/A")
		} else {
			if deviceInfo.CPUAffinity != "" {
				row.WriteString("  ")
				row.WriteString(deviceInfo.CPUAffinity)
			} else {
				row.WriteString("  N/A")
			}
		}

		t.Log(row.String())
	}
}

// TestDeviceTopologySingle demonstrates getting topology for a single GPU
func TestDeviceTopologySingle(t *testing.T) {
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
	topo, err := dcgm.GetDeviceTopology(gpu)
	if err != nil {
		t.Fatalf("Failed to get topology for GPU %d: %v", gpu, err)
	}

	t.Logf("Topology information for GPU %d:", gpu)
	if len(topo) == 0 {
		t.Logf("  No topology connections found")
	} else {
		for i, topoInfo := range topo {
			t.Logf("  Connection %d:", i+1)
			t.Logf("    Remote GPU: %d", topoInfo.GPU)
			t.Logf("    Remote Bus ID: %s", topoInfo.BusID)
			t.Logf("    Link Type: %s", topoInfo.Link.PCIPaths())
		}
	}

	// Get CPU affinity for this GPU
	deviceInfo, err := dcgm.GetDeviceInfo(gpu)
	if err != nil {
		t.Errorf("Failed to get device info for GPU %d: %v", gpu, err)
	} else {
		t.Logf("  CPU Affinity: %s", deviceInfo.CPUAffinity)
	}
}

// TestTopologyConnections demonstrates analyzing topology connections
func TestTopologyConnections(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}

	if len(gpus) < 2 {
		t.Skip("Need at least 2 GPUs to test topology connections")
	}

	connectionTypes := make(map[string]int)
	totalConnections := 0

	for _, gpu := range gpus {
		topo, err := dcgm.GetDeviceTopology(gpu)
		if err != nil {
			t.Errorf("Failed to get topology for GPU %d: %v", gpu, err)
			continue
		}

		for _, topoInfo := range topo {
			linkType := topoInfo.Link.PCIPaths()
			connectionTypes[linkType]++
			totalConnections++
		}
	}

	t.Logf("Topology Analysis:")
	t.Logf("  Total connections: %d", totalConnections)
	t.Logf("  Connection types:")
	for linkType, count := range connectionTypes {
		t.Logf("    %s: %d connections", linkType, count)
	}

	// Basic validation
	if totalConnections == 0 {
		t.Log("No topology connections found between GPUs")
	} else {
		t.Logf("Found %d topology connections", totalConnections)
	}
}

// TestTopologyConsistency demonstrates checking topology consistency
func TestTopologyConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping topology consistency test in short mode")
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

	if len(gpus) < 2 {
		t.Skip("Need at least 2 GPUs to test topology consistency")
	}

	// Check if topology is symmetric
	inconsistencies := 0
	for i, gpu1 := range gpus {
		topo1, err := dcgm.GetDeviceTopology(gpu1)
		if err != nil {
			t.Errorf("Failed to get topology for GPU %d: %v", gpu1, err)
			continue
		}

		for j, gpu2 := range gpus {
			if i >= j {
				continue // Skip self and avoid double-checking
			}

			// Check if gpu1 -> gpu2 connection exists
			found1to2 := false
			var link1to2 string
			for _, topoInfo := range topo1 {
				if topoInfo.GPU == gpu2 {
					found1to2 = true
					link1to2 = topoInfo.Link.PCIPaths()
					break
				}
			}

			// Check if gpu2 -> gpu1 connection exists
			topo2, err := dcgm.GetDeviceTopology(gpu2)
			if err != nil {
				t.Errorf("Failed to get topology for GPU %d: %v", gpu2, err)
				continue
			}

			found2to1 := false
			var link2to1 string
			for _, topoInfo := range topo2 {
				if topoInfo.GPU == gpu1 {
					found2to1 = true
					link2to1 = topoInfo.Link.PCIPaths()
					break
				}
			}

			// Check consistency
			if found1to2 != found2to1 {
				inconsistencies++
				t.Logf("Inconsistency: GPU%d->GPU%d exists: %t, GPU%d->GPU%d exists: %t",
					gpu1, gpu2, found1to2, gpu2, gpu1, found2to1)
			} else if found1to2 && found2to1 && link1to2 != link2to1 {
				inconsistencies++
				t.Logf("Link type inconsistency: GPU%d->GPU%d: %s, GPU%d->GPU%d: %s",
					gpu1, gpu2, link1to2, gpu2, gpu1, link2to1)
			}
		}
	}

	if inconsistencies == 0 {
		t.Log("Topology is consistent across all GPUs")
	} else {
		t.Logf("Found %d topology inconsistencies", inconsistencies)
	}
}
