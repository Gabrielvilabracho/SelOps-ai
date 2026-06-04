package tui

import (
	"testing"
)

// TestNavigationGoldens snapshots each reachable OPS navigation screen.
// Each case uses newOpsTestModel to guarantee deterministic OPS defaults
// (PresetSelOpsOperational + PersonaOperator, Version="dev", empty Detection).
// SpinnerFrame defaults to 0 and all async state is left at zero-value.
func TestNavigationGoldens(t *testing.T) {
	tests := []struct {
		name   string
		screen Screen
		cursor int
		golden string
	}{
		{
			name:   "welcome screen",
			screen: ScreenWelcome,
			cursor: 0,
			golden: "navigation-welcome.golden",
		},
		{
			name:   "detection screen",
			screen: ScreenDetection,
			cursor: 0,
			golden: "navigation-detection.golden",
		},
		{
			name:   "agents screen",
			screen: ScreenAgents,
			cursor: 0,
			golden: "navigation-agents.golden",
		},
		{
			name:   "persona screen",
			screen: ScreenPersona,
			cursor: 0,
			golden: "navigation-persona.golden",
		},
		{
			name:   "preset screen",
			screen: ScreenPreset,
			cursor: 0,
			golden: "navigation-preset.golden",
		},
		{
			name:   "sdd-mode screen",
			screen: ScreenSDDMode,
			cursor: 0,
			golden: "navigation-sdd-mode.golden",
		},
		{
			name:   "strict-tdd screen",
			screen: ScreenStrictTDD,
			cursor: 0,
			golden: "navigation-strict-tdd.golden",
		},
		{
			name:   "opencode-plugins screen",
			screen: ScreenOpenCodePlugins,
			cursor: 0,
			golden: "navigation-opencode-plugins.golden",
		},
		{
			name:   "skill-picker screen",
			screen: ScreenSkillPicker,
			cursor: 0,
			golden: "navigation-skill-picker.golden",
		},
		{
			name:   "review screen",
			screen: ScreenReview,
			cursor: 0,
			golden: "navigation-review.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newOpsTestModel(t, tt.screen, tt.cursor)
			assertTUIGolden(t, tt.golden, m.View())
		})
	}
}
