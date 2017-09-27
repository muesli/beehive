/*
 *    Copyright (C) 2017 Dominik Schmidt
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
 *      Dominik Schmidt <dev@dominik-schmidt.de>
 */

package horizonboxbee

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/muesli/beehive/bees"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type HorizonBoxBeeState struct {
	online bool
	ip     string
}

type HorizonBoxBee struct {
	bees.Bee

	address  string
	user     string
	password string

	State     HorizonBoxBeeState
	eventChan chan bees.Event
}

func (mod *HorizonBoxBee) poll() {
	// Form data
	v := url.Values{}
	v.Set("username", mod.user)
	v.Set("password", mod.password)
	v.Set("page", "login")

	req, err := http.NewRequest("POST", "http://"+mod.address+"/cgi-bin/sendResult.cgi?section=login", strings.NewReader(v.Encode()))
	if err != nil {
		mod.LogErrorf("http.NewRequest() error: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	cookieJar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: cookieJar,
	}

	resp, err := c.Do(req)
	if err != nil {
		mod.LogErrorf("Could not connect to: %s", mod.address, err)
		return
	}

	state := HorizonBoxBeeState{}
	doc, err := goquery.NewDocumentFromResponse(resp)
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(j int, t *goquery.Selection) {
			// Unfortunately keys here are translated, thus we rely on the order instead of keys
			// key := t.Find("td").First().Text()
			value := t.Find("td").Last().Text()

			switch j {
			case 8:
				state.online = false
				if value == "Enabled" {
					state.online = true
				}
			case 9:
				state.ip = value
			}
		})
	})

	if mod.State.online != state.online {
		mod.announceOnlineStateChange(state.online)
	}

	if mod.State.ip != state.ip {
		mod.announceIpChange(state.ip)
	}

	mod.State = state
}

func (mod *HorizonBoxBee) announceIpChange(ip string) {
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "external_ip_change",
		Options: []bees.Placeholder{
			{
				Name:  "new_external_ip",
				Type:  "string",
				Value: ip,
			},
		},
	}
	mod.eventChan <- event
}

func (mod *HorizonBoxBee) announceOnlineStateChange(online bool) {
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "connection_status_change",
		Options: []bees.Placeholder{
			{
				Name:  "online",
				Type:  "bool",
				Value: online,
			},
		},
	}
	mod.eventChan <- event
}

// Run executes the Bee's event loop.
func (mod *HorizonBoxBee) Run(cin chan bees.Event) {
	mod.eventChan = cin
	for {
		select {
		case <-mod.SigChan:
			return

		case <-time.After(time.Duration(10 * time.Second)):
			mod.poll()
		}
	}
}

// Action triggers the action passed to it.
func (mod *HorizonBoxBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *HorizonBoxBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.address)
	options.Bind("user", &mod.user)
	options.Bind("password", &mod.password)
}
