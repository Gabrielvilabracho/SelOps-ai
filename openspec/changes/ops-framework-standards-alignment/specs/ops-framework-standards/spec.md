# OPS Framework Standards Specification

## Purpose

Define the content, citation, categorization, and wiring requirements for aligning the SelOps OPS framework with external governance standards while preserving the existing 5-phase pipeline shape.

---

## Requirements

### Requirement: Citation Discipline

All standards-facing OPS skills MUST inherit a cross-cutting citation discipline. Any specific control, clause, article, tactic, or requirement identifier cited in skill content MUST either resolve to a published source or be marked inline as `[UNVERIFIED: <claim>]`. Each standards-facing **domain** skill MUST include a short `## References` section listing the standards, version/year, and cited identifiers. Pipeline phase skills surface their standards through inline citations and their `Cross-References` section rather than a `## References` section (see House-Style Conformance). Fabricated identifiers MUST NOT appear as fact in any skill type.

#### Scenario: Verified identifiers are reviewable
- GIVEN a reviewer reads any expanded or new standards-facing OPS skill
- WHEN the skill cites a specific control identifier
- THEN that identifier is either traceable in a `## References` section or marked `[UNVERIFIED: <claim>]`

#### Scenario: Uncertain identifiers are not stated as fact
- GIVEN the author cannot confirm a cited control ID
- WHEN the skill is written or updated
- THEN the uncertain claim is marked `[UNVERIFIED: <claim>]`
- AND no fabricated identifier appears as an unqualified citation

#### Scenario: Verification can be executed in review
- GIVEN the verify phase inspects the changed skill set
- WHEN it checks cited identifiers
- THEN every specific control ID either resolves to a published source or carries an `[UNVERIFIED]` marker

---

### Requirement: House-Style Conformance

The OPS framework has TWO distinct skill structures, and each modified or new skill MUST conform to the structure native to its TYPE:

- **Domain knowledge skills** (e.g. `ops-governance`, `ops-adversarial-security`, all 6 new skills, and the existing domain skills) MUST use the five-section domain house style: `When to Use`, `Core Principles`, `Patterns`, `Checklist`, and `SelOps-Specific Context`. Domain skills that cite standards MUST also include a `## References` section (per Citation Discipline).
- **Pipeline phase skills** (`ops-brief`, `ops-structure`, `ops-produce`, `ops-review`, `ops-deliver`) MUST use the pipeline-phase structure: `Role in the Pipeline`, `When to Use`, `Procedure`, `Inputs / Outputs`, `Gates & Escalation`, and `Cross-References`. This structure is intentional and shared across all five phases; a pipeline phase MUST NOT be converted to the domain five-section style, as that would break consistency with its sibling phases. A pipeline phase is NOT required to add a `## References` section; standards it references are cited inline (with `[UNVERIFIED]` discipline) and surfaced through its `Cross-References` section.

All instructions added by this change MUST produce verifiable review outcomes and MUST NOT introduce open-ended execution guidance, regardless of skill type.

#### Scenario: Expanded domain skills keep the domain house template
- GIVEN a reviewer opens any modified DOMAIN skill covered by this change
- WHEN the structure is inspected
- THEN the file contains the five required domain sections in the established house style

#### Scenario: Modified pipeline phases keep the pipeline-phase structure
- GIVEN a reviewer opens a modified PIPELINE PHASE skill covered by this change (e.g. `ops-produce`)
- WHEN the structure is inspected
- THEN the file retains the pipeline-phase structure (Role in the Pipeline / When to Use / Procedure / Inputs / Outputs / Gates & Escalation / Cross-References)
- AND it is NOT converted to the domain five-section style

#### Scenario: New skills use reviewable instructions
- GIVEN a reviewer reads a newly created OPS skill
- WHEN the patterns and checklist are inspected
- THEN each added instruction yields a checkable outcome
- AND no step is left as open-ended guidance without a verifiable result

---

### Requirement: Domain Categorization Preservation

The six new skills introduced by this change MUST be categorized as domain knowledge skills, not as new pipeline phases. The existing `brief → structure → produce → review → deliver` pipeline MUST remain unchanged.

#### Scenario: New skills are registered as domain knowledge
- GIVEN the new skill set is reviewed
- WHEN `inject.go`, presets, and skill files are compared
- THEN the six new skills are added to domain-skill registration paths
- AND no new pipeline phase is introduced

#### Scenario: Pipeline shape remains unchanged
- GIVEN the OPS pipeline assets are reviewed after the change
- WHEN phase ordering is inspected
- THEN the framework still exposes exactly five phases: `brief`, `structure`, `produce`, `review`, and `deliver`

---

### Requirement: Wiring Coherence

Every new `SkillID` registered in `internal/model/types.go`, `internal/catalog/skills.go`, `internal/components/sddops/inject.go`, or `internal/components/skills/presets.go` MUST correspond to a present `internal/assets/skills/<skill-id>/SKILL.md` file so install-time post-checks pass without orphan registrations.

