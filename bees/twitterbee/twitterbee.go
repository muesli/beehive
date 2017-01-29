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
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"

	"github.com/muesli/beehive/bees"
)

// TwitterBee is a Bee that can interface with Twitter.
type TwitterBee struct {
	bees.Bee

	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string

	twitterAPI      *anaconda.TwitterApi
	twitterMentions []anaconda.Tweet

	evchan chan bees.Event
}

func handleAnacondaError(err error, msg string) {
	if err != nil {
		isRateLimitError, nextWindow := err.(*anaconda.ApiError).RateLimitCheck()
		if isRateLimitError {
			log.Println("Oops, I exceeded the API rate limit!")
			waitPeriod := nextWindow.Sub(time.Now())
			log.Printf("waiting %f seconds to next window!\n", waitPeriod.Seconds())
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

		postedTweet := false
		for !postedTweet {
			v := url.Values{}

			for _, mention := range mod.twitterMentions {
				tmpMentionTime, _ := mention.CreatedAtTime()
				if strings.Contains(status, "@"+mention.User.ScreenName) && time.Now().Sub(tmpMentionTime).Hours() < 2 {
					log.Printf("This might be a reply to " + mention.User.ScreenName)
					v.Set("in_reply_to_status_id", mention.IdStr)
					break
				}
			}
			postedTweet = true

			log.Printf("Attempting to paste \"%s\" to Twitter", status)
			_, err := mod.twitterAPI.PostTweet(status, v)
			if err != nil {
				log.Printf("Error posting to twitter %v", err)
				handleAnacondaError(err, "")
			} else {
				postedTweet = true
			}
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

	// Test the credentials on startup
	credentialsVerified := false
	for !credentialsVerified {
		ok, err := mod.twitterAPI.VerifyCredentials()
		handleAnacondaError(err, "Could not verify Twitter API Credentials")
		credentialsVerified = ok
	}

	// populate mentions initially
	mentionsPopulated := false
	for !mentionsPopulated {
		v := url.Values{}
		v.Set("count", "30")

		log.Println("Populating Mentions...")
		mentions, err := mod.twitterAPI.GetMentionsTimeline(v)
		handleAnacondaError(err, "Could not populate mentions initially")
		if err == nil {
			mentionsPopulated = true
		}
		mod.twitterMentions = mentions
	}

	// check twitter mentions every 60 seconds
	for {
		//FIXME: don't block
		select {
		case <-mod.SigChan:
			return

		default:
		}
		time.Sleep(2 * time.Minute)

		log.Println("Checking for new mentions...")
		v := url.Values{}
		v.Set("count", "30")
		newMentions, err := mod.twitterAPI.GetMentionsTimeline(v)
		if err != nil {
			panic("Error: Could not get mentions")
		}

		// check if newest new mention is newer than newest old
		newestNewTime, _ := newMentions[0].CreatedAtTime()
		newestOldTime, _ := mod.twitterMentions[0].CreatedAtTime()

		if newestNewTime.After(newestOldTime) {
			log.Println("New mentions found!")
			for i := 0; ; i++ {
				tmpMention := newMentions[i]
				tmpMentionTime, _ := tmpMention.CreatedAtTime()
				if tmpMentionTime.After(newestOldTime) {
					ev := bees.Event{
						Bee:  mod.Name(),
						Name: "mention",
						Options: []bees.Placeholder{
							{
								Name:  "username",
								Type:  "string",
								Value: tmpMention.User.ScreenName,
							},
							{
								Name:  "text",
								Type:  "string",
								Value: tmpMention.Text,
							},
						},
					}

					mod.evchan <- ev

				} else {
					break
				}
			}
			mod.twitterMentions = newMentions
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
