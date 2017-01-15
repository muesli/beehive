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
	"log"
	"strconv"

	"github.com/muesli/beehive/bees"
	telegram "gopkg.in/telegram-bot-api.v4"
)

type TelegramBee struct {
	bees.Bee

	// Telegram bot API Key
	apiKey string
	// Bot API client
	bot *telegram.BotAPI
}

func (mod *TelegramBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		chatId := ""
		text := ""
		action.Options.Bind("chatId", &chatId)
		action.Options.Bind("text", &text)

		cid, err := strconv.Atoi(chatId)
		if err != nil {
			panic("Invalid telegram chat ID")
		}

		msg := telegram.NewMessage(int64(cid), text)
		_, err = mod.bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message %v", err)
		}
	}

	return outs
}

func (mod *TelegramBee) Run(eventChan chan bees.Event) {
	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := mod.bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

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
					Name:  "chatID",
					Type:  "string",
					Value: strconv.FormatInt(update.Message.Chat.ID, 10),
				},
				{
					Name:  "userID",
					Type:  "string",
					Value: strconv.Itoa(update.Message.From.ID),
				},
			},
		}
		eventChan <- ev
	}
}

func (mod *TelegramBee) Stop() {
	log.Println("Stopping the Telegram bee")
}
