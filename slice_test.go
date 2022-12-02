package lizt

import (
	"errors"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
)

func TestNewSliceIterator(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		want *SliceIterator
		name string
		args args
	}{
		{
			name: "TestNewSliceIterator_SetsLines_Correctly",
			args: args{
				lines: []string{"a", "b", "c"},
			},
			want: &SliceIterator{
				lines:   []string{"a", "b", "c"},
				pointer: new(atomic.Uint64),
				name:    NameNumbers,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSliceIterator(NameNumbers, tt.args.lines, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSliceIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSliceIterator_Next(t *testing.T) {
	m := NewManager()
	f, err := ReadFromFile(FilenameOneMillion)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = m.AddIter(NewSliceIterator(NameNumbers, f, false))
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	first, err := m.Get(NameNumbers).Next(10)
	if err != nil {
		t.Errorf("SliceIterator.Next() error = %v", err)
	}
	second, err := m.Get(NameNumbers).Next(10)
	if err != nil {
		t.Errorf("SliceIterator.Next() error = %v", err)
	}

	if reflect.DeepEqual(first, second) {
		t.Errorf("SliceIterator: expected next `%v` to be different", first)
	}
}

func TestSliceIterator_Next_RoundRobin(t *testing.T) {
	m := NewManager()
	f, err := ReadFromFile(FilenameTen)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = m.AddIter(NewSliceIterator(NameNumbers, f, true))
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	first, err := m.Get(NameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	second, err := m.Get(NameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Errorf("expected %s to be %s", strings.Join(first, ","), strings.Join(second, ","))
	}
}

func TestSliceIterator_Next_RoundRobin_NoMoreLines(t *testing.T) {
	m := NewManager()
	f, err := ReadFromFile(FilenameTen)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = m.AddIter(NewSliceIterator(NameNumbers, f, false))
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	_, err = m.Get(NameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	_, err = m.Get(NameNumbers).Next(10)
	if !errors.Is(err, ErrNoMoreLines) {
		t.Errorf("wanted ErrNoMoreLines, got error = %v", err)
	}
}
