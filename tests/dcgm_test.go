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

func TestCpuQuery(t *testing.T) {
	t.Setenv("DCGM_SKIP_SYSMON_HARDWARE_CHECK", "1")

	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)

	defer cleanup()

	hierarchy, err := dcgm.GetCPUHierarchy()
	check(t, err)

	if hierarchy.NumCPUs == 0 {
		t.Errorf("Found no CPUs")
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

func TestDeviceStatus(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(t, err)
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	check(t, err)

	fields := []string{
		"power.draw",
		"temperature.gpu",
		"utilization.gpu",
		"utilization.memory",
		"encoder.stats.averageFps",
		"clocks.current.sm",
		"clocks.current.memory",
	}

	for _, gpu := range gpus {
		status, err := dcgm.GetDeviceStatus(gpu)
		check(t, err)

		id := strconv.FormatUint(uint64(gpu), 10)

		for _, val := range fields {
			var msg, output string

			res := Query(id, val)
			if res == "[N/A]" {
				continue
			}

			switch val {
			case "power.draw":
				msg = "Device power utilization"
				output = strconv.FormatFloat(math.Round(status.Power), 'f', -1, 64)
				power, err := strconv.ParseFloat(res, 64)
				check(t, err)

				res = strconv.FormatFloat(math.Round(power), 'f', -1, 64)
			case "temperature.gpu":
				msg = "Device temperature"
				output = strconv.FormatInt(status.Temperature, 10)
			case "utilization.gpu":
				msg = "Device gpu utilization"
				output = strconv.FormatInt(status.Utilization.GPU, 10)
			case "utilization.memory":
				msg = "Device memory utilization"
				output = strconv.FormatInt(status.Utilization.Memory, 10)
			case "encoder.stats.averageFps":
				msg = "Device encoder utilization"
				output = strconv.FormatInt(status.Utilization.Encoder, 10)
			case "clocks.current.sm":
				msg = "Device sm clock"
				output = strconv.FormatInt(status.Clocks.Cores, 10)
			case "clocks.current.memory":
				msg = "Device mem clock"
				output = strconv.FormatInt(status.Clocks.Memory, 10)
			}

			if strings.Compare(res, output) != 0 {
				t.Errorf("%v from dcgm is wrong, got: %v, want: %v", msg, output, res)
			}
		}
	}
}
