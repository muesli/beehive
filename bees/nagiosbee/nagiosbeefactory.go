/*
 *    Copyright (C) 2014      Daniel 'grindhold' Brendle
 *                  2014-2017 Christian Muehlhaeuser
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package nagiosbee

import (
	"github.com/muesli/beehive/bees"
)

type NagiosBeeFactory struct {
	bees.BeeFactory
}

func (factory *NagiosBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := NagiosBee{
		Bee:      bees.NewBee(name, factory.Name(), description, options),
		url:      options.GetValue("url").(string),
		user:     options.GetValue("user").(string),
		password: options.GetValue("password").(string),
	}
	bee.services = make(map[string]map[string]service)
	return &bee
}

func (factory *NagiosBeeFactory) Name() string {
	return "nagiosbee"
}

func (factory *NagiosBeeFactory) Description() string {
	return "A bee that fetches status changes from nagios-monitors."
}

func (factory *NagiosBeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *NagiosBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "url",
			Description: "URL to the statusJson.php-script typically http://domain.com/nagios3/statusJson.php",
			Type:        "string",
		},
		{
			Name:        "user",
			Description: "The username of the nagios-user",
			Type:        "string",
		},
		{
			Name:        "password",
			Description: "Password of the nagios-user's account",
			Type:        "string",
		},
	}
	return opts
}

func (factory *NagiosBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "statuschange",
			Description: "The status of a Service has changed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "host",
					Description: "Name of the system the changed server resides on",
					Type:        "string",
				},
				{
					Name:        "service",
					Description: "Name of the service that has changed",
					Type:        "string",
				},
				{
					Name:        "message",
					Description: "Message that the NRPE-service returned",
					Type:        "string",
				},
				{
					Name:        "status",
					Description: "New status number",
					Type:        "int",
				},
			},
		},
	}
	return events
}

func init() {
	f := NagiosBeeFactory{}
	bees.RegisterFactory(&f)
}
