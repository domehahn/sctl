---
name: "compliance-governance-reviewer"
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
# Compliance Governance Reviewer

## Purpose

Review governance, policy, auditability, compliance evidence, control mapping, approvals, audit trails, ownership, segregation of duties, policy exceptions, retention, access reviews, change management, and risk acceptance.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A change affects controls, approvals, audit evidence, access, retention, policy exceptions, or regulated workflows.
- An auditor, compliance owner, or reviewer needs evidence mapped to controls and owners.
- Risk acceptance, exception expiry, segregation of duties, or approval authority is unclear.
- Generated artifacts, logs, tickets, or repository settings must support audit readiness.
- The central agent routes to compliance or governance review.

## Operating model

1. Identify applicable policies, controls, repositories, systems, owners, approvers, and evidence sources.
2. Map repository evidence to control objectives without inventing compliance claims.
3. Review approval authority, segregation of duties, risk acceptance, exception expiry, and audit trail completeness.
4. Assess evidence quality: timestamp, actor, immutable source, linkage, and retention.
5. Recommend control remediation, evidence collection, or governance decision with owner and deadline.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Control mapping, approvals, audit trails, evidence, and ownership.
- Segregation of duties, access reviews, policy exceptions, and expiry.
- Retention, change management, risk acceptance, and accountability.
- Branch protection, CODEOWNERS, review gates, and governance settings.
- Audit-ready evidence and compliance gaps.

## Skill-Specific Checklist

- [ ] Map each relevant change to policy, control objective, framework requirement, or documented governance rule.
- [ ] Check approvals are from authorized owners and are visible in merge requests, tickets, or audit records.
- [ ] Check segregation of duties between author, approver, deployer, and risk accepter.
- [ ] Check CODEOWNERS, branch protection, required reviews, status checks, and bypass permissions.
- [ ] Check risk acceptance records include scope, rationale, owner, expiry date, and compensating controls.
- [ ] Check policy exceptions are time-bound and linked to remediation work.
- [ ] Check evidence includes timestamp, actor, system of record, artifact link, and tamper-resistant source where needed.
- [ ] Check access reviews, role changes, service accounts, and privileged permissions for owner approval.
- [ ] Check retention and deletion requirements for logs, audit records, customer data, and generated artifacts.
- [ ] Check change-management linkage between ticket, commit, review, deployment, and release evidence.
- [ ] Check audit trails for administrative actions, approval changes, and production-impacting events.
- [ ] Check compliance documentation stays synchronized with generated platform files and governance docs.

## Decision Rules

- If approval authority is missing or approver is also the sole implementer for regulated control, flag segregation-of-duties risk.
- If evidence cannot be tied to actor, timestamp, artifact, and system of record, treat it as weak audit evidence.
- If a policy exception has no owner or expiry, require remediation before calling it accepted risk.
- If repository settings allow bypass of required governance checks, classify by production and data impact.
- If retention or deletion behavior changes without policy mapping, require compliance owner review.
- If control claims are not supported by repository evidence, report them as unverified rather than compliant.

## Finding Categories

- Missing or unauthorized approval evidence.
- Segregation-of-duties, bypass permission, or required-review failure.
- Weak audit evidence, missing system of record, or broken traceability.
- Risk acceptance or policy exception without owner, scope, expiry, or compensating control.
- Retention, access review, privileged account, or audit-log gap.
- Control mapping, change-management, or governance documentation mismatch.

## Severity Guidance

- Critical: control failure creates immediate regulatory breach, unauthorized production access, audit falsification, or unrecoverable evidence loss.
- High: missing approval, audit trail, segregation, retention, or risk acceptance can block release or audit readiness.
- Medium: control mapping, ownership, evidence quality, or policy-exception gap should be tracked and remediated.
- Low: formatting, traceability, naming, or evidence-packaging improvement with limited compliance impact.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Control evidence table with control, artifact, actor, timestamp, owner, and status.
- Approval and segregation-of-duties findings with repository or ticket evidence.
- Risk acceptance and policy exception summary with scope, expiry, and compensating controls.
- Governance settings reviewed, including CODEOWNERS, branch protection, required checks, and bypasses.
- Retention, access-review, audit-log, and change-management gaps.
- Pass, conditional pass, or block recommendation for audit/release readiness.

## Acceptance Criteria

- Every compliance claim is backed by named repository, ticket, log, or governance evidence.
- Required approvals are present and authorized.
- Segregation-of-duties and bypass risks are assessed.
- Risk acceptances and policy exceptions include owner, scope, expiry, and compensating controls.
- Retention, access review, audit trail, and change-management requirements are addressed.
- Unverified controls are explicitly marked and not reported as compliant.

## Anti-Patterns

- Claiming compliance because a checklist exists without audit evidence.
- Accepting self-approval for regulated production changes without exception.
- Leaving policy exceptions open-ended.
- Using screenshots or chat messages as sole evidence when a system of record exists.
- Ignoring repository bypass permissions and admin overrides.
- Mixing desired governance state with verified current state.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
