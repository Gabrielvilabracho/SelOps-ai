package skills

import "github.com/Gabrielvilabracho/selops-ai/internal/model"

// foundationSkills are baseline learning skills for the "recommended" tier.
// These are neutral, domain-agnostic skills — not sdd-* and not ops-*.
var foundationSkills = []model.SkillID{
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

// opsSkills are the SelOps operational domain knowledge skills — included for PresetSelOpsOperational.
var opsSkills = []model.SkillID{
	model.SkillOpsStandardDocumentation,
	model.SkillOpsModularArchitecture,
	model.SkillOpsDataContracts,
	model.SkillOpsGovernance,
	model.SkillOpsObservability,
	model.SkillOpsGraduatedAutonomy,
}

// opsPipelineSkills are the SelOps OPS pipeline phase agents — execution roles that
// form the 5-phase pipeline: brief → structure → produce → review → deliver.
// Kept separate from opsSkills because they are EXECUTION ROLES, not KNOWLEDGE.
// Domain skills inform what to do; pipeline agents define how to do it, phase by phase.
var opsPipelineSkills = []model.SkillID{
	model.SkillOpsBrief,
	model.SkillOpsStructure,
	model.SkillOpsProduce,
	model.SkillOpsReview,
	model.SkillOpsDeliver,
}

// SkillsForPreset returns which skills should be installed for a given preset.
//
//   - "minimal" / PresetMinimal:              foundation skills only
//   - "ecosystem-only" / PresetEcosystemOnly: foundation skills
//   - "full-gentleman" / PresetFullGentleman: foundation skills
//   - "selops-operational" / PresetSelOpsOperational: ops domain knowledge skills only
//     (pipeline phase agents are injected separately by the sddops component)
//   - "custom" / PresetCustom:                empty (caller should provide explicit list)
func SkillsForPreset(preset model.PresetID) []model.SkillID {
	switch preset {
	case model.PresetMinimal:
		return copySkills(foundationSkills)
	case model.PresetEcosystemOnly:
		return copySkills(foundationSkills)
	case model.PresetFullGentleman:
		return copySkills(foundationSkills)
	case model.PresetSelOpsOperational:
		return copySkills(opsSkills)
	case model.PresetCustom:
		return nil
	default:
		// Unknown preset — default to foundation skills.
		return copySkills(foundationSkills)
	}
}

// KnowledgeBaseSkills returns the 10 domain-agnostic foundation skills used by
// ComponentKnowledgeBase. These are neutral skills — not sdd-* and not ops-*.
// ComponentKnowledgeBase exposes them as an optional add-on for any operator.
func KnowledgeBaseSkills() []model.SkillID {
	return copySkills(foundationSkills)
}

// AllSkillIDs returns every known skill ID (foundation skills only; sdd-* removed in Phase 0e).
func AllSkillIDs() []model.SkillID {
	return copySkills(foundationSkills)
}

func copySkills(src []model.SkillID) []model.SkillID {
	dst := make([]model.SkillID, len(src))
	copy(dst, src)
	return dst
}
