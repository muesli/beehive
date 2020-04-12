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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"

	"github.com/muesli/beehive/api"
	"github.com/muesli/beehive/app"
	"github.com/muesli/beehive/cfg"
	_ "github.com/muesli/beehive/filters"
	_ "github.com/muesli/beehive/filters/template"

	"github.com/muesli/beehive/bees"
)

var (
	configURI   string
	versionFlag bool
	debugFlag   bool
)

func main() {
	app.AddFlags([]app.CliFlag{
		{
			V:     &configURI,
			Name:  "config",
			Value: cfg.DefaultPath(),
			Desc:  "Default configuration path",
		},
		{
			V:     &versionFlag,
			Name:  "version",
			Value: false,
			Desc:  "Beehive version",
		},
		{
			V:     &debugFlag,
			Name:  "debug",
			Value: false,
			Desc:  "Turn on debugging",
		},
	})

	// Parse command-line args for all registered bees
	app.Run()

	if versionFlag {
		fmt.Printf("Beehive %s (%s)\n", Version, CommitSHA)
		os.Exit(0)
	}

	api.Run()

	if debugFlag {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Println()
	log.Println("Beehive is buzzing...")

	config := cfg.Config{}
	var err error
	cfgBackend, err := cfg.NewBackend(configURI)
	if err != nil {
		log.Fatalf("Error loading the configuration backend %s", err)
	}

	if cfgBackend.URI() != cfg.DefaultPath() { // the user specified a custom config path or URI
		config, err = cfgBackend.Load()
		if err != nil {
			log.Fatalf("Error loading configuration file from %s. err: %v", cfgBackend.URI(), err)
		}
		log.Infof("Loading configuration from %s", cfgBackend.URI())
	} else { // try to load default config from user paths
		path := cfg.Lookup()
		if path == "" {
			log.Info("No config file found, loading defaults")
		} else {
			cfgBackend.SetURI("file://" + path)
			log.Infof("Loading config file from %s", path)
			config, err = cfgBackend.Load()
			if err != nil {
				log.Fatalf("Error loading user config file %s. err: %v", path, err)
			}
		}
	}

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
			config, err := cfgBackend.Load()
			if err != nil {
				log.Panicf("Error loading config from %s: %v", cfgBackend.URI(), err)
			}
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

	// Save actions & chains to config
	log.Printf("Saving config to %s", cfgBackend.URI())
	config.Bees = bees.BeeConfigs()
	config.Chains = bees.GetChains()
	config.Actions = bees.GetActions()
	err = cfgBackend.Save(config)
	if err != nil {
		log.Printf("Error saving config file to %s! %v", cfgBackend.URI(), err)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
}
