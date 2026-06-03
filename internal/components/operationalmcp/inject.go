// Package operationalmcp installs connection wiring for the SelOps operational
// MCP servers (e.g. a RAG server for organizational knowledge). It writes only
// configuration entries — it does NOT manage any runtime, storage, or RAG service.
package operationalmcp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Gabrielvilabracho/selops-ai/internal/agents"
	"github.com/Gabrielvilabracho/selops-ai/internal/components/filemerge"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
)

// InjectionResult mirrors the shape used by other component packages.
type InjectionResult struct {
	Changed bool
	Files   []string
}

// ServerSpec is a type alias for model.OperationalMCPServerSpec so callers
// within this package and tests can use the short name.
// Defined in model to avoid an import cycle between selection.go and this package.
type ServerSpec = model.OperationalMCPServerSpec

// disabledPlaceholderJSON is written when no servers are configured.
// The key matches the documented placeholder shape.
var disabledPlaceholderJSON = []byte(`{
  "mcpServers": {
    "selops-rag-placeholder": {
      "disabled": true,
      "note": "Configure external RAG MCP server URL or command before enabling."
    }
  }
}
`)

// Inject writes MCP server connection entries for the given adapter.
// When servers is nil or empty (or all entries are disabled), the documented
// disabled placeholder is written instead.
//
// The merge is performed via filemerge.MergeJSONObjects so existing keys in
// the target config file are preserved.
func Inject(homeDir string, adapter agents.Adapter, servers []ServerSpec) (InjectionResult, error) {
	if !adapter.SupportsMCP() {
		return InjectionResult{}, nil
	}

	settingsPath := adapter.SettingsPath(homeDir)
	if settingsPath == "" {
		return InjectionResult{}, nil
	}

	overlay, err := buildOverlay(servers)
	if err != nil {
		return InjectionResult{}, fmt.Errorf("build operationalmcp overlay: %w", err)
	}

	baseJSON, err := osReadFile(settingsPath)
	if err != nil {
		return InjectionResult{}, err
	}

	merged, err := filemerge.MergeJSONObjects(baseJSON, overlay)
	if err != nil {
		return InjectionResult{}, fmt.Errorf("merge operationalmcp overlay: %w", err)
	}

	writeResult, err := filemerge.WriteFileAtomic(settingsPath, merged, 0o644)
	if err != nil {
		return InjectionResult{}, err
	}

	return InjectionResult{Changed: writeResult.Changed, Files: []string{settingsPath}}, nil
}

// buildOverlay produces the JSON overlay to merge. When no active servers are
// provided it returns the disabled placeholder. Otherwise it builds an
// mcpServers map from the provided specs.
func buildOverlay(servers []ServerSpec) ([]byte, error) {
	active := activeServers(servers)
	if len(active) == 0 {
		return disabledPlaceholderJSON, nil
	}

	mcpServers := make(map[string]any, len(active))
	for _, s := range active {
		entry := buildServerEntry(s)
		mcpServers[s.Name] = entry
	}

	root := map[string]any{"mcpServers": mcpServers}
	encoded, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal operationalmcp overlay: %w", err)
	}
	return append(encoded, '\n'), nil
}

// buildServerEntry converts a ServerSpec to a JSON-serialisable map.
func buildServerEntry(s ServerSpec) map[string]any {
	entry := map[string]any{}
	if s.URL != "" {
		entry["url"] = s.URL
	}
	if s.Command != "" {
		entry["command"] = s.Command
	}
	if len(s.EnvRefs) > 0 {
		env := make(map[string]any, len(s.EnvRefs))
		for _, ref := range s.EnvRefs {
			env[ref] = "${" + ref + "}"
		}
		entry["env"] = env
	}
	return entry
}

// activeServers filters out disabled entries.
func activeServers(servers []ServerSpec) []ServerSpec {
	out := make([]ServerSpec, 0, len(servers))
	for _, s := range servers {
		if !s.Disabled {
			out = append(out, s)
		}
	}
	return out
}

var osReadFile = func(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read json file %q: %w", path, err)
	}
	return content, nil
}
