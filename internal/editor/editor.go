package editor

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	root         *pane
	focused      *pane
	prefix       bool
	confirmation string
	message      string
	width        int
	height       int
	buffers      map[string]*buffer
	opening      bool
	pathInput    textinput.Model
}

func New(path string) (*Model, error) {
	path, err := filePath(path)

	if err != nil {
		return nil, err
	}

	document, err := open(path)

	if err != nil {
		return nil, err
	}

	root := &pane{buffer: document, viewport: viewport.New(80, 23)}
	model := &Model{
		root:      root,
		focused:   root,
		width:     80,
		height:    24,
		buffers:   map[string]*buffer{path: document},
		pathInput: textinput.New(),
	}

	model.pathInput.Prompt = "Open: "
	model.pathInput.CharLimit = 0
	root.resize(model.width, model.height, root)

	return model, nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	if message, ok := message.(tea.WindowSizeMsg); ok {
		m.width = max(0, message.Width)
		m.height = max(0, message.Height)
		m.pathInput.Width = max(1, m.width-8)
		m.root.resize(m.width, m.height, m.focused)

		return m, nil
	}

	if m.opening {
		return m, m.openFile(message)
	}

	switch message := message.(type) {
	case tea.KeyMsg:
		input := message.String()

		if input == "ctrl+x" {
			m.prefix = true
			m.message = "C-x"

			return m, nil
		}

		confirmed := m.prefix && m.confirmation == input
		prefixed := m.prefix
		m.prefix = false
		m.confirmation = ""
		m.message = ""

		if input == "ctrl+g" || input == "esc" {
			m.message = "Canceled"

			return m, nil
		}

		if !prefixed {
			m.focused.edit(message, m.root.panes())

			m.root.resize(m.width, m.height, m.focused)

			return m, nil
		}

		panes := m.root.panes()
		index := slices.Index(panes, m.focused)
		next := panes[(index+1)%len(panes)]

		switch input {
		case "ctrl+f":
			m.opening = true
			m.pathInput.SetValue("")
			m.pathInput.Width = max(1, m.width-8)

			return m, m.pathInput.Focus()

		case "ctrl+s":
			err := m.focused.buffer.save()

			if err != nil {
				m.message = fmt.Sprintf("Save failed: %v", err)

				return m, nil
			}

			m.message = "Saved"

		case "ctrl+c", "0":
			modified := false

			if input == "ctrl+c" || len(panes) == 1 {
				for _, buffer := range m.buffers {
					modified = modified || buffer.modified
				}
			}

			if modified && !confirmed {
				m.confirmation = input
				m.message = "Unsaved changes - repeat C-x " + input + " to discard"

				return m, nil
			}

			if input == "ctrl+c" || len(panes) == 1 {
				return m, tea.Quit
			}

			m.root = m.root.remove(m.focused)
			m.focused = next

		case "1":
			m.root = m.focused

		case "o":
			m.focused = next

		case "2", "3":
			direction := DirectionHorizontal

			if input == "3" {
				direction = DirectionVertical
			}

			tooNarrow := direction == DirectionVertical && m.focused.width < 3
			tooShort := direction == DirectionHorizontal && m.focused.height < 5

			if tooNarrow || tooShort {
				m.message = "Not enough room to split"

				return m, nil
			}

			first := *m.focused
			second := *m.focused
			m.focused.buffer = nil
			m.focused.direction = direction
			m.focused.children = [2]*pane{&first, &second}
			m.focused = &second

		default:
			m.message = "Unknown C-x binding"
		}

		m.root.resize(m.width, m.height, m.focused)
	}

	return m, nil
}

func (m *Model) openFile(message tea.Msg) tea.Cmd {
	if message, ok := message.(tea.KeyMsg); ok {
		switch message.String() {
		case "ctrl+g", "esc":
			m.opening = false
			m.pathInput.Blur()
			m.message = "Canceled"

			return nil

		case "enter":
			path := m.pathInput.Value()

			if strings.TrimSpace(path) == "" {
				return nil
			}

			path, err := filePath(path)

			if err != nil {
				m.message = fmt.Sprintf("Open failed: %v", err)

				return nil
			}

			buffer := m.buffers[path]

			if buffer == nil {
				buffer, err = open(path)

				if err != nil {
					m.message = fmt.Sprintf("Open failed: %v", err)

					return nil
				}

				m.buffers[path] = buffer
			}

			m.focused.buffer = buffer
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

	m.pathInput, command = m.pathInput.Update(message)

	return command
}
