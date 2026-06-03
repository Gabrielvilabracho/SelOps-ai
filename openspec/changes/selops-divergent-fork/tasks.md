# Tasks: SelOps-ai Divergent Fork

Hierarchical, phased. Each phase MUST end with `go build ./... && go test ./...` green.
Strict TDD active (openspec/config.yaml). Work on a feature branch.

IMPORTANT: Always run the FULL `go test ./...` (including golden tests for embedded assets)
before merging any PR touching `internal/assets/` — partial test runs miss golden failures.

---

## Phase 0 — Fork mechanics ✅ COMPLETE (PRs #5–#9)

### 0a. Rename module + binary ✅
- [x] 0a.1 Rename Go module `github.com/gentleman-programming/gentle-ai` → `github.com/Gabrielvilabracho/selops-ai`
- [x] 0a.2 Rename binary dir `cmd/gentle-ai` → `cmd/selops`
- [x] 0a.3 `go build ./...` green; smoke-test: `selops 0.0.0-...+dirty`

### 0b. Repoint defaults ✅
- [x] 0b.1 `normalizePreset` default → `selops-operational`
- [x] 0b.2 TUI default preset/persona → `selops-operational` / `selops-operator`
- [x] 0b.3 Fix bug: `componentsForPreset(PresetSelOpsOperational)` was silently falling through to full-gentleman
- [x] 0b.4 Tests updated (3 new TDD tests)

### 0c. Context7 ALWAYS-ON in OPS preset ✅
- [x] 0c.1 `ComponentContext7` added to `componentsForPreset(PresetSelOpsOperational)` in both validate.go and model.go
- [x] 0c.2 Non-optional — not behind a flag in OPS path
- [x] 0c.3 Tests updated

### 0d. ComponentKnowledgeBase (10 NEUTRAL skills, optional) ✅
- [x] 0d.1 `KnowledgeBaseSkills()` wrapping `foundationSkills` in presets.go
- [x] 0d.2 `ComponentKnowledgeBase` constant + catalog + planner + run/sync cases
- [x] 0d.3 Graph node with nil deps
- [x] 0d.4 Optional — user can request via `--component selops-knowledge-base`
- [x] 0d.5 9 new TDD tests

### 0e. Strip DEV ✅
- [x] 0e.1 Deleted `internal/components/sdd/` (~9.1k lines)
- [x] 0e.2 Deleted `internal/assets/skills/sdd-*/` (~2.4k lines)
- [x] 0e.3 Removed ComponentSDD from 6 lockstep breakage points + shared registries
- [x] 0e.4 TUI SDD screens → always false; `normalizePersona` rejects "gentleman"
- [x] 0e.5 `go build ./... && go test ./...` green (49 packages after sdd/ deletion)

---

## Phase 1 — CONTENT ✅ COMPLETE (PR #10)

### 1a. Operator persona ✅
- [x] 1a.1 `internal/assets/generic/persona-operator.md` — 48 lines, dual scope, 10 concrete rules, "You are not an assistant. You are an operator."

### 1b. The 6 domain bodies ✅
- [x] 1b.1 ops-standard-documentation (77 lines)
- [x] 1b.2 ops-modular-architecture (70 lines)
- [x] 1b.3 ops-data-contracts (68 lines)
- [x] 1b.4 ops-governance (79 lines)
- [x] 1b.5 ops-observability (71 lines)
- [x] 1b.6 ops-graduated-autonomy (77 lines)

---

## Phase 2 — Conditional OPS pipeline ✅ COMPLETE (PRs #11, #12)

### 2a. Phase agents ✅ (PR #11)
- [x] 2a.1 5 pipeline phase agent SKILL.md files created
- [x] 2a.2 ops-brief, ops-structure, ops-produce, ops-review, ops-deliver (67–75 lines each)
- [x] 2a.3 Wired as separate `opsPipelineSkillIDs` in sddops/inject.go (knowledge vs execution roles)
- [x] 2a.4 Registries updated: model/types.go, catalog/skills.go, presets.go
- [x] 2a.5 Fixed pre-existing golden test debt from Phase 1 (18→23 asset count)

### 2b. OPS orchestrator + threshold ✅ (PR #12)
- [x] 2b.1 `internal/assets/generic/ops-orchestrator.md` (161 lines, capable + small variants)
- [x] 2b.2 `internal/assets/opencode/ops-orchestrator.md` (177 lines, OpenCode preflight UX)
- [x] 2b.3 OPS threshold: hybrid veto-gate + weighted-score
  - 7 veto gates (incl. AI-native: time-to-detect >24h, eval coverage)
  - 5-dim score: env/reversibility/data-mutation/systems/ttd
  - Route: 0-3 inline / 4-6 pipeline Supervised / 7+ pipeline Suggest
  - Data mutation 4-level scale: read-only / bounded-write / unbounded-write / destructive
  - Grounded in ITIL4, Google SRE, AWS WAR OPS6, NIST AI RMF/GenAI, SAE J3016, FMEA
- [x] 2b.4 Injection via `filemerge.InjectMarkdownSection` in sddops/inject.go
- [x] 2b.5 4 TDD tests (injection, content, adapter selection, idempotency)

---

## Phase 3 — Visual rebrand ✅ COMPLETE
- [x] 3.1 **Design decision (user)**: define SelOps visual identity
  - Colors: Operator palette (midnight blue base, cyan blueprint accent, orange safety, green success)
  - Logo: industrial SELOPS wordmark framed in box-drawing panel badge
  - Tone: industrial/blueprint — dry engineer humor, operational identity
- [x] 3.2 Replace Rose Pine palette in internal/tui/styles/styles.go with Operator palette (13 color vars, same names)
- [x] 3.3 Replace Braille-rose logo in internal/tui/styles/logo.go with SELOPS industrial badge; update gradient to cyan→blue
- [x] 3.4 Update TUI welcome screen: Tagline → SelOps identity; add personality line "On shift. Risk gates armed. Let's keep prod boring."
- [x] 3.4b Update catalog description for ComponentOpenCodeGentleLogo: "Braille rose" → "SelOps Operator badge"
- [x] 3.4c Update plugin.go roseArt/compactArt to SELOPS badge art
- [x] 3.5 `go build ./... && go test ./...` green — 49/49 packages pass, golden tests regenerated

---

## Phase 4 — Hermes integration
- [ ] 4.1 SelOps-ai → Hermes via MCP (`hermes mcp serve`)
- [ ] 4.2 Verify SelOps skills load in Hermes as Markdown
- [ ] 4.3 Test the ops-orchestrator routing works end-to-end via Hermes
