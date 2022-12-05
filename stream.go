package lizt

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

// StreamIterator is an iterator that reads from a file.
type StreamIterator struct {
	reader     *bufio.Reader
	pointer    *atomic.Uint64
	filename   string
	name       string
	fileLines  int
	roundRobin bool
	mu         sync.RWMutex
}

// NewStreamIterator returns a new stream iterator.
func NewStreamIterator(filename string, roundRobin bool) (*StreamIterator, error) {
	count, err := FileLineCount(filename)
	if err != nil {
		return nil, err
	}

	rdr, err := newFileReader(filename)
	if err != nil {
		return nil, err
	}

	name := makeNameFromFilename(filename)

	return &StreamIterator{
		filename:   filename,
		name:       name,
		reader:     rdr,
		fileLines:  count,
		pointer:    new(atomic.Uint64),
		roundRobin: roundRobin,
	}, nil
}

// Next returns the next line from the iterator.
func (si *StreamIterator) Next(count int) ([]string, error) {
	si.mu.Lock()
	defer si.mu.Unlock()

	var lines []string
	for i := 1; i <= count; i++ {
		txt, err := si.reader.ReadString('\n')
		if err != nil {
			if si.roundRobin {
				sr, err := newFileReader(si.filename)
				if err != nil {
					return nil, fmt.Errorf("newFileReader: %s -> %w", si.filename, err)
				}
				si.reader = sr

				si.SetPointer(0)

				txt, err = si.reader.ReadString('\n')
				if err != nil {
					return nil, fmt.Errorf("ReadString(): %s -> %w", si.filename, err)
				}
			} else {
				if len(lines) == 0 {
					return nil, fmt.Errorf("file: %s -> %w", si.filename, ErrNoMoreLines)
				}
				return lines, nil
			}
		}

		lines = append(lines, strings.TrimSpace(txt))
		si.Inc()
	}
	return lines, nil
}

// MustNext returns the next lines, of a given count, from the iterator. Panics on error.
func (si *StreamIterator) MustNext(count int) []string {
	lines, err := si.Next(count)
	if err != nil {
		panic(err)
	}
	return lines
}

// NextOne returns the next line from the iterator.
func (si *StreamIterator) NextOne() (string, error) {
	lines, err := si.Next(1)
	if err != nil {
		return "", err
	}
	return lines[0], nil
}

// MustNextOne returns the next line from the iterator. Panics on error.
func (si *StreamIterator) MustNextOne() string {
	line, err := si.NextOne()
	if err != nil {
		panic(err)
	}
	return line
}

// Pointer returns the current pointer.
func (si *StreamIterator) Pointer() uint64 {
	return si.pointer.Load()
}

// SetPointer sets the current pointer.
func (si *StreamIterator) SetPointer(p uint64) {
	if p > uint64(si.Len()) {
		si.pointer.Store(0)
		return
	}

	// we can assume if it reaches this line that the reader is not nil and the pointer is in a valid range.
	si.unsafePointerPairing(p)
	si.pointer.Store(p)
}

// unsafePointerPairing create a new reader and iterate through the lines until it reaches the pointer.
func (si *StreamIterator) unsafePointerPairing(p uint64) {
	sr, err := newFileReader(si.filename)
	// even though this is unsafe, we'll just do nothing if there is an error.
	if err != nil {
		return
	}
	si.reader = sr
	var i uint64
	for i = 0; i < p; i++ {
		_, _ = si.reader.ReadString('\n')
	}
}

// Inc increments the pointer.
func (si *StreamIterator) Inc() {
	si.pointer.Add(1)
}

// Len returns the length of the iterator.
func (si *StreamIterator) Len() int {
	return si.fileLines
}

// Name returns the name of the iterator.
func (si *StreamIterator) Name() string {
	return si.name
}
