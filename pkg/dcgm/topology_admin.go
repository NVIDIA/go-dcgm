package dcgm

import (
	/*
		#include "dcgm_agent.h"
	*/
	"C"
	"fmt"
)

// GroupTopology describes the CPU affinity and slowest GPU path in a group.
type GroupTopology struct {
	CPUAffinityMask []uint64
	NUMAOptimal     bool
	// SlowestPath retains DCGM's topology bitmask, including combined PCI and NVLink bits.
	SlowestPath uint64
}

// GetGroupTopology reads topology for a DCGM group.
func GetGroupTopology(group GroupHandle) (GroupTopology, error) {
	return cgoAdapter{}.getGroupTopology(group)
}

// getGroupTopology uses the native adapter to read group topology.
func (api cgoAdapter) getGroupTopology(group GroupHandle) (GroupTopology, error) {
	return getGroupTopologyWithCalls(api, group)
}

// getGroupTopologyWithCalls initializes the topology request and preserves native failures.
func getGroupTopologyWithCalls(api administrationNativeCalls, group GroupHandle) (GroupTopology, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmGetGroupTopology"); err != nil {
		return GroupTopology{}, err
	}

	response := C.dcgmGroupTopology_v2{version: C.dcgmGroupTopology_version2}
	if err := administrationNativeResult(api, "dcgmGetGroupTopology", api.groupTopologyNative(group, &response)); err != nil {
		return GroupTopology{}, err
	}

	return groupTopologyFromNative(response), nil
}

// groupTopologyFromNative copies the full CPU affinity mask and topology bits.
func groupTopologyFromNative(response C.dcgmGroupTopology_v2) GroupTopology {
	topology := GroupTopology{
		CPUAffinityMask: make([]uint64, len(response.groupCpuAffinityMask)),
		NUMAOptimal:     response.numaOptimalFlag != 0,
		SlowestPath:     uint64(response.slowestPath),
	}
	for i, word := range response.groupCpuAffinityMask {
		topology.CPUAffinityMask[i] = uint64(word)
	}

	return topology
}

// TopologyHint controls DCGM's topology selection.
type TopologyHint uint64

// TopologyIgnoreHealth allows DCGM to consider unhealthy GPUs.
const TopologyIgnoreHealth TopologyHint = C.DCGM_TOPO_HINT_F_IGNOREHEALTH

// SelectGPUsByTopology selects up to count GPUs near each other. Empty input
// considers every GPU. IDs must fit DCGM's 64-bit mask (0 through 63).
func SelectGPUsByTopology(input []uint, count uint, hints TopologyHint) ([]uint, error) {
	return selectGPUsByTopologyWithCalls(cgoAdapter{}, input, count, hints)
}

// selectGPUsByTopologyWithCalls validates the requested mask and count before native selection.
func selectGPUsByTopologyWithCalls(api gpuSelector, input []uint, count uint, hints TopologyHint) ([]uint, error) {
	mask, err := gpuMask(input)
	if err != nil {
		return nil, err
	}
	if count == 0 || count > 64 {
		return nil, newBadParameterError()
	}

	selected, err := api.selectGPUsByTopology(mask, uint32(count), hints)
	if err != nil {
		return nil, err
	}

	return gpuIDsFromMask(selected), nil
}

// gpuMask converts GPU IDs into the native 64-bit selection mask.
func gpuMask(input []uint) (uint64, error) {
	var mask uint64
	for _, id := range input {
		if id > 63 {
			return 0, fmt.Errorf("GPU ID %d exceeds topology mask limit 63: %w", id, newBadParameterError())
		}
		mask |= uint64(1) << id
	}

	return mask, nil
}

// gpuIDsFromMask returns the GPU IDs represented by a native selection mask.
func gpuIDsFromMask(mask uint64) []uint {
	ids := make([]uint, 0)
	for id := uint(0); id < 64; id++ {
		if mask&(uint64(1)<<id) != 0 {
			ids = append(ids, id)
		}
	}

	return ids
}

// selectGPUsByTopology uses the native adapter for topology selection.
func (api cgoAdapter) selectGPUsByTopology(input uint64, count uint32, hints TopologyHint) (uint64, error) {
	return selectGPUMaskWithCalls(api, input, count, hints)
}

// selectGPUMaskWithCalls checks the session before requesting a native topology selection.
func selectGPUMaskWithCalls(api administrationNativeCalls, input uint64, count uint32, hints TopologyHint) (uint64, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmSelectGpusByTopology"); err != nil {
		return 0, err
	}

	output, status := api.topologySelectionNative(input, count, hints)
	if err := administrationNativeResult(api, "dcgmSelectGpusByTopology", status); err != nil {
		return 0, err
	}

	return output, nil
}

// ChipArchitecture identifies the GPU architecture reported by DCGM.
type ChipArchitecture uint32

// DCGM chip architectures.
const (
	ChipOlder     ChipArchitecture = C.DCGM_CHIP_ARCH_OLDER
	ChipKepler    ChipArchitecture = C.DCGM_CHIP_ARCH_KEPLER
	ChipMaxwell   ChipArchitecture = C.DCGM_CHIP_ARCH_MAXWELL
	ChipPascal    ChipArchitecture = C.DCGM_CHIP_ARCH_PASCAL
	ChipVolta     ChipArchitecture = C.DCGM_CHIP_ARCH_VOLTA
	ChipTuring    ChipArchitecture = C.DCGM_CHIP_ARCH_TURING
	ChipAmpere    ChipArchitecture = C.DCGM_CHIP_ARCH_AMPERE
	ChipAda       ChipArchitecture = C.DCGM_CHIP_ARCH_ADA
	ChipHopper    ChipArchitecture = C.DCGM_CHIP_ARCH_HOPPER
	ChipBlackwell ChipArchitecture = C.DCGM_CHIP_ARCH_BLACKWELL
	ChipUnknown   ChipArchitecture = C.DCGM_CHIP_ARCH_UNKNOWN
)

// GetGPUChipArchitecture reads a GPU's chip architecture.
func GetGPUChipArchitecture(gpuID uint) (ChipArchitecture, error) {
	if uint64(gpuID) > uint64(^uint32(0)) {
		return 0, newBadParameterError()
	}

	return cgoAdapter{}.getGPUChipArchitecture(gpuID)
}

// getGPUChipArchitecture uses the native adapter to read a GPU architecture.
func (api cgoAdapter) getGPUChipArchitecture(gpuID uint) (ChipArchitecture, error) {
	return getGPUChipArchitectureWithCalls(api, gpuID)
}

// getGPUChipArchitectureWithCalls checks the session before converting the native architecture.
func getGPUChipArchitectureWithCalls(api administrationNativeCalls, gpuID uint) (ChipArchitecture, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmGetGpuChipArchitecture"); err != nil {
		return 0, err
	}

	architecture, status := api.chipArchitectureNative(gpuID)
	if err := administrationNativeResult(api, "dcgmGetGpuChipArchitecture", status); err != nil {
		return 0, err
	}

	return architecture, nil
}
