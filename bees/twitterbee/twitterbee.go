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
 *      Johannes Fürmann <johannes@weltraumpflege.org>
 */

// Package twitterbee is a Bee that can interface with Twitter.
package twitterbee

import (
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"

	"github.com/muesli/beehive/bees"
)

// TwitterBee is a Bee that can interface with Twitter.
type TwitterBee struct {
	bees.Bee

	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string

	twitterAPI *anaconda.TwitterApi
	self       anaconda.User

	evchan chan bees.Event
}

func (mod *TwitterBee) handleAnacondaError(err error, msg string) {
	if err != nil {
		isRateLimitError, nextWindow := err.(*anaconda.ApiError).RateLimitCheck()
		if isRateLimitError {
			mod.Logln("Oops, I exceeded the API rate limit!")
			waitPeriod := nextWindow.Sub(time.Now())
			mod.Logf("waiting %f seconds to next window!\n", waitPeriod.Seconds())
			time.Sleep(waitPeriod)
		} else {
			if msg != "" {
				panic(msg)
			}
		}
	}
}

// Action triggers the action passed to it.
func (mod *TwitterBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "tweet":
		status := ""
		action.Options.Bind("status", &status)

		v := url.Values{}

		mod.Logf("Attempting to post \"%s\" to Twitter", status)
		_, err := mod.twitterAPI.PostTweet(status, v)
		if err != nil {
			mod.Logf("Error posting to twitter %v", err)
			mod.handleAnacondaError(err, "")
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "call_finished",
			Options: []bees.Placeholder{
				{
					Name:  "success",
					Type:  "bool",
					Value: true,
				},
			},
		}
		mod.evchan <- ev

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *TwitterBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	anaconda.SetConsumerKey(mod.consumerKey)
	anaconda.SetConsumerSecret(mod.consumerSecret)
	mod.twitterAPI = anaconda.NewTwitterApi(mod.accessToken, mod.accessTokenSecret)
	mod.twitterAPI.ReturnRateLimitError(true)
	defer mod.twitterAPI.Close()

	// Test the credentials on startup
	credentialsVerified := false
	for !credentialsVerified {
		ok, err := mod.twitterAPI.VerifyCredentials()
		mod.handleAnacondaError(err, "Could not verify Twitter API Credentials")
		credentialsVerified = ok
	}

	var err error
	mod.self, err = mod.twitterAPI.GetSelf(url.Values{})
	mod.handleAnacondaError(err, "Could not get own user object from Twitter API")

	mod.handleStream()
}

func (mod *TwitterBee) handleStreamEvent(item interface{}) {
	switch status := item.(type) {
	case anaconda.DirectMessage:
		// mod.Logf("DM: %s %s\n", status.Text, status.Sender.ScreenName)
	case anaconda.Tweet:
		// mod.Logf("Tweet: %+v %s %s\n", status, status.Text, status.User.ScreenName)

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "tweet",
			Options: []bees.Placeholder{
				{
					Name:  "username",
					Type:  "string",
					Value: status.User.ScreenName,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: status.Text,
				},
			},
		}

		for _, mention := range status.Entities.User_mentions {
			if mention.Screen_name == mod.self.ScreenName {
				ev.Name = "mention"
			}
		}

		mod.evchan <- ev

	case anaconda.EventTweet:
		// mod.Logf("Event Tweet: %+v\n", status)

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "",
			Options: []bees.Placeholder{
				{
					Name:  "username",
					Type:  "string",
					Value: status.Source.ScreenName,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: status.TargetObject.Text,
				},
			},
		}

		switch status.Event.Event {
		case "favorite":
			ev.Name = "like"
		case "unfavorite":
			ev.Name = "unlike"
		default:
			mod.Logln("Unhandled event type", status.Event.Event)
		}

		if ev.Name != "" {
			mod.evchan <- ev
		}

	case anaconda.LimitNotice:
		mod.Logf("Limit: %+v\n", status)
	case anaconda.DisconnectMessage:
		mod.Logf("Disconnect: %+v\n", status)
	case anaconda.StatusWithheldNotice:
		mod.Logf("Status Withheld: %+v\n", status)
	case anaconda.Event:
		mod.Logf("Event: %+v\n", status)
	default:
		// mod.Logf("Unhandled type %v\n", item)
	}
}

func (mod *TwitterBee) handleStream() {
	s := mod.twitterAPI.UserStream(url.Values{})

	for {
		select {
		case <-mod.SigChan:
			return
		case item := <-s.C:
			mod.handleStreamEvent(item)
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TwitterBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("consumer_key", &mod.consumerKey)
	options.Bind("consumer_secret", &mod.consumerSecret)
	options.Bind("access_token", &mod.accessToken)
	options.Bind("access_token_secret", &mod.accessTokenSecret)
}
