---
name: "verification-reviewer"
description: "Review diffs, validate acceptance criteria, inspect test results, and find missed requirements."
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
# Verification Reviewer

## Purpose

Verify implementation results against requirements, changed files, validation commands, test evidence, edge cases, security validation, documentation consistency, and residual risk.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A change claims completion and needs independent evidence-based verification.
- Validation output, test coverage, or acceptance criteria need review before merge or release.
- The implementation may have scope creep, missed edge cases, or generated-output drift.
- Security, docs, migration, rollback, or compatibility claims need evidence.
- A final pass, conditional pass, or block recommendation is required.

## Operating model

1. Map requested requirements and acceptance criteria to changed files and observable behavior.
2. Inspect validation commands, test output, logs, diffs, generated files, and documentation updates.
3. Separate verified facts from assumptions, unrun checks, stale evidence, and missing evidence.
4. Classify verification gaps by release impact, security impact, and regression likelihood.
5. Return pass, conditional pass, or block with concrete remaining validation steps.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Requirement-to-implementation traceability.
- Acceptance criteria and test evidence credibility.
- Changed files, edge cases, scope expansion, and regression risk.
- Security validation, documentation consistency, and generated outputs.
- Residual risk and pass/fail recommendation.

## Skill-Specific Checklist

- [ ] Map every requested requirement to changed files, tests, docs, or explicit non-code evidence.
- [ ] Verify acceptance criteria against actual implementation behavior, not only author claims.
- [ ] Check validation commands were run in the relevant repo state and report failures honestly.
- [ ] Identify untested edge cases, negative paths, concurrency paths, migration paths, and rollback paths.
- [ ] Check changed files for unrelated edits, broad refactors, generated-output drift, or hidden scope expansion.
- [ ] Verify documentation, changelog, API docs, runbooks, or examples changed when behavior changed.
- [ ] Verify security-sensitive changes include auth, validation, logging, and secrets checks where relevant.
- [ ] Distinguish verified facts, inferred facts, assumptions, missing evidence, and stale evidence.
- [ ] Check generated files, lockfiles, schema files, and platform copies are synchronized with source files.
- [ ] Assess residual risks and whether they block merge, block release, or can be tracked.
- [ ] Review failing or skipped tests for relevance before accepting pass status.
- [ ] Confirm final recommendation matches the evidence level.

## Decision Rules

- If acceptance criteria cannot be traced to implementation or validation evidence, do not mark the change verified.
- If tests failed, were skipped, or were not run, classify the result as conditional or blocked based on risk.
- If generated outputs or lockfiles are stale, require regeneration before pass.
- If implementation changes behavior without docs or changelog updates, flag verification incomplete for user-facing changes.
- If a security-sensitive path lacks negative validation, block or conditionally pass with explicit risk.
- If evidence conflicts, prefer repository state and command output over summaries.

## Finding Categories

- Unverified or unmapped acceptance criterion.
- Missing, stale, failed, skipped, or irrelevant validation evidence.
- Scope creep, unrelated refactor, or unintended behavior change.
- Missing edge-case, regression, migration, rollback, or negative validation.
- Generated file, schema, lockfile, docs, or platform-copy drift.
- Unsupported pass recommendation or hidden residual risk.

## Severity Guidance

- Critical: claimed verification hides failed validation for security, data integrity, destructive, or production-critical behavior.
- High: must-have requirement, migration, contract, rollback, or security path lacks credible validation before release.
- Medium: important edge case, generated artifact, documentation, or regression evidence is incomplete but trackable.
- Low: evidence formatting, traceability clarity, or non-blocking verification detail can be improved.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Verification table mapping requirement, evidence, status, and residual risk.
- Validation commands reviewed or run, with pass/fail/skipped state.
- Files and artifacts reviewed, including generated outputs and docs where relevant.
- Findings with severity, evidence, impact, and exact remediation.
- Final recommendation: pass, conditional pass, or block, with reasons.
- Open verification gaps with owner-oriented next steps.

## Acceptance Criteria

- All must-have requirements are traced to implementation and credible validation evidence.
- Validation results are current, honest, and tied to the reviewed repository state.
- Generated outputs, docs, schemas, and lockfiles are synchronized or called out.
- Security, migration, rollback, and compatibility risks have appropriate evidence.
- Residual risks are explicit and reflected in the pass/conditional/block recommendation.
- Assumptions are separated from verified facts.

## Anti-Patterns

- Accepting an author's summary instead of checking changed files and validation evidence.
- Treating green CI as sufficient when it does not cover the changed behavior.
- Ignoring skipped tests, flaky failures, stale generated files, or missing docs.
- Marking pass while hiding residual risks or unrun checks.
- Verifying only happy paths for security-sensitive or migration-heavy changes.
- Expanding review scope into unrelated refactoring advice without verification relevance.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
