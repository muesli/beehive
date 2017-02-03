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

package hellobee

import (
	"github.com/muesli/beehive/bees"
)

// HelloBeeFactory is a factory for HelloBees.
type HelloBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *HelloBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := HelloBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *HelloBeeFactory) ID() string {
	return "hellobee"
}

// Name returns the name of this Bee.
func (factory *HelloBeeFactory) Name() string {
	return "Hello"
}

// Description returns the description of this Bee.
func (factory *HelloBeeFactory) Description() string {
	return "A 'Hello World' module for beehive"
}

func init() {
	f := HelloBeeFactory{}
	bees.RegisterFactory(&f)
}
