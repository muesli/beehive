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

// Package jabberbee is a Bee that can connect to a Jabber/XMPP server.
package jabberbee

import (
	"strings"

	"github.com/mattn/go-xmpp"

	"github.com/muesli/beehive/bees"
)

// JabberBee is a Bee that can connect to a Jabber/XMPP server.
type JabberBee struct {
	bees.Bee

	client *xmpp.Client

	server   string
	user     string
	password string
	notls    bool
}

// Action triggers the action passed to it.
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

func (mod *JabberBee) handleEvents(eventChan chan bees.Event) error {
	chat, err := mod.client.Recv()
	if err != nil {
		return err
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

	return nil
}

// Run executes the Bee's event loop.
func (mod *JabberBee) Run(eventChan chan bees.Event) {
	if len(mod.server) == 0 {
		return
	}

	options := xmpp.Options{
		Host:     mod.server,
		User:     mod.user,
		Password: mod.password,
		NoTLS:    mod.notls,
		Debug:    false,
	}

	var err error
	mod.client, err = options.NewClient()
	if err != nil {
		mod.LogErrorf("Connection error: %s", err)
		return
	}
	defer mod.client.Close()

	go func() {
		for {
			if err := mod.handleEvents(eventChan); err != nil {
				return
			}
		}
	}()

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *JabberBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.server)
	options.Bind("user", &mod.user)
	options.Bind("password", &mod.password)
	options.Bind("notls", &mod.notls)
}
