/*
 *    Copyright (C) 2014 Stefan 'glaxx' Luecke     
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
 *      Stefan 'glaxx' Luecke <glaxx@glaxx.net>            
 */

package timebee

import (
	"github.com/muesli/beehive/modules"
)

type TimeBeeFactory struct {}

func (factory *TimeBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := TimeBee{
		time: options.GetValue("beeTime").(string),
	}
	bee.Module = modules.Module{name, factory.Name(), description}
	return &bee
}

func (factory *TimeBeeFactory) Name() string {
	return "timebee"
}

func (factory *TimeBeeFactory) Description() string {
	return "A bee that triggers an event at a given time"
}

func (factory *TimeBeeFactory) Options() []modules.BeeOptionDescriptor {
	opts := []modules.BeeOptionDescriptor{
		modules.BeeOptionDescriptor{
			Name:		"beeTime",
			Description:	"The time when the event triggers",
			Type:		"string",
		},
	}
	return opts
}

func (factory *TimeBeeFactory) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:	factory.Name(),
			Name:		"time_event",
			Description:	"The time has come ...",
			Options: []modules.PlaceholderDescriptor{},
			},
		}
	return events
}

func (factory *TimeBeeFactory) Actions() []modules.ActionDescriptor {
        actions := []modules.ActionDescriptor{}
        return actions
}

func init() {
	f := TimeBeeFactory{}
	modules.RegisterFactory(&f)
}
