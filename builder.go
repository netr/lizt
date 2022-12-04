package lizt

import (
	"fmt"
	"reflect"
)

var (
	IterKeySeeds = "seeds"
)

type IteratorBuilder struct{}

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

func (ib *PointerIteratorBuilder) Build() (PointerIterator, error) {
	if ib.seedIter != nil {
		ib.seedIter.PointerIterator = ib.listIter
		return ib.seedIter, nil
	}

	return ib.listIter, nil
}

type SeedingIteratorBuilder struct {
	path     string
	seedIter *SeedingIterator
	listIter PointerIterator
}

func (ib *PointerIteratorBuilder) WithSeeds(every int, seeds interface{}) *SeedingIteratorBuilder {
	if ib.listIter == nil {
		panic("no iterator")
	}

	si := &SeedingIteratorBuilder{
		path:     ib.path,
		listIter: ib.listIter,
	}
	switch reflect.TypeOf(seeds).Kind() {
	case reflect.Slice:
		si.seedIter = NewSeedingIterator(SeedingIteratorConfig{
			PointerIter: nil,
			SeedIter:    NewSliceIterator(IterKeySeeds, seeds.([]string), true),
			PlantEvery:  every,
		})
		break
	case reflect.String:
		stream, err := NewStreamIterator(fmt.Sprintf("%s", seeds), true)
		if err != nil {
			panic(err)
		}
		si.seedIter = NewSeedingIterator(SeedingIteratorConfig{
			PointerIter: nil,
			SeedIter:    stream,
			PlantEvery:  every,
		})
	}
	return si
}

func (sb *SeedingIteratorBuilder) Build() (*SeedingIterator, error) {
	sb.seedIter.PointerIterator = sb.listIter
	return sb.seedIter, nil
}
