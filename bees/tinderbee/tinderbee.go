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

// Package tinderbee is a Bee that can post blogs & quotes on Tinder.
package tinderbee

import (
	"github.com/mnzt/tinder"
	"github.com/muesli/beehive/bees"
)

// TinderBee is a Bee that can post blogs & quotes on Tinder.
type TinderBee struct {
	bees.Bee

	client *tinder.Tinder

	userID    string
	userToken string
	evchan    chan bees.Event
}

// Action triggers the action passed to it.
func (mod *TinderBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "get_updates":
		var limit int
		action.Options.Bind("limit", &limit)

		updt, err := mod.client.GetUpdates(limit)
		if err != nil {
			mod.LogErrorf("Failed to fetch updates: %v", err)
			return nil
		}

		mod.TriggerUpdateEvent(&updt)

	case "get_user":
		var userID string
		action.Options.Bind("user_id", &userID)

		user, err := mod.client.GetUser(userID)
		if err != nil {
			mod.LogErrorf("Failed to get user: %v", err)
			return nil
		}

		mod.TriggerUserEvents(&user)

	case "send_message":
		var userID string
		var text string
		action.Options.Bind("user_id", &userID)
		action.Options.Bind("text", &text)

		err := mod.client.SendMessage(userID, text)
		if err != nil {
			mod.LogErrorf("Failed to send message: %v", err)
			return nil
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "message_sent",
		}
		mod.evchan <- ev

	case "get_recommendations":
		var limit int
		action.Options.Bind("limit", &limit)

		rec, err := mod.client.GetRecommendations(limit)
		if err != nil {
			mod.LogErrorf("Failed to fetch recommendations: %v", err)
			return nil
		}

		mod.TriggerRecommendationsEvent(&rec)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *TinderBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	mod.client = tinder.Init(mod.userID, mod.userToken)
	if err := mod.client.Auth(); err != nil {
		mod.LogErrorf("Failed to Authenticate: %v", err)
		return
	}

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TinderBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("user_id", &mod.userID)
	options.Bind("user_token", &mod.userToken)
}
