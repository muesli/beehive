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

// beehive's equals-filter.
package equalsfilter

import (
	"github.com/muesli/beehive/filters"
)

type EqualsFilter struct {
}

func (filter *EqualsFilter) Name() string {
	return "equals"
}

func (filter *EqualsFilter) Description() string {
	return "This filter passes when a placeholder equals a specific thing"
}

func (filter *EqualsFilter) Passes(data interface{}, value interface{}) bool {
	switch v := data.(type) {
	case string:
		return v == value.(string)
	case bool:
		return v == value.(bool)
	}

	return false
}

func init() {
	f := EqualsFilter{}

	filters.RegisterFilter(&f)
}
