/*
 *    Copyright (C) 2019 Christian Muehlhaeuser
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

package socketbee

import (
	"github.com/muesli/beehive/bees"
)

// SocketBeeFactory is a factory for SocketBees.
type SocketBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *SocketBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SocketBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *SocketBeeFactory) ID() string {
	return "socketbee"
}

// Name returns the name of this Bee.
func (factory *SocketBeeFactory) Name() string {
	return "UDP Client"
}

// Description returns the description of this Bee.
func (factory *SocketBeeFactory) Description() string {
	return "Lets you transmit data via UDP sockets"
}

// Image returns the filename of an image for this Bee.
func (factory *SocketBeeFactory) Image() string {
	return "socketbee.png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *SocketBeeFactory) LogoColor() string {
	return "#223f5e"
}

// Events describes the available events provided by this Bee.
func (factory *SocketBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *SocketBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends data via a UDP socket",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "address",
					Description: "Address to connect to",
					Type:        "address",
					Mandatory:   true,
				},
				{
					Name:        "port",
					Description: "Port to connect to",
					Type:        "int",
					Mandatory:   true,
				},
				{
					Name:        "data",
					Description: "Data to send",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := SocketBeeFactory{}
	bees.RegisterFactory(&f)
}
