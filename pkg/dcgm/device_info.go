package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"slices"
	"unsafe"

	"github.com/bits-and-blooms/bitset"
)

// PCIInfo contains PCI bus related information for a GPU device
type PCIInfo struct {
	BusID     string
	BAR1      uint  // MB
	FBTotal   uint  // MB
	Bandwidth int64 // MB/s
}

// DeviceIdentifiers contains various identification information for a GPU device
type DeviceIdentifiers struct {
	Brand               string
	Model               string
	Serial              string
	Vbios               string
	InforomImageVersion string
	DriverVersion       string
}

// Device represents a GPU device and its properties
type Device struct {
	GPU           uint
	DCGMSupported string
	UUID          string
	Power         uint // W
	PCI           PCIInfo
	Identifiers   DeviceIdentifiers
	Topology      []P2PLink
	CPUAffinity   string
}

// NvLinkP2PStatus represents the state of NvLinks between the GPU pairs
type NvLinkP2PStatus struct {
	numGpus uint
	Gpus    [][]Link_State
}

// getAllDeviceCount counts all GPUs on the system
func getAllDeviceCount() (gpuCount uint, err error) {
	var (
		gpuIDList [C.DCGM_MAX_NUM_DEVICES]C.uint
		count     C.int
	)

	result := C.dcgmGetAllDevices(handle.handle, &gpuIDList[0], &count)
	if err = errorString(result); err != nil {
		return gpuCount, fmt.Errorf("error getting devices count: %w", err)
	}
	if err = validateNativeCount(int(count), len(gpuIDList), "device"); err != nil {
		return gpuCount, fmt.Errorf("error getting devices count: %w", err)
	}
	gpuCount = uint(count)
	return
}

// getEntityGroupEntities returns the entity IDs in one native entity group.
func getEntityGroupEntities(entityGroup Field_Entity_Group) ([]uint, error) {
	var err error
	var pEntities [C.DCGM_ENTITY_GROUP_ENTITIES_MAX]C.uint
	count := C.int(len(pEntities))

	result := C.dcgmGetEntityGroupEntities(handle.handle, C.dcgm_field_entity_group_t(entityGroup), &pEntities[0], &count, 0)
	if err = errorString(result); err != nil {
		return nil, fmt.Errorf("error getting entity count: %w", err)
	}

	return entityGroupEntitiesFromNative(&pEntities, count)
}

// entityGroupEntitiesFromNative validates and copies native entity IDs into Go-owned memory.
func entityGroupEntitiesFromNative(pEntities *[C.DCGM_ENTITY_GROUP_ENTITIES_MAX]C.uint, count C.int) ([]uint, error) {
	entityCount := int(count)
	if err := validateNativeCount(entityCount, len(pEntities), "entity"); err != nil {
		return nil, err
	}

	entities := make([]uint, entityCount)
	for i := 0; i < entityCount; i++ {
		entities[i] = uint(pEntities[i])
	}
	return entities, nil
}

// getSupportedDevices returns DCGM supported GPUs
func getSupportedDevices() (gpus []uint, err error) {
	return cgoAdapter{}.getSupportedDevices()
}

// getSupportedDevices queries DCGM for GPU IDs supported by the current session.
func (cgoAdapter) getSupportedDevices() (gpus []uint, err error) {
	var gpuIDList [C.DCGM_MAX_NUM_DEVICES]C.uint
	var count C.int

	result := C.dcgmGetAllSupportedDevices(handle.handle, &gpuIDList[0], &count)
	if err := errorString(result); err != nil {
		return gpus, err
	}

	return supportedDevicesFromNative(&gpuIDList, count)
}

// supportedDevicesFromNative validates and copies supported GPU IDs into Go-owned memory.
func supportedDevicesFromNative(gpuIDList *[C.DCGM_MAX_NUM_DEVICES]C.uint, count C.int) ([]uint, error) {
	numGPUs := int(count)
	if err := validateNativeCount(numGPUs, len(gpuIDList), "supported GPU"); err != nil {
		return nil, err
	}

	gpus := make([]uint, numGPUs)
	for i := 0; i < numGPUs; i++ {
		gpus[i] = uint(gpuIDList[i])
	}
	return gpus, nil
}

// validateNativeCount rejects negative counts and counts larger than a native fixed array.
func validateNativeCount(count, capacity int, description string) error {
	if count < 0 || count > capacity {
		return fmt.Errorf("invalid %s count %d, expected 0 through %d", description, count, capacity)
	}
	return nil
}

// getPciBandwidth calculates a GPU's PCIe bandwidth using the native adapter.
func (api cgoAdapter) getPciBandwidth(gpuID uint) (bandwidth int64, cleanupErr, err error) {
	return getPciBandwidthWithOps(api, gpuID)
}

