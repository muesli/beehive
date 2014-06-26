/*
 *    Copyright (C) 2014 Michael Wendland
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
 */

// hellobee dummy module for beehive
package hellobee

import (
	"github.com/muesli/beehive/modules"
)

type HelloBee struct {
	name        string
	namespace   string
	description string
}

func (mod *HelloBee) Name() string {
	return mod.name
}

func (mod *HelloBee) Namespace() string {
	return mod.namespace
}

func (mod *HelloBee) Description() string {
	return mod.description
}

func (mod *HelloBee) Run(eventChan chan modules.Event) {
	/*	ev := modules.Event{
			Bee: mod.Name(),
			Name:      "hello",
			Options:   []modules.Placeholder{},
		}

		eventChan <- ev*/
}

func (mod *HelloBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}
