package lizt_test

import (
	"testing"

	"git.faze.center/netr/lizt"
)

func TestSeeder(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	seeds := []string{"seeder1", "seeder2", "seeder3", "seeder4", "seeder5", "seeder6", "seeder7", "seeder8", "seeder9", "seeder10"}
	seed := lizt.NewSeedingIterator(
		lizt.NewSliceIterator(nameNumbers, numbers, false),
		lizt.NewSliceIterator("seeds", seeds, true), 2,
	)
	next, err := seed.Next(2)
	if err != nil {
		return
	}
	if next[0] != "1" {
		t.Errorf("expected 1, got %s", next[0])
	}
	if next[1] != "seeder1" {
		t.Errorf("expected seeder, got %s", next[1])
	}
}
