package tui

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
	"github.com/Gabrielvilabracho/selops-ai/internal/pipeline"
	"github.com/Gabrielvilabracho/selops-ai/internal/system"
	"github.com/Gabrielvilabracho/selops-ai/internal/tui/screens"
)

var updateTUIGoldens = flag.Bool("update", false, "update TUI golden files")

type flowAction struct {
	key       tea.KeyMsg
	cursor    int
	setCursor bool
	prepare   func(Model) Model
}

func TestPresetSelectionNextScreenFlowMatrix(t *testing.T) {
	tests := []struct {
		name       string
		agents     []model.AgentID
		preset     model.PresetID
		wantScreen Screen
		golden     string
	}{
		{
			// OPS fork (Phase 0e): SDD mode screen removed — flows directly to OpenCode plugins.
			name:       "full gentleman with opencode enters plugins (SDD mode skipped)",
			agents:     []model.AgentID{model.AgentOpenCode},
			preset:     model.PresetFullGentleman,
			wantScreen: ScreenOpenCodePlugins,
			golden:     "preset-full-gentleman-opencode-next.golden",
		},
		{
			// OPS fork (Phase 0e): SDD mode screen removed — flows directly to OpenCode plugins.
			name:       "ecosystem only with opencode enters plugins (SDD mode skipped)",
			agents:     []model.AgentID{model.AgentOpenCode},
			preset:     model.PresetEcosystemOnly,
			wantScreen: ScreenOpenCodePlugins,
			golden:     "preset-ecosystem-only-opencode-next.golden",
		},
		{
			name:       "minimal with opencode enters plugin selection",
			agents:     []model.AgentID{model.AgentOpenCode},
			preset:     model.PresetMinimal,
			wantScreen: ScreenOpenCodePlugins,
			golden:     "preset-minimal-opencode-next.golden",
		},
		{
			name:       "custom with opencode enters component selection before plugins",
			agents:     []model.AgentID{model.AgentOpenCode},
			preset:     model.PresetCustom,
			wantScreen: ScreenDependencyTree,
			golden:     "preset-custom-opencode-next.golden",
		},
		{
			// OPS fork (Phase 0e): Strict TDD screen removed — flows directly to dependency tree.
			name:       "full gentleman without opencode enters dependency tree (strict TDD skipped)",
			agents:     []model.AgentID{model.AgentCursor},
			preset:     model.PresetFullGentleman,
			wantScreen: ScreenDependencyTree,
			golden:     "preset-full-gentleman-no-opencode-next.golden",
		},
		{
			// OPS fork (Phase 0e): Strict TDD screen removed — flows directly to dependency tree.
			name:       "ecosystem only without opencode enters dependency tree (strict TDD skipped)",
			agents:     []model.AgentID{model.AgentCursor},
			preset:     model.PresetEcosystemOnly,
			wantScreen: ScreenDependencyTree,
			golden:     "preset-ecosystem-only-no-opencode-next.golden",
		},
		{
			name:       "minimal without opencode enters dependency plan",
			agents:     []model.AgentID{model.AgentCursor},
			preset:     model.PresetMinimal,
			wantScreen: ScreenDependencyTree,
			golden:     "preset-minimal-no-opencode-next.golden",
		},
		{
			name:       "custom without opencode enters component selection",
			agents:     []model.AgentID{model.AgentCursor},
			preset:     model.PresetCustom,
			wantScreen: ScreenDependencyTree,
			golden:     "preset-custom-no-opencode-next.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(system.DetectionResult{}, "dev")
			m.Screen = ScreenPreset
			m.Selection.Agents = tt.agents
			m.Cursor = presetCursor(t, tt.preset)

			updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
			state := updated.(Model)

			if state.Screen != tt.wantScreen {
				t.Fatalf("screen = %v, want %v", state.Screen, tt.wantScreen)
			}
			assertTUIGolden(t, tt.golden, state.View())
		})
	}
}

