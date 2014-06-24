// hellobee dummy module for beehive
package hellobee

import (
	"github.com/muesli/beehive/app"
	"github.com/muesli/beehive/modules"
)

type HelloBee struct {
	some_flag string
}

func (sys *HelloBee) Name() string {
	return "hellobee"
}

func (sys *HelloBee) Events() []modules.Event {
	events := []modules.Event{}
	return events
}

func (sys *HelloBee) Actions() []modules.Action {
	actions := []modules.Action{}
	return actions
}

func (sys *HelloBee) Run(MyChannel chan modules.Event) {
	hello_event := modules.Event{
		Name:    "Say Hello",
		Options: []modules.Placeholder{},
	}

	MyChannel <- hello_event
}

func init() {
	hello := HelloBee{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&hello.some_flag, "foo", "", "some text"},
	})

	modules.RegisterModule(&hello)
}
