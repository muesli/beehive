/*
 *    Copyright (C) 2017 Christian Muehlhaeuser
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
 *      Nicolas Martin <penguwingithub@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package gitterbee is a Bee that can interface with Gitter
package gitterbee

import (
	gitter "github.com/sromku/go-gitter"

	"github.com/muesli/beehive/bees"
)

// GitterBee is a Bee that can interface with Gitter
type GitterBee struct {
	bees.Bee

	eventChan chan bees.Event
	client    *gitter.Gitter
	userID    string

	accessToken string
	rooms       []string
	roomChans   map[string]chan interface{}
}

// Action triggers the actions passed to it.
func (mod *GitterBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {

	case "send":
		var room string
		var message string

		action.Options.Bind("room", &room)
		action.Options.Bind("message", &message)

		roomID, err := mod.client.GetRoomId(room)
		if err != nil {
			mod.LogErrorf("Failed to fetch room ID from uri:", err)
			return outs
		}

		if err = mod.client.SendMessage(roomID, message); err != nil {
			mod.LogErrorf("Failed to send message:", err)
			return outs
		}

	case "join":
		var room string
		action.Options.Bind("room", &room)

		mod.join(room)

	case "leave":
		var room string
		action.Options.Bind("room", &room)

		ch, ok := mod.roomChans[room]
		if !ok {
			mod.LogErrorf("Can't leave this room: %s", room)
			return outs
		}

		mod.Logln("Closing room stream", room)
		close(ch)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *GitterBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan

	mod.client = gitter.New(mod.accessToken)
	user, err := mod.client.GetUser()
	if err != nil {
		mod.LogErrorf("Failed to fetch current user: %v", err)
		return
	}
	mod.userID = user.ID

	for _, room := range mod.rooms {
		mod.join(room)
	}

	select {
	case <-mod.SigChan:
		for room, ch := range mod.roomChans {
			mod.Logln("Closing room stream", room)
			close(ch)
		}
		return
	}
}

func (mod *GitterBee) join(room string) error {
	mod.Logln("Joining room", room)

	roomID, err := mod.client.GetRoomId(room)
	if err != nil {
		mod.LogErrorf("Failed to fetch room ID: %v", err)
		return err
	}
	r, err := mod.client.JoinRoom(roomID, mod.userID)
	if err != nil {
		mod.LogErrorf("Failed to join room: %v", err)
		return err
	}

	sigchan := make(chan interface{})
	mod.roomChans[room] = sigchan

	stream := mod.client.Stream(r.ID)
	defer stream.Close()
	go mod.client.Listen(stream)

	go func() {
		for {
			select {
			case <-sigchan:
				mod.Logln("Exiting", room)
				return

			case event := <-stream.Event:
				switch ev := event.Data.(type) {
				case *gitter.MessageReceived:
					mod.handleMessage(room, &ev.Message)
				case *gitter.GitterConnectionClosed:
					// connection was closed
					mod.LogErrorf("Connection closed")
					go mod.join(room)
					return
				}
			}
		}
	}()

	return nil
}

func (mod *GitterBee) handleMessage(room string, v *gitter.Message) {
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "message",
		Options: []bees.Placeholder{
			{
				Name:  "id",
				Type:  "string",
				Value: v.ID,
			},
			{
				Name:  "text",
				Type:  "string",
				Value: v.Text,
			},
			{
				Name:  "username",
				Type:  "string",
				Value: v.From.Username,
			},
			{
				Name:  "room",
				Type:  "string",
				Value: room,
			},
			{
				Name:  "read_by",
				Type:  "int",
				Value: v.ReadBy,
			},
		},
	}
	mod.eventChan <- ev
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *GitterBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("access_token", &mod.accessToken)
	options.Bind("rooms", &mod.rooms)

	mod.roomChans = make(map[string]chan interface{})
}
