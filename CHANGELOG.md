# Changelog

All notable changes to qrcli are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions follow
[Semantic Versioning](https://semver.org). For a CLI, the public API is the
CLI surface: flags, arguments, exit codes and the output contract.

## [Unreleased]

### Added

- `-a/--ascii`: pure 7-bit ASCII rendering (`##`), a fallback for terminals
  without Unicode block glyphs; composes with `-i/--invert` ([#15]).

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
[#2]: https://github.com/CaporalDead/qrcli/issues/2
[#3]: https://github.com/CaporalDead/qrcli/issues/3
[#4]: https://github.com/CaporalDead/qrcli/issues/4
[Unreleased]: https://github.com/CaporalDead/qrcli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/CaporalDead/qrcli/releases/tag/v0.1.0
