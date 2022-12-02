package lizt

import (
	"bufio"
	"errors"
	"path"
	"strings"
)

var (
	ErrNoMoreLines = errors.New("no more lines")
)

// Manager manages iterators
type Manager struct {
	files map[string]Iterator
}

// NewManager returns a new manager
func NewManager() *Manager {
	return &Manager{
		files: map[string]Iterator{},
	}
}

// AddIter adds an iterator to the manager
func (m *Manager) AddIter(i Iterator) error {
	m.files[i.Name()] = i
	return nil
}

// Get returns the next line from the iterator
func (m *Manager) Get(name string) Iterator {
	if m.files[name] == nil {
		return nil
	}
	return m.files[name]
}

// makeNameFromFilename takes a filename and returns a name
func makeNameFromFilename(filename string) string {
	p := path.Clean(filename)
	ps := strings.Split(p, "/")
	p = ps[len(ps)-1]
	ps = strings.Split(p, ".")
	return ps[0]
}

// newFileReader returns a new file reader
func newFileReader(filename string) (*bufio.Reader, error) {
	file, err := OpenFile(filename)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(file), nil
}
