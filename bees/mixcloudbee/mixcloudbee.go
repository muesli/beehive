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

	client       mixcloud.Client
	lastUpdate   time.Time

	eventChan chan bees.Event
}

// Poll a Mixcloud Cloudcasts feed
func (mod *MixcloudBee) pollFeed(feed string) {
	for {
		time.Sleep(20*time.Second)
		var allCloudcastsData []mixcloud.CloudcastData
		var opt mixcloud.ListOptions
		opt.Since = mod.lastUpdate
		opt.Until = time.Now()
		mod.lastUpdate = opt.Until
		cloudcasts, err := mod.client.GetCloudcasts("zanjradio", &opt)
		if err != nil {
			continue
		}
		allCloudcastsData = append(allCloudcastsData, cloudcasts.Data...)
		nextUrl := cloudcasts.Paging.NextURL
		for {
			// the following line is necessary to always create a new object, else just some values would be overwritten
			// and missing values would stay the same as before
			cloudcasts = mixcloud.Cloudcasts{}
			err := mod.client.GetPage(nextUrl, &cloudcasts)
			if err != nil {
				continue
			}
			nextUrl = cloudcasts.Paging.NextURL
			if nextUrl == "" {
				break
			}
			allCloudcastsData = append(allCloudcastsData, cloudcasts.Data...)
		}
	}
}


// Run executes the Bee's event loop.
func (mod *MixcloudBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	mod.pollFeed(mod.feed)
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *MixcloudBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("baseUrl", &mod.baseUrl)
	options.Bind("feed", &mod.feed)

	&mod.client = mixcloud.NewClient(nil)
}
