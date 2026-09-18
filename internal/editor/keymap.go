package editor

import (
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *pane) edit(message tea.KeyMsg, panes []*pane, killed *string) {
	text := m.buffer.Text()
	start, end := m.buffer.LineBounds(m.cursor)

	switch message.String() {
	case "enter", "ctrl+j":
		m.replace(m.cursor, m.cursor, "\n", panes)

	case "tab":
		m.replace(m.cursor, m.cursor, "\t", panes)

	case "left", "ctrl+b":
		m.cursor = max(0, m.cursor-1)

	case "right", "ctrl+f":
		m.cursor = min(len(text), m.cursor+1)

	case "home", "ctrl+a":
		m.cursor = start

	case "end", "ctrl+e":
		m.cursor = end

	case "up", "ctrl+p":
		m.cursor = m.buffer.MoveRows(m.cursor, -1)

	case "down", "ctrl+n":
		m.cursor = m.buffer.MoveRows(m.cursor, 1)

	case "pgup", "alt+v":
		m.cursor = m.buffer.MoveRows(m.cursor, -m.viewport.Height)

	case "pgdown", "ctrl+v":
		m.cursor = m.buffer.MoveRows(m.cursor, m.viewport.Height)

	case "alt+<":
		m.cursor = 0

	case "alt+>":
		m.cursor = len(text)

	case "backspace", "ctrl+h":
		if m.cursor > 0 {
			m.replace(m.cursor-1, m.cursor, "", panes)
		}

	case "delete", "ctrl+d":
		if m.cursor < len(text) {
			m.replace(m.cursor, m.cursor+1, "", panes)
		}

	case "ctrl+k":
		if m.cursor == end && end < len(text) {
			end++
		}

		*killed = string(text[m.cursor:end])
		m.replace(m.cursor, end, "", panes)

	case "ctrl+y":
		m.replace(m.cursor, m.cursor, *killed, panes)

	default:
		textInput := message.Type == tea.KeyRunes || message.Type == tea.KeySpace

		if !textInput || message.Alt {
			return
		}

		text := strings.Map(func(character rune) rune {
			allowed := character == '\n' || character == '\t'

			if unicode.IsControl(character) && !allowed {
				return -1
			}

			return character
		}, string(message.Runes))

		m.replace(m.cursor, m.cursor, text, panes)
	}
}

func (m *pane) replace(start, end int, text string, panes []*pane) {
	inserted := m.buffer.Replace(start, end, text)

	for _, pane := range panes {
		if pane.buffer != m.buffer {
			continue
		}

		if pane == m || pane.cursor > start && pane.cursor <= end {
			pane.cursor = start + inserted
		} else if pane.cursor > end {
			pane.cursor += inserted - (end - start)
		}
	}
}
