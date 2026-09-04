package main

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"
)

// exec runs the CLI with injected I/O and returns exit code, stdout, stderr.
func exec(t *testing.T, args []string, stdin string) (int, string, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := run(args, strings.NewReader(stdin), &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestRun(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		stdin        string
		wantCode     int
		wantInStdout string
		wantInStderr string
	}{
		{
			name:         "version flag",
			args:         []string{"-v"},
			wantCode:     0,
			wantInStdout: "qrcli dev",
		},
		{
			name:         "long version flag",
			args:         []string{"--version"},
			wantCode:     0,
			wantInStdout: "qrcli dev",
		},
		{
			name:         "help flag",
			args:         []string{"-h"},
			wantCode:     0,
			wantInStderr: "Usage:",
		},
		{
			name:         "unknown flag",
			args:         []string{"--nope"},
			wantCode:     2,
			wantInStderr: "not defined",
		},
		{
			name:         "no input at all",
			args:         nil,
			stdin:        "",
			wantCode:     2,
			wantInStderr: "empty input",
		},
		{
			name:         "empty argument",
			args:         []string{""},
			wantCode:     2,
			wantInStderr: "empty input",
		},
		{
			name:         "invalid level",
			args:         []string{"-l", "X", "hello"},
			wantCode:     2,
			wantInStderr: "invalid error correction level",
		},
		{
			// Beyond the capacity of ANY QR code: rejected as a usage
			// error before encoding (issue #22).
			name:         "payload beyond any QR capacity",
			args:         []string{"-l", "H", strings.Repeat("a", 3000)},
			wantCode:     2,
			wantInStderr: "exceeds the maximum QR capacity",
		},
		{
			// Fits a QR in principle, but not at level H: still the
			// encoder's call, exit 1.
			name:         "payload too long for the chosen level",
			args:         []string{"-l", "H", strings.Repeat("a", 2000)},
			wantCode:     1,
			wantInStderr: "qrcli:",
		},
		{
			name:         "stdin at exact capacity",
			args:         []string{"-l", "L"},
			stdin:        strings.Repeat("a", 2953),
			wantCode:     0,
			wantInStdout: "█",
		},
		{
			name:         "stdin at exact capacity with trailing newline",
			args:         []string{"-l", "L"},
			stdin:        strings.Repeat("a", 2953) + "\n",
			wantCode:     0,
			wantInStdout: "█",
		},
		{
			name:         "stdin one byte over capacity",
			args:         nil,
			stdin:        strings.Repeat("a", 2954),
			wantCode:     2,
			wantInStderr: "exceeds the maximum QR capacity",
		},
		{
			// The read is bounded: a huge pipe fails fast without being
			// buffered (issue #22 measured 787 MB RSS before the fix).
			name:         "huge stdin fails fast",
			args:         nil,
			stdin:        strings.Repeat("a", 100_000),
			wantCode:     2,
			wantInStderr: "exceeds the maximum QR capacity",
		},
		{
			name:         "happy path with argument",
			args:         []string{"hello"},
			wantCode:     0,
			wantInStdout: "█",
		},
		{
			name:         "happy path with stdin",
			args:         nil,
			stdin:        "hello\n",
			wantCode:     0,
			wantInStdout: "█",
		},
		{
			name:         "explicit level",
			args:         []string{"--level", "q", "hello"},
			wantCode:     0,
			wantInStdout: "█",
		},
		{
			name:         "wifi missing ssid",
			args:         []string{"wifi"},
			wantCode:     2,
			wantInStderr: "--ssid is required",
		},
		{
			name:         "wifi help",
			args:         []string{"wifi", "-h"},
			wantCode:     0,
			wantInStderr: "qrcli wifi",
		},
		{
			name:         "wifi unexpected positional",
			args:         []string{"wifi", "--ssid", "x", "--pass", "p", "oops"},
			wantCode:     2,
			wantInStderr: "unexpected argument",
		},
		{
			name:         "wifi pass conflicts with nopass",
			args:         []string{"wifi", "--ssid", "x", "--pass", "p", "--type", "nopass"},
			wantCode:     2,
			wantInStderr: "conflicts",
		},
		{
			name:         "wifi with render flags",
			args:         []string{"wifi", "--ssid", "x", "--pass", "p", "-l", "H", "-i"},
			wantCode:     0,
			wantInStdout: "█",
		},
		{
			name:         "vcard missing name",
			args:         []string{"vcard"},
			wantCode:     2,
			wantInStderr: "--name is required",
		},
		{
			name:         "vcard help",
			args:         []string{"vcard", "-h"},
			wantCode:     0,
			wantInStderr: "qrcli vcard",
		},
		{
			name:         "subcommand name after flags fails loudly",
			args:         []string{"-i", "wifi", "--ssid", "x"},
			wantCode:     2,
			wantInStderr: "subcommand",
		},
		{
			name:         "literal subcommand word still works via stdin",
			args:         nil,
			stdin:        "wifi",
			wantCode:     0,
			wantInStdout: "█",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := exec(t, tt.args, tt.stdin)
			if code != tt.wantCode {
				t.Errorf("exit code = %d, want %d (stderr: %s)", code, tt.wantCode, stderr)
			}
			if tt.wantInStdout != "" && !strings.Contains(stdout, tt.wantInStdout) {
				t.Errorf("stdout %q does not contain %q", stdout, tt.wantInStdout)
			}
			if tt.wantInStderr != "" && !strings.Contains(stderr, tt.wantInStderr) {
				t.Errorf("stderr %q does not contain %q", stderr, tt.wantInStderr)
			}
		})
	}
}

