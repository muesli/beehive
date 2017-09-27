/*
 *    Copyright (C) 2017 Dominik Schmidt
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
 *      Dominik Schmidt <dev@dominik-schmidt.de>
 */

package horizonboxbee

import (
	"github.com/muesli/beehive/bees"
)

// HorizonBoxBeeFactory is a factory for HorizonBoxBees.
type HorizonBoxBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *HorizonBoxBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := HorizonBoxBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *HorizonBoxBeeFactory) ID() string {
	return "horizonboxbee"
}

// Name returns the name of this Bee.
func (factory *HorizonBoxBeeFactory) Name() string {
	return "HorizonBox"
}

// Description returns the description of this Bee.
func (factory *HorizonBoxBeeFactory) Description() string {
	return "A module observing the state of a UnityMedia HorizonBox for beehive"
}

// // Image returns the filename of an image for this Bee.
// func (factory *HorizonBoxBeeFactory) Image() string {
// 	return factory.ID() + ".png"
// }

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *HorizonBoxBeeFactory) LogoColor() string {
	return "#212727"
}

func (factory *HorizonBoxBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "address",
			Description: "Address of the HorizonBox, eg: 192.168.192.1",
			Type:        "address",
			Mandatory:   true,
		},
		{
			Name:        "user",
			Description: "Username to login to the HorizonBox",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password to login to the HorizonBox",
			Type:        "password",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *HorizonBoxBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "connection_status_change",
			Description: "The internet connection changed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "online",
					Description: "The new connection status",
					Type:        "bool",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "external_ip_change",
			Description: "The external ip changed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "new_external_ip",
					Description: "The new external ip",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := HorizonBoxBeeFactory{}
	bees.RegisterFactory(&f)
}
