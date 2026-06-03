package skills

import (
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

func TestSkillsForPresetMinimalReturnsSDDOnly(t *testing.T) {
	skills := SkillsForPreset(model.PresetMinimal)
	if len(skills) == 0 {
		t.Fatalf("SkillsForPreset(minimal) returned empty")
	}

	// Orchestration skills that are always bundled with SDD.
	orchestrationSkills := map[model.SkillID]bool{
		model.SkillJudgmentDay: true,
	}

	for _, skill := range skills {
		isSDD := len(skill) >= 4 && skill[:3] == "sdd"
		if !isSDD && !orchestrationSkills[skill] {
			t.Fatalf("minimal preset should only contain SDD/orchestration skills, got %q", skill)
		}
	}
}

func TestSkillsForPresetEcosystemIncludesFrameworks(t *testing.T) {
	skills := SkillsForPreset(model.PresetEcosystemOnly)

	hasGoTesting := false
	hasSkillCreator := false
	hasSDDInit := false
	for _, skill := range skills {
		if skill == model.SkillGoTesting {
			hasGoTesting = true
		}
		if skill == model.SkillCreator {
			hasSkillCreator = true
		}
		if skill == model.SkillSDDInit {
			hasSDDInit = true
		}
	}

	if !hasGoTesting {
		t.Fatalf("ecosystem preset should include go-testing")
	}
	if !hasSDDInit {
		t.Fatalf("ecosystem preset should include sdd-init")
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

	if len(skills) != 6 {
		t.Fatalf("PresetSelOpsOperational should return exactly 6 skills, got %d: %v", len(skills), skills)
	}

	expected := map[model.SkillID]struct{}{
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

func TestSkillsForPresetSelOpsOperationalContainsNoDevSkills(t *testing.T) {
	skills := SkillsForPreset(model.PresetSelOpsOperational)

	devSkills := map[model.SkillID]struct{}{
		model.SkillSDDInit:         {},
		model.SkillSDDApply:        {},
		model.SkillSDDVerify:       {},
		model.SkillSDDExplore:      {},
		model.SkillSDDPropose:      {},
		model.SkillSDDSpec:         {},
		model.SkillSDDDesign:       {},
		model.SkillSDDTasks:        {},
		model.SkillSDDArchive:      {},
		model.SkillSDDOnboard:      {},
		model.SkillGoTesting:       {},
		model.SkillCreator:         {},
		model.SkillImprover:        {},
		model.SkillJudgmentDay:     {},
		model.SkillBranchPR:        {},
		model.SkillIssueCreation:   {},
		model.SkillSkillRegistry:   {},
		model.SkillChainedPR:       {},
		model.SkillCognitiveDoc:    {},
		model.SkillCommentWriter:   {},
		model.SkillWorkUnitCommits: {},
	}

	for _, skill := range skills {
		if _, ok := devSkills[skill]; ok {
			t.Fatalf("PresetSelOpsOperational must NOT contain DEV skill %q", skill)
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
// contains no sdd-* skills (those belong to ComponentSDD, not ComponentKnowledgeBase).
func TestKnowledgeBaseSkillsContainsNoSDDSkills(t *testing.T) {
	skills := KnowledgeBaseSkills()

	for _, skill := range skills {
		if len(skill) >= 4 && skill[:3] == "sdd" {
			t.Errorf("KnowledgeBaseSkills() must not contain SDD skill %q", skill)
		}
	}
}

// TestKnowledgeBaseSkillsContainsNoOpsSkills verifies the knowledge base
// contains no ops-* skills (those belong to ComponentSDDOps, not ComponentKnowledgeBase).
func TestKnowledgeBaseSkillsContainsNoOpsSkills(t *testing.T) {
	skills := KnowledgeBaseSkills()

	for _, skill := range skills {
		if len(skill) >= 4 && skill[:4] == "ops-" {
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

func TestAllSkillIDsIncludesEveryKnownSkill(t *testing.T) {
	all := AllSkillIDs()

	required := []model.SkillID{
		model.SkillSDDInit,
		model.SkillGoTesting,
		model.SkillCreator,
		model.SkillJudgmentDay,
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
