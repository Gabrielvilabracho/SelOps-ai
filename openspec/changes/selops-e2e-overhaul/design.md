# Design: selops-e2e-overhaul

## Executive Summary

Extend the existing Go golden-file pattern in `internal/tui/` to cover every reachable OPS TUI screen via deterministic `Model.View()` snapshots. The work is delivered as 4 chained PR slices grouped by screen family (Navigation, Install, Operation, Flow), each under 400 changed lines. A single tiny test helper (`newOpsTestModel`) is introduced in `internal/tui/preset_flow_test.go` to canonicalize OPS defaults across all new tests; the existing `-update` golden mechanism and `assertTUIGolden` helper are preserved verbatim. No new dependencies, no new CI jobs, no `teatest` — `go test ./...` covers everything.

## Helper / Harness Design

### Test model construction

`NewModel(system.DetectionResult{}, "dev")` already returns a model pre-configured with `PresetSelOpsOperational` + `PersonaOperator` (see `internal/tui/model.go:471-503`). That is the canonical OPS starting state.

A single shared helper is added to `internal/tui/preset_flow_test.go` (alongside `assertTUIGolden`, `applyFlowAction`, `presetCursor`):

```go
// newOpsTestModel returns a Model initialised with OPS defaults
// (SelOpsOperational preset + PersonaOperator persona). It is the
// canonical starting point for every new TUI golden test.
//
// Optional overrides:
//   - screen: target Screen to land on (default ScreenWelcome)
//   - cursor: cursor position (default 0)
func newOpsTestModel(t *testing.T, screen Screen, cursor int) Model {
    t.Helper()
    m := NewModel(system.DetectionResult{}, "dev")
    if screen != ScreenUnknown {
        m.Screen = screen
    }
    m.Cursor = cursor
    return m
}
```

Rationale:
- Keeps test files thin and self-documenting.
- Centralizes the OPS-defaults assumption (FR-2) in one location so future preset changes propagate uniformly.
- Avoids a new `testutil` package — the helper lives in the same package as the tests that use it.

### Key-sequence + assertion pattern

All new tests follow the established pattern already used in `TestPresetSelectionNextScreenFlowMatrix`:

```go
m := newOpsTestModel(t, ScreenXXX, cursor)
// optional: pre-populate m fields the screen reads
updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
state := updated.(Model)
if state.Screen != wantScreen { t.Fatalf(...) }
assertTUIGolden(t, "name.golden", state.View())
```

For multi-step flows, reuse the existing `applyFlowAction` helper.

### Golden assertion mechanism

Unchanged. `assertTUIGolden` (preset_flow_test.go:243) and the `-update` flag (preset_flow_test.go:15) are reused. No new mechanism is introduced.

## Screen Grouping per PR Slice

`internal/tui/model.go:240-284` defines 43 screen constants. Of these, 13 are already covered by goldens (preset/custom variants). The remaining ~30 screens are partitioned into 4 slices below.

### Slice 1 — Navigation & Menu Screens (~220-320 LOC)

Screens added:
- `ScreenWelcome`
- `ScreenDetection`
- `ScreenAgents`
- `ScreenPersona`
- `ScreenPreset` (current-state snapshot, distinct from existing next-screen goldens)
- `ScreenSDDMode`
- `ScreenStrictTDD`
- `ScreenOpenCodePlugins` (default state, not post-toggle)
- `ScreenSkillPicker`
- `ScreenReview`

Test file:
- New: `internal/tui/navigation_golden_test.go` (table-driven, one `t.Run` per screen)

Goldens added in `internal/tui/testdata/`:
- `navigation-welcome.golden`
- `navigation-detection.golden`
- `navigation-agents.golden`
- `navigation-persona.golden`
- `navigation-preset.golden`
- `navigation-sdd-mode.golden`
- `navigation-strict-tdd.golden`
- `navigation-opencode-plugins.golden`
- `navigation-skill-picker.golden`
- `navigation-review.golden`

### Slice 2 — Install Flow Screens (~260-340 LOC)

Screens added:
- `ScreenClaudeModelPicker`
- `ScreenKiroModelPicker`
- `ScreenModelPicker`
- `ScreenDependencyTree` (full deterministic plan; distinct from existing custom-flow goldens)
- `ScreenInstalling` (in-progress + done states)
- `ScreenComplete` (success + with-failed-steps variants)
- `ScreenModelConfig`
- `ScreenOpenCodePluginResult`

Test file:
- New: `internal/tui/install_golden_test.go`

Goldens added:
- `install-claude-model-picker.golden`
- `install-kiro-model-picker.golden`
- `install-model-picker.golden`
- `install-dependency-tree.golden`
- `install-installing-running.golden`
- `install-installing-done.golden`
- `install-complete-success.golden`
- `install-complete-with-failures.golden`
- `install-model-config.golden`
- `install-opencode-plugin-result.golden`

