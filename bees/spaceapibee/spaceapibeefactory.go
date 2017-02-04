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
 *      Johannes Fürmann <johannes@weltraumpflege.org>
 */

package spaceapibee

import (
	"github.com/muesli/beehive/bees"
)

// SpaceAPIBeeFactory is a factory for SpaceAPIBees.
type SpaceAPIBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *SpaceAPIBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SpaceAPIBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *SpaceAPIBeeFactory) ID() string {
	return "spaceapibee"
}

// Name returns the name of this Bee.
func (factory *SpaceAPIBeeFactory) Name() string {
	return "SpaceAPI"
}

// Description returns the description of this Bee.
func (factory *SpaceAPIBeeFactory) Description() string {
	return "Reacts to SpaceAPI status changes"
}

// Image returns the filename of an image for this Bee.
func (factory *SpaceAPIBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *SpaceAPIBeeFactory) LogoColor() string {
	return "#edb112"
}

// Options returns the options available to configure this Bee.
func (factory *SpaceAPIBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "url",
			Description: "URL to the SpaceAPI endpoint",
			Type:        "string",
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *SpaceAPIBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "status",
			Description: "Gets the Status of a LabAPI instance",
			Options:     []bees.PlaceholderDescriptor{},
		},
	}
	return actions
}

// Events describes the available events provided by this Bee.
func (factory *SpaceAPIBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "result",
			Description: "is triggered as soon as the query has been executed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "open",
					Description: "open-state of the spaceapi instance that was queried",
					Type:        "bool",
				},
			},
		},
	}
	return events
}

func init() {
	f := SpaceAPIBeeFactory{}
	bees.RegisterFactory(&f)
}
