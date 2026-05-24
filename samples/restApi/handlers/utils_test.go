package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
	"github.com/stretchr/testify/assert"
)

// TestPrinter exercises every category of input across the package's
// pre-parsed *template.Template values. Cases are tagged by category so
// failures surface what kind of behavior regressed.
func TestPrinter(t *testing.T) {
	cases := []struct {
		name       string
		category   string // positive | negative | boundary | corner | attacker
		tmpl       *template.Template
		stats      any
		wantCode   int
		wantSubstr string // must appear in body
		notSubstr  string // must NOT appear in body (proves no template re-evaluation)
	}{
		// positive — happy path; known-good stats render expected substring.
		{
			"positive_hostengine_renders_memory_and_cpu", "positive", hostengineTmpl,
			dcgm.Status{Memory: 1024, CPU: 1.5},
			http.StatusOK, "Memory(KB)      : 1024", "",
		},
		{
			"positive_hostengine_renders_two_decimal_cpu", "positive", hostengineTmpl,
			dcgm.Status{Memory: 1, CPU: 99.875},
			http.StatusOK, "CPU(%)          : 99.88", "",
		},
		{
			"positive_healthStatus_renders_overall", "positive", healthStatusTmpl,
			dcgm.DeviceHealth{GPU: 0, Status: "OK"},
			http.StatusOK, "Status             : OK", "",
		},

		// negative — template can't render the stats. Note: Execute() writes
		// the static prefix before hitting the failing action, which
		// implicitly commits HTTP 200 to the recorder. The subsequent
		// http.Error() in printer() then appends the error message to the
		// body but cannot rewrite the status code. We assert what actually
		// happens (status stays 200, error text appears in body) — fixing
		// the "should return 500" wart is a separate change, out of scope
		// for the G708 hardening.
		{
			"negative_unknown_field_appends_error_in_body", "negative", hostengineTmpl,
			struct{ Unrelated string }{"x"},
			http.StatusOK, "can't evaluate field Memory", "",
		},

		// boundary — empty/zero stats render zero values, no panic.
		{
			"boundary_empty_hostengine_renders_zeros", "boundary", hostengineTmpl,
			dcgm.Status{},
			http.StatusOK, "Memory(KB)      : 0", "",
		},
		{
			"boundary_health_no_watches_renders_header_only", "boundary", healthStatusTmpl,
			dcgm.DeviceHealth{GPU: 0, Status: "OK"},
			http.StatusOK, "GPU                : 0", "",
		},
		{
			"boundary_health_with_one_watch_renders_block", "boundary", healthStatusTmpl,
			dcgm.DeviceHealth{
				GPU:    1,
				Status: "Warn",
				Watches: []dcgm.SystemWatch{
					{Type: "PCIe", Status: "Warn", Error: "link down"},
				},
			},
			http.StatusOK, "Type               : PCIe", "",
		},

		// corner — unicode and very large stats fields render through as data.
		{
			"corner_unicode_in_status_field", "corner", healthStatusTmpl,
			dcgm.DeviceHealth{GPU: 0, Status: "Stätüs—✓"},
			http.StatusOK, "Stätüs—✓", "",
		},
		{
			"corner_oversized_status_string", "corner", healthStatusTmpl,
			dcgm.DeviceHealth{GPU: 0, Status: strings.Repeat("A", 10_000)},
			http.StatusOK, strings.Repeat("A", 100), "",
		},

		// attacker — values that LOOK like template directives must be rendered
		// as data, NOT executed. This is the core security assertion: pre-parsed
		// templates do not re-parse the stats payload, so {{...}} embedded in
		// data round-trips verbatim into the response body.
		{
			"attacker_action_in_status_is_literal_data", "attacker", healthStatusTmpl,
			dcgm.DeviceHealth{GPU: 0, Status: "{{print 1337}}"},
			http.StatusOK, "{{print 1337}}", "1337\n",
		},
		{
			"attacker_pipeline_in_status_is_literal_data", "attacker", healthStatusTmpl,
			dcgm.DeviceHealth{GPU: 0, Status: "{{ .Foo | call }}"},
			http.StatusOK, "{{ .Foo | call }}", "",
		},
		{
			"attacker_html_script_tag_passes_through_as_text", "attacker", healthStatusTmpl,
			dcgm.DeviceHealth{GPU: 0, Status: "<script>alert(1)</script>"},
			// text/template does not HTML-escape — the script tag travels through
			// as plain text. The assertion is that it is NOT executed by the
			// templating engine (no eval of the inner expression).
			http.StatusOK, "<script>alert(1)</script>", "",
		},
		{
			"attacker_watch_error_field_with_action_is_literal", "attacker", healthStatusTmpl,
			dcgm.DeviceHealth{
				GPU:    0,
				Status: "OK",
				Watches: []dcgm.SystemWatch{
					{Type: "X", Status: "X", Error: `{{exec "whoami"}}`},
				},
			},
			http.StatusOK, `{{exec "whoami"}}`, "uid=",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			printer(rec, req, c.stats, c.tmpl)

			assert.Equal(t, c.wantCode, rec.Code, c.category+": status code")
			if c.wantSubstr != "" {
				assert.Contains(t, rec.Body.String(), c.wantSubstr,
					c.category+": expected substring missing")
			}
			if c.notSubstr != "" {
				assert.NotContains(t, rec.Body.String(), c.notSubstr,
					c.category+": stats data must NOT be re-evaluated as template")
			}
		})
	}
}

