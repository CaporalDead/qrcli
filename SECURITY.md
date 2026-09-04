# Security Policy

## Supported versions

Only the **latest release** receives security fixes. qrcli is a small,
offline, network-free CLI — its attack surface is the payload text you feed
it and the supply chain that builds it.

## Reporting a vulnerability

Please report vulnerabilities **privately** via GitHub:
[Report a vulnerability](https://github.com/CaporalDead/qrcli/security/advisories/new)
(Security → Report a vulnerability). Do not open a public issue for
exploitable problems.

You can expect an acknowledgement within a week. Fixes ship as a PATCH
release with credit in the release notes unless you prefer otherwise.

## Verifying what you run

- Release archives are covered by `checksums.txt`, a per-archive SPDX SBOM,
  and a build-provenance attestation:
  `gh attestation verify <asset> --repo CaporalDead/qrcli`
- Binaries are fully static (`CGO_ENABLED=0`) with exactly one direct,
  transitively-empty dependency; `govulncheck` runs weekly in CI.
