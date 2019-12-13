/*
 *    Copyright (C) 2019 Adaptant Solutions AG
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
 *      Paul Mundt <paul.mundt@adaptant.io>
 */

package knowgobee

import (
	"context"
	"encoding/json"
	"github.com/knowgoio/knowgo-pubsub/api"
	"github.com/muesli/beehive/bees"
	"net"
	"net/url"
	"strconv"
	"time"
)

type KnowGoBee struct {
	bees.Bee

	apiKey   string
	server   string
	client   *api.ClientConfig
	sub      *api.Subscription
	interval int

	eventChan chan bees.Event
}

type CountryChangeNotification struct {
	DriverID  int    `json:"driverId,omitempty"`
	Entering  string `json:"entering"`
	Leaving   string `json:"leaving"`
	Timestamp string `json:"timestamp"`
}

func (mod *KnowGoBee) handleEvent(eventData []byte) {
	var notification CountryChangeNotification

	err := json.Unmarshal(eventData, &notification)
	if err != nil {
		mod.LogErrorf("Failed to unmarshal event")
		return
	}

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "country-change",
		Options: []bees.Placeholder{
			{
				Name:  "driverId",
				Type:  "int",
				Value: notification.DriverID,
			},
			{
				Name:  "entering",
				Type:  "string",
				Value: notification.Entering,
			},
			{
				Name:  "leaving",
				Type:  "string",
				Value: notification.Leaving,
			},
			{
				Name:  "timestamp",
				Type:  "timestamp",
				Value: notification.Timestamp,
			},
		},
	}

	mod.eventChan <- ev
}

func (mod *KnowGoBee) newAPIClient() *api.ClientConfig {
	client := api.DefaultClientConfig()

	if mod.apiKey != "" {
		client.APIKey = mod.apiKey
	}

	if mod.interval > 0 {
		client.PollInterval = time.Duration(mod.interval) * time.Second
	}

	if mod.server != "" {
		u, err := url.Parse(mod.server)
		if err == nil {
			client.Host = u.Hostname()
			client.Port, _ = strconv.Atoi(u.Port())
		} else {
			host, port, err := net.SplitHostPort(mod.server)
			if err != nil {
				client.Host = mod.server
			} else {
				client.Host = host
				client.Port, _ = strconv.Atoi(port)
			}
		}
	}

	return client
}

// Run executes the Bee's event loop.
func (mod *KnowGoBee) Run(eventChan chan bees.Event) {
	client := mod.newAPIClient()

	sub, err := client.Subscribe(&api.SubscriptionRequest{
		Event: "country-change",
	})
	if err != nil {
		mod.LogErrorf("Failed to subscribe:", err)
		return
	}

	mod.Logf("Successfully subscribed")

	mod.client = client
	mod.sub = sub
	mod.eventChan = eventChan

	ticker := time.NewTicker(client.PollInterval).C
	for {
		select {
		case <-mod.SigChan:
			return
		case <-ticker:
			b := client.Receive(context.Background(), mod.sub)
			if b != nil {
				mod.handleEvent(b)
			}
		}
	}
}

// Action triggers the action passed to it.
func (mod *KnowGoBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *KnowGoBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("polling_interval", &mod.interval)
	options.Bind("server", &mod.server)
	options.Bind("api_key", &mod.apiKey)
}
