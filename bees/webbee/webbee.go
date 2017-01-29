/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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
 */

// Package webbee is a Bee that starts an HTTP server and fires events for
// incoming requests.
package webbee

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hoisie/web"
	"github.com/muesli/beehive/bees"
)

type WebBee struct {
	bees.Bee

	addr string
	path string

	eventChan chan bees.Event
}

func (mod *WebBee) triggerJsonEvent(resp *[]byte) {
	var payload interface{}
	err := json.Unmarshal(*resp, &payload)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "post",
		Options: []bees.Placeholder{
			{
				Name:  "json",
				Type:  "map",
				Value: payload,
			},
			{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}

	j := make(map[string]interface{})
	err = json.Unmarshal(*resp, &j)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	for k, v := range j {
		log.Printf("POST JSON param: %s = %+v\n", k, v)

		ph := bees.Placeholder{
			Name:  k,
			Type:  "string",
			Value: v,
		}
		ev.Options = append(ev.Options, ph)
	}

	mod.eventChan <- ev
}

func (mod *WebBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	web.Get(mod.path, mod.GetRequest)
	web.Post(mod.path, mod.PostRequest)

	web.Run(mod.addr)

	for {
		select {
		case <-mod.SigChan:
			web.Close()
			return

		default:
		}
	}
}

// Action triggers the action passed to it.
func (mod *WebBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "post":
		url := ""
		j := ""
		action.Options.Bind("url", &url)
		action.Options.Bind("json", &j)

		buf := strings.NewReader(j)
		resp, err := http.Post(url, "application/json", buf)
		if err != nil {
			log.Println("Error:", err)
			return outs
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error:", err)
			return outs
		}

		mod.triggerJsonEvent(&b)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *WebBee) GetRequest(ctx *web.Context) {
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "get",
		Options: []bees.Placeholder{
			{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}

	for k, v := range ctx.Params {
		log.Println("GET param:", k, "=", v)

		ph := bees.Placeholder{
			Name:  k,
			Type:  "string",
			Value: v,
		}
		ev.Options = append(ev.Options, ph)
	}

	mod.eventChan <- ev
}

func (mod *WebBee) PostRequest(ctx *web.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	mod.triggerJsonEvent(&b)
}

func (mod *WebBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("addr", &mod.addr)
	options.Bind("path", &mod.path)
}
