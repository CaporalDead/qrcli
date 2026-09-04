# qrcli

[![CI](https://github.com/CaporalDead/qrcli/actions/workflows/ci.yml/badge.svg)](https://github.com/CaporalDead/qrcli/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/CaporalDead/qrcli)](https://github.com/CaporalDead/qrcli/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/CaporalDead/qrcli.svg)](https://pkg.go.dev/github.com/CaporalDead/qrcli)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![SemVer](https://img.shields.io/badge/versioning-SemVer-informational)](https://semver.org)

> Generate QR codes in your terminal. One standalone binary, zero configuration, no ANSI escapes.

```console
$ qrcli "hello"
█████████████████████████
██ ▄▄▄▄▄ ██▀█▀▀█ ▄▄▄▄▄ ██
██ █   █ █ █ ▄ █ █   █ ██
██ █▄▄▄█ █ █▄▀▄█ █▄▄▄█ ██
██▄▄▄▄▄▄▄█ █ █▄█▄▄▄▄▄▄▄██
██ ▀ ▄ ▄▄▀██  ▀▀ ▄▄  █▀██
███▄▄█ ▀▄▀███▀ ▀ ▄█  ▄███
███▄█▄██▄▄  ▀█ ██ █▀██ ██
██ ▄▄▄▄▄ █▀▄▀▄█▄█▀ ▀  ███
██ █   █ █  █ ▀ ▀▀▄▀█▄███
██ █▄▄▄█ █▄▀█▀ ▀ ▄▀▀ ████
██▄▄▄▄▄▄▄█▄███▄███▄██▄███
█████████████████████████
```

Point a phone camera at your terminal — that's it.

## Install

### Prebuilt binaries (Linux, macOS, Windows — including NixOS)

Binaries are **fully static** (`CGO_ENABLED=0`): they run on any Linux distribution
including NixOS, with no dependencies. Grab the archive for your platform from the
[latest release](https://github.com/CaporalDead/qrcli/releases/latest), extract, done:

```console
$ tar xzf qrcli_*_linux_amd64.tar.gz
$ ./qrcli --version
```

Every release ships a `checksums.txt` for verification.

### Nix / NixOS

```console
$ nix run github:CaporalDead/qrcli -- "hello"      # try without installing
$ nix profile install github:CaporalDead/qrcli     # imperative install
```

Or declaratively, as a flake input in your NixOS / home-manager configuration:

```nix
{
  inputs.qrcli.url = "github:CaporalDead/qrcli";
  # then add to your packages:
  # inputs.qrcli.packages.${pkgs.system}.default
}
```

### From source (Go 1.24+)

```console
$ go install github.com/CaporalDead/qrcli/cmd/qrcli@latest
```

## Usage

```text
qrcli [options] <text>
<command> | qrcli [options]

  -l, --level L|M|Q|H  error correction level (default M)
  -i, --invert         invert colors, for light terminal backgrounds
  -a, --ascii          pure-ASCII output, for terminals without Unicode blocks
  -v, --version        print version and exit
  -h, --help           show this help
```

### Examples

```console
$ qrcli "https://example.org"                    # share a URL
$ echo -n "some secret" | qrcli                  # pipe from stdin
$ qrcli -l H "survives 30% damage"               # highest error correction
$ qrcli -i "hello"                               # on a light terminal theme
$ qrcli "WIFI:T:WPA;S:MyNetwork;P:hunter2;;"     # Wi-Fi credentials
$ qrcli -a "hello"                               # 7-bit ASCII (##), no Unicode
$ qrcli -- "-starts-with-a-dash"                 # payload starting with '-'
```

Error correction levels trade capacity for damage resistance:
`L` ≈ 7 %, `M` ≈ 15 % (default), `Q` ≈ 25 %, `H` ≈ 30 % of the code may be
unreadable and still scan.

### Exit codes

| Code | Meaning |
|---|---|
| 0 | success |
| 1 | generation error (e.g. payload too large for a QR code) |
| 2 | usage error (bad flag, bad level, empty input) |

### Troubleshooting

- **The code doesn't scan on a light terminal** → use `-i`. The default polarity
  is tuned for dark backgrounds (see the
  [rendering decision](https://github.com/CaporalDead/qrcli/issues/2)).
- **Mangled characters** (legacy `cmd.exe`, serial console, exotic fonts) → use
  `-a/--ascii` for pure 7-bit output, or switch to a UTF-8 terminal
  (`chcp 65001` on legacy Windows consoles). Note that ASCII mode is twice as
  large in both directions (a small QR needs 50 columns).
- **Huge QR overflows the terminal** → lower the error correction (`-l L`) or
  shorten the payload (URL shortener); a QR code's size grows with content.

## Design

The whole tool is ~200 lines on top of a single dependency
([`skip2/go-qrcode`](https://github.com/skip2/go-qrcode), itself dependency-free).
Architecture, diagrams and the error-handling contract live in
[docs/architecture.md](docs/architecture.md).

Every non-obvious choice is recorded in the issue tracker as an ADR — browse the
[`decision` label](https://github.com/CaporalDead/qrcli/issues?q=label%3Adecision)
for the full log (language choice, rendering strategy, quiet zone size, …), and the
[`pitfall` label](https://github.com/CaporalDead/qrcli/issues?q=label%3Apitfall) for
the traps we already fell into so you don't have to.

## Versioning

qrcli follows [Semantic Versioning](https://semver.org). For a CLI, the public API
is the CLI surface: flags, arguments, exit codes and the output contract.
Breaking any of those bumps MAJOR; new capabilities bump MINOR; fixes bump PATCH.
See [CHANGELOG.md](CHANGELOG.md).

## Contributing

Issues and PRs welcome — start with [CONTRIBUTING.md](CONTRIBUTING.md) and the
[Definition of Done](docs/definition-of-done.md). AI agents: read
[AGENTS.md](AGENTS.md) first.

## License

[MIT](LICENSE)
