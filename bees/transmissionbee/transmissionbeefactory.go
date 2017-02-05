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
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *TransmissionBeeFactory) ID() string {
	return "transmissionbee"
}

// Name returns the name of this Bee.
func (factory *TransmissionBeeFactory) Name() string {
	return "Transmission"
}

// Image returns the filename of an image for this Bee.
func (factory *TransmissionBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// Description returns the description of this Bee.
func (factory *TransmissionBeeFactory) Description() string {
	return "Lets you add torrents to a Transmission server"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *TransmissionBeeFactory) LogoColor() string {
	return "#111111"
}

// Options returns the options available to configure this Bee.
func (factory *TransmissionBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "url",
			Description: "Server URL",
			Type:        "url",
			Mandatory:   true,
			Default:     "http://localhost:9091/transmission/rpc",
		},
		{
			Name:        "username",
			Description: "Username",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password",
			Type:        "string",
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
				Description: "URL of the Torrent to download",
				Type:        "string",
				Mandatory:   true,
			},
		},
	}}
	return actions
}

func init() {
	f := TransmissionBeeFactory{}
	bees.RegisterFactory(&f)
}
