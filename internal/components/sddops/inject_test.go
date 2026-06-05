package sddops

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/agents"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/claude"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/cursor"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/kimi"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/kilocode"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/kiro"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/opencode"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

var updateGoldens = flag.Bool("update", false, "update sddops golden files")

// goldenAdapter returns a named adapter for golden test parameterisation.
type goldenAdapter struct {
	name    string
	adapter agents.Adapter
}

func goldenAdapters() []goldenAdapter {
	return []goldenAdapter{
		{name: "opencode", adapter: opencode.NewAdapter()},
		{name: "claude", adapter: claude.NewAdapter()},
	}
}

// assertGolden reads (or writes, when -update is set) the golden fixture at
// testdata/<name>.golden and compares it to actual.
func assertGolden(t *testing.T, name string, actual string) {
	t.Helper()
	goldenPath := filepath.Join("testdata", name+".golden")
	if *updateGoldens {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(goldenPath), err)
		}
		if err := os.WriteFile(goldenPath, []byte(actual), 0o644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", goldenPath, err)
		}
		return
	}
	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v — run with -args -update to generate", goldenPath, err)
	}
	if string(expected) != actual {
		t.Fatalf("golden mismatch for %s\n\nwant:\n%s\n\ngot:\n%s", name, string(expected), actual)
	}
}

// TestSddOpsSkillIDsMembership verifies the exact set of sddOpsSkillIDs (domain knowledge skills).
func TestSddOpsSkillIDsMembership(t *testing.T) {
	want := []model.SkillID{
		"ops-standard-documentation",
		"ops-modular-architecture",
		"ops-data-contracts",
		"ops-governance",
		"ops-observability",
		"ops-graduated-autonomy",
	}
	if len(sddOpsSkillIDs) != len(want) {
		t.Fatalf("sddOpsSkillIDs len = %d, want %d", len(sddOpsSkillIDs), len(want))
	}
	wantSet := make(map[model.SkillID]struct{}, len(want))
	for _, id := range want {
		wantSet[id] = struct{}{}
	}
	for _, id := range sddOpsSkillIDs {
		if _, ok := wantSet[id]; !ok {
			t.Errorf("sddOpsSkillIDs contains unexpected ID %q", id)
		}
	}
}

// TestOpsPipelineSkillIDsMembership verifies the exact set of opsPipelineSkillIDs
// (execution-role agents that form the 5-phase operational pipeline).
func TestOpsPipelineSkillIDsMembership(t *testing.T) {
	want := []model.SkillID{
		"ops-brief",
		"ops-structure",
		"ops-produce",
		"ops-review",
		"ops-deliver",
	}
	if len(opsPipelineSkillIDs) != len(want) {
		t.Fatalf("opsPipelineSkillIDs len = %d, want %d", len(opsPipelineSkillIDs), len(want))
	}
	wantSet := make(map[model.SkillID]struct{}, len(want))
	for _, id := range want {
		wantSet[id] = struct{}{}
	}
	for _, id := range opsPipelineSkillIDs {
		if _, ok := wantSet[id]; !ok {
			t.Errorf("opsPipelineSkillIDs contains unexpected ID %q", id)
		}
	}
}

// TestAllOpsSkillIDsContainsBothLists verifies that allOpsSkillIDs() returns the
// union of sddOpsSkillIDs and opsPipelineSkillIDs with no overlap.
func TestAllOpsSkillIDsContainsBothLists(t *testing.T) {
	all := allOpsSkillIDs()
	expectedLen := len(sddOpsSkillIDs) + len(opsPipelineSkillIDs)
	if len(all) != expectedLen {
		t.Fatalf("allOpsSkillIDs() len = %d, want %d (sddOps=%d + pipeline=%d)",
			len(all), expectedLen, len(sddOpsSkillIDs), len(opsPipelineSkillIDs))
	}
	// No duplicates.
	seen := make(map[model.SkillID]struct{}, len(all))
	for _, id := range all {
		if _, ok := seen[id]; ok {
			t.Errorf("allOpsSkillIDs() contains duplicate ID %q", id)
		}
		seen[id] = struct{}{}
	}
	// All domain knowledge skills present.
	for _, id := range sddOpsSkillIDs {
		if _, ok := seen[id]; !ok {
			t.Errorf("allOpsSkillIDs() missing domain knowledge skill %q", id)
		}
	}
	// All pipeline phase skills present.
	for _, id := range opsPipelineSkillIDs {
		if _, ok := seen[id]; !ok {
			t.Errorf("allOpsSkillIDs() missing pipeline skill %q", id)
		}
	}
}

// TestSddOpsSkillIDsNoIntersectionWithSddSkills verifies zero intersection between
// the sddops skill IDs and the SDD orchestrator skill IDs (those starting with "sdd-").
func TestSddOpsSkillIDsNoIntersectionWithSddSkills(t *testing.T) {
	for _, id := range allOpsSkillIDs() {
		s := string(id)
		if len(s) >= 4 && s[:4] == "sdd-" {
			t.Errorf("allOpsSkillIDs() contains SDD skill ID %q — must have zero intersection", id)
		}
	}
}

// TestInjectWritesSkillFilesWithMinimumSize verifies that after a successful
// Inject call each ops SKILL.md (domain knowledge + pipeline phase agents) exists
// on disk and is ≥100 bytes. Mirrors the guard in internal/components/sdd/inject.go:713.
func TestInjectWritesSkillFilesWithMinimumSize(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	result, err := Inject(home, adapter, InjectOptions{})
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	skillDir := adapter.SkillsDir(home)
	for _, id := range allOpsSkillIDs() {
		path := filepath.Join(skillDir, string(id), "SKILL.md")
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Errorf("skill %q SKILL.md not found: %v", id, statErr)
			continue
		}
		if info.Size() < 100 {
			t.Errorf("skill %q SKILL.md is %d bytes; want ≥100", id, info.Size())
		}
	}

	if !result.Changed {
		t.Error("Inject() first changed = false; want true (new files)")
	}

	// Idempotency: second call must not change files.
	result2, err2 := Inject(home, adapter, InjectOptions{})
	if err2 != nil {
		t.Fatalf("Inject() second error = %v", err2)
	}
	if result2.Changed {
		t.Error("Inject() second changed = true; want false (idempotent)")
	}
}

