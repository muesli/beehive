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

// beehive's Tumblr module.
package tumblrbee

import (
	_ "log"

	"github.com/MariaTerzieva/gotumblr"
	"github.com/muesli/beehive/bees"
)

type TumblrBee struct {
	bees.Bee

	client *gotumblr.TumblrRestClient

	blogname string

	callbackUrl    string
	consumerKey    string
	consumerSecret string
	token          string
	tokenSecret    string
}

// Interface impl

func (mod *TumblrBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "postText":
		text := ""
		action.Options.Bind("text", &text)

		state := "published"
		mod.client.CreateText(mod.blogname, map[string]string{
			"body":  text,
			"state": state})

	case "postQuote":
		quote := ""
		source := ""
		action.Options.Bind("quote", &quote)
		action.Options.Bind("source", &source)

		state := "published"
		mod.client.CreateQuote(mod.blogname, map[string]string{
			"quote":  quote,
			"source": source,
			"state":  state})

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *TumblrBee) Run(eventChan chan bees.Event) {
	mod.client = gotumblr.NewTumblrRestClient(mod.consumerKey, mod.consumerSecret,
		mod.token, mod.tokenSecret,
		mod.callbackUrl, "http://api.tumblr.com")
}

func (mod *TumblrBee) ReloadOptions(options bees.BeeOptions) {
	//FIXME: implement this
	mod.SetOptions(options)
}
