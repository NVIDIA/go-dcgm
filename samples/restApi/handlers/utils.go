package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

const (
	base    = 10
	bitsize = 32

	// DeviceInfo is the template for formatting device information output
	deviceInfo = `Driver Version         : {{.Identifiers.DriverVersion}}
GPU                    : {{.GPU}}
DCGMSupported          : {{.DCGMSupported}}
UUID                   : {{.UUID}}
Brand                  : {{.Identifiers.Brand}}
Model                  : {{.Identifiers.Model}}
Serial Number          : {{.Identifiers.Serial}}
Vbios                  : {{or .Identifiers.Vbios "N/A"}}
InforomImage Version   : {{.Identifiers.InforomImageVersion}}
Bus ID                 : {{.PCI.BusID}}
BAR1 (MB)              : {{or .PCI.BAR1 "N/A"}}
FrameBuffer Memory (MB): {{or .PCI.FBTotal "N/A"}}
Bandwidth (MB/s)       : {{or .PCI.Bandwidth "N/A"}}
Power (W)              : {{or .Power "N/A"}}
CPUAffinity            : {{or .CPUAffinity "N/A"}}
P2P Available          : {{if not .Topology}}None{{else}}{{range .Topology}}
    GPU{{.GPU}} - (BusID){{.BusID}} - {{.Link.PCIPaths}}{{end}}{{end}}
---------------------------------------------------------------------
`
	// DeviceStatus is the template for formatting device status output
	deviceStatus = `Power (W)		: {{.Power}}
Temperature (°C)        : {{.Temperature}}
Sm Utilization (%)      : {{.Utilization.GPU}}
Memory Utilization (%)  : {{.Utilization.Memory}}
Encoder Utilization (%) : {{.Utilization.Encoder}}
Decoder Utilization (%) : {{.Utilization.Decoder}}
Memory Clock (MHz       : {{.Clocks.Memory}}
SM Clock (MHz)          : {{.Clocks.Cores}}
`

	// ProcessInfo is the template for formatting process information output
	processInfo = `----------------------------------------------------------------------
GPU ID                       : {{.GPU}}
----------Execution Stats---------------------------------------------
PID                          : {{.PID}}
Name                         : {{or .Name "N/A"}}
Start Time                   : {{.ProcessUtilization.StartTime.String}}
End Time                     : {{.ProcessUtilization.EndTime.String}}
----------Performance Stats-------------------------------------------
Energy Consumed (Joules)     : {{or .ProcessUtilization.EnergyConsumed "N/A"}}
Max GPU Memory Used (bytes)  : {{or .Memory.GlobalUsed "N/A"}}
Avg SM Clock (MHz)           : {{or .Clocks.Cores "N/A"}}
Avg Memory Clock (MHz)       : {{or .Clocks.Memory "N/A"}}
Avg SM Utilization (%)       : {{or .GpuUtilization.Memory "N/A"}}
Avg Memory Utilization (%)   : {{or .GpuUtilization.GPU "N/A"}}
Avg PCIe Rx Bandwidth (MB)   : {{or .PCI.Throughput.Rx "N/A"}}
Avg PCIe Tx Bandwidth (MB)   : {{or .PCI.Throughput.Tx "N/A"}}
----------Event Stats-------------------------------------------------
Single Bit ECC Errors        : {{or .Memory.ECCErrors.SingleBit "N/A"}}
Double Bit ECC Errors        : {{or .Memory.ECCErrors.DoubleBit "N/A"}}
Critical XID Errors          : {{.XIDErrors.NumErrors}}
----------Slowdown Stats----------------------------------------------
Due to - Power (%)           : {{.Violations.Power}}
       - Thermal (%)         : {{.Violations.Thermal}}
       - Reliability (%)     : {{.Violations.Reliability}}
       - Board Limit (%)     : {{.Violations.BoardLimit}}
       - Low Utilization (%) : {{.Violations.LowUtilization}}
       - Sync Boost (%)      : {{.Violations.SyncBoost}}
----------Process Utilization-----------------------------------------
Avg SM Utilization (%)       : {{or .ProcessUtilization.SmUtil "N/A"}}
Avg Memory Utilization (%)   : {{or .ProcessUtilization.MemUtil "N/A"}}
----------------------------------------------------------------------
`
	// HealthStatus is the template for formatting health status output
	healthStatus = `GPU                : {{.GPU}}
Status             : {{.Status}}
{{range .Watches}}
Type               : {{.Type}}
Status             : {{.Status}}
Error              : {{.Error}}
{{end}}`

	// HostEngine is the template for formatting DCGM host engine status
	hostengine = `Memory(KB)      : {{.Memory}}
CPU(%)          : {{printf "%.2f" .CPU}}
`
)

