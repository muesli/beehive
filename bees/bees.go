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
package bees

import (
	"log"
	"sync"
)

// Interface which all modules need to implement
type BeeInterface interface {
	// Name of the module
	Name() string
	// Namespace of the module
	Namespace() string
	// Description of the module
	Description() string

	// Activates the module
	Run(eventChannel chan Event)
	// Stop the module
	Stop()

	WaitGroup() *sync.WaitGroup

	// Handles an action
	Action(action Action) []Placeholder
}

// An instance of a module is called a Bee
type Bee struct {
	Name        string
	Class       string
	Description string
	Options     []BeeOption
}

// An Event
type Event struct {
	Bee     string
	Name    string
	Options []Placeholder
}

// An Action
type Action struct {
	Bee     string
	Name    string
	Options []Placeholder
}

// A Filter
type Filter struct {
	Name    string
	Options []FilterOption
}

// A FilterOption used by filters
type FilterOption struct {
	Name            string
	Type            string
	Inverse         bool
	CaseInsensitive bool
	Trimmed         bool
	Value           interface{}
}

// A BeeOption is used to configure bees
type BeeOptions []BeeOption
type BeeOption struct {
	Name  string
	Type  string
	Value interface{}
}

// A Placeholder used by ins & outs of a module.
type Placeholder struct {
	Name  string
	Type  string
	Value interface{}
}

var (
	eventsIn                                     = make(chan Event)
	bees   map[string]*BeeInterface        = make(map[string]*BeeInterface)
	factories map[string]*BeeFactoryInterface = make(map[string]*BeeFactoryInterface)
	chains    []Chain
)

func (opts BeeOptions) GetValue(name string) interface{} {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.Value
		}
	}

	return nil
}

// Handles incoming events and executes matching Chains.
func handleEvents() {
	for {
		event, ok := <-eventsIn
		if !ok {
			log.Println()
			log.Println("Stopped event handler!")
			break
		}

		log.Println()
		log.Println("Event received:", event.Bee, "/", event.Name, "-", GetEventDescriptor(&event).Description)
		for _, v := range event.Options {
			log.Println("\tOptions:", v)
		}

		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Println("Fatal chain event:", e)
				}
			}()

			execChains(&event)
		}()
	}
}

// Bees need to call this method to register themselves
func RegisterBee(mod BeeInterface) {
	log.Println("Worker bee ready:", mod.Name(), "-", mod.Description())

	bees[mod.Name()] = &mod
}

// Returns module with this name
func GetBee(identifier string) *BeeInterface {
	bee, ok := bees[identifier]
	if ok {
		return bee
	}

	return nil
}

func startBee(mod *BeeInterface, fatals int) {
	if fatals >= 3 {
		log.Println("Terminating evil bee", (*mod).Name(), "after", fatals, "failed tries!")
		return
	}

	defer func(mod *BeeInterface) {
		if e := recover(); e != nil {
			log.Println("Fatal bee event:", e, fatals)
			startBee(mod, fatals+1)
		}
	}(mod)

	defer (*mod).WaitGroup().Done()
	(*mod).Run(eventsIn)
}

// Starts all registered bees
func StartBees(beeList []Bee) {
	eventsIn = make(chan Event)
	go handleEvents()

	for _, bee := range beeList {
		factory := GetFactory(bee.Class)
		if factory == nil {
			panic("Unknown bee-class in config file: " + bee.Class)
		}
		mod := (*factory).New(bee.Name, bee.Description, bee.Options)
		RegisterBee(mod)
	}

	for _, m := range bees {
		go func(mod *BeeInterface) {
			startBee(mod, 0)
		}(m)
	}
}

func StopBees() {
	for _, bee := range bees {
		log.Println("Stopping bee:", (*bee).Name())
		(*bee).Stop()
	}

	close(eventsIn)
	bees = make(map[string]*BeeInterface)
}

func RestartBees(bees []Bee) {
	StopBees()
	StartBees(bees)
}

func NewBee(name, factoryName, description string) Module {
	b := Module{
		ModName:        name,
		ModNamespace:   factoryName,
		ModDescription: description,
		SigChan:        make(chan bool),
		waitGroup:      &sync.WaitGroup{},
	}
	b.waitGroup.Add(1)

	return b
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
}
