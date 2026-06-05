---
name: ops-model-validation
description: "SelOps AI model validation: independent review and ongoing monitoring. Trigger: When validating AI models or assessing conceptual soundness."
---

# Model Validation

## When to Use

Load this skill when an AI model is being considered for production deployment, when a model update requires validation before it ships, when a client in a regulated industry (finance, healthcare) requires independent validation evidence, when designing a validation plan for a new model or use case, or when setting up ongoing monitoring protocols for a deployed model.

## Core Principles

- **Validation must be independent from development.** The team that built the model cannot be the sole validators of that model. At minimum, the validator must have no development stake in the outcome. For finance (SR 11-7) and high-risk AI under the EU AI Act, independent validation by a separate function or external party is required — not optional.
- **Conceptual soundness precedes performance metrics.** Before measuring accuracy, AUC, or F1, confirm the model's approach is conceptually valid for the problem domain. A model can achieve high benchmark scores on the wrong framing. SR 11-7 calls this the conceptual soundness review — it asks: "Does the model make sense for this problem, given the data and assumptions used?" Performance metrics are evaluated after conceptual soundness is confirmed.
- **Back-testing performance is not training performance.** A model validated only on the data it was trained on (or tuned on) provides no real validation signal. Validation requires a held-out test set the model has never seen, plus adversarial/stress inputs that explore boundary conditions and failure modes. The gap between in-sample and out-of-sample performance is a key validation signal.
- **Ongoing monitoring is not optional.** Validation at deployment is the starting gate, not the finish line. Models degrade over time as the data distribution shifts, as user behavior changes, or as the environment evolves. Ongoing monitoring defines triggers that require re-validation — not arbitrary schedules.
- **Model Inventory is the anchor.** Every validated model must have an entry in the Model Inventory (see ops-governance). The validation record is attached to the inventory entry: validation date, validator identity, findings, approval status, and re-validation schedule.

## Patterns

### Validation Plan Structure
Every model validation engagement begins with a written Validation Plan before any validation activity starts. The plan includes: (1) model identifier and version being validated, (2) intended use and deployment context, (3) validation scope — what is in scope and what is explicitly out of scope, (4) validation approach and methods, (5) data sources for validation (held-out set, adversarial set, benchmark set), (6) performance thresholds and acceptance criteria, (7) validator identity and independence statement, (8) timeline and review milestones. The Validation Plan is approved before validation begins. Changes to scope require an amendment, not a verbal agreement.

### Conceptual Soundness Review (SR 11-7 Effective Challenge)
The Conceptual Soundness Review asks whether the model's design is valid for the stated problem before examining its performance. Steps: (1) review the problem framing — is this the right type of model for this task? (2) review input variables and feature engineering — are the chosen features causally or statistically appropriate for the outcome, or are they proxies that introduce bias? (3) review training data — is it representative of the deployment population? Does it contain historical biases that the model will amplify? (4) review the modelling assumptions — does the mathematical structure (linear, tree-based, neural, etc.) match the data generating process? (5) document findings and the reviewer's judgment on soundness. A model with conceptual soundness gaps must not proceed to performance evaluation until those gaps are resolved or explicitly accepted with documented rationale.

### Performance Evaluation Protocol
After conceptual soundness is confirmed, evaluate performance across three data regimes:
- **In-sample evaluation**: run against training data to understand the baseline. High in-sample performance alone is not a validation signal — it confirms the model learned the training set.
- **Held-out test set evaluation**: run against data the model has never seen. Report primary metrics (accuracy, precision, recall, AUC, RMSE, etc. as appropriate for the model type) plus calibration, fairness metrics across demographic groups if applicable, and confidence distribution. The held-out set must not have been used for hyperparameter tuning.
- **Adversarial and stress evaluation**: probe boundary conditions, edge cases, out-of-distribution inputs, and adversarial inputs designed to surface failure modes. Report performance under these conditions and assess whether degradation is acceptable. For LLMs: include prompt-injection probes, context-length stress, and domain-boundary tests.

