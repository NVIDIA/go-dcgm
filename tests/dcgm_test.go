package tests

import (
	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestDeviceCount(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(err, t)
	defer cleanup()

	count, err := dcgm.GetAllDeviceCount()
	check(err, t)

	query := "count"
	c := DeviceCount(query)

	if c != count {
		t.Errorf("Device Count from dcgm is wrong, got %d, want: %d", count, c)
	}
}

func BenchmarkDeviceCount1(b *testing.B) {
	dcgm.Init(dcgm.Embedded)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		dcgm.GetAllDeviceCount()
	}
	b.StopTimer()

	dcgm.Shutdown()
}

func TestCpuQuery(t *testing.T) {
	os.Setenv("DCGM_SKIP_SYSMON_HARDWARE_CHECK", "1")
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(err, t)
	defer cleanup()

	hierarchy, err := dcgm.GetCpuHierarchy()
	check(err, t)

	if hierarchy.NumCpus == 0 {
		t.Errorf("Found no CPUs")
	}

	for i := uint(0); i < hierarchy.NumCpus; i++ {
		var coresFound = false
		for j := uint(0); j < dcgm.MAX_CPU_CORE_BITMASK_COUNT; j++ {
			if hierarchy.Cpus[i].OwnedCores[j] != 0 {
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
	check(err, t)
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
	check(err, t)

	for _, gpu := range gpus {
		info, err := dcgm.GetDeviceInfo(gpu)
		check(err, t)

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
				check(err, t)
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
	dcgm.Init(dcgm.Embedded)

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		// assuming there will be atleast 1 GPU attached
		dcgm.GetDeviceInfo(uint(0))
	}
	b.StopTimer()

	dcgm.Shutdown()
}

func TestDeviceStatus(t *testing.T) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	check(err, t)
	defer cleanup()

	gpus, err := dcgm.GetSupportedDevices()
	check(err, t)

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
		check(err, t)

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
				output = strconv.FormatUint(uint64(math.Round(status.Power)), 10)
				power, err := strconv.ParseFloat(res, 64)
				check(err, t)
				res = strconv.FormatUint(uint64(math.Round(power)), 10)
			case "temperature.gpu":
				msg = "Device temperature"
				output = strconv.FormatUint(uint64(status.Temperature), 10)
			case "utilization.gpu":
				msg = "Device gpu utilization"
				output = strconv.FormatUint(uint64(status.Utilization.GPU), 10)
			case "utilization.memory":
				msg = "Device memory utilization"
				output = strconv.FormatUint(uint64(status.Utilization.Memory), 10)
			case "encoder.stats.averageFps":
				msg = "Device encoder utilization"
				output = strconv.FormatUint(uint64(status.Utilization.Encoder), 10)
			case "clocks.current.sm":
				msg = "Device sm clock"
				output = strconv.FormatUint(uint64(status.Clocks.Cores), 10)
			case "clocks.current.memory":
				msg = "Device mem clock"
				output = strconv.FormatUint(uint64(status.Clocks.Memory), 10)
			}

			if strings.Compare(res, output) != 0 {
				t.Errorf("%v from dcgm is wrong, got: %v, want: %v", msg, output, res)
			}
		}
	}
}
