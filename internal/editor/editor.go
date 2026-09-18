package editor

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	path        string
	buffer      []rune
	saved       string
	cursor      int
	left        int
	modified    bool
	confirmQuit bool
	message     string
	prefix      bool
	killed      string
	viewport    viewport.Model
}

func New(path string) (Model, error) {
	text, err := readFile(path)

	if err != nil {
		return Model{}, err
	}

	model := Model{
		path:     path,
		buffer:   []rune(text),
		saved:    text,
		viewport: viewport.New(80, 23),
	}

	model.refresh()

	return model, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = max(1, message.Width)
		m.viewport.Height = max(1, message.Height-1)
		m.refresh()

	case tea.KeyMsg:
		input := message.String()

		if input == "ctrl+g" || input == "esc" {
			m.prefix = false
			m.confirmQuit = false
			m.message = "Canceled"

			return m, nil
		}

		if input == "ctrl+x" {
			m.prefix = true
			m.message = "C-x"

			return m, nil
		}

		if m.prefix {
			m.prefix = false

			switch input {
			case "ctrl+c":
				if !m.modified || m.confirmQuit {
					return m, tea.Quit
				}

				m.confirmQuit = true
				m.message = "Unsaved changes - repeat C-x C-c to discard"

			case "ctrl+s":
				m.confirmQuit = false

				err := m.save()

				if err != nil {
					m.message = fmt.Sprintf("Save failed: %v", err)

					return m, nil
				}

				m.saved = string(m.buffer)
				m.modified = false
				m.message = "Saved"

			default:
				m.confirmQuit = false
				m.message = "Unknown C-x binding"
			}

			return m, nil
		}

		m.confirmQuit = false
		m.edit(message)
		m.modified = string(m.buffer) != m.saved
		m.message = ""
		m.refresh()
	}

	return m, nil
}
