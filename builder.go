package lizt

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
)

var IterKeySeeds = "seeds"

type PointerIteratorBuilder struct {
	path      string
	seedIter  *SeedingIterator
	listIter  PointerIterator
	persister Persister
}

func NewBuilder() *PointerIteratorBuilder {
	return &PointerIteratorBuilder{}
}

func B() *PointerIteratorBuilder {
	return NewBuilder()
}

func (ib *PointerIteratorBuilder) Stream(path string, roundRobin bool) *PointerIteratorBuilder {
	stream, err := NewStreamIterator(path, roundRobin)
	if err != nil {
		panic(err)
	}
	ib.listIter = stream
	return ib
}

// Slice creates a new SliceIterator. Note that this randomizes the name and won't work while using a Manager. Use SliceWithName instead.
func (ib *PointerIteratorBuilder) Slice(lines []string, roundRobin bool) *PointerIteratorBuilder {
	ib.listIter = NewSliceIterator(randomString(8), lines, roundRobin)
	return ib
}

// SliceWithName creates a new SliceIterator with a name.
func (ib *PointerIteratorBuilder) SliceWithName(name string, lines []string, roundRobin bool) *PointerIteratorBuilder {
	ib.listIter = NewSliceIterator(name, lines, roundRobin)
	return ib
}

var (
	ErrNoIterator      = errors.New("no iterator")
	ErrInvalidSeedType = errors.New("invalid seed type")
)

func (ib *PointerIteratorBuilder) BuildWithSeeds(every int, seeds interface{}) (*SeedingIterator, error) {
	if ib.listIter == nil {
		return nil, fmt.Errorf("builder: %w", ErrNoIterator)
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

func (ib *PointerIteratorBuilder) Build() (PointerIterator, error) {
	if ib.seedIter != nil {
		ib.seedIter.PointerIterator = ib.listIter
		return ib.seedIter, nil
	}

	return ib.listIter, nil
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
