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
Every audit log entry includes: timestamp (UTC, millisecond precision), actor (agent ID or human identifier), action type, target system, target resource identifier, action outcome (success/failure), and a correlation ID linking to the change approval record. Entries are written to an append-only log. No entry is modified or deleted. Retention period is documented per engagement based on regulatory requirements.

### Compliance Checkpoint Registry
Maintain a registry of compliance checkpoints per engagement: what the check is, which regulatory requirement it satisfies, where in the pipeline it runs, who owns the check, and the last passing result. This registry is reviewed at each engagement milestone and updated when regulations or client requirements change.

### Model Update Governance Record
Every production model update produces a governance record: model name and version (old and new), reason for update, test results comparing old and new on benchmark datasets, rollback procedure, approval chain with timestamps, and post-deployment verification results. This record is stored in the engagement artifact store alongside the audit trail.

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

## SelOps-Specific Context

Governance for AI systems is not the same as governance for traditional software. A model update can change system behavior in ways that are impossible to fully enumerate in advance — no test suite covers every possible input. This means AI governance requires probabilistic risk assessment and sampling-based verification, not just deterministic test pass/fail gates.

At SelOps, we operate AI systems in production for clients who may themselves be in regulated industries. Our governance obligations are therefore composite: SelOps's own internal standards plus whatever the client's regulatory environment requires. When these conflict, the stricter requirement applies. When they are unclear, engage the client's legal or compliance team before proceeding.

The audit trail is a client asset, not a SelOps internal record. At engagement close, the full audit trail is transferred to the client. Design it from day one to be readable and useful by the client's own compliance and legal teams, not just by SelOps engineers.

AI incident classification matters because clients (and their legal teams) will ask: "Did the AI make a wrong decision, or did it fail to respond?" These have different legal implications in regulated industries. Classify accurately, document clearly, and never characterize an output quality failure as a service issue to minimize apparent severity.
