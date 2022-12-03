package lizt

import (
	"bufio"
	"fmt"
	"strings"
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

				err = si.SetPointer(0)
				if err != nil {
					return nil, fmt.Errorf("SetPointer: %s -> %w", si.filename, err)
				}

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

// Pointer returns the current pointer.
func (si *StreamIterator) Pointer() uint64 {
	return si.pointer.Load()
}

// SetPointer sets the current pointer.
func (si *StreamIterator) SetPointer(p uint64) error {
	if p > uint64(si.Len()) {
		return fmt.Errorf("pointer: %d / len: %d -> %w", p, si.Len(), ErrPointerOutOfRange)
	}

	// we can assume if it reaches this line that the reader is not nil and the pointer is in a valid range.
	si.unsafePointerPairing(p)
	si.pointer.Store(p)
	return nil
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
