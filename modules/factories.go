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

package modules

import (
	"log"
	"sync"
)

type Module struct {
	ModName        string
	ModNamespace   string
	ModDescription string

	SigChan		chan bool
	waitGroup	*sync.WaitGroup
}

func (mod *Module) Name() string {
	return mod.ModName
}

func (mod *Module) Namespace() string {
	return mod.ModNamespace
}

func (mod *Module) Description() string {
	return mod.ModDescription
}

func (mod *Module) WaitGroup() *sync.WaitGroup {
	return mod.waitGroup
}

func (mod *Module) Run(chan Event) {
}

func (mod *Module) Stop() {
	close(mod.SigChan)
	mod.waitGroup.Wait()
	log.Println(mod.Name(), "stopped gracefully!")
}

type ModuleFactory struct {
}

func (factory *ModuleFactory) Image() string {
	return ""
}

func (factory *ModuleFactory) Options() []BeeOptionDescriptor {
	return []BeeOptionDescriptor{}
}

func (factory *ModuleFactory) Events() []EventDescriptor {
	return []EventDescriptor{}
}

func (factory *ModuleFactory) Actions() []ActionDescriptor {
	return []ActionDescriptor{}
}

type ModuleFactoryInterface interface {
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

	New(name, description string, options BeeOptions) ModuleInterface
}

// ModuleFactories need to call this method to register themselves
func RegisterFactory(factory ModuleFactoryInterface) {
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
func GetFactory(identifier string) *ModuleFactoryInterface {
	factory, ok := factories[identifier]
	if ok {
		return factory
	}

	return nil
}
