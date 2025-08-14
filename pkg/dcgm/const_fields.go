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
	// DCGM_FI_CUDA_DRIVER_VERSION represents the CUDA driver version. Retrieves a number with the major value in the thousands place and the minor value in the hundreds place. (e.g. CUDA 11.1 = 11100)
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
	// DCGM_FI_DEV_P2P_NVLINK_STATUS represents the NVLINK P2P status for the device
	DCGM_FI_DEV_P2P_NVLINK_STATUS Short = 64
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
	DCGM_FI_DEV_FABRIC_MANAGER_STATUS Short = 170
	// DCGM_FI_DEV_FABRIC_MANAGER_ERROR_CODE is the value for fabric manager error code
	// NOTE: this is not populated unless the fabric manager completed startup
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
	DCGM_FI_DEV_ECC_CURRENT Short = 300
	// DCGM_FI_DEV_ECC_PENDING is the value for ECC pending
	DCGM_FI_DEV_ECC_PENDING Short = 301
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
	// DCGM_FI_DEV_THRESHOLD_SRM represents the threshold for SRM ECC errors
	DCGM_FI_DEV_THRESHOLD_SRM Short = 346
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
	// DCGM_FI_DEV_GPU_NVLINK_ERRORS is the value for GPU NVLink error information
	DCGM_FI_DEV_GPU_NVLINK_ERRORS Short = 450
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
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13 Short = 407
	// DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14 is the value for ECC NVLink CRC FLIT error count L14
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14 Short = 408
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
	// DCGM_FI_DEV_NVSWITCH_POWER_VDD represents the NVSwitch VDD power consumption in watts
	DCGM_FI_DEV_NVSWITCH_POWER_VDD Short = 705
	// DCGM_FI_DEV_NVSWITCH_POWER_DVDD represents the NVSwitch DVDD power consumption in watts
	DCGM_FI_DEV_NVSWITCH_POWER_DVDD Short = 706
	// DCGM_FI_DEV_NVSWITCH_POWER_HVDD represents the NVSwitch HVDD power consumption in watts
	DCGM_FI_DEV_NVSWITCH_POWER_HVDD Short = 707
	// DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX represents the NVSwitch Tx Throughput Counter for ports 0-17 in KB/s
	DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX Short = 780
	// DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX represents the NVSwitch Rx Throughput Counter for ports 0-17 in KB/s
	DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX Short = 781
	// DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS represents the number of fatal errors for ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS Short = 782
	// DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS represents the number of non-fatal errors for ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS Short = 783
	// DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS represents the number of replay errors for ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS Short = 784
	// DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS represents the number of recovery errors for ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS Short = 785
	// DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS represents the number of FLIT errors for ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS Short = 786
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS represents the number of CRC errors for ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS Short = 787
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS represents the number of ECC errors for ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS Short = 788
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0 is the value for Nvlink lane latency low lane0 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0 Short = 789
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1 is the value forNvlink lane latency low lane1 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1 Short = 790
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2 is the value for Nvlink lane latency low lane2 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2 Short = 791
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3 is the value for Nvlink lane latency low lane3 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3 Short = 792
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0 is the value for Nvlink lane latency medium lane0 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0 Short = 793
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1 is the value for Nvlink lane latency medium lane1 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1 Short = 794
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2 is the value for Nvlink lane latency medium lane2 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2 Short = 795
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3 is the value for Nvlink lane latency medium lane3 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3 Short = 796
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0 is the value for Nvlink lane latency high lane0 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0 Short = 797
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1 is the value for Nvlink lane latency high lane1 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1 Short = 798
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2 is the value for Nvlink lane latency high lane2 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2 Short = 799
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3 is the value for Nvlink lane latency high lane3 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3 Short = 800
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0 is the value for Nvlink lane latency panic lane0 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0 Short = 801
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1 is the value for Nvlink lane latency panic lane1 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1 Short = 802
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2 is the value for Nvlink lane latency panic lane2 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2 Short = 803
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3 is the value for Nvlink lane latency panic lane3 counter
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3 Short = 804
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0 represents the latency counter for virtual channel 0 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0 Short = 805
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1 represents the latency counter for virtual channel 1 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1 Short = 806
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2 represents the latency counter for virtual channel 2 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2 Short = 807
	// DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3 represents the latency counter for virtual channel 3 on the NVSwitch link
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3 Short = 808
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0 represents the number of CRC errors on lane 0 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0 Short = 809
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1 represents the number of CRC errors on lane 1 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1 Short = 810
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2 represents the number of CRC errors on lane 2 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2 Short = 811
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3 represents the number of CRC errors on lane 3 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3 Short = 812
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0 represents the number of ECC errors on lane 0 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0 Short = 813
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1 represents the number of ECC errors on lane 1 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1 Short = 814
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2 represents the number of ECC errors on lane 2 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2 Short = 815
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3 represents the number of ECC errors on lane 3 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3 Short = 816
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4 represents the number of CRC errors on lane 4 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4 Short = 817
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5 represents the number of CRC errors on lane 5 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5 Short = 818
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6 represents the number of CRC errors on lane 6 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6 Short = 819
	// DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7 represents the number of CRC errors on lane 7 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7 Short = 820
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4 represents the number of ECC errors on lane 4 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4 Short = 821
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5 represents the number of ECC errors on lane 5 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5 Short = 822
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6 represents the number of ECC errors on lane 6 on ports 0-17
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6 Short = 823
	// DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE7 represents the number of ECC errors on lane 7 on ports 0-17
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
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L12 represents the NV Link TX Bandwidth Counter for Lane 12
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L12 Short = 837
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L13 represents the NV Link TX Bandwidth Counter for Lane 13
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L13 Short = 838
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L14 represents the NV Link TX Bandwidth Counter for Lane 14
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L14 Short = 839
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L15 represents the NV Link TX Bandwidth Counter for Lane 15
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L15 Short = 840
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L16 represents the NV Link TX Bandwidth Counter for Lane 16
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L16 Short = 841
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L17 represents the NV Link TX Bandwidth Counter for Lane 17
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L17 Short = 842
	// DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL represents the NV Link Bandwidth Counter total for all TX Lanes
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL Short = 843
	// DCGM_FI_DEV_NVSWITCH_FATAL_ERRORS represents the NVSwitch fatal error information.
	// Note: value field indicates the specific SXid reported
	DCGM_FI_DEV_NVSWITCH_FATAL_ERRORS Short = 856
	// DCGM_FI_DEV_NVSWITCH_NON_FATAL_ERRORS represents the NVSwitch non fatal error information.
	DCGM_FI_DEV_NVSWITCH_NON_FATAL_ERRORS Short = 857
	// DCGM_FI_DEV_NVSWITCH_TEMPERATURE_CURRENT represents the NVSwitch current temperature.
	DCGM_FI_DEV_NVSWITCH_TEMPERATURE_CURRENT Short = 858
	// DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SLOWDOWN represents the NVSwitch limit slowdown temperature
	DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SLOWDOWN Short = 859
	// DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SHUTDOWN represents the NVSwitch limit shutdown temperature
	DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SHUTDOWN Short = 860
	// DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX represents the NVSwitch throughput Tx
	DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX Short = 861
	// DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX represents the NVSwitch throughput Rx
	DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX Short = 862
	// DCGM_FI_DEV_NVSWITCH_PHYS_ID represents the NVSwitch physical ID
	DCGM_FI_DEV_NVSWITCH_PHYS_ID Short = 863
	// DCGM_FI_DEV_NVSWITCH_RESET_REQUIRED represents the NVSwitch reset required
	DCGM_FI_DEV_NVSWITCH_RESET_REQUIRED Short = 864
	// DCGM_FI_DEV_NVSWITCH_LINK_ID represents the NVSwitch link ID
	DCGM_FI_DEV_NVSWITCH_LINK_ID Short = 865
	// DCGM_FI_DEV_NVSWITCH_PCIE_DOMAIN represents the NVSwitch PCIe domain
	DCGM_FI_DEV_NVSWITCH_PCIE_DOMAIN Short = 866
	// DCGM_FI_DEV_NVSWITCH_PCIE_BUS represents the NVSwitch PCIe bus
	DCGM_FI_DEV_NVSWITCH_PCIE_BUS Short = 867
	// DCGM_FI_DEV_NVSWITCH_PCIE_DEVICE represents the NVSwitch PCIe device
	DCGM_FI_DEV_NVSWITCH_PCIE_DEVICE Short = 868
	// DCGM_FI_DEV_NVSWITCH_PCIE_FUNCTION represents the NVSwitch PCIe function
	DCGM_FI_DEV_NVSWITCH_PCIE_FUNCTION Short = 869
	// DCGM_FI_DEV_NVSWITCH_LINK_STATUS represents the NVSwitch link status UNKNOWN:-1 OFF:0 SAFE:1 ACTIVE:2 ERROR:3
	DCGM_FI_DEV_NVSWITCH_LINK_STATUS Short = 870
	// DCGM_FI_DEV_NVSWITCH_LINK_TYPE represents the NVSwitch link type GPU/Switch
	DCGM_FI_DEV_NVSWITCH_LINK_TYPE Short = 871
	// DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN represents the NVSwitch remote PCIe domain
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN Short = 872
	// DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS represents the NVSwitch remote PCIe bus
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS Short = 873
	// DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE represents the NVSwitch remote PCIe device
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE Short = 874
	// DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION represents the NVSwitch remote PCIe function
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION Short = 875
	// DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_ID represents the NVSwitch link device link ID
	DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_ID Short = 876
	// DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_SID represents the NVSwitch link device link SID
	DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_SID Short = 877
	// DCGM_FI_DEV_NVSWITCH_DEVICE_UUID represents the NVSwitch device UUID
	DCGM_FI_DEV_NVSWITCH_DEVICE_UUID Short = 878
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0 represents the receive bandwidth for NVLink lane 0 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0 Short = 879
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1 represents the receive bandwidth for NVLink lane 1 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1 Short = 880
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2 represents the receive bandwidth for NVLink lane 2 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2 Short = 881
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3 represents the receive bandwidth for NVLink lane 3 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3 Short = 882
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4 represents the receive bandwidth for NVLink lane 4 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4 Short = 883
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5 represents the receive bandwidth for NVLink lane 5 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5 Short = 884
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6 represents the receive bandwidth for NVLink lane 6 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6 Short = 885
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7 represents the receive bandwidth for NVLink lane 7 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7 Short = 886
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8 represents the receive bandwidth for NVLink lane 8 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8 Short = 887
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9 represents the receive bandwidth for NVLink lane 9 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9 Short = 888
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10 represents the receive bandwidth for NVLink lane 10 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10 Short = 889
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11 represents the receive bandwidth for NVLink lane 11 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11 Short = 890
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L12 represents the receive bandwidth for NVLink lane 12 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L12 Short = 891
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L13 represents the receive bandwidth for NVLink lane 13 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L13 Short = 892
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L14 represents the receive bandwidth for NVLink lane 14 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L14 Short = 893
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L15 represents the receive bandwidth for NVLink lane 15 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L15 Short = 894
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L16 represents the receive bandwidth for NVLink lane 16 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L16 Short = 895
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L17 represents the receive bandwidth for NVLink lane 17 in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L17 Short = 896
	// DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL represents the total receive bandwidth for all NVLink lanes in KB/s
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL Short = 897

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
	// DCGM_FI_PROF_PIPE_INT_ACTIVE represents the ratio of cycles the integer pipe is active
	DCGM_FI_PROF_PIPE_INT_ACTIVE Short = 1016
	// DCGM_FI_PROF_NVDEC0_ACTIVE represents the ratio of cycles the NVDEC engine 0 is active
	DCGM_FI_PROF_NVDEC0_ACTIVE Short = 1017
	// DCGM_FI_PROF_NVDEC1_ACTIVE represents the ratio of cycles the NVDEC engine 1 is active
	DCGM_FI_PROF_NVDEC1_ACTIVE Short = 1018
	// DCGM_FI_PROF_NVDEC2_ACTIVE represents the ratio of cycles the NVDEC engine 2 is active
	DCGM_FI_PROF_NVDEC2_ACTIVE Short = 1019
	// DCGM_FI_PROF_NVDEC3_ACTIVE represents the ratio of cycles the NVDEC engine 3 is active
	DCGM_FI_PROF_NVDEC3_ACTIVE Short = 1020
	// DCGM_FI_PROF_NVDEC4_ACTIVE represents the ratio of cycles the NVDEC engine 4 is active
	DCGM_FI_PROF_NVDEC4_ACTIVE Short = 1021
	// DCGM_FI_PROF_NVDEC5_ACTIVE represents the ratio of cycles the NVDEC engine 5 is active
	DCGM_FI_PROF_NVDEC5_ACTIVE Short = 1022
	// DCGM_FI_PROF_NVDEC6_ACTIVE represents the ratio of cycles the NVDEC engine 6 is active
	DCGM_FI_PROF_NVDEC6_ACTIVE Short = 1023
	// DCGM_FI_PROF_NVDEC7_ACTIVE represents the ratio of cycles the NVDEC engine 7 is active
	DCGM_FI_PROF_NVDEC7_ACTIVE Short = 1024

	// DCGM_FI_PROF_NVJPG0_ACTIVE represents the ratio of cycles the NVJPG engine 0 is active
	DCGM_FI_PROF_NVJPG0_ACTIVE Short = 1025
	// DCGM_FI_PROF_NVJPG1_ACTIVE represents the ratio of cycles the NVJPG engine 1 is active
	DCGM_FI_PROF_NVJPG1_ACTIVE Short = 1026
	// DCGM_FI_PROF_NVJPG2_ACTIVE represents the ratio of cycles the NVJPG engine 2 is active
	DCGM_FI_PROF_NVJPG2_ACTIVE Short = 1027
	// DCGM_FI_PROF_NVJPG3_ACTIVE represents the ratio of cycles the NVJPG engine 3 is active
	DCGM_FI_PROF_NVJPG3_ACTIVE Short = 1028
	// DCGM_FI_PROF_NVJPG4_ACTIVE represents the ratio of cycles the NVJPG engine 4 is active
	DCGM_FI_PROF_NVJPG4_ACTIVE Short = 1029
	// DCGM_FI_PROF_NVJPG5_ACTIVE represents the ratio of cycles the NVJPG engine 5 is active
	DCGM_FI_PROF_NVJPG5_ACTIVE Short = 1030
	// DCGM_FI_PROF_NVJPG6_ACTIVE represents the ratio of cycles the NVJPG engine 6 is active
	DCGM_FI_PROF_NVJPG6_ACTIVE Short = 1031
	// DCGM_FI_PROF_NVJPG7_ACTIVE represents the ratio of cycles the NVJPG engine 7 is active
	DCGM_FI_PROF_NVJPG7_ACTIVE Short = 1032

	// DCGM_FI_PROF_NVOFA0_ACTIVE represents the ratio of cycles the NVOFA engine 0 is active
	DCGM_FI_PROF_NVOFA0_ACTIVE Short = 1033
	// DCGM_FI_PROF_NVOFA1_ACTIVE represents the ratio of cycles the NVOFA engine 1 is active
	DCGM_FI_PROF_NVOFA1_ACTIVE Short = 1034

	// DCGM_FI_PROF_NVLINK_L0_TX_BYTES represents the number of bytes transmitted through NVLink lane 0 in KB/s
	DCGM_FI_PROF_NVLINK_L0_TX_BYTES Short = 1040
	// DCGM_FI_PROF_NVLINK_L0_RX_BYTES represents the number of bytes received through NVLink lane 0 in KB/s
	DCGM_FI_PROF_NVLINK_L0_RX_BYTES Short = 1041
	// DCGM_FI_PROF_NVLINK_L1_TX_BYTES represents the number of bytes transmitted through NVLink lane 1 in KB/s
	DCGM_FI_PROF_NVLINK_L1_TX_BYTES Short = 1042
	// DCGM_FI_PROF_NVLINK_L1_RX_BYTES represents the number of bytes received through NVLink lane 1 in KB/s
	DCGM_FI_PROF_NVLINK_L1_RX_BYTES Short = 1043
	// DCGM_FI_PROF_NVLINK_L2_TX_BYTES represents the number of bytes transmitted through NVLink lane 2 in KB/s
	DCGM_FI_PROF_NVLINK_L2_TX_BYTES Short = 1044
	// DCGM_FI_PROF_NVLINK_L2_RX_BYTES represents the number of bytes received through NVLink lane 2 in KB/s
	DCGM_FI_PROF_NVLINK_L2_RX_BYTES Short = 1045
	// DCGM_FI_PROF_NVLINK_L3_TX_BYTES represents the number of bytes transmitted through NVLink lane 3 in KB/s
	DCGM_FI_PROF_NVLINK_L3_TX_BYTES Short = 1046
	// DCGM_FI_PROF_NVLINK_L3_RX_BYTES represents the number of bytes received through NVLink lane 3 in KB/s
	DCGM_FI_PROF_NVLINK_L3_RX_BYTES Short = 1047
	// DCGM_FI_PROF_NVLINK_L4_TX_BYTES represents the number of bytes transmitted through NVLink lane 4 in KB/s
	DCGM_FI_PROF_NVLINK_L4_TX_BYTES Short = 1048
	// DCGM_FI_PROF_NVLINK_L4_RX_BYTES represents the number of bytes received through NVLink lane 4 in KB/s
	DCGM_FI_PROF_NVLINK_L4_RX_BYTES Short = 1049
	// DCGM_FI_PROF_NVLINK_L5_TX_BYTES represents the number of bytes transmitted through NVLink lane 5 in KB/s
	DCGM_FI_PROF_NVLINK_L5_TX_BYTES Short = 1050
	// DCGM_FI_PROF_NVLINK_L5_RX_BYTES represents the number of bytes received through NVLink lane 5 in KB/s
	DCGM_FI_PROF_NVLINK_L5_RX_BYTES Short = 1051
	// DCGM_FI_PROF_NVLINK_L6_TX_BYTES represents the number of bytes transmitted through NVLink lane 6 in KB/s
	DCGM_FI_PROF_NVLINK_L6_TX_BYTES Short = 1052
	// DCGM_FI_PROF_NVLINK_L6_RX_BYTES represents the number of bytes received through NVLink lane 6 in KB/s
	DCGM_FI_PROF_NVLINK_L6_RX_BYTES Short = 1053
	// DCGM_FI_PROF_NVLINK_L7_TX_BYTES represents the number of bytes transmitted through NVLink lane 7 in KB/s
	DCGM_FI_PROF_NVLINK_L7_TX_BYTES Short = 1054
	// DCGM_FI_PROF_NVLINK_L7_RX_BYTES represents the number of bytes received through NVLink lane 7 in KB/s
	DCGM_FI_PROF_NVLINK_L7_RX_BYTES Short = 1055
	// DCGM_FI_PROF_NVLINK_L8_TX_BYTES represents the number of bytes transmitted through NVLink lane 8 in KB/s
	DCGM_FI_PROF_NVLINK_L8_TX_BYTES Short = 1056
	// DCGM_FI_PROF_NVLINK_L8_RX_BYTES represents the number of bytes received through NVLink lane 8 in KB/s
	DCGM_FI_PROF_NVLINK_L8_RX_BYTES Short = 1057
	// DCGM_FI_PROF_NVLINK_L9_TX_BYTES represents the number of bytes transmitted through NVLink lane 9 in KB/s
	DCGM_FI_PROF_NVLINK_L9_TX_BYTES Short = 1058
	// DCGM_FI_PROF_NVLINK_L9_RX_BYTES represents the number of bytes received through NVLink lane 9 in KB/s
	DCGM_FI_PROF_NVLINK_L9_RX_BYTES Short = 1059
	// DCGM_FI_PROF_NVLINK_L10_TX_BYTES represents the number of bytes transmitted through NVLink lane 10 in KB/s
	DCGM_FI_PROF_NVLINK_L10_TX_BYTES Short = 1060
	// DCGM_FI_PROF_NVLINK_L10_RX_BYTES represents the number of bytes received through NVLink lane 10 in KB/s
	DCGM_FI_PROF_NVLINK_L10_RX_BYTES Short = 1061
	// DCGM_FI_PROF_NVLINK_L11_TX_BYTES represents the number of bytes transmitted through NVLink lane 11 in KB/s
	DCGM_FI_PROF_NVLINK_L11_TX_BYTES Short = 1062
	// DCGM_FI_PROF_NVLINK_L11_RX_BYTES represents the number of bytes received through NVLink lane 11 in KB/s
	DCGM_FI_PROF_NVLINK_L11_RX_BYTES Short = 1063
	// DCGM_FI_PROF_NVLINK_L12_TX_BYTES represents the number of bytes transmitted through NVLink lane 12 in KB/s
	DCGM_FI_PROF_NVLINK_L12_TX_BYTES Short = 1064
	// DCGM_FI_PROF_NVLINK_L12_RX_BYTES represents the number of bytes received through NVLink lane 12 in KB/s
	DCGM_FI_PROF_NVLINK_L12_RX_BYTES Short = 1065
	// DCGM_FI_PROF_NVLINK_L13_TX_BYTES represents the number of bytes transmitted through NVLink lane 13 in KB/s
	DCGM_FI_PROF_NVLINK_L13_TX_BYTES Short = 1066
	// DCGM_FI_PROF_NVLINK_L13_RX_BYTES represents the number of bytes received through NVLink lane 13 in KB/s
	DCGM_FI_PROF_NVLINK_L13_RX_BYTES Short = 1067
	// DCGM_FI_PROF_NVLINK_L14_TX_BYTES represents the number of bytes transmitted through NVLink lane 14 in KB/s
	DCGM_FI_PROF_NVLINK_L14_TX_BYTES Short = 1068
	// DCGM_FI_PROF_NVLINK_L14_RX_BYTES represents the number of bytes received through NVLink lane 14 in KB/s
	DCGM_FI_PROF_NVLINK_L14_RX_BYTES Short = 1069
	// DCGM_FI_PROF_NVLINK_L15_TX_BYTES represents the number of bytes transmitted through NVLink lane 15 in KB/s
	DCGM_FI_PROF_NVLINK_L15_TX_BYTES Short = 1070
	// DCGM_FI_PROF_NVLINK_L15_RX_BYTES represents the number of bytes received through NVLink lane 15 in KB/s
	DCGM_FI_PROF_NVLINK_L15_RX_BYTES Short = 1071
	// DCGM_FI_PROF_NVLINK_L16_TX_BYTES represents the number of bytes transmitted through NVLink lane 16 in KB/s
	DCGM_FI_PROF_NVLINK_L16_TX_BYTES Short = 1072
	// DCGM_FI_PROF_C2C_TX_ALL_BYTES represents C2C (Chip-to-Chip) interface metric
	DCGM_FI_PROF_C2C_TX_ALL_BYTES Short = 1076
	// DCGM_FI_PROF_C2C_TX_DATA_BYTES represents C2C (Chip-to-Chip) interface metric
	DCGM_FI_PROF_C2C_TX_DATA_BYTES Short = 1077
	// DCGM_FI_PROF_C2C_RX_ALL_BYTES represents C2C (Chip-to-Chip) interface metric
	DCGM_FI_PROF_C2C_RX_ALL_BYTES Short = 1078
	// DCGM_FI_PROF_C2C_RX_DATA_BYTES represents C2C (Chip-to-Chip) interface metric
	DCGM_FI_PROF_C2C_RX_DATA_BYTES Short = 1079

	// DCGM_FI_DEV_CPU_UTIL_TOTAL represents the total CPU utilization, total
	DCGM_FI_DEV_CPU_UTIL_TOTAL Short = 1100
	// DCGM_FI_DEV_CPU_UTIL_USER represents the CPU utilization, user
	DCGM_FI_DEV_CPU_UTIL_USER Short = 1101
	// DCGM_FI_DEV_CPU_UTIL_NICE represents the CPU utilization, nice
	DCGM_FI_DEV_CPU_UTIL_NICE Short = 1102
	// DCGM_FI_DEV_CPU_UTIL_SYS represents the CPU utilization, system time
	DCGM_FI_DEV_CPU_UTIL_SYS Short = 1103
	// DCGM_FI_DEV_CPU_UTIL_IRQ represents the CPU utilization, interrupt servicing
	DCGM_FI_DEV_CPU_UTIL_IRQ Short = 1104
	// DCGM_FI_DEV_CPU_TEMP_CURRENT represents the current CPU temperature in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_CURRENT Short = 1110
	// DCGM_FI_DEV_CPU_TEMP_WARNING represents the CPU temperature warning threshold in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_WARNING Short = 1111
	// DCGM_FI_DEV_CPU_TEMP_SHUTDOWN represents the CPU temperature shutdown threshold in degrees Celsius
	DCGM_FI_DEV_CPU_TEMP_SHUTDOWN Short = 1112
	// DCGM_FI_DEV_CPU_CLOCK_CURRENT represents the current CPU clock frequency in MHz
	DCGM_FI_DEV_CPU_CLOCK_CURRENT Short = 1120
	// DCGM_FI_DEV_CPU_POWER_CURRENT represents the current CPU power usage
	DCGM_FI_DEV_CPU_POWER_CURRENT Short = 1130
	// DCGM_FI_DEV_CPU_POWER_LIMIT represents the GPU power limit
	DCGM_FI_DEV_CPU_POWER_LIMIT Short = 1131
	// DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT represents the SoC power utilization
	DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT Short = 1132
	// DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT represents the Module power utilization
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
	// DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER_FLOAT represents BER for symbol errors - decoded float (derived from DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER)
	DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER_FLOAT Short = 1216
	// DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER represents Effective BER for effective errors - raw value
	DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER Short = 1217
	// DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER_FLOAT represents Effective BER for effective errors - decoded float (derived from DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER)
	DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER_FLOAT Short = 1218
	// DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_ERRORS represents Sum of the number of errors in each Nvlink packet
	DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_ERRORS Short = 1219
	// DCGM_FI_DEV_CONNECTX_HEALTH represents a health state of ConnectX
	DCGM_FI_DEV_CONNECTX_HEALTH Short = 1300
	// DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH is the value of an active PCIe link width
	DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH Short = 1301
	// DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED is the value of an active PCIe link speed
	DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED Short = 1302
	// DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH is the value of an expected PCIe link width
	DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH Short = 1303
	// DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED is the value of an expected PCIe link speed
	DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED Short = 1304
	// DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS is the value of a correctable error status
	DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS Short = 1305
	// DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK is the value of a correctable error mask
	DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK Short = 1306
	// DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS is the value of an uncorrectable error status
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS Short = 1307
	// DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK is the value of an uncorrectable error mask
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK Short = 1308
	// DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY is the value of an uncorrectable error severity
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY Short = 1309
	// DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE is the value of a device temperature
	DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE Short = 1310
	// DCGM_FI_DEV_LAST_CONNECTX_FIELD_ID represents the last field ID for ConnectX fields
	DCGM_FI_DEV_LAST_CONNECTX_FIELD_ID Short = 1399
	// DCGM_FI_DEV_C2C_LINK_ERROR_INTR represents C2C Link CRC Error Counter
	DCGM_FI_DEV_C2C_LINK_ERROR_INTR Short = 1400
	// DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY represents C2C Link Replay Error Counter
	DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY Short = 1401
	// DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY_B2B represents C2C Link Back to Back Replay Error Counter
	DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY_B2B Short = 1402
	// DCGM_FI_DEV_C2C_LINK_POWER_STATE represents C2C Link Power state. See NVML_C2C_POWER_STATE_*
	DCGM_FI_DEV_C2C_LINK_POWER_STATE Short = 1403
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_0 represents Count of symbol errors that are corrected - bin 0
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_0 Short = 1404
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_1 represents Count of symbol errors that are corrected - bin 1
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_1 Short = 1405
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_2 represents Count of symbol errors that are corrected - bin 2
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_2 Short = 1406
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_3 represents Count of symbol errors that are corrected - bin 3
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_3 Short = 1407
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_4 represents Count of symbol errors that are corrected - bin 4
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_4 Short = 1408
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_5 represents Count of symbol errors that are corrected - bin 5
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_5 Short = 1409
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_6 represents Count of symbol errors that are corrected - bin 6
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_6 Short = 1410
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_7 represents Count of symbol errors that are corrected - bin 7
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_7 Short = 1411
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_8 represents Count of symbol errors that are corrected - bin 8
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_8 Short = 1412
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_9 represents Count of symbol errors that are corrected - bin 9
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_9 Short = 1413
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_10 represents Count of symbol errors that are corrected - bin 10
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_10 Short = 1414
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_11 represents Count of symbol errors that are corrected - bin 11
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_11 Short = 1415
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_12 represents Count of symbol errors that are corrected - bin 12
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_12 Short = 1416
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_13 represents Count of symbol errors that are corrected - bin 13
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_13 Short = 1417
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_14 represents Count of symbol errors that are corrected - bin 14
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_14 Short = 1418
	// DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_15 represents Count of symbol errors that are corrected - bin 15
	DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_15 Short = 1419
	// DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_POWER_CAP_NS represents Count, in nanoseconds, of slowdown or shutdown in sampling interval.
	DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_POWER_CAP_NS Short = 1420
	// DCGM_FI_DEV_CLOCKS_EVENT_REASON_SYNC_BOOST_NS represents Throttling to not exceed currently set power limits in ns
	DCGM_FI_DEV_CLOCKS_EVENT_REASON_SYNC_BOOST_NS Short = 1421
	// DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_THERM_SLOWDOWN_NS represents Throttling to ensure ((GPU temp < GPU Max Operating Temp) && (Memory Temp < Memory Max Operating Temp)) in ns
	DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_THERM_SLOWDOWN_NS Short = 1422
	// DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_THERM_SLOWDOWN_NS represents Throttling due to temperature being too high (reducing core clocks by a factor of 2 or more) in ns
	DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_THERM_SLOWDOWN_NS Short = 1423
	// DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_POWER_BRAKE_SLOWDOWN_NS represents Throttling due to external power brake assertion trigger (reducing core clocks by a factor of 2 or more) in ns
	DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_POWER_BRAKE_SLOWDOWN_NS Short = 1424
	// DCGM_FI_MAX_FIELDS represents 1 greater than maximum fields above. This is the 1 greater than the maximum field id that could be allocated.
	DCGM_FI_MAX_FIELDS Short = 1425
)

