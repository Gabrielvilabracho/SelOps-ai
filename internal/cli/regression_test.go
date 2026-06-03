package cli

// regression_test.go — DEV byte-for-byte regression for the SelOps operational layer.
//
// Task 4.2 (PR3): Proves that installing PresetSelOpsOperational (OPS) does not
// mutate any asset written by PresetFullGentleman (DEV). Every DEV file must be
// byte-identical before and after an OPS preset install. The two presets must also
// write to disjoint asset path namespaces.

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
	"github.com/Gabrielvilabracho/selops-ai/internal/system"
)

// snapshotDir walks rootDir and returns a map from relative path → sha256 hex digest.
func snapshotDir(t *testing.T, rootDir string) map[string]string {
	t.Helper()
	snap := make(map[string]string)
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(rootDir, path)
		if relErr != nil {
			return relErr
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		sum := sha256.Sum256(content)
		snap[rel] = fmt.Sprintf("%x", sum)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshotDir(%q) error = %v", rootDir, err)
	}
	return snap
}

// devInstallArgs returns the CLI args for a minimal DEV preset install using
// OpenCode only (avoids heavy Kimi/Pi bootstrapping) with a stable, testable
// component set. The --preset flag drives the full-gentleman component resolution.
func devInstallArgs() []string {
	return []string{
		"--agent", "opencode",
		"--preset", "full-gentleman",
		"--persona", "neutral", // avoid persona-specific file variance
	}
}

// opsInstallArgs returns the CLI args for the shipped selops-operational preset.
// This uses --preset selops-operational (the exact path a real user runs) so the
// regression test proves the actual shipped preset preserves the DEV invariant,
// not just a manually curated component subset.
func opsInstallArgs() []string {
	return []string{
		"--agent", "opencode",
		"--preset", string(model.PresetSelOpsOperational),
	}
}

// devPaths returns the set of relative path prefixes expected to be written by
// the DEV preset (PresetFullGentleman). These are the canonical DEV namespaces.
var devPaths = []string{
	"skills/sdd-",
	"skills/_shared/",
	"skills/go-testing/",
	"skills/branch-pr/",
	"skills/chained-pr/",
	"skills/cognitive-doc-design/",
	"skills/comment-writer/",
	"skills/issue-creation/",
	"skills/judgment-day/",
	"skills/skill-creator/",
	"skills/skill-improver/",
	"skills/skill-registry/",
	"skills/work-unit-commits/",
}

// opsPaths returns the set of relative path prefixes expected to be written by
// the OPS preset (PresetSelOpsOperational). These are the canonical OPS namespaces.
var opsPaths = []string{
	"skills/ops-",
	"generic/persona-operator.md",
}

