// Package processoutput formats DCGM process readings for the text samples.
package processoutput

import (
	"io"
	"strconv"
	"text/template"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

const text = `----------------------------------------------------------------------
GPU ID                       : {{.GPU}}
----------Execution Stats---------------------------------------------
PID                          : {{.PID}}
Name                         : {{or .Name "N/A"}}
Start Time                   : {{timestamp .ProcessUtilization.StartTime}}
End Time                     : {{timestamp .ProcessUtilization.EndTime}}
----------Performance Stats-------------------------------------------
Energy Consumed (Joules)     : {{with .ProcessUtilization.EnergyConsumed}}{{.}}{{else}}N/A{{end}}
Max GPU Memory Used (bytes)  : {{int64 .Memory.GlobalUsed}}
Avg SM Clock (MHz)           : {{int32 .Clocks.Cores}}
Avg Memory Clock (MHz)       : {{int32 .Clocks.Memory}}
Avg SM Utilization (%)       : {{int32 .GpuUtilization.GPU}}
Avg Memory Utilization (%)   : {{int32 .GpuUtilization.Memory}}
Avg PCIe Rx Bandwidth (MB)   : {{int64 .PCI.Throughput.Rx}}
Avg PCIe Tx Bandwidth (MB)   : {{int64 .PCI.Throughput.Tx}}
----------Event Stats-------------------------------------------------
Single Bit ECC Errors        : {{int32 .Memory.ECCErrors.SingleBit}}
Double Bit ECC Errors        : {{int32 .Memory.ECCErrors.DoubleBit}}
Critical XID Errors          : {{.XIDErrors.NumErrors}}
----------Slowdown Stats----------------------------------------------
Due to - Power (%)           : {{unsigned .Violations.Power}}
       - Thermal (%)         : {{unsigned .Violations.Thermal}}
       - Reliability (%)     : {{unsigned .Violations.Reliability}}
       - Board Limit (%)     : {{unsigned .Violations.BoardLimit}}
       - Low Utilization (%) : {{unsigned .Violations.LowUtilization}}
       - Sync Boost (%)      : {{unsigned .Violations.SyncBoost}}
----------Process Utilization-----------------------------------------
Avg SM Utilization (%)       : {{float .ProcessUtilization.SmUtil}}
Avg Memory Utilization (%)   : {{float .ProcessUtilization.MemUtil}}
----------------------------------------------------------------------
`

var report = template.Must(template.New("processInfo").Funcs(template.FuncMap{
	"int32":     formatInt32,
	"int64":     formatInt64,
	"unsigned":  formatUnsigned,
	"float":     formatFloat,
	"timestamp": formatTime,
}).Parse(text))

// Write renders one process report, using N/A for blank markers.
// ProcessInfo represents unavailable energy as zero, so the report shows 0 in that case.
func Write(w io.Writer, info dcgm.ProcessInfo) error {
	// #nosec G708 -- report is a package-local template parsed at startup.
	return report.Execute(w, info)
}

// Process clocks and utilization are widened native 32-bit measurements.
func formatInt32(value int64) string {
	if value >= dcgm.DCGM_FT_INT32_BLANK {
		return "N/A"
	}
	return strconv.FormatInt(value, 10)
}

func formatInt64(value int64) string {
	if dcgm.IsInt64Blank(value) {
		return "N/A"
	}
	return strconv.FormatInt(value, 10)
}

func formatUnsigned(value *uint64) string {
	if value == nil || *value >= uint64(dcgm.DCGM_FT_INT64_BLANK) {
		return "N/A"
	}
	return strconv.FormatUint(*value, 10)
}

func formatFloat(value *float64) string {
	if value == nil || dcgm.IsFloat64Blank(*value) {
		return "N/A"
	}
	return strconv.FormatFloat(*value, 'g', -1, 64)
}

func formatTime(value dcgm.Time) string {
	// ProcessInfo divides the native microsecond timestamp by one million.
	if value == dcgm.Time(dcgm.DCGM_FT_INT64_BLANK/1_000_000) {
		return "N/A"
	}
	return value.String()
}
