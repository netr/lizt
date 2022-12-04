package lizt_test

import (
	"errors"
	"reflect"
	"testing"

	"git.faze.center/netr/lizt"
)

func TestPersistentIterator_Next(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	iter := lizt.NewSliceIterator(nameNumbers, numbers, false)

	mem := NewInMemoryPersister()

	p, err := lizt.NewPersistentIterator(
		lizt.PersistentIteratorConfig{
			PointerIter: iter,
			Persister:   mem,
		},
	)
	if err != nil {
		t.Errorf("NewPersistentIterator expected no error, got %v", err)
	}

	next, err := p.Next(5)
	if err != nil {
		t.Errorf("Next expected no error, got %v", err)
	}
	if !reflect.DeepEqual(next, numbers[:5]) {
		t.Errorf("Expected %v, got %v", numbers[:5], next)
	}

	if mem.pointers[nameNumbers] != 5 {
		t.Errorf("Expected %d, got %d", 5, mem.pointers[nameNumbers])
	}
}

func TestPersistentIterator_Slice_Next_ShouldBeOffsetBy_Two_AfterPersistenceInit(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	iter := lizt.NewSliceIterator(nameNumbers, numbers, false)

	mem := NewInMemoryPersister()
	mem.pointers[nameNumbers] = 2

	p, err := lizt.NewPersistentIterator(
		lizt.PersistentIteratorConfig{
			PointerIter: iter,
			Persister:   mem,
		},
	)
	if err != nil {
		t.Errorf("NewPersistentIterator expected no error, got %v", err)
	}

	next, err := p.Next(5)
	if err != nil {
		t.Errorf("Next expected no error, got %v", err)
	}
	if !reflect.DeepEqual(next, numbers[2:7]) {
		t.Errorf("Expected %v, got %v", numbers[2:7], next)
	}

	if mem.pointers[nameNumbers] != 7 {
		t.Errorf("Expected %d, got %d", 7, mem.pointers[nameNumbers])
	}
}

func TestPersistentIterator_Stream_Next_ShouldBeOffsetBy_Two_AfterPersistenceInit(t *testing.T) {
	iter, err := lizt.NewStreamIterator(filenameTen, false)
	if err != nil {
		t.Errorf("NewStreamIterator expected no error, got %v", err)
	}

	iterKey := "10"
	mem := NewInMemoryPersister()
	mem.pointers[iterKey] = 2

	p, err := lizt.NewPersistentIterator(
		lizt.PersistentIteratorConfig{
			PointerIter: iter,
			Persister:   mem,
		},
	)
	if err != nil {
		t.Errorf("NewPersistentIterator expected no error, got %v", err)
	}

	next, err := p.Next(5)
	if err != nil {
		t.Errorf("Next expected no error, got %v", err)
	}

	expected := []string{"c", "d", "e", "f", "g"}
	if !reflect.DeepEqual(next, expected) {
		t.Errorf("Expected %v, got %v", expected, next)
	}

	if mem.pointers[iterKey] != 7 {
		t.Errorf("Expected %d, got %d", 7, mem.pointers[iterKey])
	}
}

func TestNewPersistentIterator_UsingSeedingIterator_Next(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	numberIter := lizt.NewSliceIterator(nameNumbers, numbers, false)

	seeds := []string{"seeder1", "seeder2", "seeder3", "seeder4", "seeder5", "seeder6", "seeder7", "seeder8", "seeder9", "seeder10"}
	seedSliceIter := lizt.NewSliceIterator("seedIter", seeds, true)

	seedIter := lizt.NewSeedingIterator(
		lizt.SeedingIteratorConfig{
			PointerIter: numberIter,
			SeedIter:    seedSliceIter,
			PlantEvery:  2,
		},
	)

	mem := NewInMemoryPersister()
	persistIter, err := lizt.NewPersistentIterator(
		lizt.PersistentIteratorConfig{
			PointerIter: seedIter,
			Persister:   mem,
		},
	)
	if err != nil {
		t.Errorf("NewPersistentIterator expected no error, got %v", err)
	}

	next, err := persistIter.Next(6)
	if err != nil {
		t.Errorf("Next expected no error, got %v", err)
	}

	if !reflect.DeepEqual(next, []string{"seeder1", "1", "seeder2", "2", "seeder3", "3"}) {
		t.Errorf("Expected %v, got %v", []string{"seeder1", "1", "seeder2", "2", "seeder3", "3"}, next)
	}

	// Only 3 iterator items were pulled, the other 3 are seeds. So the pointer should be 3, not 6.
	if val, err := mem.Get(nameNumbers); err != nil || val != 3 {
		t.Errorf("Expected %d, got %d", 3, val)
	}
}

// USED AS "MOCK" FOR TESTING.
type InMemoryPersister struct {
	lizt.PersistentIterator
	pointers map[string]uint64
}

func NewInMemoryPersister() *InMemoryPersister {
	return &InMemoryPersister{
		pointers: make(map[string]uint64, 0),
	}
}

func (i *InMemoryPersister) Set(key string, value uint64) error {
	i.pointers[key] = value
	return nil
}

var ErrNotFound = errors.New("not found")

func (i *InMemoryPersister) Get(key string) (uint64, error) {
	if val, ok := i.pointers[key]; ok {
		return val, nil
	}

	return 0, ErrNotFound
}
