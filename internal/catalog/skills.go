package catalog

import "github.com/Gabrielvilabracho/selops-ai/internal/model"

type Skill struct {
	ID       model.SkillID
	Name     string
	Category string
	Priority string
}

var mvpSkills = []Skill{
	// Foundation skills (neutral, domain-agnostic)
	{ID: model.SkillGoTesting, Name: "go-testing", Category: "testing", Priority: "p0"},
	{ID: model.SkillCreator, Name: "skill-creator", Category: "workflow", Priority: "p0"},
	{ID: model.SkillImprover, Name: "skill-improver", Category: "workflow", Priority: "p0"},
	{ID: model.SkillBranchPR, Name: "branch-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillIssueCreation, Name: "issue-creation", Category: "workflow", Priority: "p0"},
	{ID: model.SkillSkillRegistry, Name: "skill-registry", Category: "workflow", Priority: "p0"},
	// Sustainable review skills
	{ID: model.SkillChainedPR, Name: "chained-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCognitiveDoc, Name: "cognitive-doc-design", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCommentWriter, Name: "comment-writer", Category: "workflow", Priority: "p0"},
	{ID: model.SkillWorkUnitCommits, Name: "work-unit-commits", Category: "workflow", Priority: "p0"},
	// SelOps operational skills
	{ID: model.SkillOpsStandardDocumentation, Name: "ops-standard-documentation", Category: "operational", Priority: "p0"},
	{ID: model.SkillOpsModularArchitecture, Name: "ops-modular-architecture", Category: "operational", Priority: "p0"},
	{ID: model.SkillOpsDataContracts, Name: "ops-data-contracts", Category: "operational", Priority: "p0"},
	{ID: model.SkillOpsGovernance, Name: "ops-governance", Category: "operational", Priority: "p0"},
	{ID: model.SkillOpsObservability, Name: "ops-observability", Category: "operational", Priority: "p0"},
	{ID: model.SkillOpsGraduatedAutonomy, Name: "ops-graduated-autonomy", Category: "operational", Priority: "p0"},
	// SelOps OPS pipeline phase agents (execution roles, not domain knowledge)
	{ID: model.SkillOpsBrief, Name: "ops-brief", Category: "pipeline", Priority: "p0"},
	{ID: model.SkillOpsStructure, Name: "ops-structure", Category: "pipeline", Priority: "p0"},
	{ID: model.SkillOpsProduce, Name: "ops-produce", Category: "pipeline", Priority: "p0"},
	{ID: model.SkillOpsReview, Name: "ops-review", Category: "pipeline", Priority: "p0"},
	{ID: model.SkillOpsDeliver, Name: "ops-deliver", Category: "pipeline", Priority: "p0"},
}

func MVPSkills() []Skill {
	skills := make([]Skill, len(mvpSkills))
	copy(skills, mvpSkills)
	return skills
}
