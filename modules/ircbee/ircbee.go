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
 */

// beehive's IRC module.
package ircbee

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/muesli/beehive/modules"
	"log"
	"strings"
	"time"
)

type IrcBee struct {
	name        string
	namespace   string
	description string

	// channel signaling irc connection status
	connectedState chan bool

	// setup IRC client:
	client   *irc.Conn
	channels []string

	Server   string
	Nick     string
	Password string
	SSL      bool
	Channel  string
}

// Interface impl

func (mod *IrcBee) Name() string {
	return mod.name
}

func (mod *IrcBee) Namespace() string {
	return mod.namespace
}

func (mod *IrcBee) Description() string {
	return mod.description
}

func (mod *IrcBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}

	switch action.Name {
	case "send":
		tos := []string{}
		text := ""

		for _, opt := range action.Options {
			if opt.Name == "channel" {
				tos = append(tos, opt.Value.(string))
			}
			if opt.Name == "text" {
				text = opt.Value.(string)
			}
		}

		for _, recv := range tos {
			if recv == "*" {
				// special: send to all joined channels
				for _, to := range mod.channels {
					mod.client.Privmsg(to, text)
				}
			} else {
				// needs stripping hostname when sending to user!host
				if strings.Index(recv, "!") > 0 {
					recv = recv[0:strings.Index(recv, "!")]
				}

				mod.client.Privmsg(recv, text)
			}
		}

	case "join":
		for _, opt := range action.Options {
			if opt.Name == "channel" {
				mod.Join(opt.Value.(string))
			}
		}
	case "part":
		for _, opt := range action.Options {
			if opt.Name == "channel" {
				mod.Part(opt.Value.(string))
			}
		}
	default:
		// unknown action
		return outs
	}

	return outs
}

// ircbee specific impl

func (mod *IrcBee) Rejoin() {
	for _, channel := range mod.channels {
		mod.client.Join(channel)
	}
}

func (mod *IrcBee) Join(channel string) {
	channel = strings.TrimSpace(channel)
	mod.client.Join(channel)

	mod.channels = append(mod.channels, channel)
}

func (mod *IrcBee) Part(channel string) {
	channel = strings.TrimSpace(channel)
	mod.client.Part(channel)

	for k, v := range mod.channels {
		if v == channel {
			mod.channels = append(mod.channels[:k], mod.channels[k+1:]...)
			return
		}
	}
}

func (mod *IrcBee) Run(eventChan chan modules.Event) {
	if len(mod.Server) == 0 {
		return
	}

	// channel signaling IRC connection status
	mod.connectedState = make(chan bool)

	// setup IRC client:
	mod.client = irc.SimpleClient(mod.Nick, "beehive", "beehive")
	mod.client.SSL = mod.SSL

	mod.client.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		mod.connectedState <- true
	})
	mod.client.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		mod.connectedState <- false
	})
	mod.client.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		if channel == mod.client.Me.Nick {
			channel = line.Src // replies go via PM too.
		}
		msg := ""
		if len(line.Args) > 1 {
			msg = line.Args[1]
		}

		ev := modules.Event{
			Bee:  mod.Name(),
			Name: "message",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:  "channel",
					Type:  "string",
					Value: channel,
				},
				modules.Placeholder{
					Name:  "user",
					Type:  "string",
					Value: line.Src,
				},
				modules.Placeholder{
					Name:  "text",
					Type:  "string",
					Value: msg,
				},
			},
		}
		eventChan <- ev
	})

	// loop on IRC dis/connected events
	go func() {
		for {
			log.Println("Connecting to IRC:", mod.Server)
			err := mod.client.Connect(mod.Server, mod.Password)
			if err != nil {
				log.Println("Failed to connect to IRC:", mod.Server)
				log.Println(err)
				continue
			}
			for {
				status := <-mod.connectedState
				if status {
					log.Println("Connected to IRC:", mod.Server)

					if len(mod.channels) == 0 {
						// join default channel
						mod.Join(mod.Channel)
					} else {
						// we must have been disconnected, rejoin channels
						mod.Rejoin()
					}
				} else {
					log.Println("Disconnected from IRC:", mod.Server)
					break
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
