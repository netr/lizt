package lizt

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"reflect"
)

var IterKeySeeds = "seeds"

type PointerIteratorBuilder struct {
	seedIter      *SeedingIterator
	blacklistIter *BlacklistingIterator
	listIter      PointerIterator
}

func NewBuilder() *PointerIteratorBuilder {
	return &PointerIteratorBuilder{}
}

func B() *PointerIteratorBuilder {
	return NewBuilder()
}

// Stream creates a new StreamIterator.
func (ib *PointerIteratorBuilder) Stream(path string) *PointerIteratorBuilder {
	stream, err := NewStreamIterator(path, false)
	if err != nil {
		panic(err)
	}
	ib.listIter = stream
	return ib
}

// StreamRR creates a new StreamIterator with round-robin.
func (ib *PointerIteratorBuilder) StreamRR(path string) *PointerIteratorBuilder {
	stream, err := NewStreamIterator(path, true)
	if err != nil {
		panic(err)
	}
	ib.listIter = stream
	return ib
}

// Slice creates a new SliceIterator. Note that this randomizes the name and won't work while using a Manager. Use SliceNamed instead.
func (ib *PointerIteratorBuilder) Slice(lines []string) *PointerIteratorBuilder {
	ib.listIter = NewSliceIterator(randomString(8), lines, false)
	return ib
}

// SliceRR creates a new SliceIterator with round-robin. Note that this randomizes the name and won't work while using a Manager. Use SliceNamed instead.
func (ib *PointerIteratorBuilder) SliceRR(lines []string) *PointerIteratorBuilder {
	ib.listIter = NewSliceIterator(randomString(8), lines, true)
	return ib
}

// SliceNamed creates a new SliceIterator with a name.
func (ib *PointerIteratorBuilder) SliceNamed(name string, lines []string, roundRobin bool) *PointerIteratorBuilder {
	ib.listIter = NewSliceIterator(name, lines, roundRobin)
	return ib
}

// SliceNamedRR creates a new SliceIterator with a name and round-robin.
func (ib *PointerIteratorBuilder) SliceNamedRR(name string, lines []string) *PointerIteratorBuilder {
	ib.listIter = NewSliceIterator(name, lines, true)
	return ib
}

// Blacklist creates a new BlacklistingIterator
func (ib *PointerIteratorBuilder) Blacklist(bl *BlacklistManager) *PointerIteratorBuilder {
	var err error
	ib.blacklistIter, err = NewBlacklistingIterator(BlacklistingIteratorConfig{
		PointerIter: ib.listIter,
		Blacklisted: bl,
	})
	if err != nil {
		panic(err)
	}
	return ib
}

var (
	ErrNoIterator      = errors.New("no iterator")
	ErrInvalidSeedType = errors.New("invalid seed type")
)

// BuildWithSeeds will build a pointer iterator with the given iterator and seeds.
func (ib *PointerIteratorBuilder) BuildWithSeeds(every int, seeds interface{}) (*SeedingIterator, error) {
	if ib.listIter == nil {
		return nil, fmt.Errorf("builder: %w", ErrNoIterator)
	}

	if ib.blacklistIter != nil {
		ib.blacklistIter.PointerIterator = ib.listIter
		ib.listIter = ib.blacklistIter
	}

	switch reflect.TypeOf(seeds).Kind() {
	case reflect.Slice:
		return NewSeedingIterator(SeedingIteratorConfig{
			PointerIter: ib.listIter,
			SeedIter:    NewSliceIterator(IterKeySeeds, seeds.([]string), true),
			PlantEvery:  every,
		}), nil
	case reflect.String:
		stream, err := NewStreamIterator(fmt.Sprintf("%s", seeds), true)
		if err != nil {
			panic(err)
		}
		return NewSeedingIterator(SeedingIteratorConfig{
			PointerIter: ib.listIter,
			SeedIter:    stream,
			PlantEvery:  every,
		}), nil
	}

	return nil, fmt.Errorf("builder: %w", ErrInvalidSeedType)
}

// MustBuildWithSeeds will build a pointer iterator with the given iterator and seeds. Panics.
func (ib *PointerIteratorBuilder) MustBuildWithSeeds(every int, seeds interface{}) *SeedingIterator {
	iter, err := ib.BuildWithSeeds(every, seeds)
	if err != nil {
		panic(err)
	}
	return iter
}

// Build will build a pointer iterator with the given iterators.
func (ib *PointerIteratorBuilder) Build() (PointerIterator, error) {
	if ib.blacklistIter != nil {
		ib.blacklistIter.PointerIterator = ib.listIter
		ib.listIter = ib.blacklistIter
	}

	if ib.seedIter != nil {
		ib.seedIter.PointerIterator = ib.listIter
		return ib.seedIter, nil
	}

	return ib.listIter, nil
}

