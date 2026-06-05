---
name: ops-standard-documentation
description: "SelOps operational documentation standards. Trigger: When writing or reviewing docs, READMEs, ADRs, or API contracts for the SelOps platform."
---

# Standard Documentation

## When to Use

Load this skill when writing or reviewing any documentation for an AI system: architecture docs, operational runbooks, API contracts, ADRs, client handoff packages, or model cards. Use it when documentation is a required deliverable or when existing documentation needs to be brought into compliance with SelOps standards.

## Core Principles

- **Documentation is a system artifact, not an afterthought.** Docs live in the same repository as the system they describe. They are reviewed in PRs. They go stale when the system changes and they are not updated. Stale docs are a liability, not a resource.
- **Operational runbooks are executable.** A runbook that says "investigate the database" is not a runbook. A runbook must specify: which tool, which query, what a healthy result looks like, what an unhealthy result looks like, and what to do next in each case.
- **ADRs capture decisions, not implementations.** An ADR answers: what was the decision, why was this option chosen, what alternatives were considered, and what are the consequences we accepted. ADRs are written at decision time, not reconstructed later.
- **Client handoff documentation must be self-sufficient.** After SelOps hands off a system, the client must be able to operate it, debug common issues, and make routine changes without calling SelOps. If they cannot do that with the docs provided, the handoff is incomplete.
- **Model cards are required for every deployed model.** A model card documents what the model is, what it was trained on, what it is good at, what it fails at, and what the known biases and limitations are. It is the foundational document for AI system governance.
- **Model cards serve two audiences.** Internal governance audiences need technical detail: training data, evaluation benchmarks, known failure modes, and update policy (ISO 42001 A.8.2 [UNVERIFIED: exact clause number]). External or end-user audiences need plain-language disclosure: what the system does, what it cannot do, how to escalate, and when a human is in the loop (EU AI Act transparency requirements [UNVERIFIED: exact article in adopted Regulation (EU) 2024/1689]).

## Patterns

### AI System Architecture Document
Every AI system has a top-level architecture document covering: (1) system purpose and scope, (2) component diagram showing module boundaries, (3) data flow from input to output including model calls, (4) external dependencies (model providers, data sources, client systems), (5) deployment architecture, (6) known limitations and constraints. This document is updated whenever the architecture changes — not at release, at change.

### Model Card Standard
Each model card contains: model identifier (name, version, provider), intended use cases, training data summary (sources, date range, known gaps), evaluation results on benchmark datasets, known failure modes and edge cases, known biases with assessment methodology, and update/retirement policy. Model cards are stored with the system they support and versioned alongside the system.

### User-Facing Model Card
For systems deployed to external users or clients, produce a second model card variant in plain language. Required fields: (1) what this system does in one sentence, (2) what it cannot do or where it is unreliable, (3) whether a human reviews its outputs and under what conditions, (4) how to report a problem or escalate, (5) when and how the system is updated, (6) how to opt out or request human review. This card is part of the client handoff package. Grounded in ISO 42001 A.8.2 (user information and transparency) [UNVERIFIED: exact clause number] and EU AI Act transparency provisions for high-risk and general-purpose AI systems (Regulation (EU) 2024/1689) [UNVERIFIED: exact article number in adopted text; Art.50 in near-final draft].

### Operational Runbook Structure
Every runbook follows this structure: (1) purpose — what operational scenario this covers, (2) prerequisites — what access, tools, and context you need before starting, (3) step-by-step procedure — concrete commands with expected outputs, (4) decision points — what each possible output means and which step to go to next, (5) escalation path — when to stop and who to contact, (6) post-procedure checklist — what to verify when done. No open-ended steps.

### ADR Template
Title: `[number]-[short-slug].md` (e.g., `005-prompt-registry-approach.md`)
Sections: **Status** (proposed/accepted/deprecated/superseded), **Context** (why this decision was needed), **Decision** (what was decided, stated simply), **Alternatives considered** (what was evaluated and why rejected), **Consequences** (what we accept by choosing this — both positive and negative).

### Predetermined Change Control Plan (FDA PCCP)
For AI systems deployed in regulated healthcare contexts, maintain a Predetermined Change Control Plan documenting: (1) scope of permitted changes — which parameters, thresholds, or model components may change without full revalidation, (2) validation methodology — the test suite and performance criteria that must be met before a permitted change is deployed, (3) performance criteria — quantitative thresholds that define acceptable operation after each change, (4) monitoring plan — how the system is observed post-change to detect drift or degradation. This document is submitted to or reviewed by the regulatory authority (FDA AI/ML-Based SaMD guidance: "Predetermined Change Control Plan") and must be kept current with each permitted change cycle.

