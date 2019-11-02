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

// Package watchdogbee implements a Systemd watchdog for Linux systems
package watchdogbee

import (
	"net/http"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/muesli/beehive/api"
	"github.com/muesli/beehive/bees"
)

// WatchdogBee struct
type WatchdogBee struct {
	bees.Bee
}

// Run executes the Bee's event loop.
func (mod *WatchdogBee) Run(eventChan chan bees.Event) {
	interval, err := daemon.SdWatchdogEnabled(false)
	// systemd's service unit interval in microseconds
	if err != nil || interval == 0 {
		return
	}
	interval = interval / 3 / 1000000000
	mod.LogDebugf("Watchdog interval: %d", interval)

	for {
		select {
		case <-mod.SigChan:
			return
		case <-time.After(time.Duration(interval) * time.Second):
			resp, err := http.Get(api.CanonicalURL().String())
			if err == nil {
				resp.Body.Close()
				mod.LogDebugf("Notify Systemd's watchdog")
				daemon.SdNotify(false, "WATCHDOG=1")
			}
		}
	}
}

// Action triggers the action passed to it.
func (mod *WatchdogBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *WatchdogBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}
