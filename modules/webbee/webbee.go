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
)

var (
	eventChan chan modules.Event
)

type WebBee struct {
	addr string
}

func (mod *WebBee) Name() string {
	return "webbee"
}

func (mod *WebBee) Description() string {
	return "A RESTful HTTP module for beehive"
}

func (mod *WebBee) Run(cin chan modules.Event) {
	eventChan = cin
	go web.Run(mod.addr)
}

func (mod *WebBee) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   mod.Name(),
			Name:        "post",
			Description: "A POST call was received by the HTTP server",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "json",
					Description: "JSON map received from caller",
					Type:        "map",
				},
				modules.PlaceholderDescriptor{
					Name:        "ip",
					Description: "IP of the caller",
					Type:        "string",
				},
			},
		},
		modules.EventDescriptor{
			Namespace:   mod.Name(),
			Name:        "get",
			Description: "A GET call was received by the HTTP server",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "query_params",
					Description: "Map of query parameters received from caller",
					Type:        "map",
				},
				modules.PlaceholderDescriptor{
					Name:        "ip",
					Description: "IP of the caller",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (mod *WebBee) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{}
	return actions
}

func (mod *WebBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	return outs
}

func GetRequest(ctx *web.Context) {
	//FIXME
	ms := make(map[string]string)
	ev := modules.Event{
		Namespace: "webbee",
		Name:      "get",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "query_params",
				Type:  "map",
				Value: ms,
			},
			modules.Placeholder{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}
	eventChan <- ev
}

func PostRequest(ctx *web.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	var payload interface{}
	err = json.Unmarshal(b, &payload)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	ev := modules.Event{
		Namespace: "webbee",
		Name:      "post",
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
	eventChan <- ev
}

func init() {
	w := WebBee{
		addr: "0.0.0.0:12345",
	}
	web.Get("/event", GetRequest)
	web.Post("/event", PostRequest)

	modules.RegisterModule(&w)
}
