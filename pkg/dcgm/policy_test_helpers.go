//go:build linux && cgo

/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

package dcgm

/*
#include "dcgm_structs.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// fireFakePolicyCallback invokes ViolationRegistration with a synthetic
// dcgmPolicyCallbackResponse_t whose `condition` field matches the
// supplied key. Valid keys: "dbe", "pcie", "maxrtpg", "thermal",
// "power", "nvlink", "xid". Unknown keys panic.
//
// Test-only helper. Go does not allow `import "C"` inside _test.go
// files, so the cgo struct construction lives here. Returns the int
// returned by ViolationRegistration.
func fireFakePolicyCallback(condition string) int {
	var resp C.dcgmPolicyCallbackResponse_t
	switch condition {
	case "dbe":
		resp.condition = C.DCGM_POLICY_COND_DBE
	case "pcie":
		resp.condition = C.DCGM_POLICY_COND_PCI
	case "maxrtpg":
		resp.condition = C.DCGM_POLICY_COND_MAX_PAGES_RETIRED
	case "thermal":
		resp.condition = C.DCGM_POLICY_COND_THERMAL
	case "power":
		resp.condition = C.DCGM_POLICY_COND_POWER
	case "nvlink":
		resp.condition = C.DCGM_POLICY_COND_NVLINK
	case "xid":
		resp.condition = C.DCGM_POLICY_COND_XID
	default:
		panic("fireFakePolicyCallback: unknown condition " + condition)
	}
	rc := ViolationRegistration(unsafe.Pointer(&resp))
	runtime.KeepAlive(&resp)
	return rc
}
