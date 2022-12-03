package lizt

import (
	"fmt"
	"sync/atomic"
)

// SeedingIterator is an iterator that reads from a slice.
type SeedingIterator struct {
	Seeder
	PointerIterator
	seeds       *SliceIterator
	totalSeeded *atomic.Int64
	seedEvery   int
}

type SeedingIteratorConfig struct {
	PointerIterator PointerIterator
	Seeds           *SliceIterator
	SeedEvery       int
}

// NewSeedingIterator returns a new slice iterator.
func NewSeedingIterator(cfg SeedingIteratorConfig) *SeedingIterator {
	return &SeedingIterator{
		PointerIterator: cfg.PointerIterator,
		seeds:           cfg.Seeds,
		seedEvery:       cfg.SeedEvery,
		totalSeeded:     new(atomic.Int64),
	}
}

// Seeds returns the seeds.
func (si *SeedingIterator) Seeds() []string {
	return si.seeds.lines
}

// Planted returns how many seeds have been planted.
func (si *SeedingIterator) Planted() int64 {
	return si.totalSeeded.Load()
}

// PlantEvery returns how often the seeds are planted.
func (si *SeedingIterator) PlantEvery() int {
	return si.seedEvery
}

func (si *SeedingIterator) IncPlanted() {
	si.totalSeeded.Add(1)
}

// Next returns the next line from the iterator and will automatically seed every PlantEvery() lines.
func (si *SeedingIterator) Next(count int) ([]string, error) {
	var lines []string
	for i := 0; i < count; i++ {
		sent := si.Pointer() + uint64(si.Planted())
		if sent%uint64(si.PlantEvery()) == 0 {
			seed, err := si.seeds.Next(1)
			if err != nil {
				return nil, err
			}
			si.IncPlanted()
			lines = append(lines, seed[0])
		} else {
			next, err := si.PointerIterator.Next(1)
			if err != nil {
				return nil, fmt.Errorf("file: %s -> %w", si.Name(), err)
			}
			lines = append(lines, next...)
		}
	}

	return lines, nil
}
