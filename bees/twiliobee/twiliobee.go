/*
 *    Copyright (C) 2017 James Vaughan
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
 *      James Vaughan <james@jamesbvaughan.com>
 */

// Package twiliobee is a Bee that is able to send SMS messages.
package twiliobee

import (
	twilio "github.com/carlosdp/twiliogo"
	"github.com/muesli/beehive/bees"
)

// TwilioBee is a Bee that is able to send SMS messages.
type TwilioBee struct {
	bees.Bee

	client     *twilio.TwilioClient
	accountsid string
	authtoken  string
	fromNumber string
	toNumber   string
}

// Action triggers the action passed to it.
func (mod *TwilioBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		body := ""
		action.Options.Bind("body", &body)

		_, err := twilio.NewMessage(mod.client, mod.fromNumber, mod.toNumber, twilio.Body(body))
		if err != nil {
			mod.LogErrorf("Error sending twilio SMS: %s", err)
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *TwilioBee) Run(eventChan chan bees.Event) {
	mod.client = twilio.NewClient(mod.accountsid, mod.authtoken)
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TwilioBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("account_sid", &mod.accountsid)
	options.Bind("auth_token", &mod.authtoken)
	options.Bind("from_number", &mod.fromNumber)
	options.Bind("to_number", &mod.toNumber)
}
