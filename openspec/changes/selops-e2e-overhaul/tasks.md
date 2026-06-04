# Tasks: selops-e2e-overhaul

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | ~1,420 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR1 Navigation → PR2 Install → PR3 Operation/Builder → PR4 Flow |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: pending
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|---|---|---|---|
| 1 | Add OPS helper + navigation goldens | PR 1 | validates naming and fixture pattern |
| 2 | Cover install/result screens | PR 2 | builds on helper and slice-1 conventions |
| 3 | Cover operation + agent-builder screens | PR 3 | split 3a/3b if review diff grows |
| 4 | Add deterministic happy-path flow | PR 4 | final integration proof |

## PR Slice 1 — Navigation screens
- [ ] T1.1 Add `newOpsTestModel`; `internal/tui/preset_flow_test.go`; centralize OPS defaults plus screen/cursor overrides; AC: helper returns deterministic OPS model and is reused by new tests; ~15 LOC; deps: none.
- [ ] T1.2 Create navigation golden suite; `internal/tui/navigation_golden_test.go`; add table-driven coverage for 10 navigation screens via `Update()`/`View()` and `assertTUIGolden`; AC: each target screen has one deterministic case; ~110 LOC; deps: T1.1.
- [ ] T1.3 Generate navigation goldens; `internal/tui/testdata/navigation-*.golden`; run `go test ./internal/tui -run TestNavigation -update`; AC: 10 committed goldens match kebab-case naming; ~180 LOC; deps: T1.2.
- [ ] T1.4 Verify slice-1 green; `internal/tui/navigation_golden_test.go`, `internal/tui/testdata/navigation-*.golden`; run targeted test then `go test ./...`; AC: no flaky snapshots and full suite passes; ~5 LOC; deps: T1.3.

## PR Slice 2 — Install screens
- [ ] T2.1 Create install golden suite; `internal/tui/install_golden_test.go`; cover model pickers, dependency tree, installing variants, complete variants, model config, plugin result using OPS helper; AC: 10 deterministic install/result cases exist; ~120 LOC; deps: T1.4.
- [ ] T2.2 Generate install goldens; `internal/tui/testdata/install-*.golden`; run `go test ./internal/tui -run TestInstall -update`; AC: 10 committed install goldens regenerate without manual edits; ~190 LOC; deps: T2.1.
- [ ] T2.3 Verify slice-2 green; `internal/tui/install_golden_test.go`, `internal/tui/testdata/install-*.golden`; run targeted tests and `go test ./...`; AC: suite stays green with no new deps; ~5 LOC; deps: T2.2.

## PR Slice 3 — Operation/progress screens
- [ ] T3.1 Create operation golden suite; `internal/tui/operation_golden_test.go`; cover upgrade/sync/backups/restore/delete/uninstall/profiles variants with deterministic state fixtures; AC: all non-builder operation screens from design are asserted; ~150 LOC; deps: T2.3.
- [ ] T3.2 Create agent-builder golden suite; `internal/tui/agent_builder_golden_test.go`; cover all 7 agent-builder screens with fixed cursor/state setup; AC: one committed case per builder screen; ~95 LOC; deps: T2.3.
- [ ] T3.3 Generate operation goldens; `internal/tui/testdata/operation-*.golden`; run targeted `-update`; AC: operation goldens commit cleanly and use success/error suffixes; ~220 LOC; deps: T3.1.
- [ ] T3.4 Generate agent-builder goldens; `internal/tui/testdata/agent-builder-*.golden`; run targeted `-update`; AC: 7 builder goldens commit with ordered names; ~120 LOC; deps: T3.2.
- [ ] T3.5 Verify slice-3 size and tests; `internal/tui/operation_golden_test.go`, `internal/tui/agent_builder_golden_test.go`, new goldens; run `go test ./...` and split 3a/3b if diff exceeds budget; AC: green suite and reviewable PR plan captured; ~10 LOC; deps: T3.3, T3.4.

## PR Slice 4 — Happy-path flow test
- [ ] T4.1 Add OPS happy-path test; `internal/tui/preset_flow_test.go`; create `TestInstallHappyPathFlow_OpsDefaults` with snapshot checkpoints only at major transitions; AC: test names 9 ordered snapshots from welcome to complete; ~70 LOC; deps: T1.1, T2.3, T3.5.
- [ ] T4.2 Define deterministic flow actions; `internal/tui/preset_flow_test.go`; encode key sequence, cursor moves, and synthetic `PipelineDoneMsg` payload matching `handlePipelineDone`; AC: no goroutines, no `Init()`, no time/random inputs; ~35 LOC; deps: T4.1.
- [ ] T4.3 Generate flow goldens; `internal/tui/testdata/flow-install-*.golden`; run `go test ./internal/tui -run TestInstallHappyPathFlow_OpsDefaults -update`; AC: 9 zero-padded flow goldens commit in order; ~90 LOC; deps: T4.2.
- [ ] T4.4 Verify deterministic flow; `internal/tui/preset_flow_test.go`, `internal/tui/testdata/flow-install-*.golden`; rerun targeted test and `go test ./...`; AC: snapshots are stable and final suite passes; ~5 LOC; deps: T4.3.
