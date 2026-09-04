package main

import (
	"bytes"
	"strings"
	"testing"
)

// FuzzRun asserts the CLI never panics and always honors the documented
// exit-code contract (0/1/2) for arbitrary arguments. Adopted from the
// 2026-09-05 audit (issue #22).
func FuzzRun(f *testing.F) {
	f.Add("hello", "world")
	f.Add("-l", "H")
	f.Add("wifi", "--ssid")
	f.Add("--", "-x")
	f.Fuzz(func(t *testing.T, a, b string) {
		var stdout, stderr bytes.Buffer
		code := run([]string{a, b}, strings.NewReader(""), &stdout, &stderr)
		if code != 0 && code != 1 && code != 2 {
			t.Fatalf("exit code %d outside the documented contract for args %q %q", code, a, b)
		}
	})
}
