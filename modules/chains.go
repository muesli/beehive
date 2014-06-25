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
	"bytes"
	"log"
	"text/template"

	"github.com/muesli/beehive/filters"
)

// An element in a Chain
type ChainElement struct {
	Action  Action
	Filter	Filter
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
			m := make(map[string]interface{})
			for _, opt := range event.Options {
				m[opt.Name] = opt.Value
			}

			if el.Filter.Name != "" {
				filter := *filters.GetFilter(el.Filter.Name)
				passes := true

				log.Println("\tExecuting filter:", filter.Name(), "-", filter.Description())
				for _, opt := range el.Filter.Options {
					log.Println("\t\tOptions:", opt)
					if filter.Passes(m[opt.Name], opt.Value) == opt.Inverse {
						log.Println("\t\tDid not pass filter!")
						passes = false
						break
					}
				}

				if !passes {
					break
				}
				log.Println("\t\tPassed filter!")
			}
			if el.Action.Name != "" {
				action := Action{
					Namespace: el.Action.Namespace,
					Name: el.Action.Name,
				}

				for _, opt := range el.Action.Options {
					var value bytes.Buffer
					tmpl, err := template.New(el.Action.Namespace + "_" + el.Action.Name + "_" + opt.Name).Parse(opt.Value.(string))
					if err == nil {
						err = tmpl.Execute(&value, m)
					}
					if err != nil {
						panic(err)
					}

					ph := Placeholder{
						Name:  opt.Name,
						Type:  "string", //FIXME
						Value: value.String(),
					}
					action.Options = append(action.Options, ph)
				}

				log.Println("\tExecuting action:", action.Namespace, "/", action.Name, "-", GetActionDescriptor(&action).Description)
				for _, v := range action.Options {
					log.Println("\t\tOptions:", v)
				}
				(*GetModule(action.Namespace)).Action(action)
			}
		}
	}
}
