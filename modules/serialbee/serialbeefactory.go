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

package serialbee

import (
	"github.com/muesli/beehive/modules"
)

type SerialBeeFactory struct {
	modules.ModuleFactory
}

// Interface impl

func (factory *SerialBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := SerialBee{
		Module: modules.NewBee(name, factory.Name(), description),
		device:   options.GetValue("device").(string),
		baudrate: int(options.GetValue("baudrate").(float64)),
	}

	return &bee
}

func (factory *SerialBeeFactory) Name() string {
	return "serialbee"
}

func (factory *SerialBeeFactory) Description() string {
	return "A bee that talks serially"
}

func (factory *SerialBeeFactory) Options() []modules.BeeOptionDescriptor {
	opts := []modules.BeeOptionDescriptor{
		modules.BeeOptionDescriptor{
			Name:        "device",
			Description: "Serial device to use",
			Type:        "string",
			Mandatory:   true,
		},
		modules.BeeOptionDescriptor{
			Name:        "baudrate",
			Description: "The baudrate you want to use",
			Type:        "int",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *SerialBeeFactory) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over the serial port",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "port",
					Description: "The port the message was received on",
					Type:        "int",
				},
			},
		},
	}
	return events
}

func (factory *SerialBeeFactory) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message via the serial port",
			Options: []modules.PlaceholderDescriptor{
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
	f := SerialBeeFactory{}
	modules.RegisterFactory(&f)
}
