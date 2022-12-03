package lizt

import "fmt"

type PersistentIterator struct {
	Persister
	PointerIterator
}

type PersistentIteratorConfig struct {
	PointerIter PointerIterator
	Persister   Persister
}

func NewPersistentIterator(cfg PersistentIteratorConfig) (*PersistentIterator, error) {
	if val, err := cfg.Persister.Get(cfg.PointerIter.Name()); err == nil {
		err = cfg.PointerIter.SetPointer(val)
		if err != nil {
			return nil, fmt.Errorf("error setting pointer: name: %s / pointer: %d -> %w", cfg.PointerIter.Name(), val, err)
		}
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
		return nil, err
	}

	err = pi.Set(pi.Name(), pi.Pointer())
	if err != nil {
		return nil, err
	}
	return next, nil
}
