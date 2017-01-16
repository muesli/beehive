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
 *      Johannes Fürmann <johannes@weltraumpflege.org>
 */

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/muesli/beehive/api"
	"github.com/muesli/beehive/app"
	_ "github.com/muesli/beehive/filters"
	_ "github.com/muesli/beehive/filters/contains"
	_ "github.com/muesli/beehive/filters/endswith"
	_ "github.com/muesli/beehive/filters/equals"
	_ "github.com/muesli/beehive/filters/matches"
	_ "github.com/muesli/beehive/filters/startswith"
	_ "github.com/muesli/beehive/filters/template"

	"github.com/muesli/beehive/bees"
	_ "github.com/muesli/beehive/bees/anelpowerctrlbee"
	_ "github.com/muesli/beehive/bees/cronbee"
	_ "github.com/muesli/beehive/bees/efabee"
	_ "github.com/muesli/beehive/bees/emailbee"
	_ "github.com/muesli/beehive/bees/execbee"
	_ "github.com/muesli/beehive/bees/htmlextractbee"
	_ "github.com/muesli/beehive/bees/huebee"
	_ "github.com/muesli/beehive/bees/ircbee"
	_ "github.com/muesli/beehive/bees/jabberbee"
	_ "github.com/muesli/beehive/bees/jenkinsbee"
	_ "github.com/muesli/beehive/bees/nagiosbee"
	_ "github.com/muesli/beehive/bees/notificationbee"
	_ "github.com/muesli/beehive/bees/rssbee"
	_ "github.com/muesli/beehive/bees/serialbee"
	_ "github.com/muesli/beehive/bees/slackbee"
	_ "github.com/muesli/beehive/bees/spaceapibee"
	_ "github.com/muesli/beehive/bees/telegrambee"
	_ "github.com/muesli/beehive/bees/timebee"
	_ "github.com/muesli/beehive/bees/transmissionbee"
	_ "github.com/muesli/beehive/bees/tumblrbee"
	_ "github.com/muesli/beehive/bees/twitterbee"
	_ "github.com/muesli/beehive/bees/webbee"
)

var (
	configFile string
)

// This is where we unmarshal our beehive.conf into
type Config struct {
	Bees   []bees.BeeInstance
	Chains []bees.Chain
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
	app.AddFlags([]app.CliFlag{
		{&configFile, "config", "./beehive.conf", "Config-file to use"},
	})

	// Parse command-line args for all registered bees
	app.Run()
	api.Run()

	log.Println()
	log.Println("Beehive is buzzing...")

	config := loadConfig()

	// Initialize bees
	bees.StartBees(config.Bees)
	// Load chains from config
	bees.SetChains(config.Chains)

	// Wait for signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	for s := range ch {
		log.Println("Got signal:", s)

		abort := false
		switch s {
		case syscall.SIGHUP:
			config = loadConfig()
			bees.RestartBees(config.Bees)
			bees.SetChains(config.Chains)

		case syscall.SIGTERM:
			fallthrough
		case syscall.SIGKILL:
			fallthrough
		case syscall.SIGINT:
			abort = true
			break
		}

		if abort {
			break
		}
	}

	// Save chains to config
	log.Println("Storing config...")
	config.Chains = bees.Chains()
	saveConfig(config)
}
