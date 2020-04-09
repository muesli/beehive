package cfg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	_, err := Load("foobar")
	if err == nil {
		t.Error("Loading an invalid config file should return an error")
	}

	conf, err := Load(filepath.Join("testdata", "beehive.conf"))
	if err != nil {
		t.Error("Error loading config file fixture")
	}

	if conf.Bees[0].Name != "echo" {
		t.Error("The first bee should be an exec bee named echo")
	}
}

func TestSave(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error("Could not create temp file")
	}
	defer os.Remove(tmpfile.Name()) // clean up

	testConfPath := filepath.Join("testdata", "beehive.conf")
	testConf, err := Load(testConfPath)
	if err != nil {
		t.Errorf("Failed to load config fixture %s: %v", testConfPath, err)
	}

	configFile := tmpfile.Name()
	err = Save(configFile, testConf)
	if err != nil {
		t.Errorf("Failed to save the config to %s", configFile)
	}

	if !exist(configFile) {
		t.Error("Configuration file wasn't saved")
	}

	err = Save(filepath.Join(os.TempDir(), "fooconf/beehive.conf"), testConf)
	if err != nil {
		t.Errorf("Failed to create intermediate directories when saving config")
	}
}

func TestSaveCurrent(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error("Could not create temp file")
	}
	defer os.Remove(tmpfile.Name()) // clean up

	t.Run("configFile empty", func(t *testing.T) {
		err = SaveCurrent("")
		if err == nil {
			t.Error("Configuration file should not be saved")
		}
	})

	t.Run("configFile tmpfile", func(t *testing.T) {
		configFile := tmpfile.Name()
		err = SaveCurrent(configFile)
		if err != nil {
			t.Error("Configuration file should have been saved")
		}
	})
}
