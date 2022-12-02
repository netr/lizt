package lizt

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

var (
	FilenameOneMillion = "test/1000000.txt"
	FilenameTen        = "test/10.txt"
	NameNumbers        = "numbers"
)

func TestNewStreamIterator(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		want     *StreamIterator
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "TestNewStreamIterator_SetsLineCount_Correctly",
			filename: FilenameOneMillion,
			want: &StreamIterator{
				fileLines: 999998,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStreamIterator(tt.filename, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStreamIterator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.fileLines, tt.want.fileLines) {
				t.Errorf("NewStreamIterator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStreamIterator_ShouldAddBothIteratorsCorrectly(t *testing.T) {
	m := NewManager()
	f, err := ReadFromFile(FilenameOneMillion)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	err = m.AddIter(NewSliceIterator(NameNumbers, f, false))
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	fs, err := NewStreamIterator(FilenameOneMillion, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = m.AddIter(fs)
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	if len(m.files) != 2 {
		t.Errorf("wanted %d files, got: %d", 2, len(m.files))
	}
}

func TestStreamIterator_Next(t *testing.T) {
	m := NewManager()
	fs, err := NewStreamIterator(FilenameOneMillion, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = m.AddIter(fs)
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	first, err := m.Get("1000000").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	second, err := m.Get("1000000").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	if reflect.DeepEqual(first, second) {
		t.Errorf("StreamIterator: expected next `%v` to be different", first)
	}
}

func TestStreamIterator_Next_RoundRobin(t *testing.T) {
	m := NewManager()
	fs, err := NewStreamIterator(FilenameTen, true)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = m.AddIter(fs)
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	first, err := m.Get("10").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	second, err := m.Get("10").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Errorf("expected %s to be %s", strings.Join(first, ","), strings.Join(second, ","))
	}
}

func TestStreamIterator_Next_RoundRobin_NoMoreLines(t *testing.T) {
	m := NewManager()
	fs, err := NewStreamIterator(FilenameTen, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	err = m.AddIter(fs)
	if err != nil {
		t.Errorf("AddPointerIter() error = %v", err)
	}

	_, err = m.Get("10").Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	_, err = m.Get("10").Next(10)
	if !errors.Is(err, ErrNoMoreLines) {
		t.Errorf("wanted ErrNoMoreLines, got error = %v", err)
	}
}
