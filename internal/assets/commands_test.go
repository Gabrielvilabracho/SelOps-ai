package assets

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// ops phase names used across multiple assertions.
var opsPhases = []string{"ops-brief", "ops-structure", "ops-produce", "ops-review", "ops-deliver"}

// TestOpsCommandsAssetDir verifies that OpsCommandsAssetDir returns the
// correct embedded-FS directory path per agent (mirrors SDDCommandsAssetDir).
func TestOpsCommandsAssetDir(t *testing.T) {
	tests := []struct {
		agent model.AgentID
		want  string
	}{
		{agent: model.AgentClaudeCode, want: "claude/ops-commands"},
		{agent: model.AgentOpenCode, want: "opencode/ops-commands"},
		{agent: model.AgentKiroIDE, want: "opencode/ops-commands"},
		{agent: model.AgentKimi, want: "opencode/ops-commands"},
		{agent: model.AgentCursor, want: "opencode/ops-commands"},
	}

	for _, tt := range tests {
		t.Run(string(tt.agent), func(t *testing.T) {
			got := OpsCommandsAssetDir(tt.agent)
			if got != tt.want {
				t.Errorf("OpsCommandsAssetDir(%q) = %q, want %q", tt.agent, got, tt.want)
			}
		})
	}
}

// TestOpsClaudeSubAgentsExistAndContainCLAUDE_MODEL verifies scenario 5 and 6:
// all 5 Claude OPS sub-agent .md files exist and each contains the
// {{CLAUDE_MODEL}} placeholder (which will be resolved at inject time).
func TestOpsClaudeSubAgentsExistAndContainCLAUDE_MODEL(t *testing.T) {
	for _, phase := range opsPhases {
		phase := phase
		t.Run(phase, func(t *testing.T) {
			path := "claude/agents/" + phase + ".md"
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v (file must exist)", path, err)
			}
			if !strings.Contains(content, "{{CLAUDE_MODEL}}") {
				t.Errorf("%s missing {{CLAUDE_MODEL}} placeholder", path)
			}
			if len(content) < 100 {
				t.Errorf("%s suspiciously short (%d bytes)", path, len(content))
			}
		})
	}
}

// TestOpsKiroSubAgentsExistAndContainKIRO_MODEL verifies scenario 6 for Kiro:
// all 5 Kiro OPS sub-agent .md files exist and contain {{KIRO_MODEL}}.
func TestOpsKiroSubAgentsExistAndContainKIRO_MODEL(t *testing.T) {
	for _, phase := range opsPhases {
		phase := phase
		t.Run(phase, func(t *testing.T) {
			path := "kiro/agents/" + phase + ".md"
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v (file must exist)", path, err)
			}
			if !strings.Contains(content, "{{KIRO_MODEL}}") {
				t.Errorf("%s missing {{KIRO_MODEL}} placeholder", path)
			}
			if len(content) < 100 {
				t.Errorf("%s suspiciously short (%d bytes)", path, len(content))
			}
		})
	}
}

// TestOpsKimiDualFormatMdAndYaml verifies scenario 7:
// Kimi has BOTH .md and .yaml companion files for all 5 phases.
func TestOpsKimiDualFormatMdAndYaml(t *testing.T) {
	for _, phase := range opsPhases {
		phase := phase
		t.Run(phase+".md", func(t *testing.T) {
			path := "kimi/agents/" + phase + ".md"
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v (file must exist)", path, err)
			}
			if len(content) < 100 {
				t.Errorf("%s suspiciously short (%d bytes)", path, len(content))
			}
		})
		t.Run(phase+".yaml", func(t *testing.T) {
			path := "kimi/agents/" + phase + ".yaml"
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v (file must exist)", path, err)
			}
			if len(strings.TrimSpace(content)) == 0 {
				t.Errorf("%s is empty", path)
			}
		})
	}
}

// TestOpsClaudeCommandsHaveDescriptionNotAgentField verifies scenario 11:
// Claude ops-commands have a `description` frontmatter field but NOT
// `agent: ops-orchestrator` (Claude uses a different command format).
func TestOpsClaudeCommandsHaveDescriptionNotAgentField(t *testing.T) {
	for _, phase := range opsPhases {
		phase := phase
		t.Run(phase, func(t *testing.T) {
			path := "claude/ops-commands/" + phase + ".md"
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v (file must exist)", path, err)
			}
			if !strings.Contains(content, "description:") {
				t.Errorf("%s missing 'description:' frontmatter field", path)
			}
			if strings.Contains(content, "agent: ops-orchestrator") {
				t.Errorf("%s must NOT contain 'agent: ops-orchestrator' (Claude-native format)", path)
			}
		})
	}
}

// TestOpsOpenCodeCommandsHaveAgentAndSubtask verifies scenario 12:
// OpenCode ops-commands have `agent: ops-orchestrator` and `subtask: true`
// in their frontmatter.
func TestOpsOpenCodeCommandsHaveAgentAndSubtask(t *testing.T) {
	for _, phase := range opsPhases {
		phase := phase
		t.Run(phase, func(t *testing.T) {
			path := "opencode/ops-commands/" + phase + ".md"
			content, err := Read(path)
			if err != nil {
				t.Fatalf("Read(%q) error = %v (file must exist)", path, err)
			}
			if !strings.Contains(content, "agent: ops-orchestrator") {
				t.Errorf("%s missing 'agent: ops-orchestrator' frontmatter field", path)
			}
			if !strings.Contains(content, "subtask: true") {
				t.Errorf("%s missing 'subtask: true' frontmatter field", path)
			}
		})
	}
}

// TestOpsOverlayJSONSentinelAndSubAgentKeys verifies scenario 15:
// ops-overlay.json contains the {{OPS_ORCHESTRATOR_PROMPT}} sentinel,
// and all 5 sub-agent keys + the orchestrator key are present under "agent".
func TestOpsOverlayJSONSentinelAndSubAgentKeys(t *testing.T) {
	content, err := Read("opencode/ops-overlay.json")
	if err != nil {
		t.Fatalf("Read(opencode/ops-overlay.json) error = %v", err)
	}

	// Must contain the prompt sentinel so injectOpsOpenCodeOverlay can replace it.
	if !strings.Contains(content, "{{OPS_ORCHESTRATOR_PROMPT}}") {
		t.Errorf("ops-overlay.json missing {{OPS_ORCHESTRATOR_PROMPT}} sentinel")
	}

	// Parse JSON to assert structural integrity.
	var root map[string]any
	if err := json.Unmarshal([]byte(content), &root); err != nil {
		t.Fatalf("ops-overlay.json is not valid JSON: %v", err)
	}

	agents, ok := root["agent"].(map[string]any)
	if !ok {
		t.Fatalf("ops-overlay.json missing top-level 'agent' map")
	}

	// Orchestrator must be present.
	if _, ok := agents["ops-orchestrator"]; !ok {
		t.Errorf("ops-overlay.json missing 'ops-orchestrator' key under 'agent'")
	}

	// All 5 sub-agent phase keys must be present.
	for _, phase := range opsPhases {
		if _, ok := agents[phase]; !ok {
			t.Errorf("ops-overlay.json missing %q key under 'agent'", phase)
		}
	}
}
