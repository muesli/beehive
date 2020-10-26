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
		switch e := err.(type) {
		case *anaconda.ApiError:
			isRateLimitError, nextWindow := e.RateLimitCheck()
			if isRateLimitError {
				mod.Logln("Oops, I exceeded the API rate limit!")
				waitPeriod := nextWindow.Sub(time.Now())
				mod.Logf("waiting %f seconds to next window!", waitPeriod.Seconds())
				time.Sleep(waitPeriod)
			} else {
				if msg != "" {
					mod.LogErrorf("Error: %s (%+v)", msg, err)
					panic(msg)
				}
			}
		default:
			mod.LogErrorf("Error: %s (%+v)", msg, err)
			panic(msg)
		}
	}
}

// Action triggers the action passed to it.
func (mod *TwitterBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "tweet":
		var status string
		action.Options.Bind("status", &status)
		mod.Logf("Attempting to post \"%s\" to Twitter", status)

		_, err := mod.twitterAPI.PostTweet(status, url.Values{})
		if err != nil {
			mod.Logf("Error posting to Twitter %v", err)
			mod.handleAnacondaError(err, "")
		}

	case "follow":
		var username string
		action.Options.Bind("username", &username)
		mod.Logf("Attempting to follow \"%s\" to Twitter", username)

		_, err := mod.twitterAPI.FollowUser(username)
		if err != nil {
			mod.Logf("Error following user on Twitter %v", err)
			mod.handleAnacondaError(err, "")
		}

	case "unfollow":
		var username string
		action.Options.Bind("username", &username)
		mod.Logf("Attempting to unfollow \"%s\" to Twitter", username)

		_, err := mod.twitterAPI.UnfollowUser(username)
		if err != nil {
			mod.Logf("Error unfollowing user on Twitter %v", err)
			mod.handleAnacondaError(err, "")
		}

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

	// wait for the bee to be shut down
	<-mod.SigChan
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TwitterBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("consumer_key", &mod.consumerKey)
	options.Bind("consumer_secret", &mod.consumerSecret)
	options.Bind("access_token", &mod.accessToken)
	options.Bind("access_token_secret", &mod.accessTokenSecret)
}
