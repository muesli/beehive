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

package mastodonbee

import "github.com/muesli/beehive/bees"

// mastodonBeeFactory is a factory for mastodonBees.
type mastodonBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *mastodonBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := mastodonBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *mastodonBeeFactory) ID() string {
	return "mastodonbee"
}

// Name returns the name of this Bee.
func (factory *mastodonBeeFactory) Name() string {
	return "mastodon"
}

// Description returns the description of this Bee.
func (factory *mastodonBeeFactory) Description() string {
	return "Interact with mastodon"
}

// Image returns the filename of an image for this Bee.
func (factory *mastodonBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *mastodonBeeFactory) LogoColor() string {
	return "#003b66"
}

// Options returns the options available to configure this Bee.
func (factory *mastodonBeeFactory) Options() []bees.BeeOptionDescriptor {
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
func (factory *mastodonBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "toot_sent",
			Description: "A toot has been sent",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "Text of the toot that has been sent",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *mastodonBeeFactory) Actions() []bees.ActionDescriptor {
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
	}
	return actions
}

func init() {
	f := mastodonBeeFactory{}
	bees.RegisterFactory(&f)
}
