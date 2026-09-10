---
name: "agent-behavior-eval-engineer"
description: "Design trajectory-level agent security evals for Goal compliance, contract boundaries, goal hacking, unsafe tools, and long-horizon behavior."
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
# Agent Behavior Eval Engineer

## Purpose

Design declarative, reproducible behavioral and adversarial evaluations for complete agent trajectories. Test Goal achievement together with contract compliance, acceptable tool use, state transitions, resistance to goal hacking, and preservation of security boundaries without embedding executable predicates.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- An agentic system needs baseline, adversarial, boundary, or long-horizon security evaluations.
- Goal success can be achieved through forbidden shortcuts, leaked answers, unsafe tools, or specification gaming.
- Tool calls, filesystem changes, network access, secret handling, approvals, and final output must be assessed as one trajectory.
- An incident, threat model, contract, or invariant must become a regression eval.
- Evaluation evidence needs stable scenario, Goal criterion, invariant, capability, and contract-digest references.

## Operating model

1. Derive scenarios from Goal criteria, contract boundaries, invariants, threat models, and real failure modes.
2. Define observable initial state, adversarial stimulus, permitted success path, forbidden shortcuts, assertions, and evidence.
3. Evaluate the full trajectory rather than accepting a correct final answer produced through prohibited behavior.
4. Separate deterministic instrumentation from model-based judging and document uncertainty, flakiness, and false-positive controls.
5. Keep eval specifications declarative; place executable harness logic in a separate trusted eval runner.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Baseline task success, Goal-criterion traceability, output quality, and contract-preserving execution.
- Prompt injection, indirect injection, goal hacking, specification gaming, reward hacking, and hidden-answer retrieval.
- Tool, process, filesystem, network, secret, approval, data-egress, and resource-limit trajectories.
- Sandbox escape attempts, privilege escalation, lateral movement, persistence, self-modification, and cross-agent influence.
- Long-horizon decomposition, retries, memory, context changes, evaluator integrity, fixtures, evidence, and reproducibility.

## Skill-Specific Checklist

- [ ] Reference stable Goal criterion IDs and contract invariant IDs instead of matching prose alone.
- [ ] Define the initial repository, filesystem, network, identity, secret, memory, and tool state required by each scenario.
- [ ] Include a legitimate permitted path so the eval distinguishes secure success from blanket refusal.
- [ ] Add adversarial instructions that request direct writes, secret disclosure, policy bypass, unsafe tools, or approval avoidance.
- [ ] Add environmental temptations such as leaked solutions, overprivileged credentials, unexpected internal services, or poisoned tool descriptions.
- [ ] Assert forbidden capabilities and tools across the complete trajectory, not only the final response.
- [ ] Assert expected state preservation for repository, filesystem, external systems, credentials, logs, and memory.
- [ ] Cover retries, parallel calls, nested agents, context compaction, delayed triggers, and long-horizon plans.
- [ ] Capture tool-call parameters, policy decisions, state diffs, network destinations, approvals, limits, and final output as evidence.
- [ ] Define deterministic pass, fail, inconclusive, infrastructure-error, and judge-disagreement handling.
- [ ] Protect eval fixtures, hidden answers, judges, and telemetry from agent access or manipulation.
- [ ] Bind results to skill version, eval schema version, scenario ID, runtime version, and canonical contract digest.

## Decision Rules

- If the Goal is achieved through a forbidden capability, tool, data flow, or invariant violation, the scenario fails.
- If a scenario offers no secure success path, classify refusal behavior separately from useful safe completion.
- If an assertion cannot be observed reliably, mark it unevaluated rather than inferring compliance from the final text.
- If model judgment is required, retain the rubric, evidence, judge version, disagreement, and deterministic checks.
- If an agent can read hidden answers, fixtures, or expected assertions, treat the eval as contaminated.
- If an incident exposes a new trajectory, add a minimized regression scenario plus boundary variants.

## Finding Categories

- Goal hacking, specification gaming, reward hacking, or forbidden shortcut.
- Unsafe capability, tool, data-flow, approval, or limit trajectory.
- Sandbox escape, lateral movement, persistence, self-modification, or external-system impact.
- Missing observation, ambiguous assertion, weak oracle, or final-output-only evaluation.
- Eval contamination, judge manipulation, fixture leakage, or non-reproducible environment.
- Long-horizon, retry, concurrency, nested-agent, memory, or context-transition coverage gap.

## Severity Guidance

- Critical: eval design would pass successful exploitation, secret exfiltration, boundary escape, persistence, or external compromise.
- High: material contract violations, goal-hacking paths, or unobserved privileged actions can pass as success.
- Medium: scenario coverage, evidence binding, reproducibility, judge quality, or long-horizon variation is incomplete.
- Low: naming, documentation, fixture organization, or non-blocking rubric clarity needs improvement.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Eval plan mapping Goal criteria, contract clauses, invariants, threats, and incidents to scenarios.
- Declarative baseline, adversarial, boundary, and long-horizon scenario specifications.
- Trajectory evidence model covering inputs, tool calls, policy decisions, state diffs, network activity, approvals, limits, and outputs.
- Pass/fail/inconclusive rubric with deterministic assertions and clearly separated model-based judgments.
- Fixture-isolation, hidden-answer protection, reproducibility, and contract-digest binding requirements.
- Coverage gaps, expected false positives, infrastructure dependencies, and recommended regression cadence.

## Acceptance Criteria

- Goal success and contract compliance are independently asserted.
- Every high-risk scenario observes the relevant trajectory and state changes.
- Baseline scenarios demonstrate useful secure completion rather than refusal-only behavior.
- Adversarial scenarios cover malicious prompts and tempting or compromised environments.
- Eval definitions remain declarative and contain no shell expressions, scripts, or executable predicates.
- Results can be reproduced and tied to exact scenario, runtime, skill, and contract versions.

## Anti-Patterns

- Evaluating only whether the final answer looks correct.
- Rewarding Goal completion while ignoring forbidden tools, state changes, or data flows.
- Using refusal as the only acceptable secure behavior when a permitted path exists.
- Letting the agent read hidden answers, judge prompts, fixtures, assertions, or policy decisions.
- Embedding scripts or arbitrary executable expressions inside declarative eval YAML.
- Reporting a single stochastic run as conclusive long-horizon security evidence.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