#### Scenario: No orphan registration exists
- GIVEN a reviewer compares Go registrations against skill files
- WHEN each of the six new SkillIDs is traced
- THEN each registered skill has a matching `SKILL.md` file
- AND no Go registration points to a missing skill asset

---

### Requirement: Graduated Autonomy Standards Mapping

`ops-graduated-autonomy` MUST, in addition to **Citation Discipline**, document its academic foundations and map SelOps Suggest/Supervised/Autonomous levels to Sheridan-Verplank or Parasuraman-derived automation framing plus the EU AI Act Art.14.4 supervisory sub-articles required at each level.

#### Scenario: Autonomy levels map to EU supervision articles
- GIVEN a reviewer reads `ops-graduated-autonomy`
- WHEN the standards mapping content is inspected
- THEN Suggest, Supervised, and Autonomous are mapped to the required Art.14.4 sub-articles
- AND the foundational taxonomy is acknowledged with traceable or `[UNVERIFIED]` citations

---

### Requirement: Governance Standards Expansion

`ops-governance` MUST, in addition to **Citation Discipline**, add standards-grounded content for Model Inventory, third-party model governance, Effective Challenge, and regulated-industry audit-log expectations, grounded in NIST AI RMF GOVERN, SR 11-7, SP 800-53 AU-12, and EU AI Act Art.12.

#### Scenario: Governance includes challenge and inventory patterns
- GIVEN a reviewer reads `ops-governance`
- WHEN the added patterns and checklist are inspected
- THEN the skill defines an Effective Challenge protocol and a Model Inventory record
- AND it includes third-party model governance and audit-log checks with cited or `[UNVERIFIED]` identifiers

---

### Requirement: Observability Security and FinOps Expansion

`ops-observability` MUST, in addition to **Citation Discipline**, add security observability, FinOps, and environmental-impact guidance grounded in OWASP LLM Top 10, MITRE ATLAS, NIST AI 600-1, ISO 42001, and NIST MEASURE 1.1.

#### Scenario: Observability covers security, cost, and environment
- GIVEN a reviewer reads `ops-observability`
- WHEN the new principles, patterns, and checklist items are inspected
- THEN the skill includes security observability signals, per-request cost metrics, token-budget alerts, and environmental tracking guidance
- AND the cited standards are traceable or marked `[UNVERIFIED]`

---

### Requirement: Data Contracts Privacy Touchpoint

`ops-data-contracts` MUST, in addition to **Citation Discipline**, add a targeted privacy pattern for PII handling at contract boundaries, including explicit PII annotation, masking or tokenisation expectations, legal-basis documentation, and pre-downstream scanning for LLM outputs.

#### Scenario: Data contracts define PII boundary handling
- GIVEN a reviewer reads `ops-data-contracts`
- WHEN the added privacy pattern and checklist are inspected
- THEN the skill requires PII field identification and contract-boundary protections
- AND it references privacy controls with traceable or `[UNVERIFIED]` citations

---

### Requirement: Standard Documentation Regulatory Expansion

`ops-standard-documentation` MUST, in addition to **Citation Discipline**, add user-facing model card, FDA PCCP, and DPIA documentation requirements grounded in ISO 42001, EU AI Act Art.52, FDA PCCP guidance, and GDPR Art.35.

#### Scenario: Documentation covers users, healthcare, and privacy
- GIVEN a reviewer reads `ops-standard-documentation`
- WHEN the added patterns and handoff checklist are inspected
- THEN the skill requires a user-facing model card, healthcare PCCP guidance, and DPIA documentation triggers
- AND the cited standards are traceable or marked `[UNVERIFIED]`

---

### Requirement: Modular Architecture Supply-Chain Expansion

`ops-modular-architecture` MUST, in addition to **Citation Discipline**, add model supply-chain security guidance, including model SBOM/provenance records, version pinning, checksum verification, update gating, isolation, and fallback-provider expectations grounded in OWASP LLM05, MITRE ATLAS, NIST GOVERN 6.1, and NIST SP 800-161.

#### Scenario: Architecture includes model supply-chain controls
- GIVEN a reviewer reads `ops-modular-architecture`
- WHEN the added principles, patterns, and checklist are inspected
- THEN the skill requires model provenance or SBOM records plus concrete supply-chain controls
- AND the cited standards are traceable or marked `[UNVERIFIED]`

---

### Requirement: Adversarial Security Skill

The system MUST add `ops-adversarial-security` as a domain knowledge skill. In addition to **Citation Discipline**, the skill MUST be grounded in OWASP LLM Top 10 and MITRE ATLAS and MUST contain patterns for AI threat modeling, prompt-injection defense layers, excessive-agency prevention, red-team protocol, and poisoning detection.

#### Scenario: Adversarial security skill contains required sections and patterns
- GIVEN a reviewer opens `internal/assets/skills/ops-adversarial-security/SKILL.md`
- WHEN the file is inspected
- THEN it follows the five-section house style
- AND it includes the required standards grounding, patterns, and `## References`

---

### Requirement: Privacy Governance Skill

