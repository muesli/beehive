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
 *      Martin Schlierf <go@koma666.de>
 */

package mumblebee

import (
	"github.com/muesli/beehive/bees"
)

// MumbleBeeFactory is a factory for MumbleBees.
type MumbleBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *MumbleBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := MumbleBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *MumbleBeeFactory) ID() string {
	return "mumblebee"
}

// Name returns the name of this Bee.
func (factory *MumbleBeeFactory) Name() string {
	return "Mumble"
}

// Description returns the description of this Bee.
func (factory *MumbleBeeFactory) Description() string {
	return "Connects to Mumble"
}

// Image returns the filename of an image for this Bee.
func (factory *MumbleBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *MumbleBeeFactory) LogoColor() string {
	return "#f42929"
}

// Options returns the options available to configure this Bee.
func (factory *MumbleBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "address",
			Description: "Hostname of Mumble server, eg: de.clanwarz.com:64738",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "channel",
			Description: "Channel to join on the Mumble server, blank for root",
			Type:        "string",
			Mandatory:   false,
		},
		{
			Name:        "user",
			Description: "Username to authenticate with Mumble server",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password to use to connect to Mumble server",
			Type:        "password",
			Mandatory:   false,
		},
		{
			Name:        "insecure",
			Description: "Do not check server certificate",
			Default:     true,
			Type:        "bool",
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *MumbleBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over Mumble",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				{
					Name:        "user",
					Description: "The user that sent the message",
					Type:        "string",
				},
				{
					Name:        "channels",
					Description: "The channels where the message was received",
					Type:        "[]string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "user_connected",
			Description: "A user connected to the mumble server",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "The name of the user",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "user_disconnected",
			Description: "A user disconnected from the mumble server",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "The name of the user",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "user_registered",
			Description: "A user registered to mumble server",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "The name of the user",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "user_kicked",
			Description: "An admin kicked a user from the mumble server",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "The name of the user",
					Type:        "string",
				},
				{
					Name:        "admin",
					Description: "The name of the admin",
					Type:        "string",
				},
				{
					Name:        "message",
					Description: "The message from the admin",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "user_banned",
			Description: "An admin banned a user from connecting to the mumble server",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "The name of the user",
					Type:        "string",
				},
				{
					Name:        "admin",
					Description: "The name of the admin",
					Type:        "string",
				},
				{
					Name:        "message",
					Description: "The message from the admin",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "user_changed_channel",
			Description: "An user changed the channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "The name of the user",
					Type:        "string",
				},
				{
					Name:        "channel",
					Description: "The name of the channel the user moved to",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *MumbleBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "move",
			Description: "Moves to a new Channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "channel",
					Description: "The channel to move in",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "send_to_user",
			Description: "Sends a message to a user",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "Which user to send the message to",
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
		{
			Namespace:   factory.Name(),
			Name:        "send_to_channel",
			Description: "Sends a message to a channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "channel",
					Description: "Which channel to send the message to",
					Type:        "string",
					Mandatory:   false,
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
	f := MumbleBeeFactory{}
	bees.RegisterFactory(&f)
}
