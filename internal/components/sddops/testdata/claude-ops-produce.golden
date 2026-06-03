---
name: ops-produce
description: "SelOps pipeline phase: execution. Trigger: When the OPS router enters the produce phase of the operational pipeline."
---

# Produce (Pipeline Phase 3 of 5)

## Role in the Pipeline

Receives: the execution plan from ops-structure (ordered steps, rollback procedures, human approval gates, post-step verifications).
Produces: an execution log with result, pre-state, and post-state for each completed step.

## When to Use

The router enters this phase after ops-structure produces a confirmed execution plan. Do not enter this phase without a complete plan — improvising steps during execution is prohibited.

## Procedure

For each step in the execution plan, in order:

1. Verify pre-state: confirm the system is in the expected state before executing this step. If pre-state does not match, halt immediately — do not execute the step.
2. Check for a human approval gate on this step. If present, pause and request explicit operator confirmation before continuing.
3. Execute the step exactly as specified in the plan. Do not improvise, expand scope, or apply shortcuts not in the plan.
4. Verify post-state: confirm the observable condition that indicates this step succeeded.
5. Log the result: step identifier, action taken, pre-state observed, post-state observed, timestamp, outcome (success / failure / partial).
6. If the step fails: halt. Do not proceed to the next step. Report the failure with the pre- and post-state observed and the planned rollback procedure for this step.
7. If the step succeeds: check whether any escalation trigger has been activated (from ops-graduated-autonomy). If triggered, halt and report before proceeding.

After all steps complete: verify that the overall success criteria from the brief are satisfied.

## Inputs / Outputs

**Consumes**: execution plan (from ops-structure) — ordered steps, rollback procedures, approval gates, post-step verifications
**Produces**: execution log — per-step: action taken, pre-state, post-state, outcome, timestamp; overall: success criteria evaluation result

## Gates & Escalation

Halt immediately and report to the operator if:
- Pre-state verification fails for any step — the system is not in the expected state.
- A step fails — do not attempt recovery or workarounds not in the plan.
- An escalation trigger activates during execution (scope expansion detected, unexpected system state, rollback path invalidated).
- The task would take materially longer than estimated (more than 2x) — stop and report before continuing.
- Any irreversible step is reached and the rollback procedure has become invalid since the plan was written.

Never improvise outside the plan. If the plan is wrong, halt and return to ops-structure.

## Cross-References

- ops-graduated-autonomy — escalation trigger list, scope verification protocol, rollback documentation rule
- ops-governance — immutable audit log format, change approval workflow
- ops-observability — structured logging, structured error classification
