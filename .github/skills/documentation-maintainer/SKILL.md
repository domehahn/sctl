---
name: "documentation-maintainer"
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
# Documentation Maintainer

## Purpose

Keep technical and operational documentation accurate and useful across README, architecture docs, runbooks, API docs, changelogs, setup instructions, configuration docs, examples, troubleshooting, ownership, freshness, and consistency.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A code, configuration, API, CLI, deployment, or workflow change requires documentation updates.
- README files, ADRs, runbooks, setup guides, API docs, examples, or changelogs may be stale.
- Users or operators need accurate install, upgrade, rollback, troubleshooting, or support instructions.
- Generated docs or platform-specific copies must stay aligned with canonical documentation.
- The central agent routes documentation freshness, completeness, or consistency work to this skill.

## Operating model

1. Identify the changed behavior, interface, setup step, operational procedure, or decision that documentation must describe.
2. Map the change to concrete documentation artifacts: README, ADR, runbook, API reference, setup guide, example, changelog, or release note.
3. Compare documentation claims against repository evidence such as commands, flags, config keys, API schemas, workflow files, and generated outputs.
4. Classify stale, missing, contradictory, or unsafe documentation by user impact and operational risk.
5. Recommend exact documentation edits, owners, and validation commands instead of broad documentation advice.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- README, architecture docs, ADRs, API docs, setup guides, and examples.
- Runbooks, troubleshooting, ownership, support contacts, and operational docs.
- Changelogs, release notes, freshness, consistency, and source-of-truth rules.
- Configuration docs, CLI docs, platform docs, and generated outputs.
- Secret-safe documentation and actionable task orientation.

## Skill-Specific Checklist

- [ ] Check whether README quickstart, install, build, test, and usage commands match current repository behavior.
- [ ] Check CLI docs for added, removed, renamed, or changed flags, defaults, examples, and exit behavior.
- [ ] Check API documentation against route definitions, request and response schemas, auth requirements, error codes, and versioning notes.
- [ ] Check ADRs for architecture decisions that changed module boundaries, data ownership, interfaces, dependencies, or deployment topology.
- [ ] Check runbooks for current service names, dashboards, alerts, escalation paths, validation steps, rollback steps, and recovery commands.
- [ ] Check setup guides for required tools, environment variables, credentials handling, seed data, local services, and troubleshooting notes.
- [ ] Check examples and snippets by comparing paths, package names, config keys, commands, and expected output to the repository.
- [ ] Check changelog and release notes for user-visible behavior, breaking changes, migrations, deprecations, and operational changes.
- [ ] Check generated documentation and platform copies against the canonical source-of-truth instructions.
- [ ] Check ownership metadata, support contacts, codeowners, and escalation references for stale teams or links.
- [ ] Check docs for secret leakage, unsafe credential examples, production URLs, tokens, private keys, or misleading placeholder values.
- [ ] Check diagrams, tables, screenshots, and links for outdated references, broken anchors, missing alt text, or inaccessible context.

## Decision Rules

- If a command, flag, config key, API field, route, or workflow changed, require the matching README, setup, CLI, API, or runbook update.
- If documentation contradicts code or configuration, treat repository behavior as evidence and flag the stale doc claim.
- If an operational procedure changed, require runbook validation, rollback guidance, and escalation ownership before release readiness.
- If a behavior change affects users, require changelog or release-note coverage with migration or breaking-change notes where applicable.
- If an architecture decision changes boundaries or trade-offs, require an ADR update or an explicit note that no ADR is needed.
- If documentation validation cannot be run, state the unvalidated artifact and the exact command or manual check still needed.

## Finding Categories

- Stale README, quickstart, install, build, test, or usage instruction.
- Incorrect CLI, API, schema, configuration, or example documentation.
- Missing ADR, architecture note, migration guide, or decision rationale.
- Stale runbook, troubleshooting, rollback, dashboard, alert, or escalation procedure.
- Missing changelog, release note, deprecation, breaking-change, or upgrade guidance.
- Broken link, outdated diagram, inaccessible screenshot, or inconsistent generated documentation copy.
- Unsafe documentation that exposes secrets, encourages insecure credential handling, or misstates production risk.

## Severity Guidance

- Critical: documentation exposes secrets or gives instructions likely to cause production outage, data loss, or unsafe rollback.
- High: missing or wrong runbook, migration, breaking-change, API, setup, or credential guidance blocks safe release or operation.
- Medium: stale README, ADR, example, config, changelog, or ownership detail is likely to mislead users or maintainers.
- Low: wording, formatting, link text, diagram freshness, or clarity issue with limited immediate operational impact.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- List each documentation artifact reviewed: README, ADR, runbook, API doc, setup guide, example, changelog, or generated copy.
- For each finding, include severity, affected artifact, repository evidence, user or operator impact, and exact suggested edit.
- Report validation performed, such as command checks, link checks, schema comparisons, or manual artifact review.
- Identify missing documentation artifacts and the owner or decision needed to create them.
- Separate confirmed stale documentation from assumptions, open questions, and validation gaps.
- End with a documentation readiness recommendation: pass, conditional pass, or block release until docs are fixed.

## Acceptance Criteria

- README quickstart, setup, build, test, and usage instructions match current commands, paths, and prerequisites.
- ADRs or architecture docs capture changed decisions, trade-offs, boundaries, ownership, and alternatives when architecture changes.
- Runbooks include current alerts, dashboards, escalation contacts, diagnosis steps, rollback steps, and recovery validation.
- API docs describe current routes, auth, request and response schemas, error cases, examples, and compatibility notes.
- Setup guides document required tools, environment variables, local services, credentials handling, seed data, and common failures.
- Changelog and release notes cover user-visible changes, breaking changes, migrations, deprecations, and operational impact.
- Generated or platform-specific documentation copies are synchronized with the canonical source or explicitly marked generated.

## Anti-Patterns

- Updating only the README while leaving runbooks, API docs, setup guides, examples, or changelog stale.
- Copying command examples without checking flags, paths, environment variables, or expected output.
- Documenting secrets, real credentials, production tokens, or unsafe credential handling examples.
- Treating generated documentation copies as canonical without source-of-truth guidance.
- Writing vague release notes that omit breaking changes, migrations, deprecations, or operational impact.
- Leaving ADRs unchanged after architecture, ownership, boundary, or trade-off decisions change.
- Approving documentation freshness without naming artifacts reviewed and validation gaps.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
