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
		brightness := 254
		action.Options.Bind("light", &lightId)
		action.Options.Bind("color", &color)
		action.Options.Bind("brightness", &brightness)

		light, err := mod.client.FindLightById(strconv.Itoa(lightId))
		if err != nil {
			panic(err)
		}

		state := hue.SetLightState{
			On:  "true",
			Bri: strconv.FormatInt(int64(brightness), 10),
			Sat: "254",
		}

		switch strings.ToLower(color) {
		case "coolwhite":
			state.Hue = strconv.FormatInt(150, 10)
		case "warmwhite":
			state.Hue = strconv.FormatInt(500, 10)
		case "green":
			state.Hue = strconv.FormatInt(182*140, 10)
		case "red":
			state.Hue = strconv.FormatInt(0, 10)
		case "blue":
			state.Hue = strconv.FormatInt(182*250, 10)
		case "orange":
			state.Hue = strconv.FormatInt(182*25, 10)
		case "yellow":
			state.Hue = strconv.FormatInt(182*85, 10)
		case "pink":
			state.Hue = strconv.FormatInt(182*300, 10)
		case "purple":
			state.Hue = strconv.FormatInt(182*270, 10)
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
