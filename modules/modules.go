// beehive's central module system. 
package modules

import (
	"fmt"
	_ "strings"
)

// Interface which all modules need to implement
type ModuleInterface interface {
	// Name of the module
	Name() 		string
	Events()	[]Event
	Actions()	[]Action
//	Outs()		[]Placeholder

	Handle(event Event) bool
	Run()
}

type Event struct {
	Name	string
//	Options	[]Placeholder
}

type Action struct {
	Name	string
//	Options	[]Placeholder
}

type Message struct {
	To     []string
	Msg    string
	Source string
	Authed bool
}

var (
	EventsIn  = make(chan Event)
	ActionsOut = make(chan Action)

	modules map[string]*ModuleInterface = make(map[string]*ModuleInterface)
)

func init() {
	fmt.Println("Waking the bees...")

	go func() {
		for {
			action := <-ActionsOut
			fmt.Println("Action:", action.Name)
		}
	}()
}

// Sub-systems need to call this method to register themselves
func RegisterModule(mod ModuleInterface) {
	fmt.Println("Registering bee:", mod.Name())

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
		(*mod).Run()
	}
}
