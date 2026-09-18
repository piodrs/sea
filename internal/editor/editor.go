package editor

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/piodrs/sea/internal/buffer"
	"github.com/piodrs/sea/internal/ui"

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
	buffers      map[string]*buffer.Buffer
	opening      bool
	pathInput    textinput.Model
	files        []string
	selected     int
	killed       string
}

func New(path string) (*Model, error) {
	document, err := buffer.Open(path)

	if err != nil {
		return nil, err
	}

	path = document.Path()

	root := &pane{buffer: document, viewport: viewport.New(80, 23)}
	model := &Model{
		root:      root,
		focused:   root,
		width:     80,
		height:    24,
		buffers:   map[string]*buffer.Buffer{path: document},
		pathInput: textinput.New(),
	}

	model.pathInput.Prompt = "Open: "
	model.pathInput.CharLimit = 0
	model.pathInput.PromptStyle = ui.Accent
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
			m.focused.edit(message, m.root.panes(), &m.killed)

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
			m.completePath()

			return m, m.pathInput.Focus()

		case "ctrl+s":
			err := m.focused.buffer.Save()

			if err != nil {
				m.message = fmt.Sprintf("Save failed: %v", err)

				return m, nil
			}

			m.message = "Saved"

		case "ctrl+c", "0":
			modified := false

			if input == "ctrl+c" || len(panes) == 1 {
				for _, document := range m.buffers {
					modified = modified || document.Modified()
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
