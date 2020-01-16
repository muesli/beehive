/*
 *    Copyright (C) 2019 Adaptant Solutions AG
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
 *      Paul Mundt <paul.mundt@adaptant.io>
 */

package knowgobee

import (
	"github.com/muesli/beehive/bees"
)

// KnowGoBeeFactory is a factory for KnowGoBees.
type KnowGoBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *KnowGoBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := KnowGoBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *KnowGoBeeFactory) ID() string {
	return "knowgobee"
}

// Name returns the name of this Bee.
func (factory *KnowGoBeeFactory) Name() string {
	return "KnowGo"
}

// Description returns the description of this Bee.
func (factory *KnowGoBeeFactory) Description() string {
	return "React to various KnowGo platform events"
}

// Image returns the filename of an image for this Bee.
func (factory *KnowGoBeeFactory) Image() string {
	return factory.ID() + ".jpg"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *KnowGoBeeFactory) LogoColor() string {
	return "#7ace56"
}

// Options returns the options available to configure this Bee.
func (factory *KnowGoBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "server",
			Description: "URL for the desired KnowGo server",
			Type:        "string",
		},
		{
			Name:        "polling_interval",
			Description: "Number of seconds to wait between polling intervals",
			Type:        "int",
			Default:     10,
		},
		{
			Name:        "api_key",
			Description: "API Key for the KnowGo server",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *KnowGoBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "country-change",
			Description: "A country change event has taken place",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "driverId",
					Description: "ID of Driver",
					Type:        "int",
				},
				{
					Name:        "entering",
					Description: "ISO 3166-2 country code of country entered",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "leaving",
					Description: "ISO 3166-2 country code of country left",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "timestamp",
					Description: "Time when event took place",
					Type:        "timestamp",
					Mandatory:   true,
				},
			},
		},
	}
	return events
}

func init() {
	f := KnowGoBeeFactory{}
	bees.RegisterFactory(&f)
}
