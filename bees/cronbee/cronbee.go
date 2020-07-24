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

// Package cronbee is a Bee that acts like a time-based job scheduler (cron).
package cronbee

import (
	"strings"
	"time"

	"github.com/muesli/beehive/bees"
	"github.com/robfig/cron/v3"
)

// CronBee is a Bee that acts like a time-based job scheduler (cron).
type CronBee struct {
	bees.Bee
	input     [6]string
	eventChan chan bees.Event
}

func (mod *CronBee) runTask() {
	mod.LogDebugf("Running cron task " + strings.Join(mod.input[:], " "))
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "time",
		Options: []bees.Placeholder{
			{
				Name:  "timestamp",
				Type:  "string",
				Value: time.Now().Format(time.RFC3339),
			},
		},
	}
	mod.eventChan <- event
}

// Run executes the Bee's event loop.
func (mod *CronBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan
	c := cron.New(cron.WithSeconds())
	mod.LogDebugf("Scheduling " + strings.Join(mod.input[:], " "))
	_, err := c.AddFunc(strings.Join(mod.input[:], " "), mod.runTask)
	if err != nil {
		mod.LogFatal(err)
	}
	c.Start()

	for {
		select {
		case <-mod.SigChan:
			return
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *CronBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("second", &mod.input[0])
	options.Bind("minute", &mod.input[1])
	options.Bind("hour", &mod.input[2])
	options.Bind("day_of_week", &mod.input[3])
	options.Bind("day_of_month", &mod.input[4])
	options.Bind("month", &mod.input[5])
}
