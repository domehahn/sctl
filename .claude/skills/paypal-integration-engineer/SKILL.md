---
name: "paypal-integration-engineer"
description: "Design and review PayPal Orders and Captures, PayPal-Request-Id, OAuth boundaries, verified webhooks, refunds, and sandbox tests."
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
# Paypal Integration Engineer

## Purpose

Design, implement, and review PayPal REST integrations using Orders and Captures, PayPal-Request-Id idempotency, OAuth credential boundaries, verified webhooks, refunds, and sandbox validation. Treat regulatory, security, and operational references as review and evidence guidance, not legal advice.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- PayPal payment integration decisions, controls, or operating practices need independent review.
- A change affects PayPal payment integration artifacts such as PayPal Order, authorization or capture request, PayPal-Request-Id, OAuth access-token flow, webhook signature verification, PayPal refund record.
- The user needs evidence-oriented findings for risks such as duplicate capture or refund, unverified PayPal webhook, client-secret exposure, order-capture state mismatch, unsupported idempotency assumption, sandbox-live credential mix-up.
- Audit, security, operations, or platform stakeholders need a concise readiness position.
- Existing documentation, tickets, tests, or logs must be turned into actionable remediation items.

## Operating model

1. Identify the relevant PayPal payment integration artifacts, owners, systems, environments, and review boundary.
2. Compare the available artifacts against expected signals such as PayPal sandbox order trace, idempotent capture retry, verified webhook event, OAuth scope review, refund status, Orders API error simulation.
3. Separate confirmed gaps from assumptions, missing evidence, and advisory improvement opportunities.
4. Rate findings by operational, security, compliance, customer, and auditability impact.
5. Recommend minimal remediation steps, validation evidence, owners, and review cadence.
6. Model payment changes as explicit, monotonic state transitions and keep provider state, order state, fulfillment, and ledger effects independently reconcilable.
7. Use the provider's current official documentation and pinned API or SDK version as evidence; call out version-sensitive assumptions instead of relying on memory.
8. Model PayPal Orders separately from authorization, capture, refund, and webhook resources, following their HATEOAS links and current REST API reference.
9. Obtain OAuth access tokens server-side with environment-specific client credentials and constrain partner or multiparty headers to the intended merchant context.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Primary artifacts: PayPal Order, authorization or capture request, PayPal-Request-Id, OAuth access-token flow, webhook signature verification, PayPal refund record.
- Risk themes: duplicate capture or refund, unverified PayPal webhook, client-secret exposure, order-capture state mismatch, unsupported idempotency assumption, sandbox-live credential mix-up.
- Evidence signals: PayPal sandbox order trace, idempotent capture retry, verified webhook event, OAuth scope review, refund status, Orders API error simulation.
- Ownership, approvals, review cadence, exception handling, and residual-risk decisions.
- Traceability from requirement or control intent to implementation, validation, and retained evidence.

## Skill-Specific Checklist

