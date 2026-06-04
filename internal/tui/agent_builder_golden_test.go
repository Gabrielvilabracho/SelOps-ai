package tui

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/Gabrielvilabracho/selops-ai/internal/agentbuilder"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// TestAgentBuilderGoldens snapshots each OPS agent-builder screen.
// Each case uses newOpsTestModel for deterministic defaults with AgentBuilder
// state set directly to cover all meaningful render paths.
func TestAgentBuilderGoldens(t *testing.T) {
	tests := []struct {
		name    string
		screen  Screen
		cursor  int
		golden  string
		prepare func(m Model) Model
	}{
		// ── Engine selection ──────────────────────────────────────────────────────
		{
			name:   "agent-builder engine screen",
			screen: ScreenAgentBuilderEngine,
			cursor: 0,
			golden: "agent-builder-engine.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.AvailableEngines = []model.AgentID{model.AgentClaudeCode}
				return m
			},
		},

		// ── Prompt input ──────────────────────────────────────────────────────────
		{
			name:   "agent-builder prompt screen (empty)",
			screen: ScreenAgentBuilderPrompt,
			cursor: 0,
			golden: "agent-builder-prompt-empty.golden",
			prepare: func(m Model) Model {
				ta := textarea.New()
				m.AgentBuilder.Textarea = ta
				return m
			},
		},
		{
			name:   "agent-builder prompt screen (with text)",
			screen: ScreenAgentBuilderPrompt,
			cursor: 0,
			golden: "agent-builder-prompt-filled.golden",
			prepare: func(m Model) Model {
				ta := textarea.New()
				ta.SetValue("review CSS for accessibility issues")
				m.AgentBuilder.Textarea = ta
				return m
			},
		},

		// ── SDD integration mode ──────────────────────────────────────────────────
		{
			name:   "agent-builder sdd screen",
			screen: ScreenAgentBuilderSDD,
			cursor: 0,
			golden: "agent-builder-sdd.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.SDDMode = agentbuilder.SDDStandalone
				return m
			},
		},

		// ── SDD phase selection ───────────────────────────────────────────────────
		{
			name:   "agent-builder sdd-phase screen",
			screen: ScreenAgentBuilderSDDPhase,
			cursor: 0,
			golden: "agent-builder-sdd-phase.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.SDDMode = agentbuilder.SDDNewPhase
				return m
			},
		},

		// ── Generating ────────────────────────────────────────────────────────────
		{
			name:   "agent-builder generating screen (running)",
			screen: ScreenAgentBuilderGenerating,
			cursor: 0,
			golden: "agent-builder-generating.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.Generating = true
				m.AgentBuilder.SelectedEngine = model.AgentClaudeCode
				m.AgentBuilder.GenerationErr = nil
				m.SpinnerFrame = 0
				return m
			},
		},
		{
			name:   "agent-builder generating screen (error)",
			screen: ScreenAgentBuilderGenerating,
			cursor: 0,
			golden: "agent-builder-generating-error.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.Generating = false
				m.AgentBuilder.SelectedEngine = model.AgentClaudeCode
				m.AgentBuilder.GenerationErr = errors.New("timeout: engine did not respond")
				m.SpinnerFrame = 0
				return m
			},
		},

		// ── Preview ───────────────────────────────────────────────────────────────
		{
			name:   "agent-builder preview screen",
			screen: ScreenAgentBuilderPreview,
			cursor: 0,
			golden: "agent-builder-preview.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.Generated = &agentbuilder.GeneratedAgent{
					Name:        "a11y-reviewer",
					Title:       "A11y Reviewer",
					Description: "Reviews CSS and HTML for accessibility issues",
					Trigger:     "When reviewing frontend code for accessibility",
					Content:     "# A11y Reviewer\n\nReviews CSS for a11y compliance.\n",
				}
				m.AgentBuilder.SelectedEngine = model.AgentClaudeCode
				m.AgentBuilder.PreviewScroll = 0
				m.AgentBuilder.InstallErr = nil
				m.AgentBuilder.ConflictWarning = ""
				m.Height = 40
				return m
			},
		},

		// ── Installing ────────────────────────────────────────────────────────────
		{
			name:   "agent-builder installing screen (running)",
			screen: ScreenAgentBuilderInstalling,
			cursor: 0,
			golden: "agent-builder-installing.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.Installing = true
				m.AgentBuilder.SelectedEngine = model.AgentClaudeCode
				m.AgentBuilder.InstallErr = nil
				m.SpinnerFrame = 0
				return m
			},
		},

		// ── Complete ──────────────────────────────────────────────────────────────
		{
			name:   "agent-builder complete screen",
			screen: ScreenAgentBuilderComplete,
			cursor: 0,
			golden: "agent-builder-complete.golden",
			prepare: func(m Model) Model {
				m.AgentBuilder.Generated = &agentbuilder.GeneratedAgent{
					Name:    "a11y-reviewer",
					Title:   "A11y Reviewer",
					Trigger: "When reviewing frontend code for accessibility",
				}
				m.AgentBuilder.InstallResults = []agentbuilder.InstallResult{
					{AgentID: model.AgentClaudeCode, Path: "/home/user/.claude/skills/a11y-reviewer/SKILL.md", Success: true},
				}
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
