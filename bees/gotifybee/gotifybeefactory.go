/*
 *    Copyright (C) 2020      deranjer
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
 *      deranjer <deranjer@gmail.com>
 */

// Package gotifybee is able to send notifications on Gotify.
package gotifybee

import (
	"github.com/muesli/beehive/bees"
)

// GotifyBeeFactory is a factory for GotifyBees.
type GotifyBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *GotifyBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := GotifyBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *GotifyBeeFactory) ID() string {
	return "gotifybee"
}

// Name returns the name of this Bee.
func (factory *GotifyBeeFactory) Name() string {
	return "Gotify"
}

// Description returns the description of this Bee.
func (factory *GotifyBeeFactory) Description() string {
	return "Lets you push notifications on Gotify"
}

// Image returns the filename of an image for this Bee.
func (factory *GotifyBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *GotifyBeeFactory) LogoColor() string {
	return "#448CCB"
}

// Options returns the options available to configure this Bee.
func (factory *GotifyBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "token",
			Description: "The gotify token for the Application to send messages",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name: "serverURL",
			Description: "The URL to the gotify server, eg: http://gotify.com/ (trailing slash required!)",
			Type: "string",
			Mandatory: true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *GotifyBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *GotifyBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "title",
					Description: "Title of the message",
					Type:        "string",
				},
				{
					Name:        "message",
					Description: "Content of the message",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "priority",
					Description: "Priority of the message, see https://github.com/gotify/android#message-priorities",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := GotifyBeeFactory{}
	bees.RegisterFactory(&f)
}
