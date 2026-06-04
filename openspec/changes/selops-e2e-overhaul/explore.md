# Exploration: selops-e2e-overhaul

## Current State

### What exists today

#### Unit/TUI tests (154 Go test files — `go test ./...`)
The dominant testing strategy. Two distinct flavours:

1. **Pure unit tests** — Logic-only packages: `internal/model`, `internal/planner`, `internal/pipeline`, `internal/verify`, `internal/catalog`, `internal/update`, `internal/backup`, `internal/system`, etc.

2. **Bubbletea model-level TUI tests** — `internal/tui/model_test.go`, `internal/tui/preset_flow_test.go`, `internal/tui/restore_test.go`, `internal/tui/agent_builder_nav_test.go`, and all `internal/tui/screens/*_test.go` (≈40+ files).
   - Pattern: construct a `Model`, send `tea.KeyMsg`/`tea.Msg` via `m.Update(msg)`, assert on resulting `Screen`, `Selection`, or `Progress` fields.
   - **No `teatest` library** — zero `github.com/charmbracelet/teatest` import in `go.mod`.
   - Uses **golden files** for screen `View()` output (`internal/tui/testdata/*.golden`), 13 files currently.
   - These tests are **fast** (no real TTY, no I/O), deterministic, and already cover most navigation flows.

#### E2E shell tests (`e2e/`)
- `e2e_test.sh` + `lib.sh`: a bash test harness with 70+ named test functions.
- **Tier 1** (always runs): binary exists, `--dry-run` output format, agent/preset/component flags, invalid input rejection. Fast, no side-effects.
- **Tier 2** (`RUN_FULL_E2E=1`): actual install invocations that write to `$HOME` — verifies Claude Code, OpenCode, Qwen injections for every component (engram, SDD, persona, skills, context7, permissions, theme), preset integration, idempotency.
- **Tier 3** (`RUN_BACKUP_TESTS=1`): backup/restore cycle.
- Runs inside Docker on Ubuntu, Arch, Fedora (`e2e/Dockerfile.ubuntu`, `.arch`, `.fedora`).
- CI: PR → Tier 1 only; push to main / nightly → Tier 1+2+3.

#### Integration tests (limited)
- `internal/agentbuilder/integration_test.go`: AI generation+install cycle with real file system.
- `internal/app/parity_test.go`: flag parity between CLI and TUI flows.
- No HTTP-level integration tests.

### What is NOT tested

1. **Full TUI user journeys end-to-end** — no test drives Welcome → Detection → Agents → … → Complete as a single run of the binary. The model-level tests cover individual screen transitions but they bypass the real `tea.Program` loop, the `Init()` update-check goroutine, the real pipeline `ExecuteFn`, and the actual file-system side-effects.

2. **Real `tea.Program` execution** — no test starts `bubbletea.NewProgram(model)` and feeds simulated key events through the program's event loop. `teatest` (the charmbracelet library for this) is absent from `go.mod`.

3. **TUI → pipeline integration** — the `ExecuteFn` is always swapped out in tests. No test drives a real install through the TUI *and* verifies the resulting file writes.

4. **TUI → backup/restore integration** — `RestoreFn`, `DeleteBackupFn`, `ListBackupsFn` are all faked. No test runs the full restore flow and confirms disk state.

5. **TUI output fidelity at scale** — only 13 golden files; many screens (Installing, Complete, Upgrade, Sync, UninstallResult, Profiles, ModelPicker…) have no golden snapshot.

6. **Cross-screen regression** — a change to one screen's `View()` can silently corrupt an adjacent screen with no test catching it.

---

## Affected Areas

- `e2e/e2e_test.sh` — main E2E harness; would be extended (Tier 4 Go-based tests) or restructured
- `e2e/lib.sh` — shared assertion helpers for shell harness
- `e2e/Dockerfile.ubuntu|arch|fedora` — Docker E2E environment
- `internal/tui/model.go` — main TUI model; all TUI tests touch this
- `internal/tui/model_test.go` — existing model-level tests; golden extension target
- `internal/tui/preset_flow_test.go` — golden file + flow matrix tests; extend for missing screens
- `internal/tui/screens/*.go` + `*_test.go` — individual screen renderers and their unit tests
- `internal/tui/testdata/*.golden` — 13 golden files; extend to cover remaining screens
- `go.mod` — would need `teatest` added if Approach B is chosen
- `.github/workflows/ci.yml` — CI configuration; may need new job for Go-based E2E

---

## Approaches

### Approach 1 — Golden File Coverage Expansion (Low effort, high ROI)

**What**: Extend the existing golden-file pattern to cover all screens that currently have no snapshot. Add flow-level multi-step golden tests for the full happy path (Welcome → … → Complete) and the main error paths.

