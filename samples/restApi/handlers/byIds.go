package handlers

import (
	"net/http"
)

// DeviceInfo handles HTTP requests for device information by device ID
// It returns either JSON or formatted text output based on the request URL
func DeviceInfo(resp http.ResponseWriter, req *http.Request) {
	device := getDeviceInfo(resp, req)
	if device == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, device)
		return
	}

	printer(resp, req, device, deviceInfo)
}

// DeviceStatus handles HTTP requests for device status by device ID
// It returns either JSON or formatted text output based on the request URL
func DeviceStatus(resp http.ResponseWriter, req *http.Request) {
	st := getDeviceStatus(resp, req)
	if st == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, st)
		return
	}

	printer(resp, req, st, deviceStatus)
}

// ProcessInfo handles HTTP requests for process information by PID
// It returns either JSON or formatted text output based on the request URL
func ProcessInfo(resp http.ResponseWriter, req *http.Request) {
	pInfo := getProcessInfo(resp, req)
	if len(pInfo) == 0 {
		return
	}

	if isJson(req) {
		encode(resp, req, pInfo)
		return
	}

	processPrint(resp, req, pInfo)
}

// Health handles HTTP requests for device health status by device ID
// It returns either JSON or formatted text output based on the request URL
func Health(resp http.ResponseWriter, req *http.Request) {
	h := getHealth(resp, req)
	if h == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, h)
		return
	}

	printer(resp, req, h, healthStatus)
}

// Status handles HTTP requests for DCGM daemon status
// It returns either JSON or formatted text output based on the request URL
func Status(resp http.ResponseWriter, req *http.Request) {
	st := getStatus(resp, req)
	if st == nil {
		return
	}

	if isJson(req) {
		encode(resp, req, st)
		return
	}

	printer(resp, req, st, hostengine)
}