Setup notes:
- `ScreenInstalling`: pre-populate `m.Progress = NewProgressState([...])` and call `m.Progress.Start(0)`; assert with `SpinnerFrame=0` for determinism.
- `ScreenComplete`: set `m.Execution` to a deterministic `pipeline.ExecutionResult` for both happy and failed cases.

### Slice 3 — Operation / Result Screens (~220-300 LOC)

Screens added:
- `ScreenUpgrade` (idle + check-done + post-run report variants)
- `ScreenSync` (confirm + post-run variants)
- `ScreenUpgradeSync`
- `ScreenBackups`
- `ScreenRestoreConfirm`
- `ScreenRestoreResult` (success + error)
- `ScreenDeleteConfirm`
- `ScreenDeleteResult` (success + error)
- `ScreenRenameBackup`
- `ScreenUninstallMode`
- `ScreenUninstall`
- `ScreenUninstallComponents`
- `ScreenUninstallProfiles`
- `ScreenUninstallConfirm`
- `ScreenUninstallResult` (success + error)
- `ScreenProfiles`
- `ScreenProfileCreate`
- `ScreenProfileDelete`
- `ScreenAgentBuilderEngine`
- `ScreenAgentBuilderPrompt`
- `ScreenAgentBuilderSDD`
- `ScreenAgentBuilderSDDPhase`
- `ScreenAgentBuilderGenerating`
- `ScreenAgentBuilderPreview`
- `ScreenAgentBuilderInstalling`
- `ScreenAgentBuilderComplete`

Test files:
- New: `internal/tui/operation_golden_test.go` (backups + upgrade + sync + uninstall + profiles)
- New: `internal/tui/agent_builder_golden_test.go` (all `ScreenAgentBuilder*`)

Goldens added in `internal/tui/testdata/` (naming convention `<family>-<screen>-<variant>.golden`):
- `upgrade-idle.golden`, `upgrade-checking.golden`, `upgrade-report.golden`
- `sync-confirm.golden`, `sync-result.golden`
- `upgrade-sync-report.golden`
- `backups-list.golden`
- `restore-confirm.golden`, `restore-result-success.golden`, `restore-result-error.golden`
- `delete-confirm.golden`, `delete-result-success.golden`, `delete-result-error.golden`
- `rename-backup.golden`
- `uninstall-mode.golden`, `uninstall-agents.golden`, `uninstall-components.golden`, `uninstall-profiles.golden`, `uninstall-confirm.golden`, `uninstall-result-success.golden`, `uninstall-result-error.golden`
- `profiles-list.golden`, `profile-create.golden`, `profile-delete.golden`
- `agent-builder-engine.golden`, `agent-builder-prompt.golden`, `agent-builder-sdd-mode.golden`, `agent-builder-sdd-phase.golden`, `agent-builder-generating.golden`, `agent-builder-preview.golden`, `agent-builder-installing.golden`, `agent-builder-complete.golden`

Note: Slice 3 is the largest in screen count. If it exceeds 400 LOC during implementation, split as 3a (backups/restore/uninstall) and 3b (profiles + agent builder). The spec already allows this kind of split under NFR-3.

### Slice 4 — Install Happy-Path Multi-Step Flow (~140-220 LOC)

Test file:
- Extension of `internal/tui/preset_flow_test.go` — add `TestInstallHappyPathFlow_OpsDefaults`.

Goldens added (one per major transition):
- `flow-install-01-welcome.golden`
- `flow-install-02-detection.golden`
- `flow-install-03-agents.golden`
- `flow-install-04-persona.golden`
- `flow-install-05-preset.golden`
- `flow-install-06-dependency-tree.golden`
- `flow-install-07-review.golden`
- `flow-install-08-installing.golden`
- `flow-install-09-complete.golden`

The flow test design is detailed in the next section.

## Multi-Step Install Flow Test Design

### Test location

Extend `internal/tui/preset_flow_test.go` rather than create a new file. The file already hosts flow-style tests (`TestPresetSelectionNextScreenFlowMatrix`, `TestCustomPresetPostComponentFlowMatrix`) and owns the golden helper.

### Test structure

```go
func TestInstallHappyPathFlow_OpsDefaults(t *testing.T) {
    m := newOpsTestModel(t, ScreenWelcome, 0)

    // Stage 1: Welcome — snapshot then advance
    assertTUIGolden(t, "flow-install-01-welcome.golden", m.View())
    m = stepEnter(t, m, ScreenDetection)

    // Stage 2: Detection
    assertTUIGolden(t, "flow-install-02-detection.golden", m.View())
    m = stepEnter(t, m, ScreenAgents)

    // ... and so on through Complete
}
```

### Key sequence

Driven entirely by `tea.KeyMsg{Type: tea.KeyEnter}` with cursor positioning between stages where required (preset selection, agent selection). No mouse events, no async messages.

### Installing → Complete transition

`ScreenInstalling` is the only stage that normally requires an async `PipelineDoneMsg`. To keep the test deterministic and avoid touching `ExecuteFn`:

