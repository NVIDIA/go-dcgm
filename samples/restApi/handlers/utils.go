package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/NVIDIA/go-dcgm/samples/internal/processoutput"
)

const (
	base                       = 10
	bitsize                    = 32
	internalServerErrorMessage = "internal server error"

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

var (
	deviceInfoTemplate   = template.Must(template.New("deviceInfo").Parse(deviceInfo))
	deviceStatusTemplate = template.Must(template.New("deviceStatus").Parse(deviceStatus))
	healthStatusTemplate = template.Must(template.New("healthStatus").Parse(healthStatus))
	hostengineTemplate   = template.Must(template.New("hostengine").Parse(hostengine))
)

func logRequestError(req *http.Request, err error) {
	if err != nil {
		logRequestMessage(req, err.Error())
	}
}

func writeInternalError(resp http.ResponseWriter, req *http.Request, err error) {
	logRequestError(req, err)
	http.Error(resp, internalServerErrorMessage, http.StatusInternalServerError)
}

func logRequestStatus(req *http.Request, status int) {
	logRequestMessage(req, http.StatusText(status))
}

func logRequestMessage(req *http.Request, message string) {
	method, path := "", ""
	if req != nil {
		method = req.Method
		if req.URL != nil {
			path = req.URL.Path
		}
	}

	log.Printf(
		"error: method=%s path=%s message=%s",
		strconv.Quote(method),
		strconv.Quote(path),
		strconv.Quote(message),
	)
}

// getId converts a string key to a GPU ID
// Returns math.MaxUint32 if the conversion fails
func getId(resp http.ResponseWriter, req *http.Request, key string) uint {
	id, err := strconv.ParseUint(key, base, bitsize)
	if err == nil && id == math.MaxUint32 {
		err = fmt.Errorf("invalid ID %q", key)
	}
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		logRequestError(req, err)

		return math.MaxUint32
	}

	return uint(id)
}

func getIdByUuid(resp http.ResponseWriter, req *http.Request, key string) uint {
	id, exists := uuids[key]
	if !exists {
		http.NotFound(resp, req)
		logRequestStatus(req, http.StatusNotFound)

		return math.MaxUint32
	}

	return id
}

func isValidIdWithDeps(id uint, resp http.ResponseWriter, req *http.Request, deps deviceDeps) bool {
	count, err := deps.getAllDeviceCount()
	if err != nil {
		writeInternalError(resp, req, err)
		return false
	}

	if id >= count {
		http.NotFound(resp, req)
		logRequestStatus(req, http.StatusNotFound)
		return false
	}

	return true
}

func isDcgmSupportedWithDeps(gpuID uint, resp http.ResponseWriter, req *http.Request, deps deviceDeps) bool {
	gpus, err := deps.getSupportedDevices()
	if err != nil {
		writeInternalError(resp, req, err)
		return false
	}

	for _, gpu := range gpus {
		if gpuID == gpu {
			return true
		}
	}

	err = fmt.Errorf("error adding gpu %d to group: This gpu is not supported by dcgm", gpuID)
	writeInternalError(resp, req, err)
	return false
}

// isJson checks if the request URL ends with "/json" to determine output format
// Returns true if JSON output is requested
func isJson(req *http.Request) bool {
	return strings.HasSuffix(req.URL.Path, "/json")
}

// printer formats and writes templated text output to the response.
func printer(resp http.ResponseWriter, req *http.Request, stats any, t *template.Template) {
	// Render first so a template error cannot commit a partial 200 response.
	var output bytes.Buffer
	// #nosec G708 -- t is parsed from package-local templates at startup.
	if err := t.Execute(&output, stats); err != nil {
		writeInternalError(resp, req, err)
		return
	}

	if _, err := output.WriteTo(resp); err != nil {
		logRequestError(req, err)
	}
}

// encode writes JSON-formatted output to the response
func encode(resp http.ResponseWriter, req *http.Request, stats any) {
	resp.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(resp).Encode(stats); err != nil {
		writeInternalError(resp, req, err)
	}
}

// processPrint formats and writes process information to the response
func processPrint(resp http.ResponseWriter, req *http.Request, pInfo []dcgm.ProcessInfo) {
	// Render every entry before writing any of them to the response.
	var output bytes.Buffer
	for i := range pInfo {
		if err := processoutput.Write(&output, pInfo[i]); err != nil {
			writeInternalError(resp, req, err)
			return
		}
	}

	if _, err := output.WriteTo(resp); err != nil {
		logRequestError(req, err)
	}
}