// TestInjectWritesExpectedOpsHeadings verifies that each SKILL.md
// (domain knowledge + pipeline phase agents) contains at least one heading.
func TestInjectWritesExpectedOpsHeadings(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	if _, err := Inject(home, adapter, InjectOptions{}); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	skillDir := adapter.SkillsDir(home)
	for _, id := range allOpsSkillIDs() {
		path := filepath.Join(skillDir, string(id), "SKILL.md")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("skill %q ReadFile error = %v", id, err)
			continue
		}
		text := string(content)
		// Each skill must have a top-level heading.
		if !strings.Contains(text, "#") {
			t.Errorf("skill %q SKILL.md missing any heading", id)
		}
	}
}

// TestInjectGoldenPerAdapter covers the injected SKILL.md content per adapter
// strategy. One golden file is produced per (adapter, skill) pair, capturing
// the exact bytes written to disk. This serves as a byte-for-byte regression
// guard: any unintended change to an ops skill file will cause a golden mismatch.
// Covers both domain knowledge skills (sddOpsSkillIDs) and pipeline phase agents
// (opsPipelineSkillIDs) via allOpsSkillIDs().
func TestInjectGoldenPerAdapter(t *testing.T) {
	for _, ga := range goldenAdapters() {
		ga := ga
		t.Run(ga.name, func(t *testing.T) {
			if !ga.adapter.SupportsSkills() {
				t.Skipf("%s does not support skills", ga.name)
			}

			home := t.TempDir()
			result, err := Inject(home, ga.adapter, InjectOptions{})
			if err != nil {
				t.Fatalf("Inject(%s) error = %v", ga.name, err)
			}
			if !result.Changed {
				t.Fatalf("Inject(%s) changed = false; want true", ga.name)
			}

			skillDir := ga.adapter.SkillsDir(home)
			for _, id := range allOpsSkillIDs() {
				id := id
				t.Run(string(id), func(t *testing.T) {
					path := filepath.Join(skillDir, string(id), "SKILL.md")
					content, readErr := os.ReadFile(path)
					if readErr != nil {
						t.Fatalf("ReadFile(%q) error = %v", path, readErr)
					}
					goldenName := ga.name + "-" + string(id)
					assertGolden(t, goldenName, string(content))
				})
			}
		})
	}
}

// TestInjectGoldenResultFiles verifies that result.Files lists every SKILL.md
// written for a fresh inject (no pre-existing files). Covers both domain knowledge
// skills and pipeline phase agents via allOpsSkillIDs().
func TestInjectGoldenResultFiles(t *testing.T) {
	for _, ga := range goldenAdapters() {
		ga := ga
		t.Run(ga.name, func(t *testing.T) {
			if !ga.adapter.SupportsSkills() {
				t.Skipf("%s does not support skills", ga.name)
			}

			home := t.TempDir()
			result, err := Inject(home, ga.adapter, InjectOptions{})
			if err != nil {
				t.Fatalf("Inject(%s) error = %v", ga.name, err)
			}

			// Every skill ID (domain + pipeline) must appear in result.Files.
			skillDir := ga.adapter.SkillsDir(home)
			for _, id := range allOpsSkillIDs() {
				expectedPath := filepath.Join(skillDir, string(id), "SKILL.md")
				found := false
				for _, f := range result.Files {
					if f == expectedPath {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Inject(%s) result.Files missing %q; got %v", ga.name, expectedPath, result.Files)
				}
			}
		})
	}
}

// TestInjectNonSkillsAdapterIsNoop verifies that adapters without skills support
// return an empty result without error.
func TestInjectNonSkillsAdapterIsNoop(t *testing.T) {
	// Pi adapter explicitly returns false for SupportsSkills.
	piAdapter, err := agents.NewAdapter(model.AgentPi)
	if err != nil {
		t.Fatalf("NewAdapter(pi) error = %v", err)
	}
	home := t.TempDir()
	result, injectErr := Inject(home, piAdapter, InjectOptions{})
	if injectErr != nil {
		t.Fatalf("Inject(pi) error = %v; want nil", injectErr)
	}
	if result.Changed {
		t.Fatal("Inject(pi) changed = true; want false for non-skills adapter")
	}
	if len(result.Files) != 0 {
		t.Fatalf("Inject(pi) files = %v; want empty", result.Files)
	}
}

// TestInjectWritesOpsOrchestratorSection verifies that Inject writes the
// <!-- gentle-ai:ops-orchestrator --> section into the adapter's system prompt.
// This proves the orchestrator is injected for the selops-operational preset.
func TestInjectWritesOpsOrchestratorSection(t *testing.T) {
	for _, ga := range goldenAdapters() {
		ga := ga
		t.Run(ga.name, func(t *testing.T) {
			if !ga.adapter.SupportsSystemPrompt() {
				t.Skipf("%s does not support system prompt", ga.name)
			}

			home := t.TempDir()
			if _, err := Inject(home, ga.adapter, InjectOptions{}); err != nil {
				t.Fatalf("Inject(%s) error = %v", ga.name, err)
			}

			promptPath := ga.adapter.SystemPromptFile(home)
			content, err := os.ReadFile(promptPath)
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", promptPath, err)
			}
			text := string(content)

			// Must contain the ops-orchestrator section markers.
			if !strings.Contains(text, "<!-- gentle-ai:ops-orchestrator -->") {
				t.Errorf("%s: system prompt missing ops-orchestrator open marker", ga.name)
			}
			if !strings.Contains(text, "<!-- /gentle-ai:ops-orchestrator -->") {
				t.Errorf("%s: system prompt missing ops-orchestrator close marker", ga.name)
			}
		})
	}
}

