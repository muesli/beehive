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

import log "github.com/Sirupsen/logrus"

// ChainElement is an element in a Chain
type ChainElement struct {
	Action Action
	Filter Filter
}

// Chain is a user defined chain
type Chain struct {
	Name        string
	Description string
	Event       *Event
	Filters     []string
	Actions     []string
	Elements    []ChainElement `json:"Elements,omitempty"`
}

var (
	chains []Chain
)

// GetChains returns all chains
func GetChains() []Chain {
	return chains
}

// GetChain returns a chain with a specific id
func GetChain(id string) *Chain {
	for _, c := range chains {
		if c.Name == id {
			return &c
		}
	}

	return nil
}

// SetChains sets the currently configured chains
func SetChains(cs []Chain) {
	newcs := []Chain{}
	// migrate old chain style
	for _, c := range cs {
		for _, el := range c.Elements {
			if el.Action.Name != "" {
				el.Action.ID = UUID()
				c.Actions = append(c.Actions, el.Action.ID)
				actions = append(actions, el.Action)
			}
			if el.Filter.Name != "" {
				//FIXME: migrate old style filters
				c.Filters = append(c.Filters, el.Filter.Options.Value.(string))
			}
		}
		c.Elements = []ChainElement{}

		newcs = append(newcs, c)
	}

	chains = newcs
}

// execChains executes chains for an event we received
func execChains(event *Event) {
	for _, c := range chains {
		if c.Event.Name != event.Name || c.Event.Bee != event.Bee {
			continue
		}

		m := make(map[string]interface{})
		for _, opt := range event.Options {
			m[opt.Name] = opt.Value
		}

		failed := false
		log.Println("Executing chain:", c.Name, "-", c.Description)
		for _, el := range c.Filters {
			if execFilter(el, m) {
				log.Println("\t\tPassed filter!")
			} else {
				log.Println("\t\tDid not pass filter!")
				failed = true
				break
			}
		}
		if failed {
			continue
		}

		for _, el := range c.Actions {
			action := GetAction(el)
			if action == nil {
				log.Println("\t\tERROR: Unknown action referenced!")
				continue
			}
			execAction(*action, m)
		}
	}
}
