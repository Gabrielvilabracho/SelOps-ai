package cli

import (
	"path/filepath"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)









func TestComponentPathsContext7KimiIncludesMCPConfig(t *testing.T) {
	home := t.TempDir()
	adapters := resolveAdapters([]model.AgentID{model.AgentKimi})

	paths := componentPaths(home, model.Selection{}, adapters, model.ComponentContext7)

	want := filepath.Join(home, ".kimi", "mcp.json")
	if !containsPath(paths, want) {
		t.Fatalf("componentPaths(context7,kimi) missing %q\npaths=%v", want, paths)
	}
}

// TestComponentPathsEngramCodexIncludesConfigTOML verifies that componentPaths
// for ComponentEngram + Codex reports ~/.codex/config.toml as a backup target.
func TestComponentPathsEngramCodexIncludesConfigTOML(t *testing.T) {
	home := t.TempDir()
	adapters := resolveAdapters([]model.AgentID{model.AgentCodex})

	paths := componentPaths(home, model.Selection{}, adapters, model.ComponentEngram)

	want := filepath.Join(home, ".codex", "config.toml")
	if !containsPath(paths, want) {
		t.Fatalf("componentPaths(engram,codex) missing %q\npaths=%v", want, paths)
	}
}

// TestComponentPathsEngramOpenClawUsesCanonicalSettingsPath asserts that the
// engram component path for OpenClaw always resolves to the canonical
// ~/.openclaw/openclaw.json and never to a workspace-scoped copy.
//
// This is a regression test for issue #522: the verifier used to call
// SettingsPath(workspaceDir) which produced
// <workspace>/.openclaw/openclaw.json, causing post-sync verification to
// fail even when the file at the canonical path existed.
func TestComponentPathsEngramOpenClawUsesCanonicalSettingsPath(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	adapters := resolveAdapters([]model.AgentID{model.AgentOpenClaw})

	paths := componentPathsWithWorkspace(home, workspace, model.Selection{}, adapters, model.ComponentEngram)

	canonical := filepath.Join(home, ".openclaw", "openclaw.json")
	if !containsPath(paths, canonical) {
		t.Fatalf("componentPathsWithWorkspace(engram,openclaw) missing canonical path %q\npaths=%v", canonical, paths)
	}

	wrongPath := filepath.Join(workspace, ".openclaw", "openclaw.json")
	if containsPath(paths, wrongPath) {
		t.Fatalf("componentPathsWithWorkspace(engram,openclaw) must not include workspace-scoped path %q\npaths=%v", wrongPath, paths)
	}
}

func containsPath(paths []string, want string) bool {
	for _, p := range paths {
		if p == want {
			return true
		}
	}
	return false
}

// TestComponentApplyStepRunSDDOpsNoError verifies that componentApplyStep.Run
// for ComponentSDDOps succeeds without error when adapters are provided.
// The new component must NOT fall into the default error branch.
func TestComponentApplyStepRunSDDOpsNoError(t *testing.T) {
	home := t.TempDir()
	step := componentApplyStep{
		id:        "component:sddops",
		component: model.ComponentSDDOps,
		homeDir:   home,
		agents:    []model.AgentID{model.AgentOpenCode},
		selection: model.Selection{},
	}
	// Run must not return the "not supported" error.
	err := step.Run()
	if err != nil {
		t.Fatalf("componentApplyStep.Run(ComponentSDDOps) error = %v; want nil", err)
	}
}

// TestComponentApplyStepRunOperationalMCPNoError verifies that componentApplyStep.Run
// for ComponentOperationalMCP succeeds without error.
func TestComponentApplyStepRunOperationalMCPNoError(t *testing.T) {
	home := t.TempDir()
	step := componentApplyStep{
		id:        "component:operationalmcp",
		component: model.ComponentOperationalMCP,
		homeDir:   home,
		agents:    []model.AgentID{model.AgentOpenCode},
		selection: model.Selection{},
	}
	err := step.Run()
	if err != nil {
		t.Fatalf("componentApplyStep.Run(ComponentOperationalMCP) error = %v; want nil", err)
	}
}
