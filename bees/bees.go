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

// Package bees is Beehive's central module system.
package bees

import (
	"fmt"
	"sync"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
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

	Logln(args ...interface{})
	Logf(format string, args ...interface{})
	LogErrorf(format string, args ...interface{})
	LogFatal(args ...interface{})

	SetSigChan(c chan bool)
	WaitGroup() *sync.WaitGroup

	// Handles an action
	Action(action Action) []Placeholder
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
	bees      = make(map[string]*BeeInterface)
	factories = make(map[string]*BeeFactoryInterface)
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
		(*bee).Stop()
		return
	}

	(*bee).WaitGroup().Add(1)
	defer (*bee).WaitGroup().Done()

	defer func(bee *BeeInterface) {
		if e := recover(); e != nil {
			log.Println("Fatal bee event:", e, fatals)
			go startBee(bee, fatals+1)
		}
	}(bee)

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

	return b
}

// Name returns the configured name for a bee.
func (bee *Bee) Name() string {
	return bee.config.Name
}

// Namespace returns the namespace for a bee.
func (bee *Bee) Namespace() string {
	return bee.config.Class
}

// Description returns the description for a bee.
func (bee *Bee) Description() string {
	return bee.config.Description
}

// SetDescription sets the description for a bee.
func (bee *Bee) SetDescription(s string) {
	bee.config.Description = s
}

// Config returns the config for a bee.
func (bee *Bee) Config() BeeConfig {
	return bee.config
}

// Options returns the options for a bee.
func (bee *Bee) Options() BeeOptions {
	return bee.config.Options
}

// SetOptions sets the options for a bee.
func (bee *Bee) SetOptions(options BeeOptions) {
	bee.config.Options = options
}

// SetOption sets one option for a bee.
func (bee *Bee) SetOption(name string, value string) bool {
	for i := 0 ; i < len(bee.config.Options); i++ {
		if bee.config.Options[i].Name == name {
			bee.config.Options[i].Value = value

			return true
		}
	}

	return false
}

// SetSigChan sets the signaling channel for a bee.
func (bee *Bee) SetSigChan(c chan bool) {
	bee.SigChan = c
}

// WaitGroup returns the WaitGroup for a bee.
func (bee *Bee) WaitGroup() *sync.WaitGroup {
	return bee.waitGroup
}

// Run is the default, empty implementation of a Bee's Run method.
func (bee *Bee) Run(chan Event) {
	select {
	case <-bee.SigChan:
		return
	}
}

// Action is the default, empty implementation of a Bee's Action method.
func (bee *Bee) Action(action Action) []Placeholder {
	return []Placeholder{}
}

// IsRunning returns whether a Bee is currently running.
func (bee *Bee) IsRunning() bool {
	return bee.Running
}

// Start gets called when a Bee gets started.
func (bee *Bee) Start() {
	bee.Running = true
}

// Stop gracefully stops a Bee.
func (bee *Bee) Stop() {
	if !bee.IsRunning() {
		return
	}
	log.Println(bee.Name(), "stopping gracefully!")

	close(bee.SigChan)
	bee.waitGroup.Wait()
	bee.Running = false
	log.Println(bee.Name(), "stopped gracefully!")
}

// LastEvent returns the timestamp of the last triggered event.
func (bee *Bee) LastEvent() time.Time {
	return bee.lastEvent
}

// LastAction returns the timestamp of the last triggered action.
func (bee *Bee) LastAction() time.Time {
	return bee.lastAction
}

// LogEvent logs the last triggered event.
func (bee *Bee) LogEvent() {
	bee.lastEvent = time.Now()
}

// LogAction logs the last triggered action.
func (bee *Bee) LogAction() {
	bee.lastAction = time.Now()
}

// Logln logs args
func (bee *Bee) Logln(args ...interface{}) {
	a := []interface{}{"[" + bee.Name() + "]:"}
	for _, v := range args {
		a = append(a, v)
	}

	log.Println(a...)
	Log(bee.Name(), fmt.Sprintln(args...), LogInfo)
}

// Logf logs a formatted string
func (bee *Bee) Logf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Printf("[%s]: %s", bee.Name(), s)
	Log(bee.Name(), s, LogInfo)
}

// LogErrorf logs a formatted error string
func (bee *Bee) LogErrorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Errorf("[%s]: %s", bee.Name(), s)
	Log(bee.Name(), s, LogError)
}

// LogDebugf logs a formatted debug string
func (bee *Bee) LogDebugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Debugf("[%s]: %s", bee.Name(), s)
	Log(bee.Name(), s, LogDebug)
}

// LogFatal logs a fatal error
func (bee *Bee) LogFatal(args ...interface{}) {
	a := []interface{}{"[" + bee.Name() + "]:"}
	for _, v := range args {
		a = append(a, v)
	}
	log.Panicln(a...)
	Log(bee.Name(), fmt.Sprintln(args...), LogFatal)
}

// UUID generates a new unique ID.
func UUID() string {
	u, _ := uuid.NewV4()
	return u.String()
}
