package cfg

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

// FileBackend implements a filesystem backend for the configuration
type FileBackend struct{}

// NewFileBackend returns a FileBackend that handles loading and
// saving files from the local filesytem.
func NewFileBackend() *FileBackend {
	return &FileBackend{}
}

// Load loads chains from config
func (fs *FileBackend) Load(u *url.URL) (*Config, error) {
	var config Config
	fixWinURL(u)

	if !exist(u.Path) {
		return &Config{url: u}, nil
	}

	j, err := ioutil.ReadFile(u.Path)
	if err != nil {
		return &config, err
	}

	err = json.Unmarshal(j, &config)
	if err != nil {
		return nil, err
	}
	config.backend = fs
	config.url = u

	return &config, nil
}

// Save saves chains to config
func (fs *FileBackend) Save(config *Config) error {
	fixWinURL(config.URL())
	cfgDir := filepath.Dir(config.URL().Path)
	if !exist(cfgDir) {
		os.MkdirAll(cfgDir, 0755)
	}

	j, err := json.MarshalIndent(config, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(config.URL().Path, j, 0644)
	}

	return err
}
