package styles

import "github.com/charmbracelet/lipgloss"

// Operator color palette — industrial/blueprint identity.
var (
	ColorBase     = lipgloss.Color("#0D1B2A") // midnight blue background
	ColorSurface  = lipgloss.Color("#152A3E") // slightly lighter panel
	ColorOverlay  = lipgloss.Color("#3A5A78") // muted borders/dim
	ColorText     = lipgloss.Color("#E0E6ED") // primary near-white
	ColorSubtext  = lipgloss.Color("#8AA0B4") // dim labels
	ColorLavender = lipgloss.Color("#48CAE4") // cyan blueprint — titles, selected, frame
	ColorGreen    = lipgloss.Color("#7FD962") // success / all OK
	ColorPeach    = lipgloss.Color("#FF6B35") // orange safety — risk/warning/danger accent
	ColorRed      = lipgloss.Color("#F26D78") // errors
	ColorBlue     = lipgloss.Color("#2BA8C9") // deeper blueprint cyan
	ColorMauve    = lipgloss.Color("#48CAE4") // heading accent — same cyan family
	ColorYellow   = lipgloss.Color("#FFD23F") // warnings
	ColorTeal     = lipgloss.Color("#7FD962") // duplicate of green for consistency
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "SelOps " + version + " — Your AI operations engineer"
}

// Pre-built reusable styles.
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorMauve).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorLavender).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorOverlay)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorPeach).
			Bold(true)
)
