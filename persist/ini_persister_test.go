package persist

import (
	"git.faze.center/netr/lizt"
	"os"
	"strings"
	"testing"
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

	if doesFileExist(path) == false {
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

func doesFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
