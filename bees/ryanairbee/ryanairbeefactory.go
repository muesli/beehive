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
 *      Nicolas Martin <penguwingit@gmail.com>
 */

package ryanairbee

import (
	"github.com/muesli/beehive/bees"
)

// RyanairBeeFactory is a factory for TumblrBees.
type RyanairBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *RyanairBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := RyanairBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *RyanairBeeFactory) ID() string {
	return "ryanairbee"
}

// Name returns the name of this Bee.
func (factory *RyanairBeeFactory) Name() string {
	return "Ryanair"
}

// Description returns the description of this Bee.
func (factory *RyanairBeeFactory) Description() string {
	return "interacts with the RyanairAPI"
}

// Image returns the filename of an image for this Bee.
func (factory *RyanairBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *RyanairBeeFactory) LogoColor() string {
	return "#35465c"
}

// Options returns the options available to configure this Bee.
func (factory *RyanairBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{}
	// No options needed :)
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *RyanairBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "flight_schedule",
			Description: "is triggered after fetching schedules",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "number",
					Description: "Flight number",
					Type:        "string",
				},
				{
					Name:        "departure_time",
					Description: "Departure Time",
					Type:        "string",
				},
				{
					Name:        "arrival_time",
					Description: "Arrival time",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *RyanairBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "get_schedules",
			Description: "Fetches flight schedules from the ryanairAPI",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "city_from",
					Description: "City flight starts",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "city_to",
					Description: "Arrival city",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "date",
					Description: "Date to fetch (Layout: 2023-04-20)",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := RyanairBeeFactory{}
	bees.RegisterFactory(&f)
}
