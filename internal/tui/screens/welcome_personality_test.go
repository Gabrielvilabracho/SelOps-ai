package screens_test

import (
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/tui/screens"
)

// TestRenderWelcome_ContainsOperatorPersonalityLine verifies that the welcome
// screen includes the Operator personality status line beneath the tagline.
func TestRenderWelcome_ContainsOperatorPersonalityLine(t *testing.T) {
	output := screens.RenderWelcome(0, "1.0.0", "", nil, true, false, 0, true)
	const wantLine = "On shift. Risk gates armed. Let's keep prod boring."
	if !strings.Contains(output, wantLine) {
		snippet := output
		if len(snippet) > 500 {
			snippet = snippet[:500]
		}
		t.Errorf("RenderWelcome() missing operator personality line %q; output snippet: %q", wantLine, snippet)
	}
}

// TestRenderWelcome_PersonalityLineAppearsBeforeMenu verifies ordering:
// personality line comes before the "Menu" heading.
func TestRenderWelcome_PersonalityLineAppearsBeforeMenu(t *testing.T) {
	output := screens.RenderWelcome(0, "1.0.0", "", nil, true, false, 0, true)
	personalityIdx := strings.Index(output, "On shift.")
	menuIdx := strings.Index(output, "Menu")
	if personalityIdx < 0 {
		t.Fatal("personality line not found in output")
	}
	if menuIdx < 0 {
		t.Fatal("'Menu' heading not found in output")
	}
	if personalityIdx >= menuIdx {
		t.Errorf("personality line (at %d) should appear before 'Menu' (at %d)", personalityIdx, menuIdx)
	}
}
