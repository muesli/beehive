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

package watchdogbee

import (
	"github.com/muesli/beehive/bees"
)

// WatchdogBeeFactory is a factory for Watchdog bees.
type WatchdogBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *WatchdogBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := WatchdogBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}

	return &bee
}

// ID returns the ID of this Bee.
func (factory *WatchdogBeeFactory) ID() string {
	return "watchdogbee"
}

// Name returns the name of this Bee.
func (factory *WatchdogBeeFactory) Name() string {
	return "Systemd Watchdog"
}

// Description returns the description of this Bee.
func (factory *WatchdogBeeFactory) Description() string {
	return "Notifies Systemd's watchdog"
}

// Image returns the asset name of this Bee (in the assets/bees folder)
func (factory *WatchdogBeeFactory) Image() string {
	return factory.ID() + ".png"
}

func init() {
	f := WatchdogBeeFactory{}
	bees.RegisterFactory(&f)
}
