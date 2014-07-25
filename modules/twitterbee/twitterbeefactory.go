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
	"github.com/muesli/beehive/modules"
	
)

type TwitterBeeFactory struct {
	modules.ModuleFactory
}

// Interface impl

func (factory *TwitterBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := TwitterBee{
		consumer_key:         options.GetValue("consumer_key").(string),
		consumer_secret:      options.GetValue("consumer_secret").(string),
		access_token:         options.GetValue("access_token").(string),
		access_token_secret:  options.GetValue("access_token_secret").(string),
	}

	bee.Module = modules.Module{name, factory.Name(), description}

	return &bee
}

func (factory *TwitterBeeFactory) Name() string {
	return "twitterbee"
}

func (factory *TwitterBeeFactory) Description() string {
	return "Tweet and receive Tweets."
}

func (factory *TwitterBeeFactory) Options() []modules.BeeOptionDescriptor {
	opts := []modules.BeeOptionDescriptor{
		modules.BeeOptionDescriptor{
			Name:        "consumer_key",
			Description: "Consumer key for Twitter API",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "consumer_secret",
			Description: "Consumer secret for Twitter API",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "access_token",
			Description: "Access token Twitter API",
			Type:        "string",
		},
		modules.BeeOptionDescriptor{
			Name:        "access_token_secret",
			Description: "API secret for Twitter API",
			Type:        "string",
		},
	}
	return opts
}

func (factory *TwitterBeeFactory) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "tweet",
			Description: "Update your status according to twitter",
			Options:     []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor {
					Name: "status",
					Description: "Text of the Status to tweet, may be no longer than 140 characters",
					Type: "String",
				},
			},
		},
	}
	return actions
}

func (factory *TwitterBeeFactory) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "call_finished",
			Description: "is triggered as soon as the API call has been executed",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name: "success",
					Description: "Result of the API call",
					Type: "bool",
				},
			},
		},
	}
	return events
}

func init() {
	f := TwitterBeeFactory{}
	modules.RegisterFactory(&f)
}
