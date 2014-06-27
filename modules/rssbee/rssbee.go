/*
 *	  Copyright (C) 2014 Daniel 'grindhold' Brendle
 *                  2014 Christian Muehlhaeuser
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

// RSS module for beehive.
package rssbee

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/muesli/beehive/modules"
	"os"
	"time"
)

type RSSBee struct {
	name        string
	namespace   string
	description string
	url         string

	eventChan chan modules.Event
}

func (mod *RSSBee) Name() string {
	return mod.name
}

func (mod *RSSBee) Namespace() string {
	return mod.namespace
}

func (mod *RSSBee) Description() string {
	return mod.description
}

func (mod *RSSBee) PollFeed(uri string, timeout int) {
	feed := rss.New(timeout, true, mod.chanHandler, mod.itemHandler)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

func (mod *RSSBee) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	//fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (mod *RSSBee) itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
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

		newitemEvent := modules.Event{
			Bee:  mod.Name(),
			Name: "newitem",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:  "title",
					Type:  "string",
					Value: newitems[i].Title,
				},
				modules.Placeholder{
					Name:  "links",
					Type:  "[]string",
					Value: links,
				},
				modules.Placeholder{
					Name:  "description",
					Type:  "string",
					Value: newitems[i].Description,
				},
				modules.Placeholder{
					Name:  "author",
					Type:  "string",
					Value: newitems[i].Author.Name,
				},
				modules.Placeholder{
					Name:  "categories",
					Type:  "[]string",
					Value: categories,
				},
				modules.Placeholder{
					Name:  "comments",
					Type:  "string",
					Value: newitems[i].Comments,
				},
				modules.Placeholder{
					Name:  "enclosures",
					Type:  "[]string",
					Value: enclosures,
				},
				modules.Placeholder{
					Name:  "guid",
					Type:  "string",
					Value: newitems[i].Guid,
				},
				modules.Placeholder{
					Name:  "pubdate",
					Type:  "string",
					Value: newitems[i].PubDate,
				},
			},
		}
		if newitems[i].Source != nil {
			ph := modules.Placeholder{
				Name:  "source",
				Type:  "string",
				Value: newitems[i].Source.Url,
			}

			newitemEvent.Options = append(newitemEvent.Options, ph)
		}

		mod.eventChan <- newitemEvent
	}
	fmt.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
}

func (mod *RSSBee) Run(cin chan modules.Event) {
	mod.eventChan = cin

	go func(){
		time.Sleep(10 * time.Second)
		mod.PollFeed(mod.url, 5)
	}()
}

func (mod *RSSBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}
