/*
 *    Copyright (C) 2016 Sergio Rubio
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
 *      Sergio Rubio <rubiojr@frameos.org>
 */

package slackbee

import (
	"log"

	"github.com/muesli/beehive/bees"
	"github.com/nlopes/slack"
)

type SlackBee struct {
	bees.Bee

	client   *slack.Client
	channels map[string]string
	apiKey   string
}

func (mod *SlackBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		tos := []string{}
		text := ""
		action.Options.Bind("text", &text)

		for _, opt := range action.Options {
			if opt.Name == "channel" {
				cid := mod.findChannelID(opt.Value.(string), false)
				if cid != "" {
					tos = append(tos, cid)
				} else {
					log.Printf("Slack: channel ID for %s not found\n", opt.Value.(string))
				}
			}
		}

		msgParams := slack.NewPostMessageParameters()
		for _, to := range tos {
			_, _, err := mod.client.PostMessage(to, text, msgParams)
			if err != nil {
				log.Println("Slack: error posting message to the slack channel " + to)
			}
		}
	default:
		log.Printf("Slack: unknown action triggered in %s: %s\n", mod.Name(), action.Name)
	}
	return outs
}

func stringInMap(a string, list map[string]string) bool {
	for _, v := range list {
		if v == a {
			return true
		}
	}
	return false
}

func (mod *SlackBee) findChannelID(name string, cache bool) string {
	cid := mod.channels[name]

	if cid != "" {
		return cid
	}

	channels, err := mod.client.GetChannels(true)
	if err != nil {
		panic(err)
	}

	for _, ch := range channels {
		if ch.Name == name {
			cid = ch.ID
		}
	}

	// Channel not found, try to find a group
	groups, err := mod.client.GetGroups(true)
	if err != nil {
		panic(err)
	}
	for _, grp := range groups {
		if grp.Name == name {
			cid = grp.ID
		}
	}

	if cache {
		mod.channels[name] = cid
	}
	log.Println("Channel map " + name + " " + cid)

	return cid
}

func sendEvent(bee string, channel string, user string, text string, eventChan chan bees.Event) {
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
				Name:  "user",
				Type:  "string",
				Value: user,
			},
			{
				Name:  "text",
				Type:  "string",
				Value: text,
			},
		},
	}
	eventChan <- event
}

func (mod *SlackBee) Run(eventChan chan bees.Event) {
	rtm := mod.client.NewRTM()

	go rtm.ManageConnection()

	// Cache the channels
	for k := range mod.channels {
		mod.findChannelID(k, true)
	}

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				if stringInMap(ev.Channel, mod.channels) {
					u := ev.Msg.User
					if u == "" {
						u = ev.Msg.Username
					}
					t := ev.Msg.Text
					if t == "" {
						for _, v := range ev.Msg.Attachments {
							sendEvent(mod.Name(), ev.Channel, u, v.Text, eventChan)
						}
					} else {
						sendEvent(mod.Name(), ev.Channel, u, t, eventChan)
					}
				}
			case *slack.RTMError:
				log.Printf("Slack: error %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				log.Println("Slack: invalid credentials")
				break Loop
			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
