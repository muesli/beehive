/*
 *    Copyright (C) 2018 Stefan Derkits
 *                  2018 Christian Muehlhaeuser
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
 *      Stefan Derkits <stefan@derkits.at>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package mixcloudbee

import (
	"time"

	"github.com/horrendus/go-mixcloud/mixcloud"

	"github.com/muesli/beehive/bees"
)

// MixcloudBee is a Bee that can interface with Mixcloud.
type MixcloudBee struct {
	bees.Bee

	baseUrl      string
	feed         string

	client       *mixcloud.Client
	lastUpdate   time.Time

	eventChan chan bees.Event
}

// Poll a Mixcloud Cloudcasts feed
func (mod *MixcloudBee) pollFeed(feed string) {
	for {
		mod.Logln("Parsing feed")
		var allCloudcastsData []mixcloud.CloudcastData
		var opt mixcloud.ListOptions
		opt.Since = mod.lastUpdate
		opt.Until = time.Now()
		mod.lastUpdate = opt.Until
		mod.Logln("Since", opt.Since, "Until", opt.Until)
		cloudcasts, err := mod.client.GetCloudcasts(mod.feed, &opt)
		if err != nil {
			mod.LogErrorf("Error:", err)
			break
		}
		allCloudcastsData = append(allCloudcastsData, cloudcasts.Data...)
		nextUrl := cloudcasts.Paging.NextURL
		for {
			if nextUrl == "" {
				break
			}
			// the following line is necessary to always create a new object, else just some values would be overwritten
			// and missing values would stay the same as before
			cloudcasts = mixcloud.Cloudcasts{}
			err := mod.client.GetPage(nextUrl, &cloudcasts)
			allCloudcastsData = append(allCloudcastsData, cloudcasts.Data...)
			if err != nil {
				mod.LogErrorf("Error:", err)
				break
			}
			nextUrl = cloudcasts.Paging.NextURL
		}

		for _, cloudcast := range allCloudcastsData {
			newCloudcastEvent := bees.Event{
				Bee:  mod.Name(),
				Name: "new_cloudcast",
				Options: []bees.Placeholder{
					{
						Name:  "name",
						Type:  "string",
						Value: cloudcast.Name,
					},
					{
						Name:  "url",
						Type:  "string",
						Value: cloudcast.URL,
					},
				},
			}
			mod.Logln("Event new_cloudcast", newCloudcastEvent)
			mod.eventChan <- newCloudcastEvent
		}
		time.Sleep(25 * time.Second)
	}
}


// Run executes the Bee's event loop.
func (mod *MixcloudBee) Run(cin chan bees.Event) {
	mod.eventChan = cin
	mod.Logln("Mixcloudbee starting to run")
	mod.pollFeed(mod.feed)
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *MixcloudBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("baseUrl", &mod.baseUrl)
	options.Bind("feed", &mod.feed)

	mod.client = mixcloud.NewClient(nil)
}
