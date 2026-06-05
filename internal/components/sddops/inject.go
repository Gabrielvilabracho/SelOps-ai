// Package sddops installs the SelOps operational SDD skill set.
// It is separate from the core sdd package so the operational layer can
// be added or removed without touching the DEV (Gentleman) preset.
package sddops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gabrielvilabracho/selops-ai/internal/agents"
	"github.com/Gabrielvilabracho/selops-ai/internal/assets"
	"github.com/Gabrielvilabracho/selops-ai/internal/components/filemerge"
	"github.com/Gabrielvilabracho/selops-ai/internal/components/skills"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// InjectionResult mirrors the shape used by other component packages.
type InjectionResult struct {
	Changed bool
	Files   []string
}

// InjectOptions carries parameters for the injection.
type InjectOptions struct {
	// WorkspaceDir is the root of the current workspace (e.g. os.Getwd()).
	WorkspaceDir string

	// Capability is the model capability ("capable" or "small") used to
	// extract the appropriate section from skill files.
	Capability string

	// ClaudeModelAssignments optionally maps OPS phase names (or "default")
	// to ClaudeModelAlias values. Used to resolve {{CLAUDE_MODEL}} placeholders
	// in Claude sub-agent files. Zero-value safe: nil maps fall back to sonnet.
	ClaudeModelAssignments map[string]model.ClaudeModelAlias

	// KiroModelAssignments optionally maps OPS phase names (or "default")
	// to ClaudeModelAlias values. Used to resolve {{KIRO_MODEL}} placeholders
	// in Kiro sub-agent files. Zero-value safe: nil falls back to ClaudeModelAssignments,
	// then to sonnet.
	KiroModelAssignments map[string]model.ClaudeModelAlias
}

// kiroModelResolver is an optional adapter capability. When implemented,
// injectOpsSubAgents resolves ClaudeModelAlias values to native Kiro model IDs
// and stamps them into the {{KIRO_MODEL}} sentinel in agent frontmatter.
// Adapters that do not implement this interface are unaffected.
type kiroModelResolver interface {
	KiroModelID(alias model.ClaudeModelAlias) string
}

// claudeModelResolver is an optional adapter capability. When implemented,
// injectOpsSubAgents stamps the resolved ClaudeModelAlias into the {{CLAUDE_MODEL}}
// sentinel in agent frontmatter. Claude Code accepts alias strings directly, so
// the resolver is effectively identity over alias.String().
type claudeModelResolver interface {
	ClaudeModelID(alias model.ClaudeModelAlias) string
}

// resolveClaudeModelAlias returns the alias to use for a given OPS phase,
// falling back through assignments["phase"] → assignments["default"] → sonnet.
func resolveClaudeModelAlias(assignments map[string]model.ClaudeModelAlias, phase string) model.ClaudeModelAlias {
	merged := model.ClaudeModelPresetBalanced()
	for key, alias := range assignments {
		if alias.Valid() {
			merged[key] = alias
		}
	}
	if alias, ok := merged[phase]; ok && alias.Valid() {
		return alias
	}
	if alias, ok := merged["default"]; ok && alias.Valid() {
		return alias
	}
	return model.ClaudeModelSonnet
}

// sddOpsSkillIDs is the canonical list of SelOps operational DOMAIN KNOWLEDGE skill IDs.
// These are intentionally distinct from the SDD orchestrator skill IDs
// (which start with "sdd-") to avoid any overlap.
// These skills encode WHAT the operator knows — principles, patterns, checklists.
var sddOpsSkillIDs = []model.SkillID{
	"ops-standard-documentation",
	"ops-modular-architecture",
	"ops-data-contracts",
	"ops-governance",
	"ops-observability",
	"ops-graduated-autonomy",
}

