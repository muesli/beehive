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
	"net"
	"net/http"
	"net/url"

	"github.com/muesli/beehive/bees"
)

// WebBee is a Bee that starts an HTTP server and fires events for incoming
// requests.
type WebBee struct {
	bees.Bee

	addr string

	eventChan chan bees.Event
}

// Run executes the Bee's event loop.
func (mod *WebBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	srv := &http.Server{Addr: mod.addr, Handler: mod}
	l, err := net.Listen("tcp", mod.addr)
	if err != nil {
		mod.LogErrorf("Can't listen on %s", mod.addr)
		return
	}
	defer l.Close()

	go func() {
		err := srv.Serve(l)
		if err != nil {
			mod.LogErrorf("Server error: %v", err)
		}
		// Go 1.8+: srv.Close()
	}()

	select {
	case <-mod.SigChan:
		return
	}
}

func (mod *WebBee) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ev := bees.Event{
		Bee: mod.Name(),
		Options: []bees.Placeholder{
			{
				Name:  "remote_addr",
				Type:  "address",
				Value: req.RemoteAddr,
			},
			{
				Name:  "url",
				Type:  "url",
				Value: req.URL.String(),
			},
		},
	}

	u, err := url.ParseRequestURI(req.RequestURI)
	if err == nil {
		params := u.Query()
		ev.Options.SetValue("query_params", "map", params)
	}

	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err == nil {
		ev.Options.SetValue("data", "string", string(b))
	}

	var payload interface{}
	err = json.Unmarshal(b, &payload)
	if err == nil {
		ev.Options.SetValue("json", "map", payload)
	}

	switch req.Method {
	case "GET":
		ev.Name = "get"
	case "POST":
		ev.Name = "post"
	case "PUT":
		ev.Name = "put"
	case "PATCH":
		ev.Name = "patch"
	case "DELETE":
		ev.Name = "delete"
	}

	mod.eventChan <- ev
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *WebBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.addr)
}
