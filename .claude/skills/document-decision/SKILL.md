---
name: document-decision
description: Record a design decision (ADR) in the qrcli knowledge base — when a non-obvious choice was made, a settled decision needs superseding, or a pitfall/idea should be captured.
---

# Documenting a decision in the knowledge base

qrcli's issue tracker is its architecture decision log. Code review answers
"is this right?"; the KB answers "why is it like this?" two years later.

## When to write one

- You chose between ≥2 plausible approaches and the loser would tempt a future
  contributor (library choice, algorithm, format, scope call).
- You are about to contradict an existing `decision` issue.
- You hit a trap that cost real time → that's a `pitfall`, lighter format.
- You had an out-of-scope idea → comment on the idea backlog (#5) instead.

## How

1. Search first: `gh issue list --label decision --state all --search "<topic>"`.
   Superseding an old decision? Say so explicitly and cross-link both ways.
2. Create the issue with the **Decision (ADR)** form (or `gh issue create
   --label decision`) with exactly these sections:
   - **Context** — constraints and forces, not the solution.
   - **Options considered** — table with honest pros/cons, including the boring option.
   - **Decision** — one bold sentence.
   - **Consequences** — what gets easier, what gets harder, what we accept.
3. Link it: the implementing PR says `Closes #N`; docs/architecture.md's
   decision-log table gets a row if the decision is structural.
4. Decisions are closed by the PR that implements them. A closed decision is
   still binding — reopening the debate means a *new* issue that supersedes it.

## Style

Write for the contributor who has none of today's context. Name the options
you rejected and *why* — "we picked X" without the losers is a changelog,
not a decision record.
