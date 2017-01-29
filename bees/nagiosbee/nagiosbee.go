/*
 *    Copyright (C) 2014      Daniel 'grindhold' Brendle
 *                  2014-2017 Christian Muehlhaeuser
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

/*
   Please note that, in order to run this bee on a nagios-server, you
   have to provide the nagios status-script found at

   https://github.com/lizell/php-nagios-json/blob/master/statusJson.php

   just drop this script in the htdocs-folder of your nagios-installation
   and change the variable $statusFile to where the status.dat-file of your
   installation resides
*/

// Package nagiosbee is a Bee that can interface with a Nagios instance.
package nagiosbee

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/muesli/beehive/bees"
)

// NagiosBee is a Bee that can interface with a Nagios instance.
type NagiosBee struct {
	bees.Bee

	url      string
	user     string
	password string
	services map[string]map[string]service // services[hostname][servicename]

	eventChan chan bees.Event
}

type report struct {
	Services map[string]map[string]service `json:"services"` // services[hostname][servicename]
}

type service struct {
	HostName           string `json:"host_name"`
	ServiceDescription string `json:"service_description"`
	CurrentState       string `json:"current_state"`
	LastHardState      string `json:"last_hard_state"`
	PluginOutput       string `json:"plugin_output"`
}

func (mod *NagiosBee) announceStatuschange(s service) {
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "status_change",
		Options: []bees.Placeholder{
			{
				Name:  "host",
				Type:  "string",
				Value: s.HostName,
			},
			{
				Name:  "service",
				Type:  "string",
				Value: s.ServiceDescription,
			},
			{
				Name:  "message",
				Type:  "string",
				Value: s.PluginOutput,
			},
			{
				Name:  "status",
				Type:  "string",
				Value: s.CurrentState,
			},
		},
	}
	mod.eventChan <- event
}

// Run executes the Bee's event loop.
func (mod *NagiosBee) Run(cin chan bees.Event) {
	mod.eventChan = cin
	for {
		select {
		case <-mod.SigChan:
			return

		default:
		}
		time.Sleep(10 * time.Second)

		request, err := http.NewRequest("GET", mod.url, nil)
		if err != nil {
			log.Println("Could not build request")
			break
		}
		request.SetBasicAuth(mod.user, mod.password)

		client := http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Println("Couldn't find status-JSON at " + mod.url)
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Could not read data from URL")
			continue
		}
		log.Println(string(body))
		rep := new(report)
		err = json.Unmarshal(body, &rep)
		if err != nil {
			log.Println("Failed to unmarshal JSON")
			continue
		}

		log.Println("Start crawling map", len(rep.Services))
		var oldService service
		var ok bool
		for hn, mp := range rep.Services {
			snmap := make(map[string]service)
			for sn, s := range mp {
				log.Println(s)
				if oldService, ok = mod.services[hn][sn]; !ok {
					log.Println("jedesmaldarein")
					mod.announceStatuschange(s)
				} else {
					if s.CurrentState != oldService.CurrentState {
						log.Println("statuschange")
						mod.announceStatuschange(s)
					}
				}
				if s.CurrentState != s.LastHardState {
					log.Println("hardstate_changed")
					//TODO: Evaluate if good enough
				}
				snmap[sn] = rep.Services[hn][sn]
			}
			mod.services[hn] = snmap
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *NagiosBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("url", &mod.url)
	options.Bind("user", &mod.user)
	options.Bind("password", &mod.password)
}
