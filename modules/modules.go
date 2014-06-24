// beehive's central module system.
package modules

import (
	"log"
)

// Interface which all modules need to implement
type ModuleInterface interface {
	// Name of the module
	Name() string
	// Events defined by module
	Events() []Event
	// Actions supported by module
	Actions() []Action

	// Activates the module
	Run(eventChannel chan Event, actionChannel chan Action)
	// Handles an action
	Action(action Action) bool
}

type Event struct {
	Name    string
	Options []Placeholder
}

type Action struct {
	Name    string
	Options []Placeholder
}

type Placeholder struct {
	Name  string
	Type  string
	Value string
}

var (
	EventsIn   = make(chan Event)
	ActionsOut = make(chan Action)

	modules map[string]*ModuleInterface = make(map[string]*ModuleInterface)
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

// Sub-systems need to call this method to register themselves
func RegisterModule(mod ModuleInterface) {
	log.Println("Registering bee:", mod.Name())

	modules[mod.Name()] = &mod
}

// Returns sub-system with this name
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
}
