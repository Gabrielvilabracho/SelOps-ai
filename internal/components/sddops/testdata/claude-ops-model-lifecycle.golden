---
name: ops-model-lifecycle
description: "SelOps AI model lifecycle governance. Trigger: When versioning, retiring, or decommissioning a production AI model."
---

# Model Lifecycle

## When to Use

Load this skill when making any decision about a model's production status: introducing a new model version, planning deprecation of an existing version, executing a decommission, or defining how long model records must be kept. Use it when the question is "what happens to this model next?" or "what do we owe the client and regulators when this model changes?"

## Core Principles

- **Every production model has a documented lifecycle.** A model does not exist in production without a version record that says when it arrived, what it replaced, what its approval basis was, and what will trigger its replacement. Undocumented models are ungovernable. (NIST GOVERN 1.6 — mechanisms to inventory AI systems.)
- **Versioning requires behavioral baselines.** Assigning a version number is not enough. Each model version must be accompanied by a behavioral baseline — a set of benchmark results, evaluation outputs, or regression test outcomes that define what "this version" means. Without a baseline, version changes cannot be assessed and rollbacks cannot be validated.
- **Decommission is a planned event, not a shutdown.** A model is not simply switched off when it is no longer needed. Decommission is a defined procedure: stakeholders are notified, a deprecation window is enforced, data flows are redirected, and records are preserved before the model is removed. (NIST GOVERN 1.7 — processes for decommissioning and phasing out AI systems safely.)
- **Record preservation is a regulatory obligation.** The records that accompany a model — training data summaries, evaluation results, approval decisions, incident logs — must be preserved after decommission for a defined retention period. These records are required for post-incident investigations, audits, and regulatory reviews. Destroying them prematurely is a compliance failure. (ISO/IEC 42001:2023 A.6.2.8 [UNVERIFIED: confirm clause against purchased ISO standard] — AI system records and retention.)

## Patterns

### Model Version Control Record
For each model version in production, maintain a version record containing: (1) model identifier — name, version string, provider, (2) deployment date and approving party, (3) what this version replaced and why, (4) behavioral baseline — linked benchmark results or evaluation report, (5) known limitations or deviations from the previous version, (6) expected replacement trigger — the condition or date that will initiate the next lifecycle event, (7) current lifecycle status — active / deprecated / decommissioned. The record lives in source control or a governed artifact store; it is not a wiki entry that can be silently edited. Grounded in NIST GOVERN 1.6 model inventory requirements and ISO/IEC 42001:2023 A.6.2.3 [UNVERIFIED: confirm clause against purchased ISO standard] (AI system documentation).

### Retirement Criteria Definition
Before a model is deployed, define in writing the criteria that will trigger its retirement. Acceptable trigger types: (1) performance degradation — a named metric falls below a defined threshold for a defined period, (2) provider end-of-life — the provider announces the model version is discontinued, (3) compliance change — a regulatory or contractual requirement can no longer be met by this model, (4) scheduled review — a calendar-based review reveals the model no longer meets current standards, (5) incident — a safety, accuracy, or bias incident of defined severity automatically opens a retirement review. Document which party has authority to make the retirement decision and within what timeframe they must act once a trigger condition is met.

### Decommission Procedure
When a model version is to be decommissioned (NIST GOVERN 1.7):
1. **Notify stakeholders** — inform all teams and clients that consume the model of the planned decommission date at least the required deprecation window in advance (minimum 30 days unless a security incident requires immediate action).
2. **Enforce the deprecation window** — run the old and new versions in parallel or with a shadow-mode setup during the window; document the migration path and offer support.
3. **Redirect data flows** — update all consumers, integrations, and monitoring to the replacement model or to a documented no-model fallback.
4. **Verify redirection** — confirm in production logs that no traffic is reaching the decommissioned model before proceeding.
5. **Archive records** — move all model-related records (version record, evaluation reports, incident logs, approval decisions) to the long-term archive store with a retention label before removing the model.
6. **Remove model artifacts** — delete or deactivate model weights, endpoints, and references from production systems; document the removal date and approving party.
7. **Close the lifecycle record** — update the version record status to "decommissioned" with the removal date and a pointer to the archived records.

### Record Preservation Policy
Define and document the retention period for model lifecycle records: (1) version records — kept for the longer of 5 years after decommission or the engagement contract's audit period, (2) evaluation and benchmark reports — same period as version records, (3) incident logs — kept for the longer of 5 years or any applicable regulatory requirement, (4) approval decisions and sign-offs — same period as incident logs, (5) training data summaries — kept for the model's active lifetime plus the retention period above. Records are stored in an immutable or write-protected archive; they must not be modified after the model is decommissioned. Grounded in ISO/IEC 42001:2023 A.6.2.8 [UNVERIFIED: confirm clause against purchased ISO standard] (AI system records).

## Checklist

**When deploying a new model version:**
- [ ] Version record created with all required fields (identifier, deployment date, approver, behavioral baseline, replacement trigger)
- [ ] Behavioral baseline documented — benchmark results or evaluation report linked from the version record
- [ ] Retirement criteria defined and signed off before or at deployment
- [ ] Previous version status updated (e.g., from "active" to "deprecated" with deprecation window start date)

**When deprecating a model version:**
- [ ] Stakeholders notified with the decommission date and migration path
- [ ] Deprecation window duration recorded; must meet the minimum unless an emergency exception is approved
- [ ] Replacement model's version record and baseline verified before starting the deprecation window

**When decommissioning a model version:**
- [ ] All consumers and data flows redirected to replacement or fallback
- [ ] Zero traffic to the old model confirmed in production logs
- [ ] All lifecycle records archived before removal
- [ ] Model artifacts removed from production; removal date and approver documented
- [ ] Lifecycle record status set to "decommissioned" with pointer to archive

**Ongoing:**
- [ ] Retention schedule reviewed annually or when regulatory requirements change
- [ ] Archived records confirmed readable and integrity-checked at least annually

## SelOps-Specific Context

AI model lifecycle governance is frequently absent in AI consultancy engagements. Models are deployed, quietly replaced, and eventually forgotten — with no record of what was running in production, when it changed, or what triggered the change. This becomes a serious problem when clients face regulatory audits, post-incident investigations, or contract disputes about system behavior.

At SelOps, lifecycle governance is non-negotiable because we operate in client production environments. When we deploy a model to a client's system, we become accountable for that model's behavior until it is formally decommissioned and its records are handed off. A client who cannot produce model version records in an audit is a client who cannot demonstrate they governed their AI system — and SelOps will be implicated in that failure.

The decommission procedure matters especially for long-running engagements. Clients sometimes run models for years. Without a formal deprecation window and stakeholder notification, a "routine model update" can create a production incident when downstream teams are not ready for the behavioral change. The deprecation window exists to protect the client, not to slow us down.

## References

- NIST AI RMF 1.0 (2023) — GOVERN function; GOVERN 1.6 (mechanisms to inventory AI systems), GOVERN 1.7 (processes for decommissioning and phasing out AI systems safely)
- ISO/IEC 42001:2023 — AI management system standard; A.6.2.3 (AI system documentation), A.6.2.8 (AI system records and retention) [UNVERIFIED: confirm clause numbers against purchased ISO standard]
