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
 *		CalmBit <calmbit@posteto.net>
 */

// Package httpbee is a Bee that lets you trigger HTTP requests.
package httpbee

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/muesli/beehive/bees"
)

// HTTPBee is a Bee that lets you trigger HTTP requests.
type HTTPBee struct {
	bees.Bee

	addr string
	path string

	eventChan chan bees.Event
}

// Run executes the Bee's event loop.
func (mod *HTTPBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	select {
	case <-mod.SigChan:
		return
	}
}

// Action triggers the action passed to it.
func (mod *HTTPBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	u := ""
	h := []string{}
	action.Options.Bind("url", &u)
	action.Options.Bind("headers", &h)

	switch action.Name {
	case "get":
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}

		mod.parseHeaders(h, req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}

		// reforming the request headers slice here for absolute
		// conformity with what was taken in by the request

		reqHeaders, respHeaders := mod.serializeHeaders(req, resp)

		ev, err := mod.prepareResponseEvent(b)
		if err == nil {
			ev.Name = "get"
			ev.Options.SetValue("url", "url", u)
			ev.Options.SetValue("reqHeaders", "[]string", reqHeaders)
			ev.Options.SetValue("respHeaders", "[]string", respHeaders)
			mod.eventChan <- ev
		}

	case "post":
		var err error
		var b []byte
		var j string
		var form url.Values
		var resp *http.Response
		var h []string
		var ctype string

		action.Options.Bind("json", &j)
		action.Options.Bind("form", &form)
		action.Options.Bind("headers", &h)

		var buf *strings.Reader

		if j != "" {
			buf = strings.NewReader(j)
			ctype = "application/json"
		} else {
			buf = strings.NewReader(form.Encode())
			ctype = "application/x-www-form-urlencoded"
		}

		req, err := http.NewRequest("POST", u, buf)
		req.Header.Set("Content-Type", ctype)

		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}

		mod.parseHeaders(h, req)

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}
		defer resp.Body.Close()

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			mod.LogErrorf("Error: %s", err)
			return outs
		}

		ev, err := mod.prepareResponseEvent(b)

		// reforming the request headers slice here for absolute
		// conformity with what was taken in by the request

		reqHeaders, respHeaders := mod.serializeHeaders(req, resp)

		if err == nil {
			ev.Name = "post"
			ev.Options.SetValue("url", "url", u)
			ev.Options.SetValue("reqHeaders", "[]string", reqHeaders)
			ev.Options.SetValue("respHeaders", "[]string", respHeaders)
			mod.eventChan <- ev
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *HTTPBee) prepareResponseEvent(resp []byte) (bees.Event, error) {
	ev := bees.Event{
		Bee: mod.Name(),
		Options: []bees.Placeholder{
			{
				Name:  "data",
				Type:  "string",
				Value: string(resp),
			},
		},
	}

	var payload interface{}
	err := json.Unmarshal(resp, &payload)
	if err == nil {
		ev.Options.SetValue("json", "map", payload)

		// this is a bit of a funny hack:
		// each parameter in the JSON response will get mapped to an (undocumented) event parameter
		// TODO: decide if this a good idea or not (probably not)
		j := make(map[string]interface{})
		err = json.Unmarshal(resp, &j)

		if err == nil {
			for k, v := range j {
				mod.LogDebugf("JSON param: %s = %+v\n", k, v)
				if k == "json" || k == "data" {
					continue
				}

				// FIXME: hard-coded 'string'
				ev.Options.SetValue(k, "string", v)
			}
		}
	}

	return ev, nil
}

func (mod *HTTPBee) parseHeaders(headers []string, req *http.Request) {
	for _, header := range headers {
		comp := strings.SplitN(header, ":", 2)
		if len(comp) != 2 {
			mod.LogErrorf("Warning: header '%s' was not formatted correctly - ignoring", header)
			continue
		}
		req.Header.Set(comp[0], strings.TrimSpace(comp[1]))
	}
}

func (mod *HTTPBee) serializeHeaders(req *http.Request, resp *http.Response) ([]string, []string) {
	reqHeaders := make([]string, len(req.Header))

	for name, header := range req.Header {
		reqHeaders = append(reqHeaders, name+": "+strings.Join(header, ", "))
	}

	respHeaders := make([]string, len(resp.Header))

	for name, header := range resp.Header {
		respHeaders = append(respHeaders, name+": "+strings.Join(header, ", "))
	}

	return reqHeaders, respHeaders

}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *HTTPBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}
