package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575")).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	docStyle = lipgloss.NewStyle().
			Margin(1, 2)
)
