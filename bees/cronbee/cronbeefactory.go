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

package cronbee

import (
	"github.com/muesli/beehive/bees"
)

type CronBeeFactory struct {
	bees.BeeFactory
}

func (factory *CronBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := CronBee{
		Module: bees.NewBee(name, factory.Name(), description),
	}
	bee.input[0] = options.GetValue("Second").(string)
	bee.input[1] = options.GetValue("Minute").(string)
	bee.input[2] = options.GetValue("Hour").(string)
	bee.input[3] = options.GetValue("DayOfWeek").(string)
	bee.input[4] = options.GetValue("DayOfMonth").(string)
	bee.input[5] = options.GetValue("Month").(string)
	return &bee
}

func (factory *CronBeeFactory) Name() string {
	return "cronbee"
}

func (factory *CronBeeFactory) Description() string {
	return "A bee that triggers an event at a given time"
}

func (factory *CronBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		bees.BeeOptionDescriptor{
			Name:		"Second",
			Description:	"00-59 for a specific time; * for ignore",
			Type:		"string",
		},
		bees.BeeOptionDescriptor{
			Name:		"Minute",
			Description:	"00-59 for a specific time; * for ignore",
			Type:		"string",
		},
		bees.BeeOptionDescriptor{
			Name:		"Hour",
			Description:	"00-23 for a specific time; * for ignore",
			Type:		"string",
		},
		bees.BeeOptionDescriptor{
			Name:		"DayOfWeek",
			Description:	"0-6 0 = Sunday 6 = Saturday; * for ignore",
			Type:		"string",
		},
		bees.BeeOptionDescriptor{
			Name:		"DayOfMonth",
			Description:	"01-31 for a specific time; * for ignore)",
			Type:		"string",
		},
		bees.BeeOptionDescriptor{
			Name:		"Month",
			Description:	"01 - 12 for a specific time; * for ignore)",
			Type:		"string",
		},
	}
	return opts
}

func (factory *CronBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		bees.EventDescriptor{
			Namespace:	factory.Name(),
			Name:		"time_event",
			Description:	"The time has come ...",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:			"timestamp", // For the lulz & future
					Description:	"Timestamp of the next event",
					Type:			"string",
				},
			},
		},
	}
	return events
}

func init() {
	f := CronBeeFactory{}
	bees.RegisterFactory(&f)
}
