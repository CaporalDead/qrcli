# Changelog

All notable changes to qrcli are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[Semantic Versioning](https://semver.org). For a CLI, the public API is the
CLI surface: flags, arguments, exit codes and the output contract.

## [Unreleased]

### Fixed

- `--version` now reports the real module version for
  `go install .../cmd/qrcli@vX.Y.Z` builds via `debug.ReadBuildInfo`,
  instead of `dev`; ldflags-stamped release binaries are unaffected ([#26]).
- `wifi`: SSIDs and passwords made solely of hex digits are now
  double-quoted in the payload, so strict parsers don't misread them as
  hex-encoded (ZXing recommendation) ([#25]).
- `vcard`: `--name "Family, Given"` now maps the family name correctly
  instead of splitting on the last word ([#25]).

### Security

- Releases now ship a build-provenance attestation
  (`gh attestation verify <asset> --repo CaporalDead/qrcli`) and one SPDX
  SBOM per archive; CI actions are pinned to commit SHAs and a weekly lane
  runs `govulncheck` plus a fuzz smoke ([#23]).
- Standard-input ingestion is now bounded: streams beyond the maximum QR
  capacity (2953 bytes) fail fast with constant memory, instead of being
  buffered whole (a 200 MB pipe was measured at 787 MB peak RSS) before the
  encoder rejected them ([#22]).

### Changed

- Inputs larger than 2953 bytes — impossible for any QR code — now exit 2
  (usage error) with a clear message; inputs that fit a QR but not the chosen
  error-correction level still exit 1 ([#22]).

## [0.2.0] - 2026-09-05

### Added

- `-a/--ascii`: pure 7-bit ASCII rendering (`##`), a fallback for terminals
  without Unicode block glyphs; composes with `-i/--invert` ([#15]).
- `wifi` subcommand: encode Wi-Fi credentials (`--ssid`, `--pass`, `--type`,
  `--hidden`) in the standard `WIFI:` format with proper escaping ([#17]).
- `vcard` subcommand: encode a vCard 3.0 contact card (`--name`, `--tel`,
  `--email`, `--org`, `--url`) with proper escaping ([#17]).
- Rendering flags (`-l`, `-i`, `-a`) are accepted after a subcommand.

### Changed

- The words `wifi` and `vcard` are now reserved as first argument (subcommand
  dispatch). To encode those literal words, pipe them on stdin
  (`echo -n wifi | qrcli`); a subcommand name appearing after flags is a loud
  usage error instead of silently joining the payload ([#17]).

## [0.1.0] - 2026-09-04

### Added

- QR code generation in the terminal from an argument or stdin
  (`qrcli "text"`, `cmd | qrcli`), rendered with Unicode half-blocks,
  no ANSI escapes ([#1], [#2]).
- `-l/--level` error correction selection (L/M/Q/H, default M).
- `-i/--invert` polarity flip for light terminal backgrounds.
- `-v/--version`, `-h/--help`; exit-code contract 0/1/2.
- Static release binaries for linux/darwin/windows × amd64/arm64 — run on
  any distro including NixOS ([#3]).
- Nix flake: `nix run github:CaporalDead/qrcli`, package for 4 systems,
  `nix develop` dev shell ([#4]).

[#1]: https://github.com/CaporalDead/qrcli/issues/1
[#15]: https://github.com/CaporalDead/qrcli/issues/15
[#22]: https://github.com/CaporalDead/qrcli/issues/22
[#23]: https://github.com/CaporalDead/qrcli/issues/23
[#25]: https://github.com/CaporalDead/qrcli/issues/25
[#26]: https://github.com/CaporalDead/qrcli/issues/26
[#17]: https://github.com/CaporalDead/qrcli/issues/17
[#2]: https://github.com/CaporalDead/qrcli/issues/2
[#3]: https://github.com/CaporalDead/qrcli/issues/3
[#4]: https://github.com/CaporalDead/qrcli/issues/4
[Unreleased]: https://github.com/CaporalDead/qrcli/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/CaporalDead/qrcli/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/CaporalDead/qrcli/releases/tag/v0.1.0
