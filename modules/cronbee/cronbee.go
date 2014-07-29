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

package cronbee

import (
	//"fmt"
	"github.com/muesli/beehive/modules"
	//"strings"
	"time"
	"github.com/muesli/beehive/modules/cronbee/cron"
)

type CronBee struct {
	modules.Module
	input [6]string
	eventChan chan modules.Event
}

func (mod *CronBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}

func (mod *CronBee) Run(eventChan chan modules.Event) {
	mod.eventChan = eventChan
	timer := cron.ParseInput(mod.input)
	for {
		//fmt.Println(timer.NextEvent())
		time.Sleep(timer.NextEvent())
		event := modules.Event{
			Bee: mod.Name(),
			Name: "time_event",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:  "timestamp", // Will be handy in future versions
					Type:  "string",
					Value: timer.CalculatedTime.String(),
				},
			},
		}
		mod.eventChan <- event
	}
}
