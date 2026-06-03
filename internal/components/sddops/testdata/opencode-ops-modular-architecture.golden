---
name: ops-modular-architecture
description: "SelOps modular architecture patterns and boundaries. Trigger: When designing or reviewing module boundaries, service contracts, or data store ownership."
---

# Modular Architecture

## When to Use

Load this skill when designing a new AI system component, reviewing proposed module boundaries, evaluating whether a change crosses architectural lines, or deciding how to structure a new integration between SelOps systems and client systems. Use it when the question is "where does this belong?" or "how tightly coupled is this?"

## Core Principles

- **AI inference and business logic are separate modules.** The code that calls a model and the code that applies the result to business rules must not be in the same module. Mixing them makes both untestable and unmaintainable.
- **Prompt management is a first-class module.** Prompts are not strings embedded in business logic. They are versioned artifacts managed in their own module, with their own change lifecycle and their own tests.
- **Interfaces are model-agnostic.** The interface between your AI inference layer and the rest of the system must not expose model-specific details (provider names, model IDs, token counts) in its public contract. Caller code must be swappable between models without changes.
- **Each module owns its data.** No module reads or writes directly to another module's data store. Cross-module data access goes through the owning module's API. No exceptions.
- **Pipeline boundaries are explicit.** In data pipelines, the boundary between each stage is a defined data shape (see ops-data-contracts), not an implicit convention. A pipeline stage that receives a raw dict and returns a raw dict has no boundary.

## Patterns

### Inference-Logic Separation
Structure: `inference/` module (calls the model, handles retries and rate limits, returns a typed response) ↔ `domain/` module (applies business rules to the typed response, makes decisions, mutates state). The domain module never imports the inference module directly — it depends on an interface. The inference module never imports domain logic.

### Prompt Registry Pattern
Manage prompts in a dedicated `prompts/` module: versioned template files (e.g., `prompts/classify_v3.txt`), a loader that selects the correct version, and a test suite that validates each template against a set of known inputs and expected output structures. Prompt changes go through the same review process as code changes.

### Model-Agnostic Interface
Define a provider interface (e.g., `ModelProvider { complete(prompt: Prompt) -> TypedResponse }`). Each supported model (OpenAI, Anthropic, local) is an implementation of that interface. Callers import the interface, not a specific implementation. Model selection is configuration, not code.

### Pipeline Stage Contract
Each stage in a data pipeline is a function with a typed input and typed output, validated against the schema at the stage boundary. Stage functions are independently testable. The pipeline runner wires stages together but contains no business logic itself.

### Integration Adapter Pattern
When connecting to a client system (their CRM, their data warehouse, their API), isolate the integration behind an adapter interface in an `adapters/` module. The rest of the system depends on the adapter interface, not on the client's API structure. When the client's system changes, only the adapter changes.

### Dependency Direction Rule
Dependencies always flow inward: `adapters → domain → inference`, never the reverse. The domain module has no knowledge of how it is called or which adapter wraps it. If a lower-level module needs to call a higher-level one, extract the shared concern into a shared utilities module.

## Checklist

**When designing a new component:**
- [ ] AI inference code identified and isolated in its own module
- [ ] Prompt templates in a separate versioned module, not embedded in logic
- [ ] Public interface of the inference module is model-agnostic
- [ ] Each module has exactly one data store, and no other module accesses it directly
- [ ] Pipeline stage inputs and outputs have typed contracts
- [ ] Integration with external systems goes through an adapter, not direct coupling

**When reviewing a proposed change:**
- [ ] Does this change cross a module boundary? If yes, is the crossing going through the module's public API?
- [ ] Does this change embed a prompt string into business logic? If yes, reject it — move to prompt registry.
- [ ] Does this change reference a specific model name or provider in non-configuration code? If yes, it violates the model-agnostic interface rule.
- [ ] Does this change add a direct read/write to another module's data store? If yes, reject it.

**Before handing off to client:**
- [ ] Architecture diagram reflects actual module boundaries
- [ ] Each module boundary documented with its public interface
- [ ] Prompt registry documented: how to version, how to test, how to deploy a prompt change
- [ ] Dependency direction documented and verified against the actual codebase

## SelOps-Specific Context

AI consultancies often build systems quickly to validate hypotheses. The architectural debt from those early builds accumulates at inference boundaries — model calls end up entangled with business rules, prompt strings live inside application logic, and replacing one model with another requires touching dozens of files.

At SelOps, every system we build or operate must be designed to survive a model replacement. The model is a dependency, not an architecture assumption. If a client's LLM provider discontinues a model or raises prices, we must be able to swap the provider without rewriting the system.

This is not an academic concern. It has happened on client engagements. Systems without a model-agnostic interface took 3–5x longer to migrate than systems with one.

Prompt management is the other common failure mode. Prompts that live as raw strings in application code get changed without review, without testing, and without version history. The result is silent regressions — the system appears to work but produces subtly different outputs. SelOps treats prompt changes as code changes: reviewed, versioned, tested.
