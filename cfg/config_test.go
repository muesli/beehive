package cfg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig("foobar")
	if err == nil {
		t.Error("Loading an invalid config file should return an error")
	}

	conf, err := LoadConfig(filepath.Join("testdata", "beehive.conf"))
	if err != nil {
		t.Error("Error loading config file fixture")
	}

	if conf.Bees[0].Name != "echo" {
		t.Error("The first bee should be an exec bee named echo")
	}
}

func TestLoad(t *testing.T) {
	// override package Lookup func so it always returns an empty array
	oldLookupFunc := lookupFunc
	lookupFunc = func() []string {
		return []string{}
	}
	_, _, err := Load()
	if err != nil {
		t.Error("Should not return an error when there are no configs")
	}
	lookupFunc = oldLookupFunc

	oldCfgFileName := cfgFileName
	cfgFileName = "testdata/beehive.conf"

	path, cfg, err := Load()
	cwd, _ := os.Getwd()
	if path != filepath.Join(cwd, cfgFileName) {
		t.Error("")
	}
	if err != nil || cfg.Bees[0].Name != "echo" {
		t.Error("Should not return an error when there are no configs")
	}
	cfgFileName = oldCfgFileName
}

func TestSaveConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error("Could not create temp file")
	}
	defer os.Remove(tmpfile.Name()) // clean up

	testConfPath := filepath.Join("testdata", "beehive.conf")
	testConf, err := LoadConfig(testConfPath)
	if err != nil {
		t.Errorf("Failed to load config fixture %s: %v", testConfPath, err)
	}

	configFile := tmpfile.Name()
	err = SaveConfig(configFile, testConf)
	if err != nil {
		t.Errorf("Failed to save the config to %s", configFile)
	}

	if !Exist(configFile) {
		t.Error("Configuration file wasn't saved")
	}

	err = SaveConfig(filepath.Join(os.TempDir(), "fooconf/beehive.conf"), testConf)
	if err != nil {
		t.Errorf("Failed to create intermediate directories when saving config")
	}
}

func TestSaveCurrentConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error("Could not create temp file")
	}
	defer os.Remove(tmpfile.Name()) // clean up

	t.Run("configFile empty", func(t *testing.T) {
		err = SaveCurrentConfig("")
		if err == nil {
			t.Error("Configuration file should not be saved")
		}
	})

	t.Run("configFile tmpfile", func(t *testing.T) {
		configFile := tmpfile.Name()
		err = SaveCurrentConfig(configFile)
		if err != nil {
			t.Error("Configuration file should have been saved")
		}
	})
}

func TestFindUserConfigPath(t *testing.T) {
	// override package Lookup func so it always returns an empty array
	oldLookupFunc := lookupFunc
	lookupFunc = func() []string {
		return []string{}
	}
	path := FindUserConfigPath()
	if path != "" {
		t.Error("Should return an empty string when Lookup fails")
	}
	lookupFunc = oldLookupFunc

	// override package internal variable so we can test this without writing
	// to the filesystem
	cfgFileName = "testdata/beehive.conf"
	cwd, _ := os.Getwd()
	if FindUserConfigPath() != filepath.Join(cwd, cfgFileName) {
		t.Error("Should return $CWD/beehive.conf when available")
	}
	cfgFileName = "beehive.conf"
}
