package cricketbee

/*
 *    Copyright (C) 2017 Akash Shinde
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
 *      Akash Shinde <akashshinde159@gmail.com>
 */

import (
	"github.com/akashshinde/go_cricket"
	"github.com/muesli/beehive/bees"
)

type CricketBee struct {
	bees.Bee
	favTeam string
	event   chan gocricket.ResponseEvent
	beeEvt  chan bees.Event
}

func (c *CricketBee) placeholderOptions(r gocricket.ResponseEvent) []bees.Placeholder {
	return []bees.Placeholder{
		{
			Name:  "batting_team",
			Type:  "string",
			Value: r.BtTeamName,
		},
		{
			Name:  "score",
			Type:  "string",
			Value: r.Runs,
		},
		{
			Name:  "wickets",
			Type:  "string",
			Value: r.Wickets,
		},
		{
			Name:  "overs",
			Type:  "string",
			Value: r.Overs,
		},
	}
}

func (c *CricketBee) announceCricketEvent(response gocricket.ResponseEvent) {
	switch response.EventType {
	case gocricket.EVENT_OUT:
		c.beeEvt <- bees.Event{
			Name:    "out",
			Options: c.placeholderOptions(response),
			Bee:     c.Name(),
		}
	case gocricket.EVENT_OVER_CHANGED:
		c.beeEvt <- bees.Event{
			Name:    "over_changed",
			Options: c.placeholderOptions(response),
			Bee:     c.Name(),
		}
	case gocricket.EVENT_RUN_CHANGE:
		c.beeEvt <- bees.Event{
			Name:    "run_change",
			Options: c.placeholderOptions(response),
			Bee:     c.Name(),
		}
	}
}

func (c *CricketBee) Run(cin chan bees.Event) {
	evt := make(chan gocricket.ResponseEvent)
	// Start Cricket GoRoutine to poll cricket score
	cricket := gocricket.NewCricketWatcher(c.favTeam, evt)
	cricket.Start()
	c.beeEvt = cin
	for {
		select {
		case e := <-evt:
			c.announceCricketEvent(e)

		case <-c.SigChan:
			return
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *CricketBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	options.Bind("favourite_team", &mod.favTeam)
}
