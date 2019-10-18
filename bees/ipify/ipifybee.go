/*
 *    Copyright (C) 2019 Sergio Rubio
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
 *      Sergio Rubio <sergio@rubio.im>
 */

package ipify

import (
	"time"

	"github.com/muesli/beehive/bees"
	"github.com/rdegges/go-ipify"
)

type IpifyBee struct {
	bees.Bee
	interval int
}

func (mod *IpifyBee) getIP(oldIP string, eventChan chan bees.Event) string {
	ip, err := ipify.GetIp()
	if err != nil {
		ip, err = ipify.GetIp()
		if err != nil {
			panic(err)
		}
	}

	if oldIP != ip {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "ip",
			Options: []bees.Placeholder{
				{
					Name:  "ip",
					Type:  "string",
					Value: ip,
				},
			},
		}
		eventChan <- ev
		return ip
	}

	return oldIP
}

// Run executes the Bee's event loop.
func (mod *IpifyBee) Run(eventChan chan bees.Event) {
	// protects us against a user setting the wrong value here
	if mod.interval < 1 {
		mod.interval = defaultUpdateInterval
	}

	oldIP := mod.getIP("", eventChan)

	for {
		select {
		case <-mod.SigChan:
			return
		case <-time.After(time.Duration(mod.interval) * time.Minute):
			mod.LogDebugf("Retrieving public IP from ipify.com")
			oldIP = mod.getIP(oldIP, eventChan)
		}
	}
}

// Action triggers the action passed to it.
func (mod *IpifyBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *IpifyBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	options.Bind("interval", &mod.interval)
}
