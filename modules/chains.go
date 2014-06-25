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
