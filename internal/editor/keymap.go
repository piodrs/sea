package editor

import (
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) edit(message tea.KeyMsg) {
	start, end := m.lineBounds()

	switch message.String() {
	case "enter", "ctrl+j":
		m.insert("\n")

	case "tab":
		m.insert("\t")

	case "left", "ctrl+b":
		m.cursor = max(0, m.cursor-1)

	case "right", "ctrl+f":
		m.cursor = min(len(m.buffer), m.cursor+1)

	case "home", "ctrl+a":
		m.cursor = start

	case "end", "ctrl+e":
		m.cursor = end

	case "up", "ctrl+p":
		m.moveRows(-1)

	case "down", "ctrl+n":
		m.moveRows(1)

	case "pgup", "alt+v":
		m.moveRows(-m.viewport.Height)

	case "pgdown", "ctrl+v":
		m.moveRows(m.viewport.Height)

	case "alt+<":
		m.cursor = 0

	case "alt+>":
		m.cursor = len(m.buffer)

	case "backspace", "ctrl+h":
		if m.cursor > 0 {
			m.buffer = append(m.buffer[:m.cursor-1], m.buffer[m.cursor:]...)
			m.cursor--
		}

	case "delete", "ctrl+d":
		if m.cursor < len(m.buffer) {
			m.buffer = append(m.buffer[:m.cursor], m.buffer[m.cursor+1:]...)
		}

	case "ctrl+k":
		if m.cursor == end && end < len(m.buffer) {
			end++
		}

		m.killed = string(m.buffer[m.cursor:end])
		m.buffer = append(m.buffer[:m.cursor], m.buffer[end:]...)

	case "ctrl+y":
		m.insert(m.killed)

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

		m.insert(text)
	}
}