// TestOpsOrchestratorContainsThresholdContent verifies that the injected
// ops-orchestrator section contains the required threshold model content:
// veto gates, weighted score table, and route-by-score mapping.
func TestOpsOrchestratorContainsThresholdContent(t *testing.T) {
	for _, ga := range goldenAdapters() {
		ga := ga
		t.Run(ga.name, func(t *testing.T) {
			if !ga.adapter.SupportsSystemPrompt() {
				t.Skipf("%s does not support system prompt", ga.name)
			}

			home := t.TempDir()
			if _, err := Inject(home, ga.adapter, InjectOptions{}); err != nil {
				t.Fatalf("Inject(%s) error = %v", ga.name, err)
			}

			promptPath := ga.adapter.SystemPromptFile(home)
			content, err := os.ReadFile(promptPath)
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", promptPath, err)
			}
			text := string(content)

			// Required threshold model content — must be present in EVERY adapter's injected orchestrator.
			for _, required := range []string{
				// Veto gates section
				"Veto Gates",
				"destructive",
				// Weighted score table
				"Weighted Score",
				"Reversibility",
				"Data mutation",
				"Systems affected",
				// Route-by-score mapping
				"Route by Score",
				"Inline",
				"Pipeline in Supervised mode",
				"Pipeline in Suggest mode",
				// OPS Pipeline phases
				"ops-brief",
				"ops-structure",
				"ops-produce",
				"ops-review",
				"ops-deliver",
				// Preflight gate
				"OPS Session Preflight",
				// Init guard
				"OPS Init Guard",
				// Result contract
				"Result Contract",
			} {
				if !strings.Contains(text, required) {
					t.Errorf("%s: ops-orchestrator missing required threshold content %q", ga.name, required)
				}
			}
		})
	}
}

// TestOpsOrchestratorContentSelectionByAdapter verifies that:
// - OpenCode gets the opencode-specific variant (has OpenCode model assignments section).
// - Claude gets the generic variant (has the section:model-capable markup).
func TestOpsOrchestratorContentSelectionByAdapter(t *testing.T) {
	tests := []struct {
		name     string
		adapter  agents.Adapter
		required string
		absent   string
	}{
		{
			name:     "opencode",
			adapter:  opencode.NewAdapter(),
			// OpenCode variant has model assignments section and coded preflight options.
			required: "Model Assignments",
			// Generic variant has section:model-capable markup; opencode does not.
			absent: "<!-- section:model-capable -->",
		},
		{
			name:     "claude",
			adapter:  claude.NewAdapter(),
			// Generic variant has section:model-capable markup.
			required: "<!-- section:model-capable -->",
			// OpenCode model assignments are in the opencode variant only.
			absent: "Model Assignments",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			if _, err := Inject(home, tc.adapter, InjectOptions{}); err != nil {
				t.Fatalf("Inject(%s) error = %v", tc.name, err)
			}

			promptPath := tc.adapter.SystemPromptFile(home)
			content, err := os.ReadFile(promptPath)
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", promptPath, err)
			}
			text := string(content)

			if !strings.Contains(text, tc.required) {
				t.Errorf("%s: expected %q in injected orchestrator content", tc.name, tc.required)
			}
			if tc.absent != "" && strings.Contains(text, tc.absent) {
				t.Errorf("%s: did not expect %q in injected orchestrator content", tc.name, tc.absent)
			}
		})
	}
}

// TestOpsOrchestratorSectionIsIdempotent verifies that a second Inject call
// does not change the system prompt when the ops-orchestrator section is already present.
func TestOpsOrchestratorSectionIsIdempotent(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	// First inject.
	result1, err := Inject(home, adapter, InjectOptions{})
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !result1.Changed {
		t.Fatal("Inject() first changed = false; want true (new install)")
	}

	// Read the system prompt after first inject.
	promptPath := adapter.SystemPromptFile(home)
	content1, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile after first inject error = %v", err)
	}

	// Second inject.
	result2, err := Inject(home, adapter, InjectOptions{})
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}

	// Read the system prompt after second inject.
	content2, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile after second inject error = %v", err)
	}

	// Content must be identical across both injects.
	if string(content1) != string(content2) {
		t.Error("Inject() second call changed the system prompt; want idempotent result")
	}
	_ = result2 // Changed may be false or true for individual files; content idempotency is the invariant.
}

// opsPhases reuses the production opsPhaseNames slice (same package) to keep
// the two in sync automatically — no separate copy maintained in tests.
var opsPhases = opsPhaseNames

// ── injectOpsSubAgents tests (scenarios 1, 5, 6, 7, 8, 9) ─────────────────────

// TestInjectOpsSubAgents_CapabilityGate — Scenario 1
// Only capable adapters receive sub-agent files; OpenCode (SupportsSubAgents=false) does not.
func TestInjectOpsSubAgents_CapabilityGate(t *testing.T) {
	tests := []struct {
		name          string
		adapter       agents.Adapter
		expectFiles   bool
	}{
		{name: "claude", adapter: claude.NewAdapter(), expectFiles: true},
		{name: "opencode", adapter: opencode.NewAdapter(), expectFiles: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			changed, files, err := injectOpsSubAgents(home, tc.adapter, InjectOptions{})
			if err != nil {
				t.Fatalf("injectOpsSubAgents(%s) error = %v", tc.name, err)
			}
			if tc.expectFiles {
				if !changed {
					t.Errorf("injectOpsSubAgents(%s) changed = false; want true (new files)", tc.name)
				}
				if len(files) == 0 {
					t.Errorf("injectOpsSubAgents(%s) files is empty; want sub-agent files", tc.name)
				}
			} else {
				if changed {
					t.Errorf("injectOpsSubAgents(%s) changed = true; want false (no sub-agent support)", tc.name)
				}
				if len(files) != 0 {
					t.Errorf("injectOpsSubAgents(%s) files = %v; want empty", tc.name, files)
				}
			}
		})
	}
}

// TestInjectOpsSubAgents_ClaudeFiveFiles — Scenario 5
// Claude receives exactly five OPS sub-agent .md files, each >= 10 bytes.
func TestInjectOpsSubAgents_ClaudeFiveFiles(t *testing.T) {
	home := t.TempDir()
	adapter := claude.NewAdapter()

	changed, files, err := injectOpsSubAgents(home, adapter, InjectOptions{})
	if err != nil {
		t.Fatalf("injectOpsSubAgents(claude) error = %v", err)
	}
	if !changed {
		t.Fatal("injectOpsSubAgents(claude) changed = false; want true")
	}

	agentsDir := adapter.SubAgentsDir(home)
	for _, phase := range opsPhases {
		path := filepath.Join(agentsDir, phase+".md")
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Errorf("claude sub-agent %q not found: %v", phase, statErr)
			continue
		}
		if info.Size() < 10 {
			t.Errorf("claude sub-agent %q is %d bytes; want >= 10", phase, info.Size())
		}
	}
	// Exactly 5 files written.
	if len(files) != len(opsPhases) {
		t.Errorf("injectOpsSubAgents(claude) wrote %d files; want %d", len(files), len(opsPhases))
	}
}

