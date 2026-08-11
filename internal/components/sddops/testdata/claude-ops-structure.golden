---
name: ops-structure
description: "SelOps pipeline phase: planning and decomposition. Trigger: When the OPS router enters the structure phase of the operational pipeline."
---

# Structure (Pipeline Phase 2 of 5)

## Role in the Pipeline

Receives: the structured brief from ops-brief (scope, risk classification, autonomy level, success criteria, rollback expectation).
Produces: a concrete, ordered execution plan with explicit checkpoints that ops-produce executes step by step.

## When to Use

The router enters this phase after ops-brief produces a confirmed brief with no unresolved questions. Do not enter this phase without a confirmed brief.

## Procedure

1. Decompose the task into discrete, atomic steps. Each step must produce a verifiable state change or observable output.
2. Identify dependencies between steps and enforce ordering. Steps with unresolved dependencies must block the dependent step, not skip it.
3. Classify each step as reversible or irreversible. Irreversible steps require explicit notation and a per-step confirmation gate at the autonomy level in effect.
4. Write the rollback procedure for each step:
   - Supervised level: one rollback action per step, documented before ops-produce executes that step.
   - Autonomous level: one full-scope rollback procedure covering all steps, written before any step executes.
5. Identify human approval gates — steps that require explicit operator confirmation before execution regardless of autonomy level.
6. Map affected systems for each step: name the system, the type of interaction (read / write / configuration / deployment), and the expected impact radius.
7. Flag contract-boundary changes: for any step that modifies a schema, module interface, or data contract shared across system boundaries, mark the step as a contract-boundary change in the execution plan. These steps require contract validation in their post-step verification (see ops-data-contracts).
8. Define a post-step verification for each step: the observable condition that confirms the step succeeded. For contract-boundary steps, the verification must include the relevant contract test or schema validation.

## Inputs / Outputs

**Consumes**: structured brief (from ops-brief) — scope, risk classification, autonomy level, success criteria, rollback expectation
**Produces**: execution plan — ordered steps, dependency map, reversibility classification per step, rollback procedures, human approval gates, affected system map, post-step verifications

## Gates & Escalation

Halt and escalate to the operator if:
- A step cannot be decomposed into a verifiable action (the step is too vague to execute safely).
- A rollback procedure cannot be written for an irreversible step and the autonomy level is Supervised or higher.
- The affected system map reveals systems outside the approved engagement scope from the brief.
- The number of irreversible steps exceeds what the registered autonomy level permits without additional approval.

Do not hand the plan to ops-produce until every step has a post-step verification, every irreversible step has a documented rollback or explicit exception approval, and every contract-boundary step is flagged with its required contract validation.

## Cross-References

- ops-graduated-autonomy — rollback documentation rule, escalation trigger list, three-level model
- ops-governance — change approval workflow, immutable audit log format
- ops-data-contracts — contract validation requirements for boundary-crossing steps
