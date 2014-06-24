// beehive's web-module.
package web

import (
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/muesli/beehive/modules"
	"io/ioutil"
	"log"
)

var (
	cIn  chan modules.Event
	cOut chan modules.Action
)

type WebBee struct {
	Addr string
}

func (sys *WebBee) Name() string {
	return "webbee"
}

func (sys *WebBee) Description() string {
	return "A RESTful HTTP module for beehive"
}

func (sys *WebBee) Run(channelIn chan modules.Event, channelOut chan modules.Action) {
	cIn = channelIn
	cOut = channelOut
	go web.Run(sys.Addr)
}

func (sys *WebBee) Events() []modules.Event {
	events := []modules.Event{
		modules.Event{
			Name:        "post",
			Description: "A POST call was received by the HTTP server",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:        "json",
					Description: "JSON map received from caller",
					Type:        "json",
				},
				modules.Placeholder{
					Name:        "ip",
					Description: "IP of the caller",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (sys *WebBee) Actions() []modules.Action {
	actions := []modules.Action{}
	return actions
}

func (sys *WebBee) Action(action modules.Action) bool {
	return false
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
		Name: "post",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "json",
				Type:  "json",
				Value: payload,
			},
			modules.Placeholder{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}
	cIn <- ev
}

func init() {
	w := WebBee{
		Addr: "0.0.0.0:12345",
	}
	web.Post("/event", PostRequest)

	modules.RegisterModule(&w)
}
