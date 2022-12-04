package lizt_test

import (
	"errors"
	"reflect"
	"sort"
	"strings"
	"sync"
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

	m.AddIter(lizt.NewSliceIterator(nameNumbers, f, false))

	fs, err := lizt.NewStreamIterator(filenameOneMillion, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}

	m.AddIter(fs)

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

	mgr.AddIter(fs)

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

	m.AddIter(fs)

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

	mgr.AddIter(fs)

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

func TestStreamIterator_Next_Parallel(t *testing.T) {
	b, err := lizt.B().StreamRR("test/10.txt").Build()
	if err != nil {
		t.Errorf("Build() error = %v", err)
	}

	res := make(chan string, b.Len()*3)
	wg := &sync.WaitGroup{}
	wg.Add(30)
	for i := 0; i < b.Len()*3; i++ {
		go func(w *sync.WaitGroup, ch chan string) {
			x, err := b.Next(1)
			if err != nil {
				t.Errorf("Next() error = %v", err)
			}
			ch <- x[0]
			w.Done()
		}(wg, res)
	}
	wg.Wait()
	close(res)

	var results []string
	for msg := range res {
		results = append(results, msg)
	}
	sort.Strings(results)
	expected := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	sort.Strings(expected)

	if !reflect.DeepEqual(results, expected) {
		t.Errorf("expected %v, got %v", expected, results)
	}
}
