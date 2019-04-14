package cfg

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/muesli/beehive/bees"
)

var configFile string

// Config contains an entire configuration set for Beehive
type Config struct {
	Bees    []bees.BeeConfig
	Actions []bees.Action
	Chains  []bees.Chain
}

// Loads chains from config
func LoadConfig(file string) Config {
	configFile = file
	config := Config{}

	j, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(j, &config)
		if err != nil {
			log.Fatal("Error parsing config file: ", err)
		}
	}

	return config
}

// Saves chains to config
func SaveConfig(c Config) {
	j, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(configFile, j, 0644)
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Configuration saved")
}

func SaveCurrentConfig() {
	config := Config{}
	config.Bees = bees.BeeConfigs()
	config.Chains = bees.GetChains()
	config.Actions = bees.GetActions()
	SaveConfig(config)
}

func GetConfig() Config {
	return LoadConfig(configFile)
}
