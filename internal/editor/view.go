package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/piodrs/sea/internal/ui"
)

func (m pane) viewBuffer(focused bool, height int, message string) string {
	if height < 2 {
		return m.statusline(focused, message)
	}

	return m.viewport.View() + "\n" + m.statusline(focused, message)
}

func (m *pane) refresh(focused bool) {
	text := m.buffer.Text()
	m.cursor = min(m.cursor, len(text))

	replacer := strings.NewReplacer("\t", "    ", "\r", " ")

	row, column := m.buffer.Position(m.cursor)
	lines := strings.Split(string(text), "\n")
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
		caret = ui.Selection.Render(cell)
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

func (m *Model) View() string {
	if !m.opening || m.height == 0 || m.width == 0 {
		return m.root.view(m.focused, m.message)
	}

	lines := strings.Split(m.root.view(m.focused, m.message), "\n")
	detailed := m.width >= 16 && m.height >= 4
	width := m.width
	available := m.height - 1

	if detailed {
		available -= 2
	}

	count := min(6, available, max(1, len(m.files)))
	start := max(0, m.selected-count+1)
	rows := []string{}

	if detailed {
		title := fmt.Sprintf(" Open file · %d matches", len(m.files))
		rows = append(rows, ui.Line(title, width, ui.Accent))
	}

	for index := range count {
		label := " No matches · Enter to create"
		style := ui.Muted

		if start+index < len(m.files) {
			style = ui.Accent.Bold(false)
			path := m.files[start+index]
			name := filepath.Base(strings.TrimRight(path, "/"))

			if strings.HasSuffix(path, "/") {
				name += "/"
				style = ui.Accent
			}

			label = "   " + name

			if start+index == m.selected {
				label = " › " + name
				style = ui.Selection
			}
		}

		label = strings.Map(func(character rune) rune {
			if unicode.IsControl(character) {
				return ' '
			}

			return character
		}, label)

		rows = append(rows, ui.Line(label, width, style))
	}

	if detailed {
		rows = append(rows, ui.Line(" ↑↓ select · Tab complete · Enter open", width, ui.Muted))
	}

	copy(lines[len(lines)-1-len(rows):], rows)

	prompt := m.pathInput.View()

	if m.message != "" {
		prompt = m.message + " · " + prompt
	}

	lines[len(lines)-1] = ui.Line(prompt, m.width, lipgloss.NewStyle())

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

	style := ui.Muted

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
