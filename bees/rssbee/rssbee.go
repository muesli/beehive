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
	"github.com/muesli/beehive/bees"
	"os"
	"time"
)

type RSSBee struct {
	bees.Module

	url         string

	eventChan chan bees.Event
}

func (mod *RSSBee) pollFeed(uri string, timeout int) {
	feed := rss.New(timeout, true, mod.chanHandler, mod.itemHandler)

	for {
		select {
			case <-mod.SigChan:
				return

			default:
		}

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

		newitemEvent := bees.Event{
			Bee:  mod.Name(),
			Name: "newitem",
			Options: []bees.Placeholder{
				bees.Placeholder{
					Name:  "title",
					Type:  "string",
					Value: newitems[i].Title,
				},
				bees.Placeholder{
					Name:  "links",
					Type:  "[]string",
					Value: links,
				},
				bees.Placeholder{
					Name:  "description",
					Type:  "string",
					Value: newitems[i].Description,
				},
				bees.Placeholder{
					Name:  "author",
					Type:  "string",
					Value: newitems[i].Author.Name,
				},
				bees.Placeholder{
					Name:  "categories",
					Type:  "[]string",
					Value: categories,
				},
				bees.Placeholder{
					Name:  "comments",
					Type:  "string",
					Value: newitems[i].Comments,
				},
				bees.Placeholder{
					Name:  "enclosures",
					Type:  "[]string",
					Value: enclosures,
				},
				bees.Placeholder{
					Name:  "guid",
					Type:  "string",
					Value: newitems[i].Guid,
				},
				bees.Placeholder{
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
	}
	fmt.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
}

func (mod *RSSBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	time.Sleep(10 * time.Second)
	mod.pollFeed(mod.url, 5)
}

func (mod *RSSBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}
