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
// also permits the listed NVLink fields to be read for FE_LINK entities.
//
// Callers must pass the field's declared entity group from FieldMeta. This
// function does not replace that metadata or make other GPU fields available
// for link reads. FE_LINK support is coarse: use FieldSupportsEntityGroupWithParent
// to distinguish GPU-owned links from switch-owned links.
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

// FieldSupportsEntityGroupWithParent reports structural field compatibility for
// an entity group and, for FE_LINK, its parent group. Pass declaredLevel from
// FieldMeta. Non-link entities retain FieldSupportsEntityGroup behavior and
// ignore parentGroup. Links require FE_GPU or FE_SWITCH parents, including for
// global fields. Unknown fields and unclassified link fields are rejected.
//
// FieldMeta supplies DCGM's declared scope, but not the valid link parents.
// The parent rules below fill that gap using DCGM 4.8.0's native readers.
// Compatibility does not guarantee device, driver, permission, or hostengine
// capability; native errors still apply.
func FieldSupportsEntityGroupWithParent(fieldID Short, declaredLevel, entityGroup, parentGroup Field_Entity_Group) bool {
	if entityGroup != FE_LINK {
		return FieldSupportsEntityGroup(fieldID, declaredLevel, entityGroup)
	}

	if (parentGroup != FE_GPU && parentGroup != FE_SWITCH) || !isKnownFieldID(fieldID) {
		return false
	}

	if declaredLevel == FE_NONE {
		return true
	}

	if declaredLevel == FE_GPU {
		return parentGroup == FE_GPU && isNVLinkLinkField(fieldID)
	}

	if declaredLevel != FE_LINK {
		return false
	}

	switch parentGroup {
	case FE_GPU:
		// DcgmCacheManager::ReadAndCacheNvLinkFields routes these fields to
		// GPU-only state, per-link NVML, GPM, and PRM readers. Their field lists
		// and dcgmlib/src/dcgm_fields.cpp are the source of this classification.
		switch fieldID {
		case DCGM_FI_DEV_NVLINK_LINK_STATUS,
			DCGM_FI_DEV_NVLINK_PPCNT_IBPC_TX_WAIT_TOTAL,
			DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_PER_LINK_TOTAL,
			DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_PER_LINK_TOTAL,
			DCGM_FI_DEV_NVLINK_REPLAY_ERROR_PER_LINK_TOTAL,
			DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_PER_LINK_TOTAL,
			DCGM_FI_DEV_NVLINK_THROUGHPUT_PER_LINK,
			DCGM_FI_DEV_NVLINK_TX_THROUGHPUT_PER_LINK,
			DCGM_FI_DEV_NVLINK_RX_THROUGHPUT_PER_LINK,
			DCGM_FI_PROF_NVLINK_TX_BYTES_PER_LINK,
			DCGM_FI_PROF_NVLINK_RX_BYTES_PER_LINK,
			DCGM_FI_DEV_NVLINK_PPRM_OPER_RECOVERY,
			DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_SINCE_LAST_SECONDS,
			DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_TIME_BETWEEN_LAST_TWO,
			DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_SUCCESSFUL_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PHYSICAL_RECOVERY_SUCCESSFUL_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PHYSICAL_LINK_DOWN_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PLR_RX_CODE_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PLR_RX_CODE_ERROR_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PLR_RX_CODE_UNCORRECTABLE_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PLR_TX_CODE_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PLR_TX_RETRY_CODE_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PLR_TX_RETRY_EVENT_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_PLR_SYNC_EVENT_TOTAL,
			DCGM_FI_DEV_NVLINK_PPCNT_RECOVERY_LAST_GAP_SECONDS:
			return true
		}
	case FE_SWITCH:
		// The NVSwitch module owns these declared-link fields. NSCQ
		// GetNvSwitchObject and NVSDM port lookup accept switch links only.
		switch fieldID {
		case DCGM_FI_DEV_NVSWITCH_LINK_ID,
			DCGM_FI_DEV_NVSWITCH_LINK_STATUS,
			DCGM_FI_DEV_NVSWITCH_LINK_TYPE,
			DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS,
			DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE,
			DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN,
			DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION,
			DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_LINK_ID,
			DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_LINK_SID,
			DCGM_FI_DEV_NVSWITCH_UUID,
			DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX,
			DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX,
			DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS,
			DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS,
			DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERROR_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERROR_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERROR_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC0_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC1_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC2_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_SAMPLE_VC3_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L0_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L1_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L2_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L3_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L4_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L5_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L6_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERROR_L7_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L0_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L1_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L2_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L3_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L4_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L5_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L6_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERROR_L7_TOTAL:
			return true
		}
	}

	return false
}

// isKnownFieldID rejects fields absent from the generated binding metadata.
func isKnownFieldID(fieldID Short) bool {
	_, ok := knownFieldIDs[fieldID]
	return ok
}

// isNVLinkLinkField matches the GPU-declared link fields accepted by
// DcgmFieldIsNvLinkCountField and ReadAndCacheNvLinkFields.
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
