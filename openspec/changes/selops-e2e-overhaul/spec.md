# Spec: selops-e2e-overhaul

## Overview

This change defines full spec coverage for OPS TUI golden testing because the proposal declares no existing capability mapping to modify. The system MUST expand the current Go golden-test pattern so every reachable OPS screen has deterministic snapshot coverage without adding dependencies or changing CI topology.

## Functional Requirements

### Requirement: FR-1 Full OPS screen golden coverage
The test suite MUST provide a deterministic golden-file test for every reachable OPS screen, including navigation, install, operation, profile, uninstall, model-picker, and agent-builder flows.

#### Scenario: Navigation screens are snapshotted
- GIVEN an OPS model configured for a reachable navigation screen
- WHEN the screen is rendered through the existing `m.Update()`-driven test pattern
- THEN a committed golden file SHALL exist for that screen state

#### Scenario: Error/result screens are snapshotted
- GIVEN an OPS model configured with an error or result state
- WHEN the screen is rendered
- THEN the matching golden file MUST capture the deterministic output

### Requirement: FR-2 OPS defaults are fixed for new goldens
Every new golden and flow test MUST use `SelOpsOperational` preset assumptions and `PersonaOperator` defaults unless the scenario explicitly verifies a different reachable OPS branch.

#### Scenario: Test setup uses OPS defaults
- GIVEN a new golden or flow test
- WHEN its model fixture is created
- THEN the fixture MUST initialize OPS-specific preset/persona defaults before rendering

### Requirement: FR-3 Install happy-path flow coverage
The suite MUST include a multi-step happy-path flow test covering install progression from welcome through completion, snapshotting each major transition.

#### Scenario: Install flow snapshots every major transition
- GIVEN a deterministic install path with mocked update and pipeline completion messages
- WHEN the test advances the model through welcome, detection, selection, review, installing, and complete states
- THEN each major state SHALL be asserted by committed golden output

### Requirement: FR-4 Suite remains runnable by standard Go tests
All added coverage MUST pass through `go test ./...` and MUST NOT require a second runner, Docker-only path, or runtime-only harness.

#### Scenario: Local and CI execution stay unchanged
- GIVEN the repository test command
- WHEN contributors or CI run `go test ./...`
- THEN the new golden and flow tests MUST run within that command

### Requirement: FR-5 Goldens remain regeneratable
The suite MUST preserve the existing `-update` regeneration pattern for any newly added golden files.

#### Scenario: Golden regeneration updates fixtures
- GIVEN intentional UI output changes
- WHEN a contributor runs the targeted Go test with `-update`
- THEN the committed golden files SHALL regenerate in `internal/tui/testdata/` without extra tooling

## Non-Functional Requirements

### Requirement: NFR-1 No new dependencies
The change MUST NOT add external modules or new `go.mod` entries.

### Requirement: NFR-2 CI compatibility
The change MUST keep existing `ci.yml` green without adding jobs.

### Requirement: NFR-3 Review-budget enforcement
Implementation delivery MUST ship as chained PR slices, each under 400 changed lines, ordered by dependency.

## BDD Coverage Matrix

### Scenario: Navigation screen group
- GIVEN OPS navigation states such as welcome, detection, agents, persona, preset, dependency tree, review, and skill/plugin pickers
- WHEN each state is rendered from deterministic model setup
- THEN each screen MUST have a stable golden assertion

### Scenario: Install screen group
- GIVEN install-path states such as model selection, installing, and complete
- WHEN the install path is advanced with deterministic messages
- THEN each install screen SHALL have committed golden coverage

### Scenario: Operation screen group
- GIVEN operation-path states such as upgrade, sync, upgrade-sync, backups, restore, delete, uninstall, profiles, and agent-builder progress screens
- WHEN each state is rendered in a deterministic fixture
- THEN the output MUST be protected by golden files

### Scenario: Error screen group
- GIVEN restore, delete, uninstall, plugin, or install result states with failure data
- WHEN the state is rendered
- THEN the screen MUST show deterministic error output captured by a golden

### Scenario: Multi-step install flow
- GIVEN `SelOpsOperational` and `PersonaOperator` defaults with deterministic mocked side effects
- WHEN the user path advances from welcome to completion
- THEN the flow test MUST snapshot each major transition in order

### Scenario: Golden regeneration workflow
- GIVEN approved UI copy or layout changes
- WHEN a contributor runs the targeted test with `-update`
- THEN only the intended goldens update and `go test ./...` returns green afterward

## PR Slice Plan

| Slice | Scope | Files Touched | Est. Lines | Depends On |
|---|---|---|---:|---|
| 1 | Core navigation + menu screens | `internal/tui/model_test.go`, `internal/tui/screens/*_test.go`, `internal/tui/testdata/navigation-*.golden` | 220-320 | none |
| 2 | Install-flow screens | `internal/tui/preset_flow_test.go`, targeted `screens/*_test.go`, install-related goldens | 260-340 | Slice 1 |
| 3 | Operation/progress screens | `internal/tui/model_test.go`, operation screen tests, progress/result goldens | 220-300 | Slice 2 |
| 4 | Happy-path multi-step flow | `internal/tui/preset_flow_test.go`, flow goldens, minimal helper reuse | 140-220 | Slices 1-3 |

## Acceptance Criteria

- All new golden files are committed, deterministic, and rooted in OPS defaults.
- `go test ./...` passes with no new dependencies and no new CI jobs.
- Existing 13 golden files remain non-regressed unless intentionally updated via reviewed fixture diffs.
- CI is green for every chained slice and for the final integrated branch.
