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

// FieldSupportsEntityGroup reports whether a known field can be read for an
// entity group. It preserves DCGM's declared scope and global fields. DCGM
// 4.7 also permits the listed NVLink fields to be read for FE_LINK entities.
//
// Callers must pass the field's declared entity group from FieldMeta. This
// function does not replace that metadata or make other GPU fields available
// for link reads.
func FieldSupportsEntityGroup(fieldID Short, declaredLevel, entityGroup Field_Entity_Group) bool {
	if !isKnownFieldID(fieldID) {
		return false
	}

	if declaredLevel == entityGroup || declaredLevel == FE_NONE {
		return true
	}

	if declaredLevel != FE_GPU || entityGroup != FE_LINK {
		return false
	}

	return isNVLinkLinkField(fieldID)
}

func isKnownFieldID(fieldID Short) bool {
	_, ok := knownFieldIDs[fieldID]
	return ok
}

func isNVLinkLinkField(fieldID Short) bool {
	return fieldID == DCGM_FI_PROF_NVLINK_TX_BYTES ||
		fieldID == DCGM_FI_PROF_NVLINK_RX_BYTES ||
		(fieldID >= DCGM_FI_DEV_NVLINK_TX_PACKET_TOTAL &&
			fieldID <= DCGM_FI_DEV_NVLINK_EFFECTIVE_ERROR_TOTAL) ||
		(fieldID >= DCGM_FI_DEV_NVLINK_FEC_HISTORY_0_TOTAL &&
			fieldID <= DCGM_FI_DEV_NVLINK_FEC_HISTORY_15_TOTAL)
}

var knownFieldIDs = func() map[Short]struct{} {
	fieldIDs := make(map[Short]struct{}, len(dcgmFields))
	for _, fieldID := range dcgmFields {
		fieldIDs[fieldID] = struct{}{}
	}
	return fieldIDs
}()
