package lizt

import "sync/atomic"

// SeedingIterator is an iterator that reads from a slice
type SeedingIterator struct {
	Seeder
	PointerIterator
	seeds       *SliceIterator
	totalSeeded *atomic.Int64
	seedAfter   int
}

// NewSeedingIterator returns a new slice iterator
func NewSeedingIterator(p PointerIterator, seeds *SliceIterator, seedAfter int) *SeedingIterator {
	return &SeedingIterator{
		PointerIterator: p,
		seeds:           seeds,
		seedAfter:       seedAfter,
		totalSeeded:     new(atomic.Int64),
	}
}

// Seeds returns the seeds
func (si *SeedingIterator) Seeds() []string {
	return si.seeds.lines
}

// Planted returns how many seeds have been planted
func (si *SeedingIterator) Planted() int64 {
	return si.totalSeeded.Load()
}

// PlantEvery returns how often the seeds are planted
func (si *SeedingIterator) PlantEvery() int {
	return si.seedAfter
}

// Next returns the next line from the iterator and will automatically seed every PlantEvery() lines
func (si *SeedingIterator) Next(count int) ([]string, error) {
	var lines []string
	for i := 0; i < count; i++ {
		sent := si.Pointer()
		if sent > 0 && sent%uint64(si.PlantEvery()) == 0 {
			seed, err := si.seeds.Next(1)
			if err != nil {
				return nil, err
			}
			lines = append(lines, seed[0])
		} else {
			next, err := si.PointerIterator.Next(1)
			if err != nil {
				return nil, err
			}
			si.Inc()
			lines = append(lines, next...)
		}
	}

	return lines, nil

}
