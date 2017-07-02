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
 *      Nicolas Martin <penguwingit@gmail.com>
 */

package cleverbotbee

import (
	"github.com/muesli/beehive/bees"
)

// CleverbotBeeFactory is a factory for CleverbotBees.
type CleverbotBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *CleverbotBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := CleverbotBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *CleverbotBeeFactory) ID() string {
	return "cleverbotbee"
}

// Name returns the name of this Bee.
func (factory *CleverbotBeeFactory) Name() string {
	return "Cleverbot"
}

// Description returns the description of this Bee.
func (factory *CleverbotBeeFactory) Description() string {
	return "Chat with cleverbot"
}

// Image returns the filename of an image for this Bee.
func (factory *CleverbotBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *CleverbotBeeFactory) LogoColor() string {
	return "#33f3ff"
}

// Options returns the options available to configure this Bee.
func (factory *CleverbotBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "api_user",
			Description: "Your cleverbot api username",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "api_key",
			Description: "Your cleverbot api key",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "session_nick",
			Description: "Optionally set a nickname for the session",
			Type:        "string",
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *CleverbotBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "answer",
			Description: "is triggerd when you receive a message/answer from cleverbot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "answer",
					Description: "cleverbots message",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *CleverbotBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send_message",
			Description: "Send a message to the cleverbot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "Contents of the message/question",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := CleverbotBeeFactory{}
	bees.RegisterFactory(&f)
}