// TestInjectOpsSubAgents_PlaceholderResolved — Scenario 6
// Written files must NOT contain {{CLAUDE_MODEL}} or {{KIRO_MODEL}} literals.
func TestInjectOpsSubAgents_PlaceholderResolved(t *testing.T) {
	tests := []struct {
		name        string
		adapter     agents.Adapter
		placeholder string
	}{
		{name: "claude", adapter: claude.NewAdapter(), placeholder: "{{CLAUDE_MODEL}}"},
		{name: "kiro", adapter: kiro.NewAdapter(), placeholder: "{{KIRO_MODEL}}"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			_, _, err := injectOpsSubAgents(home, tc.adapter, InjectOptions{})
			if err != nil {
				t.Fatalf("injectOpsSubAgents(%s) error = %v", tc.name, err)
			}

			agentsDir := tc.adapter.SubAgentsDir(home)
			entries, readErr := os.ReadDir(agentsDir)
			if readErr != nil {
				t.Fatalf("ReadDir(%q) error = %v", agentsDir, readErr)
			}
			for _, entry := range entries {
				content, rErr := os.ReadFile(filepath.Join(agentsDir, entry.Name()))
				if rErr != nil {
					t.Fatalf("ReadFile(%q) error = %v", entry.Name(), rErr)
				}
				if strings.Contains(string(content), tc.placeholder) {
					t.Errorf("%s sub-agent %q still contains %q; placeholder not resolved", tc.name, entry.Name(), tc.placeholder)
				}
			}
		})
	}
}

// TestInjectOpsSubAgents_KimiDualFormat — Scenario 7
// Kimi gets both .md and .yaml files for every OPS phase.
func TestInjectOpsSubAgents_KimiDualFormat(t *testing.T) {
	home := t.TempDir()
	kimiAdapter := kimi.NewAdapter()

	_, _, err := injectOpsSubAgents(home, kimiAdapter, InjectOptions{})
	if err != nil {
		t.Fatalf("injectOpsSubAgents(kimi) error = %v", err)
	}

	agentsDir := kimiAdapter.SubAgentsDir(home)
	for _, phase := range opsPhases {
		for _, ext := range []string{".md", ".yaml"} {
			path := filepath.Join(agentsDir, phase+ext)
			if _, statErr := os.Stat(path); statErr != nil {
				t.Errorf("kimi sub-agent %q not found: %v", phase+ext, statErr)
			}
		}
	}
}

// TestInjectOpsSubAgents_OpenCodeNoFiles — Scenario 8
// OpenCode and Kilocode have SupportsSubAgents=false → no native sub-agent files.
func TestInjectOpsSubAgents_OpenCodeNoFiles(t *testing.T) {
	tests := []struct {
		name    string
		adapter agents.Adapter
	}{
		{name: "opencode", adapter: opencode.NewAdapter()},
		{name: "kilocode", adapter: kilocode.NewAdapter()},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			changed, files, err := injectOpsSubAgents(home, tc.adapter, InjectOptions{})
			if err != nil {
				t.Fatalf("injectOpsSubAgents(%s) error = %v", tc.name, err)
			}
			if changed {
				t.Errorf("injectOpsSubAgents(%s) changed = true; want false", tc.name)
			}
			if len(files) != 0 {
				t.Errorf("injectOpsSubAgents(%s) files = %v; want empty", tc.name, files)
			}
		})
	}
}

// TestInjectOpsSubAgents_Idempotent — Scenario 9
// Running injectOpsSubAgents a second time reports changed=false and files unchanged.
func TestInjectOpsSubAgents_Idempotent(t *testing.T) {
	home := t.TempDir()
	adapter := claude.NewAdapter()

	if _, _, err := injectOpsSubAgents(home, adapter, InjectOptions{}); err != nil {
		t.Fatalf("injectOpsSubAgents first run error = %v", err)
	}

	// Capture content after first run.
	agentsDir := adapter.SubAgentsDir(home)
	contentsBefore := map[string]string{}
	for _, phase := range opsPhases {
		path := filepath.Join(agentsDir, phase+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", path, err)
		}
		contentsBefore[phase] = string(data)
	}

	// Second run — must report changed=false.
	changed2, _, err2 := injectOpsSubAgents(home, adapter, InjectOptions{})
	if err2 != nil {
		t.Fatalf("injectOpsSubAgents second run error = %v", err2)
	}
	if changed2 {
		t.Error("injectOpsSubAgents second run changed = true; want false (idempotent)")
	}

	// Content must be byte-identical.
	for _, phase := range opsPhases {
		path := filepath.Join(agentsDir, phase+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", path, err)
		}
		if string(data) != contentsBefore[phase] {
			t.Errorf("injectOpsSubAgents second run changed content for %q", phase)
		}
	}
}

// ── injectOpsSlashCommands tests (scenarios 2, 10, 11, 12, 13, 14) ─────────────

// TestInjectOpsSlashCommands_CapabilityGate — Scenario 2
// Only capable adapters receive command files; Kimi (SupportsSlashCommands=false) does not.
func TestInjectOpsSlashCommands_CapabilityGate(t *testing.T) {
	tests := []struct {
		name        string
		adapter     agents.Adapter
		expectFiles bool
	}{
		{name: "claude", adapter: claude.NewAdapter(), expectFiles: true},
		{name: "kimi", adapter: kimi.NewAdapter(), expectFiles: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			changed, files, err := injectOpsSlashCommands(home, tc.adapter)
			if err != nil {
				t.Fatalf("injectOpsSlashCommands(%s) error = %v", tc.name, err)
			}
			if tc.expectFiles {
				if !changed {
					t.Errorf("injectOpsSlashCommands(%s) changed = false; want true", tc.name)
				}
				if len(files) == 0 {
					t.Errorf("injectOpsSlashCommands(%s) files empty; want command files", tc.name)
				}
			} else {
				if changed {
					t.Errorf("injectOpsSlashCommands(%s) changed = true; want false", tc.name)
				}
				if len(files) != 0 {
					t.Errorf("injectOpsSlashCommands(%s) files = %v; want empty", tc.name, files)
				}
			}
		})
	}
}

