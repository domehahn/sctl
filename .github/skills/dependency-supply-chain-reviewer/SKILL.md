---
name: "dependency-supply-chain-reviewer"
description: "Describe what this skill helps an agent do."
version: "0.1.0"
since: "2026-09-09"
last_modified: "2026-09-09"
authors:
  - "platform-engineering"
stability: "experimental"
min_platform_version:
  codex: "unknown"
  amazon-q: "unknown"
  antigravity: "unknown"
  auggie: "unknown"
  bob: "unknown"
  claude-code: "unknown"
  cline: "unknown"
  codebuddy: "unknown"
  continue: "unknown"
  costrict: "unknown"
  crush: "unknown"
  github-copilot: "unknown"
  gitlab-duo: "unknown"
  factory: "unknown"
  forgecode: "unknown"
  opencode: "unknown"
  openhands: "unknown"
  cursor: "unknown"
  roo-code: "unknown"
  kiro: "unknown"
  junie: "unknown"
  gemini-cli: "unknown"
  iflow: "unknown"
  kilocode: "unknown"
  kimi: "unknown"
  lingma: "unknown"
  pi: "unknown"
  qoder: "unknown"
  qwen: "unknown"
  windsurf: "unknown"
  ollama: "unknown"
deprecated_since:
replaces:
supersedes: []
changelog:
  - version: "0.1.0"
    date: "2026-09-09"
    change: "Initial generated production-ready SDLC / DevSecOps skill"
---
# Dependency Supply Chain Reviewer

## Purpose

Review dependencies, lockfiles, package managers, SBOMs, provenance, vulnerability reachability, license risk, transitive dependencies, update strategy, artifact integrity, and supply-chain controls.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- Dependencies, lockfiles, package registries, images, build tools, or update bots change.
- A vulnerability, license, provenance, or package integrity question must be triaged.
- Transitive dependency reachability or exploitability is unclear.
- SBOM, signing, checksums, or registry trust needs review.
- The central agent routes to dependency or supply-chain review.

## Operating model

1. Identify ecosystems, manifests, lockfiles, registries, package managers, build steps, and artifacts.
2. Compare manifest and lockfile changes for unexpected transitive movement.
3. Assess vulnerability reachability, exploit preconditions, and available fixed versions.
4. Review package provenance, signatures, checksums, registry source, and maintainer trust.
5. Recommend upgrade, pin, replace, isolate, or risk-accept actions.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Manifests, lockfiles, container base images, build plugins, and generated dependency metadata.
- Direct and transitive dependencies, vulnerable code paths, exploitability, and fixed versions.
- License, provenance, signing, checksums, SBOM, registry, and artifact integrity.
- Update automation, pinning, version ranges, vendoring, and reproducibility.
- Dependency removal, replacement, or isolation options.

## Skill-Specific Checklist

- [ ] Check manifest and lockfile consistency.
- [ ] Identify newly added, removed, upgraded, downgraded, or transitive dependencies.
- [ ] Check whether vulnerable dependency code is reachable in this application.
- [ ] Check fixed versions, breaking changes, and upgrade notes.
- [ ] Check version ranges, floating tags, branch dependencies, and unpinned plugins.
- [ ] Check package source, registry, maintainer, download URL, and namespace confusion risk.
- [ ] Check license compatibility and policy exceptions.
- [ ] Check SBOM, signatures, checksums, provenance, and reproducible build evidence.
- [ ] Check container base images and OS packages for update hygiene.
- [ ] Check dependency update automation and review gates.
- [ ] Check whether dependency removal or replacement is safer than upgrade.
- [ ] Document residual risk and owner when risk is accepted.

## Decision Rules

- If a vulnerable package is reachable and fix exists, require upgrade or compensating control.
- If a package is unpinned, registry-sourced unexpectedly, or from an unknown maintainer, flag supply-chain risk.
- If lockfile changes are unexplained, require review before merge.
- If license policy is violated, require legal or governance review.
- If artifact integrity cannot be verified for release input, block release readiness.
- If vulnerability is not reachable, document evidence instead of relying only on CVSS.

## Finding Categories

- Reachable vulnerable dependency or base image.
- Suspicious, unpinned, typosquatted, abandoned, or provenance-weak package.
- Manifest/lockfile drift or unexpected transitive dependency change.
- License, policy, or attribution violation.
- Missing SBOM, signature, checksum, or artifact integrity evidence.
- Unsafe update automation, registry trust, or reproducibility gap.

## Severity Guidance

- Critical: immediate exploitability or operational failure can expose secrets, regulated data, production safety, or release integrity.
- High: credible security, reliability, compliance, rollback, or user-impact risk requires owner action before merge or release.
- Medium: meaningful maintainability, validation, documentation, or process gap should be tracked and resolved.
- Low: advisory improvement, clarity issue, or hardening opportunity with limited immediate impact.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Dependency change summary with direct/transitive and lockfile impact.
- Vulnerability table with reachability, fixed version, exploitability, and remediation.
- Supply-chain integrity findings covering source, signatures, checksums, SBOM, and registry trust.
- License and policy findings with required owner review.
- Upgrade, pin, replace, remove, or risk-accept recommendation.
- Validation commands and residual risk.

## Acceptance Criteria

- Manifest and lockfile are consistent.
- Reachable vulnerabilities have fixes or documented compensating controls.
- Unpinned or suspicious dependencies are resolved or risk-accepted.
- License and policy risks are reviewed.
- SBOM/provenance/integrity expectations are met for release artifacts.
- Residual risk includes owner and deadline.

## Anti-Patterns

- Triage by CVSS only without reachability.
- Accepting broad version ranges for production-critical packages.
- Ignoring lockfile diffs because manifest diff is small.
- Trusting package names without checking source or maintainer.
- Treating SBOM generation as artifact integrity.
- Upgrading major versions without compatibility validation.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
