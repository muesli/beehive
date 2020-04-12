package cfg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLoad(t *testing.T) {
	backend, err := NewBackend("file://foobar")
	if err != nil {
		t.Error("Error loading the file system backend")
	}

	_, err = backend.Load()
	if err == nil {
		t.Error("Loading an invalid config file should return an error")
	}

	// try to load the config from a relative path
	cfgURI := filepath.Join("testdata", "beehive.conf")
	backend, err = NewBackend(cfgURI)
	conf, err := backend.Load()
	if err != nil {
		t.Errorf("Error loading config file fixture from relative path %s. %v", cfgURI, err)
	}
	if conf.Bees[0].Name != "echo" {
		t.Error("The first bee should be an exec bee named echo")
	}

	// try to load the config from an absolute path using a URI
	cwd, _ := os.Getwd()
	cfgURI = filepath.Join("file://", cwd, "testdata", "beehive.conf")
	backend, err = NewBackend(cfgURI)
	conf, err = backend.Load()
	if err != nil {
		t.Errorf("Error loading config file fixture from absolute path %s. %v", cfgURI, err)
	}
	if conf.Bees[0].Name != "echo" {
		t.Error("The first bee should be an exec bee named echo")
	}
}

func TestFileSave(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "beehivetest")
	if err != nil {
		t.Error("Could not create temp directory")
	}

	testConfPath := filepath.Join("testdata", "beehive.conf")
	backend, _ := NewBackend(testConfPath)
	testConf, err := backend.Load()
	if err != nil {
		t.Errorf("Failed to load config fixture from relative path %s: %v", testConfPath, err)
	}

	// Save the config file to a new absolute path using a URI
	configFile := filepath.Join(tmpdir, "beehive.conf")
	backend.SetURI("file://" + configFile)
	err = backend.Save(testConf)
	if err != nil {
		t.Errorf("Failed to save the config to %s", configFile)
	}
	if !exist(configFile) {
		t.Errorf("Configuration file wasn't saved to %s", configFile)
	}

	// Save the config file to a new absolute path using a regular path
	configFile = filepath.Join(tmpdir, "beehive.conf")
	backend.SetURI(configFile)
	err = backend.Save(testConf)
	if err != nil {
		t.Errorf("Failed to save the config to %s", configFile)
	}
	if !exist(configFile) {
		t.Errorf("Configuration file wasn't saved to %s", configFile)
	}
}
