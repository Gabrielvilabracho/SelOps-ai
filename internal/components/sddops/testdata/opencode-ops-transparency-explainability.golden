---
name: ops-transparency-explainability
description: "SelOps AI transparency: disclosure notices and automation-bias controls. Trigger: When disclosing AI use, explaining decisions, or auditing explainability."
---

# Transparency and Explainability

## When to Use

Load this skill when deploying an AI system that interacts with end users and requires a disclosure notice, when a client or user asks why the AI produced a particular output and a decision explanation is needed, when designing safeguards against automation bias in a system where AI outputs influence human decisions, when conducting an explainability audit for a regulated deployment, or when a client's legal or compliance team requires documentation of AI transparency practices.

## Core Principles

- **AI disclosure is a legal requirement in regulated contexts.** Under the EU AI Act [UNVERIFIED: specific article number in the adopted Regulation (EU) 2024/1689 — transparency obligations for certain AI systems appear in Art.50 of the adopted text; verify against the published Official Journal version], operators of AI systems must inform users when they are interacting with an AI — not as a courtesy, but as a compliance obligation. For high-risk AI systems, additional transparency requirements apply: users must be informed of the system's capabilities and limitations. Failing to disclose AI use when required is a regulatory violation, not a design preference.
- **Explainability is context-dependent.** What counts as a sufficient explanation depends on the deployment context, the user's technical literacy, and the regulatory regime. A credit decision explanation must meet the GDPR Art.22 and sector-specific requirements. An AI-assisted recommendation in a consumer application requires a different level of detail. Define the explainability requirement before choosing an explainability technique — not the other way around.
- **Automation bias is a measurable risk.** Automation bias is the tendency of human operators to accept AI-generated outputs uncritically, even when those outputs are wrong. It is not a theoretical concern — it is documented in healthcare, aviation, criminal justice, and financial services. A system that presents AI outputs without clear uncertainty signals, without override mechanisms, and without regular audits of how often humans override AI recommendations is creating automation bias systematically. Mitigating automation bias is a design requirement, not a user-training problem.
- **Transparency documentation must be maintained.** A disclosure notice written at deployment and never updated does not satisfy ongoing transparency obligations. When the model changes, when the system's capabilities expand, or when the regulatory environment changes, the transparency documentation must be updated. Transparency docs are living artifacts, not one-time deliverables.

## Patterns

### AI Disclosure Notice
Every AI system that interacts with end users requires a disclosure notice. The notice must be: (1) visible before the user relies on any AI output — not buried in a terms of service document, (2) written in plain language appropriate for the intended user population, (3) specific about what the AI does and does not do (capabilities and limitations), (4) clear that the user is interacting with an AI, not a human. For high-risk AI systems under the EU AI Act [UNVERIFIED: exact article reference for high-risk system transparency requirements], the disclosure must also include information about the system's intended purpose, the operator responsible for the system, and how users can seek human review of AI decisions. Template elements: system name, brief description of what the AI does, what data it uses, what decisions or recommendations it makes or supports, known limitations, how to request human review or escalation, and contact for questions or concerns. The disclosure notice is reviewed at each model or capability update.

### Decision Explanation Template
When a user, client, or regulator requests an explanation of an AI decision or recommendation, use a structured template:
1. **What the AI decided or recommended**: state the output plainly (e.g., "The system flagged this transaction as high-risk").
2. **Key factors that influenced the output**: list the most influential inputs or signals. For interpretable models (logistic regression, decision trees), these are the feature contributions. For LLMs, these are the context signals and retrieval results used. For opaque models, use post-hoc explanation methods (SHAP, LIME [UNVERIFIED: specific tool names — use whichever is applicable to the model type]) and note that the explanation is approximate.
3. **Confidence or uncertainty**: state the model's confidence in the output and any relevant uncertainty estimate.
4. **What the AI cannot tell you**: be explicit about what is outside the model's scope, what data was not used, and what assumptions underlie the output.
5. **How to escalate**: tell the user how to request human review of the decision.
The explanation template is calibrated for the audience. A technical reviewer gets feature importance scores. An end user gets plain-language factors. Both get the escalation path.

### Automation-Bias Mitigation Protocol
Design AI systems to actively resist automation bias, not just tolerate it:
1. **Uncertainty signaling**: every AI output surfaced to a human decision-maker must carry a visible uncertainty or confidence signal. High-confidence outputs are marked differently from low-confidence outputs. Do not display all outputs with equal visual weight.
2. **Override mechanism**: every AI recommendation must have a clearly accessible, frictionless override mechanism. Overriding the AI must not require more steps than accepting it. Override rates are logged.
3. **Audit of override rates**: regularly review the rate at which human operators override AI recommendations. A very low override rate (< 1%) in a system where the AI is not perfectly accurate may indicate automation bias — humans are not checking outputs. A very high override rate may indicate the AI adds no value. Both extremes require investigation.
4. **Decision support framing**: present AI outputs as inputs to human judgment, not as conclusions. Use language like "the AI suggests…" rather than "the answer is…". Avoid UX patterns that imply the AI has decided.
5. **Human-in-the-loop gates for high-stakes outputs**: for irreversible or high-stakes outputs (see ops-governance Reversibility Classification), require explicit human confirmation before the AI-recommended action is taken — regardless of the autonomy level registered for the system.

