package persist

import (
	"errors"
	"fmt"
	"git.faze.center/netr/lizt"
	"gopkg.in/ini.v1"
	"os"
	"sync"
)

type IniPersister struct {
	lizt.PersistentIterator
	iniPath string
	mu      sync.RWMutex
	iniFile *ini.File
}

func getCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.New("failed to get current directory")
	}
	return dir, nil
}

func NewIniPersister(iniPath string) (*IniPersister, error) {
	if lizt.DoesFileExist(iniPath) == false {
		_, err := os.Create(iniPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create ini file: %s -> %w", iniPath, err)
		}
	}
	cfg, err := ini.Load(iniPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ini file: %s -> %w", iniPath, err)
	}

	return &IniPersister{
		iniFile: cfg,
		iniPath: iniPath,
	}, nil
}

func (i *IniPersister) Set(key string, value uint64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.iniFile.Section("pointers").Key(key).SetValue(fmt.Sprintf("%d", value))
	err := i.iniFile.SaveTo(i.iniPath)
	if err != nil {
		return fmt.Errorf("failed to save ini file: %s -> %w", i.iniPath, err)
	}
	return nil
}

var ErrNotFound = errors.New("not found")

func (i *IniPersister) Get(key string) (uint64, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	val, err := i.iniFile.Section("pointers").Key(key).Uint64()
	if err != nil {
		return 0, ErrNotFound
	}

	return val, nil
}
