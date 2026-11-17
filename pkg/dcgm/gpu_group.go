package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"context"
	"fmt"
)

// DCGM_GROUP_MAX_ENTITIES represents the maximum number of entities allowed in a group
const (
	DCGM_GROUP_MAX_ENTITIES int = C.DCGM_GROUP_MAX_ENTITIES_V2
)

// GroupHandle represents a handle to a DCGM GPU group
type GroupHandle struct{ handle C.dcgmGpuGrp_t }

// SetHandle sets the internal group handle value
func (g *GroupHandle) SetHandle(val uintptr) {
	g.handle = C.dcgmGpuGrp_t(val)
}

// GetHandle returns the internal group handle value
func (g *GroupHandle) GetHandle() uintptr {
	return uintptr(g.handle)
}

// GroupAllGPUs returns a GroupHandle representing all GPUs in the system
func GroupAllGPUs() GroupHandle {
	return GroupHandle{C.DCGM_GROUP_ALL_GPUS}
}

// CreateGroup creates a new empty GPU group with the specified name.
//
// Important: Groups must be destroyed using DestroyGroup when no longer needed
// to prevent resource leaks in the DCGM library.
//
// Example:
//
//	group, err := dcgm.CreateGroup("myGroup")
//	if err != nil {
//	    return err
//	}
//	defer dcgm.DestroyGroup(group)
//
//	// Use the group...
func CreateGroup(groupName string) (goGroupId GroupHandle, err error) {
	return cgoAdapter{}.createGroup(groupName, false)
}

// NewDefaultGroup creates a new group with default GPUs and the specified name
func NewDefaultGroup(groupName string) (GroupHandle, error) {
	return cgoAdapter{}.createGroup(groupName, true)
}

// createGroup creates a native GPU group, either empty or populated with the default GPUs.
func (cgoAdapter) createGroup(groupName string, defaultGroup bool) (GroupHandle, error) {
	var cGroupID C.dcgmGpuGrp_t
	cname := C.CString(groupName)
	defer freeCString(cname)
	groupType := C.dcgmGroupType_t(C.DCGM_GROUP_EMPTY)
	if defaultGroup {
		groupType = C.DCGM_GROUP_DEFAULT
	}
	result := C.dcgmGroupCreate(handle.handle, groupType, cname, &cGroupID)
	if err := errorString(result); err != nil {
		return GroupHandle{}, fmt.Errorf("error creating group: %w", &Error{msg: err.Error(), Code: result})
	}

	return GroupHandle{cGroupID}, nil
}

// AddToGroup adds a GPU to an existing group
func AddToGroup(groupID GroupHandle, gpuID uint) (err error) {
	return cgoAdapter{}.addToGroup(groupID, gpuID)
}

// addToGroup adds one GPU to a native group owned by the caller.
func (cgoAdapter) addToGroup(groupID GroupHandle, gpuID uint) (err error) {
	result := C.dcgmGroupAddDevice(handle.handle, groupID.handle, C.uint(gpuID))
	if err = errorString(result); err != nil {
		return fmt.Errorf("error adding GPU %v to group: %w", gpuID, &Error{msg: err.Error(), Code: result})
	}

	return
}

// AddLinkEntityToGroup adds a link entity to the group. It checks the native
// GPU/switch parent and dcgm_link_t bounds through LinkEntityPair.
// Native failures wrap *Error, preserving the native status in Error.Code.
func AddLinkEntityToGroup(groupID GroupHandle, index uint, entityGroupID Field_Entity_Group, parentID uint) (err error) {
	return addLinkEntityToGroup(groupID, index, entityGroupID, parentID, AddEntityToGroup)
}

// addLinkEntityToGroup shares checked native identity conversion with link reads.
func addLinkEntityToGroup(groupID GroupHandle, index uint, parentType Field_Entity_Group, parentID uint,
	add func(GroupHandle, Field_Entity_Group, uint) error,
) error {
	entity, err := LinkEntityPair(index, parentType, parentID)
	if err != nil {
		return err
	}

	return add(groupID, entity.EntityGroupId, entity.EntityId)
}

// AddEntityToGroup adds an entity to an existing group.
// Native failures wrap *Error, preserving the native status in Error.Code.
func AddEntityToGroup(groupID GroupHandle, entityGroupID Field_Entity_Group, entityID uint) (err error) {
	result := C.dcgmGroupAddEntity(handle.handle, groupID.handle, C.dcgm_field_entity_group_t(entityGroupID),
		C.uint(entityID))
	if result != C.DCGM_ST_OK {
		err = &Error{msg: dcgmErrorText(result), Code: result}
		return fmt.Errorf("error adding entity group type %v, entity %v to group: %w", entityGroupID, entityID, err)
	}

	return
}

// RemoveEntityFromGroup removes an entity from an existing group.
// It returns an error if the group or entity is not found.
func RemoveEntityFromGroup(groupID GroupHandle, entityGroupID Field_Entity_Group, entityID uint) (err error) {
	result := C.dcgmGroupRemoveEntity(handle.handle, groupID.handle, C.dcgm_field_entity_group_t(entityGroupID),
		C.uint(entityID))
	if err = errorString(result); err != nil {
		return fmt.Errorf("error removing entity group type %v, entity %v from group: %w", entityGroupID, entityID, &Error{msg: err.Error(), Code: result})
	}

	return
}

// DestroyGroup destroys an existing GPU group
func DestroyGroup(groupID GroupHandle) (err error) {
	return cgoAdapter{}.destroyGroup(groupID)
}

// destroyGroup releases a native GPU group owned by the caller.
func (cgoAdapter) destroyGroup(groupID GroupHandle) (err error) {
	result := C.dcgmGroupDestroy(handle.handle, groupID.handle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("error destroying group: %w", &Error{msg: err.Error(), Code: result})
	}

	return
}

// GroupInfo contains information about a DCGM group
type GroupInfo struct {
	Version    uint32
	GroupName  string
	EntityList []GroupEntityPair
}

// GetGroupInfo retrieves information about a DCGM group
func GetGroupInfo(groupID GroupHandle) (*GroupInfo, error) {
	return cgoAdapter{}.getGroupInfo(groupID)
}

// getGroupInfo reads native group metadata and converts its entity membership.
func (cgoAdapter) getGroupInfo(groupID GroupHandle) (*GroupInfo, error) {
	response := C.dcgmGroupInfo_v3{
		version: C.dcgmGroupInfo_version3,
	}

	result := C.dcgmGroupGetInfo(handle.handle, groupID.handle, &response)
	if err := errorString(result); err != nil {
		return nil, &Error{msg: err.Error(), Code: result}
	}

	ret := GroupInfo{
		Version:    uint32(response.version),
		GroupName:  C.GoString(&response.groupName[0]),
		EntityList: make([]GroupEntityPair, response.count),
	}

	for i := 0; i < int(response.count); i++ {
		ret.EntityList[i].EntityId = uint(response.entityList[i].entityId)
		ret.EntityList[i].EntityGroupId = Field_Entity_Group(response.entityList[i].entityGroupId)
	}

	return &ret, nil
}

// CreateGroupWithContext creates a new group with a context
func CreateGroupWithContext(ctx context.Context, groupName string) (GroupHandle, error) {
	select {
	case <-ctx.Done():
		return GroupHandle{}, ctx.Err()
	default:
		return CreateGroup(groupName)
	}
}
