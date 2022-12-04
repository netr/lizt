package lizt_test

import (
	"reflect"
	"testing"

	"git.faze.center/netr/lizt"
)

func Test_NewIteratorBuilder_Slice_WithSeeds(t *testing.T) {
	si, err := lizt.B().
		SliceRR([]string{"a", "b", "c", "d", "e"}).
		BuildWithSeeds(2, []string{"seed1", "seed2", "seed3"})
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

func Test_NewIteratorBuilder_Slice_Build(t *testing.T) {
	si, err := lizt.B().
		SliceRR([]string{"a", "b", "c", "d", "e"}).
		Build()
	if err != nil {
		t.Errorf("Builder() error = %v", err)
	}

	next, err := si.Next(3)
	if err != nil {
		t.Errorf("Builder.Next() error = %v", err)
	}

	if len(next) != 3 {
		t.Errorf("Builder.Next() len = %v, want %v", len(next), 3)
	}

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(next, expected) {
		t.Errorf("Builder.Next() = %v, want %v", next, expected)
	}
}

func Test_NewIteratorBuilder_Stream_WithSeeds(t *testing.T) {
	si, err := lizt.B().
		StreamRR("test/10.txt").
		BuildWithSeeds(2, []string{"seed1", "seed2", "seed3"})
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

func Test_NewIteratorBuilder_Stream_Build(t *testing.T) {
	si, err := lizt.B().
		StreamRR("test/10.txt").
		Build()
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

	expected := []string{"a", "b", "c", "d", "e"}
	if !reflect.DeepEqual(next, expected) {
		t.Errorf("Builder.Next() = %v, want %v", next, expected)
	}
}
