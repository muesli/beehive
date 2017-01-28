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

type SpaceApiBeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *SpaceApiBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SpaceApiBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

func (factory *SpaceApiBeeFactory) Name() string {
	return "spaceapibee"
}

func (factory *SpaceApiBeeFactory) Description() string {
	return "A bee that echoes the status of a SpaceAPI instance"
}

func (factory *SpaceApiBeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *SpaceApiBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "url",
			Description: "URL to the SpaceAPI endpoint",
			Type:        "string",
		},
	}
	return opts
}

func (factory *SpaceApiBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "get_status",
			Description: "Gets the Status of a LabAPI instance",
			Options:     []bees.PlaceholderDescriptor{},
		},
	}
	return actions
}

func (factory *SpaceApiBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "query_result",
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
	f := SpaceApiBeeFactory{}
	bees.RegisterFactory(&f)
}
