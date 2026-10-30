package tests

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// TestHealthCheck demonstrates GPU health checking functionality
// This is equivalent to the health sample but runs for a limited time
func TestHealthCheck(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}

	if len(gpus) == 0 {
		t.Skip("No supported GPUs found for health checking")
	}

	// Monitor health for a few seconds instead of indefinitely
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	timeout := time.After(6 * time.Second)
	checkCount := 0

	for {
		select {
		case <-ticker.C:
			for _, gpu := range gpus {
				h, err := dcgm.HealthCheckByGpuId(gpu)
				if err != nil {
					t.Errorf("Failed to get health status for GPU %d: %v", gpu, err)
					continue
				}

				t.Logf("GPU %d Health Check:", gpu)
				t.Logf("  Status: %s", h.Status)

				for _, watch := range h.Watches {
					t.Logf("  Watch Type: %s", watch.Type)
					t.Logf("  Watch Status: %s", watch.Status)
					if watch.Error != "" {
						t.Logf("  Watch Error: %s", watch.Error)
					}
				}

				// Basic validation
				if h.Status == "" {
					t.Errorf("GPU %d has empty health status", gpu)
				}
			}
			checkCount++

		case <-timeout:
			t.Logf("Health monitoring completed after %d checks", checkCount)
			return
		}
	}
}

// TestHealthCheckSingle demonstrates a single health check
func TestHealthCheckSingle(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}

	if len(gpus) == 0 {
		t.Skip("No supported GPUs found")
	}

	// Test first GPU
	gpu := gpus[0]
	h, err := dcgm.HealthCheckByGpuId(gpu)
	if err != nil {
		t.Fatalf("Failed to get health status for GPU %d: %v", gpu, err)
	}

	t.Logf("GPU %d Health Status: %s", gpu, h.Status)

	if len(h.Watches) == 0 {
		t.Logf("No health watches configured for GPU %d", gpu)
	} else {
		t.Logf("Health watches for GPU %d:", gpu)
		for i, watch := range h.Watches {
			t.Logf("  Watch %d:", i+1)
			t.Logf("    Type: %s", watch.Type)
			t.Logf("    Status: %s", watch.Status)
			if watch.Error != "" {
				t.Logf("    Error: %s", watch.Error)
			}
		}
	}

	// Basic assertions
	if h.Status == "" {
		t.Error("Health status is empty")
	}
}

// TestHealthCheckAllGPUs demonstrates health checking for all GPUs
func TestHealthCheckAllGPUs(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}

	if len(gpus) == 0 {
		t.Skip("No supported GPUs found")
	}

	healthyGPUs := 0
	unhealthyGPUs := 0

	for _, gpu := range gpus {
		h, err := dcgm.HealthCheckByGpuId(gpu)
		if err != nil {
			t.Errorf("Failed to get health status for GPU %d: %v", gpu, err)
			continue
		}

		t.Logf("GPU %d: %s", gpu, h.Status)

		// Count healthy vs unhealthy
		if h.Status == "Healthy" || h.Status == "OK" {
			healthyGPUs++
		} else {
			unhealthyGPUs++
			t.Logf("GPU %d is not healthy: %s", gpu, h.Status)

			// Log any watch errors
			for _, watch := range h.Watches {
				if watch.Error != "" {
					t.Logf("  Watch %s error: %s", watch.Type, watch.Error)
				}
			}
		}
	}

	t.Logf("Health summary: %d healthy, %d unhealthy GPUs", healthyGPUs, unhealthyGPUs)

	// We expect at least some GPUs to be available
	if healthyGPUs == 0 && unhealthyGPUs == 0 {
		t.Error("No GPU health status could be determined")
	}
}

var dcgmVersionPattern = regexp.MustCompile(`version:\s*(\d+)\.(\d+)\.(\d+)`)

// requireHealthV3API stops the test when the DCGM 4.7.1 health APIs are unavailable.
// DCGM 4.7.1 and later must provide them, so the test fails there; older DCGM skips.
func requireHealthV3API(t *testing.T, err error) {
	t.Helper()

	if !dcgm.IsFeatureUnavailable(err) {
		return
	}

	version := hostengineVersion(t, err)
	if version[0] > 4 || (version[0] == 4 && (version[1] > 7 || (version[1] == 7 && version[2] >= 1))) {
		t.Fatalf("DCGM %d.%d.%d must provide the health API: %v", version[0], version[1], version[2], err)
	}

	t.Skipf("DCGM %d.%d.%d does not provide the health API: %v", version[0], version[1], version[2], err)
}

// hostengineVersion returns the connected host engine's major, minor and patch
// version from its build info. It fails the test if the version can't be read;
// apiErr is the unavailable-API error, included in that failure message.
func hostengineVersion(t *testing.T, apiErr error) [3]int {
	t.Helper()

	info, err := dcgm.GetHostengineVersionInfo()
	if err != nil {
		t.Fatalf("health API unavailable (%v) and the host engine version is unknown: %v", apiErr, err)
	}

	match := dcgmVersionPattern.FindStringSubmatch(info.RawBuildInfoString)
	if match == nil {
		t.Fatalf("health API unavailable (%v) and no version in build info %q", apiErr, info.RawBuildInfoString)
	}

	var version [3]int
	for i := range version {
		version[i], _ = strconv.Atoi(match[i+1])
	}

	return version
}

