// hellobee dummy module for beehive
package hellobee

import (
	"github.com/muesli/beehive/app"
	"github.com/muesli/beehive/modules"
)

type HelloBee struct {
	some_flag string
}

func (mod *HelloBee) Name() string {
	return "hellobee"
}

func (mod *HelloBee) Description() string {
	return "A 'Hello World' module for beehive"
}

func (mod *HelloBee) Events() []modules.Event {
	events := []modules.Event{}
	return events
}

func (mod *HelloBee) Actions() []modules.Action {
	actions := []modules.Action{}
	return actions
}

func (mod *HelloBee) Run(MyChannel chan modules.Event) {
	hello_event := modules.Event{
		Namespace: mod.Name(),
		Name:    "Say Hello",
		Options: []modules.Placeholder{},
	}

	MyChannel <- hello_event
}

func (mod *HelloBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}

func init() {
	hello := HelloBee{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&hello.some_flag, "foo", "default value", "some option"},
	})

	modules.RegisterModule(&hello)
}
