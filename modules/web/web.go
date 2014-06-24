// beehive's web-module.
package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/hoisie/web"
	"github.com/muesli/beehive/modules"
)

var(
	cIn chan modules.Event
	cOut chan modules.Action
)

type WebBee struct {
	Addr string
}

func (sys *WebBee) Name() string {
	return "webbee"
}

func (sys *WebBee) Run(channelIn chan modules.Event, channelOut chan modules.Action) {
	cIn = channelIn
	cOut = channelOut
	go web.Run(sys.Addr)
}

func (sys *WebBee) Events() []modules.Event {
	events := []modules.Event{}
	return events
}

func (sys *WebBee) Actions() []modules.Action {
	actions := []modules.Action{}
	return actions
}

func (sys *WebBee) Action(action modules.Action) bool {
	return false
}

func ActionRequest(ctx *web.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Params:", string(b))

	var payload interface{}
	err = json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("JSON'd:", payload)

	data := payload.(map[string]interface{})
	name := data["name"].(string)
	channel := data["channel"].(string)
	text := data["text"].(string)

		action := modules.Action{
			Name: name,
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name: "channel",
					Type: "string",
					Value: channel,
				},
				modules.Placeholder{
					Name: "text",
					Type: "string",
					Value: text,
				},
			},
		}

	cOut <- action
}

func init() {
	w := WebBee{
		Addr: "0.0.0.0:12345",
	}
	web.Post("/action", ActionRequest)

	modules.RegisterModule(&w)
}
