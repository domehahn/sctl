---
name: "cncf-novo-nordisk-member-reviewer"
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
# Cncf Novo Nordisk Member Reviewer

## Purpose

Review Novo Nordisk (member) as an organization listed in the CNCF Landscape using current first-party evidence and explicit procurement, governance, security, portability, and operational criteria.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- The organization, its cloud-native products, support, partnership, procurement, outsourcing, or exit strategy is being assessed.
- Claims about CNCF membership or capabilities need current verification.

## Operating model

1. Treat CNCF Landscape inclusion as classification metadata, not endorsement, certification, security assurance, or proof of product fitness.
2. Verify current claims against the organization's first-party material and the official CNCF Landscape entry.
3. Separate organization-level membership from the maturity or CNCF status of individual projects.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- CNCF classification: CNCF Members / Silver
- First-party source: https://www.novonordisk.com
- product scope, ownership, support, security, compliance, portability, commercial dependencies, concentration risk, and exit strategy

## Skill-Specific Checklist

- [ ] Confirm legal and product identity, current CNCF classification, offered services, support boundaries, regions, data handling, and shared responsibilities.
- [ ] Check security documentation, incident handling, vulnerability disclosure, support SLAs, audit evidence, subcontractors, portability, and termination assistance.
- [ ] Identify proprietary dependencies, lock-in, pricing or licensing constraints, migration paths, and operational ownership.
- [ ] Trace Novo Nordisk (member) inputs, outputs, identities, trust boundaries, external dependencies, and persistent state.
- [ ] Check Novo Nordisk (member) defaults, configuration precedence, environment separation, and drift from reviewed source.
- [ ] Check Novo Nordisk (member) least privilege, credential rotation, audit events, policy enforcement, and break-glass behavior.
- [ ] Check Novo Nordisk (member) resource ownership, cleanup, quotas, rate limits, timeouts, retries, and backpressure.
- [ ] Check Novo Nordisk (member) release notes, supported upgrade paths, schema or API compatibility, and rollback constraints.
- [ ] Check Novo Nordisk (member) dashboards, alerts, health signals, incident procedures, backup evidence, and recovery exercises.
- [ ] Check Novo Nordisk (member) license, maintenance state, vulnerability handling, provenance, signatures, and dependency pinning.

## Decision Rules

- Do not infer CNCF project status, certification, security, or endorsement from membership.
- Do not make procurement or legal conclusions without the applicable current contract and accountable reviewers.
- Require an exit path and evidence for material capability or compliance claims.
- If Novo Nordisk (member) identity, version, edition, or deployment model cannot be established, report the uncertainty before recommending a change.
- If a Novo Nordisk (member) change can affect availability, security boundaries, persistent data, or external consumers, require staged validation and explicit rollback criteria.

## Finding Categories

- Identity, scope, classification, or claim mismatch.
- Security, compliance, support, responsibility, or incident-management gap.
- Portability, concentration, licensing, cost, contract, or exit-strategy risk.
- Novo Nordisk (member) documentation, ownership, maintenance, licensing, provenance, or evidence gap.
- Novo Nordisk (member) integration, dependency, compatibility, migration, or decommissioning risk.

## Severity Guidance

- Critical: Novo Nordisk (member) can enable broad compromise, destructive production mutation, tenant escape, secret exposure, or material data loss.
- High: Novo Nordisk (member) can cause exploitable access, sustained outage, corrupt state, unsafe upgrade, or loss of recovery capability.
- Medium: Novo Nordisk (member) has a bounded correctness, reliability, compatibility, observability, maintenance, or performance risk.
- Low: Novo Nordisk (member) needs documentation, classification, idiomatic configuration, or validation-evidence improvement.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- State the Landscape snapshot classification and sources checked.
- Return an evidence-backed capability and risk assessment with open questions, owners, and expiry dates.
- Include a Novo Nordisk (member) architecture and dependency summary with trust boundaries and persistent state.
- Include Novo Nordisk (member) validation commands or tests that are safe for the stated environment.
- Include Novo Nordisk (member) ownership, rollout, rollback, monitoring, and follow-up actions.

## Acceptance Criteria

- Membership is not presented as endorsement or project maturity.
- Material claims, responsibilities, dependencies, and exit assumptions are traceable to current evidence.
- Novo Nordisk (member) identity, version, configuration, dependencies, and deployment assumptions are explicit.
- Novo Nordisk (member) failure, upgrade, rollback, observability, and recovery paths have proportionate evidence.
- Novo Nordisk (member) security boundaries, credentials, network exposure, data handling, and supply-chain risks are addressed.

## Anti-Patterns

- Equating a Landscape logo with CNCF certification or technical approval.
- Inventing product features, contract terms, regions, certifications, or support commitments.
- Assuming Landscape inclusion or popularity proves that Novo Nordisk (member) is secure, supported, compatible, or suitable.
- Using an unpinned latest version of Novo Nordisk (member) as the basis for migration or production guidance.
- Ignoring Novo Nordisk (member) operational ownership, data lifecycle, upgrade constraints, rollback, and decommissioning.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
