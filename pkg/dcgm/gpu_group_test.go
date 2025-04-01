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

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupHandle(t *testing.T) {
	gh := GroupHandle{}
	assert.Equal(t, uintptr(0), gh.GetHandle(), "value mismatch")

	inputs := []uintptr{1000, 0, 1, 10, 11, 50, 100, 1939902, 9992932938239, 999999999999999999}

	for _, input := range inputs {
		gh.SetHandle(input)
		assert.Equal(t, input, gh.GetHandle(), "values mismatch")
	}
}

func TestGetGroupInfo(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	runOnlyWithLiveGPUs(t)
	gpus, err := withInjectionGPUs(t, 1)
	require.NoError(t, err)

	gpuID := gpus[0]

	groupID, err := CreateGroup("test1")
	require.NoError(t, err)

	defer func() {
		_ = DestroyGroup(groupID)
	}()

	err = AddEntityToGroup(groupID, FE_GPU, gpuID)
	require.NoError(t, err)

	grInfo, err := GetGroupInfo(groupID)
	require.NoError(t, err)

	assert.Equal(t, "test1", grInfo.GroupName)
	assert.Len(t, grInfo.EntityList, 1)
	assert.Equal(t, FE_GPU, grInfo.EntityList[0].EntityGroupId)
	assert.Equal(t, gpuID, grInfo.EntityList[0].EntityId)
}

func TestCreateGroupWithContext(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	runOnlyWithLiveGPUs(t)

	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()
		groupName := "test_group"

		group, err := CreateGroupWithContext(ctx, groupName)
		require.NoError(t, err)
		require.NotZero(t, group.GetHandle())

		// Clean up
		err = DestroyGroup(group)
		require.NoError(t, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		group, err := CreateGroupWithContext(ctx, "test_group")
		require.Error(t, err)
		require.Equal(t, context.Canceled, err)
		require.Zero(t, group.GetHandle())
	})
}