// isDevPath reports whether a relative file path belongs to the DEV namespace.
func isDevPath(rel string) bool {
	rel = filepath.ToSlash(rel)
	for _, prefix := range devPaths {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	return false
}

// isOpsPath reports whether a relative file path belongs to the OPS namespace.
func isOpsPath(rel string) bool {
	rel = filepath.ToSlash(rel)
	for _, prefix := range opsPaths {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	// Also match standalone ops skill directories written under the skills dir.
	// e.g. .config/opencode/skills/ops-*/SKILL.md
	if strings.HasPrefix(rel, ".config/opencode/skills/ops-") {
		return true
	}
	return false
}

// TestDEVPresetByteForByteRegressionAfterOPSInstall is the critical guard.
// It verifies that:
//  1. DEV preset produces a stable set of files in a temp home dir.
//  2. OPS preset, installed in the SAME temp home dir, does NOT mutate any
//     DEV file (every DEV file is byte-identical before and after OPS install).
//  3. OPS preset does NOT add any files to the DEV-owned path namespaces.
//  4. DEV files list is unchanged (OPS must not touch any DEV path).
//  5. The second DEV install (re-plan) is still byte-identical to the first snapshot.
func TestDEVPresetByteForByteRegressionAfterOPSInstall(t *testing.T) {
	// --- Shared test wiring ---
	home := t.TempDir()

	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	// --- Step 1: Install DEV preset and snapshot every produced file ---
	_, err := RunInstall(devInstallArgs(), system.DetectionResult{})
	if err != nil {
		t.Fatalf("RunInstall(DEV preset) step 1 error = %v", err)
	}

	snap1 := snapshotDir(t, home)
	if len(snap1) == 0 {
		t.Fatal("DEV preset install produced zero files — snapshot is empty")
	}

	// --- Step 2: Install OPS preset into the SAME home dir ---
	_, err = RunInstall(opsInstallArgs(), system.DetectionResult{})
	if err != nil {
		t.Fatalf("RunInstall(OPS preset) step 2 error = %v", err)
	}

	snap2 := snapshotDir(t, home)

	// 2a. Every non-infra file in snap1 must still exist and be byte-identical after OPS.
	for rel, hash1 := range snap1 {
		relSlash := filepath.ToSlash(rel)
		if isInfraPath(relSlash) {
			// Shared settings files are intentionally merged across installs.
			// Byte-identity is not required for these paths.
			continue
		}
		hash2, exists := snap2[rel]
		if !exists {
			t.Errorf("OPS install removed DEV file %q", rel)
			continue
		}
		if hash1 != hash2 {
			t.Errorf("OPS install mutated DEV file %q (sha256 before=%s, after=%s)", rel, hash1, hash2)
		}
	}

	// 2b. Any NEW file added by OPS must NOT be in the DEV namespace.
	for rel := range snap2 {
		if _, existedBefore := snap1[rel]; existedBefore {
			continue
		}
		relSlash := filepath.ToSlash(rel)
		if isDevPath(relSlash) {
			t.Errorf("OPS install wrote a new file into the DEV namespace: %q", rel)
		}
	}

	// --- Step 3: Re-install DEV preset — non-infra files must be byte-identical to snap1 ---
	_, err = RunInstall(devInstallArgs(), system.DetectionResult{})
	if err != nil {
		t.Fatalf("RunInstall(DEV preset) step 3 error = %v", err)
	}

	snap3 := snapshotDir(t, home)

	for rel, hash1 := range snap1 {
		relSlash := filepath.ToSlash(rel)
		if isInfraPath(relSlash) {
			continue
		}
		hash3, exists := snap3[rel]
		if !exists {
			t.Errorf("re-install of DEV preset removed file %q that was present in snap1", rel)
			continue
		}
		if hash1 != hash3 {
			t.Errorf("re-install of DEV preset mutated file %q (sha256 snap1=%s, snap3=%s)", rel, hash1, hash3)
		}
	}
}

// TestOPSAndDEVPresetPathNamespacesAreDisjoint verifies that the set of files
// produced by a fresh DEV install and a fresh OPS install share no file paths.
// This is done with two independent home dirs.
func TestOPSAndDEVPresetPathNamespacesAreDisjoint(t *testing.T) {
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	// DEV home
	devHome := t.TempDir()
	osUserHomeDir = func() (string, error) { return devHome, nil }
	if _, err := RunInstall(devInstallArgs(), system.DetectionResult{}); err != nil {
		t.Fatalf("RunInstall(DEV) error = %v", err)
	}
	devSnap := snapshotDir(t, devHome)

	// OPS home
	opsHome := t.TempDir()
	osUserHomeDir = func() (string, error) { return opsHome, nil }
	if _, err := RunInstall(opsInstallArgs(), system.DetectionResult{}); err != nil {
		t.Fatalf("RunInstall(OPS) error = %v", err)
	}
	opsSnap := snapshotDir(t, opsHome)

	// Normalize relative paths for both (they are already relative to respective homes).
	// Any path that appears in BOTH snap sets is a conflict.
	//
	// We only flag collisions when both namespaces write the SAME relative path.
	// Infrastructure paths shared intentionally (e.g. opencode.json, backups/)
	// are allowed — we check only skill-content paths.
	for rel := range opsSnap {
		relSlash := filepath.ToSlash(rel)
		// Skip infrastructure / settings paths that are expected to be shared.
		if isInfraPath(relSlash) {
			continue
		}
		if isDevPath(relSlash) {
			if _, conflict := devSnap[rel]; conflict {
				t.Errorf("OPS path %q conflicts with DEV path namespace", rel)
			}
		}
	}

	for rel := range devSnap {
		relSlash := filepath.ToSlash(rel)
		if isInfraPath(relSlash) {
			continue
		}
		if isOpsPath(relSlash) {
			if _, conflict := opsSnap[rel]; conflict {
				t.Errorf("DEV path %q conflicts with OPS path namespace", rel)
			}
		}
	}
}

// isInfraPath returns true for paths that are expected to be shared / mutated
// by both DEV and OPS installs (settings files, backup manifests, etc.).
// These paths are intentionally excluded from the byte-identical regression check
// because they are designed to be merged across installs.
func isInfraPath(rel string) bool {
	rel = filepath.ToSlash(rel)
	for _, prefix := range []string{
		".gentle-ai/",                    // backup root
		".config/opencode/opencode.json", // settings file merged by all installs
		".config/opencode/mcp/",          // MCP separate-file configs
		".config/opencode/AGENTS.md",     // agent instructions (persona-written)
		".config/opencode/plugins/",      // opencode plugins
	} {
		if strings.HasPrefix(rel, prefix) || rel == prefix {
			return true
		}
	}
	return false
}
