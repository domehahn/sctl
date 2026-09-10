---
name: "observability-reviewer"
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
# Observability Reviewer

## Purpose

Review logging, metrics, tracing, alerting, dashboards, operational visibility, SLOs, SLIs, runbooks, audit logs, sensitive data in logs, on-call usability, and incident detection.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A service, feature, or deployment changes behavior that operators must detect, debug, or support.
- Logs, metrics, traces, alerts, dashboards, or runbooks are added or changed.
- SLOs, SLIs, audit logs, or incident detection coverage is unclear.
- Sensitive data may enter logs or telemetry.
- The central agent routes to observability review.

## Operating model

1. Map critical user journeys, failure modes, dependencies, and operational questions.
2. Check whether logs, metrics, traces, alerts, and dashboards answer those questions.
3. Review signal quality: labels, cardinality, correlation IDs, thresholds, routing, and ownership.
4. Assess privacy and security of telemetry.
5. Recommend concrete telemetry, alert, dashboard, or runbook changes.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Logs, metrics, traces, alerts, dashboards, and runbooks.
- SLOs, SLIs, critical journeys, audit logs, and incident detection.
- Sensitive data in logs, correlation IDs, retention, and cardinality.
- On-call usability, alert routing, alert fatigue, and ownership.
- Operational decision support and recovery verification.

## Skill-Specific Checklist

- [ ] Identify critical user journeys and failure modes requiring visibility.
- [ ] Check error, latency, throughput, saturation, dependency, and queue metrics.
- [ ] Check logs include event names, correlation IDs, actor/resource IDs, and safe context.
- [ ] Check traces connect ingress, service calls, database, queues, and external dependencies.
- [ ] Check alerts have actionable thresholds, severity, runbook, owner, and routing.
- [ ] Check dashboards answer deploy health, customer impact, dependency health, and rollback decisions.
- [ ] Check sensitive data, secrets, tokens, PII, and payloads are excluded or redacted.
- [ ] Check metric label cardinality, retention, cost, and aggregation safety.
- [ ] Check SLO/SLI coverage for critical journeys.
- [ ] Check audit logs for security-relevant actions and tamper-resistant retention.
- [ ] Check runbooks include diagnosis, mitigation, rollback, and verification steps.
- [ ] Identify noisy, duplicate, missing, or unactionable alerts.

## Decision Rules

- If operators cannot detect failure of a critical journey, require metrics or alerts before release.
- If logs can expose secrets or PII, require redaction before approval.
- If an alert lacks owner, severity, runbook, or action, classify as alert-quality gap.
- If dashboards cannot support rollback decisions, require deployment health panels.
- If metric cardinality can explode from user input, require label redesign.
- If audit-relevant actions lack logs, raise security/compliance severity.

## Finding Categories

- Missing signal for critical user journey or dependency failure.
- Unsafe logging of secrets, PII, tokens, or sensitive payloads.
- Unactionable, noisy, duplicate, or ownerless alert.
- Dashboard gap for deploy health, customer impact, or rollback decision.
- Trace, correlation ID, or context propagation gap.
- SLO/SLI, audit log, retention, or cardinality risk.

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

- Observability coverage map for journeys, failure modes, signals, alerts, dashboards, and runbooks.
- Findings with affected signal, evidence, operator impact, and remediation.
- Concrete metric, log, trace, alert, dashboard, or runbook recommendation.
- Sensitive telemetry and cardinality risk assessment.
- Post-deploy monitoring and incident-detection recommendation.
- Residual blind spots and owner actions.

## Acceptance Criteria

- Critical journeys have metrics, logs/traces, alerts, and dashboard support.
- Telemetry excludes secrets and sensitive payloads.
- Alerts are actionable, owned, routed, and tied to runbooks.
- Dashboards support diagnosis and rollback decisions.
- SLOs/SLIs or audit logs exist where required.
- Known blind spots are explicitly documented.

## Anti-Patterns

- Adding logs without deciding what operator question they answer.
- Alerting on every error without severity, owner, or action.
- Using high-cardinality user input as metric labels.
- Logging full payloads to debug production issues.
- Building dashboards that cannot support rollback or incident triage.
- Assuming tracing solves missing metrics or alerting.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
