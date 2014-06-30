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

type TimeBeeFactory struct {
	modules.ModuleFactory
}

func (factory *TimeBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := TimeBee{
		MyTime.second: options.GetValue("Second").(int),
		MyTime.minute: options.GetValue("Minute").(int),
		MyTime.hour: options.GetValue("Hour").(int),
		MyTime.dayofweek: options.GetValue("DayOfWeek").(int),
		MyTime.dayofmonth: options.GetValue("DayOfMonth").(int),
		MyTime.month: options.GetValue("Month").(int),
		MyTime.year: options.GetValue("Year").(int),
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
			Name:		"Second",
			Description:	"00-59 for a specific time; -1 for ignore",
			Type:		"int",
		},
		modules.BeeOptionDescriptor{
			Name:		"Minute",
			Description:	"00-59 for a specific time; -1 for ignore",
			Type:		"int",
		},
		modules.BeeOptionDescriptor{
			Name:		"Hour",
			Description:	"00-23 for a specific time; -1 for ignore",
			Type:		"int",
		},
		modules.BeeOptionDescriptor{
			Name:		"DayOfWeek",
			Description:	"0-6 0 = Sunday 6 = Saturday; -1 for ignore",
			Type:		"int",
		},
		modules.BeeOptionDescriptor{
			Name:		"DayOfMonth",
			Description:	"01-31 for a specific time; -1 for ignore)",
			Type:		"int",
		},
		modules.BeeOptionDescriptor{
			Name:		"Month",
			Description:	"01 - 12 for a specific time; -1 for ignore)",
			Type:		"int",
		},
		modules.BeeOptionDescriptor{
			Name:		"Year",
			Description:	"2014 - 9999 for specific time (non-reoccuring); -1 for ignore (recommended)",
			Type:		"int",
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
/*
func (factory *TimeBeeFactory) Actions() []modules.ActionDescriptor {
        actions := []modules.ActionDescriptor{}
        return actions
}

func (factory *TimeBeeFactory) Image() string {
	return ""
}*/

func init() {
	f := TimeBeeFactory{}
	modules.RegisterFactory(&f)
}
