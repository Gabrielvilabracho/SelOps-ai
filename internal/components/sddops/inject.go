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

// sddOpsSkillIDs is the canonical list of SelOps operational skill IDs.
// These are intentionally distinct from the SDD orchestrator skill IDs
// (which start with "sdd-") to avoid any overlap.
var sddOpsSkillIDs = []model.SkillID{
	"ops-standard-documentation",
	"ops-modular-architecture",
	"ops-data-contracts",
	"ops-governance",
	"ops-observability",
	"ops-graduated-autonomy",
}

// Inject writes the operational SDD skill files for all provided adapters.
// It calls skills.InjectWithCapability, which handles adapter capability checks
// and per-skill directory creation.
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

	result, err := skills.InjectWithCapability(targetDir, adapter, sddOpsSkillIDs, capability)
	if err != nil {
		return InjectionResult{}, fmt.Errorf("inject sddops skills: %w", err)
	}

	// Post-check: verify each skill SKILL.md exists and is ≥100 bytes.
	// Mirrors the guard in internal/components/sdd/inject.go:713.
	for _, id := range sddOpsSkillIDs {
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
