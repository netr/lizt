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

	si := lizt.NewSliceIterator(nameNumbers, f, false)

	mgr.AddIter(si)

	first, err := mgr.MustGet(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("SliceIterator.Next() error = %v", err)
	}
	second, err := mgr.MustGet(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("SliceIterator.Next() error = %v", err)
	}

	if reflect.DeepEqual(first, second) {
		t.Errorf("SliceIterator: expected next `%v` to be different", first)
	}

	var expected uint64 = 20
	if si.Pointer() != expected {
		t.Errorf("expected pointer to be %d, got %d", expected, si.Pointer())
	}
}

func TestSliceIterator_Next_RoundRobin(t *testing.T) {
	mgr := lizt.NewManager()
	f, err := lizt.ReadFromFile(filenameTen)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	si := lizt.NewSliceIterator(nameNumbers, f, true)
	mgr.AddIter(si)

	first, err := mgr.MustGet(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}
	second, err := mgr.MustGet(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Errorf("expected %s to be %s", strings.Join(first, ","), strings.Join(second, ","))
	}

	var expected uint64 = 10
	if si.Pointer() != expected {
		t.Errorf("expected pointer to be %d, got %d", expected, si.Pointer())
	}
}

func TestSliceIterator_Next_RoundRobin_NoMoreLines(t *testing.T) {
	mgr := lizt.NewManager()
	f, err := lizt.ReadFromFile(filenameTen)
	if err != nil {
		t.Errorf("ReadFromFile() error = %v", err)
	}

	si := lizt.NewSliceIterator(nameNumbers, f, false)
	mgr.AddIter(si)

	_, err = mgr.MustGet(nameNumbers).Next(10)
	if err != nil {
		t.Errorf("StreamIterator.Next() error = %v", err)
	}

	_, err = mgr.MustGet(nameNumbers).Next(10)
	if !errors.Is(err, lizt.ErrNoMoreLines) {
		t.Errorf("wanted ErrNoMoreLines, got error = %v", err)
	}

	var expected uint64 = 10
	if si.Pointer() != expected {
		t.Errorf("expected pointer to be %d, got %d", expected, si.Pointer())
	}
}

func TestSliceIterator_Next_Parallel(t *testing.T) {
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	b, err := lizt.B().Slice(letters).Build()
	if err != nil {
		t.Errorf("Build() error = %v", err)
	}

	res := make(chan string, b.Len())
	wg := &sync.WaitGroup{}
	wg.Add(26)
	for i := 0; i < b.Len(); i++ {
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

	if !reflect.DeepEqual(results, letters) {
		t.Errorf("expected %v, got %v", letters, results)
	}
}
