package handlers

import (
	"log"
	"net/http"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// map of uuids and device id
var uuids map[string]uint

// DevicesUuids initializes a global map of GPU UUIDs to device IDs
// This must be called before using UUID-based endpoints
func DevicesUuids() {
	uuids = make(map[string]uint)

	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		log.Printf("(DCGM) Error getting devices: %s", err)
		return
	}

	for i := uint(0); i < count; i++ {
		deviceInfo, err := dcgm.GetDeviceInfo(i)
		if err != nil {
			log.Printf("(DCGM) Error getting device information: %s", err)
			return
		}

		uuids[deviceInfo.UUID] = i
	}
}

// DeviceInfoByUuid handles HTTP requests for device information by GPU UUID
// It returns either JSON or formatted text output based on the request URL
func DeviceInfoByUuid(resp http.ResponseWriter, req *http.Request) {
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

// DeviceStatusByUuid handles HTTP requests for device status by GPU UUID
// It returns either JSON or formatted text output based on the request URL
func DeviceStatusByUuid(resp http.ResponseWriter, req *http.Request) {
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

// HealthByUuid handles HTTP requests for device health status by GPU UUID
// It returns either JSON or formatted text output based on the request URL
func HealthByUuid(resp http.ResponseWriter, req *http.Request) {
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
