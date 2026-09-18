package buffer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

func ResolvePath(path string) (string, error) {
	absolute, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	resolved, err := filepath.EvalSymlinks(absolute)

	if errors.Is(err, os.ErrNotExist) {
		parent, parentErr := filepath.EvalSymlinks(filepath.Dir(absolute))

		if parentErr == nil {
			absolute = filepath.Join(parent, filepath.Base(absolute))
		}

		return absolute, nil
	}

	if err != nil {
		return "", err
	}

	return resolved, nil
}

func Open(path string) (*Buffer, error) {
	path, err := ResolvePath(path)

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}

	if !utf8.Valid(data) {
		return nil, fmt.Errorf("open %q: expected UTF-8 text", path)
	}

	for _, character := range string(data) {
		allowed := character == '\n' || character == '\t' || character == '\r'

		if unicode.IsControl(character) && !allowed {
			return nil, fmt.Errorf("open %q: unsupported control character", path)
		}
	}

	return &Buffer{
		path:  path,
		text:  []rune(string(data)),
		saved: string(data),
	}, nil
}

func (m *Buffer) Save() error {
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

	if _, err := file.WriteString(string(m.text)); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	err = os.Rename(file.Name(), path)

	if err != nil {
		return err
	}

	m.saved = string(m.text)
	m.modified = false

	return nil
}
