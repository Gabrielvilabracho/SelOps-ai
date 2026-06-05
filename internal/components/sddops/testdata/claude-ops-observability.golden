---
name: ops-observability
description: "SelOps observability: metrics, tracing, and logging standards. Trigger: When implementing or reviewing health endpoints, distributed traces, or log levels."
---

# Observability

## When to Use

Load this skill when adding instrumentation to an AI system, debugging a production issue, reviewing log quality, implementing health checks, or setting up alerting for an AI pipeline. Use it whenever the question is "how do we know this system is working correctly?"

## Core Principles

- **Observe AI-specific failure modes.** Traditional error rates and latencies are necessary but not sufficient for AI systems. Track token consumption, output quality signals, and model degradation indicators alongside standard service metrics.
- **Structured logs everywhere.** Unstructured log lines are unsearchable at scale. Every log entry is a JSON object with a consistent schema: timestamp, level, service, trace_id, and a message field. No printf-style logs in production.
- **Traces cross component boundaries.** A trace that ends at the LLM API call tells you nothing about the full request lifecycle. Propagate trace context through every hop: API → pipeline → model call → response parsing → downstream delivery.
- **Alerts fire on user-visible impact, not on internal metrics thresholds.** Alert when customers are affected or will be within minutes. Do not alert on CPU spikes that do not degrade response quality.
- **Health checks reflect real capability.** A /health endpoint that only checks if the process is running is not a health check. It must verify the system can do its job: model reachable, data sources connected, output parseable.
- **Security observability is distinct from reliability observability.** Track anomalous input patterns (prompt injection signals, unusual token velocity, repeated adversarial probes) alongside standard quality metrics. A system can have 99.9% uptime and zero latency alerts while being actively attacked through its prompt interface. OWASP LLM01:2025 (Prompt Injection) and MITRE ATLAS provide the adversarial threat vocabulary; observability is the detection layer. (OWASP LLM01:2025; MITRE ATLAS [UNVERIFIED: specific technique IDs])
- **Cost and environmental impact are observable system properties.** Token consumption is not just a quality signal — it is a cost signal and an environmental signal. Track actual cost per request, cost per pipeline stage, and cumulative cost per engagement. For production systems, track or estimate compute carbon footprint when clients have ESG reporting requirements. (NIST AI 600-1 Environmental Impacts)

## Patterns

### LLM Metrics Set
Instrument every LLM call with: request latency (p50, p95, p99), input token count, output token count, model name and version, and error type (timeout / rate-limit / malformed-output / model-error). Aggregate these per model, per pipeline stage, and per client engagement.

### Hallucination Signal Logging
When a system has a validation layer (schema validation, confidence scoring, human review), log validation failures with the full input context and the specific validation rule that failed. These are your hallucination signals. Aggregate them. Track them over time. A rising failure rate is a model degradation signal.

### Agent Chain Tracing
In multi-agent or multi-step pipelines, each step must receive and propagate a trace context (W3C TraceContext or OpenTelemetry). The trace must be explorable end-to-end from a single trace ID. Never create a new trace ID mid-pipeline unless starting an intentionally separate operation.

### Structured Error Classification
Classify every error at the point it is logged: `model_error`, `schema_validation_error`, `timeout`, `rate_limit`, `upstream_failure`, `downstream_delivery_failure`, `internal_error`. Use these classifications in dashboards and alerting rules. Generic "error" classifications are not actionable.

### Degradation Threshold Alerting
Define and document degradation thresholds for each AI pipeline: acceptable error rate, acceptable p95 latency, acceptable token budget per request. Set alerts at 80% of threshold (warning) and 100% (critical). Review and adjust thresholds after each major model or prompt change.

### Health Check Layers
Implement three-layer health checks: (1) liveness — is the process running? (2) readiness — can the service accept traffic? (verify model reachable, DB connected); (3) capability — can it produce correct output? (run a known-good input through the pipeline and validate the output). Only liveness is required for restart decisions. Readiness gates traffic routing. Capability signals operational health.

### Security Observability Signals (OWASP LLM Top 10 / MITRE ATLAS)
Instrument the following security signals alongside standard metrics:
- **LLM01:2025 / Prompt Injection**: log input length outliers, inputs containing instruction-override patterns (common jailbreak phrases), and inputs that cause unexpected output format shifts. Set an alert threshold for prompt injection pattern density.
- **LLM06:2025 / Excessive Agency**: log every out-of-scope tool call, every permission escalation attempt, and every action that would exceed the registered autonomy level. These should never occur — even one is an alert-worthy event.
- **LLM10:2025 / Unbounded Consumption / Model Extraction**: log query rate per client/session, flag sessions with systematic input variation (adversarial probing patterns), and alert on sustained high query volume from a single source.
- **MITRE ATLAS [UNVERIFIED: AML.TA0013 Exfiltration]**: log unusual output lengths, outputs containing unexpected data structures, and outputs that reference system internals.

