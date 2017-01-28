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

// beehive's Jabber module.
package jabberbee

import (
	"log"
	"strings"

	"github.com/mattn/go-xmpp"
	"github.com/muesli/beehive/bees"
)

type JabberBee struct {
	bees.Bee

	client *xmpp.Client

	server   string
	user     string
	password string
	notls    bool
}

func (mod *JabberBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		chat := xmpp.Chat{Type: "chat"}

		action.Options.Bind("user", &chat.Remote)
		action.Options.Bind("text", &chat.Text)

		mod.client.Send(chat)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *JabberBee) Run(eventChan chan bees.Event) {
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
		select {
		case <-mod.SigChan:
			mod.client.Close()
			return

		default:
		}

		chat, err := mod.client.Recv()
		if err != nil {
			log.Fatal(err)
		}
		switch v := chat.(type) {
		case xmpp.Chat:
			if len(v.Text) > 0 {
				text := strings.TrimSpace(v.Text)

				ev := bees.Event{
					Bee:  mod.Name(),
					Name: "message",
					Options: []bees.Placeholder{
						{
							Name:  "user",
							Type:  "string",
							Value: v.Remote,
						},
						{
							Name:  "text",
							Type:  "string",
							Value: text,
						},
					},
				}
				eventChan <- ev
			}

		case xmpp.Presence:
			//fmt.Println(v.From, v.Show)
		}
	}
}

func (mod *JabberBee) ReloadOptions(options bees.BeeOptions) {
	//FIXME: implement this
	mod.SetOptions(options)
}
