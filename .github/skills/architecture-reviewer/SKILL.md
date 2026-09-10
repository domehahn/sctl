---
name: "architecture-reviewer"
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
# Architecture Reviewer

## Purpose

Review architecture for module boundaries, service boundaries, coupling, cohesion, dependency direction, data ownership, resilience, security boundaries, and ADR-worthy decisions.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A design, PR, or refactor changes architecture, module boundaries, or service boundaries.
- API contracts, shared libraries, data ownership, or deployment topology change.
- Scalability, resilience, runtime coupling, or security boundaries need review.
- Circular dependencies, layering violations, or unclear ownership are suspected.
- An ADR should be created or updated.

## Operating model

1. Map components, modules, services, APIs, data stores, queues, and deployment units.
2. Verify dependency direction, layering, cohesion, and ownership boundaries.
3. Trace data flows and ownership across tables, topics, buckets, and integrations.
4. Review runtime coupling, fan-out, retry behavior, scalability, and cascading failure risk.
5. Identify decisions that deserve ADRs and distinguish architecture risk from style preference.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Module boundaries, service boundaries, layering, coupling, cohesion, and circular dependencies.
- API contracts, interface stability, compatibility, and versioning.
- Data ownership, data flows, cross-boundary writes, and direct database access.
- Runtime and deployment coupling, scalability, resilience, retry storms, and fan-out.
- Security boundaries, trust boundaries, ownership, and ADR candidates.

## Skill-Specific Checklist

- [ ] Identify architectural layers and verify dependency direction.
- [ ] Detect circular dependencies between modules, packages, services, or libraries.
- [ ] Check whether module boundaries align with business capabilities or clear technical responsibilities.
- [ ] Review coupling between services, APIs, databases, queues, shared libraries, and deployment units.
- [ ] Identify shared mutable state, shared database writes, hidden dependencies, and temporal coupling.
- [ ] Check public interfaces for ownership, versioning, compatibility expectations, and tests.
- [ ] Verify explicit data ownership for domain objects, tables, topics, buckets, and integrations.
- [ ] Identify cross-boundary writes or direct database access across service boundaries.
- [ ] Review synchronous call chains for fan-out, latency amplification, retry storms, and cascading failures.
- [ ] Identify ADR-worthy decisions and missing architecture documentation.
- [ ] Check security and trust boundaries between components.
- [ ] Assess deployment coupling and independent rollback ability.

## Decision Rules

- If a dependency violates the intended layer direction, classify it as architecture risk.
- If circular dependencies exist, recommend interface inversion or boundary redesign.
- If data ownership is unclear, block cross-boundary writes until an owner is named.
- If synchronous fan-out can cascade failures, recommend async, timeout, circuit breaker, or fallback design.
- If a public contract lacks versioning or compatibility tests, flag interface stability risk.
- If a decision changes long-term structure, recommend an ADR.

## Finding Categories

- Circular dependencies and dependency direction violations.
- Excessive coupling or weak cohesion.
- Unclear module, service, or data ownership.
- Unstable API contracts and compatibility risk.
- Runtime coupling, fan-out, retry storm, and cascading failure risk.
- Missing ADR or architecture documentation.

## Severity Guidance

- Critical: architecture enables unsafe deployment, data corruption, or unavoidable cascading failure.
- High: coupling, ownership, or boundary issue blocks safe evolution or security isolation.
- Medium: maintainability, scalability, or compatibility risk is likely but controllable.
- Low: documentation or ADR gap that does not currently block delivery.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Architecture summary with components, boundaries, and data flows.
- Findings table with severity, component, evidence, impact, and recommendation.
- Dependency and coupling analysis with circular dependencies called out.
- API contract and data ownership review.
- Runtime resilience, scalability, and deployment coupling notes.
- ADR recommendations with decision question and rationale.

## Acceptance Criteria

- Module and service boundaries are identified or marked unknown.
- Circular dependencies and layering violations are explicitly assessed.
- Data ownership and cross-boundary access are reviewed.
- Critical and High findings include concrete mitigation.
- ADR candidates are listed for broad decisions.
- Findings distinguish structural risk from style preference.

## Anti-Patterns

- Reporting style preferences as architecture risk.
- Recommending service extraction without ownership and operational cost analysis.
- Ignoring data ownership while reviewing module boundaries.
- Approving circular dependencies because they compile today.
- Treating ADRs as optional for precedent-setting decisions.
- Suggesting broad rewrites for localized coupling issues.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
