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
	"slices"
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

var deviceStatusFields = [deviceStatusFieldCount]Short{
	devicePower:              C.DCGM_FI_DEV_POWER_USAGE,
	deviceTemperature:        C.DCGM_FI_DEV_GPU_TEMP,
	deviceGPUUtilization:     C.DCGM_FI_DEV_GPU_UTIL,
	deviceMemoryUtilization:  C.DCGM_FI_DEV_MEM_COPY_UTIL,
	deviceEncoderUtilization: C.DCGM_FI_DEV_ENC_UTIL,
	deviceDecoderUtilization: C.DCGM_FI_DEV_DEC_UTIL,
	deviceSMClock:            C.DCGM_FI_DEV_SM_CLOCK,
	deviceMemoryClock:        C.DCGM_FI_DEV_MEM_CLOCK,
	deviceBAR1Used:           C.DCGM_FI_DEV_BAR1_USED,
	devicePCIeRxThroughput:   C.DCGM_FI_DEV_PCIE_RX_THROUGHPUT,
	devicePCIeTxThroughput:   C.DCGM_FI_DEV_PCIE_TX_THROUGHPUT,
	devicePCIeReplay:         C.DCGM_FI_DEV_PCIE_REPLAY_COUNTER,
	deviceFBUsed:             C.DCGM_FI_DEV_FB_USED,
	deviceSingleBitErrors:    C.DCGM_FI_DEV_ECC_SBE_AGG_TOTAL,
	deviceDoubleBitErrors:    C.DCGM_FI_DEV_ECC_DBE_AGG_TOTAL,
	devicePerformanceState:   C.DCGM_FI_DEV_PSTATE,
	deviceFanSpeed:           C.DCGM_FI_DEV_FAN_SPEED,
}

// orderDeviceStatusValues puts samples in the order used to build DeviceStatus.
// DCGM can return samples in a different order from the request.
func orderDeviceStatusValues(values []FieldValue_v1) ([deviceStatusFieldCount]*FieldValue_v1, error) {
	var ordered [deviceStatusFieldCount]*FieldValue_v1
	for i := range values {
		value := &values[i]
		index := slices.Index(deviceStatusFields[:], value.FieldID)
		if index < 0 {
			return ordered, fmt.Errorf("unexpected GPU status field %d", value.FieldID)
		}
		if ordered[index] != nil {
			return ordered, fmt.Errorf("duplicate GPU status field %d", value.FieldID)
		}
		ordered[index] = value
	}

	return ordered, nil
}

// deviceStatusFromValues checks that every requested field is present. It
// keeps valid readings when DCGM reports a failure for another field.
func deviceStatusFromValues(values []FieldValue_v1) (DeviceStatus, error) {
	ordered, err := orderDeviceStatusValues(values)
	if err != nil {
		return DeviceStatus{}, err
	}

	var fieldErrors error
	var integers [deviceStatusFieldCount]int64
	var power float64
	for i, fieldID := range deviceStatusFields {
		value := ordered[i]
		if value == nil {
			return DeviceStatus{}, fmt.Errorf("missing GPU status field %d", fieldID)
		}
		wantType := DCGM_FT_INT64
		if i == devicePower {
			wantType = DCGM_FT_DOUBLE
		}
		if value.Status == DCGM_ST_OK && value.FieldType != wantType {
			return DeviceStatus{}, fmt.Errorf("unexpected type %d for GPU status field %d, expected %d", value.FieldType, fieldID, wantType)
		}
		integer, floating, fieldErr := deviceStatusFieldValue(value, wantType)
		if fieldErr != nil {
			fieldErrors = errors.Join(fieldErrors, fmt.Errorf("read GPU status field %d: %w", fieldID, fieldErr))
		}
		integers[i] = integer
		if i == devicePower {
			power = floating
		}
	}

	return DeviceStatus{
		Power:       power,
		Temperature: integers[deviceTemperature],
		Utilization: UtilizationInfo{
			GPU:     integers[deviceGPUUtilization],
			Memory:  integers[deviceMemoryUtilization],
			Encoder: integers[deviceEncoderUtilization],
			Decoder: integers[deviceDecoderUtilization],
		},
		Memory: MemoryInfo{ECCErrors: ECCErrorsInfo{
			SingleBit: integers[deviceSingleBitErrors],
			DoubleBit: integers[deviceDoubleBitErrors],
		}},
		Clocks: ClockInfo{
			Cores:  integers[deviceSMClock],
			Memory: integers[deviceMemoryClock],
		},
		PCI: PCIStatusInfo{
			BAR1Used: integers[deviceBAR1Used],
			Throughput: PCIThroughputInfo{
				Rx:      integers[devicePCIeRxThroughput],
				Tx:      integers[devicePCIeTxThroughput],
				Replays: integers[devicePCIeReplay],
			},
			FBUsed: integers[deviceFBUsed],
		},
		Performance: PerfState(integers[devicePerformanceState]),
		FanSpeed:    integers[deviceFanSpeed],
	}, fieldErrors
}

