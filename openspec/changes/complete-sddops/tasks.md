# Tasks: Complete sddops

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | 520-680 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR1 assets → PR2 sub-agents+commands → PR3 overlay |
| Delivery strategy | pending |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: pending
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|---|---|---|---|
| 1 | Drop OPS assets only | PR 1 | buildable no-op, size:exception |
| 2 | Add sub-agent + command injection | PR 2 | depends on PR1, TDD |
| 3 | Add OpenCode overlay merge | PR 3 | depends on PR1, TDD |

## PR1 — assets only

### Infrastructure
- [x] 1.1 Add `OpsCommandsAssetDir` in `internal/assets/commands.go` for Claude/OpenCode routing.

### Implementation
- [x] 1.2 Create `internal/assets/claude/agents/ops-*.md` and `internal/assets/cursor/agents/ops-*.md`.
- [x] 1.3 Create `internal/assets/kimi/agents/ops-*.md`, `internal/assets/kimi/agents/ops-*.yaml`, and `internal/assets/kiro/agents/ops-*.md`.
- [x] 1.4 Create `internal/assets/claude/ops-commands/ops-*.md` and `internal/assets/opencode/ops-commands/ops-*.md`.
- [x] 1.5 Create `internal/assets/opencode/ops-overlay.json` with `{{OPS_ORCHESTRATOR_PROMPT}}` and 5 subagents.

### Testing
- [x] 1.6 Add asset-shape test in `internal/assets/commands_test.go` for frontmatter/sentinel inventory (scenarios 5,10,11,12,15).

### Definition of Done
- [x] PR1: `go build ./...`, `go test ./...`, `go vet ./...`.
- [x] PR1: 36 asset files + `OpsCommandsAssetDir` present.

## PR2 — inject features 1+2

### Infrastructure
- [ ] 2.1 Extend `internal/components/sddops/inject.go` `InjectOptions` with Claude/Kiro model-assignment maps and local resolver interfaces.

### Implementation
- [ ] 2.2 RED: add failing `injectOpsSubAgents` tests in `internal/components/sddops/inject_test.go` for scenarios 1,5,6,7,8,9.
- [ ] 2.3 GREEN: implement `injectOpsSubAgents` in `internal/components/sddops/inject.go` using atomic writes and placeholder replacement.
- [ ] 2.4 REFACTOR: trim shared sub-agent fixtures/helpers in `internal/components/sddops/inject_test.go` if duplication appears.
- [ ] 2.5 RED: add failing `injectOpsSlashCommands` tests in `internal/components/sddops/inject_test.go` for scenarios 2,10,11,12,13,14.
- [ ] 2.6 GREEN: implement `injectOpsSlashCommands` in `internal/components/sddops/inject.go` via `assets.OpsCommandsAssetDir`.
- [ ] 2.7 REFACTOR: normalize command-copy assertions/helpers in `internal/components/sddops/inject.go` and `inject_test.go`.

### Testing
- [ ] 2.8 Wire both helpers into `internal/components/sddops/inject.go` and add cross-cutting `Inject()` tests in `inject_test.go` for scenarios 3 and 4.

### Definition of Done
- [ ] PR2: RED→GREEN→REFACTOR order preserved for both functions.
- [ ] PR2: `go test ./...`, `go vet ./...`, scenarios 1-14 covered.

## PR3 — inject feature 3

### Infrastructure
- [ ] 3.1 Prep JSON merge fixtures/helpers in `internal/components/sddops/inject_test.go` for existing-settings cases.

### Implementation
- [ ] 3.2 RED: add failing overlay tests in `internal/components/sddops/inject_test.go` for scenarios 15,16,17,18,19.
- [ ] 3.3 GREEN: add private `mergeJSONFile` to `internal/components/sddops/inject.go`, returning merged bytes.
- [ ] 3.4 GREEN: implement `injectOpsOpenCodeOverlay` in `internal/components/sddops/inject.go` with prompt inlining and semantic post-check.
- [ ] 3.5 REFACTOR: collapse JSON/assert helpers in `internal/components/sddops/inject.go` and `inject_test.go` if needed.

### Testing
- [ ] 3.6 Wire overlay into `internal/components/sddops/inject.go` and extend `Inject()` tests in `inject_test.go` for scenarios 3,4,16,19.

### Definition of Done
- [ ] PR3: `go test ./...`, `go vet ./...`, scenarios 15-19 covered.
- [ ] PR3: idempotency and preserve-keys checks pass for overlay.
