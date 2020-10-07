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
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/muesli/beehive/filters"
)

// Filter describes a user configured event filter.
type Filter struct {
	ID      string
	Name    string
	Options FilterOption
}

// execFilter executes a filter. Returns whether the filter passed or not.
func execFilter(source string, opts map[string]interface{}) bool {
	name := "template"
	if strings.Contains(source, "def main(") {
		name = "starlark"
	}

	f := *filters.GetFilter(name)
	log.Println("\tExecuting filter:", source)

	defer func() {
		if e := recover(); e != nil {
			log.Println("Fatal filter event:", e)
		}
	}()

	return f.Passes(opts, source)
}
