---
name: ops-review
description: "SelOps pipeline phase: verification and quality. Trigger: When the OPS router enters the review phase of the operational pipeline."
---

# Review (Pipeline Phase 4 of 5)

## Role in the Pipeline

Receives: the execution log from ops-produce (per-step results, pre/post states, overall success criteria evaluation).
Produces: a review verdict (PASS / PASS-WITH-WARNINGS / FAIL) with a structured quality report that ops-deliver uses to prepare the handoff.

## When to Use

The router enters this phase after ops-produce completes all steps and produces a full execution log. Do not enter this phase with a partial log — a mid-execution halt requires a different path (rollback decision, not review).

## Procedure

1. Check the overall success criteria from the brief against the execution log. Each criterion must be satisfied by an observed post-state in the log. Unsatisfied criteria are failures.
2. Run the applicable quality gates:
   - Evaluation coverage: every step in the plan has a logged result.
   - Observability checks: structured log entries exist for all executed steps per ops-observability standards (timestamp, action type, target, outcome).
   - Governance checkpoints: audit trail entries cover all privileged operations per ops-governance standards (actor, action, target, outcome, correlation ID).
3. Verify no scope expansion occurred: compare the set of systems touched in the execution log against the approved scope from the brief. Any system not in the approved scope is an automatic FAIL.
4. Confirm the rollback procedure is still valid: verify that the documented rollback can still be applied to return to the pre-execution state.
5. Classify any incidents using the ops-governance four-category system (output quality degradation / safety violation / service unavailability / data exposure). If any incident is present, include it in the verdict.
6. Produce the review verdict:
   - **PASS**: all success criteria met, no scope expansion, audit trail complete, rollback valid, no incidents.
   - **PASS-WITH-WARNINGS**: success criteria met but minor quality gate gaps or low-severity warnings present. Document each warning explicitly.
   - **FAIL**: one or more success criteria unmet, scope expansion detected, audit trail incomplete, or an incident present.

## Inputs / Outputs

**Consumes**: execution log (from ops-produce), brief (success criteria, approved scope), execution plan (rollback procedures)
**Produces**: review verdict (PASS / PASS-WITH-WARNINGS / FAIL), quality report listing criteria evaluation, gate results, scope check, rollback status, and any incidents

## Gates & Escalation

Halt and escalate to the operator if:
- Scope expansion is detected — this is not a warning, it is a FAIL requiring human review before any delivery.
- An incident of category safety/compliance violation or data exposure is present — immediate escalation, do not proceed to delivery.
- The execution log is incomplete (missing steps) — cannot produce a valid verdict without a complete log.

A FAIL verdict blocks ops-deliver. The operator must decide: rollback, re-execute, or accept partial delivery with documented exceptions.

## Cross-References

- ops-governance — AI-specific incident classification, immutable audit log format, compliance checkpoint registry
- ops-graduated-autonomy — escalation trigger list, scope verification protocol
- ops-observability — structured error classification, LLM metrics set, agent chain tracing
