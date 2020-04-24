package cfg

import (
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestAESBackendLoad(t *testing.T) {
	u, _ := url.Parse("crypt://foobar")
	backend := NewAESBackend()

	_, err := backend.Load(u)
	if err != nil {
		t.Error("Loading an non-existing config file should not return an error")
	}

	// try to load the config from an absolute path using a URI
	cwd, _ := os.Getwd()
	u, err = url.Parse("crypto://x:foo@" + filepath.Join(cwd, "testdata", "beehive-crypto.conf"))
	backend = NewAESBackend()
	conf, err := backend.Load(u)
	if err != nil {
		t.Errorf("Error loading config file fixture from absolute path %s. %v", u, err)
	}
	if conf.Bees[0].Name != "echo" {
		t.Error("The first bee should be an exec bee named echo")
	}

	// try to load the config with an invalid password
	u, err = url.Parse("crypto://x:bar@" + filepath.Join(cwd, "testdata", "beehive-crypto.conf"))
	backend = NewAESBackend()
	conf, err = backend.Load(u)
	if err.Error() != "cipher: message authentication failed" {
		t.Errorf("Loading the config file with an invalid password should fail. %v", err)
	}
}

func TestAESBackendSave(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "beehivetest")
	if err != nil {
		t.Error("Could not create temp directory")
	}

	cwd, _ := os.Getwd()
	u, err := url.Parse(filepath.Join("crypto://x:foo", cwd, "testdata", "beehive-crypto.conf"))
	backend := NewAESBackend()
	c, err := backend.Load(u)
	if err != nil {
		t.Errorf("Failed to load config fixture from relative path %s: %v", u, err)
	}

	// Save the config file to a new absolute path using a URL
	p := filepath.Join(tmpdir, "beehive-crypto.conf")
	u, err = url.Parse("crypto://x:foo@" + p)
	c.SetURL(u.String())
	backend = NewAESBackend()
	err = backend.Save(c)
	if err != nil {
		t.Errorf("Failed to save the config to %s", u)
	}
	if !exist(p) {
		t.Errorf("Configuration file wasn't saved to %s", p)
	}
}
