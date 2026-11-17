package main

import (
	"reflect"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestTopologyCells keeps peer links and the self marker under their GPU headers.
func TestTopologyCells(t *testing.T) {
	for _, test := range []struct {
		name  string
		gpus  []uint
		gpu   uint
		links []dcgm.P2PLink
		want  []string
	}{
		{
			name: "sparse IDs",
			gpus: []uint{1, 3},
			gpu:  1,
			links: []dcgm.P2PLink{
				{GPU: 3, Link: dcgm.TwoNVLINKLinks},
			},
			want: []string{"X", "NV2"},
		},
		{
			name: "reordered IDs and links",
			gpus: []uint{2, 0, 1},
			gpu:  0,
			links: []dcgm.P2PLink{
				{GPU: 1, Link: dcgm.P2PLinkHostBridge},
				{GPU: 2, Link: dcgm.P2PLinkCrossCPU},
			},
			want: []string{"SYS", "X", "PHB"},
		},
		{
			name: "missing peer and undisplayed peer",
			gpus: []uint{1, 3},
			gpu:  3,
			links: []dcgm.P2PLink{
				{GPU: 9, Link: dcgm.SingleNVLINKLink},
			},
			want: []string{"", "X"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := topologyCells(test.gpus, test.gpu, test.links)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("topologyCells() = %q, want %q", got, test.want)
			}
		})
	}
}

// TestFormatTopologyRow aligns actual header and data output, including multi-digit GPU IDs.
func TestFormatTopologyRow(t *testing.T) {
	got := formatTopologyRow("", []string{"GPU1", "GPU123"}, "CPUAffinity") +
		formatTopologyRow("GPU1", []string{"X", "NV2"}, "0-15")
	want := "                GPU1    GPU123  CPUAffinity\n" +
		"      GPU1         X       NV2  0-15\n"
	if got != want {
		t.Fatalf("topology table = %q, want %q", got, want)
	}
}
