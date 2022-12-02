package lizt

import (
	"reflect"
	"testing"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name string
		want *Manager
	}{
		{
			name: "TestNewManager_SetsFiles_Correctly",
			want: &Manager{
				files: map[string]Iterator{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_AddIter(t *testing.T) {
	type fields struct {
		files map[string]Iterator
		name  string
	}
	type args struct {
		file PointerIterator
	}
	fs, _ := NewStreamIterator(FilenameOneMillion, false)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestManager_AddFile_SliceIterator_Correctly",
			fields: fields{
				files: map[string]Iterator{},
				name:  NameNumbers,
			},
			args: args{
				file: NewSliceIterator(NameNumbers, []string{"a", "b", "c"}, false),
			},
			wantErr: false,
		},
		{
			name: "TestManager_AddFile_StreamIterator_Correctly",
			fields: fields{
				files: map[string]Iterator{},
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
			m := &Manager{
				files: tt.fields.files,
			}
			if err := m.AddIter(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("AddSliceIter() error = %v, wantErr %v", err, tt.wantErr)
			}
			if m.Get(tt.fields.name) != tt.args.file {
				t.Errorf("AddSliceIter() error = %v, wantErr %v", m.Get(NameNumbers), tt.args.file)
			}
		})
	}
}

func TestManager_AddSeeder(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	seeds := []string{"seeder1", "seeder2", "seeder3", "seeder4", "seeder5", "seeder6", "seeder7", "seeder8", "seeder9", "seeder10"}
	seed := NewSeedingIterator(
		NewSliceIterator(NameNumbers, numbers, false),
		NewSliceIterator("seeds", seeds, true), 2,
	)

	m := NewManager()
	err := m.AddIter(seed)
	if err != nil {
		t.Errorf("AddIter() error = %v", err)
	}
	if m.Get(NameNumbers).Len() != len(numbers) {
		t.Errorf("expected %d, got %d", len(numbers), m.Get(NameNumbers).Len())
	}
}

func Test_makeName(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should extract file name properly", args{"test/test/test/test.txt"}, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeNameFromFilename(tt.args.filename); got != tt.want {
				t.Errorf("makeNameFromFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	numbers := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	shuffled := Shuffle(numbers)
	if reflect.DeepEqual(numbers, shuffled) {
		t.Errorf("expected %v to be different from %v", numbers, shuffled)
	}
}

func TestManager_AddDirIter(t *testing.T) {
	m := NewManager()
	err := m.AddDirIter("test/", false)
	if err != nil {
		t.Errorf("AddDirIter() error = %v", err)
	}
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
}
