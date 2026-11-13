package main

import (
	"errors"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func TestReadConfigurationsContinuesAfterError(t *testing.T) {
	currentErr := errors.New("current read failed")
	targetErr := errors.New("target read failed")
	var states []dcgm.ConfigType
	err := readConfigurations(func(_ dcgm.GroupHandle, kind dcgm.ConfigType) ([]dcgm.GPUConfig, error) {
		states = append(states, kind)
		if kind == dcgm.ConfigCurrentState {
			return nil, currentErr
		}
		return nil, targetErr
	})
	if len(states) != 2 || states[0] != dcgm.ConfigCurrentState || states[1] != dcgm.ConfigTargetState {
		t.Errorf("read states %v; want current followed by target", states)
	}
	if !errors.Is(err, currentErr) || !errors.Is(err, targetErr) {
		t.Errorf("got error %v; want both read errors", err)
	}
}
