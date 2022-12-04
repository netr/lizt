package lizt_test

import (
	"reflect"
	"testing"

	"git.faze.center/netr/lizt"
)

var (
	filenameOneMillion = "test/1000000.txt"
	filenameTen        = "test/10.txt"
	nameNumbers        = "numbers"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		want *lizt.Manager
		name string
	}{
		{
			name: "TestNewManager_SetsFiles_Correctly",
			want: lizt.NewManager(),
		},
	}
	for _, tt := range tests {
		t.Parallel()
		t.Run(tt.name, func(t *testing.T) {
			if got := lizt.NewManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_AddIter(t *testing.T) {
	type fields struct {
		files map[string]lizt.Iterator
		name  string
	}
	type args struct {
		file lizt.PointerIterator
	}
	fs, err := lizt.NewStreamIterator(filenameOneMillion, false)
	if err != nil {
		t.Errorf("NewStreamIterator() error = %v", err)
	}
	tests := []struct {
		args    args
		fields  fields
		name    string
		wantErr bool
	}{
		{
			name: "TestManager_AddFile_SliceIterator_Correctly",
			fields: fields{
				files: map[string]lizt.Iterator{},
				name:  nameNumbers,
			},
			args: args{
				file: lizt.NewSliceIterator(nameNumbers, []string{"a", "b", "c"}, false),
			},
			wantErr: false,
		},
		{
			name: "TestManager_AddFile_StreamIterator_Correctly",
			fields: fields{
				files: make(map[string]lizt.Iterator, 0),
				name:  "1000000",
			},
			args: args{
				file: fs,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := lizt.NewManager()
			m.AddIter(tt.args.file)
			if m.MustGet(tt.fields.name) != tt.args.file {
				t.Errorf("AddSliceIter() error = %v, wantErr %v", nameNumbers, tt.args.file)
			}
		})
	}
}

func TestManager_AddSeeder(t *testing.T) {
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

	mgr := lizt.NewManager().AddIter(seed)
	if mgr.MustGet(nameNumbers).Len() != len(numbers) {
		t.Errorf("expected %d, got %d", len(numbers), mgr.MustGet(nameNumbers).Len())
	}
}

func TestShuffle(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	shuffled := lizt.Shuffle(numbers)
	if reflect.DeepEqual(numbers, shuffled) {
		t.Errorf("expected %v to be different from %v", numbers, shuffled)
	}
}

func TestManager_AddDirIter(t *testing.T) {
	mgr := lizt.NewManager()
	err := mgr.AddDirIter("test/", false)
	if err != nil {
		t.Errorf("AddDirIter() error = %v", err)
	}
	if mgr.Len() != 2 {
		t.Errorf("expected 2, got %d", mgr.Len())
	}
}

func TestManager_SmartAddDirIter(t *testing.T) {
	mgr := lizt.NewManager()
	err := mgr.SmartAddDirIter("test/", false)
	if err != nil {
		t.Errorf("SmartAddDirIter() error = %v", err)
	}
	if mgr.Len() != 2 {
		t.Errorf("expected 2, got %d", mgr.Len())
	}

	millionIter, err := mgr.Get("1000000")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}

	if reflect.TypeOf(millionIter).Elem().Name() != "StreamIterator" {
		t.Errorf("expected StreamIterator, got %s", reflect.TypeOf(millionIter).Name())
	}

	tenIter, err := mgr.Get("10")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if reflect.TypeOf(tenIter).Elem().Name() != "SliceIterator" {
		t.Errorf("expected SliceIterator, got %s", reflect.TypeOf(tenIter).Name())
	}
}
