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

	"github.com/muesli/beehive/modules"
	//_ "github.com/muesli/beehive/modules/hellobee"
	_ "github.com/muesli/beehive/modules/ircbee"
	_ "github.com/muesli/beehive/modules/jabberbee"
	_ "github.com/muesli/beehive/modules/webbee"
)

var (
	config = "./beehive.conf"
)

// Loads chains from config
func loadConfig() []modules.Chain {
	chains := []modules.Chain{}

	j, err := ioutil.ReadFile(config)
	if err == nil {
		err = json.Unmarshal(j, &chains)
	}

	if err != nil {
		log.Fatal(err)
	}

	return chains
}

// Saves chains to config
func saveConfig(chains []modules.Chain) {
	j, err := json.MarshalIndent(chains, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(config, j, 0644)
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

	// Initialize modules
	modules.StartModules()
	// Load chains from config
	modules.SetChains(loadConfig())

	// Keep app alive
	ch := make(chan bool)
	<-ch

	// Save chains to config
	saveConfig(modules.Chains())
}
