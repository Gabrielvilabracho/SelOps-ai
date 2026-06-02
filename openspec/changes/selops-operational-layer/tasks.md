# Tasks: SelOps Operational Layer

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 700–1,100 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR1 registration+preset, PR2 new packages+wiring, PR3 assets+regression |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

Decision needed before apply: Yes
Chained PRs recommended: Yes
Chain strategy: pending
400-line budget risk: High

## MVP decisions fixed in tasks

| Topic | MVP decision |
|-------|--------------|
| TUI surfacing | CLI `--preset selops-operational` only; no TUI preset row in this pass. |
| OperationalMCPServers source | `selection.OperationalMCPServers` defaults empty; write disabled placeholder entry. |
| Sync semantics | Reuse install-only `sddops.Inject` path for MVP; no `InjectForSync` yet. |
| Empty-asset guard | Add ≥100-byte post-check coverage so hollow `ops-*` stubs fail tests. |

## Phase 1: Infrastructure / registration

- [ ] 1.1 (TEST) Add planner/preset RED tests in `internal/planner/resolver_test.go`, `internal/planner/order_test.go`, and `internal/components/skills/presets_test.go` for `PresetSelOpsOperational`, `ComponentSDDOps`, `ComponentOperationalMCP`, and Persona-before-SDDOps ordering.
- [ ] 1.2 Add registration constants and catalogs in `internal/model/types.go`, `internal/catalog/components.go`, and `internal/catalog/skills.go`, including the 6 `ops-*` skill IDs.
- [ ] 1.3 Wire the preset and graph in `internal/components/skills/presets.go`, `internal/planner/graph.go`, `internal/cli/validate.go`, and preset/help surfaces for CLI-only MVP; document TUI row as deferred.

## Phase 2: Implementation / new packages and runtime wiring

- [ ] 2.1 (TEST) Add RED unit tests for `internal/components/sddops/` covering exact `sddOpsSkillIDs` membership and zero intersection with `internal/components/sdd` skill IDs.
- [ ] 2.2 Create `internal/components/sddops/inject.go` and helpers to call `skills.InjectWithCapability(...)`; keep MVP install-only and record sync follow-up.
- [ ] 2.3 (TEST) Add RED unit tests for `internal/components/operationalmcp/` verifying empty `OperationalMCPServers` writes the documented disabled placeholder and preserves existing JSON keys via merge.
- [ ] 2.4 Extend `internal/model/selection.go` with `OperationalMCPServers`, then create `internal/components/operationalmcp/inject.go` using existing MCP-strategy/filemerge paths.
- [ ] 2.5 (TEST) Add RED CLI/runtime tests in `internal/cli/run_component_paths_test.go`, `internal/cli/run_integration_test.go`, and `internal/cli/sync_test.go` for new component cases, placeholder path, and MVP sync exclusion.
- [ ] 2.6 Add `componentApplyStep` and path/backup wiring in `internal/cli/run.go`; add sync behavior in `internal/cli/sync.go` only as explicit non-managed/no-op for `ComponentSDDOps`.

## Phase 3: Assets / placeholders only

- [ ] 3.1 (TEST) Add RED asset/post-check coverage in `internal/components/sddops/inject_test.go` and/or `internal/assets/assets_test.go` for the ≥100-byte `SKILL.md` guard mirroring `internal/components/sdd/inject.go:713`.
- [ ] 3.2 Add 6 layout stubs at `internal/assets/skills/ops-*/SKILL.md` with headings/body placeholders that clear the size guard; content research stays deferred.
- [ ] 3.3 Add `internal/assets/generic/persona-operator.md` stub and the additive `PersonaOperator` dispatch case in `internal/components/persona/inject.go`.

## Phase 4: Testing / regression

- [ ] 4.1 Add golden tests for operational asset injection per supported strategy in `internal/components/sddops/inject_test.go`, `internal/components/operationalmcp/inject_test.go`, and relevant `testdata/` goldens.
- [ ] 4.2 Add DEV byte-for-byte regression coverage proving `PresetFullGentleman` assets/config remain unchanged before vs. after planning/installing `selops-operational`.
- [ ] 4.3 Run `gofmt`, `go test ./...`, and `go vet ./...`; capture any fixture/golden refreshes required by the new preset.

## Phase 5: Follow-ups (deferred, not scheduled in this pass)

- [ ] 5.1 Add TUI surfacing for the operational preset after CLI MVP proves stable.
- [ ] 5.2 Design real `OperationalMCPServers` input UX (env/config/prompted source) beyond the placeholder default.
- [ ] 5.3 Revisit sync semantics only if managed refresh of SDD-OPS assets becomes necessary; then consider `InjectForSync`.
- [ ] 5.4 Research and author the 6 operational domains plus the final operator persona body.