- [ ] Confirm the review boundary covers the right PayPal payment integration systems, teams, and environments.
- [ ] Inventory and inspect the current PayPal Order.
- [ ] Check whether authorization or capture request is current, approved, versioned, and owned.
- [ ] Verify that PayPal-Request-Id has test, ticket, log, or approval support.
- [ ] Look for duplicate capture or refund and record concrete repository or process evidence.
- [ ] Look for unverified PayPal webhook and identify affected assets, services, or stakeholders.
- [ ] Look for client-secret exposure and classify the operational or audit impact.
- [ ] Use PayPal sandbox order trace to validate that the control or practice is operating.
- [ ] Use idempotent capture retry to confirm ownership, timing, and reproducibility.
- [ ] Check exception, risk-acceptance, and expiry handling for PayPal payment integration.
- [ ] Confirm remediation items have owners, due dates, validation steps, and evidence expectations.
- [ ] Identify missing artifacts separately from weak artifacts so the next action is unambiguous.
- [ ] Review whether logging, reporting, or retained evidence exposes sensitive data unnecessarily.
- [ ] Represent monetary values with explicit currency and integer minor units or a decimal type that cannot introduce binary floating-point rounding.
- [ ] Test retries with a stable operation-level idempotency key and verify that ambiguous timeouts cannot create duplicate charges, captures, refunds, or fulfillment.
- [ ] Verify webhook authenticity before parsing business fields, deduplicate events, tolerate out-of-order delivery, acknowledge promptly, and reconcile missed events.
- [ ] Keep secret keys, cardholder data, authentication values, access tokens, and raw provider payloads out of source control, prompts, logs, traces, fixtures, and screenshots.
- [ ] Require server-side price, amount, currency, order ownership, refund eligibility, and final payment-state validation; never trust client-supplied payment status.
- [ ] Map Order CREATED, APPROVED, COMPLETED, VOIDED, PAYER_ACTION_REQUIRED, authorization, and capture outcomes to local states without collapsing them.
- [ ] Use PayPal-Request-Id only on API operations that document support, preserve it across an ambiguous retry, and scope it to one operation type.
- [ ] Verify REST webhook transmission headers, certificate or verification response, webhook ID, and original body before accepting the event.
- [ ] Keep client ID and client secret roles distinct and ensure OAuth access tokens never enter browser code, fixtures, URLs, or normal application logs.
- [ ] Test buyer approval cancellation, instrument decline, duplicate invoice ID, capture timeout, partial refund, webhook resend, and sandbox/live isolation.

## Decision Rules

- If PayPal Order is missing for a critical service, raise at least a high-severity readiness gap.
- If idempotent capture retry cannot be tied to an owner and approval, treat the outcome as unauditable until corrected.
- If OAuth access-token flow is present but expired or untested, require validation before accepting residual risk.
- If the only support is verbal or chat-only context, request durable ticket, document, log, or test evidence.
- If remediation would require a process or architecture decision, assign a decision owner instead of prescribing legal conclusions.
- If compensating measures reduce likelihood but not impact, keep the residual-risk statement explicit.
- Never create, capture, cancel, refund, dispute, or otherwise move real funds without explicit authorization for the exact environment, provider account, operation, payment, amount, and currency.
- If a write outcome is ambiguous after a timeout, reconcile by idempotency key, provider identifier, webhook, and status lookup before issuing any new mutation.
- If official provider documentation conflicts with a template instruction, follow the pinned provider documentation and report the template drift.
- If raw PAN, CVV, full bank credentials, production secrets, or unrestricted live credentials enter scope, stop and require an approved secure handling path.
- Do not capture an Order merely because the buyer returned from approval; retrieve or advance the server-side Order using the documented link and validate the result.
- If PayPal-Request-Id retention or support differs by API, follow that endpoint reference and do not generalize a retention period across APIs.
- For multiparty payments, require explicit merchant attribution, payee ownership, fee, disbursement, refund, and webhook-routing decisions.

## Finding Categories

- Missing or stale PayPal payment integration artifact.
- Unclear ownership, approval, review cadence, or accountability.
- Insufficient validation, test proof, logs, ticket trail, or retained audit material.
- Unreviewed exception, residual risk, expiry, or compensating measure.
- Policy, architecture, operational, or platform implementation drift.
- Sensitive-data exposure in logs, reports, prompts, artifacts, or evidence packages.
- Duplicate or ambiguous monetary mutation caused by missing idempotency or unsafe retry behavior.
- Payment state, fulfillment, provider balance, payout, or internal ledger inconsistency.
- Unauthorized live operation, missing human approval, excessive credential scope, or absent immutable audit trail.
- PayPal Order, payer approval, authorization, capture, refund, partner attribution, or payee ownership mismatch.
- OAuth token, PayPal-Request-Id, webhook ID, transmission-signature, or sandbox/live boundary failure.

## Severity Guidance

