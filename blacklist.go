package lizt

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// BlacklistingIterator is an iterator that skips blacklists while iterating.
type BlacklistingIterator struct {
	Blacklister
	PointerIterator
	blacklist *BlacklistManager
}

// BlacklistingIteratorConfig is the config for a blacklisting iterator.
type BlacklistingIteratorConfig struct {
	PointerIter PointerIterator
	Blacklisted *BlacklistManager
}

// NewBlacklistingIterator returns a new persistent iterator. It will set the pointer to the last known pointer.
func NewBlacklistingIterator(cfg BlacklistingIteratorConfig) (*BlacklistingIterator, error) {
	blkIter := &BlacklistingIterator{
		PointerIterator: cfg.PointerIter,
		blacklist:       cfg.Blacklisted,
	}

	return blkIter, nil
}

// Next returns the next line from the iterator.
func (bi *BlacklistingIterator) Next(count int) ([]string, error) {
	var clean []string
	for len(clean) < count {
		next, err := bi.PointerIterator.Next(count - len(clean))
		if err != nil {
			return nil, fmt.Errorf("next: name: %s -> %w", bi.Name(), err)
		}

		for _, n := range next {
			if !bi.IsBlacklisted(n) {
				clean = append(clean, n)
			}
		}
	}
	return clean, nil
}

// IsBlacklisted returns true if the given line is blacklisted.
func (bi *BlacklistingIterator) IsBlacklisted(line string) bool {
	return bi.blacklist.Has(line)
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

// ScrubFileWithBlacklist iterates over every line in a file and saves to a new file with the blacklisted lines removed.
func ScrubFileWithBlacklist(blkMap map[string]struct{}, sourcePath, destPath string) (n int, err error) {
	// Read from source file
	source, err := ReadFromFile(sourcePath)
	if err != nil {
		return 0, fmt.Errorf("read source: %w", err)
	}

	// Write to dest file
	dest, err := os.Create(destPath)
	if err != nil {
		return 0, fmt.Errorf("create dest: %w", err)
	}

	sb := strings.Builder{}
	for _, line := range source {
		if _, ok := blkMap[line]; !ok {
			sb.WriteString(line + "\n")
		} else {
			n++
		}
	}

	_, err = dest.WriteString(sb.String())
	if err != nil {
		return 0, fmt.Errorf("write line: %w", err)
	}

	return n, nil
}

type BlacklistManager struct {
	mu    sync.Mutex
	items map[string]struct{}
}

func NewBlacklistManager(items map[string]struct{}) *BlacklistManager {
	return &BlacklistManager{
		items: items,
	}
}

// Has returns true if the given string is in the list
func (l *BlacklistManager) Has(who string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, ok := l.items[who]
	return ok
}

// Map returns the map of items in the list
func (l *BlacklistManager) Map() map[string]struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.items
}

// Add adds a string to the list and appends to a file at the given path
func (l *BlacklistManager) Add(who string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.items[who]; ok {
		return fmt.Errorf("already in list")
	}

	l.items[who] = struct{}{}
	return nil
}

// ToStringSlice returns a slice of strings from the list
func (l *BlacklistManager) ToStringSlice() []string {
	l.mu.Lock()
	defer l.mu.Unlock()

	var items []string
	for item := range l.items {
		items = append(items, item)
	}

	return items
}
