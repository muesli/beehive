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

// beehive's SpaceAPI module.
package spaceapibee

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/muesli/beehive/bees"
)

type SpaceApiBee struct {
	bees.Bee

	url string

	evchan chan bees.Event
}

func (mod *SpaceApiBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "get_status":
		type SpaceApiResult struct {
			State struct {
				Open bool `json:"open"`
			} `json:"state"`
		}
		api_state := new(SpaceApiResult)

		// get json data
		resp, err := http.Get(mod.url)
		if err != nil {
			log.Println("Error: SpaceAPI instance @ " + mod.url + " not reachable")
		} else {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, api_state)
			if err != nil {
				log.Println("Sorry, couldn't unmarshal the JSON data from SpaceAPI Instance @ " + mod.url)
				api_state.State.Open = false
			}
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "query_result",
			Options: []bees.Placeholder{
				{
					Name:  "open",
					Type:  "bool",
					Value: api_state.State.Open,
				},
			},
		}
		mod.evchan <- ev

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *SpaceApiBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan
}

func (mod *SpaceApiBee) ReloadOptions(options bees.BeeOptions) {
	//FIXME: implement this
	mod.SetOptions(options)
}
