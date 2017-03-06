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

// TumblrBee is a Bee that can post blogs & quotes on Tumblr.
type TumblrBee struct {
	bees.Bee

	client *gotumblr.TumblrRestClient

	blogname string

	callbackURL    string
	consumerKey    string
	consumerSecret string
	token          string
	tokenSecret    string

	evchan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *TumblrBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "post_text":
		text := ""
		action.Options.Bind("text", &text)

		state := "published"
		err := mod.client.CreateText(mod.blogname, map[string]string{
			"body":  text,
			"state": state})
		if err != nil {
			mod.LogErrorf("Failed to post text: %v", err)
			return outs
		}
		ev := bees.Event{
			Bee:     mod.Name(),
			Name:    "posted",
			Options: []bees.Placeholder{},
		}
		mod.evchan <- ev

	case "post_quote":
		quote := ""
		source := ""
		action.Options.Bind("quote", &quote)
		action.Options.Bind("source", &source)

		state := "published"
		err := mod.client.CreateQuote(mod.blogname, map[string]string{
			"quote":  quote,
			"source": source,
			"state":  state})
		if err != nil {
			mod.LogErrorf("Failed to post quote: %v", err)
			return outs
		}
		ev := bees.Event{
			Bee:     mod.Name(),
			Name:    "posted",
			Options: []bees.Placeholder{},
		}
		mod.evchan <- ev

	case "follow":
		blogname := ""
		action.Options.Bind("blogname", &blogname)

		if err := mod.client.Follow(blogname); err != nil {
			mod.LogErrorf("Failed to follow blog: %v", err)
			return outs
		}

	case "unfollow":
		blogname := ""
		action.Options.Bind("blogname", &blogname)

		if err := mod.client.Unfollow(blogname); err != nil {
			mod.LogErrorf("Failed to unfollow blog: %v", err)
			return outs
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *TumblrBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TumblrBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("blogname", &mod.blogname)
	options.Bind("callback_url", &mod.callbackURL)
	options.Bind("consumer_key", &mod.consumerKey)
	options.Bind("consumer_secret", &mod.consumerSecret)
	options.Bind("token", &mod.token)
	options.Bind("token_secret", &mod.tokenSecret)

	mod.client = gotumblr.NewTumblrRestClient(mod.consumerKey, mod.consumerSecret,
		mod.token, mod.tokenSecret,
		mod.callbackURL, "http://api.tumblr.com")
}
