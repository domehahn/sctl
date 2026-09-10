---
name: "agent-runtime-enforcement-reviewer"
description: "Compare declared agent contracts with independent runtime enforcement across tools, files, processes, networks, secrets, approvals, and limits."
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
# Agent Runtime Enforcement Reviewer

## Purpose

Compare an agent's declared Goal and contract.yaml boundary with the independent runtime mechanisms that actually mediate tools, filesystem, repository, process, network, secrets, data egress, approvals, and limits. Identify declaration-to-enforcement gaps without treating the declaration itself as enforcement.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A contract.yaml or equivalent capability policy is deployed with an agent runtime.
- Tool gateways, sandboxes, firewalls, secret brokers, approval services, or resource governors are expected to enforce agent restrictions.
- Runtime architecture must demonstrate that undeclared tools and capabilities fail closed.
- Policy translation, identity binding, bypass resistance, or enforcement telemetry needs review.
- An organization needs evidence that the evaluated contract matches the contract enforced in production.

## Operating model

1. Normalize the declared contract into enforceable decisions and trace each decision to an independent runtime control.
2. Verify identity, contract digest, policy version, invocation scope, and fail-closed behavior at every enforcement point.
3. Exercise negative paths for undeclared destinations, tools, files, processes, secrets, approvals, and exhausted limits.
4. Search for alternate interfaces and confused-deputy paths that bypass the primary gateway.
5. Separate authoring validation, runtime enforcement, behavioral observation, and post-hoc detection in every conclusion.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Contract loading, schema validation, canonical digest binding, policy compilation, distribution, and cache invalidation.
- Tool mediation, repository and filesystem authorization, process sandboxing, network gateways, secret brokers, and data-loss prevention.
- Agent and invocation identity, tenant binding, delegated credentials, approval decisions, and replay resistance.
- Limit enforcement for tool calls, runtime, network requests, concurrency, retries, and nested agents.
- Denied-action telemetry, tamper resistance, failover behavior, policy drift, and emergency revocation.

## Skill-Specific Checklist

- [ ] Map every Contract v1 required and allowed field to a named runtime enforcement point or an explicit unsupported gap.
- [ ] Verify the runtime binds the intended skill, invocation, principal, and canonical contract digest before granting capabilities.
- [ ] Confirm malformed, missing, unsupported, stale, or unverifiable contracts fail closed.
- [ ] Test undeclared tool invocation and explicit tool deny precedence through every available tool interface.
- [ ] Test repository and filesystem scopes against alternate paths, symlinks, archives, generated files, and helper services.
- [ ] Test process.execute false against shells, interpreters, package hooks, build tools, plugins, and indirect execution.
- [ ] Test network destinations against raw IPs, DNS rebinding, redirects, proxies, tunnels, alternate protocols, and nested tools.
- [ ] Test secrets.read false against environment variables, mounted files, metadata services, logs, caches, and delegated tools.
- [ ] Verify approval modes are independently recorded, scoped, non-replayable, and re-evaluated when parameters change.
- [ ] Verify zero and finite limits are enforced across retries, parallel calls, nested agents, and gateway failover.
- [ ] Check that no direct SDK, shell, browser, MCP, plugin, or administrative path bypasses the policy gateway.
- [ ] Verify denied actions, policy errors, digest mismatches, overrides, and emergency revocations produce tamper-resistant telemetry.

## Decision Rules

- If a declared restriction has no independent mediator, classify it as unenforced even when instructions repeat the restriction.
- If any alternate interface reaches the protected resource without the same policy decision, treat the control as bypassable.
- If the runtime cannot bind policy to the exact invocation and contract digest, report policy-confusion and evaluation-integrity risk.
- If policy loading or a gateway outage falls back to allow, classify the affected boundary as a critical fail-open condition.
- If approval can be reused after action parameters, destination, capability, or invocation changes, report approval-scope bypass.
- If only post-hoc logs exist, label the control detective rather than preventive.

## Finding Categories

- Declared capability without an independent enforcement point.
- Tool, filesystem, process, network, secret, or data-egress gateway bypass.
- Contract digest, skill identity, invocation, tenant, or policy-version binding failure.
- Fail-open parsing, distribution, caching, outage, or fallback behavior.
- Approval replay, parameter substitution, limit evasion, or nested-agent accounting gap.
- Missing, mutable, or incomplete denied-action and policy-decision telemetry.

## Severity Guidance

- Critical: a forbidden capability is reachable, enforcement fails open, or policy identity can be substituted across tenants or invocations.
- High: a practical alternate interface, approval replay, stale policy, or accounting gap bypasses a material restriction.
- Medium: enforcement exists but coverage, telemetry, recovery, or negative-test evidence is incomplete.
- Low: naming, documentation, ownership, or non-critical observability needs improvement without changing effective authorization.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Contract-to-enforcement matrix for capabilities, tools, data boundaries, approvals, limits, and output restrictions.
- Runtime trust-boundary diagram naming policy decision and enforcement points.
- Negative-test evidence for allowed, denied, missing, malformed, zero-limit, and exhausted-limit cases.
- Bypass analysis covering alternate interfaces, delegated tools, nested agents, proxies, and failure modes.
- Findings with declared rule, effective behavior, evidence, impact, and remediation owner.
- Clear statement of what is preventive, detective, declared-only, untested, or delegated downstream.

## Acceptance Criteria

- Every security-relevant contract field maps to a tested control or explicit enforcement gap.
- The exact contract digest, principal, tenant, and invocation are bound at runtime.
- Missing, malformed, unsupported, or unavailable policy fails closed.
- Alternate access paths receive equivalent authorization and limit decisions.
- Approvals and limits preserve their declared scope under retries, parameter changes, concurrency, and nesting.
- Conclusions distinguish declaration, validation, observation, detection, and technical enforcement.

## Anti-Patterns

- Calling contract validation runtime enforcement.
- Testing only allowed operations and never proving deny paths.
- Reviewing the primary tool gateway while ignoring shell, SDK, browser, MCP, plugin, or administrative alternatives.
- Binding policy to a mutable filename or skill name without a canonical digest and invocation identity.
- Treating approval UI presence as proof of action scope, freshness, and replay resistance.
- Claiming prevention when the only control is logging or after-the-fact alerting.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
