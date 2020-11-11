/*
 *    Copyright (C) 2014      Daniel 'grindhold' Brendle
 *                  2014-2017 Christian Muehlhaeuser
 *
 *	  This program is free software: you can redistribute it and/or modify
 *	  it under the terms of the GNU Affero General Public License as published
 *	  by the Free Software Foundation, either version 3 of the License, or
 *	  (at your option) any later version.
 *
 *	  This program is distributed in the hope that it will be useful,
 *	  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *	  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.	See the
 *	  GNU Affero General Public License for more details.
 *
 *	  You should have received a copy of the GNU Affero General Public License
 *	  along with this program.	If not, see <http://www.gnu.org/licenses/>.
 *
 *	  Authors:
 *		Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package rssbee is a Bee for handling RSS feeds.
package rssbee

import (
	"time"

	rss "github.com/muesli/go-pkg-rss"

	"github.com/muesli/beehive/bees"
)

// RSSBee is a Bee for handling RSS feeds.
type RSSBee struct {
	bees.Bee

	url string
	// decides whether the next fetch should be skipped
	skipNextFetch bool

	// decides whether only the newest item from the next fetch should to get through
	skipNextFetchAllowNewest bool

	eventChan chan bees.Event
}

func (mod *RSSBee) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	//fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (mod *RSSBee) itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	if mod.skipNextFetch == true {
		mod.skipNextFetch = false
		return
	}
	for i := range newitems {
		var links []string
		var categories []string
		var enclosures []string

		for j := range newitems[i].Links {
			links = append(links, newitems[i].Links[j].Href)
		}

		for j := range newitems[i].Categories {
			categories = append(categories, newitems[i].Categories[j].Text)
		}

		for j := range newitems[i].Enclosures {
			enclosures = append(enclosures, newitems[i].Enclosures[j].Url)
		}

		newitemEvent := bees.Event{
			Bee:  mod.Name(),
			Name: "new_item",
			Options: []bees.Placeholder{
				{
					Name:  "title",
					Type:  "string",
					Value: newitems[i].Title,
				},
				{
					Name:  "links",
					Type:  "[]string",
					Value: links,
				},
				{
					Name:  "description",
					Type:  "string",
					Value: newitems[i].Description,
				},
				{
					Name:  "author",
					Type:  "string",
					Value: newitems[i].Author.Name,
				},
				{
					Name:  "categories",
					Type:  "[]string",
					Value: categories,
				},
				{
					Name:  "comments",
					Type:  "string",
					Value: newitems[i].Comments,
				},
				{
					Name:  "enclosures",
					Type:  "[]string",
					Value: enclosures,
				},
				{
					Name:  "guid",
					Type:  "string",
					Value: newitems[i].Guid,
				},
				{
					Name:  "pubdate",
					Type:  "string",
					Value: newitems[i].PubDate,
				},
			},
		}
		if newitems[i].Source != nil {
			ph := bees.Placeholder{
				Name:  "source",
				Type:  "string",
				Value: newitems[i].Source.Url,
			}

			newitemEvent.Options = append(newitemEvent.Options, ph)
		}

		mod.eventChan <- newitemEvent

		if mod.skipNextFetchAllowNewest == true {
			// the first time only let the newest item pass
			mod.skipNextFetchAllowNewest = false
			break
		}
	}
	mod.Logf("%d new item(s) in %s", len(newitems), feed.Url)
}

func (mod *RSSBee) pollFeed(uri string, timeout int) {
	feed := rss.New(timeout, true, mod.chanHandler, mod.itemHandler)

	wait := time.Duration(0)
	for {
		select {
		case <-mod.SigChan:
			return

		case <-time.After(wait):
			if err := feed.Fetch(uri, nil); err != nil {
				mod.LogErrorf("%s: %s", uri, err)
			}
		}

		wait = time.Duration(feed.SecondsTillUpdate() * 1e9)
	}
}

// Run executes the Bee's event loop.
func (mod *RSSBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	time.Sleep(10 * time.Second)
	mod.pollFeed(mod.url, 5)
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *RSSBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("skip_first", &mod.skipNextFetch)
	options.Bind("skip_first_allow_newest", &mod.skipNextFetchAllowNewest)
	options.Bind("url", &mod.url)
}
