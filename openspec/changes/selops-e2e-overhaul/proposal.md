# Proposal: selops-e2e-overhaul

## Intent

Raise TUI regression confidence by snapshot-testing every currently unsnapshotted OPS screen with the existing Go golden-file workflow.

## Problem Statement

Today only 13 golden files protect ~43 TUI screens, leaving ~30 screens without snapshot coverage; regressions in rendering and cross-screen happy-path flow can ship while `go test ./...` stays green.

## Scope

### In Scope
- Expand `internal/tui/testdata/*.golden` from 13 files to cover every currently unsnapshotted screen reached by the OPS fork, including install/complete, upgrade/sync, uninstall result, profiles, model pickers, and agent-builder screens.
- Add/extend model-level tests in `internal/tui/preset_flow_test.go`, `internal/tui/model_test.go`, and targeted `internal/tui/screens/*_test.go` files so each new screen snapshot is produced from `m.Update()` flows.
- Add one multi-step happy-path golden flow using `SelOpsOperational` preset and `PersonaOperator` defaults, snapshotting each major transition from welcome through completion.
- Keep the existing `-update` golden workflow and verify with `go test ./...`.

### Out of Scope
- Adding `teatest` or real `tea.Program` event-loop tests.
- Changing `e2e/e2e_test.sh`, Docker E2E, or bash harness behavior.
- Building Go subprocess/binary E2E coverage or real filesystem side-effect tests.

## Capabilities

### New Capabilities
None.

### Modified Capabilities
None.

## Proposed Solution

Use Approach 1 only: extend the proven golden-file pattern already implemented in `internal/tui/preset_flow_test.go`, add shared helpers only if needed, and make every new screen test deterministic through model state setup or explicit `tea.KeyMsg` progression.

## Affected Files

- `internal/tui/preset_flow_test.go`
- `internal/tui/model_test.go`
- `internal/tui/screens/*_test.go`
- `internal/tui/testdata/*.golden`
- `openspec/changes/selops-e2e-overhaul/proposal.md`

## Success Criteria

- [ ] Golden fixtures increase from 13 to cover all currently unsnapshotted OPS TUI screens.
- [ ] Every new screen snapshot is backed by a deterministic Go test runnable via `go test ./...`.
- [ ] A full happy-path OPS flow test snapshots each major transition through completion.
- [ ] CI stays green with no new test dependency added.

## Estimated Size

Forecast: ~700-1100 changed lines total (mostly new `.golden` fixtures); exceeds the default 400-line review budget, so chained review slices are likely.

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Golden churn from UI copy/style changes | Med | Update via `-update`; review fixture diffs carefully |
| Oversized review due to many goldens | High | Split by screen family / flow slice |
| Hidden async/runtime gaps remain | Med | Keep scope explicit; defer `teatest` to follow-up |

## Rollback Plan

Revert new tests and golden fixtures for the affected slice; no runtime code or user data migrations are introduced.

## Dependencies

- No technical blocker found.
- Input baseline is the existing OPS fork behavior, especially `SelOpsOperational` + `PersonaOperator` defaults.
