package cfg

import "net/url"

// MemBackend implements a dummy memory backend for the configuration
type MemBackend struct {
	conf *Config
}

// Load loads chains from config
func NewMemBackend() *MemBackend {
	return &MemBackend{conf: &Config{}}
}

func (m *MemBackend) Load(u *url.URL) (*Config, error) {
	return m.conf, nil
}

// Save saves chains to config
func (m *MemBackend) Save(config *Config) error {
	return nil
}
