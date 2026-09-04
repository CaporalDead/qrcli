# Contributing to qrcli

Thanks for stopping by! This project is intentionally tiny — the bar for
changes is not size, it's *traceability*: we should still understand in two
years why every line exists.

## Dev environment

Two equivalent options:

```console
$ nix develop        # NixOS/Nix users: go, golangci-lint, goreleaser, make — done
```

or install Go 1.24+ and [golangci-lint v2](https://golangci-lint.run) yourself.

Daily driver commands (see `Makefile`):

```console
$ make build    # static binary in ./qrcli
$ make test     # go test -race ./...
$ make cover    # coverage report
$ make lint     # golangci-lint run
$ make fmt      # gofmt + go mod tidy
```

`nix build` also runs the full test suite inside the sandbox.

## Workflow

1. **Start from an issue.** Bug, feature or decision — if none exists, open one
   (forms are provided). Non-obvious design choices get the `decision` label and
   an ADR-style body: context → options considered → decision → consequences.
   Search the [decision log](https://github.com/CaporalDead/qrcli/issues?q=label%3Adecision)
   first; a closed decision is overturned by a *new* issue, not by silently
   diverging from it.
2. **Branch** from `main`: `feat/...`, `fix/...`, `docs/...`, `ci/...`.
3. **Commit** using [Conventional Commits](https://www.conventionalcommits.org):
   `feat:` (MINOR), `fix:` (PATCH), `feat!:` / `BREAKING CHANGE:` (MAJOR),
   `docs:`/`ci:`/`test:`/`chore:` (no release).
4. **Open a PR** — the template asks for linked issues, decisions made,
   pitfalls encountered, and semver impact. That content *is* the knowledge
   base; "small fix, no details" PRs are the only thing we push back on.
5. **CI must be green** (test matrix, lint, nix). PRs are **squash-merged**, so
   the PR title becomes the conventional commit on `main`.

`main` is protected: the five CI checks are required, branches must be up to
date, direct pushes are disabled (PRs only, maintainers included) and history
is linear (squash is the only enabled merge method). Repository admins can
bypass in emergencies — doing so warrants a `pitfall` issue explaining why.

The full checklist lives in [docs/definition-of-done.md](docs/definition-of-done.md).

## What lands, what doesn't

- The product statement is **console-only, one static binary, minimal flags**.
  Features that fight this (image export, config files, daemons) need a
  `decision` issue making the case first — see the
  [idea backlog](https://github.com/CaporalDead/qrcli/issues/5) for what has
  already been parked and why.
- New dependencies are a decision, not a convenience: the current tree is one
  direct dependency with zero transitive ones, and we like it that way.

## Releasing (maintainers)

The pipeline is described in [#3](https://github.com/CaporalDead/qrcli/issues/3);
the step-by-step runbook is `.claude/skills/release/SKILL.md` (readable by
humans too). Short version:

1. Update `CHANGELOG.md`, bump `version` in `flake.nix`, land via PR.
2. `git tag -a vX.Y.Z -m "qrcli vX.Y.Z" && git push origin vX.Y.Z`.
3. GoReleaser publishes the release; verify the assets and smoke-test one binary.

## AI agents

Machine-oriented guidance (project map, guardrails, KB conventions) is in
[AGENTS.md](AGENTS.md). Humans are welcome to read it; agents are required to.