// opsPipelineSkillIDs is the canonical list of SelOps OPS pipeline phase agent skill IDs.
// Kept separate from sddOpsSkillIDs because these are EXECUTION ROLES, not KNOWLEDGE.
// Domain skills (sddOpsSkillIDs) inform WHAT to do; pipeline skills define HOW to do it,
// phase by phase: brief → structure → produce → review → deliver.
// Both lists are injected by this package so the full operational layer is self-contained.
var opsPipelineSkillIDs = []model.SkillID{
	"ops-brief",
	"ops-structure",
	"ops-produce",
	"ops-review",
	"ops-deliver",
}

// allOpsSkillIDs returns the combined list of all ops skill IDs managed by this
// package: domain knowledge skills (sddOpsSkillIDs) and pipeline phase agents
// (opsPipelineSkillIDs). Both are injected together so the full operational
// layer — knowledge AND execution roles — is self-contained in this package.
func allOpsSkillIDs() []model.SkillID {
	combined := make([]model.SkillID, 0, len(sddOpsSkillIDs)+len(opsPipelineSkillIDs))
	combined = append(combined, sddOpsSkillIDs...)
	combined = append(combined, opsPipelineSkillIDs...)
	return combined
}

// opsOrchestratorContent returns the ops-orchestrator markdown for the given adapter.
// OpenCode gets the opencode-specific variant; all other adapters get the generic variant.
// The opencode variant has OpenCode-specific preflight UX (coded options A1/B1/C1/D1,
// localized Spanish block, and model assignment instructions).
func opsOrchestratorContent(agent model.AgentID) string {
	if agent == model.AgentOpenCode || agent == model.AgentKilocode {
		return assets.MustRead("opencode/ops-orchestrator.md")
	}
	return assets.MustRead("generic/ops-orchestrator.md")
}

// readFileOrEmpty reads a file and returns its content as a string.
// Returns "" if the file does not exist.
func readFileOrEmpty(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("read file %q: %w", path, err)
	}
	return string(data), nil
}

// opsPhaseNames is the ordered list of the five OPS pipeline phase names,
// used as the post-check critical phase set for sub-agent injection.
var opsPhaseNames = []string{"ops-brief", "ops-structure", "ops-produce", "ops-review", "ops-deliver"}

