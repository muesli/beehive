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
 */

// Package httpbee is a Bee that lets you trigger HTTP requests.
package httpbee

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/muesli/beehive/bees"
)

// HTTPBee is a Bee that lets you trigger HTTP requests.
type HTTPBee struct {
	bees.Bee

	addr string
	path string

	eventChan chan bees.Event
}

// Run executes the Bee's event loop.
func (mod *HTTPBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	select {
	case <-mod.SigChan:
		return
	}
}

// Action triggers the action passed to it.
func (mod *HTTPBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	url := ""
	action.Options.Bind("url", &url)

	switch action.Name {
	case "get":
		resp, err := http.Get(url)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}

		ev, err := mod.prepareResponseEvent(b)
		if err == nil {
			ev.Name = "get"
			ev.Options.SetValue("url", "url", url)
			mod.eventChan <- ev
		}

	case "post":
		j := ""
		action.Options.Bind("json", &j)

		buf := strings.NewReader(j)
		resp, err := http.Post(url, "application/json", buf)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}

		ev, err := mod.prepareResponseEvent(b)
		if err == nil {
			ev.Name = "post"
			ev.Options.SetValue("url", "url", url)
			mod.eventChan <- ev
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *HTTPBee) prepareResponseEvent(resp []byte) (bees.Event, error) {
	ev := bees.Event{
		Bee: mod.Name(),
		Options: []bees.Placeholder{
			{
				Name:  "data",
				Type:  "string",
				Value: string(resp),
			},
		},
	}

	var payload interface{}
	err := json.Unmarshal(resp, &payload)
	if err == nil {
		ev.Options.SetValue("json", "map", payload)

		// this is a bit of a funny hack:
		// each parameter in the JSON response will get mapped to an (undocumented) event parameter
		// TODO: decide if this a good idea or not (probably not)
		j := make(map[string]interface{})
		err = json.Unmarshal(resp, &j)

		if err == nil {
			for k, v := range j {
				mod.Logf("JSON param: %s = %+v\n", k, v)
				if k == "json" || k == "data" {
					continue
				}

				// FIXME: hard-coded 'string'
				ev.Options.SetValue(k, "string", v)
			}
		}
	}

	return ev, nil
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *HTTPBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}
