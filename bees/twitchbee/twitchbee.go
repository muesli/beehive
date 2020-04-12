/*
 *    Copyright (C) 2014-2019 Christian Muehlhaeuser
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

// Package twitchbee is a Bee that can connect to Twitch.
package twitchbee

import (
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v2"
	"github.com/nicklaw5/helix"

	"github.com/muesli/beehive/bees"
)

// TwitchBee is a Bee that can connect to Twitch.
type TwitchBee struct {
	bees.Bee

	// channel signaling twitch connection status
	connectedState chan bool

	// setup Twitch client:
	chat     *twitch.Client
	client   *helix.Client
	channels []string

	username string
	password string

	clientId string
}

// Action triggers the action passed to it.
func (mod *TwitchBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		tos := []string{}
		text := ""
		action.Options.Bind("text", &text)

		for _, opt := range action.Options {
			if opt.Name == "channel" {
				tos = append(tos, opt.Value.(string))
			}
		}

		for _, recv := range tos {
			if recv == "*" {
				// special: send to all joined channels
				for _, to := range mod.channels {
					mod.chat.Say(to, text)
				}
			} else {
				// needs stripping hostname when sending to user!host
				if strings.Index(recv, "!") > 0 {
					recv = recv[0:strings.Index(recv, "!")]
				}

				mod.chat.Say(recv, text)
			}
		}

	case "join":
		for _, opt := range action.Options {
			if opt.Name == "channel" {
				mod.join(opt.Value.(string))
			}
		}
	case "part":
		for _, opt := range action.Options {
			if opt.Name == "channel" {
				mod.part(opt.Value.(string))
			}
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *TwitchBee) rejoin() {
	for _, channel := range mod.channels {
		mod.chat.Join(channel)
	}
}

func (mod *TwitchBee) join(channel string) {
	channel = strings.TrimSpace(channel)
	mod.chat.Join(channel)

	mod.channels = append(mod.channels, channel)
}

func (mod *TwitchBee) part(channel string) {
	channel = strings.TrimSpace(channel)
	mod.chat.Depart(channel)

	for k, v := range mod.channels {
		if v == channel {
			mod.channels = append(mod.channels[:k], mod.channels[k+1:]...)
			return
		}
	}
}

func (mod *TwitchBee) monitorFollows(eventChan chan bees.Event) {
	var twitchId string
	{
		resp, err := mod.client.GetUsers(&helix.UsersParams{
			Logins: []string{mod.username},
		})
		if err != nil || len(resp.Data.Users) != 1 {
			mod.LogErrorf("Failed retrieving user info from Twitch API: %v", err)
		}
		twitchId = resp.Data.Users[0].ID
	}

	follows := make(map[string]helix.UserFollow)
	var seeded bool
	for {
		var pagination string
		for {
			resp, err := mod.client.GetUsersFollows(&helix.UsersFollowsParams{
				After: pagination,
				First: 40,
				ToID:  twitchId,
			})
			if err != nil {
				mod.LogErrorf("Failed retrieving follows from Twitch API: %v", err)
				break
			}
			pagination = resp.Data.Pagination.Cursor

			for _, f := range resp.Data.Follows {
				if _, ok := follows[f.FromID]; !ok {
					follows[f.FromID] = f

					if seeded {
						ev := bees.Event{
							Bee:  mod.Name(),
							Name: "follow",
							Options: []bees.Placeholder{
								{
									Name:  "user",
									Type:  "string",
									Value: f.FromName,
								},
							},
						}
						eventChan <- ev
					}
				}
			}

			if len(resp.Data.Follows) == 0 || pagination == "" {
				// end of paginated list
				seeded = true
				break
			}
		}

		// poll once a minute
		time.Sleep(30 * time.Second)
	}
}

// Run executes the Bee's event loop.
func (mod *TwitchBee) Run(eventChan chan bees.Event) {
	// channel signaling Twitch connection status
	mod.connectedState = make(chan bool)

	// setup Twitch c:
	mod.chat = twitch.NewClient(mod.username, mod.password)

	mod.chat.OnConnect(func() {
		mod.connectedState <- true
	})

	mod.chat.OnUserJoinMessage(func(msg twitch.UserJoinMessage) {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "join",
			Options: []bees.Placeholder{
				{
					Name:  "channel",
					Type:  "string",
					Value: msg.Channel,
				},
				{
					Name:  "user",
					Type:  "string",
					Value: msg.User,
				},
			},
		}
		eventChan <- ev
	})
	mod.chat.OnUserPartMessage(func(msg twitch.UserPartMessage) {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "part",
			Options: []bees.Placeholder{
				{
					Name:  "channel",
					Type:  "string",
					Value: msg.Channel,
				},
				{
					Name:  "user",
					Type:  "string",
					Value: msg.User,
				},
			},
		}
		eventChan <- ev
	})

	mod.chat.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "message",
			Options: []bees.Placeholder{
				{
					Name:  "channel",
					Type:  "string",
					Value: msg.Channel,
				},
				{
					Name:  "user",
					Type:  "string",
					Value: msg.User,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: msg.Message,
				},
			},
		}
		eventChan <- ev
	})

	connected := false
	mod.ContextSet("connected", &connected)

	var err error
	mod.client, err = helix.NewClient(&helix.Options{
		ClientID: mod.clientId,
	})
	if err != nil {
		mod.LogErrorf("Failed connecting to Twitch API: %v", err)
		return
	}
	go mod.monitorFollows(eventChan)

	go func() {
		mod.Logln("Connecting to Twitch")
		mod.rejoin()
		err := mod.chat.Connect()
		if err != nil {
			mod.LogErrorf("Failed to connect to Twitch: %v", err)
		}
	}()

	for {
		select {
		case status := <-mod.connectedState:
			if status {
				mod.Logln("Connected to Twitch")
				connected = true
			} else {
				mod.Logln("Disconnected from Twitch")
				connected = false
			}

		case <-mod.SigChan:
			mod.chat.Disconnect()
			return

		default:
			time.Sleep(time.Second)
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TwitchBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("client_id", &mod.clientId)
	options.Bind("username", &mod.username)
	options.Bind("password", &mod.password)
	options.Bind("channels", &mod.channels)

	mod.ContextSet("channels", &mod.channels)
}
