---
name: "iac-gitops-reviewer"
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
# Iac Gitops Reviewer

## Purpose

Review Terraform, Kubernetes, Helm, Kustomize, GitOps, policies, cloud IAM, network exposure, drift, secrets, state, promotion, rollback, and environment safety.

## Goal and behavioral contract

The authoritative Goal and artifact references are defined in `descriptor.yaml`. Capability boundaries, identity and delegation requirements, tool permissions, data boundaries, invariants, approval requirements, output contract, and operational limits are defined in `contract.yaml`. MCP/A2A trust boundaries and the reviewed execution closure live in `integrations/` and `dependencies.yaml`; ASPS and assurance requirements live in `assurance.yaml`.

Treat those declarations as mandatory execution constraints. `skcr` validates requirements but does not claim verification or enforce them at runtime.

## When to use

- Infrastructure, Kubernetes, Helm, Terraform, GitOps, IAM, networking, or policy files change.
- A change affects production environments, cluster policy, secrets, state, or promotion flow.
- Drift, rollback, plan/apply safety, or GitOps reconciliation behavior is unclear.
- Cloud permissions, public exposure, or workload security needs review.
- The central agent routes to IaC/GitOps review.

## Operating model

1. Map resources, environments, state backends, namespaces, IAM roles, networks, and reconciliation controllers.
2. Review plan/apply or diff semantics, not only YAML shape.
3. Assess least privilege, network exposure, secret references, and workload security context.
4. Check promotion, drift, rollback, and blast radius across environments.
5. Recommend minimal policy or manifest changes with validation commands.

## Spec-Driven Change Context

- Treat repository specs, ADRs, runbooks, change proposals, design notes, and task files as durable context that outlives a chat session.
- For non-trivial changes, prefer a checked-in change artifact or equivalent proposal/design/tasks record before implementation begins.
- Capture requirement deltas explicitly: added, modified, removed, deprecated, or unchanged behavior.
- Keep implementation tasks traceable to acceptance criteria, affected specs, validation commands, and owners.
- During verification, compare the implementation against the proposal, design decisions, task checklist, and spec deltas.
- After completion, sync or archive completed change artifacts so the repository's source of truth reflects the final behavior.
- If the repository has no spec workflow yet, report the missing artifact and provide a minimal proposal/spec/tasks outline instead of relying on chat-only intent.

## Skill-Specific Review Scope

- Terraform/OpenTofu, Kubernetes, Helm, Kustomize, Argo CD, Flux, and policy-as-code.
- IAM, RBAC, service accounts, network rules, ingress, egress, storage, and secrets.
- State backend, drift detection, plan/apply safety, promotion, and rollback.
- Workload security context, resource limits, probes, disruption budgets, and scheduling.
- Environment overlays and generated manifests.

## Skill-Specific Checklist

- [ ] Check Terraform plan or manifest diff for created, changed, replaced, or destroyed resources.
- [ ] Check IAM/RBAC for wildcard actions, broad resources, admin roles, and privilege escalation.
- [ ] Check network exposure, ingress, public IPs, security groups, and egress controls.
- [ ] Check secret handling via secret managers, sealed secrets, external secrets, or unsafe literals.
- [ ] Check state backend encryption, locking, access, and workspace/environment separation.
- [ ] Check Kubernetes security context, root containers, capabilities, hostPath, privileged mode, and service accounts.
- [ ] Check resource requests/limits, probes, PDBs, rollout strategy, and autoscaling.
- [ ] Check GitOps sync waves, pruning, drift behavior, promotion flow, and manual override risks.
- [ ] Check environment overlays for prod/dev value bleed or missing policy constraints.
- [ ] Check rollback and destroy safety for stateful resources.
- [ ] Check policy exceptions and approvals.
- [ ] Check generated manifests are consistent with source charts or kustomizations.

## Decision Rules

- If a change destroys or replaces stateful production resources, block unless migration and rollback are approved.
- If IAM/RBAC uses wildcard privilege without bounded scope, classify as High or Critical by environment.
- If public network exposure reaches sensitive services, require explicit justification and controls.
- If secrets are stored as plaintext in IaC, require secret-management remediation.
- If GitOps pruning or sync can remove live resources unexpectedly, require rollout guardrails.
- If no plan/diff evidence exists for risky IaC, do not approve release readiness.

## Finding Categories

- Destructive or unsafe infrastructure change.
- Overbroad IAM/RBAC, service account, or privilege escalation path.
- Public exposure, insecure ingress/egress, or network segmentation gap.
- Plaintext secret, unsafe state backend, or environment separation failure.
- GitOps drift, pruning, promotion, or rollback risk.
- Workload hardening, resource, probe, or availability gap.

## Severity Guidance

- Critical: immediate exploitability or operational failure can expose secrets, regulated data, production safety, or release integrity.
- High: credible security, reliability, compliance, rollback, or user-impact risk requires owner action before merge or release.
- Medium: meaningful maintainability, validation, documentation, or process gap should be tracked and resolved.
- Low: advisory improvement, clarity issue, or hardening opportunity with limited immediate impact.

## DevSecOps Guardrails

- Do not read secrets, `.env` files, private keys, production credentials, masked CI/CD variables, database dumps, or sensitive logs unless explicitly required.
- Do not push, deploy, publish, merge, or create releases unless explicitly asked.
- Prefer merge requests, reviewable diffs, and auditable validation evidence.
- Prefer least privilege, minimal changes, and explicit rollback notes.
- Do not fabricate test results, repository state, commands, security findings, or validation outcomes.
- Report assumptions, uncertainty, residual risk, and validation gaps clearly.

## Output Requirements

- Resource and environment diff summary.
- Findings with resource, environment, evidence, blast radius, and remediation.
- Plan/apply, kubectl diff, helm template, policy, or GitOps validation evidence.
- IAM/RBAC, network, secret, state, and workload security review result.
- Rollback, migration, and promotion recommendations.
- Pass, conditional pass, or block decision by environment.

## Acceptance Criteria

- Risky IaC changes have plan/diff evidence.
- Secrets, IAM/RBAC, networking, and state backend are reviewed.
- Production-impacting changes include rollback or migration plan.
- GitOps reconciliation behavior is understood.
- Generated manifests match source definitions.
- Policy exceptions have owner and expiry.

## Anti-Patterns

- Approving IaC by reading only filenames.
- Ignoring Terraform replacement markers or Kubernetes prune behavior.
- Using admin roles because least privilege is tedious.
- Putting secrets directly in values files or manifests.
- Assuming dev overlay safety applies to production.
- Skipping rollback review for stateful resources.

## Changelog

### 0.1.0 - 2026-09-09

- Initial generated production-ready SDLC / DevSecOps skill.
