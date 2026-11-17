package dcgm

/*
#include <stdbool.h>
#include <dlfcn.h>
#include "dcgm_test_apis.h"

// The test header spells this dcgmEntityInjectFieldValue, but DCGM exports
// dcgmInjectEntityFieldValue. Resolve the exported entry point without adding
// an unresolved direct call against the header-only spelling.
static dcgmReturn_t injectEntityValueDynamic(void *library,
                                            dcgmHandle_t handle,
                                            dcgm_field_entity_group_t entityGroup,
                                            dcgm_field_eid_t entityId,
                                            dcgmInjectFieldValue_t *value)
{
    if (library == NULL)
    {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }
    typedef dcgmReturn_t (*entityInjectionFunction)(dcgmHandle_t,
                                                  dcgm_field_entity_group_t,
                                                  dcgm_field_eid_t,
                                                  dcgmInjectFieldValue_t *);

    entityInjectionFunction injectValue =
        (entityInjectionFunction)dlsym(library, "dcgmInjectEntityFieldValue");
    if (injectValue == NULL)
    {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }
    return injectValue(handle, entityGroup, entityId, value);
}
*/
import "C"

import (
	"time"
	"unsafe"
)

// injectEntityValueForTest populates a private host engine's field cache.
func injectEntityValueForTest(entity GroupEntityPair, fieldID Short, fieldType byte, value int64) error {
	const symbol = "dcgmInjectEntityFieldValue"
	if !dcgmSymbolAvailable(symbol) {
		return &FeatureUnavailableError{Feature: "test entity injection", Symbol: symbol}
	}

	native := C.dcgmInjectFieldValue_t{
		version:   C.dcgmInjectFieldValue_version1,
		fieldId:   C.ushort(fieldID),
		fieldType: C.ushort(fieldType),
		ts:        C.int64_t(time.Now().UnixMicro()),
	}

	// Store the sample in the union member selected by its native field type.
	if uint(fieldType) == DCGM_FT_DOUBLE {
		*(*C.double)(unsafe.Pointer(&native.value[0])) = C.double(value)
	} else {
		*(*C.longlong)(unsafe.Pointer(&native.value[0])) = C.longlong(value)
	}

	result := C.injectEntityValueDynamic(
		dcgmLibHandle,
		handle.handle,
		C.dcgm_field_entity_group_t(entity.EntityGroupId),
		C.dcgm_field_eid_t(entity.EntityId),
		&native,
	)
	if result == C.DCGM_ST_OK {
		return nil
	}

	return dcgmFeatureResult("test entity injection", symbol, result, dcgmErrorText(result))
}
