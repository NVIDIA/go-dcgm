package dcgm

/*
#include "dcgm_errors.h"
*/
import "C"

import "errors"

// ErrInvalidMode represents an error indicating that an invalid mode was used
var ErrInvalidMode = errors.New("invalid mode")

// ErrorSeverity describes the action required for a DCGM health or diagnostic error.
type ErrorSeverity int

const (
	// DCGM_ERROR_NONE indicates that no action is required.
	DCGM_ERROR_NONE ErrorSeverity = C.DCGM_ERROR_NONE
	// DCGM_ERROR_MONITOR indicates that the GPU can run workloads but should be monitored.
	DCGM_ERROR_MONITOR ErrorSeverity = C.DCGM_ERROR_MONITOR
	// DCGM_ERROR_ISOLATE indicates that the GPU should be isolated from workloads.
	DCGM_ERROR_ISOLATE ErrorSeverity = C.DCGM_ERROR_ISOLATE
	// DCGM_ERROR_UNKNOWN indicates that the error code is not recognized.
	DCGM_ERROR_UNKNOWN ErrorSeverity = C.DCGM_ERROR_UNKNOWN
	// DCGM_ERROR_TRIAGE indicates that the error should be triaged.
	DCGM_ERROR_TRIAGE ErrorSeverity = C.DCGM_ERROR_TRIAGE
	// DCGM_ERROR_CONFIG indicates that the error may be resolved through configuration.
	DCGM_ERROR_CONFIG ErrorSeverity = C.DCGM_ERROR_CONFIG
	// DCGM_ERROR_RESET indicates that the GPU should be drained and reset.
	DCGM_ERROR_RESET ErrorSeverity = C.DCGM_ERROR_RESET
)

// ErrorCategory identifies the subsystem associated with a DCGM health or diagnostic error.
type ErrorCategory int

const (
	// DCGM_FR_EC_NONE indicates no error category.
	DCGM_FR_EC_NONE ErrorCategory = C.DCGM_FR_EC_NONE
	// DCGM_FR_EC_PERF_THRESHOLD indicates a performance-threshold error.
	DCGM_FR_EC_PERF_THRESHOLD ErrorCategory = C.DCGM_FR_EC_PERF_THRESHOLD
	// DCGM_FR_EC_PERF_VIOLATION indicates a performance-violation error.
	DCGM_FR_EC_PERF_VIOLATION ErrorCategory = C.DCGM_FR_EC_PERF_VIOLATION
	// DCGM_FR_EC_SOFTWARE_CONFIG indicates a software-configuration error.
	DCGM_FR_EC_SOFTWARE_CONFIG ErrorCategory = C.DCGM_FR_EC_SOFTWARE_CONFIG
	// DCGM_FR_EC_SOFTWARE_LIBRARY indicates a software-library error.
	DCGM_FR_EC_SOFTWARE_LIBRARY ErrorCategory = C.DCGM_FR_EC_SOFTWARE_LIBRARY
	// DCGM_FR_EC_SOFTWARE_XID indicates a software XID error.
	DCGM_FR_EC_SOFTWARE_XID ErrorCategory = C.DCGM_FR_EC_SOFTWARE_XID
	// DCGM_FR_EC_SOFTWARE_CUDA indicates a CUDA software error.
	DCGM_FR_EC_SOFTWARE_CUDA ErrorCategory = C.DCGM_FR_EC_SOFTWARE_CUDA
	// DCGM_FR_EC_SOFTWARE_EUD indicates an EUD software error.
	DCGM_FR_EC_SOFTWARE_EUD ErrorCategory = C.DCGM_FR_EC_SOFTWARE_EUD
	// DCGM_FR_EC_SOFTWARE_OTHER indicates another software error.
	DCGM_FR_EC_SOFTWARE_OTHER ErrorCategory = C.DCGM_FR_EC_SOFTWARE_OTHER
	// DCGM_FR_EC_HARDWARE_THERMAL indicates a thermal hardware error.
	DCGM_FR_EC_HARDWARE_THERMAL ErrorCategory = C.DCGM_FR_EC_HARDWARE_THERMAL
	// DCGM_FR_EC_HARDWARE_MEMORY indicates a memory hardware error.
	DCGM_FR_EC_HARDWARE_MEMORY ErrorCategory = C.DCGM_FR_EC_HARDWARE_MEMORY
	// DCGM_FR_EC_HARDWARE_NVLINK indicates an NVLink hardware error.
	DCGM_FR_EC_HARDWARE_NVLINK ErrorCategory = C.DCGM_FR_EC_HARDWARE_NVLINK
	// DCGM_FR_EC_HARDWARE_NVSWITCH indicates an NVSwitch hardware error.
	DCGM_FR_EC_HARDWARE_NVSWITCH ErrorCategory = C.DCGM_FR_EC_HARDWARE_NVSWITCH
	// DCGM_FR_EC_HARDWARE_PCIE indicates a PCIe hardware error.
	DCGM_FR_EC_HARDWARE_PCIE ErrorCategory = C.DCGM_FR_EC_HARDWARE_PCIE
	// DCGM_FR_EC_HARDWARE_POWER indicates a power hardware error.
	DCGM_FR_EC_HARDWARE_POWER ErrorCategory = C.DCGM_FR_EC_HARDWARE_POWER
	// DCGM_FR_EC_HARDWARE_OTHER indicates another hardware error.
	DCGM_FR_EC_HARDWARE_OTHER ErrorCategory = C.DCGM_FR_EC_HARDWARE_OTHER
	// DCGM_FR_EC_INTERNAL_OTHER indicates an internal DCGM error.
	DCGM_FR_EC_INTERNAL_OTHER ErrorCategory = C.DCGM_FR_EC_INTERNAL_OTHER
)

// ErrorMeta contains the metadata associated with a DCGM health or diagnostic error code.
type ErrorMeta struct {
	ErrorID       HealthCheckErrorCode
	MessageFormat string
	Suggestion    string
	Severity      ErrorSeverity
	Category      ErrorCategory
}

func getErrorMeta(code HealthCheckErrorCode) *ErrorMeta {
	if code >= HealthCheckErrorCode(C.DCGM_FR_ERROR_SENTINEL) {
		return nil
	}

	meta := C.dcgmGetErrorMeta(C.dcgmError_t(code))
	if meta == nil {
		return nil
	}

	return &ErrorMeta{
		ErrorID:       HealthCheckErrorCode(meta.errorId),
		MessageFormat: goStringOrEmpty(meta.msgFormat),
		Suggestion:    goStringOrEmpty(meta.suggestion),
		Severity:      ErrorSeverity(meta.severity),
		Category:      ErrorCategory(meta.category),
	}
}

func goStringOrEmpty(value *C.char) string {
	if value == nil {
		return ""
	}

	return C.GoString(value)
}
