/*
 *    Copyright (C) 2014      Michael Wendland
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
 *      Michael Wendland <michiwend@michiwend.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package hellobee is an example for a Bee skeleton, designed to help you get
// started with writing your own Bees.
package hellobee

import (
	"github.com/muesli/beehive/bees"
)

type HelloBee struct {
	bees.Bee
}

func (mod *HelloBee) Run(eventChan chan bees.Event) {
	/*	ev := bees.Event{
			Bee: mod.Name(),
			Name:      "hello",
			Options:   []bees.Placeholder{},
		}

		eventChan <- ev*/
}

func (mod *HelloBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

func (mod *HelloBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}
