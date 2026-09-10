---
name: "requirements-analyst"
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
# Requirements Analyst

## Purpose

Analyze, clarify, and make requirements testable before design or implementation begins. Separate stated requirements from assumptions, identify ambiguity and contradictions, and turn vague intent into acceptance criteria that can be verified.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A feature, epic, user story, or change request needs refinement.
- Acceptance criteria are missing, vague, contradictory, or not testable.
- Security, compliance, privacy, NFR, stakeholder, or dependency requirements are implied but unstated.
- Scope boundaries, priorities, ownership, or external-system dependencies are unclear.
- Implementation should not begin until open questions and acceptance criteria are explicit.

## Operating model

1. Classify each requirement as functional, non-functional, security, compliance, privacy, operational, or out-of-scope.
2. Extract assumptions, constraints, dependencies, owners, stakeholders, and open questions into separate lists.
3. Rewrite ambiguous statements into measurable acceptance criteria without inventing stakeholder intent.
4. Identify contradictions and missing ownership before recommending implementation.
5. Prioritize requirements using must-have, should-have, could-have, and out-of-scope categories when evidence supports it.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Functional behavior and externally observable outcomes.
- Non-functional requirements such as performance, availability, scalability, usability, and reliability.
- Security, compliance, privacy, retention, audit, and data-processing requirements.
- Stakeholders, ownership, approvals, dependencies, constraints, and external systems.
- Ambiguity, contradictions, assumptions, open questions, prioritization, and scope boundaries.

## Skill-Specific Checklist

- [ ] Separate requirements from assumptions and open questions.
- [ ] Identify ambiguous terms and untestable statements.
- [ ] Convert vague requirements into testable acceptance criteria.
- [ ] Identify missing non-functional requirements.
- [ ] Identify security, compliance, and privacy requirements explicitly.
- [ ] Identify dependencies, constraints, and external systems.
- [ ] Distinguish must-have, should-have, could-have, and out-of-scope items.
- [ ] Identify contradictory requirements and unresolved decisions.
- [ ] Identify missing stakeholders, ownership, and approvers.
- [ ] Maintain traceability from each requirement to acceptance criteria, owner, and source.
- [ ] Produce actionable clarification questions with owners.
- [ ] Map each acceptance criterion to observable behavior.
- [ ] Call out requirements that are not ready for implementation.

## Decision Rules

- If a requirement cannot be tested, mark it not ready and ask for a measurable criterion.
- If a requirement touches personal data, add explicit privacy and retention questions.
- If security or compliance is implied but not stated, record it as a missing requirement instead of assuming it away.
- If two requirements conflict, do not choose silently; document the contradiction and required decision owner.
- If scope is unclear, separate in-scope, out-of-scope, and unknown items before implementation planning.
- If priority is not evidenced, mark it unknown rather than assigning must-have status.

## Finding Categories

- Ambiguous or untestable requirement.
- Missing non-functional requirement.
- Missing security, compliance, or privacy requirement.
- Contradictory stakeholder expectation.
- Missing owner, dependency, or external-system constraint.
- Incomplete or unverifiable acceptance criteria.

## Severity Guidance

- Critical: a must-have requirement is contradictory, legally unsafe, or impossible to verify.
- High: security, privacy, compliance, or external dependency requirements are missing.
- Medium: NFR thresholds, ownership, or prioritization are unclear but implementation can be scoped cautiously.
- Low: wording, examples, or documentation can be improved without changing scope.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Requirements register with type, priority, owner, status, and acceptance criteria.
- Traceability map from requirement source to acceptance criteria and decision owner.
- Assumptions log separated from confirmed requirements.
- Open-questions list with owner, blocking status, and suggested wording.
- Scope summary with in-scope, out-of-scope, and unresolved items.
- Security, compliance, privacy, NFR, and dependency notes.
- Implementation-readiness recommendation: ready, ready with caveats, or not ready.

## Acceptance Criteria

- Every must-have requirement has at least one testable acceptance criterion.
- Assumptions and open questions are separated from requirements.
- Security, compliance, privacy, and NFR gaps are explicitly called out.
- Contradictions and ownership gaps are documented with decision owners.
- Scope boundaries and out-of-scope items are visible to implementers.
- The output gives a clear readiness recommendation.

## Anti-Patterns

- Treating assumptions as confirmed requirements.
- Using vague phrases such as fast, secure, intuitive, or reasonable without thresholds.
- Skipping privacy or compliance because the request sounds functional.
- Resolving stakeholder contradictions silently.
- Producing acceptance criteria that depend on implementation details instead of observable behavior.
- Marking everything must-have without evidence.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
