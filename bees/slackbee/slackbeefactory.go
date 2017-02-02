/*
 *    Copyright (C) 2016 Sergio Rubio
 *                  2017 Christian Muehlhaeuser
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
 *      Sergio Rubio <rubiojr@frameos.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package slackbee

import "github.com/muesli/beehive/bees"

// SlackBeeFactory is a factory for SlackBees.
type SlackBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *SlackBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SlackBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// Name returns the name of this Bee.
func (factory *SlackBeeFactory) Name() string {
	return "slackbee"
}

// Description returns the description of this Bee.
func (factory *SlackBeeFactory) Description() string {
	return "A Slack module for beehive"
}

// Image returns the filename of an image for this Bee.
func (factory *SlackBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *SlackBeeFactory) LogoColor() string {
	return "#4b4b4b"
}

// Options returns the options available to configure this Bee.
func (factory *SlackBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "api_key",
			Description: "Slack API key",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "channels",
			Description: "Slack channels to listen on",
			Type:        "[]string",
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *SlackBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over Slack",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				{
					Name:        "channel",
					Description: "The channel the message was received in",
					Type:        "string",
				},
				{
					Name:        "user",
					Description: "The user that sent the message",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *SlackBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message to a channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "channel",
					Description: "Which channel to send the message to",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "text",
					Description: "Content of the message",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := SlackBeeFactory{}
	bees.RegisterFactory(&f)
}
