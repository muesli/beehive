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
	room        string
}

// Action triggers the actions passed to it.
func (mod *GitterBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {

	case "sendMessage":
		var room string
		var message string

		action.Options.Bind("room", &room)
		action.Options.Bind("message", &message)

		roomID, err := mod.client.GetRoomId(room)
		if err != nil {
			mod.LogErrorf("Failed to fetch room ID from uri:", err)
		}

		if err = mod.client.SendMessage(roomID, message); err != nil {
			mod.LogErrorf("Failed to send message:", err)
		}

	case "joinRoom":
		var room string
		action.Options.Bind("room", &room)

		roomID, err := mod.client.GetRoomId(room)
		if err != nil {
			mod.LogErrorf("Failed to fetch room ID from uri: %v", err)
		}

		if _, err := mod.client.JoinRoom(roomID, mod.userID); err != nil {
			mod.LogErrorf("Failed to join room: %v", err)
		}
		mod.Logln("Successfully joined room:", room)

	case "leaveRoom":
		var room string
		action.Options.Bind("room", &room)

		roomID, err := mod.client.GetRoomId(room)
		if err != nil {
			mod.LogErrorf("Failed to fetch room ID from uri: %v", err)
		}

		if err := mod.client.LeaveRoom(roomID, mod.userID); err != nil {
			mod.LogErrorf("Failed to leave Room: %v", err)
		}
		mod.Logln("Successfully leaved room:", room)

	// case "getRoomMessages":
	// 	var room string
	// 	action.Options.Bind("room", &room)

	// 	roomID, err := mod.client.GetRoomId(room)
	// 	if err != nil {
	// 		mod.LogErrorf("Failed to fetch room ID from uri: %v", err)
	// 	}

	// 	mod.getRoomMessages(roomID)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *GitterBee) Run(eventChan chan bees.Event) {

	mod.client = gitter.New(mod.accessToken)
	user, err := mod.client.GetUser()
	if err != nil {
		mod.LogErrorf("Failed to fetch current user: %v", err)
	}
	mod.userID = user.ID

	// Getting the roomID from the uri as a string in order to pass it to getRoomMessages()
	// roomID, err := mod.client.GetRoomId(mod.room)
	// if err != nil {
	// 	mod.LogErrorf("Failed to fetch room ID: %v", err)
	// }

	// mod.eventChan = eventChan

	// timeout := time.Duration(time.Second * 10) // TODO: Look after api limits!
	// for {
	// 	select {
	// 	case <-mod.SigChan:
	// 		return
	// 	case <-time.After(timeout):
	// 		mod.getRoomMessages(roomID)
	// 	}
	// 	timeout = time.Duration(time.Second * 10)
	// }
}

// getRoomMessages receives unread messages
// func (mod *GitterBee) getRoomMessages(room string) {

// 	messages, err := mod.client.GetMessages(room, nil) // TODO: Look after params, maybe theres something useful
// 	if err != nil {
// 		mod.LogErrorf("Failed to fetch messages: %v", err)
// 	}
// 	// Parsing the messages into the bees event chan
// 	for _, v := range messages {
// 		if v.Unread == true {
// 			mod.Logln("Unread messages in room:", room)
// 			ev := bees.Event{
// 				Bee:  mod.Name(),
// 				Name: "roomMessages",
// 				Options: []bees.Placeholder{
// 					{
// 						Name:  "ID",
// 						Type:  "string",
// 						Value: v.ID,
// 					},
// 					{
// 						Name:  "test",
// 						Type:  "string",
// 						Value: v.Text,
// 					},
// 					{
// 						Name:  "username",
// 						Type:  "string",
// 						Value: v.From.Username,
// 					},
// 					{
// 						Name:  "readBy",
// 						Type:  "int",
// 						Value: v.ReadBy,
// 					},
// 				},
// 			}
// 			mod.eventChan <- ev

// 			if v.Mentions != nil {
// 				for _, mention := range v.Mentions {
// 					ev := bees.Event{
// 						Bee:  mod.Name(),
// 						Name: "mention",
// 						Options: []bees.Placeholder{
// 							{
// 								Name:  "mention",
// 								Type:  "string",
// 								Value: mention.ScreenName,
// 							},
// 						},
// 					}
// 					mod.eventChan <- ev
// 				}
// 			}

// 			if v.Issues != nil {
// 				for _, issue := range v.Issues {
// 					ev := bees.Event{
// 						Bee:  mod.Name(),
// 						Name: "issue",
// 						Options: []bees.Placeholder{
// 							{
// 								Name:  "issue",
// 								Type:  "int",
// 								Value: issue.Number,
// 							},
// 						},
// 					}
// 					mod.eventChan <- ev
// 				}
// 			}
// 		}
// 	}
// }

// ReloadOptions parses the config options and initializes the Bee.
func (mod *GitterBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("accessToken", &mod.accessToken)
	options.Bind("room", &mod.room)
}
