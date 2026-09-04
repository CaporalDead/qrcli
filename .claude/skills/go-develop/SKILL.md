---
name: go-develop
description: Develop Go code in qrcli the project way — design gates, purity rules, table-driven exact-string tests, the local quality loop and PR conventions. Use when implementing, refactoring or reviewing any Go change in this repo.
---

# Developing Go in qrcli

The codebase is ~600 lines and intends to stay legible. These rules are how it
reached 97 % coverage with zero mocking libraries — work with them, not around them.

## 0. Before writing code

1. `gh issue list --label decision --state all --search "<topic>"` — settled
   decisions are binding; supersede with a new ADR (see the `document-decision`
   skill), never by silently diverging.
2. Re-read the guardrails in AGENTS.md: console-only scope, single-dependency
   rule, CLI stability = public API (flags, args, exit codes, output contract).
3. A non-obvious design choice ahead? ADR issue first, code second.

## 1. Architecture rules

- **Purity boundary**: `internal/*` packages are pure functions of their
  inputs — no I/O, no globals, no clocks. All I/O, flag parsing and exit codes
  live in `cmd/qrcli`. Put new logic on the pure side whenever possible; that
  is what keeps exact-string testing viable.
- **stdlib first**: a new dependency is an ADR, not an import. The tree is
  1 direct / 0 transitive and the audit story depends on it.
- **Errors**: user-shaped problems (flags, fields, input) → exit 2 with usage;
  encoding failures → exit 1; success → 0. Wrap with `fmt.Errorf("...: %w", err)`.
- **Comments** explain *why* (constraints, spec references, issue links),
  never *what* the next line does.

## 2. Testing patterns (the pitfalls here are already paid for)

- **Pure functions → table-driven exact-string tests.** Hand-compute the
  expected string in the table; compare with `%q` in the error message.
  Models: `internal/render/render_test.go`, `internal/payload/payload_test.go`.
- **CLI → through `run()` with injected I/O.** Never exec the binary, never
  touch a real TTY. Use the `exec(t, args, stdin)` helper.
- **Dimension assertions count runes, never bytes** — `█` is 3 bytes and a
  space is 1 (pitfall, PR #7).
- **Sugar must equal raw**: any convenience path gets an equivalence test
  pinning byte-identical output against the underlying raw form (pattern and
  guardrail from PR #18).
- **Method values copy value receivers** at evaluation time — before
  `flag.Parse` fills the fields. Pass `func() (string, error) { return x.Encode() }`,
  never `x.Encode` (pitfall, PR #18).
- Every `return err` branch appears in some test table.

## 3. Quality loop — all green before opening the PR

```console
$ make fmt        # gofmt + go mod tidy — must produce zero diff
$ make test       # go test -race ./...
$ make cover      # keep total ≳ 95%; justify any drop in the PR body
$ make lint       # golangci-lint v2 (.golangci.yml)
$ nix build       # hermetic build + full test suite (CI runs it anyway)
```

## 4. Shipping

Conventional Commit title (it becomes the squash commit on main). PR template
filled: Decisions and Pitfalls are prose — "None" is acceptable, empty is not.
State the SemVer impact. Checklist: docs/definition-of-done.md.
