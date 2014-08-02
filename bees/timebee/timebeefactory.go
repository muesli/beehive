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
	"github.com/muesli/beehive/bees"
)

type TimeBeeFactory struct {
	bees.BeeFactory
}

func (factory *TimeBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TimeBee{
		Module: bees.NewBee(name, factory.Name(), description),
		second: int(options.GetValue("Second").(float64)),
		minute: int(options.GetValue("Minute").(float64)),
		hour: int(options.GetValue("Hour").(float64)),
		dayofweek: int(options.GetValue("DayOfWeek").(float64)),
		dayofmonth: int(options.GetValue("DayOfMonth").(float64)),
		month: int(options.GetValue("Month").(float64)),
		year: int(options.GetValue("Year").(float64)),
	}
	return &bee
}

func (factory *TimeBeeFactory) Name() string {
	return "timebee"
}

func (factory *TimeBeeFactory) Description() string {
	return "A bee that triggers an event at a given time"
}

func (factory *TimeBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		bees.BeeOptionDescriptor{
			Name:		"Second",
			Description:	"00-59 for a specific time; -1 for ignore",
			Type:		"int",
		},
		bees.BeeOptionDescriptor{
			Name:		"Minute",
			Description:	"00-59 for a specific time; -1 for ignore",
			Type:		"int",
		},
		bees.BeeOptionDescriptor{
			Name:		"Hour",
			Description:	"00-23 for a specific time; -1 for ignore",
			Type:		"int",
		},
		bees.BeeOptionDescriptor{
			Name:		"DayOfWeek",
			Description:	"0-6 0 = Sunday 6 = Saturday; -1 for ignore",
			Type:		"int",
		},
		bees.BeeOptionDescriptor{
			Name:		"DayOfMonth",
			Description:	"01-31 for a specific time; -1 for ignore)",
			Type:		"int",
		},
		bees.BeeOptionDescriptor{
			Name:		"Month",
			Description:	"01 - 12 for a specific time; -1 for ignore)",
			Type:		"int",
		},
		bees.BeeOptionDescriptor{
			Name:		"Year",
			Description:	"2014 - 9999 for specific time (non-reoccuring); -1 for ignore (recommended)",
			Type:		"int",
		},
	}
	return opts
}

func (factory *TimeBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		bees.EventDescriptor{
			Namespace:	factory.Name(),
			Name:		"time_event",
			Description:	"The time has come ...",
			Options: []bees.PlaceholderDescriptor{},
			},
		}
	return events
}
/*
func (factory *TimeBeeFactory) Actions() []bees.ActionDescriptor {
        actions := []bees.ActionDescriptor{}
        return actions
}

func (factory *TimeBeeFactory) Image() string {
	return ""
}*/

func init() {
	f := TimeBeeFactory{}
	bees.RegisterFactory(&f)
}
