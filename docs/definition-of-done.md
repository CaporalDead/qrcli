# Definition of Done

A change is *done* when every applicable line below holds. The PR template
mirrors this list; check items off there.

## Implementation
- [ ] Behavior matches the linked issue; scope creep split into new issues.
- [ ] Errors follow the exit-code contract (0 / 1 / 2 — see docs/architecture.md).
- [ ] No new dependency without a `decision` issue justifying it.

## Tests
- [ ] New behavior is covered by table-driven tests (exact-string for the
      renderer, injected-I/O for the CLI).
- [ ] `make test` (race detector) passes locally.
- [ ] Dimension/geometry assertions count **runes**, not bytes.

## Quality
- [ ] `make lint` is clean, `gofmt` produces no diff.
- [ ] CI green on all jobs: ubuntu / macos / windows, lint, nix.

## Documentation
- [ ] User-facing change → README updated (flags table, examples, exit codes).
- [ ] Structural change → docs/architecture.md and its diagrams updated.
- [ ] CHANGELOG.md updated under *Unreleased* for anything user-visible.

## Traceability (the knowledge base)
- [ ] PR links its issue(s) (`Closes #N`).
- [ ] Decisions made along the way are written down (PR body, or a `decision`
      issue when they outlive the PR).
- [ ] Pitfalls encountered are recorded (PR body, or a `pitfall` issue if
      others will hit it).
- [ ] SemVer impact stated (MAJOR / MINOR / PATCH / none).

## Release readiness (when tagging)
- [ ] CHANGELOG entry moved from *Unreleased* to the version section.
- [ ] `version` in `flake.nix` matches the tag.
- [ ] Release workflow green, assets present, one binary smoke-tested.
