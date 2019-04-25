// +build !darwin

// FIXME: This bee doesn't build on macOS
// XXX: https://gist.github.com/prologic/b5fa148410a26f917b4d944b72d847c0

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

package serialbee

import (
	"github.com/muesli/beehive/bees"
)

// SerialBeeFactory is a factory for SerialBees.
type SerialBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *SerialBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SerialBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *SerialBeeFactory) ID() string {
	return "serialbee"
}

// Name returns the name of this Bee.
func (factory *SerialBeeFactory) Name() string {
	return "Serial Port"
}

// Description returns the description of this Bee.
func (factory *SerialBeeFactory) Description() string {
	return "Sends and receives data over a serial device"
}

// Image returns the filename of an image for this Bee.
func (factory *SerialBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *SerialBeeFactory) LogoColor() string {
	return "#73d44c"
}

// Options returns the options available to configure this Bee.
func (factory *SerialBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "device",
			Description: "Serial device to use",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "baudrate",
			Description: "The baudrate you want to use",
			Type:        "int",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *SerialBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over the serial port",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				{
					Name:        "port",
					Description: "The port the message was received on",
					Type:        "int",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *SerialBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message via the serial port",
			Options: []bees.PlaceholderDescriptor{
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
	f := SerialBeeFactory{}
	bees.RegisterFactory(&f)
}
