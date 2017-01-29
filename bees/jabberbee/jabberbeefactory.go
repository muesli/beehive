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

package jabberbee

import (
	"github.com/muesli/beehive/bees"
)

// JabberBeeFactory is a factory for JabberBees.
type JabberBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *JabberBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := JabberBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// Name returns the name of this Bee.
func (factory *JabberBeeFactory) Name() string {
	return "jabberbee"
}

// Description returns the description of this Bee.
func (factory *JabberBeeFactory) Description() string {
	return "A Jabber module for beehive"
}

// Image returns the filename of an image for this Bee.
func (factory *JabberBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// Options returns the options available to configure this Bee.
func (factory *JabberBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "server",
			Description: "Hostname of Jabber server, eg: talk.google.com:443",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "user",
			Description: "Username to authenticate with Jabber server",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password to use to connect to Jabber server",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "notls",
			Description: "Avoid using TLS for authentication",
			Type:        "bool",
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *JabberBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over Jabber",
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
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *JabberBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message to a remote",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user",
					Description: "Which remote to send the message to",
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
	f := JabberBeeFactory{}
	bees.RegisterFactory(&f)
}