// TestProcessPrint mirrors TestPrinter but for the slice-driven processPrint
// helper, which iterates and renders processInfoTmpl per element.
func TestProcessPrint(t *testing.T) {
	cases := []struct {
		name       string
		category   string
		stats      []dcgm.ProcessInfo
		wantCode   int
		wantSubstr string
		notSubstr  string
	}{
		// positive — single process renders identifying fields.
		{
			"positive_single_process_renders_pid_and_gpu", "positive",
			[]dcgm.ProcessInfo{{GPU: 0, PID: 42, Name: "nvidia-smi"}},
			http.StatusOK, "PID                          : 42", "",
		},

		// negative — empty slice produces empty body, status stays 200
		// (the handler short-circuits before reaching processPrint in real flow;
		// here we exercise the helper directly).
		{
			"negative_nil_slice_no_writes", "negative",
			nil,
			http.StatusOK, "", "GPU ID",
		},

		// boundary — multiple processes render multiple blocks.
		{
			"boundary_two_processes_render_both", "boundary",
			[]dcgm.ProcessInfo{
				{GPU: 0, PID: 1, Name: "a"},
				{GPU: 0, PID: 2, Name: "b"},
			},
			http.StatusOK, "PID                          : 2", "",
		},

		// corner — unicode in process name.
		{
			"corner_unicode_in_process_name", "corner",
			[]dcgm.ProcessInfo{{GPU: 0, PID: 1, Name: "進程"}},
			http.StatusOK, "進程", "",
		},

		// attacker — template-action lookalike in the Name field must not be
		// re-evaluated by the pre-parsed processInfoTmpl.
		{
			"attacker_action_in_name_is_literal_data", "attacker",
			[]dcgm.ProcessInfo{{GPU: 0, PID: 1, Name: "{{print 9000}}"}},
			http.StatusOK, "{{print 9000}}", "9000\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			processPrint(rec, req, c.stats)

			assert.Equal(t, c.wantCode, rec.Code, c.category+": status code")
			if c.wantSubstr != "" {
				assert.Contains(t, rec.Body.String(), c.wantSubstr,
					c.category+": expected substring missing")
			}
			if c.notSubstr != "" {
				assert.NotContains(t, rec.Body.String(), c.notSubstr,
					c.category+": stats data must NOT be re-evaluated as template")
			}
		})
	}
}
