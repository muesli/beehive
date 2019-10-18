/*
 *    Copyright (C) 2016 Sergio Rubio
 *                  2017 Christian Muehlhaeuser
 *                  2019 David Schneider
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
 *      David Schneider <dsbrng25b@gmail.com>
 */

package rocketchatbee

import "github.com/muesli/beehive/bees"

// RocketchatBeeFactory is a factory for RocketchatBees.
type RocketchatBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *RocketchatBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := RocketchatBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *RocketchatBeeFactory) ID() string {
	return "rocketchatbee"
}

// Name returns the name of this Bee.
func (factory *RocketchatBeeFactory) Name() string {
	return "Rocket.Chat"
}

// Description returns the description of this Bee.
func (factory *RocketchatBeeFactory) Description() string {
	return "Connects to Rocket.Chat"
}

// Image returns the filename of an image for this Bee.
func (factory *RocketchatBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *RocketchatBeeFactory) LogoColor() string {
	return "#2c3e50"
}

// Options returns the options available to configure this Bee.
func (factory *RocketchatBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "url",
			Description: "Rocket.Chat URL",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "user_id",
			Description: "Rocket.Chat user id",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "auth_token",
			Description: "Rocket.Chat auth token",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *RocketchatBeeFactory) Actions() []bees.ActionDescriptor {
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
				{
					Name:        "alias",
					Description: "The name to show as sender. If empty, the name to which the user_id belongs to.",
					Type:        "string",
					Mandatory:   false,
				},
			},
		},
	}
	return actions
}

func init() {
	f := RocketchatBeeFactory{}
	bees.RegisterFactory(&f)
}
