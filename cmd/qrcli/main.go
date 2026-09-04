// Command qrcli renders QR codes in the terminal.
//
// Usage is intentionally minimal: the payload comes from the arguments or
// from stdin, and the QR code is written to stdout as plain UTF-8 text.
// Two subcommands (wifi, vcard) build standard payloads from fields so
// users don't have to remember those formats' escaping rules.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/CaporalDead/qrcli/internal/payload"
	"github.com/CaporalDead/qrcli/internal/render"
)

// version is stamped at build time via -ldflags "-X main.version=...".
var version = "dev"

// quietZone is the number of light modules drawn around the code.
// See the rendering decision (issue #2) for why it is 2 and not the
// spec-mandated 4.
const quietZone = 2

// maxPayload is the byte capacity of the largest QR code (version 40,
// level L, byte mode). Anything longer can be rejected before buffering
// or encoding — the bound that keeps stdin ingestion constant-memory
// (issue #22).
const maxPayload = 2953

var errTooLong = fmt.Errorf("input exceeds the maximum QR capacity (%d bytes)", maxPayload)

const usage = `qrcli — generate QR codes in your terminal

Usage:
  qrcli [options] <text>
  <command> | qrcli [options]
  qrcli wifi --ssid <name> [--pass <secret>] [options]
  qrcli vcard --name <full name> [options]

Options:
  -l, --level L|M|Q|H  error correction level (default M)
  -i, --invert         invert colors, for light terminal backgrounds
  -a, --ascii          pure-ASCII output, for terminals without Unicode blocks
  -v, --version        print version and exit
  -h, --help           show this help

Subcommands (run "qrcli <subcommand> -h" for details):
  wifi   encode Wi-Fi credentials in the standard WIFI: format
  vcard  encode a contact card (vCard 3.0)

The words "wifi" and "vcard" are reserved as first argument; to encode such
literal text, pipe it on stdin instead.

Examples:
  qrcli "https://example.org"
  qrcli -l H "critical payload"
  echo -n "hello" | qrcli -i
  qrcli wifi --ssid Home --pass hunter2

Exit codes: 0 success, 1 generation error, 2 usage error.
`

const wifiUsage = `qrcli wifi — encode Wi-Fi credentials as a QR code

Usage:
  qrcli wifi --ssid <name> [--pass <secret>] [options]

Options:
  --ssid <name>          network name (required)
  --pass <secret>        password (omit for an open network)
  --type WPA|WEP|nopass  security (default: WPA with --pass, nopass without)
  --hidden               the network does not broadcast its SSID
  -l, --level L|M|Q|H    error correction level (default M)
  -i, --invert           invert colors, for light terminal backgrounds
  -a, --ascii            pure-ASCII output
  -h, --help             show this help

Example:
  qrcli wifi --ssid "Home" --pass "hunter2" --hidden
`

