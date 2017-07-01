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

package tinderbee

import (
	"github.com/muesli/beehive/bees"
)

// TinderBeeFactory is a factory for TinderBees.
type TinderBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TinderBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TinderBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *TinderBeeFactory) ID() string {
	return "tinderbee"
}

// Name returns the name of this Bee.
func (factory *TinderBeeFactory) Name() string {
	return "Tinder"
}

// Description returns the description of this Bee.
func (factory *TinderBeeFactory) Description() string {
	return "Interacts with Tinder"
}

// Image returns the filename of an image for this Bee.
func (factory *TinderBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *TinderBeeFactory) LogoColor() string {
	return "#35465c"
}

// Options returns the options available to configure this Bee.
func (factory *TinderBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "userID",
			Description: "Your facebook user ID",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "userToken",
			Description: "Your facebook user token",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *TinderBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "update",
			Description: "is triggered after fetching updates",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "ID",
					Description: "Update ID",
					Type:        "string",
				},
				{
					Name:        "common_friend_count",
					Description: "Common friend count",
					Type:        "int",
				},
				{
					Name:        "common_like_count",
					Description: "Common like count:",
					Type:        "int",
				},
				{
					Name:        "message_count",
					Description: "Message count",
					Type:        "int",
				},
				{
					Name:        "person_ID",
					Description: "Tinder User ID",
					Type:        "string",
				},
				{
					Name:        "person_bio",
					Description: "persons biography on tinder",
					Type:        "string",
				},
				{
					Name:        "person_birth",
					Description: "Persons birth date",
					Type:        "string",
				},
				{
					Name:        "person_name",
					Description: "name from the person",
					Type:        "string",
				},
				{
					Name:        "person_ping_time",
					Description: "Persons ping time",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "user",
			Description: "returns user informations",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "status",
					Description: "user status",
					Type:        "int",
				},
				{
					Name:        "connection_count",
					Description: "users connection count",
					Type:        "int",
				},
				{
					Name:        "common_like_count",
					Description: "common like count",
					Type:        "int",
				},
				{
					Name:        "common_friend_count",
					Description: "common friend count",
					Type:        "int",
				},
				{
					Name:        "common_likes",
					Description: "common likes",
					Type:        "[]string",
				},
				{
					Name:        "common_interests",
					Description: "common interests",
					Type:        "[]string",
				},
				{
					Name:        "uncommon_interests",
					Description: "common interests",
					Type:        "[]string",
				},
				{
					Name:        "common_friends",
					Description: "common friends",
					Type:        "[]string",
				},
				{
					Name:        "ID",
					Description: "user ID",
					Type:        "string",
				},
				{
					Name:        "bio",
					Description: "users biography on tinder",
					Type:        "string",
				},
				{
					Name:        "birthdate",
					Description: "users birth date",
					Type:        "string",
				},
				{
					Name:        "gender",
					Description: "users gender",
					Type:        "int",
				},
				{
					Name:        "name",
					Description: "username",
					Type:        "string",
				},
				{
					Name:        "ping_time",
					Description: "ping time",
					Type:        "string",
				},
				{
					Name:        "distance_miles",
					Description: "distance ins miles",
					Type:        "int",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "recommendation",
			Description: "is triggered after you've fetch recommendations",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "ID",
					Description: "recommendations user ID",
					Type:        "string",
				},
				{
					Name:        "bio",
					Description: "recommendations biography",
					Type:        "string",
				},
				{
					Name:        "birth",
					Description: "recommendations birth date",
					Type:        "string",
				},
				{
					Name:        "gender",
					Description: "recommendations gender",
					Type:        "int",
				},
				{
					Name:        "name",
					Description: "recommendations name",
					Type:        "string",
				},
				{
					Name:        "distance_miles",
					Description: "recommendations distance in miles",
					Type:        "int",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "message_sent",
			Description: "is triggered after you've successfully sent a message",
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *TinderBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "get_updates",
			Description: "Fetches your tinder updates",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "limit",
					Description: "limit of updates to fetch",
					Type:        "int",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "get_user",
			Description: "Fetches user information",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "userID",
					Description: "User id of the user to search",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "send_message",
			Description: "Sends messag eto desired userID",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "userID",
					Description: "User id from the recipient",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "test",
					Description: "Contents of the message",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := TinderBeeFactory{}
	bees.RegisterFactory(&f)
}
