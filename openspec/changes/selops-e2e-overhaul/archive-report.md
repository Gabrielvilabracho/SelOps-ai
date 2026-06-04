# Archive Report — selops-e2e-overhaul

## Status
**completed** | Archived: 2026-06-04

## Summary
Expanded TUI golden file coverage for SelOps-ai from 13 files (covering only preset-flow and custom-flow screens) to 80 files covering all 44 reachable OPS screens plus a 9-step Install happy-path flow test. The change introduced a shared `newOpsTestModel` test helper, established deterministic rendering conventions for OPS-fork screens, and delivered 6 chained PRs each under the 400-line review budget.

## Deliverables

### Pull Requests
| PR | Slice | Scope | LOC (test only) |
|----|-------|-------|-----------------|
| #18 | 1 | Navigation screens (10 goldens) | 100 |
| #19 | 2 | Install flow screens (10 goldens) | 169 |
| #23 | 3a | Upgrade/Sync/Backups/Restore (15 goldens) | 234 |
| #24 | 3b | Uninstall/Profiles (15 goldens) | 221 |
| #25 | 3c | AgentBuilder screens (10 goldens) | 179 |
| #22 | 4 | Happy-path Install flow (9 goldens) | 116 |

### Files Created
- `internal/tui/navigation_golden_test.go`
- `internal/tui/install_golden_test.go`
- `internal/tui/operation_upgrade_sync_golden_test.go`
- `internal/tui/operation_uninstall_profiles_golden_test.go`
- `internal/tui/agent_builder_golden_test.go`
- `internal/tui/testdata/navigation-*.golden` (10 files)
- `internal/tui/testdata/install-*.golden` (10 files)
- `internal/tui/testdata/operation-*.golden` (28 files)
- `internal/tui/testdata/agent-builder-*.golden` (10 files)
- `internal/tui/testdata/flow-install-*.golden` (9 files)

### Files Modified
- `internal/tui/preset_flow_test.go` — added `newOpsTestModel` helper + `TestInstallHappyPathFlow_OpsDefaults`

## Spec Compliance

### Functional Requirements
| ID | Requirement | Result |
|----|-------------|--------|
| FR-1 | All OPS screens have deterministic golden test | ✅ PASS — 44/44 screens covered |
| FR-2 | SelOpsOperational + PersonaOperator defaults in all tests | ✅ PASS — via newOpsTestModel helper |
| FR-3 | Multi-step happy-path flow test with 9 snapshots | ✅ PASS — TestInstallHappyPathFlow_OpsDefaults |
| FR-4 | go test ./... green, no new dependencies | ✅ PASS — 52 packages, 0 failures |
| FR-5 | Golden files regeneratable with -update flag | ✅ PASS — existing mechanism preserved |

### Non-Functional Requirements
| ID | Requirement | Result |
|----|-------------|--------|
| NFR-1 | No new external dependencies in go.mod | ✅ PASS |
| NFR-2 | CI stays green, no new jobs needed | ✅ PASS |
| NFR-3 | Each PR slice under 400 LOC | ✅ PASS — required splitting operation slice into 3a/3b |

## Lessons Learned

### Patterns Established
- `newOpsTestModel(t testing.TB, screen Screen, cursor int) Model` — canonical OPS test setup; lives in `preset_flow_test.go`
- `-update` flag must scope to `go test ./internal/tui` — sub-packages (`screens/`, `styles/`) don't define the flag and fail when it propagates via `./internal/tui/...`
- OPS-fork screens (`PersonaOperator`, `PresetSelOpsOperational`, `ScreenDependencyTree`) are not reachable via `Update()` keypresses because they are not in the standard picker options — use direct model field injection and document with a comment
- `PipelineDoneMsg` can be injected synthetically to simulate install completion without goroutines or HTTP calls
- Never call `Init()` in TUI tests — it spawns goroutines and HTTP calls; use `Update()` and `View()` only
- Installing screen state: `m.Progress.Start(0)` for running; `m.Progress.Mark(i, "succeeded")` for done
- Complete screen failure path: set `m.Execution.Apply.Steps` with a `StepStatusFailed` entry

### Process Lessons
- Golden content (fixture files) inflates GitHub diff line counts — define review budget in terms of test code LOC, not total diff
- Splitting a 28-case table-driven test across two files by domain (upgrade/sync vs uninstall/profiles) is a natural seam that keeps review cohesion
- `testManifest()` helper in one file is accessible cross-file within the same Go package — no need to duplicate

## Follow-on Changes (deferred)
1. **teatest integration** (Approach 2 from exploration) — add `charmbracelet/x/exp/teatest` to cover the real `tea.Program` event loop, async messages, and goroutine safety. Deferred from this change.
2. **Go-based binary E2E** (Approach 3) — replace or augment the bash E2E harness with `os/exec` Go tests. Lower priority.

## Golden File Count
- Before: 13
- After: **80** (67 new)
