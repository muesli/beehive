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

package mastodonbee

import "github.com/muesli/beehive/bees"

// MastodonBeeFactory is a factory for mastodonBees.
type MastodonBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *MastodonBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := MastodonBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *MastodonBeeFactory) ID() string {
	return "mastodonbee"
}

// Name returns the name of this Bee.
func (factory *MastodonBeeFactory) Name() string {
	return "mastodon"
}

// Description returns the description of this Bee.
func (factory *MastodonBeeFactory) Description() string {
	return "Interact with mastodon"
}

// Image returns the filename of an image for this Bee.
func (factory *MastodonBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *MastodonBeeFactory) LogoColor() string {
	return "#003b66"
}

// Options returns the options available to configure this Bee.
func (factory *MastodonBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "server",
			Description: "URL for the desired mastodon server",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "client_id",
			Description: "Client id for the mastodon client",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "client_secret",
			Description: "Client secret for the mastodon client",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "email",
			Description: "User account email",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "User account password",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *MastodonBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "deleted",
			Description: "is triggered when a toot has been deleted",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of the deleted toot",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "toot_fetched",
			Description: "is triggered when a toot has been fetched",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of the toot",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "Text of the toot that has been sent",
					Type:        "string",
				},
				{
					Name:        "user_id",
					Description: "Mastodon ID if the toots author",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "Mastodon handle of the toots author",
					Type:        "string",
				},
				{
					Name:        "reblogs",
					Description: "reblogs count",
					Type:        "int64",
				},
				{
					Name:        "favourites",
					Description: "favourites count",
					Type:        "int64",
				},
				{
					Name:        "url",
					Description: "The url for the toot",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "follow",
			Description: "is triggered when someone wants to follow you",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user_id",
					Description: "Mastodon ID of the user which triggered the follow event",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "Mastodon handle of the user which triggered the follow event",
					Type:        "string",
				},
				{
					Name:        "following",
					Description: "Indicates if your're following the user",
					Type:        "bool",
				},
				{
					Name:        "followed_by",
					Description: "Indicates if you're followed by the user",
					Type:        "bool",
				},
				{
					Name:        "followers",
					Description: "Number of followers for the user which triggered the follow request",
					Type:        "int64",
				},
				{
					Name:        "follows",
					Description: "Number of follows for the user which triggered the follow request",
					Type:        "int64",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "followed",
			Description: "is triggered when you followed someone on Mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user_id",
					Description: "Mastodon ID of the user to follow",
					Type:        "string",
				},
				{
					Name:        "following",
					Description: "Indicates if your're following the user",
					Type:        "bool",
				},
				{
					Name:        "requested",
					Description: "Indicates if you've requested following the user",
					Type:        "bool",
				},
				{
					Name:        "followed_by",
					Description: "Indicates if you're followed by the user",
					Type:        "bool",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "unfollowed",
			Description: "is triggered when you unfollow someone on Mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user_id",
					Description: "Mastodon ID of the user to follow",
					Type:        "string",
				},
				{
					Name:        "following",
					Description: "Indicates if your're following the user",
					Type:        "bool",
				},
				{
					Name:        "followed_by",
					Description: "Indicates if you're followed by the user",
					Type:        "bool",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "favourite",
			Description: "is triggered when someone favourites one of your toots",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of toot",
					Type:        "string",
				},
				{
					Name:        "user_id",
					Description: "Mastodon ID of the user that favourited your toot",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "The Mastodon handle of the user that favourited your toot",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "text content of the favourited toot",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL of the favourited toot",
					Type:        "string",
				},
				{
					Name:        "favourites",
					Description: "The count of favourites for this toot",
					Type:        "int64",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "favourited",
			Description: "is triggered when you favourite someones toot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of toot",
					Type:        "string",
				},
				{
					Name:        "user_id",
					Description: "Mastodon ID of the toots author",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "The Mastodon handle of the toots author",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "text content of the favourited toot",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL of the favourited toot",
					Type:        "string",
				},
				{
					Name:        "favourites",
					Description: "The count of favourites for this toot",
					Type:        "int64",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "reblog",
			Description: "is triggered when someone reblogs on of your toots",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user_id",
					Description: "Mastodon ID of the user tht reblogged your toot",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "Mastodon handle of the user that reblogged your toot",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "text content of the reblog",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL of the reblog",
					Type:        "string",
				},
				{
					Name:        "reblogs",
					Description: "Number of reblogs for the post",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "reblogged",
			Description: "is triggered when you reblog a toot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "user_id",
					Description: "Mastodon ID of the reblogged toots author",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "Mastodon handle of the reblogged toots author",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "text content of the reblog",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL of the reblog",
					Type:        "string",
				},
				{
					Name:        "reblogs",
					Description: "Number of reblogs for the post",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "mention",
			Description: "is triggered whenever someone mentions you on Mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of toot",
					Type:        "string",
				},
				{
					Name:        "user_id",
					Description: "Mastodon ID if the mention's author",
					Type:        "string",
				},
				{
					Name:        "username",
					Description: "The Mastodon handle of the mention's author",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "text content of the mention",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL of the mention",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *MastodonBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "toot",
			Description: "Post a new status toot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "Text of the status to toot, may not be longer than 500 characters",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "delete_toot",
			Description: "Delete a toot from mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of the to toot to delete",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "get_toots",
			Description: "Returns the current user's toots",
			Options:     []bees.PlaceholderDescriptor{},
		},
		{
			Namespace:   factory.Name(),
			Name:        "follow",
			Description: "Follow an user on Mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of the to user to follow",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "unfollow",
			Description: "unfollow an user on Mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of the to user to unfollow",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "reblog",
			Description: "reblog a toot on Mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of the to toot to reblog",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "favourite",
			Description: "favourite a toot on Mastodon",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The ID of the to toot to favourite",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := MastodonBeeFactory{}
	bees.RegisterFactory(&f)
}