func getPciBandwidthWithOps(api deviceStatusOps, gpuID uint) (bandwidth int64, cleanupErr, err error) {
	fields := []Short{C.DCGM_FI_DEV_PCIE_MAX_LINK_GEN, C.DCGM_FI_DEV_PCIE_MAX_LINK_WIDTH}
	values, cleanupErr, err := readTemporaryDeviceFields(api, gpuID, fields, "pciBandwidth", "PCIe bandwidth")
	if err != nil {
		return 0, cleanupErr, err
	}

	integers, err := requiredInt64Values(values, fields)
	if err != nil {
		return 0, cleanupErr, fmt.Errorf("get PCIe bandwidth for GPU %d: %w", gpuID, err)
	}
	for i, value := range integers {
		if IsInt64Blank(value) {
			return 0, cleanupErr, fmt.Errorf("get PCIe bandwidth for GPU %d: field %d: value is unavailable", gpuID, fields[i])
		}
	}

	return pcieBandwidth(integers[0], integers[1]), cleanupErr, nil
}

// pcieBandwidth converts a PCIe generation and lane width into megabytes per second.
func pcieBandwidth(generation, width int64) int64 {
	perLane := map[int64]int64{
		1: 250, // MB/s
		2: 500,
		3: 985,
		4: 1969,
	}
	return perLane[generation] * width
}

// getCPUAffinity formats a GPU's CPU affinity using the native adapter.
func (api cgoAdapter) getCPUAffinity(gpuID uint) (affinity string, cleanupErr, err error) {
	return getCPUAffinityWithOps(api, gpuID)
}

func getCPUAffinityWithOps(api deviceStatusOps, gpuID uint) (affinity string, cleanupErr, err error) {
	fields := []Short{
		C.DCGM_FI_DEV_CPU_AFFINITY_0,
		C.DCGM_FI_DEV_CPU_AFFINITY_1,
		C.DCGM_FI_DEV_CPU_AFFINITY_2,
		C.DCGM_FI_DEV_CPU_AFFINITY_3,
	}
	values, cleanupErr, err := readTemporaryDeviceFields(api, gpuID, fields, "cpuAff", "CPU affinity")
	if err != nil {
		return "N/A", cleanupErr, err
	}

	integers, err := requiredInt64Values(values, fields)
	if err != nil {
		return "N/A", cleanupErr, fmt.Errorf("get CPU affinity for GPU %d: %w", gpuID, err)
	}
	bits := make([]uint64, len(integers))
	for i, value := range integers {
		bits[i] = uint64(value)
	}

	return bitset.From(bits).String(), cleanupErr, nil
}

// readTemporaryDeviceFields releases the watch group before its field group.
// Separate errors let GetDeviceInfo keep valid data when only cleanup failed.
func readTemporaryDeviceFields(api deviceStatusOps, gpuID uint, fields []Short, groupPrefix, operation string) (
	values []FieldValue_v1, cleanupErr, err error,
) {
	fieldsID, err := api.fieldGroupCreate(groupPrefix+"Fields"+rand.Text(), fields)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if destroyErr := api.fieldGroupDestroy(fieldsID); destroyErr != nil {
			cleanupErr = errors.Join(cleanupErr, fmt.Errorf("destroy %s field group: %w", operation, destroyErr))
		}
	}()

	groupID, err := api.watchFields(gpuID, fieldsID, groupPrefix+rand.Text())
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if destroyErr := api.destroyGroup(groupID); destroyErr != nil {
			cleanupErr = fmt.Errorf("destroy %s GPU group: %w", operation, destroyErr)
		}
	}()

	values, err = api.getLatestValuesForFields(gpuID, fields)
	if err != nil {
		return nil, nil, fmt.Errorf("get %s for GPU %d: %w", operation, gpuID, err)
	}
	return values, nil, nil
}

// requiredInt64Values checks field metadata without interpreting integer bit patterns.
func requiredInt64Values(values []FieldValue_v1, fields []Short) ([]int64, error) {
	result := make([]int64, len(fields))
	for i, field := range fields {
		index := slices.IndexFunc(values, func(value FieldValue_v1) bool { return value.FieldID == field })
		if index < 0 {
			return nil, fmt.Errorf("field %d was not returned by DCGM", field)
		}
		value := &values[index]
		if value.Status != DCGM_ST_OK {
			return nil, fmt.Errorf("read field %d: %w", field,
				&Error{msg: fmt.Sprintf("field status %d", value.Status), Code: C.dcgmReturn_t(value.Status)})
		}
		if value.FieldType != DCGM_FT_INT64 {
			return nil, fmt.Errorf("field %d: expected an integer value", field)
		}
		result[i] = value.Int64()
	}
	return result, nil
}

