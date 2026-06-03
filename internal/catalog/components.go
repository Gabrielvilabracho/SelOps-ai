package catalog

import "github.com/Gabrielvilabracho/selops-ai/internal/model"

type Component struct {
	ID          model.ComponentID
	Name        string
	Description string
}

var mvpComponents = []Component{
	{ID: model.ComponentEngram, Name: "Engram", Description: "Persistent cross-session memory"},
	{ID: model.ComponentSDD, Name: "SDD", Description: "Spec-driven development workflow"},
	{ID: model.ComponentSkills, Name: "Skills", Description: "Curated coding skill library"},
	{ID: model.ComponentContext7, Name: "Context7", Description: "Latest framework and library docs"},
	{ID: model.ComponentPersona, Name: "Persona", Description: "Gentleman, neutral or custom behavior"},
	{ID: model.ComponentPermission, Name: "Permissions", Description: "Security-first defaults and guardrails"},
	{ID: model.ComponentGGA, Name: "GGA", Description: "Gentleman Guardian Angel — AI provider switcher"},
	{ID: model.ComponentTheme, Name: "Theme", Description: "Gentleman Kanagawa theme overlay"},
	{ID: model.ComponentClaudeTheme, Name: "Claude Gentleman Theme", Description: "Claude Code Gentleman custom theme"},
	{ID: model.ComponentOpenCodeGentleLogo, Name: "OpenCode Gentle Logo", Description: "OpenCode home logo TUI plugin with Braille rose"},
	// SelOps operational layer components.
	{ID: model.ComponentSDDOps, Name: "SDD Ops", Description: "SelOps operational SDD workflow layer"},
	{ID: model.ComponentOperationalMCP, Name: "Operational MCP", Description: "SelOps operational MCP server configuration"},
	{ID: model.ComponentPersonaOperator, Name: "Operator Persona", Description: "SelOps operator persona (replaces Gentleman persona for OPS agents)"},
	// ComponentKnowledgeBase is optional — not included in any preset by default.
	{ID: model.ComponentKnowledgeBase, Name: "Knowledge Base", Description: "10 domain-agnostic foundation skills (opt-in for any operator)"},
}

func MVPComponents() []Component {
	components := make([]Component, len(mvpComponents))
	copy(components, mvpComponents)
	return components
}
