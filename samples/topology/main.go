package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

const (
	legend = `
Legend:
 X    = Self
 SYS  = Connection traversing PCIe as well as the SMP interconnect between NUMA nodes (e.g., QPI/UPI)
 NODE = Connection traversing PCIe as well as the interconnect between PCIe Host Bridges within a NUMA node
 PHB  = Connection traversing PCIe as well as a PCIe Host Bridge (typically the CPU)
 PXB  = Connection traversing multiple PCIe switches (without traversing the PCIe Host Bridge)
 PIX  = Connection traversing a single PCIe switch
 PSB  = Connection traversing a single on-board PCIe switch
 NV#  = Connection traversing a bonded set of # NVLinks`
)

// main prints a topology matrix similar to nvidia-smi topo -m.
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run reads supported GPUs and prints their peer links and CPU affinities.
func run() error {
	// choose dcgm hostengine running mode
	// 1. dcgm.Embedded
	// 2. dcgm.Standalone
	// 3. dcgm.StartHostengine
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		return fmt.Errorf("initialize DCGM: %w", err)
	}
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		return fmt.Errorf("discover supported GPUs: %w", err)
	}

	headers := make([]string, len(gpus))
	for column, gpu := range gpus {
		headers[column] = fmt.Sprintf("GPU%d", gpu)
	}

	fmt.Print(formatTopologyRow("", headers, "CPUAffinity"))

	for _, gpu := range gpus {
		links, err := dcgm.GetDeviceTopology(gpu)
		if err != nil {
			return fmt.Errorf("read topology for GPU %d: %w", gpu, err)
		}

		deviceInfo, err := dcgm.GetDeviceInfo(gpu)
		if err != nil {
			return fmt.Errorf("read device information for GPU %d: %w", gpu, err)
		}

		fmt.Print(formatTopologyRow(fmt.Sprintf("GPU%d", gpu),
			topologyCells(gpus, gpu, links), deviceInfo.CPUAffinity))
	}

	fmt.Println(legend)
	return nil
}

// topologyCells orders peer links by the displayed GPU IDs, which need not be contiguous.
func topologyCells(gpus []uint, gpu uint, links []dcgm.P2PLink) []string {
	paths := make(map[uint]string, len(links))
	for _, link := range links {
		paths[link.GPU] = link.Link.PCIPaths()
	}

	cells := make([]string, len(gpus))
	for column, peer := range gpus {
		cells[column] = paths[peer]
		if peer == gpu {
			cells[column] = "X"
		}
	}

	return cells
}

// formatTopologyRow gives headers and data the same label, GPU and affinity columns.
func formatTopologyRow(label string, cells []string, affinity string) string {
	var row strings.Builder
	fmt.Fprintf(&row, "%10s", label)
	for _, cell := range cells {
		fmt.Fprintf(&row, "%10s", cell)
	}

	fmt.Fprintf(&row, "  %s\n", affinity)
	return row.String()
}
