/*
 *    Copyright (C) 2020      deranjer
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
 *      deranjer <deranjer@gmail.com>
 */

// Package gotifybee is able to send notifications on Gotify.
package gotifybee

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/muesli/beehive/bees"
)

// GotifyBee is a Bee that is able to send notifications on Gotify.
type GotifyBee struct {
	bees.Bee
	token string
	serverURL string
}

// Action triggers the action passed to it.
func (mod *GotifyBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		var title, message, priority string
		action.Options.Bind("title", &title)
		action.Options.Bind("message", &message)
		action.Options.Bind("priority", &priority)

		if priority == "" {
			priority = "0"
		}
		if title == "" {
			title = "Gotify"
		}

		// the message must be plain text, so
		// remove the HTML tags, such as <html></html> and so on
		re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
		message = re.ReplaceAllString(message, "\n")

		data := url.Values{
			"title":    {title},
			"message":  {message},
			"priority": {priority},
		}
		// Build the URL for sending the message
		rawURL := mod.serverURL + "message?token=" + mod.token
		resp, err := http.PostForm(rawURL, data)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			mod.Logln("Gotify send message success.")
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *GotifyBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	options.Bind("token", &mod.token)
	options.Bind("serverURL", &mod.serverURL)
}