// MustBuild will build a persistent iterator with the given persister. Panics.
func (ib *PointerIteratorBuilder) MustBuild() PointerIterator {
	iter, err := ib.Build()
	if err != nil {
		panic(err)
	}
	return iter
}

// PersistentIteratorBuilder is a builder for a PersistentIterator.
type PersistentIteratorBuilder struct {
	*PointerIteratorBuilder
	persister Persister
}

// PersistTo creates a new PersistentIteratorBuilder.
func (ib *PointerIteratorBuilder) PersistTo(p Persister) *PersistentIteratorBuilder {
	if ib.listIter == nil {
		log.Fatal("builder: no iterator")
	}

	pib := &PersistentIteratorBuilder{
		PointerIteratorBuilder: ib,
		persister:              p,
	}
	return pib
}

// Blacklist creates a new BlacklistingIterator
func (ib *PersistentIteratorBuilder) Blacklist(bl *BlacklistManager) *PersistentIteratorBuilder {
	var err error
	ib.blacklistIter, err = NewBlacklistingIterator(BlacklistingIteratorConfig{
		PointerIter: ib.listIter,
		Blacklisted: bl,
	})
	if err != nil {
		panic(err)
	}
	return ib
}

// BuildWithSeeds will build a persistent iterator with the given persister and seeds.
func (ib *PersistentIteratorBuilder) BuildWithSeeds(every int, seeds interface{}) (*PersistentIterator, error) {
	if ib.listIter == nil {
		return nil, fmt.Errorf("builder: %w", ErrNoIterator)
	}

	if ib.blacklistIter != nil {
		ib.blacklistIter.PointerIterator = ib.listIter
		ib.listIter = ib.blacklistIter
	}

	switch reflect.TypeOf(seeds).Kind() {
	case reflect.Slice:
		per, err := NewPersistentIterator(PersistentIteratorConfig{
			PointerIter: NewSeedingIterator(SeedingIteratorConfig{
				PointerIter: ib.listIter,
				SeedIter:    NewSliceIterator(IterKeySeeds, seeds.([]string), true),
				PlantEvery:  every,
			}),
			Persister: ib.persister,
		})
		if err != nil {
			return nil, fmt.Errorf("builder: %w", err)
		}

		return per, nil
	case reflect.String:
		stream, err := NewStreamIterator(fmt.Sprintf("%s", seeds), true)
		if err != nil {
			panic(err)
		}

		per, err := NewPersistentIterator(PersistentIteratorConfig{
			PointerIter: NewSeedingIterator(SeedingIteratorConfig{
				PointerIter: ib.listIter,
				SeedIter:    stream,
				PlantEvery:  every,
			}),
			Persister: ib.persister,
		})
		if err != nil {
			return nil, fmt.Errorf("builder: %w", err)
		}

		return per, nil
	}

	return nil, fmt.Errorf("builder: %w", ErrInvalidSeedType)
}

// MustBuildWithSeeds will build a persistent iterator with the given persister and seeds. Panics.
func (ib *PersistentIteratorBuilder) MustBuildWithSeeds(every int, seeds interface{}) *PersistentIterator {
	iter, err := ib.BuildWithSeeds(every, seeds)
	if err != nil {
		panic(err)
	}
	return iter
}

// Build will build a persistent iterato with the given persister.
func (ib *PersistentIteratorBuilder) Build() (*PersistentIterator, error) {
	if ib.blacklistIter != nil {
		ib.blacklistIter.PointerIterator = ib.listIter
		ib.listIter = ib.blacklistIter
	}

	if ib.seedIter != nil {
		ib.seedIter.PointerIterator = ib.listIter
		per, err := NewPersistentIterator(PersistentIteratorConfig{
			PointerIter: ib.seedIter,
			Persister:   ib.persister,
		})
		if err != nil {
			return nil, fmt.Errorf("builder: %w", err)
		}

		return per, nil
	}

	per, err := NewPersistentIterator(PersistentIteratorConfig{
		PointerIter: ib.listIter,
		Persister:   ib.persister,
	})
	if err != nil {
		return nil, fmt.Errorf("builder: %w", err)
	}

	return per, nil
}

// MustBuild will build a persistent iterator with the given persister. Panics.
func (ib *PersistentIteratorBuilder) MustBuild() *PersistentIterator {
	iter, err := ib.Build()
	if err != nil {
		panic(err)
	}
	return iter
}

// randomString ty copilot.
func randomString(count int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, count)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
