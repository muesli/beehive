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

// A Filter
type Filter struct {
	Name      string
	Options   []Placeholder
}

// A Placeholder used by ins & outs of a module.
type Placeholder struct {
	Name  string
	Type  string
	Value interface{}
}

var (
	eventsIn                             = make(chan Event)
	modules  map[string]*ModuleInterface = make(map[string]*ModuleInterface)
	chains   []Chain
)

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

// Starts all registered modules
func StartModules() {
	for _, mod := range modules {
		(*mod).Run(eventsIn)
	}
}

// Getter for chains
func Chains() []Chain {
	return chains
}

// Setter for chains
func SetChains(cs []Chain) {
	chains = cs
}

func init() {
	log.Println("Waking the bees...")

	go handleEvents()
}
