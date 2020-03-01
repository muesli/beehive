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
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-redis/redis"
	"github.com/muesli/beehive/bees"
	log "github.com/sirupsen/logrus"
)

// RedisBee is able to pubsub to Redis and store key/values
type RedisBee struct {
	bees.Bee
	client  *redis.Client
	channel string
	bkoff   *backoff.ExponentialBackOff
}

// Run executes the Bee's event loop.
func (mod *RedisBee) Run(eventChan chan bees.Event) {

	if mod.channel == "" {
		log.Debugf("Redis channel not configured, disabling pubsub")
		return
	}

	log.Debugf("Redis: subscribed to channel '%s'", mod.channel)

	operation := func() error {
		pubsub := mod.client.Subscribe(mod.channel)
		_, err := pubsub.Receive()
		if err != nil {
			mod.LogErrorf("Redis: error subscribing to channel '%s', disabling pubsub. Redis error: %v", mod.channel, err)
			return err
		}
		ch := pubsub.Channel()
		for {
			select {
			case <-mod.SigChan:
				return nil
			case msg := <-ch:
				sendEvent(mod.Name(), msg.Channel, msg.Payload, eventChan)
			}
		}
	}

	err := backoff.Retry(operation, backoff.WithMaxRetries(mod.bkoff, 10))
	if err != nil {
		mod.LogErrorf("Redis connection failed. Exhausted reconnection attempts.")
		return
	}
}

func (mod *RedisBee) publishWithBackoff(message string) error {
	operation := func() error {
		_, err := mod.client.Publish(mod.channel, message).Result()
		return err
	}

	err := backoff.Retry(operation, backoff.WithMaxRetries(mod.bkoff, 10))
	if err != nil {
		mod.LogErrorf("Redis connection failed. Exhausted reconnection attempts.")
		return err
	}

	return nil
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
	case "publish":
		//_, err := mod.client.Publish(mod.channel, action.Options.Value("message").(string)).Result()
		err := mod.publishWithBackoff(action.Options.Value("message").(string))
		if err != nil {
			mod.LogErrorf("Redis: error publishing message to channel. Redis error: %v", err)
		}
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

	var channel string
	options.Bind("channel", &channel)
	mod.channel = channel

	mod.bkoff = backoff.NewExponentialBackOff()
	mod.bkoff.InitialInterval = 1 * time.Second
	mod.bkoff.Multiplier = 2
}
