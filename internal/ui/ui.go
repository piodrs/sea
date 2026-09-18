package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	Accent = lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true)

	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	Selection = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("15")).
			Bold(true)

	ActiveStatus = Selection

	InactiveStatus = Muted
)

func Line(text string, width int, style lipgloss.Style) string {
	if width < 1 {
		return ""
	}

	return style.Width(width).Render(ansi.Truncate(text, width, "…"))
}
