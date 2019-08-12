/*
 *    Copyright (C) 2016 Sergio Rubio
 *                  2017 Christian Muehlhaeuser
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
 *      Sergio Rubio <rubiojr@frameos.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package rocketchatbee is a Bee that can connect to Rocketchat.
package rocketchatbee

import (
	"github.com/muesli/beehive/bees"
)

// RocketchatBee is a Bee that can connect to Rocketchat.
type RocketchatBee struct {
	bees.Bee

	client *Client
}

// Action triggers the action passed to it.
func (mod *RocketchatBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		var (
			text    string
			channel string
			alias   string
		)

		action.Options.Bind("text", &text)
		action.Options.Bind("channel", &channel)
		action.Options.Bind("alias", &alias)

		err := mod.client.SendMessage(channel, text, alias)
		if err != nil {
			mod.LogErrorf("Failed to send message: %s", err)
			return outs
		}

	default:
		mod.LogErrorf("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *RocketchatBee) ReloadOptions(options bees.BeeOptions) {
	var (
		url       string
		userID    string
		authToken string
	)

	mod.SetOptions(options)

	options.Bind("url", &url)
	options.Bind("user_id", &userID)
	options.Bind("auth_token", &authToken)

	mod.client = NewClient(url, userID, authToken)

	err := mod.client.TestConnection()
	if err != nil {
		mod.LogErrorf("Connection to Rocket.Chat failed: %s", err)
	}

}
