//go:build linux && cgo

/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
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
	"errors"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestWatchRollbackRetainsBothErrors checks that watch setup and group cleanup failures both reach the caller.
func TestWatchRollbackRetainsBothErrors(t *testing.T) {
	primary, cleanup := errors.New("watch failed"), errors.New("destroy failed")
	for _, stage := range []string{"add", "watch", "update"} {
		t.Run(stage, func(t *testing.T) {
			m := NewMockWatchGroupOps(gomock.NewController(t))
			group := GroupHandle{}
			expectSetupSteps(watchSetupSteps(m, group, FieldHandle{}, primary), stage)
			m.EXPECT().destroyGroup(group).Return(cleanup)
			_, err := watchFieldsWithOps(m, 5, FieldHandle{}, "watch")
			require.ErrorIs(t, err, primary)
			require.ErrorIs(t, err, cleanup)
		})
	}
	for _, stage := range []string{"enumerate", "add", "watch", "update"} {
		t.Run("pid/"+stage, func(t *testing.T) {
			m := NewMockPIDWatchOps(gomock.NewController(t))
			group := GroupHandle{}
			expectSetupSteps(pidWatchSetupSteps(m, group, primary, true, time.Second, time.Minute, 1), stage)
			m.EXPECT().destroyGroup(group).Return(cleanup)
			_, err := watchPidFieldsWithOps(m, time.Second, time.Minute, 1)
			require.ErrorIs(t, err, primary)
			require.ErrorIs(t, err, cleanup)
		})
	}
}

// TestHealthCleanupRetainsBothErrors checks that cleanup failures preserve the health result or the original error.
func TestHealthCleanupRetainsBothErrors(t *testing.T) {
	primary, cleanup := errors.New("health failed"), errors.New("destroy failed")
	for _, failed := range []bool{false, true} {
		t.Run(strconv.FormatBool(failed), func(t *testing.T) {
			m := NewMockHealthGroupOps(gomock.NewController(t))
			group := GroupHandle{}
			m.EXPECT().createGroup(gomock.Any(), false).Return(group, nil)
			m.EXPECT().addToGroup(group, uint(5)).Return(nil)
			m.EXPECT().healthSet(group, DCGM_HEALTH_WATCH_ALL).Return(nil)
			m.EXPECT().healthCheck(group).Return(HealthResponse{}, pick(failed, primary, nil))
			m.EXPECT().destroyGroup(group).Return(cleanup)
			got, err := healthCheckByGPUWithOps(m, 5)
			require.ErrorIs(t, err, cleanup)
			if failed {
				require.ErrorIs(t, err, primary)
			} else {
				require.Equal(t, uint(5), got.GPU)
			}
		})
	}
}

// TestInitRejectsInvalidMode checks that an unknown mode leaves the session state intact.
func TestInitRejectsInvalidMode(t *testing.T) {
	for _, count := range []int{0, 1} {
		t.Run(strconv.Itoa(count), func(t *testing.T) {
			previous := dcgmInitCounter
			dcgmInitCounter = count
			t.Cleanup(func() { dcgmInitCounter = previous })
			library, sessionMode := dcgmLibHandle, stopMode
			cleanup, err := Init(mode(99))
			require.ErrorIs(t, err, ErrInvalidMode)
			require.Nil(t, cleanup)
			require.Equal(t, count, dcgmInitCounter)
			require.Equal(t, library, dcgmLibHandle)
			require.Equal(t, sessionMode, stopMode)
		})
	}
}

// TestHostengineLookupPreservesCause checks that callers can identify a missing nv-hostengine executable.
func TestHostengineLookupPreservesCause(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	err := cgoAdapter{}.startHostengine()
	require.ErrorIs(t, err, exec.ErrNotFound)
	require.Contains(t, err.Error(), "nv-hostengine")
}

// TestDerivedCleanupFailures checks that both temporary group cleanup failures are reported independently of the field read.
func TestDerivedCleanupFailures(t *testing.T) {
	readErr, fieldErr, groupErr := errors.New("read failed"), errors.New("field release failed"), errors.New("group release failed")
	for _, operation := range []struct {
		name   string
		fields []Short
		check  func(*testing.T, *MockDeviceStatusOps, bool) (error, error)
	}{
		{
			name:   "PCIe bandwidth",
			fields: []Short{DCGM_FI_DEV_PCIE_MAX_LINK_GEN, DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH},
			check: func(t *testing.T, m *MockDeviceStatusOps, failedRead bool) (error, error) {
				got, cleanup, err := getPciBandwidthWithOps(m, 5)
				require.Equal(t, pick(failedRead, int64(0), int64(7876)), got)
				return cleanup, err
			},
		},
		{
			name:   "CPU affinity",
			fields: []Short{DCGM_FI_DEV_CPU_AFFINITY_0, DCGM_FI_DEV_CPU_AFFINITY_1, DCGM_FI_DEV_CPU_AFFINITY_2, DCGM_FI_DEV_CPU_AFFINITY_3},
			check: func(t *testing.T, m *MockDeviceStatusOps, failedRead bool) (error, error) {
				got, cleanup, err := getCPUAffinityWithOps(m, 5)
				require.Equal(t, pick(failedRead, "N/A", "{2,66,130,194}"), got)
				return cleanup, err
			},
		},
	} {
		for _, read := range []struct {
			name  string
			stage string
		}{
			{name: "read succeeds"},
			{name: "read fails", stage: "collect"},
		} {
			for _, cleanup := range []struct {
				name  string
				field error
				group error
			}{
				{name: "field cleanup fails", field: fieldErr},
				{name: "group cleanup fails", group: groupErr},
				{name: "both cleanups fail", field: fieldErr, group: groupErr},
			} {
				t.Run(operation.name+"/"+read.name+"/"+cleanup.name, func(t *testing.T) {
					m := NewMockDeviceStatusOps(gomock.NewController(t))
					values := derivedInt64Samples(operation.fields, []int64{4, 4, 4, 4}[:len(operation.fields)]...)
					expectDerivedFieldRead(m, 5, operation.fields, values, derivedReadErrors{
						stage:        read.stage,
						operation:    readErr,
						fieldCleanup: cleanup.field,
						groupCleanup: cleanup.group,
					})
					failedRead := read.stage != ""
					cleanupErr, err := operation.check(t, m, failedRead)
					require.Equal(t, failedRead, errors.Is(err, readErr))
					require.Equal(t, cleanup.field != nil, errors.Is(cleanupErr, fieldErr))
					require.Equal(t, cleanup.group != nil, errors.Is(cleanupErr, groupErr))
				})
			}
		}
	}
}
