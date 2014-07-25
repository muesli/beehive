/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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
 *      Christian Muehlhaeuser <muesli@gmail.com>
 *      Johannes Fürmann <johannes@weltraumpflege.org>
 */

// beehive's SpaceAPI module.
package twitterbee

import (
	"github.com/muesli/beehive/modules"
	"github.com/ChimeraCoder/anaconda"
)

type TwitterBee struct {
	modules.Module

	consumer_key string
	consumer_secret string
	access_token string
	access_token_secret string
	
	twitter_api *anaconda.TwitterApi

	evchan chan modules.Event
}

func (mod *TwitterBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	switch action.Name {
	case "tweet":
		


		ev := modules.Event{
			Bee:  mod.Name(),
			Name: "call_finished",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:  "success",
					Type:  "bool",
					Value: true,
				},
			},
		}
		mod.evchan <- ev
		
	default:
		panic("Unknown action triggered in " +mod.Name()+": "+action.Name)
	}
	
	return outs
}

func (mod *TwitterBee) Run(eventChan chan modules.Event) {
	mod.evchan = eventChan
	
	anaconda.SetConsumerKey(mod.consumer_key)
	anaconda.SetConsumerSecret(mod.consumer_secret)
	mod.twitter_api = anaconda.NewTwitterApi(mod.access_token, mod.access_token_secret)
}
