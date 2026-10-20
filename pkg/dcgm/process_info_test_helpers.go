/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dcgm

/*
#include "dcgm_structs.h"
*/
import "C"

// processGPUInputForTest describes one GPU's process metrics in plain Go values
// so unit tests do not need to populate the native response fields directly.
type processGPUInputForTest struct {
	gpuID                    uint
	energyConsumed           int64
	startTime                int64
	endTime                  int64
	processSMUtilization     float64
	processMemoryUtilization float64
	pcieRxBandwidth          int64
	pcieTxBandwidth          int64
	pcieReplays              int64
	maxGPUMemoryUsed         int64
	eccSingleBit             uint
	eccDoubleBit             uint
	smUtilization            int32
	memoryUtilization        int32
	memoryClock              int32
	smClock                  int32
	powerViolationTime       int64
	thermalViolationTime     int64
	reliabilityViolationTime int64
	boardLimitViolationTime  int64
	lowUtilizationTime       int64
	syncBoostTime            int64
	reportedXIDCount         int
	xidTimestamps            []int64
}

// processInfoFromValuesForTest builds a native process response from Go test
// values, then passes it through the production response converter.
func processInfoFromValuesForTest(
	pid uint,
	name string,
	reportedGPUCount int,
	gpus []processGPUInputForTest,
) ([]ProcessInfo, error) {
	var response C.dcgmPidInfo_t
	response.pid = C.uint(pid)
	response.numGpus = C.int(reportedGPUCount)

	gpuCount := min(len(gpus), len(response.gpus))
	for i := 0; i < gpuCount; i++ {
		input := gpus[i]
		gpu := &response.gpus[i]
		gpu.gpuId = C.uint(input.gpuID)
		gpu.energyConsumed = C.longlong(input.energyConsumed)
		gpu.startTime = C.longlong(input.startTime)
		gpu.endTime = C.longlong(input.endTime)
		gpu.processUtilization.smUtil = C.double(input.processSMUtilization)
		gpu.processUtilization.memUtil = C.double(input.processMemoryUtilization)
		gpu.pcieRxBandwidth.average = C.longlong(input.pcieRxBandwidth)
		gpu.pcieTxBandwidth.average = C.longlong(input.pcieTxBandwidth)
		gpu.pcieReplays = C.longlong(input.pcieReplays)
		gpu.maxGpuMemoryUsed = C.longlong(input.maxGPUMemoryUsed)
		gpu.eccSingleBit = C.uint(input.eccSingleBit)
		gpu.eccDoubleBit = C.uint(input.eccDoubleBit)
		gpu.smUtilization.average = C.int(input.smUtilization)
		gpu.memoryUtilization.average = C.int(input.memoryUtilization)
		gpu.memoryClock.average = C.int(input.memoryClock)
		gpu.smClock.average = C.int(input.smClock)
		gpu.powerViolationTime = C.longlong(input.powerViolationTime)
		gpu.thermalViolationTime = C.longlong(input.thermalViolationTime)
		gpu.reliabilityViolationTime = C.longlong(input.reliabilityViolationTime)
		gpu.boardLimitViolationTime = C.longlong(input.boardLimitViolationTime)
		gpu.lowUtilizationTime = C.longlong(input.lowUtilizationTime)
		gpu.syncBoostTime = C.longlong(input.syncBoostTime)
		gpu.numXidCriticalErrors = C.int(input.reportedXIDCount)

		xidCount := min(len(input.xidTimestamps), len(gpu.xidCriticalErrorsTs))
		for j := 0; j < xidCount; j++ {
			gpu.xidCriticalErrorsTs[j] = C.longlong(input.xidTimestamps[j])
		}
	}

	return processInfoFromResponse(&response, name)
}

// processInfoLimitsForTest reports the fixed native array sizes used to test
// valid boundary counts and rejected overflow counts.
func processInfoLimitsForTest() (gpuCount, xidCount int) {
	var response C.dcgmPidInfo_t
	return len(response.gpus), len(response.gpus[0].xidCriticalErrorsTs)
}
