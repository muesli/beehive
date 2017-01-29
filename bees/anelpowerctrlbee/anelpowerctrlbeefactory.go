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

package anelpowerctrlbee

import (
	"github.com/muesli/beehive/bees"
)

// AnelPowerCtrlBeeFactory is a factory for AnelPowerCtrlBees.
type AnelPowerCtrlBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *AnelPowerCtrlBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := AnelPowerCtrlBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// Name returns the name of this Bee.
func (factory *AnelPowerCtrlBeeFactory) Name() string {
	return "anelpowerctrlbee"
}

// Description returns the description of this Bee.
func (factory *AnelPowerCtrlBeeFactory) Description() string {
	return "A bee that controls Anel's PowerCtrl"
}

// Image returns the filename of an image for this Bee.
func (factory *AnelPowerCtrlBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *AnelPowerCtrlBeeFactory) LogoColor() string {
	return "#73d44c"
}

// Options returns the options available to configure this Bee.
func (factory *AnelPowerCtrlBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "server",
			Description: "Hostname of Anel PowerCtrl device, eg: 192.168.0.2",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "user",
			Description: "Username to authenticate with Anel PowerCtrl",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password to use to connect to Anel PowerCtrl",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *AnelPowerCtrlBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "switch",
			Description: "Switches a socket on or off",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "socket",
					Description: "Which socket to switch",
					Type:        "int",
				},
				{
					Name:        "state",
					Description: "True to activate the socket, false to cut the power",
					Type:        "bool",
				},
			},
		},
	}
	return actions
}

func init() {
	f := AnelPowerCtrlBeeFactory{}
	bees.RegisterFactory(&f)
}