// TestInjectOpsSlashCommands_ClaudeFiveCommands — Scenario 10
// Claude receives exactly five OPS command files.
func TestInjectOpsSlashCommands_ClaudeFiveCommands(t *testing.T) {
	home := t.TempDir()
	adapter := claude.NewAdapter()

	changed, files, err := injectOpsSlashCommands(home, adapter)
	if err != nil {
		t.Fatalf("injectOpsSlashCommands(claude) error = %v", err)
	}
	if !changed {
		t.Fatal("injectOpsSlashCommands(claude) changed = false; want true")
	}

	cmdsDir := adapter.CommandsDir(home)
	for _, phase := range opsPhases {
		path := filepath.Join(cmdsDir, phase+".md")
		if _, statErr := os.Stat(path); statErr != nil {
			t.Errorf("claude command %q not found: %v", phase, statErr)
		}
	}
	if len(files) != len(opsPhases) {
		t.Errorf("injectOpsSlashCommands(claude) wrote %d files; want %d", len(files), len(opsPhases))
	}
}

// TestInjectOpsSlashCommands_ClaudeNativeFrontmatter — Scenario 11
// Claude command files have description but NOT agent: or subtask: fields.
func TestInjectOpsSlashCommands_ClaudeNativeFrontmatter(t *testing.T) {
	home := t.TempDir()
	adapter := claude.NewAdapter()

	if _, _, err := injectOpsSlashCommands(home, adapter); err != nil {
		t.Fatalf("injectOpsSlashCommands(claude) error = %v", err)
	}

	cmdsDir := adapter.CommandsDir(home)
	for _, phase := range opsPhases {
		path := filepath.Join(cmdsDir, phase+".md")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", path, err)
		}
		text := string(content)
		if !strings.Contains(text, "description") {
			t.Errorf("claude command %q missing 'description' in frontmatter", phase)
		}
		if strings.Contains(text, "agent: ops-orchestrator") {
			t.Errorf("claude command %q contains 'agent: ops-orchestrator'; want Claude-native frontmatter only", phase)
		}
	}
}

// TestInjectOpsSlashCommands_OpenCodeOrchestratorFrontmatter — Scenario 12
// OpenCode and Kilocode commands have agent: ops-orchestrator and subtask: true.
func TestInjectOpsSlashCommands_OpenCodeOrchestratorFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		adapter agents.Adapter
	}{
		{name: "opencode", adapter: opencode.NewAdapter()},
		{name: "kilocode", adapter: kilocode.NewAdapter()},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			if _, _, err := injectOpsSlashCommands(home, tc.adapter); err != nil {
				t.Fatalf("injectOpsSlashCommands(%s) error = %v", tc.name, err)
			}

			cmdsDir := tc.adapter.CommandsDir(home)
			for _, phase := range opsPhases {
				path := filepath.Join(cmdsDir, phase+".md")
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile(%q) error = %v", path, err)
				}
				text := string(content)
				if !strings.Contains(text, "agent: ops-orchestrator") {
					t.Errorf("%s command %q missing 'agent: ops-orchestrator'", tc.name, phase)
				}
				if !strings.Contains(text, "subtask: true") {
					t.Errorf("%s command %q missing 'subtask: true'", tc.name, phase)
				}
			}
		})
	}
}

// TestInjectOpsSlashCommands_NoCommandsForUnsupported — Scenario 13
// Cursor, Kimi, and Kiro do not receive OPS command files.
func TestInjectOpsSlashCommands_NoCommandsForUnsupported(t *testing.T) {
	tests := []struct {
		name    string
		adapter agents.Adapter
	}{
		{name: "cursor", adapter: cursor.NewAdapter()},
		{name: "kimi", adapter: kimi.NewAdapter()},
		{name: "kiro", adapter: kiro.NewAdapter()},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			home := t.TempDir()
			changed, files, err := injectOpsSlashCommands(home, tc.adapter)
			if err != nil {
				t.Fatalf("injectOpsSlashCommands(%s) error = %v", tc.name, err)
			}
			if changed {
				t.Errorf("injectOpsSlashCommands(%s) changed = true; want false", tc.name)
			}
			if len(files) != 0 {
				t.Errorf("injectOpsSlashCommands(%s) files = %v; want empty", tc.name, files)
			}
		})
	}
}

// TestInjectOpsSlashCommands_Idempotent — Scenario 14
// Second run of injectOpsSlashCommands reports changed=false AND every command
// file is byte-for-byte identical to the content written on the first run.
func TestInjectOpsSlashCommands_Idempotent(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	if _, _, err := injectOpsSlashCommands(home, adapter); err != nil {
		t.Fatalf("injectOpsSlashCommands first run error = %v", err)
	}

	// Snapshot bytes of all five command files after first run.
	cmdsDir := adapter.CommandsDir(home)
	bytesBefore := make(map[string][]byte, len(opsPhases))
	for _, phase := range opsPhases {
		path := filepath.Join(cmdsDir, phase+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) after first run error = %v", path, err)
		}
		bytesBefore[phase] = data
	}

	// Second run — must report changed=false.
	changed2, _, err2 := injectOpsSlashCommands(home, adapter)
	if err2 != nil {
		t.Fatalf("injectOpsSlashCommands second run error = %v", err2)
	}
	if changed2 {
		t.Error("injectOpsSlashCommands second run changed = true; want false (idempotent)")
	}

	// Byte-for-byte equality: command content must remain unchanged.
	for _, phase := range opsPhases {
		path := filepath.Join(cmdsDir, phase+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) after second run error = %v", path, err)
		}
		if string(data) != string(bytesBefore[phase]) {
			t.Errorf("injectOpsSlashCommands second run changed content for command %q; idempotency violated", phase)
		}
	}
}

// ── Cross-cutting Inject() tests (scenarios 3, 4) ──────────────────────────────

