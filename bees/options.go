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

import "errors"

// A FilterOption used by filters.
type FilterOption struct {
	Name            string
	Type            string
	Inverse         bool
	CaseInsensitive bool
	Trimmed         bool
	Value           interface{}
}

// BeeOptions is an array of BeeOption.
type BeeOptions []BeeOption

// A BeeOption is used to configure bees.
type BeeOption struct {
	Name  string
	Value interface{}
}

// Value retrieves a value from a BeeOptions slice.
func (opts BeeOptions) Value(name string) interface{} {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.Value
		}
	}

	return nil
}

// Bind a value from a BeeOptions slice.
func (opts BeeOptions) Bind(name string, dst interface{}) error {
	v := opts.Value(name)
	if v == nil {
		return errors.New("Option with name " + name + " not found")
	}

	return ConvertValue(v, dst)
}
