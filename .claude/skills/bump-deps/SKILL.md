---
name: bump-deps
description: Land dependency, GitHub Action and toolchain bumps safely — vendorHash recompute for gomod changes, held-merge strategy for release-only actions, single-source Go version, manual nixpkgs refresh. Use when handling Dependabot PRs or upgrading Go, modules or actions.
---

# Bumping dependencies

## Go modules (`gomod` ecosystem)

1. A green PR is necessary but **not sufficient**: any `go.mod`/`go.sum`
   change invalidates the flake's `vendorHash`. The CI `nix` job is the
   backstop — if it is red on a gomod bump, this is why.
2. Recompute: set `vendorHash = nixpkgs.lib.fakeHash;` in flake.nix,
   `git add` (flakes only see tracked files), `nix build .#qrcli`, copy the
   `got:` hash back. The comment above the field documents this dance.
3. User-visible dependency changes get a `fix(deps)` CHANGELOG entry.

## GitHub Actions

- **CI-exercised actions** (checkout, setup-go, golangci-lint-action, nix
  installer): the PR's own checks prove them — merge on green.
- **Release-only actions** (`goreleaser-action`): tag-triggered workflows
  never run on PRs, so a major bump merges *blind*. **Hold it** until just
  after the next release ships, and say so in a PR comment — a failed release
  must bisect trivially between "our config" and "their major".
  Precedent: #11, held through v0.1.0, validated by v0.2.0.

## Go toolchain

- `go.mod` is the **single source of truth**: CI reads it via
  `go-version-file`. Never pin a Go version inside a workflow.
- Before raising the directive, check nixpkgs still satisfies it
  (`nix eval nixpkgs#go.version`) or `nix build` breaks for every Nix user.

## nixpkgs input

Dependabot does not cover it. Refresh deliberately: `nix flake update`,
`nix build`, commit `flake.lock` in its own `chore:` PR.

## Always

Squash-merge with the conventional title. A bump that revealed a trap feeds
the KB: `pitfall` issue or the PR's Pitfalls section.