// injectOpsSubAgents copies ops-*.md (and ops-*.yaml for Kimi) from
// adapter.EmbeddedSubAgentsDir() into adapter.SubAgentsDir(homeDir),
// resolving {{CLAUDE_MODEL}} / {{KIRO_MODEL}} placeholders for adapters
// that implement claudeModelResolver / kiroModelResolver.
// Post-check: at least one of {ops-brief, ops-deliver} present as .md
// or .yaml with Size() >= 10 bytes.
func injectOpsSubAgents(homeDir string, adapter agents.Adapter, opts InjectOptions) (changed bool, files []string, err error) {
	if !adapter.SupportsSubAgents() {
		return false, nil, nil
	}

	agentsDir := adapter.SubAgentsDir(homeDir)
	if agentsDir == "" {
		return false, nil, nil
	}
	if mkErr := os.MkdirAll(agentsDir, 0o755); mkErr != nil {
		return false, nil, fmt.Errorf("create ops agents dir: %w", mkErr)
	}

	embeddedDir := adapter.EmbeddedSubAgentsDir()
	entries, rdErr := assets.FS.ReadDir(embeddedDir)
	if rdErr != nil {
		return false, nil, fmt.Errorf("read embedded ops agents dir %q: %w", embeddedDir, rdErr)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// Only copy ops-* files. The embedded dir may contain SDD/JD agents too.
		if !strings.HasPrefix(entry.Name(), "ops-") {
			continue
		}
		contentStr := assets.MustRead(embeddedDir + "/" + entry.Name())

		// Derive phase name from filename (strip one extension).
		phase := strings.TrimSuffix(strings.TrimSuffix(entry.Name(), ".yaml"), ".md")

		// Resolve {{KIRO_MODEL}} for adapters that implement kiroModelResolver.
		if kmr, ok := adapter.(kiroModelResolver); ok {
			alias := model.ClaudeModelSonnet
			if opts.KiroModelAssignments != nil {
				if a, ok2 := opts.KiroModelAssignments[phase]; ok2 {
					alias = a
				} else if d, ok2 := opts.KiroModelAssignments["default"]; ok2 {
					alias = d
				}
			} else if opts.ClaudeModelAssignments != nil {
				alias = resolveClaudeModelAlias(opts.ClaudeModelAssignments, phase)
			}
			contentStr = strings.ReplaceAll(contentStr, "{{KIRO_MODEL}}", kmr.KiroModelID(alias))
		}

		// Resolve {{CLAUDE_MODEL}} for adapters that implement claudeModelResolver.
		if cmr, ok := adapter.(claudeModelResolver); ok {
			alias := resolveClaudeModelAlias(opts.ClaudeModelAssignments, phase)
			contentStr = strings.ReplaceAll(contentStr, "{{CLAUDE_MODEL}}", cmr.ClaudeModelID(alias))
		}

		outPath := filepath.Join(agentsDir, entry.Name())
		writeResult, wErr := filemerge.WriteFileAtomic(outPath, []byte(contentStr), 0o644)
		if wErr != nil {
			return false, nil, fmt.Errorf("write ops sub-agent %s: %w", entry.Name(), wErr)
		}
		if writeResult.Changed {
			changed = true
			files = append(files, outPath)
		}
	}

	// Post-check: verify at least one critical phase file exists as .md or .yaml
	// with Size() >= 10 bytes (matches pre-strip post-check pattern).
	for _, phase := range []string{"ops-brief", "ops-deliver"} {
		found := false
		for _, ext := range []string{".md", ".yaml"} {
			checkPath := filepath.Join(agentsDir, phase+ext)
			if info, statErr := os.Stat(checkPath); statErr == nil && info.Size() >= 10 {
				found = true
				break
			}
		}
		if !found {
			return false, nil, fmt.Errorf("post-check: ops sub-agent %q not written correctly (missing or truncated)", phase)
		}
	}

	return changed, files, nil
}

// injectOpsSlashCommands copies every file from assets.OpsCommandsAssetDir(adapter.Agent())
// into adapter.CommandsDir(homeDir). No size post-check (matches pre-strip behavior).
func injectOpsSlashCommands(homeDir string, adapter agents.Adapter) (changed bool, files []string, err error) {
	if !adapter.SupportsSlashCommands() {
		return false, nil, nil
	}

	cmdsDir := adapter.CommandsDir(homeDir)
	if cmdsDir == "" {
		return false, nil, nil
	}
	if mkErr := os.MkdirAll(cmdsDir, 0o755); mkErr != nil {
		return false, nil, fmt.Errorf("create ops commands dir: %w", mkErr)
	}

	embeddedDir := assets.OpsCommandsAssetDir(adapter.Agent())
	entries, rdErr := assets.FS.ReadDir(embeddedDir)
	if rdErr != nil {
		return false, nil, fmt.Errorf("read embedded ops commands dir %q: %w", embeddedDir, rdErr)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		contentStr := assets.MustRead(embeddedDir + "/" + entry.Name())
		outPath := filepath.Join(cmdsDir, entry.Name())
		writeResult, wErr := filemerge.WriteFileAtomic(outPath, []byte(contentStr), 0o644)
		if wErr != nil {
			return false, nil, fmt.Errorf("write ops command %s: %w", entry.Name(), wErr)
		}
		if writeResult.Changed {
			changed = true
			files = append(files, outPath)
		}
	}

	return changed, files, nil
}