The system MUST add `ops-privacy-governance` as a domain knowledge skill. In addition to **Citation Discipline**, the skill MUST be grounded in GDPR, ISO 42001 A.7, and NIST AI 600-1 and MUST contain patterns for PII detection and masking, legal-basis register, DPIA, data-subject rights handling, and cross-border transfer controls.

#### Scenario: Privacy governance skill contains required sections and patterns
- GIVEN a reviewer opens `internal/assets/skills/ops-privacy-governance/SKILL.md`
- WHEN the file is inspected
- THEN it follows the five-section house style
- AND it includes the required privacy patterns and `## References`

---

### Requirement: Model Validation Skill

The system MUST add `ops-model-validation` as a domain knowledge skill. In addition to **Citation Discipline**, the skill MUST be grounded in NIST MEASURE, SR 11-7, and ISO 42001 verification or validation guidance and MUST contain patterns for validation-plan structure, conceptual soundness review, performance evaluation protocol, and re-validation triggers.

#### Scenario: Model validation skill contains required sections and patterns
- GIVEN a reviewer opens `internal/assets/skills/ops-model-validation/SKILL.md`
- WHEN the file is inspected
- THEN it follows the five-section house style
- AND it includes independent validation patterns plus `## References`

---

### Requirement: FinOps Governance Skill

The system MUST add `ops-finops-governance` as a domain knowledge skill. In addition to **Citation Discipline**, the skill MUST be grounded in NIST AI 600-1, ISO 42001 resource-planning guidance, and the FinOps Foundation framework and MUST contain patterns for token-budget enforcement, cost-versus-performance evaluation, operational cost projection, and environmental-impact reporting.

#### Scenario: FinOps governance skill contains required sections and patterns
- GIVEN a reviewer opens `internal/assets/skills/ops-finops-governance/SKILL.md`
- WHEN the file is inspected
- THEN it follows the five-section house style
- AND it includes the required FinOps patterns and `## References`

---

### Requirement: Transparency and Explainability Skill

The system MUST add `ops-transparency-explainability` as a domain knowledge skill. In addition to **Citation Discipline**, the skill MUST be grounded in EU AI Act Art.52 or Art.53, NIST transparency or explainability guidance, and ISO 42001 user-information guidance and MUST contain patterns for AI disclosure notices, decision explanation templates, automation-bias mitigation, and explainability audits.

#### Scenario: Transparency skill contains required sections and patterns
- GIVEN a reviewer opens `internal/assets/skills/ops-transparency-explainability/SKILL.md`
- WHEN the file is inspected
- THEN it follows the five-section house style
- AND it includes the required transparency patterns and `## References`

---

### Requirement: Model Lifecycle Skill

The system MUST add `ops-model-lifecycle` as a domain knowledge skill. In addition to **Citation Discipline**, the skill MUST be grounded in NIST GOVERN lifecycle controls and ISO 42001 lifecycle-record guidance and MUST contain patterns for model version records, retirement criteria, decommission procedure, and record-preservation policy.

#### Scenario: Lifecycle skill contains required sections and patterns
- GIVEN a reviewer opens `internal/assets/skills/ops-model-lifecycle/SKILL.md`
- WHEN the file is inspected
- THEN it follows the five-section house style
- AND it includes the required lifecycle patterns and `## References`

---

### Requirement: Real-Time Incident Detection in Produce

`ops-produce` MUST add a real-time incident-detection step that inspects post-state and step output for safety, compliance, PII, or data-exposure signals; classifies incidents using `ops-governance`; halts immediately on safety/compliance violation or data exposure; and escalates with evidence.

#### Scenario: Produce halts on critical incident signals
- GIVEN a reviewer reads `ops-produce`
- WHEN the execution procedure and gates are inspected
- THEN a real-time incident-detection step exists after post-state verification
- AND the step requires halt-plus-escalation for safety/compliance or data-exposure signals

---

### Requirement: Data-Contract Touchpoints Across Structure and Produce

`ops-structure` and `ops-produce` MUST add lightweight contract-boundary touchpoints so tasks that modify schemas or module interfaces explicitly flag the affected boundary and require contract validation in post-step verification.

#### Scenario: Structure and produce reference contract validation only when boundaries change
- GIVEN a reviewer reads `ops-structure` and `ops-produce`
- WHEN the affected-systems and post-step verification guidance is inspected
- THEN `ops-structure` flags contract or interface changes in the plan
- AND `ops-produce` requires the relevant contract test when a boundary-changing step executes

---

### Requirement: Reversibility Governance Foundation

`ops-governance` MUST add a reversibility-classification pattern defining reversible versus irreversible actions, required approvals, compensating controls where rollback is impossible, default-to-irreversible uncertainty handling, and binding use of that classification by pipeline phases, including `ops-brief`.

#### Scenario: Governance defines reversibility and brief can reference it
- GIVEN a reviewer reads `ops-governance` and `ops-brief`
- WHEN reversibility guidance is inspected
- THEN `ops-governance` defines a reversibility classification and rollback-governance pattern
- AND the concept is available as the domain foundation referenced by pipeline risk classification
