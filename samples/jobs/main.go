package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

// main reports a failed sample after the owned job, group and DCGM session are cleaned up.
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run owns the embedded session and a group containing the first supported GPU.
func run() (err error) {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		return err
	}
	defer cleanup()

	// Keep the sample group private so cleanup cannot affect another client.
	group, err := dcgm.CreateGroup(fmt.Sprintf("jobs-%d", time.Now().UnixNano()))
	if err != nil {
		return err
	}
	defer func() {
		if e := dcgm.DestroyGroup(group); e != nil {
			err = errors.Join(err, fmt.Errorf("destroy job group: %w", e))
		}
	}()

	// Job watches cover the selected GPU; they must start before the job interval.
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		return err
	}

	if len(gpus) == 0 {
		return errors.New("no supported GPUs")
	}

	if addErr := dcgm.AddToGroup(group, gpus[0]); addErr != nil {
		return addErr
	}
	if watchErr := dcgm.WatchJobFields(group, 100*time.Millisecond, time.Minute, 0); watchErr != nil {
		return watchErr
	}

	return recordSampleJob(group)
}

// recordSampleJob records and prints one interval, then removes its owned job record.
func recordSampleJob(group dcgm.GroupHandle) (err error) {
	id := fmt.Sprintf("sample-job-%d", time.Now().UnixNano())
	if startErr := dcgm.JobStart(group, id); startErr != nil {
		return startErr
	}
	defer func() {
		if e := dcgm.JobRemove(id); e != nil {
			err = errors.Join(err, fmt.Errorf("remove sample job: %w", e))
		}
	}()

	// Run your workload on the selected GPU during this interval. The sample
	// itself only records activity; it does not generate GPU work.
	time.Sleep(3 * time.Second)
	if stopErr := dcgm.JobStop(id); stopErr != nil {
		return stopErr
	}

	stats, err := dcgm.JobGetStats(id)
	if err != nil {
		return err
	}

	printJobGPUs(stats.GPUs)
	return nil
}

// printJobGPUs keeps unavailable values distinct from real zero utilization or energy.
func printJobGPUs(gpus []dcgm.GPUJobStats) {
	for i := range gpus {
		gpu := &gpus[i]
		smAverage, energy := "N/A", "N/A"
		if !dcgm.IsInt32Blank(int(gpu.SMUtilization.Average)) {
			smAverage = fmt.Sprintf("%d%%", gpu.SMUtilization.Average)
		}

		if !dcgm.IsInt64Blank(gpu.EnergyConsumed) {
			energy = fmt.Sprintf("%d mWs", gpu.EnergyConsumed)
		}

		fmt.Printf("GPU %d: SM utilization avg %s, energy %s, interval %d-%d us\n",
			gpu.GPU, smAverage, energy, gpu.StartTime, gpu.EndTime)
	}
}
