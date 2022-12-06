package lizt

import "fmt"

// PersistentIterator is an iterator that persists the pointer.
type PersistentIterator struct {
	Persister
	PointerIterator
}

// PersistentIteratorConfig is the config for a persistent iterator.
type PersistentIteratorConfig struct {
	PointerIter PointerIterator
	Persister   Persister
}

// NewPersistentIterator returns a new persistent iterator. It will set the pointer to the last known pointer.
func NewPersistentIterator(cfg PersistentIteratorConfig) (*PersistentIterator, error) {
	if val, err := cfg.Persister.Get(cfg.PointerIter.Name()); err == nil {
		cfg.PointerIter.SetPointer(val)
	}

	return &PersistentIterator{
		PointerIterator: cfg.PointerIter,
		Persister:       cfg.Persister,
	}, nil
}

// Next returns the next line from the iterator.
func (pi *PersistentIterator) Next(count int) ([]string, error) {
	next, err := pi.PointerIterator.Next(count)
	if err != nil {
		return nil, fmt.Errorf("next: name: %s -> %w", pi.Name(), err)
	}

	err = pi.Set(pi.Name(), pi.Pointer())
	if err != nil {
		return nil, err
	}
	return next, nil
}

// MustNext returns the next lines, of a given count, from the iterator. Panics on error.
func (pi *PersistentIterator) MustNext(count int) []string {
	lines, err := pi.Next(count)
	if err != nil {
		panic(err)
	}
	return lines
}

// NextOne returns the next line from the iterator.
func (pi *PersistentIterator) NextOne() (string, error) {
	lines, err := pi.Next(1)
	if err != nil {
		return "", err
	}
	return lines[0], nil
}

// MustNextOne returns the next line from the iterator. Panics on error.
func (pi *PersistentIterator) MustNextOne() string {
	line, err := pi.NextOne()
	if err != nil {
		panic(err)
	}
	return line
}
