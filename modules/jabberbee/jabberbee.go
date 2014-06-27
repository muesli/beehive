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
	"github.com/mattn/go-xmpp"
	"github.com/muesli/beehive/modules"
	"log"
	"strings"
)

type JabberBee struct {
	name        string
	namespace   string
	description string

	client *xmpp.Client

	server   string
	user     string
	password string
	notls    bool
}

func (mod *JabberBee) Name() string {
	return mod.name
}

func (mod *JabberBee) Namespace() string {
	return mod.namespace
}

func (mod *JabberBee) Description() string {
	return mod.description
}

func (mod *JabberBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}

	switch action.Name {
	case "send":
		chat := xmpp.Chat{Type: "chat"}
		for _, opt := range action.Options {
			if opt.Name == "user" {
				chat.Remote = opt.Value.(string)
			}
			if opt.Name == "text" {
				chat.Text = opt.Value.(string)
			}
		}

		mod.client.Send(chat)

	default:
		// unknown action
		return outs
	}

	return outs
}

func (mod *JabberBee) Run(eventChan chan modules.Event) {
	if len(mod.server) == 0 {
		return
	}

	var err error
	options := xmpp.Options{
		Host:     mod.server,
		User:     mod.user,
		Password: mod.password,
		NoTLS:    mod.notls,
		Debug:    false,
	}

	mod.client, err = options.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	for {
		chat, err := mod.client.Recv()
		if err != nil {
			log.Fatal(err)
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			if len(v.Text) > 0 {
				text := strings.TrimSpace(v.Text)

				ev := modules.Event{
					Bee:  mod.Name(),
					Name: "message",
					Options: []modules.Placeholder{
						modules.Placeholder{
							Name:  "user",
							Type:  "string",
							Value: v.Remote,
						},
						modules.Placeholder{
							Name:  "text",
							Type:  "string",
							Value: text,
						},
					},
				}
				eventChan <- ev
			}

		case xmpp.Presence:
			//				fmt.Println(v.From, v.Show)
		}
	}
}
