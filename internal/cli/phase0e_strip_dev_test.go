package cli

// phase0e_strip_dev_test.go — GREEN verification tests for Phase 0e: Strip DEV sdd package.
//
// These tests verify that after Phase 0e strip:
// 1. catalog.MVPComponents() contains no DEV-only components.
// 2. catalog.MVPSkills() contains no sdd-* skills.
// 3. skills.AllSkillIDs() contains no sdd-* skill IDs.
// 4. normalizePersona rejects PersonaGentleman — DEV persona not valid in OPS CLI.

import (
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/catalog"
	"github.com/Gabrielvilabracho/selops-ai/internal/components/skills"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// TestPhase0EComponentSDDRemovedFromCatalog verifies that "sdd" ComponentID is NOT
// present in catalog.MVPComponents() after Phase 0e strip.
func TestPhase0EComponentSDDRemovedFromCatalog(t *testing.T) {
	components := catalog.MVPComponents()
	for _, c := range components {
		if c.ID == model.ComponentID("sdd") {
			t.Errorf("ComponentSDD (%q) must NOT appear in catalog.MVPComponents() after Phase 0e strip; found in %v", c.ID, components)
		}
	}
}

// TestPhase0ESDDSkillsRemovedFromCatalog verifies that no sdd-* skills appear
// in catalog.MVPSkills() after Phase 0e strip.
func TestPhase0ESDDSkillsRemovedFromCatalog(t *testing.T) {
	mvpSkills := catalog.MVPSkills()
	for _, s := range mvpSkills {
		if strings.HasPrefix(string(s.ID), "sdd-") {
			t.Errorf("sdd-* skill %q must NOT appear in catalog.MVPSkills() after Phase 0e strip", s.ID)
		}
		if string(s.ID) == "judgment-day" {
			t.Errorf("skill %q (was sddSkills member) must NOT appear in catalog.MVPSkills() after Phase 0e strip", s.ID)
		}
	}
}

// TestPhase0EAllSkillIDsContainsNoSDDSkills verifies that skills.AllSkillIDs()
// returns no sdd-* skill IDs after Phase 0e strip.
func TestPhase0EAllSkillIDsContainsNoSDDSkills(t *testing.T) {
	allIDs := skills.AllSkillIDs()
	for _, id := range allIDs {
		if strings.HasPrefix(string(id), "sdd-") {
			t.Errorf("skills.AllSkillIDs() must NOT contain sdd-* skill %q after Phase 0e strip", id)
		}
	}
}

// TestPhase0ENormalizePersonaRejectsGentleman verifies that normalizePersona
// rejects "gentleman" — the DEV persona is not valid in the OPS fork's CLI.
func TestPhase0ENormalizePersonaRejectsGentleman(t *testing.T) {
	_, err := normalizePersona(string(model.PersonaGentleman))
	if err == nil {
		t.Errorf("normalizePersona(%q) should return an error in OPS fork after Phase 0e strip; PersonaGentleman is DEV-only", model.PersonaGentleman)
	}
}