func TestCustomPresetPostComponentFlowMatrix(t *testing.T) {
	tests := []struct {
		name       string
		agents     []model.AgentID
		components []model.ComponentID
		actions    []flowAction
		wantScreen Screen
		golden     string
	}{
		{
			name:       "opencode with Engram only shows plugins after component selection",
			agents:     []model.AgentID{model.AgentOpenCode},
			components: []model.ComponentID{model.ComponentEngram},
			actions:    []flowAction{{key: tea.KeyMsg{Type: tea.KeyEnter}}},
			wantScreen: ScreenOpenCodePlugins,
			golden:     "custom-opencode-engram-next.golden",
		},
		{
			name:       "opencode with SDD reaches plugins after SDD and strict TDD stages",
			agents:     []model.AgentID{model.AgentOpenCode},
			components: []model.ComponentID{model.ComponentEngram},
			actions: []flowAction{
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // DependencyTree Continue -> SDDMode
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // SDDMode single -> StrictTDD
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // StrictTDD enable -> OpenCode plugins
			},
			wantScreen: ScreenOpenCodePlugins,
			golden:     "custom-opencode-sdd-after-strict-next.golden",
		},
		{
			name:       "opencode with SDD and Skills reaches skill picker after plugins",
			agents:     []model.AgentID{model.AgentOpenCode},
			components: []model.ComponentID{model.ComponentEngram, model.ComponentSkills},
			actions: []flowAction{
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // DependencyTree Continue -> SDDMode
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // SDDMode single -> StrictTDD
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // StrictTDD enable -> OpenCode plugins
				{key: tea.KeyMsg{Type: tea.KeyEnter}, cursor: len(opencodepluginDefinitions()) * 2, setCursor: true}, // OpenCode plugins Continue -> SkillPicker
			},
			wantScreen: ScreenSkillPicker,
			golden:     "custom-opencode-sdd-skills-after-plugins-next.golden",
		},
		{
			name:       "no opencode with SDD and Skills reaches skill picker after strict TDD",
			agents:     []model.AgentID{model.AgentCursor},
			components: []model.ComponentID{model.ComponentEngram, model.ComponentSkills},
			actions: []flowAction{
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // DependencyTree Continue -> StrictTDD
				{key: tea.KeyMsg{Type: tea.KeyEnter}}, // StrictTDD enable -> SkillPicker
			},
			wantScreen: ScreenSkillPicker,
			golden:     "custom-no-opencode-sdd-skills-next.golden",
		},
		{
			name:       "no opencode with Engram only reaches review",
			agents:     []model.AgentID{model.AgentCursor},
			components: []model.ComponentID{model.ComponentEngram},
			actions:    []flowAction{{key: tea.KeyMsg{Type: tea.KeyEnter}}},
			wantScreen: ScreenReview,
			golden:     "custom-no-opencode-engram-next.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(system.DetectionResult{}, "dev")
			m.Screen = ScreenDependencyTree
			m.Selection.Preset = model.PresetCustom
			m.Selection.Agents = tt.agents
			m.Selection.Components = tt.components
			m.Cursor = len(screens.AllComponents())

			state := m
			for _, action := range tt.actions {
				if action.setCursor {
					state.Cursor = action.cursor
				}
				updated, _ := state.Update(action.key)
				state = updated.(Model)
			}

			if state.Screen != tt.wantScreen {
				t.Fatalf("screen = %v, want %v", state.Screen, tt.wantScreen)
			}
			assertTUIGolden(t, tt.golden, state.View())
		})
	}
}


func TestPiOnlyDependencyTreeBackRowReturnsToAgentSelection(t *testing.T) {
	m := NewModel(system.DetectionResult{}, "dev")
	m.Screen = ScreenAgents
	m.Selection.Agents = []model.AgentID{model.AgentPi}
	m.Selection.Components = componentsForPreset(model.PresetFullGentleman, model.PersonaOperator)
	m.Cursor = len(screens.AgentOptions())

	state := applyFlowAction(t, m, flowAction{key: tea.KeyMsg{Type: tea.KeyEnter}})
	if state.Screen != ScreenDependencyTree {
		t.Fatalf("screen = %v, want %v", state.Screen, ScreenDependencyTree)
	}

	state = applyFlowAction(t, state, flowAction{key: tea.KeyMsg{Type: tea.KeyEnter}, cursor: 1, setCursor: true})
	if state.Screen != ScreenAgents {
		t.Fatalf("screen = %v, want %v", state.Screen, ScreenAgents)
	}
}

func applyFlowAction(t *testing.T, state Model, action flowAction) Model {
	t.Helper()
	if action.prepare != nil {
		state = action.prepare(state)
	}
	if action.setCursor {
		state.Cursor = action.cursor
	}
	updated, _ := state.Update(action.key)
	return updated.(Model)
}

func presetCursor(t *testing.T, preset model.PresetID) int {
	t.Helper()
	for idx, option := range screens.PresetOptions() {
		if option == preset {
			return idx
		}
	}
	t.Fatalf("preset %q not found", preset)
	return 0
}

// newOpsTestModel creates a deterministic OPS Model for golden tests.
// It uses NewModel defaults (PresetSelOpsOperational + PersonaOperator),
// overrides Screen to the given screen (if not ScreenUnknown), and sets Cursor.
// Version is always "dev" and Detection is empty for determinism.
//
// Note: hasAgentBuilderEngines() scans PATH at render time and always returns
// false in CI (no agent binaries installed). Goldens capture this state:
// "Create your own Agent (no agents)" is the deterministic label in a clean env.
func newOpsTestModel(t testing.TB, screen Screen, cursor int) Model {
	t.Helper()
	m := NewModel(system.DetectionResult{}, "dev")
	if screen != ScreenUnknown {
		m.Screen = screen
	}
	m.Cursor = cursor
	return m
}