Route security signals to a dedicated security dashboard, not just the standard ops dashboard.

### FinOps Metrics Set
Track per-request cost metrics: estimated USD cost (input tokens × price + output tokens × price per model), cost per pipeline stage, running total per engagement per billing period. Define a token budget per request type (documented in the engagement record). Alert at 80% of budget (warning) and 100% (critical). Compare cost-vs-quality tradeoffs when evaluating model upgrades: a 20% quality improvement at 3× the cost requires explicit client approval. Log cost alongside quality metrics so degradation and cost can be correlated.

### Environmental Impact Tracking (NIST AI 600-1)
For engagements where clients have ESG or sustainability reporting requirements: (1) instrument total compute hours per model (GPU/TPU hours if trackable via provider APIs), (2) record model efficiency (tokens/second, tasks completed per GPU-hour), (3) include environmental impact summary in engagement reports alongside cost. If direct measurement is not available, document the model provider's carbon reporting methodology and reference it. This is a reporting requirement, not a blocker — but it must be addressed before handoff to clients with ESG obligations.

## Checklist

**When adding a new AI component to production:**
- [ ] LLM metrics set instrumented (latency, token counts, model name/version, error type)
- [ ] Trace context propagated through all pipeline stages
- [ ] Structured logging in place with consistent schema
- [ ] /health endpoint implements all three layers (liveness, readiness, capability)
- [ ] Degradation thresholds defined and documented
- [ ] Alerts configured at 80% (warning) and 100% (critical) of each threshold
- [ ] Hallucination/validation failure logging in place if a validation layer exists
- [ ] Security observability signals instrumented (prompt injection patterns, excessive agency attempts, query rate anomalies)
- [ ] Cost metrics tracked per request and per pipeline stage
- [ ] Token budget defined and alert configured
- [ ] Environmental impact tracking in place if client has ESG requirements

**During a production incident:**
- [ ] Trace ID retrieved from the failing request
- [ ] Full trace explored across all pipeline stages
- [ ] Error classification confirmed (not just "error" — what type?)
- [ ] Token budget and model version checked against baseline
- [ ] Validation failure rate compared to pre-incident baseline

**After resolving an incident:**
- [ ] Root cause documented with the specific observability gap that delayed detection
- [ ] Missing instrumentation added
- [ ] Alert thresholds reviewed and adjusted if threshold was correct but alert didn't fire

## SelOps-Specific Context

Generic observability tooling is designed for deterministic systems. AI systems degrade in non-deterministic ways: a model that was accurate last week may be less accurate today due to distribution shift, prompt template drift, or a silent model update from the provider. None of these produce an HTTP 500.

At SelOps, observability for AI systems must account for quality degradation, not only availability degradation. This means output validation must be instrumented, not just service uptime.

When handing off an AI system to a client, the observability setup is part of the deliverable. The client must be able to detect model degradation, validate output quality, and investigate incidents without SelOps on call. If they cannot, the system is not production-ready.

Each client engagement gets a documented observability baseline at handoff: current thresholds, what each alert means, how to investigate each alert type, and which team (client or SelOps support) owns each response action.

## References

- OWASP Top 10 for LLM Applications 2025 — LLM01:2025 Prompt Injection, LLM06:2025 Excessive Agency, LLM10:2025 Unbounded Consumption (https://genai.owasp.org/llm-top-10/). Security observability signals are mapped to these categories.
- MITRE ATLAS — adversarial threat taxonomy for AI/ML systems (atlas.mitre.org). Specific technique IDs (e.g., AML.TA0013 Exfiltration) cited as [UNVERIFIED] where exact numbers are not confirmed against the published ATLAS matrix.
- NIST AI 600-1 (2024) — Generative AI Profile. Cited for environmental impact tracking guidance and compute cost observability obligations. Specific subcategory action IDs marked [UNVERIFIED].
- NIST AI RMF 1.0 (2023) — MEASURE function, including MEASURE 1.1 [UNVERIFIED: exact subcategory] on measurement methods for AI risk monitoring.
