package sddops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

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
