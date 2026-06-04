package assets

import "github.com/Gabrielvilabracho/selops-ai/internal/model"

// SDDCommandsAssetDir returns the embedded slash-command asset directory for an
// agent. Claude uses Claude-native frontmatter under claude/commands; agents
// without a dedicated command set fall back to the OpenCode-compatible assets.
func SDDCommandsAssetDir(agent model.AgentID) string {
	switch agent {
	case model.AgentClaudeCode:
		return "claude/commands"
	default:
		return "opencode/commands"
	}
}

// OpsCommandsAssetDir returns the embedded OPS slash-command asset directory
// for an agent. Claude uses Claude-native frontmatter under claude/ops-commands;
// agents without a dedicated command set fall back to OpenCode-compatible assets.
func OpsCommandsAssetDir(agent model.AgentID) string {
	switch agent {
	case model.AgentClaudeCode:
		return "claude/ops-commands"
	default:
		return "opencode/ops-commands"
	}
}