var dcgmFields = map[string]Short{
	// Field types
	"DCGM_FT_BINARY":    Short('b'),
	"DCGM_FT_DOUBLE":    Short('d'),
	"DCGM_FT_INT64":     Short('i'),
	"DCGM_FT_STRING":    Short('s'),
	"DCGM_FT_TIMESTAMP": Short('t'),

	"DCGM_FI_UNKNOWN":                                            DCGM_FI_UNKNOWN,                                            // 0
	"DCGM_FI_DRIVER_VERSION":                                     DCGM_FI_DRIVER_VERSION,                                     // 1
	"DCGM_FI_NVML_VERSION":                                       DCGM_FI_NVML_VERSION,                                       // 2
	"DCGM_FI_PROCESS_NAME":                                       DCGM_FI_PROCESS_NAME,                                       // 3
	"DCGM_FI_DEV_COUNT":                                          DCGM_FI_DEV_COUNT,                                          // 4
	"DCGM_FI_CUDA_DRIVER_VERSION":                                DCGM_FI_CUDA_DRIVER_VERSION,                                // 5
	"DCGM_FI_DEV_NAME":                                           DCGM_FI_DEV_NAME,                                           // 50
	"DCGM_FI_DEV_BRAND":                                          DCGM_FI_DEV_BRAND,                                          // 51
	"DCGM_FI_DEV_NVML_INDEX":                                     DCGM_FI_DEV_NVML_INDEX,                                     // 52
	"DCGM_FI_DEV_SERIAL":                                         DCGM_FI_DEV_SERIAL,                                         // 53
	"DCGM_FI_DEV_UUID":                                           DCGM_FI_DEV_UUID,                                           // 54
	"DCGM_FI_DEV_MINOR_NUMBER":                                   DCGM_FI_DEV_MINOR_NUMBER,                                   // 55
	"DCGM_FI_DEV_OEM_INFOROM_VER":                                DCGM_FI_DEV_OEM_INFOROM_VER,                                // 56
	"DCGM_FI_DEV_PCI_BUSID":                                      DCGM_FI_DEV_PCI_BUSID,                                      // 57
	"DCGM_FI_DEV_PCI_COMBINED_ID":                                DCGM_FI_DEV_PCI_COMBINED_ID,                                // 58
	"DCGM_FI_DEV_PCI_SUBSYS_ID":                                  DCGM_FI_DEV_PCI_SUBSYS_ID,                                  // 59
	"DCGM_FI_GPU_TOPOLOGY_PCI":                                   DCGM_FI_GPU_TOPOLOGY_PCI,                                   // 60
	"DCGM_FI_GPU_TOPOLOGY_NVLINK":                                DCGM_FI_GPU_TOPOLOGY_NVLINK,                                // 61
	"DCGM_FI_GPU_TOPOLOGY_AFFINITY":                              DCGM_FI_GPU_TOPOLOGY_AFFINITY,                              // 62
	"DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY":                        DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY,                        // 63
	"DCGM_FI_DEV_P2P_NVLINK_STATUS":                              DCGM_FI_DEV_P2P_NVLINK_STATUS,                              // 64
	"DCGM_FI_DEV_COMPUTE_MODE":                                   DCGM_FI_DEV_COMPUTE_MODE,                                   // 65
	"DCGM_FI_DEV_PERSISTENCE_MODE":                               DCGM_FI_DEV_PERSISTENCE_MODE,                               // 66
	"DCGM_FI_DEV_MIG_MODE":                                       DCGM_FI_DEV_MIG_MODE,                                       // 67
	"DCGM_FI_DEV_CUDA_VISIBLE_DEVICES_STR":                       DCGM_FI_DEV_CUDA_VISIBLE_DEVICES_STR,                       // 68
	"DCGM_FI_DEV_MIG_MAX_SLICES":                                 DCGM_FI_DEV_MIG_MAX_SLICES,                                 // 69
	"DCGM_FI_DEV_CPU_AFFINITY_0":                                 DCGM_FI_DEV_CPU_AFFINITY_0,                                 // 70
	"DCGM_FI_DEV_CPU_AFFINITY_1":                                 DCGM_FI_DEV_CPU_AFFINITY_1,                                 // 71
	"DCGM_FI_DEV_CPU_AFFINITY_2":                                 DCGM_FI_DEV_CPU_AFFINITY_2,                                 // 72
	"DCGM_FI_DEV_CPU_AFFINITY_3":                                 DCGM_FI_DEV_CPU_AFFINITY_3,                                 // 73
	"DCGM_FI_DEV_CC_MODE":                                        DCGM_FI_DEV_CC_MODE,                                        // 74
	"DCGM_FI_DEV_MIG_ATTRIBUTES":                                 DCGM_FI_DEV_MIG_ATTRIBUTES,                                 // 75
	"DCGM_FI_DEV_MIG_GI_INFO":                                    DCGM_FI_DEV_MIG_GI_INFO,                                    // 76
	"DCGM_FI_DEV_MIG_CI_INFO":                                    DCGM_FI_DEV_MIG_CI_INFO,                                    // 77
	"DCGM_FI_DEV_ECC_INFOROM_VER":                                DCGM_FI_DEV_ECC_INFOROM_VER,                                // 80
	"DCGM_FI_DEV_POWER_INFOROM_VER":                              DCGM_FI_DEV_POWER_INFOROM_VER,                              // 81
	"DCGM_FI_DEV_INFOROM_IMAGE_VER":                              DCGM_FI_DEV_INFOROM_IMAGE_VER,                              // 82
	"DCGM_FI_DEV_INFOROM_CONFIG_CHECK":                           DCGM_FI_DEV_INFOROM_CONFIG_CHECK,                           // 83
	"DCGM_FI_DEV_INFOROM_CONFIG_VALID":                           DCGM_FI_DEV_INFOROM_CONFIG_VALID,                           // 84
	"DCGM_FI_DEV_VBIOS_VERSION":                                  DCGM_FI_DEV_VBIOS_VERSION,                                  // 85
	"DCGM_FI_DEV_MEM_AFFINITY_0":                                 DCGM_FI_DEV_MEM_AFFINITY_0,                                 // 86
	"DCGM_FI_DEV_MEM_AFFINITY_1":                                 DCGM_FI_DEV_MEM_AFFINITY_1,                                 // 87
	"DCGM_FI_DEV_MEM_AFFINITY_2":                                 DCGM_FI_DEV_MEM_AFFINITY_2,                                 // 88
	"DCGM_FI_DEV_MEM_AFFINITY_3":                                 DCGM_FI_DEV_MEM_AFFINITY_3,                                 // 89
	"DCGM_FI_DEV_BAR1_TOTAL":                                     DCGM_FI_DEV_BAR1_TOTAL,                                     // 90
	"DCGM_FI_SYNC_BOOST":                                         DCGM_FI_SYNC_BOOST,                                         // 91
	"DCGM_FI_DEV_BAR1_USED":                                      DCGM_FI_DEV_BAR1_USED,                                      // 92
	"DCGM_FI_DEV_BAR1_FREE":                                      DCGM_FI_DEV_BAR1_FREE,                                      // 93
	"DCGM_FI_DEV_GPM_SUPPORT":                                    DCGM_FI_DEV_GPM_SUPPORT,                                    // 94
	"DCGM_FI_DEV_SM_CLOCK":                                       DCGM_FI_DEV_SM_CLOCK,                                       // 100
	"DCGM_FI_DEV_MEM_CLOCK":                                      DCGM_FI_DEV_MEM_CLOCK,                                      // 101
	"DCGM_FI_DEV_VIDEO_CLOCK":                                    DCGM_FI_DEV_VIDEO_CLOCK,                                    // 102
	"DCGM_FI_DEV_APP_SM_CLOCK":                                   DCGM_FI_DEV_APP_SM_CLOCK,                                   // 110
	"DCGM_FI_DEV_APP_MEM_CLOCK":                                  DCGM_FI_DEV_APP_MEM_CLOCK,                                  // 111
	"DCGM_FI_DEV_CLOCKS_EVENT_REASONS":                           DCGM_FI_DEV_CLOCKS_EVENT_REASONS,                           // 112
	"DCGM_FI_DEV_CLOCK_THROTTLE_REASONS":                         DCGM_FI_DEV_CLOCK_THROTTLE_REASONS,                         // 112
	"DCGM_FI_DEV_MAX_SM_CLOCK":                                   DCGM_FI_DEV_MAX_SM_CLOCK,                                   // 113
	"DCGM_FI_DEV_MAX_MEM_CLOCK":                                  DCGM_FI_DEV_MAX_MEM_CLOCK,                                  // 114
	"DCGM_FI_DEV_MAX_VIDEO_CLOCK":                                DCGM_FI_DEV_MAX_VIDEO_CLOCK,                                // 115
	"DCGM_FI_DEV_AUTOBOOST":                                      DCGM_FI_DEV_AUTOBOOST,                                      // 120
	"DCGM_FI_DEV_SUPPORTED_CLOCKS":                               DCGM_FI_DEV_SUPPORTED_CLOCKS,                               // 130
	"DCGM_FI_DEV_MEMORY_TEMP":                                    DCGM_FI_DEV_MEMORY_TEMP,                                    // 140
	"DCGM_FI_DEV_GPU_TEMP":                                       DCGM_FI_DEV_GPU_TEMP,                                       // 150
	"DCGM_FI_DEV_MEM_MAX_OP_TEMP":                                DCGM_FI_DEV_MEM_MAX_OP_TEMP,                                // 151
	"DCGM_FI_DEV_GPU_MAX_OP_TEMP":                                DCGM_FI_DEV_GPU_MAX_OP_TEMP,                                // 152
	"DCGM_FI_DEV_GPU_TEMP_LIMIT":                                 DCGM_FI_DEV_GPU_TEMP_LIMIT,                                 // 153
	"DCGM_FI_DEV_POWER_USAGE":                                    DCGM_FI_DEV_POWER_USAGE,                                    // 155
	"DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION":                       DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION,                       // 156
	"DCGM_FI_DEV_POWER_USAGE_INSTANT":                            DCGM_FI_DEV_POWER_USAGE_INSTANT,                            // 157
	"DCGM_FI_DEV_SLOWDOWN_TEMP":                                  DCGM_FI_DEV_SLOWDOWN_TEMP,                                  // 158
	"DCGM_FI_DEV_SHUTDOWN_TEMP":                                  DCGM_FI_DEV_SHUTDOWN_TEMP,                                  // 159
	"DCGM_FI_DEV_POWER_MGMT_LIMIT":                               DCGM_FI_DEV_POWER_MGMT_LIMIT,                               // 160
	"DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN":                           DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN,                           // 161
	"DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX":                           DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX,                           // 162
	"DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF":                           DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF,                           // 163
	"DCGM_FI_DEV_ENFORCED_POWER_LIMIT":                           DCGM_FI_DEV_ENFORCED_POWER_LIMIT,                           // 164
	"DCGM_FI_DEV_REQUESTED_POWER_PROFILE_MASK":                   DCGM_FI_DEV_REQUESTED_POWER_PROFILE_MASK,                   // 165
	"DCGM_FI_DEV_ENFORCED_POWER_PROFILE_MASK":                    DCGM_FI_DEV_ENFORCED_POWER_PROFILE_MASK,                    // 166
	"DCGM_FI_DEV_VALID_POWER_PROFILE_MASK":                       DCGM_FI_DEV_VALID_POWER_PROFILE_MASK,                       // 167
	"DCGM_FI_DEV_FABRIC_MANAGER_STATUS":                          DCGM_FI_DEV_FABRIC_MANAGER_STATUS,                          // 170
	"DCGM_FI_DEV_FABRIC_MANAGER_ERROR_CODE":                      DCGM_FI_DEV_FABRIC_MANAGER_ERROR_CODE,                      // 171
	"DCGM_FI_DEV_FABRIC_CLUSTER_UUID":                            DCGM_FI_DEV_FABRIC_CLUSTER_UUID,                            // 172
	"DCGM_FI_DEV_FABRIC_CLIQUE_ID":                               DCGM_FI_DEV_FABRIC_CLIQUE_ID,                               // 173
	"DCGM_FI_DEV_PSTATE":                                         DCGM_FI_DEV_PSTATE,                                         // 190
	"DCGM_FI_DEV_FAN_SPEED":                                      DCGM_FI_DEV_FAN_SPEED,                                      // 191
	"DCGM_FI_DEV_PCIE_TX_THROUGHPUT":                             DCGM_FI_DEV_PCIE_TX_THROUGHPUT,                             // 200
	"DCGM_FI_DEV_PCIE_RX_THROUGHPUT":                             DCGM_FI_DEV_PCIE_RX_THROUGHPUT,                             // 201
	"DCGM_FI_DEV_PCIE_REPLAY_COUNTER":                            DCGM_FI_DEV_PCIE_REPLAY_COUNTER,                            // 202
	"DCGM_FI_DEV_GPU_UTIL":                                       DCGM_FI_DEV_GPU_UTIL,                                       // 203
	"DCGM_FI_DEV_MEM_COPY_UTIL":                                  DCGM_FI_DEV_MEM_COPY_UTIL,                                  // 204
	"DCGM_FI_DEV_ACCOUNTING_DATA":                                DCGM_FI_DEV_ACCOUNTING_DATA,                                // 205
	"DCGM_FI_DEV_ENC_UTIL":                                       DCGM_FI_DEV_ENC_UTIL,                                       // 206
	"DCGM_FI_DEV_DEC_UTIL":                                       DCGM_FI_DEV_DEC_UTIL,                                       // 207
	"DCGM_FI_DEV_XID_ERRORS":                                     DCGM_FI_DEV_XID_ERRORS,                                     // 230
	"DCGM_FI_DEV_PCIE_MAX_LINK_GEN":                              DCGM_FI_DEV_PCIE_MAX_LINK_GEN,                              // 235
	"DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH":                            DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH,                            // 236
	"DCGM_FI_DEV_PCIE_LINK_GEN":                                  DCGM_FI_DEV_PCIE_LINK_GEN,                                  // 237
	"DCGM_FI_DEV_PCIE_LINK_WIDTH":                                DCGM_FI_DEV_PCIE_LINK_WIDTH,                                // 238
	"DCGM_FI_DEV_POWER_VIOLATION":                                DCGM_FI_DEV_POWER_VIOLATION,                                // 240
	"DCGM_FI_DEV_THERMAL_VIOLATION":                              DCGM_FI_DEV_THERMAL_VIOLATION,                              // 241
	"DCGM_FI_DEV_SYNC_BOOST_VIOLATION":                           DCGM_FI_DEV_SYNC_BOOST_VIOLATION,                           // 242
	"DCGM_FI_DEV_BOARD_LIMIT_VIOLATION":                          DCGM_FI_DEV_BOARD_LIMIT_VIOLATION,                          // 243
	"DCGM_FI_DEV_LOW_UTIL_VIOLATION":                             DCGM_FI_DEV_LOW_UTIL_VIOLATION,                             // 244
	"DCGM_FI_DEV_RELIABILITY_VIOLATION":                          DCGM_FI_DEV_RELIABILITY_VIOLATION,                          // 245
	"DCGM_FI_DEV_TOTAL_APP_CLOCKS_VIOLATION":                     DCGM_FI_DEV_TOTAL_APP_CLOCKS_VIOLATION,                     // 246
	"DCGM_FI_DEV_TOTAL_BASE_CLOCKS_VIOLATION":                    DCGM_FI_DEV_TOTAL_BASE_CLOCKS_VIOLATION,                    // 247
	"DCGM_FI_DEV_FB_TOTAL":                                       DCGM_FI_DEV_FB_TOTAL,                                       // 250
	"DCGM_FI_DEV_FB_FREE":                                        DCGM_FI_DEV_FB_FREE,                                        // 251
	"DCGM_FI_DEV_FB_USED":                                        DCGM_FI_DEV_FB_USED,                                        // 252
	"DCGM_FI_DEV_FB_RESERVED":                                    DCGM_FI_DEV_FB_RESERVED,                                    // 253
	"DCGM_FI_DEV_FB_USED_PERCENT":                                DCGM_FI_DEV_FB_USED_PERCENT,                                // 254
	"DCGM_FI_DEV_C2C_LINK_COUNT":                                 DCGM_FI_DEV_C2C_LINK_COUNT,                                 // 285
	"DCGM_FI_DEV_C2C_LINK_STATUS":                                DCGM_FI_DEV_C2C_LINK_STATUS,                                // 286
	"DCGM_FI_DEV_C2C_MAX_BANDWIDTH":                              DCGM_FI_DEV_C2C_MAX_BANDWIDTH,                              // 287
	"DCGM_FI_DEV_ECC_CURRENT":                                    DCGM_FI_DEV_ECC_CURRENT,                                    // 300
	"DCGM_FI_DEV_ECC_PENDING":                                    DCGM_FI_DEV_ECC_PENDING,                                    // 301
	"DCGM_FI_DEV_ECC_SBE_VOL_TOTAL":                              DCGM_FI_DEV_ECC_SBE_VOL_TOTAL,                              // 310
	"DCGM_FI_DEV_ECC_DBE_VOL_TOTAL":                              DCGM_FI_DEV_ECC_DBE_VOL_TOTAL,                              // 311
	"DCGM_FI_DEV_ECC_SBE_AGG_TOTAL":                              DCGM_FI_DEV_ECC_SBE_AGG_TOTAL,                              // 312
	"DCGM_FI_DEV_ECC_DBE_AGG_TOTAL":                              DCGM_FI_DEV_ECC_DBE_AGG_TOTAL,                              // 313
	"DCGM_FI_DEV_ECC_SBE_VOL_L1":                                 DCGM_FI_DEV_ECC_SBE_VOL_L1,                                 // 314
	"DCGM_FI_DEV_ECC_DBE_VOL_L1":                                 DCGM_FI_DEV_ECC_DBE_VOL_L1,                                 // 315
	"DCGM_FI_DEV_ECC_SBE_VOL_L2":                                 DCGM_FI_DEV_ECC_SBE_VOL_L2,                                 // 316
	"DCGM_FI_DEV_ECC_DBE_VOL_L2":                                 DCGM_FI_DEV_ECC_DBE_VOL_L2,                                 // 317
	"DCGM_FI_DEV_ECC_SBE_VOL_DEV":                                DCGM_FI_DEV_ECC_SBE_VOL_DEV,                                // 318
	"DCGM_FI_DEV_ECC_DBE_VOL_DEV":                                DCGM_FI_DEV_ECC_DBE_VOL_DEV,                                // 319
	"DCGM_FI_DEV_ECC_SBE_VOL_REG":                                DCGM_FI_DEV_ECC_SBE_VOL_REG,                                // 320
	"DCGM_FI_DEV_ECC_DBE_VOL_REG":                                DCGM_FI_DEV_ECC_DBE_VOL_REG,                                // 321
	"DCGM_FI_DEV_ECC_SBE_VOL_TEX":                                DCGM_FI_DEV_ECC_SBE_VOL_TEX,                                // 322
	"DCGM_FI_DEV_ECC_DBE_VOL_TEX":                                DCGM_FI_DEV_ECC_DBE_VOL_TEX,                                // 323
	"DCGM_FI_DEV_ECC_SBE_AGG_L1":                                 DCGM_FI_DEV_ECC_SBE_AGG_L1,                                 // 324
	"DCGM_FI_DEV_ECC_DBE_AGG_L1":                                 DCGM_FI_DEV_ECC_DBE_AGG_L1,                                 // 325
	"DCGM_FI_DEV_ECC_SBE_AGG_L2":                                 DCGM_FI_DEV_ECC_SBE_AGG_L2,                                 // 326
	"DCGM_FI_DEV_ECC_DBE_AGG_L2":                                 DCGM_FI_DEV_ECC_DBE_AGG_L2,                                 // 327
	"DCGM_FI_DEV_ECC_SBE_AGG_DEV":                                DCGM_FI_DEV_ECC_SBE_AGG_DEV,                                // 328
	"DCGM_FI_DEV_ECC_DBE_AGG_DEV":                                DCGM_FI_DEV_ECC_DBE_AGG_DEV,                                // 329
	"DCGM_FI_DEV_ECC_SBE_AGG_REG":                                DCGM_FI_DEV_ECC_SBE_AGG_REG,                                // 330
	"DCGM_FI_DEV_ECC_DBE_AGG_REG":                                DCGM_FI_DEV_ECC_DBE_AGG_REG,                                // 331
	"DCGM_FI_DEV_ECC_SBE_AGG_TEX":                                DCGM_FI_DEV_ECC_SBE_AGG_TEX,                                // 332
	"DCGM_FI_DEV_ECC_DBE_AGG_TEX":                                DCGM_FI_DEV_ECC_DBE_AGG_TEX,                                // 333
	"DCGM_FI_DEV_ECC_SBE_VOL_SHM":                                DCGM_FI_DEV_ECC_SBE_VOL_SHM,                                // 334
	"DCGM_FI_DEV_ECC_DBE_VOL_SHM":                                DCGM_FI_DEV_ECC_DBE_VOL_SHM,                                // 335
	"DCGM_FI_DEV_ECC_SBE_VOL_CBU":                                DCGM_FI_DEV_ECC_SBE_VOL_CBU,                                // 336
	"DCGM_FI_DEV_ECC_DBE_VOL_CBU":                                DCGM_FI_DEV_ECC_DBE_VOL_CBU,                                // 337
	"DCGM_FI_DEV_ECC_SBE_AGG_SHM":                                DCGM_FI_DEV_ECC_SBE_AGG_SHM,                                // 338
	"DCGM_FI_DEV_ECC_DBE_AGG_SHM":                                DCGM_FI_DEV_ECC_DBE_AGG_SHM,                                // 339
	"DCGM_FI_DEV_ECC_SBE_AGG_CBU":                                DCGM_FI_DEV_ECC_SBE_AGG_CBU,                                // 340
	"DCGM_FI_DEV_ECC_DBE_AGG_CBU":                                DCGM_FI_DEV_ECC_DBE_AGG_CBU,                                // 341
	"DCGM_FI_DEV_ECC_SBE_VOL_SRM":                                DCGM_FI_DEV_ECC_SBE_VOL_SRM,                                // 342
	"DCGM_FI_DEV_ECC_DBE_VOL_SRM":                                DCGM_FI_DEV_ECC_DBE_VOL_SRM,                                // 343
	"DCGM_FI_DEV_ECC_SBE_AGG_SRM":                                DCGM_FI_DEV_ECC_SBE_AGG_SRM,                                // 344
	"DCGM_FI_DEV_ECC_DBE_AGG_SRM":                                DCGM_FI_DEV_ECC_DBE_AGG_SRM,                                // 345
	"DCGM_FI_DEV_THRESHOLD_SRM":                                  DCGM_FI_DEV_THRESHOLD_SRM,                                  // 346
	"DCGM_FI_DEV_DIAG_MEMORY_RESULT":                             DCGM_FI_DEV_DIAG_MEMORY_RESULT,                             // 350
	"DCGM_FI_DEV_DIAG_DIAGNOSTIC_RESULT":                         DCGM_FI_DEV_DIAG_DIAGNOSTIC_RESULT,                         // 351
	"DCGM_FI_DEV_DIAG_PCIE_RESULT":                               DCGM_FI_DEV_DIAG_PCIE_RESULT,                               // 352
	"DCGM_FI_DEV_DIAG_TARGETED_STRESS_RESULT":                    DCGM_FI_DEV_DIAG_TARGETED_STRESS_RESULT,                    // 353
	"DCGM_FI_DEV_DIAG_TARGETED_POWER_RESULT":                     DCGM_FI_DEV_DIAG_TARGETED_POWER_RESULT,                     // 354
	"DCGM_FI_DEV_DIAG_MEMORY_BANDWIDTH_RESULT":                   DCGM_FI_DEV_DIAG_MEMORY_BANDWIDTH_RESULT,                   // 355
	"DCGM_FI_DEV_DIAG_MEMTEST_RESULT":                            DCGM_FI_DEV_DIAG_MEMTEST_RESULT,                            // 356
	"DCGM_FI_DEV_DIAG_PULSE_TEST_RESULT":                         DCGM_FI_DEV_DIAG_PULSE_TEST_RESULT,                         // 357
	"DCGM_FI_DEV_DIAG_EUD_RESULT":                                DCGM_FI_DEV_DIAG_EUD_RESULT,                                // 358
	"DCGM_FI_DEV_DIAG_CPU_EUD_RESULT":                            DCGM_FI_DEV_DIAG_CPU_EUD_RESULT,                            // 359
	"DCGM_FI_DEV_DIAG_SOFTWARE_RESULT":                           DCGM_FI_DEV_DIAG_SOFTWARE_RESULT,                           // 360
	"DCGM_FI_DEV_DIAG_NVBANDWIDTH_RESULT":                        DCGM_FI_DEV_DIAG_NVBANDWIDTH_RESULT,                        // 361
	"DCGM_FI_DEV_DIAG_STATUS":                                    DCGM_FI_DEV_DIAG_STATUS,                                    // 362
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_MAX":                     DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_MAX,                     // 385
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_HIGH":                    DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_HIGH,                    // 386
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_PARTIAL":                 DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_PARTIAL,                 // 387
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_LOW":                     DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_LOW,                     // 388
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_NONE":                    DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_NONE,                    // 389
	"DCGM_FI_DEV_RETIRED_SBE":                                    DCGM_FI_DEV_RETIRED_SBE,                                    // 390
	"DCGM_FI_DEV_RETIRED_DBE":                                    DCGM_FI_DEV_RETIRED_DBE,                                    // 391
	"DCGM_FI_DEV_RETIRED_PENDING":                                DCGM_FI_DEV_RETIRED_PENDING,                                // 392
	"DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS":                    DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS,                    // 393
	"DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS":                      DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS,                      // 394
	"DCGM_FI_DEV_ROW_REMAP_FAILURE":                              DCGM_FI_DEV_ROW_REMAP_FAILURE,                              // 395
	"DCGM_FI_DEV_ROW_REMAP_PENDING":                              DCGM_FI_DEV_ROW_REMAP_PENDING,                              // 396
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L0":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L0,                 // 400
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L1":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L1,                 // 401
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L2":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L2,                 // 402
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L3":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L3,                 // 403
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L4":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L4,                 // 404
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L5":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L5,                 // 405
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL":              DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL,              // 409
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L0":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L0,                 // 410
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L1":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L1,                 // 411
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L2":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L2,                 // 412
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L3":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L3,                 // 413
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L4":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L4,                 // 414
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L5":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L5,                 // 415
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_TOTAL":              DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_TOTAL,              // 419
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L0":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L0,                   // 420
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L1":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L1,                   // 421
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L2":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L2,                   // 422
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L3":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L3,                   // 423
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L4":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L4,                   // 424
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L5":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L5,                   // 425
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_TOTAL":                DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_TOTAL,                // 429
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L0":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L0,                 // 430
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L1":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L1,                 // 431
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L2":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L2,                 // 432
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L3":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L3,                 // 433
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L4":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L4,                 // 434
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L5":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L5,                 // 435
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_TOTAL":              DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_TOTAL,              // 439
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L0":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L0,                            // 440
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L1":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L1,                            // 441
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L2":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L2,                            // 442
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L3":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L3,                            // 443
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L4":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L4,                            // 444
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L5":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L5,                            // 445
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL":                         DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL,                         // 449
	"DCGM_FI_DEV_GPU_NVLINK_ERRORS":                              DCGM_FI_DEV_GPU_NVLINK_ERRORS,                              // 450
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L6":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L6,                 // 451
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L7":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L7,                 // 452
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L8":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L8,                 // 453
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L9":                 DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L9,                 // 454
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L10":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L10,                // 455
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L11":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L11,                // 456
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L6":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L6,                 // 457
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L7":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L7,                 // 458
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L8":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L8,                 // 459
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L9":                 DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L9,                 // 460
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L10":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L10,                // 461
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L11":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L11,                // 462
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L6":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L6,                   // 463
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L7":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L7,                   // 464
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L8":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L8,                   // 465
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L9":                   DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L9,                   // 466
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L10":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L10,                  // 467
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L11":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L11,                  // 468
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L6":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L6,                 // 469
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L7":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L7,                 // 470
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L8":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L8,                 // 471
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L9":                 DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L9,                 // 472
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L10":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L10,                // 473
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L11":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L11,                // 474
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L6":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L6,                            // 475
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L7":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L7,                            // 476
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L8":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L8,                            // 477
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L9":                            DCGM_FI_DEV_NVLINK_BANDWIDTH_L9,                            // 478
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L10":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L10,                           // 479
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L11":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L11,                           // 480
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L12":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L12,                // 406
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13,                // 407
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14,                // 408
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L15":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L15,                // 481
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L16":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L16,                // 482
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L17":                DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L17,                // 483
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L12":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L12,                // 416
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L13":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L13,                // 417
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L14":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L14,                // 418
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L15":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L15,                // 484
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L16":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L16,                // 485
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L17":                DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L17,                // 486
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L12":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L12,                  // 426
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L13":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L13,                  // 427
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L14":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L14,                  // 428
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L15":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L15,                  // 487
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L16":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L16,                  // 488
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L17":                  DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L17,                  // 489
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L12":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L12,                // 436
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L13":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L13,                // 437
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L14":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L14,                // 438
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L15":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L15,                // 491
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L16":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L16,                // 492
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L17":                DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L17,                // 493
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L12":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L12,                           // 446
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L13":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L13,                           // 447
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L14":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L14,                           // 448
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L15":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L15,                           // 494
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L16":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L16,                           // 495
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L17":                           DCGM_FI_DEV_NVLINK_BANDWIDTH_L17,                           // 496
	"DCGM_FI_DEV_NVLINK_ERROR_DL_CRC":                            DCGM_FI_DEV_NVLINK_ERROR_DL_CRC,                            // 497
	"DCGM_FI_DEV_NVLINK_ERROR_DL_RECOVERY":                       DCGM_FI_DEV_NVLINK_ERROR_DL_RECOVERY,                       // 498
	"DCGM_FI_DEV_NVLINK_ERROR_DL_REPLAY":                         DCGM_FI_DEV_NVLINK_ERROR_DL_REPLAY,                         // 499
	"DCGM_FI_DEV_VIRTUAL_MODE":                                   DCGM_FI_DEV_VIRTUAL_MODE,                                   // 500
	"DCGM_FI_DEV_SUPPORTED_TYPE_INFO":                            DCGM_FI_DEV_SUPPORTED_TYPE_INFO,                            // 501
	"DCGM_FI_DEV_CREATABLE_VGPU_TYPE_IDS":                        DCGM_FI_DEV_CREATABLE_VGPU_TYPE_IDS,                        // 502
	"DCGM_FI_DEV_VGPU_INSTANCE_IDS":                              DCGM_FI_DEV_VGPU_INSTANCE_IDS,                              // 503
	"DCGM_FI_DEV_VGPU_UTILIZATIONS":                              DCGM_FI_DEV_VGPU_UTILIZATIONS,                              // 504
	"DCGM_FI_DEV_VGPU_PER_PROCESS_UTILIZATION":                   DCGM_FI_DEV_VGPU_PER_PROCESS_UTILIZATION,                   // 505
	"DCGM_FI_DEV_ENC_STATS":                                      DCGM_FI_DEV_ENC_STATS,                                      // 506
	"DCGM_FI_DEV_FBC_STATS":                                      DCGM_FI_DEV_FBC_STATS,                                      // 507
	"DCGM_FI_DEV_FBC_SESSIONS_INFO":                              DCGM_FI_DEV_FBC_SESSIONS_INFO,                              // 508
	"DCGM_FI_DEV_SUPPORTED_VGPU_TYPE_IDS":                        DCGM_FI_DEV_SUPPORTED_VGPU_TYPE_IDS,                        // 509
	"DCGM_FI_DEV_VGPU_TYPE_INFO":                                 DCGM_FI_DEV_VGPU_TYPE_INFO,                                 // 510
	"DCGM_FI_DEV_VGPU_TYPE_NAME":                                 DCGM_FI_DEV_VGPU_TYPE_NAME,                                 // 511
	"DCGM_FI_DEV_VGPU_TYPE_CLASS":                                DCGM_FI_DEV_VGPU_TYPE_CLASS,                                // 512
	"DCGM_FI_DEV_VGPU_TYPE_LICENSE":                              DCGM_FI_DEV_VGPU_TYPE_LICENSE,                              // 513
	"DCGM_FI_DEV_VGPU_VM_ID":                                     DCGM_FI_DEV_VGPU_VM_ID,                                     // 520
	"DCGM_FI_DEV_VGPU_VM_NAME":                                   DCGM_FI_DEV_VGPU_VM_NAME,                                   // 521
	"DCGM_FI_DEV_VGPU_TYPE":                                      DCGM_FI_DEV_VGPU_TYPE,                                      // 522
	"DCGM_FI_DEV_VGPU_UUID":                                      DCGM_FI_DEV_VGPU_UUID,                                      // 523
	"DCGM_FI_DEV_VGPU_DRIVER_VERSION":                            DCGM_FI_DEV_VGPU_DRIVER_VERSION,                            // 524
	"DCGM_FI_DEV_VGPU_MEMORY_USAGE":                              DCGM_FI_DEV_VGPU_MEMORY_USAGE,                              // 525
	"DCGM_FI_DEV_VGPU_LICENSE_STATUS":                            DCGM_FI_DEV_VGPU_LICENSE_STATUS,                            // 526
	"DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT":                          DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT,                          // 527
	"DCGM_FI_DEV_VGPU_ENC_STATS":                                 DCGM_FI_DEV_VGPU_ENC_STATS,                                 // 528
	"DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO":                         DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO,                         // 529
	"DCGM_FI_DEV_VGPU_FBC_STATS":                                 DCGM_FI_DEV_VGPU_FBC_STATS,                                 // 530
	"DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO":                         DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO,                         // 531
	"DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATE":                    DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATE,                    // 532
	"DCGM_FI_DEV_VGPU_PCI_ID":                                    DCGM_FI_DEV_VGPU_PCI_ID,                                    // 533
	"DCGM_FI_DEV_VGPU_VM_GPU_INSTANCE_ID":                        DCGM_FI_DEV_VGPU_VM_GPU_INSTANCE_ID,                        // 534
	"DCGM_FI_FIRST_VGPU_FIELD_ID":                                DCGM_FI_FIRST_VGPU_FIELD_ID,                                // 520
	"DCGM_FI_LAST_VGPU_FIELD_ID":                                 DCGM_FI_LAST_VGPU_FIELD_ID,                                 // 570
	"DCGM_FI_DEV_PLATFORM_INFINIBAND_GUID":                       DCGM_FI_DEV_PLATFORM_INFINIBAND_GUID,                       // 571
	"DCGM_FI_DEV_PLATFORM_CHASSIS_SERIAL_NUMBER":                 DCGM_FI_DEV_PLATFORM_CHASSIS_SERIAL_NUMBER,                 // 572
	"DCGM_FI_DEV_PLATFORM_CHASSIS_SLOT_NUMBER":                   DCGM_FI_DEV_PLATFORM_CHASSIS_SLOT_NUMBER,                   // 573
	"DCGM_FI_DEV_PLATFORM_TRAY_INDEX":                            DCGM_FI_DEV_PLATFORM_TRAY_INDEX,                            // 574
	"DCGM_FI_DEV_PLATFORM_HOST_ID":                               DCGM_FI_DEV_PLATFORM_HOST_ID,                               // 575
	"DCGM_FI_DEV_PLATFORM_PEER_TYPE":                             DCGM_FI_DEV_PLATFORM_PEER_TYPE,                             // 576
	"DCGM_FI_DEV_PLATFORM_MODULE_ID":                             DCGM_FI_DEV_PLATFORM_MODULE_ID,                             // 577
	"DCGM_FI_FIRST_NVSWITCH_FIELD_ID":                            DCGM_FI_FIRST_NVSWITCH_FIELD_ID,                            // 700
	"DCGM_FI_DEV_NVSWITCH_VOLTAGE_MVOLT":                         DCGM_FI_DEV_NVSWITCH_VOLTAGE_MVOLT,                         // 701
	"DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ":                          DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ,                          // 702
	"DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_REV":                      DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_REV,                      // 703
	"DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_DVDD":                     DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_DVDD,                     // 704
	"DCGM_FI_DEV_NVSWITCH_POWER_VDD":                             DCGM_FI_DEV_NVSWITCH_POWER_VDD,                             // 705
	"DCGM_FI_DEV_NVSWITCH_POWER_DVDD":                            DCGM_FI_DEV_NVSWITCH_POWER_DVDD,                            // 706
	"DCGM_FI_DEV_NVSWITCH_POWER_HVDD":                            DCGM_FI_DEV_NVSWITCH_POWER_HVDD,                            // 707
	"DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX":                    DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX,                    // 780
	"DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX":                    DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX,                    // 781
	"DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS":                     DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS,                     // 782
	"DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS":                 DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS,                 // 783
	"DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS":                    DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS,                    // 784
	"DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS":                  DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS,                  // 785
	"DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS":                      DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS,                      // 786
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS":                       DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS,                       // 787
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS":                       DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS,                       // 788
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0":                  DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0,                  // 789
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1":                  DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1,                  // 790
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2":                  DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2,                  // 791
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3":                  DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3,                  // 792
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0":               DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0,               // 793
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1":               DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1,               // 794
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2":               DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2,               // 795
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3":               DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3,               // 796
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0":                 DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0,                 // 797
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1":                 DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1,                 // 798
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2":                 DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2,                 // 799
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3":                 DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3,                 // 800
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0,                // 801
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1,                // 802
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2,                // 803
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3,                // 804
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0,                // 805
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1,                // 806
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2,                // 807
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3":                DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3,                // 808
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0,                 // 809
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1,                 // 810
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2,                 // 811
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3,                 // 812
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0,                 // 813
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1,                 // 814
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2,                 // 815
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3,                 // 816
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4,                 // 817
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5,                 // 818
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6,                 // 819
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7":                 DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7,                 // 820
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4,                 // 821
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5,                 // 822
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6,                 // 823
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE7":                 DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE7,                 // 824
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L0":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L0,                         // 825
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L1":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L1,                         // 826
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L2":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L2,                         // 827
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L3":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L3,                         // 828
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L4":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L4,                         // 829
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L5":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L5,                         // 830
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L6":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L6,                         // 831
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L7":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L7,                         // 832
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L8":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L8,                         // 833
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L9":                         DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L9,                         // 834
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L10":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L10,                        // 835
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L11":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L11,                        // 836
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L12":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L12,                        // 837
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L13":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L13,                        // 838
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L14":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L14,                        // 839
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L15":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L15,                        // 840
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L16":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L16,                        // 841
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L17":                        DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L17,                        // 842
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL":                      DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL,                      // 843
	"DCGM_FI_DEV_NVSWITCH_FATAL_ERRORS":                          DCGM_FI_DEV_NVSWITCH_FATAL_ERRORS,                          // 856
	"DCGM_FI_DEV_NVSWITCH_NON_FATAL_ERRORS":                      DCGM_FI_DEV_NVSWITCH_NON_FATAL_ERRORS,                      // 857
	"DCGM_FI_DEV_NVSWITCH_TEMPERATURE_CURRENT":                   DCGM_FI_DEV_NVSWITCH_TEMPERATURE_CURRENT,                   // 858
	"DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SLOWDOWN":            DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SLOWDOWN,            // 859
	"DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SHUTDOWN":            DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SHUTDOWN,            // 860
	"DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX":                         DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX,                         // 861
	"DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX":                         DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX,                         // 862
	"DCGM_FI_DEV_NVSWITCH_PHYS_ID":                               DCGM_FI_DEV_NVSWITCH_PHYS_ID,                               // 863
	"DCGM_FI_DEV_NVSWITCH_RESET_REQUIRED":                        DCGM_FI_DEV_NVSWITCH_RESET_REQUIRED,                        // 864
	"DCGM_FI_DEV_NVSWITCH_LINK_ID":                               DCGM_FI_DEV_NVSWITCH_LINK_ID,                               // 865
	"DCGM_FI_DEV_NVSWITCH_PCIE_DOMAIN":                           DCGM_FI_DEV_NVSWITCH_PCIE_DOMAIN,                           // 866
	"DCGM_FI_DEV_NVSWITCH_PCIE_BUS":                              DCGM_FI_DEV_NVSWITCH_PCIE_BUS,                              // 867
	"DCGM_FI_DEV_NVSWITCH_PCIE_DEVICE":                           DCGM_FI_DEV_NVSWITCH_PCIE_DEVICE,                           // 868
	"DCGM_FI_DEV_NVSWITCH_PCIE_FUNCTION":                         DCGM_FI_DEV_NVSWITCH_PCIE_FUNCTION,                         // 869
	"DCGM_FI_DEV_NVSWITCH_LINK_STATUS":                           DCGM_FI_DEV_NVSWITCH_LINK_STATUS,                           // 870
	"DCGM_FI_DEV_NVSWITCH_LINK_TYPE":                             DCGM_FI_DEV_NVSWITCH_LINK_TYPE,                             // 871
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN":               DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN,               // 872
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS":                  DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS,                  // 873
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE":               DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE,               // 874
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION":             DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION,             // 875
	"DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_ID":                   DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_ID,                   // 876
	"DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_SID":                  DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_SID,                  // 877
	"DCGM_FI_DEV_NVSWITCH_DEVICE_UUID":                           DCGM_FI_DEV_NVSWITCH_DEVICE_UUID,                           // 878
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0,                         // 879
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1,                         // 880
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2,                         // 881
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3,                         // 882
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4,                         // 883
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5,                         // 884
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6,                         // 885
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7,                         // 886
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8,                         // 887
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9":                         DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9,                         // 888
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10,                        // 889
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11,                        // 890
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L12":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L12,                        // 891
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L13":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L13,                        // 892
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L14":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L14,                        // 893
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L15":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L15,                        // 894
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L16":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L16,                        // 895
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L17":                        DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L17,                        // 896
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL":                      DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL,                      // 897
	"DCGM_FI_PROF_GR_ENGINE_ACTIVE":                              DCGM_FI_PROF_GR_ENGINE_ACTIVE,                              // 1001
	"DCGM_FI_PROF_SM_ACTIVE":                                     DCGM_FI_PROF_SM_ACTIVE,                                     // 1002
	"DCGM_FI_PROF_SM_OCCUPANCY":                                  DCGM_FI_PROF_SM_OCCUPANCY,                                  // 1003
	"DCGM_FI_PROF_PIPE_TENSOR_ACTIVE":                            DCGM_FI_PROF_PIPE_TENSOR_ACTIVE,                            // 1004
	"DCGM_FI_PROF_DRAM_ACTIVE":                                   DCGM_FI_PROF_DRAM_ACTIVE,                                   // 1005
	"DCGM_FI_PROF_PIPE_FP64_ACTIVE":                              DCGM_FI_PROF_PIPE_FP64_ACTIVE,                              // 1006
	"DCGM_FI_PROF_PIPE_FP32_ACTIVE":                              DCGM_FI_PROF_PIPE_FP32_ACTIVE,                              // 1007
	"DCGM_FI_PROF_PIPE_FP16_ACTIVE":                              DCGM_FI_PROF_PIPE_FP16_ACTIVE,                              // 1008
	"DCGM_FI_PROF_PCIE_TX_BYTES":                                 DCGM_FI_PROF_PCIE_TX_BYTES,                                 // 1009
	"DCGM_FI_PROF_PCIE_RX_BYTES":                                 DCGM_FI_PROF_PCIE_RX_BYTES,                                 // 1010
	"DCGM_FI_PROF_NVLINK_TX_BYTES":                               DCGM_FI_PROF_NVLINK_TX_BYTES,                               // 1011
	"DCGM_FI_PROF_NVLINK_RX_BYTES":                               DCGM_FI_PROF_NVLINK_RX_BYTES,                               // 1012
	"DCGM_FI_PROF_PIPE_TENSOR_IMMA_ACTIVE":                       DCGM_FI_PROF_PIPE_TENSOR_IMMA_ACTIVE,                       // 1013
	"DCGM_FI_PROF_PIPE_TENSOR_HMMA_ACTIVE":                       DCGM_FI_PROF_PIPE_TENSOR_HMMA_ACTIVE,                       // 1014
	"DCGM_FI_PROF_PIPE_TENSOR_DFMA_ACTIVE":                       DCGM_FI_PROF_PIPE_TENSOR_DFMA_ACTIVE,                       // 1015
	"DCGM_FI_PROF_PIPE_INT_ACTIVE":                               DCGM_FI_PROF_PIPE_INT_ACTIVE,                               // 1016
	"DCGM_FI_PROF_NVDEC0_ACTIVE":                                 DCGM_FI_PROF_NVDEC0_ACTIVE,                                 // 1017
	"DCGM_FI_PROF_NVDEC1_ACTIVE":                                 DCGM_FI_PROF_NVDEC1_ACTIVE,                                 // 1018
	"DCGM_FI_PROF_NVDEC2_ACTIVE":                                 DCGM_FI_PROF_NVDEC2_ACTIVE,                                 // 1019
	"DCGM_FI_PROF_NVDEC3_ACTIVE":                                 DCGM_FI_PROF_NVDEC3_ACTIVE,                                 // 1020
	"DCGM_FI_PROF_NVDEC4_ACTIVE":                                 DCGM_FI_PROF_NVDEC4_ACTIVE,                                 // 1021
	"DCGM_FI_PROF_NVDEC5_ACTIVE":                                 DCGM_FI_PROF_NVDEC5_ACTIVE,                                 // 1022
	"DCGM_FI_PROF_NVDEC6_ACTIVE":                                 DCGM_FI_PROF_NVDEC6_ACTIVE,                                 // 1023
	"DCGM_FI_PROF_NVDEC7_ACTIVE":                                 DCGM_FI_PROF_NVDEC7_ACTIVE,                                 // 1024
	"DCGM_FI_PROF_NVJPG0_ACTIVE":                                 DCGM_FI_PROF_NVJPG0_ACTIVE,                                 // 1025
	"DCGM_FI_PROF_NVJPG1_ACTIVE":                                 DCGM_FI_PROF_NVJPG1_ACTIVE,                                 // 1026
	"DCGM_FI_PROF_NVJPG2_ACTIVE":                                 DCGM_FI_PROF_NVJPG2_ACTIVE,                                 // 1027
	"DCGM_FI_PROF_NVJPG3_ACTIVE":                                 DCGM_FI_PROF_NVJPG3_ACTIVE,                                 // 1028
	"DCGM_FI_PROF_NVJPG4_ACTIVE":                                 DCGM_FI_PROF_NVJPG4_ACTIVE,                                 // 1029
	"DCGM_FI_PROF_NVJPG5_ACTIVE":                                 DCGM_FI_PROF_NVJPG5_ACTIVE,                                 // 1030
	"DCGM_FI_PROF_NVJPG6_ACTIVE":                                 DCGM_FI_PROF_NVJPG6_ACTIVE,                                 // 1031
	"DCGM_FI_PROF_NVJPG7_ACTIVE":                                 DCGM_FI_PROF_NVJPG7_ACTIVE,                                 // 1032
	"DCGM_FI_PROF_NVOFA0_ACTIVE":                                 DCGM_FI_PROF_NVOFA0_ACTIVE,                                 // 1033
	"DCGM_FI_PROF_NVOFA1_ACTIVE":                                 DCGM_FI_PROF_NVOFA1_ACTIVE,                                 // 1034
	"DCGM_FI_PROF_NVLINK_L0_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L0_TX_BYTES,                            // 1040
	"DCGM_FI_PROF_NVLINK_L0_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L0_RX_BYTES,                            // 1041
	"DCGM_FI_PROF_NVLINK_L1_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L1_TX_BYTES,                            // 1042
	"DCGM_FI_PROF_NVLINK_L1_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L1_RX_BYTES,                            // 1043
	"DCGM_FI_PROF_NVLINK_L2_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L2_TX_BYTES,                            // 1044
	"DCGM_FI_PROF_NVLINK_L2_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L2_RX_BYTES,                            // 1045
	"DCGM_FI_PROF_NVLINK_L3_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L3_TX_BYTES,                            // 1046
	"DCGM_FI_PROF_NVLINK_L3_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L3_RX_BYTES,                            // 1047
	"DCGM_FI_PROF_NVLINK_L4_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L4_TX_BYTES,                            // 1048
	"DCGM_FI_PROF_NVLINK_L4_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L4_RX_BYTES,                            // 1049
	"DCGM_FI_PROF_NVLINK_L5_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L5_TX_BYTES,                            // 1050
	"DCGM_FI_PROF_NVLINK_L5_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L5_RX_BYTES,                            // 1051
	"DCGM_FI_PROF_NVLINK_L6_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L6_TX_BYTES,                            // 1052
	"DCGM_FI_PROF_NVLINK_L6_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L6_RX_BYTES,                            // 1053
	"DCGM_FI_PROF_NVLINK_L7_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L7_TX_BYTES,                            // 1054
	"DCGM_FI_PROF_NVLINK_L7_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L7_RX_BYTES,                            // 1055
	"DCGM_FI_PROF_NVLINK_L8_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L8_TX_BYTES,                            // 1056
	"DCGM_FI_PROF_NVLINK_L8_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L8_RX_BYTES,                            // 1057
	"DCGM_FI_PROF_NVLINK_L9_TX_BYTES":                            DCGM_FI_PROF_NVLINK_L9_TX_BYTES,                            // 1058
	"DCGM_FI_PROF_NVLINK_L9_RX_BYTES":                            DCGM_FI_PROF_NVLINK_L9_RX_BYTES,                            // 1059
	"DCGM_FI_PROF_NVLINK_L10_TX_BYTES":                           DCGM_FI_PROF_NVLINK_L10_TX_BYTES,                           // 1060
	"DCGM_FI_PROF_NVLINK_L10_RX_BYTES":                           DCGM_FI_PROF_NVLINK_L10_RX_BYTES,                           // 1061
	"DCGM_FI_PROF_NVLINK_L11_TX_BYTES":                           DCGM_FI_PROF_NVLINK_L11_TX_BYTES,                           // 1062
	"DCGM_FI_PROF_NVLINK_L11_RX_BYTES":                           DCGM_FI_PROF_NVLINK_L11_RX_BYTES,                           // 1063
	"DCGM_FI_PROF_NVLINK_L12_TX_BYTES":                           DCGM_FI_PROF_NVLINK_L12_TX_BYTES,                           // 1064
	"DCGM_FI_PROF_NVLINK_L12_RX_BYTES":                           DCGM_FI_PROF_NVLINK_L12_RX_BYTES,                           // 1065
	"DCGM_FI_PROF_NVLINK_L13_TX_BYTES":                           DCGM_FI_PROF_NVLINK_L13_TX_BYTES,                           // 1066
	"DCGM_FI_PROF_NVLINK_L13_RX_BYTES":                           DCGM_FI_PROF_NVLINK_L13_RX_BYTES,                           // 1067
	"DCGM_FI_PROF_NVLINK_L14_TX_BYTES":                           DCGM_FI_PROF_NVLINK_L14_TX_BYTES,                           // 1068
	"DCGM_FI_PROF_NVLINK_L14_RX_BYTES":                           DCGM_FI_PROF_NVLINK_L14_RX_BYTES,                           // 1069
	"DCGM_FI_PROF_NVLINK_L15_TX_BYTES":                           DCGM_FI_PROF_NVLINK_L15_TX_BYTES,                           // 1070
	"DCGM_FI_PROF_NVLINK_L15_RX_BYTES":                           DCGM_FI_PROF_NVLINK_L15_RX_BYTES,                           // 1071
	"DCGM_FI_PROF_NVLINK_L16_TX_BYTES":                           DCGM_FI_PROF_NVLINK_L16_TX_BYTES,                           // 1072
	"DCGM_FI_PROF_C2C_TX_ALL_BYTES":                              DCGM_FI_PROF_C2C_TX_ALL_BYTES,                              // 1076
	"DCGM_FI_PROF_C2C_TX_DATA_BYTES":                             DCGM_FI_PROF_C2C_TX_DATA_BYTES,                             // 1077
	"DCGM_FI_PROF_C2C_RX_ALL_BYTES":                              DCGM_FI_PROF_C2C_RX_ALL_BYTES,                              // 1078
	"DCGM_FI_PROF_C2C_RX_DATA_BYTES":                             DCGM_FI_PROF_C2C_RX_DATA_BYTES,                             // 1079
	"DCGM_FI_DEV_CPU_UTIL_TOTAL":                                 DCGM_FI_DEV_CPU_UTIL_TOTAL,                                 // 1100
	"DCGM_FI_DEV_CPU_UTIL_USER":                                  DCGM_FI_DEV_CPU_UTIL_USER,                                  // 1101
	"DCGM_FI_DEV_CPU_UTIL_NICE":                                  DCGM_FI_DEV_CPU_UTIL_NICE,                                  // 1102
	"DCGM_FI_DEV_CPU_UTIL_SYS":                                   DCGM_FI_DEV_CPU_UTIL_SYS,                                   // 1103
	"DCGM_FI_DEV_CPU_UTIL_IRQ":                                   DCGM_FI_DEV_CPU_UTIL_IRQ,                                   // 1104
	"DCGM_FI_DEV_CPU_TEMP_CURRENT":                               DCGM_FI_DEV_CPU_TEMP_CURRENT,                               // 1110
	"DCGM_FI_DEV_CPU_TEMP_WARNING":                               DCGM_FI_DEV_CPU_TEMP_WARNING,                               // 1111
	"DCGM_FI_DEV_CPU_TEMP_SHUTDOWN":                              DCGM_FI_DEV_CPU_TEMP_SHUTDOWN,                              // 1112
	"DCGM_FI_DEV_CPU_CLOCK_CURRENT":                              DCGM_FI_DEV_CPU_CLOCK_CURRENT,                              // 1120
	"DCGM_FI_DEV_CPU_POWER_CURRENT":                              DCGM_FI_DEV_CPU_POWER_CURRENT,                              // 1130
	"DCGM_FI_DEV_CPU_POWER_LIMIT":                                DCGM_FI_DEV_CPU_POWER_LIMIT,                                // 1131
	"DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT":                       DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT,                       // 1132
	"DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT":                      DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT,                      // 1133
	"DCGM_FI_DEV_CPU_VENDOR":                                     DCGM_FI_DEV_CPU_VENDOR,                                     // 1140
	"DCGM_FI_DEV_CPU_MODEL":                                      DCGM_FI_DEV_CPU_MODEL,                                      // 1141
	"DCGM_FI_DEV_NVLINK_COUNT_TX_PACKETS":                        DCGM_FI_DEV_NVLINK_COUNT_TX_PACKETS,                        // 1200
	"DCGM_FI_DEV_NVLINK_COUNT_TX_BYTES":                          DCGM_FI_DEV_NVLINK_COUNT_TX_BYTES,                          // 1201
	"DCGM_FI_DEV_NVLINK_COUNT_RX_PACKETS":                        DCGM_FI_DEV_NVLINK_COUNT_RX_PACKETS,                        // 1202
	"DCGM_FI_DEV_NVLINK_COUNT_RX_BYTES":                          DCGM_FI_DEV_NVLINK_COUNT_RX_BYTES,                          // 1203
	"DCGM_FI_DEV_NVLINK_COUNT_RX_MALFORMED_PACKET_ERRORS":        DCGM_FI_DEV_NVLINK_COUNT_RX_MALFORMED_PACKET_ERRORS,        // 1204
	"DCGM_FI_DEV_NVLINK_COUNT_RX_BUFFER_OVERRUN_ERRORS":          DCGM_FI_DEV_NVLINK_COUNT_RX_BUFFER_OVERRUN_ERRORS,          // 1205
	"DCGM_FI_DEV_NVLINK_COUNT_RX_ERRORS":                         DCGM_FI_DEV_NVLINK_COUNT_RX_ERRORS,                         // 1206
	"DCGM_FI_DEV_NVLINK_COUNT_RX_REMOTE_ERRORS":                  DCGM_FI_DEV_NVLINK_COUNT_RX_REMOTE_ERRORS,                  // 1207
	"DCGM_FI_DEV_NVLINK_COUNT_RX_GENERAL_ERRORS":                 DCGM_FI_DEV_NVLINK_COUNT_RX_GENERAL_ERRORS,                 // 1208
	"DCGM_FI_DEV_NVLINK_COUNT_LOCAL_LINK_INTEGRITY_ERRORS":       DCGM_FI_DEV_NVLINK_COUNT_LOCAL_LINK_INTEGRITY_ERRORS,       // 1209
	"DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS":                       DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS,                       // 1210
	"DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_SUCCESSFUL_EVENTS":   DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_SUCCESSFUL_EVENTS,   // 1211
	"DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_FAILED_EVENTS":       DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_FAILED_EVENTS,       // 1212
	"DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_EVENTS":              DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_EVENTS,              // 1213
	"DCGM_FI_DEV_NVLINK_COUNT_RX_SYMBOL_ERRORS":                  DCGM_FI_DEV_NVLINK_COUNT_RX_SYMBOL_ERRORS,                  // 1214
	"DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER":                        DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER,                        // 1215
	"DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER_FLOAT":                  DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER_FLOAT,                  // 1216
	"DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER":                     DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER,                     // 1217
	"DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER_FLOAT":               DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_BER_FLOAT,               // 1218
	"DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_ERRORS":                  DCGM_FI_DEV_NVLINK_COUNT_EFFECTIVE_ERRORS,                  // 1219
	"DCGM_FI_DEV_CONNECTX_HEALTH":                                DCGM_FI_DEV_CONNECTX_HEALTH,                                // 1300
	"DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH":                DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH,                // 1301
	"DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED":                DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED,                // 1302
	"DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH":                DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH,                // 1303
	"DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED":                DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED,                // 1304
	"DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS":                DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS,                // 1305
	"DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK":                  DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK,                  // 1306
	"DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS":              DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS,              // 1307
	"DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK":                DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK,                // 1308
	"DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY":            DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY,            // 1309
	"DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE":                    DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE,                    // 1310
	"DCGM_FI_DEV_LAST_CONNECTX_FIELD_ID":                         DCGM_FI_DEV_LAST_CONNECTX_FIELD_ID,                         // 1399
	"DCGM_FI_DEV_C2C_LINK_ERROR_INTR":                            DCGM_FI_DEV_C2C_LINK_ERROR_INTR,                            // 1400
	"DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY":                          DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY,                          // 1401
	"DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY_B2B":                      DCGM_FI_DEV_C2C_LINK_ERROR_REPLAY_B2B,                      // 1402
	"DCGM_FI_DEV_C2C_LINK_POWER_STATE":                           DCGM_FI_DEV_C2C_LINK_POWER_STATE,                           // 1403
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_0":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_0,                     // 1404
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_1":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_1,                     // 1405
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_2":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_2,                     // 1406
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_3":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_3,                     // 1407
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_4":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_4,                     // 1408
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_5":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_5,                     // 1409
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_6":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_6,                     // 1410
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_7":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_7,                     // 1411
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_8":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_8,                     // 1412
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_9":                     DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_9,                     // 1413
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_10":                    DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_10,                    // 1414
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_11":                    DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_11,                    // 1415
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_12":                    DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_12,                    // 1416
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_13":                    DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_13,                    // 1417
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_14":                    DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_14,                    // 1418
	"DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_15":                    DCGM_FI_DEV_NVLINK_COUNT_FEC_HISTORY_15,                    // 1419
	"DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_POWER_CAP_NS":            DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_POWER_CAP_NS,            // 1420
	"DCGM_FI_DEV_CLOCKS_EVENT_REASON_SYNC_BOOST_NS":              DCGM_FI_DEV_CLOCKS_EVENT_REASON_SYNC_BOOST_NS,              // 1421
	"DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_THERM_SLOWDOWN_NS":       DCGM_FI_DEV_CLOCKS_EVENT_REASON_SW_THERM_SLOWDOWN_NS,       // 1422
	"DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_THERM_SLOWDOWN_NS":       DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_THERM_SLOWDOWN_NS,       // 1423
	"DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_POWER_BRAKE_SLOWDOWN_NS": DCGM_FI_DEV_CLOCKS_EVENT_REASON_HW_POWER_BRAKE_SLOWDOWN_NS, // 1424
	"DCGM_FI_MAX_FIELDS":                                         DCGM_FI_MAX_FIELDS,                                         // 1425
}

