/*
 *    Copyright (C) 2020 Christian Muehlhaeuser
 *                  2020 Nicolas Martin
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
 *      Nicolas Martin <penguwin@penguwin.eu>
 */
package mastodonbee

import (
	"context"

	mastodon "github.com/mattn/go-mastodon"
	"github.com/muesli/beehive/bees"
)

// handleUpdateEvent handles incoming Toot updates from mastodon from yourself
// and from people you follow.
func (mod *MastodonBee) handleStatus(status *mastodon.Status) {
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "toot_fetched",
		Options: []bees.Placeholder{
			{
				Name:  "id",
				Value: string(status.ID),
				Type:  "string",
			},
			{
				Name:  "text",
				Value: status.Content,
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
				Name:  "reblogs",
				Value: status.ReblogsCount,
				Type:  "int64",
			},
			{
				Name:  "favourites",
				Value: status.FavouritesCount,
				Type:  "int64",
			},
			{
				Name:  "url",
				Value: status.URL,
				Type:  "string",
			},
			{
				Name:  "created",
				Value: status.CreatedAt,
				Type:  "time.Time",
			},
		},
	}
	mod.evchan <- ev
}

func (mod *MastodonBee) handleNotification(notif *mastodon.Notification) {
	switch notif.Type {

	case "follow":
		rel, err := mod.client.GetAccountRelationships(context.Background(), []string{string(notif.Account.ID)})
		if err != nil {
			mod.LogErrorf("Failed to fetch account relationship at follow event: %v", err)
			return
		}

		if len(rel) == 0 {
			// NOTE: Does this even happen?
			// TODO: Investigate
			mod.LogErrorf("No relationship information was fetched")
			return
		}

		// NOTE: there should be only one relationship entity fetched
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "follow",
			Options: []bees.Placeholder{
				{
					Name:  "user_id",
					Value: string(notif.Account.ID),
					Type:  "string",
				},
				{
					Name:  "username",
					Value: notif.Account.DisplayName,
					Type:  "string",
				},
				{
					Name:  "following",
					Value: rel[0].Following,
					Type:  "bool",
				},
				{
					Name:  "followed_by",
					Value: rel[0].FollowedBy,
					Type:  "bool",
				},
				{
					Name:  "followers",
					Value: notif.Account.FollowersCount,
					Type:  "int64",
				},
				{
					Name:  "follows",
					Value: notif.Account.FollowingCount,
					Type:  "int64",
				},
			},
		}
		mod.evchan <- ev

	case "favourite":
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "favourite",
			Options: []bees.Placeholder{
				{
					Name:  "id",
					Value: string(notif.Status.ID),
					Type:  "string",
				},
				{
					Name:  "user_id",
					Value: string(notif.Account.ID),
					Type:  "string",
				},
				{
					Name:  "username",
					Value: notif.Account.DisplayName,
					Type:  "string",
				},
				{
					Name:  "text",
					Value: notif.Status.Content,
					Type:  "string",
				},
				{
					Name:  "url",
					Value: notif.Status.URL,
					Type:  "string",
				},
				{
					Name:  "favourites",
					Value: notif.Status.FavouritesCount,
					Type:  "int64",
				},
			},
		}
		mod.evchan <- ev

	case "reblog":
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "reblog",
			Options: []bees.Placeholder{
				{
					Name:  "user_id",
					Value: string(notif.Account.ID),
					Type:  "string",
				},
				{
					Name:  "username",
					Value: notif.Status.Account.DisplayName,
					Type:  "string",
				},
				{
					Name:  "text",
					Value: notif.Status.Content,
					Type:  "string",
				},
				{
					Name:  "url",
					Value: notif.Status.URL,
					Type:  "string",
				},
				{
					Name:  "reblogs",
					Value: notif.Status.ReblogsCount,
					Type:  "int64",
				},
			},
		}
		mod.evchan <- ev

	case "mention":
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "mention",
			Options: []bees.Placeholder{
				{
					Name:  "id",
					Value: string(notif.Status.ID),
					Type:  "string",
				},
				{
					Name:  "user_id",
					Value: string(notif.Account.ID),
					Type:  "string",
				},
				{
					Name:  "username",
					Value: notif.Status.Account.DisplayName,
					Type:  "string",
				},
				{
					Name:  "text",
					Value: notif.Status.Content,
					Type:  "string",
				},
				{
					Name:  "url",
					Value: notif.Status.URL,
					Type:  "string",
				},
			},
		}
		mod.evchan <- ev

	default:
		mod.LogErrorf("Unkown notification type: %s", notif.Type)
	}
}
