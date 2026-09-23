package tests

import (
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func TestIntegrationPreflight(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("start embedded hostengine: %v", err)
	}
	defer cleanup()

	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		t.Fatalf("query GPUs through DCGM: %v", err)
	}
	if count == 0 {
		t.Fatal("DCGM did not discover a GPU")
	}
	liveGPU := false
	for id := range count {
		if dcgm.GetGPUStatus(id) == dcgm.EntityStatusOk {
			liveGPU = true
			break
		}
	}
	if !liveGPU {
		t.Fatal("DCGM did not discover a live GPU")
	}
	if _, err := dcgm.GetHostengineVersionInfo(); err != nil {
		t.Fatalf("query embedded hostengine version: %v", err)
	}
}
