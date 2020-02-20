/*
 *    Copyright (C) 2019 Sergio Rubio
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
 *      Sergio Rubio <sergio@rubio.im>
 */

package redisbee

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/muesli/beehive/bees"
	log "github.com/sirupsen/logrus"
)

// RedisBee  updates a Cloudflare domain name
type RedisBee struct {
	bees.Bee
	client *redis.Client
	domain string
}

// Run executes the Bee's event loop.
func (mod *RedisBee) Run(eventChan chan bees.Event) {

	rchann := mod.Config().Options.Value("channel").(string)
	if rchann == "" {
		log.Debugf("Redis channel not configured, disabling pubsub")
		return
	}

	log.Debugf("Redis: subscribed to channel '%s'", rchann)

	pubsub := mod.client.Subscribe(rchann)
	_, err := pubsub.Receive()
	if err != nil {
		mod.LogErrorf("Redis: error subscribing to channel '%s', disabling pubsub", rchann)
	}
	ch := pubsub.Channel()

	for {
		select {
		case <-mod.SigChan:
			return

		case msg := <-ch:
			sendEvent(mod.Name(), msg.Channel, msg.Payload, eventChan)
		}
	}
}

func sendEvent(bee string, channel string, msg string, eventChan chan bees.Event) {
	event := bees.Event{
		Bee:  bee,
		Name: "message",
		Options: []bees.Placeholder{
			{
				Name:  "channel",
				Type:  "string",
				Value: channel,
			},
			{
				Name:  "message",
				Type:  "string",
				Value: msg,
			},
		},
	}
	eventChan <- event
}

// Action triggers the action passed to it.
func (mod *RedisBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "set":
		mod.client.Set(action.Options.Value("key").(string), action.Options.Value("value").(string), 0).Err()
	default:
		mod.LogDebugf("Unknown action triggered in %s: %s", mod.Name(), action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *RedisBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	var host, port, password string
	options.Bind("host", &host)
	if host == "" {
		host = "localhost"
	}
	options.Bind("port", &port)
	if port == "" {
		port = "6379"
	}
	options.Bind("password", &password)
	var db int
	options.Bind("db", &db)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	mod.client = client
}
