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

// beehive's endsWith-filter.
package endswithfilter

import (
	"github.com/muesli/beehive/filters"
	"strings"
)

type EndsWithFilter struct {
}

func (filter *EndsWithFilter) Name() string {
	return "endswith"
}

func (filter *EndsWithFilter) Description() string {
	return "This filter passes when a placeholder starts with a specific thing"
}

func (filter *EndsWithFilter) Passes(data interface{}, value interface{}) bool {
	switch v := data.(type) {
	case string:
		return strings.HasSuffix(v, value.(string))
		//FIXME: support maps
	}

	return false
}

func init() {
	f := EndsWithFilter{}

	filters.RegisterFilter(&f)
}
