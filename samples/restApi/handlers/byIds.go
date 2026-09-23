package handlers

import (
	"net/http"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// DeviceInfo handles HTTP requests for device information by device ID
// It returns either JSON or formatted text output based on the request URL
func DeviceInfo(resp http.ResponseWriter, req *http.Request) {
	handleDeviceInfo(resp, req, getDeviceInfo)
}

func handleDeviceInfo(
	resp http.ResponseWriter,
	req *http.Request,
	get func(http.ResponseWriter, *http.Request) *dcgm.Device,
) {
	device := get(resp, req)
	if device == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, device)
		return
	}

	printer(resp, req, device, deviceInfoTemplate)
}

// DeviceStatus handles HTTP requests for device status by device ID
// It returns either JSON or formatted text output based on the request URL
func DeviceStatus(resp http.ResponseWriter, req *http.Request) {
	handleDeviceStatus(resp, req, getDeviceStatus)
}

func handleDeviceStatus(
	resp http.ResponseWriter,
	req *http.Request,
	get func(http.ResponseWriter, *http.Request) *dcgm.DeviceStatus,
) {
	status := get(resp, req)
	if status == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, status)
		return
	}

	printer(resp, req, status, deviceStatusTemplate)
}

// ProcessInfo handles HTTP requests for process information by PID
// It returns either JSON or formatted text output based on the request URL
func ProcessInfo(resp http.ResponseWriter, req *http.Request) {
	handleProcessInfo(resp, req, getProcessInfo)
}

func handleProcessInfo(
	resp http.ResponseWriter,
	req *http.Request,
	get func(http.ResponseWriter, *http.Request) []dcgm.ProcessInfo,
) {
	processes := get(resp, req)
	if len(processes) == 0 {
		return
	}

	if isJson(req) {
		encode(resp, req, processes)
		return
	}

	processPrint(resp, req, processes)
}

// Health handles HTTP requests for device health status by device ID
// It returns either JSON or formatted text output based on the request URL
func Health(resp http.ResponseWriter, req *http.Request) {
	handleHealth(resp, req, getHealth)
}

func handleHealth(
	resp http.ResponseWriter,
	req *http.Request,
	get func(http.ResponseWriter, *http.Request) *dcgm.DeviceHealth,
) {
	health := get(resp, req)
	if health == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, health)
		return
	}

	printer(resp, req, health, healthStatusTemplate)
}

// Status handles HTTP requests for DCGM daemon status
// It returns either JSON or formatted text output based on the request URL
func Status(resp http.ResponseWriter, req *http.Request) {
	handleStatus(resp, req, getStatus)
}

func handleStatus(
	resp http.ResponseWriter,
	req *http.Request,
	get func(http.ResponseWriter, *http.Request) *dcgm.Status,
) {
	status := get(resp, req)
	if status == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, status)
		return
	}

	printer(resp, req, status, hostengineTemplate)
}
