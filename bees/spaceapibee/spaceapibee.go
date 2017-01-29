/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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
 *      Johannes Fürmann <johannes@weltraumpflege.org>
 */

// Package spaceapibee is a Bee that can query a spaceapi server.
package spaceapibee

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/muesli/beehive/bees"
)

// SpaceAPIBee is a Bee that can query a spaceapi server.
type SpaceAPIBee struct {
	bees.Bee

	url string

	evchan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *SpaceAPIBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "get_status":
		type SpaceAPIResult struct {
			State struct {
				Open bool `json:"open"`
			} `json:"state"`
		}
		apiState := new(SpaceAPIResult)

		// get json data
		resp, err := http.Get(mod.url)
		if err != nil {
			log.Println("Error: SpaceAPI instance @ " + mod.url + " not reachable")
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, apiState)
			if err != nil {
				log.Println("Sorry, couldn't unmarshal the JSON data from SpaceAPI Instance @ " + mod.url)
				apiState.State.Open = false
			}
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "query_result",
			Options: []bees.Placeholder{
				{
					Name:  "open",
					Type:  "bool",
					Value: apiState.State.Open,
				},
			},
		}
		mod.evchan <- ev

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *SpaceAPIBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *SpaceAPIBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("url", &mod.url)
}
