package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Layout
	App = lipgloss.NewStyle().Margin(1, 2)

	// Logo
	LogoStyle = lipgloss.NewStyle().
			Foreground(Guac).
			Bold(true)

	// Title bar
	Title = lipgloss.NewStyle().
		Foreground(White).
		Background(Charple).
		Padding(0, 1).
		Bold(true)

	// Input
	InputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Malibu).
			Padding(1).
			Width(60)

	InputPrompt = lipgloss.NewStyle().Foreground(Charple)
	InputText   = lipgloss.NewStyle().Foreground(Ash)

	// Cards
	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Charcoal).
		Padding(0, 1).
		MarginBottom(0)

	CardSelected = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Charple).
			Padding(0, 1).
			MarginBottom(0)

	// Video info
	VideoTitle = lipgloss.NewStyle().
			Foreground(Ash).
			Bold(true)

	VideoTitleSelected = lipgloss.NewStyle().
				Foreground(Dolly).
				Bold(true)

	Channel = lipgloss.NewStyle().
		Foreground(Malibu)

	Duration = lipgloss.NewStyle().
			Foreground(Guac)

	Index = lipgloss.NewStyle().
		Foreground(Squid).
		Width(3)

	IndexSelected = lipgloss.NewStyle().
			Foreground(Zest).
			Bold(true).
			Width(3)

	// Detail view
	DetailTitle = lipgloss.NewStyle().
			Foreground(Dolly).
			Bold(true).
			MarginBottom(1)

	DetailHeader = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Charple).
			Padding(1, 2)

	// Messages
	Success = lipgloss.NewStyle().
		Foreground(Guac).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(Coral).
		Bold(true)

	// Help
	Help = lipgloss.NewStyle().
		Foreground(Squid).
		MarginTop(1)

	// Meta
	Meta   = lipgloss.NewStyle().Foreground(Smoke)
	Muted  = lipgloss.NewStyle().Foreground(Squid)
	Subtle = lipgloss.NewStyle().Foreground(Squid).Italic(true)

	// Download
	DownloadBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Salmon).
			Padding(1, 3).
			Align(lipgloss.Center).
			Width(50)

	// Spinner
	Spinner = lipgloss.NewStyle().Foreground(Salmon)
)

func Separator() string {
	return Muted.Render("  •  ")
}
