/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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

// Package bees is Beehive's central module system
package bees

import (
	"log"
	"sync"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

// BeeInterface is an interface all bees implement.
type BeeInterface interface {
	// Name of the bee
	Name() string
	// Namespace of the bee
	Namespace() string

	// Description of the bee
	Description() string
	// SetDescription sets a description
	SetDescription(s string)

	// Config returns this bees config
	Config() BeeConfig
	// Options of the bee
	Options() BeeOptions
	// SetOptions to configure the bee
	SetOptions(options BeeOptions)

	// ReloadOptions gets called after a bee's options get updated
	ReloadOptions(options BeeOptions)

	// Activates the bee
	Run(eventChannel chan Event)
	// Running returns the current state of the bee
	IsRunning() bool
	// Start the bee
	Start()
	// Stop the bee
	Stop()

	LastEvent() time.Time
	LogEvent()
	LastAction() time.Time
	LogAction()

	SetSigChan(c chan bool)
	WaitGroup() *sync.WaitGroup

	// Handles an action
	Action(action Action) []Placeholder
}

// BeeConfig contains all settings for a single Bee.
type BeeConfig struct {
	Name        string
	Class       string
	Description string
	Options     BeeOptions
}

// Bee is the base-struct to be embedded by bee implementations.
type Bee struct {
	config BeeConfig

	lastEvent  time.Time
	lastAction time.Time

	Running   bool
	SigChan   chan bool
	waitGroup *sync.WaitGroup
}

var (
	bees      map[string]*BeeInterface        = make(map[string]*BeeInterface)
	factories map[string]*BeeFactoryInterface = make(map[string]*BeeFactoryInterface)
)

// RegisterBee gets called by Bees to register themselves.
func RegisterBee(bee BeeInterface) {
	log.Println("Worker bee ready:", bee.Name(), "-", bee.Description())

	bees[bee.Name()] = &bee
}

// GetBee returns a bee with a specific name.
func GetBee(identifier string) *BeeInterface {
	bee, ok := bees[identifier]
	if ok {
		return bee
	}

	return nil
}

// GetBees returns all known bees.
func GetBees() []*BeeInterface {
	r := []*BeeInterface{}
	for _, bee := range bees {
		r = append(r, bee)
	}

	return r
}

// startBee starts a bee and recovers from panics.
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

// NewBeeInstance sets up a new Bee with supplied config.
func NewBeeInstance(bee BeeConfig) *BeeInterface {
	factory := GetFactory(bee.Class)
	if factory == nil {
		panic("Unknown bee-class in config file: " + bee.Class)
	}
	mod := (*factory).New(bee.Name, bee.Description, bee.Options)
	RegisterBee(mod)

	return &mod
}

// DeleteBee removes a Bee instance.
func DeleteBee(bee *BeeInterface) {
	(*bee).Stop()

	delete(bees, (*bee).Name())
}

// StartBee starts a bee.
func StartBee(bee BeeConfig) *BeeInterface {
	b := NewBeeInstance(bee)

	(*b).Start()
	go func(mod *BeeInterface) {
		startBee(mod, 0)
	}(b)

	return b
}

// StartBees starts all registered bees.
func StartBees(beeList []BeeConfig) {
	eventsIn = make(chan Event)
	go handleEvents()

	for _, bee := range beeList {
		StartBee(bee)
	}
}

// StopBees stops all bees gracefully.
func StopBees() {
	for _, bee := range bees {
		log.Println("Stopping bee:", (*bee).Name())
		(*bee).Stop()
	}

	close(eventsIn)
	bees = make(map[string]*BeeInterface)
}

// RestartBee restarts a Bee.
func RestartBee(bee *BeeInterface) {
	(*bee).Stop()

	(*bee).SetSigChan(make(chan bool))
	(*bee).WaitGroup().Add(1)
	(*bee).Start()
	go func(mod *BeeInterface) {
		startBee(mod, 0)
	}(bee)
}

// RestartBees stops all running bees and restarts a new set of bees.
func RestartBees(bees []BeeConfig) {
	StopBees()
	StartBees(bees)
}

// NewBee returns a new bee and sets up sig-channel & waitGroup.
func NewBee(name, factoryName, description string, options []BeeOption) Bee {
	c := BeeConfig{
		Name:        name,
		Class:       factoryName,
		Description: description,
		Options:     options,
	}
	b := Bee{
		config:    c,
		SigChan:   make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	b.waitGroup.Add(1)

	return b
}

// BeeConfigs returns configs for all Bees.
func BeeConfigs() []BeeConfig {
	bs := []BeeConfig{}
	for _, b := range bees {
		bs = append(bs, (*b).Config())
	}

	return bs
}

// UUID generates a new unique ID.
func UUID() string {
	u, _ := uuid.NewV4()
	return u.String()
}

func init() {
	log.Println("Waking the bees...")
}
