package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/muesli/beehive/bees"
	gap "github.com/muesli/go-app-paths"
	log "github.com/sirupsen/logrus"
)

const appName = "beehive"

var cfgFileName = "beehive.conf"
var lookupFunc = Lookup

// Config contains an entire configuration set for Beehive
type Config struct {
	Bees    []bees.BeeConfig
	Actions []bees.Action
	Chains  []bees.Chain
}

// DefaultPath returns Beehive's default config path
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

// Lookup an existing user config file
//
// Returns an empty array when no configs are found.
func Lookup() []string {
	paths := []string{}
	defaultPath := DefaultPath()
	if Exist(defaultPath) {
		paths = append(paths, defaultPath)
	}

	// Prepend ./beehive.conf to the search path if exists, takes priority
	// over the rest
	cwd, err := os.Getwd()
	if err != nil {
		log.Errorf("Error getting current working directory. err: %v", err)
		cwd = "."
	}
	cwdCfg := filepath.Join(cwd, cfgFileName)
	if Exist(cwdCfg) {
		paths = append([]string{cwdCfg}, paths...)
	}

	return paths
}

// FindUserConfigPath tries to find the config file.
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
func FindUserConfigPath() string {
	paths := lookupFunc()
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// Load tries to load the default config file, returning an empty config object
// if not found
//
// Returns the path of the config found and the Config object, returning an error
// also if loading the config from the filesystem fails.
func Load() (string, Config, error) {
	path := FindUserConfigPath()
	if path == "" {
		log.Info("No config file found, loading defaults")
		return DefaultPath(), Config{}, nil
	}
	log.Infof("Loading config file from %s", path)
	cfg, err := LoadConfig(path)
	return path, cfg, err
}

// Loads chains from config
func LoadConfig(file string) (Config, error) {
	config := Config{}

	j, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(j, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Saves chains to config
func SaveConfig(file string, c Config) error {
	cfgDir := filepath.Dir(file)
	if !Exist(cfgDir) {
		os.MkdirAll(cfgDir, 0755)
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(file, j, 0644)
	}

	return err
}

func SaveCurrentConfig(file string) error {
	config := Config{}
	config.Bees = bees.BeeConfigs()
	config.Chains = bees.GetChains()
	config.Actions = bees.GetActions()
	return SaveConfig(file, config)
}

func Exist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	return false
}
