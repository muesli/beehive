package cfg

import (
	"path/filepath"
	"testing"
)

func TestMemLoad(t *testing.T) {
	backend, err := NewBackend("mem://mem")
	if err != nil {
		t.Error("Error loading the mem backend")
	}

	_, err = backend.Load()
	if err != nil {
		t.Error("Loading an invalid config file should return an error")
	}
}

func TestMemSave(t *testing.T) {
	testConfPath := filepath.Join("testdata", "foobar")
	backend, _ := NewBackend("mem://" + testConfPath)

	conf := Config{}
	err := backend.Save(conf)
	if err != nil {
		t.Errorf("Failed to save the config to memory")
	}

	if exist(testConfPath) {
		t.Error("Configuration file should not exist")
	}
}
