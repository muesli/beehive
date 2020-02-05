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

// Package huebee is a Bee that can talk to Philips Hue bridges.
package huebee

import (
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut/palette"
	hue "github.com/muesli/go.hue"

	"github.com/muesli/beehive/bees"
)

// HueBee is a Bee that can talk to Philips Hue bridges.
type HueBee struct {
	bees.Bee

	client *hue.Bridge

	key    string
	bridge string
}

// Action triggers the action passed to it.
func (mod *HueBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "setcolor":
		var lightID int
		var color string
		brightness := 254
		action.Options.Bind("light", &lightID)
		action.Options.Bind("color", &color)
		action.Options.Bind("brightness", &brightness)

		light, err := mod.client.FindLightById(strconv.Itoa(lightID))
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
			state.Ct = strconv.FormatInt(150, 10)
			state.Sat = ""
		case "warmwhite":
			state.Ct = strconv.FormatInt(500, 10)
			state.Sat = ""
		case "green":
			state.Hue = strconv.FormatInt(182*140, 10)
			state.Sat = "254"
		case "red":
			state.Hue = strconv.FormatInt(0, 10)
			state.Sat = "254"
		case "blue":
			state.Hue = strconv.FormatInt(182*250, 10)
			state.Sat = "254"
		case "orange":
			state.Hue = strconv.FormatInt(182*25, 10)
			state.Sat = "254"
		case "yellow":
			state.Hue = strconv.FormatInt(182*85, 10)
			state.Sat = "254"
		case "pink":
			state.Hue = strconv.FormatInt(182*300, 10)
			state.Sat = "254"
		case "purple":
			state.Hue = strconv.FormatInt(182*270, 10)
			state.Sat = "254"
		case "":
			// ignore
		default:
			var hex colorful.Color
			if strings.HasPrefix(color, "#") {
				// RGB color
				hex, err = colorful.Hex(color)
				if err != nil {
					return outs
				}
			} else {
				cols := palette.Wikipedia.Filter(color)

				var ok bool
				hex, ok = colorful.MakeColor(cols[0].Color)
				if !ok {
					return outs
				}
			}
			cx, cy, cz := hex.Xyz()
			cx = cx / (cx + cy + cz)
			cy = cy / (cx + cy + cz)
			state.Xy = append(state.Xy, float32(cx))
			state.Xy = append(state.Xy, float32(cy))
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

// ReloadOptions parses the config options and initializes the Bee.
func (mod *HueBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.bridge)
	options.Bind("key", &mod.key)

	mod.client = hue.NewBridge(mod.bridge, mod.key)
}
