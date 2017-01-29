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
 */

package tumblrbee

import (
	"github.com/muesli/beehive/bees"
)

// TumblrBeeFactory is a factory for TumblrBees.
type TumblrBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TumblrBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TumblrBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// Name returns the name of this Bee.
func (factory *TumblrBeeFactory) Name() string {
	return "tumblrbee"
}

// Description returns the description of this Bee.
func (factory *TumblrBeeFactory) Description() string {
	return "A Tumblr module for beehive"
}

// Image returns the filename of an image for this Bee.
func (factory *TumblrBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// Options returns the options available to configure this Bee.
func (factory *TumblrBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "blogname",
			Description: "Name of the Tumblr blog",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "consumer_key",
			Description: "Consumer Key",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "consumer_secret",
			Description: "Consumer Secret",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "token",
			Description: "Token",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "token_secret",
			Description: "Token Secret",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "callback_url",
			Description: "Callback URL",
			Type:        "url",
			Mandatory:   false,
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *TumblrBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "postText",
			Description: "Posts a text on Tumblr",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "Content of the Tumblr post",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "postQuote",
			Description: "Posts a quote on Tumblr",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "quote",
					Description: "Content of the Tumblr quote",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "source",
					Description: "Source of the Tumblr quote",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := TumblrBeeFactory{}
	bees.RegisterFactory(&f)
}
