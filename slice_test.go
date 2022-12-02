package lizt_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"git.faze.center/netr/lizt"
)

func TestNewSliceIterator(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		want *lizt.SliceIterator
		name string
		args args
	}{
		{
			name: "TestNewSliceIterator_SetsLines_Correctly",
			args: args{
				lines: []string{"a", "b", "c"},
			},
			want: lizt.NewSliceIterator(nameNumbers, []string{"a", "b", "c"}, false),
		},
	}
	for _, tt := range tests {
		t.Parallel()
		t.Run(tt.name, func(t *testing.T) {
			if got := lizt.NewSliceIterator(nameNumbers, tt.args.lines, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSliceIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSliceIterator_Next(t *testing.T) {
	mgr := lizt.NewManager()
	f, err := lizt.ReadFromFile(filenameOneMillion)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = mgr.AddIter(lizt.NewSliceIterator(nameNumbers, f, false))
	if err != nil {
		t.Errorf("AddIter() error = %v", err)
	}

	first, err := mgr.Get(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("SliceIterator.Next() error = %v", err)
	}
	second, err := mgr.Get(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("SliceIterator.Next() error = %v", err)
	}

	if reflect.DeepEqual(first, second) {
		t.Errorf("SliceIterator: expected next `%v` to be different", first)
	}
}

func TestSliceIterator_Next_RoundRobin(t *testing.T) {
	mgr := lizt.NewManager()
	f, err := lizt.ReadFromFile(filenameTen)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = mgr.AddIter(lizt.NewSliceIterator(nameNumbers, f, true))
	if err != nil {
		t.Errorf("AddIter() error = %v", err)
	}

	first, err := mgr.Get(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	second, err := mgr.Get(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Errorf("expected %s to be %s", strings.Join(first, ","), strings.Join(second, ","))
	}
}

func TestSliceIterator_Next_RoundRobin_NoMoreLines(t *testing.T) {
	mgr := lizt.NewManager()
	f, err := lizt.ReadFromFile(filenameTen)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = mgr.AddIter(lizt.NewSliceIterator(nameNumbers, f, false))
	if err != nil {
		t.Errorf("NewSliceIterator() error = %v", err)
	}

	_, err = mgr.Get(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	_, err = mgr.Get(nameNumbers).Next(10)
	if !errors.Is(err, lizt.ErrNoMoreLines) {
		t.Errorf("wanted ErrNoMoreLines, got error = %v", err)
	}
}
