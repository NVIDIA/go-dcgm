package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"math/rand"
)

// EntityStatus represents the status of a GPU entity
type EntityStatus uint

const (
	// EntityStatusUnknown - Entity has not been referenced yet
	EntityStatusUnknown EntityStatus = 0
	// EntityStatusOk - Entity is known and OK
	EntityStatusOk EntityStatus = 1
	// EntityStatusUnsupported - Entity is unsupported by DCGM
	EntityStatusUnsupported EntityStatus = 2
	// EntityStatusInaccessible - Entity is inaccessible, usually due to cgroups
	EntityStatusInaccessible EntityStatus = 3
	// EntityStatusLost - Entity has been lost. Usually set from NVML returning NVML_ERROR_GPU_IS_LOST
	EntityStatusLost EntityStatus = 4
	// EntityStatusFake - Entity is a fake, injection-only entity for testing
	EntityStatusFake EntityStatus = 5
	// EntityStatusDisabled - Don't collect values from this GPU
	EntityStatusDisabled EntityStatus = 6
	// EntityStatusDetached - Entity is detached, not good for any uses
	EntityStatusDetached EntityStatus = 7
)

// String returns a string representation of the entity status
func (e EntityStatus) String() string {
	switch e {
	case EntityStatusUnknown:
		return "Unknown"
	case EntityStatusOk:
		return "OK"
	case EntityStatusUnsupported:
		return "Unsupported"
	case EntityStatusInaccessible:
		return "Inaccessible"
	case EntityStatusLost:
		return "Lost"
	case EntityStatusFake:
		return "Fake"
	case EntityStatusDisabled:
		return "Disabled"
	case EntityStatusDetached:
		return "Detached"
	default:
		return fmt.Sprintf("Unknown(%d)", e)
	}
}

// PerfState represents the performance state (P-state) of a GPU
type PerfState uint

const (
	// PerfStateMax represents the highest performance state (P0)
	PerfStateMax = 0

	// PerfStateMin represents the lowest performance state (P15)
	PerfStateMin = 15

	// PerfStateUnknown represents an unknown performance state
	PerfStateUnknown = 32
)

// String returns a string representation of the performance state
func (p PerfState) String() string {
	if p >= PerfStateMax && p <= PerfStateMin {
		return fmt.Sprintf("P%d", p)
	}
	return "Unknown"
}

// UtilizationInfo contains GPU utilization metrics
type UtilizationInfo struct {
	GPU     int64 // %
	Memory  int64 // %
	Encoder int64 // %
	Decoder int64 // %
}

// ECCErrorsInfo contains ECC memory error counts
type ECCErrorsInfo struct {
	SingleBit int64
	DoubleBit int64
}

// MemoryInfo contains GPU memory usage and error information
type MemoryInfo struct {
	GlobalUsed int64
	ECCErrors  ECCErrorsInfo
}

// ClockInfo contains GPU clock frequencies
type ClockInfo struct {
	Cores  int64 // MHz
	Memory int64 // MHz
}

// PCIThroughputInfo contains PCI bus transfer metrics
type PCIThroughputInfo struct {
	Rx      int64 // MB
	Tx      int64 // MB
	Replays int64
}

// PCIStatusInfo contains PCI bus status information
type PCIStatusInfo struct {
	BAR1Used   int64 // MB
	Throughput PCIThroughputInfo
	FBUsed     int64
}

// DeviceStatus contains comprehensive GPU device status information
type DeviceStatus struct {
	Power       float64 // W
	Temperature int64   // °C
	Utilization UtilizationInfo
	Memory      MemoryInfo
	Clocks      ClockInfo
	PCI         PCIStatusInfo
	Performance PerfState
	FanSpeed    int64 // %
}

const (
	devicePower = iota
	deviceTemperature
	deviceGPUUtilization
	deviceMemoryUtilization
	deviceEncoderUtilization
	deviceDecoderUtilization
	deviceSMClock
	deviceMemoryClock
	deviceBAR1Used
	devicePCIeRxThroughput
	devicePCIeTxThroughput
	devicePCIeReplay
	deviceFBUsed
	deviceSingleBitErrors
	deviceDoubleBitErrors
	devicePerformanceState
	deviceFanSpeed
	deviceStatusFieldCount
)

