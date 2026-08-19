package dcgm

// This file contains helpers for benchmarking field value operations.
// These functions expose internal implementation details for performance testing only.
// They should not be used in production code.

/*
#include "dcgm_structs.h"
*/
import "C"
import "unsafe"

type testFieldValueSpec struct {
	fieldID   Short
	fieldType uint
	status    int
	timestamp int64
	payload   []byte
}

// makeTestCFields creates test C field values for benchmarking purposes only.
func makeTestCFields(count int) []C.dcgmFieldValue_v1 {
	return makeTestCFieldsFromSpec(count, testFieldValueSpec{
		fieldType: DCGM_FT_INT64,
		timestamp: 1000000,
	})
}

func makeTestCFieldsFromSpec(count int, spec testFieldValueSpec) []C.dcgmFieldValue_v1 {
	cfields := make([]C.dcgmFieldValue_v1, count)
	for i := range cfields {
		cfields[i].version = C.dcgmFieldValue_version1
		cfields[i].fieldId = C.ushort(spec.fieldID + Short(i))
		cfields[i].fieldType = C.ushort(spec.fieldType)
		cfields[i].status = C.int(spec.status)
		cfields[i].ts = C.int64_t(spec.timestamp + int64(i))
		copy(cfields[i].value[:], spec.payload)
	}
	return cfields
}

func makeTestCFieldsV2FromSpec(entities, fieldsPerEntity int, spec testFieldValueSpec) []C.dcgmFieldValue_v2 {
	cfields := make([]C.dcgmFieldValue_v2, entities*fieldsPerEntity)
	for entityID := range entities {
		for field := range fieldsPerEntity {
			i := entityID*fieldsPerEntity + field
			cfields[i].version = C.dcgmFieldValue_version2
			cfields[i].entityGroupId = C.dcgm_field_entity_group_t(FE_GPU)
			cfields[i].entityId = C.dcgm_field_eid_t(entityID)
			cfields[i].fieldId = C.ushort(spec.fieldID + Short(field))
			cfields[i].fieldType = C.ushort(spec.fieldType)
			cfields[i].status = C.int(spec.status)
			cfields[i].ts = C.int64_t(spec.timestamp + int64(entityID))
			copy(cfields[i].value[:], spec.payload)
		}
	}
	return cfields
}

// oldAppendApproach implements the pre-optimization approach for benchmark comparison.
// It creates an intermediate slice before appending, which causes an extra allocation.
func oldAppendApproach(dst []FieldValue_v2, entityGroup Field_Entity_Group, entityID uint, cfields []C.dcgmFieldValue_v1) []FieldValue_v2 {
	intermediate := make([]FieldValue_v2, len(cfields))
	for i := range cfields {
		intermediate[i] = FieldValue_v2{
			Version:       C.dcgmFieldValue_version2,
			EntityGroupId: entityGroup,
			EntityID:      entityID,
			FieldID:       Short(cfields[i].fieldId),
			FieldType:     uint(cfields[i].fieldType),
			Status:        int(cfields[i].status),
			TS:            int64(cfields[i].ts),
			Value:         cfields[i].value,
			StringValue:   nil,
		}
		if uint(cfields[i].fieldType) == DCGM_FT_STRING {
			intermediate[i].StringValue = stringPtr((*C.char)(unsafe.Pointer(&cfields[i].value[0])))
		}
	}
	return append(dst, intermediate...)
}
