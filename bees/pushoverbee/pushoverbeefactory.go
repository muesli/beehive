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
 *      Raphael Mutschler <info@raphaelmutschler.de>
 */

// Package pushoverbee is a Bee that can send pushover notifications.
package pushoverbee

import (
	"github.com/muesli/beehive/bees"
)

// PushoverBeeFactory is a factory for PushoverBees.
type PushoverBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *PushoverBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := PushoverBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *PushoverBeeFactory) ID() string {
	return "pushoverbee"
}

// Name returns the name of this Bee.
func (factory *PushoverBeeFactory) Name() string {
	return "Pushover Notifications"
}

// Description returns the description of this Bee.
func (factory *PushoverBeeFactory) Description() string {
	return "Send pushover notifications"
}

// Image returns the filename of an image for this Bee.
func (factory *PushoverBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *PushoverBeeFactory) LogoColor() string {
	return "#2d2d2d"
}

// Options returns the options available to configure this Bee.
func (factory *PushoverBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "token",
			Description: "Pushover APP/API Token",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "user_token",
			Description: "Pushover User Token",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *PushoverBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *PushoverBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "message",
					Description: "The Message Body",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "title",
					Description: "The Title of the Message",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "url",
					Description: "An URL to link to",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "url_title",
					Description: "Alternative title for the url",
					Type:        "string",
					Mandatory:   false,
				},
			},
		},
	}
	return actions
}

func init() {
	f := PushoverBeeFactory{}
	bees.RegisterFactory(&f)
}
