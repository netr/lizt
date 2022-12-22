package lizt_test

import (
	"testing"

	"git.faze.center/netr/lizt"
)

func TestSeeder_Next(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	iter := lizt.NewSliceIterator(nameNumbers, numbers, false)

	seeds := []string{"seeder1", "seeder2", "seeder3", "seeder4", "seeder5", "seeder6", "seeder7", "seeder8", "seeder9", "seeder10"}
	seedIter := lizt.NewSliceIterator("seedIter", seeds, true)

	seed := lizt.NewSeedingIterator(
		lizt.SeedingIteratorConfig{
			PointerIter: iter,
			SeedIter:    seedIter,
			PlantEvery:  2,
		},
	)
	next, err := seed.Next(6)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := map[int]string{
		0: "seeder1",
		1: "1",
		2: "seeder2",
		3: "2",
		4: "seeder3",
		5: "3",
	}

	for k, v := range expected {
		if next[k] != v {
			t.Errorf("Expected %s, got %s", v, next[k])
		}
	}
}

func TestSeeder_Next_RoundRobin(t *testing.T) {
	numbers := []string{"1", "2"}
	iter := lizt.NewSliceIterator(nameNumbers, numbers, true)

	seeds := []string{"seeder1", "seeder2"}
	seedIter := lizt.NewSliceIterator("seeds", seeds, true)

	seed := lizt.NewSeedingIterator(
		lizt.SeedingIteratorConfig{
			PointerIter: iter,
			SeedIter:    seedIter,
			PlantEvery:  2,
		},
	)
	next, err := seed.Next(8)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := map[int]string{
		0: "seeder1",
		1: "1",
		2: "seeder2",
		3: "2",
		4: "seeder1",
		5: "1",
		6: "seeder2",
		7: "2",
	}

	for k, v := range expected {
		if next[k] != v {
			t.Errorf("Expected %s, got %s", v, next[k])
		}
	}
}

func TestSeeder_NextSeed(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	iter := lizt.NewSliceIterator(nameNumbers, numbers, false)

	seeds := []string{"seeder1", "seeder2", "seeder3", "seeder4", "seeder5", "seeder6", "seeder7", "seeder8", "seeder9", "seeder10"}
	seedIter := lizt.NewSliceIterator("seedIter", seeds, true)

	seed := lizt.NewSeedingIterator(
		lizt.SeedingIteratorConfig{
			PointerIter: iter,
			SeedIter:    seedIter,
			PlantEvery:  2,
		},
	)
	next, seeded, err := seed.NextSeed(6)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !seeded {
		t.Errorf("Expected seeded to be true")
	}

	expected := map[int]string{
		0: "seeder1",
		1: "1",
		2: "seeder2",
		3: "2",
		4: "seeder3",
		5: "3",
	}

	for k, v := range expected {
		if next[k] != v {
			t.Errorf("Expected %s, got %s", v, next[k])
		}
	}
}
