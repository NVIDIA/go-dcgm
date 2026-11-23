package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/NVIDIA/go-dcgm/samples/internal/processoutput"
)

// NOTE: The "WatchPidFields()" function must be initially called (as root) BEFORE starting the process to be monitored:
//  1. Run as root, for enabling health watches
//     sudo dcgmi stats -e
//  2. Start process to be monitored
//  3. Run processInfo. This is equivalent to "dcgmi stats --pid ENTERPID -v"
//     bazel run //samples/processInfo -- -pid PID
func main() {
	process := flag.Uint("pid", 0, "Provide pid to get this process information.")

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	// Request DCGM to start recording stats for GPU process fields
	group, err := dcgm.WatchPidFields()
	if err != nil {
		log.Panicln(err)
	}

	// Before retrieving process stats, wait few seconds for watches to be enabled and collect data
	log.Println("Enabling DCGM watches to start collecting process stats. This may take a few seconds....")
	time.Sleep(3000 * time.Millisecond)

	flag.Parse()

	pidInfo, err := dcgm.GetProcessInfo(group, *process)
	if err != nil {
		log.Panicln(err)
	}

	for i := range pidInfo {
		if err = processoutput.Write(os.Stdout, pidInfo[i]); err != nil {
			log.Panicln("Template error:", err)
		}
	}
}
