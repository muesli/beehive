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

// TimeBee is a Bee that can fire events at a specific time.
type TimeBee struct {
	bees.Bee
	curTime, lastEvent                                       intTime
	second, minute, hour, dayofweek, dayofmonth, month, year int
	eventChan                                                chan bees.Event
}

type intTime struct {
	second, minute, hour, dayofweek, dayofmonth, month, year int
}

func (mod *TimeBee) timer() {
	fail := false
	mod.curTime.second = int(time.Now().Second())
	mod.curTime.minute = int(time.Now().Minute())
	mod.curTime.hour = int(time.Now().Hour())
	mod.curTime.dayofweek = int(time.Now().Weekday())
	mod.curTime.dayofmonth = int(time.Now().Day())
	mod.curTime.month = int(time.Now().Month())
	mod.curTime.year = int(time.Now().Year())
	if mod.second > 59 || mod.minute > 59 || mod.dayofweek > 6 || mod.dayofmonth > 31 || mod.month > 12 || mod.year > 9999 {
		fmt.Println("Error: Date is invalid")
		return
	}
	if mod.curTime.second != mod.second && mod.second != -1 {
		fail = true
	}
	if mod.curTime.minute != mod.minute && mod.minute != -1 {
		fail = true
	}
	if mod.curTime.hour != mod.hour && mod.hour != -1 {
		fail = true
	}
	if mod.curTime.dayofweek != mod.dayofweek && mod.dayofweek != -1 {
		fail = true
	}
	if mod.curTime.dayofmonth != mod.dayofmonth && mod.dayofmonth != -1 {
		fail = true
	}
	if mod.curTime.month != mod.month && mod.month != -1 {
		fail = true
	}
	if mod.curTime.year != mod.year && mod.year != -1 {
		fail = true
	}

	if fail == true || mod.curTime == mod.lastEvent {
		return
	}

	mod.lastEvent = mod.curTime
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "time_event",
	}
	mod.eventChan <- event
}

// Run executes the Bee's event loop.
func (mod *TimeBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan
	for {
		select {
		case <-mod.SigChan:
			return

		default:
		}

		mod.timer()
		time.Sleep(500 * time.Millisecond)
	}
}

// ReloadOptions parses the config options and initializes the Bee.
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
