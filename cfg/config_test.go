package cfg

import (
	"testing"
)

func TestNewBackend(t *testing.T) {
	backend, err := NewBackend("/foobar")
	if err != nil {
		panic(err)
	}
	if _, ok := backend.(*FileBackend); !ok {
		t.Error("Backend for '/foobar' should be a FileBackend")
	}

	backend, err = NewBackend("file:///foobar")
	if err != nil {
		panic(err)
	}
	if _, ok := backend.(*FileBackend); !ok {
		t.Error("Backend for 'file:///foobar' should be a FileBackend")
	}

	backend, err = NewBackend("mem:")
	if err != nil {
		panic(err)
	}
	if _, ok := backend.(*MemBackend); !ok {
		t.Error("Backend for 'mem:' should be a MemoryBackend")
	}
}
