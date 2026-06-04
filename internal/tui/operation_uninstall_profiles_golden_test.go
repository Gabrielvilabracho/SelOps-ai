package tui

import (
	"errors"
	"testing"

	componentuninstall "github.com/Gabrielvilabracho/selops-ai/internal/components/uninstall"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// TestOperationUninstallProfilesGoldens snapshots Delete, Uninstall,
// Profiles, ProfileCreate, and ProfileDelete screens.
func TestOperationUninstallProfilesGoldens(t *testing.T) {
	tests := []struct {
		name    string
		screen  Screen
		cursor  int
		golden  string
		prepare func(m Model) Model
	}{
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