- Add ~20–30 golden files for: Installing, Complete, Upgrade/Sync/UpgradeSync, UninstallResult, Profiles, ModelPicker, AgentBuilder screens.
- Add multi-step flow tests (e.g., `TestFullInstallHappyPath_GoldenFlow`) that simulate the whole screen progression and snapshot each intermediate view.
- NO new dependencies; NO `teatest`; runs in `go test ./...`.

**Pros**:
- Zero new dependencies.
- Fits existing test infrastructure perfectly.
- Regression safety for View() changes.
- Fast — sub-second.
- Strict TDD compatible.

**Cons**:
- Still does not test the real `tea.Program` event loop or real filesystem side-effects.
- Golden files can become stale; need `-update` flag discipline.
- Does not catch goroutine/concurrency bugs in `Init()` or async message flows.

**Effort**: Low (1–2 sessions)

---

### Approach 2 — Add `teatest` for Program-Level TUI E2E

**What**: Add `github.com/charmbracelet/teatest` to `go.mod`. Write `*_test.go` files that start a real `tea.Program` with a fake TTY, send keyboard events via `teatest.WaitFor` + `tt.Send()`, and assert on both output snapshots and model state.

- Covers the real program lifecycle: `Init()`, async `UpdateCheckResultMsg`, `PipelineDoneMsg`, `UpgradeDoneMsg`, etc.
- Can test concurrent goroutine flows (spinner ticking, pipeline progress).

**Pros**:
- Tests the actual event loop, not just `Update()` in isolation.
- Catches race conditions and timer-based bugs.
- Official charmbracelet recommendation for TUI E2E.

**Cons**:
- Adds an external dependency (`teatest`).
- Slightly more complex test setup (fake TTY dimensions, timing).
- Some tests may be flaky on slow CI without careful `WaitFor` timeouts.
- Learning curve on `teatest` API.

**Effort**: Medium (2–3 sessions)

---

### Approach 3 — Go-Based Binary E2E (Replace / Augment Shell Harness)

**What**: Write a Go-based E2E test package (e.g., `e2e_go/`) that builds the binary and invokes it as a subprocess with controlled args, capturing stdout/stderr. Port or complement the existing bash Tier 1+2 tests into Go using `os/exec` + `testing`.

- Makes E2E tests first-class Go citizens: runnable with `go test ./e2e_go/...`, integrated into `go test ./...`, usable with coverage tools.
- Can use `t.TempDir()`, `t.Setenv()`, parallel subtests, and Go assertions instead of bash regex matching.

**Pros**:
- Unified test command (`go test ./...` covers everything).
- Better error messages than bash assertions.
- Easier to maintain than bash.
- Can be extended to cover real filesystem assertions per component.

**Cons**:
- Does not test TUI interaction (keyboard/screen) — only CLI/binary behavior.
- Duplicates or replaces the existing shell harness; migration cost.
- Docker-based cross-platform execution still needed separately.

**Effort**: Medium (2–3 sessions)

---

## Recommendation

**Start with Approach 1 (golden file coverage expansion)** as an immediate win. It is zero-risk, builds on proven infrastructure, and systematically closes the largest gap (missing View() regression coverage for ≈20 screens). Estimate: 1–2 sessions.

**Then** evaluate Approach 2 (teatest) for the next phase: it solves the one gap Approach 1 cannot — the real program event loop with async messages. The OPS fork's `Init()` fires a real HTTP update-check; `teatest` is the right tool to test that path without network calls.

**Approach 3** is low priority: the bash harness already covers binary/CLI behavior well; a Go port would be useful debt reduction but is not blocking anything.

---

## Risks

- **Golden file drift**: if screen output changes frequently (new features, styling updates), golden files need updating. Must enforce `-update` flag discipline + PR review for golden diffs.
- **`teatest` flakiness**: timing-sensitive tests can flake on slow CI. Mitigate with generous `WaitFor` timeouts and `t.Parallel()` isolation.
- **Scope creep**: "overhaul" risks turning into a multi-week effort. Bound it to Approach 1 first; gate Approach 2 behind a separate SDD change.
- **Screen count**: `model.go` defines 43 screen constants. Full golden coverage is comprehensive; prioritize the critical-path screens (install flow + operation screens) first.
- **OPS fork specifics**: several screens are OPS-specific (`SelOpsOperational` preset, `SDDOps` component, operator persona). Tests must not assume upstream `gentle-ai` defaults.

---

## Ready for Proposal

Yes — ready for `sdd-propose` phase.

Recommended scope for the proposal: **Approach 1 only** (golden file coverage expansion for all unsnapshotted screens + multi-step happy-path flow tests). Defer `teatest` and Go-based binary E2E to follow-on changes.

Key decision for the proposal to resolve: which screens are in-scope for the first iteration vs. deferred.
