package cfg

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLoad(t *testing.T) {
	u, _ := url.Parse("file://foobar")
	backend := NewFileBackend()

	_, err := backend.Load(u)
	if err != nil {
		t.Error("Loading an non-existing config file should not return an error")
	}

	// try to load the config from a relative path
	u, err = url.Parse(filepath.Join("testdata", "beehive.conf"))
	if err != nil {
		t.Error("cannot parse config path")
	}
	backend = NewFileBackend()
	conf, err := backend.Load(u)
	if err != nil {
		t.Errorf("Error loading config file fixture from relative path %s. %v", u, err)
	}
	if conf.Bees[0].Name != "echo" {
		t.Error("The first bee should be an exec bee named echo")
	}

	// try to load the config from an absolute path using a URI
	cwd, _ := os.Getwd()
	u, err = url.Parse(filepath.Join("file://", cwd, "testdata", "beehive.conf"))
	if err != nil {
		t.Error("cannot parse config path")
	}
	backend = NewFileBackend()
	conf, err = backend.Load(u)
	if err != nil {
		t.Errorf("Error loading config file fixture from absolute path %s. %v", u, err)
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

	u, err := url.Parse(filepath.Join("testdata", "beehive.conf"))
	if err != nil {
		t.Error("cannot parse config path")
	}
	backend := NewFileBackend()
	c, err := backend.Load(u)
	if err != nil {
		t.Errorf("Failed to load config fixture from relative path %s: %v", u, err)
	}

	// Save the config file to a new absolute path using a URL
	p := filepath.Join(tmpdir, "beehive.conf")
	u, err = url.Parse("file://" + p)
	if err != nil {
		t.Error("cannot parse config path")
	}
	err = c.SetURL(u.String())
	if err != nil {
		t.Error("cannot set url")
	}
	backend = NewFileBackend()
	err = backend.Save(c)
	if err != nil {
		t.Errorf("Failed to save the config to %s", u)
	}
	if !exist(p) {
		t.Errorf("Configuration file wasn't saved to %s", p)
	}
	c, err = backend.Load(u)
	if err != nil {
		t.Errorf("Failed to load config fixture from absolute path %s: %v", u, err)
	}

	// Save the config file to a new absolute path using a regular path
	p = filepath.Join(tmpdir, "beehive.conf")
	_, err = url.Parse(p)
	if err != nil {
		t.Error("cannot parse config path")
	}
	err = backend.Save(c)
	if err != nil {
		t.Errorf("Failed to save the config to %s", p)
	}
	if !exist(p) {
		t.Errorf("Configuration file wasn't saved to %s", p)
	}
}
