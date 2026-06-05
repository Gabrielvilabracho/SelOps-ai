---
name: ops-graduated-autonomy
description: "SelOps graduated autonomy framework for AI agent execution permissions. Trigger: When setting agent autonomy levels or escalating scope boundaries."
---

# Graduated Autonomy

## When to Use

Load this skill when determining how much authority the operator agent has for a given task, when a task's scope appears to exceed what was approved, when deciding whether to proceed or escalate, or when setting up a new client engagement and defining the operating boundaries for this agent.

## Core Principles

- **The autonomy level is set per engagement, per environment, not globally.** An agent operating at Autonomous level in a staging environment may be at Suggest level in production for the same client. These settings are explicit, documented, and not assumed.
- **Escalate on scope creep, not on difficulty.** Difficulty is not a reason to escalate. Scope expansion is. If the task you are executing requires touching systems, data, or configurations outside what was originally approved, stop and report the expansion.
- **Autonomy upgrades require explicit approval.** You cannot self-promote from Supervised to Autonomous because a task would be faster that way. Autonomy level changes are human decisions, documented in the engagement record.
- **Every autonomy level has a rollback requirement.** At Suggest: the human approves before any execution, so rollback is trivial. At Supervised: every action must be individually reversible. At Autonomous: the full scope must have a documented rollback procedure before the task begins.
- **When in doubt, default to Suggest.** The cost of asking for confirmation is lower than the cost of an unintended production change.
- **Graduated autonomy is a formal adaptation of established human-automation interaction science.** The SelOps three-level model maps to the Sheridan-Verplank / Parasuraman-Sheridan-Wickens levels of automation (Parasuraman, Sheridan & Wickens, 2000, IEEE Trans. SMC 30(3):286-297): Suggest corresponds to HITL (Human-in-the-Loop) / Levels 1–4, where the computer offers alternatives and the human decides; Supervised corresponds to HOTL (Human-on-the-Loop) / Levels 5–7, where the computer executes and necessarily informs the human; Autonomous corresponds to HOOTL (Human-out-of-the-Loop) / Levels 8–10, where the computer decides and informs only if asked. The per-engagement + per-environment framing and the Scope Verification Protocol are SelOps-original additions to this taxonomy. Cite the foundational model when the framework is used in regulated contexts.

## Patterns

### Three-Level Model

**Suggest (Level 0)**
The agent identifies what to do and produces a detailed proposal. No action is taken until a human explicitly approves. Used when: the task is high-risk (production changes, data mutations, client-facing updates), the scope has not been previously reviewed, or the client engagement is new and trust is being established.

**Supervised (Level 1)**
The agent executes low-risk operations autonomously and reports a structured summary after each batch. The human reviews summaries and can halt at any point. Used when: the task type is well-understood and has been executed before, the environment is non-production or isolated, and rollback for each individual action is confirmed possible.

**Autonomous (Level 2)**
The agent handles the full defined scope end-to-end and escalates only on scope expansion or unexpected blockers. The human reviews the final result. Used when: the scope is tightly defined, the task type has been executed multiple times with no incidents, a full rollback procedure exists and is documented, and the client has explicitly approved autonomous operation for this task category.

### Scope Verification Protocol
Before starting any task, confirm: (1) which system and environment this affects, (2) which client engagement this belongs to, (3) what the approved scope is for this session, (4) what the current autonomy level is for this engagement and environment. Write these four answers down in the task log before executing.

### Escalation Trigger List
Escalate immediately if any of the following occur: the task requires touching a system not in the approved scope, the task would affect more data records than expected (define "expected" before starting), a dependent system is in an unexpected state, the rollback path is no longer valid, or the task would take materially longer than estimated (more than 2x).

### Autonomy-Level Gate at Task Start
For every task: check the registered autonomy level for this engagement + environment combination. If the task requires a higher level than registered, do not proceed. Return: "This task requires Autonomous-level approval for [environment]. Current registered level is [level]. Request approval before proceeding."

### Rollback Documentation Rule
At Supervised level: document rollback for each action before executing it. At Autonomous level: document the full rollback procedure for the entire scope before the task starts. If you cannot write the rollback procedure in concrete steps, you cannot start the task.

### EU AI Act Art.14.4 Compliance Mapping

For AI systems classified as high-risk under EU AI Act Annex III, Art.14.4 specifies supervisory mechanisms that must be enabled. Map the current SelOps autonomy level to the required sub-articles:

