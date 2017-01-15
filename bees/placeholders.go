/*
 *    Copyright (C) 2017 Christian Muehlhaeuser
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
package bees

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// A Placeholder used by ins & outs of a bee.
type Placeholder struct {
	Name  string
	Type  string
	Value interface{}
}

type PlaceholderSlice []Placeholder

// Retrieve a value from a Placeholder slice
func (ph PlaceholderSlice) Value(name string) interface{} {
	for _, p := range ph {
		if p.Name == name {
			return p.Value
		}
	}

	return nil
}

// Bind a value from a Placeholder slice
func (ph PlaceholderSlice) Bind(name string, dst interface{}) error {
	v := ph.Value(name)
	if v == nil {
		return errors.New("Placeholder with name " + name + " not found")
	}

	switch d := dst.(type) {
	case *string:
		switch vt := v.(type) {
		case string:
			*d = vt
		case bool:
			*d = strconv.FormatBool(vt)
		case int64:
			*d = strconv.FormatInt(vt, 10)
		case int:
			*d = strconv.FormatInt(int64(vt), 10)
		default:
			panic(fmt.Sprintf("Unhandled type %+v", vt))
		}

	case *bool:
		switch vt := v.(type) {
		case bool:
			*d = vt
		case string:
			vt = strings.ToLower(vt)
			if vt == "true" || vt == "on" || vt == "yes" || vt == "1" {
				*d = true
			}
		case int64:
			*d = vt > 0
		case int:
			*d = vt > 0
		case float64:
			*d = vt > 0
		default:
			panic(fmt.Sprintf("Unhandled type %+v", vt))
		}

	case *float64:
		switch vt := v.(type) {
		case float64:
			*d = vt
		case string:
			x, _ := strconv.Atoi(vt)
			*d = float64(x)
		default:
			panic(fmt.Sprintf("Unhandled type %+v", vt))
		}

	case *int:
		switch vt := v.(type) {
		case int:
			*d = vt
		case uint:
			*d = int(vt)
		case float64:
			*d = int(vt)
		case string:
			*d, _ = strconv.Atoi(vt)
		default:
			panic(fmt.Sprintf("Unhandled type %+v", vt))
		}

	default:
		panic(fmt.Sprintf("Unhandled dst type %+v", dst))
	}

	return nil
}
