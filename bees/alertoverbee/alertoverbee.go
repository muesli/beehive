/*
 *    Copyright (C) 2017      Henson Lu
 *                  2015-2017 Christian Muehlhaeuser
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
 *      Henson Lu <henson.lu@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package alertoverbee is able to send notifications on AlertOver.
package alertoverbee

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/muesli/beehive/bees"
)

// AlertOverBee is a Bee that is able to send notifications on AlertOver.
type AlertOverBee struct {
	bees.Bee
	source string
}

// Action triggers the action passed to it.
func (mod *AlertOverBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		var receiver, title, content, weburl, priority string
		action.Options.Bind("receiver", &receiver)
		action.Options.Bind("title", &title)
		action.Options.Bind("content", &content)
		action.Options.Bind("url", &weburl)
		action.Options.Bind("priority", &priority)

		if priority == "" {
			priority = "0"
		}

		// the message must be plain text, so
		// remove the HTML tags, such as <html></html> and so on
		re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
		content = re.ReplaceAllString(content, "")

		data := url.Values{
			"source":   {mod.source},
			"receiver": {receiver},
			"title":    {title},
			"content":  {content},
			"url":      {weburl},
			"priority": {priority},
		}
		resp, err := http.PostForm("https://api.alertover.com/v1/alert", data)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			mod.Logln("AlertOver send message success.")
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *AlertOverBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	options.Bind("source", &mod.source)
}
