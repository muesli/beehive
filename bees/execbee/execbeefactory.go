/*
 *    Copyright (C) 2015 Dominik Schmidt
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
 *      Dominik Schmidt <domme@tomahawk-player.org>
 */

package execbee

import (
	"github.com/muesli/beehive/bees"
)

type ExecBeeFactory struct {
	bees.BeeFactory
}

// Interface impl
func (factory *ExecBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := ExecBee{
		Bee: bees.NewBee(name, factory.Name(), description),
	}

	return &bee
}

func (factory *ExecBeeFactory) Name() string {
	return "execbee"
}

func (factory *ExecBeeFactory) Description() string {
	return "A bee that allows executing commands"
}

// func (factory *ExecBeeFactory) Image() string {
// 	return factory.Name() + ".png"
// }

func (factory *ExecBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		bees.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "cmd",
			Description: "A command was executed",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "stdout",
					Description: "std output of the executed command",
					Type:        "string",
				},
				bees.PlaceholderDescriptor{
					Name:        "stderr",
					Description: "stderr output of the executed command",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (factory *ExecBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		bees.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "localCommand",
			Description: "Executes a command on the local host",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "command",
					Description: "command to be executed",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := ExecBeeFactory{}
	bees.RegisterFactory(&f)
}
