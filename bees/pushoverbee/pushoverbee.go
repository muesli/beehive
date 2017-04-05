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
 *      Raphael Mutschler <info@raphaelmutschler.de>
 */

// Package pushoverbee is a Bee that can send pushover notifications.
package pushoverbee

import (
	"net/http"
	"net/url"

	"github.com/muesli/beehive/bees"
)

// PushoverBee is a Bee that sends Pushover notifications
type PushoverBee struct {
	bees.Bee

	token     string
	userToken string
}

// Run executes the Bee's event loop.
func (mod *PushoverBee) Run(cin chan bees.Event) {
	select {
	case <-mod.SigChan:
		return
	}
}

// Action triggers the action passed to it.
func (mod *PushoverBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		var message, title, nurl, urlTitle string
		action.Options.Bind("message", &message)
		action.Options.Bind("title", &title)
		action.Options.Bind("url", &nurl)
		action.Options.Bind("url_title", &urlTitle)

		msg := url.Values{}
		msg.Set("token", mod.token)
		msg.Set("user", mod.userToken)
		msg.Set("message", message)

		if nurl != "" {
			msg.Set("url", nurl)
		}
		if urlTitle != "" {
			msg.Set("url_title", urlTitle)
		}
		if title != "" {
			msg.Set("title", title)
		}

		mod.Logln(msg)
		resp, err := http.PostForm("https://api.pushover.net/1/messages.json", msg)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			mod.Logln("Pushover send message success.")
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *PushoverBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	options.Bind("token", &mod.token)
	options.Bind("user_token", &mod.userToken)
}
