/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

// beehive's Philips Hue module.
package huebee

import (
	_ "log"
	"strconv"
	"strings"

	"github.com/muesli/beehive/bees"
	"github.com/muesli/go.hue"
)

type HueBee struct {
	bees.Bee

	client *hue.Bridge

	key    string
	bridge string
}

// Interface impl

func (mod *HueBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "setcolor":
		var lightId int
		var color string
		action.Options.Bind("light", &lightId)
		action.Options.Bind("color", &color)

		light, err := mod.client.FindLightById(strconv.Itoa(lightId))
		if err != nil {
			panic(err)
		}

		state := hue.SetLightState{
			On:  "true",
			Bri: "254",
		}

		switch strings.ToLower(color) {
		case "green":
			state.Hue = strconv.FormatInt(182*140, 10)
			state.Sat = strconv.FormatInt(254, 10)
		case "red":
			state.Hue = strconv.FormatInt(0, 10)
			state.Sat = strconv.FormatInt(254, 10)
		case "blue":
			state.Hue = strconv.FormatInt(182*250, 10)
			state.Sat = strconv.FormatInt(254, 10)
		}
		light.SetState(state)

	case "switch":
		var lightId int
		var state bool
		action.Options.Bind("light", &lightId)
		action.Options.Bind("state", &state)

		light, err := mod.client.FindLightById(strconv.Itoa(lightId))
		if err != nil {
			panic(err)
		}

		if state {
			light.On()
		} else {
			light.Off()
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *HueBee) Run(eventChan chan bees.Event) {
	mod.client = hue.NewBridge(mod.bridge, mod.key)
}
