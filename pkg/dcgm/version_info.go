/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
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
#include <dlfcn.h>
#include "dcgm_agent.h"
#include "dcgm_structs.h"

// Call the client version function from the selected library if it exports the symbol.
static dcgmReturn_t loadedClientVersion(void *library, dcgmVersionInfo_t *info) {
    if (library == NULL) {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }
    typedef dcgmReturn_t (*clientVersionFn)(dcgmVersionInfo_t *);
    clientVersionFn query = (clientVersionFn)dlsym(library, "dcgmVersionInfo");
    if (query == NULL) {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }
    return query(info);
}

// Query the connected hostengine through the client library that opened the session.
static dcgmReturn_t loadedHostengineVersion(void *library, dcgmHandle_t handle, dcgmVersionInfo_t *info) {
    if (library == NULL) {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }
    typedef dcgmReturn_t (*hostengineVersionFn)(dcgmHandle_t, dcgmVersionInfo_t *);
    hostengineVersionFn query = (hostengineVersionFn)dlsym(library, "dcgmHostengineVersionInfo");
    if (query == NULL) {
        return DCGM_ST_FUNCTION_NOT_FOUND;
    }
    return query(handle, info);
}

// Read status text from the same library as the failed version query.
static const char *loadedVersionErrorText(void *library, dcgmReturn_t status) {
    typedef const char *(*errorStringFn)(dcgmReturn_t);
    errorStringFn describe = (errorStringFn)dlsym(library, "errorString");
    return describe == NULL ? NULL : describe(status);
}
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

// VersionInfo holds DCGM build environment information.
// RawBuildInfoString contains key-value pairs (e.g. version, arch, buildid, commit)
// separated by semicolons; each pair is "key:value".
type VersionInfo struct {
	RawBuildInfoString string
}

func versionInfo() (VersionInfo, error) {
	var info C.dcgmVersionInfo_t
	info.version = makeVersion2(unsafe.Sizeof(info))
	status := C.loadedClientVersion(dcgmLibHandle, &info)
	return versionInfoFromNative(&info, status, "DCGM client version", "dcgmVersionInfo")
}

func hostengineVersionInfo() (VersionInfo, error) {
	const symbol = "dcgmHostengineVersionInfo"
	if dcgmInitCounter <= 0 {
		return VersionInfo{}, &FeatureUnavailableError{Feature: "DCGM hostengine version (requires an initialized session)", Symbol: symbol}
	}
	var info C.dcgmVersionInfo_t
	info.version = makeVersion2(unsafe.Sizeof(info))
	status := C.loadedHostengineVersion(dcgmLibHandle, handle.handle, &info)
	return versionInfoFromNative(&info, status, "DCGM hostengine version", symbol)
}

// versionInfoFromNative converts native status codes and copies the bounded build-info string.
func versionInfoFromNative(info *C.dcgmVersionInfo_t, status C.dcgmReturn_t, feature, symbol string) (VersionInfo, error) {
	if status == C.DCGM_ST_FUNCTION_NOT_FOUND {
		return VersionInfo{}, &FeatureUnavailableError{Feature: feature, Symbol: symbol}
	}
	if status != C.DCGM_ST_OK {
		// Use the same handle as the version query, including in binaries
		// that do not link DCGM at build time.
		message := C.GoString(C.loadedVersionErrorText(dcgmLibHandle, status))
		if message == "" {
			message = fmt.Sprintf("DCGM version query failed with status %d", status)
		}
		return VersionInfo{}, &Error{msg: message, Code: status}
	}
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&info.rawBuildInfoString[0])), len(info.rawBuildInfoString))
	buildInfo, _, terminated := bytes.Cut(raw, []byte{0})
	if !terminated {
		return VersionInfo{}, errors.New("DCGM returned an unterminated version string")
	}
	return VersionInfo{RawBuildInfoString: string(buildInfo)}, nil
}
