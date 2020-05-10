package cfg

import (
	"path/filepath"
	"testing"
)

func TestMemLoad(t *testing.T) {
	u, _ := ParseURL("mem://")
	backend := NewMemBackend()
	_, err := backend.Load(u)
	if err != nil {
		t.Error("Loading an invalid config file should return an error")
	}
}

func TestMemSave(t *testing.T) {
	path := filepath.Join("testdata", "foobar")
	u, _ := ParseURL(filepath.Join("testdata", "foobar"))
	backend := NewMemBackend()
	conf := &Config{url: u}
	err := backend.Save(conf)
	if err != nil {
		t.Errorf("Failed to save the config to memory")
	}

	if exist(path) {
		t.Error("Configuration file should not exist")
	}
}
