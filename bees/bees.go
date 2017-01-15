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

// Interface which all bees need to implement
type BeeInterface interface {
	// Name of the bee
	Name() string
	// Namespace of the bee
	Namespace() string
	// Description of the bee
	Description() string

	// Activates the bee
	Run(eventChannel chan Event)
	// Stop the bee
	Stop()

	WaitGroup() *sync.WaitGroup

	// Handles an action
	Action(action Action) []Placeholder
}

// Base-struct to be embedded by bee implementations
type Bee struct {
	BeeName        string
	BeeNamespace   string
	BeeDescription string

	SigChan   chan bool
	waitGroup *sync.WaitGroup
}

// An instance of a bee
type BeeInstance struct {
	Name        string
	Class       string
	Description string
	Options     []BeeOption
}

// An Event
type Event struct {
	Bee     string
	Name    string
	Options PlaceholderSlice
}

// An Action
type Action struct {
	Bee     string
	Name    string
	Options PlaceholderSlice
}

// A Filter
type Filter struct {
	Name    string
	Options []FilterOption
}

var (
	eventsIn                                  = make(chan Event)
	bees      map[string]*BeeInterface        = make(map[string]*BeeInterface)
	factories map[string]*BeeFactoryInterface = make(map[string]*BeeFactoryInterface)
	chains    []Chain
)

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
func RegisterBee(bee BeeInterface) {
	log.Println("Worker bee ready:", bee.Name(), "-", bee.Description())

	bees[bee.Name()] = &bee
}

// Returns bee with this name
func GetBee(identifier string) *BeeInterface {
	bee, ok := bees[identifier]
	if ok {
		return bee
	}

	return nil
}

// Returns all known bees
func GetBees() []*BeeInterface {
	r := []*BeeInterface{}
	for _, bee := range bees {
		r = append(r, bee)
	}

	return r
}

// Returns all known bee factories
func GetBeeFactories() []*BeeFactoryInterface {
	r := []*BeeFactoryInterface{}
	for _, factory := range factories {
		r = append(r, factory)
	}

	return r
}

// Starts a bee and recovers from panics
func startBee(bee *BeeInterface, fatals int) {
	if fatals >= 3 {
		log.Println("Terminating evil bee", (*bee).Name(), "after", fatals, "failed tries!")
		return
	}

	defer func(bee *BeeInterface) {
		if e := recover(); e != nil {
			log.Println("Fatal bee event:", e, fatals)
			startBee(bee, fatals+1)
		}
	}(bee)

	defer (*bee).WaitGroup().Done()
	(*bee).Run(eventsIn)
}

// Starts all registered bees
func StartBees(beeList []BeeInstance) {
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

// Stops all bees gracefully
func StopBees() {
	for _, bee := range bees {
		log.Println("Stopping bee:", (*bee).Name())
		(*bee).Stop()
	}

	close(eventsIn)
	bees = make(map[string]*BeeInterface)
}

// Stops all running bees and restarts a new set of bees
func RestartBees(bees []BeeInstance) {
	StopBees()
	StartBees(bees)
}

// Returns a new bee and sets up sig-channel & waitGroup
func NewBee(name, factoryName, description string) Bee {
	b := Bee{
		BeeName:        name,
		BeeNamespace:   factoryName,
		BeeDescription: description,
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
