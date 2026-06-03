package screens

import (
	"strings"

	"github.com/Gabrielvilabracho/selops-ai/internal/components/skills"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
	"github.com/Gabrielvilabracho/selops-ai/internal/tui/styles"
)

// skillLabels maps each SkillID to a human-readable display label.
// OPS fork (Phase 0e): sdd-* skill IDs removed.
var skillLabels = map[model.SkillID]string{
	model.SkillGoTesting:     "Go Testing",
	model.SkillCreator:       "Skill Creator",
	model.SkillImprover:      "Skill Improver",
	model.SkillBranchPR:      "Branch & PR",
	model.SkillIssueCreation: "Issue Creation",
	model.SkillSkillRegistry: "Skill Registry",
	model.SkillChainedPR:     "Chained PR",
	model.SkillCognitiveDoc:  "Cognitive Doc",
	model.SkillCommentWriter: "Comment Writer",
	model.SkillWorkUnitCommits: "Work Unit Commits",
}

// SkillPickerOptions returns the action buttons shown after the skill checkboxes.
func SkillPickerOptions() []string {
	return []string{"Continue", "Back"}
}

// AllSkillsOrdered returns all skills in display order.
// OPS fork (Phase 0e): sdd-* skills removed; foundation skills only.
func AllSkillsOrdered() []model.SkillID {
	return skills.AllSkillIDs()
}

// SkillPickerOptionCount returns the total number of navigable rows on the skill picker screen.
func SkillPickerOptionCount() int {
	return len(AllSkillsOrdered()) + len(SkillPickerOptions())
}

// RenderSkillPicker renders the skill selection screen for custom preset mode.
func RenderSkillPicker(selectedSkills []model.SkillID, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("Select Skills"))
	b.WriteString("\n\n")
	b.WriteString(styles.SubtextStyle.Render("Toggle skills with enter or space. All are pre-selected by default."))
	b.WriteString("\n\n")

	selectedSet := make(map[model.SkillID]struct{}, len(selectedSkills))
	for _, s := range selectedSkills {
		selectedSet[s] = struct{}{}
	}

	allSkills := AllSkillsOrdered()

	// ── Foundation Skills group ───────────────────────────────────────────────
	b.WriteString(styles.HeadingStyle.Render("Foundation Skills"))
	b.WriteString("\n")

	for idx, skillID := range allSkills {
		_, checked := selectedSet[skillID]
		focused := idx == cursor
		label := skillLabelFor(skillID)
		b.WriteString(renderCheckbox(label, checked, focused))
	}

	b.WriteString("\n")

	// ── Action buttons ────────────────────────────────────────────────────────
	actionOffset := cursor - len(allSkills)
	b.WriteString(renderOptions(SkillPickerOptions(), actionOffset))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • space/enter: toggle • esc: back"))

	return b.String()
}

// skillLabelFor returns the human-readable label for a skill ID.
func skillLabelFor(id model.SkillID) string {
	if label, ok := skillLabels[id]; ok {
		return label
	}
	return string(id)
}
