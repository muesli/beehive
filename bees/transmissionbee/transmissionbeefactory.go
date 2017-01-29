/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
 *                  2017 Christian Muehlhaeuser
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
 *      Gonzalo Izquierdo <lalotone@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package transmissionbee

import "github.com/muesli/beehive/bees"

// TransmissionBeeFactory is a factory for TransmissionBees.
type TransmissionBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TransmissionBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TransmissionBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// Options returns the options available to configure this Bee.
func (factory *TransmissionBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "server_url",
			Description: "Transmission server URL",
			Type:        "url",
			Mandatory:   true,
		},
		{
			Name:        "username",
			Description: "Transmission server username",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Transmission server password",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *TransmissionBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{{
		Namespace:   factory.Name(),
		Name:        "add_torrent",
		Description: "Torrent URL or magnet",
		Options: []bees.PlaceholderDescriptor{
			{
				Name:        "torrent",
				Description: "Telegram chat/group to send the message to",
				Type:        "string",
				Mandatory:   true,
			},
			{
				Name:        "command_prefix",
				Description: "String that precedes the torrent URL/magnet (will be removed)",
				Type:        "string",
			},
		},
	}}
	return actions
}

// Name returns the name of this Bee.
func (factory *TransmissionBeeFactory) Name() string {
	return "transmissionbee"
}

// Image returns the filename of an image for this Bee.
func (factory *TransmissionBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// Description returns the description of this Bee.
func (factory *TransmissionBeeFactory) Description() string {
	return "A bee for adding torrents to a transmission server"
}

func init() {
	f := TransmissionBeeFactory{}
	bees.RegisterFactory(&f)
}