// getDeviceInfo assembles GPU details using the production native adapter.
func getDeviceInfo(gpuID uint) (Device, error) {
	return getDeviceInfoWithOps(cgoAdapter{}, gpuID)
}

// readDeviceAttributes converts native GPU attributes into values used to build a Device.
func (api cgoAdapter) readDeviceAttributes(gpuID uint) (deviceAttributes, error) {
	var device C.dcgmDeviceAttributes_v3
	device.version = makeVersion3(unsafe.Sizeof(device))

	result := C.dcgmGetDeviceAttributes(handle.handle, C.uint(gpuID), (*C.dcgmDeviceAttributes_t)(unsafe.Pointer(&device)))
	if err := errorString(result); err != nil {
		return deviceAttributes{}, err
	}

	attrs := deviceAttributesFromNative(&device)
	if attrs.bar1 == 0 || attrs.fbTotal == 0 {
		if err := api.readLegacyDeviceMemory(gpuID, &attrs); err != nil {
			return deviceAttributes{}, fmt.Errorf("get memory capacities for GPU %d: %w", gpuID, err)
		}
	}
	return attrs, nil
}

// readLegacyDeviceMemory reads missing capacities from a 4.6.1 hostengine.
// The 4.8.0 client uses newer fields to fill device attributes and reports
// zeros when the older hostengine does not recognize those fields.
func (api cgoAdapter) readLegacyDeviceMemory(gpuID uint, attrs *deviceAttributes) error {
	mux.Lock()
	defer mux.Unlock()

	if !usesLegacyDeviceMemory() {
		return nil
	}

	// Keep the live read in the same DCGM session as the version checks.
	return readLegacyDeviceMemoryFields(gpuID, attrs, api.entitiesGetLatestValues)
}

// readLegacyDeviceMemoryFields fills the missing capacities from one live field read.
func readLegacyDeviceMemoryFields(gpuID uint, attrs *deviceAttributes,
	read func([]GroupEntityPair, []Short, uint) ([]FieldValue_v2, error),
) error {
	var fields []Short
	var capacities []*uint
	if attrs.bar1 == 0 {
		fields = append(fields, DCGM_FI_DEV_BAR1_TOTAL)
		capacities = append(capacities, &attrs.bar1)
	}
	if attrs.fbTotal == 0 {
		fields = append(fields, DCGM_FI_DEV_FB_TOTAL)
		capacities = append(capacities, &attrs.fbTotal)
	}
	values, err := read(
		[]GroupEntityPair{{EntityGroupId: FE_GPU, EntityId: gpuID}}, fields, DCGM_FV_FLAG_LIVE_DATA,
	)
	if err != nil {
		return err
	}
	if len(values) != len(fields) {
		return fmt.Errorf("DCGM returned %d of %d memory fields", len(values), len(fields))
	}
	for i, field := range fields {
		capacity, err := deviceMemoryCapacity(&values[i], field, gpuID)
		if err != nil {
			return err
		}
		*capacities[i] = capacity
	}
	return nil
}

// deviceMemoryCapacity reads a legacy memory field in MiB and rejects unavailable or out-of-range values.
func deviceMemoryCapacity(value *FieldValue_v2, field Short, gpuID uint) (uint, error) {
	if value.FieldID != field || value.EntityGroupId != FE_GPU || value.EntityID != gpuID {
		return 0, fmt.Errorf("DCGM did not return memory field %d for GPU %d", field, gpuID)
	}
	if value.Status != DCGM_ST_OK {
		code := C.dcgmReturn_t(value.Status)
		return 0, fmt.Errorf("read memory field %d: %w", field, &Error{msg: dcgmErrorText(code), Code: code})
	}
	if value.FieldType != DCGM_FT_INT64 {
		return 0, fmt.Errorf("memory field %d: expected an integer value", field)
	}
	mib := value.Int64()
	if IsInt64Blank(mib) {
		return 0, fmt.Errorf("memory field %d: DCGM returned a blank capacity", field)
	}
	// The legacy attributes store memory capacities in 32-bit unsigned fields.
	if mib < 0 || mib > math.MaxUint32 {
		return 0, fmt.Errorf("memory field %d: capacity %d MiB out of range", field, mib)
	}
	return uint(mib), nil
}

// usesLegacyDeviceMemory reports whether the client and hostengine need the
// legacy memory fields. The caller holds mux.
func usesLegacyDeviceMemory() bool {
	client, err := versionInfo()
	if err != nil {
		return false
	}
	clientVersion, known := runtimeVersion(client)
	if !known || clientVersion != [3]uint64{4, 8, 0} {
		return false
	}
	peer, err := hostengineVersionInfo()
	if err != nil {
		return false
	}
	peerVersion, known := runtimeVersion(peer)
	return known && peerVersion == [3]uint64{4, 6, 1}
}

