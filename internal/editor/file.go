package editor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("open %q: %w", path, err)
	}

	if !utf8.Valid(data) {
		return "", fmt.Errorf("open %q: expected UTF-8 text", path)
	}

	for _, character := range string(data) {
		allowed := character == '\n' || character == '\t' || character == '\r'

		if unicode.IsControl(character) && !allowed {
			return "", fmt.Errorf("open %q: unsupported control character", path)
		}
	}

	return string(data), nil
}

func (m Model) save() error {
	path, err := filepath.EvalSymlinks(m.path)

	if errors.Is(err, os.ErrNotExist) {
		path = m.path
	} else if err != nil {
		return err
	}

	mode := os.FileMode(0644)
	info, err := os.Stat(path)

	if err == nil {
		mode = info.Mode().Perm()
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	file, err := os.CreateTemp(filepath.Dir(path), ".sea-*")

	if err != nil {
		return err
	}

	defer os.Remove(file.Name())
	defer file.Close()

	if err := file.Chmod(mode); err != nil {
		return err
	}

	if _, err := file.WriteString(string(m.buffer)); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.Rename(file.Name(), path)
}
