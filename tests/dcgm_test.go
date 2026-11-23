package tests

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeviceCount(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)
	defer cleanup()

	count, err := dcgm.GetAllDeviceCount()
	check(t, err)
	if count == 0 {
		t.Fatal("DCGM found no GPUs")
	}

	c, err := physicalDeviceCount()
	check(t, err)

	if c != count {
		t.Errorf("Device Count from dcgm is wrong, got %d, want: %d", count, c)
	}
}

func BenchmarkDeviceCount1(b *testing.B) {
	_, _ = dcgm.Init(dcgm.Embedded)

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		_, _ = dcgm.GetAllDeviceCount()
	}

	b.StopTimer()

	_ = dcgm.Shutdown()
}

// TODO: We need a way to determine if we have an NVIDIA CPU
func TestCpuQuery(t *testing.T) {
	t.Setenv("DCGM_SKIP_SYSMON_HARDWARE_CHECK", "1")

	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)

	defer cleanup()

	hierarchy, err := dcgm.GetCPUHierarchy()
	if err != nil {
		if strings.Contains(err.Error(), "not currently loaded") {
			t.Skip("CPU hierarchy module not loaded, skipping")
		}
		t.Fatalf("Failed to get CPU hierarchy: %v", err)
	}

	if hierarchy.NumCPUs == 0 {
		t.Skip("No CPUs found in hierarchy")
	}

	for i := uint(0); i < hierarchy.NumCPUs; i++ {
		coresFound := false

		for j := uint(0); j < dcgm.MAX_CPU_CORE_BITMASK_COUNT; j++ {
			if hierarchy.CPUs[i].OwnedCores[j] != 0 {
				coresFound = true
			}
		}

		if coresFound == false {
			t.Errorf("Cpu %d has no cores", i)
		}
	}
}

func TestDeviceInfo(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	check(t, err)
	if len(gpus) == 0 {
		t.Fatal("DCGM found no supported GPUs")
	}

	for _, gpu := range gpus {
		info, err := dcgm.GetDeviceInfo(gpu)
		check(t, err)

		// DCGM GPU IDs and nvidia-smi indexes can refer to different devices.
		busID := info.PCI.BusID
		if busID == "" {
			t.Fatalf("DCGM GPU %d has no PCI bus ID", gpu)
		}

		fields := []struct {
			query string
			value string
		}{
			{"driver_version", info.Identifiers.DriverVersion},
			{"name", info.Identifiers.Model},
			{"serial", info.Identifiers.Serial},
			{"uuid", info.UUID},
			{"pci.bus_id", info.PCI.BusID},
			{"vbios_version", info.Identifiers.Vbios},
			{"inforom.img", info.Identifiers.InforomImageVersion},
			{"power.default_limit", strconv.FormatUint(uint64(info.Power), 10)},
		}
		for _, field := range fields {
			checkDeviceInfoField(t, busID, field.query, field.value)
		}
	}
}

func checkDeviceInfoField(t *testing.T, busID, query, value string) {
	t.Helper()

	out, err := exec.Command(bin, gpuArg+busID, queryArg+query, formatArg).Output() //nolint:gosec // Passes the PCI bus ID and field as values to fixed options.
	if err != nil {
		t.Fatalf("query GPU at %s for %s with nvidia-smi: %v", busID, query, err)
	}

	reference := strings.TrimSpace(string(out))
	if err := compareDeviceInfoField(query, value, reference); err != nil {
		t.Errorf("GPU at %s, %s: %v", busID, query, err)
	}
}

func compareDeviceInfoField(query, value, reference string) error {
	if query == "uuid" {
		if !gpuUUIDPattern.MatchString(reference) {
			return fmt.Errorf("nvidia-smi returned an invalid UUID: %q", reference)
		}
	} else if reference == "[N/A]" || strings.Contains(value, "NOT_SUPPORTED") {
		return nil
	}

	want := reference
	if query == "power.default_limit" {
		power, err := strconv.ParseFloat(reference, 64)
		if err != nil {
			return fmt.Errorf("parse default power from nvidia-smi: %w", err)
		}

		// DCGM exposes the boot default as whole watts, discarding the fraction.
		want = strconv.FormatUint(uint64(power), 10)
	}

	if value != want {
		return fmt.Errorf("DCGM returned %q, want %q (nvidia-smi=%q)", value, want, reference)
	}
	return nil
}

func TestDeviceInfoFieldComparison(t *testing.T) {
	const uuid = "GPU-12345678-1234-1234-1234-123456789abc"
	for _, tc := range []struct {
		name, query, value, reference string
		wantError                     bool
	}{
		{name: "matching UUID", query: "uuid", value: uuid, reference: uuid},
		{name: "unavailable UUID", query: "uuid", value: "[N/A]", reference: "[N/A]", wantError: true},
		{name: "invalid UUID", query: "uuid", value: "0", reference: "0", wantError: true},
		{name: "wrong UUID", query: "uuid", value: "GPU-87654321-4321-4321-4321-cba987654321", reference: uuid, wantError: true},
		{name: "optional field unavailable", query: "serial", value: "unknown", reference: "[N/A]"},
		{name: "default power rounded down", query: "power.default_limit", value: "125", reference: "125.75"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := compareDeviceInfoField(tc.query, tc.value, tc.reference)
			if (err != nil) != tc.wantError {
				t.Fatalf("compareDeviceInfoField() = %v, want error=%v", err, tc.wantError)
			}
		})
	}
}

