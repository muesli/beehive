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
 *      Johannes Fürmann <johannes@weltraumpflege.org>
 */

package twitterbee

import (
	"github.com/muesli/beehive/bees"
)

// TwitterBeeFactory is a factory for TwitterBees.
type TwitterBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TwitterBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TwitterBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *TwitterBeeFactory) ID() string {
	return "twitterbee"
}

// Name returns the name of this Bee.
func (factory *TwitterBeeFactory) Name() string {
	return "Twitter"
}

// Description returns the description of this Bee.
func (factory *TwitterBeeFactory) Description() string {
	return "Tweets and reacts to events in your Twitter timeline"
}

// Image returns the filename of an image for this Bee.
func (factory *TwitterBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *TwitterBeeFactory) LogoColor() string {
	return "#00abec"
}

// Options returns the options available to configure this Bee.
func (factory *TwitterBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "consumer_key",
			Description: "Consumer key for Twitter API",
			Type:        "string",
		},
		{
			Name:        "consumer_secret",
			Description: "Consumer secret for Twitter API",
			Type:        "string",
		},
		{
			Name:        "access_token",
			Description: "Access token for Twitter API",
			Type:        "string",
		},
		{
			Name:        "access_token_secret",
			Description: "API secret for Twitter API",
			Type:        "string",
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *TwitterBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "tweet",
			Description: "Posts a new tweet",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "status",
					Description: "Text of the Status to tweet, may be no longer than 140 characters",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "follow",
			Description: "Follows another user",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "username",
					Description: "Twitter handle of the user you want to follow",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "unfollow",
			Description: "Unfollows another user",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "username",
					Description: "Twitter handle of the user you want to unfollow",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := TwitterBeeFactory{}
	bees.RegisterFactory(&f)
}
