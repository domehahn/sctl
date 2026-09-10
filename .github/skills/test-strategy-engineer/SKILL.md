---
name: "test-strategy-engineer"
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
# Test Strategy Engineer

## Purpose

Design a risk-based test strategy for features, fixes, migrations, and releases across unit, integration, contract, E2E, regression, negative, security, performance, test data, mocking, CI gates, and coverage risk.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A change needs a test plan before implementation or release.
- Acceptance criteria must be mapped to concrete test types and validation evidence.
- Risky migrations, external integrations, contracts, or security-sensitive behavior need targeted coverage.
- Existing tests are flaky, slow, missing, or not aligned with the changed behavior.
- CI gates need a must-have versus optional validation strategy.

## Operating model

1. Map user-visible behavior, invariants, edge cases, and failure modes to test levels.
2. Classify test scope as unit, integration, contract, E2E, regression, negative, security, performance, or manual verification.
3. Prioritize tests by production risk, blast radius, frequency of change, and cost to execute.
4. Identify fixtures, test data, mocks, stubs, and cleanup needed for deterministic results.
5. Define CI gating order so fast deterministic checks block before expensive or optional suites.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Unit, integration, contract, E2E, regression, negative, security, and performance coverage.
- Test data, fixtures, mocking strategy, determinism, isolation, and cleanup.
- CI gates, coverage risks, flaky tests, and validation sequencing.
- Boundary cases, abuse cases, migrations, and external dependencies.
- Must-have versus optional test scope.

## Skill-Specific Checklist

- [ ] Map each acceptance criterion to at least one concrete test or explicit manual verification step.
- [ ] Identify pure logic, validators, parsers, and branching behavior that require unit tests.
- [ ] Identify database, queue, filesystem, network, cache, auth provider, or external-service boundaries that require integration tests.
- [ ] Identify public APIs, events, CLI output, schemas, or SDK contracts that require contract tests.
- [ ] Add negative tests for invalid input, authorization failure, malformed payloads, timeouts, and abuse cases.
- [ ] Add regression tests that fail on the reported bug before accepting the fix.
- [ ] Specify fixtures, seed data, factories, mocks, and cleanup needed for deterministic isolated tests.
- [ ] Identify flaky-test risks caused by time, randomness, ordering, retries, shared state, or external services.
- [ ] Define CI gates with fast required checks before slow E2E, performance, or exploratory checks.
- [ ] Separate release-blocking tests from optional confidence-building tests with rationale.
- [ ] Include migration, rollback, compatibility, and data-loss test scenarios when persistence changes.
- [ ] Call out coverage gaps that remain after the proposed test plan.

## Decision Rules

- If a changed behavior has acceptance criteria but no automated or named manual verification, classify the strategy as incomplete.
- If an external API, event, or schema changes, require contract tests or documented consumer compatibility verification.
- If data migration or rollback is involved, require forward migration, rollback, idempotency, and corrupted-input scenarios.
- If a test depends on real time, randomness, global state, network services, or shared data, require determinism controls or mark flaky risk.
- If security-sensitive behavior changes, require negative tests for authn, authz, validation, and unsafe input paths.
- If CI runtime is high, split blocking smoke coverage from scheduled exhaustive coverage instead of dropping critical tests.

## Finding Categories

- Untested acceptance criterion or user-visible behavior.
- Missing contract coverage for API, event, schema, CLI, or SDK compatibility.
- Missing negative, authorization, validation, abuse-case, or error-path coverage.
- Flaky, nondeterministic, order-dependent, or environment-dependent test design.
- Missing migration, rollback, idempotency, or data-integrity validation.
- Weak CI gate that allows high-risk changes without required validation.

## Severity Guidance

- Critical: release can corrupt data, bypass security, or break core flows with no blocking validation.
- High: key acceptance criteria, contracts, migrations, or auth paths lack required tests before merge or release.
- Medium: meaningful edge cases, fixtures, determinism, or CI sequencing gaps reduce confidence but can be tracked.
- Low: naming, organization, coverage reporting, or optional confidence checks can improve maintainability.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Test matrix mapping requirements or changed behaviors to test type, file/location, owner, and gate status.
- Must-have versus optional test list with risk rationale.
- Negative, regression, contract, migration, and rollback scenarios where applicable.
- Fixture, mock, seed data, cleanup, and determinism requirements.
- CI gate recommendation with blocking, non-blocking, scheduled, and manual checks separated.
- Residual coverage gaps and explicit release risk if tests are deferred.

## Acceptance Criteria

- Every acceptance criterion has a mapped automated test or named manual verification step.
- Changed contracts have compatibility or consumer verification coverage.
- Critical negative, auth, validation, migration, and rollback paths are covered or explicitly risk-accepted.
- Required tests are deterministic, isolated, and suitable for CI gating.
- The strategy distinguishes required release blockers from optional confidence checks.
- Residual test gaps include owner, impact, and follow-up recommendation.

## Anti-Patterns

- Counting coverage percentage without mapping tests to changed behavior and risk.
- Using only happy-path E2E tests while missing unit, contract, and negative coverage.
- Relying on live external services or wall-clock timing for blocking CI tests.
- Treating manual QA as sufficient without named scenarios and evidence.
- Skipping rollback, migration, or compatibility tests because deployment tooling exists.
- Adding broad slow tests that make CI unusable instead of targeted gates.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
