/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

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
	Events() []EventDescriptor
	// Actions supported by module
	Actions() []ActionDescriptor

	// Activates the module
	Run(eventChannel chan Event)
	// Handles an action
	Action(action Action) []Placeholder
}

// An Event
type Event struct {
	Namespace string
	Name      string
	Options   []Placeholder
}

// An Action
type Action struct {
	Namespace string
	Name      string
	Options   []Placeholder
}

// A Placeholder used by ins & outs of a module.
type Placeholder struct {
	Name  string
	Type  string
	Value interface{}
}

// Descriptors

// Modules provide events, which are described in an EventDescriptor.
type EventDescriptor struct {
	Namespace   string
	Name        string
	Description string
	Options     []PlaceholderDescriptor
}

// Modules offer actions, which are described in an ActionDescriptor.
type ActionDescriptor struct {
	Namespace   string
	Name        string
	Description string
	Options     []PlaceholderDescriptor
}

// A PlaceholderDescriptor shows which in & out values a module expects and returns.
type PlaceholderDescriptor struct {
	Name        string
	Description string
	Type        string
}

// An element in a Chain
type ChainElement struct {
	Action  Action
	Mapping map[string]string
}

// A user defined Chain
type Chain struct {
	Name        string
	Description string
	Event       *Event
	Elements    []ChainElement
}

var (
	eventsIn                             = make(chan Event)
	modules  map[string]*ModuleInterface = make(map[string]*ModuleInterface)
	chains   []Chain
)

// Returns the ActionDescriptor matching an action.
func GetActionDescriptor(action *Action) ActionDescriptor {
	mod := (*GetModule(action.Namespace))
	for _, ac := range mod.Actions() {
		if ac.Name == action.Name {
			return ac
		}
	}

	return ActionDescriptor{}
}

// Returns the EventDescriptor matching an event.
func GetEventDescriptor(event *Event) EventDescriptor {
	mod := (*GetModule(event.Namespace))
	for _, ev := range mod.Events() {
		if ev.Name == event.Name {
			return ev
		}
	}

	return EventDescriptor{}
}

// Execute chains for an event we received.
func execChains(event *Event) {
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

// Handles incoming events and executes matching Chains.
func handleEvents() {
	for {
		event := <-eventsIn

		log.Println()
		log.Println("Event received:", event.Namespace, "/", event.Name, "-", GetEventDescriptor(&event).Description)
		for _, v := range event.Options {
			log.Println("\tOptions:", v)
		}

		execChains(&event)
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

// Getter for chains
func Chains() []Chain {
	return chains
}

// Setter for chains
func SetChains(cs []Chain) {
	chains = cs
}

// Starts all registered modules
func StartModules() {
	for _, mod := range modules {
		(*mod).Run(eventsIn)
	}
}

func init() {
	log.Println("Waking the bees...")

	go handleEvents()
}
