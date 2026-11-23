package tests

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const (
	bin       = "nvidia-smi"
	gpuArg    = "--id="
	queryArg  = "--query-gpu="
	formatArg = "--format=csv,noheader,nounits"
)

// Query executes nvidia-smi with the specified GPU ID and query parameters.
// It returns the query result as a trimmed string.
//
// Parameters:
//   - id: The GPU ID to query (e.g., "0" for the first GPU)
//   - query: The nvidia-smi query parameter (e.g., "temperature.gpu")
//
// Returns:
//
//	A string containing the query result with whitespace trimmed
func Query(id, query string) string {
	var out bytes.Buffer

	gpu_args := gpuArg + id
	query_args := queryArg + query

	cmd := exec.Command(bin, gpu_args, query_args, formatArg)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Printf("nvsmi exec error: %v\n", err)
	}

	return strings.TrimSpace(out.String())
}

// DeviceCount counts the output lines from the specified nvidia-smi GPU query.
//
// Deprecated: Use dcgm.GetAllDeviceCount for a GPU count that reports errors.
func DeviceCount(query string) uint {
	var out bytes.Buffer
	cmd := exec.Command(bin, queryArg+query, formatArg) //nolint:gosec // The query is one value for the fixed --query-gpu option.
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("nvsmi exec error: %v\n", err)
	}
	return uint(len(strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")))
}
