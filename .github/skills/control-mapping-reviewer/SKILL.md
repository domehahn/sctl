---
name: "control-mapping-reviewer"
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
# Control Mapping Reviewer

## Purpose

Map technical measures to DORA, VAIT or BAIT migration needs, ISO 27001, BSI, internal policies, or MaRisk review expectations. Treat regulatory, security, and operational references as review and evidence guidance, not legal advice.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- control mapping decisions, controls, or operating practices need independent review.
- A change affects control mapping artifacts such as technical control, DORA chapter map, VAIT or BAIT mapping, ISO 27001 annex, BSI baseline, internal policy.
- The user needs evidence-oriented findings for risks such as unmapped technical measure, duplicated control claim, obsolete VAIT reference, weak MaRisk linkage, missing implementation proof, policy conflict.
- Audit, security, operations, or platform stakeholders need a concise readiness position.
- Existing documentation, tickets, tests, or logs must be turned into actionable remediation items.

## Operating model

1. Identify the relevant control mapping artifacts, owners, systems, environments, and review boundary.
2. Compare the available artifacts against expected signals such as control owner, implementation ticket, test evidence, policy clause, audit note, exception record.
3. Separate confirmed gaps from assumptions, missing evidence, and advisory improvement opportunities.
4. Rate findings by operational, security, compliance, customer, and auditability impact.
5. Recommend minimal remediation steps, validation evidence, owners, and review cadence.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Primary artifacts: technical control, DORA chapter map, VAIT or BAIT mapping, ISO 27001 annex, BSI baseline, internal policy.
- Risk themes: unmapped technical measure, duplicated control claim, obsolete VAIT reference, weak MaRisk linkage, missing implementation proof, policy conflict.
- Evidence signals: control owner, implementation ticket, test evidence, policy clause, audit note, exception record.
- Ownership, approvals, review cadence, exception handling, and residual-risk decisions.
- Traceability from requirement or control intent to implementation, validation, and retained evidence.

## Skill-Specific Checklist

- [ ] Confirm the review boundary covers the right control mapping systems, teams, and environments.
- [ ] Inventory and inspect the current technical control.
- [ ] Check whether DORA chapter map is current, approved, versioned, and owned.
- [ ] Verify that VAIT or BAIT mapping has test, ticket, log, or approval support.
- [ ] Look for unmapped technical measure and record concrete repository or process evidence.
- [ ] Look for duplicated control claim and identify affected assets, services, or stakeholders.
- [ ] Look for obsolete VAIT reference and classify the operational or audit impact.
- [ ] Use control owner to validate that the control or practice is operating.
- [ ] Use implementation ticket to confirm ownership, timing, and reproducibility.
- [ ] Check exception, risk-acceptance, and expiry handling for control mapping.
- [ ] Confirm remediation items have owners, due dates, validation steps, and evidence expectations.
- [ ] Identify missing artifacts separately from weak artifacts so the next action is unambiguous.
- [ ] Review whether logging, reporting, or retained evidence exposes sensitive data unnecessarily.

## Decision Rules

- If technical control is missing for a critical service, raise at least a high-severity readiness gap.
- If implementation ticket cannot be tied to an owner and approval, treat the outcome as unauditable until corrected.
- If ISO 27001 annex is present but expired or untested, require validation before accepting residual risk.
- If the only support is verbal or chat-only context, request durable ticket, document, log, or test evidence.
- If remediation would require a process or architecture decision, assign a decision owner instead of prescribing legal conclusions.
- If compensating measures reduce likelihood but not impact, keep the residual-risk statement explicit.

## Finding Categories

- Missing or stale control mapping artifact.
- Unclear ownership, approval, review cadence, or accountability.
- Insufficient validation, test proof, logs, ticket trail, or retained audit material.
- Unreviewed exception, residual risk, expiry, or compensating measure.
- Policy, architecture, operational, or platform implementation drift.
- Sensitive-data exposure in logs, reports, prompts, artifacts, or evidence packages.

## Severity Guidance

- Critical: a gap in control mapping creates immediate outage, data-loss, privilege, regulatory-reporting, or irreversible business risk.
- High: technical control is missing, unowned, untested, or unauditable for a critical service or material change.
- Medium: DORA chapter map exists but is stale, incomplete, inconsistently enforced, or weakly evidenced.
- Low: wording, metadata, formatting, link freshness, or minor traceability improvements are needed.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Findings ordered by severity with affected control mapping artifacts and evidence references.
- Coverage note for reviewed artifacts: technical control, DORA chapter map, VAIT or BAIT mapping, ISO 27001 annex, BSI baseline, internal policy.
- Risk note covering relevant themes: unmapped technical measure, duplicated control claim, obsolete VAIT reference, weak MaRisk linkage, missing implementation proof, policy conflict.
- Evidence request list using expected signals: control owner, implementation ticket, test evidence, policy clause, audit note, exception record.
- Deliverables or updates needed: control mapping matrix, framework coverage gap list, migration notes, duplicate-control cleanup list, evidence alignment report.
- Residual-risk, assumptions, missing-context, and validation-gap summary.

## Acceptance Criteria

- Relevant control mapping artifacts are identified, current, owned, and versioned where applicable.
- Each high-impact finding includes evidence, impact, likelihood, owner, and remediation guidance.
- Missing evidence is separated from failed controls or weak implementation.
- Exceptions and risk acceptances include owner, rationale, expiry, and compensating measures.
- Recommendations are review-oriented and avoid presenting regulatory interpretation as legal advice.
- Final output states pass, conditional pass, or blocked readiness with validation gaps.

## Anti-Patterns

- Treating a policy title or control name as proof that the practice operates effectively.
- Collapsing missing evidence and failed implementation into one vague finding.
- Accepting open-ended exceptions without owner, expiry, impact, likelihood, and compensating measures.
- Making legal, regulatory, or audit conclusions beyond the available evidence and review scope.
- Recommending broad process rewrites when a targeted owner, test, ticket, or evidence fix is enough.
- Copying sensitive production data into examples, evidence packages, prompts, or reports.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
