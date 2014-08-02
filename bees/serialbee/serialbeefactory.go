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
	"github.com/muesli/beehive/bees"
)

type SerialBeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *SerialBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SerialBee{
		Bee: bees.NewBee(name, factory.Name(), description),
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

func (factory *SerialBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		bees.BeeOptionDescriptor{
			Name:        "device",
			Description: "Serial device to use",
			Type:        "string",
			Mandatory:   true,
		},
		bees.BeeOptionDescriptor{
			Name:        "baudrate",
			Description: "The baudrate you want to use",
			Type:        "int",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *SerialBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		bees.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over the serial port",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				bees.PlaceholderDescriptor{
					Name:        "port",
					Description: "The port the message was received on",
					Type:        "int",
				},
			},
		},
	}
	return events
}

func (factory *SerialBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		bees.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message via the serial port",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
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
	bees.RegisterFactory(&f)
}
