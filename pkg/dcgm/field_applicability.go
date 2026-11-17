/*
 * Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 */

package dcgm

import "fmt"

// FieldSupportsTarget reports structural native applicability using the complete
// metadata returned by FieldGetByID and a native entity key. For FE_LINK, parent
// must identify the GPU or switch encoded in that key. Other parents are optional.
// Invalid metadata/identities and unclassified link fields return an error;
// a known incompatible target returns false without an error.
//
// Fully described non-link fields need not exist in the generated catalog.
// Native MIG aliases are accepted, but this does not select reporting scopes,
// change source identities, or guarantee hardware, permission or module support.
// Native watch/read errors and sample statuses still apply.
func FieldSupportsTarget(meta FieldMeta, entity, parent GroupEntityPair) (bool, error) {
	if err := validateNativeTarget(entity, parent); err != nil {
		return false, err
	}

	// dcgm_fields.h defines GLOBAL=0 and ENTITY=1. A global field has no
	// entity level; entity metadata must name a real entity group.
	if meta.FieldID == DCGM_FI_UNKNOWN || !validEntityGroup(meta.EntityLevel) ||
		(meta.Scope != 0 && meta.Scope != 1) ||
		(meta.Scope == 0) != (meta.EntityLevel == FE_NONE) {
		return false, fmt.Errorf("invalid metadata for field %d (scope %d, entity level %d): %w", meta.FieldID, meta.Scope, meta.EntityLevel, newBadParameterError())
	}

	if entity.EntityGroupId == FE_LINK && meta.EntityLevel != FE_NONE &&
		linkFieldParent(meta.FieldID, meta.EntityLevel) == FE_NONE &&
		(meta.EntityLevel == FE_LINK || !isKnownFieldID(meta.FieldID)) {
		return false, fmt.Errorf("field %d has no source-verified link classification", meta.FieldID)
	}

	return resolveFieldScope(meta, entity.EntityGroupId, parent.EntityGroupId, true), nil
}

// nativeFieldAlias records secondary targets verified in DCGM 4.7 and 4.8.
// SetPracticalEntityInfo maps GPU fields on GI/CI to the GPU and GI fields
// on CI to the GI. GPU/GI/CI fields can also be requested on their ancestors;
// XID samples retain three independent native histories.
func nativeFieldAlias(meta FieldMeta, entityGroup Field_Entity_Group) bool {
	if isGPUEntityGroup(entityGroup) {
		return nativeGPUFieldAlias(meta)
	}

	switch {
	// Sysmon CalculateUtilizationForEntity supports both CPUs and cores;
	// ReadCoreSpeed explicitly does not support a CPU aggregate.
	case entityGroup == FE_CPU && meta.EntityLevel == FE_CPU_CORE:
		switch meta.FieldID {
		case DCGM_FI_DEV_CPU_UTIL_RATIO, DCGM_FI_DEV_CPU_UTIL_USER_RATIO,
			DCGM_FI_DEV_CPU_UTIL_NICE_RATIO, DCGM_FI_DEV_CPU_UTIL_SYS_RATIO,
			DCGM_FI_DEV_CPU_UTIL_IRQ_RATIO:
			return true
		}

	// Sysmon resolves socket-owned temperature and power through
	// GetSocketIdForEntity, and emits vendor/model values for either target.
	case entityGroup == FE_CPU_CORE && meta.EntityLevel == FE_CPU:
		switch meta.FieldID {
		case DCGM_FI_DEV_CPU_TEMP_CELSIUS, DCGM_FI_DEV_CPU_TEMP_WARNING_CELSIUS,
			DCGM_FI_DEV_CPU_TEMP_CRITICAL_CELSIUS, DCGM_FI_DEV_CPU_POWER_WATTS,
			DCGM_FI_DEV_SYSIO_POWER_WATTS, DCGM_FI_DEV_MODULE_POWER_WATTS,
			DCGM_FI_DEV_CPU_POWER_LIMIT_WATTS, DCGM_FI_DEV_CPU_VENDOR, DCGM_FI_DEV_CPU_MODEL:
			return true
		}

	// NVSDM's watched-field update path aggregates these port fields when
	// requested on a switch. Other backends may still report no support.
	case entityGroup == FE_SWITCH && meta.EntityLevel == FE_LINK:
		switch meta.FieldID {
		case DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX, DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX,
			DCGM_FI_DEV_NVSWITCH_LINK_TX_BYTES_TOTAL, DCGM_FI_DEV_NVSWITCH_LINK_RX_BYTES_TOTAL,
			DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERROR_TOTAL, DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERROR_TOTAL,
			DCGM_FI_DEV_NVSWITCH_UUID:
			return true
		}

	// NVSDM maps these switch-declared counters to the same port telemetry
	// for switch-owned links. The caller must check the link parent.
	case entityGroup == FE_LINK && meta.EntityLevel == FE_SWITCH:
		switch meta.FieldID {
		case DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX, DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX,
			DCGM_FI_DEV_NVSWITCH_TX_BYTES_TOTAL, DCGM_FI_DEV_NVSWITCH_RX_BYTES_TOTAL:
			return true
		}
	}

	return false
}

// nativeGPUFieldAlias identifies fields the GPU updater accepts for native
// GPU/GI/CI targets. The vGPU-declared exceptions retain the requested cache key.
func nativeGPUFieldAlias(meta FieldMeta) bool {
	if isGPUEntityGroup(meta.EntityLevel) {
		return true
	}
	if meta.EntityLevel != FE_VGPU {
		return false
	}

	// DCGM 4.7/4.8's GPU updater reads these fields using the physical
	// GPU handle, including for valid GI/CI keys. Others need a vGPU instance.
	switch meta.FieldID {
	case DCGM_FI_DEV_VGPU_SUPPORTED_INFO, DCGM_FI_DEV_VGPU_CREATABLE_IDS,
		DCGM_FI_DEV_VGPU_INSTANCE_INFO, DCGM_FI_DEV_VGPU_UTIL_INFO,
		DCGM_FI_DEV_VGPU_PROCESS_UTIL_INFO, DCGM_FI_DEV_VGPU_SUPPORTED_IDS,
		DCGM_FI_DEV_VGPU_TYPE_INFO, DCGM_FI_DEV_VGPU_TYPE_NAME,
		DCGM_FI_DEV_VGPU_TYPE_CLASS, DCGM_FI_DEV_VGPU_TYPE_LICENSE,
		DCGM_FI_DEV_VGPU_LICENSE_STATUS:
		return true
	default:
		return false
	}
}

// isGPUEntityGroup identifies native GPU, GPU-instance and compute-instance keys.
func isGPUEntityGroup(group Field_Entity_Group) bool {
	return group == FE_GPU || group == FE_GPU_I || group == FE_GPU_CI
}
