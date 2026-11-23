/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
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

/*
 * Test-only stand-in for libdcgm.so.4, used by error_native_test.go.
 *
 * Tests can omit either version function (-DNO_CLIENT, -DNO_HOST) or both to
 * check that the Go wrappers return an error for missing symbols. A version
 * call can also pause while a test checks that the wrapper holds the session lock.
 *
 * The unit tests compile against dcgm_agent.h to check the signatures of
 * dcgmVersionInfo, dcgmHostengineVersionInfo, and errorString.
 * dcgmInjectFieldValue and dcgmCreateFakeEntities in dcgm_test_apis.h
 * control cached samples and entities. The loader tests control whether a
 * symbol exists, what it returns, and when the native call finishes.
 */

#include "dcgm_agent.h"
#include <errno.h>
#include <limits.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#ifndef VERSION_ERROR
// Signal native entry and wait for the Go test to release the call.
static void waitForVersionRelease(void) {
    const char *fd = getenv("GO_DCGM_VERSION_FD");
    if (fd == NULL || *fd == '\0') {
        return;
    }
    char *end;
    errno = 0;
    long descriptor = strtol(fd, &end, 10);
    if (errno == ERANGE || end == fd || *end != '\0' || descriptor <= STDERR_FILENO || descriptor > INT_MAX) {
        return;
    }

    char token = 1;
    int socket = (int)descriptor;
    if (write(socket, &token, 1) != 1 || read(socket, &token, 1) != 1) {
        abort();
    }
}
#endif

#if defined(VERSION_ERROR) && !defined(NO_ERROR_TEXT)
const char *errorString(dcgmReturn_t status) {
    return status == DCGM_ST_VER_MISMATCH ? "fixture version mismatch" : "unexpected status";
}
#endif

static dcgmReturn_t writeVersion(dcgmVersionInfo_t *info, const char *version) {
#ifdef VERSION_ERROR
    return DCGM_ST_VER_MISMATCH;
#else
#ifdef UNTERMINATED_VERSION
    memset(info->rawBuildInfoString, 'x', sizeof(info->rawBuildInfoString));
#else
    strcpy(info->rawBuildInfoString, version);
#endif
    waitForVersionRelease();
    return DCGM_ST_OK;
#endif
}

#ifndef CLIENT_VERSION
#define CLIENT_VERSION "4.8.0"
#endif
#ifndef HOST_VERSION
#define HOST_VERSION "4.7.0"
#endif

#ifndef NO_CLIENT
dcgmReturn_t dcgmVersionInfo(dcgmVersionInfo_t *info) {
    return writeVersion(info, "version:" CLIENT_VERSION ";");
}
#endif
#ifndef NO_HOST
dcgmReturn_t dcgmHostengineVersionInfo(dcgmHandle_t h, dcgmVersionInfo_t *info) {
    return writeVersion(info, "version:" HOST_VERSION ";");
}
#endif
