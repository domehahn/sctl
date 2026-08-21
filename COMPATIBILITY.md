# Cross-repository compatibility matrix

skpm is one of four sibling repos in an agentic skill supply chain:

```
skcr (author/compile)  →  skil (scan/attest)  →  skpm (package/publish)  →  SkillForge (registry)
```

Each pairing that matters for skpm is enforced by a CI job, not just
documented — a version bump on either side that breaks the pairing fails
CI, not silently drifts. This file records *what's currently pinned* so
the pairing is auditable without reading workflow YAML.

| Pairing                                    | Enforced by                                                    | Currently pinned to | Status |
|---------------------------------------------|------------------------------------------------------------------|----------------------|--------|
| skpm `main` (current) × skil stable          | `.github/workflows/ci.yml` → `skil-interop`                     | [skil v0.2.0](https://github.com/domehahn/skil/releases/tag/v0.2.0) | ✅ enforced |
| skpm `main` (current) × SkillForge `main` (current) | `.github/workflows/ci.yml` → `skillforge-e2e`             | [SkillForge @ 5e80b37](https://github.com/domehahn/SkillForge/commit/5e80b37ce5c7d6d4d51c5ab578e2bbfb23e5600a) | ✅ enforced |
| skpm `main` (current) × SkillForge stable    | —                                                                 | — | ⏳ not yet available: SkillForge has no tagged release yet |

## What each job actually checks

- **`skil-interop`**: builds skil at the pinned tag, then runs
  `internal/attestation`'s `TestVerifyInteropsWithRealSkilBinary` — it
  shells out to that real `skil` binary's `key generate`/`attest
  --signing-key`, then verifies the result using only this repo's
  `internal/attestation.Verify` (no import of skil as a library). Proves
  skpm's independent attestation verifier actually understands what skil
  stable produces, not just what skpm's own tests assume it produces.
- **`skillforge-e2e`**: checks out skpm and a pinned SkillForge commit as
  sibling directories and runs the real cross-repo contract test
  (`tests/integration/skillforge_e2e_test.go`, tag `e2e`) against a real
  SkillForge server built from that commit.

## Bumping a pin

When the sibling repo cuts a release (or, for SkillForge until it has
tags, a commit) that skpm's interop job should track:

1. Build it locally at that version.
2. Run the exact test the relevant CI job runs against it — for
   `skil-interop`, `SKIL_BINARY=<path> go test -tags e2e
   ./internal/attestation/... -run TestVerifyInteropsWithRealSkilBinary
   -v`; for `skillforge-e2e`, point `SKILLFORGE_REPO` at a checkout of the
   candidate commit and run `go test -tags e2e ./tests/integration/... -run
   TestSkillForgeE2E -v`.
3. Update the pin and this file together, in the same commit, once it
   passes.

Don't bump a pin reflexively on every upstream release — only when
you've actually verified the pairing still holds, or deliberately want to
prove it doesn't (and fix skpm's side or file an issue upstream).
