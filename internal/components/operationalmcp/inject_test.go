package operationalmcp

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents/claude"
	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
)

var updateGoldens = flag.Bool("update", false, "update operationalmcp golden files")

// opencodeAdapter returns an OpenCode adapter for tests.
func opencodeAdapter() *opencode.Adapter { return opencode.NewAdapter() }

// claudeAdapter returns a Claude adapter for tests.
func claudeAdapter() *claude.Adapter { return claude.NewAdapter() }

// assertGolden reads (or writes, when -update is set) the golden fixture at
// testdata/<name>.golden and compares it to actual.
func assertGolden(t *testing.T, name string, actual string) {
	t.Helper()
	goldenPath := filepath.Join("testdata", name+".golden")
	if *updateGoldens {
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
		t.Fatalf("ReadFile(%q) error = %v — run with -args -update to generate", goldenPath, err)
	}
	if string(expected) != actual {
		t.Fatalf("golden mismatch for %s\n\nwant:\n%s\n\ngot:\n%s", name, string(expected), actual)
	}
}

// TestInjectEmptyServersWritesDisabledPlaceholder verifies that Inject with an
// empty []ServerSpec writes the documented disabled placeholder JSON and that
// the result is idempotent.
func TestInjectEmptyServersWritesDisabledPlaceholder(t *testing.T) {
	home := t.TempDir()
	adapter := opencodeAdapter()
	configPath := adapter.SettingsPath(home)

	first, err := Inject(home, adapter, nil)
	if err != nil {
		t.Fatalf("Inject(empty servers) first error = %v", err)
	}
	if !first.Changed {
		t.Fatal("Inject(empty servers) first changed = false; want true")
	}

	// Verify placeholder is present.
	assertPlaceholder(t, configPath)

	// Second call must be idempotent.
	second, err := Inject(home, adapter, nil)
	if err != nil {
		t.Fatalf("Inject(empty servers) second error = %v", err)
	}
	if second.Changed {
		t.Fatal("Inject(empty servers) second changed = true; want false (idempotent)")
	}
}

