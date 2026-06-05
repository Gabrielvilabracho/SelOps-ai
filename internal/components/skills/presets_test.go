package skills

import (
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// TestSkillsForPresetMinimalReturnsFoundationSkills verifies that the minimal preset
// returns foundation skills (OPS fork: sdd skills removed in Phase 0e).
func TestSkillsForPresetMinimalReturnsFoundationSkills(t *testing.T) {
	skills := SkillsForPreset(model.PresetMinimal)
	if len(skills) == 0 {
		t.Fatalf("SkillsForPreset(minimal) returned empty")
	}
	// No sdd-* skills in the OPS fork.
	for _, skill := range skills {
		if strings.HasPrefix(string(skill), "sdd-") {
			t.Fatalf("minimal preset must not contain SDD skills in OPS fork, got %q", skill)
		}
	}
}

func TestSkillsForPresetEcosystemIncludesFrameworks(t *testing.T) {
	skills := SkillsForPreset(model.PresetEcosystemOnly)

	hasGoTesting := false
	hasSkillCreator := false
	for _, skill := range skills {
		if skill == model.SkillGoTesting {
			hasGoTesting = true
		}
		if skill == model.SkillCreator {
			hasSkillCreator = true
		}
	}

	if !hasGoTesting {
		t.Fatalf("ecosystem preset should include go-testing")
	}
	if !hasSkillCreator {
		t.Fatalf("ecosystem preset should include skill-creator")
	}
}

func TestSkillsForPresetFullIncludesAll(t *testing.T) {
	skills := SkillsForPreset(model.PresetFullGentleman)
	all := AllSkillIDs()

	if len(skills) != len(all) {
		t.Fatalf("full preset skills len = %d, all skills len = %d", len(skills), len(all))
	}
}

func TestSkillsForPresetCustomReturnsNil(t *testing.T) {
	skills := SkillsForPreset(model.PresetCustom)
	if skills != nil {
		t.Fatalf("custom preset should return nil, got %v", skills)
	}
}

func TestSkillsForPresetSelOpsOperationalReturnsExactlySixOpsSkills(t *testing.T) {
	skills := SkillsForPreset(model.PresetSelOpsOperational)

	if len(skills) != 8 {
		t.Fatalf("PresetSelOpsOperational should return exactly 8 skills, got %d: %v", len(skills), skills)
	}

	expected := map[model.SkillID]struct{}{
		model.SkillOpsAdversarialSecurity:   {},
		model.SkillOpsPrivacyGovernance:     {},
		model.SkillOpsStandardDocumentation: {},
		model.SkillOpsModularArchitecture:   {},
		model.SkillOpsDataContracts:         {},
		model.SkillOpsGovernance:            {},
		model.SkillOpsObservability:         {},
		model.SkillOpsGraduatedAutonomy:     {},
	}

	for _, skill := range skills {
		if _, ok := expected[skill]; !ok {
			t.Fatalf("unexpected skill %q in PresetSelOpsOperational, expected only ops-* skills", skill)
		}
	}
}

// TestSkillsForPresetSelOpsOperationalContainsNoDevSkills verifies OPS preset
// has no DEV-only skills. OPS fork (Phase 0e): sdd-* IDs removed from model.
func TestSkillsForPresetSelOpsOperationalContainsNoDevSkills(t *testing.T) {
	skills := SkillsForPreset(model.PresetSelOpsOperational)

	for _, skill := range skills {
		if strings.HasPrefix(string(skill), "sdd-") {
			t.Fatalf("PresetSelOpsOperational must NOT contain SDD skill %q", skill)
		}
		// Foundation skills are neutral — they may appear in any preset.
		// Only ops-* skills should be in the OPS preset.
		if !strings.HasPrefix(string(skill), "ops-") {
			t.Fatalf("PresetSelOpsOperational should only contain ops-* skills, got %q", skill)
		}
	}
}

// TestKnowledgeBaseSkillsReturnsExactlyTenNeutralSkills verifies that
// KnowledgeBaseSkills() returns all 10 domain-agnostic foundation skills —
// neither sdd-* skills nor ops-* skills.
func TestKnowledgeBaseSkillsReturnsExactlyTenNeutralSkills(t *testing.T) {
	skills := KnowledgeBaseSkills()

	if len(skills) != 10 {
		t.Fatalf("KnowledgeBaseSkills() returned %d skills, want exactly 10; got %v", len(skills), skills)
	}
}

// TestKnowledgeBaseSkillsContainsNoSDDSkills verifies the knowledge base
// contains no sdd-* skills (those belonged to ComponentSDD, removed in Phase 0e).
func TestKnowledgeBaseSkillsContainsNoSDDSkills(t *testing.T) {
	skills := KnowledgeBaseSkills()

	for _, skill := range skills {
		if strings.HasPrefix(string(skill), "sdd-") {
			t.Errorf("KnowledgeBaseSkills() must not contain SDD skill %q", skill)
		}
	}
}

// TestKnowledgeBaseSkillsContainsNoOpsSkills verifies the knowledge base
// contains no ops-* skills (those belong to ComponentSDDOps, not ComponentKnowledgeBase).
func TestKnowledgeBaseSkillsContainsNoOpsSkills(t *testing.T) {
	skills := KnowledgeBaseSkills()

	for _, skill := range skills {
		if strings.HasPrefix(string(skill), "ops-") {
			t.Errorf("KnowledgeBaseSkills() must not contain OPS skill %q", skill)
		}
	}
}

// TestKnowledgeBaseSkillsContainsExpectedNeutralSkills verifies all 10 known
// neutral skills are present — triangulation to force real logic.
func TestKnowledgeBaseSkillsContainsExpectedNeutralSkills(t *testing.T) {
	got := KnowledgeBaseSkills()

	expected := []model.SkillID{
		model.SkillGoTesting,
		model.SkillCreator,
		model.SkillImprover,
		model.SkillBranchPR,
		model.SkillIssueCreation,
		model.SkillSkillRegistry,
		model.SkillChainedPR,
		model.SkillCognitiveDoc,
		model.SkillCommentWriter,
		model.SkillWorkUnitCommits,
	}

	skillSet := make(map[model.SkillID]struct{}, len(got))
	for _, s := range got {
		skillSet[s] = struct{}{}
	}

	for _, want := range expected {
		if _, ok := skillSet[want]; !ok {
			t.Errorf("KnowledgeBaseSkills() missing expected neutral skill %q; got %v", want, got)
		}
	}
}

// TestAllSkillIDsIncludesFoundationSkills verifies AllSkillIDs() includes
// the foundation skills (sdd-* removed in Phase 0e).
func TestAllSkillIDsIncludesFoundationSkills(t *testing.T) {
	all := AllSkillIDs()

	required := []model.SkillID{
		model.SkillGoTesting,
		model.SkillCreator,
	}

	skillSet := make(map[model.SkillID]struct{}, len(all))
	for _, skill := range all {
		skillSet[skill] = struct{}{}
	}

	for _, req := range required {
		if _, ok := skillSet[req]; !ok {
			t.Fatalf("AllSkillIDs() missing %q", req)
		}
	}
}

// TestAllSkillIDsContainsNoSDDSkills verifies AllSkillIDs() has no sdd-* skills
// after Phase 0e strip.
func TestAllSkillIDsContainsNoSDDSkills(t *testing.T) {
	all := AllSkillIDs()
	for _, id := range all {
		if strings.HasPrefix(string(id), "sdd-") {
			t.Errorf("AllSkillIDs() must NOT contain sdd-* skill %q after Phase 0e strip", id)
		}
	}
}
