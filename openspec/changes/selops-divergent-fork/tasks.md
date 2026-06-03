# Tasks: SelOps-ai Divergent Fork

Hierarchical, phased. Each phase MUST end with `go build ./... && go test ./...` green.
Strict TDD active (openspec/config.yaml). Work on a feature branch.

---

## Phase 0 — Fork mechanics (start here in the new session)

### 0a. Rename module + binary
- [x] 0a.1 Rename Go module `github.com/gentleman-programming/gentle-ai` → `github.com/Gabrielvilabracho/selops-ai` (go.mod + all imports)
- [x] 0a.2 Rename binary dir `cmd/gentle-ai` → `cmd/selops`
- [x] 0a.3 `go build ./...` green; smoke-test the binary runs

### 0b. Repoint defaults (no deletion yet)
- [ ] 0b.1 `normalizePreset` default → `selops-operational` (cli/validate.go ~line 90)
- [ ] 0b.2 TUI default preset → `selops-operational` (tui/model.go ~line 478)
- [ ] 0b.3 Update `screens/preset.go` copy to describe the OPS stack
- [ ] 0b.4 Fix affected TUI/CLI tests

### 0c. Context7 ALWAYS-ON in OPS preset
- [ ] 0c.1 Add `ComponentContext7` to `componentsForPreset(PresetSelOpsOperational)` (validate.go ~184)
- [ ] 0c.2 Ensure it is non-optional (not behind a flag) in the OPS path
- [ ] 0c.3 Add/adjust golden + regression tests

### 0d. ComponentKnowledgeBase (10 NEUTRAL skills, optional)
- [ ] 0d.1 New SkillID list (the 10 NEUTRAL) in skills/presets.go
- [ ] 0d.2 New `ComponentKnowledgeBase` injector (mirror sddops.Inject pattern)
- [ ] 0d.3 Graph node with nil deps (avoid Skills→SDD edge)
- [ ] 0d.4 Expose via optional preset/flag (NOT in the minimal default)
- [ ] 0d.5 Tests

### 0e. Review + strip DEV content (highest blast radius — LAST)
- [ ] 0e.1 Read the 11 sdd-* + go-testing SKILL.md bodies; confirm nothing reusable before deleting
- [ ] 0e.2 Remove DEV presets (full-gentleman, ecosystem-only, minimal) from consts + componentsForPreset
- [ ] 0e.3 Remove `ComponentSDD` branch from the 6 dispatch files (run.go, sync.go, tui/model.go, upgrade/executor.go, profile_delete.go, uninstall/service.go)
- [ ] 0e.4 Delete `internal/components/sdd` package
- [ ] 0e.5 Delete DEV skill assets (sdd-*, go-testing) + gentleman persona/output-style assets
- [ ] 0e.6 Trim model/types.go, catalog/skills.go, planner/graph.go of DEV IDs
- [ ] 0e.7 Update assets.go embed list + assets_test.go skill counts
- [ ] 0e.8 `go build ./... && go test ./...` green

---

## Phase 1 — CONTENT (the real product)

### 1a. Operator persona
- [ ] 1a.1 Define operator behavior (senior business operator equivalent of el Gentleman)
- [ ] 1a.2 Write `internal/assets/generic/persona-operator.md` (replace 15-line stub)
- [ ] 1a.3 Per-agent variants if needed (claude/opencode)

### 1b. The 6 domain bodies (follow-up 5.4)
For each domain define: what the agent concretely does, real triggers, outputs produced.
- [ ] 1b.1 ops-standard-documentation
- [ ] 1b.2 ops-modular-architecture
- [ ] 1b.3 ops-data-contracts
- [ ] 1b.4 ops-governance
- [ ] 1b.5 ops-observability
- [ ] 1b.6 ops-graduated-autonomy

---

## Phase 2 — Conditional OPS pipeline
- [ ] 2a.1 Define phase agents: ops-brief, ops-structure, ops-produce, ops-review, ops-deliver
- [ ] 2a.2 Write phase agent assets
- [ ] 2b.1 OPS orchestrator AGENTS.md with conditional router (small inline / large → pipeline)
- [ ] 2b.2 Reuse SDD preflight for OPS

---

## Phase 3 — Visual rebrand
- [ ] 3.1 Decide SelOps colors + logo (design decision — user)
- [ ] 3.2 Replace Kanagawa theme + Braille-rose logo assets

---

## Phase 4 — Hermes integration
- [ ] 4.1 SelOps-ai → Hermes via MCP (hermes mcp serve)
- [ ] 4.2 Verify SelOps skills load in Hermes as Markdown
