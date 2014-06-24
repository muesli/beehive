// beehive's central module system.
package modules

import (
	"log"
)

// Interface which all modules need to implement
type ModuleInterface interface {
	// Name of the module
	Name() string
	// Description of the module
	Description() string
	// Events defined by module
	Events() []Event
	// Actions supported by module
	Actions() []Action

	// Activates the module
	Run(eventChannel chan Event, actionChannel chan Action)
	// Handles an action
	Action(action Action) []Placeholder
}

type Event struct {
	Name        string
	Description string
	Options     []Placeholder
}

type Action struct {
	Name        string
	Description string
	Options     []Placeholder
}

type Placeholder struct {
	Name        string
	Description string
	Type        string
	Value       interface{}
}

type ChainElement struct {
	Action Action
	PlaceholderMap map[string]string
}

type Chain struct {
	Name        string
	Description string
	Event       *Event
	Elements []ChainElement
}

var (
	EventsIn   = make(chan Event)
	ActionsOut = make(chan Action)

	modules map[string]*ModuleInterface = make(map[string]*ModuleInterface)
	chains []Chain
)

func init() {
	log.Println("Waking the bees...")

	go func() {
		for {
			event := <-EventsIn
			log.Println("Event received:", event.Name)
			for _, v := range event.Options {
				log.Println("\tOptions:", v)
			}

			for _, c := range chains {
				if c.Event.Name != event.Name {
					continue
				}

				log.Println("Executing chain:", c.Name, "-", c.Description)
				for _, el := range c.Elements {
					action := el.Action
					for k, v := range el.PlaceholderMap {
						for _, ov := range event.Options {
							if ov.Name == k {
								opt := Placeholder{
									Name: v,
									Type: ov.Type,
									Value: ov.Value,
								}
								action.Options = append(action.Options, opt)
							} 
						}
					}

					(*GetModule("ircbee")).Action(action)
				}
			}
		}
	}()

	go func() {
		for {
			action := <-ActionsOut
			log.Println("Action:", action.Name)
			for _, v := range action.Options {
				log.Println("\tOptions:", v)
			}
			for _, mod := range modules {
				(*mod).Action(action)
			}
		}
	}()
}

// Modules need to call this method to register themselves
func RegisterModule(mod ModuleInterface) {
	log.Println("Worker bee ready:", mod.Name(), "-", mod.Description())
	for _, ev := range mod.Events() {
		log.Println("\tProvides event:", ev.Name, "-", ev.Description)
		for _, opt := range ev.Options {
			log.Println("\t\tPlaceholder:", opt.Name, "-", opt.Description)
		}
	}
	for _, ac := range mod.Actions() {
		log.Println("\tOffers action:", ac.Name, "-", ac.Description)
		for _, opt := range ac.Options {
			log.Println("\t\tPlaceholder:", opt.Name, "-", opt.Description)
		}
	}

	log.Println()
	modules[mod.Name()] = &mod
}

// Returns module with this name
func GetModule(identifier string) *ModuleInterface {
	mod, ok := modules[identifier]
	if ok {
		return mod
	}

	return nil
}

// Starts all registered modules
func StartModules() {
	for _, mod := range modules {
		(*mod).Run(EventsIn, ActionsOut)
	}

	// Create a fake sample chain
	event := Event{}
	for _, ev := range (*GetModule("ircbee")).Events() {
		if ev.Name == "message" {
			event = ev
		}
	}
	action := Action{}
	actionTest := Action{}
	for _, ac := range (*GetModule("ircbee")).Actions() {
		if ac.Name == "send" {
			action = ac
			action.Options = []Placeholder{}
			actionTest = ac
			actionTest.Options = []Placeholder{
				Placeholder{
					Name: "channel",
					Type: "string",
					Value: "muesli",
				},
			}
		}
	}

	// does the in/out placeholder mapping
	ma := make(map[string]string)
	ma["text"] = "text"
	ma["channel"] = "channel"
	mb := make(map[string]string)
	mb["text"] = "text"

	chains = []Chain{
		Chain{
			Name: "Parrot",
			Description: "Echoes everything you say on IRC",
			Event: &event,
			Elements: []ChainElement{
				ChainElement{
					Action: action,
					PlaceholderMap: ma,
				},
				ChainElement{
					Action: actionTest,
					PlaceholderMap: mb,
				},
			},
		},
	}
}
