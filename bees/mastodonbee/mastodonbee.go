/*
 *    Copyright (C) 2018 Nicolas Martin
 *                  2018 Christian Muehlhaeuser
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
 *      Nicolas Martin <penguwin@systemli.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package mastodonbee is a Bee that can connect to mastodon.
package mastodonbee

import (
	"context"

	mastodon "github.com/mattn/go-mastodon"

	"github.com/muesli/beehive/bees"
)

// mastodonBee is a Bee that can connect to mastodon.
type mastodonBee struct {
	bees.Bee

	server       string
	clientID     string
	clientSecret string
	email        string
	password     string

	client *mastodon.Client

	evchan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *mastodonBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "toot":
		var text string
		action.Options.Bind("text", &text)

		// Post status toot on mastodon
		status, err := mod.client.PostStatus(context.Background(), &mastodon.Toot{
			Status: text,
		})
		if err != nil {
			mod.LogErrorf("Error sending toot: %v", err)
		}

		// Handle back 'toot_sent' event
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "toot_sent",
			Options: []bees.Placeholder{
				{
					Name:  "text",
					Value: status.Content,
					Type:  "string",
				},
			},
		}
		mod.evchan <- ev
	}
	return outs
}

// Run executes the Bee's event loop.
func (mod *mastodonBee) Run(eventChan chan bees.Event) {
	// Create the new api client
	c := mastodon.NewClient(&mastodon.Config{
		Server:       mod.server,
		ClientID:     mod.clientID,
		ClientSecret: mod.clientSecret,
	})
	// authorize it
	err := c.Authenticate(context.Background(), mod.email, mod.password)
	if err != nil {
		mod.LogErrorf("Authorization failed, make sure the mastodon credentials are correct: %s", err)
		return
	}
	// try to get user account
	acc, err := c.GetAccountCurrentUser(context.Background())
	if err != nil {
		mod.LogErrorf("Failed to get current user account: %v", err)
	}
	mod.Logf("Successfully logged in: %s", acc.URL)

	// set client
	mod.client = c

	// set eventchan
	mod.evchan = eventChan

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *mastodonBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("server", &mod.server)
	options.Bind("client_id", &mod.clientID)
	options.Bind("client_secret", &mod.clientSecret)
	options.Bind("email", &mod.email)
	options.Bind("password", &mod.password)
}
