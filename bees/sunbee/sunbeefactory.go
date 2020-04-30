/*
 *    Copyright (C) 2020 Sergio Rubio
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

package sunbee

import (
	"github.com/muesli/beehive/bees"
)

// SunBeeFactory is a factory for SunBees.
type SunBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *SunBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SunBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *SunBeeFactory) ID() string {
	return "sunbee"
}

// Name returns the name of this Bee.
func (factory *SunBeeFactory) Name() string {
	return "Sunset/Sunrise"
}

// Description returns the description of this Bee.
func (factory *SunBeeFactory) Description() string {
	return "Send an event when the Sun raises or goes down"
}

// Image returns the filename of an image for this Bee.
func (factory *SunBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// Options returns the options available to configure this Bee.
func (factory *SunBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "query",
			Description: "The name of the city where the sunrise/sunset will happen",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "offset",
			Description: "Fire the event this number of seconds before sunset/sunrise (default: 2 min)",
			Type:        "int",
			Default:     0,
			Mandatory:   false,
		},
	}
	return opts
}

func (factory *SunBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "sunrise",
			Description: "Sunrise is happening",
			Options:     []bees.PlaceholderDescriptor{},
		},
		{
			Namespace:   factory.Name(),
			Name:        "sunset",
			Description: "Sunset is happening",
			Options:     []bees.PlaceholderDescriptor{},
		},
	}
	return events
}

func init() {
	f := SunBeeFactory{}
	bees.RegisterFactory(&f)
}
