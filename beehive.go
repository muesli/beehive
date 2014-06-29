/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
 *                  2014 Michael Wendland
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Christian Muehlhaeuser <muesli@gmail.com>
 *      Michael Wendland <michiwend@michiwend.com>
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/muesli/beehive/app"
	_ "github.com/muesli/beehive/filters"
	_ "github.com/muesli/beehive/filters/contains"
	_ "github.com/muesli/beehive/filters/endswith"
	_ "github.com/muesli/beehive/filters/startswith"

	"github.com/muesli/beehive/modules"
	//_ "github.com/muesli/beehive/modules/hellobee"
	_ "github.com/muesli/beehive/modules/anelpowerctrlbee"
	_ "github.com/muesli/beehive/modules/ircbee"
	_ "github.com/muesli/beehive/modules/jabberbee"
	_ "github.com/muesli/beehive/modules/jenkinsbee"
	_ "github.com/muesli/beehive/modules/nagiosbee"
	_ "github.com/muesli/beehive/modules/rssbee"
	_ "github.com/muesli/beehive/modules/webbee"
	_ "github.com/muesli/beehive/modules/timebee"
)

var (
	configFile = "./beehive.conf"
)

type Config struct {
	Bees   []modules.Bee
	Chains []modules.Chain
}

// Loads chains from config
func loadConfig() Config {
	config := Config{}

	j, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(j, &config)
	}

	if err != nil {
		log.Fatal(err)
	}

	return config
}

// Saves chains to config
func saveConfig(c Config) {
	j, err := json.MarshalIndent(c, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(configFile, j, 0644)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Parse command-line args for all registered modules
	app.Run()

	log.Println()
	log.Println("Beehive is buzzing...")

	config := loadConfig()

	// Initialize modules
	modules.StartModules(config.Bees)
	// Load chains from config
	modules.SetChains(config.Chains)

	// Keep app alive
	ch := make(chan bool)
	<-ch

	// Save chains to config
	config.Chains = modules.Chains()
	saveConfig(config)
}
