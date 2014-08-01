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

package efabee

import (
	"github.com/muesli/beehive/modules"
)

type EFABeeFactory struct {
	modules.ModuleFactory
}

// Interface impl

func (factory *EFABeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := EFABee{
		Module: modules.NewBee(name, factory.Name(), description),
		baseURL:      options.GetValue("baseurl").(string),
	}

	return &bee
}

func (factory *EFABeeFactory) Name() string {
	return "efabee"
}

func (factory *EFABeeFactory) Description() string {
	return "An EFA module for beehive"
}

func (factory *EFABeeFactory) Options() []modules.BeeOptionDescriptor {
	opts := []modules.BeeOptionDescriptor{
		modules.BeeOptionDescriptor{
			Name:        "baseurl",
			Description: "Base-url of the EFA API, e.g.: http://efa.avv-augsburg.de/avv/",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *EFABeeFactory) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "departure",
			Description: "Departure for a stop has been retrieved",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "stop",
					Description: "Which stop the departures are for",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "eta",
					Description: "Expected time for arrival in minutes",
					Type:        "int",
				},
				modules.PlaceholderDescriptor{
					Name:        "route",
					Description: "Route number",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "destination",
					Description: "Destination",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (factory *EFABeeFactory) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "departures",
			Description: "Retrieves next departures from a specific stop",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
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
	modules.RegisterFactory(&f)
}
