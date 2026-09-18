package buffer

import (
	"slices"
	"strings"
)

type Buffer struct {
	path     string
	text     []rune
	saved    string
	modified bool
}

func (b *Buffer) Path() string {
	return b.path
}

func (b *Buffer) Text() []rune {
	return slices.Clone(b.text)
}

func (b *Buffer) Modified() bool {
	return b.modified
}

func (b *Buffer) Replace(start, end int, text string) int {
	runes := []rune(text)
	tail := slices.Clone(b.text[end:])
	b.text = append(b.text[:start], runes...)
	b.text = append(b.text, tail...)
	b.modified = string(b.text) != b.saved

	return len(runes)
}

func (m *Buffer) LineBounds(cursor int) (int, int) {
	text := m.text
	start := cursor
	end := cursor

	for start > 0 && text[start-1] != '\n' {
		start--
	}

	for end < len(text) && text[end] != '\n' {
		end++
	}

	return start, end
}

func (m *Buffer) Position(cursor int) (int, int) {
	text := m.text
	row := strings.Count(string(text[:cursor]), "\n")
	start, _ := m.LineBounds(cursor)

	return row, cursor - start
}

func (m *Buffer) MoveRows(cursor, distance int) int {
	text := m.text
	_, column := m.Position(cursor)

	for distance != 0 {
		start, end := m.LineBounds(cursor)

		if distance < 0 {
			if start == 0 {
				break
			}

			cursor = start - 1
			distance++
		} else {
			if end == len(text) {
				break
			}

			cursor = end + 1
			distance--
		}

		start, end = m.LineBounds(cursor)
		cursor = min(start+column, end)
	}

	return cursor
}
