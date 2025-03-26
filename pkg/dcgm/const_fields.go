/*
 * Copyright (c) 2025, NVIDIA CORPORATION.  All rights reserved.
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

const (
	// DCGM_FI_UNKNOWN represents a NULL field
	DCGM_FI_UNKNOWN Short = 0
	// DCGM_FI_DRIVER_VERSION represents the driver version string
	DCGM_FI_DRIVER_VERSION Short = 1
	// DCGM_FI_NVML_VERSION represents the underlying NVML version string
	DCGM_FI_NVML_VERSION Short = 2
	// DCGM_FI_PROCESS_NAME represents the process name
	DCGM_FI_PROCESS_NAME Short = 3
	// DCGM_FI_DEV_COUNT represents the number of devices on the node
	DCGM_FI_DEV_COUNT Short = 4
	// DCGM_FI_CUDA_DRIVER_VERSION represents the CUDA driver version. Returns a number with the major value in the thousands place and the minor value in the hundreds place (e.g. CUDA 11.1 = 11100)
	DCGM_FI_CUDA_DRIVER_VERSION Short = 5
	// DCGM_FI_DEV_NAME represents the name of the GPU device
	DCGM_FI_DEV_NAME Short = 50
	// DCGM_FI_DEV_BRAND represents the device brand
	DCGM_FI_DEV_BRAND Short = 51
	// DCGM_FI_DEV_NVML_INDEX represents the NVML index of this GPU
	DCGM_FI_DEV_NVML_INDEX Short = 52
	// DCGM_FI_DEV_SERIAL represents the device serial number
	DCGM_FI_DEV_SERIAL Short = 53
	// DCGM_FI_DEV_UUID represents the UUID corresponding to the device
	DCGM_FI_DEV_UUID Short = 54
	// DCGM_FI_DEV_MINOR_NUMBER represents the device node minor number (/dev/nvidia#)
	DCGM_FI_DEV_MINOR_NUMBER Short = 55
	// DCGM_FI_DEV_OEM_INFOROM_VER represents the OEM inforom version
	DCGM_FI_DEV_OEM_INFOROM_VER Short = 56
	// DCGM_FI_DEV_PCI_BUSID represents the PCI attributes for the device
	DCGM_FI_DEV_PCI_BUSID Short = 57
	// DCGM_FI_DEV_PCI_COMBINED_ID represents the combined 16-bit device id and 16-bit vendor id
	DCGM_FI_DEV_PCI_COMBINED_ID Short = 58
	// DCGM_FI_DEV_PCI_SUBSYS_ID represents the 32-bit Sub System Device ID
	DCGM_FI_DEV_PCI_SUBSYS_ID Short = 59
	// DCGM_FI_GPU_TOPOLOGY_PCI represents the topology of all GPUs on the system via PCI (static)
	DCGM_FI_GPU_TOPOLOGY_PCI Short = 60
	// DCGM_FI_GPU_TOPOLOGY_NVLINK represents the topology of all GPUs on the system via NVLINK (static)
	DCGM_FI_GPU_TOPOLOGY_NVLINK Short = 61
	// DCGM_FI_GPU_TOPOLOGY_AFFINITY represents the affinity of all GPUs on the system (static)
	DCGM_FI_GPU_TOPOLOGY_AFFINITY Short = 62
	// DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY represents the CUDA compute capability for the device. The major version is the upper 32 bits and the minor version is the lower 32 bits
	DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY Short = 63
	// DCGM_FI_DEV_COMPUTE_MODE represents the compute mode for the device
	DCGM_FI_DEV_COMPUTE_MODE Short = 65
	// DCGM_FI_DEV_PERSISTENCE_MODE represents the persistence mode for the device. Boolean: 0 is disabled, 1 is enabled
	DCGM_FI_DEV_PERSISTENCE_MODE Short = 66
	// DCGM_FI_DEV_MIG_MODE represents the MIG mode for the device. Boolean: 0 is disabled, 1 is enabled
	DCGM_FI_DEV_MIG_MODE Short = 67
	// DCGM_FI_DEV_CUDA_VISIBLE_DEVICES_STR represents the string that CUDA_VISIBLE_DEVICES should be set to for this entity (including MIG)
	DCGM_FI_DEV_CUDA_VISIBLE_DEVICES_STR Short = 68
	// DCGM_FI_DEV_MIG_MAX_SLICES represents the maximum number of MIG slices supported by this GPU
	DCGM_FI_DEV_MIG_MAX_SLICES Short = 69
	// DCGM_FI_DEV_CPU_AFFINITY_0 represents the device CPU affinity for CPUs 0-63
	DCGM_FI_DEV_CPU_AFFINITY_0 Short = 70
	// DCGM_FI_DEV_CPU_AFFINITY_1 represents the device CPU affinity for CPUs 64-127
	DCGM_FI_DEV_CPU_AFFINITY_1 Short = 71
	// DCGM_FI_DEV_CPU_AFFINITY_2 represents the device CPU affinity for CPUs 128-191
	DCGM_FI_DEV_CPU_AFFINITY_2 Short = 72
	// DCGM_FI_DEV_CPU_AFFINITY_3 represents the device CPU affinity for CPUs 192-255
	DCGM_FI_DEV_CPU_AFFINITY_3 Short = 73
	// DCGM_FI_DEV_CC_MODE represents the ConfidentialCompute/AmpereProtectedMemory status. 0 = disabled, 1 = enabled
	DCGM_FI_DEV_CC_MODE Short = 74
	// DCGM_FI_DEV_MIG_ATTRIBUTES represents the attributes for the given MIG device handles
	DCGM_FI_DEV_MIG_ATTRIBUTES Short = 75
	// DCGM_FI_DEV_MIG_GI_INFO represents the GPU instance profile information
	DCGM_FI_DEV_MIG_GI_INFO Short = 76
	// DCGM_FI_DEV_MIG_CI_INFO represents the compute instance profile information
	DCGM_FI_DEV_MIG_CI_INFO Short = 77
	// DCGM_FI_DEV_ECC_INFOROM_VER represents the ECC inforom version
	DCGM_FI_DEV_ECC_INFOROM_VER Short = 80
	// DCGM_FI_DEV_POWER_INFOROM_VER represents the power management object inforom version
	DCGM_FI_DEV_POWER_INFOROM_VER Short = 81
	// DCGM_FI_DEV_INFOROM_IMAGE_VER represents the inforom image version
	DCGM_FI_DEV_INFOROM_IMAGE_VER Short = 82
	// DCGM_FI_DEV_INFOROM_CONFIG_CHECK represents the inforom configuration checksum
	DCGM_FI_DEV_INFOROM_CONFIG_CHECK Short = 83
	// DCGM_FI_DEV_INFOROM_CONFIG_VALID represents whether the inforom configuration is valid. Reads the infoROM from the flash and verifies the checksums
	DCGM_FI_DEV_INFOROM_CONFIG_VALID Short = 84
	// DCGM_FI_DEV_VBIOS_VERSION represents the VBIOS version of the device
	DCGM_FI_DEV_VBIOS_VERSION Short = 85
	// DCGM_FI_DEV_MEM_AFFINITY_0 represents the device memory node affinity for nodes 0-63
	DCGM_FI_DEV_MEM_AFFINITY_0 Short = 86
	// DCGM_FI_DEV_MEM_AFFINITY_1 represents the device memory node affinity for nodes 64-127
	DCGM_FI_DEV_MEM_AFFINITY_1 Short = 87
	// DCGM_FI_DEV_MEM_AFFINITY_2 represents the device memory node affinity for nodes 128-191
	DCGM_FI_DEV_MEM_AFFINITY_2 Short = 88
	// DCGM_FI_DEV_MEM_AFFINITY_3 represents the device memory node affinity for nodes 192-255
	DCGM_FI_DEV_MEM_AFFINITY_3 Short = 89
	// DCGM_FI_DEV_BAR1_TOTAL represents the total BAR1 memory of the GPU in MB
	DCGM_FI_DEV_BAR1_TOTAL Short = 90
	// DCGM_FI_SYNC_BOOST represents the sync boost settings on the node (Deprecated)
	DCGM_FI_SYNC_BOOST Short = 91
	// DCGM_FI_DEV_BAR1_USED represents the used BAR1 memory of the GPU in MB
	DCGM_FI_DEV_BAR1_USED Short = 92
	// DCGM_FI_DEV_BAR1_FREE represents the free BAR1 memory of the GPU in MB
	DCGM_FI_DEV_BAR1_FREE Short = 93
	// DCGM_FI_DEV_GPM_SUPPORT represents the GPM support for the device
	DCGM_FI_DEV_GPM_SUPPORT Short = 94
	// DCGM_FI_DEV_SM_CLOCK represents the SM clock for the device
	DCGM_FI_DEV_SM_CLOCK Short = 100
	// DCGM_FI_DEV_MEM_CLOCK represents the memory clock for the device
	DCGM_FI_DEV_MEM_CLOCK Short = 101
	// DCGM_FI_DEV_VIDEO_CLOCK represents the video encoder/decoder clock for the device
	DCGM_FI_DEV_VIDEO_CLOCK Short = 102
	// DCGM_FI_DEV_APP_SM_CLOCK represents the SM application clocks
	DCGM_FI_DEV_APP_SM_CLOCK Short = 110
	// DCGM_FI_DEV_APP_MEM_CLOCK represents the memory application clocks
	DCGM_FI_DEV_APP_MEM_CLOCK Short = 111
	// DCGM_FI_DEV_CLOCKS_EVENT_REASONS represents the current clock event reasons (bitmask of DCGM_CLOCKS_EVENT_REASON_*)
	DCGM_FI_DEV_CLOCKS_EVENT_REASONS Short = 112
	// DCGM_FI_DEV_CLOCK_THROTTLE_REASONS represents the current clock throttle reasons (Deprecated: Use DCGM_FI_DEV_CLOCKS_EVENT_REASONS instead)
	DCGM_FI_DEV_CLOCK_THROTTLE_REASONS Short = DCGM_FI_DEV_CLOCKS_EVENT_REASONS
	// DCGM_FI_DEV_MAX_SM_CLOCK represents the maximum supported SM clock for the device
	DCGM_FI_DEV_MAX_SM_CLOCK Short = 113
	// DCGM_FI_DEV_MAX_MEM_CLOCK represents the maximum supported memory clock for the device
	DCGM_FI_DEV_MAX_MEM_CLOCK Short = 114
	// DCGM_FI_DEV_MAX_VIDEO_CLOCK represents the maximum supported video encoder/decoder clock for the device
	DCGM_FI_DEV_MAX_VIDEO_CLOCK Short = 115
	// DCGM_FI_DEV_AUTOBOOST represents the auto-boost setting for the device (1 = enabled, 0 = disabled)
	DCGM_FI_DEV_AUTOBOOST Short = 120
	// DCGM_FI_DEV_SUPPORTED_CLOCKS represents the supported clocks for the device
	DCGM_FI_DEV_SUPPORTED_CLOCKS Short = 130
	// DCGM_FI_DEV_MEMORY_TEMP represents the memory temperature for the device
	DCGM_FI_DEV_MEMORY_TEMP Short = 140
	// DCGM_FI_DEV_GPU_TEMP represents the current temperature readings for the device, in degrees C
	DCGM_FI_DEV_GPU_TEMP Short = 150
	// DCGM_FI_DEV_MEM_MAX_OP_TEMP represents the maximum operating temperature for the memory of this GPU
	DCGM_FI_DEV_MEM_MAX_OP_TEMP Short = 151
	// DCGM_FI_DEV_GPU_MAX_OP_TEMP represents the maximum operating temperature for this GPU
	DCGM_FI_DEV_GPU_MAX_OP_TEMP Short = 152
	// DCGM_FI_DEV_GPU_TEMP_LIMIT represents the thermal margin temperature (distance to nearest slowdown threshold) for this GPU
	DCGM_FI_DEV_GPU_TEMP_LIMIT Short = 153
	// DCGM_FI_DEV_POWER_USAGE represents the power usage for the device in Watts
	DCGM_FI_DEV_POWER_USAGE Short = 155
	// DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION represents the total energy consumption for the GPU in mJ since the driver was last reloaded
	DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION Short = 156
	// DCGM_FI_DEV_POWER_USAGE_INSTANT represents the current instantaneous power usage of the device in Watts
	DCGM_FI_DEV_POWER_USAGE_INSTANT Short = 157
	// DCGM_FI_DEV_SLOWDOWN_TEMP represents the slowdown temperature for the device
	DCGM_FI_DEV_SLOWDOWN_TEMP Short = 158
	// DCGM_FI_DEV_SHUTDOWN_TEMP represents the shutdown temperature for the device
	DCGM_FI_DEV_SHUTDOWN_TEMP Short = 159
	// DCGM_FI_DEV_POWER_MGMT_LIMIT represents the current power limit for the device
	DCGM_FI_DEV_POWER_MGMT_LIMIT Short = 160
	// DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN represents the minimum power management limit for the device
	DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN Short = 161
	// DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX represents the maximum power management limit for the device
	DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX Short = 162
	// DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF represents the default power management limit for the device
	DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF Short = 163
	// DCGM_FI_DEV_ENFORCED_POWER_LIMIT represents the effective power limit that the driver enforces after taking into account all limiters
	DCGM_FI_DEV_ENFORCED_POWER_LIMIT Short = 164
	// DCGM_FI_DEV_REQUESTED_POWER_PROFILE_MASK represents the requested workload power profile mask (Blackwell and newer)
	DCGM_FI_DEV_REQUESTED_POWER_PROFILE_MASK Short = 165
	// DCGM_FI_DEV_ENFORCED_POWER_PROFILE_MASK represents the enforced workload power profile mask (Blackwell and newer)
	DCGM_FI_DEV_ENFORCED_POWER_PROFILE_MASK Short = 166
	// DCGM_FI_DEV_VALID_POWER_PROFILE_MASK represents the valid workload power profile mask (Blackwell and newer)
	DCGM_FI_DEV_VALID_POWER_PROFILE_MASK Short = 167
	// DCGM_FI_DEV_FABRIC_MANAGER_STATUS is the value for fabric manager status
	DCGM_FI_DEV_FABRIC_MANAGER_STATUS Short = 168
	// DCGM_FI_DEV_FABRIC_MANAGER_ERROR_CODE is the value for fabric manager error code
	DCGM_FI_DEV_FABRIC_MANAGER_ERROR_CODE Short = 171
	// DCGM_FI_DEV_FABRIC_CLUSTER_UUID is the value for fabric cluster UUID
	DCGM_FI_DEV_FABRIC_CLUSTER_UUID Short = 172
	// DCGM_FI_DEV_FABRIC_CLIQUE_ID is the value for fabric clique ID
	DCGM_FI_DEV_FABRIC_CLIQUE_ID Short = 173
	// DCGM_FI_DEV_PSTATE is the value for P-state
	DCGM_FI_DEV_PSTATE Short = 190
	// DCGM_FI_DEV_FAN_SPEED is the value for fan speed
	DCGM_FI_DEV_FAN_SPEED Short = 191
	// DCGM_FI_DEV_PCIE_TX_THROUGHPUT represents the PCIe transmit throughput in KB/s
	DCGM_FI_DEV_PCIE_TX_THROUGHPUT Short = 200
	// DCGM_FI_DEV_PCIE_RX_THROUGHPUT represents the PCIe receive throughput in KB/s
	DCGM_FI_DEV_PCIE_RX_THROUGHPUT Short = 201
	// DCGM_FI_DEV_PCIE_REPLAY_COUNTER represents the PCIe replay counter value
	DCGM_FI_DEV_PCIE_REPLAY_COUNTER Short = 202
	// DCGM_FI_DEV_GPU_UTIL represents the GPU utilization in percent
	DCGM_FI_DEV_GPU_UTIL Short = 203
	// DCGM_FI_DEV_MEM_COPY_UTIL represents the memory copy utilization in percent
	DCGM_FI_DEV_MEM_COPY_UTIL Short = 204
	// DCGM_FI_DEV_ACCOUNTING_DATA represents the process accounting information
	DCGM_FI_DEV_ACCOUNTING_DATA Short = 205
	// DCGM_FI_DEV_ENC_UTIL represents the encoder utilization in percent
	DCGM_FI_DEV_ENC_UTIL Short = 206
	// DCGM_FI_DEV_DEC_UTIL represents the decoder utilization in percent
	DCGM_FI_DEV_DEC_UTIL Short = 207
	// DCGM_FI_DEV_XID_ERRORS is the value for XID errors
	DCGM_FI_DEV_XID_ERRORS Short = 230
	// DCGM_FI_DEV_PCIE_MAX_LINK_GEN is the value for PCIe max link generation
	DCGM_FI_DEV_PCIE_MAX_LINK_GEN Short = 235
	// DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH is the value for PCIe max link width
	DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH Short = 236
	// DCGM_FI_DEV_PCIE_LINK_GEN is the value for PCIe link generation
	DCGM_FI_DEV_PCIE_LINK_GEN Short = 237
	// DCGM_FI_DEV_PCIE_LINK_WIDTH is the value for PCIe link width
	DCGM_FI_DEV_PCIE_LINK_WIDTH Short = 238
	// DCGM_FI_DEV_POWER_VIOLATION is the value for power violation time in microseconds
	DCGM_FI_DEV_POWER_VIOLATION Short = 240
	// DCGM_FI_DEV_THERMAL_VIOLATION is the value for thermal violation time in microseconds
	DCGM_FI_DEV_THERMAL_VIOLATION Short = 241
	// DCGM_FI_DEV_SYNC_BOOST_VIOLATION is the value for sync boost violation time in microseconds
	DCGM_FI_DEV_SYNC_BOOST_VIOLATION Short = 242
	// DCGM_FI_DEV_BOARD_LIMIT_VIOLATION is the value for board limit violation time in microseconds
	DCGM_FI_DEV_BOARD_LIMIT_VIOLATION Short = 243
	// DCGM_FI_DEV_LOW_UTIL_VIOLATION is the value for low utilization violation time in microseconds
	DCGM_FI_DEV_LOW_UTIL_VIOLATION Short = 244
	// DCGM_FI_DEV_RELIABILITY_VIOLATION is the value for reliability violation time in microseconds
	DCGM_FI_DEV_RELIABILITY_VIOLATION Short = 245
	// DCGM_FI_DEV_TOTAL_APP_CLOCKS_VIOLATION is the value for total application clocks violation time in microseconds
	DCGM_FI_DEV_TOTAL_APP_CLOCKS_VIOLATION Short = 246
	// DCGM_FI_DEV_TOTAL_BASE_CLOCKS_VIOLATION is the value for total base clocks violation time in microseconds
	DCGM_FI_DEV_TOTAL_BASE_CLOCKS_VIOLATION Short = 247
	// DCGM_FI_DEV_FB_TOTAL is the value for framebuffer total
	DCGM_FI_DEV_FB_TOTAL Short = 250
	// DCGM_FI_DEV_FB_FREE is the value for framebuffer free
	DCGM_FI_DEV_FB_FREE Short = 251
	// DCGM_FI_DEV_FB_USED is the value for framebuffer used
	DCGM_FI_DEV_FB_USED Short = 252
	// DCGM_FI_DEV_FB_RESERVED is the value for framebuffer reserved
	DCGM_FI_DEV_FB_RESERVED Short = 253
	// DCGM_FI_DEV_FB_USED_PERCENT is the value for framebuffer used percent
	DCGM_FI_DEV_FB_USED_PERCENT Short = 254
	// DCGM_FI_DEV_C2C_LINK_COUNT is the value for C2C link count
	DCGM_FI_DEV_C2C_LINK_COUNT Short = 285
	// DCGM_FI_DEV_C2C_LINK_STATUS is the value for C2C link status
	DCGM_FI_DEV_C2C_LINK_STATUS Short = 286
	// DCGM_FI_DEV_C2C_MAX_BANDWIDTH is the value for C2C max bandwidth
	DCGM_FI_DEV_C2C_MAX_BANDWIDTH Short = 287
	// DCGM_FI_DEV_ECC_CURRENT is the value for ECC current
	DCGM_FI_DEV_ECC_CURRENT Short = 220
	// DCGM_FI_DEV_ECC_PENDING is the value for ECC pending
	DCGM_FI_DEV_ECC_PENDING Short = 221
	// DCGM_FI_DEV_ECC_SBE_VOL_TOTAL represents the total number of single-bit ECC errors detected since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_TOTAL Short = 310
	// DCGM_FI_DEV_ECC_DBE_VOL_TOTAL represents the total number of double-bit ECC errors detected since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_TOTAL Short = 311
	// DCGM_FI_DEV_ECC_SBE_AGG_TOTAL represents the total number of single-bit ECC errors detected since the last counter reset (aggregate)
	DCGM_FI_DEV_ECC_SBE_AGG_TOTAL Short = 312
	// DCGM_FI_DEV_ECC_DBE_AGG_TOTAL represents the total number of double-bit ECC errors detected since the last counter reset (aggregate)
	DCGM_FI_DEV_ECC_DBE_AGG_TOTAL Short = 313
	// DCGM_FI_DEV_ECC_SBE_VOL_L1 represents the number of single-bit ECC errors detected in L1 cache since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_L1 Short = 314
	// DCGM_FI_DEV_ECC_DBE_VOL_L1 represents the number of double-bit ECC errors detected in L1 cache since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_L1 Short = 315
	// DCGM_FI_DEV_ECC_SBE_VOL_L2 represents the number of single-bit ECC errors detected in L2 cache since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_L2 Short = 316
	// DCGM_FI_DEV_ECC_DBE_VOL_L2 represents the number of double-bit ECC errors detected in L2 cache since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_L2 Short = 317
	// DCGM_FI_DEV_ECC_SBE_VOL_DEV represents the number of single-bit ECC errors detected in device memory since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_DEV Short = 318
	// DCGM_FI_DEV_ECC_DBE_VOL_DEV represents the number of double-bit ECC errors detected in device memory since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_DEV Short = 319
	// DCGM_FI_DEV_ECC_SBE_VOL_REG represents the number of single-bit ECC errors detected in register file since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_REG Short = 320
	// DCGM_FI_DEV_ECC_DBE_VOL_REG represents the number of double-bit ECC errors detected in register file since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_REG Short = 321
	// DCGM_FI_DEV_ECC_SBE_VOL_TEX represents the number of single-bit ECC errors detected in texture memory since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_TEX Short = 322
	// DCGM_FI_DEV_ECC_DBE_VOL_TEX represents the number of double-bit ECC errors detected in texture memory since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_TEX Short = 323
	// DCGM_FI_DEV_ECC_SBE_AGG_L1 represents the aggregate number of single-bit ECC errors detected in L1 cache
	DCGM_FI_DEV_ECC_SBE_AGG_L1 Short = 324
	// DCGM_FI_DEV_ECC_DBE_AGG_L1 represents the aggregate number of double-bit ECC errors detected in L1 cache
	DCGM_FI_DEV_ECC_DBE_AGG_L1 Short = 325
	// DCGM_FI_DEV_ECC_SBE_AGG_L2 represents the aggregate number of single-bit ECC errors detected in L2 cache
	DCGM_FI_DEV_ECC_SBE_AGG_L2 Short = 326
	// DCGM_FI_DEV_ECC_DBE_AGG_L2 represents the aggregate number of double-bit ECC errors detected in L2 cache
	DCGM_FI_DEV_ECC_DBE_AGG_L2 Short = 327
	// DCGM_FI_DEV_ECC_SBE_AGG_DEV represents the aggregate number of single-bit ECC errors detected in device memory
	DCGM_FI_DEV_ECC_SBE_AGG_DEV Short = 328
	// DCGM_FI_DEV_ECC_DBE_AGG_DEV represents the aggregate number of double-bit ECC errors detected in device memory
	DCGM_FI_DEV_ECC_DBE_AGG_DEV Short = 329
	// DCGM_FI_DEV_ECC_SBE_AGG_REG represents the aggregate number of single-bit ECC errors detected in register file
	DCGM_FI_DEV_ECC_SBE_AGG_REG Short = 330
	// DCGM_FI_DEV_ECC_DBE_AGG_REG represents the aggregate number of double-bit ECC errors detected in register file
	DCGM_FI_DEV_ECC_DBE_AGG_REG Short = 331
	// DCGM_FI_DEV_ECC_SBE_AGG_TEX represents the aggregate number of single-bit ECC errors detected in texture memory
	DCGM_FI_DEV_ECC_SBE_AGG_TEX Short = 332
	// DCGM_FI_DEV_ECC_DBE_AGG_TEX represents the aggregate number of double-bit ECC errors detected in texture memory
	DCGM_FI_DEV_ECC_DBE_AGG_TEX Short = 333
	// DCGM_FI_DEV_ECC_SBE_VOL_SHM represents the number of single-bit ECC errors detected in shared memory since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_SHM Short = 334
	// DCGM_FI_DEV_ECC_DBE_VOL_SHM represents the number of double-bit ECC errors detected in shared memory since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_SHM Short = 335
	// DCGM_FI_DEV_ECC_SBE_VOL_CBU represents the number of single-bit ECC errors detected in CBU since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_CBU Short = 336
	// DCGM_FI_DEV_ECC_DBE_VOL_CBU represents the number of double-bit ECC errors detected in CBU since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_CBU Short = 337
	// DCGM_FI_DEV_ECC_SBE_AGG_SHM represents the aggregate number of single-bit ECC errors detected in shared memory
	DCGM_FI_DEV_ECC_SBE_AGG_SHM Short = 338
	// DCGM_FI_DEV_ECC_DBE_AGG_SHM represents the aggregate number of double-bit ECC errors detected in shared memory
	DCGM_FI_DEV_ECC_DBE_AGG_SHM Short = 339
	// DCGM_FI_DEV_ECC_SBE_AGG_CBU represents the aggregate number of single-bit ECC errors detected in CBU
	DCGM_FI_DEV_ECC_SBE_AGG_CBU Short = 340
	// DCGM_FI_DEV_ECC_DBE_AGG_CBU represents the aggregate number of double-bit ECC errors detected in CBU
	DCGM_FI_DEV_ECC_DBE_AGG_CBU Short = 341
	// DCGM_FI_DEV_ECC_SBE_VOL_SRM represents the number of single-bit ECC errors detected in SRM since the last counter reset
	DCGM_FI_DEV_ECC_SBE_VOL_SRM Short = 342
	// DCGM_FI_DEV_ECC_DBE_VOL_SRM represents the number of double-bit ECC errors detected in SRM since the last counter reset
	DCGM_FI_DEV_ECC_DBE_VOL_SRM Short = 343
	// DCGM_FI_DEV_ECC_SBE_AGG_SRM represents the aggregate number of single-bit ECC errors detected in SRM
	DCGM_FI_DEV_ECC_SBE_AGG_SRM Short = 344
	// DCGM_FI_DEV_ECC_DBE_AGG_SRM represents the aggregate number of double-bit ECC errors detected in SRM
	DCGM_FI_DEV_ECC_DBE_AGG_SRM Short = 345
	// DCGM_FI_DEV_DIAG_MEMORY_RESULT is the value for ECC memory result
	DCGM_FI_DEV_DIAG_MEMORY_RESULT Short = 350
	// DCGM_FI_DEV_DIAG_DIAGNOSTIC_RESULT is the value for ECC diagnostic result
	DCGM_FI_DEV_DIAG_DIAGNOSTIC_RESULT Short = 351
	// DCGM_FI_DEV_DIAG_PCIE_RESULT is the value for ECC PCIe result
	DCGM_FI_DEV_DIAG_PCIE_RESULT Short = 352
	// DCGM_FI_DEV_DIAG_TARGETED_STRESS_RESULT is the value for ECC targeted stress result
	DCGM_FI_DEV_DIAG_TARGETED_STRESS_RESULT Short = 353
	// DCGM_FI_DEV_DIAG_TARGETED_POWER_RESULT is the value for ECC targeted power result
	DCGM_FI_DEV_DIAG_TARGETED_POWER_RESULT Short = 354
	// DCGM_FI_DEV_DIAG_MEMORY_BANDWIDTH_RESULT is the value for ECC memory bandwidth result
	DCGM_FI_DEV_DIAG_MEMORY_BANDWIDTH_RESULT Short = 355
	// DCGM_FI_DEV_DIAG_MEMTEST_RESULT is the value for ECC memtest result
	DCGM_FI_DEV_DIAG_MEMTEST_RESULT Short = 356
	// DCGM_FI_DEV_DIAG_PULSE_TEST_RESULT is the value for ECC pulse test result
	DCGM_FI_DEV_DIAG_PULSE_TEST_RESULT Short = 357
	// DCGM_FI_DEV_DIAG_EUD_RESULT is the value for ECC EUD result
	DCGM_FI_DEV_DIAG_EUD_RESULT Short = 358
	// DCGM_FI_DEV_DIAG_CPU_EUD_RESULT is the value for ECC CPU EUD result
	DCGM_FI_DEV_DIAG_CPU_EUD_RESULT Short = 359
	// DCGM_FI_DEV_DIAG_SOFTWARE_RESULT is the value for ECC software result
	DCGM_FI_DEV_DIAG_SOFTWARE_RESULT Short = 360
	// DCGM_FI_DEV_DIAG_NVBANDWIDTH_RESULT is the value for ECC NVBandwidth result
	DCGM_FI_DEV_DIAG_NVBANDWIDTH_RESULT Short = 361
	// DCGM_FI_DEV_DIAG_STATUS is the value for ECC status
	DCGM_FI_DEV_DIAG_STATUS Short = 362
	// DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_MAX is the value for ECC banks remap rows avail max
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_MAX Short = 385
	// DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_HIGH is the value for ECC banks remap rows avail high
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_HIGH Short = 386
	// DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_PARTIAL is the value for ECC banks remap rows avail partial
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_PARTIAL Short = 387
	// DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_LOW is the value for ECC banks remap rows avail low
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_LOW Short = 388
	// DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_NONE is the value for ECC banks remap rows avail none
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_NONE Short = 389
	// DCGM_FI_DEV_RETIRED_SBE is the value for ECC retired SBE
	DCGM_FI_DEV_RETIRED_SBE Short = 390
	// DCGM_FI_DEV_RETIRED_DBE is the value for ECC retired DBE
	DCGM_FI_DEV_RETIRED_DBE Short = 391
	// DCGM_FI_DEV_RETIRED_PENDING is the value for ECC retired pending
	DCGM_FI_DEV_RETIRED_PENDING Short = 392
	// DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS is the value for ECC uncorrectable remapped rows
	DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS Short = 393
	// DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS is the value for ECC correctable remapped rows
	DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS Short = 394
	// DCGM_FI_DEV_ROW_REMAP_FAILURE is the value for ECC row remap failure
	DCGM_FI_DEV_ROW_REMAP_FAILURE Short = 395
	// DCGM_FI_DEV_ROW_REMAP_PENDING is the value for ECC row remap pending
	DCGM_FI_DEV_ROW_REMAP_PENDING Short = 396
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L0 is the value for ECC NVLink CRC FLIT error count L0
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L0 Short = 400
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L1 is the value for ECC NVLink CRC FLIT error count L1
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L1 Short = 401
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L2 is the value for ECC NVLink CRC FLIT error count L2
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L2 Short = 402
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L3 is the value for ECC NVLink CRC FLIT error count L3
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L3 Short = 403
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L4 is the value for ECC NVLink CRC FLIT error count L4
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L4 Short = 404
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L5 is the value for ECC NVLink CRC FLIT error count L5
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L5 Short = 405
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL is the value for ECC NVLink CRC FLIT error count total
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL Short = 409
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L0 is the value for ECC NVLink CRC DATA error count L0
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L0 Short = 410
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L1 is the value for ECC NVLink CRC DATA error count L1
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L1 Short = 411
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L2 is the value for ECC NVLink CRC DATA error count L2
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L2 Short = 412
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L3 is the value for ECC NVLink CRC DATA error count L3
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L3 Short = 413
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L4 is the value for ECC NVLink CRC DATA error count L4
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L4 Short = 414
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L5 is the value for ECC NVLink CRC DATA error count L5
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L5 Short = 415
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_TOTAL is the value for ECC NVLink CRC DATA error count total
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_TOTAL Short = 419
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L0 is the value for ECC NVLink replay error count L0
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L0 Short = 420
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L1 is the value for ECC NVLink replay error count L1
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L1 Short = 421
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L2 is the value for ECC NVLink replay error count L2
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L2 Short = 422
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L3 is the value for ECC NVLink replay error count L3
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L3 Short = 423
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L4 is the value for ECC NVLink replay error count L4
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L4 Short = 424
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L5 is the value for ECC NVLink replay error count L5
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L5 Short = 425
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_TOTAL is the value for ECC NVLink replay error count total
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_TOTAL Short = 429
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L0 is the value for ECC NVLink recovery error count L0
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L0 Short = 430
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L1 is the value for ECC NVLink recovery error count L1
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L1 Short = 431
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L2 is the value for ECC NVLink recovery error count L2
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L2 Short = 432
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L3 is the value for ECC NVLink recovery error count L3
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L3 Short = 433
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L4 is the value for ECC NVLink recovery error count L4
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L4 Short = 434
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L5 is the value for ECC NVLink recovery error count L5
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L5 Short = 435
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_TOTAL is the value for ECC NVLink recovery error count total
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_TOTAL Short = 439
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L0 is the value for ECC NVLink bandwidth L0
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L0 Short = 440
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L1 is the value for ECC NVLink bandwidth L1
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L1 Short = 441
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L2 is the value for ECC NVLink bandwidth L2
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L2 Short = 442
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L3 is the value for ECC NVLink bandwidth L3
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L3 Short = 443
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L4 is the value for ECC NVLink bandwidth L4
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L4 Short = 444
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L5 is the value for ECC NVLink bandwidth L5
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L5 Short = 445
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL is the value for ECC NVLink bandwidth total
	DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL Short = 449
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L6 is the value for ECC NVLink CRC FLIT error count L6
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L6 Short = 451
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L7 is the value for ECC NVLink CRC FLIT error count L7
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L7 Short = 452
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L8 is the value for ECC NVLink CRC FLIT error count L8
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L8 Short = 453
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L9 is the value for ECC NVLink CRC FLIT error count L9
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L9 Short = 454
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L10 is the value for ECC NVLink CRC FLIT error count L10
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L10 Short = 455
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L11 is the value for ECC NVLink CRC FLIT error count L11
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L11 Short = 456
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L6 is the value for ECC NVLink CRC DATA error count L6
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L6 Short = 457
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L7 is the value for ECC NVLink CRC DATA error count L7
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L7 Short = 458
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L8 is the value for ECC NVLink CRC DATA error count L8
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L8 Short = 459
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L9 is the value for ECC NVLink CRC DATA error count L9
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L9 Short = 460
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L10 is the value for ECC NVLink CRC DATA error count L10
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L10 Short = 461
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L11 is the value for ECC NVLink CRC DATA error count L11
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L11 Short = 462
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L6 is the value for ECC NVLink replay error count L6
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L6 Short = 463
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L7 is the value for ECC NVLink replay error count L7
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L7 Short = 464
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L8 is the value for ECC NVLink replay error count L8
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L8 Short = 465
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L9 is the value for ECC NVLink replay error count L9
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L9 Short = 466
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L10 is the value for ECC NVLink replay error count L10
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L10 Short = 467
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L11 is the value for ECC NVLink replay error count L11
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L11 Short = 468
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L6 is the value for ECC NVLink recovery error count L6
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L6 Short = 469
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L7 is the value for ECC NVLink recovery error count L7
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L7 Short = 470
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L8 is the value for ECC NVLink recovery error count L8
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L8 Short = 471
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L9 is the value for ECC NVLink recovery error count L9
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L9 Short = 472
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L10 is the value for ECC NVLink recovery error count L10
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L10 Short = 473
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L11 is the value for ECC NVLink recovery error count L11
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L11 Short = 474
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L6 is the value for ECC NVLink bandwidth L6
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L6 Short = 475
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L7 is the value for ECC NVLink bandwidth L7
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L7 Short = 476
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L8 is the value for ECC NVLink bandwidth L8
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L8 Short = 477
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L9 is the value for ECC NVLink bandwidth L9
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L9 Short = 478
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L10 is the value for ECC NVLink bandwidth L10
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L10 Short = 479
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L11 is the value for ECC NVLink bandwidth L11
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L11 Short = 480
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L12 is the value for ECC NVLink CRC FLIT error count L12
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L12 Short = 406
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13 is the value for ECC NVLink CRC FLIT error count L13
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13 Short = 408
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14 is the value for ECC NVLink CRC FLIT error count L14
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14 Short = 409
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L15 is the value for ECC NVLink CRC FLIT error count L15
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L15 Short = 481
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L16 is the value for ECC NVLink CRC FLIT error count L16
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L16 Short = 482
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L17 is the value for ECC NVLink CRC FLIT error count L17
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L17 Short = 483
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L12 is the value for ECC NVLink CRC DATA error count L12
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L12 Short = 416
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L13 is the value for ECC NVLink CRC DATA error count L13
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L13 Short = 417
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L14 is the value for ECC NVLink CRC DATA error count L14
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L14 Short = 418
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L15 is the value for ECC NVLink CRC DATA error count L15
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L15 Short = 484
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L16 is the value for ECC NVLink CRC DATA error count L16
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L16 Short = 485
	// DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L17 is the value for ECC NVLink CRC DATA error count L17
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L17 Short = 486
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L12 is the value for ECC NVLink replay error count L12
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L12 Short = 426
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L13 is the value for ECC NVLink replay error count L13
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L13 Short = 427
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L14 is the value for ECC NVLink replay error count L14
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L14 Short = 428
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L15 is the value for ECC NVLink replay error count L15
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L15 Short = 487
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L16 is the value for ECC NVLink replay error count L16
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L16 Short = 488
	// DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L17 is the value for ECC NVLink replay error count L17
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L17 Short = 489
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L12 is the value for ECC NVLink recovery error count L12
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L12 Short = 436
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L13 is the value for ECC NVLink recovery error count L13
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L13 Short = 437
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L14 is the value for ECC NVLink recovery error count L14
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L14 Short = 438
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L15 is the value for ECC NVLink recovery error count L15
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L15 Short = 491
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L16 is the value for ECC NVLink recovery error count L16
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L16 Short = 492
	// DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L17 is the value for ECC NVLink recovery error count L17
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L17 Short = 493
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L12 is the value for ECC NVLink bandwidth L12
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L12 Short = 446
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L13 is the value for ECC NVLink bandwidth L13
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L13 Short = 447
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L14 is the value for ECC NVLink bandwidth L14
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L14 Short = 448
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L15 is the value for ECC NVLink bandwidth L15
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L15 Short = 494
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L16 is the value for ECC NVLink bandwidth L16
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L16 Short = 495
	// DCGM_FI_DEV_NVLINK_BANDWIDTH_L17 is the value for ECC NVLink bandwidth L17
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L17 Short = 496
	// DCGM_FI_DEV_NVLINK_ERROR_DL_CRC is the value for ECC NVLink error DL CRC
	DCGM_FI_DEV_NVLINK_ERROR_DL_CRC Short = 497
	// DCGM_FI_DEV_NVLINK_ERROR_DL_RECOVERY is the value for ECC NVLink error DL recovery
	DCGM_FI_DEV_NVLINK_ERROR_DL_RECOVERY Short = 498
	// DCGM_FI_DEV_NVLINK_ERROR_DL_REPLAY is the value for ECC NVLink error DL replay
	DCGM_FI_DEV_NVLINK_ERROR_DL_REPLAY Short = 499
	// DCGM_FI_DEV_VIRTUAL_MODE is the value for ECC virtual mode
	DCGM_FI_DEV_VIRTUAL_MODE Short = 500
	// DCGM_FI_DEV_SUPPORTED_TYPE_INFO is the value for ECC supported type info
	DCGM_FI_DEV_SUPPORTED_TYPE_INFO Short = 501
	// DCGM_FI_DEV_CREATABLE_VGPU_TYPE_IDS is the value for ECC creatable VGPU type IDs
	DCGM_FI_DEV_CREATABLE_VGPU_TYPE_IDS Short = 502
	// DCGM_FI_DEV_VGPU_INSTANCE_IDS is the value for ECC VGPU instance IDs
	DCGM_FI_DEV_VGPU_INSTANCE_IDS Short = 503
	// DCGM_FI_DEV_VGPU_UTILIZATIONS is the value for ECC VGPU utilizations
	DCGM_FI_DEV_VGPU_UTILIZATIONS Short = 504
	// DCGM_FI_DEV_VGPU_PER_PROCESS_UTILIZATION is the value for ECC VGPU per process utilization
	DCGM_FI_DEV_VGPU_PER_PROCESS_UTILIZATION Short = 505
	// DCGM_FI_DEV_ENC_STATS is the value for ECC enc stats
	DCGM_FI_DEV_ENC_STATS Short = 506
	// DCGM_FI_DEV_FBC_STATS is the value for ECC FBC stats
	DCGM_FI_DEV_FBC_STATS Short = 507
	// DCGM_FI_DEV_FBC_SESSIONS_INFO is the value for ECC FBC sessions info
	DCGM_FI_DEV_FBC_SESSIONS_INFO Short = 508
	// DCGM_FI_DEV_SUPPORTED_VGPU_TYPE_IDS is the value for ECC supported VGPU type IDs
	DCGM_FI_DEV_SUPPORTED_VGPU_TYPE_IDS Short = 509
	// DCGM_FI_DEV_VGPU_TYPE_INFO is the value for ECC VGPU type info
	DCGM_FI_DEV_VGPU_TYPE_INFO Short = 510
	// DCGM_FI_DEV_VGPU_TYPE_NAME is the value for ECC VGPU type name
	DCGM_FI_DEV_VGPU_TYPE_NAME Short = 511
	// DCGM_FI_DEV_VGPU_TYPE_CLASS is the value for ECC VGPU type class
	DCGM_FI_DEV_VGPU_TYPE_CLASS Short = 512
	// DCGM_FI_DEV_VGPU_TYPE_LICENSE is the value for ECC VGPU type license
	DCGM_FI_DEV_VGPU_TYPE_LICENSE Short = 513
	// DCGM_FI_DEV_VGPU_VM_ID represents the VGPU VM ID
	DCGM_FI_DEV_VGPU_VM_ID Short = 520
	// DCGM_FI_DEV_VGPU_VM_NAME represents the VGPU VM name
	DCGM_FI_DEV_VGPU_VM_NAME Short = 521
	// DCGM_FI_DEV_VGPU_TYPE represents the VGPU type
	DCGM_FI_DEV_VGPU_TYPE Short = 522
	// DCGM_FI_DEV_VGPU_UUID represents the VGPU UUID
	DCGM_FI_DEV_VGPU_UUID Short = 523
	// DCGM_FI_DEV_VGPU_DRIVER_VERSION represents the VGPU driver version
	DCGM_FI_DEV_VGPU_DRIVER_VERSION Short = 524
	// DCGM_FI_DEV_VGPU_MEMORY_USAGE represents the VGPU memory usage
	DCGM_FI_DEV_VGPU_MEMORY_USAGE Short = 525
	// DCGM_FI_DEV_VGPU_LICENSE_STATUS represents the VGPU license status
	DCGM_FI_DEV_VGPU_LICENSE_STATUS Short = 526
	// DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT represents the VGPU frame rate limit
	DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT Short = 527
	// DCGM_FI_DEV_VGPU_ENC_STATS represents the VGPU encoder statistics
	DCGM_FI_DEV_VGPU_ENC_STATS Short = 528
	// DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO represents the VGPU encoder sessions information
	DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO Short = 529
	// DCGM_FI_DEV_VGPU_FBC_STATS represents the VGPU frame buffer capture statistics
	DCGM_FI_DEV_VGPU_FBC_STATS Short = 530
	// DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO represents the VGPU frame buffer capture sessions information
	DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO Short = 531
	// DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATE represents the VGPU instance license state
	DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATE Short = 532
	// DCGM_FI_DEV_VGPU_PCI_ID represents the VGPU PCI ID
	DCGM_FI_DEV_VGPU_PCI_ID Short = 533
	// DCGM_FI_DEV_VGPU_VM_GPU_INSTANCE_ID represents the VGPU VM GPU instance ID
	DCGM_FI_DEV_VGPU_VM_GPU_INSTANCE_ID Short = 534
	// DCGM_FI_FIRST_VGPU_FIELD_ID is the value for ECC first VGPU field ID
	DCGM_FI_FIRST_VGPU_FIELD_ID Short = 520
	// DCGM_FI_LAST_VGPU_FIELD_ID is the value for ECC last VGPU field ID
	DCGM_FI_LAST_VGPU_FIELD_ID Short = 570
	// DCGM_FI_DEV_PLATFORM_INFINIBAND_GUID is the value for ECC platform InfiniBand GUID
	DCGM_FI_DEV_PLATFORM_INFINIBAND_GUID Short = 571
	// DCGM_FI_DEV_PLATFORM_CHASSIS_SERIAL_NUMBER is the value for ECC platform chassis serial number
	DCGM_FI_DEV_PLATFORM_CHASSIS_SERIAL_NUMBER Short = 572
	// DCGM_FI_DEV_PLATFORM_CHASSIS_SLOT_NUMBER is the value for ECC platform chassis slot number
	DCGM_FI_DEV_PLATFORM_CHASSIS_SLOT_NUMBER Short = 573
	// DCGM_FI_DEV_PLATFORM_TRAY_INDEX is the value for ECC platform tray index
	DCGM_FI_DEV_PLATFORM_TRAY_INDEX Short = 574
	// DCGM_FI_DEV_PLATFORM_HOST_ID is the value for ECC platform host ID
	DCGM_FI_DEV_PLATFORM_HOST_ID Short = 575
	// DCGM_FI_DEV_PLATFORM_PEER_TYPE is the value for ECC platform peer type
	DCGM_FI_DEV_PLATFORM_PEER_TYPE Short = 576
	// DCGM_FI_DEV_PLATFORM_MODULE_ID is the value for ECC platform module ID
	DCGM_FI_DEV_PLATFORM_MODULE_ID Short = 577
	// DCGM_FI_FIRST_NVSWITCH_FIELD_ID is the value for ECC first NVSwitch field ID
	DCGM_FI_FIRST_NVSWITCH_FIELD_ID Short = 700
	// DCGM_FI_DEV_NVSWITCH_VOLTAGE_MVOLT represents the NVSwitch voltage in millivolts
	DCGM_FI_DEV_NVSWITCH_VOLTAGE_MVOLT Short = 701
	// DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ represents the NVSwitch IDDQ current
	DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ Short = 702
	// DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_REV represents the NVSwitch IDDQ current revision
	DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_REV Short = 703
	// DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_DVDD represents the NVSwitch IDDQ current for DVDD
	DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_DVDD Short = 704
	// DCGM_FI_DEV_NVSWITCH_POWER_VDD represents the NVSwitch VDD power consumption
	DCGM_FI_DEV_NVSWITCH_POWER_VDD Short = 705
	// DCGM_FI_DEV_NVSWITCH_POWER_DVDD represents the NVSwitch DVDD power consumption
	DCGM_FI_DEV_NVSWITCH_POWER_DVDD Short = 706
	// DCGM_FI_DEV_NVSWITCH_POWER_HVDD represents the NVSwitch HVDD power consumption
	DCGM_FI_DEV_NVSWITCH_POWER_HVDD Short = 707
	// DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX represents the NVSwitch link transmit throughput in KB/s
	DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX Short = 780
	// DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX represents the NVSwitch link receive throughput in KB/s
	DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX Short = 781
	// DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS represents the number of fatal errors on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS Short = 782
	// DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS represents the number of non-fatal errors on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS Short = 783
	// DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS represents the number of replay errors on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS Short = 784
	// DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS represents the number of recovery errors on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS Short = 785
	// DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS represents the number of FLIT errors on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS Short = 786
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS represents the number of CRC errors on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS Short = 787
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS represents the number of ECC errors on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS Short = 788
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0 is the value for ECC NVSwitch link latency low VC0
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0 Short = 789
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1 is the value for ECC NVSwitch link latency low VC1
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1 Short = 790
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2 is the value for ECC NVSwitch link latency low VC2
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2 Short = 791
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3 is the value for ECC NVSwitch link latency low VC3
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3 Short = 792
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0 is the value for ECC NVSwitch link latency medium VC0
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0 Short = 793
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1 is the value for ECC NVSwitch link latency medium VC1
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1 Short = 794
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2 is the value for ECC NVSwitch link latency medium VC2
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2 Short = 795
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3 is the value for ECC NVSwitch link latency medium VC3
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3 Short = 796
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0 is the value for ECC NVSwitch link latency high VC0
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0 Short = 797
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1 is the value for ECC NVSwitch link latency high VC1
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1 Short = 798
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2 is the value for ECC NVSwitch link latency high VC2
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2 Short = 799
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3 is the value for ECC NVSwitch link latency high VC3
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3 Short = 800
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0 is the value for ECC NVSwitch link latency panic VC0
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0 Short = 801
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1 is the value for ECC NVSwitch link latency panic VC1
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1 Short = 802
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2 is the value for ECC NVSwitch link latency panic VC2
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2 Short = 803
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3 is the value for ECC NVSwitch link latency panic VC3
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3 Short = 804
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0 represents the latency counter for virtual channel 0 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0 Short = 805
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1 represents the latency counter for virtual channel 1 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1 Short = 806
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2 represents the latency counter for virtual channel 2 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2 Short = 807
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3 represents the latency counter for virtual channel 3 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3 Short = 808
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0 represents the number of CRC errors on lane 0 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0 Short = 809
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1 represents the number of CRC errors on lane 1 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1 Short = 810
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2 represents the number of CRC errors on lane 2 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2 Short = 811
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3 represents the number of CRC errors on lane 3 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3 Short = 812
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0 represents the number of ECC errors on lane 0 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0 Short = 813
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1 represents the number of ECC errors on lane 1 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1 Short = 814
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2 represents the number of ECC errors on lane 2 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2 Short = 815
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3 represents the number of ECC errors on lane 3 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3 Short = 816
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4 represents the number of CRC errors on lane 4 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4 Short = 817
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5 represents the number of CRC errors on lane 5 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5 Short = 818
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6 represents the number of CRC errors on lane 6 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6 Short = 819
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7 represents the number of CRC errors on lane 7 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7 Short = 820
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4 represents the number of ECC errors on lane 4 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4 Short = 821
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5 represents the number of ECC errors on lane 5 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5 Short = 822
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6 represents the number of ECC errors on lane 6 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6 Short = 823
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE7 represents the number of ECC errors on lane 7 of the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE7 Short = 824
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L0 represents the transmit bandwidth for NVLink lane 0 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L0 Short = 825
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L1 represents the transmit bandwidth for NVLink lane 1 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L1 Short = 826
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L2 represents the transmit bandwidth for NVLink lane 2 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L2 Short = 827
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L3 represents the transmit bandwidth for NVLink lane 3 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L3 Short = 828
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L4 represents the transmit bandwidth for NVLink lane 4 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L4 Short = 829
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L5 represents the transmit bandwidth for NVLink lane 5 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L5 Short = 830
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L6 represents the transmit bandwidth for NVLink lane 6 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L6 Short = 831
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L7 represents the transmit bandwidth for NVLink lane 7 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L7 Short = 832
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L8 represents the transmit bandwidth for NVLink lane 8 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L8 Short = 833
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L9 represents the transmit bandwidth for NVLink lane 9 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L9 Short = 834
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L10 represents the transmit bandwidth for NVLink lane 10 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L10 Short = 835
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L11 represents the transmit bandwidth for NVLink lane 11 in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L11 Short = 836
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0 represents the receive bandwidth for NVLink lane 0 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0 Short = 837
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1 represents the receive bandwidth for NVLink lane 1 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1 Short = 838
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2 represents the receive bandwidth for NVLink lane 2 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2 Short = 839
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3 represents the receive bandwidth for NVLink lane 3 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3 Short = 840
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4 represents the receive bandwidth for NVLink lane 4 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4 Short = 841
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5 represents the receive bandwidth for NVLink lane 5 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5 Short = 842
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6 represents the receive bandwidth for NVLink lane 6 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6 Short = 843
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7 represents the receive bandwidth for NVLink lane 7 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7 Short = 844
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8 represents the receive bandwidth for NVLink lane 8 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8 Short = 845
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9 represents the receive bandwidth for NVLink lane 9 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9 Short = 846
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10 represents the receive bandwidth for NVLink lane 10 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10 Short = 847
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11 represents the receive bandwidth for NVLink lane 11 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11 Short = 848
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL represents the total transmit bandwidth for all NVLink lanes in KB/s
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL Short = 849
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL represents the total receive bandwidth for all NVLink lanes in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL Short = 850
	// DCGM_FI_PROF_GR_ENGINE_ACTIVE represents the percentage of time the graphics engine was active
	DCGM_FI_PROF_GR_ENGINE_ACTIVE Short = 1001
	// DCGM_FI_PROF_SM_ACTIVE represents the percentage of time the streaming multiprocessors (SM) were active
	DCGM_FI_PROF_SM_ACTIVE Short = 1002
	// DCGM_FI_PROF_SM_OCCUPANCY represents the percentage of streaming multiprocessors (SM) warps residency
	DCGM_FI_PROF_SM_OCCUPANCY Short = 1003
	// DCGM_FI_PROF_PIPE_TENSOR_ACTIVE represents the percentage of time the tensor (HMMA) pipe was active
	DCGM_FI_PROF_PIPE_TENSOR_ACTIVE Short = 1004
	// DCGM_FI_PROF_DRAM_ACTIVE represents the percentage of time the device memory interface was active
	DCGM_FI_PROF_DRAM_ACTIVE Short = 1005
	// DCGM_FI_PROF_PIPE_FP64_ACTIVE represents the percentage of time the FP64 pipe was active
	DCGM_FI_PROF_PIPE_FP64_ACTIVE Short = 1006
	// DCGM_FI_PROF_PIPE_FP32_ACTIVE represents the percentage of time the FP32 pipe was active
	DCGM_FI_PROF_PIPE_FP32_ACTIVE Short = 1007
	// DCGM_FI_PROF_PIPE_FP16_ACTIVE represents the percentage of time the FP16 pipe was active
	DCGM_FI_PROF_PIPE_FP16_ACTIVE Short = 1008
	// DCGM_FI_PROF_PCIE_TX_BYTES represents the number of bytes transmitted through PCIe TX (in bytes)
	DCGM_FI_PROF_PCIE_TX_BYTES Short = 1009
	// DCGM_FI_PROF_PCIE_RX_BYTES represents the number of bytes received through PCIe RX (in bytes)
	DCGM_FI_PROF_PCIE_RX_BYTES Short = 1010
	// DCGM_FI_PROF_NVLINK_TX_BYTES represents the number of bytes transmitted through NVLink TX (in bytes)
	DCGM_FI_PROF_NVLINK_TX_BYTES Short = 1011
	// DCGM_FI_PROF_NVLINK_RX_BYTES represents the number of bytes received through NVLink RX (in bytes)
	DCGM_FI_PROF_NVLINK_RX_BYTES Short = 1012
	// DCGM_FI_PROF_PIPE_TENSOR_IMMA_ACTIVE represents the percentage of time the IMMA tensor pipe was active
	DCGM_FI_PROF_PIPE_TENSOR_IMMA_ACTIVE Short = 1013
	// DCGM_FI_PROF_PIPE_TENSOR_HMMA_ACTIVE represents the percentage of time the HMMA tensor pipe was active
	DCGM_FI_PROF_PIPE_TENSOR_HMMA_ACTIVE Short = 1014
	// DCGM_FI_PROF_PIPE_TENSOR_DFMA_ACTIVE represents the percentage of time the DFMA tensor pipe was active
	DCGM_FI_PROF_PIPE_TENSOR_DFMA_ACTIVE Short = 1015
	// DCGM_FI_PROF_PCIE_REPLAY_COUNTER represents the number of PCIe replays that occurred
	DCGM_FI_PROF_PCIE_REPLAY_COUNTER Short = 1016
	// DCGM_FI_PROF_REMAPPED_COR represents the number of correctable remapped memory errors
	DCGM_FI_PROF_REMAPPED_COR Short = 1017
	// DCGM_FI_PROF_REMAPPED_UNCOR represents the number of uncorrectable remapped memory errors
	DCGM_FI_PROF_REMAPPED_UNCOR Short = 1018
	// DCGM_FI_PROF_REMAPPED_PENDING represents the number of pending remapped memory pages
	DCGM_FI_PROF_REMAPPED_PENDING Short = 1019
	// DCGM_FI_PROF_REMAPPED_FAILURE represents the number of remapped memory page failures
	DCGM_FI_PROF_REMAPPED_FAILURE Short = 1020
	// DCGM_FI_DEV_CPU_UTIL_TOTAL represents the total CPU utilization across all processes
	DCGM_FI_DEV_CPU_UTIL_TOTAL Short = 1100
	// DCGM_FI_DEV_CPU_UTIL_KERNEL represents the CPU utilization in kernel mode across all processes
	DCGM_FI_DEV_CPU_UTIL_KERNEL Short = 1101
	// DCGM_FI_DEV_CPU_UTIL_USER represents the CPU utilization in user mode across all processes
	DCGM_FI_DEV_CPU_UTIL_USER Short = 1102
	// DCGM_FI_DEV_POWER_PER_PROCESS_TOTAL represents the total power usage per process in milliwatts
	DCGM_FI_DEV_POWER_PER_PROCESS_TOTAL Short = 1103
	// DCGM_FI_DEV_POWER_PER_PROCESS_MAX represents the maximum power usage per process in milliwatts
	DCGM_FI_DEV_POWER_PER_PROCESS_MAX Short = 1104
	// DCGM_FI_DEV_POWER_PER_PROCESS_MIN represents the minimum power usage per process in milliwatts
	DCGM_FI_DEV_POWER_PER_PROCESS_MIN Short = 1105
	// DCGM_FI_DEV_CPU_TEMP_CURRENT represents the current CPU temperature in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_CURRENT Short = 1110
	// DCGM_FI_DEV_CPU_TEMP_SLOWDOWN represents the CPU temperature slowdown threshold in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_SLOWDOWN Short = 1111
	// DCGM_FI_DEV_CPU_TEMP_SHUTDOWN represents the CPU temperature shutdown threshold in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_SHUTDOWN Short = 1112
	// DCGM_FI_DEV_CPU_TEMP_MAX represents the maximum CPU temperature in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_MAX Short = 1113
	// DCGM_FI_DEV_CPU_TEMP_MIN represents the minimum CPU temperature in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_MIN Short = 1114
	// DCGM_FI_DEV_CPU_TEMP_AVG represents the average CPU temperature in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_AVG Short = 1115
	// DCGM_FI_DEV_CPU_CLOCK_CURRENT represents the current CPU clock frequency in MHz
	DCGM_FI_DEV_CPU_CLOCK_CURRENT Short = 1120
	// DCGM_FI_DEV_CPU_CLOCK_MIN represents the minimum CPU clock frequency in MHz
	DCGM_FI_DEV_CPU_CLOCK_MIN Short = 1121
	// DCGM_FI_DEV_CPU_CLOCK_MAX represents the maximum CPU clock frequency in MHz
	DCGM_FI_DEV_CPU_CLOCK_MAX Short = 1122
	// DCGM_FI_DEV_CPU_CLOCK_THROTTLE_REASONS represents the CPU clock throttling reason bitmask
	DCGM_FI_DEV_CPU_CLOCK_THROTTLE_REASONS Short = 1123
	// DCGM_FI_DEV_CPU_POWER_CURRENT represents the current CPU power usage in milliwatts
	DCGM_FI_DEV_CPU_POWER_CURRENT Short = 1130
	// DCGM_FI_DEV_CPU_POWER_MIN represents the minimum CPU power usage in milliwatts
	DCGM_FI_DEV_CPU_POWER_MIN Short = 1131
	// DCGM_FI_DEV_CPU_POWER_MAX represents the maximum CPU power usage in milliwatts
	DCGM_FI_DEV_CPU_POWER_MAX Short = 1132
	// DCGM_FI_DEV_CPU_POWER_LIMIT represents the CPU power limit in milliwatts
	DCGM_FI_DEV_CPU_POWER_LIMIT Short = 1133
	// DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT is the value for ECC DEV SysIO Power Util Current
	DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT Short = 1132
	// DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT is the value for ECC DEV Module Power Util Current
	DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT Short = 1133
	// DCGM_FI_DEV_CPU_VENDOR is the value for ECC DEV CPU Vendor
	DCGM_FI_DEV_CPU_VENDOR Short = 1140
	// DCGM_FI_DEV_CPU_MODEL is the value for ECC DEV CPU Model
	DCGM_FI_DEV_CPU_MODEL Short = 1141
	// DCGM_FI_DEV_NVLINK_COUNT_TX_PACKETS is the value for ECC DEV NVLink Count TX Packets
	DCGM_FI_DEV_NVLINK_COUNT_TX_PACKETS Short = 1200
	// DCGM_FI_DEV_NVLINK_COUNT_TX_BYTES is the value for ECC DEV NVLink Count TX Bytes
	DCGM_FI_DEV_NVLINK_COUNT_TX_BYTES Short = 1201
	// DCGM_FI_DEV_NVLINK_COUNT_RX_PACKETS is the value for ECC DEV NVLink Count RX Packets
	DCGM_FI_DEV_NVLINK_COUNT_RX_PACKETS Short = 1202
	// DCGM_FI_DEV_NVLINK_COUNT_RX_BYTES is the value for ECC DEV NVLink Count RX Bytes
	DCGM_FI_DEV_NVLINK_COUNT_RX_BYTES Short = 1203
	// DCGM_FI_DEV_NVLINK_COUNT_RX_MALFORMED_PACKET_ERRORS is the value for ECC DEV NVLink Count RX Malformed Packet Errors
	DCGM_FI_DEV_NVLINK_COUNT_RX_MALFORMED_PACKET_ERRORS Short = 1204
	// DCGM_FI_DEV_NVLINK_COUNT_RX_BUFFER_OVERRUN_ERRORS is the value for ECC DEV NVLink Count RX Buffer Overrun Errors
	DCGM_FI_DEV_NVLINK_COUNT_RX_BUFFER_OVERRUN_ERRORS Short = 1205
	// DCGM_FI_DEV_NVLINK_COUNT_RX_ERRORS is the value for ECC DEV NVLink Count RX Errors
	DCGM_FI_DEV_NVLINK_COUNT_RX_ERRORS Short = 1206
	// DCGM_FI_DEV_NVLINK_COUNT_RX_REMOTE_ERRORS is the value for ECC DEV NVLink Count RX Remote Errors
	DCGM_FI_DEV_NVLINK_COUNT_RX_REMOTE_ERRORS Short = 1207
	// DCGM_FI_DEV_NVLINK_COUNT_RX_GENERAL_ERRORS is the value for ECC DEV NVLink Count RX General Errors
	DCGM_FI_DEV_NVLINK_COUNT_RX_GENERAL_ERRORS Short = 1208
	// DCGM_FI_DEV_NVLINK_COUNT_LOCAL_LINK_INTEGRITY_ERRORS is the value for ECC DEV NVLink Count Local Link Integrity Errors
	DCGM_FI_DEV_NVLINK_COUNT_LOCAL_LINK_INTEGRITY_ERRORS Short = 1209
	// DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS is the value for ECC DEV NVLink Count TX Discards
	DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS Short = 1210
	// DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_SUCCESSFUL_EVENTS is the value for ECC DEV NVLink Count Link Recovery Successful Events
	DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_SUCCESSFUL_EVENTS Short = 1211
	// DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_FAILED_EVENTS is the value for ECC DEV NVLink Count Link Recovery Failed Events
	DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_FAILED_EVENTS Short = 1212
	// DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_EVENTS is the value for ECC DEV NVLink Count Link Recovery Events
	DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_EVENTS Short = 1213
	// DCGM_FI_DEV_NVLINK_COUNT_RX_SYMBOL_ERRORS is the value for ECC DEV NVLink Count RX Symbol Errors
	DCGM_FI_DEV_NVLINK_COUNT_RX_SYMBOL_ERRORS Short = 1214
	// DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER is the value for ECC DEV NVLink Count Symbol BER
	DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER Short = 1215
	// DCGM_FI_DEV_CONNECTX_HEALTH is the value for ECC DEV ConnectX Health
	DCGM_FI_DEV_CONNECTX_HEALTH Short = 1300
	// DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH is the value for ECC DEV ConnectX Active PCIe Link Width
	DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH Short = 1301
	// DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED is the value for ECC DEV ConnectX Active PCIe Link Speed
	DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED Short = 1302
	// DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH is the value for ECC DEV ConnectX Expect PCIe Link Width
	DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH Short = 1303
	// DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED is the value for ECC DEV ConnectX Expect PCIe Link Speed
	DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED Short = 1304
	// DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS is the value for ECC DEV ConnectX Correctable Err Status
	DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS Short = 1305
	// DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK is the value for ECC DEV ConnectX Correctable Err Mask
	DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK Short = 1306
	// DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS is the value for ECC DEV ConnectX Uncorrectable Err Status
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS Short = 1307
	// DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK is the value for ECC DEV ConnectX Uncorrectable Err Mask
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK Short = 1308
	// DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY is the value for ECC DEV ConnectX Uncorrectable Err Severity
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY Short = 1309
	// DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE is the value for ECC DEV ConnectX Device Temperature
	DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE Short = 1310
	// DCGM_FI_DEV_LAST_CONNECTX_FIELD_ID represents the last field ID for ConnectX fields
	DCGM_FI_DEV_LAST_CONNECTX_FIELD_ID Short = 1399
	// DCGM_FI_MAX_FIELDS represents 1 greater than maximum fields above. This is the 1 greater than the maximum field id that could be allocated
	DCGM_FI_MAX_FIELDS Short = 1311
)
