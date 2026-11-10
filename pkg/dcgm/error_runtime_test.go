//go:build linux && cgo

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
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSelectedErrorMeta(t *testing.T) {
	cleanup, err := Init(Embedded)
	require.NoError(t, err)
	defer cleanup()

	meta := GetErrorMeta(DCGM_FR_PCI_REPLAY_RATE)
	require.NotNil(t, meta)
	require.Equal(t, DCGM_FR_PCI_REPLAY_RATE, meta.ErrorID)
	require.NotEmpty(t, meta.MessageFormat)
	require.NotEmpty(t, meta.Suggestion)
	// Keep this known entry exact so a runtime DCGM metadata change is detected.
	require.Equal(t, DCGM_ERROR_ISOLATE, meta.Severity)
	require.Equal(t, DCGM_FR_EC_HARDWARE_PCIE, meta.Category)

	require.Nil(t, GetErrorMeta(DCGM_FR_ERROR_SENTINEL))
	require.Nil(t, GetErrorMeta(HealthCheckErrorCode(math.MaxUint)))
}

// TestNativeFailuresRetainStatus checks native status and text for a zero session handle.
// Loading the library without initializing it makes the failure independent of GPU availability.
func TestNativeFailuresRetainStatus(t *testing.T) {
	require.Zero(t, dcgmInitCounter)
	previousHandle := handle
	handle = dcgmHandle{}
	t.Cleanup(func() { handle = previousHandle })
	require.NoError(t, cgoAdapter{}.loadLibrary())
	t.Cleanup(func() { cgoAdapter{}.closeLibrary() })

	for _, tt := range []struct {
		name   string
		prefix string
		call   func() error
	}{
		{"device count", "error getting devices count: ", func() error { _, err := GetAllDeviceCount(); return err }},
		{"entity enumeration", "error getting entity count: ", func() error { _, err := GetEntityGroupEntities(FE_GPU); return err }},
		{"hierarchy", "error retrieving DCGM MIG hierarchy: ", func() error { _, err := GetGPUInstanceHierarchy(); return err }},
		{"group create", "error creating group: ", func() error { _, err := CreateGroup("status-test"); return err }},
		{"group add GPU", "error adding GPU 0 to group: ", func() error { return AddToGroup(GroupHandle{}, 0) }},
		{"group remove entity", "error removing entity group type GPU, entity 0 from group: ", func() error { return RemoveEntityFromGroup(GroupHandle{}, FE_GPU, 0) }},
		{"group destroy", "error destroying group: ", func() error { return DestroyGroup(GroupHandle{}) }},
		{"group info", "", func() error { _, err := GetGroupInfo(GroupHandle{}); return err }},
		{"field group create", "error creating DCGM fields group: ", func() error { _, err := FieldGroupCreate("status-test", []Short{DCGM_FI_DEV_GPU_TEMP}); return err }},
		{"field group destroy", "error destroying DCGM fields group: ", func() error { return FieldGroupDestroy(FieldHandle{}) }},
		{"watch fields", "error watching fields: ", func() error { return WatchFieldsWithGroupExNoRefresh(FieldHandle{}, GroupHandle{}, 1000000, 1, 0) }},
		{"unwatch fields", "error unwatching fields: ", func() error { return UnwatchFields(FieldHandle{}, GroupHandle{}) }},
		{"watch field value", "error watching field value: ", func() error { return WatchFieldValue(0, DCGM_FI_DEV_GPU_TEMP, time.Second, time.Second, 0) }},
		{"unwatch field value", "error unwatching field value: ", func() error { return UnwatchFieldValue(0, DCGM_FI_DEV_GPU_TEMP, false) }},
		{"latest values", "error watching fields: ", func() error { _, err := GetLatestValuesForFields(0, []Short{DCGM_FI_DEV_GPU_TEMP}); return err }},
		{"refresh fields", "", UpdateAllFields},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			var native *Error
			require.ErrorAs(t, err, &native)
			require.Equal(t, DCGM_ST_BADPARAM, int(native.Code))
			require.Equal(t, dcgmErrorText(DCGM_ST_BADPARAM), native.Error())
			require.Equal(t, tt.prefix+native.Error(), err.Error())
		})
	}
}
