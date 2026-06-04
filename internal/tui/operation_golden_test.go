package tui

import (
	"errors"
	"testing"
	"time"

	"github.com/Gabrielvilabracho/selops-ai/internal/backup"
	componentuninstall "github.com/Gabrielvilabracho/selops-ai/internal/components/uninstall"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
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

// TestOperationGoldens snapshots each OPS operation/maintenance screen.
// Each case uses newOpsTestModel for deterministic defaults, with model state
// set directly to cover all meaningful render variants.
func TestOperationGoldens(t *testing.T) {
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

		// ── Delete ────────────────────────────────────────────────────────────────
		{
			name:   "delete-confirm screen",
			screen: ScreenDeleteConfirm,
			cursor: 0,
			golden: "operation-delete-confirm.golden",
			prepare: func(m Model) Model {
				m.SelectedBackup = testManifest()
				return m
			},
		},
		{
			name:   "delete-result success",
			screen: ScreenDeleteResult,
			cursor: 0,
			golden: "operation-delete-result-success.golden",
			prepare: func(m Model) Model {
				m.SelectedBackup = testManifest()
				m.DeleteErr = nil
				return m
			},
		},
		{
			name:   "delete-result error",
			screen: ScreenDeleteResult,
			cursor: 0,
			golden: "operation-delete-result-error.golden",
			prepare: func(m Model) Model {
				m.SelectedBackup = testManifest()
				m.DeleteErr = errors.New("no such file or directory")
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

		// ── Uninstall Mode ────────────────────────────────────────────────────────
		{
			name:   "uninstall-mode screen",
			screen: ScreenUninstallMode,
			cursor: 0,
			golden: "operation-uninstall-mode.golden",
		},

		// ── Uninstall (agent selection) ───────────────────────────────────────────
		{
			name:   "uninstall agents screen",
			screen: ScreenUninstall,
			cursor: 0,
			golden: "operation-uninstall-agents.golden",
			prepare: func(m Model) Model {
				m.UninstallAgents = []model.AgentID{model.AgentClaudeCode}
				return m
			},
		},

		// ── Uninstall Components ─────────────────────────────────────────────────
		{
			name:   "uninstall-components screen",
			screen: ScreenUninstallComponents,
			cursor: 0,
			golden: "operation-uninstall-components.golden",
			prepare: func(m Model) Model {
				m.UninstallComponents = []model.ComponentID{model.ComponentEngram}
				return m
			},
		},

		// ── Uninstall Profiles ────────────────────────────────────────────────────
		{
			name:   "uninstall-profiles screen",
			screen: ScreenUninstallProfiles,
			cursor: 0,
			golden: "operation-uninstall-profiles.golden",
			prepare: func(m Model) Model {
				m.UninstallProfilesAvailable = []string{"work", "personal"}
				m.UninstallProfilesToRemove = []string{}
				m.UninstallEngramProjectScopeAvailable = false
				m.UninstallEngramScope = model.EngramUninstallScopeGlobal
				return m
			},
		},

		// ── Uninstall Confirm ─────────────────────────────────────────────────────
		{
			name:   "uninstall-confirm screen",
			screen: ScreenUninstallConfirm,
			cursor: 0,
			golden: "operation-uninstall-confirm.golden",
			prepare: func(m Model) Model {
				m.UninstallMode = model.UninstallModeFull
				m.UninstallAgents = []model.AgentID{model.AgentClaudeCode}
				m.UninstallComponents = []model.ComponentID{model.ComponentEngram}
				m.UninstallProfilesToRemove = []string{}
				m.UninstallEngramScope = model.EngramUninstallScopeGlobal
				m.UninstallEngramProjectScopeAvailable = false
				m.OperationRunning = false
				m.SpinnerFrame = 0
				return m
			},
		},

		// ── Uninstall Result ──────────────────────────────────────────────────────
		{
			name:   "uninstall-result success",
			screen: ScreenUninstallResult,
			cursor: 0,
			golden: "operation-uninstall-result-success.golden",
			prepare: func(m Model) Model {
				m.UninstallResult = componentuninstall.Result{
					RemovedFiles: []string{"/home/user/.config/claude/AGENTS.md"},
				}
				m.UninstallErr = nil
				m.UninstallMode = model.UninstallModeFull
				m.UninstallProfilesToRemove = []string{}
				m.UninstallEngramScope = model.EngramUninstallScopeGlobal
				m.UninstallEngramProjectScopeAvailable = false
				m.SyncCleanInstallFiles = nil
				m.SyncCleanInstallErr = nil
				return m
			},
		},
		{
			name:   "uninstall-result error",
			screen: ScreenUninstallResult,
			cursor: 0,
			golden: "operation-uninstall-result-error.golden",
			prepare: func(m Model) Model {
				m.UninstallResult = componentuninstall.Result{}
				m.UninstallErr = errors.New("uninstall failed: permission denied")
				m.UninstallMode = model.UninstallModeFull
				m.UninstallProfilesToRemove = []string{}
				m.UninstallEngramScope = model.EngramUninstallScopeGlobal
				m.UninstallEngramProjectScopeAvailable = false
				m.SyncCleanInstallFiles = nil
				m.SyncCleanInstallErr = nil
				return m
			},
		},

		// ── Profiles ──────────────────────────────────────────────────────────────
		{
			name:   "profiles empty",
			screen: ScreenProfiles,
			cursor: 0,
			golden: "operation-profiles-empty.golden",
			prepare: func(m Model) Model {
				m.ProfileList = nil
				return m
			},
		},
		{
			name:   "profiles with entries",
			screen: ScreenProfiles,
			cursor: 0,
			golden: "operation-profiles-list.golden",
			prepare: func(m Model) Model {
				m.ProfileList = []model.Profile{
					{Name: "work", OrchestratorModel: model.ModelAssignment{ProviderID: "anthropic", ModelID: "claude-sonnet-4-5"}},
				}
				return m
			},
		},

		// ── ProfileCreate ──────────────────────────────────────────────────────────
		{
			name:   "profile-create step 0 (name input)",
			screen: ScreenProfileCreate,
			cursor: 0,
			golden: "operation-profile-create.golden",
			prepare: func(m Model) Model {
				m.ProfileCreateStep = 0
				m.ProfileNameInput = ""
				m.ProfileNamePos = 0
				m.ProfileNameErr = ""
				m.ProfileEditMode = false
				return m
			},
		},

		// ── ProfileDelete ──────────────────────────────────────────────────────────
		{
			name:   "profile-delete screen",
			screen: ScreenProfileDelete,
			cursor: 0,
			golden: "operation-profile-delete.golden",
			prepare: func(m Model) Model {
				m.ProfileDeleteTarget = "work"
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
