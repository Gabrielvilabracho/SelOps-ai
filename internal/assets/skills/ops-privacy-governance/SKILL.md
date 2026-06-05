---
name: ops-privacy-governance
description: "SelOps privacy governance for AI systems. Trigger: When handling personal data in AI pipelines, implementing GDPR controls, or conducting DPIAs for EU clients."
---

# Privacy Governance

## When to Use

Load this skill when designing or reviewing how an AI pipeline handles personal data, when a client engagement involves EU data subjects and GDPR applies, when conducting a Data Protection Impact Assessment (DPIA), when a data subject exercises rights (access, erasure, portability, objection), when classifying what data a system collects or generates, or when assessing cross-border data transfers.

## Core Principles

- **Data minimisation is the first control.** Collect and process only the personal data strictly necessary for the defined purpose (GDPR Art. 5(1)(c)). Before adding any data field to an AI pipeline, confirm it is necessary and document the justification. If uncertain, do not collect it.
- **Every processing operation requires a lawful basis.** Before an AI system processes personal data, identify and document the lawful basis under GDPR Art. 6 (or, for special categories, Art. 9). Consent is one basis — not the default. Legitimate interests, contractual necessity, and legal obligation are equally valid and often more appropriate for AI pipelines.
- **LLMs are PII amplifiers, not PII isolators.** A language model trained on or prompted with personal data may reproduce, infer, or combine it in ways that create new privacy exposures. Treat the model's output as a potential PII surface, not just its input.
- **Privacy by design is an architectural requirement.** Data protection must be embedded into the system from the start, not retrofitted (GDPR Art. 25). This means pseudonymisation by default, access controls scoped to the minimum, and data retention enforced at the pipeline level, not by manual cleanup.
- **Data subject rights are operational requirements.** The right to access, erasure, rectification, portability, and objection (GDPR Art. 15-22) must be fulfillable in practice. If an AI system processes personal data and cannot respond to these rights within the statutory timeframe, the system is non-compliant — regardless of what the privacy notice says.
- **Generative AI introduces privacy risks beyond traditional data processing.** The NIST AI 600-1 Generative AI Profile identifies data-privacy as a distinct GenAI risk: models can memorise and later reproduce personal data from training data, infer sensitive attributes that were never explicitly provided, and leak PII through outputs. These risks are not addressed by access controls alone — they require the output-screening and fine-tuning safeguards in the patterns below. Apply GDPR controls AND the GenAI-specific privacy safeguards from NIST AI 600-1.

## Patterns

### PII Detection and Masking
Before personal data enters an AI model or is logged: (1) classify the data using a defined PII taxonomy (names, emails, national IDs, health data, financial data, location, etc.); (2) apply masking or pseudonymisation appropriate to the processing purpose — tokenise identifiers that need re-linking, redact or hash identifiers that do not; (3) log the masking decision and the pseudonymisation map with restricted access; (4) scan model outputs for PII re-emergence — a model prompted with pseudonymised data may still reconstruct or infer real identifiers; (5) if unmasked PII must pass through for a specific processing step, document the exception with its lawful basis and limit its scope to that step.

### Lawful Basis Register
Maintain a lawful basis register for every personal data processing operation in the AI pipeline: (1) processing operation name; (2) categories of personal data processed; (3) lawful basis (Art. 6 sub-article, or Art. 9 condition for special categories); (4) data subjects affected; (5) purpose and necessity justification; (6) retention period; (7) cross-border transfer status. Review this register at every major system change and at least annually. A missing or expired lawful basis is a processing prohibition, not a paperwork gap.

### Data Protection Impact Assessment (DPIA)
A DPIA (GDPR Art. 35) is mandatory before processing that is likely to result in high risk to individuals. Triggers include: systematic profiling, processing special categories of data at scale, innovative technology with an untested privacy impact (including novel AI systems), or processing that could prevent individuals from exercising their rights. The DPIA must: (1) describe the processing and its purposes; (2) assess necessity and proportionality; (3) identify and assess risks to data subjects; (4) identify and document measures to mitigate those risks; (5) be reviewed and signed off by the Data Protection Officer (DPO) where one is required. A DPIA is a living document — update it when the system changes materially.

### Data Subject Rights Handling
When a data subject exercises a right: (1) confirm the identity of the requestor before disclosing or deleting any data; (2) log the request with timestamp, right claimed, requestor identity, and assigned case owner; (3) check the legal timeframe — GDPR provides 30 days extendable to 90 days for complex requests; (4) identify every system and log where this data subject's data is held, including model-embedded representations if fine-tuning was used; (5) execute the right (provide access, delete, rectify, export, restrict processing) across all identified systems; (6) confirm completion and send a written response to the data subject; (7) retain the case record with evidence of completion for audit purposes.

