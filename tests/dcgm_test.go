package tests

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestDeviceCount(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)
	defer cleanup()

	count, err := dcgm.GetAllDeviceCount()
	check(t, err)

	query := "count"
	c := DeviceCount(query)

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

	fields := []string{
		"driver_version",
		"name",
		"serial",
		"uuid",
		"pci.bus_id",
		"vbios_version",
		"inforom.img",
		"power.limit",
	}

	gpus, err := dcgm.GetSupportedDevices()
	check(t, err)

	for _, gpu := range gpus {
		info, err := dcgm.GetDeviceInfo(gpu)
		check(t, err)

		id := strconv.FormatUint(uint64(gpu), 10)

		for _, val := range fields {
			var msg, output string

			res := Query(id, val)
			if res == "[N/A]" {
				continue
			}

			switch val {
			case "driver_version":
				msg = "Driver version"
				output = info.Identifiers.DriverVersion
			case "name":
				msg = "Device name"
				output = info.Identifiers.Model
			case "serial":
				msg = "Device Serial number"
				output = info.Identifiers.Serial
			case "uuid":
				msg = "Device UUID"
				output = info.UUID
			case "pci.bus_id":
				msg = "Device PCI busId"
				output = info.PCI.BusID
			case "vbios_version":
				msg = "Device vbios version"
				output = info.Identifiers.Vbios
			case "inforom.img":
				msg = "Device inforom image"
				output = info.Identifiers.InforomImageVersion
			case "power.limit":
				msg = "Device power limit"
				output = strconv.FormatUint(uint64(info.Power), 10)
				power, err := strconv.ParseFloat(res, 64)
				check(t, err)

				res = strconv.FormatUint(uint64(math.Round(power)), 10)
			}

			if strings.Compare(res, output) != 0 {
				if strings.Contains(output, "NOT_SUPPORTED") {
					continue
				}

				t.Errorf("%v from dcgm is wrong, got: %v, want: %v", msg, output, res)
			}
		}
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

func assertInRange(t *testing.T, name string, val, min, max float64) {
	t.Helper()
	if val < min || val > max {
		t.Errorf("%s out of range: got %v, want [%v, %v]", name, val, min, max)
	}
}

func TestDeviceStatus(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	check(t, err)

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
