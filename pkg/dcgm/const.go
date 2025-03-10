package dcgm

import "C"

type Short C.ushort

type FieldValue_v1 struct {
	Version   uint
	FieldId   uint
	FieldType uint
	Status    int
	Ts        int64
	Value     [4096]byte
}

type FieldValue_v2 struct {
	Version       uint
	EntityGroupId Field_Entity_Group
	EntityId      uint
	FieldId       uint
	FieldType     uint
	Status        int
	Ts            int64
	Value         [4096]byte
	StringValue   *string
}

const (
	DCGM_FT_BINARY                 = uint('b')
	DCGM_FT_DOUBLE                 = uint('d')
	DCGM_FT_INT64                  = uint('i')
	DCGM_FT_STRING                 = uint('s')
	DCGM_FT_TIMESTAMP              = uint('t')
	DCGM_FT_INT32_BLANK            = int64(2147483632)
	DCGM_FT_INT32_NOT_FOUND        = int64(DCGM_FT_INT32_BLANK + 1)
	DCGM_FT_INT32_NOT_SUPPORTED    = int64(DCGM_FT_INT32_BLANK + 2)
	DCGM_FT_INT32_NOT_PERMISSIONED = int64(DCGM_FT_INT32_BLANK + 3)
	DCGM_FT_INT64_BLANK            = int64(9223372036854775792)
	DCGM_FT_INT64_NOT_FOUND        = int64(DCGM_FT_INT64_BLANK + 1)
	DCGM_FT_INT64_NOT_SUPPORTED    = int64(DCGM_FT_INT64_BLANK + 2)
	DCGM_FT_INT64_NOT_PERMISSIONED = int64(DCGM_FT_INT64_BLANK + 3)
	DCGM_FT_FP64_BLANK             = 140737488355328.0
	DCGM_FT_FP64_NOT_FOUND         = float64(DCGM_FT_FP64_BLANK + 1.0)
	DCGM_FT_FP64_NOT_SUPPORTED     = float64(DCGM_FT_FP64_BLANK + 2.0)
	DCGM_FT_FP64_NOT_PERMISSIONED  = float64(DCGM_FT_FP64_BLANK + 3.0)
	DCGM_FT_STR_BLANK              = "<<<NULL>>>"
	DCGM_FT_STR_NOT_FOUND          = "<<<NOT_FOUND>>>"
	DCGM_FT_STR_NOT_SUPPORTED      = "<<<NOT_SUPPORTED>>>"
	DCGM_FT_STR_NOT_PERMISSIONED   = "<<<NOT_PERM>>>"

	DCGM_FI_UNKNOWN                                          = 0
	DCGM_FI_DRIVER_VERSION                                   = 1
	DCGM_FI_NVML_VERSION                                     = 2
	DCGM_FI_PROCESS_NAME                                     = 3
	DCGM_FI_DEV_COUNT                                        = 4
	DCGM_FI_CUDA_DRIVER_VERSION                              = 5
	DCGM_FI_DEV_NAME                                         = 50
	DCGM_FI_DEV_BRAND                                        = 51
	DCGM_FI_DEV_NVML_INDEX                                   = 52
	DCGM_FI_DEV_SERIAL                                       = 53
	DCGM_FI_DEV_UUID                                         = 54
	DCGM_FI_DEV_MINOR_NUMBER                                 = 55
	DCGM_FI_DEV_OEM_INFOROM_VER                              = 56
	DCGM_FI_DEV_PCI_BUSID                                    = 57
	DCGM_FI_DEV_PCI_COMBINED_ID                              = 58
	DCGM_FI_DEV_PCI_SUBSYS_ID                                = 59
	DCGM_FI_GPU_TOPOLOGY_PCI                                 = 60
	DCGM_FI_GPU_TOPOLOGY_NVLINK                              = 61
	DCGM_FI_GPU_TOPOLOGY_AFFINITY                            = 62
	DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY                      = 63
	DCGM_FI_DEV_COMPUTE_MODE                                 = 65
	DCGM_FI_DEV_PERSISTENCE_MODE                             = 66
	DCGM_FI_DEV_MIG_MODE                                     = 67
	DCGM_FI_DEV_CUDA_VISIBLE_DEVICES_STR                     = 68
	DCGM_FI_DEV_MIG_MAX_SLICES                               = 69
	DCGM_FI_DEV_CPU_AFFINITY_0                               = 70
	DCGM_FI_DEV_CPU_AFFINITY_1                               = 71
	DCGM_FI_DEV_CPU_AFFINITY_2                               = 72
	DCGM_FI_DEV_CPU_AFFINITY_3                               = 73
	DCGM_FI_DEV_CC_MODE                                      = 74
	DCGM_FI_DEV_MIG_ATTRIBUTES                               = 75
	DCGM_FI_DEV_MIG_GI_INFO                                  = 76
	DCGM_FI_DEV_MIG_CI_INFO                                  = 77
	DCGM_FI_DEV_ECC_INFOROM_VER                              = 80
	DCGM_FI_DEV_POWER_INFOROM_VER                            = 81
	DCGM_FI_DEV_INFOROM_IMAGE_VER                            = 82
	DCGM_FI_DEV_INFOROM_CONFIG_CHECK                         = 83
	DCGM_FI_DEV_INFOROM_CONFIG_VALID                         = 84
	DCGM_FI_DEV_VBIOS_VERSION                                = 85
	DCGM_FI_DEV_MEM_AFFINITY_0                               = 86
	DCGM_FI_DEV_MEM_AFFINITY_1                               = 87
	DCGM_FI_DEV_MEM_AFFINITY_2                               = 88
	DCGM_FI_DEV_MEM_AFFINITY_3                               = 89
	DCGM_FI_DEV_BAR1_TOTAL                                   = 90
	DCGM_FI_SYNC_BOOST                                       = 91
	DCGM_FI_DEV_BAR1_USED                                    = 92
	DCGM_FI_DEV_BAR1_FREE                                    = 93
	DCGM_FI_DEV_GPM_SUPPORT                                  = 94
	DCGM_FI_DEV_SM_CLOCK                                     = 100
	DCGM_FI_DEV_MEM_CLOCK                                    = 101
	DCGM_FI_DEV_VIDEO_CLOCK                                  = 102
	DCGM_FI_DEV_APP_SM_CLOCK                                 = 110
	DCGM_FI_DEV_APP_MEM_CLOCK                                = 111
	DCGM_FI_DEV_CLOCKS_EVENT_REASONS                         = 112
	DCGM_FI_DEV_CLOCK_THROTTLE_REASONS                       = DCGM_FI_DEV_CLOCKS_EVENT_REASONS
	DCGM_FI_DEV_MAX_SM_CLOCK                                 = 113
	DCGM_FI_DEV_MAX_MEM_CLOCK                                = 114
	DCGM_FI_DEV_MAX_VIDEO_CLOCK                              = 115
	DCGM_FI_DEV_AUTOBOOST                                    = 120
	DCGM_FI_DEV_SUPPORTED_CLOCKS                             = 130
	DCGM_FI_DEV_MEMORY_TEMP                                  = 140
	DCGM_FI_DEV_GPU_TEMP                                     = 150
	DCGM_FI_DEV_MEM_MAX_OP_TEMP                              = 151
	DCGM_FI_DEV_GPU_MAX_OP_TEMP                              = 152
	DCGM_FI_DEV_GPU_TEMP_LIMIT                               = 153
	DCGM_FI_DEV_POWER_USAGE                                  = 155
	DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION                     = 156
	DCGM_FI_DEV_POWER_USAGE_INSTANT                          = 157
	DCGM_FI_DEV_SLOWDOWN_TEMP                                = 158
	DCGM_FI_DEV_SHUTDOWN_TEMP                                = 159
	DCGM_FI_DEV_POWER_MGMT_LIMIT                             = 160
	DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN                         = 161
	DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX                         = 162
	DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF                         = 163
	DCGM_FI_DEV_ENFORCED_POWER_LIMIT                         = 164
	DCGM_FI_DEV_REQUESTED_POWER_PROFILE_MASK                 = 165
	DCGM_FI_DEV_ENFORCED_POWER_PROFILE_MASK                  = 166
	DCGM_FI_DEV_VALID_POWER_PROFILE_MASK                     = 167
	DCGM_FI_DEV_FABRIC_MANAGER_STATUS                        = 170
	DCGM_FI_DEV_FABRIC_MANAGER_ERROR_CODE                    = 171
	DCGM_FI_DEV_FABRIC_CLUSTER_UUID                          = 172
	DCGM_FI_DEV_FABRIC_CLIQUE_ID                             = 173
	DCGM_FI_DEV_PSTATE                                       = 190
	DCGM_FI_DEV_FAN_SPEED                                    = 191
	DCGM_FI_DEV_PCIE_TX_THROUGHPUT                           = 200
	DCGM_FI_DEV_PCIE_RX_THROUGHPUT                           = 201
	DCGM_FI_DEV_PCIE_REPLAY_COUNTER                          = 202
	DCGM_FI_DEV_GPU_UTIL                                     = 203
	DCGM_FI_DEV_MEM_COPY_UTIL                                = 204
	DCGM_FI_DEV_ACCOUNTING_DATA                              = 205
	DCGM_FI_DEV_ENC_UTIL                                     = 206
	DCGM_FI_DEV_DEC_UTIL                                     = 207
	DCGM_FI_DEV_XID_ERRORS                                   = 230
	DCGM_FI_DEV_PCIE_MAX_LINK_GEN                            = 235
	DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH                          = 236
	DCGM_FI_DEV_PCIE_LINK_GEN                                = 237
	DCGM_FI_DEV_PCIE_LINK_WIDTH                              = 238
	DCGM_FI_DEV_POWER_VIOLATION                              = 240
	DCGM_FI_DEV_THERMAL_VIOLATION                            = 241
	DCGM_FI_DEV_SYNC_BOOST_VIOLATION                         = 242
	DCGM_FI_DEV_BOARD_LIMIT_VIOLATION                        = 243
	DCGM_FI_DEV_LOW_UTIL_VIOLATION                           = 244
	DCGM_FI_DEV_RELIABILITY_VIOLATION                        = 245
	DCGM_FI_DEV_TOTAL_APP_CLOCKS_VIOLATION                   = 246
	DCGM_FI_DEV_TOTAL_BASE_CLOCKS_VIOLATION                  = 247
	DCGM_FI_DEV_FB_TOTAL                                     = 250
	DCGM_FI_DEV_FB_FREE                                      = 251
	DCGM_FI_DEV_FB_USED                                      = 252
	DCGM_FI_DEV_FB_RESERVED                                  = 253
	DCGM_FI_DEV_FB_USED_PERCENT                              = 254
	DCGM_FI_DEV_C2C_LINK_COUNT                               = 285
	DCGM_FI_DEV_C2C_LINK_STATUS                              = 286
	DCGM_FI_DEV_C2C_MAX_BANDWIDTH                            = 287
	DCGM_FI_DEV_ECC_CURRENT                                  = 300
	DCGM_FI_DEV_ECC_PENDING                                  = 301
	DCGM_FI_DEV_ECC_SBE_VOL_TOTAL                            = 310
	DCGM_FI_DEV_ECC_DBE_VOL_TOTAL                            = 311
	DCGM_FI_DEV_ECC_SBE_AGG_TOTAL                            = 312
	DCGM_FI_DEV_ECC_DBE_AGG_TOTAL                            = 313
	DCGM_FI_DEV_ECC_SBE_VOL_L1                               = 314
	DCGM_FI_DEV_ECC_DBE_VOL_L1                               = 315
	DCGM_FI_DEV_ECC_SBE_VOL_L2                               = 316
	DCGM_FI_DEV_ECC_DBE_VOL_L2                               = 317
	DCGM_FI_DEV_ECC_SBE_VOL_DEV                              = 318
	DCGM_FI_DEV_ECC_DBE_VOL_DEV                              = 319
	DCGM_FI_DEV_ECC_SBE_VOL_REG                              = 320
	DCGM_FI_DEV_ECC_DBE_VOL_REG                              = 321
	DCGM_FI_DEV_ECC_SBE_VOL_TEX                              = 322
	DCGM_FI_DEV_ECC_DBE_VOL_TEX                              = 323
	DCGM_FI_DEV_ECC_SBE_AGG_L1                               = 324
	DCGM_FI_DEV_ECC_DBE_AGG_L1                               = 325
	DCGM_FI_DEV_ECC_SBE_AGG_L2                               = 326
	DCGM_FI_DEV_ECC_DBE_AGG_L2                               = 327
	DCGM_FI_DEV_ECC_SBE_AGG_DEV                              = 328
	DCGM_FI_DEV_ECC_DBE_AGG_DEV                              = 329
	DCGM_FI_DEV_ECC_SBE_AGG_REG                              = 330
	DCGM_FI_DEV_ECC_DBE_AGG_REG                              = 331
	DCGM_FI_DEV_ECC_SBE_AGG_TEX                              = 332
	DCGM_FI_DEV_ECC_DBE_AGG_TEX                              = 333
	DCGM_FI_DEV_ECC_SBE_VOL_SHM                              = 334
	DCGM_FI_DEV_ECC_DBE_VOL_SHM                              = 335
	DCGM_FI_DEV_ECC_SBE_VOL_CBU                              = 336
	DCGM_FI_DEV_ECC_DBE_VOL_CBU                              = 337
	DCGM_FI_DEV_ECC_SBE_AGG_SHM                              = 338
	DCGM_FI_DEV_ECC_DBE_AGG_SHM                              = 339
	DCGM_FI_DEV_ECC_SBE_AGG_CBU                              = 340
	DCGM_FI_DEV_ECC_DBE_AGG_CBU                              = 341
	DCGM_FI_DEV_ECC_SBE_VOL_SRM                              = 342
	DCGM_FI_DEV_ECC_DBE_VOL_SRM                              = 343
	DCGM_FI_DEV_ECC_SBE_AGG_SRM                              = 344
	DCGM_FI_DEV_ECC_DBE_AGG_SRM                              = 345
	DCGM_FI_DEV_DIAG_MEMORY_RESULT                           = 350
	DCGM_FI_DEV_DIAG_DIAGNOSTIC_RESULT                       = 351
	DCGM_FI_DEV_DIAG_PCIE_RESULT                             = 352
	DCGM_FI_DEV_DIAG_TARGETED_STRESS_RESULT                  = 353
	DCGM_FI_DEV_DIAG_TARGETED_POWER_RESULT                   = 354
	DCGM_FI_DEV_DIAG_MEMORY_BANDWIDTH_RESULT                 = 355
	DCGM_FI_DEV_DIAG_MEMTEST_RESULT                          = 356
	DCGM_FI_DEV_DIAG_PULSE_TEST_RESULT                       = 357
	DCGM_FI_DEV_DIAG_EUD_RESULT                              = 358
	DCGM_FI_DEV_DIAG_CPU_EUD_RESULT                          = 359
	DCGM_FI_DEV_DIAG_SOFTWARE_RESULT                         = 360
	DCGM_FI_DEV_DIAG_NVBANDWIDTH_RESULT                      = 361
	DCGM_FI_DEV_DIAG_STATUS                                  = 362
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_MAX                   = 385
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_HIGH                  = 386
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_PARTIAL               = 387
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_LOW                   = 388
	DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_NONE                  = 389
	DCGM_FI_DEV_RETIRED_SBE                                  = 390
	DCGM_FI_DEV_RETIRED_DBE                                  = 391
	DCGM_FI_DEV_RETIRED_PENDING                              = 392
	DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS                  = 393
	DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS                    = 394
	DCGM_FI_DEV_ROW_REMAP_FAILURE                            = 395
	DCGM_FI_DEV_ROW_REMAP_PENDING                            = 396
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L0               = 400
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L1               = 401
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L2               = 402
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L3               = 403
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L4               = 404
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L5               = 405
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL            = 409
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L0               = 410
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L1               = 411
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L2               = 412
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L3               = 413
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L4               = 414
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L5               = 415
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_TOTAL            = 419
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L0                 = 420
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L1                 = 421
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L2                 = 422
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L3                 = 423
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L4                 = 424
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L5                 = 425
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_TOTAL              = 429
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L0               = 430
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L1               = 431
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L2               = 432
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L3               = 433
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L4               = 434
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L5               = 435
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_TOTAL            = 439
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L0                          = 440
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L1                          = 441
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L2                          = 442
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L3                          = 443
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L4                          = 444
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L5                          = 445
	DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL                       = 449
	DCGM_FI_DEV_GPU_NVLINK_ERRORS                            = 450
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L6               = 451
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L7               = 452
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L8               = 453
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L9               = 454
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L10              = 455
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L11              = 456
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L6               = 457
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L7               = 458
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L8               = 459
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L9               = 460
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L10              = 461
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L11              = 462
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L6                 = 463
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L7                 = 464
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L8                 = 465
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L9                 = 466
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L10                = 467
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L11                = 468
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L6               = 469
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L7               = 470
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L8               = 471
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L9               = 472
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L10              = 473
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L11              = 474
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L6                          = 475
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L7                          = 476
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L8                          = 477
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L9                          = 478
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L10                         = 479
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L11                         = 480
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L12              = 406
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13              = 407
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14              = 408
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L15              = 481
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L16              = 482
	DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L17              = 483
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L12              = 416
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L13              = 417
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L14              = 418
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L15              = 484
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L16              = 485
	DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L17              = 486
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L12                = 426
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L13                = 427
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L14                = 428
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L15                = 487
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L16                = 488
	DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L17                = 489
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L12              = 436
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L13              = 437
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L14              = 438
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L15              = 491
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L16              = 492
	DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L17              = 493
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L12                         = 446
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L13                         = 447
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L14                         = 448
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L15                         = 494
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L16                         = 495
	DCGM_FI_DEV_NVLINK_BANDWIDTH_L17                         = 496
	DCGM_FI_DEV_NVLINK_ERROR_DL_CRC                          = 497
	DCGM_FI_DEV_NVLINK_ERROR_DL_RECOVERY                     = 498
	DCGM_FI_DEV_NVLINK_ERROR_DL_REPLAY                       = 499
	DCGM_FI_DEV_VIRTUAL_MODE                                 = 500
	DCGM_FI_DEV_SUPPORTED_TYPE_INFO                          = 501
	DCGM_FI_DEV_CREATABLE_VGPU_TYPE_IDS                      = 502
	DCGM_FI_DEV_VGPU_INSTANCE_IDS                            = 503
	DCGM_FI_DEV_VGPU_UTILIZATIONS                            = 504
	DCGM_FI_DEV_VGPU_PER_PROCESS_UTILIZATION                 = 505
	DCGM_FI_DEV_ENC_STATS                                    = 506
	DCGM_FI_DEV_FBC_STATS                                    = 507
	DCGM_FI_DEV_FBC_SESSIONS_INFO                            = 508
	DCGM_FI_DEV_SUPPORTED_VGPU_TYPE_IDS                      = 509
	DCGM_FI_DEV_VGPU_TYPE_INFO                               = 510
	DCGM_FI_DEV_VGPU_TYPE_NAME                               = 511
	DCGM_FI_DEV_VGPU_TYPE_CLASS                              = 512
	DCGM_FI_DEV_VGPU_TYPE_LICENSE                            = 513
	DCGM_FI_DEV_VGPU_VM_ID                                   = 520
	DCGM_FI_DEV_VGPU_VM_NAME                                 = 521
	DCGM_FI_DEV_VGPU_TYPE                                    = 522
	DCGM_FI_DEV_VGPU_UUID                                    = 523
	DCGM_FI_DEV_VGPU_DRIVER_VERSION                          = 524
	DCGM_FI_DEV_VGPU_MEMORY_USAGE                            = 525
	DCGM_FI_DEV_VGPU_LICENSE_STATUS                          = 526
	DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT                        = 527
	DCGM_FI_DEV_VGPU_ENC_STATS                               = 528
	DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO                       = 529
	DCGM_FI_DEV_VGPU_FBC_STATS                               = 530
	DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO                       = 531
	DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATE                  = 532
	DCGM_FI_DEV_VGPU_PCI_ID                                  = 533
	DCGM_FI_DEV_VGPU_VM_GPU_INSTANCE_ID                      = 534
	DCGM_FI_FIRST_VGPU_FIELD_ID                              = 520
	DCGM_FI_LAST_VGPU_FIELD_ID                               = 570
	DCGM_FI_DEV_PLATFORM_INFINIBAND_GUID                     = 571
	DCGM_FI_DEV_PLATFORM_CHASSIS_SERIAL_NUMBER               = 572
	DCGM_FI_DEV_PLATFORM_CHASSIS_SLOT_NUMBER                 = 573
	DCGM_FI_DEV_PLATFORM_TRAY_INDEX                          = 574
	DCGM_FI_DEV_PLATFORM_HOST_ID                             = 575
	DCGM_FI_DEV_PLATFORM_PEER_TYPE                           = 576
	DCGM_FI_DEV_PLATFORM_MODULE_ID                           = 577
	DCGM_FI_FIRST_NVSWITCH_FIELD_ID                          = 700
	DCGM_FI_DEV_NVSWITCH_VOLTAGE_MVOLT                       = 701
	DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ                        = 702
	DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_REV                    = 703
	DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_DVDD                   = 704
	DCGM_FI_DEV_NVSWITCH_POWER_VDD                           = 705
	DCGM_FI_DEV_NVSWITCH_POWER_DVDD                          = 706
	DCGM_FI_DEV_NVSWITCH_POWER_HVDD                          = 707
	DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX                  = 780
	DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX                  = 781
	DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS                   = 782
	DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS               = 783
	DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS                  = 784
	DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS                = 785
	DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS                    = 786
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS                     = 787
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS                     = 788
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0                = 789
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1                = 790
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2                = 791
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3                = 792
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0             = 793
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1             = 794
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2             = 795
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3             = 796
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0               = 797
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1               = 798
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2               = 799
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3               = 800
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0              = 801
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1              = 802
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2              = 803
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3              = 804
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0              = 805
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1              = 806
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2              = 807
	DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3              = 808
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0               = 809
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1               = 810
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2               = 811
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3               = 812
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0               = 813
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1               = 814
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2               = 815
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3               = 816
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4               = 817
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5               = 818
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6               = 819
	DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7               = 820
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4               = 821
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5               = 822
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6               = 823
	DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE7               = 824
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L0                       = 825
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L1                       = 826
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L2                       = 827
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L3                       = 828
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L4                       = 829
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L5                       = 830
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L6                       = 831
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L7                       = 832
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L8                       = 833
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L9                       = 834
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L10                      = 835
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L11                      = 836
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L12                      = 837
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L13                      = 838
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L14                      = 839
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L15                      = 840
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L16                      = 841
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L17                      = 842
	DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL                    = 843
	DCGM_FI_DEV_NVSWITCH_FATAL_ERRORS                        = 856
	DCGM_FI_DEV_NVSWITCH_NON_FATAL_ERRORS                    = 857
	DCGM_FI_DEV_NVSWITCH_TEMPERATURE_CURRENT                 = 858
	DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SLOWDOWN          = 859
	DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SHUTDOWN          = 860
	DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX                       = 861
	DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX                       = 862
	DCGM_FI_DEV_NVSWITCH_PHYS_ID                             = 863
	DCGM_FI_DEV_NVSWITCH_RESET_REQUIRED                      = 864
	DCGM_FI_DEV_NVSWITCH_LINK_ID                             = 865
	DCGM_FI_DEV_NVSWITCH_PCIE_DOMAIN                         = 866
	DCGM_FI_DEV_NVSWITCH_PCIE_BUS                            = 867
	DCGM_FI_DEV_NVSWITCH_PCIE_DEVICE                         = 868
	DCGM_FI_DEV_NVSWITCH_PCIE_FUNCTION                       = 869
	DCGM_FI_DEV_NVSWITCH_LINK_STATUS                         = 870
	DCGM_FI_DEV_NVSWITCH_LINK_TYPE                           = 871
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN             = 872
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS                = 873
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE             = 874
	DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION           = 875
	DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_ID                 = 876
	DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_SID                = 877
	DCGM_FI_DEV_NVSWITCH_DEVICE_UUID                         = 878
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0                       = 879
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1                       = 880
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2                       = 881
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3                       = 882
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4                       = 883
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5                       = 884
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6                       = 885
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7                       = 886
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8                       = 887
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9                       = 888
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10                      = 889
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11                      = 890
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L12                      = 891
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L13                      = 892
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L14                      = 893
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L15                      = 894
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L16                      = 895
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L17                      = 896
	DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL                    = 897
	DCGM_FI_PROF_GR_ENGINE_ACTIVE                            = 1001
	DCGM_FI_PROF_SM_ACTIVE                                   = 1002
	DCGM_FI_PROF_SM_OCCUPANCY                                = 1003
	DCGM_FI_PROF_PIPE_TENSOR_ACTIVE                          = 1004
	DCGM_FI_PROF_DRAM_ACTIVE                                 = 1005
	DCGM_FI_PROF_PIPE_FP64_ACTIVE                            = 1006
	DCGM_FI_PROF_PIPE_FP32_ACTIVE                            = 1007
	DCGM_FI_PROF_PIPE_FP16_ACTIVE                            = 1008
	DCGM_FI_PROF_PCIE_TX_BYTES                               = 1009
	DCGM_FI_PROF_PCIE_RX_BYTES                               = 1010
	DCGM_FI_PROF_NVLINK_TX_BYTES                             = 1011
	DCGM_FI_PROF_NVLINK_RX_BYTES                             = 1012
	DCGM_FI_PROF_PIPE_TENSOR_IMMA_ACTIVE                     = 1013
	DCGM_FI_PROF_PIPE_TENSOR_HMMA_ACTIVE                     = 1014
	DCGM_FI_PROF_PIPE_TENSOR_DFMA_ACTIVE                     = 1015
	DCGM_FI_PROF_PIPE_INT_ACTIVE                             = 1016
	DCGM_FI_PROF_NVDEC0_ACTIVE                               = 1017
	DCGM_FI_PROF_NVDEC1_ACTIVE                               = 1018
	DCGM_FI_PROF_NVDEC2_ACTIVE                               = 1019
	DCGM_FI_PROF_NVDEC3_ACTIVE                               = 1020
	DCGM_FI_PROF_NVDEC4_ACTIVE                               = 1021
	DCGM_FI_PROF_NVDEC5_ACTIVE                               = 1022
	DCGM_FI_PROF_NVDEC6_ACTIVE                               = 1023
	DCGM_FI_PROF_NVDEC7_ACTIVE                               = 1024
	DCGM_FI_PROF_NVJPG0_ACTIVE                               = 1025
	DCGM_FI_PROF_NVJPG1_ACTIVE                               = 1026
	DCGM_FI_PROF_NVJPG2_ACTIVE                               = 1027
	DCGM_FI_PROF_NVJPG3_ACTIVE                               = 1028
	DCGM_FI_PROF_NVJPG4_ACTIVE                               = 1029
	DCGM_FI_PROF_NVJPG5_ACTIVE                               = 1030
	DCGM_FI_PROF_NVJPG6_ACTIVE                               = 1031
	DCGM_FI_PROF_NVJPG7_ACTIVE                               = 1032
	DCGM_FI_PROF_NVOFA0_ACTIVE                               = 1033
	DCGM_FI_PROF_NVOFA1_ACTIVE                               = 1034
	DCGM_FI_PROF_NVLINK_L0_TX_BYTES                          = 1040
	DCGM_FI_PROF_NVLINK_L0_RX_BYTES                          = 1041
	DCGM_FI_PROF_NVLINK_L1_TX_BYTES                          = 1042
	DCGM_FI_PROF_NVLINK_L1_RX_BYTES                          = 1043
	DCGM_FI_PROF_NVLINK_L2_TX_BYTES                          = 1044
	DCGM_FI_PROF_NVLINK_L2_RX_BYTES                          = 1045
	DCGM_FI_PROF_NVLINK_L3_TX_BYTES                          = 1046
	DCGM_FI_PROF_NVLINK_L3_RX_BYTES                          = 1047
	DCGM_FI_PROF_NVLINK_L4_TX_BYTES                          = 1048
	DCGM_FI_PROF_NVLINK_L4_RX_BYTES                          = 1049
	DCGM_FI_PROF_NVLINK_L5_TX_BYTES                          = 1050
	DCGM_FI_PROF_NVLINK_L5_RX_BYTES                          = 1051
	DCGM_FI_PROF_NVLINK_L6_TX_BYTES                          = 1052
	DCGM_FI_PROF_NVLINK_L6_RX_BYTES                          = 1053
	DCGM_FI_PROF_NVLINK_L7_TX_BYTES                          = 1054
	DCGM_FI_PROF_NVLINK_L7_RX_BYTES                          = 1055
	DCGM_FI_PROF_NVLINK_L8_TX_BYTES                          = 1056
	DCGM_FI_PROF_NVLINK_L8_RX_BYTES                          = 1057
	DCGM_FI_PROF_NVLINK_L9_TX_BYTES                          = 1058
	DCGM_FI_PROF_NVLINK_L9_RX_BYTES                          = 1059
	DCGM_FI_PROF_NVLINK_L10_TX_BYTES                         = 1060
	DCGM_FI_PROF_NVLINK_L10_RX_BYTES                         = 1061
	DCGM_FI_PROF_NVLINK_L11_TX_BYTES                         = 1062
	DCGM_FI_PROF_NVLINK_L11_RX_BYTES                         = 1063
	DCGM_FI_PROF_NVLINK_L12_TX_BYTES                         = 1064
	DCGM_FI_PROF_NVLINK_L12_RX_BYTES                         = 1065
	DCGM_FI_PROF_NVLINK_L13_TX_BYTES                         = 1066
	DCGM_FI_PROF_NVLINK_L13_RX_BYTES                         = 1067
	DCGM_FI_PROF_NVLINK_L14_TX_BYTES                         = 1068
	DCGM_FI_PROF_NVLINK_L14_RX_BYTES                         = 1069
	DCGM_FI_PROF_NVLINK_L15_TX_BYTES                         = 1070
	DCGM_FI_PROF_NVLINK_L15_RX_BYTES                         = 1071
	DCGM_FI_PROF_NVLINK_L16_TX_BYTES                         = 1072
	DCGM_FI_PROF_NVLINK_L16_RX_BYTES                         = 1073
	DCGM_FI_PROF_NVLINK_L17_TX_BYTES                         = 1074
	DCGM_FI_PROF_NVLINK_L17_RX_BYTES                         = 1075
	DCGM_FI_PROF_C2C_TX_ALL_BYTES                            = 1076
	DCGM_FI_PROF_C2C_TX_DATA_BYTES                           = 1077
	DCGM_FI_PROF_C2C_RX_ALL_BYTES                            = 1078
	DCGM_FI_PROF_C2C_RX_DATA_BYTES                           = 1079
	DCGM_FI_DEV_CPU_UTIL_TOTAL                               = 1100
	DCGM_FI_DEV_CPU_UTIL_USER                                = 1101
	DCGM_FI_DEV_CPU_UTIL_NICE                                = 1102
	DCGM_FI_DEV_CPU_UTIL_SYS                                 = 1103
	DCGM_FI_DEV_CPU_UTIL_IRQ                                 = 1104
	DCGM_FI_DEV_CPU_TEMP_CURRENT                             = 1110
	DCGM_FI_DEV_CPU_TEMP_WARNING                             = 1111
	DCGM_FI_DEV_CPU_TEMP_CRITICAL                            = 1112
	DCGM_FI_DEV_CPU_CLOCK_CURRENT                            = 1120
	DCGM_FI_DEV_CPU_POWER_UTIL_CURRENT                       = 1130
	DCGM_FI_DEV_CPU_POWER_LIMIT                              = 1131
	DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT                     = 1132
	DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT                    = 1133
	DCGM_FI_DEV_CPU_VENDOR                                   = 1140
	DCGM_FI_DEV_CPU_MODEL                                    = 1141
	DCGM_FI_DEV_NVLINK_COUNT_TX_PACKETS                      = 1200
	DCGM_FI_DEV_NVLINK_COUNT_TX_BYTES                        = 1201
	DCGM_FI_DEV_NVLINK_COUNT_RX_PACKETS                      = 1202
	DCGM_FI_DEV_NVLINK_COUNT_RX_BYTES                        = 1203
	DCGM_FI_DEV_NVLINK_COUNT_RX_MALFORMED_PACKET_ERRORS      = 1204
	DCGM_FI_DEV_NVLINK_COUNT_RX_BUFFER_OVERRUN_ERRORS        = 1205
	DCGM_FI_DEV_NVLINK_COUNT_RX_ERRORS                       = 1206
	DCGM_FI_DEV_NVLINK_COUNT_RX_REMOTE_ERRORS                = 1207
	DCGM_FI_DEV_NVLINK_COUNT_RX_GENERAL_ERRORS               = 1208
	DCGM_FI_DEV_NVLINK_COUNT_LOCAL_LINK_INTEGRITY_ERRORS     = 1209
	DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS                     = 1210
	DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_SUCCESSFUL_EVENTS = 1211
	DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_FAILED_EVENTS     = 1212
	DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_EVENTS            = 1213
	DCGM_FI_DEV_NVLINK_COUNT_RX_SYMBOL_ERRORS                = 1214
	DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER                      = 1215
	DCGM_FI_DEV_CONNECTX_HEALTH                              = 1300
	DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH              = 1301
	DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED              = 1302
	DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH              = 1303
	DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED              = 1304
	DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS              = 1305
	DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK                = 1306
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS            = 1307
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK              = 1308
	DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY          = 1309
	DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE                  = 1310
	DCGM_FI_MAX_FIELDS                                       = 1311

	DCGM_ST_OK                          = 0
	DCGM_ST_BADPARAM                    = -1
	DCGM_ST_GENERIC_ERROR               = -3
	DCGM_ST_MEMORY                      = -4
	DCGM_ST_NOT_CONFIGURED              = -5
	DCGM_ST_NOT_SUPPORTED               = -6
	DCGM_ST_INIT_ERROR                  = -7
	DCGM_ST_NVML_ERROR                  = -8
	DCGM_ST_PENDING                     = -9
	DCGM_ST_UNINITIALIZED               = -10
	DCGM_ST_TIMEOUT                     = -11
	DCGM_ST_VER_MISMATCH                = -12
	DCGM_ST_UNKNOWN_FIELD               = -13
	DCGM_ST_NO_DATA                     = -14
	DCGM_ST_STALE_DATA                  = -15
	DCGM_ST_NOT_WATCHED                 = -16
	DCGM_ST_NO_PERMISSION               = -17
	DCGM_ST_GPU_IS_LOST                 = -18
	DCGM_ST_RESET_REQUIRED              = -19
	DCGM_ST_FUNCTION_NOT_FOUND          = -20
	DCGM_ST_CONNECTION_NOT_VALID        = -21
	DCGM_ST_GPU_NOT_SUPPORTED           = -22
	DCGM_ST_GROUP_INCOMPATIBLE          = -23
	DCGM_ST_MAX_LIMIT                   = -24
	DCGM_ST_LIBRARY_NOT_FOUND           = -25
	DCGM_ST_DUPLICATE_KEY               = -26
	DCGM_ST_GPU_IN_SYNC_BOOST_GROUP     = -27
	DCGM_ST_GPU_NOT_IN_SYNC_BOOST_GROUP = -28
	DCGM_ST_REQUIRES_ROOT               = -29
	DCGM_ST_NVVS_ERROR                  = -30
	DCGM_ST_INSUFFICIENT_SIZE           = -31
	DCGM_ST_FIELD_UNSUPPORTED_BY_API    = -32
	DCGM_ST_MODULE_NOT_LOADED           = -33
	DCGM_ST_IN_USE                      = -34
	DCGM_ST_GROUP_IS_EMPTY              = -35
	DCGM_ST_PROFILING_NOT_SUPPORTED     = -36
	DCGM_ST_PROFILING_LIBRARY_ERROR     = -37
	DCGM_ST_PROFILING_MULTI_PASS        = -38
	DCGM_ST_DIAG_ALREADY_RUNNING        = -39
	DCGM_ST_DIAG_BAD_JSON               = -40
	DCGM_ST_DIAG_BAD_LAUNCH             = -41
	DCGM_ST_DIAG_UNUSED                 = -42
	DCGM_ST_DIAG_THRESHOLD_EXCEEDED     = -43
	DCGM_ST_INSUFFICIENT_DRIVER_VERSION = -44
	DCGM_ST_INSTANCE_NOT_FOUND          = -45
	DCGM_ST_COMPUTE_INSTANCE_NOT_FOUND  = -46
	DCGM_ST_CHILD_NOT_KILLED            = -47
	DCGM_ST_3RD_PARTY_LIBRARY_ERROR     = -48
	DCGM_ST_INSUFFICIENT_RESOURCES      = -49
	DCGM_ST_PLUGIN_EXCEPTION            = -50
	DCGM_ST_NVVS_ISOLATE_ERROR          = -51
	DCGM_ST_NVVS_BINARY_NOT_FOUND       = -52
	DCGM_ST_NVVS_KILLED                 = -53
	DCGM_ST_PAUSED                      = -54
	DCGM_ST_ALREADY_INITIALIZED         = -55
	DCGM_ST_NVML_NOT_LOADED             = -56
	DCGM_ST_NVML_DRIVER_TIMEOUT         = -57
	DCGM_ST_NVVS_NO_AVAILABLE_TEST      = -58
)

