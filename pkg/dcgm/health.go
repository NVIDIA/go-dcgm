/*
 * Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"unsafe"
)

// healthSetErrorFormat wraps every error returned while setting health watches.
const healthSetErrorFormat = "error setting health watches: %w"

// SystemWatch represents a health watch system and its status
type SystemWatch struct {
	// Type identifies the type of health watch system
	Type string
	// Status indicates the current health status
	Status string
	// Error contains any error message if status is not healthy
	Error string
}

// DeviceHealth represents the health status of a GPU device
type DeviceHealth struct {
	// GPU is the ID of the GPU device
	GPU uint
	// Status indicates the overall health status of the GPU
	Status string
	// Watches contains the status of individual health watch systems
	Watches []SystemWatch
}

// HealthSet enables the DCGM health check system for the given systems.
// It configures which health watch systems should be monitored for the specified group.
func HealthSet(groupID GroupHandle, systems HealthSystem) (err error) {
	return cgoAdapter{}.healthSet(groupID, systems)
}

// healthSet enables the requested native health watches for a GPU group.
func (cgoAdapter) healthSet(groupID GroupHandle, systems HealthSystem) (err error) {
	result := C.dcgmHealthSet(handle.handle, groupID.handle, C.dcgmHealthSystems_t(systems))
	if err := errorString(result); err != nil {
		return fmt.Errorf(healthSetErrorFormat, err)
	}

	return nil
}

// HealthSetWithParams enables the given health watch systems with custom
// sampling and retention settings. updateInterval should match how often
// HealthCheck is called; maxKeepAge should cover the longest interval between
// calls.
func HealthSetWithParams(
	groupID GroupHandle,
	systems HealthSystem,
	updateInterval time.Duration,
	maxKeepAge time.Duration,
) error {
	params := healthSetParams(groupID, systems, updateInterval, maxKeepAge)

	result := C.dcgmHealthSet_v2(handle.handle, &params)
	if err := errorString(result); err != nil {
		return fmt.Errorf(healthSetErrorFormat, err)
	}

	return nil
}

// healthSetParams builds the parameters shared by dcgmHealthSet_v2 and dcgmHealthSet_v3.
func healthSetParams(
	groupID GroupHandle,
	systems HealthSystem,
	updateInterval time.Duration,
	maxKeepAge time.Duration,
) C.dcgmHealthSetParams_v2 {
	return C.dcgmHealthSetParams_v2{
		version:        makeVersion2(unsafe.Sizeof(C.dcgmHealthSetParams_v2{})),
		groupId:        groupID.handle,
		systems:        C.dcgmHealthSystems_t(systems),
		updateInterval: C.longlong(updateInterval.Microseconds()),
		maxKeepAge:     C.double(maxKeepAge.Seconds()),
	}
}

// DCGM 4.7.1 health APIs. Older libraries and host engines do not provide them.
const (
	// healthSetV3Symbol is the libdcgm symbol of the strict health watch setter.
	healthSetV3Symbol = "dcgmHealthSet_v3"
	// healthSetV3Feature names the strict health watch setter in FeatureUnavailableError.
	healthSetV3Feature = "DCGM strict health set API"
	// healthSupportedSystemsSymbol is the libdcgm symbol of the supported-systems query.
	healthSupportedSystemsSymbol = "dcgmHealthGetSupportedSystems"
	// healthSupportedSystemsFeature names the supported-systems query in FeatureUnavailableError.
	healthSupportedSystemsFeature = "DCGM supported health systems API"
)

// HealthSetWithParamsStrict is like HealthSetWithParams, but calls dcgmHealthSet_v3.
//
// If systems contains a health system that the host engine does not support, it
// returns a *Error with Code DCGM_ST_NOT_SUPPORTED and leaves the group's watches
// unchanged. DCGM_HEALTH_WATCH_ALL enables every system that
// HealthGetSupportedSystems reports. Bits above the 32-bit DCGM mask are rejected
// with a *Error with Code DCGM_ST_BADPARAM before DCGM is called.
//
// IsFeatureUnavailable(err) is true when the loaded library or the connected host
// engine lacks the API. Other DCGM failures are returned as *Error.
func HealthSetWithParamsStrict(
	groupID GroupHandle,
	systems HealthSystem,
	updateInterval time.Duration,
	maxKeepAge time.Duration,
) error {
	// DCGM's mask is 32 bits wide; a larger value would lose its high bits and
	// could turn every watch off.
	if uint64(systems) > math.MaxUint32 {
		return fmt.Errorf(healthSetErrorFormat, &Error{
			msg:  fmt.Sprintf("health systems %#x do not fit in DCGM's 32-bit mask", uint64(systems)),
			Code: C.DCGM_ST_BADPARAM,
		})
	}

	// Libraries older than 4.7.1 do not export the v3 call.
	if !dcgmSymbolAvailable(healthSetV3Symbol) {
		return &FeatureUnavailableError{Feature: healthSetV3Feature, Symbol: healthSetV3Symbol}
	}

	params := healthSetParams(groupID, systems, updateInterval, maxKeepAge)
	result := C.dcgmHealthSet_v3(handle.handle, &params)
	if err := dcgmFeatureResult(healthSetV3Feature, healthSetV3Symbol, result, dcgmErrorText(result)); err != nil {
		return fmt.Errorf(healthSetErrorFormat, err)
	}

	return nil
}

// HealthGetSupportedSystems returns the health systems that the connected host
// engine supports. A supported system may not apply to every device. Use HealthGet
// for the systems enabled on a group.
//
// IsFeatureUnavailable(err) is true when the loaded library or the connected host
// engine lacks the API. Other DCGM failures are returned as *Error.
func HealthGetSupportedSystems() (HealthSystem, error) {
	// Libraries older than 4.7.1 do not export the query.
	if !dcgmSymbolAvailable(healthSupportedSystemsSymbol) {
		return 0, &FeatureUnavailableError{Feature: healthSupportedSystemsFeature, Symbol: healthSupportedSystemsSymbol}
	}

	supported := C.dcgmHealthSupportedSystems_t{version: C.dcgmHealthSupportedSystems_version}
	result := C.dcgmHealthGetSupportedSystems(handle.handle, &supported)
	if err := dcgmFeatureResult(healthSupportedSystemsFeature, healthSupportedSystemsSymbol, result, dcgmErrorText(result)); err != nil {
		return 0, fmt.Errorf("error getting supported health systems: %w", err)
	}

	return HealthSystem(supported.systems), nil
}

// HealthGet retrieves the current state of the DCGM health check system.
// It returns which health watch systems are currently enabled for the specified group.
func HealthGet(groupID GroupHandle) (HealthSystem, error) {
	var systems C.dcgmHealthSystems_t

	result := C.dcgmHealthGet(handle.handle, groupID.handle, (*C.dcgmHealthSystems_t)(unsafe.Pointer(&systems)))
	if err := errorString(result); err != nil {
		return HealthSystem(0), err
	}

	return HealthSystem(systems), nil
}

// DiagErrorDetail contains detailed information about a health check error
type DiagErrorDetail struct {
	// Message contains a human-readable description of the error
	Message string
	// Code identifies the specific type of error
	Code HealthCheckErrorCode
}

// Incident represents a health check incident that occurred
type Incident struct {
	// System identifies which health watch system detected the incident
	System HealthSystem
	// Health indicates the severity of the incident
	Health HealthResult
	// Error contains detailed information about the incident
	Error DiagErrorDetail
	// EntityInfo identifies the GPU or component where the incident occurred
	EntityInfo GroupEntityPair
}

// HealthResponse contains the results of a health check operation
type HealthResponse struct {
	// OverallHealth indicates the aggregate health status across all watches
	OverallHealth HealthResult
	// Incidents contains details about any health issues detected
	Incidents []Incident
}

// HealthCheck checks the configured watches for any errors/failures/warnings that have occurred
// since the last time this check was invoked. On the first call, stateful information
// about all of the enabled watches within a group is created but no error results are
// provided. On subsequent calls, any error information will be returned.
// A response with more incidents than the native array can hold returns an error.
func HealthCheck(groupID GroupHandle) (HealthResponse, error) {
	return cgoAdapter{}.healthCheck(groupID)
}

// healthCheck reads native group health and converts the reported incidents.
func (cgoAdapter) healthCheck(groupID GroupHandle) (HealthResponse, error) {
	var healthResults C.dcgmHealthResponse_v5
	healthResults.version = makeVersion5(unsafe.Sizeof(healthResults))

	result := C.dcgmHealthCheck(handle.handle, groupID.handle, (*C.dcgmHealthResponse_t)(unsafe.Pointer(&healthResults)))

	if err := errorString(result); err != nil {
		return HealthResponse{}, &Error{msg: dcgmErrorText(result), Code: result}
	}

	return healthResponseFromNative(&healthResults)
}

// healthResponseFromNative validates the native incident count before copying
// incidents into Go-owned values.
func healthResponseFromNative(healthResults *C.dcgmHealthResponse_v5) (HealthResponse, error) {
	incidents := uint(healthResults.incidentCount)
	maxIncidents := uint(len(healthResults.incidents))

	// A count larger than the array means DCGM returned a bad response; reading
	// past the array would read invalid memory.
	if incidents > maxIncidents {
		return HealthResponse{}, fmt.Errorf("health response contains %d incidents, maximum %d", incidents, maxIncidents)
	}

	response := HealthResponse{
		OverallHealth: HealthResult(healthResults.overallHealth),
		Incidents:     make([]Incident, incidents),
	}

	for i := uint(0); i < incidents; i++ {
		response.Incidents[i] = incidentFromNative(&healthResults.incidents[i])
	}

	return response, nil
}

// incidentFromNative copies one native health incident into a Go Incident.
func incidentFromNative(incident *C.dcgmIncidentInfo_t) Incident {
	return Incident{
		System: HealthSystem(incident.system),
		Health: HealthResult(incident.health),
		Error: DiagErrorDetail{
			Message: *stringPtr(&incident.error.msg[0]),
			Code:    HealthCheckErrorCode(incident.error.code),
		},
		EntityInfo: GroupEntityPair{
			EntityGroupId: Field_Entity_Group(incident.entityInfo.entityGroupId),
			EntityId:      uint(incident.entityInfo.entityId),
		},
	}
}

// healthCheckByGpuId checks the health of one GPU using the real DCGM calls.
// HealthCheckByGpuId calls it.
func healthCheckByGpuId(gpuID uint) (deviceHealth DeviceHealth, err error) {
	return healthCheckByGPUWithOps(cgoAdapter{}, gpuID)
}

// healthCheckByGPUWithOps checks one GPU using a temporary group.
// It attempts to destroy the group before returning and ignores cleanup errors.
func healthCheckByGPUWithOps(api healthGroupOps, gpuID uint) (deviceHealth DeviceHealth, err error) {
	name := fmt.Sprintf("health%d", rand.Uint64())
	groupID, err := api.createGroup(name, false)
	if err != nil {
		return
	}
	defer func() {
		_ = api.destroyGroup(groupID)
	}()

	err = api.addToGroup(groupID, gpuID)
	if err != nil {
		return
	}

	// Watch every system so the check reports any problem on this GPU.
	err = api.healthSet(groupID, DCGM_HEALTH_WATCH_ALL)
	if err != nil {
		return
	}

	result, err := api.healthCheck(groupID)
	if err != nil {
		return
	}

	deviceHealth = deviceHealthFromResponse(gpuID, result)

	return
}

// deviceHealthFromResponse turns a group health check result into the
// per-GPU DeviceHealth summary, with one SystemWatch per incident.
func deviceHealthFromResponse(gpuID uint, result HealthResponse) DeviceHealth {
	watches := make([]SystemWatch, len(result.Incidents))
	for i := range result.Incidents {
		watches[i] = SystemWatch{
			Type:   systemWatch(result.Incidents[i].System),
			Status: healthStatus(result.Incidents[i].Health),
			Error:  result.Incidents[i].Error.Message,
		}
	}

	return DeviceHealth{
		GPU:     gpuID,
		Status:  healthStatus(result.OverallHealth),
		Watches: watches,
	}
}

// healthStatus returns the display name of a DCGM health result: "Healthy",
// "Warning", "Failure", or "N/A" for any other value.
func healthStatus(status HealthResult) string {
	switch status {
	case 0:
		return "Healthy"
	case 10:
		return "Warning"
	case 20:
		return "Failure"
	}
	return "N/A"
}

// systemWatch returns the display name of one health watch system, such as
// "PCIe watches". It returns "N/A" for a value that is not a single known system.
func systemWatch(watch HealthSystem) string {
	switch watch {
	case DCGM_HEALTH_WATCH_PCIE:
		return "PCIe watches"
	case DCGM_HEALTH_WATCH_NVLINK:
		return "NVLINK watches"
	case DCGM_HEALTH_WATCH_PMU:
		return "Power Management unit watches"
	case DCGM_HEALTH_WATCH_MCU:
		return "Microcontroller unit watches"
	case DCGM_HEALTH_WATCH_MEM:
		return "Memory watches"
	case DCGM_HEALTH_WATCH_SM:
		return "Streaming Multiprocessor watches"
	case DCGM_HEALTH_WATCH_INFOROM:
		return "Inforom watches"
	case DCGM_HEALTH_WATCH_THERMAL:
		return "Temperature watches"
	case DCGM_HEALTH_WATCH_POWER:
		return "Power watches"
	case DCGM_HEALTH_WATCH_DRIVER:
		return "Driver-related watches"
	case DCGM_HEALTH_WATCH_NVSWITCH_NONFATAL:
		return "NVSwitch non-fatal watches"
	case DCGM_HEALTH_WATCH_NVSWITCH_FATAL:
		return "NVSwitch fatal watches"
	case DCGM_HEALTH_WATCH_CONNECTX:
		return "ConnectX watches"
	case DCGM_HEALTH_WATCH_IMEX:
		return "IMEX watches"
	}
	return "N/A"
}
