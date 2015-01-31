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
	"github.com/muesli/beehive/bees"
	"log"
	"strings"
	"time"
)

type IrcBee struct {
	bees.Bee

	// channel signaling irc connection status
	connectedState chan bool

	// setup IRC client:
	client   *irc.Conn
	channels []string

	server   string
	nick     string
	password string
	ssl      bool
}

// Interface impl

func (mod *IrcBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

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
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
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

func (mod *IrcBee) Run(eventChan chan bees.Event) {
	if len(mod.server) == 0 {
		return
	}

	// channel signaling IRC connection status
	mod.connectedState = make(chan bool)

	// setup IRC client:
	cfg := irc.NewConfig(mod.nick, "beehive", "beehive")
	cfg.SSL = mod.ssl
	cfg.Server = mod.server
	cfg.Pass = mod.password
	cfg.NewNick = func(n string) string { return n + "_" }
	mod.client = irc.Client(cfg)

	mod.client.HandleFunc("connected", func(conn *irc.Conn, line *irc.Line) {
		mod.connectedState <- true
	})
	mod.client.HandleFunc("disconnected", func(conn *irc.Conn, line *irc.Line) {
		mod.connectedState <- false
	})
	mod.client.HandleFunc("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		if channel == mod.client.Config().Me.Nick {
			channel = line.Src // replies go via PM too.
		}
		msg := ""
		if len(line.Args) > 1 {
			msg = line.Args[1]
		}
		user := line.Src[:strings.Index(line.Src, "!")]
		hostmask := line.Src[strings.Index(line.Src, "!")+2:]

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "message",
			Options: []bees.Placeholder{
				bees.Placeholder{
					Name:  "channel",
					Type:  "string",
					Value: channel,
				},
				bees.Placeholder{
					Name:  "user",
					Type:  "string",
					Value: user,
				},
				bees.Placeholder{
					Name:  "hostmask",
					Type:  "string",
					Value: hostmask,
				},
				bees.Placeholder{
					Name:  "text",
					Type:  "string",
					Value: msg,
				},
			},
		}
		eventChan <- ev
	})

	// loop on IRC dis/connected events
	for {
		log.Println("Connecting to IRC:", mod.server)
		err := mod.client.Connect()
		if err != nil {
			log.Println("Failed to connect to IRC:", mod.server)
			log.Println(err)
		} else {
			disconnected := false
			for {
				if disconnected {
					break
				}
				select {
				case <-mod.SigChan:
					mod.client.Quit()
					return

				case status := <-mod.connectedState:
					if status {
						log.Println("Connected to IRC:", mod.server)
						mod.Rejoin()
					} else {
						log.Println("Disconnected from IRC:", mod.server)
						disconnected = true
						break
					}

				default:
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}
