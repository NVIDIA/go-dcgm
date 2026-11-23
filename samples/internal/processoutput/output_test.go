package processoutput

import (
	"bytes"
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func TestWriteDistinguishesBlankFromZero(t *testing.T) {
	zero := 0.0
	blank := dcgm.DCGM_FT_FP64_NOT_PERMISSIONED
	violation := uint64(dcgm.DCGM_FT_INT64_NOT_SUPPORTED)
	info := dcgm.ProcessInfo{
		ProcessUtilization: dcgm.ProcessUtilInfo{StartTime: dcgm.Time(dcgm.DCGM_FT_INT64_BLANK / 1_000_000), SmUtil: &zero, MemUtil: &blank},
		Clocks:             dcgm.ClockInfo{Memory: dcgm.DCGM_FT_INT32_NOT_SUPPORTED},
		GpuUtilization:     dcgm.UtilizationInfo{GPU: 21, Memory: 37},
		Memory:             dcgm.MemoryInfo{GlobalUsed: dcgm.DCGM_FT_INT64_NOT_FOUND},
		Violations:         dcgm.ViolationTime{Power: &violation},
	}
	var output bytes.Buffer
	if err := Write(&output, info); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Start Time                   : N/A\n",
		"End Time                     : Running\n",
		"Max GPU Memory Used (bytes)  : N/A\n",
		"Avg SM Clock (MHz)           : 0\n",
		"Avg Memory Clock (MHz)       : N/A\n",
		"Avg SM Utilization (%)       : 21\n",
		"Avg Memory Utilization (%)   : 37\n",
		"Due to - Power (%)           : N/A\n",
		"Energy Consumed (Joules)     : N/A\n",
		"Avg SM Utilization (%)       : 0\n",
		"Avg Memory Utilization (%)   : N/A\n",
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("report missing %q", want)
		}
	}
}

func TestWritePreservesReportedEnergyZero(t *testing.T) {
	// The binding returns a non-nil pointer to zero for blank energy.
	var energy uint64
	info := dcgm.ProcessInfo{ProcessUtilization: dcgm.ProcessUtilInfo{EnergyConsumed: &energy}}
	var output bytes.Buffer
	if err := Write(&output, info); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Energy Consumed (Joules)     : 0\n") {
		t.Error("report did not preserve the energy value returned by DCGM")
	}
}