func BenchmarkDeviceInfo1(b *testing.B) {
	_, _ = dcgm.Init(dcgm.Embedded)

	b.StartTimer()

	for n := 0; n < b.N; n++ {
		// assuming there will be atleast 1 GPU attached
		_, _ = dcgm.GetDeviceInfo(uint(0))
	}

	b.StopTimer()

	_ = dcgm.Shutdown()
}

func assertInRange(t *testing.T, name string, val, minimum, maximum float64) {
	t.Helper()
	if val < minimum || val > maximum {
		t.Errorf("%s out of range: got %v, want [%v, %v]", name, val, minimum, maximum)
	}
}

func TestDeviceStatus(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	check(t, err)
	if len(gpus) == 0 {
		t.Fatal("DCGM found no supported GPUs")
	}

	for _, gpu := range gpus {
		status, err := dcgm.GetDeviceStatus(gpu)
		check(t, err)

		t.Logf("GPU %d: Power=%.1fW Temp=%d°C GPU_Util=%d%% Mem_Util=%d%% SM=%dMHz Mem=%dMHz",
			gpu, status.Power, status.Temperature,
			status.Utilization.GPU, status.Utilization.Memory,
			status.Clocks.Cores, status.Clocks.Memory)

		assertInRange(t, "Power (W)", status.Power, 1, 1000)
		assertInRange(t, "Temperature (C)", float64(status.Temperature), 0, 110)
		assertInRange(t, "GPU Utilization (%)", float64(status.Utilization.GPU), 0, 100)
		assertInRange(t, "Memory Utilization (%)", float64(status.Utilization.Memory), 0, 100)
		assertInRange(t, "Encoder Utilization (%)", float64(status.Utilization.Encoder), 0, 100)
		assertInRange(t, "SM Clock (MHz)", float64(status.Clocks.Cores), 0, 5000)
		assertInRange(t, "Memory Clock (MHz)", float64(status.Clocks.Memory), 0, 15000)
	}
}

var gpuUUIDPattern = regexp.MustCompile(`^GPU-[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`)

// physicalDeviceCount counts unique physical GPU UUIDs reported by nvidia-smi.
func physicalDeviceCount() (uint, error) {
	out, err := exec.Command(bin, queryArg+"uuid", formatArg).Output()
	if err != nil {
		return 0, fmt.Errorf("query GPU identities with nvidia-smi: %w", err)
	}

	return countGPUUUIDs(string(out))
}

// countGPUUUIDs rejects empty, malformed, or duplicate physical GPU identities.
func countGPUUUIDs(output string) (uint, error) {
	output = strings.TrimSpace(output)
	if output == "" {
		return 0, errors.New("nvidia-smi returned no GPU identities")
	}

	seen := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		uuid := strings.TrimSpace(line)
		if !gpuUUIDPattern.MatchString(uuid) || seen[strings.ToLower(uuid)] {
			return 0, fmt.Errorf("nvidia-smi returned an invalid or duplicate GPU identity: %q", uuid)
		}
		seen[strings.ToLower(uuid)] = true
	}
	return uint(len(seen)), nil
}

// TestNvidiaSMIDeviceCount checks the identity contract without accessing a GPU.
func TestNvidiaSMIDeviceCount(t *testing.T) {
	const first = "GPU-12345678-1234-1234-1234-123456789abc"
	const second = "GPU-87654321-4321-4321-4321-cba987654321"
	for _, tc := range []struct {
		name      string
		output    string
		want      uint
		wantError bool
	}{
		{name: "empty", wantError: true},
		{name: "malformed", output: "No devices were found", wantError: true},
		{name: "duplicate", output: first + "\n" + first, wantError: true},
		{name: "blank row", output: first + "\n\n" + second, wantError: true},
		{name: "single", output: first + "\n", want: 1},
		{name: "multiple", output: first + "\n" + second + "\n", want: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := countGPUUUIDs(tc.output)
			if got != tc.want || (err != nil) != tc.wantError {
				t.Fatalf("countGPUUUIDs() = %d, %v; want %d, error=%v", got, err, tc.want, tc.wantError)
			}
		})
	}
	t.Run("executable unavailable", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		if count, err := physicalDeviceCount(); count != 0 || err == nil {
			t.Fatalf("physicalDeviceCount() = %d, %v; want zero and an execution error", count, err)
		}
		if count := DeviceCount("uuid"); count != 1 {
			t.Fatalf("legacy DeviceCount() = %d; want one for empty output", count)
		}
	})
}
