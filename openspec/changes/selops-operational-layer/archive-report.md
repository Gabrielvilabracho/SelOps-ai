# Archive Report: selops-operational-layer

**Status**: SHIPPED  
**Archived**: 2026-06-02  
**Mode**: Strict TDD  
**Verify verdict**: PASS-WITH-WARNINGS  

---

## What was built

An additive SelOps operational profile for the gentle-ai fork installer, selectable via `--preset selops-operational`. It runs fully parallel to the existing DEV profile (`--preset full-gentleman`) with zero cross-mutation.

### Delivered components

| Component | ID | Package |
|---|---|---|
| SDD-OPS skill injector | `selops-sddops` | `internal/components/sddops` |
| Operational MCP wiring | `selops-operationalmcp` | `internal/components/operationalmcp` |
| Operator persona | `selops-operator` | `internal/components/persona` (additive dispatch) |

### New IDs (fork-private namespace)
- `PresetSelOpsOperational = "selops-operational"`
- `ComponentSDDOps = "selops-sddops"`
- `ComponentOperationalMCP = "selops-operationalmcp"`
- `ComponentPersonaOperator = "selops-operator"`
- `PersonaOperator = "selops-operator"`
- Six ops domain skills: `ops-*` IDs

---

## PRs merged (in order)

| PR | Branch | Description | Lines |
|----|--------|-------------|-------|
| #1 | `selops/pr0-gofmt` | gofmt sweep on pre-existing upstream drift (whitespace only) | +56/-66 |
| #2 | `selops/pr1-registration` | Operational preset registration + resolution | +890/-10 |
| #3 | `selops/pr2-packages` | sddops + operationalmcp packages + runtime wiring | +735/-3 |
| #4 | `selops/pr3-tests-separation` | Golden + regression coverage and DEV/OPS separation fix | +888/-13 |

All merged to `main` on 2026-06-02.

---

## Test evidence (from verify report)

```
go test ./...     → 53 packages: 50 ok, 3 [no test files], 0 FAIL
go vet ./...      → clean
gofmt -l .        → empty (clean)
```

**Dry-run proof:**
```
--preset selops-operational → Persona: selops-operator
  Components: engram, selops-operationalmcp, selops-operator, selops-sddops
  Auto-added dependencies: none

--preset full-gentleman     → Persona: gentleman
  Components: claude-theme, context7, persona, engram, gga, opencode-gentle-logo, permissions, sdd, skills
  Auto-added dependencies: none
```

**Key regression test**: `TestDEVPresetByteForByteRegression` — installs the real OPS preset into the same home, asserts every DEV file is byte-identical before vs. after. Passed.

---

## Spec compliance

All 10 capability requirements in the spec matrix passed. See verify report (Engram #1799) for the full table.

---

## Known limitations / follow-ups (deferred)

| # | Item |
|---|------|
| 5.1 | TUI surfacing for `selops-operational` preset (CLI-only MVP) |
| 5.2 | Real `OperationalMCPServers` input UX (env/config/prompted source) |
| 5.3 | Sync semantics for SDD-OPS assets (`InjectForSync`) if managed refresh needed |
| 5.4 | Author the 6 operational domain SKILL.md bodies + final operator persona |

---

## Warning (process debt, not correctness)

Strict-TDD provenance is fully auditable for the DEV/OPS separation fix (RED/GREEN/TRIANGULATE/SAFETY-NET table in apply-progress), but PR1/PR2/PR3 entries are summarized at task level. Runtime evidence is strong; this is process debt only.
