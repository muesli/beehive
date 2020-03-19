package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/muesli/beehive/bees"
	gap "github.com/muesli/go-app-paths"
)

const appName = "beehive"
const cfgFileName = "beehive.conf"

// This can be easily replaced in tests so we can test without using
// the real search paths
var searchPaths = defaultSearchPaths

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

func defaultSearchPaths() []string {
	userScope := gap.NewScope(gap.User, appName)
	paths := []string{}
	userCfg, err := userScope.ConfigPath(cfgFileName)
	if err == nil {
		paths = append(paths, userCfg)
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return append(paths, filepath.Join(cwd, cfgFileName))
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
func FindUserConfigPath() string {
	path := ""
	for _, sPath := range searchPaths() {
		if Exist(sPath) {
			path = sPath
			break
		}
	}

	return path
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
