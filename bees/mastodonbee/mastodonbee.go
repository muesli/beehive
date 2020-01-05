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
 *      Nicolas Martin <penguwin@penguwin.eu>
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
		mod.Logf("Attempting to post \"%s\" to Mastodon", text)

		// Post status toot on mastodon, event gets triggered automatically via
		// the toot_fetched even.
		_, err := mod.client.PostStatus(context.Background(), &mastodon.Toot{
			Status: text,
		})
		if err != nil {
			mod.LogErrorf("Error sending toot: %v", err)
		}

	case "delete_toot":
		var id string
		action.Options.Bind("id", &id)
		mod.Logf("Attempting to delete toot \"%s\"", id)

		// Event is automatically handled in handleStreamEvent()
		err := mod.client.DeleteStatus(context.Background(), mastodon.ID(id))
		if err != nil {
			mod.LogErrorf("Error deleting toot: %v", err)
		}

	case "get_toots": // returns the current user's toots
		mod.Logf("Attempting to get current user")
		acc, err := mod.client.GetAccountCurrentUser(context.Background())
		if err != nil {
			mod.LogErrorf("Error getting current user: %v", err)
			return outs
		}

		mod.Logf("Attempting to get current user's toots")
		statuses, err := mod.client.GetAccountStatuses(context.Background(), acc.ID, &mastodon.Pagination{})
		if err != nil {
			mod.LogErrorf("Error getting current user's toots: %v", err)
			return outs
		}

		// create a toot_fetched event for every status
		for _, status := range statuses {
			mod.handleStatus(status)
		}

	case "follow":
		var id string
		action.Options.Bind("id", &id)
		mod.Logf("Attempting to follow account: %s", id)

		rel, err := mod.client.AccountFollow(context.Background(), mastodon.ID(id))
		if err != nil {
			mod.LogErrorf("Failed to follow account %s: %v", id, err)
			return outs
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "followed",
			Options: []bees.Placeholder{
				{
					Name:  "user_id",
					Value: id,
					Type:  "string",
				},
				{
					Name:  "following",
					Value: rel.Following,
					Type:  "bool",
				},
				{
					Name:  "requested",
					Value: rel.Requested,
					Type:  "bool",
				},
				{
					Name:  "followed_by",
					Value: rel.FollowedBy,
					Type:  "bool",
				},
			},
		}
		mod.evchan <- ev

	case "unfollow":
		var id string
		action.Options.Bind("id", &id)
		mod.Logf("Attempting to unfollow account: %s", id)

		rel, err := mod.client.AccountUnfollow(context.Background(), mastodon.ID(id))
		if err != nil {
			mod.LogErrorf("Failed to unfollow account %s: %v", id, err)
			return outs
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "unfollowed",
			Options: []bees.Placeholder{
				{
					Name:  "user_id",
					Value: id,
					Type:  "string",
				},
				{
					Name:  "following",
					Value: rel.Following,
					Type:  "bool",
				},
				{
					Name:  "followed_by",
					Value: rel.FollowedBy,
					Type:  "bool",
				},
			},
		}
		mod.evchan <- ev

	case "favourite":
		var id string
		action.Options.Bind("id", &id)
		mod.Logf("Attempting to favourite toot: %s", id)

		status, err := mod.client.Favourite(context.Background(), mastodon.ID(id))
		if err != nil {
			mod.LogErrorf("Failed to favourite toot: %v", err)
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "favourited",
			Options: []bees.Placeholder{
				{
					Name:  "id",
					Value: id,
					Type:  "string",
				},
				{
					Name:  "user_id",
					Value: string(status.Account.ID),
					Type:  "string",
				},
				{
					Name:  "username",
					Value: status.Account.DisplayName,
					Type:  "string",
				},
				{
					Name:  "text",
					Value: status.Content,
					Type:  "string",
				},
				{
					Name:  "url",
					Value: status.URL,
					Type:  "string",
				},
				{
					Name:  "favourites",
					Value: status.FavouritesCount,
					Type:  "int64",
				},
			},
		}
		mod.evchan <- ev

	case "reblog":
		var id string
		action.Options.Bind("id", &id)
		mod.Logf("Attempting to reblog toot %s", id)

		// reblog-Event should automatically be handled in handleStreamEvent()
		status, err := mod.client.Reblog(context.Background(), mastodon.ID(id))
		if err != nil {
			mod.LogErrorf("Failed to reblog toot: %v", err)
			return outs
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "reblogged",
			Options: []bees.Placeholder{
				{
					Name:  "user_id",
					Value: string(status.Account.ID),
					Type:  "string",
				},
				{
					Name:  "username",
					Value: status.Account.DisplayName,
					Type:  "string",
				},
				{
					Name:  "text",
					Value: status.Content,
					Type:  "string",
				},
				{
					Name:  "url",
					Value: status.URL,
					Type:  "string",
				},
				{
					Name:  "reblogs",
					Value: status.ReblogsCount,
					Type:  "int64",
				},
			},
		}
		mod.evchan <- ev

	default:
		mod.LogErrorf("Unkown action: %s", action.Name)
	}
	return outs
}

func (mod *mastodonBee) handleStreamEvent(item interface{}) {
	switch e := item.(type) {
	case *mastodon.UpdateEvent:
		mod.handleStatus(e.Status)
	case *mastodon.NotificationEvent:
		mod.handleNotification(e.Notification)
	case *mastodon.DeleteEvent:
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "deleted",
			Options: []bees.Placeholder{
				{
					Name:  "id",
					Value: string(e.ID),
					Type:  "string",
				},
			},
		}
		mod.evchan <- ev
	default:
		mod.LogErrorf("Unkown event: %+v", e)
	}
}

func (mod *mastodonBee) handleStream() {
	timeline, err := mod.client.StreamingUser(context.Background())
	if err != nil {
		mod.LogErrorf("Failed to get user stream: %+v", err)
		return
	}

	for {
		select {
		case <-mod.SigChan:
			return
		case item := <-timeline:
			mod.handleStreamEvent(item)
		}
	}
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

	// set and start eventchan
	mod.evchan = eventChan
	mod.handleStream()
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
