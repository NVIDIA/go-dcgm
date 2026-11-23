package dcgm

/*
#include <stdlib.h>
#include "dcgm_structs.h"
*/
import "C"

import (
	"math"
	"strings"
	"unsafe"
)

func uintPtr(c C.uint) *uint {
	i := uint(c)
	return &i
}

func uint64Ptr(c C.longlong) *uint64 {
	i := uint64(c)
	return &i
}

func int64Ptr(c C.longlong) *int64 {
	i := int64(c)
	return &i
}

func toInt64(c C.longlong) int64 {
	i := int64(c)
	return i
}

func dblToFloat(val C.double) *float64 {
	i := float64(val)
	return &i
}

func stringPtr(c *C.char) *string {
	s := C.GoString(c)
	return &s
}

// Error represents an error returned by the DCGM library
type Error struct {
	msg  string         // description of error
	Code C.dcgmReturn_t // dcgmReturn_t value of error
}

func (e *Error) Error() string { return e.msg }

func errorString(result C.dcgmReturn_t) error {
	if result == C.DCGM_ST_OK {
		return nil
	}
	return &Error{msg: dcgmErrorText(result), Code: result}
}

// dcgmErrorText converts DCGM's native status message into the error text used by Go callers.
func dcgmErrorText(result C.dcgmReturn_t) string {
	return C.GoString(C.errorString(result))
}

func freeCString(cStr *C.char) {
	C.free(unsafe.Pointer(cStr))
}

// IsStringBlank reports whether value starts with "<<<" and contains ">>>"
// before the first NUL byte, following DCGM_STR_IS_BLANK.
func IsStringBlank(value string) bool {
	if !strings.HasPrefix(value, "<<<") {
		return false
	}
	text, _, _ := strings.Cut(value, "\x00")
	return strings.Contains(text, ">>>")
}

// IsInt32Blank reports whether value is at or above DCGM's 32-bit blank threshold.
func IsInt32Blank(value int) bool {
	return int64(value) >= DCGM_FT_INT32_BLANK
}

// IsInt64Blank reports whether value is at or above DCGM's 64-bit blank threshold.
func IsInt64Blank(value int64) bool {
	return value >= DCGM_FT_INT64_BLANK
}

// IsFloat64Blank reports whether value is at or above DCGM's floating-point blank threshold.
func IsFloat64Blank(value float64) bool {
	return value >= DCGM_FT_FP64_BLANK
}

func makeVersion1(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 1<<24)
	return version
}

func makeVersion2(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 2<<24)
	return version
}

func makeVersion3(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 3<<24)
	return version
}

func makeVersion4(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 4<<24)
	return version
}

func makeVersion5(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 5<<24)
	return version
}

func makeVersion12(struct_type uintptr) C.uint {
	version := C.uint(struct_type | 12<<24)
	return version
}

func roundFloat(f *float64) *float64 {
	var val float64
	if f != nil {
		val = math.Round(*f)
	}
	return &val
}
