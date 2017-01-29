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

// Package filters contains Beehive's filter system.
package filters

import (
	log "github.com/Sirupsen/logrus"
)

// FilterInterface is an interface all Filters implement.
type FilterInterface interface {
	// Name of the filter
	Name() string
	// Description of the filter
	Description() string

	// Execute the filter
	Passes(data interface{}, value interface{}) bool
}

var (
	filters = make(map[string]*FilterInterface)
)

// RegisterFilter gets called by Filters to register themselves.
func RegisterFilter(filter FilterInterface) {
	log.Println("Filter bee ready:", filter.Name(), "-", filter.Description())
	filters[filter.Name()] = &filter
}

// GetFilter returns a filter with a specific name
func GetFilter(identifier string) *FilterInterface {
	filter, ok := filters[identifier]
	if ok {
		return filter
	}

	return nil
}
