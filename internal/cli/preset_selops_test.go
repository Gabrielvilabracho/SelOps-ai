package cli

// preset_selops_test.go — RED tests for selops-operational preset resolution.
//
// These tests were written FIRST (TDD RED step) to prove the bug identified in
// verify-report #1799: componentsForPreset() had no PresetSelOpsOperational case
// and normalizePersona() rejected "operator". Both caused the shipped preset to
// resolve to the DEV full-gentleman bundle instead of the operational bundle.
//
// Commit: test(cli): RED test for selops-operational preset resolution + operator persona

import (
	"slices"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
	"github.com/Gabrielvilabracho/selops-ai/internal/system"
)

// TestPresetSelOpsOperationalResolvesOperationalBundle verifies that
// --preset selops-operational produces EXACTLY the operational component set:
// ComponentEngram, ComponentSDDOps, ComponentOperationalMCP,
// ComponentPersonaOperator, and ComponentContext7 (always-on, Phase 0c).
//
// ComponentSkills MUST NOT be present: it carries a hard dep on ComponentSDD
// (DEV) via MVPGraph, which would transitively pull the DEV SDD workflow into
// every OPS install. ComponentSDDOps is the carrier for ops-* skills — it does
// not depend on ComponentSDD and calls sddops.Inject directly.
//
// ComponentContext7 has nil deps in MVPGraph (no DEV transitive pull) and is
// always-on for OPS (Phase 0c) — it provides up-to-date library docs to the
// AI agent unconditionally.
//
// This is the test that would have caught CRITICAL #1 during PR1 apply, and
// now additionally catches the DEV/OPS separation regression (OPS bundle must
// not pull DEV SDD via ComponentSkills).
func TestPresetSelOpsOperationalResolvesOperationalBundle(t *testing.T) {
	want := []model.ComponentID{
		model.ComponentEngram,
		model.ComponentSDDOps,
		model.ComponentOperationalMCP,
		model.ComponentPersonaOperator,
		model.ComponentContext7,
	}

	got := componentsForPreset(model.PresetSelOpsOperational, model.PersonaOperator)

	if len(got) != len(want) {
		t.Fatalf("componentsForPreset(PresetSelOpsOperational, PersonaOperator) = %v (len=%d), want %v (len=%d)",
			got, len(got), want, len(want))
	}
	for _, id := range want {
		if !slices.Contains(got, id) {
			t.Errorf("missing component %q in operational preset result %v", id, got)
		}
	}
}

// TestPresetSelOpsOperationalContext7AlwaysOn verifies that Context7 is
// unconditionally included in the OPS preset bundle (Phase 0c).
// Context7 has nil deps in MVPGraph — it is a standalone component that does
// not transitively pull DEV components. It provides up-to-date library/docs to
// the AI agent and must ALWAYS be present, not optional.
func TestPresetSelOpsOperationalContext7AlwaysOn(t *testing.T) {
	got := componentsForPreset(model.PresetSelOpsOperational, model.PersonaOperator)

	if !slices.Contains(got, model.ComponentContext7) {
		t.Errorf("ComponentContext7 MUST be always-on in OPS preset (Phase 0c); got %v", got)
	}
}

// TestPresetSelOpsOperationalExcludesGenericSkillsComponent verifies that
// ComponentSkills is NOT in the OPS bundle.
//
// Root cause context: ComponentSkills has a hard graph dep on ComponentSDD
// (DEV). Including it in the OPS bundle causes "Auto-added dependencies: sdd"
// in the dry-run output, leaking the DEV SDD workflow into OPS installs.
// The ops-* skills are delivered exclusively through ComponentSDDOps/sddops.Inject.
func TestPresetSelOpsOperationalExcludesGenericSkillsComponent(t *testing.T) {
	got := componentsForPreset(model.PresetSelOpsOperational, model.PersonaOperator)

	if slices.Contains(got, model.ComponentSkills) {
		t.Errorf("ComponentSkills MUST NOT be in the OPS bundle (it transitively pulls DEV ComponentSDD); got %v", got)
	}
}

