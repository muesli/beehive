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

// Package tumblrbee is a Bee that can post blogs & quotes on Tumblr.
package tumblrbee

import (
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

// Action triggers the action passed to it.
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

// Run executes the Bee's event loop.
func (mod *TumblrBee) Run(eventChan chan bees.Event) {
	mod.client = gotumblr.NewTumblrRestClient(mod.consumerKey, mod.consumerSecret,
		mod.token, mod.tokenSecret,
		mod.callbackUrl, "http://api.tumblr.com")
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TumblrBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("blogname", &mod.blogname)
	options.Bind("callback_url", &mod.callbackUrl)
	options.Bind("consumer_key", &mod.consumerKey)
	options.Bind("consumer_secret", &mod.consumerSecret)
	options.Bind("token", &mod.token)
	options.Bind("token_secret", &mod.tokenSecret)
}
