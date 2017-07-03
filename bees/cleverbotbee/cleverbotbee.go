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

// Package cleverbotbee is a Bee that can interact with cleverbot
package cleverbotbee

import (
	"github.com/CleverbotIO/go-cleverbot.io"
	"github.com/muesli/beehive/bees"
)

// CleverbotBee is a Bee that can chat with cleverbot
type CleverbotBee struct {
	bees.Bee

	client *cleverbot.Session

	api_user     string
	api_key      string
	session_nick string

	evchan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *CleverbotBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send_message":
		var text string
		action.Options.Bind("text", &text)

		resp, err := mod.client.Ask(text)
		if err != nil {
			mod.LogErrorf("Failed to ask cleverbot: %v", err)
			return nil
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "answer",
			Options: []bees.Placeholder{
				{
					Name:  "answer",
					Type:  "string",
					Value: resp,
				},
			},
		}
		mod.evchan <- ev

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *CleverbotBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	client, err := cleverbot.New(mod.api_user, mod.api_key, mod.session_nick)
	if err != nil {
		mod.LogErrorf("Failed to start session: %v", err)
		return
	}
	mod.client = client

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *CleverbotBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("api_user", &mod.api_user)
	options.Bind("api_key", &mod.api_key)
	options.Bind("session_nick", &mod.session_nick)
}
