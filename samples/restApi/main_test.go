package main

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

// TestDefaultListenAddressUsesLoopback verifies the sample's safe default bind address.
func TestDefaultListenAddressUsesLoopback(t *testing.T) {
	host, port, err := net.SplitHostPort(defaultListenAddress)
	if err != nil {
		t.Fatalf("defaultListenAddress %q is invalid: %v", defaultListenAddress, err)
	}

	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		t.Fatalf("defaultListenAddress host %q is not a loopback address", host)
	}
	if port != "8070" {
		t.Fatalf("defaultListenAddress port = %q, want 8070", port)
	}
}

// TestHTTPServerHandlerRegistersRoutes verifies every documented REST endpoint is reachable by GET.
func TestHTTPServerHandlerRegistersRoutes(t *testing.T) {
	server := &httpServer{router: mux.NewRouter()}
	server.handler()

	routes := []string{
		"/dcgm/device/info/id/0",
		"/dcgm/device/info/id/0/json",
		"/dcgm/device/info/uuid/GPU-0",
		"/dcgm/device/info/uuid/GPU-0/json",
		"/dcgm/device/status/id/0",
		"/dcgm/device/status/id/0/json",
		"/dcgm/device/status/uuid/GPU-0",
		"/dcgm/device/status/uuid/GPU-0/json",
		"/dcgm/process/info/pid/1",
		"/dcgm/process/info/pid/1/json",
		"/dcgm/health/id/0",
		"/dcgm/health/id/0/json",
		"/dcgm/health/uuid/GPU-0",
		"/dcgm/health/uuid/GPU-0/json",
		"/dcgm/status",
		"/dcgm/status/json",
	}

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, route, http.NoBody)
			match := &mux.RouteMatch{}
			if !server.router.Match(request, match) {
				t.Fatalf("GET %s did not match a registered route: %v", route, match.MatchErr)
			}
			if match.Handler == nil {
				t.Fatalf("GET %s matched without a handler", route)
			}
		})
	}

	request := httptest.NewRequest(http.MethodPost, "/dcgm/status", http.NoBody)
	match := &mux.RouteMatch{}
	if server.router.Match(request, match) {
		t.Fatal("POST /dcgm/status unexpectedly matched a GET-only route")
	}
}

// TestHTTPServerServeReportsListenErrors checks that an invalid listen address is logged.
func TestHTTPServerServeReportsListenErrors(t *testing.T) {
	var output bytes.Buffer
	previousOutput := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() {
		log.SetOutput(previousOutput)
	})

	server := &httpServer{
		server: &http.Server{
			Addr:              "127.0.0.1:65536",
			ReadHeaderTimeout: timeout,
		},
	}
	server.serve()

	if !strings.Contains(output.String(), "Error:") {
		t.Fatalf("serve log = %q, want a listen error", output.String())
	}
}

// TestHTTPServerStopInvokesShutdown checks that stop invokes the HTTP server's shutdown callback.
func TestHTTPServerStopInvokesShutdown(t *testing.T) {
	nativeServer := &http.Server{ReadHeaderTimeout: timeout}
	shutdown := make(chan struct{})
	nativeServer.RegisterOnShutdown(func() {
		close(shutdown)
	})

	server := &httpServer{server: nativeServer}
	server.stop()

	select {
	case <-shutdown:
	case <-time.After(timeout):
		t.Fatal("shutdown callback was not invoked")
	}
}
