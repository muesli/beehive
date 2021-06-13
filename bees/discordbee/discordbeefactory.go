/*
 *    Copyright (C) 2019      CalmBit
 *                  2014-2019 Christian Muehlhaeuser
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
 *      CalmBit <calmbit@posto.net>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package discordbee

import (
	"github.com/muesli/beehive/bees"
)

// DiscordBeeFactory is a factory for DiscordBees.
type DiscordBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *DiscordBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := DiscordBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *DiscordBeeFactory) ID() string {
	return "discordbee"
}

// Name returns the name of this Bee.
func (factory *DiscordBeeFactory) Name() string {
	return "Discord"
}

// Description returns the description of this Bee.
func (factory *DiscordBeeFactory) Description() string {
	return "Connects to Discord as a bot!"
}

// Image returns the filename of an image for this Bee.
func (factory *DiscordBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *DiscordBeeFactory) LogoColor() string {
	return "#7289DA"
}

// Options returns the options available to configure this Bee.
func (factory *DiscordBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "api_token",
			Description: "The Discord API token for your bot",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *DiscordBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "is triggered when a message has been received",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "contents",
					Description: "Contents of the message",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "Name of the user who sent the message",
					Type:        "string",
				},
				{
					Name:        "channel_id",
					Description: "ID of the channel the message was sent in",
					Type:        "string",
				},
				{
					Name:        "channel_name",
					Description: "Name of the channel the message was sent in",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *DiscordBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a general message to the specified channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "contents",
					Description: "Contents of the message",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "channel_id",
					Description: "ID of the channel to post in",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "send_news",
			Description: "Sends a message to a news channel and publish it",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "contents",
					Description: "Contents of the message",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "channel_id",
					Description: "ID of the channel to post in",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "set_status",
			Description: "Sets the status of the discord bot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "status",
					Description: "Playing {status}",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := DiscordBeeFactory{}
	bees.RegisterFactory(&f)
}
