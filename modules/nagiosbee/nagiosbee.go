/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
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
 */

/*
   Please note that, in order to run this bee on a nagios-server, you
   have to provide the nagios status-script found at

   https://github.com/lizell/php-nagios-json/blob/master/statusJson.php

   just drop this script in the htdocs-folder of your nagios-installation
   and change the variable $statusFile to where the status.dat-file of your
   installation resides
*/

package nagiosbee

import (
	"encoding/json"
	"fmt"
	"github.com/muesli/beehive/modules"
	"io/ioutil"
	"net/http"
)

type NagiosBee struct {
    modules.Module
	url      string
	user     string
	password string
	services map[string]map[string]service // services[hostname][servicename]

	eventChan chan modules.Event
}

type report struct {
	services map[string]map[string]service // services[hostname][servicename]
}

type service struct {
	host_name           string
	service_description string
	current_state       string
	last_hard_state     string
	plugin_output       string
}


func (mod *NagiosBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}

func (mod *NagiosBee) Run(cin chan modules.Event) {
	for {
		resp, _ := http.Get(mod.url)
		body, _ := ioutil.ReadAll(resp.Body)
		rep := new(report)
		json.Unmarshal(body, report{})

		var oldService service
		for hn, mp := range rep.services {
			for sn, s := range mp {
				oldService = mod.services[hn][sn]

				if s.current_state != oldService.current_state {
					fmt.Println("statuschange")
                    event := modules.Event{
                        Bee: mod.Name(),
                        Name: "statuschange",
                        Options: []modules.Placeholder {
                            modules.Placeholder{
                                Name:   "host",
                                Type:   "string",
                                Value:  s.host_name,
                            },
                            modules.Placeholder{
                                Name:   "service",
                                Type:   "string",
                                Value:  s.service_description,
                            },
                            modules.Placeholder{
                                Name:   "message",
                                Type:   "string",
                                Value:  s.plugin_output,
                            },
                            modules.Placeholder{
                                Name:  "status",
                                Type:   "string",
                                Value:  s.current_state,
                            },
                        },
                    }
                    mod.eventChan <- event
				}
				if s.current_state != s.last_hard_state {
					fmt.Println("hardstate_changed")
					//TODO: Evaluate if good enough
				}
				mod.services[hn][sn] = rep.services[hn][sn]
			}
		}
	}
	return
}
