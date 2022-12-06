package persist

import (
	"os"
	"strings"
	"testing"

	"git.faze.center/netr/lizt"
)

func TestNewIniPersister(t *testing.T) {
	path := "../test/pointers.ini"
	_ = os.Remove(path)

	persist, err := NewIniPersister(path)
	if err != nil {
		t.Errorf("NewIniPersister() error = %v", err)
	}

	err = persist.Set("test", 1)
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	if lizt.DoesFileExist(path) == false {
		t.Errorf("Expected file to exist")
	}

	x, err := lizt.ReadFileToString(path)
	if err != nil {
		t.Errorf("ReadFileToString() error = %v", err)
	}

	if strings.Contains(x, "test = 1") == false {
		t.Errorf("Expected file to contain test = 1")
	}

	_ = os.Remove(path)
}

func TestNewIniPersister_Builder(t *testing.T) {
	path := "../test/pointers.ini"
	_ = os.Remove(path)

	persist, err := NewIniPersister(path)
	if err != nil {
		t.Errorf("NewIniPersister() error = %v", err)
	}

	tester := lizt.B().SliceNamedRR("test", []string{"test", "test2", "test3"}).PersistTo(persist).MustBuild()
	tester.MustNextOne()

	if lizt.DoesFileExist(path) == false {
		t.Errorf("Expected file to exist")
	}

	x, err := lizt.ReadFileToString(path)
	if err != nil {
		t.Errorf("ReadFileToString() error = %v", err)
	}

	if strings.Contains(x, "test = 1") == false {
		t.Errorf("Expected file to contain test = 1")
	}

	_ = os.Remove(path)
}
