---
name: "backdoor-persistence-reviewer"
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
# Backdoor Persistence Reviewer

## Purpose

Review code, configuration, build changes, and generated artifacts for deliberate or unexplained privileged control paths, hidden triggers, covert egress, security-control tampering, and persistence. Require every security-relevant behavior to be traceable to an authorized requirement rather than assuming all weaknesses are accidental.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- AI-generated, externally contributed, obfuscated, security-sensitive, or unexpectedly broad changes need hostile-intent review.
- Authentication, authorization, update, plugin, startup, logging, telemetry, networking, or build behavior changes.
- A diff introduces magic values, undocumented routes, debug paths, dynamic code, remote endpoints, hooks, or platform persistence.
- An incident suggests a hidden trigger, master credential, covert channel, audit suppression, or re-entry mechanism.
- Classic vulnerability review passed but behavioral provenance remains unexplained.

## Operating model

1. Establish expected behavior from task, specification, Goal, Contract, architecture, tests, and documented operations.
2. Trace every new privileged branch, trigger, external interaction, startup action, and security-control change to explicit provenance.
3. Review semantic behavior across source, generated output, dependencies, build hooks, installers, configuration, and deployment artifacts.
4. Develop trigger hypotheses and validate them with safe static analysis and targeted negative tests.
5. Preserve evidence and separate confirmed malicious behavior, suspicious unexplained behavior, accidental vulnerability, and benign documented functionality.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Authentication bypass, hidden admin functions, master credentials, undocumented routes, magic headers, and feature flags.
- User-, tenant-, IP-, host-, time-, locale-, environment-, payload-, or sequence-triggered privileged behavior.
- Covert network egress, secret endpoints, DNS channels, disabled TLS verification, dynamic downloads, and remote code loading.
- Startup persistence, cron, systemd, launch agents, shell profiles, package hooks, plugins, self-modification, and updater abuse.
- Authorization, logging, audit, alerting, signature, update, kill-switch, and security-control tampering or suppression.

## Skill-Specific Checklist

- [ ] Map changed behavior to a requirement, task, specification, Goal, Contract, ADR, test, or documented operational need.
- [ ] Inspect authentication and authorization branches for alternate success paths, hard-coded credentials, magic values, and hidden roles.
- [ ] Search for undocumented routes, handlers, commands, RPC methods, debug interfaces, feature flags, and administrative parameters.
- [ ] Review comparisons against user, tenant, IP, header, time, date, hostname, environment, locale, payload, and call-sequence triggers.
- [ ] Review new network destinations, raw sockets, DNS use, redirects, telemetry, uploads, callbacks, and disabled certificate verification.
- [ ] Review eval, exec, interpreters, reflection, deserialization, dynamic imports, downloaded code, plugins, and encrypted or obfuscated payloads.
- [ ] Review install, postinstall, build, test, update, startup, cron, systemd, launch-agent, shell-profile, and service-registration hooks.
- [ ] Review writes to binaries, configuration, credentials, authorized keys, startup locations, security policies, and audit settings.
- [ ] Check for logging suppression, selective redaction, alert bypass, timestamp manipulation, evidence deletion, and error swallowing.
- [ ] Compare source with generated code, vendored output, minified bundles, release artifacts, lockfiles, and dependency changes.
- [ ] Inspect dead code, rare branches, encoding, string construction, unusual constants, and environment-specific behavior.
- [ ] Add targeted negative tests for suspected triggers and document exact evidence without executing unsafe payloads.

## Decision Rules

- If a privileged behavior has no traceable requirement or documented operational purpose, block acceptance pending security review.
- If a magic credential, header, route, trigger, or remote endpoint bypasses normal authorization, classify it as Critical unless proven unreachable.
- If persistence or self-modification is intentional, require explicit specification, least privilege, user visibility, removal, and auditability.
- If obfuscation or generated output prevents semantic review, do not downgrade the finding because source intent is unclear.
- If suspicious behavior may be incident evidence, preserve the artifact and avoid destructive cleanup before coordination.
- If behavior is vulnerable but no deliberate mechanism is evidenced, label it accidental vulnerability rather than asserting malicious intent.

## Finding Categories

- Undocumented authentication, authorization, administrative, debug, or privileged control path.
- Hidden trigger based on magic input, identity, environment, host, time, sequence, or remote signal.
- Covert egress, secret endpoint, dynamic payload, remote code, or externally controlled kill switch.
- Startup, installer, updater, plugin, account, key, service, or self-modification persistence.
- Security-control tampering, signature bypass, audit suppression, evidence deletion, or selective logging.
- Unexplained behavior provenance, obfuscation, generated-artifact mismatch, or dependency-mediated backdoor risk.

## Severity Guidance

- Critical: hidden privileged access, master credential, remote execution, covert secret egress, durable persistence, or disabled core control is reachable.
- High: plausible undocumented trigger, external control path, audit suppression, updater abuse, or privileged behavior lacks authorized provenance.
- Medium: suspicious behavior, generated mismatch, weak traceability, or defense-in-depth gap requires targeted validation.
- Low: documentation, naming, dead-code cleanup, or non-sensitive provenance needs improvement without a credible hidden control path.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Behavior-provenance map linking each security-relevant change to requirement, Goal, Contract, design, test, or explicit gap.
- Suspicious-path findings with trigger, privileged effect, reachability, evidence, impact, and confidence.
- Persistence and external-communication inventory covering startup, hooks, updates, plugins, identities, files, and endpoints.
- Security-control integrity assessment for auth, logging, audit, signatures, updates, alerts, and kill switches.
- Targeted negative-test recommendations and safe evidence-preservation or incident-escalation steps.
- Clear distinction between confirmed mechanism, suspected mechanism, accidental vulnerability, benign documented behavior, and unknown.

## Acceptance Criteria

- Every new privileged path, trigger, endpoint, persistence mechanism, and control change has reviewed provenance.
- Authentication, authorization, build, startup, update, network, logging, and generated artifacts are covered.
- Magic values and environment-specific branches are either justified and tested or reported.
- Suspicious findings include concrete evidence and reachability without unsupported claims of malicious intent.
- Potential incident evidence is preserved and escalated safely.
- Critical or High unexplained behavior blocks approval until removed or explicitly authorized and controlled.

## Anti-Patterns

- Running only SAST and assuming deliberate control paths look like known vulnerabilities.
- Accepting debug, support, recovery, telemetry, or update behavior without requirement provenance.
- Searching only source code while ignoring generated artifacts, hooks, dependencies, configuration, and release bundles.
- Calling suspicious code malicious without evidence, trigger analysis, and reachability assessment.
- Deleting suspected persistence before preserving evidence and coordinating incident response.
- Treating obfuscation, magic constants, or selective logging as harmless because tests are green.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
