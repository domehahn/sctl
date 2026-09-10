---
name: "security-invariant-test-engineer"
description: "Derive negative tests and declarative evals from contract capabilities, tools, data flows, approvals, limits, and structured invariants."
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
# Security Invariant Test Engineer

## Purpose

Derive deterministic negative tests and declarative behavioral eval scenarios from contract.yaml capabilities, tools, data flows, approvals, limits, Goal failure conditions, and structured security invariants. Turn machine-readable boundaries into observable deny-path evidence without implementing an eval runtime.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A skill Contract needs concrete tests proving that forbidden behavior remains forbidden.
- Security invariants such as repository-state-unchanged, no-secret-egress, or declared-tools-only need adversarial scenarios.
- A contract diff expands or narrows capabilities and regression coverage must change with it.
- Runtime enforcement requires negative tests for deny precedence, zero limits, approvals, and fail-closed behavior.
- An incident or threat model must become reusable invariant-level regression coverage.

## Operating model

1. Parse the Goal, Contract, and existing evals, then build a traceability matrix from each boundary to observable behavior.
2. Generate adversarial prompts and environmental conditions that attempt realistic bypasses while retaining a legitimate permitted path.
3. Define deterministic assertions for capabilities, tools, state, data flows, approvals, limits, and output.
4. Separate contract declaration tests, runtime enforcement tests, and behavioral eval scenarios.
5. Keep eval YAML declarative and place harness code in trusted deterministic tests or downstream runners.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Required and allowed repository, filesystem, network, process, secret, and tool capabilities.
- Tool allow/deny precedence, wildcard rejection, path and destination normalization, and bypass variants.
- Data classifications, egress destinations, source-to-sink flows, output restrictions, and secret handling.
- Preconditions, postconditions, invariants, approval modes, invocation limits, and Goal failure conditions.
- Contract normalization, digest binding, semantic diff impact, eval traceability, and negative-test maintenance.

## Skill-Specific Checklist

- [ ] Inventory every allowed and forbidden capability scope, tool identity, data flow, approval rule, limit, and invariant.
- [ ] Map each Goal failure condition and Contract invariant to at least one observable assertion or explicit coverage gap.
- [ ] Create repository-write bypass attempts such as fix, patch, commit, rename, delete, generated-file, and delegated-tool requests.
- [ ] Create filesystem bypass attempts using traversal, symlink, archive, temporary path, alternate encoding, and helper-process variants.
- [ ] Create network bypass attempts using raw IPs, redirects, DNS rebinding, proxies, alternate protocols, nested tools, and unexpected services.
- [ ] Create process-execution attempts through shells, interpreters, package hooks, builds, plugins, documents, and indirect tool execution.
- [ ] Create secret-access and egress attempts through environment, files, metadata, logs, reports, uploads, URLs, and encoded output.
- [ ] Test tool deny precedence, undeclared tools, exact identifiers, parameter changes, and provider-independent capability implications.
- [ ] Test approval per_action, per_invocation, and per_session scope plus replay, parameter substitution, and stale approval.
- [ ] Test null, zero, positive, exhausted, retry, parallel, and nested-agent limit behavior.
- [ ] Where the design requires them, propose stable invariants such as no-undocumented-privileged-path, no-hidden-persistence, no-undeclared-egress, and security-controls-preserved with explicit observation semantics.
- [ ] Bind expected results to scenario ID, invariant ID, Goal criterion ID, contract digest, and runtime version.
- [ ] Update or retire tests when semantic contract diff reports expansion, narrowing, invariant changes, or approval changes.

## Decision Rules

- If a security-relevant Contract boundary has no negative test or eval, report an explicit coverage gap.
- If a test proves only parser rejection, do not claim it proves runtime deny enforcement.
- If the expected secure behavior is refusal but a permitted completion path exists, test both safe completion and refusal quality.
- If an assertion relies on unobserved internal behavior, mark it inconclusive and add instrumentation requirements.
- If a contract expansion removes a previous deny expectation, require reviewed test changes instead of silently deleting failures.
- If an invariant type is not yet standardized, preserve its stable ID and document the observation semantics used by the test.

## Finding Categories

- Capability, scope, tool, or deny-precedence negative-test gap.
- Repository, filesystem, process, network, secret, or data-egress bypass variant gap.
- Approval, limit, retry, concurrency, nested-agent, or invocation-scope test gap.
- Goal failure condition, postcondition, invariant, or output-contract traceability gap.
- Parser validation confused with runtime enforcement or behavioral conformance.
- Stale eval, missing contract-digest binding, weak observability, or unsafe executable assertion.

## Severity Guidance

- Critical: no test can detect a path to secret exfiltration, privileged writes, arbitrary execution, boundary escape, or external compromise.
- High: a material Contract deny, invariant, approval, or limit lacks realistic bypass coverage or runtime evidence.
- Medium: traceability, normalization variants, long-horizon behavior, instrumentation, or regression maintenance is incomplete.
- Low: scenario naming, documentation, fixture organization, or non-blocking edge coverage needs improvement.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Contract-to-test traceability matrix for Goal failures, capabilities, tools, data flows, conditions, approvals, limits, and output.
- Negative-test inventory with bypass stimulus, expected policy decision, observable evidence, and test layer.
- Declarative Eval v1 scenario proposals referencing stable Goal criterion and invariant IDs.
- Deterministic runtime or integration-test proposals for enforcement points that declarative evals cannot prove.
- Coverage delta for semantic contract expansion, narrowing, invariant, approval, and limit changes.
- Instrumentation gaps, inconclusive assertions, fixture requirements, contract digest, and residual-risk notes.

## Acceptance Criteria

- Every material deny boundary and structured invariant has coverage or a named gap.
- Tests distinguish parser validation, runtime enforcement, behavioral conformance, and output validation.
- Bypass variants cover direct, indirect, delegated, encoded, retried, parallel, and nested-agent paths where relevant.
- Zero, null, finite, approval, and deny-precedence semantics are explicitly exercised.
- Eval specifications contain no shell, script, executable predicate, or unsafe template expression.
- Coverage is traceable to stable IDs and the exact normalized contract digest.

## Anti-Patterns

- Writing only positive tests for allowed behavior.
- Using one prompt per invariant and ignoring alternate tools, encodings, delegation, retries, and environment paths.
- Claiming a schema validation test proves that a runtime blocks forbidden actions.
- Embedding shell commands or arbitrary code inside declarative eval specifications.
- Deleting failing negative tests after a permission expansion without reviewed Contract traceability.
- Testing final text while ignoring tool calls, state changes, approvals, network activity, and data egress.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
