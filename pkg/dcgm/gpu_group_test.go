package dcgm

import (
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
