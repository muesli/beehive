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

package httpbee

import (
	"github.com/muesli/beehive/bees"
)

// HTTPBeeFactory is a factory for HTTPBees.
type HTTPBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *HTTPBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := HTTPBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *HTTPBeeFactory) ID() string {
	return "httpbee"
}

// Name returns the name of this Bee.
func (factory *HTTPBeeFactory) Name() string {
	return "HTTP Client"
}

// Description returns the description of this Bee.
func (factory *HTTPBeeFactory) Description() string {
	return "Lets you trigger HTTP requests"
}

// Image returns the filename of an image for this Bee.
func (factory *HTTPBeeFactory) Image() string {
	return "webbee.png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *HTTPBeeFactory) LogoColor() string {
	return "#223f5e"
}

// Events describes the available events provided by this Bee.
func (factory *HTTPBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "get",
			Description: "A GET request finished",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "url",
					Description: "URL of the request",
					Type:        "url",
				},
				{
					Name:        "data",
					Description: "Raw response",
					Type:        "string",
				},
				{
					Name:        "json",
					Description: "JSON response received",
					Type:        "map",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "post",
			Description: "A POST request finished",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "url",
					Description: "URL of the request",
					Type:        "url",
				},
				{
					Name:        "data",
					Description: "Raw response",
					Type:        "string",
				},
				{
					Name:        "json",
					Description: "JSON response received",
					Type:        "map",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *HTTPBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "get",
			Description: "Does a GET request",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "url",
					Description: "Where to connect to",
					Type:        "url",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "post",
			Description: "Does a POST request",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "json",
					Description: "Data to send",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "Where to connect to",
					Type:        "url",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := HTTPBeeFactory{}
	bees.RegisterFactory(&f)
}
