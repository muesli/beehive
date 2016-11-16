/*
 *    Copyright (C) 2016 Sergio Rubio
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
 *      Sergio Rubio <rubiojr@frameos.org>
 */

package slackbee

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/muesli/beehive/bees"
	"github.com/nlopes/slack"
)

type SlackBeeFactory struct {
	bees.BeeFactory
}

// Gets the API key from a file, the recipe config or the
// configured environment variable.
func getApiKey(options *bees.BeeOptions) string {
	apiKey := options.GetValue("apiKey").(string)

	if strings.HasPrefix(apiKey, "file://") {
		buf, err := ioutil.ReadFile(strings.TrimPrefix(apiKey, "file://"))
		if err != nil {
			panic("Slack: error reading API key file " + apiKey)
		}
		apiKey = string(buf)
	}

	if strings.HasPrefix(apiKey, "env://") {
		buf := strings.TrimPrefix(apiKey, "env://")
		apiKey = os.Getenv(string(buf))
	}

	return strings.TrimSpace(apiKey)
}

func (factory *SlackBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	apiKey := getApiKey(&options)
	client := slack.New(apiKey)
	_, err := client.AuthTest()
	if err != nil {
		panic("Slack: authentication failed!")
	}

	bee := SlackBee{
		Bee:      bees.NewBee(name, factory.Name(), description),
		apiKey:   apiKey,
		channels: map[string]string{},
		client:   client,
	}

	if options.GetValue("channels") != nil {
		for _, channel := range options.GetValue("channels").([]interface{}) {
			bee.channels[channel.(string)] = ""
		}
	}
	return &bee
}

func (factory *SlackBeeFactory) Name() string {
	return "slackbee"
}

func (factory *SlackBeeFactory) Description() string {
	return "A Slack module for beehive"
}

func (factory *SlackBeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *SlackBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		bees.BeeOptionDescriptor{
			Name:        "apiKey",
			Description: "Slack API key",
			Type:        "string",
			Mandatory:   true,
		},
		bees.BeeOptionDescriptor{
			Name:        "channels",
			Description: "Slack channels to listen on",
			Type:        "[]string",
			Mandatory:   false,
		},
	}
	return opts
}

func (factory *SlackBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		bees.EventDescriptor{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message was received over Slack",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				bees.PlaceholderDescriptor{
					Name:        "channel",
					Description: "The channel the message was received in",
					Type:        "string",
				},
				bees.PlaceholderDescriptor{
					Name:        "user",
					Description: "The user that sent the message",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (factory *SlackBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		bees.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a message to a channel",
			Options: []bees.PlaceholderDescriptor{
				bees.PlaceholderDescriptor{
					Name:        "channel",
					Description: "Which channel to send the message to",
					Type:        "string",
				},
				bees.PlaceholderDescriptor{
					Name:        "text",
					Description: "Content of the message",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := SlackBeeFactory{}
	bees.RegisterFactory(&f)
}