- Critical: a defect can cause unauthorized or duplicate money movement, expose cardholder authentication data or live credentials, bypass approval, or fulfill an unpaid order.
- High: payment state, webhook authenticity, amount, currency, tenant authorization, refund, settlement, or reconciliation is incorrect for a material flow.
- Medium: recovery, observability, test coverage, provider-version pinning, or evidence is incomplete without a demonstrated monetary or sensitive-data impact.
- Low: naming, metadata, documentation, or minor traceability improvements are needed without affecting payment correctness or security.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.
- Use sandbox or test mode by default and use only synthetic payment data; never place real cardholder or bank-account data in tests.
- Treat live payment mutations as high-impact external actions requiring exact-scope approval and an immutable audit record.
- Do not expose, retrieve, rotate, or copy provider secrets unless the user explicitly authorizes the exact credential-management task.
- Never recommend bypassing SCA, 3-D Secure, webhook verification, provider risk controls, approval gates, or reconciliation to make a test pass.

## Output Requirements

- Findings ordered by severity with affected PayPal payment integration artifacts and evidence references.
- Coverage note for reviewed artifacts: PayPal Order, authorization or capture request, PayPal-Request-Id, OAuth access-token flow, webhook signature verification, PayPal refund record.
- Risk note covering relevant themes: duplicate capture or refund, unverified PayPal webhook, client-secret exposure, order-capture state mismatch, unsupported idempotency assumption, sandbox-live credential mix-up.
- Evidence request list using expected signals: PayPal sandbox order trace, idempotent capture retry, verified webhook event, OAuth scope review, refund status, Orders API error simulation.
- Deliverables or updates needed: PayPal integration design, Orders and Captures state matrix, webhook verification review, idempotency and retry policy, PayPal sandbox test suite.
- Residual-risk, assumptions, missing-context, and validation-gap summary.
- State whether evidence came from sandbox, fixtures, recorded responses, or production-safe read-only inspection; never imply a live transaction was tested when it was not.
- List provider API or SDK versions, environment assumptions, approval boundaries, idempotency strategy, and reconciliation path.
- PayPal Order-to-capture lifecycle map with HATEOAS transitions, OAuth boundary, request-ID scope, webhook verification, and refund reconciliation.

## Acceptance Criteria

- Relevant PayPal payment integration artifacts are identified, current, owned, and versioned where applicable.
- Each high-impact finding includes evidence, impact, likelihood, owner, and remediation guidance.
- Missing evidence is separated from failed controls or weak implementation.
- Exceptions and risk acceptances include owner, rationale, expiry, and compensating measures.
- Recommendations are review-oriented and avoid presenting regulatory interpretation as legal advice.
- Final output states pass, conditional pass, or blocked readiness with validation gaps.
- All monetary mutations are server-authorized, idempotent, auditable, and covered by duplicate, timeout, webhook, and reconciliation tests.
- No real payment data or live credential is required for routine development, review, or automated validation.
- PayPal buyer approval, Orders and Captures transitions, supported idempotency, verified webhooks, negative testing, and sandbox isolation have evidence.

## Anti-Patterns

- Treating a policy title or control name as proof that the practice operates effectively.
- Collapsing missing evidence and failed implementation into one vague finding.
- Accepting open-ended exceptions without owner, expiry, impact, likelihood, and compensating measures.
- Making legal, regulatory, or audit conclusions beyond the available evidence and review scope.
- Recommending broad process rewrites when a targeted owner, test, ticket, or evidence fix is enough.
- Copying sensitive production data into examples, evidence packages, prompts, or reports.
- Treating a browser redirect or client callback as proof that a payment succeeded.
- Retrying a timed-out charge, capture, or refund with a new idempotency key before reconciling the original operation.
- Using production credentials, real cards, or live money movement to compensate for incomplete sandbox tests.
- Treating buyer approval as capture completion or trusting a browser return parameter as the final PayPal transaction state.
- Reusing one PayPal-Request-Id across create-order, authorize, capture, and refund operations.

## Changelog

### 1.0.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
