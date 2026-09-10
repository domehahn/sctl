---
name: "agentic-threat-modeler"
description: "Threat-model agents as untrusted principals across prompts, tools, MCP, memory, delegation, runtime infrastructure, and external systems."
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
# Agentic Threat Modeler

## Purpose

Create threat models for agentic systems by treating agents, models, tools, memory, retrieved content, MCP servers, other agents, and runtime infrastructure as distinct and potentially untrusted principals. Translate agent-specific attack paths into Goal, Contract, eval, containment, monitoring, and enforcement requirements.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A tool-using, autonomous, multi-agent, RAG, MCP, coding, browser, or operational agent is designed or materially changed.
- Agent identities, capabilities, memory, context, tools, approvals, runtime boundaries, or external systems need threat analysis.
- Prompt injection, tool poisoning, excessive agency, sandbox escape, persistence, or cross-agent compromise is plausible.
- A general threat model does not distinguish model behavior from runtime enforcement.
- Goal, contract.yaml, eval scenarios, and infrastructure controls need a common threat basis.

## Operating model

1. Model the agent as an untrusted principal whose outputs and chosen actions may be adversarial even without a malicious user.
2. Inventory assets, identities, Goals, capabilities, tools, memory, data sources, context channels, runtimes, and enforcement points.
3. Trace direct and indirect instruction flows plus authority transitions across users, content, agents, tools, and infrastructure.
4. Develop abuse cases for compromised models, poisoned context, malicious tools, unsafe autonomy, and infrastructure escape.
5. Map each threat to prevention, detection, response, Contract clauses, security invariants, and behavioral evals.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Direct and indirect prompt injection, instruction hierarchy, retrieved content, memory poisoning, and context confusion.
- Tool poisoning, MCP server compromise or rug pull, plugin updates, delegated agents, and confused deputies.
- Excessive agency, goal misgeneralization, specification gaming, reward hacking, unbounded autonomy, and approval bypass.
- Credential harvesting, secret egress, privilege escalation, sandbox escape, lateral movement, persistence, and self-modification.
- Identity, tenant, runtime, data, network, tool, memory, monitoring, kill-switch, and human-approval trust boundaries.

## Skill-Specific Checklist

- [ ] Identify users, operators, agents, models, subagents, tools, MCP servers, data providers, runtimes, and external principals.
- [ ] Inventory sensitive assets, Goals, contract boundaries, credentials, memory, context, retrieved data, outputs, and controlled systems.
- [ ] Draw trust boundaries and authority transitions for prompts, retrieved content, tool results, memory writes, approvals, and actions.
- [ ] Model direct prompt injection and indirect injection through documents, web pages, issues, logs, code, images, and tool output.
- [ ] Model tool poisoning, description manipulation, MCP compromise, rug pulls, plugin updates, and parameter substitution.
- [ ] Model memory poisoning, cross-session contamination, cross-agent injection, delegated authority, and identity confusion.
- [ ] Model goal hacking, specification gaming, leaked-solution use, reward manipulation, and unsafe shortcut selection.
- [ ] Model excessive agency, approval bypass, recursive delegation, retries, concurrency, and long-horizon autonomy.
- [ ] Model credential access, data exfiltration, covert channels, privilege escalation, sandbox escape, and lateral movement.
- [ ] Model self-modification, hidden triggers, persistence, security-control tampering, monitoring evasion, and kill-switch resistance.
- [ ] Map each high-risk threat to contract capabilities, tools, data flows, approvals, limits, and structured invariants.
- [ ] Define adversarial evals, runtime enforcement, monitoring signals, incident response, recovery, and residual-risk ownership.

## Decision Rules

- If content can influence agent decisions without equivalent authority, treat it as an untrusted instruction channel.
- If a tool or MCP server can change after approval, model update integrity and rug-pull risk separately from initial trust.
- If the agent can reach a capability through delegation or a helper, include it in effective authority even when not directly declared.
- If Goal success rewards an unsafe shortcut, require an invariant and trajectory-level eval rather than output-only acceptance.
- If a threat is mitigated only by instructions, keep the technical residual risk open.
- If a Critical or High threat lacks prevention, detection, response, and an owner, mark the design incomplete.

## Finding Categories

- Prompt, retrieved-content, memory, context, or cross-agent injection.
- Tool, MCP, plugin, delegated-agent, update, or confused-deputy compromise.
- Goal hacking, specification gaming, excessive agency, approval bypass, or unbounded autonomy.
- Credential harvesting, secret egress, covert channel, or sensitive-data misuse.
- Privilege escalation, sandbox escape, lateral movement, persistence, self-modification, or hidden trigger.
- Missing runtime enforcement, monitoring, kill switch, incident response, eval coverage, or residual-risk owner.

## Severity Guidance

- Critical: plausible agent-controlled path compromises host, privileged infrastructure, reusable credentials, sensitive data, or external systems.
- High: poisoned instructions, tools, memory, delegation, or goals can bypass a material control or approval.
- Medium: defense-in-depth, monitoring, eval coverage, context isolation, or recovery control is incomplete.
- Low: documentation, ownership, naming, diagram, or low-impact hardening needs improvement.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Agentic-system diagram showing principals, Goals, tools, MCP servers, memory, data, runtime, identities, and trust boundaries.
- Asset, authority, instruction-channel, capability, and data-flow inventory.
- Threat register with attack path, preconditions, agent behavior, affected asset, controls, and residual risk.
- Goal and Contract recommendations covering required/allowed capabilities, tools, data flows, approvals, limits, and invariants.
- Adversarial eval plan plus runtime monitoring, containment, enforcement, kill-switch, and response requirements.
- Assumptions, unsupported controls, accepted risks, decision owners, and review triggers.

## Acceptance Criteria

- The agent and every delegated component are modeled as potentially untrusted principals.
- Instruction, authority, identity, data, memory, tool, runtime, and network boundaries are explicit.
- Agent-specific threats extend beyond conventional application STRIDE analysis.
- Every Critical or High threat maps to controls, evals, monitoring, response, and an owner.
- Declared Contract boundaries are distinguished from independently enforced runtime boundaries.
- Residual risks and assumptions remain explicit where evidence or enforcement is unavailable.

## Anti-Patterns

- Treating the model as ordinary trusted application code.
- Considering only malicious users while ignoring compromised tools, content, memory, models, or infrastructure.
- Listing prompt injection without tracing the capability or asset it can reach.
- Assuming direct tool permissions capture delegated, proxied, nested-agent, or confused-deputy authority.
- Calling instructions, guardrails, or a Contract declaration an enforced mitigation.
- Producing threats without corresponding eval, monitoring, containment, response, and residual-risk decisions.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