// TestPresetSelOpsOperationalDoesNotContainDEVComponents verifies that the
// operational bundle NEVER includes DEV-only components: ComponentSDD,
// ComponentPermission, ComponentGGA, ComponentClaudeTheme,
// ComponentOpenCodeGentleLogo, ComponentPersona, and ComponentSkills.
//
// ComponentSkills is included here because it is NOT ops-specific and carries
// a hard dep on ComponentSDD via MVPGraph. Its presence in the OPS bundle
// would pull the DEV SDD workflow transitively into every OPS install.
//
// NOTE: ComponentContext7 is intentionally REMOVED from the DEV-only list
// (Phase 0c). Context7 has nil deps (no DEV transitive pull) and is
// always-on for OPS — it provides up-to-date library/docs to the AI agent.
func TestPresetSelOpsOperationalDoesNotContainDEVComponents(t *testing.T) {
	devOnly := []model.ComponentID{
		model.ComponentSDD,
		model.ComponentSkills, // transitively pulls ComponentSDD — must NOT be in OPS
		model.ComponentPermission,
		model.ComponentGGA,
		model.ComponentClaudeTheme,
		model.ComponentOpenCodeGentleLogo,
		model.ComponentPersona, // DEV persona wrapper — must NOT appear in OPS
	}

	got := componentsForPreset(model.PresetSelOpsOperational, model.PersonaOperator)

	for _, id := range devOnly {
		if slices.Contains(got, id) {
			t.Errorf("DEV component %q must NOT appear in operational preset bundle, but found in %v", id, got)
		}
	}
}

// TestPresetSelOpsOperationalExcludesKnowledgeBaseByDefault verifies that
// ComponentKnowledgeBase is NOT automatically included in the OPS bundle.
// It is OPTIONAL and must be explicitly requested by the user.
func TestPresetSelOpsOperationalExcludesKnowledgeBaseByDefault(t *testing.T) {
	got := componentsForPreset(model.PresetSelOpsOperational, model.PersonaOperator)

	if slices.Contains(got, model.ComponentKnowledgeBase) {
		t.Errorf("ComponentKnowledgeBase MUST NOT be auto-included in OPS preset (it is optional); got %v", got)
	}
}

// TestNormalizePersonaAcceptsOperator verifies that normalizePersona("selops-operator")
// (the renamed value) is accepted without error and maps to PersonaOperator.
func TestNormalizePersonaAcceptsOperator(t *testing.T) {
	got, err := normalizePersona(string(model.PersonaOperator))
	if err != nil {
		t.Fatalf("normalizePersona(%q) unexpected error = %v", model.PersonaOperator, err)
	}
	if got != model.PersonaOperator {
		t.Fatalf("normalizePersona(%q) = %q, want %q", model.PersonaOperator, got, model.PersonaOperator)
	}
}

// TestNormalizeInstallFlagsSelOpsOperationalPreset is an end-to-end wire test:
// NormalizeInstallFlags with --preset selops-operational must produce Selection
// containing the operational component set and persona = PersonaOperator.
func TestNormalizeInstallFlagsSelOpsOperationalPreset(t *testing.T) {
	flags := InstallFlags{
		Preset: string(model.PresetSelOpsOperational),
		Agents: []string{"opencode"},
	}

	input, err := NormalizeInstallFlags(flags, system.DetectionResult{})
	if err != nil {
		t.Fatalf("NormalizeInstallFlags(selops-operational) unexpected error = %v", err)
	}

	// Persona must be operator
	if input.Selection.Persona != model.PersonaOperator {
		t.Errorf("Persona = %q, want %q", input.Selection.Persona, model.PersonaOperator)
	}

	// Must contain all operational components (including Context7, always-on Phase 0c).
	wantComponents := []model.ComponentID{
		model.ComponentEngram,
		model.ComponentSDDOps,
		model.ComponentOperationalMCP,
		model.ComponentPersonaOperator,
		model.ComponentContext7,
	}
	for _, id := range wantComponents {
		if !slices.Contains(input.Selection.Components, id) {
			t.Errorf("missing component %q in Selection.Components %v", id, input.Selection.Components)
		}
	}

	// Must NOT contain DEV-only components (including ComponentSkills which
	// transitively pulls ComponentSDD via MVPGraph — not appropriate for OPS).
	// NOTE: ComponentContext7 is NOT in this list — it is always-on for OPS (Phase 0c).
	devOnly := []model.ComponentID{
		model.ComponentSDD,
		model.ComponentSkills, // has hard dep on ComponentSDD — excluded from OPS
		model.ComponentPermission,
		model.ComponentGGA,
		model.ComponentPersona,
	}
	for _, id := range devOnly {
		if slices.Contains(input.Selection.Components, id) {
			t.Errorf("DEV-only component %q must NOT appear in selops-operational selection", id)
		}
	}
}