// TestInjectEmptyServersPreservesExistingKeys verifies that MergeJSONObjects
// is used: the placeholder is merged into existing JSON without clobbering
// unrelated keys.
func TestInjectEmptyServersPreservesExistingKeys(t *testing.T) {
	home := t.TempDir()
	adapter := opencodeAdapter()
	configPath := adapter.SettingsPath(home)

	// Pre-create settings with an unrelated key.
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	existing := `{"share": "disabled", "someOtherKey": "preserved"}`
	if err := os.WriteFile(configPath, []byte(existing), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	if _, err := Inject(home, adapter, nil); err != nil {
		t.Fatalf("Inject(empty servers) error = %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}
	text := string(content)

	// Existing keys must be preserved.
	if !strings.Contains(text, `"share"`) {
		t.Error("existing 'share' key was removed by Inject")
	}
	if !strings.Contains(text, `"someOtherKey"`) {
		t.Error("existing 'someOtherKey' key was removed by Inject")
	}

	// Placeholder must be present.
	assertPlaceholder(t, configPath)
}

// TestPlaceholderShape verifies the exact JSON shape of the disabled placeholder.
func TestPlaceholderShape(t *testing.T) {
	home := t.TempDir()
	adapter := opencodeAdapter()
	configPath := adapter.SettingsPath(home)

	if _, err := Inject(home, adapter, nil); err != nil {
		t.Fatalf("Inject error = %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile error = %v", err)
	}

	var root map[string]any
	if err := json.Unmarshal(content, &root); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	mcpServers, ok := root["mcpServers"].(map[string]any)
	if !ok {
		t.Fatalf("mcpServers key missing or not an object; got %T", root["mcpServers"])
	}

	placeholder, ok := mcpServers["selops-rag-placeholder"].(map[string]any)
	if !ok {
		t.Fatalf("mcpServers.selops-rag-placeholder missing; got %#v", mcpServers)
	}

	if disabled, _ := placeholder["disabled"].(bool); !disabled {
		t.Errorf("placeholder.disabled = %v; want true", placeholder["disabled"])
	}

	note, _ := placeholder["note"].(string)
	if note == "" {
		t.Error("placeholder.note is empty; want guidance text")
	}
}

// assertPlaceholder checks that the disabled placeholder entry is present in the config file.
func assertPlaceholder(t *testing.T, configPath string) {
	t.Helper()
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", configPath, err)
	}
	text := string(content)
	if !strings.Contains(text, "selops-rag-placeholder") {
		t.Errorf("config missing selops-rag-placeholder key; got:\n%s", text)
	}
	if !strings.Contains(text, `"disabled"`) {
		t.Errorf("placeholder missing disabled key; got:\n%s", text)
	}
}

// TestInjectGoldenPlaceholderPerAdapter verifies the exact config file output
// when Inject is called with no servers (placeholder case) for each supported
// adapter strategy. This is the primary golden regression guard for operationalmcp.
func TestInjectGoldenPlaceholderPerAdapter(t *testing.T) {
	tests := []struct {
		name       string
		configPath func(home string) string
		inject     func(home string) (InjectionResult, error)
	}{
		{
			name: "opencode-placeholder",
			configPath: func(home string) string {
				return opencodeAdapter().SettingsPath(home)
			},
			inject: func(home string) (InjectionResult, error) {
				return Inject(home, opencodeAdapter(), nil)
			},
		},
		{
			name: "claude-placeholder",
			configPath: func(home string) string {
				return claudeAdapter().SettingsPath(home)
			},
			inject: func(home string) (InjectionResult, error) {
				return Inject(home, claudeAdapter(), nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			result, err := tt.inject(home)
			if err != nil {
				t.Fatalf("Inject() error = %v", err)
			}
			if !result.Changed {
				t.Fatalf("Inject() changed = false; want true")
			}

			configPath := tt.configPath(home)
			content, readErr := os.ReadFile(configPath)
			if readErr != nil {
				t.Fatalf("ReadFile(%q) error = %v", configPath, readErr)
			}
			assertGolden(t, tt.name, string(content))
		})
	}
}

// TestInjectGoldenPopulatedServersOpenCode verifies the exact config file output
// when Inject is called with a populated []ServerSpec for the OpenCode adapter.
// This exercises the non-placeholder path through buildOverlay.
func TestInjectGoldenPopulatedServersOpenCode(t *testing.T) {
	home := t.TempDir()
	adapter := opencodeAdapter()

	servers := []ServerSpec{
		{
			Name:    "selops-rag",
			URL:     "https://rag.selops.internal/mcp",
			EnvRefs: []string{"SELOPS_RAG_API_KEY"},
		},
	}

	result, err := Inject(home, adapter, servers)
	if err != nil {
		t.Fatalf("Inject(populated servers) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(populated servers) changed = false; want true")
	}

	configPath := adapter.SettingsPath(home)
	content, readErr := os.ReadFile(configPath)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", configPath, readErr)
	}

	// Structural assertion: the named server must appear in the config.
	var root map[string]any
	if err := json.Unmarshal(content, &root); err != nil {
		t.Fatalf("Unmarshal config error = %v", err)
	}
	mcpServers, ok := root["mcpServers"].(map[string]any)
	if !ok {
		t.Fatalf("mcpServers key missing or not an object; got %T", root["mcpServers"])
	}
	server, ok := mcpServers["selops-rag"].(map[string]any)
	if !ok {
		t.Fatalf("mcpServers.selops-rag missing; got %#v", mcpServers)
	}
	if url, _ := server["url"].(string); url != "https://rag.selops.internal/mcp" {
		t.Errorf("server url = %q; want https://rag.selops.internal/mcp", url)
	}
	if env, _ := server["env"].(map[string]any); env["SELOPS_RAG_API_KEY"] != "${SELOPS_RAG_API_KEY}" {
		t.Errorf("server env SELOPS_RAG_API_KEY = %v; want ${SELOPS_RAG_API_KEY}", env["SELOPS_RAG_API_KEY"])
	}

	assertGolden(t, "opencode-populated", string(content))
}

// TestInjectGoldenDisabledServerIsFiltered verifies that a disabled ServerSpec
// is excluded from the output and falls back to the placeholder.
func TestInjectGoldenDisabledServerIsFiltered(t *testing.T) {
	home := t.TempDir()
	adapter := opencodeAdapter()

	servers := []ServerSpec{
		{
			Name:     "selops-rag",
			URL:      "https://rag.selops.internal/mcp",
			Disabled: true,
		},
	}

	_, err := Inject(home, adapter, servers)
	if err != nil {
		t.Fatalf("Inject(all disabled) error = %v", err)
	}

	configPath := adapter.SettingsPath(home)
	assertPlaceholder(t, configPath)
}
