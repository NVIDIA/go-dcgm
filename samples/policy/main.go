package main

import (
	"context"
	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// dcgmi group -c "name" --default
// dcgmi policy -g GROUPID --set 0,0 -x -n -p -e -P 250 -T 100 -M 10
// dcgmi policy -g GROUPID --reg
func main() {
	ctx, done := context.WithCancel(context.Background())
	// Handle SIGINT (Ctrl+C) and SIGTERM (termination signal)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Received termination signal, exiting...")
		done()
	}()

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	// Choose policy conditions to register violation callback.
	// Note: Need to be root for some options
	// Available options are:
	// 1. dcgm.DbePolicy
	// 2. dcgm.PCIePolicy
	// 3. dcgm.MaxRtPgPolicy
	// 4. dcgm.ThermalPolicy
	// 5. dcgm.PowerPolicy
	// 6. dcgm.NvlinkPolicy
	// 7. dcgm.XidPolicy
	c, err := dcgm.ListenForPolicyViolations(ctx, dcgm.DbePolicy, dcgm.XidPolicy)
	if err != nil {
		log.Panicln(err)
	}

	for {
		select {
		case pe := <-c:
			log.Printf("PolicyViolation %6s %v\nTimestamp %2s %v\nData %7s %v",
				":", pe.Condition, ":", pe.Timestamp, ":", pe.Data)
		case <-ctx.Done():
			return
		}
	}
}
