package dcgm

/*
#include <stdlib.h>
#include "dcgm_agent.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"unsafe"
)

// JobInt64Summary contains minimum, maximum and average samples, including DCGM blank values.
type JobInt64Summary struct{ Min, Max, Average int64 }

// JobInt32Summary contains minimum, maximum and average samples, including DCGM blank values.
type JobInt32Summary struct{ Min, Max, Average int32 }

// JobFloat64Summary contains minimum, maximum and average samples, including DCGM blank values.
type JobFloat64Summary struct{ Min, Max, Average float64 }

// JobProcessStats identifies a process and its SM and memory utilization percentages.
type JobProcessStats struct {
	PID                              uint
	SMUtilization, MemoryUtilization float64
}

// JobHealthIncident records a health result for a watched system.
type JobHealthIncident struct {
	System HealthSystem
	Health HealthResult
}

// GPUJobStats contains job-lifetime measurements for one GPU, or the job summary.
// Values retain DCGM blank sentinels: use IsInt64Blank and IsInt32Blank for integer fields.
// Floating-point values at least 140737488355328.0 are DCGM blank sentinels.
// Summary min/max values for PCIe bandwidth, utilization and clocks are not applicable.
type GPUJobStats struct {
	GPU uint

	EnergyConsumed int64             // Milliwatt seconds.
	PowerUsage     JobFloat64Summary // Watts.

	PCIeRxBandwidth, PCIeTxBandwidth JobInt64Summary // KB/s.
	PCIeReplays                      int64

	// StartTime and EndTime are Unix microseconds. For a running job,
	// EndTime is the current collection endpoint; after JobStop it is the stop time.
	StartTime, EndTime               int64
	SMUtilization, MemoryUtilization JobInt32Summary // Percent.

	ECCSingleBit, ECCDoubleBit uint // ECCSingleBit is deprecated by DCGM.

	MemoryClock, SMClock JobInt32Summary // MHz.

	XIDTimestamps []int64 // Unix microseconds.

	ComputeProcesses, GraphicsProcesses []JobProcessStats

	MaxGPUMemoryUsed int64 // Bytes.

	PowerViolationTime, ThermalViolationTime, ReliabilityViolationTime int64 // Nanoseconds.

	BoardLimitViolationTime, LowUtilizationTime, SyncBoostTime int64 // Nanoseconds.

	OverallHealth   HealthResult
	HealthIncidents []JobHealthIncident
}

// JobInfo contains the summary and per-GPU statistics for a job.
type JobInfo struct {
	Summary GPUJobStats
	GPUs    []GPUJobStats
}

// WatchJobFields starts collection for jobs on group. The caller owns the group.
// updateFreq must be a positive whole number of microseconds. maxKeepAge and
// maxKeepSamples must be nonnegative; maxKeepAge is passed as fractional seconds.
// A zero retention limit leaves that limit to DCGM.
// DCGM may enable GPU accounting mode; callers that need to preserve its prior
// state must restore it after stopping collection.
func WatchJobFields(group GroupHandle, updateFreq, maxKeepAge time.Duration, maxKeepSamples int) error {
	return watchJobFieldsWithOps(cgoAdapter{}, group, updateFreq, maxKeepAge, maxKeepSamples)
}

func watchJobFieldsWithOps(api jobStatsCalls, group GroupHandle, updateFreq, maxKeepAge time.Duration, maxKeepSamples int) error {
	if updateFreq < time.Microsecond || updateFreq%time.Microsecond != 0 {
		return errors.New("watch job fields: frequency must be a positive whole number of microseconds")
	}
	if maxKeepAge < 0 || maxKeepSamples < 0 || int64(maxKeepSamples) > math.MaxInt32 {
		return errors.New("watch job fields: retention limits must fit nonnegative DCGM values")
	}
	return api.watchJobFieldsNative(group, updateFreq, maxKeepAge, maxKeepSamples)
}

// JobStart starts statistics for a unique, nonempty job ID of at most 63 bytes.
// Call WatchJobFields before starting a job. Job IDs cannot contain NUL bytes.
func JobStart(group GroupHandle, jobID string) error {
	return jobStartWithOps(cgoAdapter{}, group, jobID)
}

func jobStartWithOps(api jobStatsCalls, group GroupHandle, jobID string) error {
	if err := validateJobID(jobID); err != nil {
		return err
	}
	return api.jobStartNative(group, jobID)
}

// JobStop stops recording a job without removing its statistics.
func JobStop(jobID string) error { return jobStopWithOps(cgoAdapter{}, jobID) }

func jobStopWithOps(api jobStatsCalls, jobID string) error {
	if err := validateJobID(jobID); err != nil {
		return err
	}
	return api.jobStopNative(jobID)
}

// JobGetStats reads a running or stopped job's statistics into Go-owned values.
func JobGetStats(jobID string) (JobInfo, error) { return jobGetStatsWithOps(cgoAdapter{}, jobID) }

func jobGetStatsWithOps(api jobStatsCalls, jobID string) (JobInfo, error) {
	if err := validateJobID(jobID); err != nil {
		return JobInfo{}, err
	}
	return api.jobGetStatsNative(jobID)
}

// JobRemove removes one job's statistics, allowing its ID to be reused.
func JobRemove(jobID string) error { return jobRemoveWithOps(cgoAdapter{}, jobID) }

func jobRemoveWithOps(api jobStatsCalls, jobID string) error {
	if err := validateJobID(jobID); err != nil {
		return err
	}
	return api.jobRemoveNative(jobID)
}

// JobRemoveAll removes every job recorded by the connected hostengine, including other clients' jobs.
func JobRemoveAll() error { return cgoAdapter{}.jobRemoveAllNative() }

func validateJobID(jobID string) error {
	if jobID == "" || len(jobID) >= 64 || strings.ContainsRune(jobID, 0) {
		return errors.New("job ID must contain 1 to 63 bytes and no NUL")
	}
	return nil
}

func (cgoAdapter) watchJobFieldsNative(group GroupHandle, updateFreq, maxKeepAge time.Duration, maxKeepSamples int) error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireJobStatsSymbol("dcgmWatchJobFields"); err != nil {
		return err
	}
	return jobCallError("watch job fields", C.dcgmWatchJobFields(handle.handle, group.handle, C.longlong(updateFreq.Microseconds()), C.double(maxKeepAge.Seconds()), C.int(maxKeepSamples)))
}

func (cgoAdapter) jobStartNative(group GroupHandle, jobID string) error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireJobStatsSymbol("dcgmJobStartStats"); err != nil {
		return err
	}
	id := C.CString(jobID)
	defer C.free(unsafe.Pointer(id))
	return jobCallError("start job", C.dcgmJobStartStats(handle.handle, group.handle, id))
}

func (cgoAdapter) jobStopNative(jobID string) error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireJobStatsSymbol("dcgmJobStopStats"); err != nil {
		return err
	}
	id := C.CString(jobID)
	defer C.free(unsafe.Pointer(id))
	return jobCallError("stop job", C.dcgmJobStopStats(handle.handle, id))
}

func (cgoAdapter) jobRemoveNative(jobID string) error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireJobStatsSymbol("dcgmJobRemove"); err != nil {
		return err
	}
	id := C.CString(jobID)
	defer C.free(unsafe.Pointer(id))
	return jobCallError("remove job", C.dcgmJobRemove(handle.handle, id))
}

func (cgoAdapter) jobRemoveAllNative() error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireJobStatsSymbol("dcgmJobRemoveAll"); err != nil {
		return err
	}
	return jobCallError("remove all jobs", C.dcgmJobRemoveAll(handle.handle))
}

func (cgoAdapter) jobGetStatsNative(jobID string) (JobInfo, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireJobStatsSymbol("dcgmJobGetStats"); err != nil {
		return JobInfo{}, err
	}
	id := C.CString(jobID)
	defer C.free(unsafe.Pointer(id))
	var response C.dcgmJobInfo_t
	response.version = C.dcgmJobInfo_version
	if err := jobCallError("get job statistics", C.dcgmJobGetStats(handle.handle, id, &response)); err != nil {
		return JobInfo{}, err
	}
	return jobInfoFromNative(&response)
}

func jobCallError(operation string, result C.dcgmReturn_t) error {
	if result != C.DCGM_ST_OK {
		return fmt.Errorf("%s: %w", operation, &Error{msg: dcgmErrorText(result), Code: result})
	}
	return nil
}

func jobInfoFromNative(response *C.dcgmJobInfo_t) (JobInfo, error) {
	if err := jobCount("GPUs", int(response.numGpus), len(response.gpus)); err != nil {
		return JobInfo{}, err
	}
	summary, err := gpuJobStatsFromNative(&response.summary)
	if err != nil {
		return JobInfo{}, fmt.Errorf("job summary: %w", err)
	}
	info := JobInfo{Summary: summary, GPUs: make([]GPUJobStats, int(response.numGpus))}
	for i := range info.GPUs {
		info.GPUs[i], err = gpuJobStatsFromNative(&response.gpus[i])
		if err != nil {
			return JobInfo{}, fmt.Errorf("job GPU index %d: %w", i, err)
		}
	}
	return info, nil
}

func jobCount(name string, count, limit int) error {
	if count < 0 || count > limit {
		return fmt.Errorf("job statistics contain %d %s, expected 0 to %d", count, name, limit)
	}
	return nil
}

func gpuJobStatsFromNative(g *C.dcgmGpuUsageInfo_t) (GPUJobStats, error) {
	for _, c := range []struct {
		name         string
		count, limit int
	}{
		{"XID errors", int(g.numXidCriticalErrors), len(g.xidCriticalErrorsTs)},
		{"compute processes", int(g.numComputePids), len(g.computePidInfo)},
		{"graphics processes", int(g.numGraphicsPids), len(g.graphicsPidInfo)},
		{"health incidents", int(g.incidentCount), len(g.systems)},
	} {
		if err := jobCount(c.name, c.count, c.limit); err != nil {
			return GPUJobStats{}, err
		}
	}
	result := GPUJobStats{
		GPU:                      uint(g.gpuId),
		EnergyConsumed:           int64(g.energyConsumed),
		PowerUsage:               JobFloat64Summary{float64(g.powerUsage.minValue), float64(g.powerUsage.maxValue), float64(g.powerUsage.average)},
		PCIeReplays:              int64(g.pcieReplays),
		StartTime:                int64(g.startTime),
		EndTime:                  int64(g.endTime),
		ECCSingleBit:             uint(g.eccSingleBit),
		ECCDoubleBit:             uint(g.eccDoubleBit),
		MaxGPUMemoryUsed:         int64(g.maxGpuMemoryUsed),
		OverallHealth:            HealthResult(g.overallHealth),
		PCIeRxBandwidth:          JobInt64Summary{int64(g.pcieRxBandwidth.minValue), int64(g.pcieRxBandwidth.maxValue), int64(g.pcieRxBandwidth.average)},
		PCIeTxBandwidth:          JobInt64Summary{int64(g.pcieTxBandwidth.minValue), int64(g.pcieTxBandwidth.maxValue), int64(g.pcieTxBandwidth.average)},
		SMUtilization:            JobInt32Summary{int32(g.smUtilization.minValue), int32(g.smUtilization.maxValue), int32(g.smUtilization.average)},
		MemoryUtilization:        JobInt32Summary{int32(g.memoryUtilization.minValue), int32(g.memoryUtilization.maxValue), int32(g.memoryUtilization.average)},
		MemoryClock:              JobInt32Summary{int32(g.memoryClock.minValue), int32(g.memoryClock.maxValue), int32(g.memoryClock.average)},
		SMClock:                  JobInt32Summary{int32(g.smClock.minValue), int32(g.smClock.maxValue), int32(g.smClock.average)},
		PowerViolationTime:       int64(g.powerViolationTime),
		ThermalViolationTime:     int64(g.thermalViolationTime),
		ReliabilityViolationTime: int64(g.reliabilityViolationTime),
		BoardLimitViolationTime:  int64(g.boardLimitViolationTime),
		LowUtilizationTime:       int64(g.lowUtilizationTime),
		SyncBoostTime:            int64(g.syncBoostTime),
		XIDTimestamps:            make([]int64, int(g.numXidCriticalErrors)),
		ComputeProcesses:         make([]JobProcessStats, int(g.numComputePids)),
		GraphicsProcesses:        make([]JobProcessStats, int(g.numGraphicsPids)),
		HealthIncidents:          make([]JobHealthIncident, int(g.incidentCount)),
	}
	for i := range result.XIDTimestamps {
		result.XIDTimestamps[i] = int64(g.xidCriticalErrorsTs[i])
	}
	for i := range result.ComputeProcesses {
		p := g.computePidInfo[i]
		result.ComputeProcesses[i] = JobProcessStats{uint(p.pid), float64(p.smUtil), float64(p.memUtil)}
	}
	for i := range result.GraphicsProcesses {
		p := g.graphicsPidInfo[i]
		result.GraphicsProcesses[i] = JobProcessStats{uint(p.pid), float64(p.smUtil), float64(p.memUtil)}
	}
	for i := range result.HealthIncidents {
		s := g.systems[i]
		result.HealthIncidents[i] = JobHealthIncident{HealthSystem(s.system), HealthResult(s.health)}
	}
	return result, nil
}

func requireJobStatsSymbol(symbol string) error {
	if dcgmInitCounter > 0 && dcgmSymbolAvailable(symbol) {
		return nil
	}
	return &FeatureUnavailableError{Feature: "DCGM job statistics API", Symbol: symbol}
}
