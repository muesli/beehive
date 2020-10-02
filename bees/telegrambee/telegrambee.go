/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
 *                  2017 Christian Muehlhaeuser
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
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package telegrambee is a Bee that can connect to Telegram.
package telegrambee

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/muesli/beehive/bees"
)

const (
	ModeHTML     = telegram.ModeHTML
	ModeMarkdown = "MarkdownV2"
	ModePlain    = ""
)

// TelegramBee is a Bee that can connect to Telegram.
type TelegramBee struct {
	bees.Bee

	// Telegram bot API Key
	apiKey string
	// Bot API client
	bot *telegram.BotAPI
}

// Action triggers the action passed to it.
func (mod *TelegramBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		chatID := ""
		text := ""
		parseMode := ""
		action.Options.Bind("chat_id", &chatID)
		action.Options.Bind("text", &text)
		action.Options.Bind("parse_mode", &parseMode)

		cid, err := strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			panic("Invalid telegram chat ID")
		}

		// use newer version of markdown parser
		if parseMode == telegram.ModeMarkdown {
			parseMode = ModeMarkdown
		}
		if parseMode != ModePlain && parseMode != ModeHTML && parseMode != ModeMarkdown {
			panic("Invalid parse mode")
		}

		msg := telegram.NewMessage(cid, text)
		msg.ParseMode = parseMode
		_, err = mod.bot.Send(msg)
		if err != nil {
			mod.Logf("Error sending message %v", err)
		}
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *TelegramBee) Run(eventChan chan bees.Event) {
	var err error
	mod.bot, err = telegram.NewBotAPI(mod.apiKey)
	if err != nil {
		mod.LogErrorf("Authorization failed, make sure the Telegram API key is correct: %s", err)
		return
	}
	mod.Logf("Authorized on account %s", mod.bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := mod.bot.GetUpdatesChan(u)
	if err != nil {
		mod.LogFatal("Failed retrieving the updates channel. err: ", err)
	}

	for {
		select {
		case <-mod.SigChan:
			return
		case update := <-updates:
			if update.Message == nil || update.Message.Text == "" {
				continue
			}

			mod.LogDebugf("Message received: %+v", update.Message)

			ev := bees.Event{
				Bee:  mod.Name(),
				Name: "message",
				Options: []bees.Placeholder{
					{
						Name:  "text",
						Type:  "string",
						Value: update.Message.Text,
					},
					{
						Name:  "chat_id",
						Type:  "string",
						Value: strconv.FormatInt(update.Message.Chat.ID, 10),
					},
					{
						Name:  "user_id",
						Type:  "string",
						Value: strconv.Itoa(update.Message.From.ID),
					},
					{
						Name:  "timestamp",
						Type:  "timestamp",
						Value: time.Now(),
					},
				},
			}
			eventChan <- ev
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TelegramBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	apiKey := getAPIKey(&options)
	mod.apiKey = apiKey
}

// Gets the Bot's API key from a file, the recipe config or the
// TELEGRAM_API_KEY environment variable.
func getAPIKey(options *bees.BeeOptions) string {
	var apiKey string
	options.Bind("api_key", &apiKey)

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
