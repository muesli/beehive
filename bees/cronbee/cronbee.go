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
	"github.com/muesli/beehive/bees"
	//"strings"
	"time"

	"github.com/muesli/beehive/bees/cronbee/cron"
)

type CronBee struct {
	bees.Bee
	input     [6]string
	eventChan chan bees.Event
}

func (mod *CronBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

func (mod *CronBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan
	timer := cron.ParseInput(mod.input)
	for {
		//FIXME: don't block
		select {
		case <-mod.SigChan:
			return

		default:
		}

		//fmt.Println(timer.NextEvent())
		time.Sleep(timer.DurationUntilNextEvent())
		event := bees.Event{
			Bee:  mod.Name(),
			Name: "time_event",
			Options: []bees.Placeholder{
				{
					Name:  "timestamp",
					Type:  "string",
					Value: timer.GetNextEvent(),
				},
			},
		}
		mod.eventChan <- event
	}
}

func (mod *CronBee) SetOptions(options bees.BeeOptions) {
	//FIXME: implement this
}
