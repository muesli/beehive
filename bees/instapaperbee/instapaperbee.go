/*
 *    Copyright (C) 2019 Adam Petrovic
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
 *      Adam Petrovic <adam@petrovic.com.au>
 */

package instapaperbee

import (
	"net/http"
	"net/url"

	"github.com/muesli/beehive/bees"
)

type InstapaperBee struct {
	bees.Bee

	username string
	password string
}

func (mod *InstapaperBee) Action(action bees.Action) []bees.Placeholder {
	switch action.Name {
	case "save":
		var title, page_url string
		action.Options.Bind("title", &title)
		action.Options.Bind("url", &page_url)

		msg := url.Values{}
		msg.Set("username", mod.username)
		msg.Set("password", mod.password)
		msg.Set("url", page_url)

		if title != "" {
			msg.Set("title", title)
		}
		mod.LogDebugf("Message: %s", msg)
		resp, err := http.PostForm("https://www.instapaper.com/api/add", msg)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			mod.LogDebugf("Added article to instapaper.")
		}
	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}
	return []bees.Placeholder{}
}

func (mod *InstapaperBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	options.Bind("username", &mod.username)
	options.Bind("password", &mod.password)
}
