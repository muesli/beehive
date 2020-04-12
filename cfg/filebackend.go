package cfg

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

// FileBackend implements a filesystem backend for the configuration
type FileBackend struct{ uri *url.URL }

func (fb *FileBackend) URI() string {
	return fb.uri.String()
}

func (fb *FileBackend) SetURI(uri string) error {
	url, err := url.Parse(uri)
	if err != nil {
		return err
	}
	fb.uri = url
	return nil
}

// Load loads chains from config
func (fs *FileBackend) Load() (Config, error) {
	var config Config
	j, err := ioutil.ReadFile(fs.uri.Path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(j, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Save saves chains to config
func (fs *FileBackend) Save(config Config) error {
	cfgDir := filepath.Dir(fs.uri.Path)
	if !exist(cfgDir) {
		os.MkdirAll(cfgDir, 0755)
	}

	j, err := json.MarshalIndent(config, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(fs.uri.Path, j, 0644)
	}

	return err
}
