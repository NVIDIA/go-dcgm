/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
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

const (
	nativeNVLinkStateDisabled = C.DcgmNvLinkLinkStateDisabled
	nativeNVLinkStateDown     = C.DcgmNvLinkLinkStateDown
	nativeNVLinkStateUp       = C.DcgmNvLinkLinkStateUp
)

// deviceTopologyResponseForTest builds a native fixture with the reported peer count.
// Tests can supply an oversized count without calling DCGM.
func deviceTopologyResponseForTest(count uint) *C.dcgmDeviceTopology_v2 {
	return &C.dcgmDeviceTopology_v2{numGpus: C.uint(count)}
}

// nvLinkStatusResponseForTest builds a native fixture with independent GPU and NVSwitch counts.
// Tests can vary each count to check its array bounds.
func nvLinkStatusResponseForTest(gpuCount, switchCount uint) *C.dcgmNvLinkStatus_v5 {
	return &C.dcgmNvLinkStatus_v5{
		numGpus:       C.uint(gpuCount),
		numNvSwitches: C.uint(switchCount),
	}
}
