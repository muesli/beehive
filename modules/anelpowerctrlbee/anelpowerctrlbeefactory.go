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
	"github.com/muesli/beehive/modules"
)

type AnelPowerCtrlBeeFactory struct {
}

// Interface impl

func (factory *AnelPowerCtrlBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := AnelPowerCtrlBee{
		name:        name,
		namespace:   factory.Name(),
		description: description,
		addr:        options.GetValue("server").(string),
		user:        options.GetValue("user").(string),
		password:    options.GetValue("password").(string),
	}
	return &bee
}

func (factory *AnelPowerCtrlBeeFactory) Name() string {
	return "anelpowerctrlbee"
}

func (factory *AnelPowerCtrlBeeFactory) Description() string {
	return "A bee that controls Anel's PowerCtrl"
}

func (factory *AnelPowerCtrlBeeFactory) Options() []modules.BeeOptionDescriptor {
	opts := []modules.BeeOptionDescriptor{
		modules.BeeOptionDescriptor{
			Name:        "server",
			Description: "Hostname of Anel PowerCtrl device, eg: 192.168.0.2",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "user",
			Description: "Username to authenticate with Anel PowerCtrl",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "password",
			Description: "Password to use to connect to Anel PowerCtrl",
			Type:        "string",
		},
	}
	return opts
}

func (factory *AnelPowerCtrlBeeFactory) Events() []modules.EventDescriptor {
	return []modules.EventDescriptor{}
}

func (factory *AnelPowerCtrlBeeFactory) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "switch",
			Description: "Switches a socket on or off",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "socket",
					Description: "Which socket to switch",
					Type:        "int",
				},
				modules.PlaceholderDescriptor{
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
	modules.RegisterFactory(&f)
}
