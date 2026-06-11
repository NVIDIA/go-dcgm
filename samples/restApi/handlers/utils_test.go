package handlers

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLogRequestError pins the contract for the log-injection-safe error
// log helper. Every case asserts the blanket invariant that no raw control
// character (LF, CR, NUL, ESC) survives the helper into the final log line —
// this is what closes gosec G706. Attacker cases drive in classic
// log-forgery payloads (newline-injected fake log entries, CRLF, null-byte
// truncation, BEL/backspace terminal abuse, ANSI clear-screen, Unicode line
// separators).
func TestLogRequestError(t *testing.T) {
	cases := []struct {
		name       string
		category   string // positive | boundary | corner | attacker
		host       string
		urlPath    string
		msg        string
		wantSubstr []string // each must appear in the captured log line
	}{
		// positive — benign inputs render in quoted form, readable.
		{
			"positive_benign_request", "positive",
			"example.com", "/devices/0", "bad gateway",
			[]string{`host="example.com"`, `url="/devices/0"`, `msg="bad gateway"`},
		},

		// boundary — empty fields don't panic; quoted empty strings appear.
		{
			"boundary_empty_fields", "boundary",
			"", "", "",
			[]string{`host=""`, `url=""`, `msg=""`},
		},

		// corner — multibyte UTF-8 in legitimate input is preserved through %q
		// (Go's %q keeps printable Unicode as-is, only escapes non-printables).
		{
			"corner_unicode_in_host", "corner",
			"exämple.com", "/", "ok",
			[]string{`exämple.com`},
		},

		// attacker — LF in Host: classic log-forgery via injected fake entry.
		// %q must render the LF as the two-char escape \n.
		{
			"attacker_lf_in_host_forges_fake_entry", "attacker",
			"evil.com\n2026-05-24 ATTACKER bypassed auth", "/", "ok",
			[]string{`\n`, "ATTACKER bypassed auth"}, // escaped form survives as data
		},

		// attacker — CRLF: forges an HTTP-style header-looking line.
		{
			"attacker_crlf_in_msg", "attacker",
			"evil", "/", "boom\r\nFAKE-HEADER: x",
			[]string{`\r\n`, "FAKE-HEADER"},
		},

		// attacker — null byte: defeats parsers that null-terminate strings.
		{
			"attacker_null_byte_in_msg_does_not_truncate", "attacker",
			"evil", "/", "before\x00after",
			[]string{`\x00`, "after"},
		},

		// attacker — BEL (\a). Some terminal log viewers beep on this byte.
		// Go's %q renders BEL as \a (Go-syntax escape).
		{
			"attacker_bel_in_msg", "attacker",
			"evil", "/", "spam\a\a\a",
			[]string{`\a`},
		},

		// attacker — backspace can erase preceding log content on terminals.
		// Go's %q renders BS as \b.
		{
			"attacker_backspace_in_msg", "attacker",
			"evil", "/", "log entry\b\b\bFAKE",
			[]string{`\b`, "FAKE"},
		},

		// attacker — ESC sequence (\x1b[2J = ANSI clear-screen). Could wipe
		// a developer's terminal mid-tail.
		{
			"attacker_ansi_clear_screen_in_msg", "attacker",
			"evil", "/", "before\x1b[2Jafter",
			[]string{`\x1b`, "after"},
		},

		// attacker — Unicode line/paragraph separators (U+2028, U+2029) are
		// treated as line breaks by some log parsers / JSON.parse semantics.
		// Go's %q renders them as   /  .
		{
			"attacker_unicode_line_separator_in_msg", "attacker",
			"evil", "/", "before injected",
			[]string{`\u2028`},
		},
		{
			"attacker_unicode_paragraph_separator_in_host", "attacker",
			"evil spoof", "/", "ok",
			[]string{`\u2029`},
		},

		// attacker — embedded double quotes try to escape the quoted region.
		// Go's %q backslash-escapes them.
		{
			"attacker_embedded_quote_in_msg", "attacker",
			"evil", "/", `" host="trusted.example.com`,
			[]string{`\"`}, // the embedded quote is escaped
		},

		// attacker — long string: must not crash, helper truncates nothing
		// (caller could add their own length cap; the security property is
		// just that it doesn't break log-line invariants).
		{
			"attacker_oversized_msg_does_not_break_invariants", "attacker",
			"evil", "/", strings.Repeat("A", 10_000),
			[]string{strings.Repeat("A", 100)},
		},

		// attacker — mixed payload: every category combined in one request.
		// Note: CRLF intentionally goes in msg, not URL path. url.URL.String()
		// percent-encodes control characters in Path (\r → %0D, \n → %0A) so
		// the URL field is already cleaned before %q sees it; we exercise %q
		// directly by routing the CRLF through msg.
		{
			"attacker_mixed_lf_crlf_null_esc_in_one_request", "attacker",
			"a\nb", "/safe", "d\r\ne\x00f\x1bg",
			[]string{`\n`, `\r\n`, `\x00`, `\x1b`},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Capture log output. log.Default() is the global logger that
			// every log.Printf call in handlers/ writes to.
			var buf bytes.Buffer
			origOut := log.Default().Writer()
			origFlags := log.Default().Flags()
			log.SetOutput(&buf)
			log.SetFlags(0) // drop date/time prefix for cleaner assertion
			t.Cleanup(func() {
				log.SetOutput(origOut)
				log.SetFlags(origFlags)
			})

			req := &http.Request{
				Host: c.host,
				URL:  &url.URL{Path: c.urlPath},
			}

			logRequestError(req, c.msg)

			got := strings.TrimSuffix(buf.String(), "\n")

			// Blanket security invariant: NO raw control character survives
			// into the log line. This is what closes gosec G706 — an
			// attacker cannot inject a synthetic log entry, terminal escape,
			// null-byte truncation, etc., regardless of which field carried
			// the payload.
			for _, badByte := range []struct{ name, b string }{
				{"LF", "\n"},
				{"CR", "\r"},
				{"NUL", "\x00"},
				{"ESC", "\x1b"},
				{"BEL", "\a"},
				{"BS", "\b"},
				{"U+2028", " "},
				{"U+2029", " "},
			} {
				assert.False(t, strings.Contains(got, badByte.b),
					c.category+": raw "+badByte.name+" must not survive into log line")
			}

			// Per-case substrings.
			for _, s := range c.wantSubstr {
				assert.Contains(t, got, s,
					c.category+": expected substring %q in log line", s)
			}
		})
	}
}

// TestLogRequestError_NilURL guards the defensive nil-URL path. Synthetic
// or partially-initialized *http.Request values can have a nil URL pointer;
// the helper must not panic.
func TestLogRequestError_NilURL(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags)
	})

	req := &http.Request{Host: "host", URL: nil}
	assert.NotPanics(t, func() {
		logRequestError(req, "msg")
	})
	got := strings.TrimSuffix(buf.String(), "\n")
	assert.Contains(t, got, `url=""`, "nil URL renders as empty quoted string")
}
