package dcgm

import (
	/*
		#include "dcgm_agent.h"
	*/
	"C"
	"fmt"
)

// GroupGetAllIDs lists user-defined DCGM groups.
func GroupGetAllIDs() ([]GroupHandle, error) { return cgoAdapter{}.groupGetAllIDs() }

// groupGetAllIDs uses the native adapter to read the group inventory.
func (api cgoAdapter) groupGetAllIDs() ([]GroupHandle, error) {
	return groupGetAllIDsWithCalls(api)
}

// groupGetAllIDsWithCalls checks the session, reads bounded IDs and converts the returned count.
func groupGetAllIDsWithCalls(api administrationNativeCalls) ([]GroupHandle, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmGroupGetAllIds"); err != nil {
		return nil, err
	}

	var ids [C.DCGM_MAX_NUM_GROUPS]C.dcgmGpuGrp_t
	count := C.uint(len(ids))
	if err := administrationNativeResult(api, "dcgmGroupGetAllIds", api.groupIDsNative(ids[:], &count)); err != nil {
		return nil, err
	}

	return groupIDsFromNative(ids[:], uint(count))
}

// groupIDsFromNative copies only the IDs covered by a valid native count.
func groupIDsFromNative(ids []C.dcgmGpuGrp_t, count uint) ([]GroupHandle, error) {
	if uint64(count) > uint64(len(ids)) {
		return nil, fmt.Errorf("group count %d exceeds capacity %d", count, len(ids))
	}

	groups := make([]GroupHandle, count)
	for i := range groups {
		groups[i] = GroupHandle{ids[i]}
	}

	return groups, nil
}

// RemoveFromGroup removes a GPU from a group owned by the caller.
func RemoveFromGroup(group GroupHandle, gpuID uint) error {
	if uint64(gpuID) > uint64(^uint32(0)) {
		return newBadParameterError()
	}

	return cgoAdapter{}.removeFromGroup(group, gpuID)
}

// removeFromGroup uses the native adapter to remove a GPU from an owned group.
func (api cgoAdapter) removeFromGroup(group GroupHandle, gpuID uint) error {
	return removeFromGroupWithCalls(api, group, gpuID)
}

// removeFromGroupWithCalls checks the session before removing a GPU and preserves its native error.
func removeFromGroupWithCalls(api administrationNativeCalls, group GroupHandle, gpuID uint) error {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmGroupRemoveDevice"); err != nil {
		return err
	}

	return administrationNativeResult(api, "dcgmGroupRemoveDevice", api.removeGPUFromGroupNative(group, gpuID))
}

// FieldGroupInfo describes a field group and its field IDs.
type FieldGroupInfo struct {
	Handle   FieldHandle
	Name     string
	FieldIDs []Short
}

// FieldGroupGetInfo reads the name and fields of one field group.
func FieldGroupGetInfo(group FieldHandle) (FieldGroupInfo, error) {
	return cgoAdapter{}.fieldGroupGetInfo(group)
}

// fieldGroupGetInfo uses the native adapter to read one field group.
func (api cgoAdapter) fieldGroupGetInfo(group FieldHandle) (FieldGroupInfo, error) {
	return fieldGroupGetInfoWithCalls(api, group)
}

// fieldGroupGetInfoWithCalls initializes the native request and validates its returned field count.
func fieldGroupGetInfoWithCalls(api administrationNativeCalls, group FieldHandle) (FieldGroupInfo, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmFieldGroupGetInfo"); err != nil {
		return FieldGroupInfo{}, err
	}

	response := C.dcgmFieldGroupInfo_t{
		version:      C.dcgmFieldGroupInfo_version,
		fieldGroupId: group.handle,
	}
	if err := administrationNativeResult(api, "dcgmFieldGroupGetInfo", api.fieldGroupInfoNative(&response)); err != nil {
		return FieldGroupInfo{}, err
	}

	return fieldGroupInfoFromNative(response)
}

// fieldGroupInfoFromNative copies a bounded native field list and name into Go values.
func fieldGroupInfoFromNative(response C.dcgmFieldGroupInfo_t) (FieldGroupInfo, error) {
	if uint64(response.numFieldIds) > uint64(len(response.fieldIds)) {
		return FieldGroupInfo{}, fmt.Errorf("field count %d exceeds capacity", response.numFieldIds)
	}

	result := FieldGroupInfo{
		Handle:   FieldHandle{response.fieldGroupId},
		Name:     boundedCString(response.fieldGroupName[:]),
		FieldIDs: make([]Short, int(response.numFieldIds)),
	}
	for i := range result.FieldIDs {
		result.FieldIDs[i] = Short(response.fieldIds[i])
	}

	return result, nil
}

// FieldGroupGetAll lists all field groups in the host engine.
func FieldGroupGetAll() ([]FieldGroupInfo, error) { return cgoAdapter{}.fieldGroupGetAll() }

// fieldGroupGetAll uses the native adapter to read the field-group inventory.
func (api cgoAdapter) fieldGroupGetAll() ([]FieldGroupInfo, error) {
	return fieldGroupGetAllWithCalls(api)
}

// fieldGroupGetAllWithCalls checks the session and converts every returned field group.
func fieldGroupGetAllWithCalls(api administrationNativeCalls) ([]FieldGroupInfo, error) {
	mux.Lock()
	defer mux.Unlock()
	if err := requireAdministrationNativeSymbol(api, "dcgmFieldGroupGetAll"); err != nil {
		return nil, err
	}

	response := C.dcgmAllFieldGroup_t{version: C.dcgmAllFieldGroup_version}
	if err := administrationNativeResult(api, "dcgmFieldGroupGetAll", api.fieldGroupsNative(&response)); err != nil {
		return nil, err
	}

	return fieldGroupsFromNative(response)
}

// fieldGroupsFromNative validates group and field counts before copying native data.
func fieldGroupsFromNative(response C.dcgmAllFieldGroup_t) ([]FieldGroupInfo, error) {
	if uint64(response.numFieldGroups) > uint64(len(response.fieldGroups)) {
		return nil, fmt.Errorf("field group count %d exceeds capacity", response.numFieldGroups)
	}

	groups := make([]FieldGroupInfo, int(response.numFieldGroups))
	for i := range groups {
		var err error
		groups[i], err = fieldGroupInfoFromNative(response.fieldGroups[i])
		if err != nil {
			return nil, fmt.Errorf("field group %d: %w", i, err)
		}
	}

	return groups, nil
}
