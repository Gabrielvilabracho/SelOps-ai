# Proposal: SelOps-ai Divergent Fork

**Status**: PLANNED — ready for Phase 0 in a new session
**Author**: Gabriel + el Gentleman
**Date**: 2026-06-03
**Artifact store**: hybrid (OpenSpec files + Engram)

---

## 1. Intent

Turn this repository (currently the gentle-ai fork with an additive OPS preset)
into **SelOps-ai** — a standalone Go CLI for business/operations work, fully
divergent from gentle-ai (Ruta 2: clean divergent fork). It keeps the powerful
reusable ENGINE, strips DEV-specific CONTENT, and grows its own OPS product.

This is decided. The question was never *if* — only *how*. Ruta 2 chosen for
100% control and independence, accepting that gentle-ai improvements must be
ported by hand.

## 2. Target system architecture (the why behind the fork)

```
VPS
├── Gentle-AI CLI (Go)     → dev, when in terminal
├── SelOps-ai CLI (Go)     → ops/biz, when in terminal   ← THIS PROJECT
└── Docker: Hermes
    ├── Gateway Telegram/WhatsApp → available 24/7
    ├── Cron scheduler            → automations
    ├── MCP Server                → gentle-ai & selops-ai talk to Hermes
    ├── MCPs: filesystem, email, n8n, linear
    └── SelOps skills (Markdown)  → loaded into Hermes
```

Hermes can act as an MCP server (`hermes mcp serve`), so SelOps-ai (via the
target agent) can drive Hermes for 24/7 automation.

## 3. Two orthogonal axes

- **6 DOMAINS (the WHAT — knowledge)**: the living manual of operations.
- **CONDITIONAL PIPELINE (the HOW — delivery)**: only for large/risky
  deliverables, same routing philosophy as SDD-DEV (small inline, big → pipeline).

### The 6 domains
1. `ops-standard-documentation` — docs, READMEs, ADRs, API contracts
2. `ops-modular-architecture`   — module boundaries, structure
3. `ops-data-contracts`         — producer/consumer schemas
4. `ops-governance`             — policies, approvals
5. `ops-observability`          — metrics, traces, logs
6. `ops-graduated-autonomy`     — agent autonomy levels

### The conditional pipeline (OPS)
```
[small task]            → inline (domain answers directly)
[large/risky deliverable] → brief → structure → produce → review → deliver
```

## 4. Component decisions (verified against code)

| Component | Decision | Rationale |
|---|---|---|
| ENGINE (~90-92%) | PRESERVE verbatim | skillregistry, sdd-init engine, engram, 15 adapters, planner, pipeline, filemerge, backup, verify, agentbuilder, state |
| Context7 (components/mcp) | KEEP — ALWAYS ON, not optional | User wants up-to-date docs always active |
| GGA (components/gga) | KEEP | Provider switcher, agnostic infra (NOT security) |
| Theme + Logo | REBRAND not delete | SelOps needs its own visual identity (design pending) |
| 6 ops-* skills | KEEP | The product |
| 10 NEUTRAL skills | KEEP as optional knowledge-base bundle | skill-creator, skill-improver, skill-registry, judgment-day, branch-pr, chained-pr, issue-creation, cognitive-doc-design, comment-writer, work-unit-commits |
| 11 sdd-* + go-testing | REVIEW content, then strip | Tied to DEV SDD orchestrator being removed |
| components/sdd (DEV, ~2.9k lines) | STRIP | The big removal; OPS has ZERO dependency on it |
| gentleman persona | STRIP | Replaced by operator persona |
| DEV presets (full-gentleman, ecosystem-only, minimal) | STRIP | Replaced by selops-operational as default |

## 5. Key finding (de-risks everything)

**OPS has ZERO compile-time dependency on DEV.** `sddops` and `operationalmcp`
import only ENGINE packages (agents, filemerge, model, skills) — never the `sdd`
package. The only `sdd` references in sddops are code comments. Stripping DEV
cannot break OPS. The DEV/OPS separation is already byte-for-byte clean & tested
(`TestDEVPresetByteForByteRegression`).

## 6. Breakage points when removing `components/sdd` (edit in lockstep)

6 engine/dispatch files have a `ComponentSDD` branch to remove:
- `internal/cli/run.go`
- `internal/cli/sync.go`
- `internal/tui/model.go` (also: default preset = PresetFullGentleman, line ~478)
- `internal/update/upgrade/executor.go`
- `internal/tui/screens/profile_delete.go`
- `internal/components/uninstall/service.go`

Shared registries holding both DEV+OPS IDs (mechanical trimming, not breakage):
- `internal/model/types.go` (Component/Skill/Persona/Preset consts)
- `internal/catalog/skills.go` (mvpSkills slice) + components.go
- `internal/planner/graph.go` (MVPGraph)
- `internal/cli/validate.go` (normalizePreset/componentsForPreset)
- `internal/components/skills/presets.go` (sddSkills/foundationSkills/opsSkills)

## 7. The harness work to build (the real product)

1. **Operator persona** — currently a 15-line stub (`persona-operator.md`, DEFERRED)
2. **6 domain bodies** — currently empty stubs (follow-up 5.4)
3. **Conditional OPS pipeline** — phase agents + router in orchestrator AGENTS.md
4. **OPS orchestrator** — the routing contract (inline vs pipeline)

## 8. Rollback plan

Each phase keeps `go build ./... && go test ./...` green. Work on a feature
branch. If a strip step breaks the build, revert that commit — DEV removal is
additive-in-reverse (removing branches), not logic untangling.