// TestInject_SubAgentsAndCommandsAtomicAndPostChecked — Scenario 3
// For a capable adapter (claude), Inject writes sub-agents atomically and passes post-check.
func TestInject_SubAgentsAndCommandsAtomicAndPostChecked(t *testing.T) {
	home := t.TempDir()
	adapter := claude.NewAdapter()

	result, err := Inject(home, adapter, InjectOptions{})
	if err != nil {
		t.Fatalf("Inject(claude) error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject(claude) changed = false; want true on first install")
	}

	// Sub-agents installed.
	agentsDir := adapter.SubAgentsDir(home)
	for _, phase := range opsPhases {
		path := filepath.Join(agentsDir, phase+".md")
		info, statErr := os.Stat(path)
		if statErr != nil {
			t.Errorf("claude sub-agent %q not found after Inject: %v", phase, statErr)
			continue
		}
		if info.Size() < 10 {
			t.Errorf("claude sub-agent %q is %d bytes; want >= 10", phase, info.Size())
		}
	}

	// Slash commands installed.
	cmdsDir := adapter.CommandsDir(home)
	for _, phase := range opsPhases {
		path := filepath.Join(cmdsDir, phase+".md")
		if _, statErr := os.Stat(path); statErr != nil {
			t.Errorf("claude command %q not found after Inject: %v", phase, statErr)
		}
	}
}

// TestInject_SecondRunReportsNoChanges — Scenario 4
// Second identical Inject call reports Changed=false.
func TestInject_SecondRunReportsNoChanges(t *testing.T) {
	home := t.TempDir()
	adapter := claude.NewAdapter()

	if _, err := Inject(home, adapter, InjectOptions{}); err != nil {
		t.Fatalf("Inject first run error = %v", err)
	}

	result2, err2 := Inject(home, adapter, InjectOptions{})
	if err2 != nil {
		t.Fatalf("Inject second run error = %v", err2)
	}
	if result2.Changed {
		t.Error("Inject second run Changed = true; want false (idempotent)")
	}
}

// ── resolveClaudeModelAlias unit tests ─────────────────────────────────────────

// TestResolveClaudeModelAlias covers all four routing branches of the function:
// (a) phase-specific override from the assignments map overrides the preset,
// (b) "default" fallback is used when the phase is absent,
// (c) an invalid/unknown alias in the map is silently filtered out (falls back to preset or default),
// (d) nil assignments fall back to the balanced preset then sonnet.
func TestResolveClaudeModelAlias(t *testing.T) {
	tests := []struct {
		name        string
		assignments map[string]model.ClaudeModelAlias
		phase       string
		want        model.ClaudeModelAlias
	}{
		{
			name: "phase-specific override from assignments map",
			assignments: map[string]model.ClaudeModelAlias{
				"ops-brief": model.ClaudeModelOpus,
				"default":   model.ClaudeModelHaiku,
			},
			phase: "ops-brief",
			want:  model.ClaudeModelOpus,
		},
		{
			name: "default fallback when phase absent from assignments",
			assignments: map[string]model.ClaudeModelAlias{
				"default": model.ClaudeModelHaiku,
			},
			phase: "ops-deliver",
			want:  model.ClaudeModelHaiku,
		},
		{
			name: "invalid alias filtered out — falls back to preset default (sonnet)",
			assignments: map[string]model.ClaudeModelAlias{
				"ops-produce": "not-a-real-alias",
			},
			phase: "ops-produce",
			// ClaudeModelPresetBalanced has no "ops-produce" entry, and "default" = sonnet.
			want: model.ClaudeModelSonnet,
		},
		{
			name:        "nil assignments — balanced preset default fallback is sonnet",
			assignments: nil,
			phase:       "ops-review",
			// ClaudeModelPresetBalanced has no "ops-review" entry, "default" = sonnet.
			want: model.ClaudeModelSonnet,
		},
		{
			name: "sonnet fallback when assignments map is empty and phase unknown",
			assignments: map[string]model.ClaudeModelAlias{},
			phase:       "ops-structure",
			// Empty assignments merged over balanced preset → preset default = sonnet.
			want: model.ClaudeModelSonnet,
		},
		{
			name: "preset phase entry used when assignments nil (sdd-propose=opus in balanced)",
			assignments: nil,
			phase:       "sdd-propose",
			want:        model.ClaudeModelOpus,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := resolveClaudeModelAlias(tt.assignments, tt.phase)
			if got != tt.want {
				t.Errorf("resolveClaudeModelAlias(assignments, %q) = %q; want %q", tt.phase, got, tt.want)
			}
		})
	}
}

// ── Scenario 6 model-routing assertions ────────────────────────────────────────

