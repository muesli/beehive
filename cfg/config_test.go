package cfg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	conf, err := New("/foobar")
	if err != nil {
		t.Error("cannot create config from path")
		t.FailNow()
	}
	if _, ok := conf.Backend().(*FileBackend); !ok {
		t.Error("Backend for '/foobar' should be a FileBackend")
	}

	conf, err = New("file:///foobar")
	if err != nil {
		t.Error("cannot create config from file:// path")
		t.FailNow()
	}
	if _, ok := conf.Backend().(*FileBackend); !ok {
		t.Error("Backend for 'file:///foobar' should be a FileBackend")
	}

	cwd, _ := os.Getwd()
	p := filepath.Join(cwd, "testdata/beehive-crypto.conf")
	conf, err = New(p)
	if err != nil {
		panic(err)
	}
	if _, ok := conf.Backend().(*AESBackend); !ok {
		t.Errorf("Backend for '%s' should be an AESBackend", p)
	}

	conf, err = New("mem:")
	if err != nil {
		t.Error("cannot create config from memory")
		t.FailNow()
	}
	if _, ok := conf.Backend().(*MemBackend); !ok {
		t.Error("Backend for 'mem:' should be a MemoryBackend")
	}

	_, err = New("c:\\foobar")
	if err == nil {
		t.Error("Not a valid URL, should return an error")
	}

	_, err = New("")
	if err == nil {
		t.Error("Not a valid URL, should return an error")
	}
}

func TestLoad(t *testing.T) {
	conf, err := New("mem://")
	if err != nil {
		panic(err)
	}
	err = conf.Load()
	if err != nil {
		t.Error("Failed loading the configuration")
	}
	if conf.URL().Scheme != "mem" {
		t.Error("Config URL didn't change")
	}
}
