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

## Checklist

**When adding a new AI component to production:**
- [ ] LLM metrics set instrumented (latency, token counts, model name/version, error type)
- [ ] Trace context propagated through all pipeline stages
- [ ] Structured logging in place with consistent schema
- [ ] /health endpoint implements all three layers (liveness, readiness, capability)
- [ ] Degradation thresholds defined and documented
- [ ] Alerts configured at 80% (warning) and 100% (critical) of each threshold
- [ ] Hallucination/validation failure logging in place if a validation layer exists

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
