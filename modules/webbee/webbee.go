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
 */

// beehive's web-module.
package webbee

import (
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/muesli/beehive/modules"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type WebBee struct {
	modules.Module

	addr        string
	path		string

	eventChan chan modules.Event
}

func (mod *WebBee) triggerJsonEvent(resp *[]byte) {
	var payload interface{}
	err := json.Unmarshal(*resp, &payload)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	ev := modules.Event{
		Bee:  mod.Name(),
		Name: "post",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "json",
				Type:  "map",
				Value: payload,
			},
			modules.Placeholder{
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

		ph := modules.Placeholder{
			Name: k,
			Type: "string",
			Value: v,
		}
		ev.Options = append(ev.Options, ph)
	}

	mod.eventChan <- ev
}

func (mod *WebBee) Run(cin chan modules.Event) {
	mod.eventChan = cin

	web.Get(mod.path, mod.GetRequest)
	web.Post(mod.path, mod.PostRequest)

	web.Run(mod.addr)
}

func (mod *WebBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}

	switch action.Name {
	case "post":
		url := ""
		j := ""

		for _, opt := range action.Options {
			if opt.Name == "url" {
				url = opt.Value.(string)
			}
			if opt.Name == "json" {
				j = opt.Value.(string)
			}
		}

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

		log.Printf("Response: %+v\n", b)
		mod.triggerJsonEvent(&b)

	default:
		panic("Unknown action triggered in " +mod.Name()+": "+action.Name)
	}

	return outs
}

func (mod *WebBee) GetRequest(ctx *web.Context) {
	ev := modules.Event{
		Bee:  mod.Name(),
		Name: "get",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}

	for k, v := range ctx.Params {
		log.Println("GET param:", k, "=", v)

		ph := modules.Placeholder{
			Name: k,
			Type: "string",
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
