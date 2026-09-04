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
			name:         "content too long for a QR code",
			args:         []string{"-l", "H", strings.Repeat("a", 3000)},
			wantCode:     1,
			wantInStderr: "qrcli:",
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
