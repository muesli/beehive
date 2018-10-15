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

// Package app is Beehive's application container. Handles command-line arguments parsing.
package app

import (
	"flag"
)

// A CliFlag can be added by Beehive modules to map a command-line parameter
// to a local variable
type CliFlag struct {
	V     interface{}
	Name  string
	Value interface{}
	Desc  string
}

var (
	appflags    []CliFlag
	configFile  string
	versionFlag bool
)

// AddFlags adds CliFlags to appflags
func AddFlags(flags []CliFlag) {
	for _, flag := range flags {
		appflags = append(appflags, flag)
	}
}

// Run sets up all the cli-param mappings
func Run() {
	for _, f := range appflags {
		switch f.Value.(type) {
		case string:
			flag.StringVar((f.V).(*string), f.Name, f.Value.(string), f.Desc)
		case bool:
			flag.BoolVar((f.V).(*bool), f.Name, f.Value.(bool), f.Desc)
		}
	}

	flag.Parse()
}
