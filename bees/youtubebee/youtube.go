/*
 *    Copyright (C) 2014      Daniel 'grindhold' Brendle
 *                  2014-2017 Christian Muehlhaeuser
 *
 *	  This program is free software: you can redistribute it and/or modify
 *	  it under the terms of the GNU Affero General Public License as published
 *	  by the Free Software Foundation, either version 3 of the License, or
 *	  (at your option) any later version.
 *
 *	  This program is distributed in the hope that it will be useful,
 *	  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *	  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.	See the
 *	  GNU Affero General Public License for more details.
 *
 *	  You should have received a copy of the GNU Affero General Public License
 *	  along with this program.	If not, see <http://www.gnu.org/licenses/>.
 *
 *	  Authors:
 *		Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package youtubebee is a Bee for tunneling Youtube push notifications.
package youtubebee

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/muesli/beehive/bees"
)

// YoutubeBee is a Bee for handling Youtube push notifications.
type YoutubeBee struct {
	bees.Bee

	url string

	addr string

	eventChan chan bees.Event
}

// Run executes the Bee's event loop.
func (mod *YoutubeBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan
	subscriptionLink := "https://pubsubhubbub.appspot.com/subscribe"
	channelURLTokens := strings.Split(mod.url, "/")
	channelID := channelURLTokens[len(channelURLTokens)-1]
	topic := "https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelID
	hardcodedAddress := "0.0.0.0:5050" // should be mod.addr

	srv := &http.Server{Addr: hardcodedAddress, Handler: mod}
	l, err := net.Listen("tcp", hardcodedAddress)
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
		// send POST to Google's pubsubhubbub to subscribe
		// need to be in form-data format
		data := url.Values{}
		data.Set("hub.mode", "subscribe")
		data.Set("hub.topic", topic)
		data.Set("hub.callback", mod.addr)
		client := &http.Client{}
		r, _ := http.NewRequest("POST", subscriptionLink, strings.NewReader(data.Encode())) // URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

		resp, _ := client.Do(r)
		for resp.Status != "202 Accepted" {
			// redo until success
			resp, _ = client.Do(r)
		}

	}()

	select {
	case <-mod.SigChan:
		return
	}
}

func (mod *YoutubeBee) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		ev := bees.Event{
			Bee: mod.Name(),
		}
		ev.Name = "push"
		type Feed struct {
			XMLName xml.Name `xml:"feed"`
			Text    string   `xml:",chardata"`
			Yt      string   `xml:"yt,attr"`
			Xmlns   string   `xml:"xmlns,attr"`
			Link    []struct {
				Text string `xml:",chardata"`
				Rel  string `xml:"rel,attr"`
				Href string `xml:"href,attr"`
			} `xml:"link"`
			Title   string `xml:"title"`
			Updated string `xml:"updated"`
			Entry   struct {
				Text      string `xml:",chardata"`
				ID        string `xml:"id"`
				VideoId   string `xml:"videoId"`
				ChannelId string `xml:"channelId"`
				Title     string `xml:"title"`
				Link      struct {
					Text string `xml:",chardata"`
					Rel  string `xml:"rel,attr"`
					Href string `xml:"href,attr"`
				} `xml:"link"`
				Author struct {
					Text string `xml:",chardata"`
					Name string `xml:"name"`
					URI  string `xml:"uri"`
				} `xml:"author"`
				Published string `xml:"published"`
				Updated   string `xml:"updated"`
			} `xml:"entry"`
		}
		var feed Feed
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
		}
		xml.Unmarshal([]byte(body), &feed)
		ev.Options.SetValue("channelUrl", "string", feed.Entry.Author.URI)
		ev.Options.SetValue("vidUrl", "string", feed.Entry.Link.Href)

		mod.eventChan <- ev
	} else if req.Method == "GET" {
		challenge := req.URL.Query().Get("hub.challenge")
		if challenge != "" {
			fmt.Fprintf(w, challenge)
		}
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *YoutubeBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.addr)
	options.Bind("channel", &mod.url)
}
