package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

func TestPrinterDiscardsPartialRenderOnError(t *testing.T) {
	tmpl := template.Must(template.New("failure").Parse("prefix {{.Known}} suffix {{.Missing}}"))
	req := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
	resp := httptest.NewRecorder()

	printer(resp, req, struct{ Known string }{Known: "rendered"}, tmpl)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", resp.Code)
	}
	body := resp.Body.String()
	if !strings.Contains(body, "can't evaluate field Missing") {
		t.Fatalf("response does not explain the rendering error: %q", body)
	}
	if strings.Contains(body, "prefix") || strings.Contains(body, "rendered") {
		t.Fatalf("response contains partially rendered output: %q", body)
	}
}

func TestPrinterPreservesSuccessfulOutput(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/status", http.NoBody)
	resp := httptest.NewRecorder()

	printer(resp, req, dcgm.Status{Memory: 1024, CPU: 1.5}, hostengineTemplate)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	const want = "Memory(KB)      : 1024\nCPU(%)          : 1.50\n"
	if got := resp.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestProcessPrintRendersAllEntries(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/process/42", http.NoBody)
	resp := httptest.NewRecorder()

	processPrint(resp, req, []dcgm.ProcessInfo{{PID: 42, GPU: 0}, {PID: 42, GPU: 1}})

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	body := resp.Body.String()
	if strings.Count(body, "PID                          : 42") != 2 {
		t.Fatalf("body does not contain both process entries: %q", body)
	}
	if first, second := strings.Index(body, "GPU ID                       : 0"), strings.Index(body, "GPU ID                       : 1"); first < 0 || second <= first {
		t.Fatalf("process entries are missing or out of order: %q", body)
	}
}
