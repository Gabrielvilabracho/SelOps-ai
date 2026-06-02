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

	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

// TestPresetSelOpsOperationalResolvesOperationalBundle verifies that
// --preset selops-operational produces EXACTLY the operational component set:
// ComponentEngram, ComponentSDDOps, ComponentSkills, ComponentOperationalMCP,
// and ComponentPersonaOperator (the operator persona component).
//
// This is the test that would have caught CRITICAL #1 during PR1 apply.
func TestPresetSelOpsOperationalResolvesOperationalBundle(t *testing.T) {
	want := []model.ComponentID{
		model.ComponentEngram,
		model.ComponentSDDOps,
		model.ComponentSkills,
		model.ComponentOperationalMCP,
		model.ComponentPersonaOperator,
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

// TestPresetSelOpsOperationalDoesNotContainDEVComponents verifies that the
// operational bundle NEVER includes DEV-only components: ComponentSDD,
// ComponentContext7, ComponentPermission, ComponentGGA, ComponentClaudeTheme,
// ComponentOpenCodeGentleLogo, ComponentPersona.
func TestPresetSelOpsOperationalDoesNotContainDEVComponents(t *testing.T) {
	devOnly := []model.ComponentID{
		model.ComponentSDD,
		model.ComponentContext7,
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

	// Must contain all operational components
	wantComponents := []model.ComponentID{
		model.ComponentEngram,
		model.ComponentSDDOps,
		model.ComponentSkills,
		model.ComponentOperationalMCP,
		model.ComponentPersonaOperator,
	}
	for _, id := range wantComponents {
		if !slices.Contains(input.Selection.Components, id) {
			t.Errorf("missing component %q in Selection.Components %v", id, input.Selection.Components)
		}
	}

	// Must NOT contain DEV-only components
	devOnly := []model.ComponentID{
		model.ComponentSDD,
		model.ComponentContext7,
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
