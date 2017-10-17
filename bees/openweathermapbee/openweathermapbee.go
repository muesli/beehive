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

// Package openweathermapbee is a Bee that can interact with cleverbot
package openweathermapbee

import (
	"github.com/muesli/beehive/bees"

	owm "github.com/briandowns/openweathermap"
)

// OpenweathermapBee is a Bee that can chat with cleverbot
type OpenweathermapBee struct {
	bees.Bee

	current *owm.CurrentWeatherData

	unit     string
	language string
	key      string

	evchan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *OpenweathermapBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "get_current_weather":
		var location string
		action.Options.Bind("location", &location)

		err := mod.current.CurrentByName(location)
		if err != nil {
			mod.LogErrorf("Failed to fetch weather: %v", err)
			return nil
		}

		mod.TriggerCurrentWeatherEvent()

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *OpenweathermapBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	current, err := owm.NewCurrent(mod.unit, mod.language, mod.key)
	if err != nil {
		mod.LogErrorf("Failed to create new current service: %v", err)
		return
	}
	mod.current = current

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *OpenweathermapBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("unit", &mod.unit)
	options.Bind("language", &mod.language)
	options.Bind("key", &mod.key)
}
