/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
 *                  2014      Michael Wendland
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
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"

	"github.com/muesli/beehive/api"
	"github.com/muesli/beehive/app"
	_ "github.com/muesli/beehive/filters"
	_ "github.com/muesli/beehive/filters/template"

	"github.com/muesli/beehive/bees"
)

var (
	configFile  string
	versionFlag bool
)

// Config contains an entire configuration set for Beehive
type Config struct {
	Bees    []bees.BeeConfig
	Actions []bees.Action
	Chains  []bees.Chain
}

// Loads chains from config
func loadConfig() Config {
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
		{
			V:     &configFile,
			Name:  "config",
			Value: "./beehive.conf",
			Desc:  "Config-file to use",
		},
	})

	app.AddFlags([]app.CliFlag{
		{
			V:     &versionFlag,
			Name:  "version",
			Value: false,
			Desc:  "Beehive version",
		},
	})

	// Parse command-line args for all registered bees
	app.Run()

	if versionFlag {
		fmt.Printf("Beehive %s (%s)\n", Version, CommitSHA)
		os.Exit(0)
	}

	api.Run()

	log.Println()
	log.Println("Beehive is buzzing...")

	config := loadConfig()

	// Load actions from config
	bees.SetActions(config.Actions)
	// Load chains from config
	bees.SetChains(config.Chains)
	// Initialize bees
	bees.StartBees(config.Bees)

	// Wait for signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL)

	for s := range ch {
		log.Println("Got signal:", s)

		abort := false
		switch s {
		case syscall.SIGHUP:
			config = loadConfig()
			bees.StopBees()
			bees.SetActions(config.Actions)
			bees.SetChains(config.Chains)
			bees.StartBees(config.Bees)

		case syscall.SIGTERM:
			fallthrough
		case syscall.SIGKILL:
			fallthrough
		case os.Interrupt:
			abort = true
			break
		}

		if abort {
			break
		}
	}

	// Save chains to config
	log.Println("Storing config...")
	config.Bees = bees.BeeConfigs()
	config.Chains = bees.GetChains()
	config.Actions = bees.GetActions()
	saveConfig(config)
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
}
