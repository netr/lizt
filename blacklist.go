package lizt

import "fmt"

// BlacklistingIterator is an iterator that skips blacklists while iterating.
type BlacklistingIterator struct {
	Blacklister
	PointerIterator
	blacklist map[string]struct{}
}

// BlacklistingIteratorConfig is the config for a blacklisting iterator.
type BlacklistingIteratorConfig struct {
	PointerIter PointerIterator
	Blacklisted []string
}

// NewBlacklistingIterator returns a new persistent iterator. It will set the pointer to the last known pointer.
func NewBlacklistingIterator(cfg BlacklistingIteratorConfig) (*BlacklistingIterator, error) {
	blkIter := &BlacklistingIterator{
		PointerIterator: cfg.PointerIter,
		blacklist:       make(map[string]struct{}, len(cfg.Blacklisted)),
	}

	for _, bl := range cfg.Blacklisted {
		blkIter.blacklist[bl] = struct{}{}
	}

	return blkIter, nil
}

// Next returns the next line from the iterator.
func (bi *BlacklistingIterator) Next(count int) ([]string, error) {
	var clean []string
	for {
		next, err := bi.PointerIterator.Next(count)
		if err != nil {
			return nil, fmt.Errorf("next: name: %s -> %w", bi.Name(), err)
		}

		for _, n := range next {
			if !bi.IsBlacklisted(n) {
				clean = append(clean, n)
			}
		}

		if len(clean) > 0 {
			break
		}
	}
	return clean, nil
}

// IsBlacklisted returns true if the given line is blacklisted.
func (bi *BlacklistingIterator) IsBlacklisted(line string) bool {
	_, ok := bi.blacklist[line]
	return ok
}

func removeSlice(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// MustNext returns the next lines, of a given count, from the iterator. Panics on error.
func (bi *BlacklistingIterator) MustNext(count int) []string {
	lines, err := bi.Next(count)
	if err != nil {
		panic(err)
	}
	return lines
}

// NextOne returns the next line from the iterator.
func (bi *BlacklistingIterator) NextOne() (string, error) {
	lines, err := bi.Next(1)
	if err != nil {
		return "", err
	}
	return lines[0], nil
}

// MustNextOne returns the next line from the iterator. Panics on error.
func (bi *BlacklistingIterator) MustNextOne() string {
	line, err := bi.NextOne()
	if err != nil {
		panic(err)
	}
	return line
}
