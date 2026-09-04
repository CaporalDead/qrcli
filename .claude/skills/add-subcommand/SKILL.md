---
name: add-subcommand
description: Add a payload-helper subcommand to qrcli (like wifi/vcard) end to end — ADR, pure builder with escaping, dispatch wiring, equivalence tests, docs and changelog. Use when asked to add a new payload type, format helper or subcommand.
---

# Adding a payload subcommand

Precedent: `wifi`/`vcard` (ADR #17, PR #18). Budget ≈ 15 lines of wiring plus
the pure builder, its tests and the docs.

## 1. ADR first — a subcommand is a CLI-shape change

Open a `decision` issue covering: the payload format (link the spec), the
minimal field set, the escaping rules, and the fact that **the subcommand name
becomes a reserved first argument** — a behavior change every single time
(stdin stays the literal-text escape hatch). Cross-link #17.

## 2. Pure builder in `internal/payload`

- A struct with exported fields + `Encode() (string, error)`. Validation
  errors are user-shaped (they surface as exit 2) — phrase them in flag terms
  (`"--ssid is required"`).
- Escaping via a package-level `strings.NewReplacer` — it substitutes in a
  single pass, so inserted backslashes are never re-escaped.
- Exact-string tests: happy paths, an escaping row, the full error matrix.

## 3. Wire the command in `cmd/qrcli`

1. `runX(args, stdout, stderr)`: `newFlagSet(name, xUsage, stderr)` +
   `registerRenderFlags(fs)` + field flags, then
   `buildPayload(fs, args, func() (string, error) { return x.Encode() }, stderr)`
   and `emit`. **Closure, not the method value** — `x.Encode` copies the
   struct before Parse runs (pitfall, PR #18).
2. Add the case to `run()`'s dispatch switch **and** to `isSubcommand`.
3. Usage: an `xUsage` const mirroring `wifiUsage`, plus one line in the root
   `usage` under Subcommands and in the reserved-words sentence.

## 4. Tests

- Rows in `TestRun`: missing required field (exit 2), `-h` (exit 0),
  unexpected positional, field conflicts.
- **Equivalence test**: subcommand output byte-identical to the documented
  raw payload string — this is an AGENTS.md guardrail, not optional.
- One composition check with `-a`.

## 5. Docs & changelog

- README: example under "Payload helpers"; add the word to the reserved-words
  note.
- docs/architecture.md: decision-log row.
- AGENTS.md: extend the reserved-words guardrail.
- CHANGELOG: *Added* (the subcommand) **and** *Changed* (the new reserved word).

## 6. SemVer

MINOR, with the reserved word recorded under *Changed*. Renaming or
repurposing an existing subcommand is MAJOR.
