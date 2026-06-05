---
name: ops-finops-governance
description: "SelOps FinOps governance: token budgets and cost-performance tradeoffs. Trigger: When managing AI system costs, token budgets, or environmental impact."
---

# FinOps Governance

## When to Use

Load this skill when designing token budget policies for a production AI system, when evaluating the cost-performance tradeoff of model choices for a client engagement, when producing an operational cost projection before go-live approval, when a client asks for environmental impact reporting on AI system usage, or when a cost overrun has occurred and root-cause analysis is needed.

## Core Principles

- **Token budgets are safety controls, not optimizations.** A token budget that limits context length, output size, or total tokens per session is not a performance tuning parameter. It is a safety boundary that prevents runaway compute costs, protects against prompt-injection amplification attacks, and enables predictable unit economics. Treat budget violations as incidents, not configuration drift.
- **Cost and performance are coupled, not separable.** Choosing a smaller model because it is cheaper without measuring the performance impact is not FinOps — it is gambling. The correct decision framework is the cost-versus-performance evaluation matrix: measure actual task performance at each cost tier and choose the tier where the cost reduction does not materially degrade the outcome. "Materially degrade" is defined by the client's acceptable performance threshold, documented before the evaluation.
- **Environmental impact is reportable.** AI system compute has a measurable energy and carbon footprint. NIST AI 600-1 (2024) identifies environmental impact as a risk dimension for generative AI systems. For clients with sustainability reporting obligations or ESG commitments, environmental impact of AI system usage must be tracked and reported. This is not optional for those clients — it is a contractual deliverable.
- **Cost projections must precede production approval.** No AI system proceeds to production without a documented operational cost projection that covers: expected token consumption per request and per day, infrastructure cost, cost at peak load, and cost at failure (runaway requests, retry storms). The projection must be reviewed and approved alongside the technical design.
- **Cost visibility enables accountability.** Per-engagement, per-client, and per-model cost attribution must be trackable. Cost that cannot be attributed cannot be governed. Instrument cost metrics at the request level from day one; retrofitting cost attribution onto a running system is significantly harder.

## Patterns

### Token Budget Definition and Enforcement
Define token budgets as explicit, documented constraints — not implicit assumptions. For each AI system or workflow, document: (1) maximum context window tokens, (2) maximum output tokens per request, (3) maximum tokens per session (if applicable), (4) maximum tokens per day per user or per client (if applicable). Enforcement must be hard limits, not soft alerts: the system must refuse or truncate requests that exceed the budget, not log a warning and continue. Budget violations are logged as governance events (see ops-governance audit trail). Budget thresholds are reviewed at each engagement milestone and adjusted based on observed usage patterns — not before, not reactively only after a bill arrives.

### Cost-vs-Performance Evaluation Matrix
Before selecting a model or infrastructure configuration for production, run a structured cost-performance evaluation:
1. **Define the performance benchmark**: identify the task(s) the model must perform and the acceptable performance threshold for each (e.g., accuracy ≥ 0.92, latency p95 ≤ 2s, hallucination rate ≤ 3%).
2. **Enumerate candidate configurations**: list 2–4 candidate model/infrastructure combinations with their cost-per-request estimates.
3. **Evaluate each candidate**: run each candidate against the benchmark. Record performance metrics and cost per request.
4. **Build the matrix**: a table with rows = candidates, columns = performance metrics + cost per request + cost per 1k requests.
5. **Select and document**: choose the configuration that meets the performance threshold at the lowest cost. Document the selection rationale. If the cheapest option does not meet the threshold, document why the next tier is required — do not silently pick a more expensive option without explanation.
The matrix is a client-facing deliverable for engagements where the client funds the compute.

### Operational Cost Projection
Before production go-live, produce a written cost projection covering:
- **Request volume estimates**: expected requests per hour, per day, and at peak load. Source from client usage data or comparable engagement baselines.
- **Token consumption per request**: average and p95 input + output tokens, based on benchmark data or representative sample.
- **Unit cost calculation**: cost per 1k tokens (or per request) for the chosen model/infrastructure at the current provider pricing.
- **Total projected cost**: daily, monthly, and annual projections at average and peak load.
- **Cost at failure scenarios**: runaway request scenario (e.g., retry storm, malformed prompts triggering max output), cost impact if the token budget enforcement fails.
- **Cost-saving levers**: caching hit rates, batch processing options, model downgrade paths for lower-priority tasks.
The projection is reviewed by the client before go-live approval. It is updated at each major usage pattern change.