// newHealthTestGroup returns a group holding the first supported GPU.
func newHealthTestGroup(t *testing.T) dcgm.GroupHandle {
	t.Helper()

	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		t.Fatalf("Failed to get supported devices: %v", err)
	}
	if len(gpus) == 0 {
		t.Skip("No supported GPUs found for health checking")
	}

	group, err := dcgm.CreateGroup(fmt.Sprintf("health-v3-%d", time.Now().UnixNano()))
	if err != nil {
		t.Fatalf("Failed to create group: %v", err)
	}
	t.Cleanup(func() {
		if err := dcgm.DestroyGroup(group); err != nil {
			t.Errorf("Failed to destroy group: %v", err)
		}
	})

	if err := dcgm.AddToGroup(group, gpus[0]); err != nil {
		t.Fatalf("Failed to add GPU %d to group: %v", gpus[0], err)
	}

	return group
}

// TestHealthSetWithParamsStrict checks that the host engine reports some supported
// health systems and that DCGM_HEALTH_WATCH_ALL enables exactly those.
func TestHealthSetWithParamsStrict(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	t.Cleanup(cleanup) // runs after the group cleanup registered later

	supported, err := dcgm.HealthGetSupportedSystems()
	requireHealthV3API(t, err)
	if err != nil {
		t.Fatalf("HealthGetSupportedSystems failed: %v", err)
	}
	if supported == 0 {
		t.Fatal("HealthGetSupportedSystems returned an empty mask")
	}

	group := newHealthTestGroup(t)

	err = dcgm.HealthSetWithParamsStrict(group, dcgm.DCGM_HEALTH_WATCH_ALL, time.Second, time.Minute)
	requireHealthV3API(t, err)
	if err != nil {
		t.Fatalf("HealthSetWithParamsStrict(ALL) failed: %v", err)
	}

	enabled, err := dcgm.HealthGet(group)
	if err != nil {
		t.Fatalf("HealthGet failed: %v", err)
	}
	if enabled != supported {
		t.Fatalf("enabled systems %#x, want supported systems %#x", uint(enabled), uint(supported))
	}
}

// TestHealthSetWithParamsStrictRejectsUnsupported checks that a system the host engine
// does not support is rejected and the group's watches stay unchanged.
func TestHealthSetWithParamsStrictRejectsUnsupported(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		t.Fatalf("Failed to initialize DCGM: %v", err)
	}
	t.Cleanup(cleanup) // runs after the group cleanup registered later

	supported, err := dcgm.HealthGetSupportedSystems()
	requireHealthV3API(t, err)
	if err != nil {
		t.Fatalf("HealthGetSupportedSystems failed: %v", err)
	}

	known, unsupported := pickSupportedAndUnsupported(supported)
	if known == 0 {
		t.Fatalf("no defined health system is supported (mask %#x)", uint(supported))
	}
	if unsupported == 0 {
		t.Skipf("every defined health system is supported (mask %#x)", uint(supported))
	}

	// Start from a known state: one supported system enabled.
	group := newHealthTestGroup(t)
	if err = dcgm.HealthSetWithParamsStrict(group, known, time.Second, time.Minute); err != nil {
		t.Fatalf("HealthSetWithParamsStrict(%#x) failed: %v", uint(known), err)
	}

	// Adding an unsupported system must fail with DCGM_ST_NOT_SUPPORTED.
	err = dcgm.HealthSetWithParamsStrict(group, known|unsupported, time.Second, time.Minute)

	var dcgmErr *dcgm.Error
	if !errors.As(err, &dcgmErr) || dcgmErr.Code != dcgm.DCGM_ST_NOT_SUPPORTED {
		t.Fatalf("HealthSetWithParamsStrict(%#x) = %v, want *dcgm.Error with DCGM_ST_NOT_SUPPORTED", uint(known|unsupported), err)
	}
	if dcgm.IsFeatureUnavailable(err) {
		t.Fatalf("a rejected system must not be reported as an unavailable feature: %v", err)
	}

	enabled, err := dcgm.HealthGet(group)
	if err != nil {
		t.Fatalf("HealthGet failed: %v", err)
	}
	if enabled != known {
		t.Fatalf("enabled systems %#x after a rejected call, want %#x", uint(enabled), uint(known))
	}
}

// pickSupportedAndUnsupported returns the first defined health system that is in
// the supported mask and the first one that is not. Either is 0 when none exists.
func pickSupportedAndUnsupported(supported dcgm.HealthSystem) (known, unsupported dcgm.HealthSystem) {
	defined := []dcgm.HealthSystem{
		dcgm.DCGM_HEALTH_WATCH_PCIE, dcgm.DCGM_HEALTH_WATCH_NVLINK, dcgm.DCGM_HEALTH_WATCH_PMU,
		dcgm.DCGM_HEALTH_WATCH_MCU, dcgm.DCGM_HEALTH_WATCH_MEM, dcgm.DCGM_HEALTH_WATCH_SM,
		dcgm.DCGM_HEALTH_WATCH_INFOROM, dcgm.DCGM_HEALTH_WATCH_THERMAL, dcgm.DCGM_HEALTH_WATCH_POWER,
		dcgm.DCGM_HEALTH_WATCH_DRIVER, dcgm.DCGM_HEALTH_WATCH_NVSWITCH_NONFATAL,
		dcgm.DCGM_HEALTH_WATCH_NVSWITCH_FATAL, dcgm.DCGM_HEALTH_WATCH_CONNECTX, dcgm.DCGM_HEALTH_WATCH_IMEX,
	}

	for _, bit := range defined {
		switch {
		case supported&bit != 0 && known == 0:
			known = bit
		case supported&bit == 0 && unsupported == 0:
			unsupported = bit
		}
	}

	return known, unsupported
}
