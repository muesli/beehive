/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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
 *      Christian Muehlhaeuser <muesli@gmail.com>
 *      Nicolas Martin <penguwingit@gmail.com>
 */

// Package ryanairbee is a Bee that can post blogs & quotes on Tumblr.
package ryanairbee

import (
	"github.com/muesli/beehive/bees"
	"github.com/seblw/goryan"
	"time"
)

// RyanairBee is a Bee that can interact with the ryanairAPI
type RyanairBee struct {
	bees.Bee

	client *goryan.RyanairAPI

	evchan chan bees.Event
}

const dateLayout = "2006-01-02"

// Action triggers the action passed to it.
func (mod *RyanairBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "get_schedules":
		var cityFrom, cityTo, date string
		action.Options.Bind("city_from", &cityFrom)
		action.Options.Bind("city_to", &cityTo)
		action.Options.Bind("date", &date)

		// Parsing date && format to time.Time

		time, err := time.Parse(dateLayout, date)
		if err != nil {
			mod.LogErrorf("Failed to convert date/format: %v", err)
			return nil
		}

		// Fetching schedules
		flights, err := mod.client.GetSchedules(cityFrom, cityTo, time)
		if err != nil {
			mod.LogErrorf("Failed to fetch schedules: %v", err)
			return nil
		}

		for _, v := range flights {
			ev := bees.Event{
				Bee:  mod.Name(),
				Name: "flight_schedule",
				Options: []bees.Placeholder{
					{
						Name:  "number",
						Type:  "string",
						Value: v.Number,
					},
					{
						Name:  "departure_time",
						Type:  "string",
						Value: v.DepartureTime,
					},
					{
						Name:  "arrival_time",
						Type:  "string",
						Value: v.ArrivalTime,
					},
				},
			}
			mod.evchan <- ev
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *RyanairBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	mod.client = goryan.NewRyanairAPI(nil) // If nil {defaultHttpClient}

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *RyanairBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

}
