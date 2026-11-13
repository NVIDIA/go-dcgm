/*
 * Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 */

package dcgm

import (
	"fmt"
	"math"
)

// LinkEntityPair converts a discovered link to its native DCGM key. parentID
// must be the native GPU or switch ID, not an NVML index or an output label.
// dcgm_link_t stores an 8-bit parent type, 16-bit port index and 8-bit parent
// ID. Zero is valid; values outside those bounds fail before conversion.
func LinkEntityPair(index uint, parentType Field_Entity_Group, parentID uint) (GroupEntityPair, error) {
	if parentType != FE_GPU && parentType != FE_SWITCH {
		return GroupEntityPair{}, fmt.Errorf("link parent group %d: %w", parentType, newBadParameterError())
	}

	if index > math.MaxUint16 || parentID > math.MaxUint8 {
		return GroupEntityPair{}, fmt.Errorf("link index %d or parent ID %d exceeds dcgm_link_t bounds: %w", index, parentID, newBadParameterError())
	}

	return GroupEntityPair{
		EntityGroupId: FE_LINK,
		EntityId:      uint(parentType) | index<<8 | parentID<<24,
	}, nil
}

// validateNativeTarget checks native key bounds and the supplied relationship.
// Non-link parents are optional because applicability needs no topology lookup.
func validateNativeTarget(entity, parent GroupEntityPair) error {
	if !validEntityGroup(entity.EntityGroupId) || entity.EntityId > math.MaxUint32 {
		return fmt.Errorf("invalid native entity group %d ID %d: %w", entity.EntityGroupId, entity.EntityId, newBadParameterError())
	}

	if !validEntityGroup(parent.EntityGroupId) || parent.EntityId > math.MaxUint32 ||
		(parent.EntityGroupId == FE_NONE && parent.EntityId != 0) ||
		(entity.EntityGroupId == FE_NONE && entity.EntityId != 0) {
		return fmt.Errorf("invalid native entity or parent identity: %w", newBadParameterError())
	}

	if entity.EntityGroupId != FE_LINK {
		return nil
	}

	index := (entity.EntityId >> 8) & math.MaxUint16
	packed, err := LinkEntityPair(index, parent.EntityGroupId, parent.EntityId)
	if err != nil {
		return err
	}

	if packed != entity {
		return fmt.Errorf("link ID %d does not match parent group %d ID %d: %w", entity.EntityId, parent.EntityGroupId, parent.EntityId, newBadParameterError())
	}

	return nil
}

// validEntityGroup excludes unknown and sentinel native entity groups.
func validEntityGroup(group Field_Entity_Group) bool {
	return group < FE_COUNT
}
