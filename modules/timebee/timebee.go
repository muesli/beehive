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
)

type TimeBee struct {
	modules.Module

	time Time

	eventChan chan modules.Event
}

func (mod *TimeBee) Timer() {
	t := time.Now()
	if t == mod.time {
		event := modules.Event{
			Bee: mod.Name()
			Name: "time_event",
			Options: []modules.Placeholder{}
		}
		mod.eventChan <- event
	}
}

func (mod *TimeBee) Run(eventChan chan modules.Event) {
	for {
		mod.Timer(t)
		time.Sleep(1 * time.Second)
	}
)
