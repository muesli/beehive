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

package bees

import (
	"log"
	"sync"
)

type BeeFactory struct {
}

func (factory *BeeFactory) Image() string {
	return ""
}

func (factory *BeeFactory) Options() []BeeOptionDescriptor {
	return []BeeOptionDescriptor{}
}

func (factory *BeeFactory) Events() []EventDescriptor {
	return []EventDescriptor{}
}

func (factory *BeeFactory) Actions() []ActionDescriptor {
	return []ActionDescriptor{}
}

type BeeFactoryInterface interface {
	// Name of the module
	Name() string
	// Description of the module
	Description() string
	// An image url for the module
	Image() string

	// Options supported by module
	Options() []BeeOptionDescriptor
	// Events defined by module
	Events() []EventDescriptor
	// Actions supported by module
	Actions() []ActionDescriptor

	New(name, description string, options BeeOptions) BeeInterface
}

// ModuleFactories need to call this method to register themselves
func RegisterFactory(factory BeeFactoryInterface) {
	log.Println("Bee Factory ready:", factory.Name(), "-", factory.Description())
	for _, ev := range factory.Events() {
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
	log.Println()

	factories[factory.Name()] = &factory
}

// Returns factory with this name
func GetFactory(identifier string) *BeeFactoryInterface {
	factory, ok := factories[identifier]
	if ok {
		return factory
	}

	return nil
}

func (bee *Bee) Name() string {
	return bee.BeeName
}

func (bee *Bee) Namespace() string {
	return bee.BeeNamespace
}

func (bee *Bee) Description() string {
	return bee.BeeDescription
}

func (bee *Bee) WaitGroup() *sync.WaitGroup {
	return bee.waitGroup
}

func (bee *Bee) Run(chan Event) {
}

func (bee *Bee) IsRunning() bool {
	return bee.Running
}

func (bee *Bee) Start() {
	bee.Running = true
}

func (bee *Bee) Stop() {
	close(bee.SigChan)
	bee.waitGroup.Wait()
	bee.Running = false
	log.Println(bee.Name(), "stopped gracefully!")
}

func (bee *Bee) LastEvent() time.Time {
	return bee.lastEvent
}

func (bee *Bee) LastAction() time.Time {
	return bee.lastAction
}

func (bee *Bee) LogEvent() {
	bee.lastEvent = time.Now()
}

func (bee *Bee) LogAction() {
	bee.lastAction = time.Now()
}
