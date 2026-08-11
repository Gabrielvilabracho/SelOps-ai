# Exploration: ops-framework-standards-alignment

## Current State

The SelOps OPS framework has 11 skill files: 6 domain knowledge skills and 5 pipeline phase agents, all under `internal/assets/skills/ops-*/SKILL.md`. All files follow a consistent house style with five sections: When to Use / Core Principles / Patterns / Checklist / SelOps-Specific Context.

The framework is conceptually sound (FASE 0 validated by two independent research efforts). Its gaps are gaps of *omission*, not contradiction — it does not cite foundational standards it was already aligned with, and it is missing 6 entirely absent domains (adversarial security, privacy governance, model validation, FinOps, transparency/explainability, model lifecycle).

The registry is clean: `internal/model/types.go` has `SkillID` constants, `internal/catalog/skills.go` has catalog entries, `internal/components/sddops/inject.go` has two injection lists (`sddOpsSkillIDs` + `opsPipelineSkillIDs`), and `internal/components/skills/presets.go` has the `opsSkills` slice used by `PresetSelOpsOperational`.

---

## Affected Areas

### Skill files (content changes — the PRIMARY deliverable)
- `internal/assets/skills/ops-graduated-autonomy/SKILL.md` — ADAPT (cite academic roots + EU AI Act Art.14.4)
- `internal/assets/skills/ops-governance/SKILL.md` — EXPAND (add SR 11-7 requirements + NIST GOVERN 6.1)
- `internal/assets/skills/ops-observability/SKILL.md` — EXPAND (add security observability + FinOps + environmental metrics)
- `internal/assets/skills/ops-data-contracts/SKILL.md` — KEEP core + add PII note
- `internal/assets/skills/ops-standard-documentation/SKILL.md` — EXPAND (model cards for users, FDA PCCP, privacy docs)
- `internal/assets/skills/ops-modular-architecture/SKILL.md` — EXPAND (supply-chain security, model SBOM)
- `internal/assets/skills/ops-adversarial-security/SKILL.md` — CREATE (P0)
- `internal/assets/skills/ops-privacy-governance/SKILL.md` — CREATE (P0)
- `internal/assets/skills/ops-model-validation/SKILL.md` — CREATE (P1)
- `internal/assets/skills/ops-finops-governance/SKILL.md` — CREATE (P2)
- `internal/assets/skills/ops-transparency-explainability/SKILL.md` — CREATE (P2)
- `internal/assets/skills/ops-model-lifecycle/SKILL.md` — CREATE (P3)

### Registry/wiring files (Go code — mechanical additions only)
- `internal/model/types.go` — 6 new `SkillID` constants
- `internal/catalog/skills.go` — 6 new catalog entries
- `internal/components/sddops/inject.go` — 6 new entries in `sddOpsSkillIDs`
- `internal/components/skills/presets.go` — 6 new entries in `opsSkills`

---

## Per-Existing-Skill Delta Plans

### 1. `ops-graduated-autonomy` — ADAPT

**Target standard**: Sheridan & Verplank (1978) / Parasuraman, Sheridan & Wickens (2000) levels of automation; HITL/HOTL/HOOTL taxonomy; EU AI Act Art. 14.4(a)-(e).

**Section-by-section changes:**

**Core Principles** — Add one new principle:
> "Graduated autonomy is a formal adaptation of the Sheridan-Verplank/Parasuraman-Sheridan-Wickens levels of automation (Parasuraman et al., 2000, IEEE TSMC 30:3). The SelOps three-level model maps as: Suggest = HITL / Levels 1–4 (computer offers alternatives, human decides); Supervised = HOTL / Levels 5–7 (computer executes, necessarily informs human); Autonomous = HOOTL / Levels 8–10 (computer decides and informs only if asked). The per-engagement + per-environment framing and the Scope Verification Protocol are SelOps-original additions. Cite these when the framework is used in regulated contexts."

**Patterns** — Add new pattern: `EU AI Act Art.14.4 Compliance Mapping`:
> Map the current autonomy level to the EU AI Act Art.14.4 supervisory mechanism required:
> - Art.14.4(a) — understand capabilities/limitations: documented in engagement record → required at ALL levels
> - Art.14.4(b) — detect/avoid automation bias: explicit operator training note + Scope Verification Protocol → required at Supervised and Autonomous
> - Art.14.4(c) — interpret output correctly: model card + structured output logging → required at Supervised and Autonomous  
> - Art.14.4(d) — decide not to use / disregard output: overrule mechanism documented → required at Autonomous
> - Art.14.4(e) — interrupt via stop procedure to safe state: rollback procedure = the stop mechanism → required at Autonomous
>
> For EU-based clients or systems classified as high-risk under the EU AI Act Annex III, confirm all applicable sub-articles are satisfied before setting Autonomous level.

**Checklist** — Add 2 items to "Before starting any task":
- `[ ]` For EU-regulated systems: Art.14.4(a)-(e) compliance mapped to current autonomy level
- `[ ]` Foundational taxonomy (Sheridan/Parasuraman + HITL/HOTL/HOOTL) acknowledged in engagement record if client requests standards traceability

**SelOps-Specific Context** — Add one paragraph:
> The academic foundations matter for regulated engagements. When a financial or healthcare client's compliance team asks "on what basis does this autonomy model work?", the answer is not "it's how SelOps does it" — it is "it is a software-operations adaptation of the Parasuraman-Sheridan-Wickens (2000) levels of automation, with per-engagement configurability added as a practical extension." This framing is respected and verifiable.

---

### 2. `ops-governance` — EXPAND

**Target standards**: NIST AI RMF GOVERN 2.1, GOVERN 4.3, GOVERN 6.1; ISO 42001 A.3.2; SR 11-7 (Effective Challenge, Model Inventory, Independent Validation); SP 800-53 AU-12; EU AI Act Art.12.

**Section-by-section changes:**

**Core Principles** — Add 2 new principles:
1. > "Model Inventory is a governance requirement, not an audit artifact. Every model in production for a client must be listed in the model inventory: model ID, version, provider, purpose, risk classification, deployment date, owner, and review schedule. This applies to both SelOps-deployed models and client-owned models that SelOps operates." (SR 11-7 / NIST GOVERN 1.6)
2. > "Third-party model governance extends your governance obligations. When using foundation models from external providers (OpenAI, Anthropic, AWS Bedrock, etc.), apply NIST AI RMF GOVERN 6.1 criteria: document the provider, the model version, what due diligence was performed, what the fallback is if the provider changes the model, and what contractual protections exist. Outsourcing the inference does not outsource the governance risk." (NIST GOVERN 6.1 / ISO 42001 A.10.3)

**Patterns** — Add 2 new patterns:

