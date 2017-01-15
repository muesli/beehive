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

// beehive's startsWith-filter.
package startswithfilter

import (
	"github.com/muesli/beehive/filters"
	"strings"
)

type StartsWithFilter struct {
}

func (filter *StartsWithFilter) Name() string {
	return "startswith"
}

func (filter *StartsWithFilter) Description() string {
	return "This filter passes when a placeholder starts with a specific thing"
}

func (filter *StartsWithFilter) Passes(data interface{}, value interface{}) bool {
	switch d := data.(type) {
	case string:
		switch v := value.(type) {
		case string:
			return strings.HasPrefix(d, v)
		}
		//FIXME: support maps
	}

	return false
}

func init() {
	f := StartsWithFilter{}

	filters.RegisterFilter(&f)
}
