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

// DeviceCount returns the number of NVIDIA GPU devices available in the system
// by executing nvidia-smi with the specified query parameter.
//
// Parameters:
//   - query: The nvidia-smi query parameter to execute
//
// Returns:
//
//	The number of GPU devices as an unsigned integer
func DeviceCount(query string) uint {
	var out bytes.Buffer

	query_arg := queryArg + query
	cmd := exec.Command(bin, query_arg, formatArg)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Printf("nvsmi exec error: %v\n", err)
	}

	nvSmi := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")

	return uint(len(nvSmi))
}
