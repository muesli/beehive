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
 */

// Package efabee is a Bee that interfaces with the public EVA API.
package efabee

import (
	"log"
	"time"

	"github.com/muesli/beehive/bees"
	"github.com/muesli/goefa"
)

type EFABee struct {
	bees.Bee

	Provider string
	efa      *goefa.EFAProvider

	eventChan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *EFABee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "departures":
		stop := ""
		action.Options.Bind("stop", &stop)

		//FIXME get departures
		_, station, err := mod.efa.FindStop(stop)
		if err != nil {
			log.Println("Stop does not exist or name is not unique!")
			return outs
		}
		log.Printf("Selected stop: %s (%d)\n\n",
			station[0].Name, station[0].Id)

		departures, err := station[0].Departures(time.Now(), 3)
		if err != nil {
			log.Println("Could not retrieve departure times!")
			return outs
		}
		for _, departure := range departures {
			log.Printf("Route %-5s due in %-2d minute%s --> %s\n",
				departure.ServingLine.Number,
				departure.Countdown,
				"s",
				departure.ServingLine.Direction)

			ev := bees.Event{
				Bee:  mod.Name(),
				Name: "departure",
				Options: []bees.Placeholder{
					{
						Name:  "mottype",
						Type:  "string",
						Value: departure.ServingLine.MotType.String(),
					},
					{
						Name:  "eta",
						Type:  "int",
						Value: departure.Countdown,
					},
					{
						Name:  "etatime",
						Type:  "string",
						Value: departure.DateTime.Format("15:04"),
					},
					{
						Name:  "route",
						Type:  "string",
						Value: departure.ServingLine.Number,
					},
					{
						Name:  "destination",
						Type:  "string",
						Value: departure.ServingLine.Direction,
					},
				},
			}
			mod.eventChan <- ev
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *EFABee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *EFABee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("provider", &mod.Provider)
	mod.efa, _ = goefa.ProviderFromJson(mod.Provider)
}
