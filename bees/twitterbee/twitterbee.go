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
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/muesli/beehive/bees"
)

type TwitterBee struct {
	bees.Bee

	consumer_key        string
	consumer_secret     string
	access_token        string
	access_token_secret string

	twitter_api      *anaconda.TwitterApi
	twitter_mentions []anaconda.Tweet

	evchan chan bees.Event
}

func handle_anaconda_error(err error, msg string) {
	if err != nil {
		is_rate_limit_error, next_window := err.(*anaconda.ApiError).RateLimitCheck()
		if is_rate_limit_error {
			log.Println("Oops, I exceeded the API rate limit!")
			wait_period := next_window.Sub(time.Now())
			log.Printf("waiting %f seconds to next window!\n", wait_period.Seconds())
			time.Sleep(wait_period)
		} else {
			if msg != "" {
				panic(msg)
			}
		}
	}
}

func (mod *TwitterBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "tweet":
		status := ""
		action.Options.Bind("status", &status)

		posted_tweet := false
		for !posted_tweet {
			v := url.Values{}

			for _, mention := range mod.twitter_mentions {
				tmp_mention_time, _ := mention.CreatedAtTime()
				if strings.Contains(status, "@"+mention.User.ScreenName) && time.Now().Sub(tmp_mention_time).Hours() < 2 {
					log.Printf("This might be a reply to " + mention.User.ScreenName)
					v.Set("in_reply_to_status_id", mention.IdStr)
					break
				}
			}
			posted_tweet = true

			log.Printf("Attempting to paste \"%s\" to Twitter", status)
			_, err := mod.twitter_api.PostTweet(status, v)
			if err != nil {
				log.Printf("Error posting to twitter %v", err)
				handle_anaconda_error(err, "")
			} else {
				posted_tweet = true
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

func (mod *TwitterBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan

	anaconda.SetConsumerKey(mod.consumer_key)
	anaconda.SetConsumerSecret(mod.consumer_secret)
	mod.twitter_api = anaconda.NewTwitterApi(mod.access_token, mod.access_token_secret)
	mod.twitter_api.ReturnRateLimitError(true)

	// Test the credentials on startup
	credentials_verified := false
	for !credentials_verified {
		ok, err := mod.twitter_api.VerifyCredentials()
		handle_anaconda_error(err, "Could not verify Twitter API Credentials")
		credentials_verified = ok
	}

	// populate mentions initially
	mentions_populated := false
	for !mentions_populated {
		v := url.Values{}
		v.Set("count", "30")

		log.Println("Populating Mentions...")
		mentions, err := mod.twitter_api.GetMentionsTimeline(v)
		handle_anaconda_error(err, "Could not populate mentions initially")
		if err == nil {
			mentions_populated = true
		}
		mod.twitter_mentions = mentions
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
		new_mentions, err := mod.twitter_api.GetMentionsTimeline(v)
		if err != nil {
			panic("Error: Could not get mentions")
		}

		// check if newest new mention is newer than newest old
		newest_new_time, _ := new_mentions[0].CreatedAtTime()
		newest_old_time, _ := mod.twitter_mentions[0].CreatedAtTime()

		if newest_new_time.After(newest_old_time) {
			log.Println("New mentions found!")
			for i := 0; ; i++ {
				tmp_mention := new_mentions[i]
				tmp_mention_time, _ := tmp_mention.CreatedAtTime()
				if tmp_mention_time.After(newest_old_time) {
					ev := bees.Event{
						Bee:  mod.Name(),
						Name: "mention",
						Options: []bees.Placeholder{
							{
								Name:  "username",
								Type:  "string",
								Value: tmp_mention.User.ScreenName,
							},
							{
								Name:  "text",
								Type:  "string",
								Value: tmp_mention.Text,
							},
						},
					}

					mod.evchan <- ev

				} else {
					break
				}
			}
			mod.twitter_mentions = new_mentions
		}
	}
}

func (mod *TwitterBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("consumer_key", &mod.consumer_key)
	options.Bind("consumer_secret", &mod.consumer_secret)
	options.Bind("access_token", &mod.access_token)
	options.Bind("access_token_secret", &mod.access_token_secret)
}
