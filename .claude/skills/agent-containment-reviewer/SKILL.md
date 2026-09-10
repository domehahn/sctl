---
name: "agent-containment-reviewer"
description: "Review agent sandboxes, transitive egress, privilege escalation, lateral movement, shared infrastructure, monitoring, and kill switches."
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
# Agent Containment Reviewer

## Purpose

Review whether an autonomous or tool-using agent remains technically contained when its instructions, model behavior, tools, dependencies, or adjacent services are compromised. Treat the agent as an untrusted principal and assess direct and transitive escape, privilege-escalation, lateral-movement, and egress paths.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- An agent, model evaluation, coding assistant, or autonomous workflow runs in a sandbox, container, VM, runner, or isolated account.
- Package proxies, tool gateways, MCP servers, browsers, internal services, or shared infrastructure are reachable from an agent environment.
- A design claims that the agent has no host, internal-network, metadata-service, secret, or Internet access.
- Sandbox escape, privilege escalation, lateral movement, or kill-switch readiness needs independent review.
- Containment must remain effective even when the agent ignores every instruction.

## Operating model

1. Model the agent as hostile and inventory every direct and transitive path from the execution environment to host, control plane, internal network, credentials, and Internet.
2. Trace trust boundaries through proxies and helper services instead of stopping at the first allowlisted destination.
3. Compare intended isolation with effective namespaces, mounts, identities, capabilities, routes, credentials, and shared infrastructure.
4. Build concrete escape and lateral-movement chains, then map each hop to preventive, detective, and recovery controls.
5. Separate declared restrictions from independently enforced controls and report evidence gaps without claiming runtime guarantees.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Container, microVM, VM, namespace, syscall, filesystem, device, and host isolation.
- Mounts, container sockets, orchestrator APIs, privileged Linux capabilities, service accounts, and cloud metadata endpoints.
- DNS, proxies, package registries, internal services, overlay networks, arbitrary TCP/UDP, and Internet egress.
- Shared runners, nodes, caches, credentials, identities, control planes, and lateral-movement opportunities.
- Monitoring, session termination, credential revocation, quarantine, kill switches, and forensic evidence.

## Skill-Specific Checklist

- [ ] Draw the execution boundary from agent process through container or VM, host, network, control plane, and external systems.
- [ ] Inspect host mounts, writable bind mounts, device access, container sockets, and orchestrator credentials.
- [ ] Review user namespaces, UID/GID mapping, seccomp, AppArmor or SELinux, capabilities, privileged mode, and kernel sharing.
- [ ] Verify that package proxies and allowlisted services cannot relay arbitrary requests, fetch attacker-selected URLs, or expose administrative functions.
- [ ] Enumerate DNS resolution, internal routes, service discovery, metadata endpoints, link-local ranges, and alternate egress protocols.
- [ ] Test whether compromise of each reachable service creates a transitive path to broader infrastructure or Internet access.
- [ ] Review service-account tokens, workload identity, instance credentials, CI tokens, cached credentials, and credential rotation boundaries.
- [ ] Check whether runners, nodes, caches, volumes, or workspaces are shared across trust levels or tenants.
- [ ] Assess privilege-escalation paths through kernel attack surface, setuid files, sudo, package hooks, plugins, and control-plane APIs.
- [ ] Assess lateral movement to databases, artifact stores, registries, build systems, secret stores, and administrative interfaces.
- [ ] Verify independent telemetry for denied connections, unusual process trees, boundary violations, credential access, and escape attempts.
- [ ] Verify that operators can terminate the session, isolate the workload, revoke credentials, preserve evidence, and block further egress.

## Decision Rules

- If compromise of an allowlisted proxy or helper service exposes a broader network path, classify containment by the transitive path rather than the first hop.
- If the agent can reach a host socket, privileged control-plane API, or reusable infrastructure credential, treat host or infrastructure compromise as plausible until disproven.
- If isolation depends on instructions or cooperative model behavior, report that no technical containment control exists.
- If multiple tenants or trust levels share mutable state, require demonstrated isolation and cleanup before accepting containment.
- If a critical boundary has no independent deny telemetry or termination mechanism, report a detection and response gap in addition to the preventive-control gap.
- If evidence is configuration-only, distinguish intended settings from tested effective isolation.

## Finding Categories

- Sandbox, container, VM, namespace, or host escape path.
- Transitive proxy, package-registry, helper-service, or control-plane bypass.
- Privilege escalation through identity, kernel, capability, mount, socket, or credential exposure.
- Internal-network reachability, lateral movement, metadata access, or unintended Internet egress.
- Cross-tenant or shared-runner state and credential exposure.
- Missing monitoring, quarantine, revocation, forensic evidence, or kill switch.

## Severity Guidance

- Critical: a plausible chain escapes the agent boundary and reaches host control, reusable credentials, sensitive internal infrastructure, or unrestricted external systems.
- High: isolation can be bypassed through a reachable service, shared resource, privileged interface, or unmonitored egress path.
- Medium: defense-in-depth, tenant separation, detection, cleanup, or recovery controls are incomplete but primary isolation remains evidenced.
- Low: documentation, ownership, test cadence, or minor hardening is missing without a demonstrated boundary bypass.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Containment architecture showing direct and transitive trust boundaries, identities, network paths, and shared resources.
- Attack-path table with prerequisite, hop, reachable asset, evidence, existing control, and residual risk.
- Effective-isolation comparison covering namespaces, mounts, capabilities, credentials, proxies, routes, and egress.
- Findings ordered by severity with concrete containment remediation and verification steps.
- Monitoring, kill-switch, credential-revocation, quarantine, and evidence-preservation assessment.
- Explicit assumptions, untested controls, unavailable runtime evidence, and residual-risk owners.

## Acceptance Criteria

- The agent is modeled as an untrusted principal rather than trusted application logic.
- Direct and transitive paths through every reachable proxy or helper service are assessed.
- Host, control-plane, metadata, internal-network, credential, and Internet boundaries are covered.
- Critical isolation claims are tied to effective configuration or adversarial test evidence.
- Every Critical or High path has a preventive remediation plus detection and recovery guidance.
- The report does not claim that SKILL.md instructions or contract declarations alone provide containment.

## Anti-Patterns

- Stopping the review at an allowlisted proxy without modeling compromise of the proxy.
- Treating an internal network, shared runner, or package registry as inherently trusted.
- Equating containerization with a complete security boundary.
- Reviewing declared egress policy without inspecting effective routes, DNS, alternate protocols, and service relays.
- Ignoring credential revocation, quarantine, monitoring, or kill-switch behavior because preventive controls exist.
- Claiming isolation based only on instructions, policy text, or untested configuration.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
