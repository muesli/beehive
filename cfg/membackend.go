package cfg

import "net/url"

// MemBackend implements a dummy memory backend for the configuration
type MemBackend struct{ uri *url.URL }

func (mb *MemBackend) URI() string {
	return mb.uri.String()
}

func (mb *MemBackend) SetURI(uri string) error {
	url, err := url.Parse(uri)
	if err != nil {
		return err
	}

	mb.uri = url
	return nil
}

// Load loads chains from config
func (m *MemBackend) Load() (Config, error) {
	return Config{}, nil
}

// Save saves chains to config
func (m *MemBackend) Save(config Config) error {
	return nil
}
