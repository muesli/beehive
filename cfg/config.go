package cfg

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/muesli/beehive/bees"
	gap "github.com/muesli/go-app-paths"
	log "github.com/sirupsen/logrus"
)

const appName = "beehive"

var cfgFileName = "beehive.conf"

// Config contains an entire configuration set for Beehive
type Config struct {
	Bees    []bees.BeeConfig
	Actions []bees.Action
	Chains  []bees.Chain
	backend ConfigBackend
	url     *url.URL
}

// ConfigBackend is the interface implemented by the configuration backends.
//
// Backends are responsible for loading and saving the Config struct to
// memory, the local filesystem, the network, etc.
type ConfigBackend interface {
	Load(*url.URL) (*Config, error)
	Save(*Config) error
}

// Save the current configuration.
//
// The backend loaded will be responsible for saving it
// to the given URL
func (c *Config) Save() error {
	return c.backend.Save(c)
}

// Load the configuration.
//
// The backend loaded will be responsible for loading it
// from the given URL.
func (c *Config) Load() error {
	config, err := c.backend.Load(c.url)
	if err != nil {
		return err
	}
	c.Bees = config.Bees
	c.Actions = config.Actions
	c.Chains = config.Chains
	return nil
}

// Backend currently being used.
func (c *Config) Backend() ConfigBackend {
	return c.backend
}

// SetURL updates the configuration URL.
//
// Next time the config is loaded or saved
// the new URL will be used.
func (c *Config) SetURL(u string) error {
	url, err := url.Parse(u)
	if err != nil {
		return err
	}

	c.url = url

	return nil
}

// URL currently being used.
func (c *Config) URL() *url.URL {
	return c.url
}

// New returns a new Config struct.
//
// The URL will be matched against all the supported
// backends and the first backend that can handle the
// URL scheme will be loaded.
//
// A UNIX style path is also accepted, and will be handled
// by the default FileBackend.
func New(url string) (*Config, error) {
	config := &Config{}
	var backend ConfigBackend

	if url == "" {
		return nil, fmt.Errorf("Empty URL provided but not supported")
	}

	err := config.SetURL(url)
	if err != nil {
		return nil, err
	}

	switch config.url.Scheme {
	case "", "file":
		if ok, _ := IsEncrypted(config.url); ok {
			log.Debugf("Loading encrypted configuration file")
			backend, err = NewAESBackend(config.url)
			if err != nil {
				log.Fatalf("error loading the AES configuration backend. err: %v", err)
			}
		} else {
			backend = NewFileBackend()
		}
	case "mem":
		backend = NewMemBackend()
	case "crypto":
		backend, err = NewAESBackend(config.url)
		if err != nil {
			log.Fatalf("error loading the AES configuration backend. err: %v", err)
		}
	default:
		return nil, fmt.Errorf("Configuration backend '%s' not supported", config.url.Scheme)
	}

	config.backend = backend

	return config, nil
}

// DefaultPath returns Beehive's default config path.
//
// The path returned is OS dependant. If there's an error
// while trying to figure out the OS dependant path, "beehive.conf"
// in the current working dir is returned.
func DefaultPath() string {
	userScope := gap.NewScope(gap.User, appName)
	path, err := userScope.ConfigPath(cfgFileName)
	if err != nil {
		return cfgFileName
	}

	return path
}

// Lookup tries to find the config file.
//
// If a config file is found in the current working directory, that's returned.
// Otherwise we try to locate it following an OS dependant:
//
// Unix:
//   - ~/.config/app/filename.conf
// macOS:
//   - ~/Library/Preferences/app/filename.conf
// Windows:
//   - %LOCALAPPDATA%/app/Config/filename.conf
//
// If no valid config file is found, an empty string is returned.
func Lookup() string {
	paths := []string{}
	defaultPath := DefaultPath()
	if exist(defaultPath) {
		paths = append(paths, defaultPath)
	}

	// Prepend ./beehive.conf to the search path if exists, takes priority
	// over the rest.
	cwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Error getting current working directory. err: %v", err)
		cwd = "."
	}
	cwdCfg := filepath.Join(cwd, cfgFileName)
	if exist(cwdCfg) {
		paths = append([]string{cwdCfg}, paths...)
	}
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func exist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
