// Package sddops installs the SelOps operational SDD skill set.
// It is separate from the core sdd package so the operational layer can
// be added or removed without touching the DEV (Gentleman) preset.
package sddops

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gabrielvilabracho/selops-ai/internal/agents"
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

// Inject writes the operational SDD skill files for all provided adapters.
// It injects both domain knowledge skills (sddOpsSkillIDs) and pipeline phase
// agents (opsPipelineSkillIDs). It calls skills.InjectWithCapability, which
// handles adapter capability checks and per-skill directory creation.
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

	return InjectionResult{Changed: result.Changed, Files: result.Files}, nil
}
