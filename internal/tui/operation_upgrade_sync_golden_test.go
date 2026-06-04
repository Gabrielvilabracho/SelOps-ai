package tui

import (
	"errors"
	"testing"
	"time"

	"github.com/Gabrielvilabracho/selops-ai/internal/backup"
	"github.com/Gabrielvilabracho/selops-ai/internal/update"
	"github.com/Gabrielvilabracho/selops-ai/internal/update/upgrade"
)

// testManifest returns a deterministic backup.Manifest for golden tests.
func testManifest() backup.Manifest {
	return backup.Manifest{
		ID:               "backup-20260101-120000",
		CreatedAt:        time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
		Source:           backup.BackupSourceInstall,
		FileCount:        3,
		CreatedByVersion: "1.0.0",
	}
}

// TestOperationUpgradeSyncGoldens snapshots Upgrade, Sync, UpgradeSync,
// Backups, Restore, and Rename screens.
func TestOperationUpgradeSyncGoldens(t *testing.T) {
	tests := []struct {
		name    string
		screen  Screen
		cursor  int
		golden  string
		prepare func(m Model) Model
	}{
		// ── Upgrade ─────────────────────────────────────────────────────────────
		{
			name:   "upgrade idle (update check not done)",
			screen: ScreenUpgrade,
			cursor: 0,
			golden: "operation-upgrade-checking.golden",
			prepare: func(m Model) Model {
				m.UpdateCheckDone = false
				m.OperationRunning = false
				m.SpinnerFrame = 0
				return m
			},
		},
		{
			name:   "upgrade ready (update check done, no updates)",
			screen: ScreenUpgrade,
			cursor: 0,
			golden: "operation-upgrade-ready.golden",
			prepare: func(m Model) Model {
				m.UpdateCheckDone = true
				m.OperationRunning = false
				m.UpdateResults = []update.UpdateResult{
					{Tool: update.ToolInfo{Name: "gentle-ai"}, Status: update.UpToDate, InstalledVersion: "1.0.0"},
				}
				m.SpinnerFrame = 0
				return m
			},
		},
		{
			name:   "upgrade running",
			screen: ScreenUpgrade,
			cursor: 0,
			golden: "operation-upgrade-running.golden",
			prepare: func(m Model) Model {
				m.UpdateCheckDone = true
				m.OperationRunning = true
				m.UpgradeReport = nil
				m.SpinnerFrame = 0
				return m
			},
		},
		{
			name:   "upgrade done (report available)",
			screen: ScreenUpgrade,
			cursor: 0,
			golden: "operation-upgrade-done.golden",
			prepare: func(m Model) Model {
				m.UpdateCheckDone = true
				m.OperationRunning = false
				m.UpgradeReport = &upgrade.UpgradeReport{
					Results: []upgrade.ToolUpgradeResult{
						{ToolName: "gentle-ai", Status: upgrade.UpgradeSucceeded, OldVersion: "0.9.0", NewVersion: "1.0.0"},
					},
				}
				m.SpinnerFrame = 0
				return m
			},
		},

		// ── Sync ─────────────────────────────────────────────────────────────────
		{
			name:   "sync confirm (not yet run)",
			screen: ScreenSync,
			cursor: 0,
			golden: "operation-sync-confirm.golden",
			prepare: func(m Model) Model {
				m.OperationRunning = false
				m.HasSyncRun = false
				return m
			},
		},
		{
			name:   "sync done (success, files changed)",
			screen: ScreenSync,
			cursor: 0,
			golden: "operation-sync-done.golden",
			prepare: func(m Model) Model {
				m.OperationRunning = false
				m.HasSyncRun = true
				m.SyncFiles = []string{"~/.config/claude/AGENTS.md", "~/.config/cursor/AGENTS.md"}
				m.SyncErr = nil
				return m
			},
		},

		// ── UpgradeSync ──────────────────────────────────────────────────────────
		{
			name:   "upgrade-sync confirm (check done)",
			screen: ScreenUpgradeSync,
			cursor: 0,
			golden: "operation-upgrade-sync-confirm.golden",
			prepare: func(m Model) Model {
				m.UpdateCheckDone = true
				m.OperationRunning = false
				m.UpgradeReport = nil
				m.SpinnerFrame = 0
				return m
			},
		},
		{
			name:   "upgrade-sync done (both succeeded)",
			screen: ScreenUpgradeSync,
			cursor: 0,
			golden: "operation-upgrade-sync-done.golden",
			prepare: func(m Model) Model {
				m.UpdateCheckDone = true
				m.OperationRunning = false
				m.UpgradeReport = &upgrade.UpgradeReport{
					Results: []upgrade.ToolUpgradeResult{
						{ToolName: "gentle-ai", Status: upgrade.UpgradeSucceeded, OldVersion: "0.9.0", NewVersion: "1.0.0"},
					},
				}
				m.SyncFiles = []string{"~/.config/claude/AGENTS.md"}
				m.SyncErr = nil
				m.SpinnerFrame = 0
				return m
			},
		},

		// ── Backups ───────────────────────────────────────────────────────────────
		{
			name:   "backups empty",
			screen: ScreenBackups,
			cursor: 0,
			golden: "operation-backups-empty.golden",
			prepare: func(m Model) Model {
				m.Backups = nil
				return m
			},
		},
		{
			name:   "backups with items",
			screen: ScreenBackups,
			cursor: 0,
			golden: "operation-backups-list.golden",
			prepare: func(m Model) Model {
				m.Backups = []backup.Manifest{testManifest()}
				m.BackupScroll = 0
				return m
			},
		},

		// ── Restore ───────────────────────────────────────────────────────────────
		{
			name:   "restore-confirm screen",
			screen: ScreenRestoreConfirm,
			cursor: 0,
			golden: "operation-restore-confirm.golden",
			prepare: func(m Model) Model {
				m.SelectedBackup = testManifest()
				return m
			},
		},
		{
			name:   "restore-result success",
			screen: ScreenRestoreResult,
			cursor: 0,
			golden: "operation-restore-result-success.golden",
			prepare: func(m Model) Model {
				m.SelectedBackup = testManifest()
				m.RestoreErr = nil
				return m
			},
		},
		{
			name:   "restore-result error",
			screen: ScreenRestoreResult,
			cursor: 0,
			golden: "operation-restore-result-error.golden",
			prepare: func(m Model) Model {
				m.SelectedBackup = testManifest()
				m.RestoreErr = errors.New("permission denied")
				return m
			},
		},

		// ── Rename Backup ─────────────────────────────────────────────────────────
		{
			name:   "rename-backup screen",
			screen: ScreenRenameBackup,
			cursor: 0,
			golden: "operation-rename-backup.golden",
			prepare: func(m Model) Model {
				m.SelectedBackup = testManifest()
				m.BackupRenameText = ""
				m.BackupRenamePos = 0
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newOpsTestModel(t, tt.screen, tt.cursor)
			if tt.prepare != nil {
				m = tt.prepare(m)
			}
			assertTUIGolden(t, tt.golden, m.View())
		})
	}
}
