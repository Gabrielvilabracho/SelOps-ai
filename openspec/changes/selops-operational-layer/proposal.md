# Proposal: SelOps Operational Layer

## Intent

Add a second, installable SelOps operational profile beside the existing DEV profile. Today the fork ships developer-oriented persona/skills/SDD only; SelOps needs a parallel Company Operating Agent profile that reuses the installer engine and remains purely additive.

## Scope

### In Scope
- `PersonaOperator` plus a fork-private operator persona asset.
- `ComponentSDDOps` shipping namespaced `ops-*` six-domain operational SDD assets.
- Operational skill catalog entries with `Category: "operational"`.
- `ComponentOperationalMCP` that wires external RAG/email/drive/CRM MCP connections.
- `PresetSelOpsOperational` bundling the operational components separately from DEV.

### Out of Scope
- RAG runtime, database, vector store, server, or retrieval engine.
- Email/Drive/CRM execution logic.
- Full research and authoring of the six domains' detailed content.
- Reworking existing DEV assets, engine pipeline, or current presets.

## Capabilities

### New Capabilities
- `selops-operational-profile`: install/select a parallel operational preset with operator persona, ops skills, and external operational MCP wiring.
- `selops-operational-assets`: embed fork-private `ops-*` / `selops-*` assets without upstream collisions.

### Modified Capabilities
- None.

## Approach

Keep the fork a pure config installer. Additive only: extend model/catalog/planner/CLI registration, add one operator-persona dispatch case, add one skills preset case, and add new component packages for operational SDD and MCP wiring. The six SDD-OPS domains define content shape only: standard documentation, modular architecture, data contracts, governance, observability, and graduated autonomy. Their full definitions are deferred to later research/spec/design work.

## Preserved Invariants

- No engine rewrite: pipeline, `agents.Adapter`, `filemerge`, backup/atomic writes stay intact.
- DEV remains 100% intact.
- No existing switch case behavior changes; only additive registration/new cases.

## Scope Boundary

This repo installs and wires configuration only. RAG and operational tools are external MCP servers; this change only writes their connection entries.

## Naming Strategy

Use `ops-*` / `selops-*` prefixes for all new assets and identifiers so fork-private content never collides with upstream gentle-ai, preserving mergeability.

## Affected Areas

| Package / Area | Impact | Description |
|---|---|---|
| `internal/model` | Modified | Add new persona, component, and preset IDs. |
| `internal/catalog` | Modified | Register operational components and skills category. |
| `internal/planner` | Modified | Add graph nodes/dependencies for the new components. |
| `internal/cli` | Modified | Add additive apply-step cases for the new components. |
| `internal/components/persona` | Modified | Dispatch `PersonaOperator` to a new asset. |
| `internal/components/skills` | Modified | Add operational preset selection logic. |
| `internal/components/sddops` | New | Inject namespaced six-domain SDD-OPS assets. |
| `internal/components/operationalmcp` | New | Merge operational MCP server definitions. |
| `internal/assets/` | Modified | Add `selops-*` persona, skill, and MCP assets. |

## Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| External RAG MCP project is not ready | Med | Keep boundary explicit; ship installer wiring independently. |
| Operator persona quality is underspecified | Med | Treat persona authoring as follow-up design/content work. |
| TUI may need clearer second-profile surfacing | Med | Confirm UX during design phase before implementation. |
| Six-domain content research is still pending | High | Keep proposal at installer shape level only. |

## Rollback Plan

Rollback is low risk: do not select the new preset, or uninstall the new operational components. Because the change is additive and uses existing backup/atomic-write behavior, removing the new assets and registrations cleanly restores DEV-only behavior.

## Dependencies

- External MCP servers for RAG and operational tools.
- Follow-up research defining the six operational domains' actual content.

## Success Criteria

- [ ] DEV preset behavior remains unchanged.
- [ ] A separate `selops-operational` preset can be planned without engine changes.
- [ ] Operational RAG/tools remain explicitly external to this repo.
- [ ] Affected package scope is limited to additive registration, persona, skills, new operational components, and assets.
