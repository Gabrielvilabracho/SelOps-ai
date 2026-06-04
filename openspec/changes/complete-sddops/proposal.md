# Proposal: Complete sddops

## Intent

`sddops` is still an install-only MVP. It injects OPS skills and the orchestrator prompt, but it does not install native sub-agents, slash commands, or the OpenCode JSON overlay proven in the old `sdd` component. Without those pieces, OPS cannot be invoked through `/ops-*` commands and OpenCode/Kilocode cannot delegate into the OPS pipeline.

## Scope

### In Scope
- Add native OPS sub-agent injection for adapters that already report `SupportsSubAgents()`.
- Add OPS slash-command injection for adapters that already report `SupportsSlashCommands()`.
- Add OpenCode/Kilocode OPS JSON overlay injection into `opencode.json`.

### Out of Scope
- Extracting `mergeJSONFile` to a shared package.
- Multi-mode overlays, sync behavior, registry rewiring, or Issue #30 e2e adaptation.

## Capabilities

### New Capabilities
- `sddops-operational-assets`: installs OPS sub-agents, slash commands, and OpenCode overlay assets so the five-phase OPS pipeline is executable across supported adapters.

### Modified Capabilities
- None.

## Approach

Create the missing asset files from the exploration’s affected-areas table and extend `internal/components/sddops/inject.go` with `injectOpsSubAgents`, `injectOpsSlashCommands`, `injectOpsOpenCodeOverlay`, and a private `mergeJSONFile`. Add `internal/assets/OpsCommandsAssetDir()` and keep all gating adapter-driven via `SupportsSubAgents()` / `SupportsSlashCommands()`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/components/sddops` | Modified | New inject helpers, overlay merge, idempotent post-checks, tests |
| `internal/assets` | Modified | `OpsCommandsAssetDir`, `claude/ops-commands/`, `opencode/ops-commands/`, `opencode/ops-overlay.json`, adapter sub-agent assets |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Kimi needs dual `.md` + `.yaml` assets | Medium | Mirror the proven old `sdd` pattern |
| Overlay prompts break if they use file paths | Medium | Inline prompt content only |
| Private `mergeJSONFile` copy drifts | Low | Copy the established package-local pattern unchanged |

## Rollback Plan

This change is additive and guarded by adapter capability checks. Rollback is reverting PR1→PR3; no migrations or state cleanup are required. Re-apply must remain idempotent.

## Dependencies

- Follow-up: Issue #30 e2e adaptation is BLOCKED by this change.

## Success Criteria

- [ ] Supported sub-agent adapters install `ops-brief` through `ops-deliver` assets correctly.
- [ ] Claude/OpenCode/Kilocode receive `/ops-*` commands from dedicated OPS command directories.
- [ ] OpenCode/Kilocode `opencode.json` gains `ops-orchestrator` plus 5 hidden sub-agents via the single overlay.
- [ ] Delivery stays PR-chained: PR1 assets only (`size:exception` documented), then PR2 inject features 1+2, then PR3 inject feature 3.
