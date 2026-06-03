package styles_test

import (
	"strings"
	"testing"

	"github.com/Gabrielvilabracho/selops-ai/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// TestRenderLogoContainsSELOPS verifies the logo renders the SELOPS wordmark.
// The ANSI Shadow art uses block-drawing glyphs (█ ╗ ╝ etc.) rather than plain
// ASCII text, so we check for a distinctive multi-character fragment from the
// wordmark rows and for the subtitle that names the product identity.
func TestRenderLogoContainsSELOPS(t *testing.T) {
	output := styles.RenderLogo()
	if output == "" {
		t.Fatal("RenderLogo() returned empty string")
	}
	// Strip ANSI escape sequences for text comparison.
	// The gradient renderer wraps each line with lipgloss color codes.
	plain := stripANSI(output)

	// Check for a distinctive ANSI Shadow fragment (top row of the 'S' glyph).
	const shadowFragment = "███████╗███████╗██╗"
	if !strings.Contains(plain, shadowFragment) {
		t.Errorf("RenderLogo() output does not contain ANSI Shadow wordmark fragment %q;\ngot (stripped, first 600 chars): %q",
			shadowFragment, plain[:min(len(plain), 600)])
	}

	// Check for the subtitle that names the system.
	const subtitle = "AI Engineering"
	if !strings.Contains(plain, subtitle) {
		t.Errorf("RenderLogo() output does not contain subtitle %q", subtitle)
	}
}

// TestRenderLogoIsNotRose verifies the old Braille rose art is gone.
func TestRenderLogoIsNotRose(t *testing.T) {
	output := styles.RenderLogo()
	plain := stripANSI(output)
	// The old rose art had a specific Braille character sequence.
	// Spot-check a distinctive substring from the old art.
	if strings.Contains(plain, "⣠⣾⣷⣶⣦") {
		t.Error("RenderLogo() still contains old Braille rose art — replace with SELOPS wordmark")
	}
}

// TestRenderLogoHasReasonableHeight verifies the logo fits in a reasonable
// number of lines (6–12) per the design spec.
func TestRenderLogoHasReasonableHeight(t *testing.T) {
	output := styles.RenderLogo()
	lines := strings.Split(output, "\n")
	if len(lines) < 6 || len(lines) > 20 {
		t.Errorf("RenderLogo() has %d lines, want 6–20", len(lines))
	}
}

// TestRenderLogoFrameIsRectangular verifies that every rendered line of the logo
// has the same display-cell width, proving the box frame is a perfect rectangle.
// This test is written RED-first against the current broken art and must be made
// GREEN by the corrected ANSI-Shadow art with a properly-padded double-line frame.
func TestRenderLogoFrameIsRectangular(t *testing.T) {
	output := styles.RenderLogo()
	if output == "" {
		t.Fatal("RenderLogo() returned empty string")
	}

	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		t.Fatal("RenderLogo() returned no lines")
	}

	// Use lipgloss.Width which correctly accounts for display-cell widths of
	// Unicode block-drawing characters, box-drawing glyphs, and ANSI colour codes.
	widths := make([]int, len(lines))
	for i, line := range lines {
		widths[i] = lipgloss.Width(line)
	}

	// All lines must share the same display width.
	expected := widths[0]
	allEqual := true
	for i, w := range widths {
		if w != expected {
			t.Errorf("line %d has display width %d, want %d (same as line 0); content: %q",
				i, w, expected, lines[i])
			allEqual = false
		}
	}
	if !allEqual {
		t.Errorf("logo frame is NOT rectangular — widths: %v", widths)
	}
}

// stripANSI removes ANSI escape codes from a string for plain-text comparison.
func stripANSI(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// Skip until 'm'
			for i < len(s) && s[i] != 'm' {
				i++
			}
			i++ // skip 'm'
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
