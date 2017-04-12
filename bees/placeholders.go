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

// Package bees is Beehive's central module system.
package bees

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Placeholders is an array of Placeholder.
type Placeholders []Placeholder

// Placeholder used by ins & outs of a bee.
type Placeholder struct {
	Name  string
	Type  string
	Value interface{}
}

// SetValue sets a value in the Placeholder slice.
func (ph *Placeholders) SetValue(name string, _type string, value interface{}) {
	if ph.Value(name) == nil {
		p := Placeholder{
			Name:  name,
			Type:  _type,
			Value: value,
		}
		*ph = append(*ph, p)
	} else {
		for i := 0; i < len(*ph); i++ {
			if (*ph)[i].Name == name {
				(*ph)[i].Type = _type
				(*ph)[i].Value = value
			}
		}
	}
}

// Value retrieves a value from a Placeholder slice.
func (ph Placeholders) Value(name string) interface{} {
	for _, p := range ph {
		if p.Name == name {
			return p.Value
		}
	}

	return nil
}

// Bind a value from a Placeholder slice.
func (ph Placeholders) Bind(name string, dst interface{}) error {
	v := ph.Value(name)
	if v == nil {
		return errors.New("Placeholder with name " + name + " not found")
	}

	return ConvertValue(v, dst)
}

// ConvertValue tries to convert v to dst.
func ConvertValue(v interface{}, dst interface{}) error {
	switch d := dst.(type) {
	case *string:
		switch vt := v.(type) {
		case string:
			*d = vt
		case []string:
			*d = strings.Join(vt, ",")
		case bool:
			*d = strconv.FormatBool(vt)
		case int64:
			*d = strconv.FormatInt(vt, 10)
		case float64:
			*d = strconv.FormatFloat(vt, 'f', -1, 64)
		case int:
			*d = strconv.FormatInt(int64(vt), 10)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for string conversion", reflect.TypeOf(vt)))
		}

	case *[]string:
		switch vt := v.(type) {
		case []string:
			*d = vt
		case string:
			*d = strings.Split(vt, ",")
		default:
			panic(fmt.Sprintf("Unhandled type %+v for []string conversion", reflect.TypeOf(vt)))
		}

	case *bool:
		switch vt := v.(type) {
		case bool:
			*d = vt
		case string:
			vt = strings.ToLower(vt)
			if vt == "true" || vt == "on" || vt == "yes" || vt == "1" || vt == "t" {
				*d = true
			}
		case int64:
			*d = vt > 0
		case int:
			*d = vt > 0
		case uint64:
			*d = vt > 0
		case uint:
			*d = vt > 0
		case float64:
			*d = vt > 0
		default:
			panic(fmt.Sprintf("Unhandled type %+v for bool conversion", reflect.TypeOf(vt)))
		}

	case *float64:
		switch vt := v.(type) {
		case int64:
			*d = float64(vt)
		case int32:
			*d = float64(vt)
		case int16:
			*d = float64(vt)
		case int8:
			*d = float64(vt)
		case int:
			*d = float64(vt)
		case uint64:
			*d = float64(vt)
		case uint32:
			*d = float64(vt)
		case uint16:
			*d = float64(vt)
		case uint8:
			*d = float64(vt)
		case uint:
			*d = float64(vt)
		case float64:
			*d = vt
		case float32:
			*d = float64(vt)
		case string:
			x, _ := strconv.Atoi(vt)
			*d = float64(x)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for float64 conversion", reflect.TypeOf(vt)))
		}

	case *int:
		switch vt := v.(type) {
		case int64:
			*d = int(vt)
		case int32:
			*d = int(vt)
		case int16:
			*d = int(vt)
		case int8:
			*d = int(vt)
		case int:
			*d = vt
		case uint64:
			*d = int(vt)
		case uint32:
			*d = int(vt)
		case uint16:
			*d = int(vt)
		case uint8:
			*d = int(vt)
		case uint:
			*d = int(vt)
		case float64:
			*d = int(vt)
		case float32:
			*d = int(vt)
		case string:
			*d, _ = strconv.Atoi(vt)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for int conversion", reflect.TypeOf(vt)))
		}

	case *url.Values:
		switch vt := v.(type) {
		case string:
			*d, _ = url.ParseQuery(vt)
		default:
			panic(fmt.Sprintf("Unhandled type %+v for url.Values conversion", reflect.TypeOf(vt)))
		}

	default:
		panic(fmt.Sprintf("Unhandled dst type %+v", reflect.TypeOf(dst)))
	}

	return nil
}
