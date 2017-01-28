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

package emailbee

import (
	"github.com/muesli/beehive/bees"
)

type EmailBeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *EmailBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := EmailBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

func (factory *EmailBeeFactory) Name() string {
	return "emailbee"
}

func (factory *EmailBeeFactory) Description() string {
	return "An email module for beehive"
}

func (factory *EmailBeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *EmailBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "username",
			Description: "Username used for SMTP auth",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password used for SMTP auth",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "server",
			Description: "Address of SMTP server, eg: smtp.myserver.com:587",
			Type:        "url",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *EmailBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

func (factory *EmailBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends an email",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "Content of the email",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := EmailBeeFactory{}
	bees.RegisterFactory(&f)
}
