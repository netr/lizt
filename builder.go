package lizt

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	IterKeySeeds = "seeds"
)

type PointerIteratorBuilder struct {
	path     string
	seedIter *SeedingIterator
	listIter PointerIterator
}

func NewBuilder() *PointerIteratorBuilder {
	return &PointerIteratorBuilder{}
}

func Builder() *PointerIteratorBuilder {
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

func (ib *PointerIteratorBuilder) Slice(name string, lines []string, roundRobin bool) *PointerIteratorBuilder {
	ib.listIter = NewSliceIterator(name, lines, roundRobin)
	return ib
}

func (ib *PointerIteratorBuilder) WithSeeds(every int, seeds interface{}) (*SeedingIterator, error) {
	if ib.listIter == nil {
		panic("no iterator")
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
	return nil, errors.New("invalid seeds type")
}

func (ib *PointerIteratorBuilder) Build() (PointerIterator, error) {
	if ib.seedIter != nil {
		ib.seedIter.PointerIterator = ib.listIter
		return ib.seedIter, nil
	}

	return ib.listIter, nil
}
