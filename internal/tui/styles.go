package tui

import "github.com/charmbracelet/lipgloss"

var (
	sunGold   = lipgloss.Color("#F6C453")
	sunAmber  = lipgloss.Color("#F59E0B")
	sunGreen  = lipgloss.Color("#36D399")
	sunRed    = lipgloss.Color("#FF6B6B")
	ink       = lipgloss.Color("#F8FAFC")
	muted     = lipgloss.Color("#94A3B8")
	panel     = lipgloss.Color("#1E293B")
	panelSoft = lipgloss.Color("#334155")

	brandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(sunGold)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ink)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ink)

	errorStyle = lipgloss.NewStyle().
			Foreground(sunRed)

	successStyle = lipgloss.NewStyle().
			Foreground(sunGreen)

	subtleStyle = lipgloss.NewStyle().
			Foreground(muted)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(sunGold)

	docStyle = lipgloss.NewStyle().
			Margin(1, 2)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(panelSoft).
			Padding(1, 2)

	activePanelStyle = panelStyle.Copy().
				BorderForeground(sunGold)

	optionStyle = lipgloss.NewStyle().
			Padding(0, 1)

	activeOptionStyle = optionStyle.Copy().
				Border(lipgloss.ThickBorder(), false, false, false, true).
				BorderForeground(sunGold)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(sunGold)

	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(sunAmber)

	codeStyle = lipgloss.NewStyle().
			Foreground(ink).
			Background(panel).
			Padding(0, 1)
)
