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

// beehive's contains-filter.
package containsfilter

import (
	"github.com/muesli/beehive/filters"
	"strings"
)

type ContainsFilter struct {
}

func (filter *ContainsFilter) Name() string {
	return "contains"
}

func (filter *ContainsFilter) Description() string {
	return "This filter passes when a placeholder contains a specific thing"
}

func (filter *ContainsFilter) Passes(data interface{}, value interface{}) bool {
	switch v := data.(type) {
		case string:
			return strings.Contains(v, value.(string))
		//FIXME: support maps
	}

	return false
}

func init() {
	f := ContainsFilter{}

	filters.RegisterFilter(&f)
}
