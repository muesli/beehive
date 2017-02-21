/*
 *    Copyright (C) 2017 Sebastian Ławniczak
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
 *      Sebastian Ławniczak <seb@seblw.me>
 */

// Package pastebinbee is a Bee that can interface with Pastebin.
package pastebinbee

import (
	pastebin "github.com/glaxx/go_pastebin"
	"github.com/muesli/beehive/bees"
)

// PastebinBee is a Bee that can interface with Pastebin.
type PastebinBee struct {
	bees.Bee

	client    pastebin.Pastebin
	apiDevKey string
}

// Action triggers the action passed to it.
func (mod *PastebinBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "post":
		title := ""
		content := ""
		expire := ""
		exposure := ""

		action.Options.Bind("title", &title)
		action.Options.Bind("content", &content)
		action.Options.Bind("expire", &expire)
		action.Options.Bind("exposure", &exposure)

		ret, err := mod.client.PasteAnonymous(content, title, "text", expire, exposure)
		if err != nil {
			mod.LogErrorf("Error occurred during posting to Pastebin: %v", err)
		} else {
			mod.Logln("Paste URL:", ret)
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}
	return outs
}

// Run executes the Bee's event loop.
func (mod *PastebinBee) Run(eventChan chan bees.Event) {
	mod.client = pastebin.NewPastebin(mod.apiDevKey)
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *PastebinBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("api_dev_key", &mod.apiDevKey)
}
