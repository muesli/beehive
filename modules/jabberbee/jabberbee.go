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

// beehive's Jabber module.
package jabberbee

import (
	"fmt"
	"github.com/mattn/go-xmpp"
	"github.com/muesli/beehive/app"
	"github.com/muesli/beehive/modules"
	"log"
)

type JabberBee struct {
	client *xmpp.Client

	server   string
	username string
	password string
	notls    bool
}

func (sys *JabberBee) Name() string {
	return "jabberbee"
}

func (mod *JabberBee) Description() string {
	return "A Jabber module for beehive"
}

func (mod *JabberBee) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   mod.Name(),
			Name:        "message",
			Description: "A message was received over Jabber",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "user",
					Description: "The user that sent the message",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (mod *JabberBee) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   mod.Name(),
			Name:        "send",
			Description: "Sends a message to a remote",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "user",
					Description: "Which remote to send the message to",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "text",
					Description: "Content of the message",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func (mod *JabberBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}

func (mod *JabberBee) Run(eventChan chan modules.Event) {
	if len(mod.server) == 0 {
		return
	}

	var talk *xmpp.Client
	var err error
	if mod.notls {
		talk, err = xmpp.NewClientNoTLS(mod.server, mod.username, mod.password, false)
	} else {
		talk, err = xmpp.NewClient(mod.server, mod.username, mod.password, false)
	}
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			chat, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				if len(v.Text) > 0 {
					fmt.Println(v.Remote, v.Text)

					ev := modules.Event{
						Namespace: mod.Name(),
						Name:      "message",
						Options: []modules.Placeholder{
							modules.Placeholder{
								Name:  "user",
								Type:  "string",
								Value: v.Remote,
							},
							modules.Placeholder{
								Name:  "text",
								Type:  "string",
								Value: v.Text,
							},
						},
					}
					eventChan <- ev
				}

			case xmpp.Presence:
				fmt.Println(v.From, v.Show)
			}
		}
	}()
}

func init() {
	jabber := JabberBee{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&jabber.server, "jabberhost", "", "Hostname of Jabber server, eg: talk.google.com:443"},
		app.CliFlag{&jabber.username, "jabberuser", "beehive", "Username to authenticate with Jabber server"},
		app.CliFlag{&jabber.password, "jabberpassword", "", "Password to use to connect to Jabber server"},
		app.CliFlag{&jabber.notls, "jabbernotls", false, "If you don't want to connect with TLS"},
	})

	modules.RegisterModule(&jabber)
}