// TestInjectOpsSubAgents_ModelRoutingFromAssignments — Scenario 6 (extended)
// Verifies that ClaudeModelAssignments and KiroModelAssignments are correctly
// routed into the written frontmatter model values, not just that placeholder
// literals are removed.
func TestInjectOpsSubAgents_ModelRoutingFromAssignments(t *testing.T) {
	tests := []struct {
		name        string
		adapter     agents.Adapter
		opts        InjectOptions
		phase       string
		wantInFile  string // expected resolved model string present in the written file
		wantAbsent  string // placeholder that must NOT appear
	}{
		{
			name:    "claude phase-specific assignment routed to file — scenario 6a",
			adapter: claude.NewAdapter(),
			opts: InjectOptions{
				ClaudeModelAssignments: map[string]model.ClaudeModelAlias{
					"ops-brief": model.ClaudeModelOpus,
					"default":   model.ClaudeModelSonnet,
				},
			},
			phase:      "ops-brief",
			wantInFile: model.ClaudeModelOpus.String(), // "opus"
			wantAbsent: "{{CLAUDE_MODEL}}",
		},
		{
			name:    "claude default assignment used when no phase-specific entry — scenario 6b",
			adapter: claude.NewAdapter(),
			opts: InjectOptions{
				ClaudeModelAssignments: map[string]model.ClaudeModelAlias{
					"default": model.ClaudeModelHaiku,
				},
			},
			phase:      "ops-deliver",
			wantInFile: model.ClaudeModelHaiku.String(), // "haiku"
			wantAbsent: "{{CLAUDE_MODEL}}",
		},
		{
			name:    "kiro phase-specific assignment routed to file — scenario 6c",
			adapter: kiro.NewAdapter(),
			opts: InjectOptions{
				KiroModelAssignments: map[string]model.ClaudeModelAlias{
					"ops-review": model.ClaudeModelOpus,
					"default":    model.ClaudeModelSonnet,
				},
			},
			phase:      "ops-review",
			wantInFile: "claude-opus",  // KiroModelID(opus) contains "claude-opus"
			wantAbsent: "{{KIRO_MODEL}}",
		},
		{
			name:    "kiro nil assignments fall back to sonnet — scenario 6d",
			adapter: kiro.NewAdapter(),
			opts:    InjectOptions{},
			phase:   "ops-structure",
			// KiroModelID(sonnet) for claude-sonnet-4.6 — contains "sonnet"
			wantInFile: "sonnet",
			wantAbsent: "{{KIRO_MODEL}}",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			_, _, err := injectOpsSubAgents(home, tt.adapter, tt.opts)
			if err != nil {
				t.Fatalf("injectOpsSubAgents(%s) error = %v", tt.name, err)
			}

			agentsDir := tt.adapter.SubAgentsDir(home)
			entries, rdErr := os.ReadDir(agentsDir)
			if rdErr != nil {
				t.Fatalf("ReadDir(%q) error = %v", agentsDir, rdErr)
			}

			// Find the file for our target phase.
			var targetContent string
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), tt.phase) {
					data, rErr := os.ReadFile(filepath.Join(agentsDir, entry.Name()))
					if rErr != nil {
						t.Fatalf("ReadFile(%q) error = %v", entry.Name(), rErr)
					}
					targetContent = string(data)
					break
				}
			}
			if targetContent == "" {
				t.Fatalf("no file found for phase %q in %q", tt.phase, agentsDir)
			}

			// The placeholder must be gone.
			if strings.Contains(targetContent, tt.wantAbsent) {
				t.Errorf("file for phase %q still contains placeholder %q", tt.phase, tt.wantAbsent)
			}
			// The resolved model string must be present.
			if !strings.Contains(targetContent, tt.wantInFile) {
				t.Errorf("file for phase %q missing resolved model %q\ncontent:\n%s", tt.phase, tt.wantInFile, targetContent)
			}
		})
	}
}

// ── injectOpsOpenCodeOverlay tests (scenarios 15, 16, 17, 18, 19) ──────────────

// overlayAgentKeys extracts the set of agent keys present under root["agent"]
// in a JSON blob. Returns nil on unmarshal failure.
func overlayAgentKeys(t *testing.T, data []byte) map[string]struct{} {
	t.Helper()
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("overlayAgentKeys: unmarshal error = %v\ndata:\n%s", err, string(data))
	}
	agentRaw, ok := root["agent"]
	if !ok {
		return map[string]struct{}{}
	}
	agentMap, ok := agentRaw.(map[string]any)
	if !ok {
		return map[string]struct{}{}
	}
	keys := make(map[string]struct{}, len(agentMap))
	for k := range agentMap {
		keys[k] = struct{}{}
	}
	return keys
}

// mustReadSettingsJSON reads the settings file written by injectOpsOpenCodeOverlay.
func mustReadSettingsJSON(t *testing.T, adapter agents.Adapter, homeDir string) []byte {
	t.Helper()
	path := adapter.SettingsPath(homeDir)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	return data
}

// TestInjectOpsOpenCodeOverlay_RegistersOrchestratorAndAgents — Scenario 15
// After overlay injection, opencode.json must contain all 6 agent keys:
// ops-orchestrator, ops-brief, ops-structure, ops-produce, ops-review, ops-deliver.
func TestInjectOpsOpenCodeOverlay_RegistersOrchestratorAndAgents(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	changed, files, err := injectOpsOpenCodeOverlay(home, adapter)
	if err != nil {
		t.Fatalf("injectOpsOpenCodeOverlay(opencode) error = %v", err)
	}
	if !changed {
		t.Error("injectOpsOpenCodeOverlay changed = false; want true (new install)")
	}
	if len(files) == 0 {
		t.Error("injectOpsOpenCodeOverlay files empty; want settings path")
	}

	data := mustReadSettingsJSON(t, adapter, home)
	keys := overlayAgentKeys(t, data)

	wantKeys := []string{"ops-orchestrator", "ops-brief", "ops-structure", "ops-produce", "ops-review", "ops-deliver"}
	for _, k := range wantKeys {
		if _, ok := keys[k]; !ok {
			t.Errorf("opencode.json missing agent key %q after overlay; present: %v", k, keys)
		}
	}
}

// TestInjectOpsOpenCodeOverlay_PreservesExistingUserKeys — Scenario 16
// Unrelated pre-existing user keys in opencode.json must survive the overlay merge.
func TestInjectOpsOpenCodeOverlay_PreservesExistingUserKeys(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	// Seed opencode.json with unrelated user key.
	settingsPath := adapter.SettingsPath(home)
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	seedJSON := []byte(`{"foo":"bar","myTeam":{"theme":"dark"}}`)
	if err := os.WriteFile(settingsPath, seedJSON, 0o644); err != nil {
		t.Fatalf("WriteFile seed error = %v", err)
	}

	_, _, err := injectOpsOpenCodeOverlay(home, adapter)
	if err != nil {
		t.Fatalf("injectOpsOpenCodeOverlay error = %v", err)
	}

	data := mustReadSettingsJSON(t, adapter, home)
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}

	// Pre-existing top-level key must survive.
	if v, ok := root["foo"]; !ok || v != "bar" {
		t.Errorf("expected root[\"foo\"] = \"bar\"; got %v (ok=%v)", v, ok)
	}
	// Pre-existing nested key must survive.
	myTeam, ok := root["myTeam"].(map[string]any)
	if !ok {
		t.Fatalf("root[\"myTeam\"] not a map; got %T", root["myTeam"])
	}
	if v, ok2 := myTeam["theme"]; !ok2 || v != "dark" {
		t.Errorf("expected myTeam[\"theme\"] = \"dark\"; got %v (ok=%v)", v, ok2)
	}
	// Agent keys from overlay must also be present.
	keys := overlayAgentKeys(t, data)
	if _, ok := keys["ops-orchestrator"]; !ok {
		t.Error("ops-orchestrator key missing after merge")
	}
}

