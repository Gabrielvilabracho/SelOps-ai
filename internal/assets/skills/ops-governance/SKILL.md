---
name: ops-governance
description: "SelOps governance policies, approval workflows, and compliance. Trigger: When applying governance policies, change approvals, or compliance checkpoints."
---

# Governance

## When to Use

Load this skill when processing a change approval for a production AI system, when a client engagement requires compliance documentation, when classifying an AI system incident, when setting up audit trail infrastructure for a new client, or when reviewing whether a proposed change meets governance requirements before execution.

## Core Principles

- **Change approval is mandatory before production AI model updates.** Updating a model, modifying a prompt template in production, or changing a system configuration that affects model behavior requires explicit human approval. These are not routine deployments.
- **Compliance checkpoints happen in the pipeline, not at release.** Compliance gates (bias review, regulatory check, security scan) are CI steps that block merge. They are not review tasks deferred to a release meeting.
- **Audit trails are immutable.** Every privileged operation — model deployment, configuration change, data access, incident response action — is logged to an append-only store. Logs are never modified after the fact.
- **AI incidents have their own classification system.** A model that confidently produces wrong output is not the same class of incident as a service outage. Governance must classify both. Incident classification determines the response procedure, the client notification requirement, and the post-mortem format.
- **Regulated industries require explicit governance contracts.** Before going live with AI systems in finance, healthcare, legal, or other regulated domains, document which regulatory framework applies, what the compliance requirements are, and how the system satisfies them. This document is a deliverable, not internal notes.
- **Model Inventory is a governance requirement, not an audit artifact.** Every model in production for a client must be listed in the model inventory: model ID, version, provider, purpose, risk classification, deployment date, owner, and review schedule. This applies to both SelOps-deployed models and client-owned models that SelOps operates. (SR 11-7 / NIST AI RMF GOVERN 1.6)
- **Third-party model governance extends your governance obligations.** When using foundation models from external providers (OpenAI, Anthropic, AWS Bedrock, etc.), apply NIST AI RMF GOVERN 6.1 criteria: document the provider, the model version, what due diligence was performed, what the fallback is if the provider changes the model, and what contractual protections exist. Outsourcing the inference does not outsource the governance risk.

## Patterns

### Change Approval Workflow for AI Systems
Every production AI change follows this sequence: (1) change description submitted with impact assessment, (2) technical reviewer approves — confirms rollback exists and tests pass, (3) compliance reviewer approves if the change affects regulated functionality, (4) client owner approves for client-facing changes, (5) change executed by operator, (6) post-execution state verified and logged. Skipping any step requires explicit exception approval documented in the audit trail.

### AI-Specific Incident Classification
Classify AI incidents into four categories:
- **Output quality degradation**: model produces outputs with lower accuracy or higher error rate than baseline, no service outage. Client may not notice immediately.
- **Safety/compliance violation**: model produces outputs that violate defined safety constraints, legal requirements, or client-specific rules. Immediate escalation required.
- **Service unavailability**: model or pipeline is unreachable, returns errors consistently. Standard incident response applies.
- **Data exposure**: model reveals data that should not be accessible (PII leakage, cross-client data contamination). Immediate incident response and legal notification may apply.

Each classification has a documented response playbook. Apply the playbook for the classification, not a generic incident procedure.

### Bias and Fairness Review Gate
For systems that make decisions affecting people (credit, hiring, medical, legal), a bias review is required before any model update ships. The review must: (1) run the updated model against a defined fairness benchmark dataset, (2) compare outputs across demographic groups, (3) document the results and the reviewer's sign-off. If fairness metrics degrade beyond the defined threshold, the update is blocked.

### Immutable Audit Log Format
Every audit log entry includes: timestamp (UTC, millisecond precision), actor (agent ID or human identifier), action type, target system, target resource identifier, action outcome (success/failure), and a correlation ID linking to the change approval record. Entries are written to an append-only log. No entry is modified or deleted. Retention period is documented per engagement based on regulatory requirements. Verify log format satisfies NIST SP 800-53 AU family requirements — AU-12 (audit record generation), AU-9 (protection of audit information), AU-11 (audit record retention) — for regulated engagements.

### Compliance Checkpoint Registry
Maintain a registry of compliance checkpoints per engagement: what the check is, which regulatory requirement it satisfies, where in the pipeline it runs, who owns the check, and the last passing result. This registry is reviewed at each engagement milestone and updated when regulations or client requirements change.

### Model Update Governance Record
Every production model update produces a governance record: model name and version (old and new), reason for update, test results comparing old and new on benchmark datasets, rollback procedure, approval chain with timestamps, and post-deployment verification results. This record is stored in the engagement artifact store alongside the audit trail.

### Effective Challenge Protocol (SR 11-7)
Before any model goes to production or any significant model update ships, an Effective Challenge must be performed. Effective Challenge is an objective, critical analysis of the model by someone who did not build it. At minimum: (1) review the conceptual soundness — does the model approach the problem in a valid way? (2) review the data and training methodology if accessible, (3) run the model against adversarial inputs and edge cases prepared by the challenger, (4) document the challenger's findings and the model team's response. The challenger must have sufficient technical expertise and no development stake in the outcome. For client engagements in finance (SR 11-7), healthcare (FDA), or high-risk AI under the EU AI Act, Effective Challenge by an independent party is required before production. "We tested it and it works" is not Effective Challenge.

