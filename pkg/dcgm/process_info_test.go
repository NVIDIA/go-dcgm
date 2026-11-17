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

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPIDInfoLegacyNativeBoundary exercises the production native call without
// initializing DCGM or starting a hostengine. Race builds also check the Go
// pointer passed across this boundary against its allocation size.
func TestPIDInfoLegacyNativeBoundary(t *testing.T) {
	if dcgmInitCounter != 0 {
		t.Skip("requires an uninitialized DCGM library")
	}
	oldHandle, oldLibHandle := handle, dcgmLibHandle
	handle = dcgmHandle{}
	t.Cleanup(func() {
		handle, dcgmLibHandle = oldHandle, oldLibHandle
	})
	if err := (cgoAdapter{}).loadLibrary(); err != nil {
		t.Skipf("DCGM library not installed: %v", err)
	}
	t.Cleanup(func() { cgoAdapter{}.closeLibrary() })

	got, err := getProcessInfo(GroupHandle{}, uint(os.Getpid()))

	require.Nil(t, got)
	var nativeErr *Error
	require.ErrorAs(t, err, &nativeErr)
	assert.Equal(t, DCGM_ST_BADPARAM, int(nativeErr.Code))
}