// getId converts a string key to a GPU ID
// Returns math.MaxUint32 if the conversion fails
func getId(resp http.ResponseWriter, req *http.Request, key string) uint {
	id, err := strconv.ParseUint(key, base, bitsize)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		log.Printf("error: %v%v: %v", req.Host, req.URL, err.Error())

		return math.MaxUint32
	}

	return uint(id)
}

func getIdByUuid(resp http.ResponseWriter, req *http.Request, key string) uint {
	id, exists := uuids[key]
	if !exists {
		http.NotFound(resp, req)
		log.Printf("error: %v%v:  %v (page not found)", req.Host, req.URL, http.StatusNotFound)

		return math.MaxUint32
	}

	return id
}

// isValidId checks if the given GPU ID exists and is valid
// Returns true if the ID is valid, false otherwise
func isValidId(id uint, resp http.ResponseWriter, req *http.Request) bool {
	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		log.Printf("error: %v%v: %v", req.Host, req.URL, err.Error())

		return false
	}

	if id >= count {
		http.NotFound(resp, req)
		log.Printf("error: %v%v: %v (page not found)", req.Host, req.URL, http.StatusNotFound)

		return false
	}

	return true
}

// isDcgmSupported checks if DCGM supports the given GPU
// Returns true if supported, false otherwise
func isDcgmSupported(gpuId uint, resp http.ResponseWriter, req *http.Request) bool {
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		log.Printf("error: %v%v: %v", req.Host, req.URL, err.Error())

		return false
	}

	for _, gpu := range gpus {
		if gpuId == gpu {
			return true
		}
	}

	err = fmt.Errorf("error adding gpu %d to group: This gpu is not supported by dcgm", gpuId)
	http.Error(resp, err.Error(), http.StatusInternalServerError)
	log.Printf("error: %v%v: %v", req.Host, req.URL, err.Error())

	return false
}

// isJson checks if the request URL ends with "json" to determine output format
// Returns true if JSON output is requested
func isJson(req *http.Request) bool {
	url := (req.URL).String()
	return url[len(url)-4:] == "json"
}

// print formats and writes templated text output to the response
func printer(resp http.ResponseWriter, req *http.Request, stats any, templ string) {
	t := template.Must(template.New("").Parse(templ))
	if err := t.Execute(resp, stats); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		log.Printf("error: %v%v: %v", req.Host, req.URL, err.Error())
	}
}

// encode writes JSON-formatted output to the response
func encode(resp http.ResponseWriter, req *http.Request, stats any) {
	resp.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(resp).Encode(stats); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		log.Printf("error: %v%v: %v", req.Host, req.URL, err.Error())
	}
}

// processPrint formats and writes process information to the response
func processPrint(resp http.ResponseWriter, req *http.Request, pInfo []dcgm.ProcessInfo) {
	t := template.Must(template.New("Process").Parse(processInfo))
	for i := range pInfo {
		if err := t.Execute(resp, pInfo[i]); err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			log.Printf("error: %v%v: %v", req.Host, req.URL, err.Error())

			return
		}
	}
}
