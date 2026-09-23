package handlers

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"

	"github.com/gorilla/mux"
)

type deviceDeps struct {
	introspect          func() (dcgm.Status, error)
	getAllDeviceCount   func() (uint, error)
	getSupportedDevices func() ([]uint, error)
	getDeviceInfo       func(uint) (dcgm.Device, error)
	getDeviceStatus     func(uint) (dcgm.DeviceStatus, error)
	healthCheckByGPUID  func(uint) (dcgm.DeviceHealth, error)
}

func realDeviceDeps() deviceDeps {
	return deviceDeps{
		introspect:          dcgm.Introspect,
		getAllDeviceCount:   dcgm.GetAllDeviceCount,
		getSupportedDevices: dcgm.GetSupportedDevices,
		getDeviceInfo:       dcgm.GetDeviceInfo,
		getDeviceStatus:     dcgm.GetDeviceStatus,
		healthCheckByGPUID:  dcgm.HealthCheckByGpuId,
	}
}

func getStatus(resp http.ResponseWriter, req *http.Request) *dcgm.Status {
	return getStatusWithDeps(resp, req, realDeviceDeps())
}

func getStatusWithDeps(resp http.ResponseWriter, req *http.Request, deps deviceDeps) *dcgm.Status {
	status, err := deps.introspect()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		logRequestError(req, err)
		return nil
	}
	return &status
}

func getDeviceInfo(resp http.ResponseWriter, req *http.Request) *dcgm.Device {
	return getDeviceInfoWithDeps(resp, req, realDeviceDeps())
}

func getDeviceInfoWithDeps(resp http.ResponseWriter, req *http.Request, deps deviceDeps) *dcgm.Device {
	var id uint
	for key, value := range mux.Vars(req) {
		switch key {
		case "id":
			id = getId(resp, req, value)
		case "uuid":
			id = getIdByUuid(resp, req, value)
		}
	}

	if id == math.MaxUint32 || !isValidIdWithDeps(id, resp, req, deps) {
		return nil
	}

	device, err := deps.getDeviceInfo(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		logRequestError(req, err)
		return nil
	}
	return &device
}

func getDeviceStatus(resp http.ResponseWriter, req *http.Request) *dcgm.DeviceStatus {
	return getDeviceStatusWithDeps(resp, req, realDeviceDeps())
}

func getDeviceStatusWithDeps(resp http.ResponseWriter, req *http.Request, deps deviceDeps) *dcgm.DeviceStatus {
	var id uint
	for key, value := range mux.Vars(req) {
		switch key {
		case "id":
			id = getId(resp, req, value)
		case "uuid":
			id = getIdByUuid(resp, req, value)
		}
	}

	if id == math.MaxUint32 ||
		!isValidIdWithDeps(id, resp, req, deps) ||
		!isDcgmSupportedWithDeps(id, resp, req, deps) {
		return nil
	}

	status, err := deps.getDeviceStatus(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		logRequestError(req, err)
		return nil
	}
	return &status
}

func getHealth(resp http.ResponseWriter, req *http.Request) *dcgm.DeviceHealth {
	return getHealthWithDeps(resp, req, realDeviceDeps())
}

func getHealthWithDeps(resp http.ResponseWriter, req *http.Request, deps deviceDeps) *dcgm.DeviceHealth {
	var id uint
	for key, value := range mux.Vars(req) {
		switch key {
		case "id":
			id = getId(resp, req, value)
		case "uuid":
			id = getIdByUuid(resp, req, value)
		}
	}

	if id == math.MaxUint32 || !isValidIdWithDeps(id, resp, req, deps) {
		return nil
	}

	health, err := deps.healthCheckByGPUID(id)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		logRequestError(req, err)
		return nil
	}
	return &health
}

type processInfoDeps struct {
	watchPidFields func() (dcgm.GroupHandle, error)
	getProcessInfo func(dcgm.GroupHandle, uint) ([]dcgm.ProcessInfo, error)
	destroyGroup   func(dcgm.GroupHandle) error
	sleep          func(time.Duration)
}

func getProcessInfo(resp http.ResponseWriter, req *http.Request) (pInfo []dcgm.ProcessInfo) {
	return getProcessInfoWithDeps(resp, req, processInfoDeps{
		watchPidFields: dcgm.WatchPidFields,
		getProcessInfo: dcgm.GetProcessInfo,
		destroyGroup:   dcgm.DestroyGroup,
		sleep:          time.Sleep,
	})
}

func getProcessInfoWithDeps(resp http.ResponseWriter, req *http.Request, deps processInfoDeps) (pInfo []dcgm.ProcessInfo) {
	params := mux.Vars(req)

	pid := getId(resp, req, params["pid"])
	if pid == math.MaxUint32 {
		return
	}

	group, err := deps.watchPidFields()
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		logRequestError(req, err)

		return
	}
	defer func() {
		if destroyErr := deps.destroyGroup(group); destroyErr != nil {
			logRequestError(req, destroyErr)
		}
	}()

	// wait for watches to be enabled
	log.Printf("Enabling DCGM watches to start collecting process stats. This may take a few seconds....")
	deps.sleep(3000 * time.Millisecond)

	pInfo, err = deps.getProcessInfo(group, pid)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		logRequestError(req, err)
	}

	return
}
