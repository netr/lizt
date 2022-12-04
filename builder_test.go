package lizt_test

import (
	"git.faze.center/netr/lizt"
	"reflect"
	"testing"
)

func Test_NewIteratorBuilder(t *testing.T) {
	si, err := lizt.Builder().
		Stream("test/10.txt", true).
		WithSeeds(2, []string{"seed1", "seed2", "seed3"}).Build()
	if err != nil {
		t.Errorf("Builder() error = %v", err)
	}

	next, err := si.Next(5)
	if err != nil {
		t.Errorf("Builder.Next() error = %v", err)
	}

	if len(next) != 5 {
		t.Errorf("Builder.Next() len = %v, want %v", len(next), 5)
	}

	expected := []string{"seed1", "a", "seed2", "b", "seed3"}
	if !reflect.DeepEqual(next, expected) {
		t.Errorf("Builder.Next() = %v, want %v", next, expected)
	}
}
