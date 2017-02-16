/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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

// Package ircbee is a Bee that can connect to an IRC server.
package ircbee

import (
	"strings"

	irc "github.com/fluffle/goirc/client"

	"github.com/muesli/beehive/bees"
)

// IrcBee is a Bee that can connect to an IRC server.
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

// Action triggers the action passed to it.
func (mod *IrcBee) Action(action bees.Action) []bees.Placeholder {
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

func (mod *IrcBee) rejoin() {
	for _, channel := range mod.channels {
		mod.client.Join(channel)
	}
}

func (mod *IrcBee) join(channel string) {
	channel = strings.TrimSpace(channel)
	mod.client.Join(channel)

	mod.channels = append(mod.channels, channel)
}

func (mod *IrcBee) part(channel string) {
	channel = strings.TrimSpace(channel)
	mod.client.Part(channel)

	for k, v := range mod.channels {
		if v == channel {
			mod.channels = append(mod.channels[:k], mod.channels[k+1:]...)
			return
		}
	}
}

// Run executes the Bee's event loop.
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
			channel = line.Src // replies go via PM too
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
				{
					Name:  "channel",
					Type:  "string",
					Value: channel,
				},
				{
					Name:  "user",
					Type:  "string",
					Value: user,
				},
				{
					Name:  "hostmask",
					Type:  "string",
					Value: hostmask,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: msg,
				},
			},
		}
		eventChan <- ev
	})

	connecting := false
	disconnected := true
	waitForDisconnect := false
	for {
		// loop on IRC connection events
		if disconnected {
			if waitForDisconnect {
				return
			}

			if !connecting {
				connecting = true
				mod.Logln("Connecting to IRC:", mod.server)
				err := mod.client.Connect()
				if err != nil {
					mod.Logln("Failed to connect to IRC:", mod.server, err)
					connecting = false
				}
			}
		}
		select {
		case status := <-mod.connectedState:
			if status {
				mod.Logln("Connected to IRC:", mod.server)
				connecting = false
				disconnected = false
				mod.rejoin()
			} else {
				mod.Logln("Disconnected from IRC:", mod.server)
				connecting = false
				disconnected = true
			}

		case <-mod.SigChan:
			if !waitForDisconnect {
				mod.client.Quit()
			}
			waitForDisconnect = true
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *IrcBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.server)
	options.Bind("nick", &mod.nick)
	options.Bind("password", &mod.password)
	options.Bind("ssl", &mod.ssl)

	mod.channels = []string{}
	for _, channel := range options.Value("channels").([]interface{}) {
		mod.channels = append(mod.channels, channel.(string))
	}
}
