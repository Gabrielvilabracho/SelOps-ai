---
name: ops-deliver
description: "SelOps pipeline phase: handoff and documentation. Trigger: When the OPS router enters the deliver phase of the operational pipeline."
---

# Deliver (Pipeline Phase 5 of 5)

## Role in the Pipeline

Receives: the review verdict and quality report from ops-review (PASS or PASS-WITH-WARNINGS only — FAIL blocks this phase).
Produces: a delivery package — client-facing summary, updated living docs, decision log, confirmed audit trail, and closed engagement record.

## When to Use

The router enters this phase only when ops-review produces a PASS or PASS-WITH-WARNINGS verdict. A FAIL verdict blocks this phase entirely until the operator resolves the failure.

## Procedure

1. Produce the structured client summary with four sections:
   - **Status**: the outcome in one sentence (succeeded / partially succeeded with caveats).
   - **What changed**: a concrete list of the systems and states that were modified.
   - **Impact**: the observable effect on the client's environment or operations.
   - **Next action**: the specific action the client or operator must take, if any.
2. Update living docs and runbooks for any operational procedure that changed during execution, per ops-standard-documentation standards. Changes must be in the same delivery as the operational change — no "docs later" exceptions.
3. Record the decision log entry: what was decided, why, which alternatives were considered, and which rollback procedure was used or remains available. Use the ADR template from ops-standard-documentation for decisions that affect architecture.
4. Confirm audit trail entries are complete and accurate: all privileged operations from the execution log have corresponding audit entries per ops-governance standards (timestamp, actor, action type, target, outcome, correlation ID).
5. Close the engagement record: mark the task as complete, note any PASS-WITH-WARNINGS items and their disposition, and record the final state of all affected systems.
6. If PASS-WITH-WARNINGS: include the warning list explicitly in the client summary and the engagement record. Do not omit warnings from client communication.

## Inputs / Outputs

**Consumes**: review verdict and quality report (from ops-review), execution log (from ops-produce), brief (success criteria, engagement context), execution plan (decision rationale)
**Produces**: client-facing summary (4-section format), updated living docs and runbooks, decision log entry, confirmed audit trail, closed engagement record

## Gates & Escalation

Halt and escalate to the operator if:
- The review verdict is FAIL — do not produce a delivery package for a failed execution.
- The audit trail cannot be confirmed complete — delivery without a complete audit trail violates ops-governance standards for regulated engagements.
- A PASS-WITH-WARNINGS item requires client decision before the engagement can be closed.

Do not mark the engagement record as complete until the client summary has been delivered and the audit trail is confirmed. Delivery is not complete until documentation is updated.

## Cross-References

- ops-standard-documentation — structured client summary format, living docs synchronization rule, ADR template, client handoff package
- ops-governance — immutable audit log format, compliance checkpoint registry, model update governance record
- ops-graduated-autonomy — rollback documentation rule, scope verification (for final state confirmation)
