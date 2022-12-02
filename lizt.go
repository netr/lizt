package lizt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

var ErrNoMoreLines = errors.New("no more lines")

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

// AddIter adds an iterator to the manager.
func (m *Manager) AddIter(i Iterator) error {
	m.files[i.Name()] = i
	return nil
}

// Len returns the length of the iterator.
func (m *Manager) Len() int {
	return len(m.files)
}

// AddDirIter adds a directory of files to the manager.
func (m *Manager) AddDirIter(dir string, roundRobin bool) error {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	files, err := ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		si, err := NewStreamIterator(f, roundRobin)
		if err != nil {
			return err
		}
		err = m.AddIter(si)
		if err != nil {
			return err
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
func (m *Manager) Get(name string) Iterator {
	if m.files[name] == nil {
		return nil
	}
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
