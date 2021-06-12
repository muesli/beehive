/*
 *    Copyright (C) 2019      CalmBit
 *                  2014-2019 Christian Muehlhaeuser
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
 *      CalmBit <calmbit@posto.net>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package discordbee is a bee for sending and receiving messages with Discord
// servers.
package discordbee

import (
	"github.com/bwmarrin/discordgo"
	"github.com/muesli/beehive/bees"
)

// DiscordBee is a bee for sending and receiving messages with Discord
// servers.
type DiscordBee struct {
	bees.Bee

	apiToken      string
	triggerPhrase string
	eventChan     chan bees.Event
	discord       *discordgo.Session
}

// Run executes the Bee's event loop.
func (mod *DiscordBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan

	mod.LogDebugf("Starting discord bee with apiToken %s", mod.apiToken)
	d, err := discordgo.New("Bot " + mod.apiToken)
	if err != nil {
		mod.LogFatal("Unable to start discordbee! Error: ", err)
	}
	mod.discord = d
	mod.LogDebugf("Done!")

	mod.discord.AddHandler(mod.isReady)
	mod.discord.AddHandler(mod.onSend)

	err = mod.discord.Open()

	if err != nil {
		mod.LogFatal("Error opening websocket: ", err)
	}

	defer mod.discord.Close()

	select {
	case <-mod.SigChan:
		return
	}
}

func (mod *DiscordBee) isReady(s *discordgo.Session, event *discordgo.Ready) {
	// Default status - once the persistance layer exists:
	// TODO: Persist status, use it here.
	mod.LogDebugf("isReady called")
	s.UpdateGameStatus(0, "with bees")
}

func (mod *DiscordBee) onSend(s *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == s.State.User.ID {
		return
	}

	c, e := s.Channel(event.ChannelID)
	if e != nil {
		mod.LogErrorf("Unable to find channel with id %s", event.ChannelID)
		return
	}
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "message",
		Options: []bees.Placeholder{
			{
				Name:  "contents",
				Type:  "string",
				Value: event.Content,
			},
			{
				Name:  "username",
				Type:  "string",
				Value: event.Author.String(),
			},
			{
				Name:  "channel_id",
				Type:  "string",
				Value: event.ChannelID,
			},
			{
				Name:  "channel_name",
				Type:  "string",
				Value: c.Name,
			},
		},
	}
	mod.eventChan <- ev

}

// Action triggers the action passed to it.
func (mod *DiscordBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "send":
		var contents string
		var channelID string
		action.Options.Bind("contents", &contents)
		action.Options.Bind("channel_id", &channelID)

		_, err := mod.discord.ChannelMessageSend(channelID, contents)
		if err != nil {
			mod.LogErrorf("Unable to send message: %v", err)
		}
	case "send_news":
		var contents string
		var channelID string
		var message *discordgo.Message
		action.Options.Bind("contents", &contents)
		action.Options.Bind("channel_id", &channelID)

		message, err := mod.discord.ChannelMessageSend(channelID, contents)
		if err != nil {
			mod.LogErrorf("Unable to send message: %v", err)
		}
        message, err2 := mod.discord.ChannelMessageCrosspost(channelID, message.ID)
        if err2 != nil {
                mod.LogErrorf("Unable to publish message: %v", err2)
        }
	case "set_status":
		var status string
		action.Options.Bind("status", &status)

		err := mod.discord.UpdateGameStatus(0, status)
		if err != nil {
			mod.LogErrorf("Unable to update status: %v", err)
		}
	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}
	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *DiscordBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("api_token", &mod.apiToken)
	options.Bind("trigger_phrase", &mod.triggerPhrase)
}
