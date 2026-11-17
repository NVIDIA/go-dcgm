#ifndef GO_DCGM_FIELD_POLICY_IMPORT_H
#define GO_DCGM_FIELD_POLICY_IMPORT_H

#include <stddef.h>
#include <dlfcn.h>
#include "dcgm_agent.h"

// Resolve and invoke the same address: a direct call to an optional symbol can
// be linked to address zero when the link-time library does not provide it.
static inline dcgmReturn_t callPolicyImportV2(void *library,
                                            dcgmHandle_t handle,
                                            const char *path,
                                            unsigned int *parsed,
                                            unsigned int *skipped)
{
    if (library == NULL)
    {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }

    typedef dcgmReturn_t (*policyImportV2Function)(dcgmHandle_t,
                                                 const char *,
                                                 unsigned int *,
                                                 unsigned int *);
    policyImportV2Function importPolicies = (policyImportV2Function)dlsym(library, "dcgmPolicyImport_v2");
    if (importPolicies == NULL)
    {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }
    return importPolicies(handle, path, parsed, skipped);
}

#endif