var DCGM_FI = map[string]Short{
	"DCGM_FT_BINARY":    Short('b'),
	"DCGM_FT_DOUBLE":    Short('d'),
	"DCGM_FT_INT64":     Short('i'),
	"DCGM_FT_STRING":    Short('s'),
	"DCGM_FT_TIMESTAMP": Short('t'),

	"DCGM_FI_UNKNOWN":                                          0,
	"DCGM_FI_DRIVER_VERSION":                                   1,
	"DCGM_FI_NVML_VERSION":                                     2,
	"DCGM_FI_PROCESS_NAME":                                     3,
	"DCGM_FI_DEV_COUNT":                                        4,
	"DCGM_FI_CUDA_DRIVER_VERSION":                              5,
	"DCGM_FI_DEV_NAME":                                         50,
	"DCGM_FI_DEV_BRAND":                                        51,
	"DCGM_FI_DEV_NVML_INDEX":                                   52,
	"DCGM_FI_DEV_SERIAL":                                       53,
	"DCGM_FI_DEV_UUID":                                         54,
	"DCGM_FI_DEV_MINOR_NUMBER":                                 55,
	"DCGM_FI_DEV_OEM_INFOROM_VER":                              56,
	"DCGM_FI_DEV_PCI_BUSID":                                    57,
	"DCGM_FI_DEV_PCI_COMBINED_ID":                              58,
	"DCGM_FI_DEV_PCI_SUBSYS_ID":                                59,
	"DCGM_FI_GPU_TOPOLOGY_PCI":                                 60,
	"DCGM_FI_GPU_TOPOLOGY_NVLINK":                              61,
	"DCGM_FI_GPU_TOPOLOGY_AFFINITY":                            62,
	"DCGM_FI_DEV_CUDA_COMPUTE_CAPABILITY":                      63,
	"DCGM_FI_DEV_COMPUTE_MODE":                                 65,
	"DCGM_FI_DEV_PERSISTENCE_MODE":                             66,
	"DCGM_FI_DEV_MIG_MODE":                                     67,
	"DCGM_FI_DEV_CUDA_VISIBLE_DEVICES_STR":                     68,
	"DCGM_FI_DEV_MIG_MAX_SLICES":                               69,
	"DCGM_FI_DEV_CPU_AFFINITY_0":                               70,
	"DCGM_FI_DEV_CPU_AFFINITY_1":                               71,
	"DCGM_FI_DEV_CPU_AFFINITY_2":                               72,
	"DCGM_FI_DEV_CPU_AFFINITY_3":                               73,
	"DCGM_FI_DEV_CC_MODE":                                      74,
	"DCGM_FI_DEV_MIG_ATTRIBUTES":                               75,
	"DCGM_FI_DEV_MIG_GI_INFO":                                  76,
	"DCGM_FI_DEV_MIG_CI_INFO":                                  77,
	"DCGM_FI_DEV_ECC_INFOROM_VER":                              80,
	"DCGM_FI_DEV_POWER_INFOROM_VER":                            81,
	"DCGM_FI_DEV_INFOROM_IMAGE_VER":                            82,
	"DCGM_FI_DEV_INFOROM_CONFIG_CHECK":                         83,
	"DCGM_FI_DEV_INFOROM_CONFIG_VALID":                         84,
	"DCGM_FI_DEV_VBIOS_VERSION":                                85,
	"DCGM_FI_DEV_MEM_AFFINITY_0":                               86,
	"DCGM_FI_DEV_MEM_AFFINITY_1":                               87,
	"DCGM_FI_DEV_MEM_AFFINITY_2":                               88,
	"DCGM_FI_DEV_MEM_AFFINITY_3":                               89,
	"DCGM_FI_DEV_BAR1_TOTAL":                                   90,
	"DCGM_FI_SYNC_BOOST":                                       91,
	"DCGM_FI_DEV_BAR1_USED":                                    92,
	"DCGM_FI_DEV_BAR1_FREE":                                    93,
	"DCGM_FI_DEV_GPM_SUPPORT":                                  94,
	"DCGM_FI_DEV_SM_CLOCK":                                     100,
	"DCGM_FI_DEV_MEM_CLOCK":                                    101,
	"DCGM_FI_DEV_VIDEO_CLOCK":                                  102,
	"DCGM_FI_DEV_APP_SM_CLOCK":                                 110,
	"DCGM_FI_DEV_APP_MEM_CLOCK":                                111,
	"DCGM_FI_DEV_CLOCKS_EVENT_REASONS":                         112,
	"DCGM_FI_DEV_CLOCK_THROTTLE_REASONS":                       DCGM_FI_DEV_CLOCKS_EVENT_REASONS,
	"DCGM_FI_DEV_MAX_SM_CLOCK":                                 113,
	"DCGM_FI_DEV_MAX_MEM_CLOCK":                                114,
	"DCGM_FI_DEV_MAX_VIDEO_CLOCK":                              115,
	"DCGM_FI_DEV_AUTOBOOST":                                    120,
	"DCGM_FI_DEV_SUPPORTED_CLOCKS":                             130,
	"DCGM_FI_DEV_MEMORY_TEMP":                                  140,
	"DCGM_FI_DEV_GPU_TEMP":                                     150,
	"DCGM_FI_DEV_MEM_MAX_OP_TEMP":                              151,
	"DCGM_FI_DEV_GPU_MAX_OP_TEMP":                              152,
	"DCGM_FI_DEV_GPU_TEMP_LIMIT":                               153,
	"DCGM_FI_DEV_POWER_USAGE":                                  155,
	"DCGM_FI_DEV_TOTAL_ENERGY_CONSUMPTION":                     156,
	"DCGM_FI_DEV_POWER_USAGE_INSTANT":                          157,
	"DCGM_FI_DEV_SLOWDOWN_TEMP":                                158,
	"DCGM_FI_DEV_SHUTDOWN_TEMP":                                159,
	"DCGM_FI_DEV_POWER_MGMT_LIMIT":                             160,
	"DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN":                         161,
	"DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX":                         162,
	"DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF":                         163,
	"DCGM_FI_DEV_ENFORCED_POWER_LIMIT":                         164,
	"DCGM_FI_DEV_REQUESTED_POWER_PROFILE_MASK":                 165,
	"DCGM_FI_DEV_ENFORCED_POWER_PROFILE_MASK":                  166,
	"DCGM_FI_DEV_VALID_POWER_PROFILE_MASK":                     167,
	"DCGM_FI_DEV_FABRIC_MANAGER_STATUS":                        170,
	"DCGM_FI_DEV_FABRIC_MANAGER_ERROR_CODE":                    171,
	"DCGM_FI_DEV_FABRIC_CLUSTER_UUID":                          172,
	"DCGM_FI_DEV_FABRIC_CLIQUE_ID":                             173,
	"DCGM_FI_DEV_PSTATE":                                       190,
	"DCGM_FI_DEV_FAN_SPEED":                                    191,
	"DCGM_FI_DEV_PCIE_TX_THROUGHPUT":                           200,
	"DCGM_FI_DEV_PCIE_RX_THROUGHPUT":                           201,
	"DCGM_FI_DEV_PCIE_REPLAY_COUNTER":                          202,
	"DCGM_FI_DEV_GPU_UTIL":                                     203,
	"DCGM_FI_DEV_MEM_COPY_UTIL":                                204,
	"DCGM_FI_DEV_ACCOUNTING_DATA":                              205,
	"DCGM_FI_DEV_ENC_UTIL":                                     206,
	"DCGM_FI_DEV_DEC_UTIL":                                     207,
	"DCGM_FI_DEV_XID_ERRORS":                                   230,
	"DCGM_FI_DEV_PCIE_MAX_LINK_GEN":                            235,
	"DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH":                          236,
	"DCGM_FI_DEV_PCIE_LINK_GEN":                                237,
	"DCGM_FI_DEV_PCIE_LINK_WIDTH":                              238,
	"DCGM_FI_DEV_POWER_VIOLATION":                              240,
	"DCGM_FI_DEV_THERMAL_VIOLATION":                            241,
	"DCGM_FI_DEV_SYNC_BOOST_VIOLATION":                         242,
	"DCGM_FI_DEV_BOARD_LIMIT_VIOLATION":                        243,
	"DCGM_FI_DEV_LOW_UTIL_VIOLATION":                           244,
	"DCGM_FI_DEV_RELIABILITY_VIOLATION":                        245,
	"DCGM_FI_DEV_TOTAL_APP_CLOCKS_VIOLATION":                   246,
	"DCGM_FI_DEV_TOTAL_BASE_CLOCKS_VIOLATION":                  247,
	"DCGM_FI_DEV_FB_TOTAL":                                     250,
	"DCGM_FI_DEV_FB_FREE":                                      251,
	"DCGM_FI_DEV_FB_USED":                                      252,
	"DCGM_FI_DEV_FB_RESERVED":                                  253,
	"DCGM_FI_DEV_FB_USED_PERCENT":                              254,
	"DCGM_FI_DEV_C2C_LINK_COUNT":                               285,
	"DCGM_FI_DEV_C2C_LINK_STATUS":                              286,
	"DCGM_FI_DEV_C2C_MAX_BANDWIDTH":                            287,
	"DCGM_FI_DEV_ECC_CURRENT":                                  300,
	"DCGM_FI_DEV_ECC_PENDING":                                  301,
	"DCGM_FI_DEV_ECC_SBE_VOL_TOTAL":                            310,
	"DCGM_FI_DEV_ECC_DBE_VOL_TOTAL":                            311,
	"DCGM_FI_DEV_ECC_SBE_AGG_TOTAL":                            312,
	"DCGM_FI_DEV_ECC_DBE_AGG_TOTAL":                            313,
	"DCGM_FI_DEV_ECC_SBE_VOL_L1":                               314,
	"DCGM_FI_DEV_ECC_DBE_VOL_L1":                               315,
	"DCGM_FI_DEV_ECC_SBE_VOL_L2":                               316,
	"DCGM_FI_DEV_ECC_DBE_VOL_L2":                               317,
	"DCGM_FI_DEV_ECC_SBE_VOL_DEV":                              318,
	"DCGM_FI_DEV_ECC_DBE_VOL_DEV":                              319,
	"DCGM_FI_DEV_ECC_SBE_VOL_REG":                              320,
	"DCGM_FI_DEV_ECC_DBE_VOL_REG":                              321,
	"DCGM_FI_DEV_ECC_SBE_VOL_TEX":                              322,
	"DCGM_FI_DEV_ECC_DBE_VOL_TEX":                              323,
	"DCGM_FI_DEV_ECC_SBE_AGG_L1":                               324,
	"DCGM_FI_DEV_ECC_DBE_AGG_L1":                               325,
	"DCGM_FI_DEV_ECC_SBE_AGG_L2":                               326,
	"DCGM_FI_DEV_ECC_DBE_AGG_L2":                               327,
	"DCGM_FI_DEV_ECC_SBE_AGG_DEV":                              328,
	"DCGM_FI_DEV_ECC_DBE_AGG_DEV":                              329,
	"DCGM_FI_DEV_ECC_SBE_AGG_REG":                              330,
	"DCGM_FI_DEV_ECC_DBE_AGG_REG":                              331,
	"DCGM_FI_DEV_ECC_SBE_AGG_TEX":                              332,
	"DCGM_FI_DEV_ECC_DBE_AGG_TEX":                              333,
	"DCGM_FI_DEV_ECC_SBE_VOL_SHM":                              334,
	"DCGM_FI_DEV_ECC_DBE_VOL_SHM":                              335,
	"DCGM_FI_DEV_ECC_SBE_VOL_CBU":                              336,
	"DCGM_FI_DEV_ECC_DBE_VOL_CBU":                              337,
	"DCGM_FI_DEV_ECC_SBE_AGG_SHM":                              338,
	"DCGM_FI_DEV_ECC_DBE_AGG_SHM":                              339,
	"DCGM_FI_DEV_ECC_SBE_AGG_CBU":                              340,
	"DCGM_FI_DEV_ECC_DBE_AGG_CBU":                              341,
	"DCGM_FI_DEV_ECC_SBE_VOL_SRM":                              342,
	"DCGM_FI_DEV_ECC_DBE_VOL_SRM":                              343,
	"DCGM_FI_DEV_ECC_SBE_AGG_SRM":                              344,
	"DCGM_FI_DEV_ECC_DBE_AGG_SRM":                              345,
	"DCGM_FI_DEV_DIAG_MEMORY_RESULT":                           350,
	"DCGM_FI_DEV_DIAG_DIAGNOSTIC_RESULT":                       351,
	"DCGM_FI_DEV_DIAG_PCIE_RESULT":                             352,
	"DCGM_FI_DEV_DIAG_TARGETED_STRESS_RESULT":                  353,
	"DCGM_FI_DEV_DIAG_TARGETED_POWER_RESULT":                   354,
	"DCGM_FI_DEV_DIAG_MEMORY_BANDWIDTH_RESULT":                 355,
	"DCGM_FI_DEV_DIAG_MEMTEST_RESULT":                          356,
	"DCGM_FI_DEV_DIAG_PULSE_TEST_RESULT":                       357,
	"DCGM_FI_DEV_DIAG_EUD_RESULT":                              358,
	"DCGM_FI_DEV_DIAG_CPU_EUD_RESULT":                          359,
	"DCGM_FI_DEV_DIAG_SOFTWARE_RESULT":                         360,
	"DCGM_FI_DEV_DIAG_NVBANDWIDTH_RESULT":                      361,
	"DCGM_FI_DEV_DIAG_STATUS":                                  362,
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_MAX":                   385,
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_HIGH":                  386,
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_PARTIAL":               387,
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_LOW":                   388,
	"DCGM_FI_DEV_BANKS_REMAP_ROWS_AVAIL_NONE":                  389,
	"DCGM_FI_DEV_RETIRED_SBE":                                  390,
	"DCGM_FI_DEV_RETIRED_DBE":                                  391,
	"DCGM_FI_DEV_RETIRED_PENDING":                              392,
	"DCGM_FI_DEV_UNCORRECTABLE_REMAPPED_ROWS":                  393,
	"DCGM_FI_DEV_CORRECTABLE_REMAPPED_ROWS":                    394,
	"DCGM_FI_DEV_ROW_REMAP_FAILURE":                            395,
	"DCGM_FI_DEV_ROW_REMAP_PENDING":                            396,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L0":               400,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L1":               401,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L2":               402,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L3":               403,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L4":               404,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L5":               405,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_TOTAL":            409,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L0":               410,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L1":               411,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L2":               412,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L3":               413,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L4":               414,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L5":               415,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_TOTAL":            419,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L0":                 420,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L1":                 421,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L2":                 422,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L3":                 423,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L4":                 424,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L5":                 425,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_TOTAL":              429,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L0":               430,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L1":               431,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L2":               432,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L3":               433,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L4":               434,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L5":               435,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_TOTAL":            439,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L0":                          440,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L1":                          441,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L2":                          442,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L3":                          443,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L4":                          444,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L5":                          445,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_TOTAL":                       449,
	"DCGM_FI_DEV_GPU_NVLINK_ERRORS":                            450,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L6":               451,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L7":               452,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L8":               453,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L9":               454,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L10":              455,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L11":              456,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L6":               457,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L7":               458,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L8":               459,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L9":               460,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L10":              461,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L11":              462,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L6":                 463,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L7":                 464,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L8":                 465,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L9":                 466,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L10":                467,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L11":                468,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L6":               469,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L7":               470,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L8":               471,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L9":               472,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L10":              473,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L11":              474,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L6":                          475,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L7":                          476,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L8":                          477,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L9":                          478,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L10":                         479,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L11":                         480,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L12":              406,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L13":              407,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L14":              408,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L15":              481,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L16":              482,
	"DCGM_FI_DEV_NVLINK_CRC_FLIT_ERROR_COUNT_L17":              483,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L12":              416,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L13":              417,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L14":              418,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L15":              484,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L16":              485,
	"DCGM_FI_DEV_NVLINK_CRC_DATA_ERROR_COUNT_L17":              486,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L12":                426,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L13":                427,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L14":                428,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L15":                487,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L16":                488,
	"DCGM_FI_DEV_NVLINK_REPLAY_ERROR_COUNT_L17":                489,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L12":              436,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L13":              437,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L14":              438,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L15":              491,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L16":              492,
	"DCGM_FI_DEV_NVLINK_RECOVERY_ERROR_COUNT_L17":              493,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L12":                         446,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L13":                         447,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L14":                         448,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L15":                         494,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L16":                         495,
	"DCGM_FI_DEV_NVLINK_BANDWIDTH_L17":                         496,
	"DCGM_FI_DEV_NVLINK_ERROR_DL_CRC":                          497,
	"DCGM_FI_DEV_NVLINK_ERROR_DL_RECOVERY":                     498,
	"DCGM_FI_DEV_NVLINK_ERROR_DL_REPLAY":                       499,
	"DCGM_FI_DEV_VIRTUAL_MODE":                                 500,
	"DCGM_FI_DEV_SUPPORTED_TYPE_INFO":                          501,
	"DCGM_FI_DEV_CREATABLE_VGPU_TYPE_IDS":                      502,
	"DCGM_FI_DEV_VGPU_INSTANCE_IDS":                            503,
	"DCGM_FI_DEV_VGPU_UTILIZATIONS":                            504,
	"DCGM_FI_DEV_VGPU_PER_PROCESS_UTILIZATION":                 505,
	"DCGM_FI_DEV_ENC_STATS":                                    506,
	"DCGM_FI_DEV_FBC_STATS":                                    507,
	"DCGM_FI_DEV_FBC_SESSIONS_INFO":                            508,
	"DCGM_FI_DEV_SUPPORTED_VGPU_TYPE_IDS":                      509,
	"DCGM_FI_DEV_VGPU_TYPE_INFO":                               510,
	"DCGM_FI_DEV_VGPU_TYPE_NAME":                               511,
	"DCGM_FI_DEV_VGPU_TYPE_CLASS":                              512,
	"DCGM_FI_DEV_VGPU_TYPE_LICENSE":                            513,
	"DCGM_FI_DEV_VGPU_VM_ID":                                   520,
	"DCGM_FI_DEV_VGPU_VM_NAME":                                 521,
	"DCGM_FI_DEV_VGPU_TYPE":                                    522,
	"DCGM_FI_DEV_VGPU_UUID":                                    523,
	"DCGM_FI_DEV_VGPU_DRIVER_VERSION":                          524,
	"DCGM_FI_DEV_VGPU_MEMORY_USAGE":                            525,
	"DCGM_FI_DEV_VGPU_LICENSE_STATUS":                          526,
	"DCGM_FI_DEV_VGPU_FRAME_RATE_LIMIT":                        527,
	"DCGM_FI_DEV_VGPU_ENC_STATS":                               528,
	"DCGM_FI_DEV_VGPU_ENC_SESSIONS_INFO":                       529,
	"DCGM_FI_DEV_VGPU_FBC_STATS":                               530,
	"DCGM_FI_DEV_VGPU_FBC_SESSIONS_INFO":                       531,
	"DCGM_FI_DEV_VGPU_INSTANCE_LICENSE_STATE":                  532,
	"DCGM_FI_DEV_VGPU_PCI_ID":                                  533,
	"DCGM_FI_DEV_VGPU_VM_GPU_INSTANCE_ID":                      534,
	"DCGM_FI_FIRST_VGPU_FIELD_ID":                              520,
	"DCGM_FI_LAST_VGPU_FIELD_ID":                               570,
	"DCGM_FI_DEV_PLATFORM_INFINIBAND_GUID":                     571,
	"DCGM_FI_DEV_PLATFORM_CHASSIS_SERIAL_NUMBER":               572,
	"DCGM_FI_DEV_PLATFORM_CHASSIS_SLOT_NUMBER":                 573,
	"DCGM_FI_DEV_PLATFORM_TRAY_INDEX":                          574,
	"DCGM_FI_DEV_PLATFORM_HOST_ID":                             575,
	"DCGM_FI_DEV_PLATFORM_PEER_TYPE":                           576,
	"DCGM_FI_DEV_PLATFORM_MODULE_ID":                           577,
	"DCGM_FI_FIRST_NVSWITCH_FIELD_ID":                          700,
	"DCGM_FI_DEV_NVSWITCH_VOLTAGE_MVOLT":                       701,
	"DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ":                        702,
	"DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_REV":                    703,
	"DCGM_FI_DEV_NVSWITCH_CURRENT_IDDQ_DVDD":                   704,
	"DCGM_FI_DEV_NVSWITCH_POWER_VDD":                           705,
	"DCGM_FI_DEV_NVSWITCH_POWER_DVDD":                          706,
	"DCGM_FI_DEV_NVSWITCH_POWER_HVDD":                          707,
	"DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_TX":                  780,
	"DCGM_FI_DEV_NVSWITCH_LINK_THROUGHPUT_RX":                  781,
	"DCGM_FI_DEV_NVSWITCH_LINK_FATAL_ERRORS":                   782,
	"DCGM_FI_DEV_NVSWITCH_LINK_NON_FATAL_ERRORS":               783,
	"DCGM_FI_DEV_NVSWITCH_LINK_REPLAY_ERRORS":                  784,
	"DCGM_FI_DEV_NVSWITCH_LINK_RECOVERY_ERRORS":                785,
	"DCGM_FI_DEV_NVSWITCH_LINK_FLIT_ERRORS":                    786,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS":                     787,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS":                     788,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC0":                789,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC1":                790,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC2":                791,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_LOW_VC3":                792,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC0":             793,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC1":             794,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC2":             795,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_MEDIUM_VC3":             796,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC0":               797,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC1":               798,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC2":               799,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_HIGH_VC3":               800,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC0":              801,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC1":              802,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC2":              803,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_PANIC_VC3":              804,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC0":              805,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC1":              806,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC2":              807,
	"DCGM_FI_DEV_NVSWITCH_LINK_LATENCY_COUNT_VC3":              808,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE0":               809,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE1":               810,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE2":               811,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE3":               812,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE0":               813,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE1":               814,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE2":               815,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE3":               816,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE4":               817,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE5":               818,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE6":               819,
	"DCGM_FI_DEV_NVSWITCH_LINK_CRC_ERRORS_LANE7":               820,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE4":               821,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE5":               822,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE6":               823,
	"DCGM_FI_DEV_NVSWITCH_LINK_ECC_ERRORS_LANE7":               824,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L0":                       825,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L1":                       826,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L2":                       827,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L3":                       828,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L4":                       829,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L5":                       830,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L6":                       831,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L7":                       832,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L8":                       833,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L9":                       834,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L10":                      835,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L11":                      836,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L12":                      837,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L13":                      838,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L14":                      839,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L15":                      840,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L16":                      841,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_L17":                      842,
	"DCGM_FI_DEV_NVLINK_TX_BANDWIDTH_TOTAL":                    843,
	"DCGM_FI_DEV_NVSWITCH_FATAL_ERRORS":                        856,
	"DCGM_FI_DEV_NVSWITCH_NON_FATAL_ERRORS":                    857,
	"DCGM_FI_DEV_NVSWITCH_TEMPERATURE_CURRENT":                 858,
	"DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SLOWDOWN":          859,
	"DCGM_FI_DEV_NVSWITCH_TEMPERATURE_LIMIT_SHUTDOWN":          860,
	"DCGM_FI_DEV_NVSWITCH_THROUGHPUT_TX":                       861,
	"DCGM_FI_DEV_NVSWITCH_THROUGHPUT_RX":                       862,
	"DCGM_FI_DEV_NVSWITCH_PHYS_ID":                             863,
	"DCGM_FI_DEV_NVSWITCH_RESET_REQUIRED":                      864,
	"DCGM_FI_DEV_NVSWITCH_LINK_ID":                             865,
	"DCGM_FI_DEV_NVSWITCH_PCIE_DOMAIN":                         866,
	"DCGM_FI_DEV_NVSWITCH_PCIE_BUS":                            867,
	"DCGM_FI_DEV_NVSWITCH_PCIE_DEVICE":                         868,
	"DCGM_FI_DEV_NVSWITCH_PCIE_FUNCTION":                       869,
	"DCGM_FI_DEV_NVSWITCH_LINK_STATUS":                         870,
	"DCGM_FI_DEV_NVSWITCH_LINK_TYPE":                           871,
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DOMAIN":             872,
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_BUS":                873,
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_DEVICE":             874,
	"DCGM_FI_DEV_NVSWITCH_LINK_REMOTE_PCIE_FUNCTION":           875,
	"DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_ID":                 876,
	"DCGM_FI_DEV_NVSWITCH_LINK_DEVICE_LINK_SID":                877,
	"DCGM_FI_DEV_NVSWITCH_DEVICE_UUID":                         878,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L0":                       879,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L1":                       880,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L2":                       881,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L3":                       882,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L4":                       883,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L5":                       884,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L6":                       885,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L7":                       886,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L8":                       887,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L9":                       888,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L10":                      889,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L11":                      890,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L12":                      891,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L13":                      892,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L14":                      893,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L15":                      894,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L16":                      895,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_L17":                      896,
	"DCGM_FI_DEV_NVLINK_RX_BANDWIDTH_TOTAL":                    897,
	"DCGM_FI_PROF_GR_ENGINE_ACTIVE":                            1001,
	"DCGM_FI_PROF_SM_ACTIVE":                                   1002,
	"DCGM_FI_PROF_SM_OCCUPANCY":                                1003,
	"DCGM_FI_PROF_PIPE_TENSOR_ACTIVE":                          1004,
	"DCGM_FI_PROF_DRAM_ACTIVE":                                 1005,
	"DCGM_FI_PROF_PIPE_FP64_ACTIVE":                            1006,
	"DCGM_FI_PROF_PIPE_FP32_ACTIVE":                            1007,
	"DCGM_FI_PROF_PIPE_FP16_ACTIVE":                            1008,
	"DCGM_FI_PROF_PCIE_TX_BYTES":                               1009,
	"DCGM_FI_PROF_PCIE_RX_BYTES":                               1010,
	"DCGM_FI_PROF_NVLINK_TX_BYTES":                             1011,
	"DCGM_FI_PROF_NVLINK_RX_BYTES":                             1012,
	"DCGM_FI_PROF_PIPE_TENSOR_IMMA_ACTIVE":                     1013,
	"DCGM_FI_PROF_PIPE_TENSOR_HMMA_ACTIVE":                     1014,
	"DCGM_FI_PROF_PIPE_TENSOR_DFMA_ACTIVE":                     1015,
	"DCGM_FI_PROF_PIPE_INT_ACTIVE":                             1016,
	"DCGM_FI_PROF_NVDEC0_ACTIVE":                               1017,
	"DCGM_FI_PROF_NVDEC1_ACTIVE":                               1018,
	"DCGM_FI_PROF_NVDEC2_ACTIVE":                               1019,
	"DCGM_FI_PROF_NVDEC3_ACTIVE":                               1020,
	"DCGM_FI_PROF_NVDEC4_ACTIVE":                               1021,
	"DCGM_FI_PROF_NVDEC5_ACTIVE":                               1022,
	"DCGM_FI_PROF_NVDEC6_ACTIVE":                               1023,
	"DCGM_FI_PROF_NVDEC7_ACTIVE":                               1024,
	"DCGM_FI_PROF_NVJPG0_ACTIVE":                               1025,
	"DCGM_FI_PROF_NVJPG1_ACTIVE":                               1026,
	"DCGM_FI_PROF_NVJPG2_ACTIVE":                               1027,
	"DCGM_FI_PROF_NVJPG3_ACTIVE":                               1028,
	"DCGM_FI_PROF_NVJPG4_ACTIVE":                               1029,
	"DCGM_FI_PROF_NVJPG5_ACTIVE":                               1030,
	"DCGM_FI_PROF_NVJPG6_ACTIVE":                               1031,
	"DCGM_FI_PROF_NVJPG7_ACTIVE":                               1032,
	"DCGM_FI_PROF_NVOFA0_ACTIVE":                               1033,
	"DCGM_FI_PROF_NVOFA1_ACTIVE":                               1034,
	"DCGM_FI_PROF_NVLINK_L0_TX_BYTES":                          1040,
	"DCGM_FI_PROF_NVLINK_L0_RX_BYTES":                          1041,
	"DCGM_FI_PROF_NVLINK_L1_TX_BYTES":                          1042,
	"DCGM_FI_PROF_NVLINK_L1_RX_BYTES":                          1043,
	"DCGM_FI_PROF_NVLINK_L2_TX_BYTES":                          1044,
	"DCGM_FI_PROF_NVLINK_L2_RX_BYTES":                          1045,
	"DCGM_FI_PROF_NVLINK_L3_TX_BYTES":                          1046,
	"DCGM_FI_PROF_NVLINK_L3_RX_BYTES":                          1047,
	"DCGM_FI_PROF_NVLINK_L4_TX_BYTES":                          1048,
	"DCGM_FI_PROF_NVLINK_L4_RX_BYTES":                          1049,
	"DCGM_FI_PROF_NVLINK_L5_TX_BYTES":                          1050,
	"DCGM_FI_PROF_NVLINK_L5_RX_BYTES":                          1051,
	"DCGM_FI_PROF_NVLINK_L6_TX_BYTES":                          1052,
	"DCGM_FI_PROF_NVLINK_L6_RX_BYTES":                          1053,
	"DCGM_FI_PROF_NVLINK_L7_TX_BYTES":                          1054,
	"DCGM_FI_PROF_NVLINK_L7_RX_BYTES":                          1055,
	"DCGM_FI_PROF_NVLINK_L8_TX_BYTES":                          1056,
	"DCGM_FI_PROF_NVLINK_L8_RX_BYTES":                          1057,
	"DCGM_FI_PROF_NVLINK_L9_TX_BYTES":                          1058,
	"DCGM_FI_PROF_NVLINK_L9_RX_BYTES":                          1059,
	"DCGM_FI_PROF_NVLINK_L10_TX_BYTES":                         1060,
	"DCGM_FI_PROF_NVLINK_L10_RX_BYTES":                         1061,
	"DCGM_FI_PROF_NVLINK_L11_TX_BYTES":                         1062,
	"DCGM_FI_PROF_NVLINK_L11_RX_BYTES":                         1063,
	"DCGM_FI_PROF_NVLINK_L12_TX_BYTES":                         1064,
	"DCGM_FI_PROF_NVLINK_L12_RX_BYTES":                         1065,
	"DCGM_FI_PROF_NVLINK_L13_TX_BYTES":                         1066,
	"DCGM_FI_PROF_NVLINK_L13_RX_BYTES":                         1067,
	"DCGM_FI_PROF_NVLINK_L14_TX_BYTES":                         1068,
	"DCGM_FI_PROF_NVLINK_L14_RX_BYTES":                         1069,
	"DCGM_FI_PROF_NVLINK_L15_TX_BYTES":                         1070,
	"DCGM_FI_PROF_NVLINK_L15_RX_BYTES":                         1071,
	"DCGM_FI_PROF_NVLINK_L16_TX_BYTES":                         1072,
	"DCGM_FI_PROF_NVLINK_L16_RX_BYTES":                         1073,
	"DCGM_FI_PROF_NVLINK_L17_TX_BYTES":                         1074,
	"DCGM_FI_PROF_NVLINK_L17_RX_BYTES":                         1075,
	"DCGM_FI_PROF_C2C_TX_ALL_BYTES":                            1076,
	"DCGM_FI_PROF_C2C_TX_DATA_BYTES":                           1077,
	"DCGM_FI_PROF_C2C_RX_ALL_BYTES":                            1078,
	"DCGM_FI_PROF_C2C_RX_DATA_BYTES":                           1079,
	"DCGM_FI_DEV_CPU_UTIL_TOTAL":                               1100,
	"DCGM_FI_DEV_CPU_UTIL_USER":                                1101,
	"DCGM_FI_DEV_CPU_UTIL_NICE":                                1102,
	"DCGM_FI_DEV_CPU_UTIL_SYS":                                 1103,
	"DCGM_FI_DEV_CPU_UTIL_IRQ":                                 1104,
	"DCGM_FI_DEV_CPU_TEMP_CURRENT":                             1110,
	"DCGM_FI_DEV_CPU_TEMP_WARNING":                             1111,
	"DCGM_FI_DEV_CPU_TEMP_CRITICAL":                            1112,
	"DCGM_FI_DEV_CPU_CLOCK_CURRENT":                            1120,
	"DCGM_FI_DEV_CPU_POWER_UTIL_CURRENT":                       1130,
	"DCGM_FI_DEV_CPU_POWER_LIMIT":                              1131,
	"DCGM_FI_DEV_SYSIO_POWER_UTIL_CURRENT":                     1132,
	"DCGM_FI_DEV_MODULE_POWER_UTIL_CURRENT":                    1133,
	"DCGM_FI_DEV_CPU_VENDOR":                                   1140,
	"DCGM_FI_DEV_CPU_MODEL":                                    1141,
	"DCGM_FI_DEV_NVLINK_COUNT_TX_PACKETS":                      1200,
	"DCGM_FI_DEV_NVLINK_COUNT_TX_BYTES":                        1201,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_PACKETS":                      1202,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_BYTES":                        1203,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_MALFORMED_PACKET_ERRORS":      1204,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_BUFFER_OVERRUN_ERRORS":        1205,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_ERRORS":                       1206,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_REMOTE_ERRORS":                1207,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_GENERAL_ERRORS":               1208,
	"DCGM_FI_DEV_NVLINK_COUNT_LOCAL_LINK_INTEGRITY_ERRORS":     1209,
	"DCGM_FI_DEV_NVLINK_COUNT_TX_DISCARDS":                     1210,
	"DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_SUCCESSFUL_EVENTS": 1211,
	"DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_FAILED_EVENTS":     1212,
	"DCGM_FI_DEV_NVLINK_COUNT_LINK_RECOVERY_EVENTS":            1213,
	"DCGM_FI_DEV_NVLINK_COUNT_RX_SYMBOL_ERRORS":                1214,
	"DCGM_FI_DEV_NVLINK_COUNT_SYMBOL_BER":                      1215,
	"DCGM_FI_DEV_CONNECTX_HEALTH":                              1300,
	"DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_WIDTH":              1301,
	"DCGM_FI_DEV_CONNECTX_ACTIVE_PCIE_LINK_SPEED":              1302,
	"DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_WIDTH":              1303,
	"DCGM_FI_DEV_CONNECTX_EXPECT_PCIE_LINK_SPEED":              1304,
	"DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_STATUS":              1305,
	"DCGM_FI_DEV_CONNECTX_CORRECTABLE_ERR_MASK":                1306,
	"DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_STATUS":            1307,
	"DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_MASK":              1308,
	"DCGM_FI_DEV_CONNECTX_UNCORRECTABLE_ERR_SEVERITY":          1309,
	"DCGM_FI_DEV_CONNECTX_DEVICE_TEMPERATURE":                  1310,
	"DCGM_FI_MAX_FIELDS":                                       1311,
}

