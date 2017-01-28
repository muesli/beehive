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

type TumblrBeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *TumblrBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TumblrBee{
		Bee:            bees.NewBee(name, factory.Name(), description, options),
		blogname:       options.GetValue("blogname").(string),
		callbackUrl:    options.GetValue("callback_url").(string),
		consumerKey:    options.GetValue("consumer_key").(string),
		consumerSecret: options.GetValue("consumer_secret").(string),
		token:          options.GetValue("token").(string),
		tokenSecret:    options.GetValue("token_secret").(string),
	}

	return &bee
}

func (factory *TumblrBeeFactory) Name() string {
	return "tumblrbee"
}

func (factory *TumblrBeeFactory) Description() string {
	return "A Tumblr module for beehive"
}

func (factory *TumblrBeeFactory) Image() string {
	return factory.Name() + ".png"
}

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

func (factory *TumblrBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

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
				},
				{
					Name:        "source",
					Description: "Optional source of the Tumblr quote",
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
