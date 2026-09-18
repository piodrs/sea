package editor

import "strings"

func (m *Model) insert(text string) {
	runes := []rune(text)
	tail := append([]rune{}, m.buffer[m.cursor:]...)

	m.buffer = append(m.buffer[:m.cursor], runes...)
	m.buffer = append(m.buffer, tail...)
	m.cursor += len(runes)
}

func (m Model) lineBounds() (int, int) {
	start := m.cursor
	end := m.cursor

	for start > 0 && m.buffer[start-1] != '\n' {
		start--
	}

	for end < len(m.buffer) && m.buffer[end] != '\n' {
		end++
	}

	return start, end
}

func (m Model) position() (int, int) {
	row := strings.Count(string(m.buffer[:m.cursor]), "\n")
	start, _ := m.lineBounds()

	return row, m.cursor - start
}

func (m *Model) moveRows(distance int) {
	_, column := m.position()

	for distance != 0 {
		start, end := m.lineBounds()

		if distance < 0 {
			if start == 0 {
				break
			}

			m.cursor = start - 1
			distance++
		} else {
			if end == len(m.buffer) {
				break
			}

			m.cursor = end + 1
			distance--
		}

		start, end = m.lineBounds()
		m.cursor = min(start+column, end)
	}
}
