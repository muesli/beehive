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

// beehive's filter system.
package filters

import (
	"log"
)

// Interface which all filters need to implement
type FilterInterface interface {
	// Name of the filter
	Name() string
	// Description of the filter
	Description() string

	// Execute the filter
	Passes(data map[string]interface{}) bool

    // Handle FilterOptions
    SetOptions(options []FilterOption)
    GetOptions() []FilterOption
}

// A FilterOption used by filters
type FilterOption struct {
	Name            string
    Type            string
	Value           interface{}
}

type Filter struct {
    Name    string
    Options []FilterOption
}

var (
	filters map[string]*FilterInterface = make(map[string]*FilterInterface)
)

// Filters need to call this method to register themselves
func RegisterFilter(filter FilterInterface) {
	log.Println("Filter bee ready:", filter.Name(), "-", filter.Description())
	filters[filter.Name()] = &filter
}

// Returns filter with this name
func GetFilter(identifier string) *FilterInterface {
	filter, ok := filters[identifier]
	if ok {
		return filter
	}

	return nil
}
