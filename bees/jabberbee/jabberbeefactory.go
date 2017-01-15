/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

type JabberBeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *JabberBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := JabberBee{
		Bee:      bees.NewBee(name, factory.Name(), description),
		server:   options.GetValue("server").(string),
		user:     options.GetValue("user").(string),
		password: options.GetValue("password").(string),
	}

	if options.GetValue("notls") != nil {
		bee.notls = options.GetValue("notls").(bool)
	}

	return &bee
}

func (factory *JabberBeeFactory) Name() string {
	return "jabberbee"
}

func (factory *JabberBeeFactory) Description() string {
	return "A Jabber module for beehive"
}

func (factory *JabberBeeFactory) Image() string {
	return factory.Name() + ".png"
}

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
				},
				{
					Name:        "text",
					Description: "Content of the message",
					Type:        "string",
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