| Art.14.4 Mechanism | Suggest | Supervised | Autonomous |
|---|---|---|---|
| (a) [UNVERIFIED: understand capabilities/limitations] — documented in engagement record | Required | Required | Required |
| (b) [UNVERIFIED: detect/avoid automation bias] — Scope Verification Protocol + operator training | Not required | Required | Required |
| (c) [UNVERIFIED: interpret output correctly] — structured output logging + model card | Not required | Required | Required |
| (d) [UNVERIFIED: decide not to use / disregard output] — overrule mechanism documented | Not required | Not required | Required |
| (e) [UNVERIFIED: interrupt via stop procedure to safe state] — rollback procedure = stop mechanism | Not required | Not required | Required |

Sub-article letters (a)–(e) are marked [UNVERIFIED] because the exact letter-to-mechanism mapping has not been independently confirmed against the published regulation text. The mechanisms themselves (understand, avoid bias, interpret, override, interrupt) are real Art.14.4 requirements. For EU-based clients or systems classified as high-risk under the EU AI Act Annex III, confirm all applicable sub-articles are satisfied before setting Autonomous level.

## Checklist

**Before starting any task:**
- [ ] Autonomy level confirmed for this engagement + environment
- [ ] Scope verified: systems, environments, client engagement
- [ ] Rollback procedure written (per-action for Supervised, full-scope for Autonomous)
- [ ] Escalation triggers defined for this task (what unexpected condition would cause a stop)
- [ ] For EU-regulated high-risk systems: Art.14.4(a)–(e) compliance mapped to current autonomy level [UNVERIFIED: sub-article letters]
- [ ] Foundational taxonomy (Sheridan/Parasuraman HITL/HOTL/HOOTL levels) acknowledged in engagement record if client requests standards traceability

**During task execution (Supervised mode):**
- [ ] Each action logged before execution
- [ ] Each action logged with result after execution
- [ ] Structured summary prepared for human review after each batch
- [ ] Scope boundary checked after each action — has anything expanded?

**During task execution (Autonomous mode):**
- [ ] Full scope completed without scope expansion? If yes, proceed to final report.
- [ ] Scope expansion detected? Stop. Report to human before continuing.
- [ ] Unexpected blocker encountered? Stop. Report. Do not improvise a workaround.

**After task completion:**
- [ ] Task outcome logged: what was done, which systems were affected, what the final state is
- [ ] Any scope deviations noted (even minor ones)
- [ ] Rollback procedure confirmed still valid post-execution (or documented as consumed/voided)

## SelOps-Specific Context

Graduated autonomy is especially critical in an AI consultancy context because SelOps operates inside client systems — production databases, live pipelines, real customer data. A misstep at the wrong autonomy level is a client incident, not an internal one.

The autonomy level for a given engagement is a negotiated agreement, not a technical setting. Clients must understand and consent to the autonomy level granted to this agent. The engagement record documents this agreement. Do not treat autonomy levels as defaults — confirm them explicitly at the start of every engagement and at the start of every session in a new environment.

AI agents present a specific risk that human operators do not: they can execute at machine speed. An autonomous agent executing the wrong scope does more damage in 10 seconds than a human would in an hour. This makes the rollback requirement and the scope verification protocol non-negotiable, not best practices.

When building SelOps-internal automation, apply the same framework. Internal systems are not exempt. The staging/production distinction applies internally as much as it does for clients.

The academic foundations matter for regulated engagements. When a financial or healthcare client's compliance team asks "on what basis does this autonomy model work?", the answer is not "it's how SelOps does it" — it is "it is a software-operations adaptation of the Parasuraman-Sheridan-Wickens (2000) levels of automation, with per-engagement configurability added as a practical extension." This framing is respected and verifiable.

## References

- Parasuraman, R., Sheridan, T.B., & Wickens, C.D. (2000). A model for types and levels of human interaction with automation. *IEEE Transactions on Systems, Man, and Cybernetics — Part A*, 30(3), 286-297. Defines the 10-level human-automation interaction model underlying HITL/HOTL/HOOTL taxonomy.
- Sheridan, T.B., & Verplank, W.L. (1978). Human and computer control of undersea teleoperators. MIT Man-Machine Systems Lab Technical Report. Original source of the levels of automation concept.
- EU AI Act — Regulation (EU) 2024/1689, Art.14 (Human oversight), including Art.14.4(a)–(e) supervisory mechanism requirements for high-risk AI systems. Sub-article letter assignments marked [UNVERIFIED] — verify against published regulation text.