### Model Inventory Record (NIST AI RMF GOVERN 1.6 / SR 11-7)
Maintain a Model Inventory entry for every model in production. Each entry includes: model identifier (name, version, provider), deployment environment, purpose and use case, risk classification (using the engagement's risk framework), deployment date, approved-by (with timestamp), review schedule (at minimum: after each update, annually otherwise), known limitations and biases, and decommission date (if planned). The inventory is a living document — not a one-time registry. It is reviewed at each engagement milestone.

### Reversibility Classification and Rollback Governance
Classify every production action as reversible or irreversible before execution:

- **Reversible**: can return the system to its prior state within 15 minutes using documented steps. Examples: configuration changes with saved prior state, additive database migrations, feature flag toggles.
- **Irreversible**: cannot be undone. Examples: data deletion, external system notifications (emails sent, webhooks fired), financial transactions, regulatory filings, model retraining that overwrites prior weights without backup.

For irreversible actions, apply compensating controls:
1. **Explicit approval**: require approval at the appropriate autonomy level before executing — regardless of the registered autonomy level, irreversible actions always require human confirmation.
2. **Compensating control**: document the containment action if the irreversible step produces an unacceptable outcome (e.g., notify affected parties, restore from backup, file correction with regulator).
3. **Post-execution audit record**: log the irreversible action with an explicit `irreversible: true` field in the audit entry.

The reversibility classification is made at the structure phase (see ops-structure) and is binding — it cannot be downgraded during execution. When uncertain whether an action is reversible, default to irreversible. This classification is the domain foundation for risk classification in ops-brief and rollback requirements in ops-graduated-autonomy.

## Checklist

**Before executing a production AI change:**
- [ ] Change description submitted with impact scope
- [ ] Technical review complete and documented
- [ ] Compliance review complete (if applicable)
- [ ] Client owner approval obtained (if client-facing)
- [ ] Rollback procedure written and verified
- [ ] Audit trail entry prepared for the change action

**When responding to an AI incident:**
- [ ] Incident classified using the four-category system
- [ ] Client notified per the classification's notification requirement
- [ ] Response playbook for this classification identified and followed
- [ ] All response actions logged to the audit trail with timestamps
- [ ] Post-mortem scheduled per the classification's requirement

**When setting up a new regulated-industry engagement:**
- [ ] Applicable regulatory framework documented
- [ ] Compliance requirements listed and mapped to system components
- [ ] Compliance checkpoints added to CI pipeline
- [ ] Bias/fairness benchmark dataset defined (if applicable)
- [ ] Audit trail retention period set per regulatory requirement
- [ ] Client signed off on the governance contract
- [ ] Model Inventory created with entries for all models to be deployed
- [ ] Effective Challenge scheduled (identify the challenger before go-live, not after) — required for finance (SR 11-7), healthcare, and EU high-risk AI
- [ ] Third-party model governance documented per NIST AI RMF GOVERN 6.1 (provider, version, fallback, contractual protections)
- [ ] Audit log format verified against NIST SP 800-53 AU family requirements — AU-12 (generation), AU-9 (protection), AU-11 (retention)
- [ ] For EU high-risk AI systems: verify audit log satisfies EU AI Act Art.12 automatic logging requirements

## SelOps-Specific Context

Governance for AI systems is not the same as governance for traditional software. A model update can change system behavior in ways that are impossible to fully enumerate in advance — no test suite covers every possible input. This means AI governance requires probabilistic risk assessment and sampling-based verification, not just deterministic test pass/fail gates.

At SelOps, we operate AI systems in production for clients who may themselves be in regulated industries. Our governance obligations are therefore composite: SelOps's own internal standards plus whatever the client's regulatory environment requires. When these conflict, the stricter requirement applies. When they are unclear, engage the client's legal or compliance team before proceeding.

The audit trail is a client asset, not a SelOps internal record. At engagement close, the full audit trail is transferred to the client. Design it from day one to be readable and useful by the client's own compliance and legal teams, not just by SelOps engineers.

AI incident classification matters because clients (and their legal teams) will ask: "Did the AI make a wrong decision, or did it fail to respond?" These have different legal implications in regulated industries. Classify accurately, document clearly, and never characterize an output quality failure as a service issue to minimize apparent severity.

The SR 11-7 / SR 26-2 Effective Challenge requirement deserves emphasis. In banking and financial services, an AI model that affects credit decisions, risk scoring, or fraud detection cannot go to production without independent validation and Effective Challenge. SelOps cannot be both the builder and the only reviewer. Identify an independent challenger — whether internal (a SelOps team member with no development stake) or external — and document their findings as part of the governance record.

## References

- SR 11-7 (2011) — Federal Reserve / OCC Supervisory Guidance on Model Risk Management. Foundational source for Effective Challenge, Model Inventory, and Independent Validation requirements in banking and financial services.
- NIST AI RMF 1.0 (2023) — GOVERN function, including GOVERN 1.6 (mechanisms to inventory AI systems) and GOVERN 6.1 (policies addressing AI risks from third-party entities).
- NIST SP 800-53 Rev 5 — AU family (Audit and Accountability): AU-12 (Audit Record Generation), AU-9 (Protection of Audit Information), AU-11 (Audit Record Retention). Applies to regulated audit log format requirements.
- EU AI Act — Regulation (EU) 2024/1689, Art.12 (Record-keeping and logging for high-risk AI systems). Mandatory for EU high-risk AI deployments.
