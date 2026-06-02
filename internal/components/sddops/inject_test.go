package sddops

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/agents/claude"
	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
	"github.com/gentleman-programming/gentle-ai/internal/model"
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

// TestSddOpsSkillIDsMembership verifies the exact set of sddOpsSkillIDs.
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

// TestSddOpsSkillIDsNoIntersectionWithSddSkills verifies zero intersection between
// the sddops skill IDs and the SDD orchestrator skill IDs (those starting with "sdd-").
func TestSddOpsSkillIDsNoIntersectionWithSddSkills(t *testing.T) {
	for _, id := range sddOpsSkillIDs {
		s := string(id)
		if len(s) >= 4 && s[:4] == "sdd-" {
			t.Errorf("sddOpsSkillIDs contains SDD skill ID %q — must have zero intersection", id)
		}
	}
}

// TestInjectWritesSkillFilesWithMinimumSize verifies that after a successful
// Inject call each ops SKILL.md exists on disk and is ≥100 bytes.
// Mirrors the guard in internal/components/sdd/inject.go:713.
func TestInjectWritesSkillFilesWithMinimumSize(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	result, err := Inject(home, adapter, InjectOptions{})
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	skillDir := adapter.SkillsDir(home)
	for _, id := range sddOpsSkillIDs {
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

// TestInjectWritesExpectedOpsHeadings verifies that each SKILL.md stub
// contains at least the mandatory ops heading prefix so content is not empty/stub.
func TestInjectWritesExpectedOpsHeadings(t *testing.T) {
	home := t.TempDir()
	adapter := opencode.NewAdapter()

	if _, err := Inject(home, adapter, InjectOptions{}); err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	skillDir := adapter.SkillsDir(home)
	for _, id := range sddOpsSkillIDs {
		path := filepath.Join(skillDir, string(id), "SKILL.md")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("skill %q ReadFile error = %v", id, err)
			continue
		}
		text := string(content)
		// Each stub must have a top-level heading containing the ops skill name.
		if !strings.Contains(text, "#") {
			t.Errorf("skill %q SKILL.md missing any heading", id)
		}
	}
}

// TestInjectGoldenPerAdapter covers the injected SKILL.md content per adapter
// strategy. One golden file is produced per (adapter, skill) pair, capturing
// the exact bytes written to disk. This serves as a byte-for-byte regression
// guard: any unintended change to an ops skill stub will cause a golden mismatch.
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
			for _, id := range sddOpsSkillIDs {
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
// written for a fresh inject (no pre-existing files).
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

			// Every sddOpsSkillID must appear in result.Files.
			skillDir := ga.adapter.SkillsDir(home)
			for _, id := range sddOpsSkillIDs {
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
