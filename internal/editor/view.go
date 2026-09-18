package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (m Model) View() string {
	return m.viewport.View() + "\n" + m.statusline()
}

func (m *Model) refresh() {
	replacer := strings.NewReplacer("\t", "    ", "\r", " ")

	row, column := m.position()
	lines := strings.Split(string(m.buffer), "\n")
	current := []rune(lines[row])
	prefix := replacer.Replace(string(current[:column]))
	cell := " "
	suffix := ""

	if column < len(current) {
		cell = replacer.Replace(string(current[column]))
		suffix = replacer.Replace(string(current[column+1:]))
	}

	width := ansi.StringWidth(prefix)
	m.left = min(m.left, width)
	m.left = max(m.left, width+max(1, ansi.StringWidth(cell))-m.viewport.Width)
	caret := lipgloss.NewStyle().Reverse(true).Render(cell)

	for index, line := range lines {
		lines[index] = replacer.Replace(line)
	}

	lines[row] = prefix + caret + suffix

	for index, line := range lines {
		lines[index] = ansi.Cut(line, m.left, m.left+m.viewport.Width)
	}

	m.viewport.SetContent(strings.Join(lines, "\n"))

	if row < m.viewport.YOffset {
		m.viewport.SetYOffset(row)
	}

	if row >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(row - m.viewport.Height + 1)
	}
}

func (m Model) statusline() string {
	width := m.viewport.Width
	style := lipgloss.NewStyle().Reverse(true)
	row, column := m.position()
	position := fmt.Sprintf(" %d:%d ", row+1, column+1)

	name := strings.Map(func(character rune) rune {
		if unicode.IsControl(character) {
			return ' '
		}

		return character
	}, filepath.Base(m.path))

	marker := ""

	if m.modified {
		marker = " [+]"
	}

	left := " " + name + marker
	available := width - ansi.StringWidth(position)

	if available < 1 {
		return style.Render(ansi.Truncate(position, width, ""))
	}

	if m.message != "" {
		message := strings.Map(func(character rune) rune {
			if unicode.IsControl(character) {
				return ' '
			}

			return character
		}, m.message)

		minimumNameWidth := min(ansi.StringWidth(left), 20)
		messageWidth := available - minimumNameWidth - 3

		if messageWidth > 0 {
			message = ansi.Truncate(message, messageWidth, "…")
			nameWidth := available - ansi.StringWidth(message) - 3
			left = ansi.Truncate(left, nameWidth, "…") + " - " + message
		}
	}

	left = ansi.Truncate(left, available, "…")
	padding := strings.Repeat(" ", available-ansi.StringWidth(left))

	return style.Render(left + padding + position)
}
