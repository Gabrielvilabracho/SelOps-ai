package operationalmcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
)

// opencodeAdapter returns an OpenCode adapter for tests.
func opencodeAdapter() *opencode.Adapter { return opencode.NewAdapter() }

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
