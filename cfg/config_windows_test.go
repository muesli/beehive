package cfg

import "testing"

func TestWindowsStylePaths(t *testing.T) {
	conf, err := New("c:/foo/bar/beehive.conf")
	if err != nil {
		t.Fatalf("Error in New. %v", err)
	}
	if _, ok := conf.Backend().(*FileBackend); !ok {
		t.Errorf("Backend for %s should be a FileBackend", conf.URL().Raw)
	}
}
