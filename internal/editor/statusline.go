package editor

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/x/ansi"
	"github.com/piodrs/sea/internal/ui"
)

func (m pane) statusline(focused bool, message string) string {
	width := m.viewport.Width
	style := ui.ActiveStatus

	if !focused {
		style = ui.InactiveStatus
	}

	row, column := m.buffer.Position(m.cursor)
	position := fmt.Sprintf(" %d:%d ", row+1, column+1)

	name := filepath.Base(m.buffer.Path())

	marker := ""

	if m.buffer.Modified() {
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