### Environmental Impact Report (NIST AI 600-1)
For clients with sustainability or ESG reporting obligations, produce an environmental impact report as a standing deliverable. The report covers: (1) total tokens consumed over the reporting period (monthly or quarterly), (2) estimated energy consumption using provider-published efficiency data or industry benchmarks [UNVERIFIED: provider-specific kWh/token figures — use available published data or mark as estimated], (3) estimated CO₂-equivalent emissions using the energy mix of the data center region where the model is hosted, (4) comparison to prior period (trend), (5) efficiency improvements implemented (caching, model selection, prompt compression) and their estimated environmental impact reduction. NIST AI 600-1 (2024) identifies environmental impact — including energy use and carbon footprint — as a risk dimension for generative AI systems that operators must track and manage. [UNVERIFIED: exact action ID within NIST AI 600-1 for environmental impact reporting.] Report methodology and assumptions are documented so the client can reproduce the calculation independently.

## Checklist

**Before production go-live:**
- [ ] Token budgets defined for all system boundaries (context, output, session, daily)
- [ ] Token budget enforcement implemented as hard limits (not soft alerts)
- [ ] Cost-versus-performance evaluation matrix completed and documented
- [ ] Operational cost projection written, reviewed, and approved by client
- [ ] Cost attribution instrumented at the request level (per-engagement, per-model)
- [ ] Cost-at-failure scenarios assessed and budget failure handling tested

**For each engagement milestone:**
- [ ] Actual vs. projected cost comparison completed
- [ ] Token budget thresholds reviewed against observed usage patterns
- [ ] Cost-performance matrix updated if model or infrastructure changed

**For clients with sustainability obligations:**
- [ ] Environmental impact reporting methodology documented and agreed with client
- [ ] Token consumption and estimated energy/CO₂ metrics tracked per reporting period
- [ ] Environmental impact report delivered per the agreed cadence

**Ongoing governance:**
- [ ] Token budget violations logged as governance events in the audit trail
- [ ] Cost overrun root-cause analysis documented when projected vs. actual variance exceeds threshold (define threshold with client)

## SelOps-Specific Context

Commercial viability at scale is a SelOps responsibility, not a client afterthought. A system that works technically but costs 10× what the client expected to pay is a project failure. Build cost governance in from the design phase — not as a post-launch optimization exercise.

Token budget failures have a compounding risk: an unbounded context window does not just cost more money — it increases the model's susceptibility to context-injection attacks (see ops-adversarial-security) and makes output behavior less predictable. Budget enforcement is therefore a security and reliability control as well as a cost control.

For regulated-industry clients (finance, healthcare), cost projections and environmental reports may be required for procurement approval or regulatory filing. Produce them proactively, not on request.

The FinOps Foundation framework provides practitioner guidance on cloud cost management applicable to AI infrastructure. FinOps principles — inform, optimize, operate — map directly onto the patterns above: instrument cost visibility (inform), evaluate cost-performance tradeoffs (optimize), enforce budgets and review projections (operate).

## References

- **NIST AI 600-1 (2024)** — Artificial Intelligence Risk Management Framework: Generative AI Profile. Identifies environmental impact — energy use, carbon footprint, and resource consumption — as a risk dimension for generative AI systems that operators must track and manage. [UNVERIFIED: specific action/control ID within the document.]
- **ISO/IEC 42001:2023** — Clause A.4 [UNVERIFIED: exact clause number] addresses resource planning and management for AI systems, including compute resource governance and cost-related risk controls.
- **FinOps Foundation** — FinOps Framework (https://www.finops.org/framework/). Practitioner framework for cloud and AI cost management: inform, optimize, operate. No fabricated control IDs — cite by framework name and principle.