// deviceStatusFieldValue reads a payload only when DCGM reports success.
func deviceStatusFieldValue(value *FieldValue_v1, fieldType uint) (integer int64, floating float64, err error) {
	if value.Status != DCGM_ST_OK {
		return deviceStatusBlank(value.Status)
	}
	if fieldType == DCGM_FT_DOUBLE {
		return 0, value.Float64(), nil
	}
	return value.Int64(), 0, nil
}

// deviceStatusBlank picks a blank value for an unavailable reading. It also
// returns the native status as an error when the failure may affect other reads.
func deviceStatusBlank(status int) (integer int64, floating float64, err error) {
	switch status {
	case DCGM_ST_NO_DATA, DCGM_ST_STALE_DATA, DCGM_ST_NOT_WATCHED:
		return DCGM_FT_INT64_BLANK, DCGM_FT_FP64_BLANK, nil
	case DCGM_ST_UNKNOWN_FIELD:
		return DCGM_FT_INT64_NOT_FOUND, DCGM_FT_FP64_NOT_FOUND, nil
	case DCGM_ST_NOT_SUPPORTED, DCGM_ST_FIELD_UNSUPPORTED_BY_API:
		return DCGM_FT_INT64_NOT_SUPPORTED, DCGM_FT_FP64_NOT_SUPPORTED, nil
	case DCGM_ST_NO_PERMISSION, DCGM_ST_REQUIRES_ROOT:
		return DCGM_FT_INT64_NOT_PERMISSIONED, DCGM_FT_FP64_NOT_PERMISSIONED, nil
	default:
		return DCGM_FT_INT64_BLANK, DCGM_FT_FP64_BLANK, &Error{Code: C.dcgmReturn_t(status), msg: fmt.Sprintf("DCGM status %d", status)}
	}
}

func getGPUStatus(gpuID uint) EntityStatus {
	return cgoAdapter{}.getGPUStatus(gpuID)
}

// getGPUStatus reads native GPU status, returning unknown when the query fails.
func (cgoAdapter) getGPUStatus(gpuID uint) EntityStatus {
	var status C.DcgmEntityStatus_t
	result := C.dcgmGetGpuStatus(handle.handle, C.uint(gpuID), &status)
	if result != C.DCGM_ST_OK {
		return EntityStatusUnknown
	}
	return EntityStatus(status)
}

func latestValuesForDevice(gpuId uint) (DeviceStatus, error) {
	return latestValuesForDeviceWithOps(gpuId, cgoAdapter{})
}

func latestValuesForDeviceWithOps(gpuId uint, ops deviceStatusOps) (status DeviceStatus, err error) {
	fieldsName := fmt.Sprintf("devStatusFields%d", rand.Uint64())
	fieldsId, err := ops.fieldGroupCreate(fieldsName, deviceStatusFields[:])
	if err != nil {
		return
	}

	groupName := fmt.Sprintf("devStatus%d", rand.Uint64())
	groupId, err := ops.watchFields(gpuId, fieldsId, groupName)
	if err != nil {
		err = errors.Join(err, ops.fieldGroupDestroy(fieldsId))
		return
	}
	defer func() {
		err = errors.Join(err, ops.fieldGroupDestroy(fieldsId), ops.destroyGroup(groupId))
	}()

	values, err := ops.getLatestValuesForFields(gpuId, deviceStatusFields[:])
	if err != nil {
		return status, err
	}

	status, err = deviceStatusFromValues(values)
	if err != nil {
		err = fmt.Errorf("get status for GPU %d: %w", gpuId, err)
	}

	return
}
