package tui

import (
	"fmt"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/pipeline"
	"github.com/Gabrielvilabracho/selops-ai/internal/planner"
)

// TestInstallGoldens snapshots each install-flow OPS screen.
// Each case uses newOpsTestModel for deterministic defaults
// (PresetSelOpsOperational + PersonaOperator, Version="dev", empty Detection).
// Where a screen requires specific model state (e.g. Installing needs Progress,
// Complete needs Execution), the state is set directly on the model struct.
func TestInstallGoldens(t *testing.T) {
	tests := []struct {
		name    string
		screen  Screen
		cursor  int
		golden  string
		prepare func(m Model) Model
	}{
		{
			name:   "claude-model-picker screen",
			screen: ScreenClaudeModelPicker,
			cursor: 0,
			golden: "install-claude-model-picker.golden",
		},
		{
			name:   "kiro-model-picker screen",
			screen: ScreenKiroModelPicker,
			cursor: 0,
			golden: "install-kiro-model-picker.golden",
		},
		{
			name:   "model-picker screen",
			screen: ScreenModelPicker,
			cursor: 0,
			golden: "install-model-picker.golden",
		},
		{
			name:   "dependency-tree screen",
			screen: ScreenDependencyTree,
			cursor: 0,
			golden: "install-dependency-tree.golden",
			prepare: func(m Model) Model {
				// Build a minimal non-empty plan so the tree renders content.
				m.DependencyPlan = planner.ResolvedPlan{
					OrderedComponents: nil, // OPS preset: empty plan is the normal state
				}
				return m
			},
		},
		{
			name:   "installing-running screen",
			screen: ScreenInstalling,
			cursor: 0,
			golden: "install-installing-running.golden",
			prepare: func(m Model) Model {
				// Simulate a running install: first step running, rest pending.
				m.Progress = NewProgressState([]string{
					"Install dependencies",
					"Configure selected agents",
					"Inject ecosystem components",
				})
				m.Progress.Start(0)
				m.SpinnerFrame = 0
				return m
			},
		},
		{
			name:   "installing-done screen",
			screen: ScreenInstalling,
			cursor: 0,
			golden: "install-installing-done.golden",
			prepare: func(m Model) Model {
				// Simulate a completed install: all steps succeeded.
				m.Progress = NewProgressState([]string{
					"Install dependencies",
					"Configure selected agents",
					"Inject ecosystem components",
				})
				m.Progress.Mark(0, string(pipeline.StepStatusSucceeded))
				m.Progress.Mark(1, string(pipeline.StepStatusSucceeded))
				m.Progress.Mark(2, string(pipeline.StepStatusSucceeded))
				m.Progress.AppendLog("pipeline completed successfully")
				m.SpinnerFrame = 0
				return m
			},
		},
		{
			name:   "complete-success screen",
			screen: ScreenComplete,
			cursor: 0,
			golden: "install-complete-success.golden",
			prepare: func(m Model) Model {
				// All steps succeeded — Execution has no failed steps.
				m.Execution = pipeline.ExecutionResult{
					Prepare: pipeline.StageResult{
						Steps: []pipeline.StepResult{
							{StepID: "Install dependencies", Status: pipeline.StepStatusSucceeded},
						},
					},
					Apply: pipeline.StageResult{
						Steps: []pipeline.StepResult{
							{StepID: "Configure selected agents", Status: pipeline.StepStatusSucceeded},
							{StepID: "Inject ecosystem components", Status: pipeline.StepStatusSucceeded},
						},
					},
				}
				return m
			},
		},
		{
			name:   "complete-with-failures screen",
			screen: ScreenComplete,
			cursor: 0,
			golden: "install-complete-with-failures.golden",
			prepare: func(m Model) Model {
				// One step failed — triggers the "with failures" render path.
				m.Execution = pipeline.ExecutionResult{
					Prepare: pipeline.StageResult{
						Steps: []pipeline.StepResult{
							{StepID: "Install dependencies", Status: pipeline.StepStatusSucceeded},
						},
					},
					Apply: pipeline.StageResult{
						Steps: []pipeline.StepResult{
							{StepID: "Configure selected agents", Status: pipeline.StepStatusFailed,
								Err: fmt.Errorf("permission denied")},
							{StepID: "Inject ecosystem components", Status: pipeline.StepStatusSkipped},
						},
					},
					Err: fmt.Errorf("pipeline failed"),
				}
				return m
			},
		},
		{
			name:   "model-config screen",
			screen: ScreenModelConfig,
			cursor: 0,
			golden: "install-model-config.golden",
		},
		{
			name:   "opencode-plugin-result screen (no plugins selected)",
			screen: ScreenOpenCodePluginResult,
			cursor: 0,
			golden: "install-opencode-plugin-result.golden",
			prepare: func(m Model) Model {
				// Empty results — "no plugins selected" path for determinism.
				m.OpenCodePluginRegistrationResults = nil
				m.OpenCodePluginRegistrationErr = nil
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
