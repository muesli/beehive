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

package webbee

import (
	"github.com/muesli/beehive/modules"
)

type WebBeeFactory struct {
	modules.ModuleFactory
}

// Interface impl

func (factory *WebBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := WebBee{
		addr:        options.GetValue("addr").(string),
		path:        options.GetValue("path").(string),
	}

	bee.Module = modules.Module{name, factory.Name(), description}
	return &bee
}

func (factory *WebBeeFactory) Name() string {
	return "webbee"
}

func (factory *WebBeeFactory) Description() string {
	return "A RESTful HTTP module for beehive"
}

func (factory *WebBeeFactory) Options() []modules.BeeOptionDescriptor {
	opts := []modules.BeeOptionDescriptor{
		modules.BeeOptionDescriptor{
			Name:        "addr",
			Description: "Which addr to listen on, eg: 0.0.0.0:12345",
			Type:        "string",
			Mandatory:   true,
		},
		modules.BeeOptionDescriptor{
			Name:        "path",
			Description: "Which path to expect GET/POST requests on, eg: /foobar",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *WebBeeFactory) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "post",
			Description: "A POST call was received by the HTTP server",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "json",
					Description: "JSON map received from caller",
					Type:        "map",
				},
				modules.PlaceholderDescriptor{
					Name:        "ip",
					Description: "IP of the caller",
					Type:        "string",
				},
			},
		},
		modules.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "get",
			Description: "A GET call was received by the HTTP server",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "query_params",
					Description: "Map of query parameters received from caller",
					Type:        "map",
				},
				modules.PlaceholderDescriptor{
					Name:        "ip",
					Description: "IP of the caller",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := WebBeeFactory{}
	modules.RegisterFactory(&f)
}