```go
// At ScreenInstalling: feed a synthetic PipelineDoneMsg directly.
result := pipeline.ExecutionResult{
    Prepare: pipeline.PhaseResult{Steps: []pipeline.StepResult{...all succeeded...}},
    Apply:   pipeline.PhaseResult{Steps: []pipeline.StepResult{...all succeeded...}},
}
updated, _ := m.Update(PipelineDoneMsg{Result: result})
m = updated.(Model)
// Then press Enter to move from ScreenInstalling to ScreenComplete.
```

This mirrors `handlePipelineDone` (model.go:700) without invoking the real executor. The pattern is already used in `model_test.go` (see line 467+).

### Snapshot policy

One golden per major transition listed in Slice 4 above (9 total). NOT one per cursor move. NOT one per spinner tick.

### Determinism guarantees

- `Version` set to `"dev"` (already done by helper).
- `Detection` left at zero value — `system.DetectionResult{}` is deterministic.
- `SpinnerFrame` explicitly set to 0 before snapshotting Installing screen.
- `m.Execution` populated via synthetic result, no time-dependent fields used.
- No `time.Now()`, no random IDs, no goroutines.

## Golden File Conventions

### Naming

- Prefix by family: `navigation-`, `install-`, `upgrade-`, `sync-`, `backups-`, `restore-`, `delete-`, `rename-`, `uninstall-`, `profiles-`, `profile-`, `agent-builder-`, `flow-install-`.
- Suffix with variant when a screen has multiple states: `-success`, `-error`, `-idle`, `-running`, `-done`, `-confirm`, `-result`, `-checking`, `-report`.
- Multi-step flow goldens use a zero-padded ordinal: `flow-install-01-welcome.golden`.
- Existing `*-next.golden` files are NOT renamed — they continue to represent "next-screen-after-action" snapshots and are preserved as-is.
- All filenames are lowercase-kebab-case to match existing convention.

### Masking strategy for dynamic content

The View() pipeline is already deterministic for OPS defaults:
- Version: passed as `"dev"` via `newOpsTestModel`.
- Update banners: gated by `m.UpdateCheckDone` — left false in tests, banner is empty.
- Spinner: snapshotted at `SpinnerFrame=0` only.
- Pipeline timestamps: `pipeline.ExecutionResult` does not expose them in `RenderInstalling` / `RenderComplete` output.

If any future screen surfaces non-deterministic data (e.g. real timestamps, paths under `$HOME`), the masking rule is: **the test sets the underlying model field to a fixed deterministic value before calling `View()`**. No string post-processing, no regex stripping in the golden helper. The current 13 goldens already follow this rule.

### Update workflow

Unchanged: `go test ./internal/tui/... -update -run TestName`. The `-update` flag is package-scoped; running it regenerates only the goldens whose tests match the `-run` selector.

## CI Impact

- No new CI jobs.
- `go test ./...` (already wired in `.github/workflows/ci.yml`) covers all new tests automatically.
- Golden files are committed to the repository — CI never auto-generates them. A diff in `internal/tui/testdata/` is treated like any other source change.
- Test runtime budget: each new golden test is `< 5ms` (no I/O beyond reading the golden file). All 30+ new tests together add `< 200ms` to the existing suite.

## Risk Mitigations

| Risk | Mitigation |
|---|---|
| Golden churn from UI copy/style changes | (1) Slices are reviewed per family so a copy change touches a small bounded set of goldens. (2) PR reviewers diff golden files explicitly. (3) `-update` regeneration is per-test, never global. |
| Oversized review | Strict 4-slice partition; each slice ≤ 400 LOC. Slice 3 has a pre-planned split (3a/3b) if it overshoots. |
| Test brittleness from cursor positions | All cursor values are computed from helpers (`presetCursor`, `len(screens.AgentOptions())`) rather than hard-coded indices — already the established pattern. |
| Slice interdependency stalls review | Slices 1-3 are independent screen families that read pre-set model fields. Slice 4 depends conceptually on Slices 1-2 (it transits screens covered there) but does NOT depend on their test code — it can be authored in parallel and merged after them. |
| Hidden non-determinism in `View()` | Mitigated by the strict OPS-defaults helper (`newOpsTestModel`) plus the model-field-setting masking rule. Every new test runs from a known fixed state. |
| `Init()` HTTP update-check leaks into tests | Tests never call `m.Init()` — they construct the model via `NewModel(...)` and feed `m.Update(...)` directly. This is already how the existing 13 goldens work. |

## Next Recommended Tasks

1. **Tasks phase** (`sdd-tasks`): break each slice into individual implementable tasks with file targets and acceptance criteria, in dependency order (Slice 1 → 2 → 3 → 4).
2. Author Slice 1 first to validate the helper + naming conventions land cleanly before scaling out.
3. Land Slice 1 as a standalone PR, verify CI is green, then start Slice 2.