### Explainability Audit
Conduct an explainability audit before go-live for any high-risk or regulated AI deployment, and at each major model update. The audit covers:
1. **Disclosure coverage**: verify that all user-facing touchpoints display the AI disclosure notice as required.
2. **Explanation quality**: generate explanations for a representative sample of outputs and assess whether the explanations are accurate, understandable, and complete (using the Decision Explanation Template).
3. **Automation-bias controls**: verify that uncertainty signals, override mechanisms, and decision-support framing are implemented and functioning.
4. **Override rate baseline**: record the baseline override rate at deployment for future comparison.
5. **Documentation currency**: verify that the disclosure notice and transparency documentation reflect the current system capabilities and model version.
NIST AI RMF MEASURE [UNVERIFIED: MEASURE 2.9 is cited in the literature as covering explainability and transparency evaluation — verify the exact subcategory number] recommends measuring explainability as part of AI risk management. ISO/IEC 42001:2023 Clause A.8.2 [UNVERIFIED: exact clause number] addresses user information and transparency obligations for AI systems. Both require that transparency be operationalized, not just documented.

## Checklist

**At deployment:**
- [ ] AI disclosure notice written, reviewed, and placed at all required user-facing touchpoints
- [ ] Disclosure notice reviewed by legal or compliance for regulatory adequacy
- [ ] Decision Explanation Template defined and tested for the deployment context
- [ ] Automation-bias controls implemented: uncertainty signals, override mechanism, decision-support framing
- [ ] Explainability audit completed before go-live (high-risk or regulated deployments)
- [ ] Override rate baseline recorded

**At each model or capability update:**
- [ ] AI disclosure notice updated to reflect current capabilities and limitations
- [ ] Transparency documentation version-controlled and change log maintained
- [ ] Explanation quality verified for updated model outputs (sample audit)

**Ongoing:**
- [ ] Override rate reviewed at each engagement milestone
- [ ] Automation-bias audit findings documented and remediation tracked
- [ ] For regulated deployments: transparency documentation available for regulatory review on request

**For EU AI Act regulated deployments:**
- [ ] AI disclosure notice satisfies applicable transparency obligations [UNVERIFIED: confirm applicable article(s) in Regulation (EU) 2024/1689 against the Official Journal version]
- [ ] High-risk AI transparency requirements documented and implemented

## SelOps-Specific Context

Regulated clients require disclosure. When operating AI systems for clients in finance, healthcare, legal, or public-sector domains, the transparency obligation does not belong only to the client — SelOps, as the operator, shares responsibility for ensuring the systems it deploys meet the transparency requirements of the applicable regulatory framework. This means SelOps must understand the applicable transparency rules before deployment, not discover them during a regulatory review.

Automation bias is especially consequential in financial services (credit decisioning, fraud flagging) and healthcare (clinical decision support) where AI recommendations directly influence high-stakes decisions. SelOps operators working in these sectors must actively design against automation bias — not assume users will exercise critical judgment in the absence of structural safeguards.

Explainability is not a feature — it is a system property. It must be designed in from the architecture phase (see ops-structure), not retrofitted after go-live. Post-hoc explanation methods applied to a black-box production model are harder to validate, harder to audit, and less reliable than explanation approaches designed into the system from the start.

## References

- **EU AI Act — Regulation (EU) 2024/1689**: Transparency obligations for certain AI systems. [UNVERIFIED: Art.50 is cited as the transparency article in the adopted text — verify article numbering against the Official Journal publication before citing definitively.] High-risk AI system requirements include informing users about AI use, capabilities, and limitations.
- **NIST AI RMF 1.0 (2023)** — MEASURE function. MEASURE 2.9 [UNVERIFIED: exact subcategory number] addresses explainability and interpretability evaluation as part of AI risk measurement. The GOVERN and MAP functions also address transparency obligations.
- **ISO/IEC 42001:2023** — Clause A.8.2 [UNVERIFIED: exact clause number] covers user information and transparency obligations for AI system operators, including disclosure of AI use and explanation of AI decisions.
- **GDPR — Regulation (EU) 2016/679**, Art.22: Rights related to automated individual decision-making, including the right to explanation. Applies where AI decisions produce legal or similarly significant effects on individuals.
