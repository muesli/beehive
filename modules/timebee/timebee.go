/*
 *        Copyright (C) 2014 Stefan 'glaxx' Luecke 
 *
 *        This program is free software: you can redistribute it and/or modify
 *        it under the terms of the GNU Affero General Public License as published
 *        by the Free Software Foundation, either version 3 of the License, or
 *        (at your option) any later version.
 *
 *        This program is distributed in the hope that it will be useful,
 *        but WITHOUT ANY WARRANTY; without even the implied warranty of
 *        MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *        GNU Affero General Public License for more details.
 *
 *        You should have received a copy of the GNU Affero General Public License
 *        along with this program.      If not, see <http://www.gnu.org/licenses/>.
 *
 *        Authors:
 *		Stefan Luecke <glaxx@glaxx.net>
 */

// 

package timebee

import (
	"github.com/muesli/beehive/modules"
	"time"
	"fmt"
)

type TimeBee struct {
	modules.Module
	sched_time, cur_Time MyTime
	eventChan chan modules.Event
}

type MyTime struct {
	second, minute, hour, dayofweek, dayofmonth, month, year int
}

func (mod *TimeBee) Timer() {
	/*event := modules.Event{
		Bee: mod.Name(),
		Name: "time_event",
	}
	mod.eventChan <- event
	*/
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
	if 
}

func (mod *TimeBee) Action(action modules.Action) []modules.Placeholder {
        return []modules.Placeholder{}
}


func (mod *TimeBee) Run(eventChan chan modules.Event) {
	mod.eventChan = eventChan
	//mod.parsedtime, mod.parsererror = time.Parse(time.RFC1123, mod.time)
	//mod.parsedtime = time.Now()
	for {
		mod.Timer()
		time.Sleep(0.5 * time.Second)
	}
}
