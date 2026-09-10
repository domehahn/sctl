---
name: "universal-skill-creator"
description: "Create, adapt, validate, and optimize reusable agent skills across agentic platforms."
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
# Universal Skill Creator

## Purpose

Create new production-ready skills and prevent generic copy-paste skills by enforcing full frontmatter, SemVer, dates, authors, stability, min_platform_version, changelog, domain-specific scope, checklist, decision rules, finding categories, severity guidance, outputs, acceptance criteria, anti-patterns, and no generic body reuse.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- A user asks to create or upgrade a skill.
- A skill body must be checked for generator smell or copy-paste generic content.
- Skill metadata, versioning, compatibility, or changelog rules need enforcement.
- A skill needs domain-specific review logic, not only name and description changes.
- The central agent routes to universal skill creation.

## Operating model

1. Identify the skill domain, users, trigger conditions, outputs, risks, and non-goals.
2. Write domain-specific review scope, checklist, decisions, categories, severity, outputs, acceptance, and anti-patterns.
3. Reject bodies that only differ by name, purpose, or generic operating text.
4. Ensure frontmatter, body changelog, compatibility metadata, and versioning are consistent.
5. Validate generated examples or tests that prove the skill is not generic.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- YAML frontmatter, SemVer, since, last_modified, authors, stability, and min_platform_version.
- Body changelog, purpose, review scope, checklist, decision rules, and finding categories.
- Severity guidance, output requirements, acceptance criteria, and anti-patterns.
- Generic body reuse detection and platform compatibility honesty.
- Skill routing, generated copies, validation, and governance preservation.

## Skill-Specific Checklist

- [ ] Validate full YAML frontmatter and required metadata fields.
- [ ] Synchronize frontmatter changelog and body changelog version/date/message.
- [ ] Require min_platform_version entries for all supported platforms and mark unvalidated platforms honestly.
- [ ] Write a purpose that names the domain and concrete work products.
- [ ] Write When-to-use triggers that are specific enough for routing decisions.
- [ ] Write at least 10 checklist items that mention domain artifacts, risks, and evidence.
- [ ] Write at least 5 decision rules that decide real domain trade-offs.
- [ ] Write finding categories that are domain failure types, not generic evidence/control phrases.
- [ ] Write severity guidance with domain-specific Critical, High, Medium, and Low criteria.
- [ ] Write output requirements naming concrete artifacts the agent must produce.
- [ ] Write acceptance criteria that are testable for the domain.
- [ ] Write anti-patterns that describe misuse of this exact skill.
- [ ] Reject “structured analysis or review”, “<skill> evidence”, and “<skill> control” boilerplate.
- [ ] Verify generated platform copies use the shared renderer and stay synchronized.

## Decision Rules

- Never create a skill that only differs by name and description. Every generated skill must include domain-specific review scope, checklist items, decision rules, finding categories, severity guidance, output requirements, acceptance criteria, and anti-patterns. Generic operating-model text is allowed only as shared baseline, never as the complete skill body.
- If a checklist item could apply unchanged to most skills, rewrite it with domain artifacts and failure modes.
- If a finding category contains the skill name plus “evidence” or “control”, reject it as generator smell.
- If severity guidance does not say what is Critical/High/Medium/Low in this domain, reject production readiness.
- If compatibility versions are concrete without validation evidence, use unknown or mark compatibility unverified.
- If frontmatter and body changelogs disagree, block the skill.

## Finding Categories

- Generic copy-paste body or name-only variation.
- Missing or inconsistent versioning, changelog, or compatibility metadata.
- Non-domain checklist, decision rule, finding category, severity guidance, output, or acceptance criterion.
- Missing trigger clarity or routing ambiguity.
- Unsafe governance, secrets, release, or validation instruction.
- Generated output drift across platform copies.

## Severity Guidance

- Critical: skill instructs unsafe actions, fabricates validation, leaks secrets, or falsely claims production compatibility.
- High: skill is generic enough to misroute or produce low-quality domain work despite valid structure.
- Medium: domain content exists but lacks testable acceptance, output artifacts, or severity precision.
- Low: wording, examples, or metadata clarity needs improvement without blocking basic use.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Complete SKILL.md body with all required sections and synchronized changelog.
- Domain-specific checklist, decision rules, finding categories, severity, outputs, acceptance, and anti-patterns.
- Compatibility metadata state and validation evidence or unverified marker.
- Generator-smell review result with any rejected generic phrases.
- Tests or validation commands that enforce structure and non-generic content.
- Generated-copy synchronization notes where applicable.

## Acceptance Criteria

- Skill cannot be reduced to name, description, and shared boilerplate.
- Every required section contains domain-specific, testable content.
- Frontmatter and body changelog are synchronized.
- Compatibility metadata is honest and centrally sourced.
- Generic phrases are absent or explicitly rejected.
- Validation/tests cover required structure and non-genericness.

## Anti-Patterns

- Creating a skill by search-and-replace from another skill.
- Using “structured analysis or review” as a trigger.
- Writing finding categories like “missing <skill> evidence”.
- Using generic severity guidance unrelated to domain impact.
- Claiming production-ready because required headings exist.
- Setting concrete platform versions without validation evidence.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
