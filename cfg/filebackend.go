package cfg

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Format int

const (
	FormatJSON Format = iota
	FormatYAML        = iota
)

// FileBackend implements a filesystem backend for the configuration
type FileBackend struct {
	format Format
}

// NewFileBackend returns a FileBackend that handles loading and
// saving files from the local filesytem.
func NewFileBackend() *FileBackend {
	return &FileBackend{format: FormatJSON}
}

// Load loads chains from config
func (fs *FileBackend) Load(u *url.URL) (*Config, error) {
	var config Config

	// detect file format by extension
	if strings.HasSuffix(u.Path, ".yaml") {
		fs.format = FormatYAML
	} else if strings.HasSuffix(u.Path, ".yml") {
		fs.format = FormatYAML
	} else {
		fs.format = FormatJSON
	}

	if !exist(u.Path) {
		return &Config{url: u}, nil
	}

	content, err := ioutil.ReadFile(u.Path)
	if err != nil {
		return &config, err
	}

	if fs.format == FormatYAML {
		err = yaml.Unmarshal(content, &config)
	} else {
		err = json.Unmarshal(content, &config)
	}
	if err != nil {
		return nil, err
	}
	config.backend = fs
	config.url = u

	return &config, nil
}

// Save saves chains to config
func (fs *FileBackend) Save(config *Config) error {
	cfgDir := filepath.Dir(config.URL().Path)
	if !exist(cfgDir) {
		err := os.MkdirAll(cfgDir, 0755)
		if err != nil {
			return err
		}
	}

	var content []byte
	var err error
	if fs.format == FormatYAML {
		content, err = yaml.Marshal(config)
	} else {
		content, err = json.MarshalIndent(config, "", "  ")
	}
	if err != nil {
		return err
	}
	return ioutil.WriteFile(config.URL().Path, content, 0644)
}
