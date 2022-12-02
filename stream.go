package lizt

import (
	"bufio"
	"fmt"
	"sync/atomic"
)

// StreamIterator is an iterator that reads from a file
type StreamIterator struct {
	reader     *bufio.Reader
	pointer    *atomic.Uint64
	filename   string
	name       string
	fileLines  int
	roundRobin bool
}

// NewStreamIterator returns a new stream iterator
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

// Next returns the next line from the iterator
func (si *StreamIterator) Next(count int) ([]string, error) {
	var lines []string
	fmt.Println("calling")
	for i := 1; i <= count; i++ {
		txt, err := si.reader.ReadString('\n')
		if err != nil {
			if si.roundRobin {
				fmt.Println("round", len(lines))
				sr, err := newFileReader(si.filename)
				if err != nil {
					return nil, fmt.Errorf("newFileReader: %s -> %w", si.filename, err)
				}
				si.reader = sr
				txt, err = si.reader.ReadString('\n')
				if err != nil {
					return nil, fmt.Errorf("ReadString(): %s -> %w", si.filename, err)
				}
			} else {
				return nil, fmt.Errorf("file: %s -> %w", si.filename, ErrNoMoreLines)
			}
		}

		lines = append(lines, txt)
	}
	return lines, nil
}

// Pointer returns the current pointer
func (si *StreamIterator) Pointer() uint64 {
	return si.pointer.Load()
}

// Inc increments the pointer
func (si *StreamIterator) Inc() {
	si.pointer.Add(1)
}

// Len returns the length of the iterator
func (si *StreamIterator) Len() int {
	return si.fileLines
}

// Name returns the name of the iterator
func (si *StreamIterator) Name() string {
	return si.name
}

// ResetPointer resets the pointer
func (si *StreamIterator) ResetPointer() {
	si.pointer.Store(0)
}
