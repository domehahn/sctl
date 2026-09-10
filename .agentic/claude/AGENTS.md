# Claude Code Integration

The Claude Code adapter generates both skills and subagents.

## Generated structure

```text
CLAUDE.md
.claude/
├── skills/
│   └── <skill-name>/
│       └── SKILL.md
└── agents/
    └── <subagent-name>.md
Skills vs. Subagents

Skills are reusable instruction bundles.

.claude/skills/<skill-name>/SKILL.md

Subagents are specialized Claude Code workers with their own context, tool access, model, permissions, and system prompt.

.claude/agents/<subagent-name>.md
Generated subagents

The DevSecOps SDLC edition generates these Claude Code subagents when the claude platform is enabled:

requirements-analyst
architecture-reviewer
devsecops-reviewer
security-reviewer
ci-cd-reviewer
iac-gitops-reviewer
test-runner
release-readiness-reviewer
incident-postmortem-assistant
Invocation

Claude Code can delegate automatically when a task matches a subagent description.

You can also ask explicitly:

Use the security-reviewer subagent to review this change.

or:

Use the test-runner subagent to run relevant tests and summarize failures.
Design rules
Keep review subagents mostly read-only.
Use explicit tools.
Use skills: to preload relevant skill instructions into subagent context.
Prefer bounded, specialized subagents over generic workers.
Do not use subagents as a replacement for human review on security-sensitive changes.

---

## README-Ergänzung

Diesen Block in die Plattformübersicht deiner `README.md` einfügen:

```md
### Claude Code

Generated files:

```text
CLAUDE.md
.claude/skills/<skill-name>/SKILL.md
.claude/agents/<subagent-name>.md

Claude Code skills provide reusable instructions. Claude Code subagents provide specialized workers with their own context, tools, model, permissions, and preloaded skills.

The Claude adapter generates both:

skills under .claude/skills/<skill-name>/SKILL.md
subagents under .claude/agents/<subagent-name>.md

Subagents are useful for bounded review, testing, analysis, and DevSecOps tasks that would otherwise flood the main Claude Code conversation with logs, grep output, test output, or detailed review findings.


---

## `src/agentic_template_kit/templates/claude/CLAUDE.md.j2`

```md
# CLAUDE.md

## Claude Code Instructions

Project: skpm  
Owner team: platform-engineering  
Governance level: standard

## Required behavior

- Make real file changes when implementation is requested.
- Do not only explain how to do the work.
- Do not commit, push, deploy, publish, or merge unless explicitly asked.
- Avoid secrets and sensitive files.
- Keep changes minimal, validated, and reviewable.
- Summarize changed files, validation results, and residual risks.

## Project Skills

Claude Code project skills are stored under:

```text
.claude/skills/<skill-name>/SKILL.md

Use slash commands where supported, or apply the matching skill instructions directly.



/requirements-analyst — Analyze requirements, user stories, acceptance criteria, constraints, risks, and open questions before implementation.


/cost-based-planner — Plan coding work with minimal context, relevant file selection, risk awareness, rollback, and validation strategy.


/architecture-reviewer — Review architecture, module boundaries, interfaces, coupling, scalability, data flows, and technical risks.


/threat-modeler — Identify assets, trust boundaries, abuse cases, attack paths, threats, and required security controls.


/safe-implementer — Create or modify code, tests, configuration, and project files safely with real file changes.


/test-strategy-engineer — Design and generate unit, integration, regression, security, and end-to-end test strategies.


/verification-reviewer — Review diffs, validate acceptance criteria, inspect test results, and find missed requirements.


/security-reviewer — Review code, CI/CD, configuration, permissions, dependencies, input validation, and DevSecOps risks.


/secrets-reviewer — Detect and prevent exposure of secrets, tokens, credentials, private keys, CI variables, and sensitive logs.


/dependency-supply-chain-reviewer — Review dependencies, lockfiles, package managers, container images, actions, and supply-chain risks.


/ci-cd-reviewer — Review CI/CD pipelines, runners, permissions, artifacts, caches, deployment gates, and token exposure.


/iac-gitops-reviewer — Review Terraform, Kubernetes, Helm, Kustomize, GitOps reconciliation, promotion, and environment safety.


/compliance-governance-reviewer — Review governance controls such as CODEOWNERS, branch protection, approvals, auditability, and policy compliance.


/release-readiness-reviewer — Assess release readiness, rollback, migrations, feature flags, monitoring, documentation, and breaking changes.


/observability-reviewer — Review logging, metrics, tracing, health checks, alerts, dashboards, runbooks, and operational readiness.


/incident-postmortem-assistant — Support incident analysis, timeline creation, root cause analysis, impact assessment, corrective actions, and follow-up issues.


/documentation-maintainer — Create and update README files, ADRs, setup guides, API docs, runbooks, and operational documentation.


/universal-skill-creator — Create, adapt, validate, and optimize reusable agent skills across agentic platforms.


/dora-readiness-reviewer — Review DORA readiness for ICT risk management, resilience testing, incidents, third-party risk, roles, policies, evidence, and auditability.


/ict-risk-management-reviewer — Review ICT risks, protection needs, criticality, controls, residual risks, treatment, and recurring reassessment.


/ict-third-party-risk-reviewer — Review cloud, SaaS, outsourcing, subcontractors, contracts, exit strategies, concentration risks, and DORA information-register readiness.


/ict-incident-reporting-reviewer — Review ICT incident classification, escalation, documentation, reportability, timelines, responsibilities, templates, and communication chains.


/operational-resilience-tester — Review backup and restore, failover, disaster recovery, restart procedures, crisis exercises, scenario tests, and lessons learned.


/audit-evidence-reviewer — Review evidence, approvals, tickets, logs, test protocols, risk decisions, versioning, and accountable owners.


/control-mapping-reviewer — Map technical measures to DORA, VAIT or BAIT migration needs, ISO 27001, BSI, internal policies, or MaRisk review expectations.


/outsourcing-exit-strategy-reviewer — Review exit plans, data return, provider transitions, emergency operations, suboutsourcing, cloud dependencies, and business impact.


/documentation-governance-reviewer — Review documentation freshness, ownership, review cycles, approvals, versioning, validity, and traceability.


/runbook-playbook-maintainer — Create and review runbooks, operating instructions, incident playbooks, escalation paths, restart procedures, and checklists.


/architecture-decision-recorder — Create and maintain ADRs with context, decisions, alternatives, risks, security impact, compliance relation, and review points.


/audit-traceability-maintainer — Link requirements, controls, implementation, tests, tickets, and evidence into an auditable trace.


/policy-documentation-maintainer — Create and update policies, standards, procedures, and control descriptions.


/evidence-package-creator — Create auditable evidence packages from tickets, pipeline results, test reports, approvals, scans, and architecture information.


/devsecops-maturity-reviewer — Assess maturity across plan, code, build, test, release, deploy, and operate with automation, security gates, ownership, and feedback loops.


/pipeline-security-architect — Design and review secure CI/CD pipelines with isolated runners, minimal rights, OIDC, signed artifacts, protected environments, and approval gates.


/software-supply-chain-architect — Review SLSA, provenance, SBOM, signatures, attestations, build integrity, artifact promotion, and trusted builders.


/policy-as-code-engineer — Create and review policies for OPA/Rego, Kyverno, GitLab Policies, Conftest, Checkov, Terraform, Kubernetes, and CI/CD gates.


/secure-developer-platform-reviewer — Review Internal Developer Platforms for secure golden paths, self-service guardrails, templates, permission models, secrets handling, and auditability.


/vulnerability-management-coordinator — Assess CVE triage, prioritization, SLAs, exploitability, asset criticality, exceptions, risk acceptance, and remediation tracking.


/cloud-landing-zone-reviewer — Review cloud accounts or subscriptions, networks, IAM, logging, policies, baselines, guardrails, encryption, tagging, and tenant separation.


/cloud-governance-reviewer — Review cloud naming, tags, ownership, cost centers, allowed services, regions, data classification, policy enforcement, and audit evidence.


/finops-reviewer — Review cloud costs, budgets, rightsizing, reserved or committed usage, anomalies, showback or chargeback, and team cost transparency.


/sre-reliability-reviewer — Assess SLOs, SLIs, error budgets, capacity, degradation, timeouts, retries, circuit breakers, load shedding, and operational risks.


/kubernetes-platform-reviewer — Review Kubernetes clusters, namespaces, RBAC, NetworkPolicies, Pod Security, admission controllers, resource limits, secrets, ingress, tenancy, and upgrades.


/gitops-operations-reviewer — Review Argo CD or Flux setups, sync policies, drift detection, promotion, rollback, app-of-apps, secrets, cluster access, and deployment governance.


/aiops-signal-correlation-reviewer — Assess correlation of logs, metrics, traces, events, and incidents to reduce noise, improve root-cause analysis, and lower alert fatigue.


/alert-quality-reviewer — Review alerts for actionability, clear symptoms, runbook links, severity, ownership, SLO relation, deduplication, escalation, and remediation suitability.


/auto-remediation-reviewer — Review automated repair actions for safe limits, dry runs, approval modes, rollback, audit logs, blast radius, and loop protection.


/mlops-governance-reviewer — Review model versioning, training data, bias, drift, monitoring, approvals, reproducibility, model registry, and deployment gates.


/llmops-security-reviewer — Review GenAI workloads for prompt injection, tool permissions, data exfiltration, RAG sources, sensitive prompt logging, evals, guardrails, and model access.


/ai-change-risk-reviewer — Review AI-assisted changes before execution for automation boundaries, human approval, affected-system criticality, and audit evidence.


/agent-containment-reviewer — Review agent sandboxes, transitive egress, privilege escalation, lateral movement, shared infrastructure, monitoring, and kill switches.


/agent-runtime-enforcement-reviewer — Compare declared agent contracts with independent runtime enforcement across tools, files, processes, networks, secrets, approvals, and limits.


/agent-behavior-eval-engineer — Design trajectory-level agent security evals for Goal compliance, contract boundaries, goal hacking, unsafe tools, and long-horizon behavior.


/backdoor-persistence-reviewer — Review changes for hidden privileged paths, triggers, covert egress, persistence, security-control tampering, and unexplained behavior.


/agentic-threat-modeler — Threat-model agents as untrusted principals across prompts, tools, MCP, memory, delegation, runtime infrastructure, and external systems.


/security-invariant-test-engineer — Derive negative tests and declarative evals from contract capabilities, tools, data flows, approvals, limits, and structured invariants.


/privacy-data-protection-reviewer — Review privacy, personal data, data classification, deletion concepts, purpose limitation, GDPR risks, and sensitive-data logging.


/api-contract-reviewer — Review REST, GraphQL, OpenAPI, and gRPC contracts, breaking changes, versioning, AuthN/AuthZ, error formats, and compatibility.


/secure-design-reviewer — Review secure-by-design decisions, least privilege, Zero Trust, tenant separation, secure defaults, and abuse scenarios.


/policy-as-code-reviewer — Review GitLab Security Policies, OPA/Rego, Kyverno, Conftest, Sentinel, admission policies, compliance pipelines, and central guardrails.


/container-security-reviewer — Review Dockerfiles, base images, user rights, capabilities, SBOM, image signing, distroless or slim images, CVEs, and runtime hardening.


/identity-access-reviewer — Review IAM, roles, service accounts, groups, tokens, OIDC federation, GitLab or GitHub permissions, cloud rights, and privilege-escalation paths.


/risk-acceptance-reviewer — Document and assess conscious risk decisions, impact and likelihood, expiry dates, and compensating measures.


/secure-code-reviewer — Review code vulnerabilities such as injection, path traversal, SSRF, XSS, deserialization, crypto misuse, and race conditions.


/performance-scalability-reviewer — Review load behavior, bottlenecks, caching, database access, queue behavior, scaling, timeouts, and resource limits.


/migration-change-reviewer — Review database migrations, schema changes, breaking changes, rollback ability, backward compatibility, and zero-downtime deployments.


/sbom-vulnerability-management-reviewer — Review SBOM generation, CVE triage, VEX, exception processes, patch SLAs, and the vulnerability lifecycle.


/developer-experience-reviewer — Review setup, local development, error messages, Makefiles or scripts, onboarding, tooling consistency, and practicality for teams.


/resilience-reviewer — Review timeouts, retries, circuit breakers, failover, backpressure, degraded modes, and resilience behavior.


/backup-restore-reviewer — Review restore tests, RPO/RTO, data integrity, backup protection, recoverability, and disaster recovery.


/payment-integration-engineer — Design and implement provider-neutral payment integrations with explicit state machines, idempotency, failure recovery, and auditable order linkage.


/payment-security-reviewer — Review payment data flows, tokenization, credential boundaries, PCI DSS scope, authorization, sensitive logging, and checkout abuse cases.


/payment-webhook-reviewer — Review payment webhooks for signature verification, replay and duplicate handling, ordering, durable processing, retries, and reconciliation.


/payment-flow-tester — Test authorization, capture, settlement, cancellation, refund, timeout, decline, authentication, webhook delay, and provider outage paths.


/refund-dispute-handler — Guide controlled refunds, reversals, chargebacks, and disputes with eligibility checks, evidence, deadlines, approval, and audit trails.


/payment-reconciliation-reviewer — Reconcile orders, provider transactions, fees, refunds, chargebacks, settlements, payouts, and ledger entries.


/subscription-billing-engineer — Design subscription billing for plans, trials, invoices, proration, usage, renewals, dunning, cancellation, and entitlements.


/sca-3ds-reviewer — Review Strong Customer Authentication and 3-D Secure challenge, exemption, liability, fallback, accessibility, and state handling.


/payment-fraud-risk-reviewer — Review payment fraud signals, velocity controls, risk rules, step-up actions, manual review, false positives, and feedback loops.


/payment-observability-reviewer — Review payment conversion, authorization, webhook, refund, reconciliation, latency, provider health, alerts, and privacy-safe telemetry.


/payment-provider-migration-reviewer — Review payment-provider migrations for parity, token portability, dual processing, routing, reconciliation, rollback, and decommissioning.


/payment-compliance-reviewer — Review payment controls and evidence for PCI DSS, privacy, SCA, retention, access, auditability, outsourcing, and exceptions.


/payment-operations-agent — Support controlled capture, cancellation, refund, resend, and lookup operations with approval, limits, idempotency, and audit records.


/stripe-integration-engineer — Design and review Stripe PaymentIntents or Checkout, idempotency, signed webhooks, Connect or Billing boundaries, versions, and tests.


/paypal-integration-engineer — Design and review PayPal Orders and Captures, PayPal-Request-Id, OAuth boundaries, verified webhooks, refunds, and sandbox tests.


/adyen-integration-engineer — Design and review Adyen Checkout, merchant references, idempotency, HMAC webhooks, asynchronous results, modifications, and tests.


/java-reviewer — Review modern Java code, JVM behavior, concurrency, APIs, testing, performance, and maintainability.


/golang-reviewer — Review Go code for idioms, concurrency, context propagation, errors, interfaces, modules, tests, and performance.


/python-reviewer — Review Python code for typing, packaging, async behavior, resource safety, tests, security, and maintainability.


/ruby-reviewer — Review Ruby code for idioms, object design, metaprogramming boundaries, Bundler hygiene, tests, and performance.


/javascript-reviewer — Review modern JavaScript for async behavior, modules, runtime correctness, security, tests, and performance.


/typescript-reviewer — Review TypeScript strictness, narrowing, generics, declarations, runtime validation, and compiler configuration.


/rust-reviewer — Review Rust ownership, lifetimes, unsafe boundaries, concurrency, error handling, Cargo dependencies, and tests.


/csharp-reviewer — Review C# and .NET code for async correctness, dependency injection, resource disposal, nullable types, tests, and performance.


/kotlin-reviewer — Review Kotlin code for null safety, coroutines, sealed models, Java interop, Gradle configuration, and tests.


/php-reviewer — Review modern PHP code for type safety, Composer hygiene, framework boundaries, security, tests, and performance.


/spring-boot-reviewer — Review Spring Boot dependency injection, transactions, persistence, security, configuration, Actuator, and migrations.


/quarkus-reviewer — Review Quarkus build-time behavior, CDI, reactive paths, native images, configuration, security, and tests.


/angular-reviewer — Review Angular components, signals, RxJS, forms, routing, state, accessibility, testing, and bundle performance.


/react-reviewer — Review React components, hooks, state, rendering, accessibility, Server Components, testing, and performance.


/vuejs-reviewer — Review Vue.js Composition API, reactivity, components, state, routing, accessibility, tests, and performance.


/nodejs-reviewer — Review Node.js services for event-loop safety, async behavior, modules, streams, HTTP security, dependencies, and operations.


/nextjs-reviewer — Review Next.js routing, Server and Client Components, caching, data fetching, security, deployment, and performance.


/django-reviewer — Review Django models, migrations, ORM usage, authentication, middleware, settings, tests, and deployment safety.


/fastapi-reviewer — Review FastAPI schemas, dependency injection, async paths, authentication, validation, OpenAPI, tests, and operations.


/ruby-on-rails-reviewer — Review Rails models, controllers, jobs, migrations, Active Record behavior, security, tests, and deployment safety.


/aws-cloud-reviewer — Review AWS accounts, IAM, networking, compute, storage, databases, observability, security, cost, and resilience.


/azure-cloud-reviewer — Review Azure tenants, subscriptions, identities, networks, compute, data services, Policy, monitoring, cost, and resilience.


/gcp-cloud-reviewer — Review GCP organizations, projects, IAM, networking, compute, data services, observability, cost, and resilience.


/terraform-reviewer — Review Terraform modules, providers, plans, state, drift, lifecycle, IAM, tests, and safe apply or destroy boundaries.


/opentofu-reviewer — Review OpenTofu modules, providers, plans, state, drift, lifecycle, tests, and safe apply or destroy boundaries.


/ansible-reviewer — Review Ansible inventories, roles, playbooks, idempotency, secrets, privilege escalation, testing, and rollout safety.


/vagrant-reviewer — Review Vagrant environments, providers, networking, provisioning, shared folders, reproducibility, and isolation.


/virtualization-reviewer — Review VM and virtualization architecture, images, isolation, networking, storage, snapshots, patching, and capacity.


/helm-reviewer — Review Helm charts, templates, values, dependencies, hooks, secrets, schema validation, upgrades, and rollbacks.


/cncf-1nce-member-reviewer — Review 1NCE (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-23-technologies-member-reviewer — Review 23 Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-3-shake-kcsp-reviewer — Review 3-Shake (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category.


/cncf-3-shake-member-reviewer — Review 3-Shake (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-3scale-reviewer — Review 3Scale using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-42on-member-reviewer — Review 42on (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-6wind-member-reviewer — Review 6WIND (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-6wind-virtual-service-router-vsr-reviewer — Review 6WIND Virtual Service Router (VSR) using its official documentation and repository in the Special / Certified CNFs category. 6WIND VSR is a Virtual Service Router that leverages 6WIND&#39;s core technology to deliver network services solutions into virtual and cloud environments quickly, efficiently, and cost-effectively. 6WIND VSR enables CSPs, Enterprises, and Cloud Providers to benefit from our leading high-performance network services for building efficient, flexible, and reliable networks.


/cncf-8gears-member-reviewer — Review 8gears (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-acc-ict-kcsp-reviewer — Review ACC ICT (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ACC ICT uses Kubernetes to deliver continuity for mission-critical solutions in demanding complex multi-cloud &amp; public-cloud environments. We design, build and maintain these environments and offer training.


/cncf-acc-ict-member-reviewer — Review ACC ICT (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-accenture-federal-services-member-reviewer — Review Accenture Federal Services (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-access-quality-kcsp-reviewer — Review Access Quality (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide professional Kubernetes design, deployment, management, and security services, ensuring high availability, scalability, and efficient operation to optimize your applications.


/cncf-access-quality-member-reviewer — Review Access Quality (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-accordion-reviewer — Review Accordion using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Accordion Kubernetes Platform helps you to orchestrate containerized workloads for your DevOps practices &amp; CI/CD pipelines, delivering enhanced developer productivity and accelerating time to market.


/cncf-actualyze-ai-member-reviewer — Review Actualyze AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-adidas-supporter-reviewer — Review Adidas (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-adobe-member-reviewer — Review Adobe (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-adyen-member-reviewer — Review Adyen (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-aembit-member-reviewer — Review Aembit (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aeraki-mesh-reviewer — Review Aeraki Mesh using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Aeraki Mesh allows you to manage any layer-7 traffic in a service mesh


/cncf-aerospike-member-reviewer — Review Aerospike (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-afi-ai-reviewer — Review Afi.ai using its official documentation and repository in the Runtime / Cloud Native Storage category. Afi SaaS backup for Kubernetes provides a fully managed cloud service with secure offsite storage and an easy-to-use web UI


/cncf-afi-technologies-member-reviewer — Review Afi Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-agenda-kcsp-reviewer — Review Agenda (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide expert consulting, implementation, and support services for Kubernetes to help organizations streamline operations and  accelerate innovation.


/cncf-agenda-member-reviewer — Review Agenda (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-agent-development-kit-reviewer — Review Agent Development Kit using its official documentation and repository in the AI Agent / Agent Framework category. An open-source, code-first Python toolkit for building, evaluating, and deploying sophisticated AI agents with flexibility and control.


/cncf-agent-evaluation-reviewer — Review Agent Evaluation using its official documentation and repository in the AI Agent / Evaluation category. A generative AI-powered framework for testing virtual agents


/cncf-agent-field-member-reviewer — Review Agent Field (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-agent-sandbox-reviewer — Review agent-sandbox using its official documentation and repository in the AI Native Infra / Orchestration and Scheduling category. agent-sandbox enables easy management of isolated, stateful, singleton workloads, ideal for use cases like AI agent runtimes.


/cncf-agent2agent-reviewer — Review Agent2agent using its official documentation and repository in the AI Agent / Protocol category. Agent2Agent (A2A) is an open protocol enabling communication and interoperability between opaque agentic applications.


/cncf-agentgateway-reviewer — Review Agentgateway using its official documentation and repository in the AI Native Infra / Gateway category. Next Generation Agentic Proxy for AI Agents and MCP servers


/cncf-agentregistry-reviewer — Review agentregistry using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Agentregistry is a cloud-native registry for discovering, curating, and deploying MCP servers, agents, and skills across local development and Kubernetes environments.


/cncf-agileops-kcntp-reviewer — Review AgileOps (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. AgileOps is a recognized Kubernetes Certified Service Provider (KCSP) and Kubernetes and Cloud Native Training Partner (KCNTP) in Vietnam. Our certified instructors bring real-world, hands-on experience to deliver world-class Kubernetes and DevSecOps training programs.


/cncf-agileops-kcsp-reviewer — Review AgileOps (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. AgileOps delivers expert Kubernetes consulting, deployment, and optimization services to help teams build scalable, secure, and cloud-native infrastructure. We empower DevOps global teams to streamline operations and accelerate delivery with Kubernetes at the core.


/cncf-agileops-member-reviewer — Review AgileOps (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-agno-reviewer — Review Agno using its official documentation and repository in the AI Agent / Agent Framework category. Build, run, manage agentic software at scale


/cncf-agola-reviewer — Review Agola using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Open Source CI/CD platform with advanced features and architecture. Powerful, reproducible and containerized workflows (called Runs), git based workflow (integrates with all the primary git repositories like GitHub, GitLab, Gitea), restart Runs from failed tasks, user direct runs (test your local changes to a remote Agola server with just one command), distributed and high available by design and runs everywhere (Kubernetes, docker, IaaS, bare metal).


/cncf-agones-reviewer — Review Agones using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Agones is a library for hosting, running, and scaling dedicated game servers on Kubernetes.


/cncf-ahnlab-cloudmate-kcsp-reviewer — Review AhnLab CloudMate (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Cloudmate provides professional services for Kubernetes. Maximize your business with application modernization.


/cncf-ahnlab-cloudmate-member-reviewer — Review AhnLab CloudMate (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aibrix-reviewer — Review AIBrix using its official documentation and repository in the Inference / Framework category. Cost-efficient and pluggable Infrastructure components for GenAI inference.


/cncf-aikido-security-member-reviewer — Review Aikido Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-airbnb-member-reviewer — Review Airbnb (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-airlock-reviewer — Review Airlock using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-airship-reviewer — Review Airship using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-aiven-member-reviewer — Review Aiven (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-akamai-member-reviewer — Review Akamai (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-akamas-member-reviewer — Review Akamas (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-akamas-reviewer — Review Akamas using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-akana-reviewer — Review Akana using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-akka-member-reviewer — Review Akka (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-akka-reviewer — Review Akka using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-akri-reviewer — Review Akri using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. A Kubernetes Resource Interface for the Edge


/cncf-aks-engine-for-azure-stack-reviewer — Review AKS Engine for Azure Stack using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. AKS Engine provides convenient tooling to quickly bootstrap Kubernetes clusters on Azure Stack. By leveraging ARM (Azure Resource Manager), AKS Engine helps create, destroy and maintain clusters provisioned with basic IaaS resources in Azure Stack.


/cncf-akuity-member-reviewer — Review Akuity (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-akuity-reviewer — Review Akuity using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Akuity is the enterprise-grade company for Argo, the open source suite of cloud native application delivery tools. Akuity was founded by                 Argo originators Hong Wang, Jesse Suen and Alexander Matyushentsev, and its mission is to empower DevOps teams to deliver their apps in a               simpler, safer, and faster way. The Akuity Platform provides a best-in-class developer experience with enterprise readiness, and enables                 organizations to modernize their toolchain for the cloud-native era.


/cncf-alasca-member-reviewer — Review ALASCA (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-alauda-container-platform-acp-reviewer — Review Alauda Container Platform (ACP) using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. A Cloud-Native Application Platform based on Kubernetes


/cncf-alauda-kcntp-reviewer — Review Alauda (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Alauda is a Cloud Native Platform provider at the forefront of empowering Enterprise IT with Continuous Innovation in the Digital Age. We provide Kubernetes training and consulting services with experienced instructors.


/cncf-alauda-kcsp-reviewer — Review Alauda (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Alauda provides Kubernetes-Centric Enterprise Platform-as-a-Service offerings with a razor focus on delivering Cloud Native capabilities and DevOps best practices to enterprise customers across industries globally.


/cncf-alauda-member-reviewer — Review Alauda (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-alcide-reviewer — Review Alcide using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Alcide provides dev-to-production security for workloads running in Kubernetes &amp; Istio


/cncf-alerant-kcsp-reviewer — Review Alerant (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Alerant delivers enterprise-grade Platform Engineering, Application Modernization, and AI solutions, leveraging AI-assisted engineering techniques. We design, deploy, and support Kubernetes-based platforms, accelerate DevOps and CI/CD adoption, and guide application migration and modernization in the cloud.


/cncf-alerant-member-reviewer — Review Alerant (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-algorithmia-reviewer — Review Algorithmia using its official documentation and repository in the Serverless / Hosted Platform category.


/cncf-alibaba-cloud-application-real-time-service-reviewer — Review Alibaba Cloud Application Real-Time Service using its official documentation and repository in the Observability and Analysis / Observability category. Application Real-Time Monitoring Service (ARMS) is an end-to-end Alibaba Cloud monitoring service for Application Performance Management (APM). You can quickly develop real-time business monitoring capabilities using the frontend monitoring, application monitoring, and custom monitoring features provided by ARMS.


/cncf-alibaba-cloud-container-registry-acr-reviewer — Review Alibaba Cloud Container Registry (ACR) using its official documentation and repository in the Provisioning / Container Registry category. Alibaba Cloud Container Registry (ACR) is a cloud-native artifacts management platform that helps your team build, manage, and ship containerized applications,




 also provides vulnerability analysis, global synchronization, content trust, and more functionalities out of the box.


/cncf-alibaba-cloud-container-service-for-kubernetes-reviewer — Review Alibaba Cloud Container Service for Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Alibaba Cloud Container Service for Kubernetes (ACK) is a fully-managed service compatible with Kubernetes to help users focus on their applications rather than managing container infrastructure.


/cncf-alibaba-cloud-file-storage-cpfs-reviewer — Review Alibaba Cloud File Storage CPFS using its official documentation and repository in the Runtime / Cloud Native Storage category. Alibaba Cloud File Storage CPFS is a scale-out parallel file system for Artificial intelligence(AI), High performance computing(HPC) workload while ensuring security, reliability, data efficiency and high performance.


/cncf-alibaba-cloud-file-storage-reviewer — Review Alibaba Cloud File Storage using its official documentation and repository in the Runtime / Cloud Native Storage category. Alibaba Cloud File Storage enables you to have a distributed file system with unlimited capacity and performance scaling with a single namespace,high-performance, high reliability, high availability and scalable file storage services.


/cncf-alibaba-cloud-function-compute-reviewer — Review Alibaba Cloud Function Compute using its official documentation and repository in the Serverless / Hosted Platform category. A fully hosted and serverless running environment that takes away the need to manage infrastructure such as servers and enables developers to focus on writing and uploading code.


/cncf-alibaba-cloud-kcsp-reviewer — Review Alibaba Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Alibaba Cloud Container Service for Kubernetes is a fully-managed service compatible with Kubernetes to help users focus on their  applications rather than managing container infrastructure.


/cncf-alibaba-cloud-log-service-reviewer — Review Alibaba Cloud Log Service using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-alibaba-cloud-member-reviewer — Review Alibaba Cloud (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-alibaba-cloud-serverless-app-engine-reviewer — Review Alibaba Cloud Serverless App Engine using its official documentation and repository in the Serverless / Hosted Platform category. Serverless Application Engine (SAE) is an application-oriented Serverless PaaS platform that helps PaaS users to go to the cloud without O&amp;M, pay-as-you-go, and low-threshold microservices. It combines Serverless architecture and microservice architecture perfectly.


/cncf-alibaba-cloud-serverless-workflow-reviewer — Review Alibaba Cloud Serverless Workflow using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-alluxio-reviewer — Review Alluxio using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-alpha-bravo-member-reviewer — Review Alpha Bravo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-altinity-member-reviewer — Review Altinity (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-amadeus-supporter-reviewer — Review Amadeus (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-amazee-io-member-reviewer — Review amazee.io (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-amazon-cloudwatch-reviewer — Review Amazon CloudWatch using its official documentation and repository in the Observability and Analysis / Observability category. Amazon CloudWatch is a monitoring and management service built for developers, system operators, site reliability engineers (SRE), and IT managers. CloudWatch provides you with data and actionable insights to monitor your applications, understand and respond to system-wide performance changes, optimize resource utilization, and get a unified view of operational health.


/cncf-amazon-elastic-block-store-ebs-reviewer — Review Amazon Elastic Block Store (EBS) using its official documentation and repository in the Runtime / Cloud Native Storage category. Amazon Elastic Block Store (EBS) is an easy to use, high performance block storage service designed for use with Amazon Elastic Compute Cloud (EC2) for both throughput and transaction intensive workloads at any scale.


/cncf-amazon-elastic-container-registry-ecr-reviewer — Review Amazon Elastic Container Registry (ECR) using its official documentation and repository in the Provisioning / Container Registry category. Amazon Elastic Container Registry (ECR) is a fully-managed Docker container registry that makes it easy for developers to store, manage, and deploy Docker container images.


/cncf-amazon-elastic-container-service-ecs-reviewer — Review Amazon Elastic Container Service (ECS) using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Amazon Elastic Container Service (Amazon ECS) is a highly scalable, high-performance container orchestration service that supports Docker containers and allows you to easily run and scale containerized applications on AWS.


/cncf-amazon-elastic-container-service-for-kub-350e8a41-reviewer — Review Amazon Elastic Container Service for Kubernetes (EKS) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Amazon Elastic Container Service for Kubernetes (Amazon EKS) is a managed service that makes it easy for you to run Kubernetes on AWS without needing to install and operate your own Kubernetes clusters.


/cncf-amazon-elastic-kubernetes-service-anywhe-52237633-reviewer — Review Amazon Elastic Kubernetes Service Anywhere (Amazon EKS Anywhere) using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Amazon EKS Anywhere is a new deployment option for Amazon EKS that allows customers to create and operate Kubernetes clusters on customer-managed infrastructure, supported by AWS.


/cncf-amazon-elastic-kubernetes-service-distro-1c0ca665-reviewer — Review Amazon Elastic Kubernetes Service Distro (Amazon EKS-D) using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Amazon Elastic Kubernetes Service Distro (Amazon EKS-D) is a Kubernetes distribution based on and used by Amazon Elastic Kubernetes Service (EKS) to create reliable and secure Kubernetes clusters.


/cncf-amazon-kinesis-reviewer — Review Amazon Kinesis using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Amazon Kinesis makes it easy to collect, process, and analyze real-time, streaming data so you can get timely insights and react quickly to new information.


/cncf-amazon-web-services-member-reviewer — Review Amazon Web Services (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-ambient-it-kcntp-reviewer — Review Ambient IT (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Ambient IT is a French training organization specializing in IT and digital professions.


/cncf-ambient-it-kcsp-reviewer — Review Ambient IT (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Ambient IT offers professional training as well as technical audits of your DevOps / Kubernetes infrastructures.


/cncf-ambient-it-member-reviewer — Review Ambient IT (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ambient-reviewer — Review Ambient using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-amd-member-reviewer — Review AMD (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-american-express-member-reviewer — Review American Express (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-amnic-member-reviewer — Review Amnic (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-amoniac-o-member-reviewer — Review Amoniac OÜ (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ampere-computing-member-reviewer — Review Ampere Computing (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-amutable-member-reviewer — Review Amutable (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-andes-digital-kcsp-reviewer — Review Andes Digital (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We empower teams to succeed with Kubernetes through consulting, implementation, and 24/7 support. Our solutions accelerate cloud-native adoption, improve developer productivity, and reduce operational risk.


/cncf-andes-digital-member-reviewer — Review Andes Digital (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-anova-supporter-reviewer — Review Anova (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-ansible-reviewer — Review Ansible using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-ant-financial-member-reviewer — Review Ant Financial (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-anteon-reviewer — Review Anteon using its official documentation and repository in the Observability and Analysis / Observability category. Anteon is a platform that combines effortless Kubernetes Monitoring and Performance Testing  to provide seamless observability of your K8S infrastructure using eBPF.


/cncf-antrea-reviewer — Review Antrea using its official documentation and repository in the Runtime / Cloud Native Network category. Kubernetes networking based on Open vSwitch


/cncf-anynines-member-reviewer — Review anynines (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aoe-kcsp-reviewer — Review AOE (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. AOE is your trusted partner in navigating the dynamic world of container orchestration. With a rich history dating back to Kubernetes&#39; inception in 2014, our experts bring unparalleled expertise to the table. We offer services around Kubernetes Consulting, Knowledge Transfer and Trainings as well as Operation and Implementation.


/cncf-aoe-member-reviewer — Review AOE (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aokumo-kcsp-reviewer — Review Aokumo (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Aokumo Inc. empowers enterprises to build and operate resilient, high-performance Kubernetes environments for multi-cluster architectures and AI/LLM workloads, combining advanced observability and automation.


/cncf-aokumo-member-reviewer — Review Aokumo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-apache-age-reviewer — Review Apache AGE using its official documentation and repository in the AI Agent / Knowledge Graph category. Graph database optimized for fast analysis and real-time data processing. It is provided as an extension to PostgreSQL.


/cncf-apache-brpc-reviewer — Review Apache bRPC using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category.


/cncf-apache-camel-k-reviewer — Review Apache Camel K using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-apache-carbondata-reviewer — Review Apache CarbonData using its official documentation and repository in the App Definition and Development / Database category.


/cncf-apache-druid-reviewer — Review Apache Druid using its official documentation and repository in the Data / Data Architecture category. A high performance real-time analytics database.


/cncf-apache-flink-reviewer — Review Apache Flink using its official documentation and repository in the Data / Data Architecture category. Open source stream processing framework with powerful stream- and batch-processing capabilities.


/cncf-apache-hadoop-reviewer — Review Apache Hadoop using its official documentation and repository in the App Definition and Development / Database category.


/cncf-apache-hbase-reviewer — Review Apache HBase using its official documentation and repository in the Data / Data Architecture category. Open source, distributed, versioned, column-oriented store modeled after Google&#39;s Bigtable.


/cncf-apache-ignite-reviewer — Review Apache Ignite using its official documentation and repository in the App Definition and Development / Database category.


/cncf-apache-mesos-reviewer — Review Apache Mesos using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-apache-nifi-reviewer — Review Apache NiFi using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-apache-openserverless-reviewer — Review Apache OpenServerless using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-apache-openwhisk-reviewer — Review Apache OpenWhisk using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-apache-pinot-reviewer — Review Apache Pinot using its official documentation and repository in the Data / Data Architecture category. A realtime distributed OLAP datastore.


/cncf-apache-rocketmq-reviewer — Review Apache RocketMQ using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-apache-spark-reviewer — Review Apache Spark using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-apache-storm-reviewer — Review Apache Storm using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-apache-streampipes-reviewer — Review Apache StreamPipes using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-apache-thrift-reviewer — Review Apache Thrift using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category.


/cncf-apache-tvm-reviewer — Review Apache TVM using its official documentation and repository in the Wasm / AI/Machine Learning category.


/cncf-apache-zeppelin-reviewer — Review Apache Zeppelin using its official documentation and repository in the Data / Data Science category. Web-based notebook that enables data-driven, interactive data analytics and collaborative documents with SQL, Scala and more.


/cncf-apache-zookeeper-reviewer — Review Apache Zookeeper using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category.


/cncf-ape-factory-kcsp-reviewer — Review ape factory (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As Kubernetes Experts, we design, build and operate your Kubernetes cluster in the public cloud, on-prem, or multi/hybrid cloud environments  tailored to your needs.


/cncf-ape-factory-member-reviewer — Review ape factory (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-apiclarity-reviewer — Review APIClarity using its official documentation and repository in the Provisioning / Security &amp; Compliance category. APIClarity is a cloud API observability open-source project that provides functionality to discover and monitor any API traffic that interacts with modern applications and report suspected security weaknesses or possible abuses. Once an OpenAPI spec is either uploaded or reconstructed by APIClarity, it monitors for shadow and zombie APIs, Broken Functional-Level Authorization (BFLA), Broken Object-Level Authorization (BOLA), weak authentication, sensitive data leaks and data injection risks.


/cncf-apicurio-registry-reviewer — Review Apicurio Registry using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Apicurio Registry is a runtime server system that stores a specific set of artifacts as files.


/cncf-apiiro-member-reviewer — Review Apiiro (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-apioak-reviewer — Review APIOAK using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-apisix-reviewer — Review APISIX using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-apolicy-reviewer — Review Apolicy using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Apolicy fuses security and compliance into your cloud native development pipeline. Using policy-as-code,  we automate risk, policy and remediation from source to production


/cncf-apollo-member-reviewer — Review Apollo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-apollo-reviewer — Review Apollo using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Apollo is a reliable open-source configuration management system


/cncf-appdynamics-reviewer — Review AppDynamics using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-apple-member-reviewer — Review Apple (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-application-high-availability-service-reviewer — Review Application High Availability Service using its official documentation and repository in the Observability and Analysis / Observability category. A SaaS-based service that aims to improve the high availability of your applications. Application High Availability Service features automatic application architecture discovery, high availability assessment following principles of chaos engineering, application traffic control, and service fallback.


/cncf-application-platform-for-lke-reviewer — Review Application Platform for LKE using its official documentation and repository in the Platform / PaaS/Container Service category. Application Platform for Linode Kubernetes Engine (LKE)


/cncf-applications-manager-reviewer — Review Applications Manager using its official documentation and repository in the Observability and Analysis / Observability category. Applications Manager is a multi-platform monitoring tool built for IT admins and DevOps for tracking KPIs of 100+ technologies spanning hybrid clouds, containers as well as traditional environments (web applications, servers, VMs, application servers, databases, big data stores, middleware &amp; messaging components, web servers, web services &amp; ERP suites).


/cncf-appneta-reviewer — Review AppNeta using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-appoptics-reviewer — Review AppOptics using its official documentation and repository in the Observability and Analysis / Observability category. Simple, powerful infrastructure and application monitoring


/cncf-appscale-reviewer — Review AppScale using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-appsignal-reviewer — Review AppSignal using its official documentation and repository in the Observability and Analysis / Observability category. AppSignal is a powerful APM that works out-of-the-box for Elixir, Ruby, and Node.js.


/cncf-appveyor-reviewer — Review Appveyor using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-aqua-reviewer — Review Aqua using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-aqua-security-kcsp-reviewer — Review Aqua Security (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Aqua Security helps enterprises secure their cloud native applications from development to production, whether they run using containers, serverless, or virtual machines


/cncf-aqua-security-member-reviewer — Review Aqua Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-arangodb-reviewer — Review ArangoDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-aranya-member-reviewer — Review Aranya (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-arcfra-kubernetes-engine-reviewer — Review Arcfra Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Arcfra Kubernetes Engine (AKE) is the Kubernetes service of Arcfra Enterprise Cloud Platform (AECP). It helps enterprises automate the complete production-ready Kubernetes infrastructure, simplifying deployment, management, and usage with an out-of-the-box experience.


/cncf-arcfra-member-reviewer — Review Arcfra (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-archestra-ai-member-reviewer — Review Archestra.AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-architect-reviewer — Review Architect using its official documentation and repository in the Serverless / Framework category.


/cncf-ardc-kcsp-reviewer — Review ARDC (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ARDC&#39;s Kubernetes service offering involves engaging the Australian research community in the adoption of Kubernetes to deliver cloud native research platforms. The core pillars of the service are community engagement via the ARCOS community, consultancy and support for platform research projects; training and implementations for researchers and research software engineers.


/cncf-ardc-member-reviewer — Review ARDC (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-ardc-nectar-research-cloud-magnum-service-reviewer — Review ARDC Nectar Research Cloud Magnum Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Get started with Kubernetes easily with Nectar Research Cloud Magnum Service.


/cncf-argo-reviewer — Review Argo using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Kubernetes-native tools to run workflows, manage clusters, and do GitOps right.


/cncf-argonix-reviewer — Review Argonix using its official documentation and repository in the Observability and Analysis / Observability category. Sovereign AI-powered SRE platform combining observability, cloud security posture management, FinOps and incident response for cloud-native environments.


/cncf-arkflow-reviewer — Review ArkFlow using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. High-performance Rust stream processing engine, providing powerful data stream processing capabilities, supporting multiple input/output sources and processors.


/cncf-arks-reviewer — Review Arks using its official documentation and repository in the Inference / Framework category. Arks is an end-to-end framework for managing LLM-based applications within Kubernetes cluster. It provides a robust and extensible infrastructure tailored for deploying, orchestrating, and scaling LLM inference workloads in cloud-native environments.


/cncf-arm-member-reviewer — Review Arm (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-armada-reviewer — Review Armada using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Armada is a multi-Kubernetes cluster batch job scheduler


/cncf-armo-member-reviewer — Review Armo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-armo-reviewer — Review ARMO using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-arpio-member-reviewer — Review Arpio (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-arrikto-reviewer — Review Arrikto using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-artifact-hub-reviewer — Review Artifact Hub using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-aruba-managed-kubernetes-reviewer — Review Aruba Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Aruba Managed Kubernetes is the easiest way to manage your kubernetes clusters.


/cncf-aruba-spa-member-reviewer — Review Aruba SpA (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ascendra-networks-member-reviewer — Review Ascendra Networks (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aserto-reviewer — Review Aserto using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Fine-grained, policy-based, real-time authorization for APIs and microservices. Aserto is the maintainer of the Topaz and Open Policy Containers OSS projects.


/cncf-asml-member-reviewer — Review ASML (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aspecto-reviewer — Review Aspecto using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-astron-agent-reviewer — Review Astron Agent using its official documentation and repository in the AI Agent / Workflow Orchestration category. Enterprise-grade, commercial-friendly agentic workflow platform for building next-generation SuperAgents with multi-agent orchestration capabilities.


/cncf-asus-cloud-infra-reviewer — Review ASUS Cloud Infra using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. ASUS Cloud Infra helps you bootstrap a Kubernetes cluster.


/cncf-asyncify-reviewer — Review Asyncify using its official documentation and repository in the Wasm / Tooling category.


/cncf-aternity-reviewer — Review Aternity using its official documentation and repository in the Observability and Analysis / Observability category. Cloud-native, scalable, high-definition application performance monitoring (APM) for next gen apps.


/cncf-athenz-reviewer — Review Athenz using its official documentation and repository in the Provisioning / Key Management category. Open source platform for X.509 certificate based service authentication and fine grained access control in dynamic infrastructures


/cncf-atix-kcntp-reviewer — Review ATIX (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. ATIX plans, builds &amp; runs Kubernetes clusters. We provide operational support and teach companies how to profit from using Kubernetes.


/cncf-atix-kcsp-reviewer — Review ATIX (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ATIX plans, builds &amp; runs Kubernetes clusters. We provide operational support and teach companies how to profit from using Kubernetes.


/cncf-atix-member-reviewer — Review ATIX (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-atlantis-reviewer — Review Atlantis using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Terraform Pull Request Automation for Teams


/cncf-atlassian-member-reviewer — Review Atlassian (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-atomist-reviewer — Review Atomist using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-atos-kcsp-reviewer — Review Atos (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Atos is a leading cloud native services provider, helping enterprises develop, deploy and manage applications using Kubernetes across various landing zones such as AWS, Azure, GCP or others European Cloud Service Providers.


/cncf-atos-member-reviewer — Review Atos (member) using its official documentation and repository in the CNCF Members / Silver category. Atos Group is a global leader in digital transformation — providing end-to-end cloud, cybersecurity, data &amp; AI, and infrastructure services to industries.


/cncf-attribute-member-reviewer — Review Attribute (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-audi-supporter-reviewer — Review Audi (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-augment-code-member-reviewer — Review Augment Code (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-auristor-member-reviewer — Review AuriStor (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-authentik-reviewer — Review authentik using its official documentation and repository in the Provisioning / Security &amp; Compliance category. authentik is an open-source Identity Provider that emphasizes flexibility and versatility, with support for a wide set of protocols.


/cncf-authing-reviewer — Review Authing using its official documentation and repository in the Provisioning / Key Management category.


/cncf-authkeys-supporter-reviewer — Review AuthKeys (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-authzed-member-reviewer — Review Authzed (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-autogen-reviewer — Review Autogen using its official documentation and repository in the AI Agent / Agent Framework category. A programming framework for agentic AI


/cncf-automq-reviewer — Review AutoMQ using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. AutoMQ is a cloud-native fork of Kafka by separating storage to S3. 10x cost-effective. Autoscale in seconds. Single-digit ms latency.


/cncf-autovia-member-reviewer — Review Autovia (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-avap-member-reviewer — Review AVAP (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-avap-reviewer — Review AVAP using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-avi-networks-reviewer — Review Avi Networks using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. Avi Networks drives automation and intelligence across multi-cloud environments with intent-based application services, including LB and WAF


/cncf-aviatrix-member-reviewer — Review Aviatrix (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aviatrix-reviewer — Review Aviatrix using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-avisi-ame-reviewer — Review Avisi AME using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. AME provides a managed Kubernetes platform that delivers seamless, easy, and secure production-ready cluster management across public, private, and hybrid cloud environments.


/cncf-avisi-cloud-kcsp-reviewer — Review Avisi Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Avisi Cloud provides Managed Kubernetes for enterprises and helps them with challenges around compliance, implementation and support for hosting mission critical software on Kubernetes.


/cncf-avisi-cloud-member-reviewer — Review Avisi Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-avro-reviewer — Review Avro using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category.


/cncf-awesome-information-technology-kcsp-reviewer — Review awesome information technology (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Independent consulting, development of specialized and individual Kubernetes clusters, and fully automated deployments


/cncf-awesome-information-technology-member-reviewer — Review awesome information technology (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-aws-app-mesh-reviewer — Review AWS App Mesh using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. AWS App Mesh is a service mesh that provides application-level networking to make it easy for your services to communicate with each other across multiple types of compute infrastructure. App Mesh standardizes how your services communicate, giving you end-to-end visibility and ensuring high-availability for your applications.


/cncf-aws-cloud-map-reviewer — Review AWS Cloud Map using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category. AWS Cloud Map is a cloud resource discovery service. With Cloud Map, you can define custom names for your application resources, and it maintains the updated location of these dynamically changing resources. This increases your application availability because your web service always discovers the most up-to-date locations of its resources.


/cncf-aws-cloudformation-reviewer — Review AWS CloudFormation using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. AWS CloudFormation provides a common language for you to describe and provision all the infrastructure resources in your cloud environment.


/cncf-aws-codepipeline-reviewer — Review AWS CodePipeline using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. AWS CodePipeline is a fully managed continuous delivery service that helps you automate your release pipelines for fast and reliable application and infrastructure updates. CodePipeline automates the build, test, and deploy phases of your release process every time there is a code change, based on the release model you define.


/cncf-aws-lambda-reviewer — Review AWS Lambda using its official documentation and repository in the Serverless / Hosted Platform category. AWS Lambda lets you run code without provisioning or managing servers. You pay only for the compute time you consume - there is no charge when your code is not running.


/cncf-aws-server-application-model-sam-reviewer — Review AWS Server Application Model (SAM) using its official documentation and repository in the Serverless / Framework category.


/cncf-axmos-technologies-member-reviewer — Review AXMOS Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-axolotl-reviewer — Review Axolotl using its official documentation and repository in the Training / Post Training category. Go ahead and axolotl questions


/cncf-azure-api-management-reviewer — Review Azure API Management using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. Azure API Management is a hybrid, multi-cloud API, full lifecycle management platform for APIs across all environments. It allows customers to self-host API gateways as containers on Kubernetes.


/cncf-azure-disk-storage-reviewer — Review Azure Disk Storage using its official documentation and repository in the Runtime / Cloud Native Storage category. Get HDD/SSD durability, scalability, availability, and security you need for all your workloads—from mission-critical workloads to test scenarios.


/cncf-azure-event-hubs-reviewer — Review Azure Event Hubs using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Event Hubs is a fully managed, real-time data ingestion service that’s simple, trusted, and scalable.


/cncf-azure-functions-reviewer — Review Azure Functions using its official documentation and repository in the Serverless / Hosted Platform category. Develop more efficiently with Functions, an event-driven serverless compute platform that can also solve complex orchestration problems. Build and debug locally without additional setup, deploy and operate at scale in the cloud, and integrate services using triggers and bindings.


/cncf-azure-kubernetes-service-aks-reviewer — Review Azure Kubernetes Service (AKS) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Simplify the deployment, management, and operations of Kubernetes. Use a fully managed Kubernetes container orchestration service or choose other orchestrators.


/cncf-azure-kubernetes-service-aks-wasm-reviewer — Review Azure Kubernetes Service (AKS) (Wasm) using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-azure-monitor-reviewer — Review Azure Monitor using its official documentation and repository in the Observability and Analysis / Observability category. Collect, analyze, and act on telemetry data from your Azure and on-premises environments. Azure Monitor helps you maximize performance and availability of your applications and proactively identify problems in seconds.


/cncf-azure-pipelines-reviewer — Review Azure Pipelines using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Continuously build, test, and deploy to any platform and cloud


/cncf-azure-registry-reviewer — Review Azure Registry using its official documentation and repository in the Provisioning / Container Registry category. Azure Container Registry allows you to store images for all types of container deployments including DC/OS, Docker Swarm, Kubernetes, and Azure services such as App Service, Batch, Service Fabric, and others.


/cncf-azure-service-fabric-reviewer — Review Azure Service Fabric using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-b1-systems-kcsp-reviewer — Review B1 Systems (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. B1 Systems GmbH is a worldwide operating provider of Kubernetes consulting, training, managed service &amp; support located in Germany.


/cncf-b1-systems-member-reviewer — Review B1 Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-backend-ai-reviewer — Review Backend.AI using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Backend.AI transforms GPU complexity into operational simplicity. Our container-level GPU virtualization maximizes performance while minimizing costs. The platform supports the entire AI lifecycle from data analysis to training and inference, streamlining AI development for businesses of all scales.


/cncf-backstage-reviewer — Review Backstage using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-baidu-cloud-container-engine-reviewer — Review Baidu Cloud Container Engine using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Baidu Cloud Container Engine, aka CCE.


/cncf-baidu-cloud-function-compute-reviewer — Review Baidu Cloud Function Compute using its official documentation and repository in the Serverless / Hosted Platform category. CFC (Cloud Function Compute) is a serverless solution provided by Baidu for enterprises and developers. It provides event-based, flexible, highly available, scalable and highly responsive cloud computing capabilities, and supports a variety of function triggers meeting a variety of scenarios.


/cncf-baidu-kcsp-reviewer — Review Baidu (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Baidu Cloud Container Engine is an enterprise-level container platform, the core of CCE is leveraging on open source technologies including Kubernetes and Docker.


/cncf-baidu-member-reviewer — Review Baidu (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-bamboo-reviewer — Review Bamboo using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Bamboo Server is the choice of professional teams for continuous integration, deployment, and delivery


/cncf-banco-de-cr-dito-bcp-member-reviewer — Review Banco de Crédito BCP (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bank-vaults-reviewer — Review Bank-Vaults using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Bank-Vaults is a Vault swiss-army knife: a K8s operator, Go client with automatic token renewal, automatic configuration, multiple unseal options and more. A CLI tool to init, unseal and configure Vault (auth methods, secret engines). Direct secret injection into Pods.


/cncf-bc-cloud-kcsp-reviewer — Review BC Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. BC Cloud (Beijing Big Data) provides professional cloud native platform solutions, DevOps implementation, Service Mesh consulting and other comprehensive services for big data.


/cncf-bc-cloud-member-reviewer — Review BC Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-beam-reviewer — Review Beam using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-beats-reviewer — Review Beats using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-begasoft-kcsp-reviewer — Review BEGASOFT (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We offer Swiss companies a secure, scalable Kubernetes platform – operated in our certified Swiss data centre. Beyond the platform, our Kubernetes  Professional Services support you at every stage, from architecture and security reviews to readiness checks, strategic consulting and training  – in partnership with leading experts.


/cncf-begasoft-member-reviewer — Review BEGASOFT (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bellsoft-member-reviewer — Review BellSoft (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-beneva-contributor-reviewer — Review Beneva (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-bentoml-reviewer — Review BentoML using its official documentation and repository in the AI Native Infra / Continuous Integration and Delivery category. Build Production-Grade AI Applications


/cncf-bessystem-beijing-baolande-software-corp-beefd867-reviewer — Review BESSYSTEM - Beijing Baolande Software Corporation (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Beijing Baolande Software Corporation provides Kubernetes special technical services including consulting, training, implementation, and technical support, as well as customized research &amp; development, business operations and other advanced services. We aim to help enterprises of different scales to achieve Kubernetes community version quickly and safely, better meeting their requirements of business development and accelerating business prosperity &amp; growth.


/cncf-bessystem-member-reviewer — Review BESSYSTEM (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-better-stack-member-reviewer — Review Better Stack (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bfe-reviewer — Review BFE using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. Open-source layer 7 load balancer derived from proprietary Baidu FrontEnd


/cncf-big-ip-next-edge-firewall-cnf-reviewer — Review BIG-IP NEXT Edge Firewall CNF using its official documentation and repository in the Special / Certified CNFs category. BIG-IP NEXT Edge Firewall provides Firewall, DDoS and Intrusion Prevention capabilities and is part of F5&#39;s Consolidated suite of CNFs.


/cncf-bigchaindb-reviewer — Review BigchainDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-binaryen-reviewer — Review Binaryen using its official documentation and repository in the Wasm / Tooling category.


/cncf-bindplane-member-reviewer — Review Bindplane (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-black-duck-reviewer — Review Black Duck using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Black Duck provides a comprehensive software composition analysis (SCA) solution for managing security, quality, and license compliance risk that comes from the use of open source and third-party code in applications and containers.


/cncf-blackrock-member-reviewer — Review BlackRock (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-blacksmith-member-reviewer — Review Blacksmith (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-blakyaks-kcsp-reviewer — Review BlakYaks (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. BlakYaks drive customers cloud native outcomes. A key service is designing and building Kubernetes based enterprise container platforms with a focus on industry best practice, building with code (IaC), establishing good DevSecOps &amp; GitOps processes with supply chain security embedded. Our comprehensive solutions cover strategy, design, deploy, migrations, operations &amp; support.


/cncf-blakyaks-member-reviewer — Review BlakYaks (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-blizzard-supporter-reviewer — Review Blizzard (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-block-member-reviewer — Review Block (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bloombase-reviewer — Review Bloombase using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Bloombase is an intelligent storage firewall company providing AI-accelerated security and PQC encryption of high-bandwidth, low-latency data cloud from edge computing, physical/virtual datacenters, through HCI and CDI, to the cloud.


/cncf-bloomberg-member-reviewer — Review Bloomberg (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-blue-matador-reviewer — Review Blue Matador using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-blue-sentry-kcsp-reviewer — Review Blue Sentry (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Blue Sentry increases efficiency, security, ensures compliance and helps organizations deliver products faster with our managed DevOps services utilizing Kubernetes, Infrastructure as Code and CI/CD-enabled change management.


/cncf-blue-sentry-member-reviewer — Review Blue Sentry (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bluebricks-member-reviewer — Review Bluebricks (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bmc-helix-reviewer — Review BMC Helix using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-bocloud-beyondcontainer-reviewer — Review BoCloud BeyondContainer using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. BoCloud BeyondContainer is a production-ready PaaS platform for Internet application deployment and management.


/cncf-bocloud-kcntp-reviewer — Review BoCloud (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category.


/cncf-bocloud-kcsp-reviewer — Review BoCloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category.


/cncf-bocloud-member-reviewer — Review BoCloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-boeing-contributor-reviewer — Review Boeing (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-bootc-reviewer — Review bootc using its official documentation and repository in the Runtime / Container Runtime category. The bootc provides transactional, in-place operating system images and updates using OCI/Docker container images. This project applies the Docker container layering model to bootable host systems, using standard OCI/Docker containers as a transport and delivery format for base operating system updates.


/cncf-booz-allen-hamilton-kcsp-reviewer — Review Booz Allen Hamilton (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Booz Allen partners with public and private sector clients to solve their most difficult challenges through a combination of consulting, analytics, mission operations, technology, systems delivery, cybersecurity, engineering, and innovation expertise.


/cncf-booz-allen-hamilton-member-reviewer — Review Booz Allen Hamilton (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bosh-reviewer — Review BOSH using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-botkube-reviewer — Review Botkube using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-bouncy-castle-reviewer — Review Bouncy Castle using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Bouncy Castle is one of the most widely used FIPS-certified open-source cryptographic APIs for Java and C#, allowing developers to easily integrate  PKI security into their applications.


/cncf-boundary-reviewer — Review Boundary using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Boundary is designed to grant access to critical systems using the principle of least privilege, solving challenges organizations encounter when users need to securely access applications and machines.


/cncf-box-supporter-reviewer — Review Box (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-bpfman-reviewer — Review bpfman using its official documentation and repository in the Provisioning / Security &amp; Compliance category. An eBPF Manager for Linux and Kubernetes


/cncf-braingu-member-reviewer — Review BrainGu (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-breqwatr-bks-reviewer — Review Breqwatr BKS using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. BKS is a hosted Kubernetes service on OpenStack, leveraging Cluster API to manage Kubernetes cluster lifecycle.


/cncf-breqwatr-member-reviewer — Review Breqwatr (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-brigade-reviewer — Review Brigade using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-broadcom-member-reviewer — Review Broadcom (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-brobridge-kcsp-reviewer — Review Brobridge (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Brobridge provides professional Kubernetes support and microservice-relevant components such as workflow engine, API gateway, queueing system, and orchestration engine to assist your business to move to microservice.


/cncf-brobridge-member-reviewer — Review Brobridge (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bronto-reviewer — Review Bronto using its official documentation and repository in the Observability and Analysis / Observability category. Bronto is an AI-native observability data platform for logs, metrics and traces. It lets teams ingest telemetry at petabyte scale with sub-second search, 12-month hot data retention by default and up to 90% lower storage costs, with all data parsed, searchable and enriched with context for both human and agentic workflows.


/cncf-browser-use-reviewer — Review Browser Use using its official documentation and repository in the AI Agent / Agent Tool category. Make websites accessible for AI agents. Automate tasks online with ease.


/cncf-bub-reviewer — Review Bub using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. A common shape for agents that live alongside people.


/cncf-bucket-reviewer — Review Bucket using its official documentation and repository in the Observability and Analysis / Feature Flagging category. Feature flagging purpose-built for B2B SaaS products


/cncf-bucketeer-reviewer — Review Bucketeer using its official documentation and repository in the Observability and Analysis / Feature Flagging category. A feature flag management platform created to help teams make better decisions, reduce deployment lead time, and release risk through feature flags.


/cncf-buf-technologies-member-reviewer — Review Buf Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-buildkite-reviewer — Review Buildkite using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Buildkite is a platform for running fast, secure, and scalable continuous integration pipelines on your own infrastructure.


/cncf-buildpacks-reviewer — Review Buildpacks using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-bunkerweb-reviewer — Review BunkerWeb using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. BunkerWeb is a next-generation, open-source Web Application Firewall (WAF/WAAP), reverse proxy, ingress controller, and API gateway controller.


/cncf-bunnyshell-reviewer — Review Bunnyshell using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Bunnyshell is an EaaS platform that enables fast, dead-simple environment creation and management for teams and developers who want to release better code faster.


/cncf-buoyant-member-reviewer — Review Buoyant (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-bytebase-reviewer — Review Bytebase using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Reliable Database CI/CD for Developers and DBAs


/cncf-bytesource-kcsp-reviewer — Review ByteSource (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. One Stop Shop Professional Services for everything Kubernetes


/cncf-bytesource-member-reviewer — Review ByteSource (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-c-plus-plus-reviewer — Review C++ using its official documentation and repository in the Wasm / Languages category. Compiled language to Wasm


/cncf-c-reviewer — Review C using its official documentation and repository in the Wasm / Languages category. Compiled language to Wasm


/cncf-cablelabs-member-reviewer — Review CableLabs (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-caddy-reviewer — Review Caddy using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. Caddy is a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go


/cncf-cadence-workflow-reviewer — Review Cadence Workflow using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Cadence is a distributed, scalable, durable, and highly available fault-oblivious stateful code platform.


/cncf-caepe-reviewer — Review CAEPE using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. CAEPE is a Continuous Deployment platform for Kubernetes. With CAEPE, deploy applications on Kubernetes with confidence.


/cncf-cambia-health-solutions-member-reviewer — Review Cambia Health Solutions (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-camel-reviewer — Review CAMEL using its official documentation and repository in the AI Agent / Agent Framework category. The first and the best multi-agent framework. Finding the Scaling Law of Agents.


/cncf-camptocamp-kcntp-reviewer — Review Camptocamp (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. As part of its Infrastructure services, Camptocamp provides Kubernetes consulting and training, with a focus on OpenShift installation and integration.


/cncf-camptocamp-kcsp-reviewer — Review Camptocamp (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As part of its Infrastructure services, Camptocamp provides Kubernetes consulting and training, with a focus on OpenShift installation and integration.


/cncf-camptocamp-member-reviewer — Review Camptocamp (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cann-reviewer — Review CANN using its official documentation and repository in the AI Native Infra / Accelerator and SuperPod category. Compute Architecture for Neural Networks


/cncf-canonical-charmed-kubernetes-reviewer — Review Canonical Charmed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Deploy, scale and upgrade Kubernetes clusters across multiple physical or virtual machines with Charmed Kubernetes.


/cncf-canonical-kcsp-reviewer — Review Canonical (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. The Canonical Distribution of Kubernetes enables you to operate Kubernetes clusters on demand on any major public cloud and private infrastructure.


/cncf-canonical-kubernetes-reviewer — Review Canonical Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Performant and easy to deploy Kubernetes cluster with a zero-ops experience and everything needed for a fully functioning cluster (e.g. DNS, gateway, networking).


/cncf-canonical-member-reviewer — Review Canonical (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-capital-one-member-reviewer — Review Capital One (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-capsule-reviewer — Review Capsule using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Capsule implements a multi-tenant and policy-based environment in your Kubernetes cluster. It is designed as a micro-services-based ecosystem with the minimalist approach, leveraging only on upstream Kubernetes.


/cncf-capsule8-reviewer — Review Capsule8 using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-cardinal-health-supporter-reviewer — Review Cardinal Health (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-cargo-reviewer — Review cargo using its official documentation and repository in the Wasm / Tooling category.


/cncf-cargurus-supporter-reviewer — Review CarGurus (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-cariad-supporter-reviewer — Review CARIAD (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-carina-reviewer — Review Carina using its official documentation and repository in the Runtime / Cloud Native Storage category. Carina: an high performance and ops-free local storage for kubernetes


/cncf-cars-china-academy-of-railway-sciences-kcntp-reviewer — Review CARS - China Academy of Railway Sciences (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. CARS promotes cloud native best practices in the railway domain, helping make railway information systems development and operation modern and easy. Services around Kubernetes include, but are not limited to, installation, maintenance, service operator development, and training.


/cncf-cars-china-academy-of-railway-sciences-kcsp-reviewer — Review CARS - China Academy of Railway Sciences (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. CARS promotes cloud native best practices in the railway domain, helping make railway information systems development and operation modern and easy. Services around Kubernetes include, but are not limited to, installation, maintenance, service operator development, and training.


/cncf-cars-china-academy-of-railway-sciences-member-reviewer — Review CARS - China Academy of Railway Sciences (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cars-taichu-reviewer — Review CARS TaiChu using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. TaiChu Kubernetes engine provides an easy, powerful platform that helps enterprises to use Kubernetes faster and better, including but not limited to installation and maintenance of K8s, deployment, monitoring and management of containerized applications, etc.


/cncf-cartographer-reviewer — Review Cartographer using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Cartographer is a Kubernetes-native Choreographer providing higher modularity and scalability for the software supply chain.


/cncf-cartography-reviewer — Review Cartography using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Cartography is a Python tool that consolidates infrastructure assets and the relationships between them in an intuitive graph view.


/cncf-carvel-reviewer — Review Carvel using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Carvel provides a set of reliable, single-purpose, composable tools that aid in your application building, configuration, and deployment to Kubernetes.


/cncf-casdoor-reviewer — Review Casdoor using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Casdoor is an open-source UI-first identity access management (IAM) / Single-Sign-On (SSO) platform with web UI supporting OAuth 2.0, OIDC, SAML and CAS.


/cncf-casper-reviewer — Review Casper using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-cassandra-reviewer — Review Cassandra using its official documentation and repository in the App Definition and Development / Database category.


/cncf-cast-ai-member-reviewer — Review Cast.ai (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cast-ai-reviewer — Review cast.ai using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-catalyst-by-zoho-reviewer — Review Catalyst by Zoho using its official documentation and repository in the Serverless / Hosted Platform category. Catalyst is a full-stack cloud development platform that enables you to build feature-rich applications. From building and testing to hosting and deploying, Catalyst streamlines your development cycle with a rich suite of Serverless components, Backend tools, DevOps, AI/ML capabilities and more.


/cncf-catalyst-cloud-kcsp-reviewer — Review Catalyst Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Catalyst Cloud specializes in the delivery of cloud services (private cloud consulting, training service and public cloud) benefiting from the freedom and choice of the broad ecosystem surrounding Kubernetes and OpenStack.


/cncf-catalyst-cloud-member-reviewer — Review Catalyst Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-catalyst-kubernetes-service-reviewer — Review Catalyst Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Catalyst Kubernetes Service makes it easy for you to deploy, manage, and scale Kubernetes clusters to run containerised applications on the Catalyst Cloud.


/cncf-catchpoint-reviewer — Review Catchpoint using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-cathay-financial-holding-member-reviewer — Review Cathay Financial Holding (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-causely-member-reviewer — Review Causely (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-causely-reviewer — Review Causely using its official documentation and repository in the Observability and Analysis / Observability category. Causely&#39;s Causal Reasoning Platform continuously analyzes observed anomalies to automatically pinpoint root causes and trigger immediate remediation actions.


/cncf-cdevents-reviewer — Review CDEvents using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-cdk-for-kubernetes-cdk8s-reviewer — Review CDK for Kubernetes (CDK8s) using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. CDK8s lets you define Kubernetes apps and components using familiar programming languages and object-oriented APIs.


/cncf-cecloud-kcsp-reviewer — Review CECloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. CEcloud helps customers with cloud native best practices for modernizing applications.


/cncf-cecloud-member-reviewer — Review CECloud (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-cedar-reviewer — Review Cedar using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Cedar is an open source authorization policy language that enables developers to express fine-grained permissions as easy-to-understand policies enforced in their applications, and decouple access control from application logic. Cedar is designed to be ergonomic, fast, safe, and analyzable using automated reasoning. Cedar&#39;s simple and intuitive syntax supports common authorization use-cases with readable policies, naturally expressing concepts from role-based, attribute-based, and relation-based access control models. Cedar&#39;s policy structure enables authorization requests to be decided quickly. Its policy validator uses optional typing to help policy writers avoid mistakes, but not get in their way. Cedar&#39;s design has been finely balanced to allow for a sound, complete, and decidable logical encoding, which enables precise automated analysis of Cedar policies, e.g., to ensure that policy refactoring preserves existing permissions. Cedar&#39;s language specification has been formally verified using a theorem prover to satisfy key security properties like &quot;deny trumps allow,&quot; and its implementation in Rust undergoes rigorous differential random testing against its formal specification. By combining mathematical rigor with developer-friendly design, Cedar offers a practical approach to secure, maintainable authorization for modern applications.


/cncf-celerdata-member-reviewer — Review CelerData (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-celonis-supporter-reviewer — Review Celonis (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-centreon-reviewer — Review Centreon using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-centrica-energy-contributor-reviewer — Review Centrica Energy (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-ceph-reviewer — Review Ceph using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-cerbos-reviewer — Review Cerbos using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-cern-member-reviewer — Review CERN (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-cert-manager-reviewer — Review cert-manager using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-cfengine-reviewer — Review CFEngine using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-cfl-supporter-reviewer — Review CFL (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-chainguard-member-reviewer — Review Chainguard (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-chainloop-member-reviewer — Review Chainloop (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-chaitin-tech-reviewer — Review Chaitin Tech using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-chalice-reviewer — Review Chalice using its official documentation and repository in the Serverless / Framework category.


/cncf-chaos-mesh-reviewer — Review Chaos Mesh using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-chaos-toolkit-reviewer — Review Chaos Toolkit using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-chaosblade-reviewer — Review Chaosblade using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-chaoskube-reviewer — Review chaoskube using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-chaosmeta-reviewer — Review ChaosMeta using its official documentation and repository in the Observability and Analysis / Chaos Engineering category. Provides a one-stop drill platform for the whole life cycle of access detection, fault injection, fault measurement, recovery measurement, injection recovery, etc., and has a built-in rich risk fault injection library for business, application, resources, etc., which can quickly mine the potential risks of the application system and Check the emergency fresh-keeping ability.


/cncf-chaterm-reviewer — Review Chaterm using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. AI terminal and SSH Client for EC2, Database and Kubernetes.


/cncf-check-point-reviewer — Review Check Point using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Check Point Software Technologies Ltd. is a leading provider of cyber security solutions to corporate enterprises and governments globally. Check Point protects over 100,000 organizations of all sizes.


/cncf-checkly-member-reviewer — Review Checkly (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-checkmk-reviewer — Review Checkmk using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-checkov-reviewer — Review Checkov using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Checkov scans cloud infrastructure configurations to find misconfigurations before they are deployed. Checkov manages and analyzes infrastructure as code (IaC) scan results across platforms such as Terraform, CloudFormation, Kubernetes, Helm, ARM Templates and Serverless framework.


/cncf-chef-habitat-reviewer — Review Chef Habitat using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-chef-infra-reviewer — Review Chef Infra using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-chef-inspec-reviewer — Review Chef InSpec using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-chelsio-t6-unified-wire-100gbe-adapters-reviewer — Review Chelsio T6 Unified Wire 100GbE Adapters using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-china-mobile-cloud-cloud-native-applicat-3728c2b4-reviewer — Review China Mobile Cloud Cloud-native Application Security using its official documentation and repository in the Provisioning / Security &amp; Compliance category. China Mobile Cloud Cloud-native Application Security is an application protection platform based on container technology, providing comprehensive security detection and protection capabilities for cloud-native environments.


/cncf-china-mobile-cloud-cnp-reviewer — Review China Mobile Cloud CNP using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. China Mobile Cloud CNP is an application management platform for multi-cluster and multi-cloud scenarios.


/cncf-china-mobile-kcs-reviewer — Review China Mobile KCS using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Ecloud Kubernetes Container Service（KCS）provides high-performance and reliable container application management capabilities. It simplifies the establishment of a cloud operating environment for container applications and provides a complete cloud operating environment for applications.


/cncf-china-mobile-kcsp-reviewer — Review China Mobile (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. China Mobile (Suzhou) Software Technology provides full support and services to domestic companies for their cloud native journeys


/cncf-china-mobile-member-reviewer — Review China Mobile (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-china-mobile-panji-paas-platform-reviewer — Review China Mobile Panji PaaS Platform using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. China Mobile Panji PaaS Platform is an open, standard and self-controlled cloud native technical base. Using containerization technology, it provides core capabilities including enterprise-level application service management, elastic scalability, unified governance of heterogeneous microservice framework, and panoramic cloud native observation, etc.


/cncf-china-systems-kcsp-reviewer — Review China Systems (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. China Systems’ team of cloud technology experts provides consulting, DevOps and Production training, and support services to enterprises to fast-track their adoption of Kubernetes on the cloud-native application platform.


/cncf-china-systems-member-reviewer — Review China Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-china-unicom-member-reviewer — Review China Unicom (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-chislitel-lab-member-reviewer — Review Chislitel Lab (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-choreo-reviewer — Review Choreo using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-chroma-reviewer — Review Chroma using its official documentation and repository in the AI Agent / Vector Database category. The AI-native open-source embedding database.


/cncf-chronosphere-member-reviewer — Review Chronosphere (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-chronosphere-reviewer — Review Chronosphere using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-cielara-member-reviewer — Review Cielara (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cilium-reviewer — Review Cilium using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-circleci-member-reviewer — Review CircleCI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-circleci-reviewer — Review CircleCI using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-ciroos-member-reviewer — Review Ciroos (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cisco-member-reviewer — Review Cisco (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-citi-member-reviewer — Review Citi (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-citrix-adc-formerly-netscaler-adc-reviewer — Review Citrix ADC (formerly NetScaler ADC) using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. Citrix ADC is an application delivery and load balancing solution that provides a high-quality user experience for your web, traditional, and cloud-native applications regardless of where they are hosted.


/cncf-civo-kubernetes-reviewer — Review Civo Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. The first cloud native service provider powered only by Kubernetes


/cncf-civo-member-reviewer — Review Civo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-claion-kcsp-reviewer — Review Claion (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Claion offers professional consulting on Kubernetes, and also provides cloud migration, implements and operation for Kubernetes environments. We are  developing Cloud Management Platform based on MSA and Kubernetes.


/cncf-claion-member-reviewer — Review Claion (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-clair-reviewer — Review Clair using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-cleanstart-member-reviewer — Review CleanStart (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-clever-cloud-reviewer — Review Clever Cloud using its official documentation and repository in the Platform / PaaS/Container Service category. Clever Cloud is an IT Automation platform helping developers to focus on business value instead of ops work.


/cncf-clickhouse-member-reviewer — Review ClickHouse (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-clickhouse-reviewer — Review ClickHouse using its official documentation and repository in the App Definition and Development / Database category.


/cncf-clockwork-member-reviewer — Review Clockwork (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloud-66-member-reviewer — Review Cloud 66 (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloud-66-skycap-reviewer — Review Cloud 66 Skycap using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Skycap is the easiest way to deploy and run applications on any existing Kubernetes cluster.


/cncf-cloud-community-labs-member-reviewer — Review Cloud Community Labs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloud-custodian-reviewer — Review Cloud Custodian using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-cloud-foundry-application-runtime-reviewer — Review Cloud Foundry Application Runtime using its official documentation and repository in the Platform / PaaS/Container Service category. Cloud Foundry Application Runtime utilizes containers as part of its DNA, and has since before Docker popularized containers. The new CF Container Runtime gives you more granular control and management of containers with Kubernetes.


/cncf-cloud-foundry-foundation-member-reviewer — Review Cloud Foundry Foundation (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-cloud-native-landscape-reviewer — Review Cloud Native Landscape using its official documentation and repository in the Serverless / Tools category.


/cncf-cloud-native-texas-member-reviewer — Review Cloud Native Texas (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudark-kubeplus-reviewer — Review CloudARK KubePlus using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Build SaaS for your containerized applications


/cncf-cloudbase-solutions-member-reviewer — Review Cloudbase Solutions (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudbees-codeship-reviewer — Review Cloudbees Codeship using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Ship faster with CI/CD as a Service


/cncf-cloudbees-member-reviewer — Review CloudBees (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudbolt-software-member-reviewer — Review CloudBolt Software (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudcasa-by-catalogic-software-reviewer — Review CloudCasa by Catalogic Software using its official documentation and repository in the Runtime / Cloud Native Storage category. Backup-as-a-service for Kubernetes and cloud native applications


/cncf-cloudchipr-member-reviewer — Review Cloudchipr (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudeq-member-reviewer — Review cloudEQ (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudera-kcsp-reviewer — Review Cloudera (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide comprehensive, end-to-end Kubernetes services designed to accelerate your modernization journey and ensure long-term operational excellence.


/cncf-cloudera-member-reviewer — Review Cloudera (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudevents-reviewer — Review CloudEvents using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Standardizing common eventing metadata and their location to help with event identification and routing.


/cncf-cloudferro-member-reviewer — Review CloudFerro (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudflare-workers-reviewer — Review Cloudflare Workers using its official documentation and repository in the Serverless / Hosted Platform category. Cloudflare Workers provides a lightweight JavaScript execution environment that allows developers to augment existing applications or create entirely new ones without configuring or maintaining infrastructure.


/cncf-cloudflare-workers-wasm-reviewer — Review Cloudflare Workers (Wasm) using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-cloudgeometry-kcsp-reviewer — Review CloudGeometry (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. CloudGeometry is an expert managed cloud software engineering consultancy, cloud-native systems integrator, and Kubernetes Certified Service Provider. Our hands-on design and implementation services help our clients achieve Kubernetes success with open-source tooling and commercial cloud platform technologies, including AWS, Azure, and GCP.


/cncf-cloudgeometry-member-reviewer — Review CloudGeometry (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudification-member-reviewer — Review Cloudification (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudify-reviewer — Review Cloudify using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-cloudlink-cmp-reviewer — Review CloudLink-CMP using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. CloudLink CMP is a kubernetes enabled enterprise-ready PaaS solution, consists of features like CICD、multi-tenant, rich micro-services and components support, integrated with an easy-to-use operation/maintenance interface.


/cncf-cloudmatos-reviewer — Review CloudMatos using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Pioneering the world of self-healing, self-aware, self-sustaining, self-resilient, self-secure and intelligent remediation, MatosSphere brings a complete cloud security and governance solution for your cloud infrastructure.


/cncf-cloudnativepg-reviewer — Review CloudNativePG using its official documentation and repository in the App Definition and Development / Database category. The most popular Kubernetes Operator for PostgreSQL.


/cncf-cloudops-by-aptum-kcntp-reviewer — Review CloudOps by Aptum (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category.


/cncf-cloudops-by-aptum-kcsp-reviewer — Review CloudOps by Aptum (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Training, support, and services for DevOps, Kubernetes, and cloud native practices. We design, build and operate DevOps platforms and hybrid cloud platforms.


/cncf-cloudops-by-aptum-member-reviewer — Review CloudOps by Aptum (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudpilot-ai-reviewer — Review CloudPilot AI using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. CloudPilot AI automates Kubernetes cost optimization to reduce cloud costs up to 80%. Zero performance impact, 5-minute deployment.


/cncf-cloudscale-ch-member-reviewer — Review cloudscale.ch (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudsmith-member-reviewer — Review Cloudsmith (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cloudsmith-reviewer — Review Cloudsmith using its official documentation and repository in the Provisioning / Container Registry category. A single source of truth for every artifact and container. Cloudsmith is a powerful, cloud native, enterprise-grade artifact management solution.


/cncf-cloudtty-reviewer — Review CloudTTY using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. A Friendly Kubernetes CloudShell (Web Terminal)


/cncf-cloudwego-reviewer — Review CloudWeGo using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. CloudWeGo is ByteDance&#39;s open source Golang-centric middleware that can be used to quickly build enterprise-class cloud native architectures.


/cncf-cloudwise-synthetic-monitoring-reviewer — Review Cloudwise Synthetic Monitoring using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-cloudzero-reviewer — Review CloudZero using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-clush-eduplace-reviewer — Review Clush EduPlace using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Clush EduPlace is a cloud-native, event-driven messaging platform for enterprise collaboration and workflow automation.   It connects legacy systems through API and Bot integration, supports real-time WebSocket communication, and operates on a Container-based architecture.   Designed for secure financial and public sector environments, it enables on-premise and hybrid cloud deployments with high availability and data compliance.


/cncf-clush-kcsp-reviewer — Review Clush (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Unified Management of Kubernetes and Cloud Native Environments, Accelerating Enterprise AI through Expertise in GPUOps and MLOps Platforms.


/cncf-clush-kube-reviewer — Review Clush Kube using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Clush Kube is a next-generation enterprise cloud operating solution that provides user management, systems management, PaaS management, DevOps management, helpdesk/service support, and monitoring capabilities.


/cncf-clush-member-reviewer — Review Clush (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-clusternet-reviewer — Review Clusternet using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. [CNCF Sandbox Project] Managing your Kubernetes clusters (including public, private, edge, etc.) as easily as visiting the Internet


/cncf-clusterpedia-reviewer — Review Clusterpedia using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Clusterpedia is used for complex resources search across multiple clusters, support simultaneous search of a single kind of resource  or multiple kinds of resources existing in multiple clusters.


/cncf-clyso-member-reviewer — Review Clyso (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cme-group-member-reviewer — Review CME Group (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cni-genie-reviewer — Review CNI-Genie using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-cobank-member-reviewer — Review CoBank (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cockroachdb-reviewer — Review CockroachDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-cocktail-cloud-reviewer — Review Cocktail Cloud using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Cocktail Cloud is a container application management platform that provides a continuous development/deployment/operating environment for container-based applications.


/cncf-cocktail-io-kcntp-reviewer — Review Cocktail.io (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. As a CNCF KCSP (Kubernetes Certified Service Provider) and a developer of Cocktail Cloud Kubernetes Orchestration Platform Product, Acornsoft trains clients and prospective clients on Docker and Kubernetes.


/cncf-cocktail-io-kcsp-reviewer — Review Cocktail.io (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Cocktail.io develops and provides Cloud Native Application Platform, a computing environment for future application development and helps customers adopt container, DevOps, and microservice technologies.


/cncf-cocktail-io-member-reviewer — Review Cocktail.io (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-codefresh-reviewer — Review Codefresh using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-coder-member-reviewer — Review Coder (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-coder-reviewer — Review Coder using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Self-hosted cloud development environments consistently provisioned as code and pre-configured for developer activity on day one.


/cncf-codezero-reviewer — Review CodeZero using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Tools for local collaborative development in pre-production Kubernetes clusters


/cncf-cohdi-reviewer — Review CoHDI using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. CoHDI (Composable Hardware in Disaggregated Infrastructure) enables dynamic device scaling across next-generation architectures. As a community-driven, standards-based open ecosystem, CoHDI focuses on expanding cloud-native frameworks built on disaggregate infrastructure. Our core objective is to bridge the gap between Kubernetes and underlying hardware by actively collaborating with upstream projects to increase cloud native composability, specifically Dynamic Resource Allocation (DRA), Autoscaler, and Scheduling. By integrating these cloud-native capabilities, CoHDI empowers data center and infrastructure operators to maximize cost efficiency, achieve high availability, and drive sustainability through a seamlessly disaggregated computing system.


/cncf-commvault-reviewer — Review Commvault using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-componentize-js-reviewer — Review componentize-js using its official documentation and repository in the Wasm / Languages category. Scripting languages that support Wasm


/cncf-composefs-reviewer — Review composefs using its official documentation and repository in the Runtime / Container Runtime category. A project that combines Linux kernel features to provide read-only mountable filesystem trees stacking on top of an underlying &quot;lower&quot; Linux filesystem, particularly useful for mounting container images.


/cncf-composio-reviewer — Review Composio using its official documentation and repository in the AI Agent / Agent Tool category. Composio powers 1000+ toolkits, tool search, context management, authentication, and a sandboxed workbench to help you build AI agents that turn intent into action.


/cncf-concourse-reviewer — Review Concourse using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-confidential-containers-reviewer — Review Confidential Containers using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Confidential Containers is an open source community working to enable cloud native  confidential computing by leveraging Trusted Execution Environments to protect  containers and data.


/cncf-confidentialmind-member-reviewer — Review ConfidentialMind (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-confluentis-consulting-kcsp-reviewer — Review Confluentis Consulting (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Confluentis Consulting is a technology engineering and DevOps focused consulting firm helping enterprises build scalable, cost-efficient, and resilient digital platforms. We specialize in Cloud FinOps, multi-cloud consulting, and platform engineering, enabling organizations to gain deep cost visibility, optimize Kubernetes and cloud workloads, and align engineering decisions with financial outcomes.


/cncf-confluentis-consulting-member-reviewer — Review Confluentis Consulting (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-connect-rpc-reviewer — Review Connect RPC using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. Connect is a family of libraries for building browser and gRPC-compatible HTTP APIs.


/cncf-conoa-kcntp-reviewer — Review Conoa (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Conoa offers a wide range of training and workshops in Kubernetes, cloud native and container technologies. All our instructors have extensive hands-on experience in the technologies in which they teach.


/cncf-conoa-kcsp-reviewer — Review Conoa (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Conoa helps businesses build cloud native, container and Kubernetes infrastructure to enable modernisation. Our experts bring organisations from zero knowledge to production-ready environments quickly.


/cncf-conoa-member-reviewer — Review Conoa (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-conoa-proact-managed-container-platform-reviewer — Review Conoa Proact Managed Container Platform using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Full stack managed Kubernetes platform which lets you focus on your application, we&#39;ll do the rest!


/cncf-consul-reviewer — Review Consul using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-container-network-interface-cni-reviewer — Review Container Network Interface (CNI) using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-container-storage-interface-csi-reviewer — Review Container Storage Interface (CSI) using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-container2wasm-reviewer — Review container2wasm using its official documentation and repository in the Wasm / Orchestration &amp; Management category. A tool to run containers on Wasm-enabled environments.


/cncf-containerd-reviewer — Review containerd using its official documentation and repository in the Runtime / Container Runtime category. An open and reliable container runtime


/cncf-containerd-wasm-reviewer — Review containerd (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-containerized-routing-protocol-daemon-reviewer — Review Containerized Routing Protocol Daemon using its official documentation and repository in the Special / Certified CNFs category. The Junos® containerized routing protocol daemon (cRPD) offers deployment-hardened feature-rich routing functionality in a container for cloud-native deployments.


/cncf-containerized-security-router-and-switch-csrx-reviewer — Review Containerized Security Router and Switch (cSRX) using its official documentation and repository in the Special / Certified CNFs category. The cSRX Container Firewall protects containerized applications and environments with advanced security services, including content security and intrusion prevention system (IPS). Purpose-built for containers, the cSRX next-generation firewall can be spun up or down in less than a second for the agility needed to manage transitory container environments like Kubernetes.


/cncf-containerssh-reviewer — Review ContainerSSH using its official documentation and repository in the Provisioning / Security &amp; Compliance category. ContainerSSH launches a new container for each SSH connection in Kubernetes, Podman or Docker. The user is transparently dropped in the container and the container is removed when the user disconnects. Authentication and container configuration are dynamic using webhooks, no system users required.


/cncf-contour-reviewer — Review Contour using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-control-plane-corporation-member-reviewer — Review Control Plane Corporation (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-control-plane-managed-kubernetes-reviewer — Review Control Plane Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Control Plane Managed Kubernetes is a service that creates and manages Kubernetes clusters across multiple cloud platforms and on-prem environments.


/cncf-controlmonkey-member-reviewer — Review ControlMonkey (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-controlplane-kcntp-reviewer — Review ControlPlane (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. ControlPlane Kubernetes training covers operations, cloud native security, secure SDLC and application delivery, and penetration testing.


/cncf-controlplane-kcsp-reviewer — Review ControlPlane (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ControlPlane is a cloud native security consultancy specialising in high-compliance development and security operations for financial and government organisations, third party code risk management, and end-to-end supply chain security.


/cncf-controlplane-member-reviewer — Review ControlPlane (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-controltheory-member-reviewer — Review ControlTheory (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-controltheory-reviewer — Review ControlTheory using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. ControlTheory builds practical observability tools for modern SREs - open-source tools and AI-powered platforms that distill telemetry down to essential answers.


/cncf-copa-reviewer — Review Copa using its official documentation and repository in the Provisioning / Security &amp; Compliance category. CLI tool for directly patching container image vulnerabilities


/cncf-core-24-7-kcsp-reviewer — Review Core 24/7 (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We specialize in Kubernetes-based infrastructure design, deployment, and ongoing management. Our certified engineers deliver DevOps-as-a-Service, ensuring 24×7×365 production reliability and developer assistance.


/cncf-core-24-7-member-reviewer — Review Core 24/7 (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-coredge-kcsp-reviewer — Review Coredge (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Coredge helps to swiftly deploy, manage, and maintain Kubernetes with the help of our comprehensive Kubernetes professional services that accelerates your transformation journey with our innovative solutions.


/cncf-coredge-kubernetes-platform-reviewer — Review Coredge Kubernetes Platform using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. A custom Kubernetes distro.


/cncf-coredge-member-reviewer — Review Coredge (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-coredns-reviewer — Review CoreDNS using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category.


/cncf-coreweave-kubernetes-service-cks-reviewer — Review CoreWeave Kubernetes Service (CKS) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. CKS is a managed Kubernetes environment purpose-built for building, training, and deploying AI applications.


/cncf-coreweave-member-reviewer — Review CoreWeave (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-coroot-reviewer — Review Coroot using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-cortex-member-reviewer — Review Cortex (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cortex-reviewer — Review Cortex using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-cosmonic-member-reviewer — Review Cosmonic (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cosmonic-reviewer — Review Cosmonic using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-cosmonic-wasm-reviewer — Review Cosmonic (Wasm) using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-cosmwasm-reviewer — Review CosmWasm using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-couchbase-reviewer — Review Couchbase using its official documentation and repository in the App Definition and Development / Database category.


/cncf-couler-reviewer — Review Couler using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-cozystack-reviewer — Review Cozystack using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Cozystack is a free PaaS platform and framework for building private clouds and providing users/customers with managed Kubernetes,  KubeVirt-based VMs, databases as a service, NATS, message brokers, etc. with GPU support in VMs and Kubernetes clusters.


/cncf-crate-io-reviewer — Review Crate.io using its official documentation and repository in the App Definition and Development / Database category.


/cncf-crates-io-reviewer — Review crates.io using its official documentation and repository in the Wasm / Packaging, Registries &amp; Application Delivery category.


/cncf-crew-ai-reviewer — Review Crew AI using its official documentation and repository in the AI Agent / Agent Framework category. Framework for orchestrating role-playing, autonomous AI agents. By fostering collaborative intelligence, CrewAI empowers agents to work together seamlessly, tackling complex tasks.


/cncf-cri-o-reviewer — Review CRI-O using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-cri-o-wasm-reviewer — Review cri-o (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-crossplane-reviewer — Review Crossplane using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Crossplane is the cloud native control plane framework that allows you to build control planes without needing to write code. Crossplane has a highly extensible backend that enables you to orchestrate applications and infrastructure no matter where they run and a highly configurable frontend that lets you define the declarative API it offers.


/cncf-crun-reviewer — Review crun using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-crunchy-postgres-operator-reviewer — Review Crunchy Postgres Operator using its official documentation and repository in the App Definition and Development / Database category.


/cncf-crusoe-member-reviewer — Review Crusoe (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cto2b-member-reviewer — Review CTO2B (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cubefs-reviewer — Review CubeFS using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-cubesandbox-reviewer — Review CubeSandbox using its official documentation and repository in the AI Native Infra / Workload Runtime category. A high-performance sandbox infrastructure for AI Agent execution with 60ms startup and 5MB memory footprint.


/cncf-cubex-reviewer — Review Cubex using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. CubeX is the easiest and most powerful PaaS solution for Cloud Native Apps. CubeX adds stability to speed, enabling full DevOps.


/cncf-cue-labs-member-reviewer — Review CUE Labs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cumulus-reviewer — Review Cumulus using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-curiefense-reviewer — Review Curiefense using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-curve-reviewer — Review Curve using its official documentation and repository in the Runtime / Cloud Native Storage category. Curve is a distributed storage system designed and developed independently by NetEase,  featured with high performance, high availability, high reliability and well expansibility,  and it can serve as the basis for storage systems designed for different scenario.


/cncf-curvine-reviewer — Review Curvine using its official documentation and repository in the Runtime / Cloud Native Storage category. Curvine is a high-performance distributed multi-tier caching system written in Rust (memory/SSD/HDD), providing POSIX (FUSE), S3, and HDFS access to cloud object storage, with Kubernetes CSI integration.


/cncf-cvs-health-member-reviewer — Review CVS Health (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-cyberark-conjur-reviewer — Review CyberArk Conjur using its official documentation and repository in the Provisioning / Key Management category.


/cncf-cybozu-kubernetes-engine-reviewer — Review Cybozu Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Cybozu Kubernetes Engine, a distributed service that automates Kubernetes cluster management.


/cncf-cybozu-member-reviewer — Review Cybozu (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-cyclops-reviewer — Review Cyclops using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Cyclops is a customizable UI for Kubernetes workloads.


/cncf-cyso-kcsp-reviewer — Review Cyso (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Cyso Cloud, based in Netherlands, operating in the EU, our Managed Kubernetes service helps you effortlessly manage, scale, and optimise your applications in secure and compliant way (ISO27K1/SOC2/DORA/GDPR) in both Private and Public Cloud environments. We provide training, migration, Devops, SRE, Security, Compliancy Services and operate with 24x7 NOC and SOC.


/cncf-cyso-managed-kubernetes-reviewer — Review CYSO Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Simplify your cloud-native journey with Managed Kubernetes. Powered by Gardener, deploy, manage, and scale your containerized apps effortlessly on Cyso Cloud and beyond—all GDPR compliant.


/cncf-cyso-member-reviewer — Review Cyso (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dachs-it-kcsp-reviewer — Review DACHS IT (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. DACHS IT GMBH provides Kubernetes consulting, training, managed service &amp; support located in Germany.


/cncf-dachs-it-member-reviewer — Review DACHS IT (member) using its official documentation and repository in the CNCF Members / Silver category. DACHS IT GMBH provides Kubernetes consulting, training, managed service &amp; support located in Germany.


/cncf-daekyo-cns-kcsp-reviewer — Review Daekyo CNS (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Daekyo CNS provides consulting services by experts with professional CKA qualifications and consulting experience.


/cncf-daekyo-cns-member-reviewer — Review Daekyo CNS (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dalec-reviewer — Review Dalec using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Dalec provides a declarative format for building system packages and containers from those packages in a secure way for supply chain security.


/cncf-daocloud-enterprise-reviewer — Review DaoCloud Enterprise using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Daocloud helps you provide a reliable and consistent basic support environment to meet the high SLA requirements of enterprise critical applications.


/cncf-daocloud-kcntp-reviewer — Review DaoCloud (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. DaoCloud is an innovative leader in the cloud-native field, being committed to creating an open Cloud OS which enables enterprises to easily carry out digital transformation.  DaoCloud provides excellent courses of KCNA, CKAD, CKA and CKS, with the accumulation of cloud-native technology and examination. These courses can train technical experts in digital transformation.


/cncf-daocloud-kcsp-reviewer — Review DaoCloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide enterprise-level cloud native application platform that supports both Kubernetes and Docker Swarm.


/cncf-daocloud-member-reviewer — Review DaoCloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dapr-reviewer — Review Dapr using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-dapr-sdk-reviewer — Review Dapr SDK using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-dart-reviewer — Review Dart using its official documentation and repository in the Wasm / Languages category. Managed language


/cncf-dash0-member-reviewer — Review Dash0 (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dash0-reviewer — Review Dash0 using its official documentation and repository in the Observability and Analysis / Observability category. Dash0 is the cloud native observability solution, powered by OpenTelemetry, Prometheus and Perses. Gain deep insights into your applications and infrastructure with a solution that is fully based on open standards. Built for developers and loved by SREs - ready for your stack.


/cncf-dashbird-reviewer — Review Dashbird using its official documentation and repository in the Serverless / Tools category.


/cncf-databend-reviewer — Review Databend using its official documentation and repository in the App Definition and Development / Database category. Databend is a modern Elasticity and Performance cloud data warehouse, activate your object storage for real-time analytics.  Databend Serverless at https://app.databend.com/


/cncf-databuff-reviewer — Review DataBuff using its official documentation and repository in the Observability and Analysis / Observability category. Open-source, AI-native OpenTelemetry APM platform. OTLP trace and metrics ingest, service topology, and multi-agent SRE investigation that cites real spans. Self-hosted via Docker or Kubernetes.


/cncf-datacore-kcsp-reviewer — Review DataCore (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. DataCore offers container-native storage software solutions that simplify the deployment and management of stateful applications on Kubernetes. Implementation and support services are also available.


/cncf-datacore-member-reviewer — Review DataCore (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-datacore-puls8-reviewer — Review DataCore Puls8 using its official documentation and repository in the Runtime / Cloud Native Storage category. DataCore Puls8 gives Kubernetes the enterprise-grade storage it needs – fast, resilient, and fully container-native. It ensures stateful workloads stay persistent, protected, and production-ready across any environment.


/cncf-datadog-member-reviewer — Review Datadog (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-datadog-reviewer — Review Datadog using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-datafy-member-reviewer — Review Datafy (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-datagalaxy-supporter-reviewer — Review DataGalaxy (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-dataset-reviewer — Review DataSet using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-datenlord-reviewer — Review DatenLord using its official documentation and repository in the Runtime / Cloud Native Storage category. DatenLord is a cloud-native distributed storage platform, aiming to meet the performance-critical storage needs from next-generation cloud-native applications.


/cncf-datera-reviewer — Review Datera using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-datica-reviewer — Review Datica using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-daytona-reviewer — Review Daytona using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-db-systel-supporter-reviewer — Review DB Systel (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-de-novo-kcsp-reviewer — Review De Novo (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. De Novo, Ukraine&#39;s leading cloud provider since 2012, offers a revamped service portfolio, including classic IaaS and advanced VMware Tanzu Kubernetes Grid-based private and public Kubernetes orchestration platforms. Our professional services cover consulting, training, PoC, onboarding, migration, and commercial-grade support.


/cncf-de-novo-member-reviewer — Review De Novo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-decathlon-contributor-reviewer — Review Decathlon (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-decimal-member-reviewer — Review Decimal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-deepchecks-reviewer — Review Deepchecks using its official documentation and repository in the AI Native Infra / Observability category. Tests for Continuous Validation of ML Models &amp; Data. Deepchecks is a holistic open-source solution for all of your AI &amp; ML validation needs, enabling to thoroughly test your data and models from research to production.


/cncf-deeperthanblue-kcsp-reviewer — Review DeeperThanBlue (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As a Kubernetes Certified Service Provider, we offer cutting-edge solutions that empower your organisation to navigate the complexities of cloud migration and digital transformation seamlessly. Our professional services team help you with platform design, build, DevOps, CI/CD pipeline and operations management of your Kubernetes environments and cloud native transition.


/cncf-deeperthanblue-member-reviewer — Review DeeperThanBlue (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-deepeval-reviewer — Review DeepEval using its official documentation and repository in the AI Agent / Evaluation category. The LLM Evaluation Framework


/cncf-deepflow-reviewer — Review DeepFlow using its official documentation and repository in the AI Native Infra / Observability category. DeepFlow leverages eBPF to collect signals such as performance metrics, distributed tracing, and CPU &amp; GPU profiling. Its goal is to provide zero-code observability for complex cloud-native and AI applications, particularly those involving LLM training and inference using frameworks like PyTorch.


/cncf-deepshore-kcsp-reviewer — Review Deepshore (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Deepshore provides Kubernetes support in planning, installation and operating your container platform.


/cncf-deepshore-member-reviewer — Review Deepshore (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-deepspeed-reviewer — Review DeepSpeed using its official documentation and repository in the Training / Distributed Training category. Deep learning optimization library that makes distributed training and inference easy, efficient, and effective.


/cncf-deepstream-reviewer — Review deepstream using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-defense-unicorns-member-reviewer — Review Defense Unicorns (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dell-emc-reviewer — Review Dell EMC using its official documentation and repository in the Runtime / Cloud Native Storage category. Data storage that takes you from insights to innovation.


/cncf-dell-technologies-consulting-kcsp-reviewer — Review Dell Technologies Consulting (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Leverage our Cloud Native Application Modernization and DevOps consulting practice offerings to accelerate your adoption of K8s and containers.


/cncf-dell-technologies-member-reviewer — Review Dell Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-deployhub-reviewer — Review DeployHub using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-depot-member-reviewer — Review Depot (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-depot-reviewer — Review Depot using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-desktop-kubernetes-reviewer — Review Desktop Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Stands up a three-VM Kubernetes development cluster on the desktop using VirtualBox, by running a single Bash script.


/cncf-desotech-kcntp-reviewer — Review Desotech (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Desotech provides training and consulting services for the world’s leading containerization platform and can customize it to your business needs.


/cncf-desotech-kcsp-reviewer — Review Desotech (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Desotech provides training and consulting services for the world’s leading containerization platform and can customize it to your business needs.


/cncf-desotech-member-reviewer — Review Desotech (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-deutsche-telekom-ag-member-reviewer — Review Deutsche Telekom AG (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-devcycle-reviewer — Review DevCycle using its official documentation and repository in the Observability and Analysis / Feature Flagging category.


/cncf-develocity-reviewer — Review Develocity using its official documentation and repository in the Observability and Analysis / Observability category. Develocity speeds up builds and tests, makes troubleshooting more efficient, and increases toolchain observability with reporting and visualization. Develocity gives developers back at least one day each week in lost productivity, and supports Apache Maven, Bazel, sbt, and Gradle build systems.


/cncf-devfile-reviewer — Review Devfile using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. An open standard defining containerized development environments that enables developer tools to simplify and accelerate workflows


/cncf-devspace-reviewer — Review DevSpace using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Client-Only Developer Tool for Cloud-Native Development with Kubernetes


/cncf-devstream-reviewer — Review DevStream using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-devtron-member-reviewer — Review Devtron (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-devtron-reviewer — Review Devtron using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Open source Software delivery workflow for Kubernetes


/cncf-devzero-member-reviewer — Review DevZero (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dex-reviewer — Review Dex using its official documentation and repository in the Provisioning / Security &amp; Compliance category. OpenID Connect (OIDC) identity and OAuth 2.0 provider with pluggable connectors


/cncf-dfds-supporter-reviewer — Review DFDS (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-dfinity-reviewer — Review Dfinity using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-dg-i-kubernetes-platform-reviewer — Review DG-i Kubernetes Platform using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. The DG-i Kubernetes Platform is a curated and production-proven solution for our fully managed service for regulated environments including monitoring, logging and audit-ready operations across the DG-i private cloud, on-premise installations or in a multi-cloud environment.


/cncf-dgi-kcsp-reviewer — Review DGi (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide Managed Kubernetes Services, helping financial services providers and other organizations in highly regulated industries to plan, set up, and operate secure, high-available clusters in our datacenters or in hybrid environments across private and public clouds. We support the migration and onboarding of internal and external software solutions onto Kubernetes and ensure reliable 24/7 operations.


/cncf-dgi-member-reviewer — Review DGi (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dgraph-reviewer — Review Dgraph using its official documentation and repository in the App Definition and Development / Database category.


/cncf-dhh-member-reviewer — Review DHH (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-diagrid-member-reviewer — Review Diagrid (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-diamanti-kcsp-reviewer — Review Diamanti (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Diamanti offers support and training for the installation and operation of Diamanti’s Kubernetes Certified bare-metal platform as well as provides consulting services for onboarding 3rd party applications on-premises and in public cloud environments.


/cncf-diamanti-kubernetes-engine-reviewer — Review Diamanti Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Diamanti Kubernetes Engine is a pre-validated and pre-packaged Kubernetes distribution complete with integrated CSI and CNI plugins for modern, distributed applications across the hybrid cloud


/cncf-diamanti-member-reviewer — Review Diamanti (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-diamanti-reviewer — Review Diamanti using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-dify-reviewer — Review Dify using its official documentation and repository in the AI Agent / Workflow Orchestration category. Production-ready platform for agentic workflow development.


/cncf-digital-china-member-reviewer — Review Digital China (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-digitalocean-kubernetes-reviewer — Review DigitalOcean Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. DigitalOcean Kubernetes is designed for developers who want a simple and cost effective way to deploy container workloads to a managed Kubernetes service.


/cncf-digitalocean-member-reviewer — Review DigitalOcean (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dina-it-solutions-kcsp-reviewer — Review Dina IT Solutions (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Dina IT Solutions provides consulting on cloud native technologies such as Kubernetes. We design, implement, and operate your Kubernetes solutions.


/cncf-dina-it-solutions-member-reviewer — Review Dina IT Solutions (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-directv-supporter-reviewer — Review DIRECTV (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-direktiv-reviewer — Review Direktiv using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-distr-reviewer — Review Distr using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. The Open Source control plane for self-managed, BYOC, and on-prem deployments. Everything you need to distribute applications to self-managed customers out of the box.


/cncf-distribution-reviewer — Review Distribution using its official documentation and repository in the Provisioning / Container Registry category. The toolkit to pack, ship, store, and deliver container content


/cncf-dlocal-member-reviewer — Review Dlocal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-docker-compose-reviewer — Review Docker Compose using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-docker-hub-reviewer — Review Docker Hub using its official documentation and repository in the Wasm / Packaging, Registries &amp; Application Delivery category.


/cncf-docker-member-reviewer — Review Docker (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-docker-swarm-reviewer — Review Docker Swarm using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-docker-wasm-reviewer — Review Docker (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-docling-reviewer — Review Docling using its official documentation and repository in the AI Agent / RAG category. Get your documents ready for gen AI


/cncf-docomo-innovations-member-reviewer — Review DOCOMO Innovations (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-doctor-droid-member-reviewer — Review Doctor Droid (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-docusign-member-reviewer — Review Docusign (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-doit-kcntp-reviewer — Review DoiT (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Founded by experienced software engineers and Kubernetes contributors, DoiT delivers hands-on Kubernetes consulting and training for some of the most demanding deployments in North America and Israel.


/cncf-doit-kcsp-reviewer — Review DoiT (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. DoiT helps organizations design, deploy, and operationalize Kubernetes environments across AWS, Google Cloud, and Azure. As a Kubernetes Certified Service Provider, we specialize in production-ready architectures built on EKS, GKE, and AKS, with consistent patterns for security, scalability, and cost efficiency. With extensive experience supporting Kubernetes workloads across industries and clouds, DoIT helps customers move from pilot to production quickly while maintaining flexibility, compliance, and operational maturity across any environment.


/cncf-doit-member-reviewer — Review DoiT (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dolphinscheduler-reviewer — Review DolphinScheduler using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-dongobi-member-reviewer — Review Dongobi (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-doris-reviewer — Review Doris using its official documentation and repository in the App Definition and Development / Database category.


/cncf-dosec-kcntp-reviewer — Review Dosec (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. As the leader of cloud native security in China and the first cloud native security manufacturer in China to join CNCF, Xiaoyou Technology has a large number of Kubernetes expert lecturers who have passed CKA certification and the most professional container security R&amp;D team in China, providing Kubernetes full life cycle service solutions.


/cncf-dosec-kcsp-reviewer — Review Dosec (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Dosec provides Kubernetes services including technical support, consulting, security, and professional training.


/cncf-dosec-member-reviewer — Review Dosec (member) using its official documentation and repository in the CNCF Members / Silver category. Dosec is a cloud-native security solution provider with security protection capabilities for containers, kubernetes, microservices, serverless, and DevOps.


/cncf-dosec-reviewer — Review Dosec using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Dosec Security Platform provides security protection capabilities for containers, kubernetes, microservices, serverless, and DevOps.


/cncf-doubleword-member-reviewer — Review Doubleword (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-draftt-member-reviewer — Review Draftt (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dragonfly-reviewer — Review Dragonfly using its official documentation and repository in the Provisioning / Container Registry category. Delivers efficient, stable, and secure data distribution and acceleration powered by P2P technology, with an optional content‑addressable filesystem that accelerates OCI container launch.


/cncf-dragonflydb-member-reviewer — Review Dragonflydb (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dragonflydb-reviewer — Review DragonflyDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-drasi-reviewer — Review Drasi using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. A data change processing platform to simplify change-driven systems that need to detect, evaluate, and react to data changes quickly and efficiently at scale.


/cncf-drivescale-reviewer — Review DriveScale using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-druid-reviewer — Review Druid using its official documentation and repository in the App Definition and Development / Database category.


/cncf-dubbo-reviewer — Review Dubbo using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category.


/cncf-dx-o2-by-broadcom-reviewer — Review DX O2 by Broadcom using its official documentation and repository in the Observability and Analysis / Observability category. DX O2 is an AIOps &amp; Observability platform that offers infrastructure, applications &amp; service assurance across Mainframe to mobile in hybrid cloud data center.


/cncf-dynatrace-member-reviewer — Review Dynatrace (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-dynatrace-reviewer — Review Dynatrace using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-e-on-member-reviewer — Review E.ON (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-e2b-reviewer — Review e2b using its official documentation and repository in the AI Native Infra / Workload Runtime category. Open-source, secure environment with real-world tools for enterprise-grade agents.


/cncf-easeagent-reviewer — Review EaseAgent using its official documentation and repository in the Observability and Analysis / Observability category. EaseAgent is a Javaagent that can be integrated with the mainstream monitoring system, providing standard data formats that are fully compatible with OpenZipkin and Prometheus. EaseAgent is also very easy to extend through the Plugin Mechanism which only a minimum of three interfaces are required to be implemented to complete a plugin development.


/cncf-easegress-reviewer — Review Easegress using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. A Cloud Native traffic orchestration system


/cncf-easemesh-reviewer — Review EaseMesh using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-easy-ngo-reviewer — Review Easy-Ngo using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. easy-ngo is a fast, out-of-the-box microservice development framework based on the Go language.


/cncf-easystack-kcntp-reviewer — Review EasyStack (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category.


/cncf-easystack-kcsp-reviewer — Review EasyStack (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category.


/cncf-easystack-kubernetes-service-eks-reviewer — Review EasyStack Kubernetes Service (EKS) using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. EasyStack Kubernetes Service (EKS) is enterprise container platform contains the best of bred capability of an integrated Kubernetes and OpenStack solution.


/cncf-easystack-member-reviewer — Review EasyStack (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ebay-member-reviewer — Review eBay (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-echo-member-reviewer — Review Echo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-eclipse-che-reviewer — Review Eclipse Che using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-eclipse-foundation-member-reviewer — Review Eclipse Foundation (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-edb-kcsp-reviewer — Review EDB (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. EDB is at the forefront of running Postgres workloads on Kubernetes as the creator and primary contributor to CloudNativePG, now a CNCF Sandbox  project. In the Cloud-Native space, our services are exclusively focused on Postgres, providing specialized training, consulting, and 24/7 support to ensure reliable Kubernetes-native deployment and operations for mission-critical databases.


/cncf-edb-member-reviewer — Review EDB (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-edb-reviewer — Review EDB using its official documentation and repository in the App Definition and Development / Database category. EDB Postgres for Kubernetes is a PostgreSQL operator based on CloudNativePG


/cncf-edera-member-reviewer — Review Edera (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-edge-delta-member-reviewer — Review Edge Delta (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-edge-delta-reviewer — Review Edge Delta using its official documentation and repository in the Observability and Analysis / Observability category. Edge Delta is a platform for creating observability pipelines, monitoring Kubernetes resources, and querying large-scale datasets. Edge Delta uniquely processes telemetry as it&#39;s created at the source to enable greater scale and cost-effectiveness.


/cncf-effectual-kcsp-reviewer — Review Effectual (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Effectual delivers full stack Kubernetes professional services including automated deployments, service mesh integration, network policy security, observability, and persistent storage solutions to streamline, secure, and scale container environments. Our experience and expertise help customers accelerate adoption, optimize workloads, and manage Kubernetes platforms with proven design and implementation approaches.


/cncf-effectual-member-reviewer — Review Effectual (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-effodio-member-reviewer — Review Effodio (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-eidp-member-reviewer — Review EIDP (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ejbca-community-reviewer — Review EJBCA Community using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Issue certificates for all your Kubernetes workloads with the open-source public key infrastructure (PKI) and certificate authority (CA) software EJBCA.


/cncf-ekuiper-reviewer — Review eKuiper using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-elastic-apm-reviewer — Review Elastic APM using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-elastic-member-reviewer — Review Elastic (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-elastic-reviewer — Review Elastic using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-elastiflow-member-reviewer — Review ElastiFlow (member) using its official documentation and repository in the CNCF Members / Silver category. ElastiFlow brings network observability and insights to OpenTelemetry.


/cncf-elastiflow-reviewer — Review ElastiFlow using its official documentation and repository in the Observability and Analysis / Observability category. ElastiFlow brings network observability and insights to OpenTelemetry.


/cncf-elastisys-kcntp-reviewer — Review Elastisys (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Elastisys is a Kubernetes and Cloud Native Training Partner as well as a Linux Foundation Authorized Training Partner. We offer the official Linux Foundation courses: Kubernetes Administration (LFS458), Kubernetes for App Developers (LFD459) and Kubernetes Security Fundamentals (LFS460), both on-site and remotely. These are recommended if you want to complete the CKA, CKAD, or CKS certifications. Additionally, we offer custom courses in cloud native technologies to suit your organization&#39;s needs.


/cncf-elastisys-kcsp-reviewer — Review Elastisys (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Elastisys enables organizations to run mission-critical software securely, at scale, and in a sustainable way. With deep expertise in cloud technology and cybersecurity, Elastisys strengthens digital resilience for organizations in both the public and private sector. We provide expert DevOps and cloud-native consulting — from strategy to execution.


/cncf-elastisys-member-reviewer — Review Elastisys (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-elastisys-welkin-based-on-cluster-api-reviewer — Review Elastisys Welkin® based on Cluster API using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Welkin is a security-hardened Kubernetes distribution. It is designed to meet the tough demands on security and compliance of regulated industries.


/cncf-elastisys-welkin-based-on-kubespray-reviewer — Review Elastisys Welkin® based on Kubespray using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Welkin is a security-hardened Kubernetes distribution. It is designed to meet the tough demands on security and compliance of regulated industries. The distribution includes logging, monitoring, single-sign-on integration with common identity providers, vulnerability scanning, intrusion detection, strict security hardening of the entire platform, and additional safeguards.


/cncf-elastx-kcsp-reviewer — Review Elastx (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Elastx is a sovereign cloud provider from Sweden, delivering CNCF Certified Kubernetes platforms with full lifecycle management and automation. Our Cloud Architects act as trusted advisors — presenting the full potential of our platform, and guiding customers from concept to production while Our CloudOps Engineers provide hands-on support. Built on open standards, powered by OpenStack and Kubernetes — designed for sustainable, vendor-free cloud operations in Europe.


/cncf-elastx-member-reviewer — Review Elastx (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-elastx-private-kubernetes-reviewer — Review Elastx Private Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Elastx delivers a fully managed, production-grade Kubernetes platform built for high availability, performance, and digital sovereignty. Running distributed across three data centers in Sweden - ensuring full control, data residency, and security compliance.


/cncf-elfconv-reviewer — Review elfconv using its official documentation and repository in the Wasm / Tooling category. An AOT binary translator that converts Linux/ELF binaries to WebAssembly


/cncf-embrace-member-reviewer — Review Embrace (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-embrace-reviewer — Review Embrace using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-emissary-ingress-reviewer — Review Emissary-Ingress using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-emqx-reviewer — Review EMQX using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-emscripten-reviewer — Review emscripten using its official documentation and repository in the Wasm / Tooling category.


/cncf-encore-reviewer — Review Encore using its official documentation and repository in the Serverless / Framework category.


/cncf-endee-reviewer — Review Endee using its official documentation and repository in the App Definition and Development / Database category. High-performance open-source vector database for AI search, RAG, semantic search, and hybrid retrieval.


/cncf-endress-plus-hauser-supporter-reviewer — Review Endress+Hauser (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category. Endress+Hauser is a global leader in measurement and automation technology for process and laboratory applications.


/cncf-enroute-onestep-ingress-reviewer — Review EnRoute OneStep Ingress using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-entigo-kcntp-reviewer — Review Entigo (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Simplify your software delivery and operations.


/cncf-entigo-kcsp-reviewer — Review Entigo (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Improve your software delivery and operations by standardising and automating with Kubernetes.


/cncf-entigo-member-reviewer — Review Entigo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-enum-member-reviewer — Review enum (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-env-zero-member-reviewer — Review env zero (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-envoy-ai-gateway-reviewer — Review Envoy-ai-gateway using its official documentation and repository in the AI Native Infra / Gateway category. Manages Unified Access to Generative AI Services built on Envoy Gateway


/cncf-envoy-gateway-reviewer — Review Envoy Gateway using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. Part of the Envoy project, Envoy Gateway is the control plane for dynamically managing Envoy Proxy, aimed at significantly decreasing the barrier to entry when using it for Gateway (sometimes known as &quot;north-south&quot;) use cases.


/cncf-envoy-reviewer — Review Envoy using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-envoy-wasm-reviewer — Review Envoy (Wasm) using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-eos-reviewer — Review EOS using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-epam-systems-member-reviewer — Review EPAM Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-epsagon-reviewer — Review Epsagon using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-epsilla-reviewer — Review Epsilla using its official documentation and repository in the App Definition and Development / Database category. 10x faster vector database and one-stop RAGaaS platform for building LLM applications


/cncf-epsio-member-reviewer — Review Epsio (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-eraser-reviewer — Review Eraser using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Eraser uses vulnerability data to remove non-running images from all Kubernetes nodes in a cluster.


/cncf-ericsson-cloud-container-distribution-reviewer — Review Ericsson Cloud Container Distribution using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Ericsson Cloud Container Distribution provides container management and orchestration for the Ericsson Telco applications that have been adopted to the Cloud Native based Architecture and run in a container environment.


/cncf-ericsson-member-reviewer — Review Ericsson (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-erlang-ecosystem-foundation-member-reviewer — Review Erlang Ecosystem Foundation (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-etcd-reviewer — Review etcd using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category.


/cncf-eunomia-bpf-reviewer — Review eunomia-bpf using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-eventmesh-reviewer — Review EventMesh using its official documentation and repository in the Serverless / Framework category. EventMesh is a new generation serverless event middleware for building distributed event-driven applications.


/cncf-everquote-supporter-reviewer — Review EverQuote (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-excellion-sdn-bhd-kcsp-reviewer — Review Excellion Sdn Bhd (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Excellion Sdn Bhd is a Malaysia-based Specialised System Integrator (SSI) designing secure, cloud-native infrastructure through Enterprise Kubernetes, advanced observability, and AI-driven engineering. We champion Malaysia’s digital sovereignty by enabling public and private enterprises to adopt resilient, open-standard architectures for their mission-critical projects &amp; implementations, essential to the nation&#39;s digital ecosystem.


/cncf-excellion-sdn-bhd-member-reviewer — Review Excellion Sdn Bhd (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-exein-member-reviewer — Review Exein (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-exoscale-member-reviewer — Review Exoscale (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-exoscale-scalable-kubernetes-service-reviewer — Review Exoscale Scalable Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Scalable, On-demand Kubernetes Clusters on a privacy minded public cloud to host from simple applications to complex architectures.  Deploy a production ready cluster in 90 seconds and manage it with a simple web portal, CLI, API or your choice of tools like Terraform.


/cncf-exostellar-member-reviewer — Review Exostellar (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-expontech-reviewer — Review ExponTech using its official documentation and repository in the Runtime / Cloud Native Storage category. ExponTech is a provider of overall data infrastructure solutions based on a new generation of distributed architecture.


/cncf-external-secrets-reviewer — Review external-secrets using its official documentation and repository in the Provisioning / Security &amp; Compliance category. External Secrets Operator reads information from a third-party service like AWS Secrets Manager and automatically injects the values as Kubernetes Secrets.


/cncf-extism-reviewer — Review Extism using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-f5-member-reviewer — Review F5 (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-f5-reviewer — Review F5 using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-fabedge-reviewer — Review FabEdge using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-fabric8-kubernetes-client-reviewer — Review Fabric8 Kubernetes Client using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. The Fabric8 Kubernetes Client is a Java client that is used to interact with Kubernetes clusters through the Kubernetes API. It&#39;s one of the main building blocks for other CNCF projects like the Java Operator SDK from the Operator Framework.


/cncf-factory-member-reviewer — Review Factory (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-fairwinds-insights-reviewer — Review Fairwinds Insights using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-fairwinds-kcsp-reviewer — Review Fairwinds (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Fairwinds, the Kubernetes enablement company, offers services, open source, and software to help organizations run secure, efficient and reliable Kubernetes infrastructure.


/cncf-fairwinds-member-reviewer — Review Fairwinds (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-falco-reviewer — Review Falco using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-falcon-reviewer — Review Falcon using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-falkordb-reviewer — Review FalkorDB using its official documentation and repository in the AI Agent / Knowledge Graph category. A super fast Graph Database uses GraphBLAS under the hood for its sparse adjacency matrix graph representation. Our goal is to provide the best Knowledge Graph for LLM (GraphRAG).


/cncf-fastly-compute-at-edge-reviewer — Review Fastly Compute@Edge using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-fd-io-reviewer — Review FD.io using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-feast-reviewer — Review FEAST using its official documentation and repository in the Data / Data Architecture category. Feature Store to manage machine learning features.


/cncf-featurehub-reviewer — Review FeatureHub using its official documentation and repository in the Observability and Analysis / Feature Flagging category. FeatureHub is a Cloud Native Feature Flags, Remote Configuration and A/B Testing Platform suitable for mobile, web and server applications and available with the variety of SDKs. Self-hosted open source or SaaS options are available.


/cncf-ferymon-cloud-reviewer — Review Ferymon Cloud using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-ffmpeg-reviewer — Review FFmpeg using its official documentation and repository in the Wasm / AI/Machine Learning category.


/cncf-fidelity-investments-member-reviewer — Review Fidelity Investments (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-filecoin-reviewer — Review Filecoin using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-finout-member-reviewer — Review Finout (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-firecracker-reviewer — Review Firecracker using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-firecrawl-reviewer — Review Firecrawl using its official documentation and repository in the AI Agent / Agent Tool category. The API to search, scrape, and interact with the web for AI


/cncf-firefly-member-reviewer — Review Firefly (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-fission-reviewer — Review Fission using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-flagger-reviewer — Review Flagger using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-flagsmith-reviewer — Review Flagsmith using its official documentation and repository in the Observability and Analysis / Feature Flagging category.


/cncf-flannel-reviewer — Review Flannel using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-flatcar-container-linux-reviewer — Review Flatcar Container Linux using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. A community Linux distribution designed for container workloads, with high security and low maintenance


/cncf-flink-reviewer — Review Flink using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-flipt-reviewer — Review Flipt using its official documentation and repository in the Observability and Analysis / Feature Flagging category. Open source, self-hosted, developer first, feature flagging and dynamic configuration service


/cncf-flogo-reviewer — Review Flogo using its official documentation and repository in the Serverless / Framework category.


/cncf-flomesh-service-mesh-fsm-reviewer — Review Flomesh Service Mesh (FSM) using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Lightweight SMI compatible service mesh for Kubernetes east-west and north-south traffic management, uses ebpf for layer4 and pipy proxy for layer7 traffic management, support multi cluster network.


/cncf-flowmill-reviewer — Review Flowmill using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-flows-network-reviewer — Review flows.network using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-flox-member-reviewer — Review Flox (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-fluentd-reviewer — Review Fluentd using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-fluid-reviewer — Review Fluid using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Fluid is an orchestration platform for elastic data abstraction and acceleration in cloud native environment.


/cncf-fluvio-reviewer — Review Fluvio using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-flux-reviewer — Review Flux using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-fongcon-kcsp-reviewer — Review Fongcon (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As an NVIDIA Elite Partner, we deliver expert Kubernetes consulting, implementation, and training services specifically optimized for GPU-accelerated platforms to ensure reliable AI workload orchestration. We empower enterprises to bridge the gap between high-performance hardware and cloud-native ecosystems through professional support and tailored operational excellence.


/cncf-fongcon-member-reviewer — Review Fongcon (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-fonio-reviewer — Review Fonio using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-foreman-reviewer — Review Foreman using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-foresight-reviewer — Review Foresight using its official documentation and repository in the Observability and Analysis / Observability category. Monitor workflow and test trends over time, assess the risk of code changes, deal with flaky tests.


/cncf-form3-supporter-reviewer — Review Form3 (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-formal-member-reviewer — Review Formal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-fortio-reviewer — Review Fortio using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Fortio is an open source load testing library, command-line tool, and server application written in Go. Originated as Istio&#39;s performance characterization tool (in particular proxy overhead), Fortio has evolved into a versatile tool for load testing and performance benchmarking of HTTP, gRPC, and TCP services. Among other options it enables users to generate a specified query-per-second (QPS) load and records detailed latency histograms, facilitating the analysis of performance metrics. The server component offers a straightforward web UI and REST API, allowing users to initiate tests and visualize results through graphical representations. Beyond load testing, Fortio provides server-side features akin to httpbin, including request echoing with headers, configurable latency or error responses, TCP echoing and proxying, HTTP fan-out, and support for gRPC echo and health checks. These capabilities make Fortio a versatile tool for debugging and testing high-performance services in cloud-native environments. As an organization set of git repositories, is also a growing set of reusable libraries for writing Cloud Native Go code and CLIs.


/cncf-fossa-member-reviewer — Review FOSSA (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-fossa-reviewer — Review FOSSA using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-fossid-reviewer — Review FOSSID using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-foundationdb-reviewer — Review FoundationDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-framework-member-reviewer — Review Framework (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-frontier-kcsp-reviewer — Review Frontier (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Operate predictably. Innovate safely. It&#39;s not just a tagline, it&#39;s our engineering philosophy, and it&#39;s especially relevant when working with  advanced technologies like Kubernetes. Our architecture, implementation and consulting services provide container platforms that are reliable, secure and cost-effective. Our Kubernetes support services give our clients confidence in the platform and the automation, just as you would  expect from us.


/cncf-frontier-member-reviewer — Review Frontier (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-fugue-reviewer — Review Fugue using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Fugue empowers cloud engineering and security teams to move faster with confidence in their cloud security. Use Fugue to secure the entire development lifecycle—from infrastructure as code through CI/CD and runtime—using a unified policy engine powered by the Open Policy Agent (OPA), the open source standard for policy as code.


/cncf-fujitsu-member-reviewer — Review Fujitsu (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-fullstacks-kcntp-reviewer — Review FullStackS (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. We at FullStackS GmbH live Cloud, Automation, Container &amp; Kubernetes, CI/CD &amp; Application Performance Monitoring! We offer a wide range of trainings regarding Cloud Native solutions as well as best practices of Application Modernization.


/cncf-fullstacks-kcsp-reviewer — Review FullStackS (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We live Cloud, Automation, Container &amp; Kubernetes, CI/CD &amp; Application Performance Monitoring! We have an holistic approach to provide our customers a Application-Delivery-Platform.


/cncf-fullstacks-member-reviewer — Review FullStackS (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-g-research-member-reviewer — Review G-Research (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gaia-information-technology-member-reviewer — Review Gaia Information Technology  (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gaia-kcntp-reviewer — Review Gaia (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Gaia, a multi-cloud solutions provider, offers Kubernetes implementation and configuration services, application container migrations, and professional IT training.   We help our customers adopt DevSecOps, build CI/CD, and migrate applications to Kubernetes.


/cncf-gaia-kcsp-reviewer — Review Gaia (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Gaia, a multi-cloud solutions provider, offers Kubernetes implementation and configuration services, application container migrations, and professional IT training. We help our customers adopt DevSecOps, build CI/CD, and migrate applications to Kubernetes.


/cncf-gardener-reviewer — Review Gardener using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. The Gardener implements automated management and operation of Kubernetes clusters as a service and aims to support that service on multiple  Cloud providers.


/cncf-gcore-managed-kubernetes-reviewer — Review Gcore Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Gcore Managed Kubernetes is a service that lets you deploy Kubernetes clusters without the complexities of handling the control plane and containerized infrastructure.


/cncf-gcore-member-reviewer — Review Gcore (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gear-reviewer — Review GEAR using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-geeks-solutions-kcsp-reviewer — Review Geeks Solutions (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Revolutionize your infrastructure with our cutting-edge Kubernetes services. Seamlessly manage, deploy, and scale applications for unparalleled efficiency and agility in the ever-evolving digital landscape.


/cncf-geeks-solutions-member-reviewer — Review Geeks Solutions (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gefyra-reviewer — Review Gefyra using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Gefyra runs local code in any Kubernetes cluster without the build and push cycle. It overlays containers in the cluster making code changes immediately available.


/cncf-genode-reviewer — Review Genode using its official documentation and repository in the Wasm / Edge/Bare metal category.


/cncf-giant-swarm-kcsp-reviewer — Review Giant Swarm (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Giant Swarm provides fully-managed Kubernetes Clusters in your location of choice, so you can focus on your product.


/cncf-giant-swarm-managed-kubernetes-reviewer — Review Giant Swarm Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. The Giant Swarm platform enables users to simply and rapidly create and use 24/7 managed Kubernetes clusters on-demand.


/cncf-giant-swarm-member-reviewer — Review Giant Swarm (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gientech-container-cloud-platform-gtcp-reviewer — Review GienTech Container Cloud Platform (GTCP) using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. GienTech Container Cloud Platform (GTCP) provides financial level, high availability, and high security container service platform.  GTCP fully follows the conformance with Kubernetes and is flexible to extend more other services with GienTech PAAS products.


/cncf-gientech-kcsp-reviewer — Review GienTech (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Gientech can help customers understand and deploy Kubernetes container technology deeply. We provides customized Kubernetes services including consulting, training, implementation and technical support.


/cncf-gientech-member-reviewer — Review GienTech (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gim-uemoa-member-reviewer — Review GIM UEMOA (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gitguardian-reviewer — Review GitGuardian using its official documentation and repository in the Provisioning / Security &amp; Compliance category. GitGuardian is the code security platform for the DevOps generation that offers automated Secrets Detection, Infra as Code Security,  and Honeytoken capabilities, facilitating a Secure Software Development Lifecycle for Dev, Sec, and Ops teams.


/cncf-github-actions-reviewer — Review GitHub Actions using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. GitHub Actions makes it easy to automate all your software workflows, now with world-class CI/CD. Build, test, and deploy your code right from GitHub. Make code reviews, branch management, and issue triaging work the way you want.


/cncf-gitlab-reviewer — Review GitLab using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-gitness-reviewer — Review Gitness using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-gitpod-member-reviewer — Review Gitpod (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gitpod-reviewer — Review Gitpod using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-glasnostic-reviewer — Review Glasnostic using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Glasnostic makes modern cloud applications resilient by shaping how systems interact, automatically and in real-time.


/cncf-glassflow-reviewer — Review GlassFlow using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Open-source streaming ETL for ClickHouse with deduplication, joins, and transforms over Kafka and other sources.


/cncf-gloo-mesh-reviewer — Review Gloo Mesh using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-gloo-reviewer — Review Gloo using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-gluster-reviewer — Review Gluster using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-gmx-supporter-reviewer — Review GMX (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-go-feature-flag-reviewer — Review GO Feature Flag using its official documentation and repository in the Observability and Analysis / Feature Flagging category.


/cncf-go-zero-reviewer — Review go-zero using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. go-zero is a web and rpc framework written in Go. It&#39;s born to ensure the stability of the busy sites with resilient design. Builtin goctl greatly improves the development productivity.


/cncf-gocd-reviewer — Review GoCD using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-gocrane-reviewer — Review gocrane using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. Crane (FinOps Crane) is an opensource project which manages cloud resource on Kubernetes stack, it is inspired by FinOps concepts.


/cncf-godel-scheduler-reviewer — Review Godel-Scheduler using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Godel-Scheduler is a unified scheduler for both online and offline tasks.


/cncf-gofr-reviewer — Review GoFr using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category.


/cncf-golang-reviewer — Review Golang using its official documentation and repository in the Wasm / Languages category. Compiled language to Wasm


/cncf-golden-gate-university-member-reviewer — Review Golden Gate University (member) using its official documentation and repository in the CNCF Members / Academic category.


/cncf-goldilocks-reviewer — Review Goldilocks using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-goldman-sachs-member-reviewer — Review Goldman Sachs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-golem-cloud-wasm-reviewer — Review Golem Cloud (Wasm) using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-gonzo-reviewer — Review Gonzo using its official documentation and repository in the Observability and Analysis / Observability category. Gonzo is a powerful, real-time log analysis terminal UI inspired by k9s. Analyze log streams with beautiful charts, AI-powered insights, and advanced filtering – all from your terminal.


/cncf-google-cloud-build-reviewer — Review Google Cloud Build using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Cloud Build lets you build software quickly across all languages. Get complete control over defining custom workflows for building, testing, and deploying across multiple environments such as VMs, serverless, Kubernetes, or Firebase.


/cncf-google-cloud-dataflow-reviewer — Review Google Cloud Dataflow using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Cloud Dataflow is a fully-managed service for transforming and enriching data in stream (real time) and batch (historical) modes with equal reliability and expressiveness -- no more complex workarounds or compromises needed. And with its serverless approach to resource provisioning and management, you have access to virtually limitless capacity to solve your biggest data processing challenges, while paying only for what you use.


/cncf-google-cloud-member-reviewer — Review Google Cloud (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-google-cloud-run-reviewer — Review Google Cloud Run using its official documentation and repository in the Serverless / Hosted Platform category. Cloud Run allows you to build and deploy scalable containerized apps written in any language (including Go, Python, Java, Node.js, .NET, and Ruby) on a fully managed serverless platform.


/cncf-google-container-registry-reviewer — Review Google Container Registry using its official documentation and repository in the Provisioning / Container Registry category. Container Registry is a single place for your team to manage Docker images, perform vulnerability analysis, and decide who can access what with fine-grained access control.


/cncf-google-kubernetes-engine-gke-reviewer — Review Google Kubernetes Engine (GKE) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. GKE is an enterprise-grade platform for containerized applications, including stateful and stateless, AI and ML, Linux and Windows, complex and simple web apps, API, and backend services.


/cncf-google-persistent-disk-reviewer — Review Google Persistent Disk using its official documentation and repository in the Runtime / Cloud Native Storage category. Google Persistent Disk is durable and high performance block storage for the Google Cloud Platform. Persistent Disk provides SSD and HDD storage which can be attached to instances running in either Google Compute Engine or Google Kubernetes Engine. Storage volumes can be transparently resized, quickly backed up, and offer the ability to support simultaneous readers.


/cncf-google-stackdriver-reviewer — Review Google Stackdriver using its official documentation and repository in the Observability and Analysis / Observability category. Stackdriver aggregates metrics, logs, and events from infrastructure, giving developers and operators a rich set of observable signals that speed root-cause analysis and reduce mean time to resolution (MTTR).


/cncf-goose-reviewer — Review Goose using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-gradle-build-tool-reviewer — Review Gradle Build Tool using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Gradle is a build automation tool for multi-language software development. It controls the development process in the tasks of compilation and packaging to testing, deployment, and publishing. Supported languages include Java (Kotlin, Groovy, Scala), C/C++, and JavaScript.


/cncf-grafana-labs-member-reviewer — Review Grafana Labs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-grafana-loki-reviewer — Review Grafana Loki using its official documentation and repository in the Observability and Analysis / Observability category. Loki is a horizontally-scalable, highly-available, multi-tenant log aggregation system inspired by Prometheus. It is designed to be very cost effective and easy to operate. It does not index the contents of the logs, but rather a set of labels for each log stream.


/cncf-grafana-mimir-reviewer — Review Grafana Mimir using its official documentation and repository in the Observability and Analysis / Observability category. Grafana Mimir lets you scale to 1 billion metrics and beyond, with high availability, multi-tenancy, durable storage, and blazing fast query performance over long periods of time.


/cncf-grafana-pyroscope-reviewer — Review Grafana Pyroscope using its official documentation and repository in the Observability and Analysis / Observability category. Grafana Pyroscope is an open source continuous profiling database that provides fast, scalable, highly available, and efficient storage and querying. This helps you get a better understanding of resource usage in your applications down to the line number.


/cncf-grafana-reviewer — Review Grafana using its official documentation and repository in the Observability and Analysis / Observability category. Grafana allows you to query, visualize, alert on and understand your metrics no matter where they are stored. Create, explore, and share beautiful dashboards with your team and foster a data driven culture.


/cncf-grafana-tempo-reviewer — Review Grafana Tempo using its official documentation and repository in the Observability and Analysis / Observability category. Grafana Tempo is an open source, easy-to-use and high-scale distributed tracing backend. Tempo is cost-efficient, requiring only object storage to operate, and is deeply integrated with Grafana, Prometheus, and Loki. Tempo can be used with any of the open source tracing protocols, including Jaeger, Zipkin, and OpenTelemetry.


/cncf-grafeas-reviewer — Review Grafeas using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-grain-reviewer — Review Grain using its official documentation and repository in the Wasm / Languages category. WASM focused languages


/cncf-granulate-reviewer — Review Granulate using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-graphite-reviewer — Review Graphite using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-graphscope-reviewer — Review GraphScope using its official documentation and repository in the App Definition and Development / Database category. A One-Stop Large-Scale Graph Computing System from Alibaba


/cncf-gravitee-io-reviewer — Review Gravitee.io using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-graylog-reviewer — Review Graylog using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-gremlin-reviewer — Review Gremlin using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-grepr-member-reviewer — Review Grepr (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-greptime-member-reviewer — Review Greptime (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-greptimedb-reviewer — Review GreptimeDB using its official documentation and repository in the App Definition and Development / Database category. GreptimeDB is an open-source, cloud-native time series database which also has powerful analytical features


/cncf-greymatter-io-reviewer — Review greymatter.io using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-groundcover-member-reviewer — Review Groundcover (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-groundcover-reviewer — Review groundcover using its official documentation and repository in the Observability and Analysis / Observability category. groundcover is a cloud native observability platform powered by eBPF. It runs inside the customer’s cloud and provides complete visibility into applications, infrastructure, networks and AI systems without operational overhead.


/cncf-grpc-reviewer — Review gRPC using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. A high performance, open source universal RPC framework.


/cncf-grype-reviewer — Review Grype using its official documentation and repository in the Provisioning / Security &amp; Compliance category. A vulnerability scanner for container images and filesystems


/cncf-gthulhu-reviewer — Review Gthulhu using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Gthulhu optimizes cloud-native workloads using the Linux Scheduler Extension for different application scenarios.


/cncf-guance-cloud-reviewer — Review Guance Cloud using its official documentation and repository in the Observability and Analysis / Observability category. Guance Cloud offers a full-stack observability platform for cloud-based workflows.


/cncf-guardicore-centra-reviewer — Review Guardicore Centra using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-guardrails-ai-reviewer — Review Guardrails AI using its official documentation and repository in the AI Agent / Guardrail category. Adding guardrails to large language models.


/cncf-guida-kcsp-reviewer — Review Guida (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Guida builds and operates cloud native platforms based on Kubernetes so that software developers can focus on building innovative software.


/cncf-guida-member-reviewer — Review Guida (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-guidance-reviewer — Review Guidance using its official documentation and repository in the AI Agent / Structured Output category. A guidance language for controlling large language models.


/cncf-guidewire-member-reviewer — Review Guidewire (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-gvisor-reviewer — Review gVisor using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-h3c-technologies-member-reviewer — Review H3C Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hadoop-hdfs-reviewer — Review Hadoop HDFS using its official documentation and repository in the Data / Data Architecture category. Open source framework works by rapidly transferring data between nodes. It&#39;s often used by companies who need to handle and store big data.


/cncf-hami-reviewer — Review HAMi using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Heterogeneous AI Computing Virtualization Middleware


/cncf-hammerspace-member-reviewer — Review Hammerspace (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hango-reviewer — Review Hango using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-haproxy-reviewer — Review HAProxy using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-haproxy-technologies-member-reviewer — Review HAProxy Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-harbor-reviewer — Review Harbor using its official documentation and repository in the Provisioning / Container Registry category.


/cncf-harbor-wasm-reviewer — Review Harbor (Wasm) using its official documentation and repository in the Wasm / Packaging, Registries &amp; Application Delivery category.


/cncf-harmonycloud-kcntp-reviewer — Review HarmonyCloud (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category.


/cncf-harmonycloud-kcsp-reviewer — Review HarmonyCloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Harmony Cloud has many Kubernetes certified engineers and we provide a full range of Kubernetes consulting, professional services and training support for our customers who are committed to the construction of enterprise modernization infrastructure.


/cncf-harmonycloud-member-reviewer — Review HarmonyCloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-harmonycloud-reviewer — Review harmonycloud using its official documentation and repository in the Observability and Analysis / Observability category. Provide eBPF-based cloud native monitoring platform, which aim to achieve efficient triage of issues


/cncf-harness-io-reviewer — Review Harness.io using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-harness-member-reviewer — Review Harness (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-haskell-reviewer — Review Haskell using its official documentation and repository in the Wasm / Languages category. Purely functional language with wasm backend


/cncf-hasura-graphql-engine-reviewer — Review Hasura GraphQL Engine using its official documentation and repository in the Serverless / Tools category.


/cncf-hatchet-member-reviewer — Review Hatchet (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hatchet-reviewer — Review Hatchet using its official documentation and repository in the AI Agent / Workflow Orchestration category. Run Background Tasks at Scale


/cncf-hawkstack-member-reviewer — Review HawkStack (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-haystack-reviewer — Review Haystack using its official documentation and repository in the AI Agent / RAG category. Open-source AI orchestration framework for building context-engineered, production-ready LLM applications. Design modular pipelines and agent workflows with explicit control over retrieval, routing, memory, and generation. Built for scalable agents, RAG, multimodal applications, semantic search, and conversational systems.


/cncf-hazelcast-imdg-reviewer — Review Hazelcast IMDG using its official documentation and repository in the App Definition and Development / Database category.


/cncf-headlamp-reviewer — Review Headlamp using its official documentation and repository in the Observability and Analysis / Observability category. Extensible open source multi-cluster Kubernetes user interface


/cncf-hedgehog-member-reviewer — Review Hedgehog (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-helios-reviewer — Review Helios using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-helm-reviewer — Review Helm using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-helmwave-reviewer — Review Helmwave using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-heroku-member-reviewer — Review Heroku (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-heroku-reviewer — Review Heroku using its official documentation and repository in the Platform / PaaS/Container Service category. Developers, teams, and businesses of all sizes use Heroku to deploy, manage, and scale apps.


/cncf-hertzbeat-reviewer — Review HertzBeat using its official documentation and repository in the Observability and Analysis / Observability category. Apache HertzBeat (Incubating) is an open-source, real-time monitoring system with custom monitoring, high performance cluster and agentless capabilities.


/cncf-hexa-reviewer — Review Hexa using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-hidora-kcsp-reviewer — Review Hidora (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Hidora provides fully managed Kubernetes services, supporting cluster setup, migration, and operations. Hosted in Swiss datacenters, our solutions ensure sovereignty, reliability, and expert local support.


/cncf-hidora-member-reviewer — Review Hidora (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hightech-payment-systems-supporter-reviewer — Review Hightech Payment Systems (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-higress-reviewer — Review Higress using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. An AI-native API gateway built on Envoy for ingress, microservice, and LLM traffic


/cncf-hikube-reviewer — Review Hikube using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Hikube’s Managed Kubernetes provides a fully managed, multi-datacenter Kubernetes environment distributed across three Swiss sites, ensuring cluster redundancy, automated updates, integrated monitoring, and secure networking for production-grade workloads.


/cncf-hitachi-member-reviewer — Review Hitachi (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-hitachi-reviewer — Review Hitachi using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-holmesgpt-reviewer — Review HolmesGPT using its official documentation and repository in the Observability and Analysis / Observability category. HolmesGPT is an AI agent that automates cloud-native troubleshooting, bridging knowledge gaps by investigating alerts, executing runbooks, and correlating observability data in cloud-native platforms.


/cncf-homestar-reviewer — Review homestar using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-honeybadger-reviewer — Review Honeybadger using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-honeycomb-member-reviewer — Review Honeycomb (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-honeycomb-reviewer — Review Honeycomb using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-hoop-dev-member-reviewer — Review hoop.dev (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-horovod-reviewer — Review Horovod using its official documentation and repository in the Training / Distributed Training category. Distributed deep learning training framework for TensorFlow, Keras, PyTorch, and Apache MXNet.


/cncf-hostersi-kcsp-reviewer — Review Hostersi (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Hostersi is a sysops/devops house that helps organizations design, implement and maintain their IT infrastructure on AWS, Azure, GCP and on-premise solutions with Kubernetes, CI/CD, HA, DR and support 24/7/365.


/cncf-hostersi-member-reviewer — Review Hostersi (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hpe-enterprise-kubernetes-platform-hks-reviewer — Review HPE Enterprise Kubernetes Platform (HKS) using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. HKS delivers an enterprise-grade Kubernetes platform to provision and manage container workload across your hybrid cloud environments.


/cncf-hpe-ezmeral-runtime-enterprise-reviewer — Review HPE Ezmeral Runtime Enterprise using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. The HPE Ezmeral Runtime Enterprise is a secure, enterprise-grade platform to deploy cloud-native and non-cloud-native applications at scale across your data centers, multiple clouds, and at the edge.


/cncf-hpe-kcntp-reviewer — Review HPE (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. HPE is a leader in the hybrid cloud market offering Kubernetes solutions that are supported by professional services and education.


/cncf-hpe-kcsp-reviewer — Review HPE (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. HPE Advisory &amp; Professional Service provides various services from strategic advisory to design/implementation professional services related to Kubernetes and cloud native computing technologies.


/cncf-hpe-member-reviewer — Review HPE (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hpe-storage-reviewer — Review HPE Storage using its official documentation and repository in the Runtime / Cloud Native Storage category. The world’s most intelligent storage, built for our hybrid cloud world to help you unlock the full potential of your data.


/cncf-huatuo-reviewer — Review HUATUO using its official documentation and repository in the Observability and Analysis / Observability category. Kernel-wide Insight, Instant Observability, AutoTracing, Continuous Profiling for cloud-native and AI infrastructure using eBPF.


/cncf-huawei-cloud-container-engine-cce-reviewer — Review Huawei Cloud Container Engine (CCE) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Cloud Container Engine (CCE) is a high-performance, high-reliability service through which enterprises can manage containerized applications. CCE supports native Kubernetes applications and tools, allowing you to set up a container runtime environment on the cloud with ease.


/cncf-huawei-convergent-billing-system-cbs-reviewer — Review Huawei Convergent Billing System (CBS) using its official documentation and repository in the Special / Certified CNFs category. The next-generation convergent billing system(CBS) helps carriers reduce costs and improve efficiency with leading 5G monetization practices and cloud-native technologies, enabling carriers to monetize 5G productivity beyond connectivity.


/cncf-huawei-functionstage-reviewer — Review Huawei FunctionStage using its official documentation and repository in the Serverless / Hosted Platform category. FunctionStage is an event-driven function hosting and computing service, that compiles function code and configures running conditions without provisioning or managing servers, thus providing a scalable, maintenance-free, and reliable operation environment for functions. You only pay for what you use.


/cncf-huawei-kcntp-reviewer — Review Huawei (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. HUAWEI CLOUD provides customers with stable, reliable, secure, and sustainably growing cloud services. It helps large enterprises address challenges in cloud transformation and enables them to take better advantages of potential business opportunities. It also helps small- and medium-sized enterprises expand their business growth and rise to challenges.


/cncf-huawei-kcsp-reviewer — Review Huawei (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. HUAWEI CLOUD provides customers with stable, reliable, secure, and sustainably growing cloud services. It helps large enterprises address challenges in cloud transformation and enables them to take better advantages of potential business opportunities. It also helps small- and medium-sized enterprises expand their business growth and rise to challenges.


/cncf-huawei-member-reviewer — Review Huawei (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-huawei-reviewer — Review Huawei using its official documentation and repository in the Runtime / Cloud Native Storage category. Ubiquitous Storage Unleashes Data Power


/cncf-hubble-reviewer — Review Hubble using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-humio-reviewer — Review Humio using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-hunter-strategy-kcsp-reviewer — Review Hunter Strategy (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Hunter provides Kubernetes expertise across on-prem and cloud environments, enabling secure deployments for the Department of Defense, Federal agencies and commercial clients. We can work with existing teams implementing real world best practices and leveling up their skillsets, audit existing implementations, and work with you on securing your environment!


/cncf-hunter-strategy-member-reviewer — Review Hunter Strategy (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hush-security-member-reviewer — Review Hush Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hwameistor-reviewer — Review HwameiStor using its official documentation and repository in the Runtime / Cloud Native Storage category. Hwameistor is an HA local storage system for cloud-native stateful workloads


/cncf-hydrolix-member-reviewer — Review Hydrolix (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hyground-member-reviewer — Review Hyground (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-hyperlight-reviewer — Review Hyperlight using its official documentation and repository in the Runtime / Container Runtime category. A lightweight, secure container runtime solution designed for modern cloud-native workloads


/cncf-hyperopt-reviewer — Review Hyperopt using its official documentation and repository in the Training / Post Training category. Distributed Asynchronous Hyper-parameter Optimization


/cncf-hyscale-reviewer — Review HyScale using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. An app-centric abstraction framework over K8s and an enterprise platform for accelerating deployments to multi-cloud Kubernetes. Offers self-service app deployments, container sprawl management across clusters, and DevSecOps.


/cncf-hyve-managed-hosting-member-reviewer — Review Hyve Managed Hosting (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-i-cubed-systems-inc-supporter-reviewer — Review i Cubed Systems, Inc. (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-ibm-cloud-code-engine-reviewer — Review IBM Cloud Code Engine using its official documentation and repository in the Serverless / Hosted Platform category. Host all your Cloud Native applications in one platform. You&#39;ll get all the benefits of serverless (no servers, auto-scaling) along with batch processing and only pay for what you use.


/cncf-ibm-cloud-container-registry-reviewer — Review IBM Cloud Container Registry using its official documentation and repository in the Provisioning / Container Registry category. Detect vulnerabilities before images are ever deployed to containers. Store and distribute Docker images in your managed private registry.


/cncf-ibm-cloud-functions-reviewer — Review IBM Cloud Functions using its official documentation and repository in the Serverless / Hosted Platform category. Run your application code without servers, scale it automatically, and pay nothing when it&#39;s not in use.


/cncf-ibm-cloud-kcsp-reviewer — Review IBM Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. IBM Cloud provides Kubernetes services including training, consulting, support, and development on our Kubernetes Service delivering an intuitive user experience, and built-in security and isolation to enable rapid delivery of applications all while empowered by a full stack cloud platform including data, containers, AI, IoT, and blockchain. Get started today through our on-boarding resources and industry expertise.


/cncf-ibm-cloud-kubernetes-service-reviewer — Review IBM Cloud Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. IBM Cloud Kubernetes Service is a managed offering to create your own Kubernetes cluster of compute hosts to deploy and manage containerized apps on IBM Cloud.


/cncf-ibm-db2-reviewer — Review IBM Db2 using its official documentation and repository in the App Definition and Development / Database category. IBM Db2® is a family of hybrid data management products offering a complete suite of AI empowered capabilities designed to help you manage both structured and unstructured data on premises as well as in private and public cloud environments. Db2 is built on an intelligent common SQL engine designed for scalability and flexibility.


/cncf-ibm-member-reviewer — Review IBM (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-ibm-storage-reviewer — Review IBM Storage using its official documentation and repository in the Runtime / Cloud Native Storage category. IBM Storage for containers, Kubernetes, and Red Hat OpenShift delivers native cloud acceleration to help you build powerful, agile, and persistent storage for private cloud environments.


/cncf-ibm-turbonomic-reviewer — Review IBM Turbonomic using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. IBM Turbonomic targets and optimizes all Kubernetes and OpenShift environments on-prem and in the Cloud in both managed and un-managed clusters. Providing continuous optimization for performance, efficiency and compliance, sustainably and at the lowest cost. By scaling in and out based on historical utilization,  moving containerized workloads to help prevent pod evictions and node congestion, SLO scaling to help assure application performance and scaling in clusters to save on power and cost,  scaling out clusters proactively for performance before hitting capacity thresholds.  Also providing customized planning for future growth by running simulations against your environments.


/cncf-icinga-reviewer — Review Icinga using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-icon-business-systems-kcsp-reviewer — Review Icon Business Systems (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Embark on a transformative journey with Icon Business Systems Ltd., the beacon of SDLC innovation. Our mastery of Kubernetes services ensures effortless application modernization and containerization, while our expertise in CI/CD pipelines, AI modeling, LLM, RAG, and web crawling illuminates the path to unparalleled efficiency and creativity.


/cncf-icon-business-systems-member-reviewer — Review Icon Business Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-icubed-kcsp-reviewer — Review iCubed (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. iCubed provides strategy, transformation and managed services across Application Modernisation, DevOps and Kubernetes.


/cncf-icubed-member-reviewer — Review iCubed (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-idem-project-reviewer — Review Idem Project using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Idem reduces your cloud configuration to data – for you! Instead of Infrastructure as Code, Idem delivers Infrastructure as Data. The cloud becomes simplified, easier to maintain, easier to discover, and easier to use.


/cncf-if-information-systems-kcsp-reviewer — Review IF Information Systems (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide expert Kubernetes services, ensuring seamless container orchestration, scalability, and high availability for enterprise systems. Our team specializes in optimizing cloud-native infrastructure for maximum efficiency and reliability.


/cncf-if-information-systems-member-reviewer — Review IF Information Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-iguazio-reviewer — Review iguazio using its official documentation and repository in the App Definition and Development / Database category.


/cncf-iherb-supporter-reviewer — Review iHerb (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-iits-consulting-kcsp-reviewer — Review iits-consulting (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Put your trust in one of Germany’s most successful bootstrapped IT consulting companies of the past decade. iits-consulting delivers sovereign EU-cloud solutions, expert Time &amp; Material consultants, and CloudOps as a Service with KumoOps, our unified Developer Platform and Operations solution. From Kubernetes to AI-powered innovation, we are your trusted partner for scalable and future-ready cloud solutions.


/cncf-iits-consulting-member-reviewer — Review iits-consulting (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ilki-kcntp-reviewer — Review ILKI (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. ILKI is an IT consulting, architecture, and training company. ILKI helps companies to accelerate and secure their technological transitions. A CNCF member since 2019, ILKI provides a full range of services on Kubernetes and the CNCF ecosystem to design, build, operate and support your Cloud Native infrastructures.


/cncf-ilki-kcsp-reviewer — Review ILKI (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ILKI is an IT consulting, architecture, and training company. ILKI helps companies to accelerate and secure their technological transitions.  A CNCF member since 2019, ILKI provides a full range of services on Kubernetes and the CNCF ecosystem to design, build, operate and support your Cloud Native infrastructures.


/cncf-ilki-member-reviewer — Review ILKI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-in-toto-reviewer — Review in-toto using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-inclavare-containers-reviewer — Review Inclavare Containers using its official documentation and repository in the Runtime / Container Runtime category. A novel container runtime, aka confidential container, for cloud-native confidential computing and enclave runtime ecosystem


/cncf-incloud-member-reviewer — Review Incloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-indeed-member-reviewer — Review Indeed (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-infinidat-reviewer — Review INFINIDAT using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-infinispan-reviewer — Review Infinispan using its official documentation and repository in the App Definition and Development / Database category.


/cncf-infino-ai-member-reviewer — Review Infino AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-infinyon-reviewer — Review Infinyon using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-infisical-member-reviewer — Review Infisical (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-influxdata-reviewer — Review InfluxData using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-influxdb-reviewer — Review Influxdb using its official documentation and repository in the AI Native Infra / Observability category. InfluxDB is an open source time series database written in Rust, using Apache Arrow, Apache Parquet, and Apache DataFusion as its foundational building blocks.


/cncf-infoblox-member-reviewer — Review Infoblox (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-infoscale-for-kubernetes-reviewer — Review InfoScale for Kubernetes using its official documentation and repository in the Runtime / Cloud Native Storage category. InfoScale for Kubernetes is a Kubernetes-native platform for running stateful workloads with high availability, disaster recovery, and performance optimization across core, edge, and remote office environments. It provides automated, block-level replication and recovery with minimal RPO, enabling business continuity and fast failover during infrastructure or application failures. InfoScale delivers Kubernetes-based virtualization support, including OpenShift Virtualization, enabling unified protection and recovery for both virtual machines and containerized workloads. Through a unified, CSI-compliant architecture that is agnostic to underlying storage vendors, InfoScale integrates seamlessly with heterogeneous storage environments. Designed for hybrid and distributed Kubernetes deployments, InfoScale helps organizations operate resilient, consistent, and scalable platforms.


/cncf-infosys-kcsp-reviewer — Review Infosys (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Infosys architecture-first approach simplifies how you develop and run container-based applications on industry-standard Kubernetes infrastructure. We provide a gamut of Kubernetes services, solutions and platforms for enterprises to accelerate their cloud journey.


/cncf-infosys-member-reviewer — Review Infosys (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-infracost-reviewer — Review Infracost using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. Infracost shows engineers how their code changes will impact cloud costs. It does this by sitting in the CI/CD workflow and on changing IaC code, leaves a comment like “This change will increase your next month bill by 20%” with a detailed breakdown of resources changed and their cost.


/cncf-infros-member-reviewer — Review InfrOS (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-inlets-reviewer — Review Inlets using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-inngest-reviewer — Review Inngest using its official documentation and repository in the AI Agent / Workflow Orchestration category. The leading workflow orchestration platform. Run stateful step functions and AI workflows on serverless, servers, or the edge.


/cncf-innogrid-kcsp-reviewer — Review Innogrid (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Innogrid provides the consulting and technical support related to K8s, especially for an optimized migration for customers&#39; environment.


/cncf-innogrid-member-reviewer — Review Innogrid (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-inspektor-gadget-reviewer — Review Inspektor Gadget using its official documentation and repository in the Observability and Analysis / Observability category. Open source eBPF debugging and data collection tool for Kubernetes and Linux


/cncf-inspur-kcntp-reviewer — Review INSPUR (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category.


/cncf-inspur-kcsp-reviewer — Review Inspur (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As a CNCF Member, a Kubernetes Certified Service Provider (KCSP), and an official CNCF-authorized Kubernetes Training Partner (KTP), Inspur Cloud leverages its profound expertise in cloud-native technologies to launch the CNCF CKAD Cloud-Native Application Development training series. This program is designed to upskill cloud-native application developers, guiding them from foundational to advanced concepts in cloud-native and Kubernetes technologies. It equips participants to confidently tackle the CKAD certification exam and excel in professional challenges within the Kubernetes ecosystem.


/cncf-inspur-member-reviewer — Review Inspur (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-inspur-storage-reviewer — Review Inspur Storage using its official documentation and repository in the Runtime / Cloud Native Storage category. Cloud + AI, operate new massive data operations.


/cncf-instana-reviewer — Review Instana using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-instructor-reviewer — Review Instructor using its official documentation and repository in the AI Agent / Structured Output category. structured outputs for llms


/cncf-instruqt-member-reviewer — Review Instruqt (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-intercept-kcsp-reviewer — Review Intercept (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We help you run Kubernetes on Azure with better control, lower risk, and less overhead. Our team manages governance, security, and cost efficiency, so you can focus on building and deploying.


/cncf-intercept-member-reviewer — Review Intercept (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-interlink-reviewer — Review Interlink using its official documentation and repository in the Runtime / Container Runtime category. InterLink aims to provide an abstraction for the execution of a Kubernetes pod on any remote resource capable of managing a Container execution lifecycle thanks to the Virtual Kubelet interface. It allows you to extend your cloud environment anywhere by running Kubernetes workloads on various infrastructures, creating a seamless cloud-native experience across diverse environments.


/cncf-intersystems-iris-data-platform-reviewer — Review InterSystems IRIS Data Platform using its official documentation and repository in the App Definition and Development / Database category.


/cncf-intuit-member-reviewer — Review Intuit (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-iomesh-reviewer — Review IOMesh using its official documentation and repository in the Runtime / Cloud Native Storage category. IOMesh is an enterprise Kubernetes-native distributed storage product with extreme performance and strong reliability. It leverages Kubernetes-native DevOps system to manage storage resources in the Kubernetes cluster, and provide persistent storage for the most demanding stateful applications. IOMesh will help reduce the cost and complexity of adopting Kubernetes and accelerate your containerization journey.


/cncf-ionir-reviewer — Review Ionir using its official documentation and repository in the Runtime / Cloud Native Storage category. Ionir Liberates Your Applications from Clouds


/cncf-ionos-cloud-managed-kubernetes-reviewer — Review IONOS Cloud Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. IONOS Cloud Managed Kubernetes is the ideal platform for automating the deployment, scaling and management of highly scalable and demanding containerized applications.


/cncf-ionos-member-reviewer — Review IONOS (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-irondb-reviewer — Review IronDB using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-isovalent-member-reviewer — Review Isovalent (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-isovalent-reviewer — Review Isovalent using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-isscloud-reviewer — Review iSSCloud using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. iSoftStone Multi-cloud Management System is an enterprise-level cloud management platform applied in multiple public  clouds and public-private hybrid cloud scenarios. It provides consistent self-service and operation management capabilities  on heterogeneous resources which include centralized resource management and control, automated resource delivery, intelligent  operation analysis and O&amp;M management to reduce the costs of hybrid cloud management.


/cncf-istio-reviewer — Review Istio using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Simplify observability, traffic management, security, and policy with the Istio service mesh.


/cncf-istio-wasm-reviewer — Review Istio (Wasm) using its official documentation and repository in the Wasm / Embedded Functions category. Simplify observability, traffic management, security, and policy with the Istio service mesh.


/cncf-isulad-reviewer — Review iSulad using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-itgix-kcsp-reviewer — Review ITGix (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We specialize in building and supporting Kubernetes clusters across popular cloud providers, as well as on-premise vanilla installations. Our expertise lies in assisting customers with their transition to cloud-native environments and containerization. We can consult, implement and support our clients to adopt Kubernetes and tools from the CNCF landscape.


/cncf-itgix-member-reviewer — Review ITGix (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-itopia-reviewer — Review itopia using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Containerized developer environments in a browser - onboard developers fast and prevent code exfiltration with precise security controls. Devs launch spaces with all their tools pre-installed and start coding in seconds.


/cncf-itq-kcntp-reviewer — Review ITQ (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Master Kubernetes, Cloud Native, and Automation with ITQ — hands-on training built by practitioners to help your teams gain real-world skills they can apply immediately. Tailored or ready-to-go programs designed to accelerate your cloud-native journey.


/cncf-itq-kcsp-reviewer — Review ITQ (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ITQ, a European consulting company, specializes in deploying and managing Kubernetes for organizations of all sizes. Our services are designed to help you navigate the complexities of Kubernetes, allowing you to focus on delivering value.


/cncf-itq-member-reviewer — Review ITQ (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-jaeger-reviewer — Review Jaeger using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-javy-reviewer — Review Javy using its official documentation and repository in the Wasm / Languages category. Scripting languages that support Wasm


/cncf-jenkins-reviewer — Review Jenkins using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-jenkinsx-reviewer — Review JenkinsX using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-jetify-member-reviewer — Review Jetify (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-jetstack-kcsp-reviewer — Review Jetstack (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Jetstack is an organisation focused entirely on Kubernetes. They will help you to get the most out of Kubernetes through expert professional services and open source tooling. Get in touch, and accelerate your project.


/cncf-jetstack-member-reviewer — Review Jetstack (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-jfrog-artifactory-reviewer — Review JFrog Artifactory using its official documentation and repository in the Provisioning / Container Registry category. Shipping updates continuously and automatically has become a critical element of any successful operation. JFrog is revolutionizing the software world with the practice of Continuous Update, with a speed and continuity that forever changes the way organizations manage and release software.


/cncf-jhipster-reviewer — Review JHipster using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-jina-reader-reviewer — Review Jina Reader using its official documentation and repository in the AI Agent / Agent Tool category. Convert any URL to an LLM-friendly input with a simple prefix


/cncf-jozu-member-reviewer — Review Jozu (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-jpmorgan-chase-member-reviewer — Review JPMorgan Chase (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-jreleaser-reviewer — Review JReleaser using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-juicefs-reviewer — Review JuiceFS using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-juju-reviewer — Review Juju using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-jysk-supporter-reviewer — Review JYSK (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-k0rdent-reviewer — Review k0rdent using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. k0rdent is an open-source, Kubernetes-native distributed container management environment for platform engineers.


/cncf-k0s-reviewer — Review k0s using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-k3k-reviewer — Review K3k using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Virtual Kubernetes cluster


/cncf-k3s-reviewer — Review k3s using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Lightweight Kubernetes


/cncf-k6-reviewer — Review k6 using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. k6 is a developer-centric, free and open-source load testing tool built for making performance testing a productive and enjoyable experience. Using k6, you&#39;ll be able to catch performance regression and problems earlier, allowing you to build resilient systems and robust applications.


/cncf-k8gb-reviewer — Review k8gb using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category. A cloud native Kubernetes Global Balancer


/cncf-k8sgpt-reviewer — Review K8sGPT using its official documentation and repository in the Observability and Analysis / Observability category. Giving Kubernetes Superpowers to everyone


/cncf-k8up-reviewer — Review K8up using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-kafka-reviewer — Review Kafka using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-kagent-reviewer — Review kagent using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Kagent is an open source programming framework designed for DevOps and platform engineers to run AI agents in Kubernetes


/cncf-kai-scheduler-reviewer — Review KAI Scheduler using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. KAI Scheduler is a robust, efficient, and scalable Kubernetes scheduler that optimizes GPU resource allocation for AI workloads in large-scale clusters.


/cncf-kairos-reviewer — Review Kairos using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. The immutable Linux meta-distribution for edge Kubernetes


/cncf-kaito-reviewer — Review KAITO using its official documentation and repository in the Inference / Framework category. Kubernetes AI Toolchain Operator (KAITO) simplifies LLM inference, tuning, and RAG workloads on Kubernetes.


/cncf-kaniko-reviewer — Review kaniko using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-kanister-reviewer — Review Kanister using its official documentation and repository in the Runtime / Cloud Native Storage category. An extensible framework for application-level data management on Kubernetes


/cncf-kapa-ai-member-reviewer — Review kapa.ai (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kapeta-reviewer — Review Kapeta using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Kapeta accelerates and automates the entire software development lifecycle


/cncf-kapitan-reviewer — Review Kapitan using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. A configuration management system for platform engineering and other things


/cncf-karmada-reviewer — Review Karmada using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Open, Multi-Cloud, Multi-Cluster Kubernetes Orchestration


/cncf-karpenter-reviewer — Review Karpenter using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. Karpenter is a Kubernetes Node Autoscaler built for flexibility, performance, and simplicity.


/cncf-kasten-reviewer — Review Kasten using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-kata-containers-reviewer — Review Kata Containers using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-katalyst-reviewer — Review Katalyst using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Katalyst is a QoS-based resource management system for workload colocation on kubernetes


/cncf-kbind-reviewer — Review kbind using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. kbind aims to provide better support for service providers and consumers that reside in distinct Kubernetes clusters.


/cncf-kcl-reviewer — Review KCL using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. A constraint-based record &amp; functional language mainly used in configuration and policy scenarios.


/cncf-kcp-reviewer — Review kcp using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-keda-reviewer — Review KEDA using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-kedify-member-reviewer — Review Kedify (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-keep-reviewer — Review Keep using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-kepler-reviewer — Review Kepler using its official documentation and repository in the Observability and Analysis / Observability category. Kepler (Kubernetes-based Efficient Power Level Exporter) uses eBPF to probe energy related system stats and exports as Prometheus metrics.


/cncf-keploy-reviewer — Review Keploy using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-keptn-reviewer — Review Keptn using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Cloud-native application life-cycle orchestration. Keptn automates your SLO-driven multi-stage delivery and operations &amp; remediation of your applications.


/cncf-kestra-reviewer — Review Kestra using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Scalable, event-driven, language-agnostic orchestration and scheduling platform to manage millions of workflows declaratively in code.


/cncf-kestrel-ai-member-reviewer — Review Kestrel AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-keycloak-reviewer — Review Keycloak using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Keycloak is an open-source identity and access management solution for modern applications and services,  built on top of industry security standard protocols.


/cncf-keyfactor-member-reviewer — Review Keyfactor (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-keylime-reviewer — Review Keylime using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Bootstrap &amp; Maintain Trust on the Edge / Cloud and IoT.


/cncf-keymate-io-member-reviewer — Review Keymate.io (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kgateway-reviewer — Review Kgateway using its official documentation and repository in the AI Native Infra / Gateway category. An Envoy-powered, Kubernetes-native API Gateway that integrates Kubernetes Gateway API with a control plane for API connectivity in any cloud environment.


/cncf-kiali-reviewer — Review Kiali using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-kics-reviewer — Review KICS using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Find security vulnerabilities, compliance issues, and infrastructure misconfigurations early in the development cycle of your infrastructure-as-code with KICS by Checkmarx.


/cncf-kilo-reviewer — Review Kilo using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-kind-reviewer — Review kind using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. kind creates local multi-node Kubernetes clusters using Docker container nodes.


/cncf-king-supporter-reviewer — Review King (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-kiosk-reviewer — Review kiosk using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-kiowy-member-reviewer — Review KIOWY (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kiratech-s-p-a-kcntp-reviewer — Review Kiratech S.p.A. (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Kiratech provides cloud native training and professional services for the design, implementation and operation of Kubernetes solutions.


/cncf-kiratech-s-p-a-kcsp-reviewer — Review Kiratech S.p.A. (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Kiratech provides cloud native training and professional services for the design, implementation and operation of Kubernetes solutions.


/cncf-kiratech-s-p-a-member-reviewer — Review Kiratech S.p.A. (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kitops-reviewer — Review KitOps using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. An open standard for packaging, managing, and deploying ML models and artifacts across different systems


/cncf-kloia-kcsp-reviewer — Review kloia (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We specialize in end-to-end Kubernetes consulting and platform engineering — from cluster design, deployment, and migration to observability, security hardening, and cost optimization. Our team helps organizations adopt and scale Kubernetes with best practices, automation, and long-term operational excellence.


/cncf-kloia-member-reviewer — Review kloia (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kloudfuse-member-reviewer — Review Kloudfuse (member) using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-kmesh-reviewer — Review Kmesh using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Kmesh is a high-performance and low overhead service mesh data plane based on eBPF and programmable kernel. Kmesh brings traffic management, security and monitoring to service communication without needing application code changes. It is natively sidecarless, zero intrusion and without adding any resource cost to application container.


/cncf-knative-reviewer — Review Knative using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Knative is a developer-focused serverless application layer which is a great complement to the existing Kubernetes application constructs. Knative consists of three components: an HTTP-triggered autoscaling container runtime called “Knative Serving”, a CloudEvents-over-HTTP asynchronous routing layer called “Knative Eventing”, and a developer-focused function framework which leverages the Serving and Eventing components, called &quot;Knative Functions&quot;.


/cncf-knix-reviewer — Review Knix using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-ko-reviewer — Review ko using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. ko is a tool to simplify building Go container images


/cncf-kodekloud-kcntp-reviewer — Review KodeKloud (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. KodeKloud provides Interactive Hands-On training on Cloud Native Computing technologies using immersive learning techniques to help spread awareness of CNCF and its projects around the world.


/cncf-kodekloud-kcsp-reviewer — Review KodeKloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. KodeKloud provides Interactive Hands-On training on Cloud Native Computing technologies using immersive learning techniques to help spread awareness  of CNCF and its projects around the world.


/cncf-kodekloud-member-reviewer — Review KodeKloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-komodor-member-reviewer — Review Komodor (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kong-member-reviewer — Review Kong (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kong-reviewer — Review Kong using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-konveyor-reviewer — Review Konveyor using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Documentation for Konveyor Community


/cncf-koordinator-reviewer — Review Koordinator using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. QoS based scheduling system for hybrid orchestration workloads on Kubernetes, bringing workloads the best layout and status.


/cncf-kops-reviewer — Review kOps using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. kOps is open-source, community-maintained tooling for production-grade Kubernetes cluster installation, management and upgrades.


/cncf-korea-open-source-association-member-reviewer — Review Korea Open Source Association (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-kosko-reviewer — Review Kosko using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Kosko can help you organize your Kubernetes manifests in TypeScript, manage multiple environments, ensure type safety using OpenAPI schema, and find issues in your manifests.


/cncf-kotlin-reviewer — Review Kotlin using its official documentation and repository in the Wasm / Languages category. Managed language


/cncf-kots-reviewer — Review KOTS using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. KOTS provides the framework, tools and integrations that enable the delivery and management of 3rd-party Kubernetes applications, a.k.a. Kubernetes Off-The-Shelf (KOTS) Software.&#39;


/cncf-koyeb-reviewer — Review Koyeb using its official documentation and repository in the Serverless / Hosted Platform category. Koyeb is a developer-friendly serverless platform to deploy apps globally. Seamlessly run Docker containers, web apps, and APIs with git-based deployment, native autoscaling, a global edge network, and built-in service mesh and discovery.


/cncf-kpt-reviewer — Review kpt using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Automate Kubernetes Configuration Editing


/cncf-kraken-reviewer — Review Kraken using its official documentation and repository in the Provisioning / Container Registry category.


/cncf-krakend-reviewer — Review KrakenD using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-kratix-reviewer — Review Kratix using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-krator-reviewer — Review Krator using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-kratos-reviewer — Review kratos using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. Kratos is a microservice-oriented governance framework implements by golang, which offers convenient capabilities to help you quickly build a bulletproof application from scratch.


/cncf-krkn-reviewer — Review Krkn using its official documentation and repository in the Observability and Analysis / Chaos Engineering category. Chaos testing tool for Kubernetes to identify bottlenecks and improve resilience and performance under failure conditions.


/cncf-kruize-reviewer — Review Kruize using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. Kruize analyzes your Kubernetes workload metrics and automatically generates right-sizing recommendations for CPU, memory, and GPU resources — reducing costs and improving performance without manual tuning.


/cncf-krustlet-reviewer — Review Krustlet using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-ksat-supporter-reviewer — Review KSAT (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-kserve-reviewer — Review KServe using its official documentation and repository in the Inference / Framework category. Standardized Distributed Generative and Predictive AI Inference Platform for Scalable, Multi-Framework Deployment on Kubernetes


/cncf-kt-cloud-member-reviewer — Review kt cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kuadrant-reviewer — Review Kuadrant using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. Kuadrant combines Gateway API and Istio-based gateway controllers to enhance application connectivity. It enables platform engineers  and application developers to easily connect, secure, and protect their services and infrastructure across multiple clusters  with policies for TLS, DNS, application authentication &amp; authorization, and rate limiting.


/cncf-kuasar-reviewer — Review Kuasar using its official documentation and repository in the Runtime / Container Runtime category. A multi-sandbox container runtime that provides cloud-native, all-scenario multiple sandbox container solutions.


/cncf-kuasar-wasm-reviewer — Review Kuasar (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-kubara-reviewer — Review kubara using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. A single binary CLI tool written in Go to bootstrap, build, and package Kubernetes platforms.


/cncf-kube-bench-reviewer — Review kube-bench using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-kube-burner-reviewer — Review Kube-burner using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Kubernetes performance and scale test orchestration framework written in golang


/cncf-kube-green-reviewer — Review kube-green using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. kube-green is a simple k8s addon that automatically shuts down (some of) your resources when you don&#39;t need them.


/cncf-kube-hunter-reviewer — Review kube-hunter using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-kube-ovn-reviewer — Review Kube-OVN using its official documentation and repository in the Runtime / Cloud Native Network category. A Kubernetes Network Fabric for Enterprises that is Rich in Functions and Easy in Operations


/cncf-kube-router-reviewer — Review Kube-router using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-kube-rs-reviewer — Review kube-rs using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. kube-rs is the core Rust ecosystem for building applications against Kubernetes


/cncf-kube-scheduler-reviewer — Review Kube-scheduler using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-kube-vip-reviewer — Review kube-vip using its official documentation and repository in the Runtime / Cloud Native Network category. Kubernetes Virtual IP and Load-Balancer for both control plane and Kubernetes services


/cncf-kubeadmiral-reviewer — Review KubeAdmiral using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. KubeAdmiral is a multi-cluster scheduling and orchestration system for Kubernetes.


/cncf-kubean-reviewer — Review Kubean using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Product ready cluster lifecycle management toolchains based on kubespray and other cluster LCM engine.


/cncf-kubearmor-reviewer — Review KubeArmor using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Runtime protection for Kubernetes &amp; other cloud Workloads. Kubearmor provides a observability and policy enforcement system to restrict any unwanted, malicious behaviour of cloud-native workloads at runtime.


/cncf-kubeasz-reviewer — Review Kubeasz using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Kubeasz is a tool to deploy a Production Ready Kubernetes Cluster with ansible playbooks.


/cncf-kubeblocks-by-apecloud-reviewer — Review KubeBlocks by ApeCloud using its official documentation and repository in the App Definition and Development / Database category. KubeBlocks is an open-source Kubernetes operator that manages relational, NoSQL, vector, and streaming databases on the public cloud or on-premise.


/cncf-kubebrain-reviewer — Review KubeBrain using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category.


/cncf-kubeclipper-reviewer — Review KubeClipper using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Manage kubernetes in the most light and convenient way.


/cncf-kubecost-reviewer — Review Kubecost using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-kubectl-mcp-server-reviewer — Review kubectl-mcp-server using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. AI-native Kubernetes management through Model Context Protocol (MCP) with 220+ tools for natural language cluster operations.


/cncf-kubedb-by-appscode-reviewer — Review KubeDB by AppsCode using its official documentation and repository in the App Definition and Development / Database category.


/cncf-kubedb-member-reviewer — Review KubeDB (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kubediagrams-reviewer — Review KubeDiagrams using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. KubeDiagrams is a tool to generate Kubernetes architecture diagrams from Kubernetes manifest files, kustomization files, Helm charts, and actual cluster state. KubeDiagrams supports most of all Kubernetes built-in resources, any custom resources, and label-based resource clustering.


/cncf-kubedl-reviewer — Review KubeDL using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-kubeedge-reviewer — Review KubeEdge using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-kubeedge-wasm-reviewer — Review KubeEdge (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-kubeelasti-reviewer — Review KubeElasti using its official documentation and repository in the Serverless / Installable Platform category. Auto scale-to-zero pods when idle and scale up pods when traffic arrives, without losing any requests. KubeElasti uses a smart proxy that queues incoming requests while scaling up targets, ensuring no request loss. It works with existing Kubernetes services and deployments without requiring code changes.


/cncf-kubefirst-reviewer — Review Kubefirst using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Instant operational open source gitops platforms for platform engineering, infrastructure, and software teams


/cncf-kubefleet-reviewer — Review KubeFleet using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. A multi-cluster solution that enables users to effectively manage their applications running in multiple Kubernetes clusters.


/cncf-kubeflow-katib-reviewer — Review Kubeflow Katib using its official documentation and repository in the Training / Post Training category. Automated Machine Learning on Kubernetes


/cncf-kubeflow-model-registry-reviewer — Review Kubeflow Model Registry using its official documentation and repository in the AI Native Infra / Model Asset and Registry category. Model Registry to store and manage models, versions, and artifacts metadata.


/cncf-kubeflow-mpi-operator-reviewer — Review Kubeflow MPI Operator using its official documentation and repository in the Training / Distributed Training category. MPI Operator to manage all-reduce distributed training and HPC workloads.


/cncf-kubeflow-notebooks-reviewer — Review Kubeflow Notebooks using its official documentation and repository in the Data / Data Science category. Machine Learning Toolkit for Kubernetes


/cncf-kubeflow-pipelines-reviewer — Review Kubeflow Pipelines using its official documentation and repository in the AI Native Infra / Continuous Integration and Delivery category. Machine Learning Pipelines for Kubeflow


/cncf-kubeflow-reviewer — Review Kubeflow using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Kubeflow is the foundation of tools for AI Platforms on Kubernetes.


/cncf-kubeflow-spark-operator-reviewer — Review Kubeflow Spark Operator using its official documentation and repository in the Data / Data Architecture category. Operator to manage Apache Spark applications on Kubernetes.


/cncf-kubeflow-training-operator-reviewer — Review Kubeflow Training Operator using its official documentation and repository in the Training / Distributed Training category. Training operators on Kubernetes


/cncf-kubefwd-reviewer — Review kubefwd using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Bulk port forwarding Kubernetes services to localhost with unique loopback IPs, enabling multiple services on the same port without conflicts.


/cncf-kubegateway-reviewer — Review KubeGateway using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-kubeinvaders-reviewer — Review Kubeinvaders using its official documentation and repository in the Observability and Analysis / Chaos Engineering category. Gamified Chaos Engineering Tool for K8s


/cncf-kubekey-reviewer — Review Kubekey using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Kubekey provides a flexible, rapid and convenient way to install Kubernetes only, both Kubernetes and KubeSphere, and related cloud-native add-ons.  It is also an efficient tool to scale and upgrade your cluster.


/cncf-kubelinter-reviewer — Review KubeLinter using its official documentation and repository in the Provisioning / Security &amp; Compliance category. KubeLinter analyzes Kubernetes YAML files and Helm charts, and checks them against a variety of best practices, with a focus on production readiness and security.


/cncf-kubemq-reviewer — Review KubeMQ using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. KubeMQ is a Kubernetes Message Queue Broker


/cncf-kubeops-kcntp-reviewer — Review KubeOps (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Training from foundation to expert – with our practice-oriented training programs, you will gain all the Kubernetes knowledge you need for your daily operation.


/cncf-kubeops-kcsp-reviewer — Review KubeOps (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. KubeOps provides assistance with all Kubernetes-related issues, like setting up brand-new systems, managed services for existing Kubernetes infrastructure, our own software solutions based on Open Source technology, as well as consulting,  training and certification.


/cncf-kubeops-member-reviewer — Review KubeOps (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kubeops-reviewer — Review KubeOps using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. KubeOps Platform involves all the activities required to run, manage and maintain Kubernetes clusters in production environments, including our best practices, self-deployed tools, and strategies.


/cncf-kubeorbit-reviewer — Review KubeOrbit using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-kuberay-reviewer — Review Kuberay using its official documentation and repository in the AI Native Infra / Orchestration and Scheduling category. A toolkit to run Ray applications on Kubernetes


/cncf-kubereport-reviewer — Review KubeReport using its official documentation and repository in the Observability and Analysis / Observability category. KubeReport is an open-source tool that generates detailed Kubernetes cluster reports in PDF and CSV formats, providing insights into  resource utilization, workload status, and cluster health for easier auditing and troubleshooting.


/cncf-kuberhealthy-reviewer — Review Kuberhealthy using its official documentation and repository in the Observability and Analysis / Observability category. A Kubernetes operator for running synthetic checks as pods. Works great with Prometheus!


/cncf-kubermatic-kcntp-reviewer — Review Kubermatic (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Kubermatic provides enterprise-grade software solutions and professional services to help organizations worldwide to  fully automate their Kubernetes and cloud native operations across multi-cloud, edge, and on-prem.


/cncf-kubermatic-kcsp-reviewer — Review Kubermatic (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Kubermatic empowers organizations to fully automate their Kubernetes operations across multi-cloud, edge, and on-prem environments, simplifying the management of thousands of clusters for leading enterprises. Leveraging our deep expertise as a top 15 contributor to the Kubernetes project, we provide professional services and enterprise-grade software to safely accelerate your cloud-native transformation like we&#39;ve done for leading companies like Lufthansa, Bosch, and T-Systems.


/cncf-kubermatic-kubeone-reviewer — Review Kubermatic KubeOne using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Lifecycle management tool for Highly-Available Kubernetes clusters on any infrastructure


/cncf-kubermatic-kubernetes-platform-reviewer — Review Kubermatic Kubernetes Platform using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Kubermatic Kubernetes Platform automates Kubernetes deployments and Day 2 operations for thousands of Kubernetes clusters on any cloud, on-premises and edge.


/cncf-kubermatic-member-reviewer — Review Kubermatic (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kubernetes-reviewer — Review Kubernetes using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Kubernetes is an open-source system for automating deployment, scaling, and management of containerized applications


/cncf-kubernetes-the-easier-way-reviewer — Review Kubernetes - The Easier Way using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Kubernetes - The Easier Way enables users to create and manage an easily customizable HA Kubernetes cluster with only a couple of commands


/cncf-kubero-reviewer — Review Kubero using its official documentation and repository in the Platform / PaaS/Container Service category. Kubero is a developer friendly selfservice platform for Kubernetes.


/cncf-kubescape-reviewer — Review Kubescape using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Kubescape is an open source security and compliance platform that scans clusters, Kubernetes manifest files (YAML files, and Helm charts), code repositories, container registries and images. It detects misconfigurations according to frameworks such as the NSA-CISA,  MITRE ATT&amp;CK® and CIS, as well as software vulnerabilities, and calculates risk scores.


/cncf-kubeservice-stack-reviewer — Review KubeService Stack using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. KubeService Stack = Kubernetes + Custom Service Build a Kubernetes/KubeEdge enterprise-level peripheral ecosystem


/cncf-kubeshop-member-reviewer — Review kubeshop (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kubeskoop-reviewer — Review KubeSkoop using its official documentation and repository in the Observability and Analysis / Observability category. KubeSkoop is a network monitoring &amp; diagnosis suite for Kubernetes.


/cncf-kubeslice-reviewer — Review KubeSlice using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Applications provided by Avesha Systems, ready to launch on Kubernetes using Kubernetes Helm


/cncf-kubesphere-member-reviewer — Review Kubesphere (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kubesphere-reviewer — Review Kubesphere using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Kubesphere.io is an upstream project of the KubeSphere container management platform. Our vision is to provide an easier, more friendly and more powerful distributed management platform for individuals and enterprises based on Kubernetes, as well as meet more business demands and help more users to use Kubernetes faster and better.


/cncf-kubespray-reviewer — Review Kubespray using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Deploy a Production Ready Kubernetes Cluster


/cncf-kubestellar-reviewer — Review KubeStellar using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. KubeStellar - a flexible solution for challenges associated with multi-cluster configuration management for edge, multi-cloud, and hybrid cloud


/cncf-kubetail-reviewer — Review Kubetail using its official documentation and repository in the Observability and Analysis / Observability category. Real-time logging dashboard for Kubernetes


/cncf-kubevela-reviewer — Review KubeVela using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-kubevirt-reviewer — Review KubeVirt using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Kubernetes Virtualization API and runtime in order to define and manage virtual machines


/cncf-kubevpn-reviewer — Review KubeVPN using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. KubeVPN offers a Cloud-Native Dev Environment that seamlessly connects to your Kubernetes cluster network. Gain access to the Kubernetes cluster network effortlessly using service names or Pod IP/Service IP. Facilitate the interception of inbound traffic from remote Kubernetes cluster services to your local PC through a service mesh and more. For instance, you have the flexibility to run your Kubernetes pod within a local Docker container, ensuring an identical environment, volume, and network setup. With KubeVPN, empower yourself to develop applications entirely on your local PC!


/cncf-kubewarden-reviewer — Review Kubewarden using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Kubewarden is a Policy Engine powered by WebAssembly policies. Its policies can be written in CEL, Rego (OPA &amp; Gatekeeper flavours), Rust, Go, YAML, and others. Kubewarden simplifies Policy-As-Code by allowing policy authors and consumers to use their preferred tooling and stack, develop and test policies out of cluster.


/cncf-kubex-member-reviewer — Review Kubex (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kubex-reviewer — Review Kubex using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-kubezoo-reviewer — Review KubeZoo using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-kublr-kcsp-reviewer — Review Kublr (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Kublr offers enterprise-grade secure, scalable, highly reliable Kubernetes clusters on AWS, Azure, GCP, and on-premise. It includes out-of-the-box backup and disaster recovery, multi-cluster centralized logging and monitoring, and built-in alerting.


/cncf-kublr-member-reviewer — Review Kublr (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kublr-reviewer — Review Kublr using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Accelerate and control the deployment, scaling, monitoring and management of your containerized applications.


/cncf-kudo-reviewer — Review KUDO using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Kubernetes Universal Declarative Operator


/cncf-kueue-reviewer — Review Kueue using its official documentation and repository in the AI Native Infra / Orchestration and Scheduling category. Cloud-native job queueing system for batch, HPC, AI/ML, and similar applications in a Kubernetes cluster.


/cncf-kui-reviewer — Review Kui using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-kuma-reviewer — Review Kuma using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. The universal Envoy service mesh for distributed service connectivity


/cncf-kunobi-reviewer — Review Kunobi using its official documentation and repository in the Observability and Analysis / Observability category. Kunobi is a desktop Kubernetes management app built with Tauri/Rust, featuring multi-cluster management, real-time resource monitoring, and a built-in MCP server for AI-assisted operations.


/cncf-kured-reviewer — Review Kured using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Kured (KUbernetes REboot Daemon) is a Kubernetes daemonset that performs safe automatic node reboots when the need to do so is indicated by the package management system of the underlying OS


/cncf-kurl-reviewer — Review kURL using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Kubernetes URL (kURL) is a framework for creating custom Kubernetes distributions. These distros can then be shared as URLs (to install via curl and bash) or as downloadable packages (to install in airgapped environments).


/cncf-kusari-member-reviewer — Review Kusari (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-kusionstack-reviewer — Review KusionStack using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Declarative Intent Driven Platform Orchestrator for Internal Developer Platform (IDP)


/cncf-kusk-gateway-reviewer — Review Kusk Gateway using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. Kusk Gateway is a self-service API gateway powered by OpenAPI and Envoy. Kusk Gateway is built and maintained by Kubeshop.


/cncf-kwasm-reviewer — Review Kwasm using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-kwok-reviewer — Review kwok using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Kubernetes WithOut Kubelet - Simulates thousands of Nodes and Clusters.


/cncf-kyma-reviewer — Review Kyma using its official documentation and repository in the Platform / PaaS/Container Service category. Kyma is the opinionated set of Kubernetes based modular building blocks that includes the necessary capabilities to develop and run enterprise-grade cloud-native applications.


/cncf-kyverno-json-reviewer — Review Kyverno-JSON using its official documentation and repository in the AI Native Infra / Governance, Policy and Security category. Use Kyverno&#39;s powerful, declarative, low-code policies to validate any runtime or configuration data that can be converted to JSON.


/cncf-kyverno-reviewer — Review Kyverno using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-la-mobiliere-supporter-reviewer — Review La Mobiliere (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-lablup-member-reviewer — Review Lablup (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-lagoon-reviewer — Review Lagoon using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Build and Deploy System for OpenShift &amp; Kubernetes


/cncf-lancedb-reviewer — Review LanceDB using its official documentation and repository in the AI Agent / Vector Database category. Developer-friendly OSS embedded retrieval library for multimodal AI. Search More; Manage Less.


/cncf-langchain-reviewer — Review LangChain using its official documentation and repository in the AI Agent / Agent Framework category. The agent engineering platform


/cncf-langfuse-reviewer — Review Langfuse using its official documentation and repository in the AI Native Infra / Observability category. Open source LLM engineering platform - Observability, metrics, evals, prompt management, playground, datasets.


/cncf-langgraph-reviewer — Review LangGraph using its official documentation and repository in the AI Agent / Agent Framework category. Build resilient language agents as graphs.


/cncf-last9-member-reviewer — Review Last9 (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-last9-reviewer — Review Last9 using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-launchdarkly-reviewer — Review LaunchDarkly using its official documentation and repository in the Observability and Analysis / Feature Flagging category.


/cncf-layotto-reviewer — Review Layotto using its official documentation and repository in the Serverless / Framework category.


/cncf-leaderworkerset-reviewer — Review LeaderWorkerSet using its official documentation and repository in the AI Native Infra / Orchestration and Scheduling category. An API for deploying a group of pods as a unit of replication.


/cncf-leaseweb-managed-kubernetes-reviewer — Review Leaseweb Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Leaseweb Managed Kubernetes takes care of the Installation, Monitoring, and Maintenance of the Kubernetes control plane. Leaseweb also manages the underlying equipment and network infrastructure, making Kubernetes easy and efficient.


/cncf-leaseweb-member-reviewer — Review Leaseweb (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-legit-security-member-reviewer — Review Legit Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-legora-member-reviewer — Review Legora (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-leminnov-kcsp-reviewer — Review Leminnov (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Leminnov helps companies design, deploy, and manage resilient Kubernetes infrastructures tailored to their business needs. From cloud native transformation to automated CI/CD pipelines, we enable scalable, secure, and efficient operations.


/cncf-leminnov-member-reviewer — Review Leminnov (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-lenovo-kcsp-reviewer — Review Lenovo (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Based on the experience of enterprise-level IT governance in the dual state of stability and flexibility, Lenovo Cloud Native products provide one-stop Kubernetes services to help customers achieve a balance between development efficiency and operation and maintenance control.


/cncf-lenovo-member-reviewer — Review Lenovo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-libsql-reviewer — Review libSQL using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-ligato-reviewer — Review Ligato using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-lightbend-reviewer — Review Lightbend using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-lightrun-member-reviewer — Review Lightrun (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-lightrun-reviewer — Review Lightrun using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-lightstep-reviewer — Review LightStep using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-like-minds-consulting-kcsp-reviewer — Review Like Minds Consulting (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Transform your infrastructure with our expert Kubernetes services, ensuring seamless deployment, management, and scaling of your containerized applications for maximum efficiency and reliability.


/cncf-like-minds-consulting-member-reviewer — Review Like Minds Consulting (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-lima-reviewer — Review Lima using its official documentation and repository in the Runtime / Container Runtime category. Linux virtual machines, typically on macOS, for running containerd


/cncf-linbit-member-reviewer — Review Linbit (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-lindb-reviewer — Review LinDB using its official documentation and repository in the Observability and Analysis / Observability category. LinDB is a scalable, high performance, high availability distributed time series database.


/cncf-linkedin-supporter-reviewer — Review LinkedIn (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-linkerd-reviewer — Review Linkerd using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Ultra light, ultra simple, ultra powerful. Linkerd adds security, observability, and reliability to Kubernetes, without the complexity.


/cncf-linode-kubernetes-engine-reviewer — Review Linode Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Linode Kubernetes Engine is a fast and simple way to deploy, scale, and manage Kubernetes clusters running on Linode&#39;s global infrastructure.


/cncf-linstor-reviewer — Review LINSTOR using its official documentation and repository in the Runtime / Cloud Native Storage category. High Performance Software-Defined Block Storage for container, cloud and virtualisation. Fully integrated with Docker, Kubernetes, Openstack, Proxmox etc.


/cncf-linux-foundation-education-kcntp-reviewer — Review Linux Foundation Education (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Linux Foundation Education offers cloud native training and certification programs created in partnership with CNCF and industry-leading experts. As the host of certification exams including KCNA, CKA, CKAD and CKS, Linux Foundation Education is helping professionals demonstrate their skills with cloud native technologies. Training courses are designed to both prepare an individual to deploy, manage and use cloud native technologies in the workplace, as well as to be successful when sitting for a certification.


/cncf-linuxkit-reviewer — Review LinuxKit using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-liquibase-reviewer — Review Liquibase using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Liquibase helps release software faster by bringing DevOps to the database


/cncf-liquid-reply-kcntp-reviewer — Review Liquid Reply (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. As Enterprise Kubernetes Consultants, we design and build your Kubernetes cluster in the public cloud, on-prem, or multi/hybrid cloud environments according to your needs.


/cncf-liquid-reply-kcsp-reviewer — Review Liquid Reply (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As Enterprise Kubernetes Consultants, we design and build your Kubernetes cluster in the public cloud, on-prem, or multi/hybrid cloud environments according to your needs.


/cncf-liquid-reply-member-reviewer — Review Liquid Reply (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-litmus-reviewer — Review Litmus using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-llama-2-reviewer — Review Llama 2 using its official documentation and repository in the Wasm / AI/Machine Learning category.


/cncf-llama-cpp-reviewer — Review llama.cpp using its official documentation and repository in the Inference / Runtime category. LLM inference in C/C++


/cncf-llamafactory-reviewer — Review LlamaFactory using its official documentation and repository in the Training / Post Training category. Unified Efficient Fine-Tuning of 100+ LLMs &amp; VLMs


/cncf-llamaindex-reviewer — Review LlamaIndex using its official documentation and repository in the AI Agent / RAG category. LlamaIndex is the leading document agent and OCR platform


/cncf-llm-d-reviewer — Review llm-d using its official documentation and repository in the Inference / Framework category. llm-d is a Kubernetes-native, high-performance distributed LLM inference framework built on vLLM and the Kubernetes Gateway API Inference Extension, providing intelligent inference scheduling, prefix-cache-aware routing, prefill/decode disaggregation, hierarchical KV offloading, and traffic- and hardware-aware autoscaling across NVIDIA, AMD, Intel, and Google TPU accelerators.


/cncf-llm-guard-reviewer — Review LLM Guard using its official documentation and repository in the AI Agent / Guardrail category. The Security Toolkit for LLM Interactions


/cncf-llmaz-reviewer — Review llmaz using its official documentation and repository in the Inference / Framework category. Easy, advanced inference platform for large language models on Kubernetes.


/cncf-llvm-reviewer — Review llvm using its official documentation and repository in the Wasm / Tooling category.


/cncf-lm-evaluation-harness-reviewer — Review lm-evaluation-harness using its official documentation and repository in the Training / Evaluation category. A framework for few-shot evaluation of language models


/cncf-lmql-reviewer — Review LMQL using its official documentation and repository in the AI Agent / Structured Output category. A language for constraint-guided and efficient LLM programming


/cncf-lockheed-martin-member-reviewer — Review Lockheed Martin (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-loggie-reviewer — Review Loggie using its official documentation and repository in the Observability and Analysis / Observability category. Loggie is a lightweight, high-performance, cloud-native log collection agent and aggregator based on Golang.


/cncf-logging-operator-kube-logging-reviewer — Review Logging Operator (Kube Logging) using its official documentation and repository in the Observability and Analysis / Observability category. Logging operator for Kubernetes


/cncf-loggly-reviewer — Review Loggly using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-logicmonitor-reviewer — Review LogicMonitor using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-logiq-reviewer — Review Logiq using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-logstash-reviewer — Review Logstash using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-logz-io-reviewer — Review Logz.io using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-longhorn-reviewer — Review Longhorn using its official documentation and repository in the Runtime / Cloud Native Storage category. Cloud-native distributed storage for Kubernetes


/cncf-lovable-member-reviewer — Review Lovable (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-loxilb-reviewer — Review LoxiLB using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. eBPF based cloud-native load-balancer. Powering Kubernetes|Edge|5G|IoT|XaaS Apps.


/cncf-lpi-japan-member-reviewer — Review LPI-Japan (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-lseg-member-reviewer — Review LSEG (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-lumigo-reviewer — Review Lumigo using its official documentation and repository in the Serverless / Tools category.


/cncf-lunar-dev-reviewer — Review Lunar.dev using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. Lunar.dev’s mission is to enable optimization and control of  third-party API consumption in production environments. Lunar is a lightweight tool that empowers DevOps and engineering teams  to centralize consumption, gain insight and visibility into usage patterns and costs, and utilize out-of-the-box policies.


/cncf-lunatic-reviewer — Review Lunatic using its official documentation and repository in the Wasm / Runtimes category.


/cncf-lxd-reviewer — Review lxd using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-ly-corporation-member-reviewer — Review LY Corporation (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-m3-reviewer — Review M3 using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-maas-reviewer — Review MAAS using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Self-service, remote installation of Windows, CentOS, ESXi and Ubuntu on real servers turns your data center into a bare-metal cloud.


/cncf-mackerel-reviewer — Review Mackerel using its official documentation and repository in the Observability and Analysis / Observability category. The server monitoring platform we always wanted.


/cncf-macstadium-member-reviewer — Review MacStadium (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-manageiq-reviewer — Review ManageIQ using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-mandao-fintech-kcsp-reviewer — Review Mandao Fintech (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Mandao Technology has launched K8s operations and maintenance services around the full life cycle of Kubernetes to help you enjoy the benefits of rapid application transformation and efficiency improvement brought by K8s, greatly reducing business risks.


/cncf-mandao-fintech-member-reviewer — Review Mandao Fintech (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mantech-solution-kcsp-reviewer — Review Mantech Solution (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. With over 30 years of experience as an IT solutions provider, we provide Kubernetes technical support for an optimized and stable migration for each enterprise environment. Our engineers are Kubernetes experts and we provide end-to-end management of your Kubernetes operations.


/cncf-mantech-solution-member-reviewer — Review Mantech Solution (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mariadb-reviewer — Review MariaDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-marimo-reviewer — Review marimo using its official documentation and repository in the Data / Data Science category. A reactive notebook for Python — run reproducible experiments, execute as a script, deploy as an app, and version with git.


/cncf-marqo-reviewer — Review marqo using its official documentation and repository in the AI Agent / Vector Database category. Ecommerce Search and Discovery


/cncf-massdriver-member-reviewer — Review Massdriver (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-masterclass-supporter-reviewer — Review MasterClass (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-mastercontrol-supporter-reviewer — Review MasterControl (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-matano-reviewer — Review Matano using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Open source security lake platform for AWS


/cncf-mathworks-member-reviewer — Review Mathworks (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-matrixx-digital-commerce-reviewer — Review MATRIXX Digital Commerce using its official documentation and repository in the Special / Certified CNFs category. The MATRIXX Digital Commerce Platform (DCP) is a cloud-native, real-time monetization engine delivering industry-compliant rating and charging functionality along with a rich array of digital commerce capabilities such as subscription management, event streaming and management, personalization, and digital payments.


/cncf-maven-reviewer — Review Maven using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Apache Maven is a software project management and comprehension tool. Based on the concept of a project object model (POM), Maven can manage a project&#39;s build, reporting and documentation from a central piece of information.


/cncf-maze-member-reviewer — Review Maze (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mckinsey-and-company-member-reviewer — Review McKinsey &amp; Company (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mcp-gateway-registry-reviewer — Review MCP Gateway Registry using its official documentation and repository in the AI Native Infra / Gateway category. Enterprise-ready MCP Gateway &amp; Registry that centralizes AI development tools with secure OAuth authentication, dynamic tool discovery, and unified access for both autonomous AI agents and AI coding assistants. Transform scattered MCP server chaos into governed, auditable tool access with Keycloak/Entra integration.


/cncf-megatron-lm-reviewer — Review Megatron-LM using its official documentation and repository in the Training / Distributed Training category. GPU optimized techniques for training transformer models at-scale


/cncf-meltwater-supporter-reviewer — Review Meltwater (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-mem0-reviewer — Review Mem0 using its official documentation and repository in the AI Agent / State and Memory category. Universal memory layer for AI Agents


/cncf-membrane-reviewer — Review Membrane using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-memcached-reviewer — Review Memcached using its official documentation and repository in the Data / Data Architecture category. A high performance multithreaded event-based key/value cache store intended to be used in a distributed system.


/cncf-memphis-reviewer — Review Memphis using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-merbridge-reviewer — Review Merbridge using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Use eBPF to speed up your Service Mesh like crossing an Einstein-Rosen Bridge.


/cncf-mercedes-benz-ag-member-reviewer — Review Mercedes Benz AG (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-merck-supporter-reviewer — Review Merck (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-mergify-reviewer — Review Mergify using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Optimize your CI/CD pipeline using automation and merge queues


/cncf-mermin-reviewer — Review Mermin using its official documentation and repository in the Observability and Analysis / Observability category. Mermin is a Kubernetes-native network observability tool that uses eBPF exports OpenTelemetry Flow Traces, to provide deep visibility into your cluster’s network communications with zero application changes required.


/cncf-mesh-infra-reviewer — Review Mesh Infra using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. VS Code and JetBrains IDE extension that visualizes Terraform, Kubernetes, Docker Compose, and ArgoCD infrastructure as interactive diagrams directly inside the editor. Scan IaC files and generate live topology maps without leaving your IDE.


/cncf-meshcloud-member-reviewer — Review meshcloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-meshery-reviewer — Review Meshery using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-meshery-wasm-reviewer — Review Meshery (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-metal-stack-cloud-kubernetes-reviewer — Review Metal Stack Cloud Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Metal Stack Cloud Kubernetes is a fully managed and scalable Kubernetes service for the deployment and management of Kubernetes clusters and containerized applications.


/cncf-metal3-io-reviewer — Review metal3-io using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-metalbear-member-reviewer — Review MetalBear (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-metallb-reviewer — Review MetalLB using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. A network load-balancer implementation for Kubernetes using standard routing protocols


/cncf-metarget-reviewer — Review Metarget using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Metarget is a framework providing automatic constructions of vulnerable cloud native infrastructures.


/cncf-metatype-reviewer — Review Metatype using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-mewz-reviewer — Review Mewz using its official documentation and repository in the Wasm / Runtimes category.


/cncf-mezmo-member-reviewer — Review Mezmo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mezmo-reviewer — Review Mezmo using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-mia-platform-reviewer — Review Mia-Platform using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Mia-Platform is a Digital Platform Builder: it enables organizations to build and orchestrate their Internal Developer Platforms (IDP) to scale cloud-native development and operations. Mia-Platform also features a Data Fabric solution that can be used to build a Digital Integration Hub. This solution ingests data from different sources, aggregates it in single views, and makes it available in near real-time.


/cncf-michelin-member-reviewer — Review Michelin (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-microcks-reviewer — Review Microcks using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-microk8s-reviewer — Review MicroK8s using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. An easy to install, single node local Kubernetes.


/cncf-micrometer-reviewer — Review Micrometer using its official documentation and repository in the Observability and Analysis / Observability category. As an instrumentation facade, Micrometer allows you to instrument your code with dimensional metrics with a vendor-neutral interface and decide on the observability system as a last step. Instrumenting your core library code with Micrometer allows the libraries to be included in applications that ship data to different backends.


/cncf-microsoft-kcsp-reviewer — Review Microsoft (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Microsoft Services provides a comprehensive set of offerings to deploy and secure Kubernetes on Azure as well as containerize applications and migrate them from on-premises or other Clouds to Azure.


/cncf-microsoft-member-reviewer — Review Microsoft (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-microsoft-sql-server-reviewer — Review Microsoft SQL Server using its official documentation and repository in the App Definition and Development / Database category. Data-Driven. Faster Insights. Breakthrough Performance. In-Memory Technology. Hybrid Data Platform.


/cncf-middleware-member-reviewer — Review Middleware (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-middleware-reviewer — Review Middleware using its official documentation and repository in the Observability and Analysis / Observability category. Middleware is a full-stack cloud observability platform that consolidates telemetry data on a unified timeline, helping developers streamline issue resolution and enhance operational efficiency &amp; user experience.


/cncf-midships-member-reviewer — Review Midships (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-midway-serverless-reviewer — Review Midway Serverless using its official documentation and repository in the Serverless / Framework category.


/cncf-milvus-reviewer — Review Milvus using its official documentation and repository in the AI Agent / Vector Database category. A cloud-native vector database, storage for next generation AI applications.


/cncf-minikube-reviewer — Review minikube using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. minikube runs a local Kubernetes cluster on macOS, Linux, and Windows.


/cncf-minimus-member-reviewer — Review Minimus (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-minio-reviewer — Review MinIO using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-mirantis-kcntp-reviewer — Review Mirantis (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. We offer intensive cloud native training bootcamp(s) for IT professionals looking to develop skills in deploying and administering containerized applications in Kubernetes.


/cncf-mirantis-kcsp-reviewer — Review Mirantis (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Mirantis provides full-stack support for clouds built with Kubernetes and related open source software, using a flexible GitOps model for lifecycle management.


/cncf-mirantis-kubernetes-engine-reviewer — Review Mirantis Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Mirantis Kubernetes Engine Kubernetes distribution.


/cncf-mirantis-member-reviewer — Review Mirantis (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-miraxia-edge-technology-member-reviewer — Review Miraxia Edge Technology (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mirrord-reviewer — Review mirrord using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. mirrord is an open-source developer tool that runs local code as if it were a pod in a remote Kubernetes cluster, with real env vars, DNS, network, and inbound traffic. Lets engineers iterate on a single service locally while interacting with the rest of the cluster, and is also used by AI coding agents (Claude, Cursor, Codex).


/cncf-mist-io-reviewer — Review Mist.io using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-mistral-ai-member-reviewer — Review Mistral AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mlflow-reviewer — Review Mlflow using its official documentation and repository in the AI Native Infra / Continuous Integration and Delivery category. The open source AI engineering platform for agents, LLMs, and ML models. MLflow enables teams of all sizes to debug, evaluate, monitor, and optimize production-quality AI applications while controlling costs and managing access to models and data.


/cncf-mlrun-reviewer — Review MLRun using its official documentation and repository in the AI Native Infra / Continuous Integration and Delivery category. Machine Learning automation and tracking


/cncf-mockserver-reviewer — Review MockServer using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Mocking, debugging proxy and chaos engineering for HTTP, gRPC, GraphQL, LLM, MCP, Kafka, TCP and more.


/cncf-modelcontextprotocol-reviewer — Review Modelcontextprotocol using its official documentation and repository in the AI Agent / Protocol category. Specification and documentation for the Model Context Protocol.


/cncf-modelpack-reviewer — Review ModelPack using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. The project establishes open standards for packaging, distributing and running AI artifacts in the cloud-native environment.


/cncf-modsurfer-reviewer — Review Modsurfer using its official documentation and repository in the Wasm / Debugging &amp; Observability category.


/cncf-mogdb-reviewer — Review MogDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-mogenius-kubernetes-platform-reviewer — Review mogenius Kubernetes Platform using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. mogenius combines comprehensive Kubernetes management with developer‑focused platform services.


/cncf-mogenius-member-reviewer — Review mogenius (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mondoo-reviewer — Review Mondoo using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Cloud native infrastructure security for your entire fleet


/cncf-mongodb-reviewer — Review MongoDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-monocle-reviewer — Review Monocle using its official documentation and repository in the Observability and Analysis / Observability category. The goal of the open source Monocle project is to help GenAI developers trace their applications. Hosted in incubation as a Sandbox project in LF AI &amp; Data.


/cncf-monokle-reviewer — Review Monokle using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Monokle helps you achieve high-quality Kubernetes deployments throughout the entire application lifecycle — from code to cluster.  It enables your team to define Kubernetes configuration policies to ensure consistent, secure, and compliant application deployments every time. In addition to policy enforcement,  Monokle’s ecosystem of tools make your team’s daily YAML configuration workflows easier. Monokle is created and maintained by Kubeshop.


/cncf-monzo-supporter-reviewer — Review Monzo (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-moonbit-reviewer — Review MoonBit using its official documentation and repository in the Wasm / Languages category. WASM focused languages


/cncf-moonlight-marketing-member-reviewer — Review MoonLight Marketing (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-moosefs-reviewer — Review MooseFS using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-moresec-reviewer — Review MoreSec using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Moresec Shangfu Cloud Native Protection Platform is a security protection model based on the concept of cloud native. In view of the characteristics  of cloud native application life-cycle, it involves in many stages of security protection, including R&amp;D construction, security testing, image  control, container deployment, container operation, etc. Shangfu also integrates security with development which means providing feedback of the  application security condition through correlating operation asset, and promoting the iteration and upgrade of application. Whereupon, a closed  loop of efficient cloud native application security risk is created, covering the whole stages of DevSecOps.


/cncf-morgan-stanley-member-reviewer — Review Morgan Stanley (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mosn-reviewer — Review MOSN using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-mplat-reviewer — Review mPLAT using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. mPLAT offers cloud-type service management functions necessary for businesses.


/cncf-msap-ai-reviewer — Review MSAP.ai using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. MSAP.ai(Container Orchestration Platform) is a next-generation application operations platform that simplify complexity and support business innovation in  cloud-native environment built on Kubernetes. It maximizes the value of Kubernetes adoption for enterprises while minimizing operational costs.


/cncf-mulesoft-reviewer — Review MuleSoft using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-multigres-reviewer — Review Multigres using its official documentation and repository in the App Definition and Development / Database category.


/cncf-multus-reviewer — Review Multus using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-myfitnesspal-member-reviewer — Review MyFitnessPal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-myota-member-reviewer — Review Myota (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-mysql-reviewer — Review MySQL using its official documentation and repository in the App Definition and Development / Database category.


/cncf-nacos-reviewer — Review Nacos using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category.


/cncf-nagios-reviewer — Review Nagios using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-naic-member-reviewer — Review NAIC (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-nasdaq-supporter-reviewer — Review Nasdaq (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-national-information-society-agency-member-reviewer — Review National Information Society Agency (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-nats-reviewer — Review NATS using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. NATS.io is a connective technology for distributed systems and is a perfect fit to connect devices, edge, cloud or hybrid deployments. True multi-tenancy makes NATS ideal for SaaS and self-healing and scaling technology allows for topology changes anytime with zero downtime.


/cncf-natwest-member-reviewer — Review Natwest (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-navimentum-kcsp-reviewer — Review Navimentum (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Navimentum Software Corporation provides Kubernetes special technical services including consulting, training, implementation, and technical support, as well as customized research &amp; development, business operations and other advanced services in China.


/cncf-nccl-reviewer — Review NCCL using its official documentation and repository in the AI Native Infra / Accelerator and SuperPod category. Optimized primitives for inter-GPU communication.


/cncf-near-reviewer — Review NEAR using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-nebius-member-reviewer — Review Nebius (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nebulagraph-reviewer — Review NebulaGraph using its official documentation and repository in the App Definition and Development / Database category.


/cncf-neevcloud-member-reviewer — Review NeevCloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nemo-guardrails-reviewer — Review NeMo Guardrails using its official documentation and repository in the AI Agent / Guardrail category. NeMo Guardrails is an open-source toolkit for easily adding programmable guardrails to LLM-based conversational systems.


/cncf-neo4j-reviewer — Review Neo4j using its official documentation and repository in the App Definition and Development / Database category.


/cncf-net-reviewer — Review .NET using its official documentation and repository in the Wasm / Languages category. Managed language


/cncf-netapp-reviewer — Review NetApp using its official documentation and repository in the Runtime / Cloud Native Storage category. Build your Data Fabric on the industry’s broadest portfolio of all-flash, hybrid-flash, and object storage systems.


/cncf-netbird-member-reviewer — Review NetBird (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-netdata-reviewer — Review Netdata using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-netflix-eureka-reviewer — Review Netflix Eureka using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category.


/cncf-netflix-zuul-reviewer — Review Netflix Zuul using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-netis-reviewer — Review Netis using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-netlify-functions-reviewer — Review Netlify Functions using its official documentation and repository in the Serverless / Hosted Platform category. Run your application code without servers, scale it automatically, and pay nothing when it&#39;s not in use.


/cncf-netmatch-supporter-reviewer — Review NetMatch (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-netnod-supporter-reviewer — Review Netnod (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-netobserv-reviewer — Review NetObserv using its official documentation and repository in the Observability and Analysis / Observability category. NetObserv provides CNI-agnostic network observability of Kubernetes clusters.


/cncf-netways-managed-kubernetes-reviewer — Review NETWAYS Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Build, deploy, and scale flexibly with NETWAYS Managed Kubernetes powered by OpenStack.


/cncf-netways-managed-services-kcsp-reviewer — Review NETWAYS Managed Services (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Through its MyEngineer® service, NETWAYS Managed Services offers direct access to experienced engineers who support Kubernetes users with cluster provisioning, CI/CD integration, monitoring, upgrades, troubleshooting, and operational best practices—available 24/7 via chat or video call.


/cncf-netways-managed-services-member-reviewer — Review NETWAYS Managed Services (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-network-service-mesh-reviewer — Review Network Service Mesh using its official documentation and repository in the Runtime / Cloud Native Network category. The Hybrid/Multi-cloud IP Service Mesh


/cncf-neuvector-reviewer — Review NeuVector using its official documentation and repository in the Provisioning / Security &amp; Compliance category. NeuVector Full Lifecycle Container Security Platform delivers the only cloud-native security with uncompromising end-to-end protection from DevOps vulnerability protection to automated run-time security, and featuring a true Layer 7 container firewall.


/cncf-new-relic-member-reviewer — Review New Relic (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-new-relic-reviewer — Review New Relic using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-nexclipper-reviewer — Review NexClipper using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-nginx-member-reviewer — Review NGINX (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nginx-reviewer — Review NGINX using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-nginx-wasm-reviewer — Review NGINX (Wasm) using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-ngrok-member-reviewer — Review ngrok (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ngrok-reviewer — Review ngrok using its official documentation and repository in the Orchestration &amp; Management / API Gateway category. ngrok is a unified ingress platform that combines your firewall, API gateway and global load balancing into a versatile production service. Loved by over 6 million developers, you can deliver and monitor apps from any cloud, datacenter or your home network.


/cncf-nhn-cloud-kcsp-reviewer — Review NHN Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. NHN Cloud strongly supports customers in understanding and using Kubernetes by offering consulting and training courses run by NHN Cloud edu center to help customers design and configure Kubernetes clusters suitable for their environments and providing fully managed cloud native services.


/cncf-nhn-cloud-member-reviewer — Review NHN Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nhn-kubernetes-service-nks-reviewer — Review NHN Kubernetes Service (NKS) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. NKS(NHN Kubernetes Service) is a managed Kubernetes service that allows you to easily deploy a Kubernetes cluster in NHN Cloud.


/cncf-nielseniq-supporter-reviewer — Review NielsenIQ (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-nightingale-reviewer — Review Nightingale using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-nimbella-reviewer — Review Nimbella using its official documentation and repository in the Serverless / Hosted Platform category.


/cncf-nimtech-kcsp-reviewer — Review Nimtech (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We help organizations unlock the full potential of Kubernetes through hands-on consulting, tailored training, and end-to-end implementation, backed by ongoing support from our team of dedicated platform engineers.


/cncf-nimtech-member-reviewer — Review Nimtech (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nipa-member-reviewer — Review NIPA (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-nipr-member-reviewer — Review NIPR (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-nirmata-cloud-native-policy-manager-reviewer — Review Nirmata Cloud Native Policy Manager using its official documentation and repository in the Provisioning / Security &amp; Compliance category. The solution enables DevSecOps teams to ensure the security, compliance, and operational readiness of their Kubernetes workloads and clusters. By automating the creation, deployment, and lifecycle management of policy-based Intelligent Guardrails, customers can gain insights, alerts, and reports, and enable effective collaboration across development and operations teams.


/cncf-nirmata-kcsp-reviewer — Review Nirmata (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Nirmata simplifies kubernetes for application development teams.


/cncf-nirmata-member-reviewer — Review Nirmata (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nitric-reviewer — Review Nitric using its official documentation and repository in the Serverless / Framework category.


/cncf-nivida-nemo-reviewer — Review Nivida NeMo using its official documentation and repository in the AI Native Infra / Orchestration and Scheduling category. A framework for generative AI


/cncf-nix-member-reviewer — Review Ænix (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nkd-nestos-kubernetes-deployer-reviewer — Review NKD - NestOS Kubernetes Deployer using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. NKD (NestOS Kubernetes Deployer) is a solution specially built for deploying and maintaining Kubernetes clusters on NestOS. Its design goal is to provide a convenient cluster operation experience, allowing users to easily complete complex management tasks, thereby improving the efficiency of overall deployment and maintenance.


/cncf-nmstate-reviewer — Review NMstate using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. NMstate is a library with an accompanying command line tool that manages host networking settings in a declarative manner. When used in the Kubernetes environment it allows for declarative node network configuration through the Kubernetes API.


/cncf-no-code-reviewer — Review No Code using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-nocalhost-reviewer — Review Nocalhost using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-node-lambda-reviewer — Review Node Lambda using its official documentation and repository in the Serverless / Tools category.


/cncf-nodesource-reviewer — Review NodeSource using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-nokia-member-reviewer — Review Nokia (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nokv-reviewer — Review NoKV using its official documentation and repository in the Runtime / Cloud Native Storage category. AI-native distributed filesystem.


/cncf-nomad-reviewer — Review Nomad using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-non-public-organization-alligator-member-reviewer — Review Non-Public Organization Alligator (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-non-public-organization-bear-supporter-reviewer — Review Non-Public Organization Bear (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-camel-member-reviewer — Review Non-Public Organization Camel (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-non-public-organization-dog-supporter-reviewer — Review Non-Public Organization Dog (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-elephant-supporter-reviewer — Review Non-Public Organization Elephant (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-flamingo-supporter-reviewer — Review Non-Public Organization Flamingo (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-giraffe-supporter-reviewer — Review Non-Public Organization Giraffe (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-horse-supporter-reviewer — Review Non-Public Organization Horse (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-iguana-supporter-reviewer — Review Non-Public Organization Iguana (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-jellyfish-supporter-reviewer — Review Non-Public Organization Jellyfish (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-kiwi-supporter-reviewer — Review Non-Public Organization Kiwi (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-non-public-organization-moose-member-reviewer — Review Non-Public Organization Moose (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-non-public-organization-narwhal-member-reviewer — Review Non-Public Organization Narwhal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-non-public-organization-octopus-member-reviewer — Review Non-Public Organization Octopus (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-non-public-organization-quail-member-reviewer — Review Non-Public Organization Quail (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-nops-member-reviewer — Review nOps (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nops-reviewer — Review nOps using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. nOps is an end-to-end AWS and Kubernetes optimization platform designed to simplify and automate cloud tracking, allocation, and optimization.


/cncf-northflank-member-reviewer — Review Northflank (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-northflank-reviewer — Review Northflank using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Northflank is a unified developer platform for building, deploying, and managing applications on Kubernetes.


/cncf-notary-project-reviewer — Review Notary Project using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-novaglobal-kcsp-reviewer — Review NovaGlobal (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Our commitment to excellence is reflected in a comprehensive suite of professional services tailored for organizations seeking to adopt Kubernetes seamlessly into their computing systems.


/cncf-novaglobal-member-reviewer — Review NovaGlobal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-novo-nordisk-member-reviewer — Review Novo Nordisk (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nscale-member-reviewer — Review Nscale (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ntt-data-kcsp-reviewer — Review NTT DATA (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. NTT DATA is a global IT services provider that designs, builds and operates production-grade Kubernetes platforms across major clouds such as Amazon EKS, Azure AKS and Google Kubernetes Engine to accelerate application modernization. We combine enterprise-grade security and operations, leveraging partnerships with leading hyperscalers and tooling such as Sysdig.


/cncf-ntt-data-member-reviewer — Review NTT DATA (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nuage-networks-reviewer — Review Nuage Networks using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-nuclio-reviewer — Review Nuclio using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-numaflow-reviewer — Review Numaflow using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Language agnostic K8s native real-time data and stream processing engine


/cncf-nunet-member-reviewer — Review NuNet (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nuodb-reviewer — Review NuoDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-nutanix-data-services-for-kubernetes-ndk-reviewer — Review Nutanix Data Services for Kubernetes (NDK) using its official documentation and repository in the Runtime / Cloud Native Storage category. Protect containerized applications with enterprise grade data services to meet compliance and regulatory requirements.


/cncf-nutanix-kcsp-reviewer — Review Nutanix (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Nutanix provides comprehensive professional services, support, and training for Kubernetes.


/cncf-nutanix-kubernetes-platform-nkp-reviewer — Review Nutanix Kubernetes Platform (NKP) using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. The Nutanix Kubernetes Platform (NKP) solution is a complete, open, and enterprise-grade Kubernetes platform that combines all the best-in-class add-on services required to deliver a scalable production stack anywhere


/cncf-nutanix-member-reviewer — Review Nutanix (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nuvitek-member-reviewer — Review Nuvitek (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nuvotex-kcsp-reviewer — Review Nuvotex (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Nuvotex provides a variety of Kubernetes services, including training, consulting, technical support, and implementation.


/cncf-nuvotex-member-reviewer — Review Nuvotex (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-nuweba-reviewer — Review Nuweba using its official documentation and repository in the Serverless / Hosted Platform category. Nuweba is an ultra-fast and highly-secure FaaS Platform. Nuweba rearchitected serverless from the kernel up to enable companies to use serverless for applications that require scalability, high performance, advanced application security and deep visibility in real-time. Nuweba is compatible with leading serverless platforms, so you can start using Nuweba with only one click and without any changes to your code or configuration.


/cncf-nvidia-member-reviewer — Review NVIDIA (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-oauth2-proxy-reviewer — Review OAuth2 Proxy using its official documentation and repository in the Provisioning / Security &amp; Compliance category. A generic reverse proxy that provides authentication with Google, Azure, OpenID Connect and many more identity providers.


/cncf-observe-sdk-reviewer — Review Observe SDK using its official documentation and repository in the Wasm / Debugging &amp; Observability category.


/cncf-oceanbase-reviewer — Review OceanBase using its official documentation and repository in the App Definition and Development / Database category.


/cncf-octopus-deploy-member-reviewer — Review Octopus Deploy (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-octopus-deploy-reviewer — Review Octopus Deploy using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Octopus Deploy sets the standard for deployment automation for DevOps. We help software teams deploy freely - when and where they need, in a routine way. From modern containers and microservices to trusted legacy applications, Octopus orchestrates software delivery in data centers, multi-cloud, and hybrid IT infrastructure.


/cncf-odigos-member-reviewer — Review Odigos (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-odigos-reviewer — Review Odigos using its official documentation and repository in the Observability and Analysis / Observability category. Enterprise-Grade OpenTelemetry for superior application performance monitoring.


/cncf-okahu-member-reviewer — Review Okahu (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-okahu-reviewer — Review Okahu using its official documentation and repository in the AI Native Infra / Observability category. AI Observability as a Service to help enterprises building LLM-based apps in the cloud move to production.


/cncf-okestro-kcsp-reviewer — Review OKESTRO (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Okestro offers technology-driven open source data infrastructure management solutions in the cloud, with DevOps and development tools to increase productivity and adaptability. We provide Kubernetes distribution, expert consulting services, education services, and an enterprise cloud-native platform based on Kubernetes with containerization, implementation, and training services.


/cncf-okestro-member-reviewer — Review OKESTRO (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-okestro-viola-reviewer — Review OKESTRO VIOLA using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. VIOLA is a managed Kubernetes service that allows you to easily deploy a Kubernetes cluster in Okestro.


/cncf-okta-member-reviewer — Review Okta (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-okteto-reviewer — Review Okteto using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-ollygarden-member-reviewer — Review OllyGarden (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-omni-reviewer — Review Omni using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Omni is a Kubernetes orchestrator for bare metal and edge environments. It manages full OS and Kubernetes lifecycles with Talos Linux and remote management with wireguard communication.


/cncf-omnistrate-member-reviewer — Review Omnistrate (member) using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Omnistrate is a CoPilot for platform teams — helping them build software and agent platforms and turn cloud-native projects into production-ready  platforms without years of custom glue code. Instead of stitching together CI/CD, runtimes, infra-as-code, Day-2 operations, integrations, and  billing from scratch for every customer — and making it all work across clouds, regions, accounts, and deployment configurations — Omnistrate  provides ready-to-use modules that integrate directly with the CNCF ecosystem: Kubernetes, Argo, Flux, Helm, Operators, Kustomize, Terraform,  Prometheus, and more. By building on top of CNCF projects rather than reinventing them, Omnistrate enables modern software companies to deliver  applications in any model — OnPrem, BYOC, PaaS, SaaS, or Agent-aaS — in days, not years.


/cncf-ondat-reviewer — Review Ondat using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-onecause-supporter-reviewer — Review OneCause (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-oodle-ai-reviewer — Review Oodle AI using its official documentation and repository in the Observability and Analysis / Observability category. AI-native observability platform for infrastructure, applications and AI agents. Compatible with OpenTelemetry, drop-in replacement for Prometheus, Grafana, Elasticsearch.


/cncf-opa-gatekeeper-reviewer — Review OPA/Gatekeeper using its official documentation and repository in the AI Native Infra / Governance, Policy and Security category. An admission controller that validates requests to create and update Pods on Kubernetes clusters, using the Open Policy Agent (OPA).


/cncf-opea-reviewer — Review OPEA using its official documentation and repository in the AI Native Infra / Governance, Policy and Security category. OPEA is an ecosystem orchestration framework to integrate performant GenAI technologies &amp; workflows leading to quicker GenAI adoption and business value.


/cncf-open-application-model-reviewer — Review Open Application Model using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-open-cluster-management-reviewer — Review Open Cluster Management using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. A lightweight, open standard for multi-cluster Kubernetes management built on a hub-spoke architecture.


/cncf-open-mpi-reviewer — Review Open MPI using its official documentation and repository in the Training / Distributed Training category. A High Performance Message Passing Library


/cncf-open-policy-administration-layer-opal-reviewer — Review Open Policy Administration Layer (OPAL) using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Policy and data administration, distribution, and real-time updates on top of Policy Agents (OPA, Cedar, ...)


/cncf-open-policy-agent-opa-reviewer — Review Open Policy Agent (OPA) using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-open-policy-containers-reviewer — Review Open Policy Containers using its official documentation and repository in the Provisioning / Security &amp; Compliance category. A docker-inspired CLI for building, tagging, pushing, pulling, and signing OPA policies to and from OCI-compliant registries.


/cncf-open-service-broker-api-reviewer — Review Open Service Broker API using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-open-service-mesh-reviewer — Review Open Service Mesh using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-open-source-consulting-kcsp-reviewer — Review Open Source Consulting (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Open Source Consulting is a specialized cloud native company with extensive expertise in infrastructure, application, and architecture consulting. Our certified Kubernetes experts, who hold professional CKA qualifications, provide comprehensive consulting and technical support services for building optimal Kubernetes and cloud native platforms.


/cncf-open-source-consulting-member-reviewer — Review Open Source Consulting (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-open-vswitch-reviewer — Review Open vSwitch using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-openapi-reviewer — Review OpenAPI using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-openchoreo-reviewer — Review OpenChoreo using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. A developer platform for Kubernetes that delivers higher-level abstractions with a Backstage-powered portal, CI/CD, GitOps, and built-in observability.


/cncf-opencompass-reviewer — Review opencompass using its official documentation and repository in the Training / Evaluation category. OpenCompass is an LLM evaluation platform, supporting a wide range of models (Llama3, Mistral, InternLM2,GPT-4,LLaMa2, Qwen,GLM, Claude, etc) over 100+ datasets


/cncf-opencost-reviewer — Review OpenCost using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. OpenCost provides visibility into current and historical Kubernetes spend and resource allocation.


/cncf-opencv-reviewer — Review OpenCV using its official documentation and repository in the Wasm / AI/Machine Learning category.


/cncf-openebs-reviewer — Review OpenEBS using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-openelb-reviewer — Review OpenELB using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-openeverest-reviewer — Review OpenEverest using its official documentation and repository in the App Definition and Development / Database category. The open-source platform for automated database provisioning and management. It supports multiple database technologies and can be hosted on any Kubernetes infrastructure, in the cloud or on-premises.


/cncf-openfaas-reviewer — Review OpenFaaS using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-openfeature-reviewer — Review OpenFeature using its official documentation and repository in the Observability and Analysis / Feature Flagging category. Standardizing Feature Flagging for Everyone


/cncf-openfga-reviewer — Review OpenFGA using its official documentation and repository in the Provisioning / Security &amp; Compliance category. OpenFGA is a high performance and flexible authorization/permission system built for developers and inspired by Google Zanzibar


/cncf-openfunction-reviewer — Review OpenFunction using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-opengauss-2d2980de-reviewer — Review OpenGauss using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-opengauss-8a1a46cf-reviewer — Review openGauss using its official documentation and repository in the App Definition and Development / Database category.


/cncf-opengemini-reviewer — Review openGemini using its official documentation and repository in the App Definition and Development / Database category. openGemini is an open source distributed time series DBMS with high concurrency, high performance, and high scalability, focusing on the storage and analysis of massive observability data.


/cncf-opengitops-reviewer — Review OpenGitOps using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. OpenGitOps is a set of open-source standards, best practices, and community-focused education to help organizations adopt a structured, standardized approach to implementing GitOps


/cncf-openio-reviewer — Review OpenIO using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-openkruise-reviewer — Review OpenKruise using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-openlit-reviewer — Review OpenLIT using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-openllmetry-reviewer — Review OpenLLMetry using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-openmaru-member-reviewer — Review OPENMARU (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-openmessaging-reviewer — Review OpenMessaging using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. Cloud-oriented, simple, flexible, vendor-neutral and language-independent standards for messaging


/cncf-openmetal-io-member-reviewer — Review OpenMetal.io (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-openmetrics-reviewer — Review OpenMetrics using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-opennebula-elastic-kubernetes-service-oneks-reviewer — Review OpenNebula Elastic Kubernetes Service (OneKS) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Virtual Appliance for KVM with preinstalled Kubernetes service is available from OpenNebula&#39;s Public Marketplace.


/cncf-opennebula-member-reviewer — Review OpenNebula (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-opennebula-reviewer — Review OpenNebula using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-openobserve-member-reviewer — Review OpenObserve (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-openobserve-reviewer — Review OpenObserve using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-openops-member-reviewer — Review OpenOps (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-openresty-reviewer — Review OpenResty using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-openrun-reviewer — Review OpenRun using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. GitOps based deployment platform similar to Google Cloud Run and AWS App Runner. Easily deploy internal tools across a team.


/cncf-opensandbox-reviewer — Review OpenSandbox using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. OpenSandbox is a general-purpose sandbox platform for AI applications, offering multi-language SDKs, unified sandbox APIs, and Docker/Kubernetes runtimes for scenarios like Coding Agents, GUI Agents, Agent Evaluation, AI Code Execution, and RL Training.


/cncf-openscap-reviewer — Review OpenSCAP using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-opensdn-reviewer — Review OpenSDN using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-opensearch-reviewer — Review OpenSearch using its official documentation and repository in the Observability and Analysis / Observability category. OpenSearch is a community-driven, Apache 2.0-licensed open source search and analytics suite that makes it easy to ingest, search, visualize, and analyze data.


/cncf-opensergo-reviewer — Review OpenSergo using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-openstack-reviewer — Review OpenStack using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-opentelemetry-reviewer — Review OpenTelemetry using its official documentation and repository in the Observability and Analysis / Observability category. High-quality, ubiquitous, and portable telemetry to enable effective observability


/cncf-opentofu-reviewer — Review OpenTofu using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. OpenTofu is an open source infrastructure as code tool that enables users to safely and predictably provision and manage cloud and on-prem infrastructure. It&#39;s a community-driven fork of Terraform that maintains backward compatibility while offering enhanced features, stability.


/cncf-opentracing-reviewer — Review OpenTracing using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-opentsdb-reviewer — Review OpenTSDB using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-openvino-reviewer — Review OpenVINO using its official documentation and repository in the Inference / Runtime category. OpenVINO™ is an open source toolkit for optimizing and deploying AI inference


/cncf-openyurt-reviewer — Review OpenYurt using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. An open platform that extends upstream Kubernetes to Edge.


/cncf-operant-member-reviewer — Review Operant (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-operator-framework-reviewer — Review Operator Framework using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-opsani-reviewer — Review Opsani using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-opsmx-member-reviewer — Review OpsMx (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-opsmx-reviewer — Review OpsMx using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. OpsMx provides OSS Argo Support and services and Intelligent Software Delivery platform based on both Argo and Spinnaker.


/cncf-optro-supporter-reviewer — Review Optro (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-optuna-reviewer — Review Optuna using its official documentation and repository in the Training / Post Training category. A hyperparameter optimization framework


/cncf-oracle-ai-database-reviewer — Review Oracle AI Database using its official documentation and repository in the App Definition and Development / Database category. Oracle AI Database delivers time-tested mission-critical functionalities that bring AI to data securely, performantly, and reliably wherever it resides—in the public cloud, multicloud, hybrid cloud, private cloud, or on-premises.


/cncf-oracle-cloud-native-environment-reviewer — Review Oracle Cloud Native Environment using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Oracle Cloud Native Environment is a fully integrated suite for the development and management of cloud-native applications that uses Kubernetes to orchestrate and schedule containers.


/cncf-oracle-functions-reviewer — Review Oracle Functions using its official documentation and repository in the Serverless / Hosted Platform category. Write, deploy, and run code without having to provision, manage, or scale servers. Oracle Functions is a fully managed, highly scalable, on-demand, functions-as-a-service platform, built on enterprise-grade Oracle Cloud Infrastructure and powered by the Fn Project open source engine.


/cncf-oracle-kubernetes-engine-oke-reviewer — Review Oracle Kubernetes Engine (OKE) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Simplify operations of enterprise-grade Kubernetes at scale. Easily deploy and manage resource-intensive workloads such as AI with automatic scaling, patching, and upgrades.


/cncf-oracle-member-reviewer — Review Oracle (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-oras-reviewer — Review ORAS using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Multi-language OCI Registry SDKs and CLI


/cncf-orca-security-reviewer — Review Orca Security using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-orientdb-reviewer — Review OrientDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-orka-reviewer — Review Orka using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Orka is a platform that allows Orchestration via Kubernetes on Apple. It allows MacOS on genuine Apple hardware to be used in Kubernetes.


/cncf-orkestrix-mk8s-qks-reviewer — Review ORKESTRIX mK8s (QKS) using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. ORKESTRIX mK8s(QKS) is a multi-tenant, multi-cluster Kubernetes distribution by QUANTUM C&amp;S that enables organizations to provision and operate multiple isolated Kubernetes clusters at scale on customer on-premises infrastructure via the ORKESTRIX platform, with full automation of cluster lifecycle including provisioning, scaling, upgrades, and day-2 operations.


/cncf-orkestrix-reviewer — Review ORKESTRIX using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. ORKESTRIX is a full-stack Kubernetes distribution by QUANTUM C&amp;S, automating infrastructure provisioning, including QKS(Quantum Kubernetes Service) and various back-office components, tailored to each customer&#39;s environment.


/cncf-ortelius-reviewer — Review Ortelius using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-ory-hydra-reviewer — Review ORY Hydra using its official documentation and repository in the Provisioning / Key Management category.


/cncf-oscal-compass-reviewer — Review OSCAL-COMPASS using its official documentation and repository in the Provisioning / Security &amp; Compliance category. The OSCAL COMPASS project is set of tools that enable the creation, validation, and governance of documentation artifacts for compliance needs. It leverages NIST&#39;s OSCAL (Open Security Controls Assessment Language) as a standard data format for interchange between tools and people, and provides an opinionated approach to OSCAL SDK and adoption by policy engines.


/cncf-osso-member-reviewer — Review OSSO (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-oteemo-kcntp-reviewer — Review Oteemo (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Oteemo is a technology solutions provider specializing in helping enterprises with cloud native transformations. Our services include operationalization of Kubernetes platforms, application container migrations and mentoring of client staff on cloud native technologies.


/cncf-oteemo-kcsp-reviewer — Review Oteemo (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Oteemo is a digital product engineering firm delivering Agentic AI solutions through human-centered design. Our services include operationalization of Kubernetes platforms, application container migrations and development on cloud native architectures.


/cncf-oteemo-member-reviewer — Review Oteemo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-outlines-reviewer — Review Outlines using its official documentation and repository in the AI Agent / Structured Output category. Structured Outputs


/cncf-overops-reviewer — Review OverOps using its official documentation and repository in the Observability and Analysis / Observability category. Prevent rapid code changes from impacting customers.


/cncf-ovh-managed-kubernetes-service-reviewer — Review OVH Managed Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Benefit from free HA managed Kubernetes service by hosting your nodes and services on OVH Public Cloud.


/cncf-ovhcloud-kcsp-reviewer — Review OVHcloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Our OVHcloud Professional services can help you deploy and scale Entreprise Kubernetes workloads in our Public or Private Cloud managed services, or help you with Hybrid,  Multicloud or DRP strategies around containers.


/cncf-ovhcloud-member-reviewer — Review OVHcloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ovn-kubernetes-reviewer — Review OVN-Kubernetes using its official documentation and repository in the Runtime / Cloud Native Network category. A robust Kubernetes networking platform


/cncf-oxeye-reviewer — Review Oxeye using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-oxia-reviewer — Review Oxia using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category. Oxia is a scalable metadata store and coordination system


/cncf-oxide-computer-company-member-reviewer — Review Oxide Computer Company (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ozone-reviewer — Review Ozone using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-paasta-reviewer — Review PaaSTA using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-pachyderm-reviewer — Review Pachyderm using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-packer-reviewer — Review Packer using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-paladin-cloud-reviewer — Review Paladin Cloud using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-palark-kcsp-reviewer — Review Palark (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Palark is a European provider of DevOps and SRE services. We design, implement, and support robust Kubernetes-based infrastructure and efficient CI/CD for cloud native apps.


/cncf-palark-member-reviewer — Review Palark (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-palette-reviewer — Review Palette using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Keep control of your Kubernetes infrastructure stack with Spectro Cloud’s as-a-Service experience for public clouds, private clouds, and bare metal environments


/cncf-pandora2-0-reviewer — Review Pandora2.0 using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-panweidb-reviewer — Review PanWeiDB using its official documentation and repository in the App Definition and Development / Database category. PanWeiDB is an OLTP relational database developed by China Mobile based on the open-source kernel of openGauss. It supports multiple core business systems of the world&#39;s largest telecommunications operator. After being verified through long-term business operations, PanWeiDB has achieved leading positions in terms of feature completeness, high stable performance, and high security compared to similar products in the industry.


/cncf-paralus-reviewer — Review Paralus using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Paralus is a free, open source tool that enables controlled, audited access to Kubernetes infrastructure and Zero trust Kubernetes with zero friction.


/cncf-parseable-reviewer — Review Parseable using its official documentation and repository in the Observability and Analysis / Observability category. Parseable is a free and open source log storage and observability platform. Written in Rust, Parseable can be deployed on natively on Kubernetes. It ingests log data via HTTP POST calls and exposes a query API to search and analyze logs. It is compatible with logging agents like FluentBit, LogStash, FileBeat among others.


/cncf-parsec-reviewer — Review Parsec using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Platform AbstRaction for SECurity service


/cncf-passage-reviewer — Review Passage using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Passage makes it easy for developers to implement passwordless customer identity in any application.


/cncf-payit-supporter-reviewer — Review PayIt (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-pdt-partners-supporter-reviewer — Review PDT Partners (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-pelotech-kcsp-reviewer — Review Pelotech (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Setting clients up for success by crafting Kubernetes first cloud native solutions to their problems.


/cncf-pelotech-member-reviewer — Review Pelotech (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pengcheng-laboratory-member-reviewer — Review Pengcheng Laboratory (member) using its official documentation and repository in the CNCF Members / Academic category.


/cncf-pengyun-network-reviewer — Review Pengyun Network using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-percona-kcsp-reviewer — Review Percona (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Percona offers comprehensive support and services for cloud-native environments, including 24/7 expert assistance, managed services, consulting, and training. Our solutions help organizations optimize database operations in Kubernetes environments, reduce costs, and eliminate the complexity of cloud-native database management.


/cncf-percona-member-reviewer — Review Percona (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-percona-server-for-mysql-reviewer — Review Percona Server for MySQL using its official documentation and repository in the App Definition and Development / Database category.


/cncf-perfectscale-by-doit-reviewer — Review PerfectScale by DoiT using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. PerfectScale by DoiT is an easy-to-use solution that helps DevOps teams ensure peak Kubernetes performance at the lowest possible cost.


/cncf-perfectscale-member-reviewer — Review PerfectScale (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-permify-reviewer — Review Permify using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Permify is an open-source authorization as a service inspired by Google Zanzibar, designed to build and manage fine-grained and scalable authorization systems for any application.


/cncf-permit-io-member-reviewer — Review Permit.io (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-permit-io-reviewer — Review Permit.io using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Permit.io provides fine-grained authorization as a service (Policy-as-Code,  APIs, and customer facing UI), so developers can check this as done and focus on their core product. Permit.io is the maintainer of the OPAL (Open Policy Administration Layer) project.


/cncf-perses-reviewer — Review Perses using its official documentation and repository in the Observability and Analysis / Observability category. Perses is a dashboard tool to visualize observability data from Prometheus/Thanos/Jaeger.


/cncf-pgedge-member-reviewer — Review pgEdge (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pgvector-reviewer — Review pgvector using its official documentation and repository in the AI Agent / Vector Database category. Open-source vector similarity search for Postgres


/cncf-picnic-supporter-reviewer — Review Picnic (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-pingcap-member-reviewer — Review PingCAP (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pinniped-reviewer — Review Pinniped using its official documentation and repository in the Provisioning / Key Management category.


/cncf-pinpoint-reviewer — Review Pinpoint using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-pionative-kcsp-reviewer — Review Pionative (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Pionative helps companies succeed with their Kubernetes journey by providing enterprise-level support.


/cncf-pionative-member-reviewer — Review Pionative (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pipecd-reviewer — Review PipeCD using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. GitOps style continuous delivery platform that provides consistent deployment and operations experience for any applications


/cncf-pipekit-member-reviewer — Review Pipekit (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pipelineai-reviewer — Review PipelineAI using its official documentation and repository in the Serverless / Installable Platform category.


/cncf-pipelock-reviewer — Review Pipelock using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Pipelock is an open-source AI agent firewall. It sits between AI agents and the internet and blocks secret leaks (DLP), SSRF, unsafe tool traffic, and prompt-injection content across HTTP, WebSocket, and MCP transports. It adds MCP-specific defenses for tool-baseline drift and rug-pull detection, and supports Ed25519-signed action receipts plus Audit Packet v0 bundles that can be verified offline with a standalone Go verifier and TypeScript, Rust, and Python verifier ports. Single Go binary, fail-closed defaults, Apache 2.0.


/cncf-pipy-reviewer — Review Pipy using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. Pipy is a programmable proxy for the cloud, edge and IoT. It&#39;s written in C++, which makes it extremely lightweight and fast. It&#39;s also fully programmable by using PipyJS, a tailored version from the standard JavaScript language.


/cncf-piraeus-datastore-reviewer — Review Piraeus Datastore using its official documentation and repository in the Runtime / Cloud Native Storage category. The Piraeus Operator manages LINSTOR clusters in Kubernetes.


/cncf-pixie-reviewer — Review Pixie using its official documentation and repository in the Observability and Analysis / Observability category. Open source Kubernetes observability for developers


/cncf-plakar-member-reviewer — Review Plakar (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-planetscale-member-reviewer — Review PlanetScale (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-platform-engineering-labs-member-reviewer — Review Platform Engineering Labs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-platform-engineering-masters-member-reviewer — Review Platform Engineering Masters (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-platform-sh-reviewer — Review Platform.sh using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-platform9-elastic-machine-pool-reviewer — Review Platform9 Elastic Machine Pool using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. Elastic Machine Pool improves resource utilization, without any changes to your applications.  EMP deploys an alternate virtualization layer, using AWS Bare Metal underneath, creating Elastic VMs (EVMs), that look and feel exactly like regular EC2 VMs.


/cncf-platform9-kcsp-reviewer — Review Platform9 (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Platform9 Managed Kubernetes is SaaS-managed, infrastructure-agnostic, and works across public clouds, on-premises servers, and VMware infrastructure.


/cncf-platform9-managed-kubernetes-reviewer — Review Platform9 Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Platform9 Managed Kubernetes is SaaS-managed, infrastructure-agnostic, and works across public clouds and on-premises server infrastructure.


/cncf-platform9-member-reviewer — Review Platform9 (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-plural-member-reviewer — Review Plural (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-plural-reviewer — Review Plural using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category.


/cncf-plusserver-kcsp-reviewer — Review plusserver (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. plusserver experts offer individualized guidance covering all phases of your Kubernetes journey, including planning, migration and problem-solving of your Kubernetes clusters.


/cncf-plusserver-kubernetes-engine-pske-reviewer — Review plusserver Kubernetes Engine (PSKE) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. The plusserver Kubernetes Engine (PSKE) based on Gardener reduces the complexity in managing multi-cloud environments and enables companies to orchestrate their containers and cloud-native applications across a variety of platforms such as plusserver’s pluscloud open or hyperscalers such as AWS, either by mouseclick or via an API.


/cncf-plusserver-member-reviewer — Review plusserver (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pluto-reviewer — Review Pluto using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-podman-container-tools-reviewer — Review Podman Container Tools using its official documentation and repository in the Runtime / Container Runtime category. A set of tools providing full management of container lifecycle, including Podman, Buildah, and Skopeo,  which manage containers and images without requiring a daemon or root privileges.


/cncf-podman-desktop-reviewer — Review Podman Desktop using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. An open-source tool for developers to work with containers and Kubernetes with an intuitive and user-friendly interface to effortlessly build, manage, and deploy containers and Kubernetes — all from the desktop.


/cncf-podman-reviewer — Review Podman using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-podman-wasm-reviewer — Review podman (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-polar-signals-member-reviewer — Review Polar Signals (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-polardb-reviewer — Review PolarDB using its official documentation and repository in the App Definition and Development / Database category. PolarDB is a cloud native SQL Database.


/cncf-polaris-reviewer — Review Polaris using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-pomerium-reviewer — Review Pomerium using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Pomerium is an identity-aware access proxy that enables teams to implement true zero-trust security, eliminating the need for a VPN. Pomerium ensures that every single call is authenticated and authorized, providing a secure and seamless user experience.


/cncf-porch-financial-member-reviewer — Review Porch Financial (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-port-member-reviewer — Review Port (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-portainer-member-reviewer — Review Portainer (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-portainer-reviewer — Review Portainer using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-porter-reviewer — Review Porter using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Porter enables you to package your application artifact, client tools, configuration and deployment logic together as a versioned bundle that you can distribute, and install with a single command


/cncf-portshift-reviewer — Review Portshift using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-portworx-by-pure-storage-kcsp-reviewer — Review Portworx by Pure Storage (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Automate, protect, and unify data for modern applications across any on-premises, public, or hybrid cloud environment, with expert training, consulting, implementation, and support services. Portworx partners with customers through hands-on collaboration and robust support for running data management on Kubernetes. Portworx also offers practical training through immersive labs and Portworx certification programs


/cncf-portworx-by-pure-storage-member-reviewer — Review Portworx by Pure Storage (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-portworx-reviewer — Review Portworx using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-posit-supporter-reviewer — Review Posit (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-postfinance-supporter-reviewer — Review PostFinance (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-postgresql-reviewer — Review PostgreSQL using its official documentation and repository in the App Definition and Development / Database category.


/cncf-powerfulseal-reviewer — Review PowerfulSeal using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-pravega-reviewer — Review Pravega using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-prefect-reviewer — Review Prefect using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-presto-reviewer — Review Presto using its official documentation and repository in the App Definition and Development / Database category.


/cncf-previder-kcsp-reviewer — Review Previder (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Previder is a leading Dutch provider of cloud and IT solutions, specializing in data center services, managed services, connectivity solutions, cloud services and Kubernetes.


/cncf-previder-kubernetes-engine-reviewer — Review Previder Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Previder Kubernetes Engine is a managed Kubernetes service to help users focus on their development as we take care of the maintenance and updates of the container infrastructure


/cncf-previder-member-reviewer — Review Previder (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-prisma-cloud-by-palo-alto-networks-reviewer — Review Prisma Cloud by Palo Alto Networks using its official documentation and repository in the Serverless / Security category.


/cncf-prisma-cloud-reviewer — Review Prisma Cloud using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-prodigy-education-member-reviewer — Review Prodigy Education (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-prodyna-kcntp-reviewer — Review PRODYNA (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. PRODYNA designs, implements, and operates Kubernetes-based custom applications for mid- to large size enterprises in all of Europe.


/cncf-prodyna-kcsp-reviewer — Review PRODYNA (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. PRODYNA designs, implements, and operates Kubernetes-based custom applications for mid- to large size enterprises in all of Europe.


/cncf-prodyna-member-reviewer — Review PRODYNA (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-profisea-kcsp-reviewer — Review Profisea (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Profisea is an Israeli DevOps and cloud company providing a full spectrum of cloud and Kubernetes professional services — from infrastructure design and implementation to continuous optimization and support. Our certified experts help organizations adopt Kubernetes and cloud-native technologies through consulting, deployment, migration, CI/CD automation, DevSecOps integration, and 24×7 support.


/cncf-profisea-member-reviewer — Review Profisea (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-project-calico-reviewer — Review Project Calico using its official documentation and repository in the Runtime / Cloud Native Network category.


/cncf-project-jupyter-reviewer — Review Project Jupyter using its official documentation and repository in the Data / Data Science category. Interactive Computing


/cncf-project-syn-reviewer — Review Project Syn using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-prometheus-reviewer — Review Prometheus using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-promptfoo-reviewer — Review Promptfoo using its official documentation and repository in the AI Agent / Evaluation category. Test your prompts, agents, and RAGs. Red teaming/pentesting/vulnerability scanning for AI. Compare performance of GPT, Claude, Gemini, Llama, and more. Simple declarative configs with command line and CI/CD integration. Used by OpenAI and Anthropic.


/cncf-prosiebensat-1-supporter-reviewer — Review ProSiebenSat.1 (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-protego-reviewer — Review Protego using its official documentation and repository in the Serverless / Security category.


/cncf-ps-cloud-services-kcsp-reviewer — Review PS Cloud Services (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We design and implement fault-tolerant workstations based on Kubernetes. Scaling, management, maintenance and support 24/7.


/cncf-ps-cloud-services-member-reviewer — Review PS Cloud Services (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pubnub-functions-reviewer — Review PubNub Functions using its official documentation and repository in the Serverless / Hosted Platform category.


/cncf-puffersoft-kcsp-reviewer — Review PufferSoft (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. PufferSoft provides expert Kubernetes consulting and managed services, helping enterprises design, deploy, and operate secure, scalable, and production-grade container platforms across cloud and hybrid environments. Our certified engineers specialize in GitOps, CI/CD, observability, and Day-2 operations to accelerate your cloud-native journey.


/cncf-puffersoft-member-reviewer — Review PufferSoft (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pulsar-reviewer — Review Pulsar using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-pulumi-member-reviewer — Review Pulumi (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pulumi-reviewer — Review Pulumi using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-puppet-reviewer — Review Puppet using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-pure-storage-reviewer — Review Pure Storage using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-puzzle-itc-ag-kcsp-reviewer — Review Puzzle ITC AG (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Whether on-premises, in the cloud, or as a managed service – we help customers create the ideal Kubernetes environment and leverage the best tools from the cloud-native landscape. This results in a standardized, flexible, and automated environment that efficiently bridges the gap between development and operations.


/cncf-puzzle-itc-ag-member-reviewer — Review Puzzle ITC AG (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-pydantic-ai-reviewer — Review Pydantic AI using its official documentation and repository in the AI Agent / Agent Framework category. AI Agent Framework, the Pydantic way


/cncf-pyodide-reviewer — Review Pyodide using its official documentation and repository in the Wasm / Languages category. Pyodide is a Python distribution for the browser and Node.js base WebAssembly


/cncf-pytorch-distributeddataparallel-ddp-reviewer — Review Pytorch DistributedDataParallel (DDP) using its official documentation and repository in the Training / Distributed Training category. DistributedDataParallel (DDP) implements data parallelism at the module level which can run across multiple machines.


/cncf-pytorch-reviewer — Review PyTorch using its official documentation and repository in the Data / Data Science category. About Tensors and Dynamic neural networks in Python with strong GPU acceleration


/cncf-qaware-member-reviewer — Review QAware (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-qazaq-open-source-initiative-member-reviewer — Review Qazaq Open Source Initiative (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-qdrant-reviewer — Review Qdrant using its official documentation and repository in the App Definition and Development / Database category. Vector Database for the next generation of AI applications.


/cncf-qiming-kcsp-reviewer — Review Qiming (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. With Kubernetes as the technical base, QiMing provides full life cycle service solutions such as container orchestration, operation and maintenance, cloud native technical consultation, training, and microservice transformation.


/cncf-qiming-member-reviewer — Review Qiming (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-qingcloud-kcsp-reviewer — Review QingCloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. KubeSphere, developed by QingCloud, is an enterprise-grade Kubernetes platform that simplifies the deployment, management, and scaling of containerized applications. In addition to our platform, we offer a full range of Kubernetes professional services — including consulting, deployment, training, and ongoing support — to help organizations successfully adopt, operate, and optimize their cloud native environments.


/cncf-qingcloud-kubernetes-engine-qke-reviewer — Review QingCloud Kubernetes Engine (QKE) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. QKE can provision a K8s/KubeSphere cluster in minutes automatically by just a few clicks and inputs.


/cncf-qingstor-reviewer — Review QingStor using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-qingteng-reviewer — Review Qingteng using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Qingteng Honeycomb · Container Security Platform can display risk scenarios visually for enterprises through continuous monitoring and analysis of the container security status, and protects the container environment during the entire life cycle (build, ship, and run) of the container.


/cncf-quai-network-reviewer — Review Quai Network using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-qualitysoft-member-reviewer — Review QualitySoft (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-quantum-c-and-s-kcsp-reviewer — Review QUANTUM C&amp;S (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Quantum C&amp;S delivers Kubernetes consulting, deployment, and operations services, with a strong focus on AI infrastructure and mission-critical enterprise environments.


/cncf-quantum-c-and-s-member-reviewer — Review QUANTUM C&amp;S (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-quantum-kubernetes-service-qks-reviewer — Review Quantum Kubernetes Service (QKS) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. QKS (Quantum Kubernetes Service) is a managed Kubernetes-as-a-Service platform by QUANTUM C&amp;S, enabling on-premises multi-cluster environments with fully managed control planes — no infrastructure overhead for tenants.


/cncf-quarkus-reviewer — Review Quarkus using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Quarkus is a Kubernetes-native Java framework, designed to enable Java developers to create applications for a modern, cloud-native world. It is tailored for GraalVM and HotSpot, and crafted from best-of-breed Java libraries and standards.


/cncf-quay-reviewer — Review Quay using its official documentation and repository in the Provisioning / Container Registry category.


/cncf-qube-research-and-technologies-member-reviewer — Review Qube Research &amp; Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-qubole-reviewer — Review Qubole using its official documentation and repository in the App Definition and Development / Database category.


/cncf-quesma-member-reviewer — Review Quesma (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-quickwit-reviewer — Review Quickwit using its official documentation and repository in the Observability and Analysis / Observability category. Sub-second search &amp; analytics engine on cloud storage


/cncf-qumulo-reviewer — Review Qumulo using its official documentation and repository in the Runtime / Cloud Native Storage category. Qumulo provides file storage that delivers billion file count, petabyte scale, gigabyte throughput in the public cloud and AWS GovCloud(US) with all-flash and hybrid SSD/HDD options.


/cncf-quobyte-reviewer — Review Quobyte using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-rabbitmq-reviewer — Review RabbitMQ using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-rackner-kcsp-reviewer — Review Rackner (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Rackner is a cloud native consultancy where we work closely with our clients to build, secure, and operate on any cloud with Kubernetes.


/cncf-rackner-member-reviewer — Review Rackner (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rad-security-member-reviewer — Review RAD Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rad-security-reviewer — Review RAD Security using its official documentation and repository in the Provisioning / Security &amp; Compliance category. RAD Security maps a broad set of cluster components across the Kubernetes lifecycle  using a real-time graph, cutting noise in half through contextualized risks, highest impact remediations and Kubernetes-first  incident response.


/cncf-radar-reviewer — Review Radar using its official documentation and repository in the Observability and Analysis / Observability category. Open-source Kubernetes visibility tool with topology visualization, event timeline, Helm management, and built-in MCP server for AI tools.


/cncf-radius-reviewer — Review Radius using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Radius is a cloud-native application platform that enables developers and the platform engineers that support them to collaborate on delivering and managing cloud-native applications that follow organizational best practices for cost, operations and security, by default.


/cncf-rafay-kcsp-reviewer — Review Rafay (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Rafay Systems delivers a turnkey offering along with white-glove support to assist enterprise customers automate Kubernetes cluster management and application operations at scale across public clouds, data centers and Edge environments.


/cncf-rafay-reviewer — Review Rafay using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Rafay Systems enables DevOps and SRE teams to automate application lifecycle and cluster configuration management across on-premise and cloud-based Kubernetes clusters through a turnkey SaaS platform.


/cncf-rafay-systems-member-reviewer — Review Rafay Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-raftt-reviewer — Review Raftt using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-ragas-reviewer — Review Ragas using its official documentation and repository in the AI Agent / Evaluation category. Supercharge Your LLM Application Evaluations


/cncf-rancher-kubernetes-reviewer — Review Rancher Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Deploy Rancher’s Kubernetes distro anywhere or launch cloud Kubernetes services from Google, Amazon or Microsoft.


/cncf-rapidfort-member-reviewer — Review RapidFort (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ratify-reviewer — Review Ratify using its official documentation and repository in the Provisioning / Security &amp; Compliance category. A verification engine on Kubernetes which enables verification of artifact security metadata and admits for deployment only those that comply with policies you create.


/cncf-raydian-cloud-kcsp-reviewer — Review Raydian Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Raydian Cloud provides expert consulting, training, support, and implementation services for AI and Kubernetes environments, helping organizations build, deploy, and operate production-grade containerized workloads across leading cloud platforms. With deep expertise in Rafay Systems and GPU-backed infrastructure, we deliver comprehensive Day 1 and Day 2 operational support to ensure secure, reliable, and scalable AI and Kubernetes deployments.


/cncf-raydian-cloud-member-reviewer — Review Raydian Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-razee-reviewer — Review Razee using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-rbac-lookup-reviewer — Review RBAC Lookup using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-rbac-manager-reviewer — Review RBAC Manager using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-rbg-reviewer — Review RBG using its official documentation and repository in the Inference / Framework category. Kubernetes API for orchestrating distributed, stateful AI inference workloads with multi-role collaboration and built-in service discovery.


/cncf-reactive-interaction-gateway-reviewer — Review Reactive Interaction Gateway using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-red-hat-build-of-microshift-reviewer — Review Red Hat Build of Microshift using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. A component of Red Hat Device Edge used to provide lightweight Kubernetes container orchestration in field deployed devices at the far edge


/cncf-red-hat-member-reviewer — Review Red Hat (member) using its official documentation and repository in the CNCF Members / Platinum category.


/cncf-red-hat-openshift-dedicated-reviewer — Review Red Hat OpenShift Dedicated using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. OpenShift® Dedicated helps organizations focus on building and scaling their business with a private Kubernetes cluster fully-managed by Red Hat®.


/cncf-red-hat-openshift-on-ibm-cloud-reviewer — Review Red Hat OpenShift on IBM Cloud using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Red Hat OpenShift on IBM Cloud is a managed offering to create your own OpenShift cluster of compute hosts to deploy and manage containerized apps on IBM Cloud.


/cncf-red-hat-openshift-reviewer — Review Red Hat OpenShift using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. OpenShift® helps organizations focus on building and scaling their business with fully supported enterprise Kubernetes by Red Hat®.


/cncf-reddit-supporter-reviewer — Review Reddit (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-redeploy-kcsp-reviewer — Review Redeploy (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We know Azure AKS.


/cncf-redeploy-member-reviewer — Review Redeploy (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-redis-reviewer — Review Redis using its official documentation and repository in the App Definition and Development / Database category.


/cncf-redpanda-member-reviewer — Review Redpanda (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-redpanda-reviewer — Review Redpanda using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-redpanda-wasm-reviewer — Review Redpanda (Wasm) using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-redpill-linpro-member-reviewer — Review Redpill Linpro (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-reevo-cloud-and-cybersecurity-kcntp-reviewer — Review ReeVo Cloud &amp; CyberSecurity (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. ReeVo is a European Cloud, Cybersecurity, and Cloud Native provider with the core focus of protecting organizations&#39; and enterprises&#39; data in its digital vault, aiming to accelerate digital transformation and make performance, security, and resilience ReeVo&#39;s distinctive hallmarks. ReeVo is a key training reference in Southern Europe, and its training programs provide practical, hands-on knowledge while preparing your team for industry-recognized certifications.


/cncf-reevo-cloud-and-cybersecurity-kcsp-reviewer — Review ReeVo Cloud &amp; CyberSecurity (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ReeVo is a European Cloud, Cybersecurity, and Cloud Native provider with the core focus of protecting organizations&#39; and enterprises&#39; data in its digital vault, aiming to accelerate digital transformation and make performance, security, and resilience ReeVo&#39;s distinctive hallmarks. ReeVo is a KCSP with certified engineers that operates following industry best practices and ensures customer success by providing expert guidance, knowledge transfer, and Enterprise support.


/cncf-reevo-member-reviewer — Review ReeVo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-relvy-ai-member-reviewer — Review Relvy AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-reo-dev-member-reviewer — Review reo.dev (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-replex-reviewer — Review Replex using its official documentation and repository in the Observability and Analysis / Observability category. The central source of truth for all IT infrastructure insights. Providing one central platform that uniquely combines technical and financial metrics to manage and optimize modern containerized IT infrastructure stacks successfully.


/cncf-replicated-member-reviewer — Review Replicated (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-resolve-technology-kcsp-reviewer — Review Resolve Technology (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Resolve Technology offers comprehensive cloud-native consultation and implementation services, covering Kubernetes cluster design, cloud-native application deployment, zero-trust security, CI/CD pipeline optimization, and robust monitoring solutions for Kubernetes and related CNCF technologies. From inception to ongoing operations, we ensure seamless integration and optimal performance across the cloud-native ecosystem.


/cncf-resolve-technology-member-reviewer — Review Resolve Technology (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rethinkdb-reviewer — Review RethinkDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-rewe-international-dienstleistungsges-m-66a030c9-reviewer — Review REWE International Dienstleistungsges.m.b.h. (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-ricardo-ch-supporter-reviewer — Review Ricardo.ch (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-ripple-reviewer — Review Ripple using its official documentation and repository in the Wasm / Decentralized Platforms category.


/cncf-riptides-member-reviewer — Review Riptides (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-risc-v-member-reviewer — Review RISC-V (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-rizhiyi-reviewer — Review Rizhiyi using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-rke2-reviewer — Review RKE2 using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. A Kubernetes distribution focused on enabling Federal government compliance-based use cases.


/cncf-rkt-reviewer — Review rkt using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-roadie-member-reviewer — Review Roadie (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-robin-systems-member-reviewer — Review Robin Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-robin-systems-reviewer — Review Robin Systems using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-robinhood-supporter-reviewer — Review Robinhood (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-robusta-member-reviewer — Review Robusta (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rook-reviewer — Review Rook using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-rookout-reviewer — Review Rookout using its official documentation and repository in the Observability and Analysis / Observability category. Collect data directly from your live code using non-breaking breakpoints. Debug, observe and understand anything, anytime, anywhere.


/cncf-root-member-reviewer — Review Root (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rootly-member-reviewer — Review Rootly (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-royal-bank-of-canada-member-reviewer — Review Royal Bank of Canada (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rtx-member-reviewer — Review RTX (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rudder-reviewer — Review Rudder using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-runc-reviewer — Review runc using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-rundeck-reviewer — Review Rundeck using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-runme-notebooks-reviewer — Review Runme Notebooks using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. A toolchain that turns Markdown into interactive, cloud-native, runnable Notebook experiences for DevOps.


/cncf-runwasi-reviewer — Review runwasi using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-runwhen-member-reviewer — Review RunWhen (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-rust-reviewer — Review Rust using its official documentation and repository in the Wasm / Languages category. Compiled language to Wasm


/cncf-rx-m-kcntp-reviewer — Review RX-M (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. RX-M is a global Cloud Native &amp; AI technology training and consulting firm; we provide services and training for technologies and business practices  essential to digital transformation, ranging from ML/AI to microservice-based application design and from DevSecOps to AIOps. The RX-M team consists  of published authors, patent holders and prominent open source contributors–all focused on customer success through an unbiased, market neutral approach.


/cncf-rx-m-kcsp-reviewer — Review RX-M (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. RX-M is a global Cloud Native &amp; AI technology training and consulting firm; we provide services and training for technologies and business practices  essential to digital transformation, ranging from ML/AI to microservice-based application design and from DevSecOps to AIOps. The RX-M team consists  of published authors, patent holders and prominent open source contributors–all focused on customer success through an unbiased, market neutral approach.


/cncf-rx-m-member-reviewer — Review RX-M (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-saic-kcsp-reviewer — Review SAIC (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Application Modernization services rapidly migrate your apps to Kubernetes and incrementally transform them into cloud native apps.


/cncf-saic-member-reviewer — Review SAIC (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-salad-reviewer — Review Salad using its official documentation and repository in the Platform / PaaS/Container Service category. Deploy containers across 10,000s of GPU powered nodes on the world&#39;s largest distributed computing platform.


/cncf-salt-lake-city-devops-days-member-reviewer — Review Salt Lake City DevOps Days (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-salt-project-reviewer — Review Salt Project using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-samsung-kubernetes-engine-ske-reviewer — Review Samsung Kubernetes Engine (SKE) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. SKE offers lightweight virtual computing containers and clusters for their management. With the service, Kubernetes environments are readily available by installation, operation, and maintenance of the Kubernetes control plane.


/cncf-samsung-sds-kcsp-reviewer — Review Samsung SDS (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Samsung SDS’s Cloud Native Computing Team offers expert consulting across the range of technical aspects involved in building services targeted at a Kubernetes cluster.


/cncf-samsung-sds-member-reviewer — Review Samsung SDS (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sandstone-reviewer — Review SandStone using its official documentation and repository in the Runtime / Cloud Native Storage category. ShenZhen SandStone Data Technology Co., Ltd.  is a high-tech enterprise focusing on enterprise-level software-defined storage.


/cncf-sangfor-cnad-reviewer — Review Sangfor CNAD using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category. Sangfor CNAD is a cloud native application proxy, that provides layer7&amp;layer4 load balance service, high performance TLS encryption and decryption resource pool, and Ingress&amp;Loadbalancer for kubernetes.


/cncf-sangfor-eds-reviewer — Review Sangfor EDS using its official documentation and repository in the Runtime / Cloud Native Storage category. Sangfor Technologies is a leading global vendor of IT infrastructure solutions, specializing in Cloud Computing &amp; Network Security with a wide range of products &amp; services.


/cncf-sap-certified-gardener-reviewer — Review SAP Certified Gardener using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. The Gardener implements automated management and operation of Kubernetes clusters as a service and aims to support that service on multiple Cloud providers.


/cncf-sap-hana-reviewer — Review SAP HANA using its official documentation and repository in the App Definition and Development / Database category. SAP HANA is an in-memory, column-oriented, relational database management system


/cncf-sap-kcsp-reviewer — Review SAP (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category.


/cncf-sap-labs-supporter-reviewer — Review SAP Labs (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category. SAP Labs Network is formed by SAP’s core research and development entities focused on developing  and constantly improving SAP&#39;s enterprise solutions and cloud products. With 21 labs in 19 countries,  the network drives thought leadership globally and in local ecosystems, allowing SAP to innovate,  grow, and succeed, while leveraging and contributing to numerous CNCF cloud native projects.


/cncf-sap-member-reviewer — Review SAP (member) using its official documentation and repository in the CNCF Members / Platinum category. As a global leader in enterprise applications and business AI, SAP stands at the nexus of business and  technology. For over 50 years, organizations have trusted SAP to bring out their best by uniting  business-critical operations spanning finance, procurement, HR, supply chain, and customer experience.


/cncf-sawmills-member-reviewer — Review Sawmills (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-scalardb-reviewer — Review ScalarDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-scale-reviewer — Review Scale using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-scaleflash-reviewer — Review ScaleFlash using its official documentation and repository in the Runtime / Cloud Native Storage category. ScaleFlash is a leading provider of high-performance software-defined storage products and enterprise cloud solutions.


/cncf-scaleops-member-reviewer — Review ScaleOps (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-scaleops-reviewer — Review ScaleOps using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-scaleway-kubernetes-kapsule-reviewer — Review Scaleway Kubernetes Kapsule using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Kubernetes Kapsule provides a simple way to deploy and manage your containerized applications in the cloud. Relax and focus on your software stack while we take care of your clusters.


/cncf-scaleway-member-reviewer — Review Scaleway (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-scalingo-reviewer — Review Scalingo using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-scality-ring-reviewer — Review Scality RING using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-scar-reviewer — Review SCAR using its official documentation and repository in the Serverless / Tools category.


/cncf-schemahero-reviewer — Review SchemaHero using its official documentation and repository in the App Definition and Development / Database category.


/cncf-scitix-member-reviewer — Review SCITIX (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-score-reviewer — Review Score using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Score is an open-source workload specification designed to simplify development for cloud-native developers.


/cncf-screwdriver-reviewer — Review Screwdriver using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-scribe-security-platform-reviewer — Review Scribe Security Platform using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-scylla-reviewer — Review Scylla using its official documentation and repository in the App Definition and Development / Database category.


/cncf-scylladb-reviewer — Review ScyllaDB using its official documentation and repository in the Data / Data Architecture category. NoSQL data store using the seastar framework, compatible with Apache Cassandra.


/cncf-sdc-reviewer — Review SDC using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. SDC (schema Driven Configuration) exposes any schema based configuration as Kubernetes resources.


/cncf-seal-member-reviewer — Review Seal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-seal-security-member-reviewer — Review Seal Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sealer-reviewer — Review sealer using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-seata-reviewer — Review Seata using its official documentation and repository in the App Definition and Development / Database category.


/cncf-seatunnel-reviewer — Review SeaTunnel using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-secloudit-reviewer — Review SECloudit using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. SECloudit is a container application management solution, provides improvement of container environment operational capability and manageability by providing multiple Kubernetes-based cluster management, CI/CD management, and a user-friendly GUI portal.


/cncf-second-state-functions-reviewer — Review Second State Functions using its official documentation and repository in the Serverless / Hosted Platform category. High-performance Function as a Service (FaaS), powered by WasmEdge.


/cncf-secondfront-member-reviewer — Review Secondfront (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-securecodebox-reviewer — Review secureCodeBox using its official documentation and repository in the Provisioning / Security &amp; Compliance category. secureCodeBox is an OWASP project providing an automated and scalable open source solution that integrates multiple security scanners with a simple and lightweight interface – for continuous and automated security testing.


/cncf-sedai-member-reviewer — Review Sedai (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sedai-reviewer — Review Sedai using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. Sedai is an autonomous cloud management platform to optimize cost &amp; performance and remediate issues for Kubernetes &amp; Serverless


/cncf-sel4-reviewer — Review seL4 using its official documentation and repository in the Wasm / Edge/Bare metal category.


/cncf-seldon-reviewer — Review Seldon using its official documentation and repository in the Inference / Runtime category. An MLOps framework to package, deploy, monitor and manage thousands of production machine learning models.


/cncf-selenium-grid-reviewer — Review Selenium Grid using its official documentation and repository in the AI Agent / Agent Tool category. A browser automation framework and ecosystem.


/cncf-semantic-kernel-reviewer — Review Semantic Kernel using its official documentation and repository in the AI Agent / Agent Framework category. Integrate cutting-edge LLM technology quickly and easily into your apps.


/cncf-semaphore-reviewer — Review Semaphore using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-sematext-reviewer — Review Sematext using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-sensu-reviewer — Review Sensu using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-sentinel-reviewer — Review Sentinel using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-sentry-reviewer — Review Sentry using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-seowon-information-member-reviewer — Review Seowon Information (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sermant-reviewer — Review Sermant using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Sermant a proxyless service mesh solution based on Javaagent.


/cncf-serverless-devs-reviewer — Review Serverless Devs using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Serverless Devs developer tool ( Serverless Devs 开发者工具 )


/cncf-serverless-reviewer — Review Serverless using its official documentation and repository in the Serverless / Framework category.


/cncf-serverless-workflow-reviewer — Review Serverless Workflow using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Standards-based DSL and open-source dev tools and runtimes are at the heart of the Serverless Workflow project


/cncf-service-mesh-interface-smi-reviewer — Review Service Mesh Interface (SMI) using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-service-mesh-performance-reviewer — Review Service Mesh Performance using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-servicecomb-reviewer — Review ServiceComb using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-serviceradar-reviewer — Review ServiceRadar using its official documentation and repository in the Observability and Analysis / Observability category. Open Source network management and observability platform.


/cncf-sglang-reviewer — Review SGLang using its official documentation and repository in the Inference / Runtime category. SGLang powers fast, scalable inference for large language and multimodal models.


/cncf-shandong-cvicse-middleware-kcsp-reviewer — Review Shandong Cvicse Middleware (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Shandong Cvicse Middleware provides specialized Kubernetes technical services, encompassing consultation, training, implementation, and technical support, as well as advanced services such as customized development and business operations maintenance. Our objective is to empower enterprises to swiftly adopt Kubernetes, better meeting the demands of business growth and accelerating the prosperous expansion of their operations.


/cncf-shandong-cvicse-middleware-member-reviewer — Review Shandong Cvicse Middleware (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-shardingsphere-reviewer — Review ShardingSphere using its official documentation and repository in the App Definition and Development / Database category.


/cncf-shifu-reviewer — Review Shifu using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Shifu is the next generation cloud native IoT development framework that could 10X your IoT software development


/cncf-shinesoft-kcntp-reviewer — Review Shinesoft (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. We provide design, development and tech support include; Cloud/Cluster, Network, Server, Security, Monitoring, Data analysis, and AI/IoT


/cncf-shinesoft-kcsp-reviewer — Review Shinesoft (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Shinesoft provides Kubernetes consulting and support solutions on your choice of infrastructure.  We help our customers to adopt DevOps, build CI/CD, and migrate applications to Kubernetes. We also design and build your Kubernetes cluster in the public cloud, on-prem, or multi/hybrid cloud  environments according to your needs.


/cncf-shinesoft-member-reviewer — Review Shinesoft (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-shipfox-member-reviewer — Review Shipfox (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-shipwright-reviewer — Review Shipwright using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-shopify-contributor-reviewer — Review Shopify (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-shopify-reviewer — Review Shopify using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-shturval-reviewer — Review Shturval using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. The Shturval platform is engineered for complex air-gapped Enterprise environments with the highest security and automation requirements.


/cncf-siddhi-reviewer — Review Siddhi using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-sidekick-reviewer — Review Sidekick using its official documentation and repository in the Observability and Analysis / Observability category. Collect traces, exception stacks and generate logs on-demand without stopping &amp; redeploying your applications.


/cncf-sidero-labs-member-reviewer — Review Sidero Labs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sidero-talos-linux-reviewer — Review Sidero Talos Linux using its official documentation and repository in the Platform / Certified Kubernetes - Installer category. Talos Linux is Linux designed for Kubernetes - secure, immutable, and minimal.


/cncf-sighup-distribution-reviewer — Review SIGHUP Distribution using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. SIGHUP Distribution is a battle-tested distribution purely based on upstream Kubernetes. Deploy and manage a stable and production grade Kubernetes Cluster at scale with a comprehensive Cloud Native stack implemented with top notch CNCF components.


/cncf-siglens-reviewer — Review SigLens using its official documentation and repository in the Observability and Analysis / Observability category. SigLens is a powerful open-source observability solution designed to reduce infrastructure costs by up to 90%. It scales seamlessly from handling 8 TB/day on a single MacBook Air to 1 PB/day across just 32 EC2 instances. With a user-friendly interface and a developer-oriented, pipe-based query language, SigLens simplifies observability data exploration and monitoring at any scale.


/cncf-sigma-reviewer — Review Sigma using its official documentation and repository in the Serverless / Tools category. Sigma is a serverless application developer tool; a cloud IDE, which helps you rapidly build, test and deploy serverless applications.


/cncf-signadot-member-reviewer — Review Signadot (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-signalfx-reviewer — Review SignalFX using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-signoz-member-reviewer — Review SigNoz (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-signserver-community-reviewer — Review SignServer Community using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Sign code, containers, and attestations and create timestamps with the open-source, x509 certificate-based signing software SignServer.


/cncf-sigstore-reviewer — Review sigstore using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-simplyblock-reviewer — Review Simplyblock using its official documentation and repository in the Runtime / Cloud Native Storage category. Simplyblock is a distributed, disaggregated, cloud-native elastic block storage solution for IO-heavy and latency-sensitive container workloads, available in Kubernetes through our CSI StorageClass implementation.


/cncf-singlestore-reviewer — Review SingleStore using its official documentation and repository in the App Definition and Development / Database category.


/cncf-singlestore-wasm-reviewer — Review SingleStore (Wasm) using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-singularity-reviewer — Review Singularity using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-skaffold-reviewer — Review Skaffold using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-skipper-reviewer — Review Skipper using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-skooner-reviewer — Review Skooner using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-sky-betting-and-gaming-contributor-reviewer — Review Sky Betting &amp; Gaming (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-skyline-technology-solutions-supporter-reviewer — Review Skyline Technology Solutions (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-skyloud-kcsp-reviewer — Review Skyloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We build and administer scalable, resilient and secure Kubernetes Infrastructures for developers and end users. Our DevOps &amp; Cloud Engineers stay by your side on a daily basis for the management of your Infrastructure so that your developers remain focused on the development of your applications.


/cncf-skyloud-member-reviewer — Review Skyloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-skywalking-reviewer — Review SkyWalking using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-slime-reviewer — Review Slime using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-slimfaas-reviewer — Review SlimFaaS using its official documentation and repository in the Serverless / Installable Platform category. The slimest and simplest Function As A Service


/cncf-slimtoolkit-reviewer — Review SlimToolkit using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Inspect, Optimize and Debug Your Containers


/cncf-slurm-reviewer — Review Slurm using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Slurm is an open source, fault-tolerant, and highly scalable cluster management and job scheduling system for large and small Linux clusters.


/cncf-smartos-reviewer — Review SmartOS using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-smolagents-reviewer — Review Smolagents using its official documentation and repository in the AI Agent / Agent Framework category. A barebones library for agents that think in code.


/cncf-snapt-nova-reviewer — Review Snapt Nova using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-sncf-supporter-reviewer — Review SNCF (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-snowflake-reviewer — Review Snowflake using its official documentation and repository in the App Definition and Development / Database category.


/cncf-snyk-member-reviewer — Review Snyk (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-snyk-reviewer — Review Snyk using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-socradev-gmbh-kcsp-reviewer — Review socradev Gmbh (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide expert consulting on Kubernetes and the broader Kubernetes ecosystem, helping you implement scalable, secure, and automated infrastructure solutions tailored to your needs.


/cncf-socradev-gmbh-member-reviewer — Review socradev Gmbh (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-soda-foundation-reviewer — Review Soda Foundation using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-sofarpc-reviewer — Review SOFARPC using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category.


/cncf-sofatracer-reviewer — Review SOFATracer using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-softbank-infrinia-ai-cloud-os-reviewer — Review SoftBank Infrinia AI Cloud OS using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Softbank Infrinia AI Cloud OS is a GPU-native managed Kubernetes platform purpose-built for AI and ML workloads, providing automated GPU infrastructure provisioning, Kubernetes cluster orchestration, and MLOps tooling on GPU hardware.


/cncf-softbank-member-reviewer — Review SoftBank (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-software-ag-reviewer — Review Software AG using its official documentation and repository in the App Definition and Development / Database category. With Adabas for LUW, you can deliver extremely high transaction speeds with a fraction of the staff and system resources needed for comparable database management systems.


/cncf-software-mind-kcsp-reviewer — Review Software Mind (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Software Mind services related to Kubernetes include consulting, implementation, upgrades and migration, maintenance, training and security audits. Based on our domain knowledge, we most often provide services to the telecommunications, finance, healthcare, travel, sport betting and e-commerce industries.


/cncf-software-mind-member-reviewer — Review Software Mind (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sogou-c-plus-plus-workflow-reviewer — Review Sogou C++ Workflow using its official documentation and repository in the Platform / PaaS/Container Service category. C++ Parallel Computing and Asynchronous Networking Engine.


/cncf-sokube-kcntp-reviewer — Review SoKube (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. SoKube is pure consulting company that help companies entering the world of Containers &amp; Kubernetes, using a comprehensive SDLC approach from Dev to Production, and using best practices coming from Agile, CI/CD, DevSecOps, SRE, GitOps.


/cncf-sokube-kcsp-reviewer — Review SoKube (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. SoKube is a Swiss based company that helps companies making their digital transition to the world of containers and Kubernetes orchestration, with a strong focus on Agility and DevSecOps lifecycle.


/cncf-sokube-member-reviewer — Review SoKube (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-solanica-member-reviewer — Review Solanica (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-solo-io-member-reviewer — Review Solo.io (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sonatype-member-reviewer — Review Sonatype (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sonatype-reviewer — Review Sonatype using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Sonatype is the developer-friendly full-spectrum software supply, chain management platform helps organizations and software developers. One of the Sonatype products includes the open source software repository manager Sonatype Nexus Repository.


/cncf-sonobuoy-reviewer — Review Sonobuoy using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-sony-interactive-entertainment-supporter-reviewer — Review Sony Interactive Entertainment (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-sops-reviewer — Review SOPS using its official documentation and repository in the Provisioning / Security &amp; Compliance category. sops is an editor of encrypted files that supports YAML, JSON, ENV, INI and BINARY formats and encrypts with AWS KMS, GCP KMS, Azure Key Vault, age, and PGP.


/cncf-sosivio-reviewer — Review Sosivio using its official documentation and repository in the Observability and Analysis / Observability category. Sosivio is a non-intrusive autonomous cloud native platform delivering unmatched visibility into your cloud native environment, eventless failure prediction, and automatic resolution.


/cncf-source-allies-kcsp-reviewer — Review Source Allies (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Source Allies is a premier consultancy and your ally for delivering high-quality technical solutions built around Kubernetes. We help our partners manage their clusters, migrate their applications, and upskill their team in the Cloud Native ecosystem.


/cncf-source-allies-member-reviewer — Review Source Allies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sourcefuse-kcsp-reviewer — Review SourceFuse (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. SourceFuse delivers end-to-end Kubernetes services - from architecture consulting to production-grade managed operations - for enterprises in healthcare, fintech, and SaaS.


/cncf-sourcefuse-member-reviewer — Review SourceFuse (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-southworks-member-reviewer — Review Southworks (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-spacelift-member-reviewer — Review Spacelift (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-spacelift-reviewer — Review Spacelift using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Spacelift is a flexible management platform for Infrastructure as Code. It helps customize your workflows, automate manual tasks, reduce number of errors, improve security and auditability of your infrastructure.


/cncf-sparkfabrik-kcsp-reviewer — Review SparkFabrik (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide cloud native custom development and services on top of Kubernetes, helping our customers to transition to the cloud native era.


/cncf-sparkfabrik-member-reviewer — Review SparkFabrik (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sparta-reviewer — Review Sparta using its official documentation and repository in the Serverless / Framework category.


/cncf-spectro-cloud-kcsp-reviewer — Review Spectro Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Spectro Cloud provides high-touch assistance for enterprises and public sector organizations on their Kubernetes journey, whether in cloud, data center or edge. Our seasoned team guides you from initial project scoping and proof of concept (PoC), through pilot and full production, with knowledge transfer a priority at every stage. Our expert team of support engineers provides round-the-clock technical services covering the full cloud-native stack. We also provide enterprise support for the Kairos CNCF project.


/cncf-spectro-cloud-member-reviewer — Review Spectro Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-spegel-reviewer — Review Spegel using its official documentation and repository in the Provisioning / Container Registry category. Stateless cluster local OCI registry mirror. Spegel enables each node in a Kubernetes cluster to act as a local registry mirror, allowing nodes to share images between themselves. Any image already pulled by a node will be available for any other node in the cluster to pull.


/cncf-spicedb-reviewer — Review SpiceDB using its official documentation and repository in the App Definition and Development / Database category. SpiceDB is an open source database optimized for storing and querying authorization data. SpiceDB sets the standard for scalable authorization by remaining true to the design of the system powering permissions at Google: Zanzibar.


/cncf-spiderlightning-reviewer — Review SpiderLightning using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-spiderpool-reviewer — Review Spiderpool using its official documentation and repository in the Runtime / Cloud Native Network category. Spiderpool is the underlay and RDMA network solution of the Kubernetes, for bare metal, VM and public cloud


/cncf-spiffe-reviewer — Review SPIFFE using its official documentation and repository in the Provisioning / Key Management category.


/cncf-spin-reviewer — Review Spin using its official documentation and repository in the Wasm / Application Frameworks category. Spin is a framework for building and deploying serverless applications in WebAssembly.


/cncf-spinkube-reviewer — Review SpinKube using its official documentation and repository in the Wasm / Orchestration &amp; Management category. Open source platform for efficiently running (containerless) Spin-based WebAssembly (Wasm) applications on Kubernetes.


/cncf-spinnaker-reviewer — Review Spinnaker using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-spire-reviewer — Review SPIRE using its official documentation and repository in the Provisioning / Key Management category.


/cncf-spitzkop-member-reviewer — Review Spitzkop (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-split-reviewer — Review Split using its official documentation and repository in the Observability and Analysis / Feature Flagging category.


/cncf-splunk-reviewer — Review Splunk using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-spot-io-reviewer — Review Spot.io using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-spotify-member-reviewer — Review Spotify (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-spring-cloud-function-reviewer — Review Spring Cloud Function using its official documentation and repository in the Serverless / Framework category.


/cncf-spring-cloud-sleuth-reviewer — Review Spring Cloud Sleuth using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-springer-nature-member-reviewer — Review Springer Nature (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-spyderbat-reviewer — Review Spyderbat using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Spyderbat provides Linux runtime security, protecting dynamic environments by proactively tracking all user activities by their causal connections to detect and resolve external attacks, misconfigurations, and insider threats.


/cncf-squarespace-member-reviewer — Review Squarespace (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-squash-reviewer — Review Squash using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-squer-kcsp-reviewer — Review SQUER (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We support customers in utilizing Kubernetes best for their needs. This includes building truly cloud-native applications and Kubernetes-based developer platforms.


/cncf-squer-member-reviewer — Review SQUER (member) using its official documentation and repository in the CNCF Members / Silver category. We are SQUER: coding architects, software engineers, agile transformation experts and cloud engineers, all striving towards one goal. Solve our customer&#39;s challenges on their way to holistic digitalization while profoundly understanding their business and constraints. Our consultants are developers and vice versa. We believe in tackling challenges hands-on over concept only.


/cncf-srpc-reviewer — Review SRPC using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. High performance, low latency, lightweight enterprise-level RPC system, which supports Baidu bRPC, Tencent tRPC, thrift protocols.


/cncf-sso-reviewer — Review sso using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-sst-reviewer — Review SST using its official documentation and repository in the Serverless / Framework category.


/cncf-stack8s-member-reviewer — Review Stack8s (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stacker-reviewer — Review Stacker using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Stacker is a tool for building OCI images and related artifacts such as SBOMs natively via a declarative yaml format.


/cncf-stackery-reviewer — Review Stackery using its official documentation and repository in the Serverless / Tools category.


/cncf-stackgen-member-reviewer — Review StackGen (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stackgenie-kcsp-reviewer — Review stackgenie (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As a Kubernetes Certified Service Provider, StackGenie offers Kubernetes support, consulting and professional services for businesses and organisations wanting to migrate to cloud technologies, or improve their existing cloud infrastructure.


/cncf-stackgenie-member-reviewer — Review stackgenie (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stackguardian-member-reviewer — Review StackGuardian (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stackhawk-reviewer — Review StackHawk using its official documentation and repository in the Provisioning / Security &amp; Compliance category. StackHawk helps developers find and fix application security vulnerabilities before they hit production with CI/CD automation.


/cncf-stackit-kubernetes-engine-reviewer — Review STACKIT Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. STACKIT Kubernetes Engine (SKE) is a fully managed and scalable Kubernetes service for the deployment and management of Kubernetes  clusters and containerized applications.


/cncf-stackit-member-reviewer — Review STACKIT (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stacklet-member-reviewer — Review Stacklet (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stacklock-minder-reviewer — Review Stacklock Minder using its official documentation and repository in the AI Native Infra / Governance, Policy and Security category. Software Supply Chain Security Platform


/cncf-stacklok-member-reviewer — Review Stacklok (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stackrox-reviewer — Review StackRox using its official documentation and repository in the Provisioning / Security &amp; Compliance category. StackRox is a Kubernetes-native security platform for cloud-native applications, containers, serverless, and Kubernetes.


/cncf-stackstate-reviewer — Review StackState using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-stackstorm-reviewer — Review StackStorm using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-standard-library-reviewer — Review Standard Library using its official documentation and repository in the Serverless / Hosted Platform category.


/cncf-starlingx-reviewer — Review StarlingX using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. StarlingX is a fully integrated open source edge platform that integrates OpenStack and Kubernetes together to provide all functionality needed in the central cloud and on the edge sites. The project is developed by the StarlingX community.


/cncf-starrocks-reviewer — Review StarRocks using its official documentation and repository in the App Definition and Development / Database category. StarRocks is a next-gen sub-second MPP database for full analytics scenarios, including multi-dimensional analytics, real-time analytics and ad-hoc query.


/cncf-stash-by-appscode-reviewer — Review Stash by AppsCode using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-stateful-functions-reviewer — Review Stateful Functions using its official documentation and repository in the Serverless / Framework category.


/cncf-stclab-kcsp-reviewer — Review STCLAB (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Backed by 3+ years of operating high-traffic B2B services on Amazon EKS, we provide Kubernetes professional services — consulting, implementation, 24x7 support, and customer enablement training — that help enterprises reduce downtime, control infrastructure costs, and scale customer experience in production.


/cncf-stclab-member-reviewer — Review STCLAB (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-steadybit-reviewer — Review steadybit using its official documentation and repository in the Observability and Analysis / Chaos Engineering category.


/cncf-steamhaus-kcsp-reviewer — Review Steamhaus (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Steamhaus are AWS cloud native experts, specialising in designing, building and operating containers and serverless platforms. We help startups scale fast and securely, and we help enterprises accelerate and derisk modernisation and transformation.


/cncf-steamhaus-member-reviewer — Review Steamhaus (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stolon-reviewer — Review Stolon using its official documentation and repository in the App Definition and Development / Database category. Cloud Native PostgreSQL High Availability


/cncf-storm-reply-kcntp-reviewer — Review Storm Reply (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. We provide services along the entire cloud value chain and are specialized in the design, implementation and 24X7 operation of innovative cloud-native solutions.


/cncf-storm-reply-kcsp-reviewer — Review Storm Reply (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide services along the entire cloud value chain and are specialized in the design, implementation and 24X7 operation of innovative cloud-native solutions.


/cncf-storm-reply-member-reviewer — Review Storm Reply (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-stormforge-reviewer — Review StormForge using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-storpool-reviewer — Review StorPool using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-strands-reviewer — Review Strands using its official documentation and repository in the AI Agent / Agent Framework category. A model-driven approach to building AI agents in just a few lines of code.


/cncf-strategic-education-member-reviewer — Review Strategic Education (member) using its official documentation and repository in the CNCF Members / Academic category.


/cncf-stratovirt-reviewer — Review StratoVirt using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-strava-member-reviewer — Review Strava (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-strimzi-reviewer — Review Strimzi using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-stunner-reviewer — Review STUNner using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-submariner-reviewer — Review Submariner using its official documentation and repository in the Runtime / Cloud Native Network category. Submariner enables direct networking between Pods and Services in different Kubernetes clusters, either on-premises or in the cloud.


/cncf-sumo-logic-reviewer — Review Sumo Logic using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-sup-info-information-technology-kcsp-reviewer — Review Sup-info Information Technology (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Sup-Info is experienced in Kubernetes-based cloud building，operating，consulting，and training. We also provide professional services for cloud-related open source software.


/cncf-sup-info-information-technology-member-reviewer — Review Sup-info Information Technology (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-supabase-reviewer — Review Supabase using its official documentation and repository in the App Definition and Development / Database category. an open source Firebase alternative


/cncf-superedge-reviewer — Review SuperEdge using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. An edge-native container management system for edge computing


/cncf-superedge-wasm-reviewer — Review SuperEdge (Wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-superplane-reviewer — Review SuperPlane using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category. Open source DevOps control plane for event-driven workflows.


/cncf-surrealdb-member-reviewer — Review SurrealDB (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-suse-kcsp-reviewer — Review SUSE (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. SUSE is here to help simplify your business modernization journey offering reliable and interoperable cloud native solutions with consulting, training and support services to keep up with the pace of change and your expectations.


/cncf-suse-member-reviewer — Review SUSE (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sva-kcntp-reviewer — Review SVA (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. The services of the SVA include consulting, conception, implementation and operation of Kubernetes-based solutions as well as solution-specific and tailor-made trainings and exam preparations.


/cncf-sva-kcsp-reviewer — Review SVA (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. SVA services include consulting, implementing and optimizing Kubernetes platforms on all layers and to share our knowledge with our customers and the community.


/cncf-sva-member-reviewer — Review SVA (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sweet-security-member-reviewer — Review Sweet Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-swift-reviewer — Review Swift using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-swiss-post-contributor-reviewer — Review Swiss Post (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-swisscom-kcsp-reviewer — Review Swisscom (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Swisscom’s experienced cloud native experts support Swiss customers end to end on their journey with the Kubernetes ecosystem.


/cncf-swisscom-kubernetes-service-reviewer — Review Swisscom Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Kubernetes Service is a managed Kubernetes offering in Swisscom’s Private Cloud, based on open-source technology and CNCF tools.


/cncf-swisscom-member-reviewer — Review Swisscom (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-switch-cloud-kubernetes-sck-reviewer — Review Switch Cloud Kubernetes (SCK) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Switch Cloud Kubernetes (SCK) is a Kubernetes as a Service platform for the education, research and innovation landscape in Switzerland.


/cncf-switch-member-reviewer — Review SWITCH (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-syft-reviewer — Review Syft using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Syft is a CLI tool and library for generating a Software Bill of Materials from container images and filesystems


/cncf-symbotic-member-reviewer — Review Symbotic (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-synadia-member-reviewer — Review Synadia (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-synax-member-reviewer — Review Synax (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-syntasso-member-reviewer — Review Syntasso (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-synx-data-labs-member-reviewer — Review Synx Data Labs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sysbox-reviewer — Review Sysbox using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-sysdig-kcsp-reviewer — Review Sysdig (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Sysdig offers professional services, technical account management, support and training. These services are designed to help companies get up-and-running with Sysdig and find best practices with their Kubernetes platform and cloud native applications.


/cncf-sysdig-member-reviewer — Review Sysdig (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-sysdig-reviewer — Review sysdig using its official documentation and repository in the Observability and Analysis / Observability category. sysdig is a simple tool with deep system visibility for exploration and troubleshooting, with native support for containers.


/cncf-sysdig-secure-reviewer — Review Sysdig Secure using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Sysdig Secure embeds security and compliance into the build, run and respond stages of the Kubernetes lifecycle.


/cncf-syseleven-kcsp-reviewer — Review SysEleven (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We support our customers from start to finish, from training to the full operation of Kubernetes setups. We give regular workshops and share our know-how.


/cncf-syseleven-member-reviewer — Review SysEleven (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-syseleven-metakube-reviewer — Review SysEleven MetaKube using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. MetaKube offers carefree Multi Cloud Kubernetes as a Service on OpenStack Cloud, Amazon Web Services (AWS) and On-Premises. We provide you with SaaS solutions like a load balancer as well as solutions for backup and recovery, monitoring, etc., and various types of support - detailed documentation, tutorials and personal SRE-Support.


/cncf-t-mobile-member-reviewer — Review T-Mobile (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-t-systems-kcsp-reviewer — Review T-Systems (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Our team of CNCF certified Kubernetes experts have extensive experience working with a wide range of industries and can help you design and implement a solution that meets your specific needs.


/cncf-t-systems-member-reviewer — Review T-Systems (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-taikun-cloudworks-wasm-reviewer — Review Taikun CloudWorks (Wasm) using its official documentation and repository in the Wasm / Hosted Platforms category. Taikun is a hybrid cloud management orchestration platform.


/cncf-tailscale-member-reviewer — Review Tailscale (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-taiwan-ai-cloud-member-reviewer — Review Taiwan AI Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-talend-data-streams-reviewer — Review Talend Data Streams using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category.


/cncf-talos-linux-reviewer — Review Talos Linux using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Talos Linux is a single purpose distribution for Kubernetes. It provides an API driven—no SSH—operating system to deploy vanilla Kubernetes anywhere.


/cncf-tanka-reviewer — Review Tanka using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Tanka is a composable configuration utility for Kubernetes. It leverages the Jsonnet language to realize flexible, reusable and concise configuration.


/cncf-tarantool-reviewer — Review Tarantool using its official documentation and repository in the App Definition and Development / Database category. In-memory computing platform consisting of a database and an application server


/cncf-tarook-reviewer — Review Tarook using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Holistic life cycle management of Kubernetes clusters on bare metal or OpenStack


/cncf-tars-reviewer — Review TARS using its official documentation and repository in the Orchestration &amp; Management / Remote Procedure Call category. TARS is a high-performance RPC framework that provides an integrated solution of microservice governance.


/cncf-tata-communications-member-reviewer — Review Tata Communications (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tata-communications-vayu-kubernetes-service-reviewer — Review Tata Communications Vayu Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Tata Communications Vayu Kubernetes Service is an enterprise-grade managed service that makes it easy for you to run Kubernetes on Tata Communications  Vayu Cloud platform without needing to install and operate your own Kubernetes clusters, for running containerized applications, including stateful  and stateless, AI and ML, complex and simple web apps, API, and backend services. Service is available only India for consumption.


/cncf-taubyte-reviewer — Review Taubyte using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-tcc-consulting-limited-kcsp-reviewer — Review TCC Consulting Limited (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. TCC Consulting specializes in enterprise Kubernetes solutions, offering end-to-end architecture design, deployment optimization, and lifecycle management for containerized environments. Our certified experts ensure seamless integration with DevOps pipelines and 24/7 operational support, empowering businesses to scale securely and innovate with confidence.


/cncf-tcc-consulting-limited-member-reviewer — Review TCC Consulting Limited (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tdengine-reviewer — Review TDengine using its official documentation and repository in the App Definition and Development / Database category.


/cncf-teamcity-reviewer — Review TeamCity using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Powerful Continuous Integration out of the box


/cncf-teciem-member-reviewer — Review Teciem (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tekton-reviewer — Review Tekton using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. A powerful and flexible open source framework for creating continuous integration and delivery (CI/CD) systems that allow developers to build, test, and deploy across multiple cloud providers and on-premises systems by abstracting away the underlying implementation details.


/cncf-telemetryhub-by-scout-apm-reviewer — Review TelemetryHub by Scout APM using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-telenor-norway-member-reviewer — Review Telenor Norway (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-teleport-member-reviewer — Review Teleport (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-teleport-reviewer — Review Teleport using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-telepresence-reviewer — Review Telepresence using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Local development against a remote Kubernetes or OpenShift cluster


/cncf-teller-reviewer — Review Teller using its official documentation and repository in the Provisioning / Key Management category.


/cncf-temporal-member-reviewer — Review Temporal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-temporal-reviewer — Review Temporal using its official documentation and repository in the AI Agent / Workflow Orchestration category. Temporal service


/cncf-tencent-cloud-kcsp-reviewer — Review Tencent Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Tencent Cloud offers high performing technology with the right amount of human touch to help you build and scale the solution optimized for your business. Tencent Kubernetes Engine helps customers quickly build a native Kubernetes architecture in Tencent Cloud.


/cncf-tencent-cloud-log-service-reviewer — Review Tencent Cloud Log Service using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-tencent-cloud-member-reviewer — Review Tencent Cloud (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-tencent-cloud-serverless-cloud-function-reviewer — Review Tencent Cloud Serverless Cloud Function using its official documentation and repository in the Serverless / Hosted Platform category. Serverless Cloud Function (SCF) is a serverless execution environment provided by Tencent Cloud for enterprises and developers to help you run your code without purchasing and managing servers.


/cncf-tencent-kubernetes-engine-tke-reviewer — Review Tencent Kubernetes Engine (TKE) using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Tencent Kubernetes Engine (TKE) provides container-centric, highly scalable and high-performance container management services. Fully compatible with Kubernetes&#39; native API and capable of expanding Tencent Cloud&#39;s Kubernetes plugins such as CBS and CLB.


/cncf-tengine-reviewer — Review Tengine using its official documentation and repository in the Orchestration &amp; Management / Service Proxy category.


/cncf-tensor-security-reviewer — Review Tensor Security using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Tensor Security is a cloud-native security company focused on providing AI-based and full-stack cloud-native security solutions.


/cncf-tensor9-member-reviewer — Review Tensor9 (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tensorflow-distributed-reviewer — Review Tensorflow Distributed using its official documentation and repository in the Training / Distributed Training category. TensorFlow API to distribute training across multiple GPUs, multiple machines, or TPUs. Using this API, you can distribute your existing models and training code with minimal code changes.


/cncf-tensorflow-extended-tfx-reviewer — Review TensorFlow Extended (TFX) using its official documentation and repository in the AI Native Infra / Continuous Integration and Delivery category. TFX is an end-to-end platform for deploying production ML pipelines


/cncf-tensorflow-reviewer — Review TensorFlow using its official documentation and repository in the Data / Data Science category. An Open Source Machine Learning Framework for Everyone


/cncf-tensormesh-member-reviewer — Review Tensormesh (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tenxcloud-kcsp-reviewer — Review TenxCloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. TenxCloud is a professional cloud native application and data platform service provider. Adhering to the mission of &quot;making computing generate value and making data an asset&quot;, it provides complete digital products and services for medium and large enterprise users, and works together to promote the construction of Digital China. As a Kubernetes Certified Service Provider (KCSP), TenxCloud has internationally recognized qualifications for Kubernetes service provision, consulting, and training.


/cncf-tenxcloud-member-reviewer — Review TenxCloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-terasky-kcntp-reviewer — Review TeraSky (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. At TeraSky, our core expertise lies in delivering end-to-end Kubernetes and Cloud Native solutions, seamlessly combining hands-on implementation with customized  training. In addition to Kubernetes, we offer structured training programs on Terraform, HashiCorp Vault, and VMware, enabling our customers to not only adopt  cutting-edge technologies but also master them for long-term success and sustainability. By integrating the Linux Foundation’s official training portfolio,  we’re expanding our capabilities with industry-recognized certifications that complement our existing expertise. This holistic approach empowers organizations through a complete learning journey—from design and implementation to certified training—positioning TeraSky as a one-stop partner for enterprise Cloud Native transformation.


/cncf-terasky-kcsp-reviewer — Review TeraSky (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. TeraSky is a trusted Kubernetes Certified Service Provider (KCSP) with extensive hands-on experience in Kubernetes and Cloud-Native technologies. Our team has successfully led enterprise-scale implementations, helping organizations adopt, scale, and optimize Kubernetes-based platforms. As an active contributor to the CNCF ecosystem, led by our Principal Architect, Scott Rosenberg —a recognized community leader—we combine deep technical mastery with ongoing community engagement to deliver exceptional enablement, training, and innovation.


/cncf-terasky-member-reviewer — Review TeraSky (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-terracotta-ai-member-reviewer — Review Terracotta AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-terraform-reviewer — Review Terraform using its official documentation and repository in the Provisioning / Automation &amp; Configuration category.


/cncf-terramate-member-reviewer — Review Terramate (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-terramate-reviewer — Review Terramate using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Terramate is an Infrastructure as Code (IaC) Management Platform combining Infrastructure Delivery, Drift Management, Observability, and Collaboration in a single platform that developers are excited about. We focus on providing the best developer experience by combining workflows, best practices, and more so that developers can move beyond infrastructure and focus on developing applications.


/cncf-terranetes-reviewer — Review Terranetes using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. The Appvia Terranetes Controller manages the life-cycle of Terraform and OpenTofu resources defined and built inside Kubernetes. This allows teams running workloads inside the cluster to self-serve application dependencies and reuse the wealth of Terraform modules already written.


/cncf-terrascan-reviewer — Review Terrascan using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-tesseral-reviewer — Review Tesseral using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Tesseral is the open source platform for managing identity and access in business software.


/cncf-testkube-reviewer — Review Testkube using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Testkube provides a Kubernetes-native framework for test definition, execution and results. It decouples test artifacts and execution from CI/CD tooling and makes testing part of your cluster&#39;s state. Testkube is built and maintained by Kubeshop.


/cncf-tetragon-reviewer — Review Tetragon using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-tetrate-member-reviewer — Review Tetrate (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tetrate-service-bridge-tsb-reviewer — Review Tetrate Service Bridge (TSB) using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category. Application connectivity platform.


/cncf-teuto-net-kcsp-reviewer — Review teuto.net (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Kubernetes made in Germany - Setup - Hosting - highly skilled DevOps team


/cncf-teuto-net-managed-kubernetes-reviewer — Review teuto.net Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. teuto.net quickly deploys and manages production-ready Kubernetes clusters for you.


/cncf-teuto-net-member-reviewer — Review teuto.net (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-text-generation-inference-tgi-reviewer — Review text-generation-inference (TGI) using its official documentation and repository in the Inference / Runtime category. Large Language Model Text Generation Inference


/cncf-thales-member-reviewer — Review Thales (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-thanos-reviewer — Review Thanos using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-the-new-york-times-supporter-reviewer — Review The New York Times (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-the-scale-factory-kcsp-reviewer — Review The Scale Factory (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We help technology teams deliver great cloud native solutions on AWS.


/cncf-the-scale-factory-member-reviewer — Review The Scale Factory (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-the-update-framework-tuf-reviewer — Review The Update Framework (TUF) using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-thermo-fisher-scientific-supporter-reviewer — Review Thermo Fisher Scientific (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-thoras-ai-member-reviewer — Review Thoras.ai (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-threat-stack-reviewer — Review Threat Stack using its official documentation and repository in the Serverless / Security category.


/cncf-threatmapper-reviewer — Review ThreatMapper using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-thundra-reviewer — Review Thundra using its official documentation and repository in the Serverless / Tools category. Thundra helps answer your toughest questions about the health of your Serverless infrastructure by giving you invaluable insights into your AWS Lambda functions.


/cncf-ticketmaster-supporter-reviewer — Review TicketMaster (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-tidb-reviewer — Review TiDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-tigera-kcsp-reviewer — Review Tigera (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Calico helps solve complex use cases, accelerate your adoption process of microservices and bring your team up to speed in operating and supporting your application modernization initiative.


/cncf-tigera-member-reviewer — Review Tigera (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tigera-reviewer — Review Tigera using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-tigris-data-member-reviewer — Review Tigris Data (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tikv-reviewer — Review TiKV using its official documentation and repository in the App Definition and Development / Database category. A distributed transactional key-value database. Based on the design of Google Spanner and HBase, but simpler to manage and without dependencies on any distributed filesystem


/cncf-tilt-reviewer — Review Tilt using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category.


/cncf-timescale-reviewer — Review Timescale using its official documentation and repository in the App Definition and Development / Database category. PostgreSQL for time‑series. TimescaleDB is the leading open-source relational database for time-series data. Fully managed or self‑hosted.


/cncf-tingyun-reviewer — Review Tingyun using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-tinkerbell-reviewer — Review Tinkerbell using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Bare metal provisioning engine, supporting network and ISO booting, BMC interactions, metadata service, and workflow engine.


/cncf-tintri-member-reviewer — Review Tintri (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tinygo-reviewer — Review TinyGo using its official documentation and repository in the Wasm / Languages category. Compiled language to Wasm


/cncf-tno-member-reviewer — Review TNO (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-tokenetes-reviewer — Review Tokenetes using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Tokenetes implements Transaction Tokens (TraTs) for microservices call chains.  It&#39;s a Kubernetes-native framework providing immutable identity and context in  service-to-service communication to prevent attacks like software supply chain  or privileged user compromise.


/cncf-tokyo-gas-supporter-reviewer — Review Tokyo Gas (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-tomtom-supporter-reviewer — Review TomTom (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-tongtech-kcsp-reviewer — Review TongTech (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Beijing Tongtech Software Co., Ltd. helps all kinds of enterprises to maximize the effectiveness of Kubernetes and meet the needs of enterprise business development.


/cncf-tongtech-member-reviewer — Review TongTech (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-topaz-reviewer — Review Topaz using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Cloud-native authorization for modern applications and APIs, combining the best of Open Policy Agent and Google Zanzibar


/cncf-torchx-reviewer — Review Torchx using its official documentation and repository in the Training / Distributed Training category. Universal job launcher for PyTorch applications. TorchX is designed to have fast iteration time for training/research and support for E2E production ML pipelines when you&#39;re ready.


/cncf-toyota-motor-corporation-member-reviewer — Review Toyota Motor Corporation (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-traas-bos-reviewer — Review TRaaS BOS using its official documentation and repository in the Observability and Analysis / Observability category. Provide multi-dimensional integrated monitoring capabilities for business, application, network, and other aspects to measure business availability, gain insight into service health status, and support real-time log collection and analysis to efficiently locate problems. Provide non-invasive automatic link generation capabilities for heterogeneous applications on and off the cloud, and support defining services to aggregate massive link data into business links, achieving customized association between services and applications and intelligent diagnosis and positioning.


/cncf-traas-has-reviewer — Review TRaaS HAS using its official documentation and repository in the Observability and Analysis / Continuous Optimization category. Provides the ability to aggregate dispersed fault data into risk events, as well as emergency response, fault localization (root cause analysis, diagnostic decision trees), and fault handling (emergency self-healing, multi-active disaster recovery), to achieve a full set of emergency technology risk prevention and control capabilities in 1 minute discovery - 5 minutes localization - 10 minutes recovery.


/cncf-tracee-reviewer — Review Tracee using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Tracee helps you understand your systems and applications behavior, and apply runtime security to them. It uses Linux eBPF technology and strong security techniques to instrument the system and expose relevant events for you to consume and process.


/cncf-tracetest-reviewer — Review Tracetest using its official documentation and repository in the Observability and Analysis / Observability category. Tracetest allows easy creation of end-to-end tests via a simple user interface. It leverages your current investment in OpenTelemetry based tracing to make deep integration and E2E testing simple.


/cncf-traefik-labs-member-reviewer — Review Traefik Labs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-traefik-mesh-reviewer — Review Traefik Mesh using its official documentation and repository in the Orchestration &amp; Management / Service Mesh category.


/cncf-traefik-reviewer — Review Traefik using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-travelping-member-reviewer — Review Travelping (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-traversal-member-reviewer — Review Traversal (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-travis-ci-reviewer — Review Travis CI using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-tremor-reviewer — Review Tremor using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. An early-stage event processing system for unstructured data with rich support for structural pattern-matching, filtering and transformation


/cncf-trend-micro-cloud-one-reviewer — Review Trend Micro Cloud One using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-trend-micro-member-reviewer — Review Trend Micro (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-trickster-reviewer — Review Trickster using its official documentation and repository in the Observability and Analysis / Observability category. Open Source HTTP Reverse Proxy Cache and Time Series Dashboard Accelerator


/cncf-trilio-reviewer — Review Trilio using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-trink-io-reviewer — Review Trink.io using its official documentation and repository in the Observability and Analysis / Observability category. Trink.io is the first log shipping platform, built to replace fragmented solutions and workarounds with a single solution that takes care of the basics while offering visibility and control that only a managed solution can. Trink handles the entire operational data collection scope for you, connecting to all of your sources (e.g AWS) and enabling you to choose your preferred destinations (Datadog, Logz, etc.) in minutes instead of days.


/cncf-trino-reviewer — Review Trino using its official documentation and repository in the Data / Data Architecture category. The distributed SQL query engine for big data, formerly known as PrestoSQL.


/cncf-triton-object-storage-reviewer — Review Triton Object Storage using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-triton-reviewer — Review Triton using its official documentation and repository in the AI Native Infra / Accelerator and SuperPod category. Triton is a language and compiler for parallel programming.


/cncf-trivago-supporter-reviewer — Review trivago (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-trivy-reviewer — Review Trivy using its official documentation and repository in the Provisioning / Security &amp; Compliance category. Trivy is an all-in-one security scanner that beings vulnerability, SBOM, misconfigurations, leaked secrets, and license scanning together, across container images, code repositories, Kubernetes clusters and more.


/cncf-truefoundry-member-reviewer — Review TrueFoundry (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-truefullstaq-kcntp-reviewer — Review TrueFullstaq (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. As an internet and cloud native pioneer, we are the leading expert in cloud native technology. We offer hands-on Kubernetes trainings that prepares you for official certification. For Developers, we’ve put together a practical 3 days course for those who want to run applications on Kubernetes. You’ll learn deployment strategies, configuration, Pods, Services, Ingress, and more. Fully aligned with the CKAD certification. Expect lots of hands-on labs and best practices. Designed for experienced Kubernetes users who want to go deeper, we offer a Kubernetes for Administrators course (4 days). This training covers advanced topics required for the CKA certification, with a strong focus on safe, scalable operations through hands-on workshops.


/cncf-truefullstaq-kcsp-reviewer — Review TrueFullstaq (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As an internet and cloud native pioneer, we are the leading expert in cloud native technology. We offer customized solutions for your cloud infrastructure, including comprehensive Kubernetes services and container orchestration. We advise you on how to run your application effortlessly in the cloud, manage your Kubernetes environments, and keep your cloud costs under control. As the first certified Kubernetes organization in the Netherlands with multiple Kubestronauts on our team, we bring deep expertise in both managed Kubernetes and Azure Kubernetes Service (AKS). With our focus on innovation and reliability, we are your partner in getting the most out of the cloud!


/cncf-truefullstaq-member-reviewer — Review TrueFullstaq (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-trulens-37878c25-reviewer — Review TruLens using its official documentation and repository in the AI Agent / Evaluation category. Evaluation and Tracking for LLM Experiments and AI Agents


/cncf-trulens-743c42bc-reviewer — Review Trulens using its official documentation and repository in the AI Native Infra / Observability category. TruLens is a powerful open source library for evaluating and tracking large language model-based applications.


/cncf-tsuru-reviewer — Review Tsuru using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-twilio-functions-reviewer — Review Twilio Functions using its official documentation and repository in the Serverless / Hosted Platform category. A serverless environment to build and run your Twilio code so you can get to production faster.


/cncf-txtai-reviewer — Review txtai using its official documentation and repository in the AI Agent / RAG category. All-in-one AI framework for semantic search, LLM orchestration and language model workflows.


/cncf-tyk-member-reviewer — Review Tyk (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-tyk-reviewer — Review Tyk using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-typhoon-reviewer — Review Typhoon using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Typhoon distributes upstream Kubernetes, architectural conventions, and cluster addons, much like a GNU/Linux distribution provides the Linux kernel and userspace components.


/cncf-uber-member-reviewer — Review Uber (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ubs-member-reviewer — Review UBS (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ultimate-kronos-group-ukg-supporter-reviewer — Review Ultimate Kronos Group (UKG) (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-umb-kcsp-reviewer — Review UMB (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. UMB is offering managed OpenShift and Kubernetes services in public clouds or our own hosted cloud in Switzerland. We support our customers on their journey to cloud native applications from consulting to implementation.


/cncf-umb-member-reviewer — Review UMB (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-unicom-cloud-kcsp-reviewer — Review Unicom Cloud (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. China Unicom provides Kubernetes professional services based on its enterprise Kubernetes platform (CSK). Our services include Kubernetes cluster planning and architecture design, cluster deployment and lifecycle management, Kubernetes-based application platform setup, system integration, and operational best practices.   We also provide ongoing technical support and consulting services to help customers operate, maintain, and optimize their Kubernetes environments in production.


/cncf-unikraft-member-reviewer — Review Unikraft (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-unikraft-reviewer — Review Unikraft using its official documentation and repository in the Wasm / Edge/Bare metal category.


/cncf-union-ai-member-reviewer — Review Union.ai (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-unleash-reviewer — Review Unleash using its official documentation and repository in the Observability and Analysis / Feature Flagging category. Open-source feature management platform


/cncf-unryo-reviewer — Review unryo using its official documentation and repository in the Observability and Analysis / Observability category. Observability and cross-layer topology correlation on top of CNCF ecosystem: Prometheus, Influx, Kubernetes, Elastic, networks and 5G.


/cncf-unsloth-reviewer — Review Unsloth using its official documentation and repository in the Training / Post Training category. Unsloth Studio is a web UI for training and running open models like Gemma 4, Qwen3.5, DeepSeek, gpt-oss locally.


/cncf-unstructured-reviewer — Review Unstructured using its official documentation and repository in the AI Agent / RAG category. Convert documents to structured data effortlessly. Unstructured is open-source ETL solution for transforming complex documents into clean, structured formats for language models. Visit our website to learn more about our enterprise grade Platform product for production grade workflows, partitioning, enrichments, chunking and embedding.


/cncf-upbound-member-reviewer — Review Upbound (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-upcloud-managed-kubernetes-reviewer — Review UpCloud Managed Kubernetes using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. UpCloud Managed Kubernetes automates the deployment, scaling and management of container workloads.


/cncf-upcloud-member-reviewer — Review UpCloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-updatecli-reviewer — Review Updatecli using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Updatecli is a declarative dependency management tool for Git repositories.


/cncf-upsider-member-reviewer — Review Upsider (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-upwind-security-member-reviewer — Review Upwind Security (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-uqbar-reviewer — Review Uqbar using its official documentation and repository in the App Definition and Development / Database category.


/cncf-urunc-reviewer — Review urunc using its official documentation and repository in the Runtime / Container Runtime category. A CRI-compatible runtime for running unikernels and application kernels as containers.  urunc bridges the gap between unikernels and containerized environments, enabling seamless  integration with cloud-native architectures while maintaining OCI compatibility.


/cncf-userver-reviewer — Review userver using its official documentation and repository in the App Definition and Development / Streaming &amp; Messaging category. userver is a modern open source asynchronous framework with a rich set of abstractions and drivers for fast and comfortable creation of resilient C++ microservices, services and utilities.


/cncf-v8-reviewer — Review V8 using its official documentation and repository in the Wasm / Runtimes category.


/cncf-va-linux-systems-japan-kcsp-reviewer — Review VA Linux Systems Japan (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. VA Linux provides the following services for Kubernetes: Consulting, Development, Failure Analysis &amp; Support.


/cncf-va-linux-systems-japan-member-reviewer — Review VA Linux Systems Japan (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vald-reviewer — Review Vald using its official documentation and repository in the App Definition and Development / Database category. Vald is a highly scalable, cloud-native distributed vector search engine optimized for machine learning and AI applications. It offers efficient, near real-time search for high-dimensional vector data, ensuring robust performance and flexibility in handling large datasets. Engineered for ease of use and integration, Vald empowers developers with cutting-edge search capabilities.


/cncf-valkey-reviewer — Review Valkey using its official documentation and repository in the App Definition and Development / Database category. A flexible distributed key-value datastore that is optimized for caching and other realtime workloads.


/cncf-valve-software-member-reviewer — Review Valve Software (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vamp-io-reviewer — Review vamp.io using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-varmor-reviewer — Review vArmor using its official documentation and repository in the Provisioning / Security &amp; Compliance category. vArmor is a cloud-native container hardening system for Kubernetes. It leverages AppArmor, BPF LSM, Seccomp, and an Envoy-based NetworkProxy sidecar to enforce access control from system calls to application protocols. Its multi-enforcer architecture and built-in rules provide out-of-the-box protection for containers and AI Agent workloads without requiring deep security profile expertise.


/cncf-vault-reviewer — Review Vault using its official documentation and repository in the Provisioning / Key Management category.


/cncf-vcluster-member-reviewer — Review vCluster (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vcluster-reviewer — Review vCluster using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. vCluster lets you create fully functional but virtual Kubernetes clusters. Each vCluster runs inside a namespace of another Kubernetes cluster. Using vCluster is much cheaper than creating separate full-blown clusters and it offers better multi-tenancy and isolation compared to regular namespaces.


/cncf-vector-reviewer — Review Vector using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-veinmind-tools-reviewer — Review Veinmind Tools using its official documentation and repository in the Provisioning / Security &amp; Compliance category. veinmind-tools is self-developed by chaitin technology, cloudwalker team incubation, a container security toolset based on veinmind-sdk


/cncf-velero-reviewer — Review Velero using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-velocity-member-reviewer — Review Velocity (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vercel-reviewer — Review Vercel using its official documentation and repository in the Serverless / Hosted Platform category.


/cncf-verda-member-reviewer — Review Verda (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-veritas-automata-member-reviewer — Review Veritas Automata (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-verl-reviewer — Review verl using its official documentation and repository in the Training / Post Training category. Volcano Engine Reinforcement Learning for LLMs


/cncf-vertica-reviewer — Review Vertica using its official documentation and repository in the App Definition and Development / Database category. The core analytical platform within the OpenText software portfolio, Vertica is the Unified Analytics Platform, based on a massively scalable architecture with the broadest set of analytical capabilities and end-to-end in-database machine learning.


/cncf-vexxhost-member-reviewer — Review VEXXHOST (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vhl-technologies-member-reviewer — Review VHL Technologies (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-victoriametrics-member-reviewer — Review VictoriaMetrics (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-victoriametrics-reviewer — Review VictoriaMetrics using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-viettel-ai-platform-reviewer — Review Viettel AI Platform using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Viettel AI Platform is a container management and orchestration platform based on Cloud Native architecture, designed to support AI/ML applications running in Kubernetes environments.


/cncf-viettel-kubernetes-engine-reviewer — Review Viettel Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Viettel Kubernetes Engine is a fully managed Kubernetes service that enables you to deploy and operate Kubernetes clusters on the Viettel Cloud platform.


/cncf-viettel-member-reviewer — Review Viettel (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-vijil-member-reviewer — Review vijil (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-villagesql-member-reviewer — Review VillageSQL (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vineyard-reviewer — Review Vineyard using its official documentation and repository in the Runtime / Cloud Native Storage category. Vineyard (v6d) is an in-memory immutable data manager.


/cncf-virtasant-reviewer — Review Virtasant using its official documentation and repository in the Observability and Analysis / Observability category. Virtasant’s Cloud Optimization solution includes program setup, global cloud expertise, and a proprietary technology platform that offers complete end-to-end automation to manage Cloud FinOps programs and optimize cloud resources.


/cncf-virtual-kubelet-reviewer — Review Virtual Kubelet using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-visual-studio-code-kubernetes-tools-reviewer — Review Visual Studio Code Kubernetes Tools using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. The extension for developers building applications to run in Kubernetes clusters and for DevOps staff troubleshooting Kubernetes applications.


/cncf-vitess-reviewer — Review Vitess using its official documentation and repository in the App Definition and Development / Database category. MySQL-compatible, horizontally scalable, cloud-native database solution.


/cncf-vjsar-kcsp-reviewer — Review VJSAR (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We deliver expert Kubernetes services, including automated cluster management, zero-downtime deployments, cloud-native migrations, and integrated security, ensuring your infrastructure remains scalable, resilient, and future-ready.


/cncf-vjsar-member-reviewer — Review VJSAR (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vllm-reviewer — Review vLLM using its official documentation and repository in the Inference / Runtime category. vLLM is a fast and easy-to-use library for LLM inference and serving.


/cncf-vmware-application-catalog-reviewer — Review VMware Application Catalog using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. VMware Application Catalog is a customizable selection of trusted, pre-packaged open-source application components that are continuously maintained and verifiably tested for use in production environments. It is the enterprise version of Bitnami Application Catalog.


/cncf-vmware-aria-operations-for-applications-reviewer — Review VMware Aria Operations for Applications using its official documentation and repository in the Observability and Analysis / Observability category. VMware Aria Operations for Applications brings together metrics, traces, and logs in one solution for greater business insights and scalability


/cncf-vmware-kcsp-reviewer — Review VMware (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. VMware simplifies Kubernetes for platform teams and IT to deploy, manage, operate and secure their Kubernetes environments. VMware offers professional services to accelerate time to value for organizations of all sizes.


/cncf-vmware-nsx-reviewer — Review VMware NSX using its official documentation and repository in the Runtime / Cloud Native Network category. VMware NSX® Data Center delivers virtualized networking and security entirely in software, completing a key pillar of the Software-defined Data Center (SDDC), and enabling the virtual cloud network to connect and protect across data centers, clouds, and applications.


/cncf-vmware-tanzu-cloudhealth-reviewer — Review VMware Tanzu CloudHealth using its official documentation and repository in the Observability and Analysis / Observability category. Tanzu CloudHealth is a multi-cloud FinOps platform that helps organizations make sense of their cloud data, optimize and control cloud spend, and enhance their cloud management practice.


/cncf-vmware-tanzu-kubernetes-grid-reviewer — Review VMware Tanzu Kubernetes Grid using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. VMware Tanzu Kubernetes Grid is VMware&#39;s Kubernetes distribution - built on open source technologies, packaged for enterprise adoption and supported 24x7 by VMware Global Support Services (GSS).


/cncf-vmware-vsphere-kubernetes-service-reviewer — Review VMware vSphere Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. VMware vSphere Kubernetes Service is VMware&#39;s Kubernetes distribution - built on open source technologies, packaged for enterprise adoption and supported 24x7 by VMware Global Support Services (GSS).


/cncf-vmware-vsphere-reviewer — Review VMware vSphere using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. vSphere is the industry-leading server virtualization software and the heart of a modern SDDC, helping you run, manage, connect, and secure your applications in a common operating environment across clouds.


/cncf-volcano-engine-apmplus-reviewer — Review Volcano Engine APMPlus using its official documentation and repository in the Observability and Analysis / Observability category. APMPlus (Application Performance Management Plus) on Volcano Engine delivers end-to-end observability for modern cloud-native environments, enabling precise root cause analysis and performance tuning through unified metrics, distributed tracing, and log correlation. It provides full-stack monitoring from client devices to server-side infrastructure, with specialized support for AI-native workloads such as large language models (LLMs), ensuring visibility across hybrid architectures and alignment with cloud-native best practices.


/cncf-volcano-engine-kcntp-reviewer — Review Volcano Engine (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. Volcano Engine (available in mainland China) provides a coherent and systematic framework from the infrastructure OS level, to mid-platform, and end-to-end suites for system orchestration, data integration with various industrialized solutions to help enterprises achieve scalability and business growth in a cloud-native way. Volcano Engine also provides consulting and training services regarding various cloud-native technological integration, application and hands-on knowledge.


/cncf-volcano-engine-kcsp-reviewer — Review Volcano Engine (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Volcano Engine (available in mainland China) provides a coherent and systematic framework that works from the infrastructure OS level, to mid-platform, and end-to-end suites for system orchestration, data integration and analysis that can also be paired with various industrialized solutions to help enterprises achieve scalability and business growth in a cloud-native way.


/cncf-volcano-engine-member-reviewer — Review Volcano Engine (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-volcano-engine-vmp-reviewer — Review Volcano Engine VMP using its official documentation and repository in the Observability and Analysis / Observability category. Volcengine Managed Service for Prometheus (VMP) is a fully managed, maintenance-free cloud-native monitoring system that inherits and extends the open-source Prometheus ecosystem. The service delivers out-of-the-box observability data collection, alerting and actionable insights, enabling users to quickly establish end-to-end monitoring capabilities across Kubernetes, AI services, and multi-cloud environments.


/cncf-volcano-kthena-reviewer — Review Volcano Kthena using its official documentation and repository in the Inference / Framework category.


/cncf-volcano-reviewer — Review Volcano using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-volcengine-kubernetes-engine-reviewer — Review Volcengine Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Volcengine Kubernetes Engine(VKE) is a Kubernetes-based service that provides high-performance container cluster management. By deeply integrating the new generation of cloud-native technologies, VKE ensures high efficiency for enterprises through running containerized applications on the cloud


/cncf-voltdb-reviewer — Review VoltDB using its official documentation and repository in the App Definition and Development / Database category. VoltDB is a high-velocity decisioning engine, powering real-time applications that must react in milliseconds.


/cncf-volvo-contributor-reviewer — Review Volvo (contributor) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-vscode-wasm-extension-reviewer — Review VSCode Wasm Extension using its official documentation and repository in the Wasm / Tooling category.


/cncf-vscode-wasm-reviewer — Review vscode-wasm using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-vshn-ag-kcsp-reviewer — Review VSHN AG (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. VSHN - The DevOps Company transforms software into online services by automating and managing application workloads. Let us handle the Ops so you can build the future. We provide 24/7 platform engineering and Kubernetes support for software developers, SaaS, agencies, and ISVs.


/cncf-vshn-ag-member-reviewer — Review VSHN AG (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-vultr-kubernetes-engine-reviewer — Review Vultr Kubernetes Engine using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. Vultr Kubernetes Engine gives developers a simple way to deploy containerized workloads with predictable pricing on a fully-managed Kubernetes cluster that integrates with other Vultr managed services such as Load Balancers, Block Storage, and DNS.


/cncf-vultr-member-reviewer — Review Vultr (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wa-lang-reviewer — Review Wa-lang using its official documentation and repository in the Wasm / Languages category. Design for WebAssembly


/cncf-wabt-reviewer — Review Wabt using its official documentation and repository in the Wasm / Tooling category.


/cncf-wai-reviewer — Review WAI using its official documentation and repository in the Wasm / Tooling category.


/cncf-walmart-member-reviewer — Review Walmart (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wamr-reviewer — Review WAMR using its official documentation and repository in the Wasm / Runtimes category.


/cncf-wand-cloud-member-reviewer — Review Wand Cloud (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-warg-reviewer — Review warg using its official documentation and repository in the Wasm / Packaging, Registries &amp; Application Delivery category.


/cncf-wasi-logging-reviewer — Review wasi-logging using its official documentation and repository in the Wasm / Debugging &amp; Observability category.


/cncf-wasi-nn-for-openvino-reviewer — Review WASI NN for OpenVINO using its official documentation and repository in the Wasm / AI/Machine Learning category.


/cncf-wasi-nn-for-pytorch-reviewer — Review WASI NN for Pytorch using its official documentation and repository in the Wasm / AI/Machine Learning category.


/cncf-wasi-nn-for-tensorflow-lite-reviewer — Review WASI NN for TensorFlow Lite using its official documentation and repository in the Wasm / AI/Machine Learning category.


/cncf-wasix-reviewer — Review WASIX using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-wasm-bindgen-reviewer — Review wasm-bindgen using its official documentation and repository in the Wasm / Tooling category.


/cncf-wasm-pack-reviewer — Review wasm-pack using its official documentation and repository in the Wasm / Tooling category.


/cncf-wasm-workers-server-reviewer — Review Wasm Workers Server using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-wasm2c-reviewer — Review wasm2c using its official documentation and repository in the Wasm / Tooling category.


/cncf-wasm3-reviewer — Review wasm3 using its official documentation and repository in the Wasm / Runtimes category.


/cncf-wasmcloud-reviewer — Review wasmCloud using its official documentation and repository in the Orchestration &amp; Management / Scheduling &amp; Orchestration category.


/cncf-wasmcloud-wasm-reviewer — Review wasmCloud (Wasm) using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-wasmedge-quickjs-reviewer — Review WasmEdge-Quickjs using its official documentation and repository in the Wasm / Languages category. Scripting languages that support Wasm


/cncf-wasmedge-runtime-reviewer — Review WasmEdge Runtime using its official documentation and repository in the Runtime / Container Runtime category. WasmEdge is a lightweight, high-performance, and extensible WebAssembly runtime for cloud native, edge, and decentralized applications. It powers serverless apps, embedded functions, microservices, smart contracts, and IoT devices


/cncf-wasmer-edge-reviewer — Review Wasmer Edge using its official documentation and repository in the Wasm / Hosted Platforms category.


/cncf-wasmer-registry-reviewer — Review Wasmer Registry using its official documentation and repository in the Wasm / Packaging, Registries &amp; Application Delivery category.


/cncf-wasmer-reviewer — Review Wasmer using its official documentation and repository in the Wasm / Runtimes category.


/cncf-wasmex-reviewer — Review Wasmex using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-wasmtime-reviewer — Review Wasmtime using its official documentation and repository in the Wasm / Runtimes category.


/cncf-wavecon-member-reviewer — Review Wavecon (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wavm-reviewer — Review WAVM using its official documentation and repository in the Wasm / Runtimes category.


/cncf-wazero-reviewer — Review Wazero using its official documentation and repository in the Wasm / Runtimes category.


/cncf-weaviate-reviewer — Review Weaviate using its official documentation and repository in the App Definition and Development / Database category.


/cncf-webassembly-language-runtimes-for-python-e4048f96-reviewer — Review WebAssembly Language runtimes for Python and PHP by VMware Labs using its official documentation and repository in the Wasm / Languages category. Scripting languages that support Wasm


/cncf-webiny-reviewer — Review Webiny using its official documentation and repository in the Serverless / Framework category.


/cncf-wehkamp-retail-group-supporter-reviewer — Review Wehkamp Retail Group (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-weight-and-biases-wandb-reviewer — Review Weight and Biases (wandb) using its official documentation and repository in the AI Native Infra / Observability category. A tool for visualizing and tracking your machine learning experiments. This repo contains the CLI and Python API


/cncf-werf-reviewer — Review werf using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. werf is a solution for implementing efficient and consistent software delivery to Kubernetes. It covers the entire CI/CD lifecycle and all related artifacts, glues commonly used tools (Git, Docker/Buildah, Helm, K8s) and facilitates best practices.


/cncf-whatap-member-reviewer — Review WhaTap (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-whatap-reviewer — Review WhaTap using its official documentation and repository in the Observability and Analysis / Observability category. WhaTap Labs provides an integrated monitoring service for DevOps that analyzes application performance issue running on kubernetes in real time.


/cncf-whitesource-reviewer — Review WhiteSource using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-whitestack-member-reviewer — Review Whitestack (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-whitestack-whitemist-reviewer — Review Whitestack WhiteMist using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. WhiteMist is Whitestack&#39;s own Kubernetes distribution, geared at accelerating the adoption of containers in Cloud providers, E-commerce providers, Large organizations and Telecom Operators.


/cncf-whizus-kcntp-reviewer — Review WhizUs (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. As Enterprise Kubernetes Consultants, we support you to build up your Kubernetes cluster in your chosen environment according to your needs.


/cncf-whizus-kcsp-reviewer — Review WhizUs (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. As Enterprise Kubernetes Consultants, we support you to build up your Kubernetes cluster in your chosen environment according to your needs.


/cncf-whizus-member-reviewer — Review WhizUs (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wiit-kcsp-reviewer — Review WIIT (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We provide CNCF certified and highly compliant managed kubernetes for edge and cloud. From consulting to day 2 operations. Made in Europe - 100% data protection and ISO27001 certified.


/cncf-wiit-member-reviewer — Review WIIT (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wikimedia-member-reviewer — Review Wikimedia (member) using its official documentation and repository in the CNCF Members / Nonprofit category.


/cncf-wind-river-member-reviewer — Review Wind River (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wind-river-studio-cloud-platform-reviewer — Review Wind River Studio Cloud Platform using its official documentation and repository in the Platform / Certified Kubernetes - Distribution category. Wind River Cloud Platform is a carrier-grade Kubernetes solution that makes 5G possible by solving the operational problem of deploying and managing distributed edge networks at scale.  It is Wind River&#39;s commercial product based on the open source StarlingX project.


/cncf-winterjs-reviewer — Review WinterJS using its official documentation and repository in the Wasm / Application Frameworks category.


/cncf-wise2c-member-reviewer — Review Wise2C (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wit-binddgen-reviewer — Review Wit-binddgen using its official documentation and repository in the Wasm / Tooling category.


/cncf-witc-reviewer — Review Witc using its official documentation and repository in the Wasm / Tooling category.


/cncf-wiz-member-reviewer — Review Wiz (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-wks-wiit-kubernetes-service-reviewer — Review WKS - Wiit Kubernetes Service using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. The reliable container platform managed by WIIT.


/cncf-woodpecker-ci-reviewer — Review Woodpecker CI using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-wso2-api-microgateway-reviewer — Review WSO2 API Microgateway using its official documentation and repository in the Orchestration &amp; Management / API Gateway category.


/cncf-wso2-kcsp-reviewer — Review WSO2 (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category.


/cncf-wso2-member-reviewer — Review WSO2 (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-x-cellent-kcntp-reviewer — Review x-cellent (KCNTP) using its official documentation and repository in the Special / Kubernetes and Cloud Native Training Partner category. We design and develop cloud solutions with an agile mindset, creating value for our customers and the open source community. An important point in this respect is our customized training. So we offer not only basic and advanced training but also individual training based on the wishes of our customers. We build your cloud native future.


/cncf-x-cellent-kcsp-reviewer — Review x-cellent (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. x-cellent provides feature development and support for Kubernetes on-prem based on metal-stack.io, security audits, cloud adoption and cloud native software development.


/cncf-x-cellent-member-reviewer — Review x-cellent (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-x-ion-member-reviewer — Review x-ion (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-xgboost-reviewer — Review XGBoost using its official documentation and repository in the Data / Data Science category. An optimized distributed gradient boosting library for machine learning.


/cncf-xl-deploy-reviewer — Review XL Deploy using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category.


/cncf-xline-reviewer — Review Xline using its official documentation and repository in the Orchestration &amp; Management / Coordination &amp; Service Discovery category. Xline is a high-performance geo-distributed metadata management system, which is compatible with the ETCD interface.


/cncf-xoap-reviewer — Review XOAP using its official documentation and repository in the Provisioning / Automation &amp; Configuration category. Hybrid IT infrastructure and workplace automation platform


/cncf-xperf-member-reviewer — Review XPerf (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-xregistry-reviewer — Review xRegistry using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. The xRegistry project defines an abstract model for managing metadata about resources and provides a REST-based interface to discover, create, modify and delete those resources.


/cncf-xsky-reviewer — Review XSKY using its official documentation and repository in the Runtime / Cloud Native Storage category. XSKY (Beijing) Data Technology Co., Ltd. is a high-tech enterprise focusing on software defined infrastructure, providing software defined distributed storage product of enterprise-grade and helping customers achieve innovation in data structure.


/cncf-xtdb-reviewer — Review XTDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-yahoo-supporter-reviewer — Review Yahoo (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-ybor-ai-member-reviewer — Review Ybor.AI (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-ydb-reviewer — Review YDB using its official documentation and repository in the App Definition and Development / Database category. YDB is an open-source Distributed SQL Database that combines high availability and scalability with strong consistency and ACID transactions.


/cncf-yellowbrick-member-reviewer — Review Yellowbrick (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-yellowbrick-reviewer — Review Yellowbrick using its official documentation and repository in the App Definition and Development / Database category. Yellowbrick offers a SQL data warehouse built from the ground up on Kubernetes.


/cncf-yellowdog-reviewer — Review YellowDog using its official documentation and repository in the Platform / PaaS/Container Service category.


/cncf-yld-kcsp-reviewer — Review YLD (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. We modernise both practices and software. Applying expertise in a range of CNCF technologies and an inherent culture of innovation.


/cncf-yld-member-reviewer — Review YLD (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-yoke-reviewer — Review yoke using its official documentation and repository in the App Definition and Development / Application Definition &amp; Image Build category. Yoke is an IaC client-side package manager that deploys applications packaged as WASM executables, allowing users to leverage code to define their Applications instead of yaml templates. The yoke project includes a server-side component called the Air-Traffic-Controller allowing users to define their packages as CRDs and have them deployed natively server-side. Yoke is to helm and kro what pulumi is to terraform. Infrastructure as Code, but actually code.


/cncf-yomo-reviewer — Review YOMO using its official documentation and repository in the Wasm / Embedded Functions category.


/cncf-yonder-kcsp-reviewer — Review Yonder (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Yonder is a trusted technology partner that combines AI innovation, DevOps practices, and Kubernetes-powered cloud infrastructure to help established software companies to accelerate their digital transformation. We deliver reliable cloud-native platforms with end-to-end services that design, build, and operate the systems your business depends on.


/cncf-yonder-member-reviewer — Review Yonder (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-youki-reviewer — Review youki using its official documentation and repository in the Runtime / Container Runtime category.


/cncf-youki-wasm-reviewer — Review youki (wasm) using its official documentation and repository in the Wasm / Orchestration &amp; Management category.


/cncf-yrcloudfile-reviewer — Review YRCloudFile using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-yugabytedb-member-reviewer — Review YugabyteDB (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-yugabytedb-reviewer — Review YugabyteDB using its official documentation and repository in the App Definition and Development / Database category.


/cncf-yunikorn-reviewer — Review Yunikorn using its official documentation and repository in the AI Native Infra / Orchestration and Scheduling category. Unleash the power of resource scheduling for running Batch, Data &amp; ML on Kubernetes!


/cncf-zabbix-reviewer — Review Zabbix using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-zalando-supporter-reviewer — Review Zalando (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-zededa-inc-member-reviewer — Review ZEDEDA, Inc. (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-zendesk-supporter-reviewer — Review Zendesk (supporter) using its official documentation and repository in the CNCF Members / End User Supporter and Contributor category.


/cncf-zenko-reviewer — Review Zenko using its official documentation and repository in the Runtime / Cloud Native Storage category.


/cncf-zenml-member-reviewer — Review ZenML (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-zesty-member-reviewer — Review Zesty (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-zesty-reviewer — Review Zesty using its official documentation and repository in the Observability and Analysis / Continuous Optimization category.


/cncf-zettaset-reviewer — Review Zettaset using its official documentation and repository in the Provisioning / Security &amp; Compliance category.


/cncf-ziax-member-reviewer — Review Ziax (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-zig-reviewer — Review Zig using its official documentation and repository in the Wasm / Languages category. Compiled language to Wasm


/cncf-zipkin-reviewer — Review Zipkin using its official documentation and repository in the Observability and Analysis / Observability category.


/cncf-zitadel-reviewer — Review Zitadel using its official documentation and repository in the Provisioning / Security &amp; Compliance category. ZITADEL is the identity infrastructure, simplified for you.


/cncf-zoi-kcsp-reviewer — Review Zoi (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. Zoi is a cloud-native IT consulting company that specializes in the latest cloud technology. We are a trusted partner to manufacturing &amp; retail companies in the DACH and Iberian Peninsula regions, and together develop &amp; implement long-term digital strategies to help them remain competitive in an increasingly digital world. As a KCSP partner, we have validated expertise helping organisations successfully adopt Kubernetes. We offer Kubernetes support, consulting, and professional services for clients embarking on their Kubernetes journey.


/cncf-zoi-member-reviewer — Review Zoi (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-zot-reviewer — Review zot using its official documentation and repository in the Provisioning / Container Registry category. Zot is an OCI-native container registry for distributing container images and OCI artifacts.


/cncf-zstack-kcsp-reviewer — Review ZStack (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category. ZStack enables customers on their Kubernetes journey, helps enterprises focus on their applications rather than managing IT-infrastructure, and offers full life-cycle services including training, consulting, deployment, and pipeline to accelerate their digital transformation.


/cncf-zstack-member-reviewer — Review ZStack (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-zte-kcsp-reviewer — Review ZTE (KCSP) using its official documentation and repository in the Special / Kubernetes Certified Service Provider category.


/cncf-zte-member-reviewer — Review ZTE (member) using its official documentation and repository in the CNCF Members / Gold category.


/cncf-zte-tecs-openpalette-reviewer — Review ZTE TECS OpenPalette using its official documentation and repository in the Platform / Certified Kubernetes - Hosted category. The TECS OpenPalette platform is a unified cloud platform oriented to the ICT realm.


/cncf-zuplo-member-reviewer — Review Zuplo (member) using its official documentation and repository in the CNCF Members / Silver category.


/cncf-zuul-reviewer — Review Zuul using its official documentation and repository in the App Definition and Development / Continuous Integration &amp; Delivery category. Zuul is a system that drives continuous integration, delivery, and deployment with a focus on project gating and interrelated projects.



Project Subagents

Claude Code project subagents are stored under:

.claude/agents/<subagent-name>.md

Use subagents for bounded side tasks that would otherwise flood the main conversation with file reads, grep results, logs, test output, or security review detail.



requirements-analyst — Use proactively to analyze issues, user stories, acceptance criteria, constraints, risks, and missing requirements before implementation.


architecture-reviewer — Use proactively for architecture, module boundaries, interfaces, coupling, scalability, data flows, and technical design risks.


devsecops-reviewer — Use proactively after code, CI/CD, dependency, IaC, GitOps, or security-sensitive changes to review DevSecOps risk and merge readiness.


security-reviewer — Use proactively for authentication, authorization, input validation, file handling, permissions, secrets, and security-sensitive code changes.


ci-cd-reviewer — Use proactively for GitLab CI, GitHub Actions, runners, deployment jobs, caches, artifacts, tokens, and pipeline governance.


iac-gitops-reviewer — Use proactively for Terraform, Kubernetes, Helm, Kustomize, GitOps, environment promotion, reconciliation, and infrastructure changes.


test-runner — Use proactively to run relevant tests, analyze failures, and summarize validation results after implementation.


release-readiness-reviewer — Use proactively before release readiness claims to review tests, rollback, migrations, feature flags, monitoring, documentation, and breaking changes.


incident-postmortem-assistant — Use for incident analysis, log summaries, timeline reconstruction, root cause analysis, corrective actions, and follow-up issues.

Subagent routing

Use the matching subagent proactively when the task fits its description.

Recommended routing:

Requirements or acceptance criteria unclear → requirements-analyst
Architecture or module boundaries affected → architecture-reviewer
Security-sensitive behavior affected → security-reviewer
CI/CD or deployment logic changed → ci-cd-reviewer
Terraform, Kubernetes, Helm, Kustomize, or GitOps changed → iac-gitops-reviewer
Tests should be executed or failures analyzed → test-runner
DevSecOps merge-readiness review needed → devsecops-reviewer
Release readiness must be assessed → release-readiness-reviewer
Incident, outage, logs, or postmortem analysis needed → incident-postmortem-assistant
SDLC skill routing

For feature work, prefer:

/requirements-analyst
→ /cost-based-planner
→ /architecture-reviewer when architecture is affected
→ /threat-modeler when security-sensitive behavior is affected
→ /safe-implementer
→ /test-strategy-engineer
→ /verification-reviewer
→ /security-reviewer when needed
→ /documentation-maintainer when needed
→ /release-readiness-reviewer before release readiness claims

For CI/CD, dependency, IaC, GitOps, secrets, compliance, observability, or incident tasks, use the matching specialized skill or subagent before claiming completion.

Safety model
Prefer read-only subagents for analysis and review.
Use implementation through the main Claude Code session unless a subagent is explicitly designed to modify files.
Do not let review subagents perform broad rewrites.
Treat subagent output as review input; the main session remains responsible for final user-visible conclusions.

---

## `tests/test_templates.py` Ergänzung

In `EXPECTED_TEMPLATES` muss diese Datei enthalten sein:

```python
"claude/agent.md.j2",
