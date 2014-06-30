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
	"github.com/muesli/beehive/modules"
)

type JabberBeeFactory struct {
	modules.ModuleFactory
}

// Interface impl

func (factory *JabberBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := JabberBee{
		server:      options.GetValue("server").(string),
		user:        options.GetValue("user").(string),
		password:    options.GetValue("password").(string),
		notls:       options.GetValue("notls").(bool),
	}

	bee.Module = modules.Module{name, factory.Name(), description}
	return &bee
}

func (factory *JabberBeeFactory) Name() string {
	return "jabberbee"
}

func (factory *JabberBeeFactory) Description() string {
	return "A Jabber module for beehive"
}

func (factory *JabberBeeFactory) Options() []modules.BeeOptionDescriptor {
	opts := []modules.BeeOptionDescriptor{
		modules.BeeOptionDescriptor{
			Name:        "server",
			Description: "Hostname of Jabber server, eg: talk.google.com:443",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "user",
			Description: "Username to authenticate with Jabber server",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "password",
			Description: "Password to use to connect to Jabber server",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "notls",
			Description: "Avoid using TLS for authentication",
			Type:        "bool",
		},
	}
	return opts
}

func (factory *JabberBeeFactory) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over Jabber",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "user",
					Description: "The user that sent the message",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (factory *JabberBeeFactory) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message to a remote",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "user",
					Description: "Which remote to send the message to",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
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
	modules.RegisterFactory(&f)
}
