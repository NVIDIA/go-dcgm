package dcgm

/*
#include "dcgm_agent.h"
#include "dcgm_structs.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// MetricGroup represents a group of metrics for a specific GPU
type MetricGroup struct {
	Major    uint
	Minor    uint
	FieldIds []uint
}

// getSupportedMetricGroups asks DCGM which profiling metric groups are
// available for one GPU and converts the native response into Go values.
func getSupportedMetricGroups(gpuID uint) ([]MetricGroup, error) {
	var (
		groupInfo C.dcgmProfGetMetricGroups_t
		err       error
	)

	groupInfo.version = makeVersion3(unsafe.Sizeof(groupInfo))

	groupInfo.gpuId = C.uint(gpuID)

	result := C.dcgmProfGetSupportedMetricGroups(handle.handle, &groupInfo)

	if err = errorString(result); err != nil {
		return nil, &Error{msg: dcgmErrorText(result), Code: result}
	}

	return metricGroupsFromResponse(&groupInfo)
}

// metricGroupsFromResponse validates a native metric-group response and copies
// its groups and field IDs into Go-owned values.
func metricGroupsFromResponse(groupInfo *C.dcgmProfGetMetricGroups_t) ([]MetricGroup, error) {
	groupCount := uint(groupInfo.numMetricGroups)
	maxGroupCount := uint(len(groupInfo.metricGroups))
	if groupCount > maxGroupCount {
		return nil, fmt.Errorf("metric group response contains %d groups, maximum %d", groupCount, maxGroupCount)
	}

	groups := make([]MetricGroup, groupCount)
	for i := uint(0); i < groupCount; i++ {
		metricGroup := &groupInfo.metricGroups[i]
		groups[i].Major = uint(metricGroup.majorId)
		groups[i].Minor = uint(metricGroup.minorId)

		fieldCount := uint(metricGroup.numFieldIds)
		maxFieldCount := uint(len(metricGroup.fieldIds))
		if fieldCount > maxFieldCount {
			return nil, fmt.Errorf("metric group %d contains %d field IDs, maximum %d", i, fieldCount, maxFieldCount)
		}

		groups[i].FieldIds = make([]uint, fieldCount)
		for j := uint(0); j < fieldCount; j++ {
			groups[i].FieldIds[j] = uint(metricGroup.fieldIds[j])
		}
	}

	return groups, nil
}
