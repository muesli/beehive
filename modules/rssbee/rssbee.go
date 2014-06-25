/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 */

// RSS module for beehive 
package rssbee

import (
    "fmt"
    "os"
    "time"
	"github.com/muesli/beehive/app"
	"github.com/muesli/beehive/modules"
    rss "github.com/jteeuwen/go-pkg-rss"
)

var (
    eventChan chan modules.Event
)

type RSSBee struct {
	url string
}

func (mod *RSSBee) Name() string {
	return "rssbee"
}

func (mod *RSSBee) Description() string {
	return "A bee that manages RSS-feeds"
}

func (mod *RSSBee) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
        modules.EventDescriptor{
            Namespace: mod.Name(),
            Name: "newitem",
            Description: "A new item has been received through the Feed",
            Options: []modules.PlaceholderDescriptor{
                modules.PlaceholderDescriptor{
                    Name:           "title",
                    Description:    "Title of the Item",
                    Type:           "string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "links",
                    Description:    "Links referenced by the Item",
                    Type:           "[]string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "description",
                    Description:    "Description of the Item",
                    Type:           "string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "author",
                    Description:    "The person who wrote the Item",
                    Type:           "string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "categories",
                    Description:    "Categories that the Item belongs to",
                    Type:           "[]string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "comments",
                    Description:    "Comments of the Item",
                    Type:           "string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "enclosures",
                    Description:    "Enclosures related to Item",
                    Type:           "[]string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "guid",
                    Description:    "Global unique ID attached to the Item",
                    Type:           "string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "pubdate",
                    Description:    "Date the Item was published on",
                    Type:           "string",
                },
                modules.PlaceholderDescriptor{
                    Name:           "source",
                    Description:    "Source of the Item",
                    Type:           "string",
                },
            },
        },
    }
	return events
}

func (mod *RSSBee) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{}
	return actions
}


func PollFeed(uri string, timeout int) {
    feed := rss.New(timeout, true, chanHandler, itemHandler)

    for {
            if err := feed.Fetch(uri, nil); err != nil {
                fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
                return
            }

        <-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
    }
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
    //fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
    for i := range(newitems) {
        var links []string
        var categories []string
        var enclosures []string

        for j := range(newitems[i].Links) {
            links = append(links, newitems[i].Links[j].Href)
        }

        for j := range(newitems[i].Categories) {
            categories = append(categories, newitems[i].Categories[j].Text)
        }

        for j := range(newitems[i].Enclosures) {
            enclosures = append(enclosures, newitems[i].Enclosures[j].Url)
        }

        newitemEvent := modules.Event{
            Namespace: "rssbee",
            Name: "newitem",
            Options: []modules.Placeholder{
                modules.Placeholder {
                    Name:   "title",
                    Type:   "string",
                    Value:  newitems[i].Title,
                },
                modules.Placeholder{
                    Name:   "links",
                    Type:   "[]string",
                    Value:  links,
                },
                modules.Placeholder{
                    Name:   "description",
                    Type:   "string",
                    Value:  newitems[i].Description,
                },
                modules.Placeholder{
                    Name:   "author",
                    Type:   "string",
                    Value:  newitems[i].Author.Name,
                },
                modules.Placeholder{
                    Name:   "categories",
                    Type:   "[]string",
                    Value:  categories,
                },
                modules.Placeholder{
                    Name:   "comments",
                    Type:   "string",
                    Value:  newitems[i].Comments,
                },
                modules.Placeholder{
                    Name:   "enclosures",
                    Type:   "[]string",
                    Value:  enclosures,
                },
                modules.Placeholder{
                    Name:   "guid",
                    Type:   "string",
                    Value:  newitems[i].Guid,
                },
                modules.Placeholder{
                    Name:   "pubdate",
                    Type:   "string",
                    Value:  newitems[i].PubDate,
                },
                modules.Placeholder{
                    Name:   "source",
                    Type:   "string",
                    Value:  newitems[i].Source.Url,
                },
            },
        }
        eventChan <- newitemEvent
    }
    fmt.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
}

func (mod *RSSBee) Run(eventChan chan modules.Event) {
    go PollFeed(mod.url, 5)
}

func (mod *RSSBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}

func init() {
	rssbee := RSSBee{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&rssbee.url, "url", "", "The URL, this bee can find the RSS-Feed at"},
	})

	modules.RegisterModule(&rssbee)
}