func deviceStatusFromValues(values []FieldValue_v1) DeviceStatus {
	return DeviceStatus{
		Power:       values[devicePower].Float64(),
		Temperature: values[deviceTemperature].Int64(),
		Utilization: UtilizationInfo{
			GPU:     values[deviceGPUUtilization].Int64(),
			Memory:  values[deviceMemoryUtilization].Int64(),
			Encoder: values[deviceEncoderUtilization].Int64(),
			Decoder: values[deviceDecoderUtilization].Int64(),
		},
		Memory: MemoryInfo{ECCErrors: ECCErrorsInfo{
			SingleBit: values[deviceSingleBitErrors].Int64(),
			DoubleBit: values[deviceDoubleBitErrors].Int64(),
		}},
		Clocks: ClockInfo{
			Cores:  values[deviceSMClock].Int64(),
			Memory: values[deviceMemoryClock].Int64(),
		},
		PCI: PCIStatusInfo{
			BAR1Used: values[deviceBAR1Used].Int64(),
			Throughput: PCIThroughputInfo{
				Rx:      values[devicePCIeRxThroughput].Int64(),
				Tx:      values[devicePCIeTxThroughput].Int64(),
				Replays: values[devicePCIeReplay].Int64(),
			},
			FBUsed: values[deviceFBUsed].Int64(),
		},
		Performance: PerfState(values[devicePerformanceState].Int64()),
		FanSpeed:    values[deviceFanSpeed].Int64(),
	}
}

func getGPUStatus(gpuID uint) EntityStatus {
	var status C.DcgmEntityStatus_t
	result := C.dcgmGetGpuStatus(handle.handle, C.uint(gpuID), &status)
	if result != C.DCGM_ST_OK {
		return EntityStatusUnknown
	}
	return EntityStatus(status)
}

type deviceStatusOps struct {
	createFieldGroup  func(string, []Short) (FieldHandle, error)
	watchFields       func(uint, FieldHandle, string) (GroupHandle, error)
	getLatestValues   func(uint, []Short) ([]FieldValue_v1, error)
	destroyFieldGroup func(FieldHandle) error
	destroyGroup      func(GroupHandle) error
}

func latestValuesForDevice(gpuId uint) (DeviceStatus, error) {
	return latestValuesForDeviceWithOps(gpuId, deviceStatusOps{
		createFieldGroup:  FieldGroupCreate,
		watchFields:       WatchFields,
		getLatestValues:   GetLatestValuesForFields,
		destroyFieldGroup: FieldGroupDestroy,
		destroyGroup:      DestroyGroup,
	})
}

func latestValuesForDeviceWithOps(gpuId uint, ops deviceStatusOps) (status DeviceStatus, err error) {
	deviceFields := make([]Short, deviceStatusFieldCount)
	deviceFields[devicePower] = C.DCGM_FI_DEV_POWER_USAGE
	deviceFields[deviceTemperature] = C.DCGM_FI_DEV_GPU_TEMP
	deviceFields[deviceGPUUtilization] = C.DCGM_FI_DEV_GPU_UTIL
	deviceFields[deviceMemoryUtilization] = C.DCGM_FI_DEV_MEM_COPY_UTIL
	deviceFields[deviceEncoderUtilization] = C.DCGM_FI_DEV_ENC_UTIL
	deviceFields[deviceDecoderUtilization] = C.DCGM_FI_DEV_DEC_UTIL
	deviceFields[deviceSMClock] = C.DCGM_FI_DEV_SM_CLOCK
	deviceFields[deviceMemoryClock] = C.DCGM_FI_DEV_MEM_CLOCK
	deviceFields[deviceBAR1Used] = C.DCGM_FI_DEV_BAR1_USED
	deviceFields[devicePCIeRxThroughput] = C.DCGM_FI_DEV_PCIE_RX_THROUGHPUT
	deviceFields[devicePCIeTxThroughput] = C.DCGM_FI_DEV_PCIE_TX_THROUGHPUT
	deviceFields[devicePCIeReplay] = C.DCGM_FI_DEV_PCIE_REPLAY_COUNTER
	deviceFields[deviceFBUsed] = C.DCGM_FI_DEV_FB_USED
	deviceFields[deviceSingleBitErrors] = C.DCGM_FI_DEV_ECC_SBE_AGG_TOTAL
	deviceFields[deviceDoubleBitErrors] = C.DCGM_FI_DEV_ECC_DBE_AGG_TOTAL
	deviceFields[devicePerformanceState] = C.DCGM_FI_DEV_PSTATE
	deviceFields[deviceFanSpeed] = C.DCGM_FI_DEV_FAN_SPEED

	fieldsName := fmt.Sprintf("devStatusFields%d", rand.Uint64())
	fieldsId, err := ops.createFieldGroup(fieldsName, deviceFields)
	if err != nil {
		return
	}

	groupName := fmt.Sprintf("devStatus%d", rand.Uint64())
	groupId, err := ops.watchFields(gpuId, fieldsId, groupName)
	if err != nil {
		err = errors.Join(err, ops.destroyFieldGroup(fieldsId))
		return
	}
	defer func() {
		err = errors.Join(err, ops.destroyFieldGroup(fieldsId), ops.destroyGroup(groupId))
	}()

	values, err := ops.getLatestValues(gpuId, deviceFields)
	if err != nil {
		return status, err
	}

	status = deviceStatusFromValues(values)

	return
}
