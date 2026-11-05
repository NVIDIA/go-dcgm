package dcgm

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestWatchFieldsWithOpsOwnsGroup checks that failed setup calls destroyGroup
// if a group was created, and successful setup returns the group.
func TestWatchFieldsWithOpsOwnsGroup(t *testing.T) {
	for _, tt := range []struct {
		name  string
		stage string
	}{
		{name: "create failure", stage: "create"},
		{name: "add failure", stage: "add"},
		{name: "watch failure", stage: "watch"},
		{name: "update failure", stage: "update"},
		{name: "success"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockWatchGroupOps(gomock.NewController(t))
			var group GroupHandle
			group.SetHandle(7)
			fieldGroup := FieldHandle{}
			stageErr := errors.New(tt.stage + " failed")
			expectSetupSteps(watchSetupSteps(m, group, fieldGroup, stageErr), tt.stage)
			if tt.stage != "" && tt.stage != "create" {
				m.EXPECT().destroyGroup(group)
			}

			got, err := watchFieldsWithOps(m, 5, fieldGroup, "watch")
			if tt.stage == "" {
				require.NoError(t, err)
				require.Equal(t, group, got)
				return
			}

			require.Equal(t, GroupHandle{}, got)
			require.ErrorIs(t, err, stageErr)
			if tt.stage == "update" {
				require.ErrorContains(t, err, "error refreshing watched fields")
			}
		})
	}
}

// watchSetupSteps lists the native calls watchFieldsWithOps makes for GPU 5,
// in order. A step that fails returns stageErr.
func watchSetupSteps(m *MockWatchGroupOps, group GroupHandle, fieldGroup FieldHandle, stageErr error) []setupStep {
	return []setupStep{
		{stage: "create", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().createGroup("watch", false).Return(pick(fail, GroupHandle{}, group), pick(fail, stageErr, nil))
		}},
		{stage: "add", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().addToGroup(group, uint(5)).Return(pick(fail, stageErr, nil))
		}},
		{stage: "watch", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().watchFieldsNative(fieldGroup, group, int64(defaultUpdateFreq), float64(defaultMaxKeepAge), int32(defaultMaxKeepSamples)).
				Return(pick(fail, stageErr, nil))
		}},
		{stage: "update", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().updateAllFields().Return(pick(fail, stageErr, nil))
		}},
	}
}

// TestWatchPIDFieldsWithOpsOwnsGroup checks that failed setup calls destroyGroup
// if a group was created, and successful setup returns the group.
func TestWatchPIDFieldsWithOpsOwnsGroup(t *testing.T) {
	freq, age, samples := 30*time.Second, time.Hour, 2
	for _, tt := range []struct {
		name string
		fail string
		gpus []uint
	}{
		{name: "create failure", fail: "create", gpus: []uint{5}},
		{name: "enumeration failure", fail: "enumerate"},
		{name: "add failure", fail: "add", gpus: []uint{5}},
		{name: "watch failure", fail: "watch", gpus: []uint{5}},
		{name: "update failure", fail: "update", gpus: []uint{5}},
		{name: "success", gpus: []uint{5}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockPIDWatchOps(gomock.NewController(t))
			var group GroupHandle
			group.SetHandle(9)
			stageErr := errors.New(tt.fail + " failed")
			steps := pidWatchSetupSteps(m, group, stageErr, len(tt.gpus) == 0, freq, age, samples)
			expectSetupSteps(steps, tt.fail)
			if tt.fail != "" && tt.fail != "create" {
				m.EXPECT().destroyGroup(group)
			}

			got, err := watchPidFieldsWithOps(m, freq, age, samples, tt.gpus...)
			if tt.fail == "" {
				require.NoError(t, err)
				require.Equal(t, group, got)
				return
			}

			require.Equal(t, GroupHandle{}, got)
			require.ErrorIs(t, err, stageErr)
		})
	}
}