### Cross-Border Transfer Controls
Transferring personal data from the EU/EEA to a third country requires a lawful transfer mechanism. For AI pipelines: (1) identify every location where personal data is sent — model inference endpoints, vector databases, logging systems, monitoring services; (2) determine whether each location is in a country with an EU adequacy decision, and if not, confirm a valid transfer mechanism is in place (Standard Contractual Clauses, Binding Corporate Rules, or an applicable derogation); (3) document the mechanism in the lawful basis register; (4) assess the legal environment of the destination country — supplementary measures may be required if local laws could undermine the transfer mechanism; (5) review transfer mechanisms after any change in the destination's legal or regulatory environment.

### Privacy-by-Design Checklist for AI Systems
When architecting or reviewing an AI pipeline: (1) data minimisation — confirm each data field collected is necessary; (2) purpose limitation — confirm data collected for one purpose is not repurposed without a new lawful basis; (3) pseudonymisation by default — apply before data enters the model; (4) access controls — confirm only the principals that need personal data for their role can access it; (5) retention enforcement — confirm data is deleted or anonymised when the retention period expires, without requiring manual action; (6) output screening — confirm model outputs are scanned for PII re-emergence before logging or returning to users; (7) DPIA threshold check — confirm whether the system triggers any DPIA mandatory processing type.

## Checklist

**Before an AI system processes personal data:**
- [ ] Lawful basis identified and documented in the lawful basis register for each processing operation
- [ ] PII taxonomy defined: which data categories are collected, from whom, and why
- [ ] Pseudonymisation or masking applied before data enters the model
- [ ] DPIA threshold assessed; DPIA conducted and signed off if mandatory (GDPR Art. 35)
- [ ] Cross-border transfer mechanisms confirmed for every location where personal data is sent
- [ ] Data retention periods defined and automated enforcement configured
- [ ] Data subject rights request process defined and tested (identity verification, 30-day response window)

**When an AI system changes materially (model update, new data source, new use case):**
- [ ] Lawful basis register reviewed for continued accuracy
- [ ] DPIA reviewed and updated if the change increases risk to data subjects
- [ ] PII masking coverage verified against new data inputs
- [ ] Output PII screening verified against new model or prompt configuration
- [ ] Cross-border transfer map updated if new endpoints are introduced

**When a data subject exercises a right:**
- [ ] Requestor identity confirmed before any disclosure or deletion
- [ ] Request logged with timestamp, right claimed, and case owner
- [ ] All systems containing this data subject's personal data identified
- [ ] Right executed across all identified systems within the legal timeframe
- [ ] Written response sent to the data subject confirming completion
- [ ] Case record with evidence retained for audit

## SelOps-Specific Context

SelOps operates AI systems for EU clients. GDPR is not optional — it is a legal obligation with fines of up to 4% of global annual turnover for serious violations. Privacy governance is a client deliverable, not an internal practice.

The most common gap SelOps faces is model output PII leakage. A model prompted with pseudonymised input can still reconstruct real names, addresses, or account numbers from context. Output screening is mandatory for any pipeline that ingests personal data, regardless of whether the input was masked.

DPIAs frequently catch privacy risks that architects miss — unexpected PII fields, retention periods that are never enforced, or cross-border transfers that no one documented. Treat the DPIA as a design review tool, not a compliance formality. Conduct it early, before the system is built.

Data subject rights requests arrive at inconvenient times. Build the capability to fulfil them before the system goes live. Discovering during a live erasure request that training data cannot be deleted because the model was fine-tuned on it — and fine-tuning is not reversible — is a serious incident. Assess fine-tuning privacy implications before any fine-tuning occurs.

[UNVERIFIED: ISO/IEC 42001 Annex A section A.7 covers data governance and quality controls relevant to personal data in AI systems — confirm clause numbering against the purchased ISO standard for authoritative clause IDs.]

## References

- GDPR (Regulation (EU) 2016/679) — Art. 5(1)(c) data minimisation, Art. 6 lawful basis, Art. 9 special categories, Art. 15-22 data subject rights, Art. 25 data protection by design and by default, Art. 35 DPIA
- ISO/IEC 42001:2023 — AI management system standard; data governance controls in Annex A, section A.7 [UNVERIFIED: confirm clause against purchased ISO standard]
- NIST AI 600-1 (2024) — Artificial Intelligence Risk Management Framework: Generative Artificial Intelligence Profile; the Data Privacy risk category (training-data memorisation, sensitive-attribute inference, PII leakage in outputs) [UNVERIFIED: specific action/subcategory identifiers — refer to the published profile for authoritative IDs]
