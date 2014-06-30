/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 */

package booleanfilter

import (
	"github.com/muesli/beehive/filters"
	"log"
)

const (
	OR  = iota
	AND = iota
	XOR = iota
	NOT = iota
)

type BooleanFilter struct {
	filters   []filters.FilterInterface
	operation int
}

func New(operation int, filters ...filters.FilterInterface) {
	b := new(BooleanFilter)
	b.filters = filters
	b.operation = operation
}

func (filter *BooleanFilter) Name() string {
	return "boolean"
}

func (filter *BooleanFilter) Description() string {
	return "This filter allows boolean operations on 1 (NOT) or more (AND, OR, XOR) filters"
}

func (filter *BooleanFilter) Passes(data interface{}, value interface{}) bool {
	switch filter.operation {
	case OR:
		for f := range filter.filters {
			if filter.filters[f].Passes(data, value) {
				return true
			}
		}
		return false
	case AND:
		for f := range filter.filters {
			if !filter.filters[f].Passes(data, value) {
				return false
			}
		}
		return true
	case XOR:
		true_found := 0
		for f := range filter.filters {
			if filter.filters[f].Passes(data, value) {
				true_found++
			}
			if true_found > 1 {
				return false
			}
		}
		return true_found == 1
	case NOT:
		return !filter.filters[0].Passes(data, value)
	default:
		log.Println("Cannot join filters: No valid operation!")
	}
	return false
}

func init() {
	f := BooleanFilter{}
	filters.RegisterFilter(&f)
}