const vcardUsage = `qrcli vcard — encode a contact card (vCard 3.0) as a QR code

Usage:
  qrcli vcard --name <full name> [options]

Options:
  --name <full name>   display name (required); the last word is used as
                       the family name when saving the contact
  --tel <number>       phone number
  --email <address>    email address
  --org <name>         organization
  --url <address>      website
  -l, --level L|M|Q|H  error correction level (default M)
  -i, --invert         invert colors, for light terminal backgrounds
  -a, --ascii          pure-ASCII output
  -h, --help           show this help

Example:
  qrcli vcard --name "Ada Lovelace" --tel "+44 20 7946 0958" --email ada@example.org
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run is the testable entry point: all I/O is injected. The first argument
// selects a subcommand; anything else goes through the raw-payload path.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "wifi":
			return runWiFi(args[1:], stdout, stderr)
		case "vcard":
			return runVCard(args[1:], stdout, stderr)
		}
	}

	fs := newFlagSet("qrcli", usage, stderr)
	rf := registerRenderFlags(fs)
	var showVersion bool
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

	// A subcommand name after flags would otherwise be swallowed into the
	// payload silently; fail loudly instead (see issue #17).
	if rest := fs.Args(); len(rest) > 0 && isSubcommand(rest[0]) {
		fmt.Fprintf(stderr, "qrcli: %q is a subcommand and must come first: qrcli %s [options]\n", rest[0], rest[0])
		fmt.Fprintf(stderr, "(to encode the literal text, pipe it: echo -n %q | qrcli)\n", strings.Join(rest, " "))
		return 2
	}

	text, err := input(fs.Args(), stdin)
	if err != nil {
		fmt.Fprintf(stderr, "qrcli: %v\n\n", err)
		fs.Usage()
		return 2
	}

	return emit(text, rf, stdout, stderr)
}

func runWiFi(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("qrcli wifi", wifiUsage, stderr)
	rf := registerRenderFlags(fs)
	var w payload.WiFi
	fs.StringVar(&w.SSID, "ssid", "", "network name")
	fs.StringVar(&w.Password, "pass", "", "password")
	fs.StringVar(&w.Security, "type", "", "security type")
	fs.BoolVar(&w.Hidden, "hidden", false, "hidden network")

	// Closure, not the method value w.Encode: the latter would copy the
	// struct now, before Parse fills the fields.
	text, exit, ok := buildPayload(fs, args, func() (string, error) { return w.Encode() }, stderr)
	if !ok {
		return exit
	}
	return emit(text, rf, stdout, stderr)
}

func runVCard(args []string, stdout, stderr io.Writer) int {
	fs := newFlagSet("qrcli vcard", vcardUsage, stderr)
	rf := registerRenderFlags(fs)
	var v payload.VCard
	fs.StringVar(&v.Name, "name", "", "full display name")
	fs.StringVar(&v.Tel, "tel", "", "phone number")
	fs.StringVar(&v.Email, "email", "", "email address")
	fs.StringVar(&v.Org, "org", "", "organization")
	fs.StringVar(&v.URL, "url", "", "website")

	text, exit, ok := buildPayload(fs, args, func() (string, error) { return v.Encode() }, stderr)
	if !ok {
		return exit
	}
	return emit(text, rf, stdout, stderr)
}

// buildPayload parses a subcommand's flags and encodes its payload; the
// encode closure reads the fields bound to the flag set. ok reports
// whether rendering should proceed; when false, the caller must return
// exit as-is (0 covers the help path, which is a success that renders
// nothing).
func buildPayload(fs *flag.FlagSet, args []string, encode func() (string, error), stderr io.Writer) (text string, exit int, ok bool) {
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return "", 0, false
		}
		return "", 2, false
	}
	if rest := fs.Args(); len(rest) > 0 {
		fmt.Fprintf(stderr, "%s: unexpected argument %q\n\n", fs.Name(), rest[0])
		fs.Usage()
		return "", 2, false
	}
	text, err := encode()
	if err != nil {
		fmt.Fprintf(stderr, "qrcli: %v\n\n", err)
		fs.Usage()
		return "", 2, false
	}
	return text, 0, true
}

// renderFlags are the rendering options shared by the bare command and
// every subcommand.
type renderFlags struct {
	level  string
	invert bool
	ascii  bool
}

func registerRenderFlags(fs *flag.FlagSet) *renderFlags {
	rf := &renderFlags{}
	fs.StringVar(&rf.level, "l", "M", "error correction level")
	fs.StringVar(&rf.level, "level", "M", "error correction level")
	fs.BoolVar(&rf.invert, "i", false, "invert colors")
	fs.BoolVar(&rf.invert, "invert", false, "invert colors")
	fs.BoolVar(&rf.ascii, "a", false, "pure-ASCII output")
	fs.BoolVar(&rf.ascii, "ascii", false, "pure-ASCII output")
	return rf
}

func newFlagSet(name, usageText string, stderr io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, usageText) }
	return fs
}

func isSubcommand(s string) bool {
	return s == "wifi" || s == "vcard"
}

// emit renders an already-built payload — the shared tail of the raw
// path and of every subcommand.
func emit(text string, rf *renderFlags, stdout, stderr io.Writer) int {
	recovery, err := parseLevel(rf.level)
	if err != nil {
		fmt.Fprintf(stderr, "qrcli: %v\n", err)
		return 2
	}

	code, err := qrcode.New(text, recovery)
	if err != nil {
		fmt.Fprintf(stderr, "qrcli: %v\n", err)
		return 1
	}
	// The renderer draws its own, smaller quiet zone.
	code.DisableBorder = true

	opts := render.Options{Invert: rf.invert, QuietZone: quietZone}
	if rf.ascii {
		fmt.Fprint(stdout, render.ASCII(code.Bitmap(), opts))
	} else {
		fmt.Fprint(stdout, render.Terminal(code.Bitmap(), opts))
	}
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
		if len(text) > maxPayload {
			return "", errTooLong
		}
		return text, nil
	}

	if f, ok := stdin.(*os.File); ok {
		if fi, err := f.Stat(); err == nil && fi.Mode()&os.ModeCharDevice != 0 {
			return "", errors.New("no input: pass text as an argument or pipe it on stdin")
		}
	}

	// Read at most maxPayload plus room for a trailing newline ("\r\n")
	// plus one sentinel byte: any stream that fills the limit is over
	// capacity even after stripping the newline, so nothing is ever
	// silently truncated and memory stays constant on huge pipes.
	data, err := io.ReadAll(io.LimitReader(stdin, maxPayload+3))
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}
	text := strings.TrimSuffix(string(data), "\n")
	text = strings.TrimSuffix(text, "\r")
	if text == "" {
		return "", errors.New("empty input")
	}
	if len(text) > maxPayload {
		return "", errTooLong
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
