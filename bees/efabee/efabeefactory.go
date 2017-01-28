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

package efabee

import (
	"github.com/muesli/beehive/bees"
)

type EFABeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *EFABeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := EFABee{
		Bee:      bees.NewBee(name, factory.Name(), description, options),
		Provider: options.GetValue("provider").(string),
	}

	return &bee
}

func (factory *EFABeeFactory) Name() string {
	return "efabee"
}

func (factory *EFABeeFactory) Description() string {
	return "An EFA module for beehive"
}

func (factory *EFABeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *EFABeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "provider",
			Description: "Provider-ID for the EFA API, e.g. 'avv'",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *EFABeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "departure",
			Description: "Departure for a stop has been retrieved",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "stop",
					Description: "Which stop the departures are for",
					Type:        "string",
				},
				{
					Name:        "eta",
					Description: "Expected time of arrival in minutes",
					Type:        "int",
				},
				{
					Name:        "etatime",
					Description: "Expected departure time",
					Type:        "string",
				},
				{
					Name:        "route",
					Description: "Route number",
					Type:        "string",
				},
				{
					Name:        "destination",
					Description: "Destination",
					Type:        "string",
				},
				{
					Name:        "mottype",
					Description: "Transportation type",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (factory *EFABeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "departures",
			Description: "Retrieves next departures from a specific stop",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "stop",
					Description: "The stop you want departures for",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := EFABeeFactory{}
	bees.RegisterFactory(&f)
}
