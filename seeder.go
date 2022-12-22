package lizt

import (
	"fmt"
	"sync/atomic"
)

// SeedingIterator is an iterator that reads from a slice.
type SeedingIterator struct {
	Seeder
	PointerIterator
	seedIter     PointerIterator
	totalPlanted *atomic.Int64
	plantEvery   int
}

type SeedingIteratorConfig struct {
	PointerIter PointerIterator
	SeedIter    PointerIterator
	PlantEvery  int
}

// NewSeedingIterator returns a new slice iterator.
func NewSeedingIterator(cfg SeedingIteratorConfig) *SeedingIterator {
	return &SeedingIterator{
		PointerIterator: cfg.PointerIter,
		seedIter:        cfg.SeedIter,
		plantEvery:      cfg.PlantEvery,
		totalPlanted:    new(atomic.Int64),
	}
}

// Planted returns how many seedIter have been planted.
func (si *SeedingIterator) Planted() int64 {
	return si.totalPlanted.Load()
}

// PlantEvery returns how often the seedIter are planted.
func (si *SeedingIterator) PlantEvery() int {
	return si.plantEvery
}

// inc increments the total planted counter.
func (si *SeedingIterator) inc() {
	si.totalPlanted.Add(1)
}

// Next returns the next line from the iterator and will automatically seed every PlantEvery() lines.
func (si *SeedingIterator) Next(count int) ([]string, error) {
	lines, _, err := si.nextSeed(count)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

// NextSeed returns the next line from the iterator and will automatically seed every PlantEvery() lines.
// The difference from the interface Next() is that this returns a bool indicating if a seed was planted.
func (si *SeedingIterator) NextSeed(count int) ([]string, bool, error) {
	return si.nextSeed(count)
}

func (si *SeedingIterator) nextSeed(count int) ([]string, bool, error) {
	var lines []string
	seeded := false
	for i := 0; i < count; i++ {
		sent := si.Pointer() + uint64(si.Planted())
		if sent%uint64(si.PlantEvery()) == 0 {
			seed, err := si.seedIter.Next(1)
			if err != nil {
				return nil, seeded, fmt.Errorf("seed iter next: %w", err)
			}
			seeded = true
			si.inc()
			lines = append(lines, seed[0])
		} else {
			next, err := si.PointerIterator.Next(1)
			if err != nil {
				if len(lines) == 0 {
					return nil, seeded, fmt.Errorf("file: %s -> %w", si.Name(), err)
				}
				return lines, seeded, nil
			}
			lines = append(lines, next[0])
		}
	}

	return lines, seeded, nil
}