func TestASCIIMode(t *testing.T) {
	code, out, _ := exec(t, []string{"-a", "hello"}, "")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(out, "#") {
		t.Error("ascii output should contain '#'")
	}
	for _, r := range out {
		if r > 127 {
			t.Fatalf("ascii output contains non-ASCII rune %q", r)
		}
	}
	// One row per module, two columns per module: "hello" is a v1 QR
	// (21 modules), so 25 lines of 50 chars with the quiet zone.
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != 25 {
		t.Errorf("got %d lines, want 25", len(lines))
	}
	for i, line := range lines {
		if len(line) != 50 {
			t.Errorf("line %d is %d chars wide, want 50", i, len(line))
		}
	}
}

func TestASCIIComposesWithInvert(t *testing.T) {
	_, normal, _ := exec(t, []string{"-a", "hello"}, "")
	_, inverted, _ := exec(t, []string{"-a", "-i", "hello"}, "")
	if normal == inverted {
		t.Fatal("inverted ascii output should differ from normal ascii output")
	}
	if utf8.RuneCountInString(normal) != utf8.RuneCountInString(inverted) {
		t.Error("invert must not change ascii output dimensions")
	}
}

func TestWiFiSubcommandMatchesRawPayload(t *testing.T) {
	// The subcommand is sugar over the documented raw format: both must
	// render byte-identically (equivalence pinned by issue #17).
	code, fromSub, _ := exec(t, []string{"wifi", "--ssid", "Home", "--pass", "hunter2"}, "")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	_, fromRaw, _ := exec(t, []string{"WIFI:T:WPA;S:Home;P:hunter2;;"}, "")
	if fromSub != fromRaw {
		t.Error("wifi subcommand and raw WIFI: payload rendered different codes")
	}
}

func TestVCardSubcommandMatchesRawPayload(t *testing.T) {
	code, fromSub, _ := exec(t, []string{"vcard", "--name", "Ada Lovelace", "--email", "ada@example.org"}, "")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	raw := "BEGIN:VCARD\nVERSION:3.0\nN:Lovelace;Ada\nFN:Ada Lovelace\nEMAIL:ada@example.org\nEND:VCARD"
	_, fromRaw, _ := exec(t, []string{raw}, "")
	if fromSub != fromRaw {
		t.Error("vcard subcommand and raw vCard payload rendered different codes")
	}
}

func TestSubcommandASCIIMode(t *testing.T) {
	_, out, _ := exec(t, []string{"wifi", "--ssid", "x", "--pass", "y", "-a"}, "")
	if !strings.Contains(out, "#") {
		t.Error("ascii output should contain '#'")
	}
	for _, r := range out {
		if r > 127 {
			t.Fatalf("ascii output contains non-ASCII rune %q", r)
		}
	}
}

func TestStdinMatchesArgument(t *testing.T) {
	// Piping the payload and passing it as an argument must produce the
	// exact same QR code; the trailing newline of the pipe is stripped.
	_, fromArg, _ := exec(t, []string{"hello"}, "")
	_, fromPipe, _ := exec(t, nil, "hello\n")
	if fromArg != fromPipe {
		t.Error("stdin and argument input rendered different codes")
	}
}

func TestMultipleArgumentsAreJoined(t *testing.T) {
	_, joined, _ := exec(t, []string{"hello world"}, "")
	_, split, _ := exec(t, []string{"hello", "world"}, "")
	if joined != split {
		t.Error(`"hello world" and "hello" "world" rendered different codes`)
	}
}

func TestInvertFlipsEveryCell(t *testing.T) {
	_, normal, _ := exec(t, []string{"hello"}, "")
	_, inverted, _ := exec(t, []string{"-i", "hello"}, "")
	if normal == inverted {
		t.Fatal("inverted output should differ from normal output")
	}
	// Cell count must be preserved; byte count differs since blocks are
	// multi-byte UTF-8 while spaces are not.
	if utf8.RuneCountInString(normal) != utf8.RuneCountInString(inverted) {
		t.Error("invert must not change output dimensions")
	}
}

func TestOutputGeometry(t *testing.T) {
	// "hello" fits QR version 1 (21 modules); with a quiet zone of 2 the
	// output is 25 cells wide and ceil(25/2) = 13 lines tall.
	_, out, _ := exec(t, []string{"hello"}, "")
	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) != 13 {
		t.Errorf("got %d lines, want 13", len(lines))
	}
	for i, line := range lines {
		if n := len([]rune(line)); n != 25 {
			t.Errorf("line %d is %d cells wide, want 25", i, n)
		}
	}
}

func TestResolveVersion(t *testing.T) {
	// Test binaries carry no usable module version, so the ldflags
	// default must win in both directions.
	if got := resolveVersion(); got != "dev" {
		t.Errorf("resolveVersion() in tests = %q, want %q", got, "dev")
	}

	old := version
	defer func() { version = old }()
	version = "1.2.3"
	if got := resolveVersion(); got != "1.2.3" {
		t.Errorf("resolveVersion() with stamped version = %q, want %q", got, "1.2.3")
	}
}

func TestParseLevel(t *testing.T) {
	for _, valid := range []string{"L", "M", "Q", "H", "l", "m", "q", "h"} {
		if _, err := parseLevel(valid); err != nil {
			t.Errorf("parseLevel(%q) unexpected error: %v", valid, err)
		}
	}
	for _, invalid := range []string{"", "A", "high", "0"} {
		if _, err := parseLevel(invalid); err == nil {
			t.Errorf("parseLevel(%q) expected an error", invalid)
		}
	}
}
