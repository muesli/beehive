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
 *      Martin Schlierf <go@koma666.de>
 */

// Package mumblebee is a Bee that can connect to a Mumble/XMPP server.
package mumblebee

import (
	"crypto/tls"
	"net"
	"strconv"
	"strings"

	"github.com/muesli/beehive/bees"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleutil"
)

// MumbleBee is a Bee that can connect to a Mumble server.
type MumbleBee struct {
	bees.Bee

	// channel signaling mumble connection status
	connectedState chan bool

	client *gumble.Client

	server   string
	user     string
	password string
	insecure bool
	channel  string
}

// Action triggers the action passed to it.
func (mod *MumbleBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send_to_user":

		user, text := "", ""

		action.Options.Bind("user", &user)
		action.Options.Bind("text", &text)

		mod.client.Users.Find(user).Send(text)

	case "send_to_channel":

		channel, text := "", ""

		action.Options.Bind("channel", &channel)
		action.Options.Bind("text", &text)

		if channel != "" {
			mod.client.Channels.Find(channel).Send(text, false)
		} else {
			mod.client.Self.Channel.Send(text, false)
		}
	case "move":

		channel := ""

		action.Options.Bind("channel", &channel)
		mod.move(channel)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *MumbleBee) move(channel string) {
	channel = strings.TrimSpace(channel)
	if len(channel) > 0 {
		channels := strings.Split(channel, "/")
		result := mod.client.Channels.Find(channels...)
		if result != nil {
			mod.client.Self.Move(result)
			mod.Logln("Moved to channel:", channel)
		} else {
			mod.Logln("Moved to channel:", channel, " failed. (channel not found)")
		}
	}
}

// Run executes the Bee's event loop.
func (mod *MumbleBee) Run(eventChan chan bees.Event) {
	if len(mod.server) == 0 {
		return
	}

	// channel signaling mumble connection status
	mod.connectedState = make(chan bool)

	host, port, err := net.SplitHostPort(mod.server)
	if err != nil {
		host = mod.server
		port = strconv.Itoa(gumble.DefaultPort)
	}

	config := gumble.NewConfig()
	config.Username = mod.user
	config.Password = mod.password
	address := net.JoinHostPort(host, port)

	var tlsConfig tls.Config

	if mod.insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	config.Attach(gumbleutil.Listener{
		TextMessage: func(e *gumble.TextMessageEvent) {
			ev := bees.Event{
				Bee:  mod.Name(),
				Name: "message",
				Options: []bees.Placeholder{
					{
						Name:  "channels",
						Type:  "[]string",
						Value: e.Channels,
					},
					{
						Name:  "user",
						Type:  "string",
						Value: e.Sender.Name,
					},
					{
						Name:  "text",
						Type:  "string",
						Value: e.Message,
					},
				},
			}
			eventChan <- ev
		},
	})
	config.Attach(gumbleutil.Listener{
		Disconnect: func(e *gumble.DisconnectEvent) {
			go func() {
				mod.connectedState <- false
			}()
		},
	})
	config.Attach(gumbleutil.Listener{
		Connect: func(e *gumble.ConnectEvent) {
			go func() {
				mod.connectedState <- true
			}()
		},
	})

	/*/ ToDo: Use own client certs
	if *certificateFile != "" {
		if *keyFile == "" {
			keyFile = certificateFile
		}
		if certificate, err := tls.LoadX509KeyPair(*certificateFile, *keyFile); err != nil {
			fmt.Printf("%s: %s\n", os.Args[0], err)
			os.Exit(1)
		} else {
			tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
		}
	}
	//*/

	config.Attach(gumbleutil.AutoBitrate)

	connecting := false
	disconnected := true
	waitForDisconnect := false
	for {
		// loop on mumble connection events
		if disconnected {
			if waitForDisconnect {
				return
			}

			if !connecting {
				connecting = true
				mod.Logln("Connecting to Mumble:", mod.server)
				mod.client, err = gumble.DialWithDialer(new(net.Dialer), address, config, &tlsConfig)

				mod.Logln("Connected to:", mod.server)
				if err != nil {
					mod.Logln("Failed to connect to Mumble:", mod.server, err)
					connecting = false
				}
			}
		}
		select {
		case status := <-mod.connectedState:
			if status {
				mod.Logln("Connected to Mumble:", mod.server)
				connecting = false
				disconnected = false
				mod.move(mod.channel)
			} else {
				mod.Logln("Disconnected from Mumble:", mod.server)
				connecting = false
				disconnected = true
			}
		case <-mod.SigChan:
			if !waitForDisconnect {
				mod.client.Disconnect()
			}
			waitForDisconnect = true
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *MumbleBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.server)
	options.Bind("channel", &mod.channel)
	options.Bind("user", &mod.user)
	options.Bind("password", &mod.password)
	options.Bind("insecure", &mod.insecure)
}
