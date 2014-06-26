/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

package hellobee

import (
	"github.com/muesli/beehive/modules"
)

type HelloBeeFactory struct {
}

// Interface impl

func (factory *HelloBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := HelloBee{
		name: name,
		namespace: factory.Name(),
		description: description,
	}
	return &bee
}

func (factory *HelloBeeFactory) Name() string {
	return "hellobee"
}

func (factory *HelloBeeFactory) Description() string {
	return "A 'Hello World' module for beehive"
}

func (factory *HelloBeeFactory) Options() []modules.BeeOptionDescriptor {
	return []modules.BeeOptionDescriptor{}
}

func (factory *HelloBeeFactory) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{}
	return events
}

func (factory *HelloBeeFactory) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{}
	return actions
}

func init() {
	f := HelloBeeFactory{}
	modules.RegisterFactory(&f)
}
