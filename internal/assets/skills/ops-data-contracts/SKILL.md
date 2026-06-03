---
name: ops-data-contracts
description: "SelOps data contract definition and validation. Trigger: When defining, validating, or evolving data schemas between producers and consumers."
---

# Data Contracts

## When to Use

Load this skill when defining or modifying schemas between AI system components or between SelOps-built systems and client systems. Use it when validating API boundaries, handling model I/O versioning, or responding to a breaking change request from a producer or consumer.

## Core Principles

- **Contract first, code second.** The schema must exist and be reviewed before any implementation begins. No code ships without a validated contract.
- **Versioning is non-negotiable.** Every schema that crosses a service boundary carries a version identifier. There is no "v0" or "latest" in production.
- **Breaking changes require a migration window.** Removing or renaming a field, changing a type, or tightening a constraint is a breaking change. The old contract must remain valid for the agreed deprecation window before removal.
- **Validate at both ends.** Producers validate on output. Consumers validate on input. A contract that only one side enforces is not a contract.
- **LLM output is a contract boundary.** If a downstream system parses LLM output, that output format is a contract. Treat it with the same rigor as any other API response.

## Patterns

### Schema-First Design
Write the schema before writing the code that produces or consumes it. Share the schema with the counterpart team (or system) and get sign-off before implementation starts. Never let the schema emerge from the implementation.

### Version Envelope
Wrap every inter-service payload in a version envelope: `{ "schema_version": "1.2", "payload": {...} }`. Do this for REST responses, event payloads, and LLM prompt/response structures. Never rely on URL versioning as the only version signal.

### Backward Compatibility Rule
New fields: allowed, must be optional with a documented default. Removed fields: never from a live contract without a two-step deprecation (mark deprecated → wait for deprecation window → remove). Changed types: treated as a removal + addition.

### LLM Output Contract
When a system parses LLM responses, define the expected output schema explicitly (JSON schema, Pydantic model, or equivalent). Add a structured output validation step before the response is used downstream. Log validation failures as schema violations, not as generic errors.

### Producer-Consumer Alignment Check
Before any schema change ships: (1) identify all consumers of the schema, (2) confirm each consumer can handle the new version, (3) confirm fallback behavior for consumers that receive an unexpected version. Document this analysis in the change's decision log.

### Contract Test Gate
Every producer must have a contract test that validates its output against the agreed schema version. This test runs in CI. A contract test failure blocks merge — it does not block only deployment.

## Checklist

**Before defining a contract:**
- [ ] Identified all producers and consumers of this data shape
- [ ] Determined the version identifier format (semantic, date-based, or sequential)
- [ ] Confirmed counterpart teams have reviewed and signed off on the schema
- [ ] Schema written in a machine-readable format (JSON Schema, Protobuf, OpenAPI, Pydantic)

**Before shipping a schema change:**
- [ ] Change classified: additive (safe) or breaking (requires migration plan)
- [ ] Breaking changes: old version still valid and served for the deprecation window
- [ ] Contract tests updated to cover new schema version
- [ ] All consumers notified and verified to handle the new version
- [ ] Decision log updated with rationale and backward-compatibility analysis

**After a contract change ships:**
- [ ] Monitoring shows no unexpected validation failures at consumer boundaries
- [ ] Old schema version retirement date documented and tracked
- [ ] Runbook updated if operational behavior changed

## SelOps-Specific Context

AI consultancies build systems where the data shapes are often defined by model behavior, not by engineers. LLM outputs are variable, partially structured, and sensitive to prompt changes. This makes data contracts harder — and more necessary — than in traditional software.

At SelOps, data contracts govern three distinct boundaries: (1) the interface between SelOps-built AI components and client systems, (2) the interface between AI pipeline stages within a single system, and (3) the interface between a model's output and the downstream code that parses it.

Changes to a prompt template that alter the output structure are breaking changes to the LLM output contract. Treat them as such. Version the prompt template alongside the schema it produces.

When a client's system changes its input format, that is a contract change from the client side. Validate it explicitly, do not silently adapt.
