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
	"runtime/cgo"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestSelectedFieldValueCallbackResolvesHandle(t *testing.T) {
	cb := &callback{}
	callbackHandle := cgo.NewHandle(cb)
	defer callbackHandle.Delete()

	const entityID = 17
	cvalues := makeTestCFields(2)

	result := go_dcgmFieldValueEntityEnumeration(
		1,
		entityID,
		&cvalues[0],
		2,
		unsafe.Pointer(&callbackHandle),
	)

	require.Zero(t, result)
	require.Len(t, cb.Values, len(cvalues))
	require.Equal(t, FE_GPU, cb.Values[0].EntityGroupId)
	require.Equal(t, uint(entityID), cb.Values[0].EntityID)
	require.Equal(t, Short(0), cb.Values[0].FieldID)
	require.Equal(t, Short(1), cb.Values[1].FieldID)
}

func TestSelectedFieldValueCallbackRejectsUnexpectedHandleValue(t *testing.T) {
	callbackHandle := cgo.NewHandle("not a callback")
	defer callbackHandle.Delete()

	cb, ok := callbackFromUserData(unsafe.Pointer(&callbackHandle))

	require.False(t, ok)
	require.Nil(t, cb)
}

func TestSelectedFieldValueCallbackRejectsZeroHandle(t *testing.T) {
	var callbackHandle cgo.Handle

	cb, ok := callbackFromUserData(unsafe.Pointer(&callbackHandle))

	require.False(t, ok)
	require.Nil(t, cb)
}
