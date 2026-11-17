package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		return err
	}
	defer cleanup()
	return readConfigurations(dcgm.ConfigGet)
}

func readConfigurations(get func(dcgm.GroupHandle, dcgm.ConfigType) ([]dcgm.GPUConfig, error)) error {
	var readErrors error
	for _, kind := range []dcgm.ConfigType{dcgm.ConfigCurrentState, dcgm.ConfigTargetState} {
		configs, readErr := get(dcgm.GroupAllGPUs(), kind)
		for i := range configs {
			printConfiguration(kind, &configs[i])
		}
		readErrors = errors.Join(readErrors, readErr)
	}
	return readErrors
}

func printConfiguration(kind dcgm.ConfigType, config *dcgm.GPUConfig) {
	state := "current"
	if kind == dcgm.ConfigTargetState {
		state = "target"
	}
	fmt.Printf("GPU %d %s configuration: ", config.GPU, state)
	if config.PowerLimit == nil {
		fmt.Print("power not set or unavailable")
	} else {
		fmt.Printf("power %d W", config.PowerLimit.Watts)
	}
	if config.MemoryClockHz != nil {
		fmt.Printf(", memory clock %d Hz", *config.MemoryClockHz)
	}
	if config.SMClockHz != nil {
		fmt.Printf(", SM clock %d Hz", *config.SMClockHz)
	}
	fmt.Println()
}
