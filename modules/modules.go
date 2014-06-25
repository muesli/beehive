// beehive's central module system.
package modules

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Interface which all modules need to implement
type ModuleInterface interface {
	// Name of the module
	Name() string
	// Description of the module
	Description() string
	// Events defined by module
	Events() []EventDescriptor
	// Actions supported by module
	Actions() []ActionDescriptor

	// Activates the module
	Run(eventChannel chan Event)
	// Handles an action
	Action(action Action) []Placeholder
}

type Event struct {
	Namespace   string
	Name        string
	Options     []Placeholder
}

type Action struct {
	Namespace   string
	Name        string
	Options     []Placeholder
}

type Placeholder struct {
	Name        string
	Type        string
	Value       interface{}
}

// Descriptors

type EventDescriptor struct {
	Namespace   string
	Name        string
	Description string
	Options     []PlaceholderDescriptor
}

type ActionDescriptor struct {
	Namespace   string
	Name        string
	Description string
	Options     []PlaceholderDescriptor
}

type PlaceholderDescriptor struct {
	Name        string
	Description string
	Type        string
}

// Chain

type ChainElement struct {
	Action  Action
	Mapping map[string]string
}

type Chain struct {
	Name        string
	Description string
	Event       *Event
	Elements    []ChainElement
}

var (
	config = "./beehive.conf"

	EventsIn = make(chan Event)
	modules map[string]*ModuleInterface = make(map[string]*ModuleInterface)
	chains  []Chain
)

func GetActionDescriptor(action *Action) ActionDescriptor {
	mod := (*GetModule(action.Namespace))
	for _, ac := range mod.Actions() {
		if ac.Name == action.Name {
			return ac
		}
	}

	return ActionDescriptor{}
}

func GetEventDescriptor(event *Event) EventDescriptor {
	mod := (*GetModule(event.Namespace))
	for _, ev := range mod.Events() {
		if ev.Name == event.Name {
			return ev
		}
	}

	return EventDescriptor{}
}

func handleEvents() {
	for {
		event := <-EventsIn

		log.Println()
		log.Println("Event received:", event.Namespace, "/", event.Name, "-", GetEventDescriptor(&event).Description)
		for _, v := range event.Options {
			log.Println("\tOptions:", v)
		}

		for _, c := range chains {
			if c.Event.Name != event.Name || c.Event.Namespace != event.Namespace {
				continue
			}

			log.Println("Executing chain:", c.Name, "-", c.Description)
			for _, el := range c.Elements {
				action := el.Action
				for k, v := range el.Mapping {
					for _, ov := range event.Options {
						if ov.Name == k {
							opt := Placeholder{
								Name:  v,
								Type:  ov.Type,
								Value: ov.Value,
							}
							action.Options = append(action.Options, opt)
						}
					}
				}

				log.Println("\tExecuting action:", action.Namespace, "/", action.Name, "-", GetActionDescriptor(&action).Description)
				for _, v := range action.Options {
					log.Println("\t\tOptions:", v)
				}
				(*GetModule(action.Namespace)).Action(action)
			}
		}
	}
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

// Loads chains from config
func LoadChains() {
	j, err := ioutil.ReadFile(config)
	if err == nil {
		err = json.Unmarshal(j, &chains)
		if err != nil {
			panic(err)
		}
	}
}

// Loads chains from config
func SaveChains() {
	j, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(config, j, 0644)
}

// Starts all registered modules
func StartModules() {
	for _, mod := range modules {
		(*mod).Run(EventsIn)
	}

	LoadChains()
//	SaveChains()
}

func init() {
	log.Println("Waking the bees...")

	go handleEvents()
}