// TestInstallHappyPathFlow_OpsDefaults drives the model through the full Install
// flow from ScreenWelcome to ScreenComplete using OPS defaults
// (PresetSelOpsOperational + PersonaOperator).
//
// Key constraints:
//   - No Init(), no goroutines — all transitions via m.Update() or direct state mutation.
//   - SpinnerFrame=0 and Version="dev" for determinism.
//   - PipelineDoneMsg is synthetic with an empty successful ExecutionResult.
//
// Snapshots are taken at 9 major transition points.
//
// Design note — why steps 4–6 use direct state mutation instead of Update() keypresses:
//   - PersonaOperator is NOT in PersonaOptions() (Gentleman/Neutral/Custom only). The OPS
//     fork hardcodes it via NewModel; it is never reachable by pressing Enter on a persona
//     row. Using Update() would select the wrong persona.
//   - PresetSelOpsOperational is NOT in PresetOptions() (Full/EcosystemOnly/Minimal/Custom
//     only). Same reasoning: the OPS fork hardcodes it via NewModel.
//   - ScreenDependencyTree is reached by calling buildDependencyPlan()+setScreen() directly
//     because the Persona/Preset nav is bypassed entirely in the OPS flow.
//   These three direct mutations are intentional and correct for the OPS fork; they are not
//   test shortcuts but accurate reflections of how the OPS defaults bypass those screens.
func TestInstallHappyPathFlow_OpsDefaults(t *testing.T) {
	// --- Step 1: Welcome ---
	m := newOpsTestModel(t, ScreenWelcome, 0)
	assertTUIGolden(t, "flow-install-01-welcome.golden", m.View())

	// Welcome → Detection (cursor=0 = "Install")
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.Screen != ScreenDetection {
		t.Fatalf("expected ScreenDetection after Welcome Enter, got %v", m.Screen)
	}

	// --- Step 2: Detection ---
	assertTUIGolden(t, "flow-install-02-detection.golden", m.View())

	// Detection → Agents (cursor=0 = "Continue")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.Screen != ScreenAgents {
		t.Fatalf("expected ScreenAgents after Detection Enter, got %v", m.Screen)
	}

	// --- Step 3: Agents ---
	assertTUIGolden(t, "flow-install-03-agents.golden", m.View())

	// --- Step 4: Persona snapshot (OPS defaults — PersonaOperator pre-set) ---
	// PersonaOperator is not in PersonaOptions(), so we snapshot the Persona screen
	// with the OPS pre-configured state directly, then advance manually.
	// See design note at top of this test for why direct mutation is correct here.
	{
		snap := m
		snap.Screen = ScreenPersona
		snap.Cursor = 0
		assertTUIGolden(t, "flow-install-04-persona.golden", snap.View())
	}

	// --- Step 5: Preset snapshot (OPS defaults — PresetSelOpsOperational pre-set) ---
	// PresetSelOpsOperational is not in PresetOptions(), so we snapshot the Preset
	// screen with the OPS pre-configured state directly, then advance manually.
	// See design note at top of this test for why direct mutation is correct here.
	{
		snap := m
		snap.Screen = ScreenPreset
		snap.Cursor = 0
		assertTUIGolden(t, "flow-install-05-preset.golden", snap.View())
	}

	// Advance to DependencyTree (OPS path: build plan directly, skip Preset nav).
	// Selection already has PresetSelOpsOperational + PersonaOperator from NewModel.
	// See design note at top of this test for why direct mutation is correct here.
	m.buildDependencyPlan()
	m.setScreen(ScreenDependencyTree)
	m.Cursor = 0

	// --- Step 6: DependencyTree ---
	assertTUIGolden(t, "flow-install-06-dependency-tree.golden", m.View())

	// DependencyTree → Review (cursor=0 = "Continue")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.Screen != ScreenReview {
		t.Fatalf("expected ScreenReview after DependencyTree Enter, got %v", m.Screen)
	}

	// --- Step 7: Review ---
	assertTUIGolden(t, "flow-install-07-review.golden", m.View())

	// Review → Installing (cursor=0 = "Install")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.Screen != ScreenInstalling {
		t.Fatalf("expected ScreenInstalling after Review Enter, got %v", m.Screen)
	}

	// --- Step 8: Installing (pipeline started, first step running) ---
	assertTUIGolden(t, "flow-install-08-installing.golden", m.View())

	// Inject synthetic PipelineDoneMsg with a successful (empty) ExecutionResult.
	// An empty ExecutionResult has zero steps → ProgressFromExecution produces an
	// empty ProgressState → Percent()==100 → Done()==true.
	updated, _ = m.Update(PipelineDoneMsg{Result: pipeline.ExecutionResult{}})
	m = updated.(Model)

	// Still on ScreenInstalling — Progress.Done()==true now. Press Enter → Complete.
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.Screen != ScreenComplete {
		t.Fatalf("expected ScreenComplete after PipelineDoneMsg + Enter, got %v", m.Screen)
	}

	// --- Step 9: Complete ---
	assertTUIGolden(t, "flow-install-09-complete.golden", m.View())
}

func assertTUIGolden(t *testing.T, name string, actual string) {
	t.Helper()
	goldenPath := filepath.Join("testdata", name)

	if *updateTUIGoldens {
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
		t.Fatalf("ReadFile(%q) error = %v", goldenPath, err)
	}
	if string(expected) != actual {
		t.Fatalf("golden mismatch for %s\n\nexpected:\n%s\n\nactual:\n%s", name, string(expected), actual)
	}
}
