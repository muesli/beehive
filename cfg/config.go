package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/muesli/beehive/bees"
)

// Config contains an entire configuration set for Beehive
type Config struct {
	Bees    []bees.BeeConfig
	Actions []bees.Action
	Chains  []bees.Chain
}

// Loads chains from config
func LoadConfig(file string) (Config, error) {
	var config Config

	j, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	config = Config{}
	err = json.Unmarshal(j, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Saves chains to config
func SaveConfig(file string, c Config) error {
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
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}
