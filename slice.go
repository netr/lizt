package lizt

import (
	"fmt"
	"sync/atomic"
)

// SliceIterator is an iterator that reads from a slice.
type SliceIterator struct {
	pointer    *atomic.Uint64
	name       string
	lines      []string
	roundRobin bool
}

// NewSliceIterator returns a new slice iterator.
func NewSliceIterator(name string, lines []string, roundRobin bool) *SliceIterator {
	return &SliceIterator{
		lines:      lines,
		name:       name,
		pointer:    new(atomic.Uint64),
		roundRobin: roundRobin,
	}
}

// Next returns the next line from the iterator.
func (si *SliceIterator) Next(count int) ([]string, error) {
	var lines []string
	for i := 0; i < count; i++ {
		ptr := si.pointer.Load()
		if ptr >= uint64(len(si.lines)) {
			if si.roundRobin {
				si.pointer.Store(1)
				lines = append(lines, si.lines[0])
			} else {
				if len(lines) == 0 {
					return nil, fmt.Errorf("file: %s -> %w", si.name, ErrNoMoreLines)
				}
				return lines, nil
			}
		} else {
			lines = append(lines, si.lines[ptr])
			si.pointer.Add(1)
		}
	}
	return lines, nil
}

// Pointer returns the current pointer.
func (si *SliceIterator) Pointer() uint64 {
	return si.pointer.Load()
}

// Inc increments the pointer.
func (si *SliceIterator) Inc() {
	si.pointer.Add(1)
}

// Len returns the length of the iterator.
func (si *SliceIterator) Len() int {
	return len(si.lines)
}

// Name returns the name of the iterator.
func (si *SliceIterator) Name() string {
	return si.name
}