var OLD_DCGM_FI = map[string]Short{
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

const (
	DCGM_FV_FLAG_LIVE_DATA = uint(0x00000001)
)

type HealthSystem uint

const (
	DCGM_HEALTH_WATCH_PCIE              HealthSystem = 0x1
	DCGM_HEALTH_WATCH_NVLINK            HealthSystem = 0x2
	DCGM_HEALTH_WATCH_PMU               HealthSystem = 0x4
	DCGM_HEALTH_WATCH_MCU               HealthSystem = 0x8
	DCGM_HEALTH_WATCH_MEM               HealthSystem = 0x10
	DCGM_HEALTH_WATCH_SM                HealthSystem = 0x20
	DCGM_HEALTH_WATCH_INFOROM           HealthSystem = 0x40
	DCGM_HEALTH_WATCH_THERMAL           HealthSystem = 0x80
	DCGM_HEALTH_WATCH_POWER             HealthSystem = 0x100
	DCGM_HEALTH_WATCH_DRIVER            HealthSystem = 0x200
	DCGM_HEALTH_WATCH_NVSWITCH_NONFATAL HealthSystem = 0x400
	DCGM_HEALTH_WATCH_NVSWITCH_FATAL    HealthSystem = 0x800
	DCGM_HEALTH_WATCH_ALL               HealthSystem = 0xFFFFFFFF
)

type HealthResult uint

const (
	DCGM_HEALTH_RESULT_PASS HealthResult = 0  // All results within this system are reporting normal
	DCGM_HEALTH_RESULT_WARN HealthResult = 10 // A warning has been issued, refer to the response for more information
	DCGM_HEALTH_RESULT_FAIL HealthResult = 20 // A failure has been issued, refer to the response for more information
)

// HealthCheckErrorCode error codes for passive and active health checks.
type HealthCheckErrorCode uint

const (
	DCGM_FR_OK                              HealthCheckErrorCode = 0   // 0 No error
	DCGM_FR_UNKNOWN                         HealthCheckErrorCode = 1   // 1 Unknown error code
	DCGM_FR_UNRECOGNIZED                    HealthCheckErrorCode = 2   // 2 Unrecognized error code
	DCGM_FR_PCI_REPLAY_RATE                 HealthCheckErrorCode = 3   // 3 Unacceptable rate of PCI errors
	DCGM_FR_VOLATILE_DBE_DETECTED           HealthCheckErrorCode = 4   // 4 Uncorrectable volatile double bit error
	DCGM_FR_VOLATILE_SBE_DETECTED           HealthCheckErrorCode = 5   // 5 Unacceptable rate of volatile single bit errors
	DCGM_FR_PENDING_PAGE_RETIREMENTS        HealthCheckErrorCode = 6   // 6 Pending page retirements detected
	DCGM_FR_RETIRED_PAGES_LIMIT             HealthCheckErrorCode = 7   // 7 Unacceptable total page retirements detected
	DCGM_FR_RETIRED_PAGES_DBE_LIMIT         HealthCheckErrorCode = 8   // 8 Unacceptable total page retirements due to uncorrectable errors
	DCGM_FR_CORRUPT_INFOROM                 HealthCheckErrorCode = 9   // 9 Corrupt inforom found
	DCGM_FR_CLOCK_THROTTLE_THERMAL          HealthCheckErrorCode = 10  // 10 Clocks being throttled due to overheating
	DCGM_FR_POWER_UNREADABLE                HealthCheckErrorCode = 11  // 11 Cannot get a reading for power from NVML
	DCGM_FR_CLOCK_THROTTLE_POWER            HealthCheckErrorCode = 12  // 12 Clock being throttled due to power restrictions
	DCGM_FR_NVLINK_ERROR_THRESHOLD          HealthCheckErrorCode = 13  // 13 Unacceptable rate of NVLink errors
	DCGM_FR_NVLINK_DOWN                     HealthCheckErrorCode = 14  // 14 NVLink is down
	DCGM_FR_NVSWITCH_FATAL_ERROR            HealthCheckErrorCode = 15  // 15 Fatal errors on the NVSwitch
	DCGM_FR_NVSWITCH_NON_FATAL_ERROR        HealthCheckErrorCode = 16  // 16 Non-fatal errors on the NVSwitch
	DCGM_FR_NVSWITCH_DOWN                   HealthCheckErrorCode = 17  // 17 NVSwitch is down - NOT USED: DEPRECATED
	DCGM_FR_NO_ACCESS_TO_FILE               HealthCheckErrorCode = 18  // 18 Cannot access a file
	DCGM_FR_NVML_API                        HealthCheckErrorCode = 19  // 19 Error occurred on an NVML API - NOT USED: DEPRECATED
	DCGM_FR_DEVICE_COUNT_MISMATCH           HealthCheckErrorCode = 20  // 20 Disagreement in GPU count between /dev and NVML
	DCGM_FR_BAD_PARAMETER                   HealthCheckErrorCode = 21  // 21 Bad parameter passed to API
	DCGM_FR_CANNOT_OPEN_LIB                 HealthCheckErrorCode = 22  // 22 Cannot open a library that must be accessed
	DCGM_FR_DENYLISTED_DRIVER               HealthCheckErrorCode = 23  // 23 A driver on the denylist (nouveau) is active
	DCGM_FR_NVML_LIB_BAD                    HealthCheckErrorCode = 24  // 24 NVML library is missing expected functions - NOT USED: DEPRECATED
	DCGM_FR_GRAPHICS_PROCESSES              HealthCheckErrorCode = 25  // 25 Graphics processes are active on this GPU
	DCGM_FR_HOSTENGINE_CONN                 HealthCheckErrorCode = 26  // 26 Bad connection to nv-hostengine - NOT USED: DEPRECATED
	DCGM_FR_FIELD_QUERY                     HealthCheckErrorCode = 27  // 27 Error querying a field from DCGM
	DCGM_FR_BAD_CUDA_ENV                    HealthCheckErrorCode = 28  // 28 The environment has variables that hurt CUDA
	DCGM_FR_PERSISTENCE_MODE                HealthCheckErrorCode = 29  // 29 Persistence mode is disabled
	DCGM_FR_LOW_BANDWIDTH                   HealthCheckErrorCode = 30  // 30 The bandwidth is unacceptably low
	DCGM_FR_HIGH_LATENCY                    HealthCheckErrorCode = 31  // 31 Latency is too high
	DCGM_FR_CANNOT_GET_FIELD_TAG            HealthCheckErrorCode = 32  // 32 Cannot find a tag for a field
	DCGM_FR_FIELD_VIOLATION                 HealthCheckErrorCode = 33  // 33 The value for the specified error field is above 0
	DCGM_FR_FIELD_THRESHOLD                 HealthCheckErrorCode = 34  // 34 The value for the specified field is above the threshold
	DCGM_FR_FIELD_VIOLATION_DBL             HealthCheckErrorCode = 35  // 35 The value for the specified error field is above 0
	DCGM_FR_FIELD_THRESHOLD_DBL             HealthCheckErrorCode = 36  // 36 The value for the specified field is above the threshold
	DCGM_FR_UNSUPPORTED_FIELD_TYPE          HealthCheckErrorCode = 37  // 37 Field type cannot be supported
	DCGM_FR_FIELD_THRESHOLD_TS              HealthCheckErrorCode = 38  // 38 The value for the specified field is above the threshold
	DCGM_FR_FIELD_THRESHOLD_TS_DBL          HealthCheckErrorCode = 39  // 39 The value for the specified field is above the threshold
	DCGM_FR_THERMAL_VIOLATIONS              HealthCheckErrorCode = 40  // 40 Thermal violations detected
	DCGM_FR_THERMAL_VIOLATIONS_TS           HealthCheckErrorCode = 41  // 41 Thermal violations detected with a timestamp
	DCGM_FR_TEMP_VIOLATION                  HealthCheckErrorCode = 42  // 42 Temperature is too high
	DCGM_FR_THROTTLING_VIOLATION            HealthCheckErrorCode = 43  // 43 Non-benign clock throttling is occurring
	DCGM_FR_INTERNAL                        HealthCheckErrorCode = 44  // 44 An internal error was detected
	DCGM_FR_PCIE_GENERATION                 HealthCheckErrorCode = 45  // 45 PCIe generation is too low
	DCGM_FR_PCIE_WIDTH                      HealthCheckErrorCode = 46  // 46 PCIe width is too low
	DCGM_FR_ABORTED                         HealthCheckErrorCode = 47  // 47 Test was aborted by a user signal
	DCGM_FR_TEST_DISABLED                   HealthCheckErrorCode = 48  // 48 This test is disabled for this GPU
	DCGM_FR_CANNOT_GET_STAT                 HealthCheckErrorCode = 49  // 49 Cannot get telemetry for a needed value
	DCGM_FR_STRESS_LEVEL                    HealthCheckErrorCode = 50  // 50 Stress level is too low (bad performance)
	DCGM_FR_CUDA_API                        HealthCheckErrorCode = 51  // 51 Error calling the specified CUDA API
	DCGM_FR_FAULTY_MEMORY                   HealthCheckErrorCode = 52  // 52 Faulty memory detected on this GPU
	DCGM_FR_CANNOT_SET_WATCHES              HealthCheckErrorCode = 53  // 53 Unable to set field watches in DCGM - NOT USED: DEPRECATED
	DCGM_FR_CUDA_UNBOUND                    HealthCheckErrorCode = 54  // 54 CUDA context is no longer bound
	DCGM_FR_ECC_DISABLED                    HealthCheckErrorCode = 55  // 55 ECC memory is disabled right now
	DCGM_FR_MEMORY_ALLOC                    HealthCheckErrorCode = 56  // 56 Cannot allocate memory on the GPU
	DCGM_FR_CUDA_DBE                        HealthCheckErrorCode = 57  // 57 CUDA detected unrecovable double-bit error
	DCGM_FR_MEMORY_MISMATCH                 HealthCheckErrorCode = 58  // 58 Memory error detected
	DCGM_FR_CUDA_DEVICE                     HealthCheckErrorCode = 59  // 59 No CUDA device discoverable for existing GPU
	DCGM_FR_ECC_UNSUPPORTED                 HealthCheckErrorCode = 60  // 60 ECC memory is unsupported by this SKU
	DCGM_FR_ECC_PENDING                     HealthCheckErrorCode = 61  // 61 ECC memory is in a pending state - NOT USED: DEPRECATED
	DCGM_FR_MEMORY_BANDWIDTH                HealthCheckErrorCode = 62  // 62 Memory bandwidth is too low
	DCGM_FR_TARGET_POWER                    HealthCheckErrorCode = 63  // 63 Cannot hit the target power draw
	DCGM_FR_API_FAIL                        HealthCheckErrorCode = 64  // 64 The specified API call failed
	DCGM_FR_API_FAIL_GPU                    HealthCheckErrorCode = 65  // 65 The specified API call failed for the specified GPU
	DCGM_FR_CUDA_CONTEXT                    HealthCheckErrorCode = 66  // 66 Cannot create a CUDA context on this GPU
	DCGM_FR_DCGM_API                        HealthCheckErrorCode = 67  // 67 DCGM API failure
	DCGM_FR_CONCURRENT_GPUS                 HealthCheckErrorCode = 68  // 68 Need multiple GPUs to run this test
	DCGM_FR_TOO_MANY_ERRORS                 HealthCheckErrorCode = 69  // 69 More errors than fit in the return struct - NOT USED: DEPRECATED
	DCGM_FR_NVLINK_CRC_ERROR_THRESHOLD      HealthCheckErrorCode = 70  // 70 More than 100 CRC errors are happening per second
	DCGM_FR_NVLINK_ERROR_CRITICAL           HealthCheckErrorCode = 71  // 71 NVLink error for a field that should always be 0
	DCGM_FR_ENFORCED_POWER_LIMIT            HealthCheckErrorCode = 72  // 72 The enforced power limit is too low to hit the target
	DCGM_FR_MEMORY_ALLOC_HOST               HealthCheckErrorCode = 73  // 73 Cannot allocate memory on the host
	DCGM_FR_GPU_OP_MODE                     HealthCheckErrorCode = 74  // 74 Bad GPU operating mode for running plugin - NOT USED: DEPRECATED
	DCGM_FR_NO_MEMORY_CLOCKS                HealthCheckErrorCode = 75  // 75 No memory clocks with the needed MHz found - NOT USED: DEPRECATED
	DCGM_FR_NO_GRAPHICS_CLOCKS              HealthCheckErrorCode = 76  // 76 No graphics clocks with the needed MHz found - NOT USED: DEPRECATED
	DCGM_FR_HAD_TO_RESTORE_STATE            HealthCheckErrorCode = 77  // 77 Note that we had to restore a GPU's state
	DCGM_FR_L1TAG_UNSUPPORTED               HealthCheckErrorCode = 78  // 78 L1TAG test is unsupported by this SKU
	DCGM_FR_L1TAG_MISCOMPARE                HealthCheckErrorCode = 79  // 79 L1TAG test failed on a miscompare
	DCGM_FR_ROW_REMAP_FAILURE               HealthCheckErrorCode = 80  // 80 Row remapping failed (Ampere or newer GPUs)
	DCGM_FR_UNCONTAINED_ERROR               HealthCheckErrorCode = 81  // 81 Uncontained error - XID 95
	DCGM_FR_EMPTY_GPU_LIST                  HealthCheckErrorCode = 82  // 82 No GPU information given to plugin
	DCGM_FR_DBE_PENDING_PAGE_RETIREMENTS    HealthCheckErrorCode = 83  // 83 Pending page retirements due to a DBE
	DCGM_FR_UNCORRECTABLE_ROW_REMAP         HealthCheckErrorCode = 84  // 84 Uncorrectable row remapping
	DCGM_FR_PENDING_ROW_REMAP               HealthCheckErrorCode = 85  // 85 Row remapping is pending
	DCGM_FR_BROKEN_P2P_MEMORY_DEVICE        HealthCheckErrorCode = 86  // 86 P2P copy test detected an error writing to this GPU
	DCGM_FR_BROKEN_P2P_WRITER_DEVICE        HealthCheckErrorCode = 87  // 87 P2P copy test detected an error writing from this GPU
	DCGM_FR_NVSWITCH_NVLINK_DOWN            HealthCheckErrorCode = 88  // 88 An NvLink is down for the specified NVSwitch - NOT USED: DEPRECATED
	DCGM_FR_EUD_BINARY_PERMISSIONS          HealthCheckErrorCode = 89  // 89 EUD binary permissions are incorrect
	DCGM_FR_EUD_NON_ROOT_USER               HealthCheckErrorCode = 90  // 90 EUD plugin is not running as root
	DCGM_FR_EUD_SPAWN_FAILURE               HealthCheckErrorCode = 91  // 91 EUD plugin failed to spawn the EUD binary
	DCGM_FR_EUD_TIMEOUT                     HealthCheckErrorCode = 92  // 92 EUD plugin timed out
	DCGM_FR_EUD_ZOMBIE                      HealthCheckErrorCode = 93  // 93 EUD process remains running after the plugin considers it finished
	DCGM_FR_EUD_NON_ZERO_EXIT_CODE          HealthCheckErrorCode = 94  // 94 EUD process exited with a non-zero exit code
	DCGM_FR_EUD_TEST_FAILED                 HealthCheckErrorCode = 95  // 95 EUD test failed
	DCGM_FR_FILE_CREATE_PERMISSIONS         HealthCheckErrorCode = 96  // 96 We cannot create a file in this directory.
	DCGM_FR_PAUSE_RESUME_FAILED             HealthCheckErrorCode = 97  // 97 Pause/Resume failed
	DCGM_FR_PCIE_H_REPLAY_VIOLATION         HealthCheckErrorCode = 98  // 98 PCIe test caught correctable errors
	DCGM_FR_GPU_EXPECTED_NVLINKS_UP         HealthCheckErrorCode = 99  // 99 Expected nvlinks up per gpu
	DCGM_FR_NVSWITCH_EXPECTED_NVLINKS_UP    HealthCheckErrorCode = 100 // 100 Expected nvlinks up per nvswitch
	DCGM_FR_XID_ERROR                       HealthCheckErrorCode = 101 // 101 XID error detected
	DCGM_FR_SBE_VIOLATION                   HealthCheckErrorCode = 102 // 102 Single bit error detected
	DCGM_FR_DBE_VIOLATION                   HealthCheckErrorCode = 103 // 103 Double bit error detected
	DCGM_FR_PCIE_REPLAY_VIOLATION           HealthCheckErrorCode = 104 // 104 PCIe replay errors detected
	DCGM_FR_SBE_THRESHOLD_VIOLATION         HealthCheckErrorCode = 105 // 105 SBE threshold violated
	DCGM_FR_DBE_THRESHOLD_VIOLATION         HealthCheckErrorCode = 106 // 106 DBE threshold violated
	DCGM_FR_PCIE_REPLAY_THRESHOLD_VIOLATION HealthCheckErrorCode = 107 // 107 PCIE replay count violated
	DCGM_FR_CUDA_FM_NOT_INITIALIZED         HealthCheckErrorCode = 108 // 108 The fabricmanager is not initialized
	DCGM_FR_SXID_ERROR                      HealthCheckErrorCode = 109 // 109 NvSwitch fatal error detected
	DCGM_FR_GFLOPS_THRESHOLD_VIOLATION      HealthCheckErrorCode = 110 // 110 GPU GFLOPs threshold violated
	DCGM_FR_NAN_VALUE                       HealthCheckErrorCode = 111 // 111 NaN value detected on this GPU
	DCGM_FR_FABRIC_MANAGER_TRAINING_ERROR   HealthCheckErrorCode = 112 // 112 Fabric Manager did not finish training
	DCGM_FR_BROKEN_P2P_PCIE_MEMORY_DEVICE   HealthCheckErrorCode = 113 // 113 P2P copy test detected an error writing to this GPU over PCIE
	DCGM_FR_BROKEN_P2P_PCIE_WRITER_DEVICE   HealthCheckErrorCode = 114 // 114 P2P copy test detected an error writing from this GPU over PCIE
	DCGM_FR_BROKEN_P2P_NVLINK_MEMORY_DEVICE HealthCheckErrorCode = 115 // 115 P2P copy test detected an error writing to this GPU over NVLink
	DCGM_FR_BROKEN_P2P_NVLINK_WRITER_DEVICE HealthCheckErrorCode = 116 // 116 P2P copy test detected an error writing from this GPU over NVLink
	DCGM_FR_ERROR_SENTINEL                  HealthCheckErrorCode = 117 //!< 117 MUST BE THE LAST ERROR CODE
)
