# AGENTS.md — guide for AI agents working on qrcli

## What this is

A ~200-line Go CLI that renders QR codes **in the terminal only**. One static
binary, one dependency, four flags. Your job when touching it: keep it that
way, and leave the knowledge base richer than you found it.

## Project map

```
cmd/qrcli/        CLI surface: flags, input (args|stdin), exit codes, version
internal/render/  pure [][]bool → string half-block renderer (no ANSI)
.github/          CI (test matrix + lint + nix), release (GoReleaser on tags)
flake.nix         Nix package + devShell (NixOS is a first-class target)
docs/             architecture (Mermaid diagrams), definition-of-done
```

## Commands

```console
make build test lint fmt cover   # the whole local loop
nix develop                      # full toolchain on Nix systems
nix build                        # package + runs the test suite
```

CI reads the Go version from `go.mod` — that file is the single source of
truth for toolchain versions.

## Guardrails

- **Scope**: console output only. No image export, no config file, no network.
  Overturning this requires a new `decision` issue, not code.
- **Dependencies**: currently exactly one (`skip2/go-qrcode`, zero transitive).
  Adding any dependency is a design decision → issue first.
- **CLI stability is the public API**: flags, exit codes (0/1/2) and the
  no-ANSI output contract are semver-guarded. Renaming a flag is MAJOR.
- **Renderer stays pure**: `internal/render` must remain side-effect-free —
  that's what makes exact-string testing possible.
- Rendered-output assertions count **runes**, never bytes (`█` is 3 bytes).
- `flake.nix` gotchas: bump `vendorHash` when `go.mod`/`go.sum` change
  (build with a fake hash, copy the `got:` value); `nix build` only sees
  **git-tracked** files — `git add` new files before diagnosing "missing file".

## Knowledge-base duties (not optional)

This repo treats GitHub as its long-term memory:

1. **Before designing anything**, search closed issues:
   `gh issue list --label decision --state all`. Do not re-litigate settled
   decisions in code; open a new `decision` issue to supersede one.
2. **Non-obvious choices** → issue with the `decision` label using the ADR
   form (context / options / decision / consequences).
3. **Traps you fell into** → `pitfall` label, or the "Pitfalls" section of
   your PR. Write what the error looked like and the fix.
4. **Out-of-scope ideas** you have while working → comment on the
   [idea backlog](https://github.com/CaporalDead/qrcli/issues/5).
5. **PRs** must fill the template: linked issue, decisions, pitfalls, semver
   impact, DoD checklist. Squash-merge; PR title = Conventional Commit.

## Conventions

- Conventional Commits (`feat`/`fix`/`docs`/`ci`/`test`/`chore`; `!` = MAJOR).
- Table-driven tests; CLI is tested through `run()` with injected I/O —
  never spawn the binary in tests.
- Match existing comment density: comments explain *why*, not *what*.
- Docs are English; keep README examples copy-pasteable.

## Skills

Step-by-step runbooks live in `.claude/skills/`:
- `release` — how to cut a semver release end to end.
- `document-decision` — how to record an ADR issue properly.
