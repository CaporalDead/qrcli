---
name: release
description: Cut a semver release of qrcli — preflight checks, changelog, flake version, tag, verify GoReleaser assets. Use when asked to release, publish or tag a new version.
---

# Releasing qrcli

## 1. Decide the version

Inspect `main` since the last tag: `git log $(git describe --tags --abbrev=0)..HEAD --oneline`.
Map Conventional Commits to the bump: any `feat!`/`BREAKING CHANGE` → MAJOR,
else any `feat` → MINOR, else any `fix` → PATCH. Only `docs`/`ci`/`chore`/`test`
→ no release needed; stop and say so.

## 2. Preflight (on a PR, not on main)

- [ ] `gh run list --branch main --limit 1` → latest CI on main is green.
- [ ] Move the *Unreleased* section of `CHANGELOG.md` to `## [X.Y.Z] - YYYY-MM-DD`.
- [ ] Update `version = "X.Y.Z"` in `flake.nix` (it is NOT derived from the tag).
- [ ] If `go.mod`/`go.sum` changed since last release, confirm `nix build` still
      passes (vendorHash may need recomputing).
- [ ] Open the PR (`docs: prepare vX.Y.Z`), get CI green, squash-merge.

## 3. Tag

```console
git switch main && git pull
git tag -a vX.Y.Z -m "qrcli vX.Y.Z"
git push origin vX.Y.Z
```

## 4. Verify (do not skip)

- `gh run watch $(gh run list --workflow=release.yml --limit 1 --json databaseId -q '.[0].databaseId') --exit-status`
- `gh release view vX.Y.Z --json assets -q '.assets[].name'` → expect 6 archives
  + `checksums.txt` (linux/darwin/windows × amd64/arm64).
- Smoke-test one asset: download, extract, `./qrcli --version` must print vX.Y.Z.
- `nix run github:CaporalDead/qrcli/vX.Y.Z -- "release smoke"` renders a QR.

## 5. Aftercare

- Close the milestone if one exists; open the next one.
- If anything failed, fix forward: delete the tag only if the release never
  published (`git push --delete origin vX.Y.Z`), and record the failure as a
  `pitfall` issue.
