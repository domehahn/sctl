---
name: "release-readiness-reviewer"
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
# Release Readiness Reviewer

## Purpose

Determine whether a change or system is ready for release by reviewing test status, security findings, known issues, rollback, migrations, monitoring, alerts, runbooks, feature flags, approvals, release notes, support readiness, and go/no-go recommendation.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A release, deployment, version bump, migration, or production rollout needs go/no-go review.
- Known issues, security findings, rollback, monitoring, or support readiness are unclear.
- Feature flags, staged rollout, or compatibility risk must be assessed.
- Release notes, approvals, or operational ownership need validation.
- The central agent routes to release readiness review.

## Operating model

1. Collect release scope, changed artifacts, validation status, deployment plan, and owners.
2. Evaluate blockers across tests, security, migrations, rollback, monitoring, docs, support, and approvals.
3. Separate go/no-go criteria from follow-up work and known accepted risks.
4. Assess rollout strategy, feature flags, blast radius, and recovery time.
5. Return go, conditional go, or no-go with explicit blockers.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Test status, validation evidence, security findings, and known issues.
- Rollback, migrations, compatibility, feature flags, and staged rollout.
- Monitoring, alerts, runbooks, support readiness, and ownership.
- Approvals, release notes, communication, and go/no-go decision.
- Blockers, exceptions, residual risk, and follow-up actions.

## Skill-Specific Checklist

- [ ] Verify required tests, builds, scans, migrations, and smoke checks are complete.
- [ ] Identify unresolved blockers and classify known issues by user impact.
- [ ] Check open security findings, exceptions, owners, and expiry.
- [ ] Verify rollback plan, rollback trigger, rollback owner, and data rollback constraints.
- [ ] Check migration forward/backward compatibility, idempotency, and backup plan.
- [ ] Verify monitoring dashboards, alerts, SLO indicators, and post-deploy validation.
- [ ] Check runbooks, escalation contacts, on-call coverage, and support readiness.
- [ ] Check release notes, changelog, customer communication, and breaking-change guidance.
- [ ] Identify feature flag, canary, staged rollout, kill switch, or traffic-shaping options.
- [ ] Confirm approvals, release owner, deployment window, and freeze constraints.
- [ ] Check dependencies on external services, infra capacity, and version compatibility.
- [ ] Produce go/no-go with explicit conditions.

## Decision Rules

- If rollback is impossible or untested for high-impact change, no-go unless risk is accepted by owner.
- If Critical/High security findings are open without approved exception, no-go.
- If migrations can corrupt or lose data without backup and validation, no-go.
- If monitoring cannot detect release failure, require conditional go or no-go by impact.
- If known issues affect core user journeys, require mitigation, communication, or staged rollout.
- If release notes omit breaking changes or migrations, block external release readiness.

## Finding Categories

- Missing go/no-go evidence or owner.
- Failed, skipped, stale, or insufficient validation gate.
- Open security, compliance, privacy, or known-issue blocker.
- Rollback, migration, compatibility, or data-safety gap.
- Monitoring, alerting, runbook, on-call, or support readiness gap.
- Release notes, communication, approval, or change-management gap.

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

- Go/no-go summary with blockers, conditions, and owner.
- Release checklist covering tests, security, migrations, rollback, monitoring, docs, support, and approvals.
- Known issues table with severity, impact, mitigation, and acceptance owner.
- Rollback and post-deploy validation plan.
- Feature flag, canary, staged rollout, or kill-switch recommendation.
- Residual risks and follow-up actions with deadlines.

## Acceptance Criteria

- All release blockers are resolved or explicitly accepted.
- Rollback, migration, monitoring, and support readiness are verified.
- Security findings and known issues have owner-approved disposition.
- Release notes cover breaking changes, migrations, and operational impact.
- Go/no-go recommendation follows evidence.
- Post-release validation and escalation path are defined.

## Anti-Patterns

- Treating green CI as complete release readiness.
- Approving release with no rollback trigger or owner.
- Ignoring known issues because they are documented elsewhere.
- Skipping monitoring and support readiness until after deployment.
- Accepting security findings without expiry and owner.
- Publishing breaking changes without migration guidance.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
