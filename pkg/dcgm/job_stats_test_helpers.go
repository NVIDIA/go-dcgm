package dcgm

/*
#include "dcgm_agent.h"
*/
import "C"

// jobResponseForTest populates both the summary and GPU entry to exercise copies and bounds.
func jobResponseForTest(gpuCount, xidCount, computeCount, graphicsCount, incidentCount int) (JobInfo, error) {
	var r C.dcgmJobInfo_t
	r.numGpus = C.int(gpuCount)
	g := &r.gpus[0]
	g.gpuId = 7
	g.energyConsumed = 12345
	g.powerUsage.minValue = 1.25
	g.powerUsage.maxValue = 9.75
	g.powerUsage.average = 4.5
	g.pcieRxBandwidth.minValue = 11
	g.pcieRxBandwidth.maxValue = 22
	g.pcieRxBandwidth.average = 17
	g.pcieTxBandwidth.minValue = 41
	g.pcieTxBandwidth.maxValue = 52
	g.pcieTxBandwidth.average = 47
	g.pcieReplays = 31
	g.startTime = 1000001
	g.endTime = 2000002
	g.smUtilization.minValue = 2
	g.smUtilization.maxValue = 95
	g.smUtilization.average = 60
	g.memoryUtilization.minValue = 3
	g.memoryUtilization.maxValue = 75
	g.memoryUtilization.average = 40
	g.memoryClock.minValue = 1000
	g.memoryClock.maxValue = 2000
	g.memoryClock.average = 1500
	g.smClock.minValue = 400
	g.smClock.maxValue = 1200
	g.smClock.average = 800
	g.eccSingleBit = 3
	g.eccDoubleBit = 4
	g.maxGpuMemoryUsed = 987654321
	g.powerViolationTime = 101
	g.thermalViolationTime = 102
	g.reliabilityViolationTime = 103
	g.boardLimitViolationTime = 104
	g.lowUtilizationTime = 105
	g.syncBoostTime = 106
	g.numXidCriticalErrors = C.int(xidCount)
	g.xidCriticalErrorsTs[0] = 1234567
	g.numComputePids = C.int(computeCount)
	g.computePidInfo[0].pid = 123
	g.computePidInfo[0].smUtil = 3.5
	g.computePidInfo[0].memUtil = 4.5
	g.numGraphicsPids = C.int(graphicsCount)
	g.graphicsPidInfo[0].pid = 456
	g.graphicsPidInfo[0].smUtil = 7.5
	g.graphicsPidInfo[0].memUtil = 8.5
	g.overallHealth = C.DCGM_HEALTH_RESULT_WARN
	g.incidentCount = C.uint(incidentCount)
	g.systems[0].system = C.DCGM_HEALTH_WATCH_PCIE
	g.systems[0].health = C.DCGM_HEALTH_RESULT_FAIL
	summary := &r.summary
	summary.gpuId = 107
	summary.energyConsumed = 12445
	summary.powerUsage.minValue = 101.25
	summary.powerUsage.maxValue = 109.75
	summary.powerUsage.average = 104.5
	summary.pcieRxBandwidth.minValue = 111
	summary.pcieRxBandwidth.maxValue = 122
	summary.pcieRxBandwidth.average = 117
	summary.pcieTxBandwidth.minValue = 141
	summary.pcieTxBandwidth.maxValue = 152
	summary.pcieTxBandwidth.average = 147
	summary.pcieReplays = 131
	summary.startTime = 1000101
	summary.endTime = 2000102
	summary.smUtilization.minValue = 102
	summary.smUtilization.maxValue = 195
	summary.smUtilization.average = 160
	summary.memoryUtilization.minValue = 103
	summary.memoryUtilization.maxValue = 175
	summary.memoryUtilization.average = 140
	summary.memoryClock.minValue = 1100
	summary.memoryClock.maxValue = 2100
	summary.memoryClock.average = 1600
	summary.smClock.minValue = 500
	summary.smClock.maxValue = 1300
	summary.smClock.average = 900
	summary.eccSingleBit = 103
	summary.eccDoubleBit = 104
	summary.maxGpuMemoryUsed = 987654421
	summary.powerViolationTime = 201
	summary.thermalViolationTime = 202
	summary.reliabilityViolationTime = 203
	summary.boardLimitViolationTime = 204
	summary.lowUtilizationTime = 205
	summary.syncBoostTime = 206
	summary.numXidCriticalErrors = C.int(xidCount)
	summary.xidCriticalErrorsTs[0] = 1234667
	summary.numComputePids = C.int(computeCount)
	summary.computePidInfo[0].pid = 223
	summary.computePidInfo[0].smUtil = 103.5
	summary.computePidInfo[0].memUtil = 104.5
	summary.numGraphicsPids = C.int(graphicsCount)
	summary.graphicsPidInfo[0].pid = 556
	summary.graphicsPidInfo[0].smUtil = 107.5
	summary.graphicsPidInfo[0].memUtil = 108.5
	summary.overallHealth = C.DCGM_HEALTH_RESULT_FAIL
	summary.incidentCount = C.uint(incidentCount)
	summary.systems[0].system = C.DCGM_HEALTH_WATCH_MEM
	summary.systems[0].health = C.DCGM_HEALTH_RESULT_PASS

	return jobInfoFromNative(&r)
}

func jobBlankResponseForTest() (JobInfo, error) {
	var r C.dcgmJobInfo_t
	r.summary.energyConsumed = C.DCGM_INT64_BLANK
	r.summary.smUtilization.average = C.DCGM_INT32_BLANK
	r.summary.powerUsage.average = C.DCGM_FP64_BLANK
	return jobInfoFromNative(&r)
}

func jobNativeErrorForTest() error { return &Error{msg: "bad parameter", Code: C.DCGM_ST_BADPARAM} }

func jobMaxResponseForTest() (JobInfo, error) {
	var r C.dcgmJobInfo_t
	r.numGpus = C.int(len(r.gpus))
	g := &r.summary
	g.numXidCriticalErrors = C.int(len(g.xidCriticalErrorsTs))
	g.numComputePids = C.int(len(g.computePidInfo))
	g.numGraphicsPids = C.int(len(g.graphicsPidInfo))
	g.incidentCount = C.uint(len(g.systems))
	for i := range r.gpus {
		r.gpus[i] = *g
	}
	return jobInfoFromNative(&r)
}

func jobInvalidGPUForTest() (JobInfo, error) {
	var r C.dcgmJobInfo_t
	r.numGpus = 1
	r.gpus[0].numComputePids = -1
	return jobInfoFromNative(&r)
}
