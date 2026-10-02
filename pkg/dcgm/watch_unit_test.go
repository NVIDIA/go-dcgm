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
			create := m.EXPECT().createGroup("watch", false)
			if tt.stage == "create" {
				create.Return(GroupHandle{}, stageErr)
			} else {
				create.Return(group, nil)
				add := m.EXPECT().addToGroup(group, uint(5)).After(create)
				if tt.stage == "add" {
					add.Return(stageErr)
				} else {
					add.Return(nil)
					watch := m.EXPECT().watchFieldsNative(fieldGroup, group, int64(defaultUpdateFreq), float64(defaultMaxKeepAge), int32(defaultMaxKeepSamples)).After(add)
					if tt.stage == "watch" {
						watch.Return(stageErr)
					} else {
						watch.Return(nil)
						update := m.EXPECT().updateAllFields().After(watch)
						if tt.stage == "update" {
							update.Return(stageErr)
						} else {
							update.Return(nil)
						}
					}
				}
				if tt.stage != "" {
					m.EXPECT().destroyGroup(group)
				}
			}

			got, err := watchFieldsWithOps(m, 5, fieldGroup, "watch")
			if tt.stage == "" {
				require.NoError(t, err)
				require.Equal(t, group, got)
				return
			}
			require.Equal(t, GroupHandle{}, got)
			require.ErrorIs(t, err, stageErr)
		})
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
			create := m.EXPECT().createGroup(gomock.Any(), false)
			if tt.fail == "create" {
				create.Return(GroupHandle{}, stageErr)
			} else {
				create.Return(group, nil)
				previous := create
				if len(tt.gpus) == 0 {
					enumerate := m.EXPECT().getSupportedDevices().After(previous)
					if tt.fail == "enumerate" {
						enumerate.Return(nil, stageErr)
					} else {
						enumerate.Return([]uint{5}, nil)
					}
					previous = enumerate
				}
				if tt.fail != "enumerate" {
					add := m.EXPECT().addToGroup(group, uint(5)).After(previous)
					if tt.fail == "add" {
						add.Return(stageErr)
					} else {
						add.Return(nil)
						watch := m.EXPECT().watchPidFieldsNative(group, freq, age, samples).After(add)
						if tt.fail == "watch" {
							watch.Return(stageErr)
						} else {
							watch.Return(nil)
							update := m.EXPECT().updateAllFields().After(watch)
							if tt.fail == "update" {
								update.Return(stageErr)
							} else {
								update.Return(nil)
							}
						}
					}
				}
				if tt.fail != "" {
					m.EXPECT().destroyGroup(group)
				}
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