// TestInjectOpsOpenCodeOverlay_PromptsInlinedNoAbsolutePaths — Scenario 17
// Merged opencode.json must have inline prompts and no absolute path references.
func TestInjectOpsOpenCodeOverlay_PromptsInlinedNoAbsolutePaths(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	_, _, err := injectOpsOpenCodeOverlay(home, adapter)
	if err != nil {
		t.Fatalf("injectOpsOpenCodeOverlay error = %v", err)
	}

	data := mustReadSettingsJSON(t, adapter, home)

	// No sentinel must remain.
	if bytes.Contains(data, []byte("{{OPS_ORCHESTRATOR_PROMPT}}")) {
		t.Error("opencode.json still contains {{OPS_ORCHESTRATOR_PROMPT}} sentinel — prompt not inlined")
	}
	// No absolute path patterns (home dir or /Users/).
	for _, abs := range []string{"/Users/", "$HOME", "/root/"} {
		if bytes.Contains(data, []byte(abs)) {
			t.Errorf("opencode.json contains absolute path %q — prompts must be inline", abs)
		}
	}
	// ops-orchestrator prompt field must be a non-empty string.
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	agentMap := root["agent"].(map[string]any)
	orch := agentMap["ops-orchestrator"].(map[string]any)
	prompt, ok := orch["prompt"].(string)
	if !ok || len(prompt) < 10 {
		t.Errorf("ops-orchestrator prompt is empty or missing; got %q", prompt)
	}
}

// TestInjectOpsOpenCodeOverlay_NonOpenCodeSkipped — Scenario 18
// Non-OpenCode adapters (e.g. Claude) must not receive the overlay.
func TestInjectOpsOpenCodeOverlay_NonOpenCodeSkipped(t *testing.T) {
	home := t.TempDir()
	adapter := claude.NewAdapter()

	changed, files, err := injectOpsOpenCodeOverlay(home, adapter)
	if err != nil {
		t.Fatalf("injectOpsOpenCodeOverlay(claude) error = %v", err)
	}
	if changed {
		t.Error("injectOpsOpenCodeOverlay(claude) changed = true; want false (non-overlay adapter)")
	}
	if len(files) != 0 {
		t.Errorf("injectOpsOpenCodeOverlay(claude) files = %v; want empty", files)
	}

	// Claude's settings file must not exist (no write attempted).
	settingsPath := adapter.SettingsPath(home)
	if _, statErr := os.Stat(settingsPath); statErr == nil {
		t.Errorf("injectOpsOpenCodeOverlay(claude) wrote %q; want no write for non-overlay adapter", settingsPath)
	}
}

// TestInjectOpsOpenCodeOverlay_Idempotent — Scenario 19
// Second run must report changed=false and produce byte-identical opencode.json.
func TestInjectOpsOpenCodeOverlay_Idempotent(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	// First run.
	if _, _, err := injectOpsOpenCodeOverlay(home, adapter); err != nil {
		t.Fatalf("injectOpsOpenCodeOverlay first run error = %v", err)
	}

	// Snapshot bytes after first run.
	bytesBefore := mustReadSettingsJSON(t, adapter, home)

	// Second run — must report changed=false.
	changed2, _, err2 := injectOpsOpenCodeOverlay(home, adapter)
	if err2 != nil {
		t.Fatalf("injectOpsOpenCodeOverlay second run error = %v", err2)
	}
	if changed2 {
		t.Error("injectOpsOpenCodeOverlay second run changed = true; want false (idempotent)")
	}

	// Byte-identical.
	bytesAfter := mustReadSettingsJSON(t, adapter, home)
	if !bytes.Equal(bytesBefore, bytesAfter) {
		t.Errorf("injectOpsOpenCodeOverlay second run produced different bytes\nbefore: %s\nafter:  %s",
			string(bytesBefore), string(bytesAfter))
	}
}

// ── Cross-cutting Inject() overlay tests (scenarios 3, 4, 16, 19) ──────────────

// TestInject_OpenCodeOverlayAtomicAndPostChecked — Scenario 3 (overlay extension)
// Inject for OpenCode writes the overlay and passes semantic post-check.
func TestInject_OpenCodeOverlayAtomicAndPostChecked(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	result, err := Inject(home, adapter, InjectOptions{})
	if err != nil {
		t.Fatalf("Inject(opencode) error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject(opencode) changed = false; want true on first install")
	}

	// Settings file must exist and contain ops-orchestrator + ops-brief keys.
	data := mustReadSettingsJSON(t, adapter, home)
	keys := overlayAgentKeys(t, data)
	for _, k := range []string{"ops-orchestrator", "ops-brief"} {
		if _, ok := keys[k]; !ok {
			t.Errorf("Inject(opencode): opencode.json missing agent key %q after overlay", k)
		}
	}

	// Settings path must appear in result.Files.
	settingsPath := adapter.SettingsPath(home)
	found := false
	for _, f := range result.Files {
		if f == settingsPath {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Inject(opencode) result.Files missing settings path %q; got %v", settingsPath, result.Files)
	}
}

// TestInject_OpenCodeOverlayPreservesUserKeys — Scenario 16 (via Inject)
// Inject merges overlay without clobbering pre-existing user keys.
func TestInject_OpenCodeOverlayPreservesUserKeys(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	// Pre-seed opencode.json with a user key.
	settingsPath := adapter.SettingsPath(home)
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(settingsPath, []byte(`{"userKey":"preserved"}`), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	if _, err := Inject(home, adapter, InjectOptions{}); err != nil {
		t.Fatalf("Inject(opencode) error = %v", err)
	}

	data := mustReadSettingsJSON(t, adapter, home)
	var root map[string]any
	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("unmarshal error = %v", err)
	}
	if v, ok := root["userKey"]; !ok || v != "preserved" {
		t.Errorf("Inject clobbered userKey; got root[\"userKey\"] = %v (ok=%v)", v, ok)
	}
}

// TestInject_OpenCodeOverlayIdempotent — Scenario 19 (via Inject)
// Second Inject call for OpenCode reports Changed=false.
func TestInject_OpenCodeOverlayIdempotent(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	if _, err := Inject(home, adapter, InjectOptions{}); err != nil {
		t.Fatalf("Inject first run error = %v", err)
	}

	result2, err2 := Inject(home, adapter, InjectOptions{})
	if err2 != nil {
		t.Fatalf("Inject second run error = %v", err2)
	}
	if result2.Changed {
		t.Error("Inject second run Changed = true; want false (idempotent)")
	}
}
