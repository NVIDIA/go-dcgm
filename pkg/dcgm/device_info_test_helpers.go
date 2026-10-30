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

// supportedDevicesFromValuesForTest builds a native supported-device response
// and passes it through the production response converter.
func supportedDevicesFromValuesForTest(deviceIDs []uint, reportedCount int) ([]uint, error) {
	var nativeIDs [C.DCGM_MAX_NUM_DEVICES]C.uint
	for i := 0; i < min(len(deviceIDs), len(nativeIDs)); i++ {
		nativeIDs[i] = C.uint(deviceIDs[i])
	}
	return supportedDevicesFromNative(&nativeIDs, C.int(reportedCount))
}

// supportedDeviceCapacityForTest reports the fixed native supported-device capacity.
func supportedDeviceCapacityForTest() int {
	var nativeIDs [C.DCGM_MAX_NUM_DEVICES]C.uint
	return len(nativeIDs)
}

// entityGroupEntitiesFromValuesForTest builds a native entity response and
// passes it through the production response converter.
func entityGroupEntitiesFromValuesForTest(entityIDs []uint, reportedCount int) ([]uint, error) {
	var nativeIDs [C.DCGM_ENTITY_GROUP_ENTITIES_MAX]C.uint
	for i := 0; i < min(len(entityIDs), len(nativeIDs)); i++ {
		nativeIDs[i] = C.uint(entityIDs[i])
	}
	return entityGroupEntitiesFromNative(&nativeIDs, C.int(reportedCount))
}

// entityGroupCapacityForTest reports the fixed native entity-group capacity.
func entityGroupCapacityForTest() int {
	var nativeIDs [C.DCGM_ENTITY_GROUP_ENTITIES_MAX]C.uint
	return len(nativeIDs)
}

// nvLinkP2PStatusFromValuesForTest builds a native link-status matrix and
// passes it through the production response converter.
func nvLinkP2PStatusFromValuesForTest(statuses [][]Link_State, reportedCount uint) (NvLinkP2PStatus, error) {
	var response C.dcgmNvLinkP2PStatus_v1
	response.numGpus = C.uint(reportedCount)
	for i := 0; i < min(len(statuses), len(response.gpus)); i++ {
		for j := 0; j < min(len(statuses[i]), len(response.gpus[i].linkStatus)); j++ {
			response.gpus[i].linkStatus[j] = C.dcgmNvLinkGpuP2PStatus_t(statuses[i][j])
		}
	}
	return nvLinkP2PStatusFromNative(&response)
}

// nvLinkGPUCapacityForTest reports both dimensions of the fixed native link-status matrix.
func nvLinkGPUCapacityForTest() int {
	var response C.dcgmNvLinkP2PStatus_v1
	return min(len(response.gpus), len(response.gpus[0].linkStatus))
}

// deviceAttributesFromValuesForTest builds native GPU attributes and passes
// them through the production response converter.
func deviceAttributesFromValuesForTest(input deviceAttributes) deviceAttributes {
	var response C.dcgmDeviceAttributes_t
	copyCStringForTest(response.identifiers.pciBusId[:], input.busID)
	copyCStringForTest(response.identifiers.uuid[:], input.uuid)
	response.powerLimits.defaultPowerLimit = C.uint(input.power)
	response.memoryUsage.bar1Total = C.uint(input.bar1)
	response.memoryUsage.fbTotal = C.uint(input.fbTotal)
	copyCStringForTest(response.identifiers.brandName[:], input.identifiers.Brand)
	copyCStringForTest(response.identifiers.deviceName[:], input.identifiers.Model)
	copyCStringForTest(response.identifiers.serial[:], input.identifiers.Serial)
	copyCStringForTest(response.identifiers.vbios[:], input.identifiers.Vbios)
	copyCStringForTest(response.identifiers.inforomImageVersion[:], input.identifiers.InforomImageVersion)
	copyCStringForTest(response.identifiers.driverVersion[:], input.identifiers.DriverVersion)
	return deviceAttributesFromNative(&response)
}

// copyCStringForTest copies text into a zero-terminated native character array.
func copyCStringForTest(destination []C.char, value string) {
	for i := 0; i < min(len(value), len(destination)-1); i++ {
		destination[i] = C.char(value[i])
	}
}
