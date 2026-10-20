package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

const defaultListenAddress = "127.0.0.1:8070"

// res: curl localhost:8070/dcgm/device/info/id/0

func main() {
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)

	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	addr := defaultListenAddress
	server := newHttpServer(addr)

	go func() {
		log.Printf("Running http server on %s", addr)
		server.serve()
	}()

	defer server.stop()

	<-stopSig
}