var legacyDCGMFields = map[string]Short{
	"dcgm_sm_clock":                          100,
	"dcgm_memory_clock":                      101,
	"dcgm_memory_temp":                       140,
	"dcgm_gpu_temp":                          150,
	"dcgm_power_usage":                       155,
	"dcgm_total_energy_consumption":          156,
	"dcgm_pcie_tx_throughput":                200,
	"dcgm_pcie_rx_throughput":                201,
	"dcgm_pcie_replay_counter":               202,
	"dcgm_gpu_utilization":                   203,
	"dcgm_mem_copy_utilization":              204,
	"dcgm_enc_utilization":                   206,
	"dcgm_dec_utilization":                   207,
	"dcgm_xid_errors":                        230,
	"dcgm_power_violation":                   240,
	"dcgm_thermal_violation":                 241,
	"dcgm_sync_boost_violation":              242,
	"dcgm_board_limit_violation":             243,
	"dcgm_low_util_violation":                244,
	"dcgm_reliability_violation":             245,
	"dcgm_fb_free":                           251,
	"dcgm_fb_used":                           252,
	"dcgm_ecc_sbe_volatile_total":            310,
	"dcgm_ecc_dbe_volatile_total":            311,
	"dcgm_ecc_sbe_aggregate_total":           312,
	"dcgm_ecc_dbe_aggregate_total":           313,
	"dcgm_retired_pages_sbe":                 390,
	"dcgm_retired_pages_dbe":                 391,
	"dcgm_retired_pages_pending":             392,
	"dcgm_nvlink_flit_crc_error_count_total": 409,
	"dcgm_nvlink_data_crc_error_count_total": 419,
	"dcgm_nvlink_replay_error_count_total":   429,
	"dcgm_nvlink_recovery_error_count_total": 439,
	"dcgm_nvlink_bandwidth_total":            449,
	"dcgm_fi_prof_gr_engine_active":          1001,
	"dcgm_fi_prof_sm_active":                 1002,
	"dcgm_fi_prof_sm_occupancy":              1003,
	"dcgm_fi_prof_pipe_tensor_active":        1004,
	"dcgm_fi_prof_dram_active":               1005,
	"dcgm_fi_prof_pcie_tx_bytes":             1009,
	"dcgm_fi_prof_pcie_rx_bytes":             1010,
}

// GetFieldID returns the DCGM field ID for a given field name and whether it was found
// It first checks the current field IDs, then falls back to legacy field IDs if not found
func GetFieldID(fieldName string) (Short, bool) {
	// First check current field IDs
	if fieldID, ok := dcgmFields[fieldName]; ok {
		return fieldID, true
	}

	// Then check legacy field IDs
	if fieldID, ok := legacyDCGMFields[fieldName]; ok {
		return fieldID, true
	}

	return 0, false
}

// GetFieldIDOrPanic returns the DCGM field ID for a given field name
// It panics if the field name is not found in either current or legacy maps
func GetFieldIDOrPanic(fieldName string) Short {
	fieldID, ok := GetFieldID(fieldName)
	if !ok {
		panic("field name not found: " + fieldName)
	}
	return fieldID
}

// IsLegacyField returns true if the given field name is a legacy field
func IsLegacyField(fieldName string) bool {
	_, ok := legacyDCGMFields[fieldName]
	return ok
}

// IsCurrentField returns true if the given field name is a current field
func IsCurrentField(fieldName string) bool {
	_, ok := dcgmFields[fieldName]
	return ok
}