// pidWatchSetupSteps lists the native calls watchPidFieldsWithOps makes for
// GPU 5, in order. enumerate adds the GPU lookup that runs when the caller
// names no GPUs. A step that fails returns stageErr.
func pidWatchSetupSteps(
	m *MockPIDWatchOps, group GroupHandle, stageErr error, enumerate bool, freq, age time.Duration, samples int,
) []setupStep {
	steps := []setupStep{{stage: "create", expect: func(fail bool) *gomock.Call {
		return m.EXPECT().createGroup(gomock.Any(), false).Return(pick(fail, GroupHandle{}, group), pick(fail, stageErr, nil))
	}}}
	if enumerate {
		steps = append(steps, setupStep{stage: "enumerate", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().getSupportedDevices().Return(pick(fail, nil, []uint{5}), pick(fail, stageErr, nil))
		}})
	}

	return append(steps,
		setupStep{stage: "add", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().addToGroup(group, uint(5)).Return(pick(fail, stageErr, nil))
		}},
		setupStep{stage: "watch", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().watchPidFieldsNative(group, freq, age, samples).Return(pick(fail, stageErr, nil))
		}},
		setupStep{stage: "update", expect: func(fail bool) *gomock.Call {
			return m.EXPECT().updateAllFields().Return(pick(fail, stageErr, nil))
		}},
	)
}

// setupStep is one expected native call in a watch setup. stage names it so a
// test case can make it fail; expect adds the gomock expectation and returns
// the failing result when fail is set.
type setupStep struct {
	stage  string
	expect func(fail bool) *gomock.Call
}

// expectSetupSteps expects steps in order, each after the one before. The step
// whose stage is failAt fails and is the last one expected; an empty failAt
// expects every step to succeed.
func expectSetupSteps(steps []setupStep, failAt string) {
	var previous *gomock.Call
	for _, step := range steps {
		call := step.expect(step.stage == failAt)
		if previous != nil {
			call.After(previous)
		}

		if step.stage == failAt {
			return
		}
		previous = call
	}
}

// pick returns failed when fail is set and ok otherwise.
func pick[T any](fail bool, failed, ok T) T {
	if fail {
		return failed
	}
	return ok
}

// TestWatchFieldsWithGroupExRefreshCount checks how many refreshes each group
// watch call makes after adding the watch, that a failed watch skips the
// refresh, and that a refresh failure is wrapped.
func TestWatchFieldsWithGroupExRefreshCount(t *testing.T) {
	watchErr := errors.New("watch failed")
	updateErr := errors.New("update failed")
	refresh := func(api watchGroupOps, f FieldHandle, g GroupHandle) error {
		return watchFieldsWithGroupExWithOps(api, f, g, 1000000, 60, 5)
	}
	noRefresh := func(api watchGroupOps, f FieldHandle, g GroupHandle) error {
		return watchFieldsWithGroupExNoRefreshWithOps(api, f, g, 1000000, 60, 5)
	}
	for _, tt := range []struct {
		name      string
		call      func(watchGroupOps, FieldHandle, GroupHandle) error
		watchErr  error
		updateErr error
		wantErr   error
		wantCalls int
	}{
		{name: "WatchFieldsWithGroupEx refreshes once", call: refresh, wantCalls: 1},
		{name: "WatchFieldsWithGroupExNoRefresh does not refresh", call: noRefresh, wantCalls: 0},
		{name: "failed watch skips the refresh", call: refresh, watchErr: watchErr, wantErr: watchErr},
		{name: "failed watch without refresh", call: noRefresh, watchErr: watchErr, wantErr: watchErr},
		{name: "refresh failure is wrapped", call: refresh, updateErr: updateErr, wantErr: updateErr, wantCalls: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMockWatchGroupOps(gomock.NewController(t))
			var group GroupHandle
			group.SetHandle(3)
			fieldGroup := FieldHandle{}

			watch := m.EXPECT().watchFieldsNative(fieldGroup, group, int64(1000000), 60.0, int32(5)).Return(tt.watchErr)
			m.EXPECT().updateAllFields().After(watch).Return(tt.updateErr).Times(tt.wantCalls)

			err := tt.call(m, fieldGroup, group)
			if tt.wantErr == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, tt.wantErr)
			if tt.updateErr != nil {
				require.ErrorContains(t, err, "error refreshing watched fields")
			} else {
				require.Equal(t, tt.watchErr, err)
			}
		})
	}
}
