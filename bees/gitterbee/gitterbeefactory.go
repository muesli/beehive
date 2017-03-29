/*
 *    Copyright (C) 2017 Christian Muehlhaeuser
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
 *      Nicolas Martin <penguwingithub@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package gitterbee is a Bee that can interface with Gitter
package gitterbee

import (
	"github.com/muesli/beehive/bees"
)

// GitterBeeFactory is a factory for GitterBees
type GitterBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *GitterBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := GitterBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *GitterBeeFactory) ID() string {
	return "gitterbee"
}

// Name returns the name of this Bee.
func (factory *GitterBeeFactory) Name() string {
	return "Gitter"
}

// Description returns the desciption of this Bee.
func (factory *GitterBeeFactory) Description() string {
	return "React and interact with Gitter"
}

// Image returns the filename of an image for this Bee.
func (factory *GitterBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns ther preferred logo background color (used by the admin interface).
func (factory *GitterBeeFactory) LogoColor() string {
	return "#994499"
}

// Options returns the options available to configure this Bee.
func (factory *GitterBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "accessToken",
			Description: "Your Gitter access token",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "room",
			Description: "Room to overwatch",
			Type:        "string",
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *GitterBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
		// 	Namespace:   factory.Name(),
		// 	Name:        "roomMessages",
		// 	Description: "Receive new room messages",
		// 	Options: []bees.PlaceholderDescriptor{
		// 		{
		// 			Name:        "ID",
		// 			Description: "ID of a message",
		// 			Type:        "string",
		// 		},
		// 		{
		// 			Name:        "text",
		// 			Description: "Text contents of a message",
		// 			Type:        "string",
		// 		},
		// 		{
		// 			Name:        "username",
		// 			Description: "Username",
		// 			Type:        "string",
		// 		},
		// 		{
		// 			Name:        "readBy",
		// 			Description: "Number of users who have read the message",
		// 			Type:        "int",
		// 		},
		// 	},
		// },
		// {
		// 	Namespace:   factory.Name(),
		// 	Name:        "mention",
		// 	Description: "Mentions inside of a message",
		// 	Options: []bees.PlaceholderDescriptor{
		// 		{
		// 			Name:        "mention",
		// 			Description: "Username of the user who has been mentioned",
		// 			Type:        "string",
		// 		},
		// 	},
		// },
		// {
		// 	Namespace:   factory.Name(),
		// 	Name:        "Issue",
		// 	Description: "Issue referenced in a message",
		// 	Options: []bees.PlaceholderDescriptor{
		// 		{
		// 			Name:        "issue",
		// 			Description: "Number of the issue",
		// 			Type:        "int",
		// 		},
		// 	},
		},
	}

	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *GitterBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "sendMessage",
			Description: "Sends a message into a room",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "room",
					Description: "Message room",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "message",
					Description: "Message text",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "joinRoom",
			Description: "Join a room on gitter",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "room",
					Description: "Room to join",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "leaveRoom",
			Description: "Leave a room on gitter",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "room",
					Description: "Room to leave",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		// {
		// 	Namespace:   factory.Name(),
		// 	Name:        "getRoomMessages",
		// 	Description: "Get a rooms messages",
		// 	Options: []bees.PlaceholderDescriptor{
		// 		{
		// 			Name:        "room",
		// 			Description: "Room to check for messages",
		// 			Type:        "string",
		// 			Mandatory:   true,
		// 		},
		// 	},
		// },
	}

	return actions
}

func init() {
	f := GitterBeeFactory{}
	bees.RegisterFactory(&f)
}
