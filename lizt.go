package lizt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

var (
	ErrNoMoreLines       = errors.New("no more lines")
	ErrKeyNotFound       = errors.New("key not found")
	ErrPointerOutOfRange = errors.New("pointer out of range")
)

var (
	MaxLinesForSliceIter = 250_000
)

// Manager manages iterators.
type Manager struct {
	files map[string]Iterator
}

// NewManager returns a new manager.
func NewManager() *Manager {
	return &Manager{
		files: make(map[string]Iterator, 0),
	}
}

// List returns a list of the names of the iterators.
func (m *Manager) List() []string {
	var names []string
	for name := range m.files {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// AddIter adds an iterator to the manager.
func (m *Manager) AddIter(i Iterator) *Manager {
	m.files[i.Name()] = i
	return m
}

// AddIters adds a slice of iterators to the manager.
func (m *Manager) AddIters(iters ...Iterator) *Manager {
	for _, iter := range iters {
		m.files[iter.Name()] = iter
	}
	return m
}

// Len returns the length of the iterator.
func (m *Manager) Len() int {
	return len(m.files)
}

// AddDirIter walks a directory of files, converts the files into SliceIterators, and adds them to the manager.
// This will always be faster than SmartAddDirIter(). However, it will not take size into account.
func (m *Manager) AddDirIter(dir string, roundRobin bool) error {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	files, err := ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		name := makeNameFromFilename(f)
		lines, err := ReadFromFile(f)
		if err != nil {
			return fmt.Errorf("read from file: %s -> %w", f, err)
		}
		si := NewSliceIterator(name, lines, roundRobin)
		m.AddIter(si)
	}

	return nil
}

// SmartAddDirIter walks a directory of files, converts the files into Iterators (while taking line count into account), and adds them to the manager.
// Files with less than MaxLinesForSliceIter lines will be SliceIterators, the rest will be StreamIterators.
// This will always be slower than just running AddDirIter(), because we have to count the lines in each file.
func (m *Manager) SmartAddDirIter(dir string, roundRobin bool) error {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	files, err := ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		lines, err := FileLineCount(f)
		if err != nil {
			return fmt.Errorf("count lines from file: %s -> %w", f, err)
		}

		if lines > MaxLinesForSliceIter {
			si, err := NewStreamIterator(f, roundRobin)
			if err != nil {
				return err
			}
			m.AddIter(si)
		} else {
			name := makeNameFromFilename(f)
			lines, err := ReadFromFile(f)
			if err != nil {
				return fmt.Errorf("read from file: %s -> %w", f, err)
			}
			si := NewSliceIterator(name, lines, roundRobin)
			m.AddIter(si)
		}

	}

	return nil
}

func ReadDir(dir string) ([]string, error) {
	readDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}
	var files []string
	for _, entry := range readDir {
		if !entry.IsDir() {
			files = append(files, dir+entry.Name())
		}
	}
	return files, nil
}

// Get returns the next line from the iterator.
func (m *Manager) Get(name string) (Iterator, error) {
	if m.files[name] == nil {
		return nil, fmt.Errorf("key: %s -> %w", name, ErrKeyNotFound)
	}
	return m.files[name], nil
}

// MustGet returns the next line from the iterator.
func (m *Manager) MustGet(name string) Iterator {
	return m.files[name]
}

// makeNameFromFilename takes a filename and returns a name.
func makeNameFromFilename(filename string) string {
	p := path.Clean(filename)
	ps := strings.Split(p, "/")
	p = ps[len(ps)-1]
	ps = strings.Split(p, ".")
	return ps[0]
}

// newFileReader returns a new file reader.
func newFileReader(filename string) (*bufio.Reader, error) {
	file, err := OpenFile(filename)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(file), nil
}
