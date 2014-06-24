// hellobee dummy module for beehive
package hellobee

type HelloBee struct {
}

func (sys *IrcBee) Name() string {
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
