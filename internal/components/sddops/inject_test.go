package sddops

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/agents"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/claude"
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
		"ops-adversarial-security",
		"ops-privacy-governance",
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
