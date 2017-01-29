/*
 *    Copyright (C) 2014      Stefan 'glaxx' Luecke
 *                  2014-2017 Christian Muehlhaeuser
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
 *		Stefan Luecke <glaxx@glaxx.net>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package timebee is a Bee that can fire events at a specific time.
package timebee

import (
	"fmt"
	"time"

	"github.com/muesli/beehive/bees"
)

type TimeBee struct {
	bees.Bee
	cur_time, last_event                                     MyTime
	second, minute, hour, dayofweek, dayofmonth, month, year int
	eventChan                                                chan bees.Event
}

type MyTime struct {
	second, minute, hour, dayofweek, dayofmonth, month, year int
}

func (mod *TimeBee) Timer() {
	fail := false
	mod.cur_time.second = int(time.Now().Second())
	mod.cur_time.minute = int(time.Now().Minute())
	mod.cur_time.hour = int(time.Now().Hour())
	mod.cur_time.dayofweek = int(time.Now().Weekday())
	mod.cur_time.dayofmonth = int(time.Now().Day())
	mod.cur_time.month = int(time.Now().Month())
	mod.cur_time.year = int(time.Now().Year())
	if mod.second > 59 || mod.minute > 59 || mod.dayofweek > 6 || mod.dayofmonth > 31 || mod.month > 12 || mod.year > 9999 {
		fmt.Println("Error: Date is invalid")
		return
	}
	if mod.cur_time.second != mod.second && mod.second != -1 {
		fail = true
	}
	if mod.cur_time.minute != mod.minute && mod.minute != -1 {
		fail = true
	}
	if mod.cur_time.hour != mod.hour && mod.hour != -1 {
		fail = true
	}
	if mod.cur_time.dayofweek != mod.dayofweek && mod.dayofweek != -1 {
		fail = true
	}
	if mod.cur_time.dayofmonth != mod.dayofmonth && mod.dayofmonth != -1 {
		fail = true
	}
	if mod.cur_time.month != mod.month && mod.month != -1 {
		fail = true
	}
	if mod.cur_time.year != mod.year && mod.year != -1 {
		fail = true
	}

	if fail == true || mod.cur_time == mod.last_event {
		return
	}

	mod.last_event = mod.cur_time
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "time_event",
	}
	mod.eventChan <- event
}

func (mod *TimeBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

func (mod *TimeBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan
	for {
		select {
		case <-mod.SigChan:
			return

		default:
		}

		mod.Timer()
		time.Sleep(500 * time.Millisecond)
	}
}

func (mod *TimeBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("Second", &mod.second)
	options.Bind("Minute", &mod.minute)
	options.Bind("Hour", &mod.hour)
	options.Bind("DayOfWeek", &mod.dayofweek)
	options.Bind("DayOfMonth", &mod.dayofmonth)
	options.Bind("Month", &mod.month)
	options.Bind("Year", &mod.year)
}