// Inject writes the operational SDD skill files for all provided adapters.
// It injects both domain knowledge skills (sddOpsSkillIDs) and pipeline phase
// agents (opsPipelineSkillIDs). It calls skills.InjectWithCapability, which
// handles adapter capability checks and per-skill directory creation.
//
// In addition to skills, Inject writes the ops-orchestrator section into the
// adapter's system prompt file using the InjectMarkdownSection mechanism — the
// same approach used by the engram component for the engram-protocol section.
//
// MVP: install-only. Sync follow-up is a no-op (ComponentSDDOps is excluded
// from managed sync — see TODO comment in internal/cli/sync.go).
func Inject(targetDir string, adapter agents.Adapter, opts InjectOptions) (InjectionResult, error) {
	if !adapter.SupportsSkills() {
		return InjectionResult{}, nil
	}

	skillDir := adapter.SkillsDir(targetDir)
	if skillDir == "" {
		return InjectionResult{}, nil
	}

	capability := opts.Capability
	if capability == "" {
		capability = "capable"
	}

	allIDs := allOpsSkillIDs()
	result, err := skills.InjectWithCapability(targetDir, adapter, allIDs, capability)
	if err != nil {
		return InjectionResult{}, fmt.Errorf("inject sddops skills: %w", err)
	}

	// Post-check: verify each skill SKILL.md exists and is ≥100 bytes.
	// Covers both domain knowledge skills and pipeline phase agents.
	// Mirrors the guard in internal/components/sdd/inject.go:713.
	for _, id := range allIDs {
		path := filepath.Join(skillDir, string(id), "SKILL.md")
		info, statErr := os.Stat(path)
		if statErr != nil {
			return InjectionResult{}, fmt.Errorf("post-check: ops skill %q not found on disk: %w", id, statErr)
		}
		if info.Size() < 100 {
			return InjectionResult{}, fmt.Errorf("post-check: ops skill %q is too small (%d bytes) — content may be empty or corrupt", id, info.Size())
		}
	}

	// Inject the ops-orchestrator section into the adapter's system prompt.
	// This follows the same pattern as engram's engram-protocol injection:
	// read the existing prompt, merge the section in via InjectMarkdownSection,
	// and write atomically. Adapters that do not support a system prompt are skipped.
	if adapter.SupportsSystemPrompt() {
		orchContent := opsOrchestratorContent(adapter.Agent())
		promptPath := adapter.SystemPromptFile(targetDir)
		existing, readErr := readFileOrEmpty(promptPath)
		if readErr != nil {
			return InjectionResult{}, fmt.Errorf("read ops orchestrator prompt: %w", readErr)
		}
		updated := filemerge.InjectMarkdownSection(existing, "ops-orchestrator", orchContent)
		writeResult, writeErr := filemerge.WriteFileAtomic(promptPath, []byte(updated), 0o644)
		if writeErr != nil {
			return InjectionResult{}, fmt.Errorf("write ops orchestrator section: %w", writeErr)
		}
		if writeResult.Changed {
			result.Changed = true
		}
		result.Files = append(result.Files, promptPath)
	}

	// Step 4: Inject OPS sub-agents into the adapter's native agents directory.
	// Gate: adapter.SupportsSubAgents(). homeDir is the target root (user home).
	subAgentChanged, subAgentFiles, subAgentErr := injectOpsSubAgents(targetDir, adapter, opts)
	if subAgentErr != nil {
		return InjectionResult{}, fmt.Errorf("inject ops sub-agents: %w", subAgentErr)
	}
	if subAgentChanged {
		result.Changed = true
	}
	result.Files = append(result.Files, subAgentFiles...)

	// Step 5: Inject OPS slash commands into the adapter's commands directory.
	// Gate: adapter.SupportsSlashCommands().
	cmdChanged, cmdFiles, cmdErr := injectOpsSlashCommands(targetDir, adapter)
	if cmdErr != nil {
		return InjectionResult{}, fmt.Errorf("inject ops slash commands: %w", cmdErr)
	}
	if cmdChanged {
		result.Changed = true
	}
	result.Files = append(result.Files, cmdFiles...)

	return InjectionResult{Changed: result.Changed, Files: result.Files}, nil
}