// deviceAttributesFromNative converts native GPU attributes into Go-owned values.
func deviceAttributesFromNative(device *C.dcgmDeviceAttributes_v3) deviceAttributes {
	return deviceAttributes{
		busID:   *stringPtr(&device.identifiers.pciBusId[0]),
		uuid:    *stringPtr(&device.identifiers.uuid[0]),
		power:   *uintPtr(device.powerLimits.defaultPowerLimit),
		bar1:    *uintPtr(device.memoryUsage.bar1Total),
		fbTotal: *uintPtr(device.memoryUsage.fbTotal),
		identifiers: DeviceIdentifiers{
			Brand:               *stringPtr(&device.identifiers.brandName[0]),
			Model:               *stringPtr(&device.identifiers.deviceName[0]),
			Serial:              *stringPtr(&device.identifiers.serial[0]),
			Vbios:               *stringPtr(&device.identifiers.vbios[0]),
			InforomImageVersion: *stringPtr(&device.identifiers.inforomImageVersion[0]),
			DriverVersion:       *stringPtr(&device.identifiers.driverVersion[0]),
		},
	}
}

// getDeviceInfoWithOps builds a Device from GPU attributes.
// It reads CPU affinity, topology, and PCIe bandwidth only for supported GPUs with EntityStatusOk.
// Errors from these additional reads are returned.
func getDeviceInfoWithOps(api deviceInfoOps, gpuID uint) (_ Device, err error) {
	attrs, err := api.readDeviceAttributes(gpuID)
	if err != nil {
		return Device{}, err
	}

	gpus, err := api.getSupportedDevices()
	if err != nil {
		return Device{}, err
	}

	supported := slices.Contains(gpus, gpuID)
	if api.getGPUStatus(gpuID) != EntityStatusOk {
		supported = false
	}

	var (
		cpuAffinity string
		topology    []P2PLink
		bandwidth   int64
	)
	if supported {
		var affinityCleanup, bandwidthCleanup error
		defer func() {
			if affinityCleanup != nil || bandwidthCleanup != nil {
				err = errors.Join(err, affinityCleanup, bandwidthCleanup)
			}
		}()

		cpuAffinity, affinityCleanup, err = api.getCPUAffinity(gpuID)
		if err != nil {
			return Device{}, err
		}
		topology, err = api.getDeviceTopology(gpuID)
		if err != nil {
			return Device{}, err
		}
		bandwidth, bandwidthCleanup, err = api.getPciBandwidth(gpuID)
		if err != nil {
			return Device{}, err
		}
	}

	supportedText := "No"
	if supported {
		supportedText = "Yes"
	}

	return Device{
		GPU:           gpuID,
		DCGMSupported: supportedText,
		UUID:          attrs.uuid,
		Power:         attrs.power,
		PCI: PCIInfo{
			BusID:     attrs.busID,
			BAR1:      attrs.bar1,
			FBTotal:   attrs.fbTotal,
			Bandwidth: bandwidth,
		},
		Identifiers: attrs.identifiers,
		Topology:    topology,
		CPUAffinity: cpuAffinity,
	}, nil
}

// getNvLinkP2PStatus reads and converts the native GPU-to-GPU link matrix.
func getNvLinkP2PStatus() (NvLinkP2PStatus, error) {
	var linkStatus C.dcgmNvLinkP2PStatus_v1
	linkStatus.version = makeVersion1(unsafe.Sizeof(linkStatus))

	result := C.dcgmGetNvLinkP2PStatus(handle.handle, &linkStatus)
	if result == C.DCGM_ST_NOT_SUPPORTED {
		return NvLinkP2PStatus{}, nil
	}

	if result != C.DCGM_ST_OK {
		return NvLinkP2PStatus{}, &Error{msg: dcgmErrorText(result), Code: result}
	}

	return nvLinkP2PStatusFromNative(&linkStatus)
}

// nvLinkP2PStatusFromNative validates and converts a native GPU-to-GPU link matrix.
func nvLinkP2PStatusFromNative(linkStatus *C.dcgmNvLinkP2PStatus_v1) (NvLinkP2PStatus, error) {
	numGPUs := int(linkStatus.numGpus)
	if err := validateNativeCount(numGPUs, len(linkStatus.gpus), "NVLink GPU"); err != nil {
		return NvLinkP2PStatus{}, err
	}

	links := make([][]Link_State, numGPUs)
	for i := range links {
		links[i] = make([]Link_State, numGPUs)
	}

	for i := 0; i < numGPUs; i++ {
		for j := 0; j < numGPUs; j++ {
			links[i][j] = Link_State(linkStatus.gpus[i].linkStatus[j])
		}
	}

	return NvLinkP2PStatus{
		numGpus: uint(numGPUs),
		Gpus:    links,
	}, nil
}
