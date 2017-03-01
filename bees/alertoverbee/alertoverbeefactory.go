/*
 *    Copyright (C) 2017      Henson Lu
 *                  2015-2017 Christian Muehlhaeuser
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
 *      Henson Lu <henson.lu@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package alertoverbee is able to send notifications on AlertOver.
package alertoverbee

import (
	"github.com/muesli/beehive/bees"
)

// AlertOverBeeFactory is a factory for AlertOverBees.
type AlertOverBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *AlertOverBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := AlertOverBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *AlertOverBeeFactory) ID() string {
	return "alertoverbee"
}

// Name returns the name of this Bee.
func (factory *AlertOverBeeFactory) Name() string {
	return "AlertOver"
}

// Description returns the description of this Bee.
func (factory *AlertOverBeeFactory) Description() string {
	return "Lets you push notifications on AlertOver"
}

// Image returns the filename of an image for this Bee.
func (factory *AlertOverBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *AlertOverBeeFactory) LogoColor() string {
	return "#364150"
}

// Options returns the options available to configure this Bee.
func (factory *AlertOverBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "source",
			Description: "The alertover source code",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *AlertOverBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *AlertOverBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "receiver",
					Description: "The alertover receiver code",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "title",
					Description: "Title of the message",
					Type:        "string",
				},
				{
					Name:        "content",
					Description: "Content of the message",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "url",
					Description: "Add the url address in message",
					Type:        "string",
				},
				{
					Name:        "priority",
					Description: "Priority of the message, 0 for normal and 1 for emergency",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := AlertOverBeeFactory{}
	bees.RegisterFactory(&f)
}
