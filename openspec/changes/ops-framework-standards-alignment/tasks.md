# Tasks: OPS Framework Standards Alignment

## Review Workload Forecast

| Field | Value |
|---|---|
| Estimated changed lines | ~1,150-1,300 across 4 PRs |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR-1 P0 new + wiring; PR-2 P0/P1 expansions; PR-3 P1/P2 new + wiring; PR-4 P2/P3 expansions + wiring |
| Delivery strategy | auto-chain |
| Chain strategy | stacked-to-main |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

**Task notes**
- PR-1 MUST register only shipped P0 skills in `internal/components/sddops/inject.go:sddOpsSkillIDs` and `internal/components/skills/presets.go:opsSkills`; `Inject()` post-check iterates `allOpsSkillIDs()` and fails orphan `SKILL.md` files.
- Any PR adding `internal/assets/skills/ops-*` directories MUST update `internal/assets/assets_test.go:TestEmbeddedAssetCount`, `internal/components/sddops/inject_test.go`, related goldens, and keep build/test green in the SAME PR.
- Scenario numbering `S1-S23` follows spec order.

## PR-1

### Infrastructure
- [x] 1.1 Create `internal/assets/skills/ops-adversarial-security/SKILL.md` and `internal/assets/skills/ops-privacy-governance/SKILL.md` in 5-section house style + `## References`; satisfy `R1/R2/R11/R12`, `S1-S5/S14-S15`; cite OWASP LLM01-10, MITRE ATLAS, GDPR Art.5/6-7/15-22/25/35, ISO 42001 A.7, NIST AI 600-1.

### Implementation
- [x] 1.2 Add `SkillOpsAdversarialSecurity` and `SkillOpsPrivacyGovernance` in `internal/model/types.go`, plus matching entries in `internal/catalog/skills.go`; satisfy `R3/R4`, `S6/S8`.
- [x] 1.3 Register ONLY the two shipped P0 skills in `internal/components/sddops/inject.go:sddOpsSkillIDs` and `internal/components/skills/presets.go:opsSkills`; preserve 5-phase pipeline and avoid orphan registrations for `R3/R4`, `S6-S8`.

### Testing
- [x] 1.4 Update `internal/assets/assets_test.go:TestEmbeddedAssetCount`, `internal/components/sddops/inject_test.go`, `internal/components/skills/presets_test.go`, and add `internal/components/sddops/testdata/{claude,opencode}-ops-{adversarial-security,privacy-governance}.golden`; DoD: `go build ./... && go test ./... && go vet ./...` + citation/house-style self-check.

## PR-2

### Infrastructure
- [ ] 2.1 Expand `internal/assets/skills/ops-graduated-autonomy/SKILL.md` and `internal/assets/skills/ops-governance/SKILL.md`; satisfy `R1/R2/R5/R6/R19`, `S1-S5/S9-S10/S23`; cite Sheridan/Parasuraman, EU AI Act Art.14.4, SR 11-7, NIST GOVERN 1.6/6.1, SP 800-53 AU-12, EU AI Act Art.12.

### Implementation
- [ ] 2.2 Expand `internal/assets/skills/ops-observability/SKILL.md`; satisfy `R1/R2/R7`, `S1-S5/S11`; cite OWASP LLM01-10, MITRE ATLAS, NIST AI 600-1, ISO 42001, NIST MEASURE 1.1.
- [ ] 2.3 Expand `internal/assets/skills/ops-data-contracts/SKILL.md` and `internal/assets/skills/ops-produce/SKILL.md`; satisfy `R1/R2/R8/R17/R18`, `S1-S5/S12/S22`; cite GDPR Art.5/25 and governance incident-classification linkage.

### Testing
- [ ] 2.4 Refresh `internal/components/sddops/testdata/*ops-{graduated-autonomy,governance,observability,data-contracts,produce}.golden`; DoD: `go build ./... && go test ./... && go vet ./...` + citation/house-style self-check.

## PR-3

### Infrastructure
- [ ] 3.1 Create `internal/assets/skills/ops-model-validation/SKILL.md`, `.../ops-finops-governance/SKILL.md`, and `.../ops-transparency-explainability/SKILL.md`; satisfy `R1/R2/R13/R14/R15`, `S1-S5/S16-S18`; cite NIST MEASURE 2.x, SR 11-7, ISO 42001 A.6.2.4/A.4/A.8.2, NIST AI 600-1, FinOps Foundation, EU AI Act Art.52-53.

### Implementation
- [ ] 3.2 Add the three `SkillID` constants in `internal/model/types.go`, catalog rows in `internal/catalog/skills.go`, and matching registrations in `internal/components/sddops/inject.go:sddOpsSkillIDs` + `internal/components/skills/presets.go:opsSkills`; satisfy `R3/R4`, `S6/S8`.

### Testing
- [ ] 3.3 Update `internal/assets/assets_test.go:TestEmbeddedAssetCount`, `internal/components/sddops/inject_test.go`, `internal/components/skills/presets_test.go`, and add `internal/components/sddops/testdata/{claude,opencode}-ops-{model-validation,finops-governance,transparency-explainability}.golden`; DoD: `go build ./... && go test ./... && go vet ./...` + citation/house-style self-check.

## PR-4

### Infrastructure
- [ ] 4.1 Expand `internal/assets/skills/ops-standard-documentation/SKILL.md` and `internal/assets/skills/ops-modular-architecture/SKILL.md`; satisfy `R1/R2/R9/R10`, `S1-S5/S13`; cite ISO 42001 A.8.2, EU AI Act Art.52, FDA PCCP, GDPR Art.35, OWASP LLM05, MITRE ATLAS, NIST GOVERN 6.1, NIST SP 800-161.
- [ ] 4.2 Create `internal/assets/skills/ops-model-lifecycle/SKILL.md` and expand `internal/assets/skills/ops-structure/SKILL.md`; satisfy `R1/R2/R16/R18`, `S1-S5/S19/S22`; cite NIST GOVERN 1.6/1.7 and ISO 42001 A.6.2.3/A.6.2.8.

### Implementation
- [ ] 4.3 Add `SkillOpsModelLifecycle` in `internal/model/types.go`, catalog entry in `internal/catalog/skills.go`, and registrations in `internal/components/sddops/inject.go:sddOpsSkillIDs` + `internal/components/skills/presets.go:opsSkills`; satisfy `R3/R4`, `S6-S8`.

### Testing
- [ ] 4.4 Update `internal/assets/assets_test.go:TestEmbeddedAssetCount`, `internal/components/sddops/inject_test.go`, `internal/components/skills/presets_test.go`, and add `internal/components/sddops/testdata/{claude,opencode}-ops-model-lifecycle.golden` plus refreshed docs/architecture/structure goldens; DoD: `go build ./... && go test ./... && go vet ./...` + citation/house-style self-check.
