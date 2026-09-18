package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (m pane) viewBuffer(focused bool, height int, message string) string {
	if height < 2 {
		return m.statusline(focused, message)
	}

	return m.viewport.View() + "\n" + m.statusline(focused, message)
}

func (m *pane) refresh(focused bool) {
	m.cursor = min(m.cursor, len(m.text))

	replacer := strings.NewReplacer("\t", "    ", "\r", " ")

	row, column := m.position()
	lines := strings.Split(string(m.text), "\n")
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
	caret := cell

	if focused {
		caret = lipgloss.NewStyle().Reverse(true).Render(cell)
	}

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

func (m pane) statusline(focused bool, message string) string {
	width := m.viewport.Width
	style := lipgloss.NewStyle().Reverse(true)

	if !focused {
		style = lipgloss.NewStyle().Faint(true)
	}

	row, column := m.position()
	position := fmt.Sprintf(" %d:%d ", row+1, column+1)

	name := filepath.Base(m.path)

	marker := ""

	if m.modified {
		marker = " [+]"
	}

	left := " " + name + marker

	parts := []string{left, message}

	for index, part := range parts {
		parts[index] = strings.Map(func(character rune) rune {
			if unicode.IsControl(character) {
				return ' '
			}

			return character
		}, part)
	}

	left, message = parts[0], parts[1]
	available := width - ansi.StringWidth(position)

	if available < 1 {
		return style.Render(ansi.Truncate(position, width, ""))
	}

	if message != "" {
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

func (m *Model) View() string {
	if !m.opening || m.height == 0 || m.width == 0 {
		return m.root.view(m.focused, m.message)
	}

	lines := strings.Split(m.root.view(m.focused, m.message), "\n")
	prompt := m.pathInput.View()

	if m.message != "" {
		prompt = m.message + " · " + prompt
	}

	lines[len(lines)-1] = lipgloss.NewStyle().Width(m.width).Render(ansi.Truncate(prompt, m.width, ""))

	return strings.Join(lines, "\n")
}

func (p *pane) view(focused *pane, message string) string {
	if p.width == 0 || p.height == 0 {
		return ""
	}

	if p.buffer != nil {
		if p != focused {
			message = ""
		}

		return p.viewBuffer(p == focused, p.height, message)
	}

	first, second := p.children[0], p.children[1]

	if first.width == 0 || first.height == 0 {
		return second.view(focused, message)
	}

	if second.width == 0 || second.height == 0 {
		return first.view(focused, message)
	}

	style := lipgloss.NewStyle().Faint(true)

	if p.direction == DirectionVertical {
		divider := strings.TrimSuffix(strings.Repeat("│\n", p.height), "\n")

		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			first.view(focused, message),
			style.Render(divider),
			second.view(focused, message),
		)
	}

	divider := style.Render(strings.Repeat("―", p.width))

	return lipgloss.JoinVertical(lipgloss.Left, first.view(focused, message), divider, second.view(focused, message))
}
