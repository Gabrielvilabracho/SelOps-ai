package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the SelOps Operator badge logo.
// ANSI Shadow style "SELOPS" wordmark inside a clean double-line box frame.
// All lines are exactly 72 display cells wide (verified by TestRenderLogoFrameIsRectangular).
var logoLines = []string{
	"╔══════════════════════════════════════════════════════════════════════╗",
	"║                                                                      ║",
	"║          ███████╗███████╗██╗      ██████╗ ██████╗ ███████╗           ║",
	"║          ██╔════╝██╔════╝██║     ██╔═══██╗██╔══██╗██╔════╝           ║",
	"║          ███████╗█████╗  ██║     ██║   ██║██████╔╝███████╗           ║",
	"║          ╚════██║██╔══╝  ██║     ██║   ██║██╔═══╝ ╚════██║           ║",
	"║          ███████║███████╗███████╗╚██████╔╝██║     ███████║           ║",
	"║          ╚══════╝╚══════╝╚══════╝ ╚═════╝ ╚═╝     ╚══════╝           ║",
	"║                                                                      ║",
	"║       AI Engineering · Company Operations · Production Systems       ║",
	"║                                                                      ║",
	"╚══════════════════════════════════════════════════════════════════════╝",
}

// gradientColors defines the top-to-bottom gradient for the logo.
// Operator gradient: cyan blueprint fading to deeper blueprint blue.
var gradientColors = []lipgloss.Color{
	ColorLavender, // cyan blueprint top
	ColorLavender, // cyan blueprint mid-top
	ColorBlue,     // deeper blueprint blue mid-bottom
	ColorBlue,     // deeper blueprint blue bottom
}

// RenderLogo returns the SelOps Operator badge logo with a top-to-bottom gradient.
func RenderLogo() string {
	total := len(logoLines)
	if total == 0 {
		return ""
	}

	bands := len(gradientColors)
	var b strings.Builder

	for i, line := range logoLines {
		bandIdx := (i * bands) / total
		if bandIdx >= bands {
			bandIdx = bands - 1
		}
		style := lipgloss.NewStyle().Foreground(gradientColors[bandIdx])
		b.WriteString(style.Render(line))
		if i < total-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}
