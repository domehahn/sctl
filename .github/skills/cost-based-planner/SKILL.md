---
name: "cost-based-planner"
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
# Cost Based Planner

## Purpose

Plan work by balancing implementation cost, operational cost, maintenance cost, uncertainty, risk, and delivered value. Recommend the smallest useful implementation path and expose trade-offs before coding begins.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A request needs sizing, sequencing, or phased delivery.
- The implementation path has meaningful uncertainty, migration, infrastructure, or maintenance cost.
- Build-versus-buy, MVP scope, or cost-of-delay decisions are open.
- A broad change should be decomposed into lower-risk increments.
- The user needs a plan before implementation starts.

## Operating model

1. Identify cost drivers across implementation, operation, maintenance, migration, licensing, infrastructure, and opportunity cost.
2. Use repository evidence to estimate effort, blast radius, and validation cost.
3. Separate MVP, incremental rollout, and full-scope options.
4. Make uncertainty visible and reduce it with targeted file reads or experiments.
5. Recommend the smallest useful path that preserves rollback and validation.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Implementation effort and complexity drivers.
- Operational, maintenance, migration, license, infrastructure, and support cost.
- Uncertainty, assumptions, dependencies, and cost-of-delay.
- MVP, phased delivery, rollout, rollback, and validation cost.
- Build-versus-buy and reuse-versus-new-code trade-offs.

## Skill-Specific Checklist

- [ ] Separate one-time implementation cost from recurring operational cost.
- [ ] Identify the main cost drivers and complexity drivers.
- [ ] Identify uncertainty factors and assumptions.
- [ ] Compare MVP, incremental rollout, and full-scope implementation.
- [ ] Identify build-versus-buy and reuse trade-offs.
- [ ] Check budget constraints, funding limits, and expected cost envelope.
- [ ] Identify hidden maintenance, migration, and support costs.
- [ ] Estimate effort using repository-specific evidence when possible.
- [ ] Identify high-cost dependencies, integrations, and platform changes.
- [ ] Recommend the smallest useful implementation path.
- [ ] Call out cost risks that affect prioritization.
- [ ] Include rollback and validation cost in the plan.
- [ ] Name decisions that require stakeholder input.

## Decision Rules

- If a low-cost MVP can validate the goal, recommend it before full-scope work.
- If uncertainty dominates cost, plan a discovery step before implementation.
- If a dependency adds recurring operational burden, include it in prioritization.
- If migration cost is high, separate migration from feature delivery.
- If build-versus-buy is unresolved, compare license, integration, maintenance, and lock-in costs.
- If validation cost is high, include it as part of delivery cost rather than a footnote.

## Finding Categories

- High implementation complexity.
- Hidden operational or maintenance cost.
- Migration or rollback cost risk.
- License, infrastructure, or vendor cost exposure.
- Unclear value, priority, or cost-of-delay.
- Uncertainty requiring discovery before build.

## Severity Guidance

- Critical: cost or migration risk makes the proposed path unsafe without a different plan.
- High: major recurring cost, irreversible migration, or expensive dependency is unaccounted for.
- Medium: cost estimates rely on assumptions but can be reduced with discovery.
- Low: minor sequencing or documentation issue affects planning clarity.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Costed plan with MVP, incremental, and full-scope options.
- Budget impact summary covering expected spend, constraints, and unresolved approvals.
- List of cost drivers with evidence and assumptions.
- Risk and uncertainty register with reduction steps.
- Recommended implementation sequence and rollback points.
- Validation plan with expected command or review cost.
- Build-versus-buy or reuse rationale where applicable.

## Acceptance Criteria

- Costs are split into implementation, operation, maintenance, migration, license, and infrastructure where relevant.
- The recommended plan names the smallest useful implementation path.
- Uncertainty and assumptions are visible with next steps to reduce them.
- High-cost dependencies and rollback constraints are identified.
- Validation and release costs are included in the plan.
- The plan can be executed incrementally.

## Anti-Patterns

- Ignoring recurring operational cost.
- Treating the full solution as the only option.
- Estimating effort without reading relevant repository evidence.
- Hiding migration or rollback cost until implementation.
- Adding dependencies without considering maintenance and license cost.
- Optimizing for low upfront effort while increasing long-term support burden.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