### Privacy Impact and DPIA Documentation
Before deploying an AI system that processes personal data, complete a Data Protection Impact Assessment when: the system involves large-scale processing of sensitive personal data, systematic monitoring of individuals, or automated decision-making with significant effects (GDPR Art.35). The DPIA must cover: (1) a description of the processing and its purposes, (2) an assessment of necessity and proportionality, (3) risks to the rights and freedoms of individuals, (4) measures to address those risks. The completed DPIA is stored in the project repository and reviewed when the system's data scope changes.

### Data Lineage Document
For AI pipelines that process client data, maintain a data lineage document: source systems and their schemas, transformation steps with their logic, where data lands and in what format, what is retained and for how long, what is purged and when. This document is required for regulated industries and recommended for all engagements.

### Living Docs Synchronization Rule
When a system change is made, documentation must be updated in the same PR. No separate "docs PR later" — they are merged together. The PR checklist must include: "Have the relevant docs been updated? (architecture, runbook, ADR if a new decision was made, model card if model changed)." A PR that changes the system without updating docs requires explicit justification in the PR description.

### Client Handoff Package
The handoff package is assembled at the end of every engagement. It contains: system architecture document, model card(s) — including the user-facing variant for externally-deployed models, operational runbooks for all common scenarios, API contracts (current version), governance records and audit trail, known issues log, escalation contacts at SelOps (for warranty period), and a 30-minute walkthrough session with the client team. All items are required. Missing items delay handoff.

## Checklist

**When writing a new document:**
- [ ] Document type identified (architecture, runbook, ADR, model card, data lineage, API contract)
- [ ] Stored in the repository alongside the system it describes
- [ ] Written in English (unless the client has explicitly requested another language)
- [ ] No open-ended steps (every step produces a concrete, verifiable result)
- [ ] Reviewed by a second person before being considered complete

**When a system change is made:**
- [ ] Architecture document updated if component boundaries, data flows, or dependencies changed
- [ ] Runbook updated if operational procedures changed
- [ ] ADR written if a new design decision was made
- [ ] Model card updated if the model, prompt, or evaluation results changed
- [ ] Data lineage updated if data sources or retention policies changed
- [ ] DPIA reviewed if processing scope or data types changed (GDPR Art.35)

**Before client handoff:**
- [ ] System architecture document complete and up to date
- [ ] Model card(s) complete for all deployed models (technical variant)
- [ ] User-facing model card produced for externally-deployed models (plain-language, 6 required fields)
- [ ] Runbooks cover all common operational scenarios
- [ ] API contracts reflect current production version
- [ ] Governance records and audit trail transferred
- [ ] Known issues log is honest and complete
- [ ] DPIA completed and on file if the system processes personal data at scale (GDPR Art.35)
- [ ] PCCP document current if system is deployed in a regulated healthcare context (FDA PCCP guidance)
- [ ] Walkthrough session scheduled

## SelOps-Specific Context

Documentation for AI systems presents challenges that traditional software documentation does not. The system's behavior is partially determined by model weights and prompt templates — components that cannot be fully described by architecture diagrams alone. A model card fills the gap that a component diagram leaves.

At SelOps, documentation is part of what we sell. A client who receives an AI system without adequate documentation has received a partial deliverable. They cannot maintain it, cannot govern it, cannot comply with regulations around it, and cannot explain it to their own stakeholders. This has legal and reputational consequences for SelOps.

The "living docs" principle is especially important for AI systems because they change in ways that are less visible than traditional software changes. A prompt template update is not a code change in the traditional sense, but it changes system behavior. A model provider's silent model update is not our code change at all, but it changes our system's behavior. Documentation processes must account for these invisible changes — scheduled documentation reviews (monthly or per-engagement-milestone) catch drift that PR-triggered updates miss.

Runbooks in AI systems must cover not just infrastructure failures but output quality failures: what to do when the error rate on output validation rises, how to manually verify a model response in production, how to roll back a prompt template change, and how to declare an AI output quality incident. These are the scenarios that will actually occur. Document them before they do.

## References

- ISO/IEC 42001:2023 — AI management system standard; A.8.2 addresses AI system user information and transparency [UNVERIFIED: exact clause number]
- Regulation (EU) 2024/1689 (EU AI Act) — transparency obligations for high-risk and general-purpose AI systems [UNVERIFIED: exact article number in adopted text; Art.50 referenced in near-final draft]
- FDA AI/ML-Based SaMD Action Plan, "Predetermined Change Control Plan" (PCCP) guidance — real FDA concept for permitting scoped changes to AI/ML-based medical software
- GDPR Art.35 — Data Protection Impact Assessment (DPIA) trigger conditions for large-scale processing, systematic monitoring, and automated decision-making
