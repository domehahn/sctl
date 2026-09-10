---
name: "incident-postmortem-assistant"
description: "Support incident analysis, timeline creation, root cause analysis, impact assessment, corrective actions, and follow-up issues."
version: "1.0.0"
since: "2026-09-09"
last_modified: "2026-09-09"
authors:
  - "platform-engineering"
stability: "stable"
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
  - version: "1.0.0"
    date: "2026-09-09"
    change: "Initial generated production-ready SDLC / DevSecOps skill"
---
# Incident Postmortem Assistant

## Purpose

Support incident response, postmortems, and corrective actions across triage, severity, impact, timeline, containment, eradication, recovery, communication, evidence preservation, root cause, contributing factors, corrective actions, and prevention.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- An active incident needs structured triage, timeline, impact, containment, or communication support.
- A postmortem needs facts, contributing factors, root cause, and corrective actions.
- Follow-up actions must be owner-assigned and verifiable.
- Stakeholder communication must separate facts from assumptions.
- The central agent routes to incident or postmortem assistance.

## Operating model

1. Separate confirmed facts, assumptions, hypotheses, unknowns, and decisions.
2. Build timeline from alerts, logs, deploys, tickets, chats, and customer impact.
3. Classify severity, affected services, customer/business/security impact, and current state.
4. Guide containment, recovery, and evidence preservation without destructive shortcuts.
5. Produce blameless postmortem and corrective actions with owners and due dates.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Triage, severity, impact, timeline, containment, eradication, and recovery.
- Communication, evidence preservation, facts, assumptions, and unknowns.
- Root cause, contributing factors, corrective actions, prevention, and owners.
- Security, customer, business, compliance, and operational impact.
- Postmortem readiness and follow-up issue quality.

## Skill-Specific Checklist

- [ ] Separate confirmed facts, assumptions, hypotheses, and unknowns.
- [ ] Establish timeline with timestamps, sources, and confidence level.
- [ ] Identify customer, business, security, compliance, and operational impact.
- [ ] Classify severity and document severity changes over time.
- [ ] Recommend containment steps that preserve evidence and avoid unsafe recovery.
- [ ] Track detection, mitigation, recovery, and resolution times.
- [ ] Identify immediate cause, contributing factors, and systemic causes.
- [ ] Distinguish mitigation, eradication, recovery, prevention, and follow-up work.
- [ ] Draft stakeholder communication using confirmed facts only.
- [ ] Assign corrective actions with owner, deadline, and verification criterion.
- [ ] Identify monitoring, runbook, test, process, or architecture gaps that allowed recurrence.
- [ ] Verify recovery using telemetry, customer impact, and service health checks.

## Decision Rules

- If facts are incomplete, mark them unknown instead of inventing root cause.
- If incident may involve security compromise, preserve evidence before cleanup.
- If containment may worsen impact, state trade-offs and seek owner decision.
- If corrective action lacks owner, due date, and verification, it is not postmortem-ready.
- If customer impact is unknown, require impact assessment before final severity.
- If communication includes hypotheses, label them clearly or remove them from external messaging.

## Finding Categories

- Incomplete timeline, missing fact source, or assumption presented as fact.
- Unclear severity, impact, affected service, or customer scope.
- Unsafe containment, recovery, or evidence-destroying action.
- Weak root cause analysis or missing contributing factors.
- Corrective action without owner, deadline, or verification criterion.
- Missing communication, monitoring, runbook, test, or prevention follow-up.

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

- Incident summary with severity, impact, affected systems, status, and confidence level.
- Timeline with timestamped events, sources, and gaps.
- Facts, assumptions, hypotheses, unknowns, and decisions separated.
- Containment, recovery, and evidence-preservation recommendations.
- Root cause and contributing factors with supporting evidence.
- Corrective action table with owner, due date, verification, and priority.

## Acceptance Criteria

- Timeline covers detection, mitigation, recovery, and resolution with sources.
- Impact and severity are justified and updated when evidence changes.
- Root cause analysis distinguishes immediate cause from contributing factors.
- Corrective actions are specific, owned, dated, and verifiable.
- Communications avoid speculation and expose open questions.
- Recovery is verified by telemetry or customer-impact evidence.

## Anti-Patterns

- Declaring root cause before facts support it.
- Blaming individuals instead of analyzing systems and contributing factors.
- Deleting logs or changing systems before preserving evidence.
- Writing corrective actions like “be more careful”.
- Omitting customer impact because service health recovered.
- Closing postmortem without verifying follow-up completion criteria.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
