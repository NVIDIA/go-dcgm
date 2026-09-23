package handlers

import (
	"log"
	"net/http"
)

// map of uuids and device id
var uuids map[string]uint

// DevicesUuids initializes a global map of GPU UUIDs to device IDs
// This must be called before using UUID-based endpoints
func DevicesUuids() {
	devicesUuidsWithDeps(realDeviceDeps())
}

func devicesUuidsWithDeps(deps deviceDeps) {
	uuids = make(map[string]uint)

	count, err := deps.getAllDeviceCount()
	if err != nil {
		log.Printf("(DCGM) Error getting devices: %s", err)
		return
	}

	for i := uint(0); i < count; i++ {
		deviceInfo, err := deps.getDeviceInfo(i)
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
	handleDeviceInfo(resp, req, getDeviceInfo)
}

// DeviceStatusByUuid handles HTTP requests for device status by GPU UUID
// It returns either JSON or formatted text output based on the request URL
func DeviceStatusByUuid(resp http.ResponseWriter, req *http.Request) {
	handleDeviceStatus(resp, req, getDeviceStatus)
}

// HealthByUuid handles HTTP requests for device health status by GPU UUID
// It returns either JSON or formatted text output based on the request URL
func HealthByUuid(resp http.ResponseWriter, req *http.Request) {
	handleHealth(resp, req, getHealth)
}
