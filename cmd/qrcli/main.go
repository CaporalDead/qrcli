// Command qrcli renders QR codes in the terminal.
//
// Usage is intentionally minimal: the payload comes from the arguments or
// from stdin, and the QR code is written to stdout as plain UTF-8 text.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/CaporalDead/qrcli/internal/render"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

// quietZone is the number of light modules drawn around the code.
// See the rendering decision (issue #2) for why it is 2 and not the
// spec-mandated 4.
const quietZone = 2

const usage = `qrcli — generate QR codes in your terminal

Usage:
  qrcli [options] <text>
  <command> | qrcli [options]

Options:
  -l, --level L|M|Q|H  error correction level (default M)
  -i, --invert         invert colors, for light terminal backgrounds
  -v, --version        print version and exit
  -h, --help           show this help

Examples:
  qrcli "https://example.org"
  qrcli -l H "critical payload"
  echo -n "hello" | qrcli -i
  qrcli "WIFI:T:WPA;S:MyNetwork;P:secret;;"

Exit codes: 0 success, 1 generation error, 2 usage error.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run is the testable entry point: all I/O is injected.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("qrcli", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, usage) }

	var (
		level       string
		invert      bool
		showVersion bool
	)
	fs.StringVar(&level, "l", "M", "error correction level")
	fs.StringVar(&level, "level", "M", "error correction level")
	fs.BoolVar(&invert, "i", false, "invert colors")
	fs.BoolVar(&invert, "invert", false, "invert colors")
	fs.BoolVar(&showVersion, "v", false, "print version")
	fs.BoolVar(&showVersion, "version", false, "print version")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	if showVersion {
		fmt.Fprintf(stdout, "qrcli %s\n", version)
		return 0
	}

	recovery, err := parseLevel(level)
	if err != nil {
		fmt.Fprintf(stderr, "qrcli: %v\n", err)
		return 2
	}

	text, err := input(fs.Args(), stdin)
	if err != nil {
		fmt.Fprintf(stderr, "qrcli: %v\n\n", err)
		fs.Usage()
		return 2
	}

	code, err := qrcode.New(text, recovery)
	if err != nil {
		fmt.Fprintf(stderr, "qrcli: %v\n", err)
		return 1
	}
	// The renderer draws its own, smaller quiet zone.
	code.DisableBorder = true

	fmt.Fprint(stdout, render.Terminal(code.Bitmap(), render.Options{
		Invert:    invert,
		QuietZone: quietZone,
	}))
	return 0
}

// input returns the QR payload: positional arguments joined by spaces, or
// stdin when no argument is given. Reading from an interactive terminal
// with no argument is rejected instead of blocking forever.
func input(args []string, stdin io.Reader) (string, error) {
	if len(args) > 0 {
		text := strings.Join(args, " ")
		if text == "" {
			return "", errors.New("empty input")
		}
		return text, nil
	}

	if f, ok := stdin.(*os.File); ok {
		if fi, err := f.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
			return "", errors.New("no input: pass text as an argument or pipe it on stdin")
		}
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}
	text := strings.TrimSuffix(string(data), "\n")
	text = strings.TrimSuffix(text, "\r")
	if text == "" {
		return "", errors.New("empty input")
	}
	return text, nil
}

// parseLevel maps the CLI level names (from the QR spec) to the library's
// recovery levels.
func parseLevel(s string) (qrcode.RecoveryLevel, error) {
	switch strings.ToUpper(s) {
	case "L":
		return qrcode.Low, nil
	case "M":
		return qrcode.Medium, nil
	case "Q":
		return qrcode.High, nil
	case "H":
		return qrcode.Highest, nil
	}
	return 0, fmt.Errorf("invalid error correction level %q (want L, M, Q or H)", s)
}
