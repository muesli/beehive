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
	// over the rest
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

// Load loads chains from config
func Load(file string) (Config, error) {
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

// Save saves chains to config
func Save(file string, c Config) error {
	cfgDir := filepath.Dir(file)
	if !exist(cfgDir) {
		os.MkdirAll(cfgDir, 0755)
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(file, j, 0644)
	}

	return err
}

// SaveCurrent saves current in-memory configuration to the config file
func SaveCurrent(file string) error {
	config := Config{}
	config.Bees = bees.BeeConfigs()
	config.Chains = bees.GetChains()
	config.Actions = bees.GetActions()
	return Save(file, config)
}

func exist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	return false
}
