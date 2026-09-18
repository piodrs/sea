package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/piodrs/sea/internal/editor"
)

func Run(path string) error {
	model, err := editor.New(path)

	if err != nil {
		return err
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	_, err = program.Run()

	if err != nil {
		return fmt.Errorf("run editor: %w", err)
	}

	return nil
}
