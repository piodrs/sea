package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/piodrs/sea/internal/buffer"
)

func (m *Model) openFile(message tea.Msg) tea.Cmd {
	if message, ok := message.(tea.KeyMsg); ok {
		switch message.String() {
		case "up", "ctrl+p":
			if len(m.files) > 0 {
				if m.selected <= 0 {
					m.selected = len(m.files) - 1
				} else {
					m.selected--
				}
			}

			return nil

		case "down", "ctrl+n":
			if len(m.files) > 0 {
				m.selected = (m.selected + 1) % len(m.files)
			}

			return nil

		case "tab":
			if len(m.files) > 0 {
				m.pathInput.SetValue(m.files[max(0, m.selected)])
				m.pathInput.CursorEnd()
				m.completePath()
			}

			return nil

		case "ctrl+g", "esc":
			m.opening = false
			m.pathInput.Blur()
			m.message = "Canceled"

			return nil

		case "enter":
			path := m.pathInput.Value()

			if m.selected >= 0 && m.selected < len(m.files) {
				path = m.files[m.selected]
			}

			if strings.TrimSpace(path) == "" {
				return nil
			}

			if info, err := os.Stat(path); err == nil && info.IsDir() {
				m.pathInput.SetValue(strings.TrimRight(path, "/") + "/")
				m.pathInput.CursorEnd()
				m.completePath()

				return nil
			}

			path, err := buffer.ResolvePath(path)

			if err != nil {
				m.message = fmt.Sprintf("Open failed: %v", err)

				return nil
			}

			document := m.buffers[path]

			if document == nil {
				document, err = buffer.Open(path)

				if err != nil {
					m.message = fmt.Sprintf("Open failed: %v", err)

					return nil
				}

				m.buffers[path] = document
			}

			m.focused.buffer = document
			m.focused.cursor = 0
			m.focused.left = 0
			m.focused.viewport = viewport.New(m.focused.width, max(1, m.focused.height-1))
			m.opening = false
			m.pathInput.Blur()
			m.message = ""
			m.root.resize(m.width, m.height, m.focused)

			return nil
		}

		m.message = ""
	}

	var command tea.Cmd
	previous := m.pathInput.Value()

	m.pathInput, command = m.pathInput.Update(message)

	if m.pathInput.Value() != previous {
		m.completePath()
	}

	return command
}

func (m *Model) completePath() {
	directory, prefix := filepath.Split(m.pathInput.Value())
	m.files = []string{}
	m.selected = -1
	m.message = ""
	location := directory

	if location == "" {
		location = "."
	}

	entries, err := os.ReadDir(location)

	if err != nil {
		m.message = fmt.Sprintf("List files: %v", err)

		return
	}

	if prefix == "" {
		m.files = append(m.files, directory+"../")
	}

	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), prefix) {
			continue
		}

		path := directory + entry.Name()
		info, err := os.Stat(path)

		if err == nil && info.IsDir() {
			path += "/"
		}

		m.files = append(m.files, path)
	}

	if len(m.files) > 0 {
		m.selected = 0
	}
}
