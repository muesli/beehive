/*
 *    Copyright (C) 2019 Sergio Rubio
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
 *      Sergio Rubio <sergio@rubio.im>
 */

package redisbee

import (
	"github.com/muesli/beehive/bees"
)

// RedisBeeFactory takes care of initializing RedisBee
type RedisBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *RedisBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := RedisBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *RedisBeeFactory) ID() string {
	return "redis"
}

// Name returns the name of this Bee.
func (factory *RedisBeeFactory) Name() string {
	return "redis"
}

// Description returns the description of this Bee.
func (factory *RedisBeeFactory) Description() string {
	return "PubSub and key/value storage using a Redis server"
}

// Image returns the asset name of this Bee (in the assets/bees folder)
func (factory *RedisBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// Events describes the available events provided by this Bee.
func (factory *RedisBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over the pubsub channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "message",
					Description: "The message that was received",
					Type:        "string",
				},
				{
					Name:        "channel",
					Description: "The channel the message was received in",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *RedisBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "set",
			Description: "Key/Value pair to store",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "key",
					Description: "Redis key",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "value",
					Description: "Redis value",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "publish",
			Description: "Publish to a Redis channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "message",
					Description: "Redis message",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

// Options returns the options available to configure this Bee.
func (factory *RedisBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "host",
			Description: "Redis host",
			Type:        "string",
			Default:     "localhost",
			Mandatory:   false,
		},
		{
			Name:        "port",
			Description: "Redis port",
			Type:        "string",
			Default:     "6379",
			Mandatory:   false,
		},
		{
			Name:        "password",
			Description: "Redis password",
			Type:        "string",
			Default:     "",
			Mandatory:   false,
		},
		{
			Name:        "db",
			Description: "Redis database",
			Type:        "int",
			Default:     0,
			Mandatory:   false,
		},
		{
			Name:        "channel",
			Description: "Redis PubSub Channel",
			Type:        "string",
			Default:     "",
			Mandatory:   false,
		},
	}
	return opts
}

func init() {
	f := RedisBeeFactory{}
	bees.RegisterFactory(&f)
}
