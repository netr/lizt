package lizt

type PersistentIterator struct {
	Persister
	PointerIterator
}

type PersistentIteratorConfig struct {
	PointerIter PointerIterator
	Persister   Persister
}

func NewPersistentIterator(cfg PersistentIteratorConfig) *PersistentIterator {
	return &PersistentIterator{
		PointerIterator: cfg.PointerIter,
		Persister:       cfg.Persister,
	}
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