*Effective Challenge Protocol (SR 11-7)*:
> Before any model goes to production or any significant model update ships, an Effective Challenge must be performed. Effective Challenge is an objective, critical analysis of the model by someone who did not build it. At minimum: (1) review the conceptual soundness — does the model approach the problem in a valid way? (2) review the data and training methodology if accessible, (3) run the model against adversarial inputs and edge cases prepared by the challenger, (4) document the challenger's findings and the model team's response. The challenger must have sufficient technical expertise and no development stake in the outcome. For client engagements in finance (SR 11-7), healthcare (FDA), or legal (EU AI Act), Effective Challenge by an independent party is required before production.

*Model Inventory Record (NIST GOVERN 1.6 / SR 11-7)*:
> Maintain a Model Inventory entry for every model in production. Each entry includes: model identifier (name, version, provider), deployment environment, purpose and use case, risk classification (using the engagement's risk framework), deployment date, approved-by (with timestamp), review schedule (at minimum: after each update, annually otherwise), known limitations and biases, and decommission date (if planned). The inventory is a living document — not a one-time registry. It is reviewed at each engagement milestone.

**Checklist** — Add to "When setting up a new regulated-industry engagement":
- `[ ]` Model Inventory created with entries for all models to be deployed
- `[ ]` Effective Challenge scheduled (identify the challenger before go-live, not after)
- `[ ]` Third-party model governance documented per NIST GOVERN 6.1 (provider, version, fallback, contractual protections)
- `[ ]` Audit log format verified against SP 800-53 AU-12 requirements (generation, protection, retention)
- `[ ]` For EU high-risk systems: verify audit log satisfies EU AI Act Art.12 (automatic logging requirements)

**SelOps-Specific Context** — Add one paragraph:
> The SR 11-7 / SR 26-2 Effective Challenge requirement deserves emphasis. In banking and financial services, an AI model that affects credit decisions, risk scoring, or fraud detection cannot go to production without independent validation and Effective Challenge. SelOps cannot be both the builder and the only reviewer. Identify an independent challenger — whether internal (a SelOps team member with no development stake) or external — and document their findings as part of the governance record. "We tested it and it works" is not Effective Challenge.

---

### 3. `ops-observability` — EXPAND

**Target standards**: OWASP LLM Top 10 (2025): LLM01-LLM10; MITRE ATLAS (16 tactics, 84 techniques); NIST AI 600-1 (environmental impacts, confabulation); ISO 42001 A.6.2.6; NIST MEASURE 1.1.

**Section-by-section changes:**

**Core Principles** — Add 2 new principles:
1. > "Security observability is distinct from reliability observability. Track anomalous input patterns (prompt injection signals, unusual token velocity, repeated adversarial probes) alongside standard quality metrics. A system can have 99.9% uptime and zero latency alerts while being actively attacked through its prompt interface. OWASP LLM01 (Prompt Injection) and MITRE ATLAS provide the threat vocabulary; observability is the detection layer." (OWASP LLM01; MITRE ATLAS AML.T0051)
2. > "Cost and environmental impact are observable system properties. Token consumption is not just a quality signal — it is a cost signal and an environmental signal. Track actual cost per request, cost per pipeline stage, and cumulative cost per engagement. For production systems, track or estimate compute carbon footprint when clients have ESG reporting requirements." (NIST AI 600-1 Environmental Impacts; ISO 42001 A.4)

**Patterns** — Add 3 new patterns:

*Security Observability Signals (OWASP LLM Top 10 / MITRE ATLAS)*:
> Instrument the following security signals alongside standard metrics:
> - **LLM01 / Prompt Injection**: log input length outliers, inputs containing instruction-override patterns (common jailbreak phrases), and inputs that cause unexpected output format shifts. Set an alert threshold for prompt injection pattern density.
> - **LLM06 / Excessive Agency**: log every out-of-scope tool call, every permission escalation attempt, and every action that would exceed the registered autonomy level. These should never occur — even one is an alert-worthy event.
> - **LLM10 / Model Extraction**: log query rate per client/session, flag sessions with systematic input variation (adversarial probing patterns), and alert on sustained high query volume from a single source.
> - **MITRE ATLAS AML.TA0013 (Exfiltration)**: log unusual output lengths, outputs containing unexpected data structures, and outputs that reference system internals.
> Route security signals to a dedicated security dashboard, not just the standard ops dashboard.

*FinOps Metrics Set*:
> Track per-request cost metrics: estimated USD cost (input tokens × price + output tokens × price per model), cost per pipeline stage, running total per engagement per billing period. Define a token budget per request type (documented in engagement record). Alert at 80% of budget (warning) and 100% (critical). Compare cost-vs-quality tradeoffs when evaluating model upgrades: a 20% quality improvement at 3× the cost requires explicit client approval. Log cost alongside quality metrics so degradation and cost can be correlated.

*Environmental Impact Tracking (NIST AI 600-1)*:
> For engagements where clients have ESG or sustainability reporting requirements: (1) instrument total compute hours per model (GPU/TPU hours if trackable via provider APIs), (2) record model efficiency (tokens/second, tasks completed per GPU-hour), (3) include environmental impact summary in engagement reports alongside cost. If direct measurement is not available, document the model provider's carbon reporting methodology and reference it. This is a reporting requirement, not a blocker — but it must be addressed before handoff to clients with ESG obligations.

**Checklist** — Add to "When adding a new AI component to production":
- `[ ]` Security observability signals instrumented (prompt injection patterns, excessive agency attempts, query rate anomalies)
- `[ ]` Cost metrics tracked per request and per pipeline stage
- `[ ]` Token budget defined and alert configured
- `[ ]` Environmental impact tracking in place if client has ESG requirements

---

### 4. `ops-data-contracts` — KEEP + PII NOTE

**Target standards**: GDPR Art.5(1)(c) data minimisation; GDPR Art.25 data protection by design.

**Section-by-section changes:**

This skill is well-structured. One targeted addition only.

**Patterns** — Add one new pattern after "LLM Output Contract":

*PII in Contract Boundaries (GDPR Art.5 / Art.25)*:
> When a data contract carries or may carry personally identifiable information (names, emails, IDs, financial data, health data), apply data protection by design at the contract definition stage — not as a post-processing filter:
> - Mark PII fields explicitly in the schema with a `pii: true` annotation (or equivalent).
> - Apply field-level masking or tokenisation at the producer side before writing to the contract payload.
> - Document the legal basis for processing this PII in the contract's data lineage document.
> - LLM output contracts are especially susceptible to inadvertent PII: if training data or retrieved context includes PII, the model may reproduce it in output. Add a PII scan step before the LLM output contract is passed to downstream consumers.
> For detailed privacy controls, see `ops-privacy-governance` (when available).

**Checklist** — Add one item to "Before defining a contract":
- `[ ]` PII fields identified and annotated; legal basis for processing documented if applicable

---

### 5. `ops-standard-documentation` — EXPAND

**Target standards**: NIST AI RMF Model Cards (emerging practice); ISO 42001 A.8.2 (information for users); FDA PCCP (Marketing Submission Recommendations, 2024); GDPR Art.35 DPIA.

**Section-by-section changes:**

**Core Principles** — Add one new principle:
> "Model cards serve two audiences: internal governance and external users. The internal model card (already required) documents training data, bias evaluation, and failure modes for the governance record. A separate user-facing model card must be produced for any system where users interact with AI output: what this AI does, what it cannot do, known error patterns, how to report concerns, and who is accountable. This is required under ISO 42001 A.8.2 and EU AI Act Art.52." (ISO 42001 A.8.2; EU AI Act Art.52)

**Patterns** — Add 3 new patterns:

*User-Facing Model Card (ISO 42001 A.8.2 / EU AI Act Art.52)*:
> The user-facing model card is a separate, non-technical document. It contains: (1) what the AI system does in plain language, (2) what decisions it makes vs. what decisions humans make, (3) known limitations and error scenarios the user should watch for, (4) how to report incorrect outputs, (5) who is accountable (SelOps, the client, the model provider — document the chain), (6) date of last update and how to check for newer versions. This document is part of the client handoff package. For EU clients deploying the system to end users, it satisfies the EU AI Act Art.52 transparency obligation.

*Predetermined Change Control Plan (FDA PCCP)*:
> For AI systems used in healthcare or medical decision support: document a PCCP before deployment. The PCCP must specify: (1) the scope of changes permitted without a new regulatory submission (model updates, prompt changes, training data updates), (2) the validation methodology for each permitted change type, (3) the performance criteria that must be satisfied before any change is deployed, (4) the monitoring plan that detects performance degradation post-change. The PCCP is submitted to FDA as part of the marketing submission for AI/ML-based SaMD. For non-FDA contexts, the PCCP structure is still the best practice for governing model updates in safety-critical systems.

*Privacy/DPIA Documentation (GDPR Art.35)*:
> For systems processing personal data at scale, or processing sensitive categories (health, financial, location, biometric), produce a Data Protection Impact Assessment before deployment. The DPIA documents: (1) the nature, scope, context, and purposes of processing, (2) necessity and proportionality assessment, (3) risks to data subjects, (4) measures to address risks. The DPIA is not a one-time exercise — re-assess when the system changes materially. Store the DPIA alongside the model card and governance records. For detailed privacy controls, see `ops-privacy-governance` (when available).

**Checklist** — Add to "Before client handoff":
- `[ ]` User-facing model card produced (separate from internal governance model card)
- `[ ]` For healthcare systems: PCCP documented and reviewed
- `[ ]` For systems processing PII at scale: DPIA completed and stored with governance records

---

### 6. `ops-modular-architecture` — EXPAND

**Target standards**: OWASP LLM05 (Supply Chain Vulnerabilities); MITRE ATLAS AML.TA0003 (Initial Access — Supply Chain Compromise); NIST AI RMF GOVERN 6.1; NIST SP 800-161 (supply chain risk management).

**Section-by-section changes:**

**Core Principles** — Add one new principle:
> "Model supply chain is an attack surface. A third-party model or a downloaded model checkpoint is a supply chain dependency. Apply the same integrity controls you apply to code dependencies: verify checksums, pin versions, document provenance (who produced it, under what license, via what training data), and define a fallback for provider discontinuation. An unverified model update is an untested dependency injection." (OWASP LLM05; MITRE ATLAS AML.TA0003)

**Patterns** — Add 2 new patterns:

*Model SBOM and Provenance Record (NIST GOVERN 6.1 / OWASP LLM05)*:
> For every model used in production, maintain a Model Software Bill of Materials (SBOM) entry: model name and version, provider, source URL or artifact registry, SHA-256 checksum of the model artifact (where available), license type, training data summary (source, cutoff, known gaps), known vulnerabilities or issues (referenced by CVE or provider advisory), and SelOps risk assessment. The Model SBOM is part of the Model Inventory (see ops-governance). Update it on every model version change. For open-source or fine-tuned models, document the base model lineage (what was the base model, who fine-tuned it, on what data).

*Supply Chain Security Controls (OWASP LLM05 / MITRE ATLAS AML.TA0003)*:
> Apply these controls for every external model integration:
> - **Version pinning**: never use floating version references (e.g., `gpt-4-latest`) in production. Pin to a specific version identifier.
> - **Integrity verification**: when downloading model artifacts, verify checksums against provider-published values.
> - **Update gating**: model provider updates (even minor ones, including silent updates from API-served models) go through the change approval workflow in ops-governance before taking effect in production.
> - **Isolation**: the inference module (already required by ops-modular-architecture) also acts as the supply chain boundary — it contains the blast radius of a supply chain compromise to the inference layer.
> - **Provider fallback**: document the fallback provider or model before going to production. If the primary provider is unavailable or changes the model, the fallback is pre-validated.

**Checklist** — Add to "When designing a new component":
- `[ ]` Model SBOM entry created for every external model used
- `[ ]` Model versions pinned (no floating version references in production config)
- `[ ]` Checksum verification in place for downloadable model artifacts
- `[ ]` Provider fallback identified and documented

---

## New Skill Outlines (6 Skills)

All new skills follow the established house style exactly: frontmatter → h1 title → When to Use → Core Principles → Patterns → Checklist → SelOps-Specific Context.

---

### NEW: `ops-adversarial-security` (P0 — CRITICAL)

**File**: `internal/assets/skills/ops-adversarial-security/SKILL.md`

**Frontmatter**:
```yaml
name: ops-adversarial-security
description: "SelOps adversarial security for AI systems. Trigger: When threat modeling AI systems, implementing prompt injection defenses, or conducting red-team assessments."
```

**Foundational standards**: OWASP LLM Top 10 (2025) LLM01–LLM10; MITRE ATLAS (16 tactics, 84 techniques, esp. AML.T0051 Prompt Injection, AML.TA0013 Exfiltration, AML.TA0003 Supply Chain Compromise, AML.TA0008 Defense Evasion, AML.T0020 Training Data Poisoning, AML.T0056.000 LLM Jailbreak).

**When to Use**: Load this skill when designing security controls for an AI system, conducting threat modeling, implementing prompt injection defenses, planning red-team exercises, reviewing an AI system before production deployment, or responding to a suspected adversarial attack on an AI pipeline.

**Core Principles (5 principles)**:
1. AI systems have a distinct threat model from traditional software. Prompt injection (OWASP LLM01) is not SQL injection — it exploits the model's instruction-following behavior, not a code vulnerability. Traditional WAFs do not stop it.
2. Excessive agency (OWASP LLM06) amplifies every other vulnerability. An AI agent with tool-call access to production systems that is successfully jailbroken has direct production access. Minimize permissions before hardening prompts.
3. Red-teaming is not optional for production AI systems. NIST AI 600-1 and NIST GOVERN 4.3 require periodic adversarial testing. Schedule it before go-live and after every significant model or prompt change.
4. Defense in depth: no single control stops all LLM attacks. Layer input validation, output filtering, access controls, monitoring, and rate limiting. Assume prompt injection attempts will occur.
5. Model extraction (OWASP LLM10) is an intellectual property and confidentiality risk. An attacker who reconstructs a client's proprietary fine-tuned model has stolen a business asset.

**Patterns (5 patterns)**:
1. *Threat Model Template for AI Systems (MITRE ATLAS)*: for each AI system, produce a threat model covering: (a) entry points (API, UI, batch input, retrieved context); (b) top 3 MITRE ATLAS tactics applicable to this system; (c) specific OWASP LLM risks by system type; (d) mitigations mapped to each threat; (e) residual risks accepted and documented.
2. *Prompt Injection Defense Layers (OWASP LLM01)*: (a) input validation — reject or flag inputs containing system-override patterns before reaching the model; (b) prompt construction — separate system instructions from user content; never concatenate them directly; (c) output filtering — validate model output structure before passing downstream; reject outputs that reference system internals; (d) output sandboxing — do not execute model-suggested code without human review at Suggest level; (e) monitoring — log prompt injection signals (see ops-observability security signals pattern).
3. *Excessive Agency Prevention (OWASP LLM06)*: apply least-privilege to tool calls and system access: (a) enumerate the exact tools/permissions the model needs — nothing more; (b) enforce permissions at the tool layer, not only in the prompt; (c) require human approval (Suggest level) before any irreversible tool call from an AI agent; (d) log all tool calls with actor, tool, parameters, and outcome; (e) alert on any tool call outside the pre-approved list.
4. *Red-Team Exercise Protocol (NIST AI 600-1 / NIST GOVERN 4.3)*: before production and after major changes: (a) define scope — which OWASP LLM risks to test; (b) assign a challenger who did not build the system; (c) run jailbreak attempts, prompt injection probes, model extraction probes, and data exfiltration attempts; (d) document findings categorized by OWASP LLM risk; (e) require PASS on critical risks before deployment; (f) store the red-team report with governance records.
5. *Data Poisoning Detection (OWASP LLM03 / MITRE ATLAS AML.T0020)*: for systems that use fine-tuning or retrieval-augmented generation: (a) maintain chain-of-custody records for training/RAG data sources; (b) validate data quality and integrity before ingestion; (c) monitor for behavioral anomalies post-training that suggest poisoning (outputs that deviate from baseline on known-good inputs); (d) apply differential testing before and after data updates.

**Checklist**: Before production deployment (red-team passed, threat model complete, injection defenses layered, tool permissions minimized, monitoring in place); During operations (prompt injection signals monitored, tool call log reviewed, model extraction alerts active); After any incident (classify per OWASP LLM risk, update threat model, add detection to monitoring).

**SelOps-Specific Context**: AI consultancies are high-value targets for adversarial attacks because they operate in client production environments with broad tool access. A successfully jailbroken SelOps agent can exfiltrate client data, corrupt production systems, and execute changes outside approved scope — all at machine speed. The blast radius is not a sandbox. Treat every AI component with tool access as a security boundary.

---

### NEW: `ops-privacy-governance` (P0 — CRITICAL)

**File**: `internal/assets/skills/ops-privacy-governance/SKILL.md`

**Frontmatter**:
```yaml
name: ops-privacy-governance
description: "SelOps privacy governance for AI systems. Trigger: When handling personal data, implementing PII controls, conducting DPIAs, or building GDPR-compliant AI pipelines."
```

**Foundational standards**: GDPR Art.5 (principles), Art.6-7 (lawful basis/consent), Art.13-14 (information obligations), Art.15-22 (data subject rights), Art.25 (data protection by design), Art.35 (DPIA); ISO 42001 A.7 (data for AI systems); NIST AI 600-1 (data privacy).

**When to Use**: Load this skill when building or auditing AI systems that process personal data, when designing PII detection or masking pipelines, when a client operates in the EU or under GDPR-equivalent regulation, when producing a DPIA, or when data subject rights requests must be handled.

**Core Principles (5 principles)**:
1. Data minimisation is a design constraint, not an audit step (GDPR Art.5(1)(c)). The question "do we need this personal data field?" must be answered before the schema is defined, not after the system is built.
2. LLMs are PII amplifiers. A model trained on or given context containing PII may reproduce it in unexpected outputs. PII must be detected and masked at ingestion, before it reaches the model.
3. Legal basis must be documented before processing begins (GDPR Art.6). "We have the data" is not a legal basis. For each PII category processed, document: which GDPR Art.6 basis applies, who made the determination, and when.
4. Data subject rights (GDPR Art.15-22) apply to AI-processed data. If a user's data was used to train or evaluate a model, they have the right to access, correct, and erase it. AI systems must support these operations.
5. Privacy by design (GDPR Art.25) is a MUST for AI systems, not a SHOULD. Default to the most privacy-protective configuration. Require explicit opt-in for any processing beyond the stated purpose.

**Patterns (5 patterns)**:
1. *PII Detection and Masking Pipeline*: before any user-originated data reaches an LLM: (a) run PII detection (named entity recognition for names, emails, IDs, financial data, health data, location); (b) apply field-level masking or tokenisation (replace with consistent pseudonyms for context preservation); (c) log the original → masked mapping in an access-controlled store (needed for data subject rights); (d) configure the model to not reproduce masked tokens in outputs; (e) run a post-output PII scan before downstream delivery.
2. *Legal Basis Register (GDPR Art.6-7)*: maintain a processing register per engagement: data category, purpose, legal basis (consent / contract / legal obligation / vital interests / public task / legitimate interests), retention period, and deletion mechanism. For consent-based processing: document consent record, withdrawal mechanism, and what happens to data on withdrawal. Review the register at each engagement milestone.
3. *Data Protection Impact Assessment (DPIA, GDPR Art.35)*: required when processing is likely to result in high risk to individuals (systematic profiling, sensitive data at scale, automated decision-making). DPIA structure: (1) describe the processing (nature, scope, context, purpose); (2) assess necessity and proportionality; (3) identify and assess risks to data subjects; (4) identify measures to address risks; (5) consult the client's Data Protection Officer if required. Store DPIA with governance records. Re-run on material system changes.
4. *Data Subject Rights Handling (GDPR Art.15-22)*: document the mechanism for each right before going to production: access (Art.15) — how to retrieve all data processed about a subject; rectification (Art.16) — how to correct it; erasure (Art.17) — how to delete it, including from model training data if applicable; portability (Art.20) — how to export it in machine-readable format; objection (Art.21) — how to stop processing. For AI systems that have trained on personal data, deletion from training data is technically complex — document this limitation and the compensating control (e.g., model retraining).
5. *Cross-Border Transfer Controls (GDPR Chapter V)*: for AI systems where data crosses EU borders (e.g., US-based LLM provider processing EU data): document the transfer mechanism (Standard Contractual Clauses, Adequacy Decision, Binding Corporate Rules), verify the provider has signed SCCs or equivalent, and record this in the data lineage document.

**Checklist**: Before deployment (PII fields identified and annotated in schema, legal basis documented for each PII category, PII detection/masking pipeline implemented, DPIA completed if high-risk, data subject rights mechanism documented); During operations (PII scan on outputs active, consent withdrawals processed within 72 hours, processing register current); At engagement close (data retention period enforced, deletion procedure documented and tested, DPIA updated if system changed materially).

**SelOps-Specific Context**: EU clients assume GDPR compliance. When SelOps builds an AI system that processes EU residents' data — directly or through a client's system — GDPR applies to SelOps as a data processor. The data processing agreement (DPA) between SelOps and the client is a legal requirement, not a formality. Verify it exists before the first byte of personal data is processed.

---

### NEW: `ops-model-validation` (P1)

**File**: `internal/assets/skills/ops-model-validation/SKILL.md`

**Frontmatter**:
```yaml
name: ops-model-validation
description: "SelOps model validation framework. Trigger: When validating a model before production deployment, conducting independent validation, or applying SR 11-7 validation requirements."
```

**Foundational standards**: NIST AI RMF MEASURE 2 (MEASURE 2.1, 2.3, 2.4, 2.9); SR 11-7 Effective Challenge, Conceptual Soundness, Back-Testing, Outcome Analysis, Independent Validation; ISO 42001 A.6.2.4 (AI system verification and validation).

**When to Use**: Load this skill when validating a model before production deployment, when conducting independent model review, when a client's regulatory environment (finance, healthcare, government) requires formal model validation, or when evaluating a model against a benchmark or performance standard.

**Core Principles (4 principles)**:
1. Validation is independent from development. The validator did not build the model. This is not a convention — it is the foundational requirement of SR 11-7 and NIST MEASURE 2. Self-validation is not validation.
2. Conceptual soundness is tested before performance. A model that is technically accurate but solving the wrong problem is a worse outcome than a model with known performance gaps. Validate the approach before validating the numbers.
3. Back-testing and out-of-sample testing are distinct. Training performance tells you nothing about production performance. Require held-out test sets, temporal out-of-sample tests, and stress tests with adversarial inputs.
4. Validation is an ongoing activity, not a pre-deployment gate. Redeploy the validation protocol when: the model is updated, the data distribution shifts, usage patterns change materially, or performance thresholds are triggered by monitoring alerts.

**Patterns (4 patterns)**:
1. *Validation Plan Structure (SR 11-7 / NIST MEASURE 2.1)*: before beginning validation, produce a validation plan: model purpose and scope, validation team (independent — document who), evaluation datasets (source, size, temporal coverage, known gaps), metrics and acceptance thresholds (defined before running — not after), stress test scenarios, and sign-off requirements. The plan is approved before validation begins.
2. *Conceptual Soundness Review (SR 11-7)*: independent reviewer assesses: (a) is the problem formulation correct? (b) is the modelling approach appropriate for the problem and data? (c) are there known theoretical limitations of this approach in this domain? (d) are the training data assumptions valid? Document findings as PASS/FAIL with evidence for each question.
3. *Performance Evaluation Protocol (NIST MEASURE 2.4 / MEASURE 2.9)*: run three evaluation suites: (a) in-sample benchmark — standard metrics on the training/development set; (b) held-out test set — same metrics on data the model has never seen; (c) adversarial/stress test — edge cases, distributional shifts, inputs designed to probe failure modes. Acceptance requires passing all three suites above the pre-defined thresholds.
4. *Ongoing Monitoring Triggers (NIST MEASURE 2.3)*: define the conditions that trigger a re-validation cycle: performance metrics drop below threshold (from ops-observability), data distribution shifts detected, model version update, usage pattern changes by >20%, or any safety/compliance incident. Re-validation is not a full re-deployment — it is the specific validation protocol applied to the specific change.

**Checklist**: Before validation (validation plan approved, independent validator confirmed, evaluation datasets documented, acceptance thresholds set); During validation (conceptual soundness reviewed, all three evaluation suites run, findings documented); Before production (validation report signed off by independent validator, results stored in governance records, ongoing monitoring triggers configured); After model update (re-validation triggered per the defined trigger conditions).

---

### NEW: `ops-finops-governance` (P2)

**File**: `internal/assets/skills/ops-finops-governance/SKILL.md`

**Frontmatter**:
```yaml
name: ops-finops-governance
description: "SelOps FinOps governance for AI systems. Trigger: When managing AI compute costs, defining token budgets, or optimizing cost-vs-performance tradeoffs."
```

**Foundational standards**: NIST AI 600-1 (Environmental Impacts, compute); ISO 42001 A.4 (Resources — planning and managing AI system resources); FinOps Foundation framework.

**When to Use**: Load this skill when defining cost governance for an AI system, when token spend is approaching or exceeding budget, when evaluating model upgrade cost/benefit tradeoffs, when producing cost projections for a client engagement, or when environmental/ESG reporting is required.

**Core Principles (4 principles)**:
1. Cost governance is governance. Uncontrolled LLM spend is a production risk. A runaway agent loop, a prompt that triggers excessive output tokens, or an unexpected traffic spike can produce an invoice that exceeds the engagement value. Token budgets are not optimisations — they are safety controls.
2. Cost and performance are coupled metrics. Every evaluation of model quality must be paired with its cost. "This model is better" is incomplete. "This model is 15% more accurate at 2× the cost" is a decision.
3. Environmental impact is a reportable metric for ESG-compliant clients (NIST AI 600-1). If a client has sustainability reporting obligations, compute carbon footprint is a required output of the engagement, not an optional extra.
4. Cost projections must precede production commitments. Build the cost model at design time: estimated tokens per request × expected request volume × model price per token × buffer. Present this to the client before go-live.

**Patterns (4 patterns)**:
1. *Token Budget Definition and Enforcement*: for each request type in the AI system: (a) measure baseline token consumption (input + output) on a representative sample; (b) define a per-request budget with a 20% safety buffer; (c) enforce the budget at the inference layer — truncate input, warn on output budget approach, hard-stop on budget breach; (d) log actual vs. budget per request; (e) alert at 80% of daily budget consumption.
2. *Cost-vs-Performance Evaluation Matrix*: when evaluating model options, produce this matrix: model name / accuracy on benchmark / p95 latency / cost per 1K tokens / cost per task (average) / cost at production volume (projected) / recommendation. The client makes the cost/performance tradeoff decision explicitly, with this data. Never choose the most expensive model "because it's best" without this analysis.
3. *Operational Cost Projection (ISO 42001 A.4)*: before production: (a) model type (API-served vs. self-hosted), (b) projected request volume (p50 and p99), (c) token consumption per request (input + output, p50 and p99), (d) total monthly cost at p50 and p99 volume, (e) cost anomaly threshold (the spend level that triggers an incident response), (f) cost scaling plan (what happens if volume grows 10×). This projection is a deliverable, not an internal estimate.
4. *Environmental Impact Report (NIST AI 600-1)*: for clients with ESG obligations: (a) document the model provider's energy efficiency and carbon reporting methodology; (b) estimate compute hours per month (via provider metrics if available); (c) estimate carbon footprint using provider's disclosed emissions factor; (d) include in quarterly engagement reports; (e) evaluate lower-carbon model alternatives when ESG reporting is required.

**Checklist**: Before deployment (token budget defined per request type, cost projection produced and approved by client, cost anomaly threshold set, environmental impact methodology documented if required); During operations (actual spend vs. budget tracked, cost anomaly alerts active, monthly cost report produced); At model update (cost impact of model change assessed before go-live, cost projection revised if significant change).

---

### NEW: `ops-transparency-explainability` (P2)

**File**: `internal/assets/skills/ops-transparency-explainability/SKILL.md`

**Frontmatter**:
```yaml
name: ops-transparency-explainability
description: "SelOps AI transparency and explainability. Trigger: When implementing AI disclosure requirements, producing explanations for AI decisions, or mitigating automation bias."
```

**Foundational standards**: EU AI Act Art.52 (transparency obligations for certain AI systems); EU AI Act Art.53 (transparency obligations for general-purpose AI); NIST AI RMF MEASURE 2.9 (explainability); NIST AI 600-1 (GV-1.2-001 transparency policies); ISO 42001 A.8.2 (system documentation and information for users).

**When to Use**: Load this skill when deploying AI systems to end users (not just operators), when a client requires AI disclosure notices, when producing explanations for AI-driven decisions, when evaluating automation bias risk, or when a regulated context requires explainability documentation.

**Core Principles (4 principles)**:
1. Users have a right to know they are interacting with AI (EU AI Act Art.52). This is not a UX preference — it is a legal requirement in the EU and an ethical baseline globally. Any system that could be mistaken for a human must disclose AI involvement.
2. Explainability is context-dependent. "The model said so" is never an explanation. The explanation must be actionable: what factors drove this output, what confidence level applies, what the known limitations are. The depth of explanation scales with the stakes of the decision.
3. Automation bias is a measurable risk (EU AI Act Art.14.4(b)). When humans consistently accept AI output without verification, the AI has effectively lost its oversight. Design workflows to interrupt automation bias: require active confirmation on high-stakes outputs, surface confidence intervals, and make disagreement with the AI the default-presented option when stakes are high.
4. Transparency documentation must be maintained, not just created. A model card written at deployment and never updated is a compliance risk. Assign an owner and a review schedule.

**Patterns (4 patterns)**:
1. *AI Disclosure Notice (EU AI Act Art.52)*: for any user-facing AI system: produce a disclosure notice that states (a) this system uses AI to produce [specific outputs], (b) AI outputs may contain errors — describe the known error types, (c) this system is operated by [SelOps / client name] and the AI component is powered by [provider if disclosable], (d) how to report incorrect outputs, (e) what human review applies to AI outputs. Display this notice at first use and provide ongoing access. For EU deployments: the disclosure must satisfy Art.52 requirements before go-live.
2. *Decision Explanation Template (NIST MEASURE 2.9)*: for AI systems that drive decisions affecting people (recommendations, classifications, rankings): produce an explanation template per output type: (a) what the AI concluded; (b) the top factors that drove this conclusion (feature importance, retrieval relevance, or equivalent); (c) confidence indicator (high/medium/low with definition); (d) what the AI did not consider (known blind spots for this model on this input type); (e) recommended human review step. The template is filled per output and surfaced to the human reviewer.
3. *Automation Bias Mitigation Protocol (EU AI Act Art.14.4(b))*: for high-stakes AI outputs: (a) present confidence intervals alongside the AI recommendation — not just the recommendation; (b) require active confirmation rather than passive acceptance (make "I reviewed and agree" require a click, not just reading); (c) for sequential decisions, display disagreement rate (how often reviewers have overridden AI for this output type) to calibrate trust; (d) schedule periodic blind-review exercises where reviewers evaluate outputs without knowing the AI's answer, to calibrate automation bias.
4. *Explainability Audit (NIST AI 600-1 / ISO 42001 A.8.2)*: annually or after any major model change: (a) sample 50 production outputs, (b) evaluate whether the explanation for each output is accurate and actionable, (c) test whether end users understand the explanations (brief usability test), (d) document findings and remediation. Store audit results with governance records.

**Checklist**: Before deployment (AI disclosure notice produced and approved, decision explanation template defined, automation bias mitigations designed into workflow); During operations (explanation accuracy monitored, user feedback on explanations collected, automation bias signals tracked); Annually (explainability audit conducted, disclosure notice reviewed and updated if system changed).

---

### NEW: `ops-model-lifecycle` (P3)

**File**: `internal/assets/skills/ops-model-lifecycle/SKILL.md`

**Frontmatter**:
```yaml
name: ops-model-lifecycle
description: "SelOps model lifecycle management. Trigger: When versioning models, planning model retirement, managing decommissioning, or maintaining model inventory across a production engagement."
```

**Foundational standards**: NIST AI RMF GOVERN 1.6 (AI risk management policies and processes — versioning); NIST AI RMF GOVERN 1.7 (processes for AI system retirement and decommissioning); ISO 42001 A.6.2.3 (AI system design documentation — versioning); ISO 42001 A.6.2.8 (AI system recording / lifecycle records).

**When to Use**: Load this skill when planning a model version upgrade, when retiring a model from production, when designing version management for a new AI system, or when an engagement close requires lifecycle documentation handoff.

**Core Principles (4 principles)**:
1. Every production model has a documented lifecycle (NIST GOVERN 1.6). Deployment without a retirement plan is operational debt. Define the retirement criteria before the model goes to production, not when you want to replace it.
2. Model versioning is not the same as software versioning. A model update may change system behaviour in ways that cannot be enumerated. Version control for models requires behavioural baselines, not just artifact hashes.
3. Decommissioning is a planned operation, not a shutdown (NIST GOVERN 1.7). A model decommission must: notify affected users, migrate dependent systems, archive the model artifact and its governance records, and verify no dangling references remain in production.
4. Record preservation is a regulatory requirement in most sectors. Audit logs, model cards, validation reports, and governance records for a decommissioned model must be retained for the engagement's regulatory retention period — not deleted at decommission.

**Patterns (4 patterns)**:
1. *Model Version Control Record (NIST GOVERN 1.6 / ISO 42001 A.6.2.3)*: for every model in production, maintain a version record: current version identifier, deployment date, previous version (and retirement date), behavioural baseline (the benchmark results at deployment, used to detect drift), change log (what changed from previous version and why), and next planned review date. The version record is part of the Model Inventory (see ops-governance).
2. *Model Retirement Criteria Definition*: before going to production, document the retirement criteria: performance drop thresholds that trigger retirement consideration, support end-date from the model provider, scheduled re-evaluation date, business events that could trigger early retirement (client contract end, regulatory change, provider discontinuation). Review criteria at each engagement milestone.
3. *Decommission Procedure (NIST GOVERN 1.7)*: when a model is decommissioned: (a) notify all downstream consumers (internal and client-facing) with a deprecation window of at least 30 days; (b) execute the migration plan (validate the replacement model through the ops-model-validation protocol before traffic switches); (c) archive the model artifact and all governance records (model card, validation reports, audit logs, governance records) to the designated long-term store; (d) verify no production system references the retired model version after switchover; (e) close the Model Inventory entry with decommission date and archive location.
4. *Record Preservation Policy (ISO 42001 A.6.2.8)*: define retention periods at engagement start: operational audit logs (minimum: the longer of 3 years or the client's regulatory requirement), model validation reports (minimum: duration of model deployment + 1 year), governance records (minimum: 5 years for regulated industries — confirm with client legal). At engagement close, transfer the record archive to the client with a documented chain of custody. SelOps retains a copy for the warranty period.

**Checklist**: Before deployment (lifecycle record created, retirement criteria defined, record retention periods documented); During deployment (version record updated on each model change, behavioural baseline refreshed, retirement criteria reviewed at each milestone); At decommission (deprecation notice sent with 30-day window, migration plan validated, archive complete, Model Inventory entry closed).

---

## Registry / Wiring Impact

**The following 6 new `SkillID` constants must be added to `internal/model/types.go`**, in the SelOps operational skill IDs block (lines 86–91), maintaining the `SkillOps*` naming convention:

```go
SkillOpsAdversarialSecurity     SkillID = "ops-adversarial-security"
SkillOpsPrivacyGovernance       SkillID = "ops-privacy-governance"
SkillOpsModelValidation         SkillID = "ops-model-validation"
SkillOpsFinOpsGovernance        SkillID = "ops-finops-governance"
SkillOpsTransparencyExplainability SkillID = "ops-transparency-explainability"
SkillOpsModelLifecycle          SkillID = "ops-model-lifecycle"
```

**The following 6 entries must be added to `internal/catalog/skills.go`**, in the SelOps operational skills block (lines 26–31), with `Category: "operational"` and `Priority: "p0"` for P0 skills, `"p1"` for P1, `"p2"` for P2, `"p3"` for P3:

```go
{ID: model.SkillOpsAdversarialSecurity,      Name: "ops-adversarial-security",       Category: "operational", Priority: "p0"},
{ID: model.SkillOpsPrivacyGovernance,        Name: "ops-privacy-governance",          Category: "operational", Priority: "p0"},
{ID: model.SkillOpsModelValidation,          Name: "ops-model-validation",            Category: "operational", Priority: "p1"},
{ID: model.SkillOpsFinOpsGovernance,         Name: "ops-finops-governance",           Category: "operational", Priority: "p2"},
{ID: model.SkillOpsTransparencyExplainability, Name: "ops-transparency-explainability", Category: "operational", Priority: "p2"},
{ID: model.SkillOpsModelLifecycle,           Name: "ops-model-lifecycle",             Category: "operational", Priority: "p3"},
```

**The following 6 entries must be added to `sddOpsSkillIDs` in `internal/components/sddops/inject.go`** (lines 38–45), all new domain knowledge skills:

```go
"ops-adversarial-security",
"ops-privacy-governance",
"ops-model-validation",
"ops-finops-governance",
"ops-transparency-explainability",
"ops-model-lifecycle",
```

**The following 6 entries must be added to `opsSkills` in `internal/components/skills/presets.go`** (lines 21–28):

```go
model.SkillOpsAdversarialSecurity,
model.SkillOpsPrivacyGovernance,
model.SkillOpsModelValidation,
model.SkillOpsFinOpsGovernance,
model.SkillOpsTransparencyExplainability,
model.SkillOpsModelLifecycle,
```

**No changes are needed to `opsPipelineSkillIDs` or `opsPipelineSkills`** — all 6 new skills are DOMAIN KNOWLEDGE (see categorization below), not new pipeline phases.

---

## Coherence Gap Resolutions

### Gap 1: Incident classification defined in governance but no pipeline phase executes it in real time

**Current state**: `ops-governance` defines 4 incident categories (output quality degradation / safety violation / service unavailability / data exposure). `ops-review` references the classification system (step 5 of its procedure) but only applies it after execution completes, not during execution.

**Problem**: A safety violation or data exposure event that occurs mid-execution is not caught until `ops-review`. By that point, `ops-produce` has already completed all steps. The classification happens too late to affect the response.

**Resolution**: Add a real-time incident detection step to `ops-produce`. In `ops-produce`'s procedure, after step 4 (verify post-state), add step 4a:
> "Scan the post-state and step output for incident signals: unexpected PII in output, outputs matching safety violation patterns, data exposure indicators. If any incident signal is detected, classify it immediately using the ops-governance four-category system. For safety/compliance violation or data exposure: halt immediately, do not proceed to the next step, and escalate to the operator with the classification and evidence. For output quality degradation: log the signal with severity, continue if within acceptable threshold, halt if above threshold."

This makes `ops-produce` the real-time incident detector, and `ops-review` the post-execution incident confirmer/escalator. The two roles are distinct and complementary.

**No new phase is needed** — the fix is a targeted addition to `ops-produce`'s procedure and a cross-reference addition to `ops-produce` → `ops-governance`.

---

### Gap 2: ops-data-contracts + ops-modular-architecture are design domains no pipeline phase operationalizes

**Current state**: Both skills define design-time principles (schema-first, module boundaries). The pipeline phases reference them only loosely. `ops-structure` mentions "typed contracts" without requiring data-contract validation. `ops-produce` mentions "logging" without requiring LLM output contract validation. No phase explicitly checks schema compliance.

**Problem analysis**: This is **structurally acceptable** for the current scope. The pipeline phases operate systems that have already been designed; they do not redesign them. Data contract validation and module boundary enforcement are *pre-pipeline* concerns (design, build, and CI gates), not execution-pipeline concerns. Operationalizing them in the execution pipeline would add friction to every operation, most of which do not change schemas or module boundaries.

**Resolution**: The gap is real but the correct fix is lightweight:
- `ops-structure` step 6 (map affected systems) should explicitly note: "If this task modifies a data contract boundary or module interface, flag this in the execution plan. The post-step verification for that step must include schema validation."
- `ops-produce` step 4 (verify post-state) should note: "If a step modifies a contract boundary, run the contract test for the affected schema as the post-state verification."

This adds the touchpoint without making every execution plan a schema audit. Add these notes as targeted additions to `ops-structure` and `ops-produce` in the same PR as the domain skill updates.

---

### Gap 3: brief's "reversibility" concept lacks a domain-skill foundation

**Current state**: `ops-brief` classifies risk along 5 dimensions, one of which is "reversibility (reversible vs. irreversible)." This concept is defined and used throughout the pipeline (`ops-structure` step 3, `ops-graduated-autonomy` rollback requirements) but is not grounded in a domain skill. No skill explains what reversibility means, how to determine it, or what the regulatory/governance implications are.

**Resolution**: Ground reversibility in `ops-governance`. It fits there because reversibility is a risk governance concept — not a technical concept. Add a pattern to `ops-governance`:

*Reversibility Classification and Rollback Governance*:
> Classify every production action as reversible (can return to prior state in ≤15 minutes with documented steps) or irreversible (cannot be undone — data deletion, external system notifications, financial transactions, regulatory filings). Irreversible actions require: (a) explicit approval at the appropriate autonomy level, (b) a documented compensating control if rollback is not possible (what is the containment action?), (c) a post-execution audit record noting irreversibility. The reversibility classification is made at the structure phase and is binding — it cannot be downgraded during execution. When uncertain, default to irreversible.

This grounds the concept in ops-governance where the audit/approval logic already lives, and gives ops-brief a domain reference to cite.

---

## Domain vs. Pipeline Categorization

**All 6 new skills are DOMAIN KNOWLEDGE skills, not new pipeline phases.**

Rationale:
- Domain skills encode WHAT the operator knows: principles, patterns, standards, and checklists for a specific concern.
- Pipeline skills encode HOW to execute: sequenced procedures with gates, inputs/outputs, and escalation rules.
- None of the 6 new skills introduce new execution phases. They introduce new bodies of knowledge that existing phases (brief, structure, produce, review, deliver) apply.

Specific categorization:
| Skill | Category | Injection list |
|---|---|---|
| `ops-adversarial-security` | DOMAIN | `sddOpsSkillIDs` / `opsSkills` |
| `ops-privacy-governance` | DOMAIN | `sddOpsSkillIDs` / `opsSkills` |
| `ops-model-validation` | DOMAIN | `sddOpsSkillIDs` / `opsSkills` |
| `ops-finops-governance` | DOMAIN | `sddOpsSkillIDs` / `opsSkills` |
| `ops-transparency-explainability` | DOMAIN | `sddOpsSkillIDs` / `opsSkills` |
| `ops-model-lifecycle` | DOMAIN | `sddOpsSkillIDs` / `opsSkills` |

Cross-phase usage (for reference):
- `ops-adversarial-security` → referenced from brief (threat model), structure (security steps), produce (injection detection), review (security gate)
- `ops-privacy-governance` → referenced from brief (PII classification), structure (PII controls), produce (PII scanning), deliver (DPIA handoff)
- `ops-model-validation` → referenced from governance (Effective Challenge, pre-deployment gate)
- `ops-finops-governance` → referenced from brief (cost classification), observability (cost metrics)
- `ops-transparency-explainability` → referenced from standard-documentation (model cards), deliver (disclosure handoff)
- `ops-model-lifecycle` → referenced from governance (model inventory), standard-documentation (lifecycle docs)

---

## Sizing / Slicing Forecast

**Total work**: 6 new SKILL.md files + 6 expanded SKILL.md files + 4 Go file edits (mechanical) + 2 targeted pipeline phase additions (ops-produce + ops-structure) = ~14 file changes.

**Line estimate**:
- Each new skill: ~80–100 lines (matching house style average: 68–79 lines observed)
- Each expanded skill: +15–30 lines of net additions
- 4 Go registry files: ~30 lines total (mechanical constants and slice entries)
- ops-produce + ops-structure targeted additions: ~15 lines total
- Total: approximately 600–700 lines of new/changed content

**Review budget constraint**: 400-line PR limit → **this CANNOT ship as one PR.**

**Recommended slice boundary — by priority tier**:

| PR | Scope | Approx. lines | Rationale |
|---|---|---|---|
| **PR-1: P0 new skills + P0 registry** | Create `ops-adversarial-security` + `ops-privacy-governance` + their 4 Go registry entries + sddOpsSkillIDs additions | ~250 lines | Critical security/privacy gaps ship first. Go registry changes are mechanical. |
| **PR-2: P0-P1 expansions** | Expand `ops-graduated-autonomy` + `ops-governance` (which includes the reversibility coherence fix) + `ops-observability` + `ops-data-contracts`; add incident detection addition to `ops-produce` | ~300 lines | The highest-priority existing skill updates. Governance expansion resolves 2 coherence gaps. |
| **PR-3: P1-P2 new skills** | Create `ops-model-validation` + `ops-finops-governance` + `ops-transparency-explainability` + their 6 registry entries | ~350 lines | Mid-priority new skills. |
| **PR-4: P2-P3 expansions + remaining** | Expand `ops-standard-documentation` + `ops-modular-architecture`; create `ops-model-lifecycle`; add ops-structure data-contract touchpoint | ~300 lines | Lower-priority expansions and the lifecycle skill. |

Alternative slice: domain-vs-pipeline. But since all new skills are domain skills and pipeline changes are minimal (2 targeted additions), the priority-tier split is cleaner.

---

## Dependency Note: complete-sddops

This change is **INDEPENDENT** of `complete-sddops` (paused at PR1). The ops skill files in `internal/assets/skills/` are content files — they do not depend on the Go infrastructure in `complete-sddops`.

**However, note for future sessions**: `complete-sddops` will eventually inject these skills into agent installations. When `complete-sddops` resumes:
- The 6 new SkillID constants added here must be referenced by `complete-sddops`'s sub-agents or overlay configuration if those sub-agents are expected to know about the new skills.
- The `sddOpsSkillIDs` list in `inject.go` is already the canonical injection list — when `complete-sddops` resumes, it will pick up the new entries automatically (assuming it reads from the same list).
- No blocking coupling today. Flag for the `complete-sddops` resume session: "6 new domain skills were added in `ops-framework-standards-alignment`. Verify that the ops-orchestrator system prompt references the expanded skill set."

---

## Recommendation

Proceed to **proposal** and then **spec** phases. The scope is fully defined, the per-skill delta plans are concrete enough for spec writers, the new-skill outlines are structured enough for apply-phase execution, and the registry changes are mechanical. No open questions remain — the user locked scope in the FASE 0 decision.

Suggest beginning with PR-1 (P0 new skills) because these address regulatory-blocking gaps (GDPR, OWASP security) that affect existing client engagements immediately.

## Risks

- **House style consistency**: 6 new skill files written by different apply-phase executions may drift from the house style. Mitigation: enforce the existing 5-section template strictly; use existing skill files as reference.
- **Cross-reference accuracy**: new skills add cross-references to each other (`ops-privacy-governance` → `ops-data-contracts`, `ops-adversarial-security` → `ops-observability`). These references must be correct at time of writing. Mitigation: write skills in priority order; later skills can reference earlier ones.
- **Registry consistency**: 4 Go files must be updated in the same PR. If split across PRs, the injection list will reference SkillID constants that don't exist yet. Mitigation: always include all 4 Go registry files in the same PR as the skills they register.
- **complete-sddops coupling**: flagged above — not blocking today, must be flagged at complete-sddops resume.

## Ready for Proposal

Yes. Scope is locked, delta plans are concrete, slicing is defined, coherence resolutions are specific. The proposal phase can proceed directly.
