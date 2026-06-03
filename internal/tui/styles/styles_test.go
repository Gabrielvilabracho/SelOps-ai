package styles_test

import (
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/tui/styles"
)

// TestColorPaletteIsOperator verifies the Operator palette hex values
// replaced the old Rose Pine palette.
func TestColorPaletteIsOperator(t *testing.T) {
	// Spot-check: ColorBase must be the midnight blue of the Operator palette,
	// NOT the old Rose Pine base (#191724).
	wantBase := "#0D1B2A"
	if string(styles.ColorBase) != wantBase {
		t.Errorf("ColorBase = %q, want %q (Operator midnight blue)", styles.ColorBase, wantBase)
	}

	wantLavender := "#48CAE4"
	if string(styles.ColorLavender) != wantLavender {
		t.Errorf("ColorLavender = %q, want %q (Operator cyan blueprint)", styles.ColorLavender, wantLavender)
	}

	wantPeach := "#FF6B35"
	if string(styles.ColorPeach) != wantPeach {
		t.Errorf("ColorPeach = %q, want %q (Operator orange safety)", styles.ColorPeach, wantPeach)
	}

	// Verify that old Rose Pine values are gone.
	oldValues := []string{"#191724", "#1f1d2e", "#c4a7e7", "#f6c177", "#9ccfd8", "#ebbcba"}
	for _, old := range oldValues {
		if string(styles.ColorBase) == old ||
			string(styles.ColorSurface) == old ||
			string(styles.ColorLavender) == old ||
			string(styles.ColorPeach) == old ||
			string(styles.ColorGreen) == old ||
			string(styles.ColorMauve) == old {
			t.Errorf("Operator palette still contains old Rose Pine value %q", old)
		}
	}
}

// TestTaglineReturnsSelOpsIdentity verifies the tagline was updated from
// the old "Gentle-AI" phrasing to the SelOps operator identity.
func TestTaglineReturnsSelOpsIdentity(t *testing.T) {
	tagline := styles.Tagline("1.0.0")

	if strings.Contains(tagline, "Gentle-AI") {
		t.Errorf("Tagline() still contains old 'Gentle-AI' branding: %q", tagline)
	}
	if !strings.Contains(tagline, "SelOps") {
		t.Errorf("Tagline() missing 'SelOps' brand: %q", tagline)
	}
	if !strings.Contains(tagline, "1.0.0") {
		t.Errorf("Tagline() missing version: %q", tagline)
	}
	// Must contain the operator tagline copy.
	if !strings.Contains(tagline, "AI operations engineer") {
		t.Errorf("Tagline() missing operator copy 'AI operations engineer': %q", tagline)
	}
}
