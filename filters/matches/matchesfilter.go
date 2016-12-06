/*
 *    Copyright (C) 2016 Sergio Rubio
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
 *      Sergio Rubio <rubiojr@frameos.org>
 */

package matchesfilter

import (
	"regexp"

	"github.com/muesli/beehive/filters"
)

type MatchesFilter struct {
}

func (filter *MatchesFilter) Name() string {
	return "matches"
}

func (filter *MatchesFilter) Description() string {
	return "This filter passes when a placeholder matches a regular expression"
}

func (filter *MatchesFilter) Passes(data interface{}, value interface{}) bool {
	switch v := data.(type) {
	case string:
		matched, _ := regexp.MatchString(value.(string), v)
		return matched
	default:
		return false
	}
}

func init() {
	f := MatchesFilter{}
	filters.RegisterFilter(&f)
}
