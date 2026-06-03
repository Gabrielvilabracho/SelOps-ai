---
name: ops-brief
description: "SelOps pipeline phase: intake and framing. Trigger: When the OPS router enters the brief phase of the operational pipeline."
---

# Brief (Pipeline Phase 1 of 5)

## Role in the Pipeline

Receives: a raw operational task request from the operator or client.
Produces: a structured, unambiguous brief that the ops-structure phase consumes to build an execution plan.

## When to Use

The router enters this phase when a task clears the OPS threshold risk gate. Do not enter this phase for low-risk inline tasks — those execute directly without the pipeline.

## Procedure

1. Run the Scope Verification Protocol (from ops-graduated-autonomy): confirm (a) which system and environment this affects, (b) which client engagement this belongs to, (c) what the approved scope is for this session, (d) what the current autonomy level is for this engagement and environment.
2. Classify risk using the five OPS threshold dimensions: environment (production vs. staging), reversibility (reversible vs. irreversible), data mutation (read-only vs. write), systems affected (count and criticality), and time-to-detect (immediate vs. delayed impact).
3. Identify the autonomy level (Suggest / Supervised / Autonomous) that applies to this task and environment.
4. Define explicit success criteria: what observable state confirms the task is complete.
5. Define rollback expectations: what action or procedure restores the prior state if the task must be aborted.
6. Surface all ambiguities as specific questions. Do not assume. Do not proceed with unresolved ambiguities.
7. Produce the structured brief document.

## Inputs / Outputs

**Consumes**: raw task request (text), engagement record, registered autonomy level, environment context
**Produces**: structured brief containing scope, risk classification, autonomy level, success criteria, rollback expectation, and a list of unresolved questions (if any)

## Gates & Escalation

Halt and return to the operator if:
- The autonomy level required by this task exceeds the registered level for this engagement and environment.
- One or more scope verification answers cannot be confirmed from available context.
- Risk classification produces a dimension that the current autonomy level does not cover (e.g., irreversible action at Supervised level without explicit exception approval).
- The task involves systems outside the approved engagement scope.

Do not proceed to ops-structure until all four scope verification answers are confirmed and success criteria are written.

## Cross-References

- ops-graduated-autonomy — Scope Verification Protocol, three-level model, autonomy-level gate
- ops-governance — change approval workflow, regulated-industry context
