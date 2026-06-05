# Proposal: OPS Framework Standards Alignment

## Intent

Align the SelOps OPS framework to cited governance standards it already re-derives and close six missing governance domains. This is needed because the current framework lacks explicit adversarial-security and privacy governance guidance, which are regulatory-blocking gaps for live client production engagements.

## Scope

### In Scope
- Expand 6 existing OPS skills per explored deltas: `ops-graduated-autonomy`, `ops-governance`, `ops-observability`, `ops-data-contracts`, `ops-standard-documentation`, `ops-modular-architecture`.
- Create 6 new domain skills: adversarial security, privacy governance, model validation, FinOps governance, transparency/explainability, and model lifecycle.
- Add coherence touchpoints in `ops-produce`, `ops-structure`, and reversibility guidance in `ops-governance`.
- Wire all 6 new skills in `internal/model`, `internal/catalog`, `internal/components/sddops`, and `internal/components/skills`.

### Out of Scope
- Changing the 5-phase OPS pipeline structure (`brief → deliver`).
- Rewriting `complete-sddops` sub-agent/overlay content.
- Implementing runtime security/privacy enforcement code.
- Certifying SelOps against NIST, ISO, GDPR, EU AI Act, or FDA frameworks.

## Capabilities

### New Capabilities
- `ops-framework-standards`: Defines standards-aligned OPS domain skills, their registry wiring, and required pipeline coherence touchpoints.

### Modified Capabilities
- None.

## Approach

Use the explore-phase deltas as the source of truth: append or adapt existing skill sections, add six new skills in the current house style, and keep Go changes mechanical. Cite real control IDs only.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/assets/skills/` | Modified/New | 12 OPS skill files; primary deliverable |
| `internal/model/types.go` | Modified | 6 new `SkillID` constants |
| `internal/catalog/skills.go` | Modified | 6 skill catalog entries |
| `internal/components/sddops/inject.go` | Modified | `sddOpsSkillIDs` registration |
| `internal/components/skills/presets.go` | Modified | `opsSkills` preset registration |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Fabricated or imprecise standards mapping | Med | Use explored control IDs only; no invented citations |
| Registry/file mismatch breaks install post-check | Med | Ship wiring with referenced skill files in the same PR |
| Review overload across 12 skill files | High | Enforce 4 chained PR slices under the 400-line budget |

## Rollback Plan

Changes are content-only and mechanically wired. Roll back by reverting affected PRs. New skills are removable files; expanded sections revert cleanly; Go registry entries/constants can be removed with no migration or state cleanup. Keep each PR coherent so `inject.go` never references a missing `SKILL.md`.

## Dependencies

- Independent of `complete-sddops` today.
- When `complete-sddops` resumes, update the ops-orchestrator prompt to reference the expanded OPS skill set.

## Success Criteria

- [ ] All 6 new domain skills exist, are wired, and install cleanly.
- [ ] All 6 expanded skills include the agreed standards/control mappings and coherence fixes.
- [ ] Delivery is split into 4 chained PRs: PR-1 P0 skills + all 4 Go files; PR-2 existing-skill expansions + incident fix; PR-3 P1/P2 new skills; PR-4 remaining expansions + lifecycle + data-contract touchpoint.
