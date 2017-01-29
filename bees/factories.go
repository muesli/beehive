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
)

// A BeeFactory is the base struct to be embedded by other BeeFactories.
type BeeFactory struct {
}

// Image returns an empty image filename per default.
func (factory *BeeFactory) Image() string {
	return ""
}

// LogoColor returns the default logo color.
func (factory *BeeFactory) LogoColor() string {
	return "#35465c"
}

// Options returns the default empty options set.
func (factory *BeeFactory) Options() []BeeOptionDescriptor {
	return []BeeOptionDescriptor{}
}

// Events returns the default empty events set.
func (factory *BeeFactory) Events() []EventDescriptor {
	return []EventDescriptor{}
}

// Actions returns the default empty actions set.
func (factory *BeeFactory) Actions() []ActionDescriptor {
	return []ActionDescriptor{}
}

// A BeeFactoryInterface is the interface that gets implemented by a BeeFactory.
type BeeFactoryInterface interface {
	// Name of the module
	Name() string
	// Description of the module
	Description() string
	// An image url for the module
	Image() string
	// A logo color for the module
	LogoColor() string

	// Options supported by module
	Options() []BeeOptionDescriptor
	// Events defined by module
	Events() []EventDescriptor
	// Actions supported by module
	Actions() []ActionDescriptor

	New(name, description string, options BeeOptions) BeeInterface
}

// RegisterFactory gets called by BeeFactories to register themselves.
func RegisterFactory(factory BeeFactoryInterface) {
	log.Println("Bee Factory ready:", factory.Name(), "-", factory.Description())
	/* for _, ev := range factory.Events() {
		log.Println("\tProvides event:", ev.Name, "-", ev.Description)
		for _, opt := range ev.Options {
			log.Println("\t\tPlaceholder:", opt.Name, "-", opt.Description)
		}
	}
	for _, ac := range factory.Actions() {
		log.Println("\tOffers action:", ac.Name, "-", ac.Description)
		for _, opt := range ac.Options {
			log.Println("\t\tPlaceholder:", opt.Name, "-", opt.Description)
		}
	}
	log.Println() */

	factories[factory.Name()] = &factory
}

// GetFactory returns the factory with a specific name.
func GetFactory(identifier string) *BeeFactoryInterface {
	factory, ok := factories[identifier]
	if ok {
		return factory
	}

	return nil
}

// GetFactories returns all known bee factories.
func GetFactories() []*BeeFactoryInterface {
	r := []*BeeFactoryInterface{}
	for _, factory := range factories {
		r = append(r, factory)
	}

	return r
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
