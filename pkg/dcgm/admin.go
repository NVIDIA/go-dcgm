/*
 * Copyright (c) 2020, NVIDIA CORPORATION.  All rights reserved.
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
#cgo linux LDFLAGS: -ldl -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files
#cgo darwin LDFLAGS: -ldl -Wl,-undefined,dynamic_lookup

#include <dlfcn.h>
#include <stdio.h>
#include "dcgm_agent.h"
#include "dcgm_structs.h"

// openLibraryWithReason loads a shared library with its symbols made global.
// On failure it copies the dynamic loader's message into reason. dlopen and
// dlerror must run on the same thread, so both happen here in one C call.
static void *openLibraryWithReason(const char *name, char *reason, size_t reasonSize) {
	void *library = dlopen(name, RTLD_LAZY | RTLD_GLOBAL);
	if (library == NULL) {
		const char *loaderError = dlerror();
		if (loaderError == NULL) {
			loaderError = "unknown dynamic loader error";
		}
		snprintf(reason, reasonSize, "%s", loaderError);
	}

	return library;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type mode int

// const for DCGM hostengine running modes: Embedded, Standalone or StartHostengine
const (
	Embedded mode = iota
	Standalone
	StartHostengine
)

type dcgmHandle struct{ handle C.dcgmHandle_t }

var (
	dcgmLibHandle        unsafe.Pointer
	stopMode             mode
	handle               dcgmHandle
	hostengineAsChildPid int
)

func initDCGM(m mode, args ...string) (err error) {
	return initDCGMWithCalls(cgoAdapter{}, m, args...)
}

// initDCGMWithCalls loads the DCGM library and starts or connects to a hostengine
// according to the selected mode.
func initDCGMWithCalls(api sessionCalls, m mode, args ...string) error {
	if err := api.loadLibrary(); err != nil {
		return err
	}

	stopMode = m
	switch m {
	case Embedded:
		return api.startEmbedded()
	case Standalone:
		return api.connectStandalone(args...)
	case StartHostengine:
		return api.startHostengine()
	default:
		panic(ErrInvalidMode)
	}
}

// dcgmLibraryName is the DCGM 4 client library that Init loads.
const dcgmLibraryName = "libdcgm.so.4"

// LibraryLoadError reports that Init could not load the DCGM library, for
// example because it is not installed, is not on the library search path, or
// needs another library that is missing. go-dcgm only supports DCGM 4.
type LibraryLoadError struct {
	// Library is the file name go-dcgm tried to load.
	Library string
	// Reason is the dynamic loader's explanation of the failure.
	Reason string
}

// Error describes the failure. It starts with "<library> not found", which
// callers have matched on in earlier releases.
func (e *LibraryLoadError) Error() string {
	return fmt.Sprintf("%s not found or could not be loaded: %s; go-dcgm requires DCGM 4", e.Library, e.Reason)
}

// loadLibrary opens libdcgm.so.4 and saves its handle for native symbol lookups.
func (cgoAdapter) loadLibrary() error {
	return loadLibraryFile(dcgmLibraryName)
}

// loadLibraryFile loads the named shared library with its symbols made global,
// so later cgo calls resolve against it, and saves its handle in dcgmLibHandle.
// On failure it clears the handle and returns a *LibraryLoadError with the
// loader's reason.
func loadLibraryFile(name string) error {
	cName := C.CString(name)
	defer freeCString(cName)

	var reason [512]C.char
	dcgmLibHandle = C.openLibraryWithReason(cName, &reason[0], C.size_t(len(reason)))
	if dcgmLibHandle == nil {
		return &LibraryLoadError{Library: name, Reason: C.GoString(&reason[0])}
	}

	return nil
}

func shutdown() (err error) {
	return shutdownWithCalls(cgoAdapter{})
}

// shutdownWithCalls stops the active session and closes the library even if stopping fails.
func shutdownWithCalls(api sessionCalls) (err error) {
	switch stopMode {
	case Embedded:
		err = api.stopEmbedded()
	case Standalone:
		err = api.disconnectStandalone()
	case StartHostengine:
		err = api.stopHostengine()
	}

	api.closeLibrary()
	return
}

// closeLibrary releases the library handle acquired by loadLibrary and clears
// it. Later symbol checks then report every symbol as unavailable instead of
// looking it up through a closed handle. Calling it again does nothing.
func (cgoAdapter) closeLibrary() {
	if dcgmLibHandle == nil {
		return
	}

	C.dlclose(dcgmLibHandle)
	dcgmLibHandle = nil
}

// startEmbedded initializes DCGM and starts an in-process hostengine.
func (cgoAdapter) startEmbedded() (err error) {
	result := C.dcgmInit()
	if err = errorString(result); err != nil {
		return fmt.Errorf("error initializing DCGM: %s", err)
	}

	var cHandle C.dcgmHandle_t
	result = C.dcgmStartEmbedded(C.DCGM_OPERATION_MODE_AUTO, &cHandle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("error starting nv-hostengine: %s", err)
	}
	handle = dcgmHandle{cHandle}
	return
}

// stopEmbedded stops the in-process hostengine and shuts down DCGM.
func (cgoAdapter) stopEmbedded() (err error) {
	result := C.dcgmStopEmbedded(handle.handle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("error stopping nv-hostengine: %s", err)
	}

	result = C.dcgmShutdown()
	if err = errorString(result); err != nil {
		return fmt.Errorf("error shutting down DCGM: %s", err)
	}
	return
}

type standaloneConnection struct {
	address    string
	socketFlag string
	useV3      bool
}

func standaloneConnectionArgs(args ...string) (standaloneConnection, error) {
	if len(args) > 0 && isDCGMConnectionString(args[0]) {
		return standaloneConnection{address: args[0], useV3: true}, nil
	}
	if len(args) < 2 {
		return standaloneConnection{}, errors.New("missing dcgm address and / or port")
	}
	return standaloneConnection{address: args[0], socketFlag: args[1]}, nil
}

func isDCGMConnectionString(address string) bool {
	lower := strings.ToLower(address)
	return strings.HasPrefix(lower, "tcp://") ||
		strings.HasPrefix(lower, "unix://") ||
		strings.HasPrefix(lower, "vsock://")
}

// dcgmSymbolAvailable reports whether the loaded libdcgm exports the named C
// function. Callers use it before calling an API that older DCGM releases lack.
// It returns false when no library is loaded.
func dcgmSymbolAvailable(symbol string) bool {
	if dcgmLibHandle == nil {
		return false
	}

	cSymbol := C.CString(symbol)
	defer freeCString(cSymbol)
	return C.dlsym(dcgmLibHandle, cSymbol) != nil
}

// connectStandalone initializes DCGM and connects to an external hostengine.
func (api cgoAdapter) connectStandalone(args ...string) (err error) {
	conn, err := standaloneConnectionArgs(args...)
	if err != nil {
		return err
	}
	if conn.useV3 {
		return connectStandaloneV3(conn.address)
	}
	return connectStandaloneV2(conn.address, conn.socketFlag)
}

func connectStandaloneV2(address, socketFlag string) (err error) {
	var (
		cHandle       C.dcgmHandle_t
		connectParams C.dcgmConnectV2Params_v2
	)
	sck, err := strconv.ParseUint(socketFlag, 10, 32)
	if err != nil {
		return fmt.Errorf("error parsing %s: %w", socketFlag, err)
	}

	result := C.dcgmInit()
	if err = errorString(result); err != nil {
		return fmt.Errorf("error initializing DCGM: %s", err)
	}

	addr := C.CString(address)
	defer freeCString(addr)
	connectParams.version = makeVersion2(unsafe.Sizeof(connectParams))

	connectParams.addressIsUnixSocket = C.uint(sck)

	result = C.dcgmConnect_v2(addr, &connectParams, &cHandle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("error connecting to nv-hostengine: %s", err)
	}

	handle = dcgmHandle{cHandle}
	return nil
}

func connectStandaloneV3(connectionString string) (err error) {
	const dcgmConnectV3Symbol = "dcgmConnect_v3"
	if !dcgmSymbolAvailable(dcgmConnectV3Symbol) {
		return fmt.Errorf("%s is not available in %s; DCGM connection strings require DCGM 4.5.0 or newer", dcgmConnectV3Symbol, dcgmLibraryName)
	}

	var (
		cHandle       C.dcgmHandle_t
		connectParams C.dcgmConnectV3Params_v1
	)

	result := C.dcgmInit()
	if err = errorString(result); err != nil {
		return fmt.Errorf("error initializing DCGM: %s", err)
	}

	cConnectionString := C.CString(connectionString)
	defer freeCString(cConnectionString)
	connectParams.version = makeVersion1(unsafe.Sizeof(connectParams))

	result = C.dcgmConnect_v3(cConnectionString, &connectParams, &cHandle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("error connecting to nv-hostengine: %s", err)
	}

	handle = dcgmHandle{cHandle}
	return nil
}

// disconnectStandalone disconnects from the external hostengine and shuts down DCGM.
func (cgoAdapter) disconnectStandalone() (err error) {
	result := C.dcgmDisconnect(handle.handle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("error disconnecting from nv-hostengine: %s", err)
	}

	result = C.dcgmShutdown()
	if err = errorString(result); err != nil {
		return fmt.Errorf("error shutting down DCGM: %s", err)
	}
	return
}

const (
	// maxUnixSocketPathBytes is the size of sun_path in Linux struct sockaddr_un,
	// including the terminating NUL.
	maxUnixSocketPathBytes = 108

	hostengineSocketName = "hostengine.sock"
)

// checkSocketPathLength returns an error if the path, plus the terminating NUL,
// does not fit in the sun_path field of a Unix socket address.
func checkSocketPathLength(socketPath string) error {
	if len(socketPath) >= maxUnixSocketPathBytes {
		return fmt.Errorf(
			"nv-hostengine socket path %q is %d bytes, longer than the %d-byte limit; set TMPDIR to a shorter directory",
			socketPath, len(socketPath), maxUnixSocketPathBytes-1,
		)
	}
	return nil
}

// createHostengineSocketDir makes a private directory (mode 0700) and returns
// it with the socket path inside it. It removes the directory and returns an
// error if the socket path does not fit in a Unix socket address.
func createHostengineSocketDir() (dir, socketPath string, err error) {
	dir, err = os.MkdirTemp("", "dcgm-hostengine-")
	if err != nil {
		return "", "", fmt.Errorf("error creating socket directory: %w", err)
	}

	socketPath = filepath.Join(dir, hostengineSocketName)
	if err = checkSocketPathLength(socketPath); err != nil {
		_ = os.RemoveAll(dir)
		return "", "", err
	}
	return dir, socketPath, nil
}

// startHostengine launches a child nv-hostengine process and connects to it.
func (cgoAdapter) startHostengine() (err error) {
	var (
		procAttr      syscall.ProcAttr
		cHandle       C.dcgmHandle_t
		connectParams C.dcgmConnectV2Params_v2
	)

	bin, err := exec.LookPath("nv-hostengine")
	if err != nil {
		return fmt.Errorf("error finding nv-hostengine: %s", err)
	}
	procAttr.Files = []uintptr{
		uintptr(syscall.Stdin),
		uintptr(syscall.Stdout),
		uintptr(syscall.Stderr),
	}
	procAttr.Sys = &syscall.SysProcAttr{Setpgid: true}

	socketDir, socketPath, err := createHostengineSocketDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(socketDir)

	connectArg := "--domain-socket"
	hostengineAsChildPid, err = syscall.ForkExec(bin, []string{bin, connectArg, socketPath}, &procAttr)
	if err != nil {
		return fmt.Errorf("error fork-execing nv-hostengine: %s", err)
	}

	result := C.dcgmInit()
	if err = errorString(result); err != nil {
		return fmt.Errorf("error initializing DCGM: %s", err)
	}

	connectParams.version = makeVersion2(unsafe.Sizeof(connectParams))
	isSocket := C.uint(1)
	connectParams.addressIsUnixSocket = isSocket
	cSockPath := C.CString(socketPath)
	defer freeCString(cSockPath)
	result = C.dcgmConnect_v2(cSockPath, &connectParams, &cHandle)
	if err = errorString(result); err != nil {
		return fmt.Errorf("error connecting to nv-hostengine: %s", err)
	}

	handle = dcgmHandle{cHandle}
	return
}

// AttachDriver attaches the driver to DCGM.
// This is used to reattach the driver after a DetachDriver call, typically when updating
// the driver without restarting DCGM.
// Requires DCGM 4.5.0 or later.
func AttachDriver() error {
	result := C.dcgmAttachDriver(handle.handle)
	if result != C.DCGM_ST_OK {
		return &Error{msg: dcgmErrorText(result), Code: result}
	}
	return nil
}

// DetachDriver detaches the driver from DCGM.
// This is used when you want to update the driver without restarting DCGM.
// After detaching, GPUs will not be accessible until AttachDriver is called.
// Requires DCGM 4.5.0 or later.
func DetachDriver() error {
	result := C.dcgmDetachDriver(handle.handle)
	if result != C.DCGM_ST_OK {
		return &Error{msg: dcgmErrorText(result), Code: result}
	}
	return nil
}

// stopHostengine disconnects the session and terminates the child hostengine process.
func (api cgoAdapter) stopHostengine() (err error) {
	if err = api.disconnectStandalone(); err != nil {
		return
	}

	// terminate nv-hostengine
	cmd := exec.Command("nv-hostengine", "--term")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("error terminating nv-hostengine: %s", err)
	}

	log.Println("Successfully terminated nv-hostengine.")

	return syscall.Kill(hostengineAsChildPid, syscall.SIGKILL)
}
