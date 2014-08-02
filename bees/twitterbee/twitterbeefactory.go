/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

type TwitterBeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *TwitterBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TwitterBee{
		Module: bees.NewBee(name, factory.Name(), description),
		consumer_key:        options.GetValue("consumer_key").(string),
		consumer_secret:     options.GetValue("consumer_secret").(string),
		access_token:        options.GetValue("access_token").(string),
		access_token_secret: options.GetValue("access_token_secret").(string),
	}

	return &bee
}

func (factory *TwitterBeeFactory) Name() string {
	return "twitterbee"
}

func (factory *TwitterBeeFactory) Description() string {
	return "Tweet and receive Tweets."
}

func (factory *TwitterBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		bees.BeeOptionDescriptor{
			Name:        "consumer_key",
			Description: "Consumer key for Twitter API",
			Type:        "string",
		},
		bees.BeeOptionDescriptor{
			Name:        "consumer_secret",
			Description: "Consumer secret for Twitter API",
			Type:        "string",
		},
		bees.BeeOptionDescriptor{
			Name:        "access_token",
			Description: "Access token Twitter API",
			Type:        "string",
		},
		bees.BeeOptionDescriptor{
			Name:        "access_token_secret",
			Description: "API secret for Twitter API",
			Type:        "string",
		},
	}
	return opts
}

func (factory *TwitterBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		bees.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "tweet",
			Description: "Update your status according to twitter",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "status",
					Description: "Text of the Status to tweet, may be no longer than 140 characters",
					Type:        "String",
				},
			},
		},
	}
	return actions
}

func (factory *TwitterBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		bees.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "call_finished",
			Description: "is triggered as soon as the API call has been executed",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "success",
					Description: "Result of the API call",
					Type:        "bool",
				},
			},
		},
		bees.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "mention",
			Description: "is triggered as soon as the API call has been executed",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "username",
					Description: "Twitter handle of the mention's author",
					Type:        "string",
				},
				bees.PlaceholderDescriptor{
					Name:        "text",
					Description: "text content of the mention",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := TwitterBeeFactory{}
	bees.RegisterFactory(&f)
}
