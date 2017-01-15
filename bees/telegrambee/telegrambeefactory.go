/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
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
 *      Gonzalo Izquierdo <lalotone@gmail.com>
 */

package telegrambee

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/muesli/beehive/bees"
	telegram "gopkg.in/telegram-bot-api.v4"
)

type TelegramBeeFactory struct {
	bees.BeeFactory
}

// Gets the Bot's API key from a file, the recipe config or the
// TELEGRAM_API_KEY environment variable.
func getApiKey(options *bees.BeeOptions) string {
	apiKey := options.GetValue("apiKey").(string)

	if strings.HasPrefix(apiKey, "file://") {
		buf, err := ioutil.ReadFile(strings.TrimPrefix(apiKey, "file://"))
		if err != nil {
			panic("Error reading API key file " + apiKey)
		}
		apiKey = string(buf)
	}

	if strings.HasPrefix(apiKey, "env://") {
		buf := strings.TrimPrefix(apiKey, "env://")
		apiKey = os.Getenv(string(buf))
	}

	return strings.TrimSpace(apiKey)
}

func (factory *TelegramBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	apiKey := getApiKey(&options)

	bot, err := telegram.NewBotAPI(apiKey)
	if err != nil {
		panic("Authorization failed, make sure the Telegram API key is correct")
	}
	log.Printf("TELEGRAM: Authorized on account %s", bot.Self.UserName)

	bee := TelegramBee{
		Bee:    bees.NewBee(name, factory.Name(), description),
		apiKey: apiKey,
		bot:    bot,
	}

	return &bee
}

func (factory *TelegramBeeFactory) Name() string {
	return "telegrambee"
}

func (factory *TelegramBeeFactory) Description() string {
	return "A Telegram bot bee"
}

func (factory *TelegramBeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *TelegramBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "apiKey",
			Description: "Telegram bot API key",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *TelegramBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message received via Telegram bot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				}, {
					Name:        "chatID",
					Description: "Telegram's chat ID",
					Type:        "string",
				},
				{
					Name:        "userID",
					Description: "User ID  sending the message",
					Type:        "string",
				},
			},
		},
	}

	return events
}

func (factory *TelegramBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{{
		Namespace:   factory.Name(),
		Name:        "send",
		Description: "Sends a message to a Telegram chat or group",
		Options: []bees.PlaceholderDescriptor{
			{
				Name:        "chatId",
				Description: "Telegram chat/group to send the message to",
				Type:        "string",
			},
			{
				Name:        "text",
				Description: "Content of the message",
				Type:        "string",
			},
		},
	}}
	return actions
}

func init() {
	f := TelegramBeeFactory{}
	bees.RegisterFactory(&f)
}