Document all three sets of results and compare them explicitly. A large gap between in-sample and held-out performance signals overfitting and requires investigation before approval.

### Ongoing Monitoring Triggers
Define at deployment time the conditions that require re-validation. Standard triggers include:
- **Data distribution shift**: statistical tests (PSI, KS test, or equivalent) detect that the live data distribution has drifted significantly from the training distribution.
- **Performance degradation**: primary performance metrics drop below the defined threshold on a rolling window of live predictions.
- **Business context change**: the intended use, target population, or regulatory classification of the model changes.
- **Provider or version change**: the model is updated by its provider, fine-tuned, or replaced with a different version. Version pinning must prevent silent updates (see ops-modular-architecture).
- **Post-incident trigger**: a governance incident (see ops-governance) classified as output quality degradation or safety/compliance violation involving this model automatically triggers re-validation.
Re-validation uses the same validation plan structure. It does not require repeating the full conceptual soundness review unless the model architecture or training approach changed; it does require held-out and stress evaluation with updated data.

## Checklist

**Validation Plan:**
- [ ] Validation Plan written and approved before any validation activity starts
- [ ] Validator identity and independence statement documented
- [ ] Acceptance criteria (performance thresholds) defined before evaluation begins
- [ ] Validation scope (in/out) explicitly stated

**Conceptual Soundness:**
- [ ] Problem framing reviewed and confirmed appropriate for the model type
- [ ] Input features reviewed for causal/statistical appropriateness
- [ ] Training data reviewed for representativeness and historical bias
- [ ] Modelling assumptions reviewed and documented
- [ ] Conceptual soundness finding recorded before proceeding to performance evaluation

**Performance Evaluation:**
- [ ] In-sample, held-out, and adversarial/stress evaluation all completed
- [ ] Fairness and calibration metrics reported where applicable
- [ ] In-sample vs. held-out gap documented and assessed
- [ ] Adversarial failure modes documented and evaluated for acceptability

**Ongoing Monitoring:**
- [ ] Re-validation triggers defined at deployment time
- [ ] Monitoring metrics and thresholds configured before production launch
- [ ] Model Inventory entry updated with validation record, findings, and re-validation schedule

**Regulated Engagements:**
- [ ] For finance (SR 11-7): independent validation by a party with no development stake
- [ ] Validation record retained and available for regulatory review

## SelOps-Specific Context

Independent validation is a deliverable for finance clients. SR 11-7 requires that model risk management include independent validation — defined as evaluation by staff not involved in model development. "Self-validated" or "developer-tested" does not satisfy this requirement for models used in credit, risk scoring, fraud detection, or similar applications. SelOps must identify the independent validator before development begins, not after the model is ready.

For LLM-based agents operated by SelOps, validation must cover the full pipeline, not just the language model in isolation. An LLM that performs well in isolation can fail when integrated with a retrieval system, a tool set, or a prompt template. The validation scope must include integrated-system tests that mirror the production configuration.

Validation findings are a client asset and may be requested by the client's own compliance or legal team. Write findings in language a non-technical compliance reviewer can act on — not just metric tables.

## References

- **SR 11-7 (2011)** — Federal Reserve / OCC Supervisory Guidance on Model Risk Management. Foundational source for Conceptual Soundness, Effective Challenge, Independent Validation, ongoing monitoring, and model inventory requirements in banking and financial services.
- **NIST AI RMF 1.0 (2023)** — MEASURE function. MEASURE 2.x subcategories [UNVERIFIED: exact MEASURE 2.1, 2.3, 2.9 numbers] cover AI system testing, evaluation, verification, and monitoring, including ongoing monitoring protocols.
- **ISO/IEC 42001:2023** — Clause A.6.2.4 [UNVERIFIED: exact clause number] addresses verification and validation of AI system components including requirements for structured testing and documentation of validation outcomes.
