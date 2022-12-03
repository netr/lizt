package lizt_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"git.faze.center/netr/lizt"
)

func TestNewStreamIterator(t *testing.T) {
	stream, err := lizt.NewStreamIterator(filenameOneMillion, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}
	tests := []struct {
		want     *lizt.StreamIterator
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "TestNewStreamIterator_SetsLineCount_Correctly",
			filename: filenameOneMillion,
			want:     stream,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Parallel()
		t.Run(tt.name, func(t *testing.T) {
			got, err := lizt.NewStreamIterator(tt.filename, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStreamIterator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Len(), tt.want.Len()) {
				t.Errorf("NewStreamIterator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStreamIterator_ShouldAddBothIteratorsCorrectly(t *testing.T) {
	m := lizt.NewManager()
	f, err := lizt.ReadFromFile(filenameOneMillion)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = m.AddIter(lizt.NewSliceIterator(nameNumbers, f, false))
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	fs, err := lizt.NewStreamIterator(filenameOneMillion, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = m.AddIter(fs)
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	if m.Len() != 2 {
		t.Errorf("wanted %d files, got: %d", 2, m.Len())
	}
}

func TestStreamIterator_Next(t *testing.T) {
	mgr := lizt.NewManager()
	fs, err := lizt.NewStreamIterator(filenameOneMillion, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = mgr.AddIter(fs)
	if err != nil {
		t.Errorf("AddIter() error = %v", err)
	}

	first, err := mgr.MustGet("1000000").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	second, err := mgr.MustGet("1000000").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	if reflect.DeepEqual(first, second) {
		t.Errorf("StreamIterator: expected next `%v` to be different", first)
	}

	var expected uint64 = 20
	if fs.Pointer() != expected {
		t.Errorf("expected pointer to be %d, got %d", expected, fs.Pointer())
	}
}

func TestStreamIterator_Next_RoundRobin(t *testing.T) {
	m := lizt.NewManager()
	fs, err := lizt.NewStreamIterator(filenameTen, true)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = m.AddIter(fs)
	if err != nil {
		t.Errorf("AddIter() error = %v", err)
	}

	first, err := m.MustGet("10").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	second, err := m.MustGet("10").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Errorf("expected %s to be %s", strings.Join(first, ","), strings.Join(second, ","))
	}

	var expected uint64 = 10
	if fs.Pointer() != expected {
		t.Errorf("expected pointer to be %d, got %d", expected, fs.Pointer())
	}
}

func TestStreamIterator_Next_RoundRobin_NoMoreLines(t *testing.T) {
	mgr := lizt.NewManager()
	fs, err := lizt.NewStreamIterator(filenameTen, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = mgr.AddIter(fs)
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	_, err = mgr.MustGet("10").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	_, err = mgr.MustGet("10").Next(10)
	if !errors.Is(err, lizt.ErrNoMoreLines) {
		t.Errorf("wanted ErrNoMoreLines, got error = %v", err)
	}

	var expected uint64 = 10
	if fs.Pointer() != expected {
		t.Errorf("expected pointer to be %d, got %d", expected, fs.Pointer())
	}
}
